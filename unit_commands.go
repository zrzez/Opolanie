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

func (u *unit) addUnitCommand(cmd *command, board *boardData, resolver objectResolver) {
	log.Printf("INFO: unit.go dodano rozkaz %d.", cmd.ActionType)
	// ŁATANIE DZIURY W KOMPLETOWANIU ROZKAÓW DLA JEDNOSTEK
	// @reminder: Łatanie dziury w kompletowaniu rozkazów dla jednostek
	// @todo: ogarnij to łatanie, bo nie powinno to tutaj być! - 02.07.2026
	u.CurrentSpell = cmd.Spell
	u.AllowFriendlyFire = cmd.FriendlyFire

	var approach *point

	if cmd.ActionType.isInteraction() {
		var err error
		// Jeśli czar wymaga interakcji, to obliczamy gdzie podejść
		// Na drzewo nie da się wejść, więc trzeba znaleźć kafelek obok
		targetCoords := &point{X: cmd.TargetX, Y: cmd.TargetY}
		approach, err = u.calculateApproachTile(targetCoords, cmd.InteractionTargetID, board, resolver)
		if err != nil {
			u.setIdleWithReason("cel nieosiągalny")

			return
		}
	} else {
		// Nie wymaga interakcji, np. cmdMove, to cel jest miejscem w które się udajemy
		approach = &point{X: cmd.TargetX, Y: cmd.TargetY}
	}

	// ! tutaj się zastanawiam, co zrobić
	if cmd.ActionType == cmdUAttack {
		if !u.canAttack(cmd.InteractionTargetID, cmd.TargetX, cmd.TargetY, board, resolver) {
			fmt.Print("rozkaz nie przeszedł sprawdzenia canAttack w addUnitCommand")

			return
		}
	}

	// Przekazujemy cel oraz podejście
	u.prepareForNewCommand(cmd.ActionType, cmd.TargetX, cmd.TargetY, cmd.InteractionTargetID, approach.X, approach.Y)
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
	u.TargetX, u.TargetY = 0, 0
	u.TargetID = 0
	u.interactionTargetX, u.interactionTargetY = 0, 0

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

func (u *unit) executeCommandAction(resolver objectResolver, board *boardData, pathfindingBudget *int, bState *battleState) {
	switch u.Type {
	case unitCow:
		u.handleCowBehavior(resolver, board, pathfindingBudget, bState)
	default:
		u.executeStandardUnitCommand(resolver, board, pathfindingBudget, bState)
	}
}

func (u *unit) executeStandardUnitCommand(resolver objectResolver, board *boardData, pathfindingBudget *int, bState *battleState) {
	switch u.Command {
	case cmdUMove:
		u.move(resolver, board, pathfindingBudget, bState)
	case cmdUAttack:
		if u.canAttackTargetFromCurrentPosition(resolver, board) {
			u.State = stateAttacking
			u.clearPath()
			u.attack(resolver, board, bState)
		} else {
			u.State = stateMoving
			u.move(resolver, board, pathfindingBudget, bState)
		}
	case cmdUBuild:
		if u.canAttackTargetFromCurrentPosition(resolver, board) {
			u.State = stateBuilding

			if u.AnimationType != "fight" {
				u.AnimationType = "fight"
				u.AnimationFrame = 0
			}

			u.clearPath()
			_, targetBuilding := bState.getObjectByID(u.TargetID)

			var amount uint16

			switch u.Owner {
			case bState.PlayerID:
				amount = repairAmountPlayer
			case bState.AIPlayerID:
				amount = repairAmountAI
			}

			u.build(targetBuilding, amount)
		} else {
			u.State = stateMoving

			if u.AnimationType != "walk" {
				u.AnimationType = "walk"
			}

			u.move(resolver, board, pathfindingBudget, bState)
		}
	case cmdURepair:
		if u.canAttackTargetFromCurrentPosition(resolver, board) {
			u.State = stateRepairing

			if u.AnimationType != "fight" {
				u.AnimationType = "fight"
				u.AnimationFrame = 0
			}

			u.clearPath()
			_, targetBuilding := bState.getObjectByID(u.TargetID)

			var amount uint16

			switch u.Owner {
			case bState.PlayerID:
				amount = repairAmountPlayer
			case bState.AIPlayerID:
				amount = repairAmountAI
			}

			u.repair(targetBuilding, amount)
		} else {
			u.State = stateMoving

			if u.AnimationType != "walk" {
				u.AnimationType = "walk"
			}

			u.move(resolver, board, pathfindingBudget, bState)
		}
	case cmdUCastSpell:
		u.castSpell(resolver, board, pathfindingBudget, bState)

	case cmdUIdle, cmdUStop:
		u.actOnIdle(resolver, board, bState)
	default:
		panic("unhandled default case")
	}
}

func (u *unit) canAttackTargetFromCurrentPosition(resolver objectResolver, board *boardData) bool {
	log.Println("Sprawdzam, czy cel istnieje")

	target, err := u.validateTargetExists(resolver, board)
	if err != nil {
		return false
	}

	log.Println("Obliczam odległość")

	distance := u.calculateDistanceToTarget(target)

	log.Printf("Odległość obliczona %d", distance)
	log.Println(distance <= u.AttackRange)

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

func (u *unit) calculateApproachTile(intention *point, targetID ObjectID, board *boardData, resolver objectResolver) (*point, error) {
	if u.CurrentSpell != spellNone {
		approachTile, err := u.findApproachTileForSpell(intention, board)
		if err != nil {
			return nil, err
		}

		return approachTile, nil
	}

	// Budynki, jednostki i drzewa jako cel
	return u.findApproachTileForTarget(intention, targetID, board, resolver)
}

func (u *unit) findApproachTileForSpell(targetPosition *point, board *boardData) (*point, error) {
	switch u.CurrentSpell {
	case spellMagicShower:

		validCoords, ok := findTileForAttacking(u, nil, nil, targetPosition, board)
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
func (u *unit) findApproachTileForTarget(intention *point, targetID ObjectID, board *boardData, resolver objectResolver) (*point, error) {
	targetUnit, targetBuilding := resolver.getObjectByID(targetID)

	var targetU *unit

	var targetBld *building

	var targetTree *point

	switch {
	case targetBuilding != nil && (targetBuilding.Exists || targetBuilding.Type == buildingBridge):
		targetBld = targetBuilding
	case targetUnit != nil && targetUnit.Exists:
		targetU = targetUnit
	case intention != nil && board.Tiles[intention.X][intention.Y].isTree():
		targetTree = intention
	}

	validCoords, ok := findTileForAttacking(u, targetU, targetBld, targetTree, board)
	if !ok {
		return nil, fmt.Errorf("nie ma podejścia do celu: %t", ok)
	}

	return findBestReachableTile(u, validCoords, board)
}

// @todo: nie powinna to być metoda jednostki, bo to sprawdzanie poprawności
// ! albo przekazuję objectResolver albo combattarget? coś mi tutaj nie pasuje.
func (u *unit) canAttack(targetID ObjectID, intentionX, intentionY uint8, board *boardData, resolver objectResolver) bool {
	// Drzewa
	// dany kafelek musi istnieć więc nie robię != nil
	// ! sprawdzenie raczej bez sensu, bo nie tworzymy sobie współrzędnych bez sensu
	if targetID == 0 {
		if u.TargetX >= boardMaxX || u.TargetY >= boardMaxY {
			fmt.Print("DUPA BOARDMAX")

			return false
		}

		targetTile := &board.Tiles[intentionX][intentionY]

		if targetTile.isTree() {
			fmt.Print("DUPA TO NIE DRZEWO")

			return u.canDamageTree(targetTile)
		}

		fmt.Print("DUPA DRZEWO, ALE WYPAD")
		// Jeśli ID == 0, ale nie jest drzewem, to mamy problem…
		return false
	}

	targetUnit, targetBuilding := resolver.getObjectByID(targetID)

	// Jednostki
	// musi istnieć
	if targetUnit != nil && targetUnit.ID != u.ID {
		// musi być oznaczona jako żywa
		return targetUnit.Exists
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

		if targetBuilding.Type == buildingPalisade && !u.Type.canDamagePalisades() {
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

func (u *unit) validateAttackTarget(resolver objectResolver, board *boardData) (*combatTarget, error) {
	target, err := u.validateTargetExists(resolver, board)
	if err != nil {
		return nil, fmt.Errorf("cel zniknął")
	}

	isFriendly := (target.Unit != nil && target.Unit.Owner == u.Owner) ||
		(target.Building != nil && target.Building.Owner == u.Owner)

	if isFriendly && !u.AllowFriendlyFire {
		return nil, fmt.Errorf("atak na jednostkę sojuszniczą niedozwolony")
	}

	if target.Building != nil {
		if target.Building.Type == buildingPalisade && !u.Type.canDamagePalisades() {
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

func (u *unit) handleTargetSearch(resolver objectResolver, board *boardData, bState *battleState) {
	if u.Owner == bState.HumanPlayerState.PlayerID {
		u.handleTargetSearchForHumanPlayer(resolver, board, bState)
	} else {
		u.handleTargetSearchForAI(resolver, board, bState)
	}
}

// @reminder: szukanie celu dla gracza i SI różnią się szczegółami, np. jednostka gracza nie napadają
//    samoistnie na wrogie budynki.
func (u *unit) handleTargetSearchForHumanPlayer(resolver objectResolver, board *boardData, bState *battleState) {
	primaryTargetUnit, _, foundPrimary := findNearestEnemyExtended(u, bState)

	if !foundPrimary {
		u.setIdle()

		return
	}

	if primaryTargetUnit != nil && primaryTargetUnit.Exists {
		u.handleUnitTarget(primaryTargetUnit, resolver, board)

		return
	}

	u.setIdle()
}

// @reminder: szukanie celu dla gracza i SI różnią się szczegółami, np. jednostka gracza nie napadają
//    samoistnie na wrogie budynki.
func (u *unit) handleTargetSearchForAI(resolver objectResolver, board *boardData, bState *battleState) {
	isPalisadeBreaker := u.Type.canDamagePalisades()

	primaryTargetUnit, primaryTargetBuilding, foundPrimary := findNearestEnemyExtended(u, bState)

	if isPalisadeBreaker && foundPrimary && primaryTargetBuilding != nil && primaryTargetBuilding.Exists {
		ok := bState.Board.hasSpaceAroundBuilding(primaryTargetBuilding)

		if !ok {
			palisadeTarget := u.findNearestPalisade(bState, u.SightRange)

			if palisadeTarget != nil {
				u.handleBuildingTarget(palisadeTarget, board, resolver)

				return
			}
		}
	}

	if !foundPrimary {
		u.setIdle()

		return
	}

	if primaryTargetUnit != nil {
		u.handleUnitTarget(primaryTargetUnit, resolver, board)
	}

	if primaryTargetBuilding != nil {
		u.handleBuildingTarget(primaryTargetBuilding, board, resolver)
	}
}

func (u *unit) validateTargetExists(resolver objectResolver, board *boardData) (*combatTarget, error) {
	targetUnit, targetBuilding := resolver.getObjectByID(u.TargetID)

	// To chyba jest chodzenie po moście albo jego budowa.
	// Szkoda, że nie zostawiłem komentarzy.
	if targetBuilding != nil && targetBuilding.Type == buildingBridge {
		return &combatTarget{Building: targetBuilding}, nil
	}

	// Do atakowania drzew. Muszę poprawić!
	if u.TargetID == 0 {
		treeTile := &board.Tiles[u.TargetX][u.TargetY]
		if treeTile.isStandingTree() && !treeTile.IsBurning {
			return &combatTarget{Tile: treeTile}, nil
		}

		return nil, fmt.Errorf("cel (drzewo) nie istnieje")
	}

	// Wyłapujemy, czy budynek lub jednostka przestały istnieć
	if (targetUnit == nil || !targetUnit.Exists) && (targetBuilding == nil || !targetBuilding.Exists) {
		return nil, fmt.Errorf("cel nie istnieje")
	}

	// Jest spoko, można przepuścić cel
	return &combatTarget{Unit: targetUnit, Building: targetBuilding}, nil
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
		log.Printf("INFO: unit %d zaczyna atakować cel ID %d z miejsca.", u.ID, u.TargetID)
	} else {
		// Cel poza zasięgiem, przechodzimy w stan ruchu do wyliczonego ApproachX/Y
		u.State = stateMoving
		u.AnimationType = "walk"
		u.invalidatePathForRecalculation()
	}
}

func (u *unit) executeActionBasedOnDistance(resolver objectResolver, board *boardData) {
	target, err := u.validateTargetExists(resolver, board)
	if err != nil {
		u.setIdle()

		return
	}

	distance := u.calculateDistanceToTarget(target)
	u.executeActionByDistance(distance)
}

// @reminder: dla bezczynnych jednostek. Nie powinna się sama zadaniować.

func (u *unit) actOnIdle(resolver objectResolver, board *boardData, bState *battleState) {
	if !u.canActOnIdle() {
		return
	}

	if u.BlockedCounter > 0 {
		return
	}

	if !u.shouldSearchForTarget(resolver, board) {
		return
	}

	u.handleTargetSearch(resolver, board, bState)
}

func (u *unit) canActOnIdle() bool {
	return u.Type != unitCow && u.Type != unitShepherd
}

func (u *unit) shouldSearchForTarget(resolver objectResolver, board *boardData) bool {
	return u.isReadyToAct(resolver, board)
}

func (u *unit) isReadyToAct(resolver objectResolver, board *boardData) bool {
	if u.State == stateIdle && u.Command == cmdUIdle {
		return true
	}

	if u.Command == cmdUAttack {
		_, err := u.validateTargetExists(resolver, board)
		if err != nil {
			return true
		}
	}

	return false
}

func (u *unit) handleUnitTarget(targetUnit *unit, resolver objectResolver, board *boardData) {
	u.TargetID = ObjectID(targetUnit.ID)

	coords, err := u.findApproachTileForTarget(nil, u.TargetID, board, resolver)
	if err != nil {
		u.setIdleWithReason("nie można znaleźć drogi do wrogiej jednostki")

		return
	}

	u.ApproachX = coords.X
	u.ApproachY = coords.Y
	u.TargetX = coords.X
	u.TargetY = coords.Y

	u.executeActionBasedOnDistance(resolver, board)
}

func (u *unit) handleBuildingTarget(targetBuilding *building, board *boardData, resolver objectResolver) {
	u.TargetID = ObjectID(targetBuilding.ID)

	coords, err := u.findApproachTileForTarget(nil, u.TargetID, board, resolver)
	if err != nil {
		u.setIdleWithReason("nie można znaleźć drogi do wrogiego budynku")

		return
	}

	u.ApproachX = coords.X
	u.ApproachY = coords.Y
	u.TargetX = coords.X
	u.TargetY = coords.Y

	u.executeActionBasedOnDistance(resolver, board)
}

func (u *unit) handleTargetReached(resolver objectResolver, board *boardData, bState *battleState) {
	u.clearPath()

	switch u.Command {
	case cmdUAttack:
		log.Printf("INFO: units.go handleTargetReached cmdAttack jesteśmy u celu")

		u.State = stateAttacking
		u.attack(resolver, board, bState)
	case cmdUCastSpell:
		u.State = stateCastingSpell
	case cmdUBuild:
		u.State = stateBuilding
		_, targetBuilding := bState.getObjectByID(u.TargetID)

		var amount uint16

		switch u.Owner {
		case bState.PlayerID:
			amount = repairAmountPlayer
		case bState.AIPlayerID:
			amount = repairAmountAI
		}

		u.build(targetBuilding, amount)
	case cmdURepair:
		u.State = stateRepairing
		_, targetBuilding := bState.getObjectByID(u.TargetID)

		var amount uint16

		switch u.Owner {
		case bState.PlayerID:
			amount = repairAmountPlayer
		case bState.AIPlayerID:
			amount = repairAmountAI
		}

		u.repair(targetBuilding, amount)
	default:
		u.setIdle()
	}
}
