package main

import (
	"log"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	maxNoMoveTicks        = 200
	maxBlockedTicks       = 50
	maxPathfindingBudget  = 40
	maxPathfindingRetries = 3
	maxMovementHistory    = 6
)

func (u *unit) executeAStarMovement(resolver objectResolver, board *boardData, pathfindingBudget *int) {
	if !u.ensureValidPath(resolver, board, pathfindingBudget) {
		return
	}

	u.moveAlongPath(board)
}

func (u *unit) ensureValidPath(resolver objectResolver, board *boardData, pathfindingBudget *int) bool {
	if u.hasValidPath(resolver, board) {
		return true
	}

	if *pathfindingBudget >= maxPathfindingBudget {
		u.waitForPathfindingBudget()

		return false
	}

	*pathfindingBudget++

	return u.calculateNewPath(board)
}

func (u *unit) calculateNewPath(board *boardData) bool {
	newPath := findPath(board, u, u.X, u.Y, u.ApproachX, u.ApproachY)

	if newPath == nil {
		u.handlePathfindingFailure()

		return false
	}

	u.setPathAndState(newPath)

	return true
}

func (u *unit) executeSuccessfulMove(x, y uint8, board *boardData) {
	u.executeMove(x, y, board)
	u.resetMovementCounters()
	u.updateMovementHistory()
}

func (u *unit) waitForPathfindingBudget() {
	u.State = stateWaiting
	u.Delay = uint16(3 + rng.Intn(5))
	log.Printf("unit %d: waiting for pathfinding budget", u.ID)
}

func (u *unit) setPathAndState(path []*pathNode) {
	u.setPath(path)
	u.LastTargetX = u.TargetX
	u.LastTargetY = u.TargetY
	u.State = u.determineActiveStateFromCommand()
	u.RetryAttempts = 0
}

func (u *unit) moveAlongPath(board *boardData) {
	if u.PathIndex >= len(u.Path) {
		u.clearPath()

		return
	}

	next := u.Path[u.PathIndex]

	if u.canMoveTo(next.X, next.Y, board) {
		u.executeSuccessfulMove(next.X, next.Y, board)
	}
}

// Funkcje pomocnicze dla zarządzania ścieżką
func (u *unit) setPath(path []*pathNode) {
	u.Path = make([]*pathNode, len(path))
	for i, node := range path {
		u.Path[len(path)-1-i] = node
	}
	u.PathIndex = 1
}

func (u *unit) clearPath() {
	u.Path = nil
	u.PathIndex = 0
	u.NoMoveTicks = 0
	u.LastX, u.LastY = 0, 0
}

func (u *unit) handlePathfindingFailure() {
	u.RetryAttempts++

	if u.RetryAttempts >= maxPathfindingRetries {
		u.setIdleWithReason("pathfinding permanently failed")
		u.RetryAttempts = 0
	} else {
		u.State = stateWaiting
		u.Delay = uint16(40 + rng.Intn(20))
		log.Printf("unit %d: pathfinding failed (attempt %d/%d)",
			u.ID, u.RetryAttempts, maxPathfindingRetries)
	}
}

func (u *unit) shouldAbortMovement() bool {
	if u.NoMoveTicks > maxNoMoveTicks {
		return true
	}

	if u.detectSimpleOscillation() {
		return true
	}

	return false
}

func (u *unit) detectSimpleOscillation() bool {
	if len(u.History) < 4 {
		return false
	}

	n := len(u.History)
	a := u.History[n-1]
	b := u.History[n-2]
	c := u.History[n-3]
	d := u.History[n-4]

	return a.X == c.X && a.Y == c.Y && b.X == d.X && b.Y == d.Y
}

func (u *unit) hasValidPath(resolver objectResolver, board *boardData) bool {
	if len(u.Path) == 0 || u.PathIndex >= len(u.Path) {
		return false
	}

	target, err := u.validateTargetExists(resolver, board)
	if err != nil {
		return false
	}

	if u.TargetID == 0 {
		return true
	}

	if target.Building != nil {
		return u.TargetID == ObjectID(target.Building.ID)
	}

	if target.Unit != nil {
		return u.TargetX == target.Unit.X && u.TargetY == target.Unit.Y
	}

	return false
}

func (u *unit) invalidatePathForRecalculation() {
	u.Path = nil
	u.PathIndex = 0
}

func (u *unit) resetMovementCounters() {
	u.NoMoveTicks = 0
	u.BlockedCounter = 0
	u.LastX, u.LastY = u.X, u.Y
}

func (u *unit) updateMovementHistory() {
	u.History = append(u.History, rl.NewVector2(float32(u.X), float32(u.Y)))

	if len(u.History) > maxMovementHistory {
		u.History = u.History[1:]
	}
}

func (u *unit) handleNoMovementDetection() bool {
	if (u.X == u.LastX && u.Y == u.LastY) && u.State != stateMilking && u.State != stateRepairing {
		u.NoMoveTicks++
		if u.NoMoveTicks > maxNoMoveTicks {
			u.clearPath()
			u.setIdleWithReason("zbyt długi bezruch")

			return true
		}
	} else {
		u.NoMoveTicks = 0
	}

	u.LastX, u.LastY = u.X, u.Y

	return false
}

func (u *unit) canMoveTo(x, y uint8, board *boardData) bool {
	if x >= boardMaxX || y >= boardMaxY {
		return false
	}

	currentTile := &board.Tiles[x][y]

	// Kolizja z jednostkami (standardowo)
	if currentTile.Unit != nil && currentTile.Unit.ID != u.ID {
		return false
	}

	// Kolizja z terenem/budynkami (używamy nowej funkcji z pathfinding.go)
	// Przekazujemy 'u', aby obsłużyć wyjątek krowy wchodzącej do obory
	return isWalkableUnit(board, x, y, u)
}

// calculateMilkingSpot oblicza milking spot dla obory.
func (bld *building) calculateMilkingSpot() (uint8, uint8, bool) {
	if len(bld.OccupiedTiles) == 0 {
		return 0, 0, false
	}

	minX := uint8(math.MaxUint8)
	maxY := uint8(0)

	for _, occupiedTile := range bld.OccupiedTiles {
		if occupiedTile.X < minX {
			minX = occupiedTile.X
		}

		if occupiedTile.Y > maxY {
			maxY = occupiedTile.Y
		}
	}

	return minX, maxY, true
}

// executeMove wykonuje ruch na nową pozycję.
func (u *unit) executeMove(x, y uint8, board *boardData) {
	// ZMIANA: Używamy nowej struktury Tiles
	// Usuń z poprzedniej pozycji (jeśli to ta jednostka tam jest)
	if board.Tiles[u.X][u.Y].Unit == u {
		board.Tiles[u.X][u.Y].Unit = nil
	}

	oldX, oldY := u.X, u.Y

	u.X, u.Y = x, y

	// Ustaw na nowej pozycji
	board.Tiles[u.X][u.Y].Unit = u

	u.PathIndex++
	u.updateMovementAnimation(oldX, oldY)
}

func (u *unit) isAtTarget() bool {
	return u.X == u.ApproachX && u.Y == u.ApproachY
}

func (u *unit) move(resolver objectResolver, board *boardData, pathfindingBudget *int, bState *battleState) {
	if u.Command == cmdUAttack {
		log.Printf("INFO: units.go move rozkaz to cmdAttack")

		if u.canAttackTargetFromCurrentPosition(resolver, board) {
			log.Printf("INFO: units.go move cel osiągalny z tego miejsca")
			u.clearPath()
			u.State = stateAttacking

			return
		}

		// 25.04.2026 Dodaję bezpiecznik przerywający ruch jeśli cel przestał istnieć
		// Bez tego jednostka atakująca drzewo zaczyna się przemieszczać po jego upadku
		// szukając nowej pozycji do ataku nieistniejącego już celu.
		_, err := u.validateTargetExists(resolver, board)
		if err != nil {
			u.setIdleWithReason("cel ataku przestał istnieć")

			return
		}
	}

	if u.isAtTarget() {
		log.Printf("INFO: units.go move u celu")

		u.handleTargetReached(resolver, board, bState)

		return
	}

	if u.shouldAbortMovement() {
		u.setIdleWithReason("przerwano ruch")

		return
	}

	u.executeAStarMovement(resolver, board, pathfindingBudget)
}

func (u *unit) handleMovementTargetReached(resolver objectResolver, board *boardData, bState *battleState) {
	if u.State == stateMoving && u.X == u.TargetX && u.Y == u.TargetY {
		u.handleTargetReached(resolver, board, bState)
	}
}
