package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
)

func (bState *battleState) getPlayerState(ownerID PlayerID) *playerState {
	if ownerID == bState.HumanPlayerState.PlayerID {
		return bState.HumanPlayerState
	}

	return bState.AIEnemyState
}

func (bState *battleState) createBuilding(bldType buildingType, topLeftX, topLeftY uint8, owner PlayerID) *building {
	stats, ok := buildingDefs[bldType]
	if !ok {
		log.Printf("BŁĄD: Nieznany rodzaj budynku %d", bldType)

		return nil
	}

	newBld := &building{
		ID:            BuildingID(bState.NextUniqueObjectID),
		Type:          bldType,
		Owner:         owner,
		Exists:        true,
		Armor:         buildingArmor,
		MaxHP:         stats.MaxHP,
		MaxFood:       stats.MaxFood,
		AssignedUnits: make([]UnitID, 0, stats.MaxFood),
	}

	bState.NextUniqueObjectID++

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
		return fmt.Errorf("neznany rodzaj budynku")
	}

	ownerState := bState.getPlayerState(owner)

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

// Zakładamy, że każda palisada tworzona z mapy jest prawidłowo umiejscowiona
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
		bState.HumanPlayerState.CurrentPopulation--
	case bState.AIEnemyState.PlayerID:
		bState.AIEnemyState.CurrentPopulation--
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
