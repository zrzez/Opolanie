package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (bState *battleState) createBuilding(bldType buildingType, topLeftX, topLeftY uint8, owner PlayerID) *building {
	stats, ok := buildingDefs[bldType]
	if !ok {
		log.Printf("BŁĄD: Nieznany rodzaj budynku %d", bldType)

		return nil
	}

	newBld := &building{
		ID:            bState.NextBuildingID,
		Type:          bldType,
		Owner:         owner,
		Exists:        true,
		Armor:         buildingArmor,
		MaxHP:         stats.MaxHP,
		MaxFood:       stats.MaxFood,
		AssignedUnits: make([]UnitID, 0, stats.MaxFood),
	}

	bState.NextBuildingID++

	bState.Board.registerBuilding(newBld, topLeftX, topLeftY)
	bState.Buildings = append(bState.Buildings, newBld)

	switch owner {
	case bState.HumanPlayerState.PlayerID:
		bState.HumanPlayerState.CurrentBuildings++
	case bState.AIEnemyState.PlayerID:
		bState.AIEnemyState.CurrentBuildings++
	}

	return newBld
}

func (bState *battleState) placeConstructionSite(bldType buildingType, tileX, tileY uint8, owner PlayerID) *building {
	// 0. Walidacja, być może zbędna
	stats, ok := buildingDefs[bldType]
	if !ok {
		log.Printf("nieznany rodzaj budynku %d", bldType)

		return nil
	}

	// 1. Zasadzamy budowę
	newBld := bState.createBuilding(bldType, tileX, tileY, owner)
	if newBld == nil {
		return nil
	}

	// 2. Ustawiamy jako plac budowy
	newBld.IsUnderConstruction = true
	newBld.HP = initialConstructionHP

	// 3. Wygląd
	bState.Board.applyConstructionGraphics(newBld)

	// 4. Informacja dla gracza
	bState.CurrentMessage.Text = fmt.Sprintf("Wznoszenie: %s", stats.Name)
	bState.CurrentMessage.Duration = 60

	log.Printf("BUDOWA: Rozpoczęto %s (ID: %d) na (%d,%d).", stats.Name, newBld.ID, tileX, tileY)

	return newBld
}

func (bState *battleState) tryBuildStructure(bldType buildingType, tileX, tileY uint8, owner PlayerID) error {
	stats, ok := buildingDefs[bldType]
	if !ok {
		return fmt.Errorf("nieznany rodzaj budynku")
	}

	ownerState := bState.getPlayerState(owner)

	if ownerState == nil {
		return fmt.Errorf("nieznany właściciel budynku")
	}

	if ownerState.Milk < stats.Cost {
		return fmt.Errorf("za mało mleka")
	}

	ownerState.Milk -= stats.Cost

	if bldType == buildingRoad {
		bState.Board.placeRoad(tileX, tileY)

		return nil
	}

	// Budynki mogą być postawione na nieukończonych palisadach.
	// Taka mechanika wymaga dodatkowego wyczyszczenia planszy.
	if bldType.isRegularBuilding() {
		bState.Board.clearConstructionSite(tileX, tileY, stats)
	}

	bState.placeConstructionSite(bldType, tileX, tileY, owner)

	return nil
}

// Zakładamy, że każda palisada tworzona z mapy jest prawidłowo umiejscowiona.
func (bState *battleState) createFinishedPalisade(tileX, tileY uint8) *building {
	// 1. Tworzymy obiekt
	newBld := bState.createBuilding(buildingPalisade, tileX, tileY, colorNone)
	//                                  Palisada zawsze jest niczyja ↑↑↑↑↑
	if newBld == nil {
		return nil
	}

	// 2. Jest to gotowy budynek „zaciągnięty” z mapy
	newBld.IsUnderConstruction = false // niby wartość myślna, ale dla pewności
	newBld.HP = newBld.MaxHP

	// 3. Aktualizujemy tekstury z planszy
	joinPalisade(tileX, tileY, bState.Board)

	return newBld
}

// Zakładamy, że każdy budynek jest prawidłowo umiejscowiony
func (bState *battleState) createFinishedBuilding(bldType buildingType, tileX, tileY uint8, owner PlayerID) *building {
	// 1. Tworzymy obiekt
	newBld := bState.createBuilding(bldType, tileX, tileY, owner)
	if newBld == nil {
		return nil
	}

	// 2. Jest to gotowy budynek „zaciągnięty” z mapy
	newBld.IsUnderConstruction = false
	newBld.HP = newBld.MaxHP

	// 3. Aktualizujemy tekstury w planszy
	bState.Board.applyFinishedGraphics(newBld)

	return newBld
}

// szuka budynku w battleState.Buildings.
func (bState *battleState) getBuildingByID(bldID BuildingID) (*building, bool) {
	for _, bld := range bState.Buildings {
		if bld.ID == bldID {
			return bld, true
		}
	}

	return nil, false
}

// szuka jednostki w battleState.Units.
func (bState *battleState) getUnitByID(uID UnitID) (*unit, bool) {
	for _, currentUnit := range bState.Units {
		if currentUnit.ID == uID {
			return currentUnit, true
		}
	}

	return nil, false
}

func (bState *battleState) getObjectByID(oID ObjectID) (*unit, *building) {
	if currentUnit, ok := bState.getUnitByID(UnitID(oID)); ok {
		return currentUnit, nil
	}

	if currentBuilding, ok := bState.getBuildingByID(BuildingID(oID)); ok {
		return nil, currentBuilding
	}

	return nil, nil
}

// Zwraca listę zaznaczonych jednostek, które należą do gracza.
func (bState *battleState) getSelectedUnits() []*unit {
	var selected []*unit

	for _, currentUnit := range bState.Units {
		if currentUnit.IsSelected && currentUnit.Exists && currentUnit.Owner == bState.PlayerID {
			selected = append(selected, currentUnit)
		}
	}

	return selected
}

func (bState *battleState) decreasePopulation(owner PlayerID) {
	switch owner {
	case bState.HumanPlayerState.PlayerID:
		if bState.HumanPlayerState.CurrentPopulation > 0 {
			bState.HumanPlayerState.CurrentPopulation--
		}
	case bState.AIEnemyState.PlayerID:
		if bState.AIEnemyState.CurrentPopulation > 0 {
			bState.AIEnemyState.CurrentPopulation--
		}
	}
}

func (bState *battleState) createCorpses(u *unit) {
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

func (bState *battleState) assignGroupCommand(
	command commandType, mainTargetX, mainTargetY uint8, mainTargetID ObjectID,
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

func (bState *battleState) resolveActualTarget(mainTargetX, mainTargetY uint8, mainTargetID ObjectID) (uint8, uint8) {
	if mainTargetID == 0 {
		return mainTargetX, mainTargetY
	}

	targetUnit, targetBuilding := bState.getObjectByID(mainTargetID)

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

func (bState *battleState) assignSmallGroupTargets(units []*unit, cmdType commandType, targetX, targetY uint8, targetID ObjectID) {
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

func (bState *battleState) assignScatteredGroupTargets(units []*unit, cmdType commandType, targetX, targetY uint8, targetID ObjectID) {
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

		if x < boardMaxX && y < boardMaxY && isWalkable(bState.Board, x, y) {
			positions = append(positions, point{X: x, Y: y})
		} else {
			positions = append(positions, point{X: centerX, Y: centerY})
		}
	}

	return positions
}

func (bState *battleState) canProduceUnit(unitType unitType, bld *building) (bool, uint8) {
	// 1. Czy jest wolne miejsce w budynku?
	if !bld.hasRoom() {
		return false, produceErrNoRoom
	}

	// Ustalamy kto chce wykonać działanie
	ownerState := bState.getPlayerState(bld.Owner)

	if ownerState == nil {
		return false, produceErrInvalidOwner
	}

	// 2. Czy nie przekraczamy odgórnego ograniczenia?
	if ownerState.CurrentPopulation > maxUnitsPerPlayer {
		return false, produceErrPopulationLimit
	}

	// Pobieramy dane o jednostce do stworzenia
	uStats, ok := unitDefs[unitType]

	if !ok {
		return false, produceErrInvalidType
	}

	// 3. Czy stać gracza?
	if ownerState.Milk < uStats.Cost {
		return false, produceErrMilk
	}

	// 4. Czy jednostka może wyjść z budynku?
	if !bState.Board.hasSpaceAroundBuilding(bld) {
		return false, produceErrNoSpace
	}

	// Udało się, można tworzyć
	return true, produceErrNone
}

func (bState *battleState) cleanupDestroyedBuildings() {
	if bState.GlobalFrameCounter%6000 != 0 {
		return
	}

	log.Println("INFO: Rozpoczynam czyszczenie pamięci z budynków...")

	newBuildingsList := make([]*building, 0, len(bState.Buildings))
	for _, bld := range bState.Buildings {
		if bld.Exists {
			newBuildingsList = append(newBuildingsList, bld)
		}
	}

	removedCount := len(bState.Buildings) - len(newBuildingsList)
	bState.Buildings = newBuildingsList

	log.Printf("INFO: Wyczyszczono %d zniszczonych budynków.", removedCount)
}

// tryProduceUnit odpowiada za próbę wytworzenia jednostki.
func (bState *battleState) tryProduceUnit(newUnitType unitType, bld *building) {
	// JUŻ PO WALIDACJI!
	// Bierzemy dane jednostki, nie sprawdzamy, bo już to zrobiliśmy
	// @reminder: jeśli się przebuduje sygnaturę canProduceUnit to handleProductionComman
	// mogłoby dać już staty i oszczędzić, ale też zaciemnić obraz
	uStats := unitDefs[newUnitType]

	// 3. Ustalamy właściciela
	// mamy bld, a on owner więc bez sensu!
	ownerState := bState.getPlayerState(bld.Owner)
	if ownerState == nil {
		// Wypadałoby coś tutaj zrobić, bo to krytyczny błąd
		fmt.Print("ownerState == nil")

		return
	}

	// 4. Tworzymy jednostkę
	coords, _ := bState.Board.electSpawnTile(bld)
	fmt.Println("wchodzę w createUnit")
	bState.createUnit(newUnitType, coords, bld)

	// 5. Pobieramy mleko za jednostkę
	ownerState.Milk -= uStats.Cost
}

// To nie powinna być metoda budynku tylko bState lub board ponieważ w tej chwili
// budynek zna szczegóły tworzenia jednostek, a nie powinien.
func (bState *battleState) createUnit(unitType unitType, coords point, bld *building) {
	fmt.Print("jestem w createUnit\n")
	fmt.Printf("Chcę stworzyć jednostkę od ID %d\n", bState.NextUnitID)
	newUnit := bState.initUnit(unitType, coords.X, coords.Y, bState.NextUnitID)
	bState.NextUnitID++

	newUnit.Owner = bld.Owner
	newUnit.BelongsTo = bld

	bState.Board.Tiles[newUnit.X][newUnit.Y].Unit = newUnit
	bState.Units = append(bState.Units, newUnit)

	fmt.Print("Rejestruję jednostkę w budynku\n")
	bld.registerUnit(newUnit.ID)

	ownerState := bState.getPlayerState(bld.Owner)

	if ownerState != nil {
		ownerState.CurrentPopulation++
	}
}

// @reminder: to nie powinno przyjmować x,y uint8 a strukturę point
func (bState *battleState) initUnit(unitType unitType, x, y uint8, newUnitID UnitID) *unit {
	newUnit := &unit{}

	newUnit.ID = newUnitID
	newUnit.Exists = true
	newUnit.Type = unitType
	newUnit.X = x
	newUnit.Y = y
	newUnit.AnimationType = "walk"
	newUnit.Direction = rl.NewVector2(0, 1)
	newUnit.Wounds = make([]wound, 0, maxWoundsCount)

	if newUnit.Type != unitCow {
		newUnit.Command = cmdUIdle
	} else {
		newUnit.Command = cmdUGraze
	}

	stats, ok := unitDefs[unitType]
	if ok {
		newUnit.SightRange = stats.SightRange
		newUnit.AttackRange = stats.AttackRange
		newUnit.Damage = stats.BaseDamage
		newUnit.Armor = stats.BaseArmor
		newUnit.MaxHP = stats.MaxHP
		newUnit.MaxDelay = stats.MoveDelay
		newUnit.HP = stats.MaxHP
		newUnit.MaxMana = stats.MaxMana
		newUnit.Mana = stats.MaxMana / 2
		newUnit.ManaRegen = stats.BaseManaRegen
	} else {
		log.Printf("OSTRZEŻENIE: Nieznany rodzaj jednostki w init: %d.\n", unitType)
	}

	newUnit.Delay = newUnit.MaxDelay

	fmt.Printf("initUnit: rodzaj %d, UnitID: %d\n", newUnit.Type, newUnit.ID)

	return newUnit
}
