package main

// units.go

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	maxNoMoveTicks        = 200
	maxBlockedTicks       = 50
	maxPathfindingBudget  = 40
	maxPathfindingRetries = 3
	maxMovementHistory    = 6
)

func (u *unit) initUnit(unitType unitType, x, y uint8, command uint16, bs *battleState) {
	u.ID = bs.NextUniqueObjectID
	bs.NextUniqueObjectID++
	u.Exists = true
	u.Type = unitType
	u.X = x
	u.Y = y
	u.Command = command
	u.AnimationType = "walk"
	u.AnimationFrame = 0
	u.AnimationCounter = 0
	u.Direction = rl.NewVector2(0, 1)
	u.IsSelected = false

	stats, ok := unitDefs[unitType]
	if ok {
		u.SightRange = stats.SightRange
		u.AttackRange = stats.AttackRange
		u.Damage = stats.BaseDamage
		u.Armor = stats.BaseArmor
		u.MaxHP = stats.MaxHP
		u.MaxDelay = stats.MoveDelay
		u.HP = stats.MaxHP

		if stats.MaxMana > 0 {
			u.MaxMana = stats.MaxMana
			u.Mana = stats.MaxMana / 2
		}
	} else {
		log.Printf("OSTRZEŻENIE: Nieznany rodzaj jednostki w init: %d.", unitType)
		u.SightRange = 1
		u.AttackRange = 1
		u.Damage = 1
		u.Armor = 0
		u.MaxHP = 10
		u.MaxDelay = 20
		u.HP = u.MaxHP
		u.Mana = 0
	}

	u.Experience = 0
	u.Delay = u.MaxDelay
	u.AttackCooldown = 0
	u.Wounds = make([]wound, 0, 6)
}

// isCaster zwraca true jeżeli dana jednostka czaruje.
func (ut unitType) isCaster() bool {
	return ut == unitMage || ut == unitPriest || ut == unitPriestess
}

// hasMana zwraca true jeżeli dana jednostka może mieć więcej niż 0 Max_Mana
// robi to sprawdzając, czy isCaster lub == UNKNOWN.
func (ut unitType) hasMana() bool {
	return ut.isCaster() || ut == unitUnknown
}

// increaseManaUnit dla każdej istniejącej jednostki zwiększa manę o amount
// Pilnuje, aby u.Mana <= u.MaxMana.
func (u *unit) increaseManaUnit(amount uint16) {
	if !u.Exists {
		return
	}

	u.Mana += amount

	if u.Mana > u.MaxMana {
		u.Mana = u.MaxMana
	}
}

// dla każdej istniejącej jednostki zmniejsza manę o amount. Pilnuje, aby u.Mana >= 0.
func (u *unit) decreaseManaUnit(amount uint16) {
	if !u.Exists {
		return
	}

	u.Mana -= amount

	if u.Mana < 0 {
		u.Mana = 0
	}
}

// increaseHPUnit dla każdej istniejącej jednostki zwiększa PŻ o amount
// Pilnuje, aby u.HP <= u.MaxHP
func (u *unit) increaseHPUnit(amount uint16) {
	if !u.Exists {
		return
	}

	u.HP += amount
	if u.HP > u.MaxHP {
		u.HP = u.MaxHP
	}
}

// decreaseHPUnit dla każdej istniejącej jednostki zmniejsza PŻ o amount
// Pilnuje, aby ustawić u.Exists = false.
func (u *unit) decreaseHPUnit(amount uint16) {
	if !u.Exists {
		return
	}

	u.HP -= amount

	if u.HP <= 0 {
		u.HP = 0
		u.Exists = false
	}
}

// show umieszcza jednostkę na mapie w Tiles.
func (u *unit) show(bs *battleState) {
	if u.X < boardMaxX && u.Y < boardMaxY {
		// ZMIANA: Ustawiamy wskaźnik unit w nowej strukturze Tiles
		bs.Board.Tiles[u.X][u.Y].Unit = u
	} else {
		log.Printf("OSTRZEŻENIE: Próba umieszczenia jednostki poza mapą: (%d,%d)", u.X, u.Y)
	}
}

// ============================================================================
// LOGIKA JEDNOSTEK
// ============================================================================

func (u *unit) updateUnit(bs *battleState) {
	// Aktualizowanie ran
	// @todo przenieś do osobnej funkcji, szkoda zajmować tutaj miejsce
	nextFreeIndex := 0

	for scanIndex := range u.Wounds {
		wound := &u.Wounds[scanIndex]
		wound.Timer--

		if wound.Timer > 1 {
			if scanIndex != nextFreeIndex {
				u.Wounds[nextFreeIndex] = *wound
			}

			nextFreeIndex++
		}
	}

	u.Wounds = u.Wounds[:nextFreeIndex]

	u.handleAttackCooldown(bs)

	if u.handleNoMovementDetection() {
		return
	}

	if u.handleDelay(bs) {
		return
	}

	u.processActiveEffects()

	if u.handleBlockedCounter() {
		return
	}

	u.handleWaitingToActiveTransition()
	u.handleMovementTargetReached(bs)
	u.executeCommandAction(bs)
	u.resetDelayIfActive()
}

func (u *unit) handleAttackCooldown(bs *battleState) {
	if bs.GlobalFrameCounter%logicSpeedDivisor == 0 {
		if u.AttackCooldown > 0 {
			u.AttackCooldown--
		}
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

func (u *unit) handleDelay(bs *battleState) bool {
	if u.State == stateIdle {
		u.Delay = u.MaxDelay

		return false
	}
	if u.Delay <= 0 {
		return false
	}
	if bs.GlobalFrameCounter%logicSpeedDivisor != 0 {
		return true
	}
	u.Delay--

	return u.Delay > 0
}

func (u *unit) processActiveEffects() {
	var activeEffects []activeEffect
	for _, e := range u.Effects {
		e.Duration--
		if e.Duration > 0 {
			activeEffects = append(activeEffects, e)
		}
	}

	u.Effects = activeEffects
}

func (u *unit) handleBlockedCounter() bool {
	if u.BlockedCounter > 0 {
		u.BlockedCounter--
		u.AnimationType = "idle"
		u.AnimationFrame = 0
		u.Delay = u.MaxDelay

		return true
	}

	return false
}

func (u *unit) handleWaitingToActiveTransition() {
	if u.State == stateWaiting {
		u.State = u.determineActiveStateFromCommand()
		if u.State == stateAttacking {
			u.AnimationType = "fight"
		} else {
			u.AnimationType = "walk"
		}
	}
}

func (u *unit) handleMovementTargetReached(bs *battleState) {
	if u.State == stateMoving && u.X == u.TargetX && u.Y == u.TargetY {
		u.handleTargetReached(bs)
	}
}

func (u *unit) executeCommandAction(bs *battleState) {
	switch u.Type {
	case unitCow:
		u.handleCowBehavior(bs)
	default:
		u.executeStandardUnitCommand(bs)
	}
}

func (u *unit) executeStandardUnitCommand(bs *battleState) {
	switch u.Command {
	case cmdMove:
		u.move(bs)
	case cmdAttack:
		if u.canAttackTargetFromCurrentPosition(bs) {
			u.State = stateAttacking
			u.clearPath()
			u.attack(bs)
		} else {
			u.State = stateMoving
			u.move(bs)
		}
	case cmdRepairStructure:
		if u.canAttackTargetFromCurrentPosition(bs) {
			u.State = stateRepairing

			if u.AnimationType != "fight" {
				u.AnimationType = "fight"
				u.AnimationFrame = 0
			}

			u.clearPath()
			u.repair(bs)
		} else {
			u.State = stateMoving

			if u.AnimationType != "walk" {
				u.AnimationType = "walk"
			}

			u.move(bs)
		}
	case cmdMagicShield:
		u.castMagicShield()
	case cmdMagicFire:
		// u.castMagicFire(bs)
	case cmdMagicSight:
		u.castMagicSight(bs)
	case cmdIdle, cmdStop:
		u.actOnIdle(bs)
	default:
		panic("unhandled default case")
	}
}

func (u *unit) resetDelayIfActive() {
	if u.State != stateWaiting {
		if u.Delay < u.MaxDelay {
			u.Delay = u.MaxDelay
		}
	}
}

func (u *unit) determineActiveStateFromCommand() unitState {
	switch u.Command {
	case cmdMove:
		return stateMoving
	case cmdAttack:
		return stateAttacking
	case cmdGraze:
		return stateGrazing
	case cmdRepairStructure:
		return stateRepairing
	case cmdMagicShield, cmdMagicFire, cmdMagicSight:
		return stateCastingSpell
	default:
		return stateIdle
	}
}

func getSelectedUnits(bs *battleState) []*unit {
	var selected []*unit

	for _, unit := range bs.Units {
		if unit.IsSelected && unit.Exists && unit.Owner == bs.PlayerID {
			selected = append(selected, unit)
		}
	}

	return selected
}

func canDamagePalisades(unit *unit) bool {
	return unit.Type == unitAxeman || unit.Type == unitPriest
}

// @todo: funkcja jest zbyt rozbudowana. Niepotrzebnie zwraca „error”. -1,-1 powinno wystarczyć.
// Na tym etapie już nie potrzeba rozróżniać czemu nie da rady podejść. Chodzenie jest ogarnięte
// Wrócę do tego pod koniec. Nie jest pilne, nic nie szkodzi, ale za dużo robi.
func (u *unit) resolveApproachPosition(targetID uint, bs *battleState) (uint8, uint8, error) {
	targetUnit, targetBuilding := getObjectByID(targetID, bs)

	if targetBuilding != nil && targetBuilding.Exists {
		if u.AttackRange > 1 {
			x, y, ok := findOptimalRangedAttackTile(
				u.X, u.Y, targetBuilding, u.AttackRange, bs,
			)
			if ok {
				return x, y, nil
			}
		}

		x, y, ok := targetBuilding.getClosestWalkableTile(bs)

		if ok {
			return x, y, nil
		}
		return 0, 0, fmt.Errorf("cel (budynek) jest nieosiągalny")
	}

	if targetUnit != nil && targetUnit.Exists {
		return targetUnit.X, targetUnit.Y, nil
	}

	return 0, 0, fmt.Errorf("cel ataku ID %d nie istnieje", targetID)
}

// @todo: cmdFlee wykrzacza grę!
//func (u *unit) addUnitCommand(command uint16, targetX, targetY uint8, targetID uint, bs *battleState) {
//	// 1. Sprawdzamy, czy trzeba podejść do celu
//	isInteraction := command == cmdAttack || command == cmdRepairStructure || command == cmdBuildStructure
//	if isInteraction && targetID != 0 {
//		var finalTargetX, finalTargetY uint8
//
//		var err error
//
//		finalTargetX, finalTargetY, err = u.resolveApproachPosition(targetID, bs)
//		if err != nil {
//			// log.Printf("BŁĄD: addUnitCommand: %v", err)
//			u.setIdleWithReason("cel nieosiągalny")
//
//			return
//		}
//
//		targetX, targetY = finalTargetX, finalTargetY
//	}
//
//	// @todo: nie pamiętam po co to było. Chyba dla przemieszczania się do celu
//	if u.Command == command &&
//		u.TargetX == targetX &&
//		u.TargetY == targetY &&
//		u.TargetID == targetID {
//		return
//	}
//
//	if command == cmdAttack {
//		_, targetBuilding := getObjectByID(targetID, bs)
//		if targetBuilding != nil &&
//			targetBuilding.Exists &&
//			targetBuilding.Type == buildingPalisade &&
//			!canDamagePalisades(u) {
//			log.Printf(
//				"INFO: Jednostka %d nie może atakować palisad.", u.ID,
//			)
//
//			return
//		}
//	}
//
//	u.clearPath()
//	u.History = nil
//	u.LoopCount = 0
//	u.TicksNoProgress = 0
//	u.LastPathIndex = 0
//
//	u.Command = command
//	u.TargetX = targetX
//	u.TargetY = targetY
//	u.TargetID = targetID
//
//	u.State = stateMoving
//	u.Delay = 0
//
//	switch command {
//	case cmdAttack:
//		u.State = stateAttacking
//		u.AnimationType = "fight"
//		u.AnimationFrame = 3
//		u.AnimationCounter = 0
//	case cmdMove:
//		u.State = stateMoving
//		u.AnimationType = "walk"
//	case cmdStop:
//		u.State = stateIdle
//		u.AnimationType = "walk"
//		u.AnimationFrame = 0
//		u.Command = cmdIdle
//	case cmdMagicFire, cmdMagicSight, cmdMagicShield, cmdMagicLightning:
//		u.State = stateCastingSpell
//		// @todo: jeżeli dobrze rozumiem to tutaj mogę wywołać magiczną tarczę dla
//		// dla kapłanki, bo sama jest celem i nie wymaga to więcej niż
//		u.castMagicShield()
//		// @todo: jak ogarnąć gromobicie?!
//	case cmdGraze:
//		u.State = stateGrazing
//	case cmdRepairStructure, cmdBuildStructure:
//		u.State = stateRepairing
//	default:
//		panic("unhandled default case")
//	}
//}

func (u *unit) addUnitCommand(command uint16, targetX, targetY uint8, targetID uint, bs *battleState) {
	if u.shouldSkipDuplicate(command, targetX, targetY, targetID) {
		return
	}

	if err := u.resolveInteractionTarget(&targetX, &targetY, command, targetID, bs); err != nil {
		u.setIdleWithReason("cel nieosiągalny")

		return
	}

	if !u.validateCommand(command, targetID, bs) {
		return
	}

	u.prepareForNewCommand(command, targetX, targetY, targetID)
	u.applyCommandState(command)
}

func (u *unit) shouldSkipDuplicate(command uint16, targetX, targetY uint8, targetID uint) bool {
	return u.Command == command && u.TargetX == targetX &&
		u.TargetY == targetY && u.TargetID == targetID
}

func (u *unit) resolveInteractionTarget(targetX, targetY *uint8, command uint16, targetID uint, bs *battleState) error {
	if !isInteractionCommand(command) || targetID == 0 {
		return nil
	}

	finalX, finalY, err := u.resolveApproachPosition(targetID, bs)
	if err != nil {
		return err
	}

	*targetX, *targetY = finalX, finalY

	return nil
}

func isInteractionCommand(command uint16) bool {
	switch command {
	case cmdAttack, cmdRepairStructure, cmdBuildStructure:
		return true
	default:
		return false
	}
}

func (u *unit) validateCommand(command uint16, targetID uint, bs *battleState) bool {
	switch command {
	case cmdAttack:
		return u.canAttack(targetID, bs)
	default:
		return true
	}
}

func (u *unit) canAttack(targetID uint, bs *battleState) bool {
	_, building := getObjectByID(targetID, bs)

	if building == nil || !building.Exists {
		return true
	}

	if building.Type == buildingPalisade && !canDamagePalisades(u) {
		log.Printf("INFO: Jednostka %d nie może atakować palisad", u.ID)

		return false
	}

	return true
}

func (u *unit) prepareForNewCommand(command uint16, targetX, targetY uint8, targetID uint) {
	u.clearPath()
	u.History = nil
	u.LoopCount = 0
	u.TicksNoProgress = 0
	u.LastPathIndex = 0
	u.Command = command
	u.TargetX = targetX
	u.TargetY = targetY
	u.TargetID = targetID
	u.Delay = 0
}

// @todo: CmdFlee jest nie ogarnięta i wykrzacza program przez default
func (u *unit) applyCommandState(command uint16) {
	switch command {
	case cmdAttack:
		u.State = stateAttacking
		u.AnimationType = "fight"
		u.AnimationFrame = 3
		u.AnimationCounter = 0
	case cmdMove, cmdFlee:
		u.State = stateMoving
		u.AnimationType = "walk"
	case cmdStop:
		u.State = stateIdle
		u.AnimationType = "walk"
		u.AnimationFrame = 0
		u.Command = cmdIdle
	case cmdMagicFire, cmdMagicSight, cmdMagicShield, cmdMagicLightning:
		u.State = stateCastingSpell
		// @todo: jeżeli dobrze rozumiem to tutaj mogę wywołać magiczną tarczę dla
		// kapłanki, bo sama jest celem i nie wymaga to więcej niż.
		u.castMagicShield()
		// @todo: jak ogarnąć gromobicie?!
	case cmdGraze:
		u.State = stateGrazing
	case cmdRepairStructure, cmdBuildStructure:
		u.State = stateRepairing
	case cmdMilking:
		u.State = stateMilking
	default:
		fmt.Println("DUUUUPA NIE MA TAKIEJ KOMENDY W PRZEŁĄCZNIKU")
		panic("unhandled default case")
	}
}

func (u *unit) isAtTarget() bool {
	return u.X == u.TargetX && u.Y == u.TargetY
}

func (u *unit) move(bs *battleState) {
	if u.Command == cmdAttack {
		if u.canAttackTargetFromCurrentPosition(bs) {
			u.clearPath()
			u.State = stateAttacking

			return
		}
	}

	if u.isAtTarget() {
		u.handleTargetReached(bs)

		return
	}

	if u.shouldAbortMovement() {
		u.setIdleWithReason("przerwano ruch")

		return
	}

	u.executeAStarMovement(bs)
}

func (u *unit) canAttackTargetFromCurrentPosition(bs *battleState) bool {
	// log.Println("Waliduję, czy cel istnieje")

	target, err := u.validateTargetExists(bs)
	if err != nil {
		// log.Println("Cel nie istnieje")

		return false
	}

	// log.Println("Obliczam odległość")

	distance := u.calculateDistanceToTarget(target)

	// log.Printf("Odległość obliczona %d", distance)
	// log.Println(distance <= u.AttackRange)

	return distance <= u.AttackRange
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

func (u *unit) executeAStarMovement(bs *battleState) {
	if !u.ensureValidPath(bs) {
		return
	}

	u.moveAlongPath(bs)
}

func (u *unit) ensureValidPath(bs *battleState) bool {
	if u.hasValidPath(bs) {
		return true
	}

	if bs.pathfindingUnitsThisTick >= maxPathfindingBudget {
		u.waitForPathfindingBudget()

		return false
	}

	return u.calculateNewPath(bs)
}

func (u *unit) hasValidPath(bs *battleState) bool {
	if len(u.Path) == 0 || u.PathIndex >= len(u.Path) {
		return false
	}

	target, err := u.validateTargetExists(bs)
	if err != nil {
		return false
	}

	if target.Building != nil {
		return u.TargetID == target.Building.ID
	}

	if target.Unit != nil {
		return u.TargetX == target.Unit.X && u.TargetY == target.Unit.Y
	}

	return false
}

func (u *unit) waitForPathfindingBudget() {
	u.State = stateWaiting
	u.Delay = uint16(3 + rng.Intn(5))
	log.Printf("unit %d: waiting for pathfinding budget", u.ID)
}

func (u *unit) calculateNewPath(bs *battleState) bool {
	bs.pathfindingUnitsThisTick++

	newPath := findPath(bs, u.ID, u.X, u.Y, u.TargetX, u.TargetY)

	if newPath == nil {
		u.handlePathfindingFailure()
		return false
	}

	u.setPathAndState(newPath)
	return true
}

func (u *unit) setPathAndState(path []*pathNode) {
	u.setPath(path)
	u.LastTargetX = u.TargetX
	u.LastTargetY = u.TargetY
	u.State = u.determineActiveStateFromCommand()
	u.RetryAttempts = 0
}

func (u *unit) moveAlongPath(bs *battleState) {
	if u.PathIndex >= len(u.Path) {
		u.clearPath()
		return
	}

	next := u.Path[u.PathIndex]

	if u.canMoveTo(next.X, next.Y, bs) {
		u.executeSuccessfulMove(next.X, next.Y, bs)
	} else {
		u.handleMovementBlocked(bs, next.X, next.Y)
	}
}

func (u *unit) executeSuccessfulMove(x, y uint8, bs *battleState) {
	u.executeMove(x, y, bs)
	u.resetMovementCounters()
	u.updateMovementHistory()
}

func (u *unit) handleMovementBlocked(bs *battleState, blockedX, blockedY uint8) {
	detourX, detourY, ok := u.findLocalDetour(bs, blockedX, blockedY)

	if !ok {
		u.executeSuccessfulMove(detourX, detourY, bs)
		u.invalidatePathForRecalculation()
		return
	}

	u.handlePersistentBlock()
}

func (u *unit) findLocalDetour(bs *battleState, blockedX, blockedY uint8) (uint8, uint8, bool) {
	bestX, bestY := 0, 0
	bestScore := math.MaxFloat64

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			x, y := int(u.X)+dx, int(u.Y)+dy

			if x == int(blockedX) && y == int(blockedY) {
				continue
			}

			if !u.canMoveTo(uint8(x), uint8(y), bs) {
				continue
			}

			score := u.calculateDetourScore(uint8(x), uint8(y))
			if score < bestScore {
				bestScore = score
				bestX, bestY = x, y
			}
		}
	}

	return uint8(bestX), uint8(bestY), true
}

func (u *unit) calculateDetourScore(x, y uint8) float64 {
	distToTarget := math.Abs(float64(u.TargetX-x)) + math.Abs(float64(u.TargetY-y))
	mainDirection := u.getMainDirection()
	directionBonus := 0.0

	dx := x - u.X
	dy := y - u.Y

	if (mainDirection.X > 0 && dx > 0) || (mainDirection.X < 0 && dx < 0) {
		directionBonus -= 0.5
	}
	if (mainDirection.Y > 0 && dy > 0) || (mainDirection.Y < 0 && dy < 0) {
		directionBonus -= 0.5
	}

	return distToTarget + directionBonus
}

func (u *unit) getMainDirection() point {
	dx := 0
	dy := 0

	if u.TargetX > u.X {
		dx = 1
	} else if u.TargetX < u.X {
		dx = -1
	}

	if u.TargetY > u.Y {
		dy = 1
	} else if u.TargetY < u.Y {
		dy = -1
	}
	// @reminder: to się wywali, bo powinno być -1, a będzie zero/255
	return point{X: uint8(dx), Y: uint8(dy)}
}

func (u *unit) handlePersistentBlock() {
	u.BlockedCounter++

	if u.BlockedCounter > maxBlockedTicks {
		u.setIdleWithReason("persistently blocked")
	} else {
		u.State = stateWaiting
		u.Delay = uint16(25 + rng.Intn(15))

		if u.BlockedCounter%10 == 0 {
			u.invalidatePathForRecalculation()
		}
	}
}

func (u *unit) invalidatePathForRecalculation() {
	u.Path = nil
	u.PathIndex = 0
	// log.Printf("unit %d: path invalidated for recalculation", u.ID)
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

func (u *unit) setIdleWithReason(reason string) {
	log.Printf("unit %d going IDLE: %s", u.ID, reason)
	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdIdle
	u.clearPath()
	u.BlockedCounter = 0
	u.AllowFriendlyFire = false

	if u.State != stateWaiting {
		u.IsInQueue = false
	}
}

func (u *unit) handleTargetReached(bs *battleState) {
	u.clearPath()

	switch u.Command {
	case cmdAttack:
		u.State = stateAttacking
		u.attack(bs)
	case cmdRepairStructure:
		u.State = stateRepairing
		u.repair(bs)
	default:
		u.setIdle()
	}
}

func (bs *battleState) assignGroupCommand(
	command uint16, mainTargetX, mainTargetY uint8, mainTargetID uint,
	selectedUnits []*unit,
) {
	if len(selectedUnits) == 0 {
		return
	}

	targetX, targetY := bs.resolveActualTarget(mainTargetX, mainTargetY, mainTargetID)

	if len(selectedUnits) <= 4 {
		bs.assignSmallGroupTargets(selectedUnits, command, targetX, targetY, mainTargetID)
		return
	}

	bs.assignScatteredGroupTargets(selectedUnits, command, targetX, targetY, mainTargetID)
}

func (bs *battleState) resolveActualTarget(mainTargetX, mainTargetY uint8, mainTargetID uint) (uint8, uint8) {
	if mainTargetID == 0 {
		return mainTargetX, mainTargetY
	}

	targetUnit, targetBuilding := getObjectByID(mainTargetID, bs)

	if targetUnit != nil && targetUnit.Exists {
		return targetUnit.X, targetUnit.Y
	}

	if targetBuilding != nil && targetBuilding.Exists {
		centerX, centerY, ok := targetBuilding.getCenter()
		if ok {
			return centerX, centerY
		}
	}

	return mainTargetX, mainTargetY
}

func (bs *battleState) assignSmallGroupTargets(units []*unit, command uint16, targetX, targetY uint8, targetID uint) {
	for _, unit := range units {
		unit.addUnitCommand(command, targetX, targetY, targetID, bs)
	}
}

func (bs *battleState) assignScatteredGroupTargets(units []*unit, command uint16, targetX, targetY uint8, targetID uint) {
	positions := bs.generateFormationPositions(targetX, targetY, uint8(len(units)))

	for i, unit := range units {
		assignedX, assignedY := targetX, targetY

		if i < len(positions) {
			assignedX = positions[i].X
			assignedY = positions[i].Y
		}

		unit.addUnitCommand(command, assignedX, assignedY, targetID, bs)
	}
}

func (bs *battleState) generateFormationPositions(centerX, centerY, count uint8) []point {
	positions := make([]point, 0, count)
	cols := uint8(math.Sqrt(float64(count))) + 1

	for i := uint8(0); i < count; i++ {
		row := i / cols
		col := i % cols

		offsetX := col - cols/2
		offsetY := row - count/(cols*2)

		x := centerX + offsetX
		y := centerY + offsetY

		if x < boardMaxX && y < boardMaxY && isWalkable(bs, x, y) {
			positions = append(positions, point{X: x, Y: y})
		} else {
			positions = append(positions, point{X: centerX, Y: centerY})
		}
	}

	return positions
}

func (u *unit) canMoveTo(x, y uint8, bs *battleState) bool {
	if x >= boardMaxX || y >= boardMaxY {
		return false
	}

	tile := &bs.Board.Tiles[x][y]

	// Kolizja z jednostkami (standardowo)
	if tile.Unit != nil && tile.Unit.ID != u.ID {
		return false
	}

	// Kolizja z terenem/budynkami (używamy nowej funkcji z pathfinding.go)
	// Przekazujemy 'u', aby obsłużyć wyjątek krowy wchodzącej do obory
	return isWalkableUnit(bs, x, y, u)
}

// calculateMilkingSpot oblicza milking spot dla obory
func calculateMilkingSpot(bld *building) (uint8, uint8, bool) {
	if len(bld.OccupiedTiles) == 0 {
		return 0, 0, false
	}

	minX := uint8(math.MaxUint8)
	maxY := uint8(0)

	for _, tile := range bld.OccupiedTiles {
		if tile.X < minX {
			minX = tile.X
		}
		if tile.Y > maxY {
			maxY = tile.Y
		}
	}

	return minX, maxY, true
}

// executeMove wykonuje ruch na nową pozycję.
func (u *unit) executeMove(x, y uint8, bs *battleState) {
	// ZMIANA: Używamy nowej struktury Tiles
	// Usuń z poprzedniej pozycji (jeśli to ta jednostka tam jest)
	if bs.Board.Tiles[u.X][u.Y].Unit == u {
		bs.Board.Tiles[u.X][u.Y].Unit = nil
	}

	oldX, oldY := u.X, u.Y

	u.X, u.Y = x, y

	// Ustaw na nowej pozycji
	bs.Board.Tiles[u.X][u.Y].Unit = u

	u.PathIndex++
	u.updateMovementAnimation(oldX, oldY)
}

func (u *unit) setIdle() {
	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdIdle
	u.clearPath()
	u.BlockedCounter = 0
	u.RetryAttempts = 0
	u.PathfindingCooldown = 0
	u.AllowFriendlyFire = false
	if u.State != stateWaiting {
		u.IsInQueue = false
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

func (u *unit) updateMovementAnimation(prevX, prevY uint8) {
	dx := int(u.X) - int(prevX)
	dy := int(u.Y) - int(prevY)

	if dx != 0 || dy != 0 {
		u.Direction = rl.NewVector2(float32(dx), float32(dy))
	}

	u.AnimationType = "walk"
	u.AnimationCounter++

	if u.AnimationCounter >= animationSpeed {
		u.AnimationCounter = 0
		u.AnimationFrame++

		if u.AnimationFrame >= walkAnimationFrames {
			u.AnimationFrame = 0
		}
	}
}

// attack zadaje obrażenia celowi lub ustawia ruch w jego kierunku.
func (u *unit) attack(bs *battleState) {
	target, err := u.validateAttackTarget(bs)
	if err != nil {
		u.setIdleWithReason(err.Error())

		return
	}

	if u.canAttackTarget(target) {
		// Sprawdzanie przerwy
		if u.AttackCooldown > 0 {
			u.State = stateIdle
			u.AnimationType = "idle"

			// Obracamy jednostkę w stronę celu, żeby nie stała bokiem/tyłem
			u.faceTarget(target)

			// Ustawiamy minimalny Delay (1), aby w następnej klatce
			// znów wejść do tej funkcji i sprawdzić, czy cooldown już minął.
			u.Delay = 1
			return
		}

		u.performDirectAttack(target, bs)
	} else {
		// Jeśli cel oddalił się, gonimy go
		u.startMoveToAttack(bs)
	}
}

func (u *unit) faceTarget(target *combatTarget) {
	var tx, ty uint8

	if target.Unit != nil {
		tx, ty = target.Unit.X, target.Unit.Y
	} else if target.Building != nil {
		// Dla budynków celujemy w ich środek lub najbliższy punkt
		tx, ty, _ = target.Building.getClosestOccupiedTile(u.X, u.Y)
	} else {
		return
	}

	dx := 0
	dy := 0

	if tx > u.X {
		dx = 1
	} else if tx < u.X {
		dx = -1
	}

	if ty > u.Y {
		dy = 1
	} else if ty < u.Y {
		dy = -1
	}

	// Ustawiamy wektor kierunku
	if dx != 0 || dy != 0 {
		u.Direction = rl.NewVector2(float32(dx), float32(dy))
	}
}

func (u *unit) validateAttackTarget(
	bs *battleState,
) (*combatTarget, error) {
	target, err := u.validateTargetExists(bs)
	if err != nil {
		return nil, fmt.Errorf("cel zniknął")
	}

	isFriendly := (target.Unit != nil && target.Unit.Owner == u.Owner) ||
		(target.Building != nil && target.Building.Owner == u.Owner)

	if isFriendly && !u.AllowFriendlyFire {
		return nil, fmt.Errorf("atak na jednostkę sojuszniczą niedozwolony")
	}

	if target.Building != nil &&
		target.Building.Type == buildingPalisade &&
		!canDamagePalisades(u) {
		return nil, fmt.Errorf("jednostka nie może niszczyć palisad")
	}

	return target, nil
}

func (u *unit) canAttackTarget(target *combatTarget) bool {
	distance := u.calculateDistanceToTarget(target)
	return distance <= u.AttackRange
}

func (u *unit) performDirectAttack(target *combatTarget, bs *battleState) {
	damageBonus, _, _ := u.getExperienceBonus()
	totalDamage := u.Damage + uint16(damageBonus)

	// STAŁE ANIMACJI (czas blokady ruchu w tyknięciach)
	// 12 tyknięć to ok. 0.6 sekundy (przy 20 TPS logicznych) - wystarczy na animację strzału
	const rangedAnimationDelay = 12
	// 10 tyknięć na zamach mieczem
	const meleeAnimationDelay = 10

	if u.AttackRange > 1 {
		// === DYSTANSOWE (RANGED) ===

		var missileType string
		switch u.Type {
		case unitMage:
			missileType = missileGhost
		case unitPriest:
			missileType = missileFire
		case unitSpearman:
			missileType = missileSpear
		case unitCrossbowman:
			missileType = missileBolt
		case unitPriestess:
			missileType = missileLightning
		default:
			missileType = missileArrow
		}
		var ok bool
		targetX, targetY := uint8(0), uint8(0)

		if target.Unit != nil {
			targetX = target.Unit.X
			targetY = target.Unit.Y
		} else if target.Building != nil {
			targetX, targetY, ok = target.Building.getClosestOccupiedTile(u.X, u.Y)
		}

		if ok {
			// Stwórz pocisk
			proj := &projectile{}
			proj.initProjectile(missileType, u.Owner, u.X, u.Y, targetX, targetY, totalDamage)
			bs.Projectiles = append(bs.Projectiles, proj)

			log.Printf("unit %d fired missile %s at (%d,%d)", u.ID, missileType, targetX, targetY)

			// 1. Ustawienie czasu PRZEŁADOWANIA (Cooldown)
			// Jednostka może się ruszać, ale nie może strzelać
			switch u.Type {
			case unitCrossbowman:
				u.AttackCooldown = 70
			case unitMage, unitPriest, unitPriestess:
				u.AttackCooldown = 100
			case unitArcher, unitSpearman:
				u.AttackCooldown = 65
			default:
				u.AttackCooldown = 65
			}

			// 2. Ustawienie czasu ANIMACJI (Delay)
			// Jednostka jest "zamrożona" tylko na chwilę wypuszczenia strzału
			u.Delay = rangedAnimationDelay
		}
	} else {
		// === ATAK WRĘCZ ===

		// Zadajemy obrażenia natychmiast (lub można opóźnić o połowę Delay)
		if target.Unit != nil {
			target.Unit.takeDamage(totalDamage, bs)
			u.gainExperience(target.Unit, nil, bs)
		} else if target.Building != nil {
			target.Building.takeDamage(totalDamage)
			u.gainExperience(nil, target.Building, bs)
		}

		// 1. Ustawienie czasu ODPOCZYNKU (Cooldown)
		// Balansowanie szybkości ataku różnych jednostek wręcz
		switch u.Type {
		case unitBear:
			u.AttackCooldown = 35 // Niedźwiedź bije wolno, ale mocno
		case unitUnknown: // Strzyga
			u.AttackCooldown = 15 // Bardzo szybka
		default: // Drwal, Miecznik
			u.AttackCooldown = 25 // Standard
		}

		// 2. Ustawienie czasu ANIMACJI (Delay)
		// Krótka blokada na zamach
		u.Delay = meleeAnimationDelay
	}

	u.handleTargetPostAttack(target.Unit, target.Building)
}

func (u *unit) handleTargetPostAttack(targetUnit *unit, targetBld *building) {
	if (targetUnit != nil && !targetUnit.Exists) || (targetBld != nil && !targetBld.Exists) {
		u.setIdleWithReason("cel zniszczony")
	} else {
		u.State = stateAttacking
		u.AnimationType = "fight"
		u.AnimationFrame = 0
	}
}

func (u *unit) transitionToMovementForAttack(effectiveMoveTargetX, effectiveMoveTargetY uint8) {
	u.State = stateMoving
	u.AnimationType = "walk"
	u.Command = cmdAttack
	u.TargetX = effectiveMoveTargetX
	u.TargetY = effectiveMoveTargetY
	u.clearPath()
}

func (u *unit) gainExperience(targetUnit *unit, targetBuilding *building, bs *battleState) {
	if u.Experience >= 235 {
		return
	}
	isEnemyUnit := false
	isEnemyBuilding := false

	if targetUnit != nil && targetUnit.Owner != u.Owner {
		isEnemyUnit = true
	}
	if targetBuilding != nil && targetBuilding.Owner != u.Owner {
		isEnemyBuilding = true
	}

	canGainExp := false

	if u.Owner == bs.PlayerID {
		canGainExp = isEnemyUnit
	} else {
		canGainExp = isEnemyBuilding || isEnemyUnit
	}

	if canGainExp {
		u.Experience++

		if u.Type == unitPriestess || u.Type == unitPriest || u.Type == unitMage {
			u.Experience += 3
		}

		if u.Experience > 235 {
			u.Experience = 235
		}
	}
}

func (u *unit) getExperienceBonus() (damageBonus, armorBonus uint8, manaBonus uint16) {
	tier := u.Experience / 16
	if tier >= 15 {
		tier = 14
	}

	return dDamage[tier], dArmor[tier], dMana[tier]
}

func (u *unit) repair(bs *battleState) {
	// 1. Sprawdzamy, czy istnieje
	_, targetBuilding := getObjectByID(u.TargetID, bs)

	if targetBuilding == nil || !targetBuilding.Exists || targetBuilding.HP >= targetBuilding.MaxHP {
		u.State = stateIdle
		u.AnimationType = "idle"
		u.Command = cmdIdle

		return
	}
	// 2. Szukamy drogi do celu
	distance := targetBuilding.getDistanceToUnit(u.X, u.Y)

	var amount uint16

	switch u.Owner {
	case bs.PlayerID:
		amount = repairAmountPlayer
	case bs.AIPlayerID:
		amount = repairAmountAI
	}
	// Wydaje mi się, że powinniśmy przylegać do budynku, ale <= 2 działa
	if distance == 1 {
		targetBuilding.applyWork(amount, bs)
	}
}

//goland:noinspection SpellCheckingInspection
func (u *unit) castMagicShield() {
	if u.Mana >= spellCostMagicShield {
		u.Mana -= spellCostMagicShield
		u.Effects = append(u.Effects, activeEffect{Name: "MagicShield", Duration: 300})
		log.Printf("unit %d: Aktywowano tarczę magiczną!", u.ID)
	}
	// @todo: jak zrobić aby takedamage brało tarcze pod uwagę?
	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdIdle
}

//func (u *unit) castMagicFire(bs *battleState) {
//	if u.Mana < 80 {
//		u.Command = cmdIdle
//		return
//	}
//	u.Mana -= 80
//	log.Printf("unit %d: Rzuca Deszcz Ognia na (%d,%d)!", u.ID, u.TargetX, u.TargetY)
//
//	bonusDamage := 0
//	if u.Experience/16 < len(dDamage) {
//		bonusDamage = dDamage[u.Experience/16]
//	}
//	totalDamage := u.Damage + bonusDamage
//
//	mainTargetX := u.TargetX
//	mainTargetY := u.TargetY
//	for dx := -1; dx <= 1; dx++ {
//		applyFireDamage(bs, mainTargetX+dx, mainTargetY, totalDamage, u.Owner)
//	}
//
//	vecX, vecY := 0, 0
//	if mainTargetX > u.X {
//		vecX = 1
//	} else if mainTargetX < u.X {
//		vecX = -1
//	}
//	if mainTargetY > u.Y {
//		vecY = 1
//	} else if mainTargetY < u.Y {
//		vecY = -1
//	}
//
//	spread1X, spread1Y := mainTargetX+vecX, mainTargetY+vecY
//	applyFireDamage(bs, spread1X, spread1Y, totalDamage-10, u.Owner)
//
//	spread2X, spread2Y := spread1X+vecX, spread1Y+vecY
//	applyFireDamage(bs, spread2X, spread2Y, totalDamage-20, u.Owner)
//
//	u.State = stateIdle
//	u.Command = cmdIdle
//}

//func applyFireDamage(bs *battleState, x, y, damage int, casterOwner uint8) {
//	if x < 0 || x >= boardMaxX || y < 0 || y >= boardMaxY {
//		return
//	}
//
//	// ZMIANA: Użycie Tiles
//	tile := &bs.Board.Tiles[x][y]
//	tile.EffectID = 100
//
//	// Użycie wskaźników z Tile
//	if tile.Unit != nil && tile.Unit.Owner != casterOwner {
//		tile.Unit.takeDamage(damage, bs)
//	}
//	if tile.Building != nil && tile.Building.Owner != casterOwner {
//		tile.Building.takeDamage(damage)
//	}
//}

func (u *unit) castMagicSight(bs *battleState) {
	if u.Mana >= spellCostMagicSight {
		u.Mana -= spellCostMagicSight
		log.Printf("Jednostka %d rzuca czar widzenia", u.ID)

		revealRadius := spellRangeMagicSight
		for i := u.X - revealRadius; i <= u.X+revealRadius; i++ {
			for j := u.Y - revealRadius; j <= u.Y+revealRadius; j++ {
				if i <= boardMaxX && j <= boardMaxY {
					// @todo: czemu 18?!
					if math.Abs(float64(u.X-i))+math.Abs(float64(u.Y-j)) < 18 {
						bs.Board.Tiles[i][j].Visibility = visibilityVisible
					}
				}
			}
		}
	}

	log.Printf("unit %d: Rzucono Czar Widoczności. Cała mapa odkryta.", u.ID)

	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdIdle
}

func handleMagicShield(u *unit) {
	hasShield := false
	// @todo: ogarnij czy są jeszcze jakieś inne efekty, bo póki co tylko MT
	// jeżeli tak, to warto wyciągnąć przechodzenie po u.Effects a sprawdzenie
	// typu dać w przełączniku!
	for _, effect := range u.Effects {
		// @todo: używanie łańcucha znaków w taki sposób, to zbrodnia na programie!
		// dodaj coś w stylu „EFFECT_MAGIC_SHIELD” int.
		if effect.Name == "MagicShield" {
			hasShield = true
		}

		if hasShield {
			log.Printf("Jednostka %d: Magiczna tarcza pochłonęła obrażenia!", u.ID)

			return
		}
	}
}

func (u *unit) takeDamage(damage uint16, bs *battleState) {
	// @todo: w ogóle nie ogarnięty temat.
	if u.Type == unitPriestess {
		handleMagicShield(u)
	}

	var finalDamage uint16

	if damage <= u.Armor {
		finalDamage = 1
	} else {
		finalDamage = damage - u.Armor
	}

	if finalDamage <= u.HP {
		u.HP -= finalDamage
	} else {
		u.HP = 0
	}

	log.Printf("unit %d otrzymał %d obrażeń. HP: %d/%d", u.ID, finalDamage, u.HP, u.MaxHP)
	// === Zapisywanie informacji o odniesionych ranach

	if len(u.Wounds) < maxWoundsCount {
		offX := float32(rand.Intn(11) - 4)
		offY := float32(rand.Intn(9) - 3)

		isSevere := false
		// @todo zastąp to prostym sprawdzeniem Armor dla ranionej jednostki
		// lub % MaxHP, żeby to było wizualnie czytelniejsze w trakcie bitwy.
		// Za maxhp przemawia, że w finalDamage już jest ukryta poprawka na u.Armor
		isUnitLight := u.Type != unitCommander && u.Type != unitSpearman && u.Type != unitSwordsman &&
			u.Type != unitBear && u.Type != unitUnknown && u.Type != unitCrossbowman

		var baseScale float32 = 1.0

		if finalDamage > severeDamage && isUnitLight {
			isSevere = true
			baseScale = 1.1
		}
		// @todo: skala powinna zależeć od ilości obrażeń i tylko dla isSevere

		rotation := rand.Float32() * 120.0
		newWound := wound{
			Timer:    20,
			OffsetX:  offX,
			OffsetY:  offY,
			IsSevere: isSevere,
			Scale:    baseScale,
			Rotation: rotation,
		}

		u.Wounds = append(u.Wounds, newWound)
	}

	if u.Type == unitCow && u.Exists {
		if u.Udder < 100 && u.Command != cmdFlee {
			barnX, barnY, foundBarn := findNearestBarnMilkingSpot(u, bs)
			if foundBarn {
				u.addUnitCommand(cmdFlee, barnX, barnY, 0, bs)
				log.Printf("unit %d (COW): Otrzymała obrażenia, uciekam do obory na (%d,%d).", u.ID, barnX, barnY)
			} else {
				log.Printf("unit %d (COW): Otrzymała obrażenia, ale nie znalazła obory do ucieczki. "+
					"Pozostaję przy obecnej komendzie.", u.ID)
			}
		}
	}

	// Sprawdzamy, czy jednostka została zabita
	if u.HP == 0 {
		u.Exists = false

		// Zabita jednostka nie powinna zliczać się do górnej granicy ludności
		decreasePopulation(u, bs)

		tile := &bs.Board.Tiles[u.X][u.Y]
		if tile.Unit == u {
			tile.Unit = nil
		}

		createCorpses(u, bs)
		u.unregisterFromBuilding()
		log.Printf("Jednostka %d została zabita!", u.ID)
	}
}

func decreasePopulation(u *unit, bs *battleState) {
	switch u.Owner {
	case bs.HumanPlayerState.PlayerID:
		bs.HumanPlayerState.CurrentPopulation--
	case bs.AIEnemyState.PlayerID:
		bs.AIEnemyState.CurrentPopulation--
	}
}

func createCorpses(u *unit, bs *battleState) {
	steps := 18
	stepAngle := 10
	rotation := float32(rand.Intn(steps) * stepAngle)
	corpse := corpse{
		X:          u.X,
		Y:          u.Y,
		UnitType:   u.Type,
		DecayTimer: corpseDecayTime,
		Phase:      0,
		Rotation:   rotation,
		Owner:      u.Owner,
	}

	// @todo: ogarnąć jakie powinny być zwłoki
	// coś tam pod spriteBtnBuildAcademy o współrzędnych
	switch corpse.UnitType {
	case unitCow:
		corpse.SkeletonType = 1
	case unitBear:
		corpse.SkeletonType = 2
	default:
		corpse.SkeletonType = 0
	}

	bs.Corpses = append(bs.Corpses, corpse)
}

func (u *unit) unregisterFromBuilding() {
	if u.BelongsTo != nil {
		log.Printf("Jednostka %d próbuje się wyrejestrować z budynku %d", u.ID, u.BelongsTo.ID)
		u.BelongsTo.unregisterUnit(u.ID)
		u.BelongsTo = nil
	}
}

func (u *unit) findNearestPalisade(bs *battleState, radius uint8,
) *building {
	var best *building
	minD := math.MaxFloat64

	for _, pal := range bs.Buildings {
		if pal == nil || !pal.Exists || pal.Type != buildingPalisade {
			continue
		}

		cx, cy, _ := pal.getCenter()
		dx := math.Abs(float64(u.X - cx))
		dy := math.Abs(float64(u.Y - cy))
		d := math.Max(dx, dy)

		if d <= float64(radius) && d < minD {
			if u.isImportantPalisade(pal, bs) {
				minD = d
				best = pal
			}
		}
	}

	return best
}

func (u *unit) actOnIdle(bs *battleState) {
	if !u.canActOnIdle() {
		return
	}

	if u.BlockedCounter > 0 {
		return
	}

	if !u.shouldSearchForTarget(bs) {
		return
	}

	u.handleTargetSearch(bs)
}

func (u *unit) canActOnIdle() bool {
	return u.Type != unitCow && u.Type != unitShepherd
}

func (u *unit) shouldSearchForTarget(bs *battleState) bool {
	return u.isReadyToAct(bs)
}

func (u *unit) isReadyToAct(bs *battleState) bool {
	if u.State == stateIdle && u.Command == cmdIdle {
		return true
	}

	if u.Command == cmdAttack {
		_, err := u.validateTargetExists(bs)
		if err != nil {
			return true
		}
	}

	return false
}

const palisadeStrategicBuildingProximity = 10

func (u *unit) isImportantPalisade(palisade *building, bs *battleState) bool {
	if u.Owner != bs.AIPlayerID || !canDamagePalisades(u) {
		return false
	}

	if palisade == nil || !palisade.Exists || palisade.Type != buildingPalisade {
		return false
	}

	const prox = palisadeStrategicBuildingProximity

	palCenterX, palCenterY, _ := palisade.getCenter()

	for _, bld := range bs.Buildings {
		if bld == nil || !bld.Exists || bld.Owner == u.Owner || bld.Type == buildingPalisade || bld.ID == palisade.ID {
			continue
		}

		bldCenterX, bldCenterY, ok := bld.getCenter()
		if !ok {
			continue
		}

		distToPalisadeCenter := math.Max(math.Abs(float64(palCenterX-bldCenterX)), math.Abs(float64(palCenterY-bldCenterY)))
		if distToPalisadeCenter > float64(prox) {
			continue
		}

		_, _, ok = bld.getClosestWalkableTile(bs)

		if !ok {
			return true
		}
	}

	return false
}

func (u *unit) handleTargetSearch(bs *battleState) {
	if u.Owner == bs.HumanPlayerState.PlayerID {
		u.handleTargetSearchForHumanPlayer(bs)
	} else {
		u.handleTargetSearchForAI(bs)
	}
}

func (u *unit) handleTargetSearchForHumanPlayer(bs *battleState) {
	primaryTargetUnit, primaryTargetBuilding, foundPrimary := findNearestEnemyExtended(u, bs)

	if !foundPrimary {
		u.setIdle()

		return
	}

	if primaryTargetUnit != nil && primaryTargetUnit.Exists {
		u.handleUnitTarget(primaryTargetUnit, bs)

		return
	}

	if primaryTargetBuilding != nil && primaryTargetBuilding.Exists {
		u.setIdle()

		return
	}

	u.setIdle()
}

func (u *unit) handleTargetSearchForAI(bs *battleState) {
	isPalisadeBreaker := canDamagePalisades(u)

	primaryTargetUnit, primaryTargetBuilding, foundPrimary := findNearestEnemyExtended(u, bs)

	if isPalisadeBreaker && foundPrimary && primaryTargetBuilding != nil && primaryTargetBuilding.Exists {
		_, _, ok := primaryTargetBuilding.getClosestWalkableTile(bs)

		if !ok {
			palisadeTarget := u.findNearestPalisade(bs, u.SightRange)

			if palisadeTarget != nil {
				u.handleBuildingTarget(palisadeTarget, bs)

				return
			}
		}
	}

	if !foundPrimary {
		u.setIdle()
		return
	}

	if primaryTargetUnit != nil {
		u.handleUnitTarget(primaryTargetUnit, bs)
	} else {
		u.handleBuildingTarget(primaryTargetBuilding, bs)
	}
}

func (u *unit) handleUnitTarget(targetUnit *unit, bs *battleState) {
	u.TargetID = targetUnit.ID
	u.setMoveTargetForUnit(targetUnit, bs)

	u.executeActionBasedOnDistance(bs)
}

func (u *unit) handleBuildingTarget(targetBuilding *building, bs *battleState) {
	u.TargetID = targetBuilding.ID

	currentDistanceToBuilding := targetBuilding.getDistanceToUnit(u.X, u.Y)

	if currentDistanceToBuilding != math.MaxUint8 && currentDistanceToBuilding <= u.AttackRange {
		u.startDirectAttack(u.X, u.Y, bs)
		return
	}

	finalMoveTargetX, finalMoveTargetY := uint8(0), uint8(0)

	optimalRangedX, optimalRangedY := uint8(0), uint8(0)
	closestWalkableX, closestWalkableY := uint8(0), uint8(0)
	var ok bool

	if u.AttackRange > 1 {
		optimalRangedX, optimalRangedY, ok = findOptimalRangedAttackTile(u.X, u.Y, targetBuilding, u.AttackRange, bs)
	}

	if !ok {
		finalMoveTargetX, finalMoveTargetY = optimalRangedX, optimalRangedY
	} else {
		closestWalkableX, closestWalkableY, ok = targetBuilding.getClosestWalkableTile(bs)
		if ok {
			finalMoveTargetX, finalMoveTargetY = closestWalkableX, closestWalkableY
		}
	}

	if !ok {
		u.setIdle()
		return
	}

	u.TargetX = finalMoveTargetX
	u.TargetY = finalMoveTargetY
	u.startMoveToAttack(bs)
	log.Printf("DEBUG_AI: unit %d (Type:%d) moving to attack bld %d, target tile: (%d,%d). Current position: (%d,%d).",
		u.ID, u.Type, u.TargetID, u.TargetX, u.TargetY, u.X, u.Y)
}

func (u *unit) setMoveTargetForUnit(targetUnit *unit, bs *battleState) {
	bestX, bestY := u.findBestPositionAroundUnit(targetUnit, bs)
	u.TargetX, u.TargetY = bestX, bestY
}

func (u *unit) findBestPositionAroundUnit(targetUnit *unit, bs *battleState) (uint8, uint8) {
	bestX, bestY := int(targetUnit.X), int(targetUnit.Y)
	minDist := math.MaxFloat64
	foundFreeSpot := false

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			checkX := int(targetUnit.X) + dx
			checkY := int(targetUnit.Y) + dy

			if u.isValidMoveTarget(uint8(checkX), uint8(checkY), bs) {
				dist := math.Abs(float64(int(u.X)-checkX)) + math.Abs(float64(int(u.Y)-checkY))
				if dist < minDist {
					minDist = dist
					bestX, bestY = checkX, checkY
					foundFreeSpot = true
				}
			}
		}
	}

	if !foundFreeSpot {
		return targetUnit.X, targetUnit.Y // Fallback
	}

	return uint8(bestX), uint8(bestY)
}

func (u *unit) isValidMoveTarget(x, y uint8, bs *battleState) bool {
	return x < boardMaxX && y < boardMaxY &&
		// ZMIANA: Użycie Tiles
		bs.Board.Tiles[x][y].Unit == nil &&
		bs.Board.Tiles[x][y].Building == nil &&
		isWalkable(bs, x, y)
}

func (u *unit) validateTargetExists(bs *battleState) (*combatTarget, error) {
	targetUnit, targetBuilding := getObjectByID(u.TargetID, bs)
	if (targetUnit == nil || !targetUnit.Exists) && (targetBuilding == nil || !targetBuilding.Exists) {
		return nil, fmt.Errorf("cel nie istnieje")
	}

	return &combatTarget{Unit: targetUnit, Building: targetBuilding}, nil
}

func (u *unit) calculateDistanceToTarget(target *combatTarget) uint8 {
	if target.Unit != nil {
		return uint8(math.Max(
			math.Abs(float64(u.X-target.Unit.X)),
			math.Abs(float64(u.Y-target.Unit.Y)),
		))
	}

	// log.Println("wchodzę w target.Building.getDistanceToUnit")

	return target.Building.getDistanceToUnit(u.X, u.Y)
}

func (u *unit) executeActionByDistance(distance uint8, bs *battleState) {
	if distance > u.SightRange {
		log.Printf("DEBUG_AI: U %d: cel ID %d poza zasięgiem widzenia. IDLE", u.ID, u.TargetID)
		u.setIdle()
		return
	}

	if distance <= u.AttackRange {
		log.Printf("DEBUT_AI: U %d: Odległość <= zasięg ataku. Rozpoczynam bezpośredni atak na cel ID %d.",
			u.ID, u.TargetID)
		u.startDirectAttack(u.TargetX, u.TargetY, bs)
	} else {
		log.Printf("DEBUG_AI: U %d: odległość %d > zasięg ataku %d. Ruszam w kierunku %d.",
			u.ID, distance, u.AttackRange, u.TargetID)
		u.startMoveToAttack(bs)
	}
}

func (u *unit) executeActionBasedOnDistance(bs *battleState) {
	target, err := u.validateTargetExists(bs)
	if err != nil {
		u.setIdle()
		return
	}

	distance := u.calculateDistanceToTarget(target)
	u.executeActionByDistance(distance, bs)
}

func (u *unit) startDirectAttack(placeholderX, placeholderY uint8, bs *battleState) {
	realTargetX := placeholderX
	realTargetY := placeholderY

	if u.TargetID != 0 {
		targetUnit, targetBld := getObjectByID(u.TargetID, bs)

		if targetUnit != nil && targetUnit.Exists {
			realTargetX = targetUnit.X
			realTargetY = targetUnit.Y
		} else if targetBld != nil && targetBld.Exists {
			bx, by, ok := targetBld.getClosestOccupiedTile(u.X, u.Y)
			if ok {
				realTargetX = bx
				realTargetY = by
			}
		}
	}

	u.addUnitCommand(cmdAttack, realTargetX, realTargetY, u.TargetID, bs)

	u.State = stateAttacking
	u.AnimationType = "fight"
	u.AnimationFrame = 0

	log.Printf("DEBUG_AI: unit %d attacking target %d at real pos (%d,%d)",
		u.ID, u.TargetID, realTargetX, realTargetY)
}

func (bld *building) getClosestOccupiedTile(fromX, fromY uint8) (uint8, uint8, bool) {
	if len(bld.OccupiedTiles) == 0 {
		return 0, 0, false
	}

	closestX, closestY := uint8(0), uint8(0)
	minDistSq := math.MaxFloat64

	for _, tile := range bld.OccupiedTiles {
		dx := float64(tile.X - fromX)
		dy := float64(tile.Y - fromY)
		distSq := dx*dx + dy*dy

		if distSq < minDistSq {
			minDistSq = distSq
			closestX = tile.X
			closestY = tile.Y
		}
	}
	return closestX, closestY, true
}

func (u *unit) startMoveToAttack(bs *battleState) {
	u.addUnitCommand(cmdAttack, u.TargetX, u.TargetY, u.TargetID, bs)
	u.State = stateMoving
	u.AnimationType = "walk"

	log.Printf("DEBUG_AI: unit %d moving to attack target %d, move to (%d,%d)",
		u.ID, u.TargetID, u.TargetX, u.TargetY)
}

func getLegacyUnitIndex(t unitType) int {
	return int(t)
}

// POMOCNICZE FUNKCJE MAPY I OBIEKTÓW

func findOptimalRangedAttackTile(uCurrentX, uCurrentY uint8, bld *building, attackRange uint8, bs *battleState) (uint8, uint8, bool) {
	if len(bld.OccupiedTiles) == 0 {
		return 0, 0, false
	}

	centerX, centerY, ok := bld.getCenter()
	if !ok {
		return 0, 0, false
	}

	candidates := []point{
		{X: centerX + uint8(attackRange), Y: centerY},
		{X: centerX - uint8(attackRange), Y: centerY},
		{X: centerX, Y: centerY + uint8(attackRange)},
		{X: centerX, Y: centerY - uint8(attackRange)},
		{X: centerX + uint8(attackRange), Y: centerY + uint8(attackRange)},
		{X: centerX - uint8(attackRange), Y: centerY - uint8(attackRange)},
		{X: centerX + uint8(attackRange), Y: centerY - uint8(attackRange)},
		{X: centerX - uint8(attackRange), Y: centerY + uint8(attackRange)},
	}

	halfRange := uint8(attackRange / 2)
	if halfRange > 0 {
		candidates = append(candidates,
			point{X: centerX + halfRange, Y: centerY},
			point{X: centerX - halfRange, Y: centerY},
			point{X: centerX, Y: centerY + halfRange},
			point{X: centerX, Y: centerY - halfRange},
		)
	}

	closestX, closestY := uint8(0), uint8(0)
	minDistance := math.MaxFloat64

	for _, candidate := range candidates {
		x, y := candidate.X, candidate.Y

		if !bld.isValidWalkableTile(x, y, bs) {
			continue
		}

		distanceToBuilding := bld.getDistanceToUnit(x, y)

		if distanceToBuilding <= attackRange {
			distanceFromUnit := math.Abs(float64(uCurrentX-x)) + math.Abs(float64(uCurrentY-y))

			if distanceFromUnit < minDistance {
				minDistance = distanceFromUnit
				closestX = x
				closestY = y
			}
		}
	}

	return closestX, closestY, true
}
