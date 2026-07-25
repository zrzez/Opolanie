package main

import (
	"fmt"
	"log"
	"math"
)

const palisadeStrategicBuildingProximity = 10

// ============
// Próba rozplątania units.go, tutaj powinny trafiać funkcje związane z
// przetwarzaniem rozkazów przez jednoski.

func (u *unit) addUnitCommand(cmd *command, bState *battleState) {
	log.Printf("INFO: unit.go dodano rozkaz %d.", cmd.ActionType)
	// ŁATANIE DZIURY W KOMPLETOWANIU ROZKAÓW DLA JEDNOSTEK
	// @reminder: Łatanie dziury w kompletowaniu rozkazów dla jednostek
	// @todo: ogarnij to łatanie, bo nie powinno to tutaj być! - 02.07.2026
	u.CurrentSpell = cmd.Spell
	u.AllowFriendlyFire = cmd.FriendlyFire

	var approach *point

	if cmd.ActionType.isInteraction() {
		var err error

		approach, err = u.calculateApproachTile(cmd.Target, bState)
		if err != nil {
			u.setIdleWithReason("cel nieosiągalny")

			return
		}
	} else {
		// Nie wymaga interakcji, np. cmdMove, to cel jest miejscem w które się udajemy
		approach = &cmd.Target.Position
	}

	// ! tutaj się zastanawiam, co zrobić
	if cmd.ActionType == cmdUAttack {
		target, err := bState.resolveTarget(cmd.Target)
		if err != nil {
			log.Printf("INFO: rozkaz ataku odrzucony: %v", err)

			return
		}

		if !u.canAttack(target) {
			log.Printf("INFO: rozkaz ataku odrzucony: jdnostka %d nie może zaatakować", u.ID)

			return
		}
	}

	// Przekazujemy cel oraz podejście
	u.prepareForNewCommand(cmd.ActionType, cmd.Target, *approach)
	u.applyCommandState(cmd.ActionType)
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
	u.Target = TargetReference{}
	u.Approach = point{}

	if u.State != stateWaiting {
		u.IsInQueue = false
	}
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
		u.AnimationType = "fight"
		u.AnimationFrame = 3
		u.AnimationCounter = 0
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

func (u *unit) executeCommandAction(pathfindingBudget *int, bState *battleState) {
	switch u.Type {
	case unitCow:
		u.handleCowBehavior(pathfindingBudget, bState)
	default:
		u.executeStandardUnitCommand(pathfindingBudget, bState)
	}
}

func (u *unit) executeStandardUnitCommand(pathfindingBudget *int, bState *battleState) {
	switch u.Command {
	case cmdUMove:
		u.move(pathfindingBudget, bState)
	case cmdUAttack:
		if u.canAttackTargetFromCurrentPosition(bState) {
			u.State = stateAttacking
			u.clearPath()
			u.attack(bState)
		} else {
			u.State = stateMoving
			u.move(pathfindingBudget, bState)
		}
	case cmdUBuild, cmdURepair:
		u.handleWorkCommand(pathfindingBudget, bState)

	case cmdUCastSpell:
		u.castSpell(pathfindingBudget, bState)

	case cmdUIdle, cmdUStop:
		u.actOnIdle(bState)
	default:
		panic("unhandled default case")
	}
}

func (u *unit) handleWorkCommand(pathfindingBudget *int, bState *battleState) {
	// 1. Zasięg
	if !u.canAttackTargetFromCurrentPosition(bState) {
		u.State = stateMoving

		if u.AnimationType != "walk" {
			u.AnimationType = "walk"
		}

		u.move(pathfindingBudget, bState)

		return
	}

	// 2. Pobieramy cel
	target, err := bState.resolveTarget(u.Target)
	if err != nil || target.Building == nil {
		u.setIdleWithReason("cel prac zniknał lub nie jest budynkiem")

		return
	}

	// 3. Ustawiamy odpowiedni stan
	if u.Command == cmdUBuild {
		u.State = stateBuilding
	} else {
		u.State = stateRepairing
	}

	if u.AnimationType != "fight" {
		u.AnimationType = "fight"
		u.AnimationFrame = 0
	}

	u.clearPath()

	// 4. Ustalamy ile da jedna jednostka pracy
	var amount uint16

	switch u.Owner {
	case bState.HumanPlayerState.PlayerID:
		amount = repairAmountPlayer
	case bState.AIEnemyState.PlayerID:
		amount = repairAmountAI
	}

	// 5. Wykon
	if u.Command == cmdUBuild {
		u.build(target.Building, amount)
	} else {
		u.repair(target.Building, amount)
	}
}

func (u *unit) canAttackTargetFromCurrentPosition(bState *battleState) bool {
	target, err := bState.resolveTarget(u.Target)
	if err != nil {
		return false
	}

	distance := u.calculateDistanceToTarget(target)

	return distance <= u.AttackRange
}

// @reminder: o ile kojarzęto tutaj korzystam ze złego sposobu obliczania odległości.
//    Powinna być odległość Czebyszewa.
func (u *unit) calculateDistanceToTarget(target *combatTarget) uint8 {
	if target.Unit != nil {
		return uint8(math.Max(
			math.Abs(float64(int(u.X)-int(target.Unit.X))),
			math.Abs(float64(int(u.Y)-int(target.Unit.Y))),
		))
	}

	if target.Building != nil {
		return getDistanceToUnit(target.Building.Type, target.Building.OccupiedTiles[0], u.X, u.Y)
	}

	// Atak na drzewo
	return uint8(math.Max(
		math.Abs(float64(int(u.X)-int(target.Tile.X))),
		math.Abs(float64(int(u.Y)-int(target.Tile.Y))),
	))
}

func (u *unit) calculateApproachTile(targetRef TargetReference, bState *battleState) (*point, error) {
	if u.CurrentSpell != spellNone {
		approachTile, err := u.findApproachTileForSpell(targetRef.Position, bState.Board)
		if err != nil {
			return nil, err
		}

		return approachTile, nil
	}

	// Budynki, jednostki i drzewa jako cel
	return u.findApproachTileForTarget(targetRef, bState)
}

func (u *unit) findApproachTileForSpell(targetPosition point, board *boardData) (*point, error) {
	switch u.CurrentSpell {
	case spellMagicShower:

		validCoords, ok := findTileForAttacking(u, nil, nil, &targetPosition, board)
		if !ok {
			return nil, fmt.Errorf("nie ma podejścia do celu")
		}

		return findBestReachableTile(u, validCoords, board)

	// ↓↓↓↓↓ Poniższe przypadki nie muszą korzystać z A*
	case spellMagicShield, spellMagicSight:
		// Czary, które przyjmują rzucającego jako swój cel.
		return &point{X: u.X, Y: u.Y}, nil
	case spellNone:
		// To nigdy nie powinno mieć miejsca, bo warunkiem wywołania
		// jest u.CurrentSpell != spellNone.
		return &point{X: u.X, Y: u.Y}, fmt.Errorf("próba rzucenia spellNone")
	default:
		// To nigdy nie powinno mieć miejsca, bo wszystkie czary są obsłużone
		return &point{X: u.X, Y: u.Y}, fmt.Errorf("nieznany rodzaj czaru")
	}
}

// @reminder: nazwa dla kafelka z drzewem „intention” jest bardzo kiepska, ale nie mam teraz do tego głowy.
func (u *unit) findApproachTileForTarget(targetRef TargetReference, bState *battleState) (*point, error) {
	target, err := bState.resolveTarget(targetRef)
	if err != nil {
		return nil, err
	}

	var targetU *unit

	var targetBld *building

	var targetTree *point

	switch {
	case target.Building != nil && (target.Building.Exists || target.Building.Type == buildingBridge):
		targetBld = target.Building
	case target.Unit != nil && target.Unit.Exists:
		targetU = target.Unit
	case target.Tile != nil && target.Tile.isTree():
		targetTree = &targetRef.Position
	}

	validCoords, ok := findTileForAttacking(u, targetU, targetBld, targetTree, bState.Board)
	if !ok {
		return nil, fmt.Errorf("nie ma podejścia do celu: %t", ok)
	}

	return findBestReachableTile(u, validCoords, bState.Board)
}

// @todo: nie powinna to być metoda jednostki, bo to sprawdzanie poprawności
// ! albo przekazuję objectResolver albo combattarget? coś mi tutaj nie pasuje.
func (u *unit) canAttack(target *combatTarget) bool {
	switch {
	case target.Tile != nil:
		// Jeśli pole nie zawiera drzewa
		if !target.Tile.isTree() {
			return false
		}

		// Pole z drzewem, sprawdzamy, czy możemy je niszczyć
		return u.canDamageTree(target.Tile)

	case target.Unit != nil:
		// Nie możemy okaleczyć samego siebie
		return target.Unit.Exists && target.Unit.ID != u.ID
	case target.Building != nil:
		// Tylko istniejące budynki
		if !target.Building.Exists {
			return false
		}

		// Magowie nie mogą atakować budynków
		if u.Type == unitMage {
			return false
		}

		// Tylko niektóre jednostki mogą niszczyć palisady
		if target.Building.Type == buildingPalisade && !u.Type.canDamagePalisades() {
			return false
		}

		// Wybudowany most jest nietykalny
		if target.Building.Type == buildingBridge {
			return false
		}

		return true
	}

	// Wydaje się, że cel nie istnieje
	return false
}

func (u *unit) validateAttackTarget(target *combatTarget) error {
	isFriendly := (target.Unit != nil && target.Unit.Owner == u.Owner) ||
		(target.Building != nil && target.Building.Owner == u.Owner)

	if isFriendly && !u.AllowFriendlyFire {
		return fmt.Errorf("atak na jednostkę sojuszniczą niedozwolony")
	}

	if target.Building != nil {
		if target.Building.Type == buildingPalisade && !u.Type.canDamagePalisades() {
			return fmt.Errorf("jednostka nie może niszczyć palisad")
		}

		if target.Building.HP == 0 {
			return fmt.Errorf("budynek zburzony")
		}
	}

	if target.Unit != nil && target.Unit.HP == 0 {
		return fmt.Errorf("cel już ubity")
	}

	if target.Tile != nil && !target.Tile.isStandingTree() {
		return fmt.Errorf("drzewo zostało już ścięte")
	}

	return nil
}

func (u *unit) canAttackTarget(target *combatTarget) bool {
	distance := u.calculateDistanceToTarget(target)

	return distance <= u.AttackRange
}

// !
func (u *unit) getRangedTargetCoords(target *combatTarget) (*point, bool) {
	switch {
	case target.Unit != nil:
		return &point{X: target.Unit.X, Y: target.Unit.Y}, true
	case target.Building != nil:
		return getClosestOccupiedTile(&point{X: u.X, Y: u.Y}, &target.Building.OccupiedTiles)
	case target.Tile != nil: // drzewa dla unitPriest
		return &point{X: target.Tile.X, Y: target.Tile.Y}, true

	}

	return nil, false
}

func (u *unit) handleTargetPostAttack(target *combatTarget) {
	// Sprawdź czy cel przestał istnieć LUB ma 0 HP
	var targetDestroyed bool

	if target.Unit != nil && (!target.Unit.Exists || target.Unit.HP == 0) {
		targetDestroyed = true
	}

	if target.Building != nil {
		if target.Building.Type == buildingBridge {
			targetDestroyed = false
		} else if !target.Building.Exists || target.Building.HP == 0 {
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

func (u *unit) isImportantPalisade(palisade *building, bState *battleState) bool {
	if u.Owner != bState.AIPlayerID || !u.Type.canDamagePalisades() {
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

		//_, _, ok = bld.getClosestWalkableTile(bState)
		ok = bState.Board.hasSpaceAroundBuilding(bld)

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

// @reminder: szukanie celu dla gracza i SI różnią się szczegółami, np. jednostka gracza nie napadają
//    samoistnie na wrogie budynki.
func (u *unit) handleTargetSearchForHumanPlayer(bState *battleState) {
	primaryTargetUnit, _, foundPrimary := findNearestEnemyExtended(u, bState)

	if !foundPrimary {
		u.setIdle()

		return
	}

	if primaryTargetUnit != nil && primaryTargetUnit.Exists {
		u.handleUnitTarget(primaryTargetUnit, bState)

		return
	}

	u.setIdle()
}

// @reminder: szukanie celu dla gracza i SI różnią się szczegółami, np. jednostka gracza nie napadają
//    samoistnie na wrogie budynki.
func (u *unit) handleTargetSearchForAI(bState *battleState) {
	isPalisadeBreaker := u.Type.canDamagePalisades()

	primaryTargetUnit, primaryTargetBuilding, foundPrimary := findNearestEnemyExtended(u, bState)

	if isPalisadeBreaker && foundPrimary && primaryTargetBuilding != nil && primaryTargetBuilding.Exists {
		ok := bState.Board.hasSpaceAroundBuilding(primaryTargetBuilding)

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
	}

	if primaryTargetBuilding != nil {
		u.handleBuildingTarget(primaryTargetBuilding, bState)
	}
}

// @todo: nie pamiętam po co to, chyba tylko do odsiania obiektów, których jednostka
//   „nie widzi”. Tutaj trzeba będzie dać poprawkę na to, czy gracz odsłonił kafelek.
func (u *unit) executeActionByDistance(distance uint8) {
	if distance > u.SightRange {
		u.setIdleWithReason("cel poza zasięgiem widzenia")

		return
	}

	u.Command = cmdUAttack

	if distance <= u.AttackRange {
		// Cel w zasięgu, od razu przechodzimy do ataku
		u.State = stateAttacking
		u.AnimationType = "fight"
		u.AnimationFrame = 0
		u.AnimationCounter = 0
		log.Printf("INFO: unit %d zaczyna atakować cel ID %d z miejsca.", u.ID, u.Target.ID)
	} else {
		// Cel poza zasięgiem, przechodzimy w stan ruchu do wyliczonego ApproachX/Y
		u.State = stateMoving
		u.invalidatePathForRecalculation()
	}
}

func (u *unit) executeActionBasedOnDistance(bState *battleState) {
	target, err := bState.resolveTarget(u.Target)
	if err != nil {
		u.setIdleWithReason("cel przestał istnieć lub jest nieosiągalny")

		return
	}

	distance := u.calculateDistanceToTarget(target)
	u.executeActionByDistance(distance)
}

// @reminder: dla bezczynnych jednostek. Nie powinna się sama zadaniować.
func (u *unit) actOnIdle(bState *battleState) {
	if !u.canActOnIdle() {
		return
	}

	if u.BlockedCounter > 0 {
		return
	}

	if !u.isReadyToAct(bState) {
		return
	}

	u.handleTargetSearch(bState)
}

func (u *unit) canActOnIdle() bool {
	return u.Type != unitCow && u.Type != unitShepherd
}

func (u *unit) isReadyToAct(bState *battleState) bool {
	// 1. Siedzi i nie robi nic
	if u.State == stateIdle && u.Command == cmdUIdle {
		return true
	}

	// 2. Próbuje atakować
	if u.Command == cmdUAttack {
		// Więc sprawdzamy, czy cel jeszcze istnieje
		_, err := bState.resolveTarget(u.Target)
		if err != nil {
			// Nie istnieje więc można zacząć robić co innego
			return true
		}
	}

	// Jednostka robi coś innego, jest zajęta
	return false
}

func (u *unit) handleUnitTarget(targetedUnit *unit, bState *battleState) {
	u.Target = TargetReference{
		Kind:     targetUnit,
		ID:       uint(targetedUnit.ID),
		Position: targetedUnit.Position,
	}

	coords, err := u.findApproachTileForTarget(u.Target, bState)
	if err != nil {
		u.setIdleWithReason("nie można znaleźć drogi do wrogiej jednostki")

		return
	}

	u.Approach = *coords
	u.Target.Position = *coords

	u.executeActionBasedOnDistance(bState)
}

func (u *unit) handleBuildingTarget(targetedBuilding *building, bState *battleState) {
	// @reminder: tymczasowe, później trzeba to posprzątać, bo findApproachTileForTarget powiela
	bldCenterX, bldCenterY, ok := targetedBuilding.getCenter()
	if !ok {
		u.setIdleWithReason("cel (budynek) jest nieprawidłowy")

		return
	}

	u.Target = TargetReference{
		Kind:     targetBuilding,
		ID:       uint(targetedBuilding.ID),
		Position: point{X: bldCenterX, Y: bldCenterY},
	}

	coords, err := u.findApproachTileForTarget(u.Target, bState)
	if err != nil {
		u.setIdleWithReason("nie można znaleźć drogi do wrogiego budynku")

		return
	}

	// @reminder: wygląda to bardzo podejrzanie, że „podejście” do budynku jest tym samym, co cel ataku!
	u.Approach = *coords
	u.Target.Position = *coords

	u.executeActionBasedOnDistance(bState)
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
	case cmdUBuild, cmdURepair:
		target, err := bState.resolveTarget(u.Target)
		if err != nil || target.Building == nil {
			u.setIdleWithReason("cel budywy/naprawy przepadł")

			return
		}

		var amount uint16

		switch u.Owner {
		case bState.PlayerID:
			amount = repairAmountPlayer
		case bState.AIPlayerID:
			amount = repairAmountAI
		}

		if u.Command == cmdUBuild {
			u.State = stateBuilding
			u.build(target.Building, amount)
		} else {
			u.State = stateRepairing
			u.repair(target.Building, amount)
		}

	default:
		u.setIdle()
	}
}
