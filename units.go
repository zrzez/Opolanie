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

func (u *unit) initUnit(unitType unitType, x, y uint8, command commandType, bState *battleState) {
	u.ID = bState.NextUniqueObjectID
	bState.NextUniqueObjectID++
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
	u.Experience = 0
	u.Delay = u.MaxDelay
	u.AttackCooldown = 0
	u.Wounds = make([]wound, 0, maxWoundsCount)

	stats, ok := unitDefs[unitType]
	if ok {
		u.SightRange = stats.SightRange
		u.AttackRange = stats.AttackRange
		u.Damage = stats.BaseDamage
		u.Armor = stats.BaseArmor
		u.MaxHP = stats.MaxHP
		u.MaxDelay = stats.MoveDelay
		u.HP = stats.MaxHP
		u.MaxMana = stats.MaxMana
		u.Mana = stats.MaxMana / 2
		u.ManaRegen = stats.BaseManaRegen
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
	if !u.Exists || amount == 0 {
		return
	}

	// @reminder: najwięcej many ma unitMage i unitUnknown 15 poziomu, tj. 60+280=340.
	// Czyli nie ma możliwości, aby przyrost many przekroczył górną granicę
	// uint16 i „przekręcił się” w okolice zera. Dlatego nie ma bezpiecznika.
	u.Mana += amount

	if u.Mana > u.MaxMana {
		u.Mana = u.MaxMana
	}
}

// Jeśli to możliwe to dla każdej istniejącej jednostki zmniejsza manę o amount oraz zwraca prawda.
// Jeśli u.Mana < amount, to zwraca fałsz.
func (u *unit) tryToDecreaseMana(amount uint16) bool {
	if !u.Exists {
		return false
	}

	if u.Mana >= amount {
		u.Mana -= amount

		return true
	}

	return false
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
func (u *unit) show(bState *battleState) {
	if u.X < boardMaxX && u.Y < boardMaxY {
		// ZMIANA: Ustawiamy wskaźnik unit w nowej strukturze Tiles
		bState.Board.Tiles[u.X][u.Y].Unit = u
	} else {
		log.Printf("OSTRZEŻENIE: Próba umieszczenia jednostki poza mapą: (%d,%d)", u.X, u.Y)
	}
}

// ============================================================================
// LOGIKA JEDNOSTEK
// ============================================================================

func (u *unit) updateUnit(bState *battleState) {
	// Aktualizowanie ran
	// @todo przenieś do osobnej funkcji, szkoda zajmować tutaj miejsce
	nextFreeIndex := 0

	for scanIndex := range u.Wounds {
		currentWound := &u.Wounds[scanIndex]
		currentWound.Timer--

		if currentWound.Timer > 1 {
			if scanIndex != nextFreeIndex {
				u.Wounds[nextFreeIndex] = *currentWound
			}

			nextFreeIndex++
		}
	}

	u.Wounds = u.Wounds[:nextFreeIndex]

	u.handleAttackCooldown(bState)

	if u.handleNoMovementDetection() {
		return
	}

	if u.handleDelay(bState) {
		return
	}

	// u.processActiveEffects()

	if u.handleBlockedCounter() {
		return
	}

	u.handleWaitingToActiveTransition()
	u.handleMovementTargetReached(bState)
	u.executeCommandAction(bState)
	u.resetDelayIfActive()
}

func (u *unit) handleAttackCooldown(bState *battleState) {
	if bState.GlobalFrameCounter%logicSpeedDivisor == 0 {
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

func (u *unit) handleDelay(bState *battleState) bool {
	if u.State == stateIdle {
		u.Delay = u.MaxDelay

		return false
	}
	if u.Delay <= 0 {
		return false
	}
	if bState.GlobalFrameCounter%logicSpeedDivisor != 0 {
		return true
	}
	u.Delay--

	return u.Delay > 0
}

/*func (u *unit) processActiveEffects() {
	var activeEffects []activeEffect
	for _, e := range u.Effects {
		e.Duration--
		if e.Duration > 0 {
			activeEffects = append(activeEffects, e)
		}
	}

	u.Effects = activeEffects
}*/

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

func (u *unit) handleMovementTargetReached(bState *battleState) {
	if u.State == stateMoving && u.X == u.TargetX && u.Y == u.TargetY {
		u.handleTargetReached(bState)
	}
}

func (u *unit) executeCommandAction(bState *battleState) {
	switch u.Type {
	case unitCow:
		u.handleCowBehavior(bState)
	default:
		u.executeStandardUnitCommand(bState)
	}
}

func (u *unit) executeStandardUnitCommand(bState *battleState) {
	switch u.Command {
	case cmdUMove:
		u.move(bState)
	case cmdUAttack:
		if u.canAttackTargetFromCurrentPosition(bState) {
			u.State = stateAttacking
			u.clearPath()
			u.attack(bState)
		} else {
			u.State = stateMoving
			u.move(bState)
		}
	case cmdUBuild:
		if u.canAttackTargetFromCurrentPosition(bState) {
			u.State = stateBuilding

			if u.AnimationType != "fight" {
				u.AnimationType = "fight"
				u.AnimationFrame = 0
			}

			u.clearPath()
			u.build(bState)
		} else {
			u.State = stateMoving

			if u.AnimationType != "walk" {
				u.AnimationType = "walk"
			}

			u.move(bState)
		}
	case cmdURepair:
		if u.canAttackTargetFromCurrentPosition(bState) {
			u.State = stateRepairing

			if u.AnimationType != "fight" {
				u.AnimationType = "fight"
				u.AnimationFrame = 0
			}

			u.clearPath()
			u.repair(bState)
		} else {
			u.State = stateMoving

			if u.AnimationType != "walk" {
				u.AnimationType = "walk"
			}

			u.move(bState)
		}
	case cmdUCastSpell:
		if u.CurrentSpell == spellMagicShield {
			u.castMagicShield()
		}

		if u.CurrentSpell == spellMagicShower {
			if u.canCastSpellFromCurrentPosition() {
				u.State = stateCastingSpell
				u.AnimationType = "fight"
				u.clearPath()
				u.castSpell(bState)
			} else {
				u.State = stateMoving
				u.AnimationType = "walk"
				u.move(bState)
			}
		}

		if u.CurrentSpell == spellMagicSight {
			u.castMagicSight(bState)
		}

	case cmdUIdle, cmdUStop:
		u.actOnIdle(bState)
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
	case cmdUMove:
		return stateMoving
	case cmdUAttack:
		return stateAttacking
	case cmdUGraze:
		return stateGrazing
	case cmdUBuild:
		return stateBuilding
	case cmdURepair:
		return stateRepairing
	case cmdUCastSpell:
		return stateCastingSpell
	default:
		return stateIdle
	}
}

func getSelectedUnits(bState *battleState) []*unit {
	var selected []*unit

	for _, currentUnit := range bState.Units {
		if currentUnit.IsSelected && currentUnit.Exists && currentUnit.Owner == bState.PlayerID {
			selected = append(selected, currentUnit)
		}
	}

	return selected
}

func canDamagePalisades(unit *unit) bool {
	return unit.Type == unitAxeman || unit.Type == unitPriest
}

func (u *unit) resolveApproachPosition(targetX, targetY *uint8, targetID uint, bState *battleState) (uint8, uint8, error) {
	targetUnit, targetBuilding := getObjectByID(targetID, bState)

	// Cel jest budynkiem
	if targetBuilding != nil && (targetBuilding.Exists || targetBuilding.Type == buildingBridge) {
		if u.AttackRange > 1 {
			x, y, ok := findOptimalRangedAttackTile(u.X, u.Y, u.AttackRange, targetBuilding, bState)
			if ok {
				return x, y, nil
			}
		}

		x, y, ok := targetBuilding.getClosestWalkableTile(bState)
		if ok {
			return x, y, nil
		}

		return 0, 0, fmt.Errorf("cel (budynek) jest nieosiągalny")
	}

	// Cel jest jednostką
	if targetUnit != nil && targetUnit.Exists {
		bestX, bestY := u.findBestPositionAroundUnit(targetUnit, bState)

		// findBestPositionAroundUnit zwraca pozycję celu jako bezpiecznik.
		if bestX == targetUnit.X && bestY == targetUnit.Y {
			// Sprawdź, czy to naprawdę fallback (kafel jest zajęty przez cel)
			targetTile := &bState.Board.Tiles[bestX][bestY]
			if targetTile.Unit == targetUnit {
				return 0, 0, fmt.Errorf("brak wolnego kafelka wokół jednostki ID %d", targetID)
			}
		}

		return bestX, bestY, nil
	}

	targetTile := &bState.Board.Tiles[*targetX][*targetY]

	if isTreeStump(targetTile.TextureID) {
		bestX, bestY, ok := u.findOptimalAttackTileAroundTree(*targetX, *targetY, bState)
		if !ok {
			return 0, 0, fmt.Errorf("nie ma pozycji do ataku tego drzewa")
		}

		return bestX, bestY, nil
	}

	return 0, 0, fmt.Errorf("cel ataku ID %d nie istnieje", targetID)
}

func (u *unit) findOptimalAttackTileAroundTree(treeX, treeY uint8, bState *battleState) (uint8, uint8, bool) {
	var bestX, bestY uint8

	minDistance := math.MaxFloat64

	for column := int(treeX) - int(u.AttackRange); column <= int(treeX)+int(u.AttackRange); column++ {
		if column < 0 || column >= int(boardMaxX) {
			continue // kolumny poza planszą
		}

		for row := int(treeY) - int(u.AttackRange); row <= int(treeY)+int(u.AttackRange); row++ {
			if row < 0 || row >= int(boardMaxY) {
				continue // wiersz poza planszą
			}

			if !isWalkable(bState, uint8(column), uint8(row)) {
				continue // kafelek nieprzechodni
			}

			if column+1 == int(treeX) && row == int(treeY) {
				continue // pomijamy miejsce na które spadnie drzewo
			}

			if column == int(treeX) && row == int(treeY) {
				continue // pomijamy samo drzewo
			}

			electedTile := &bState.Board.Tiles[uint8(column)][uint8(row)]

			if electedTile.Unit != nil && electedTile.Unit.ID != u.ID {
				continue // ktoś już tam stoi
			}

			dx := int(u.X) - column
			dy := int(u.Y) - row

			distance := math.Abs(float64(dx) + math.Abs(float64(dy)))

			if distance < minDistance {
				minDistance = distance
				bestX, bestY = uint8(column), uint8(row)
			}
		}
	}

	if minDistance == math.MaxFloat64 {
		return 0, 0, false
	}

	return bestX, bestY, true
}

func (u *unit) addUnitCommand(cmd *command, bState *battleState) {
	log.Printf("INFO: unit.go dodano rozkaz %d.", cmd.ActionType)
	// ŁATANIE DZIURY W KOMPLETOWANIE ROZKAÓW DLA JEDNOSTEK
	u.CurrentSpell = cmd.Spell
	u.AllowFriendlyFire = cmd.FriendlyFire

	if u.shouldSkipDuplicate(cmd.ActionType, cmd.TargetX, cmd.TargetY, cmd.InteractionTargetID) {
		log.Printf("INFO: unit.go shouldSkipDuplicate %t.", u.shouldSkipDuplicate(cmd.ActionType, cmd.TargetX, cmd.TargetY, cmd.InteractionTargetID))

		return
	}

	if err := u.resolveInteractionTarget(&cmd.TargetX, &cmd.TargetY, cmd.ActionType, cmd.InteractionTargetID, bState); err != nil {
		u.setIdleWithReason("cel nieosiągalny")

		return
	}

	if !u.validateCommand(cmd.ActionType, cmd.InteractionTargetID, bState) {
		log.Printf("INFO: unit.go rozkaz nie przeszedł sprawdzenia %t.", u.validateCommand(cmd.ActionType, cmd.InteractionTargetID, bState))

		return
	}

	u.prepareForNewCommand(cmd.ActionType, cmd.TargetX, cmd.TargetY, cmd.InteractionTargetID)

	// @reminder: Nie powinno to tutaj być, ale lepsze to niż castle.go
	// @todo: ogarnij czemu to się tak porobiło, nie powinno tak to wyglądać!
	// ↓↓↓ NIE DZIAŁA
	// u.CurrentSpell = cmd.Spell
	// ↓↓↓↓ działa
	// u.AllowFriendlyFire = cmd.FriendlyFire

	u.applyCommandState(cmd.ActionType)
}

func (u *unit) shouldSkipDuplicate(command commandType, targetX, targetY uint8, targetID uint) bool {
	return u.Command == command && u.TargetX == targetX &&
		u.TargetY == targetY && u.TargetID == targetID
}

func (u *unit) resolveInteractionTarget(targetX, targetY *uint8, command commandType, targetID uint, bState *battleState) error {
	if !isInteractionCommand(command) {
		return nil
	}

	// Gromobicie oraz deszcz ognia
	if u.CurrentSpell == spellMagicShower {
		// Zapamiętujemy faktyczny cel czaru
		u.interactionTargetX, u.interactionTargetY = *targetX, *targetY

		// Szukamy miejsca, skąd jednostka może rzucić czar (w zasięgu AttackRange)
		finalX, finalY, ok := u.findBestPositionAroundTile(*targetX, *targetY, bState)
		if !ok {
			return fmt.Errorf("brak miejsca w zasięgu czaru")
		}

		// Nadpisujemy TargetX/Y
		*targetX, *targetY = finalX, finalY
		return nil
	}

	if u.CurrentSpell == spellMagicShield || u.CurrentSpell == spellMagicSight {
		*targetX, *targetY = u.X, u.Y
		return nil
	}

	if targetID == 0 {
		u.interactionTargetX, u.interactionTargetY = *targetX, *targetY
	}

	finalX, finalY, err := u.resolveApproachPosition(targetX, targetY, targetID, bState)
	if err != nil {
		return err
	}

	*targetX, *targetY = finalX, finalY

	return nil
}

func isInteractionCommand(command commandType) bool {
	switch command {
	case cmdUAttack, cmdUBuild, cmdURepair, cmdBPlaceConstruction, cmdUCastSpell:
		return true
	default:
		return false
	}
}

func (u *unit) validateCommand(command commandType, targetID uint, bState *battleState) bool {
	switch command {
	case cmdUAttack:
		return u.canAttack(targetID, bState)
	default:
		return true
	}
}

// caDamageTree sprawdza, czy jednostka może zaatakować dane drzewo.
func (u *unit) canDamageTree(treeX, treeY uint8, bState *battleState) bool {
	// Tylko unitAxeman i unitPriest może atakować drzewa
	if u.Type != unitPriest && u.Type != unitAxeman {
		return false
	}

	treeTile := bState.Board.Tiles[treeX][treeY]

	// Tylko stojące drzewa
	if !treeTile.isStandingTree() {
		return false
	}

	// Tylko drzewa, które nie zostały już podpalone
	if treeTile.IsBurning {
		return false
	}

	if u.Type == unitPriest {
		return true
	}

	if u.Type == unitAxeman {
		return treeTile.TextureID == spriteDryTreeStump00
	}

	return false
}

func (u *unit) canAttack(targetID uint, bState *battleState) bool {
	// Drzewa
	// Kapłan może podpalić każde drzewo
	// Drwal może ściąć suche drzewo
	if targetID == 0 {
		return u.canDamageTree(u.interactionTargetX, u.interactionTargetY, bState)
	}

	// targetTile := &bs.Board.Tiles[targetX][targetY]

	targetUnit, targetBuilding := getObjectByID(targetID, bState)

	// Jednostki
	// musi istnieć
	if targetUnit != nil && targetUnit.ID != u.ID {
		// musi być oznaczona jako żywa
		if !targetUnit.Exists {
			return false
		}
		// Istnieje i jest oznaczona jako żywa
		return true
	}

	// Budynki
	if targetBuilding != nil {
		if !targetBuilding.Exists {
			return false
		}

		// magowie nie mogą atakować budynków
		if u.Type == unitMage {
			return false
		}

		if targetBuilding.Type == buildingPalisade && !canDamagePalisades(u) {
			log.Printf("INFO: Jednostka %d nie może atakować palisad", u.ID)

			return false
		}

		if targetBuilding.Type == buildingBridge {
			return false
		}

		return true
	}

	// cel nie istnieje
	return false
}

func (u *unit) prepareForNewCommand(cmdType commandType, targetX, targetY uint8, targetID uint) {
	u.clearPath()
	u.History = nil
	u.LoopCount = 0
	u.TicksNoProgress = 0
	u.LastPathIndex = 0
	u.Command = cmdType
	u.TargetX = targetX
	u.TargetY = targetY
	u.TargetID = targetID
	u.Delay = 0
}

func (u *unit) applyCommandState(command commandType) {
	switch command {
	case cmdUAttack:
		log.Printf("INFO: units.go applyCommandState cmdAttack")

		u.State = stateAttacking
		u.AnimationType = "fight"
		u.AnimationFrame = 3
		u.AnimationCounter = 0
	case cmdUMove, cmdUFlee:
		u.State = stateMoving
		u.AnimationType = "walk"
	case cmdUStop:
		u.State = stateIdle
		u.AnimationType = "walk"
		u.AnimationFrame = 0
		u.Command = cmdUIdle
	case cmdUCastSpell:
		u.State = stateCastingSpell
	case cmdUGraze:
		u.State = stateGrazing
	case cmdUBuild:
		u.State = stateBuilding
	case cmdURepair, cmdBPlaceConstruction: // @todo: czemu do cholery metoda u ma rozkazy B?
		u.State = stateRepairing
	case cmdBMilking:
		u.State = stateMilking
	default:
		fmt.Println("DUUUUPA NIE MA TAKIEJ KOMENDY W PRZEŁĄCZNIKU")
		panic("unhandled default case")
	}
}

func (u *unit) isAtTarget() bool {
	return u.X == u.TargetX && u.Y == u.TargetY
}

func (u *unit) move(bState *battleState) {
	if u.Command == cmdUAttack {
		log.Printf("INFO: units.go move rozkaz to cmdAttack")

		if u.canAttackTargetFromCurrentPosition(bState) {
			log.Printf("INFO: units.go move cel osiągalny z tego miejsca")
			u.clearPath()
			u.State = stateAttacking

			return
		}

		// 25.04.2026 Dodaję bezpiecznik przerywający ruch jeśli cel przestał istnieć
		// Bez tego jednostka atakująca drzewo zaczyna się przemieszczać po jego upadku
		// szukając nowej pozycji do ataku nieistniejącego już celu.
		if _, err := u.validateTargetExists(bState); err != nil {
			u.setIdleWithReason("cel ataku przestał istnieć")

			return
		}
	}

	if u.isAtTarget() {
		log.Printf("INFO: units.go move u celu")

		u.handleTargetReached(bState)

		return
	}

	if u.shouldAbortMovement() {
		u.setIdleWithReason("przerwano ruch")

		return
	}

	u.executeAStarMovement(bState)
}

func (u *unit) canAttackTargetFromCurrentPosition(bState *battleState) bool {
	log.Println("Sprawdzam, czy cel istnieje")

	target, err := u.validateTargetExists(bState)
	if err != nil {
		return false
	}

	log.Println("Obliczam odległość")

	distance := u.calculateDistanceToTarget(target)

	log.Printf("Odległość obliczona %d", distance)
	log.Println(distance <= u.AttackRange)

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

func (u *unit) executeAStarMovement(bState *battleState) {
	if !u.ensureValidPath(bState) {
		return
	}

	u.moveAlongPath(bState)
}

func (u *unit) ensureValidPath(bState *battleState) bool {
	if u.hasValidPath(bState) {
		return true
	}

	if bState.pathfindingUnitsThisTick >= maxPathfindingBudget {
		u.waitForPathfindingBudget()

		return false
	}

	return u.calculateNewPath(bState)
}

func (u *unit) hasValidPath(bState *battleState) bool {
	if len(u.Path) == 0 || u.PathIndex >= len(u.Path) {
		return false
	}

	target, err := u.validateTargetExists(bState)
	if err != nil {
		return false
	}

	if u.TargetID == 0 {
		return true
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

func (u *unit) calculateNewPath(bState *battleState) bool {
	bState.pathfindingUnitsThisTick++

	newPath := findPath(bState, u.ID, u.X, u.Y, u.TargetX, u.TargetY)

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

func (u *unit) moveAlongPath(bState *battleState) {
	if u.PathIndex >= len(u.Path) {
		u.clearPath()

		return
	}

	next := u.Path[u.PathIndex]

	if u.canMoveTo(next.X, next.Y, bState) {
		u.executeSuccessfulMove(next.X, next.Y, bState)
	} else {
		u.handleMovementBlocked(bState, next.X, next.Y)
	}
}

func (u *unit) executeSuccessfulMove(x, y uint8, bState *battleState) {
	u.executeMove(x, y, bState)
	u.resetMovementCounters()
	u.updateMovementHistory()
}

func (u *unit) handleMovementBlocked(bState *battleState, blockedX, blockedY uint8) {
	detourX, detourY, ok := u.findLocalDetour(bState, blockedX, blockedY)

	if !ok {
		u.executeSuccessfulMove(detourX, detourY, bState)
		u.invalidatePathForRecalculation()
		return
	}

	u.handlePersistentBlock()
}

func (u *unit) findLocalDetour(bState *battleState, blockedX, blockedY uint8) (uint8, uint8, bool) {
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

			if !u.canMoveTo(uint8(x), uint8(y), bState) {
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
	u.Command = cmdUIdle
	u.clearPath()
	u.BlockedCounter = 0
	u.AllowFriendlyFire = false

	// 25.04.2026 dodaję czyszczenie celu, bo powoduje niespójność w stanie
	// bez tego jednostka jednocześnie jest bezczynna, jak i ma cel do ataku!
	u.TargetX, u.TargetY = 0, 0
	u.TargetID = 0
	u.interactionTargetX, u.interactionTargetY = 0, 0

	if u.State != stateWaiting {
		u.IsInQueue = false
	}
}

func (u *unit) handleTargetReached(bState *battleState) {
	u.clearPath()

	switch u.Command {
	case cmdUAttack:
		log.Printf("INFO: units.go handleTargetReached cmdAttack jesteśmy u celu")

		u.State = stateAttacking
		u.attack(bState)
	case cmdUCastSpell:
		u.State = stateCastingSpell
	case cmdUBuild:
		u.State = stateBuilding
		u.build(bState)
	case cmdURepair:
		u.State = stateRepairing
		u.repair(bState)
	default:
		u.setIdle()
	}
}

func (bState *battleState) assignGroupCommand(
	command commandType, mainTargetX, mainTargetY uint8, mainTargetID uint,
	selectedUnits []*unit,
) {
	if len(selectedUnits) == 0 {
		return
	}

	targetX, targetY := bState.resolveActualTarget(mainTargetX, mainTargetY, mainTargetID)

	if len(selectedUnits) <= 4 {
		bState.assignSmallGroupTargets(selectedUnits, command, targetX, targetY, mainTargetID)
		return
	}

	bState.assignScatteredGroupTargets(selectedUnits, command, targetX, targetY, mainTargetID)
}

func (bState *battleState) resolveActualTarget(mainTargetX, mainTargetY uint8, mainTargetID uint) (uint8, uint8) {
	if mainTargetID == 0 {
		return mainTargetX, mainTargetY
	}

	targetUnit, targetBuilding := getObjectByID(mainTargetID, bState)

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

func (bState *battleState) assignSmallGroupTargets(units []*unit, cmdType commandType, targetX, targetY uint8, targetID uint) {
	for _, currentUnit := range units {
		cmd := &command{
			ActionType:          cmdType,
			TargetX:             targetX,
			TargetY:             targetY,
			InteractionTargetID: targetID,
		}
		currentUnit.addUnitCommand(cmd, bState)
	}
}

func (bState *battleState) assignScatteredGroupTargets(units []*unit, cmdType commandType, targetX, targetY uint8, targetID uint) {
	positions := bState.generateFormationPositions(targetX, targetY, uint8(len(units)))

	for i, currentUnit := range units {
		assignedX, assignedY := targetX, targetY

		if i < len(positions) {
			assignedX = positions[i].X
			assignedY = positions[i].Y
		}
		cmd := &command{
			ActionType:          cmdType,
			TargetX:             assignedX,
			TargetY:             assignedY,
			InteractionTargetID: targetID,
		}

		currentUnit.addUnitCommand(cmd, bState)
	}
}

func (bState *battleState) generateFormationPositions(centerX, centerY, count uint8) []point {
	positions := make([]point, 0, count)
	cols := uint8(math.Sqrt(float64(count))) + 1

	for i := uint8(0); i < count; i++ {
		row := i / cols
		col := i % cols

		offsetX := col - cols/2
		offsetY := row - count/(cols*2)

		x := centerX + offsetX
		y := centerY + offsetY

		if x < boardMaxX && y < boardMaxY && isWalkable(bState, x, y) {
			positions = append(positions, point{X: x, Y: y})
		} else {
			positions = append(positions, point{X: centerX, Y: centerY})
		}
	}

	return positions
}

func (u *unit) canMoveTo(x, y uint8, bState *battleState) bool {
	if x >= boardMaxX || y >= boardMaxY {
		return false
	}

	currentTile := &bState.Board.Tiles[x][y]

	// Kolizja z jednostkami (standardowo)
	if currentTile.Unit != nil && currentTile.Unit.ID != u.ID {
		return false
	}

	// Kolizja z terenem/budynkami (używamy nowej funkcji z pathfinding.go)
	// Przekazujemy 'u', aby obsłużyć wyjątek krowy wchodzącej do obory
	return isWalkableUnit(bState, x, y, u)
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
func (u *unit) executeMove(x, y uint8, bState *battleState) {
	// ZMIANA: Używamy nowej struktury Tiles
	// Usuń z poprzedniej pozycji (jeśli to ta jednostka tam jest)
	if bState.Board.Tiles[u.X][u.Y].Unit == u {
		bState.Board.Tiles[u.X][u.Y].Unit = nil
	}

	oldX, oldY := u.X, u.Y

	u.X, u.Y = x, y

	// Ustaw na nowej pozycji
	bState.Board.Tiles[u.X][u.Y].Unit = u

	u.PathIndex++
	u.updateMovementAnimation(oldX, oldY)
}

func (u *unit) setIdle() {
	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdUIdle
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
func (u *unit) attack(bState *battleState) {
	log.Printf("units.go attack weszliśmy do metody")

	target, err := u.validateAttackTarget(bState)
	if err != nil {
		u.setIdleWithReason(err.Error())

		return
	}

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

	if u.canAttackTarget(target) {
		log.Printf("units.go attack weszliśmy do canAttackTarget")

		u.performDirectAttack(target, bState)
	} else {
		log.Printf("units.go attack weszliśmy do !canAttackTarget.")

		// Jeśli cel oddalił się, gonimy go
		u.startMoveToAttack(bState)
	}
}

func (u *unit) faceTarget(target *combatTarget) {
	var tx, ty uint8

	if target.Unit != nil {
		tx, ty = target.Unit.X, target.Unit.Y
	} else if target.Building != nil {
		// Dla budynków celujemy w ich środek lub najbliższy punkt
		tx, ty, _ = target.Building.getClosestOccupiedTile(u.X, u.Y)
	} else if u.TargetID == 0 {
		tx, ty = u.interactionTargetX, u.interactionTargetY
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

func (u *unit) validateAttackTarget(bState *battleState) (*combatTarget, error) {
	target, err := u.validateTargetExists(bState)
	if err != nil {
		return nil, fmt.Errorf("cel zniknął")
	}

	isFriendly := (target.Unit != nil && target.Unit.Owner == u.Owner) ||
		(target.Building != nil && target.Building.Owner == u.Owner)

	if isFriendly && !u.AllowFriendlyFire {
		return nil, fmt.Errorf("atak na jednostkę sojuszniczą niedozwolony")
	}

	if target.Building != nil {
		if target.Building.Type == buildingPalisade && !canDamagePalisades(u) {
			return nil, fmt.Errorf("jednostka nie może niszczyć palisad")
		}

		if target.Building.HP == 0 {
			return nil, fmt.Errorf("budynek zburzony")
		}
	}

	if target.Unit != nil && target.Unit.HP == 0 {
		return nil, fmt.Errorf("cel już ubity")
	}

	return target, nil
}

func (u *unit) canAttackTarget(target *combatTarget) bool {
	distance := u.calculateDistanceToTarget(target)

	return distance <= u.AttackRange
}

func (u *unit) performDirectAttack(target *combatTarget, bState *battleState) {
	if u.AttackRange > 1 {
		u.performRangedAttack(target, u.Damage, bState)
	} else {
		u.performMeleeAttack(target, u.Damage, bState)
	}

	u.setAttackTimings()
	u.handleTargetPostAttack(target.Unit, target.Building)
}

func (u *unit) performRangedAttack(target *combatTarget, damage uint16, bState *battleState) {
	targetX, targetY, ok := u.getRangedTargetCoords(target)
	if !ok {
		log.Printf("UWAGA: jednostka %d: nie udało się określić koordynatów celu dla pocisku", u.ID)
	}

	// Mechanizm odejmowania many za rzucenie magicznego pocisku
	if u.Type.isCaster() && !u.tryToDecreaseMana(u.getProjectileManaCost()) {
		return
	}

	proj := &projectile{}
	proj.initProjectile(
		unitTypeToMissileType(u.Type),
		u.Owner,
		uint16(u.X), uint16(u.Y),
		uint16(targetX), uint16(targetY),
		damage,
	)

	proj.AllowFriendlyFire = u.AllowFriendlyFire

	bState.Projectiles = append(bState.Projectiles, proj)

	// Za stworzenie jakiegokolwiek pocisku jest przyznawane doświadczenie.
	// Muszę dodać logikę rozdziało pomięcy celem jednostką a celem budynkiem.
	// u.gainExperience tutaj!
	u.gainExperience(target.Unit, bState)

	log.Printf("jednostka %d wystrzeliła pocisk w (%d, %d) z obrażeniami %d", u.ID, targetX, targetY, damage)
}

func (u *unit) getProjectileManaCost() uint16 {
	// @reminder: inne typy jednostek nie rzucają magicznych pocisków.
	// Dlatego mam default zwracający zawsze prawdę i nie rozpisałem
	// wszystkich typów.
	switch u.Type { //nolint:exhaustive
	case unitPriestess:
		return magicThunderManaCost
	case unitPriest:
		return magicFireManaCost
	case unitMage:
		return magicGhostManaCost
	default:
		return 0
	}
}

func (u *unit) getRangedTargetCoords(target *combatTarget) (uint8, uint8, bool) {
	if target.Unit != nil {
		return target.Unit.X, target.Unit.Y, true
	} else if target.Building != nil {
		return target.Building.getClosestOccupiedTile(u.X, u.Y)
	}

	// Drzewa
	if u.TargetID == 0 {
		return u.interactionTargetX, u.interactionTargetY, true
	}

	return 0, 0, false
}

// @reminder: zdobywanie doświadczenia jest niezależne od wyniku ataku. Wykonał atak→gainExperience()
func (u *unit) performMeleeAttack(target *combatTarget, damage uint16, bState *battleState) {
	switch {
	case target.Unit != nil && target.Unit.Exists:
		target.Unit.takeDamage(damage, bState)
		u.gainExperience(target.Unit, bState)
	case target.Building != nil && target.Building.Exists:
		target.Building.takeDamage(damage)
		u.gainExperience(nil, bState)
	case bState.Board.Tiles[u.interactionTargetX][u.interactionTargetY].isDryTree():
		bState.Board.Tiles[u.interactionTargetX][u.interactionTargetY].accumulateTreeCuts(bState)
	default:
		log.Printf("UWAGA: jednostka %d: cel ataku wręcz już nie istnieje", u.ID)
	}
}

// @todo: @reminder: wróć i zastanów się, czy to rozdzielenie ma sens, bo wygląda
// jak przekombinowane, zwykły przełącznik bez podfunkcji powinien też się sprawdzić.
func (u *unit) setAttackTimings() {
	if u.AttackRange > 1 {
		u.setRangedTimings()
	} else {
		u.setMeleeTimings()
	}
}

func (u *unit) setRangedTimings() {
	u.Delay = 12

	switch u.Type {
	case unitCrossbowman:
		u.AttackCooldown = 70
	case unitMage, unitPriest, unitPriestess:
		u.AttackCooldown = 100
	default:
		u.AttackCooldown = 65
	}
}

func (u *unit) setMeleeTimings() {
	u.Delay = 10

	switch u.Type {
	case unitBear:
		u.AttackCooldown = 35
	case unitUnknown:
		u.AttackCooldown = 15
	default:
		u.AttackCooldown = 25
	}
}

func (u *unit) handleTargetPostAttack(targetUnit *unit, targetBld *building) {
	// Sprawdź czy cel przestał istnieć LUB ma 0 HP
	var targetDestroyed bool

	if targetUnit != nil && (!targetUnit.Exists || targetUnit.HP == 0) {
		targetDestroyed = true
	}

	if targetBld != nil {
		if targetBld.Type == buildingBridge {
			targetDestroyed = false
		} else if !targetBld.Exists || targetBld.HP == 0 {
			targetDestroyed = true
		}
	}

	if targetDestroyed {
		u.setIdleWithReason("cel zniszczony")
	} else {
		u.State = stateAttacking
		u.AnimationType = "fight"
		u.AnimationFrame = 0
	}
}

func (u *unit) repair(bState *battleState) {
	// 0. Tylko drwale mogą naprawiać!
	if u.Type != unitAxeman {
		return
	}
	// 1. Sprawdzamy, czy istnieje
	_, targetBuilding := getObjectByID(u.TargetID, bState)

	if targetBuilding == nil || !targetBuilding.Exists || targetBuilding.HP >= targetBuilding.MaxHP {
		u.State = stateIdle
		u.AnimationType = "idle"
		u.Command = cmdUIdle

		return
	}
	// 2. Szukamy drogi do celu
	distance := targetBuilding.getDistanceToUnit(u.X, u.Y)

	var amount uint16

	switch u.Owner {
	case bState.PlayerID:
		amount = repairAmountPlayer
	case bState.AIPlayerID:
		amount = repairAmountAI
	}

	if distance == 1 {
		targetBuilding.repair(amount)
	}
}

func (u *unit) build(bState *battleState) {
	// 0. Tylk drwale mogą budować.
	if u.Type != unitAxeman {
		return
	}

	_, targetBuilding := getObjectByID(u.TargetID, bState)

	if targetBuilding == nil || !targetBuilding.Exists || !targetBuilding.IsUnderConstruction {
		u.State = stateIdle
		u.AnimationType = "idle"
		u.Command = cmdUIdle

		return
	}

	// 2. Szukamy drogi do celu
	distance := targetBuilding.getDistanceToUnit(u.X, u.Y)

	var amount uint16

	switch u.Owner {
	case bState.PlayerID:
		amount = repairAmountPlayer
	case bState.AIPlayerID:
		amount = repairAmountAI
	}

	if distance == 1 {
		targetBuilding.build(amount, bState)
	}
}

//goland:noinspection SpellCheckingInspection
func (u *unit) castMagicShield() {
	// 0. Jeśli już jest magiczna tarcza, to nie można rzucić nowej
	if u.hasMagicShield {
		return
	}
	// 1. Odejmujemy potrzebną manę
	// @todo: sprawdź ile rzeczywiście kosztowała
	if u.tryToDecreaseMana(spellCostMagicShield) {
		// 2. aktywujemy efekt
		u.hasMagicShield = true
		// 3. Ustalamy ile ma trwać
		// @todo: sprawdź ile rzeczywiście trwała
		u.MagicShieldCooldown = spellDurationMagicShield
	}

	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdUIdle
}

// Metoda odpowiedzialna za gromobicie i deszcz ognia.
func (u *unit) magicShower(targetX, targetY uint8, bState *battleState) bool {
	// 0. Koszt czaru
	if u.Mana < spellBufferMagicShower || !u.tryToDecreaseMana(spellCostMagicShower) {
		log.Printf("INFO: Jednostka %d nie ma wystarczająco many na rzucenie czaru", u.ID)
		return false
	}

	// 1. Tworzymy czarodziejski deszcz
	u.createMagicShower(targetX, targetY, bState)

	// 2. Skończyliśmy czarowanie, stoimy bezczynnie
	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdUIdle

	log.Printf("INFO: Jednostka %d rzuciła czar na (%d, %d)", u.ID, targetX, targetY)

	return true
}

func (u *unit) createMagicShower(targetX, targetY uint8, bState *battleState) {
	damage, missileKind, ok := u.resolveMagicShowerStats()
	if !ok {
		log.Printf("UWAGA: magicShower wywołany dla jednostki o nieobsługiwanym rodzaju %d!", u.Type)
		return
	}

	u.spawnMagicShowerProjectiles(targetX, targetY, missileKind, damage, bState)
}

func (u *unit) resolveMagicShowerStats() (damage uint16, missileKind uint8, ok bool) {
	// @reminder: inne rodzaje jednostek nie mają tego czaru więc ich nie wymieniam tutaj.
	switch u.Type {
	case unitPriest:
		return spellDamageFireShower, missileFireRain, true
	case unitPriestess:
		return spellDamageLightningShower, missileLightning, true
	default:
		return 0, 0, false
	}
}

func (u *unit) spawnMagicShowerProjectiles(targetX, targetY uint8, missileKind uint8, damage uint16, bState *battleState) {
	// 0. Bezpiecznik pozycji początkowej tworzonych pocisków
	spawnY := targetY

	if spawnY >= 4 {
		spawnY -= 4
	} else {
		spawnY = 0
	}

	// 1. Tworzenie opadów
	for offset := -1; offset <= 1; offset++ {
		spawnX := int(targetX) + offset

		if spawnX < 0 || spawnX >= int(boardMaxX) {
			continue
		}

		proj := &projectile{}
		proj.initProjectile(
			missileKind,
			u.Owner,
			uint16(spawnX), uint16(spawnY),
			uint16(spawnX), uint16(targetY),
			damage,
		)

		if proj.Exists {
			bState.Projectiles = append(bState.Projectiles, proj)
		}

		// 2. Przyzanie doświadczenia za zaatakowanie
		currentTile := &bState.Board.Tiles[spawnX][targetY]

		switch {
		case currentTile.Unit != nil && currentTile.Unit.Exists:
			u.gainExperience(currentTile.Unit, bState)
		case currentTile.Building != nil && currentTile.Building.Exists:
			u.gainExperience(nil, bState)
		default:
			// Nie przyznajemy nic doświadczenia za napaść na otoczenie
		}
	}
}

func (u *unit) canCastSpellFromCurrentPosition() bool {
	targetX, targetY := u.interactionTargetX, u.interactionTargetY

	if targetX == 0 && targetY == 0 {
		targetX, targetY = u.TargetX, u.TargetY
	}

	dist := uint8(math.Max(
		math.Abs(float64(int(u.X)-int(targetX))),
		math.Abs(float64(int(u.Y)-int(targetY))),
	))

	return dist <= u.AttackRange
}

func (u *unit) castSpell(bState *battleState) {
	if u.AttackCooldown > 0 {
		u.State = stateIdle
		u.AnimationType = "idle"
		u.Delay = 1

		return
	}

	targetX, targetY := u.interactionTargetX, u.interactionTargetY

	if targetX == 0 && targetY == 0 {
		targetX, targetY = u.TargetX, u.TargetY
	}

	if u.magicShower(targetX, targetY, bState) {
		u.setRangedTimings()
		u.setIdleWithReason("czar rzucony")
	} else {
		u.State = stateIdle
		u.AnimationType = "idle"
		u.Command = cmdUIdle
	}
}

func (u *unit) castMagicSight(bState *battleState) {
	if u.Mana >= spellCostMagicSight {
		u.Mana -= spellCostMagicSight
		log.Printf("Jednostka %d rzuca czar widzenia", u.ID)

		revealRadius := spellCostRangeMagicSight
		for i := u.X - revealRadius; i <= u.X+revealRadius; i++ {
			for j := u.Y - revealRadius; j <= u.Y+revealRadius; j++ {
				if i <= boardMaxX && j <= boardMaxY {
					// @todo: czemu 18?!
					if math.Abs(float64(u.X-i))+math.Abs(float64(u.Y-j)) < 18 {
						bState.Board.Tiles[i][j].Visibility = visibilityVisible
					}
				}
			}
		}
	}

	log.Printf("unit %d: Rzucono Czar Widoczności. Cała mapa odkryta.", u.ID)

	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdUIdle
}

func (u *unit) takeDamage(damage uint16, bState *battleState) {
	// 0. Sprawdzamy, czy jednostka jest chroniona przed obrażeniami.
	if u.hasMagicShield {
		return
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
		if u.Udder < 100 && u.Command != cmdUFlee {
			barnX, barnY, foundBarn := findNearestBarnMilkingSpot(u, bState)
			if foundBarn {
				cmd := &command{
					ActionType:          cmdUFlee,
					TargetX:             barnX,
					TargetY:             barnY,
					InteractionTargetID: 0,
				}
				u.addUnitCommand(cmd, bState)
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
		decreasePopulation(u, bState)

		occupiedTile := &bState.Board.Tiles[u.X][u.Y]
		if occupiedTile.Unit == u {
			occupiedTile.Unit = nil
		}

		createCorpses(u, bState)
		u.unregisterFromBuilding()
		log.Printf("Jednostka %d została zabita!", u.ID)
	}
}

func decreasePopulation(u *unit, bState *battleState) {
	switch u.Owner {
	case bState.HumanPlayerState.PlayerID:
		bState.HumanPlayerState.CurrentPopulation--
	case bState.AIEnemyState.PlayerID:
		bState.AIEnemyState.CurrentPopulation--
	}
}

func createCorpses(u *unit, bState *battleState) {
	steps := 18
	stepAngle := 10
	rotation := float32(rand.Intn(steps) * stepAngle)
	currentCorpse := corpse{
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
	switch currentCorpse.UnitType {
	case unitCow:
		currentCorpse.SkeletonType = 1
	case unitBear:
		currentCorpse.SkeletonType = 2
	default:
		currentCorpse.SkeletonType = 0
	}

	bState.CorpsesList = append(bState.CorpsesList, currentCorpse)
}

func (u *unit) unregisterFromBuilding() {
	if u.BelongsTo != nil {
		log.Printf("Jednostka %d próbuje się wyrejestrować z budynku %d", u.ID, u.BelongsTo.ID)
		u.BelongsTo.unregisterUnit(u.ID)
		u.BelongsTo = nil
	}
}

func (u *unit) findNearestPalisade(bState *battleState, radius uint8,
) *building {
	var best *building
	minD := math.MaxFloat64

	for _, pal := range bState.Buildings {
		if pal == nil || !pal.Exists || pal.Type != buildingPalisade {
			continue
		}

		cx, cy, _ := pal.getCenter()
		dx := math.Abs(float64(u.X - cx))
		dy := math.Abs(float64(u.Y - cy))
		d := math.Max(dx, dy)

		if d <= float64(radius) && d < minD {
			if u.isImportantPalisade(pal, bState) {
				minD = d
				best = pal
			}
		}
	}

	return best
}

func (u *unit) actOnIdle(bState *battleState) {
	if !u.canActOnIdle() {
		return
	}

	if u.BlockedCounter > 0 {
		return
	}

	if !u.shouldSearchForTarget(bState) {
		return
	}

	u.handleTargetSearch(bState)
}

func (u *unit) canActOnIdle() bool {
	return u.Type != unitCow && u.Type != unitShepherd
}

func (u *unit) shouldSearchForTarget(bState *battleState) bool {
	return u.isReadyToAct(bState)
}

func (u *unit) isReadyToAct(bState *battleState) bool {
	if u.State == stateIdle && u.Command == cmdUIdle {
		return true
	}

	if u.Command == cmdUAttack {
		_, err := u.validateTargetExists(bState)
		if err != nil {
			return true
		}
	}

	return false
}

const palisadeStrategicBuildingProximity = 10

func (u *unit) isImportantPalisade(palisade *building, bState *battleState) bool {
	if u.Owner != bState.AIPlayerID || !canDamagePalisades(u) {
		return false
	}

	if palisade == nil || !palisade.Exists || palisade.Type != buildingPalisade {
		return false
	}

	const prox = palisadeStrategicBuildingProximity

	palCenterX, palCenterY, _ := palisade.getCenter()

	for _, bld := range bState.Buildings {
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

		_, _, ok = bld.getClosestWalkableTile(bState)

		if !ok {
			return true
		}
	}

	return false
}

func (u *unit) handleTargetSearch(bState *battleState) {
	if u.Owner == bState.HumanPlayerState.PlayerID {
		u.handleTargetSearchForHumanPlayer(bState)
	} else {
		u.handleTargetSearchForAI(bState)
	}
}

func (u *unit) handleTargetSearchForHumanPlayer(bState *battleState) {
	primaryTargetUnit, primaryTargetBuilding, foundPrimary := findNearestEnemyExtended(u, bState)

	if !foundPrimary {
		u.setIdle()

		return
	}

	if primaryTargetUnit != nil && primaryTargetUnit.Exists {
		u.handleUnitTarget(primaryTargetUnit, bState)

		return
	}

	if primaryTargetBuilding != nil && primaryTargetBuilding.Exists {
		u.setIdle()

		return
	}

	u.setIdle()
}

func (u *unit) handleTargetSearchForAI(bState *battleState) {
	isPalisadeBreaker := canDamagePalisades(u)

	primaryTargetUnit, primaryTargetBuilding, foundPrimary := findNearestEnemyExtended(u, bState)

	if isPalisadeBreaker && foundPrimary && primaryTargetBuilding != nil && primaryTargetBuilding.Exists {
		_, _, ok := primaryTargetBuilding.getClosestWalkableTile(bState)

		if !ok {
			palisadeTarget := u.findNearestPalisade(bState, u.SightRange)

			if palisadeTarget != nil {
				u.handleBuildingTarget(palisadeTarget, bState)

				return
			}
		}
	}

	if !foundPrimary {
		u.setIdle()
		return
	}

	if primaryTargetUnit != nil {
		u.handleUnitTarget(primaryTargetUnit, bState)
	} else {
		u.handleBuildingTarget(primaryTargetBuilding, bState)
	}
}

func (u *unit) handleUnitTarget(targetUnit *unit, bState *battleState) {
	u.TargetID = targetUnit.ID
	u.setMoveTargetForUnit(targetUnit, bState)

	u.executeActionBasedOnDistance(bState)
}

func (u *unit) handleBuildingTarget(targetBuilding *building, bState *battleState) {
	u.TargetID = targetBuilding.ID

	currentDistanceToBuilding := targetBuilding.getDistanceToUnit(u.X, u.Y)

	if currentDistanceToBuilding != math.MaxUint8 && currentDistanceToBuilding <= u.AttackRange {
		u.startDirectAttack(u.X, u.Y, bState)
		return
	}

	finalMoveTargetX, finalMoveTargetY := uint8(0), uint8(0)

	optimalRangedX, optimalRangedY := uint8(0), uint8(0)
	closestWalkableX, closestWalkableY := uint8(0), uint8(0)
	var ok bool

	if u.AttackRange > 1 {
		optimalRangedX, optimalRangedY, ok = findOptimalRangedAttackTile(u.X, u.Y, u.AttackRange, targetBuilding, bState)
	}

	if !ok {
		finalMoveTargetX, finalMoveTargetY = optimalRangedX, optimalRangedY
	} else {
		closestWalkableX, closestWalkableY, ok = targetBuilding.getClosestWalkableTile(bState)
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
	u.startMoveToAttack(bState)
	log.Printf("DEBUG_AI: unit %d (Type:%d) moving to attack bld %d, target tile: (%d,%d). Current position: (%d,%d).",
		u.ID, u.Type, u.TargetID, u.TargetX, u.TargetY, u.X, u.Y)
}

func (u *unit) setMoveTargetForUnit(targetUnit *unit, bState *battleState) {
	bestX, bestY := u.findBestPositionAroundUnit(targetUnit, bState)
	u.TargetX, u.TargetY = bestX, bestY
}

func (u *unit) findBestPositionAroundUnit(targetUnit *unit, bState *battleState) (uint8, uint8) {
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

			if u.isValidMoveTarget(uint8(checkX), uint8(checkY), bState) {
				// log.Println("Funkcja findBestPositionAroundUnit isValidMoveTarget = true, szukam freeSpot")
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
		// log.Println("Funkcja findBestPositionAroundUnit !foundFreeSpot")

		return targetUnit.X, targetUnit.Y // Fallback
	}

	// log.Println("Funkcja findBestPositionAroundUnit foundFreeSpot, x: %d, y: %d", bestX, bestY)

	return uint8(bestX), uint8(bestY)
}

func (u *unit) findBestPositionAroundTile(tileX, tileY uint8, bState *battleState) (uint8, uint8, bool) {
	// 1. Najpierw sprawdzamy, czy sam kliknięty kafelek jest przechodni i pusty
	if isWalkable(bState, tileX, tileY) {
		currentTile := &bState.Board.Tiles[tileX][tileY]
		if currentTile.Unit == nil || currentTile.Unit.ID == u.ID {
			return tileX, tileY, true
		}
	}

	// 2. Jeśli nie, szukamy w promieniu AttackRange
	bestX, bestY := tileX, tileY
	minDist := math.MaxFloat64
	found := false

	for dx := -int(u.AttackRange); dx <= int(u.AttackRange); dx++ {
		for dy := -int(u.AttackRange); dy <= int(u.AttackRange); dy++ {
			checkX := int(tileX) + dx
			checkY := int(tileY) + dy

			if checkX < 0 || checkX >= int(boardMaxX) || checkY < 0 || checkY >= int(boardMaxY) {
				continue
			}

			if !isWalkable(bState, uint8(checkX), uint8(checkY)) {
				continue
			}

			// Pomijamy kafelki zajęte przez inne jednostki
			currentTile := &bState.Board.Tiles[checkX][checkY]
			if currentTile.Unit != nil && currentTile.Unit.ID != u.ID {
				continue
			}

			dist := math.Abs(float64(dx)) + math.Abs(float64(dy))
			if dist < minDist {
				minDist = dist
				bestX, bestY = uint8(checkX), uint8(checkY)
				found = true
			}
		}
	}

	return bestX, bestY, found
}

func (u *unit) isValidMoveTarget(x, y uint8, bState *battleState) bool {
	return x < boardMaxX && y < boardMaxY &&
		bState.Board.Tiles[x][y].Unit == nil &&
		bState.Board.Tiles[x][y].Building == nil &&
		isWalkable(bState, x, y)
}

func (u *unit) validateTargetExists(bState *battleState) (*combatTarget, error) {
	targetUnit, targetBuilding := getObjectByID(u.TargetID, bState)

	if targetBuilding != nil && targetBuilding.Type == buildingBridge {
		return &combatTarget{Unit: targetUnit, Building: targetBuilding}, nil
	}

	if u.TargetID == 0 {
		treeTile := bState.Board.Tiles[u.interactionTargetX][u.interactionTargetY]
		if treeTile.isStandingTree() && !treeTile.IsBurning {
			return &combatTarget{Unit: nil, Building: nil}, nil
		}

		return nil, fmt.Errorf("cel (drzewo) nie istnieje")
	}

	if (targetUnit == nil || !targetUnit.Exists) && (targetBuilding == nil || !targetBuilding.Exists) {
		return nil, fmt.Errorf("cel nie istnieje")
	}

	return &combatTarget{Unit: targetUnit, Building: targetBuilding}, nil
}

func (u *unit) calculateDistanceToTarget(target *combatTarget) uint8 {
	if target.Unit != nil {
		return uint8(math.Max(
			math.Abs(float64(int(u.X)-int(target.Unit.X))),
			math.Abs(float64(int(u.Y)-int(target.Unit.Y))),
		))
	}

	if target.Building != nil {
		return target.Building.getDistanceToUnit(u.X, u.Y)
	}

	// Atak na drzewo
	return uint8(math.Max(
		math.Abs(float64(int(u.X)-int(u.interactionTargetX))),
		math.Abs(float64(int(u.Y)-int(u.interactionTargetY))),
	))
}

func (u *unit) executeActionByDistance(distance uint8, bState *battleState) {
	if distance > u.SightRange {
		log.Printf("DEBUG_AI: U %d: cel ID %d poza zasięgiem widzenia. IDLE", u.ID, u.TargetID)
		u.setIdle()

		return
	}

	if distance <= u.AttackRange {
		log.Printf("DEBUT_AI: U %d: Odległość <= zasięg ataku. Rozpoczynam bezpośredni atak na cel ID %d.",
			u.ID, u.TargetID)
		u.startDirectAttack(u.TargetX, u.TargetY, bState)
	} else {
		log.Printf("DEBUG_AI: U %d: odległość %d > zasięg ataku %d. Ruszam w kierunku %d.",
			u.ID, distance, u.AttackRange, u.TargetID)
		u.startMoveToAttack(bState)
	}
}

func (u *unit) executeActionBasedOnDistance(bState *battleState) {
	target, err := u.validateTargetExists(bState)
	if err != nil {
		u.setIdle()

		return
	}

	distance := u.calculateDistanceToTarget(target)
	u.executeActionByDistance(distance, bState)
}

func (u *unit) startDirectAttack(placeholderX, placeholderY uint8, bState *battleState) {
	realTargetX := placeholderX
	realTargetY := placeholderY

	if u.TargetID != 0 {
		targetUnit, targetBld := getObjectByID(u.TargetID, bState)

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

	cmd := &command{
		ActionType:          cmdUAttack,
		TargetX:             realTargetX,
		TargetY:             realTargetY,
		InteractionTargetID: u.TargetID,
	}

	u.addUnitCommand(cmd, bState)

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

	for _, occupiedTile := range bld.OccupiedTiles {
		dx := float64(occupiedTile.X - fromX)
		dy := float64(occupiedTile.Y - fromY)
		distSq := dx*dx + dy*dy

		if distSq < minDistSq {
			minDistSq = distSq
			closestX = occupiedTile.X
			closestY = occupiedTile.Y
		}
	}

	return closestX, closestY, true
}

func (u *unit) startMoveToAttack(bState *battleState) {
	cmd := &command{
		ActionType:          cmdUAttack,
		TargetX:             u.TargetX,
		TargetY:             u.TargetY,
		InteractionTargetID: u.TargetID,
	}
	u.addUnitCommand(cmd, bState)
	u.State = stateMoving
	u.AnimationType = "walk"

	log.Printf("DEBUG_AI: unit %d moving to attack target %d, move to (%d,%d)",
		u.ID, u.TargetID, u.TargetX, u.TargetY)
}

func getLegacyUnitIndex(t unitType) int {
	return int(t)
}

func findOptimalRangedAttackTile(uCurrentX, uCurrentY, attackRange uint8, bld *building, bState *battleState) (uint8, uint8, bool) {
	if len(bld.OccupiedTiles) == 0 {
		return 0, 0, false
	}

	centerX, centerY, ok := bld.getCenter()
	if !ok {
		return 0, 0, false
	}

	candidates := []point{
		{X: centerX + attackRange, Y: centerY},
		{X: centerX - attackRange, Y: centerY},
		{X: centerX, Y: centerY + attackRange},
		{X: centerX, Y: centerY - attackRange},
		{X: centerX + attackRange, Y: centerY + attackRange},
		{X: centerX - attackRange, Y: centerY - attackRange},
		{X: centerX + attackRange, Y: centerY - attackRange},
		{X: centerX - attackRange, Y: centerY + attackRange},
	}

	halfRange := attackRange / 2
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

		if !bld.isValidWalkableTile(x, y, bState) {
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
