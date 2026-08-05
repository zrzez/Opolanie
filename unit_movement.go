package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	maxNoMoveTicks        = 100
	maxPathfindingBudget  = 40
	maxPathfindingRetries = 3
)

func (u *unit) executeAStarMovement(pathfindingBudget int, bState *battleState) {
	if !u.ensureValidPath(pathfindingBudget, bState) {
		return
	}

	u.moveAlongPath(bState.Board)
}

func (u *unit) ensureValidPath(pathfindingBudget int, bState *battleState) bool {
	if u.Command.isInteraction() {
		approachTile := &bState.Board.Tiles[u.Approach.X][u.Approach.Y]

		if approachTile.Unit != nil && approachTile.Unit.ID != u.ID {
			newApproach, err := u.calculateApproachTile(u.Target, bState)
			if err == nil {
				u.Approach = newApproach
				u.invalidatePathForRecalculation()
			} else {
				u.setIdleWithReason("brak wolnego miejsca do podejścia")

				return false
			}
		}
	}

	if u.NoMoveTicks >= 15 && u.NoMoveTicks%15 == 0 {
		newApproach, err := u.calculateApproachTile(u.Target, bState)
		if err == nil {
			u.Approach = newApproach
			u.invalidatePathForRecalculation()
			u.NoMoveTicks = 0
			u.Delay = u.MaxDelay
		} else {
			x := int(u.Target.Position.X) + rand.Intn(6) - 3
			y := int(u.Target.Position.Y) + rand.Intn(6) - 3

			if bState.Board.isValidWalkableTile(int8(x), int8(y)) {
				u.Approach = point{X: uint8(x), Y: uint8(y)}
				u.invalidatePathForRecalculation()
				u.NoMoveTicks = 0
				u.Delay = u.MaxDelay
			}
		}
	}

	if u.hasValidPath(bState) {
		return true
	}

	if pathfindingBudget >= maxPathfindingBudget {
		u.waitForPathfindingBudget()

		return false
	}

	return u.calculateNewPath(bState.Board)
}

func (u *unit) calculateNewPath(board *boardData) bool {
	// 1. Jeśli jednostka nie ma u.Path to tworzymy
	// @reminder: 100 to liczba, która wydaje mi się wystarczająca po kilku próbach.
	//   Nie stoi za tym nic więcej niż kilka prób.
	if cap(u.Path) < 100 {
		fmt.Printf("calculateNewPath powiększa pojemność u.Path, bo 100 to zbyt mało. Potrzebowałem %v\n", cap(u.Path))
		u.Path = make([]point, 0, 100)
	}
	newPath := findPath(board, u, u.X, u.Y, u.Approach.X, u.Approach.Y, u.Path)

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
	u.Delay = u.MaxDelay
}

func (u *unit) waitForPathfindingBudget() {
	u.State = stateWaiting
	u.Delay = uint16(3 + rng.Intn(5))
}

func (u *unit) setPathAndState(path []point) {
	u.setPath(path)
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
	} else {
		u.NoMoveTicks++
		u.State = stateWaiting
		u.Delay = 1
	}
}

// Funkcje pomocnicze dla zarządzania ścieżką.
func (u *unit) setPath(path []point) {
	u.Path = path

	u.PathIndex = 1
}

func (u *unit) clearPath() {
	u.Path = u.Path[:0]
	u.PathIndex = 0
	u.NoMoveTicks = 0
	u.LastX, u.LastY = 0, 0
	u.Direction = rl.Vector2{X: 0, Y: 0}
}

func (u *unit) handlePathfindingFailure() {
	u.RetryAttempts++

	if u.RetryAttempts >= maxPathfindingRetries {
		u.setIdleWithReason("odnajdywanie drogi się wyłożyło na dobre")
		u.RetryAttempts = 0
	} else {
		u.State = stateWaiting
		u.Delay = uint16(40 + rng.Intn(20))
		log.Printf("jednostka %d: odnajdywanie drogi zawiodło (podejście %d/%d)",
			u.ID, u.RetryAttempts, maxPathfindingRetries)
	}
}

func (u *unit) hasValidPath(bState *battleState) bool {
	if len(u.Path) == 0 || u.PathIndex >= len(u.Path) {
		return false
	}

	target, err := bState.resolveTarget(u.Target)
	if err != nil {
		return false
	}

	switch u.Target.Kind {
	case targetTile: // drzewa, puste pola
		return true
	case targetBuilding:
		return true
	case targetUnit:
		if target.Unit != nil {
			return u.Target.Position.X == target.Unit.X && u.Target.Position.Y == target.Unit.Y
		}
	case targetNone:
		// nie wiem co to miałoby oznaczać
		return false
	default:
		// jeśli nie wiemy co, to przypał
		return false
	}

	return false
}

func (u *unit) invalidatePathForRecalculation() {
	u.Path = u.Path[:0]
	u.PathIndex = 0
}

func (u *unit) resetMovementCounters() {
	u.NoMoveTicks = 0
	u.LastX, u.LastY = u.X, u.Y
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
	return u.X == u.Approach.X && u.Y == u.Approach.Y
}

func (u *unit) move(pathfindingBudget int, bState *battleState) {
	if u.Command == cmdUAttack {
		if u.canAttackTargetFromCurrentPosition(bState) {
			u.clearPath()
			u.State = stateAttacking
			// u.setAnimationType()

			return
		}

		// 25.04.2026 Dodaję bezpiecznik przerywający ruch jeśli cel przestał istnieć
		// Bez tego jednostka atakująca drzewo zaczyna się przemieszczać po jego upadku
		// szukając nowej pozycji do ataku nieistniejącego już celu.
		_, err := bState.resolveTarget(u.Target)
		if err != nil {
			u.setIdleWithReason("cel ataku przestał istnieć")

			return
		}
	}

	if u.isAtTarget() {
		u.handleTargetReached(bState)

		return
	}

	// Jeśli utknęliśmy, to zarzucamy ruch
	if u.NoMoveTicks > maxNoMoveTicks {
		u.setIdleWithReason("przyblokowana")

		return
	}

	u.executeAStarMovement(pathfindingBudget, bState)
}

func (u *unit) handleMovementTargetReached(bState *battleState) {
	// @todo: ogarnij, czy nie powinienem Target.position zastąpić u.Approach
	if u.State == stateMoving && u.X == u.Target.Position.X && u.Y == u.Target.Position.Y {
		u.handleTargetReached(bState)
	}
}
