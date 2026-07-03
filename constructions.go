package main

import (
	"fmt"
	"log"
	"math"
	"slices"
)

// constructions.go

// initConstruction odpowiada tylko za techniczne utworzenie obiektu w pamięci i na planszy.
// Nie mylić z tryBuildStructure.
func (bld *building) initConstruction(buildingType buildingType, owner uint8, nextUniqueObjectID uint) {
	// 1. Sprawdzamy, czy dany rodzaj budynku był określony wcześniej
	stats, ok := buildingDefs[buildingType]
	if !ok {
		log.Printf("BŁĄD KRYTYCZNY: Brak określenia dla budynku rodzaju %d!", bld.Type)

		return
	}

	// 2. Podstawowe właściwości
	bld.ID = nextUniqueObjectID
	bld.Type = buildingType
	bld.Owner = owner
	bld.Exists = true
	bld.Armor = buildingArmor
	bld.MaxHP = stats.MaxHP
	bld.HP = stats.MaxHP
	bld.MaxFood = stats.MaxFood
	bld.AssignedUnits = make([]uint, 0)
}

func (bld *building) startConstruction(bState *battleState) {
	// Pobieramy informacje o budynku
	stats, ok := buildingDefs[bld.Type]

	if !ok {
		log.Printf("BŁĄD KRYTYCZNY: Nie ma %s w buildingDefs, nie można rozpocząć budowy", stats.Name)
	}

	bld.IsUnderConstruction = true

	// Nie zaczynamy z zerem HP
	bld.HP = initialConstructionHP

	for index, tilePoint := range bld.OccupiedTiles {
		// @reminder: zakomentowałem, bo przy obecnych sprawdzeniach PRZED postawieniem budowy
		// raczej nie jest możliwe wyjście poza planszę.
		// if tilePoint.X >= boardMaxX || tilePoint.Y >= boardMaxY {
		// 	continue
		// }
		currentTile := &bState.Board.Tiles[tilePoint.X][tilePoint.Y]

		row := index / int(normalBuildingSize)
		column := index % int(normalBuildingSize)

		switch {
		case bld.Type != buildingPalisade && bld.Type != buildingBridge:
			currentTile.TextureID = constructionTemplatePhase01[column][row]
			// bs.Board.Tiles[tilePoint.X][tilePoint.Y].IsWalkable = false
		case bld.Type == buildingBridge:
			currentTile.TextureID = spriteBridgeConstruction
			// bs.Board.Tiles[tilePoint.X][tilePoint.Y].IsWalkable = false
		default:
			currentTile.TextureID = spritePalisadeDestroyed
			// bs.Board.Tiles[tilePoint.X][tilePoint.Y].IsWalkable = true
		}

		currentTile.Building = bld
		currentTile.IsWalkable = currentTile.TextureID == spritePalisadeDestroyed
	}
}

func (bld *building) repair(amount uint16) {
	if !bld.Exists || bld.HP >= bld.MaxHP {
		return
	}

	bld.increaseHPBuilding(amount)
}

func (bld *building) build(amount uint16) {
	if !bld.Exists || !bld.IsUnderConstruction {
		return
	}

	bld.increaseHPBuilding(amount)
}

func (bldType buildingType) isRegularBuilding() bool {
	return bldType != buildingRoad && bldType != buildingBridge && bldType != buildingPalisade
}

func clearConstructionSite(tileX, tileY uint8, bldStats buildingStats, bState *battleState) {
	for dx := range bldStats.Width {
		for dy := range bldStats.Height {
			cx, cy := tileX+dx, tileY+dy

			// Zabezpieczenie przed wyjściem poza mapę
			if cx >= boardMaxX || cy >= boardMaxY {
				continue
			}

			currentTile := &bState.Board.Tiles[cx][cy]

			// Czyścimy plac budowy z palisad
			if currentTile.Building != nil && currentTile.Building.Type == buildingPalisade &&
				currentTile.Building.IsUnderConstruction {

				currentTile.Building.OccupiedTiles = make([]point, 0, 1)
				currentTile.IsWalkable = false
				currentTile.Building.HP = 0
				currentTile.Building.Exists = false
			}
		}
	}
}

func placeRoad(tileX, tileY uint8, bState *battleState) {
	bState.Board.Tiles[tileX][tileY].TextureID = spriteRoadStart
	cx := int(tileX)
	cy := int(tileY)

	refreshRoadTile(cx, cy, bState.Board)

	refreshRoadTile(cx+1, cy, bState.Board) // prawo
	refreshRoadTile(cx-1, cy, bState.Board) // lewo
	refreshRoadTile(cx, cy+1, bState.Board) // góra
	refreshRoadTile(cx, cy-1, bState.Board) // dół

	log.Printf("BUDOWA: Postawiono drogę na (%d,%d).", tileX, tileY)
}

func placeConstructionSite(tileX, tileY, bldOwner uint8, bldStats buildingStats, bldType buildingType, bState *battleState) {
	newBld := &building{}

	// init teraz przyjmie tileX, tileY jako lewy górny róg
	newBld.initConstruction(bldType, bldOwner, bState.NextUniqueObjectID)
	bState.NextUniqueObjectID++

	placeConstructionOnBoard(newBld, tileX, tileY, bState.Board)

	bState.Buildings = append(bState.Buildings, newBld)
	newBld.startConstruction(bState)

	bState.CurrentMessage.Text = fmt.Sprintf("Wznoszenie: %s", bldStats.Name)
	bState.CurrentMessage.Duration = 60

	log.Printf("BUDOWA: Rozpoczęto %s (ID: %d) na (%d,%d).", bldStats.Name, newBld.ID, tileX, tileY)
}

// Stawia plac budowy we wskazanym miejscu
func tryBuildStructure(bldType buildingType, tileX, tileY uint8, owner uint8, bState *battleState) {
	bldStats := buildingDefs[bldType]
	ownerState := bState.getPlayerState(owner)

	ownerState.Milk -= bldStats.Cost

	bldOwner := colorNone

	if bldType == buildingRoad {
		placeRoad(tileX, tileY, bState)

		return
	}

	if bldType.isRegularBuilding() {
		ownerState.CurrentBuildings++
		bldOwner = colorRed

		clearConstructionSite(tileX, tileY, bldStats, bState)
	}

	placeConstructionSite(tileX, tileY, bldOwner, bldStats, bldType, bState)
}

func isObstacle(texID uint16) bool {
	switch {
	case isRockNonWalkable(texID):
		return true
	case isCompletedBridge(texID):
		return true
	case isGadget(texID):
		return true
	case isTreeStump(texID):
		return true
	case texID >= spriteTreeBurntStump00 && texID <= spriteTreeBurntStump01:
		return true
	case isSpecialTile(texID):
		return true
	case !isLandOrOther(texID):
		return true
	}

	return false
}

func (bld *building) registerUnit(unitID uint) bool {
	if uint8(len(bld.AssignedUnits)) >= bld.MaxFood {
		return false
	}
	if slices.Contains(bld.AssignedUnits, unitID) {
		return false
	}

	bld.AssignedUnits = append(bld.AssignedUnits, unitID)
	bld.Food = uint8(len(bld.AssignedUnits))
	return true
}

func (bld *building) unregisterUnit(unitID uint) bool {
	for i, id := range bld.AssignedUnits {
		if id == unitID {
			bld.AssignedUnits = append(bld.AssignedUnits[:i], bld.AssignedUnits[i+1:]...)
			bld.Food = uint8(len(bld.AssignedUnits))

			return true
		}
	}

	return false
}

func (bld *building) hasSpace() bool {
	return uint8(len(bld.AssignedUnits)) < bld.MaxFood
}

func (bld *building) getAssignedUnits(bState *battleState) []*unit {
	var units []*unit

	for _, unitID := range bld.AssignedUnits {
		currentUnit, ok := bState.getUnitByID(unitID)
		if ok && currentUnit.Exists {
			units = append(units, currentUnit)
		}
	}

	return units
}

func (bld *building) cleanupDeadUnits(bState *battleState) {
	var validUnits []uint

	for _, unitID := range bld.AssignedUnits {
		currentUnit, ok := bState.getUnitByID(unitID)
		if ok && currentUnit.Exists {
			validUnits = append(validUnits, unitID)
		}
	}

	bld.AssignedUnits = validUnits
	bld.Food = uint8(len(bld.AssignedUnits))
}

func (bld *building) getOccupiedTiles() []point {
	return bld.OccupiedTiles
}

// cx, cy, ok := bld.getCenter() if ok {}
func (bld *building) getCenter() (uint8, uint8, bool) {
	if len(bld.OccupiedTiles) == 0 {
		return 0, 0, false
	}

	minX := bld.OccupiedTiles[0].X
	minY := bld.OccupiedTiles[0].Y
	maxX := bld.OccupiedTiles[0].X
	maxY := bld.OccupiedTiles[0].Y

	for _, pt := range bld.OccupiedTiles {
		if pt.X < minX {
			minX = pt.X
		}
		if pt.Y < minY {
			minY = pt.Y
		}
		if pt.X > maxX {
			maxX = pt.X
		}
		if pt.Y > maxY {
			maxY = pt.Y
		}
	}

	centerX := minX + (maxX-minX)/2
	centerY := minY + (maxY-minY)/2

	return centerX, centerY, true
}

func (bld *building) getDistanceToUnit(unitX, unitY uint8) uint8 {
	// log.Println("Środek getDistanceToUnit")

	if len(bld.OccupiedTiles) == 0 {
		// log.Println("bld.OccupiedTiles == 0")
		// log.Println(math.MaxUint8)

		return math.MaxUint8
	}

	// log.Println("Obliczam minDist")

	minDist := uint8(math.MaxUint8)

	// log.Printf("minDist: %d", minDist)

	for _, bldTile := range bld.OccupiedTiles {
		dx := math.Abs(float64(unitX) - float64(bldTile.X))
		dy := math.Abs(float64(unitY) - float64(bldTile.Y))
		currentDist := uint8(math.Max(dx, dy))

		// log.Printf("dx: %f, dy: %f, currentDist: %d", dx, dy, currentDist)

		if currentDist < minDist {
			minDist = currentDist

			// log.Println("currentDist < minDist")
		}
	}

	// log.Printf("Wychodzę z getDistanceToUnit, obliczone minDist %d", minDist)

	return minDist
}

// Zwraca granice budynku jako int.
func (bld *building) getBounds() (int, int, int, int) {
	if len(bld.OccupiedTiles) == 0 {
		return 0, 0, 0, 0
	}

	minX, minY := math.MaxUint8, math.MaxUint8
	maxX, maxY := 0, 0

	for _, occupiedTile := range bld.OccupiedTiles {
		tileX, tileY := int(occupiedTile.X), int(occupiedTile.Y)
		if tileX < minX {
			minX = tileX
		}

		if tileY < minY {
			minY = tileY
		}

		if tileX > maxX {
			maxX = tileX
		}

		if tileY > maxY {
			maxY = tileY
		}
	}

	return minX, minY, maxX, maxY
}

// Sprawdza, czy na danym polu można postawić jednostkę
func (bld *building) isValidSpawnTile(x, y int, bState *battleState) bool {
	// 1. Czy mieści się na mapie?
	if x < 0 || x >= int(boardMaxX) || y < 0 || y >= int(boardMaxY) {
		return false
	}

	currentTile := &bState.Board.Tiles[x][y]

	// 2. Czy teren jest przechodni?
	if !currentTile.IsWalkable {
		return false
	}

	// 2. Czy pole jest zajęte przez jednostkę? (zawsze blokuje)
	if currentTile.Unit != nil {
		return false
	}

	// 3. Czy pole jest zajęte przez budynek?
	if currentTile.Building != nil {
		if currentTile.Building.Type == buildingBridge && !currentTile.Building.IsUnderConstruction {
			return true
		}

		if currentTile.Building.Type == buildingPalisade && currentTile.Building.IsUnderConstruction {
			return true
		}

		return false
	}

	return true
}

func (bld *building) getClosestWalkableTile(bState *battleState) (uint8, uint8, bool) {
	if len(bld.OccupiedTiles) == 0 {
		return 0, 0, false
	}

	// 1. Pobieramy granice
	minX, minY, maxX, maxY := bld.getBounds()

	// 2. Iterujemy z marginesem 1 kratki
	for y := minY - 1; y <= maxY+1; y++ {
		for x := minX - 1; x <= maxX+1; x++ {

			// Pomijamy wnętrze budynku (interesuje nas tylko obwódka)
			if x >= minX && x <= maxX && y >= minY && y <= maxY {
				continue
			}

			// Sprawdzamy, czy to dobre miejsce
			if bld.isValidSpawnTile(x, y, bState) {
				return uint8(x), uint8(y), true
			}
		}
	}

	return 0, 0, false
}

func (bld *building) getOptimalRangedAttackTile(unitX, unitY, attackRange uint8, board *boardData) (uint8, uint8, bool) {
	centerX, centerY, ok := bld.getCenter()
	if !ok {
		return 0, 0, false
	}
	candidates := []point{
		{centerX + attackRange, centerY},
		{centerX - attackRange, centerY},
		{centerX, centerY + attackRange},
		{centerX, centerY - attackRange},
		{centerX + attackRange, centerY + attackRange},
		{centerX - attackRange, centerY - attackRange},
		{centerX + attackRange, centerY - attackRange},
		{centerX - attackRange, centerY + attackRange},
	}

	if attackRange > 2 {
		halfRange := attackRange / 2
		candidates = append(candidates, []point{
			{centerX + halfRange, centerY},
			{centerX - halfRange, centerY},
			{centerX, centerY + halfRange},
			{centerX, centerY - halfRange},
		}...)
	}

	bestX, bestY := uint8(255), uint8(255)
	minDistance := math.MaxFloat64

	for _, candidate := range candidates {
		if !bld.isValidWalkableTile(candidate.X, candidate.Y, board) {
			continue
		}

		if bld.getDistanceToUnit(candidate.X, candidate.Y) <= attackRange {
			distance := math.Abs(float64(unitX-candidate.X)) + math.Abs(float64(unitY-candidate.Y))

			if distance < minDistance {
				minDistance = distance
				bestX, bestY = candidate.X, candidate.Y
			}
		}
	}

	return bestX, bestY, true
}

func (bld *building) isValidWalkableTile(x, y uint8, board *boardData) bool {
	if x >= boardMaxX || y >= boardMaxY {
		return false
	}

	currentTile := &board.Tiles[x][y]

	// Jest przejezdne I nie ma tam nikogo
	return currentTile.IsWalkable && currentTile.Unit == nil && currentTile.Building == nil
}

func (bld *building) takeDamage(damage uint16) {
	if !bld.Exists {
		return
	}

	bld.AccumulatedDamage += damage
	log.Printf("building %d received %d damage (accumulated: %d)", bld.ID, damage, bld.AccumulatedDamage)
}

// @todo: trzeba to wpiąć w logikę odblokowywania przycisków dla SI oraz rysowania dla ludzi.
func (bld *building) allowedUnitTypes(unitType unitType, bState *battleState) bool {
	switch bld.Type {
	case buildingBarn:
		if unitType == unitCow {
			return true
		}

		if unitType == unitShepherd {
			return shepherdLevel >= bState.CurrentLevel
		}
	case buildingTemple:
		if unitType == unitMage {
			return bState.CurrentLevel >= mageLevel
		}
	case buildingAcademy:
		if unitType == unitCrossbowman {
			return bState.CurrentLevel >= crossbowmanLevel
		}
	default:
		return true
	}

	return true
}

func (bState *battleState) getPlayerState(ownerID uint8) *playerState {
	if ownerID == bState.HumanPlayerState.PlayerID {
		return bState.HumanPlayerState
	}

	return bState.AIEnemyState
}

func (bld *building) canProduceUnit(unitType unitType, bState *battleState) bool {
	reject := func(reason string) bool {
		if bld.Owner == bState.PlayerID {
			bState.CurrentMessage.Text = reason
			bState.CurrentMessage.Duration = 60
		}

		return false
	}

	fmt.Println("sprawdzam, czy da się zrobić jednostkę.")

	// 1. Czy jest wolne miejsce w budynku?
	if !bld.hasSpace() {
		return reject("Brak miejsca w budynku!")
	}

	fmt.Println("jest miejsce w budynku")
	// Ustalamy kto chce wykonać działanie
	ownerState := bState.getPlayerState(bld.Owner)

	// 2. Czy nie przekraczamy odgórnego ograniczenia?
	if ownerState.CurrentPopulation >= maxUnitsPerPlayer {
		return reject("Ograniczenie jednostek!")
	}

	fmt.Println("nie przekraczam populacji")
	// Pobieramy dane o jednostce do stworzenia
	stats, ok := unitDefs[unitType]

	if !ok {
		log.Printf("BŁĄD KRYTYCZNY: Brak definicji dla jednostki ID %d w unitDefs. Budynek ID: %d", unitType, bld.ID)
		reject("Brak definicji jednostki!")

		return false
	}

	fmt.Println("pobrałem statystyki jednostki")

	// 3. Czy stać gracza?
	if ownerState.Milk < stats.Cost {
		return reject("Za mało mleka!")
	}

	fmt.Println("jest dość mleka na zrobienie")
	// 4. Czy jednostka może wyjść z budynku?
	if _, _, ok := bld.getClosestWalkableTile(bState); !ok {
		fmt.Println("getClosestWalkableTile FALSE")

		return reject("Wyjście z budynku zastawione!")
	}

	fmt.Println("jest miejsce obok budynku")

	fmt.Println("wszystkie warunki spełnione")

	return true
}

func (bld *building) spawnUnit(unitType unitType, spawnX, spawnY uint8, bState *battleState) {
	newUnit := &unit{}
	newUnit.initUnit(unitType, spawnX, spawnY, cmdUIdle, bState)
	newUnit.Owner = bld.Owner
	newUnit.BelongsTo = bld

	newUnit.show(bState.Board)

	bState.Units = append(bState.Units, newUnit)
	bld.registerUnit(newUnit.ID)

	if bld.Owner == bState.HumanPlayerState.PlayerID {
		bState.HumanPlayerState.CurrentPopulation++
	} else {
		bState.AIEnemyState.CurrentPopulation++
	}

	// W sumie, to może należałoby odwrócić logikę, bo != jest dużo częsciej?
	if newUnit.Type == unitCow {
		// Bez tego nowa krowa stoi bezczynnie
		newUnit.Command = cmdUGraze
	}

	log.Printf("DEBUG: Stworzono jednostkę. Populacja Gracza: %d, AI: %d",
		bState.HumanPlayerState.CurrentPopulation, bState.AIEnemyState.CurrentPopulation)
}

// produceUnit odpowiada za próbę wytworzenia jednostki.
func (bld *building) produceUnit(newUnitType unitType, bState *battleState) {
	// 1. Sprawdzamy, czy są jakieś przeszkody w stworzeniu jednostki
	if !bld.canProduceUnit(newUnitType, bState) {
		return
	}

	// 2. Weryfikujemy, czy taki rodzaj jednostki istnieje
	stats, ok := unitDefs[newUnitType]

	if !ok {
		panic(fmt.Sprintf("BŁĄD KRYTYCZNY: Brak definicji dla jednostki %d w unitDefs", newUnitType))
	}

	// 3. Ustalamy właściciela
	owner := bState.HumanPlayerState
	if bld.Owner == bState.AIPlayerID {
		owner = bState.AIEnemyState
	}

	// 4. Pobieramy mleko za jednostkę
	owner.Milk -= stats.Cost

	// 5. Tworzymy jednostkę

	spawnX, spawnY, ok := bld.getClosestWalkableTile(bState)
	if ok {
		bld.spawnUnit(newUnitType, spawnX, spawnY, bState)
		log.Printf("INFO: Budynek ID %d zrobił jednostkę typu %v. Mleka gracza: %d.", bld.ID, newUnitType, owner.Milk)
	}
}

// decreaseHPBuilding dla każdej istniejącej budowli zmniejsza PŻ o amount
// Pilnuje, aby ustawić bld.Exists = false
func (bld *building) decreaseHPBuilding(amount uint16) {
	if !bld.Exists {
		return
	}

	bld.HP -= amount
	if bld.HP <= 0 {
		bld.HP = 0
		bld.Exists = false
	}
}

// increaseHPBuilding dla każdej istniejącej budowli zwiększa PŻ o amount
// Pilnuje, aby bld.HP <= bld.MaxHP; Służy do naprawy budynków
func (bld *building) increaseHPBuilding(amount uint16) {
	if !bld.Exists {
		return
	}

	bld.HP += amount
	if bld.HP >= bld.MaxHP {
		bld.HP = bld.MaxHP
	}
}

// getButtonCommand zastępuje przestarzałe GetProductionCommand.
// Tłumaczy kliknięcie przycisku (actionIndex) na pełny rozkaz (command).
// @todo: sprawdź, czy te actionIndex muszą być zaczarodziejskie. 01.07.2026
func (bld *building) getButtonCommand(actionIndex int) command {
	// Domyślny, pusty rozkaz
	cmd := command{ActionType: cmdUIdle}

	switch bld.Type {
	case buildingMain:
		// Indeks 6: Budowa drogi/palisady (w zależności od kontekstu UI)
		if actionIndex == 6 {
			cmd.ActionType = cmdBPlaceConstruction
			cmd.InteractionTargetID = uint(buildingPalisade)
		}

	case buildingBarn:
		// Indeks 5: Wytwarzanie Krowy
		if actionIndex == 5 {
			cmd.ActionType = cmdBProduce
			cmd.CreateType = uint8(unitCow)
		}
		// Indeks 6: Budowa nowej Obory
		if actionIndex == 6 {
			cmd.ActionType = cmdBPlaceConstruction
			cmd.InteractionTargetID = uint(buildingBarn)
		}

	case buildingBarracks:
		// Indeks 4: Wytwarzanie Łucznika
		if actionIndex == 4 {
			cmd.ActionType = cmdBProduce
			cmd.CreateType = uint8(unitArcher)
		}
		// Indeks 5: Wytwarzanie Drwala
		if actionIndex == 5 {
			cmd.ActionType = cmdBProduce
			cmd.CreateType = uint8(unitAxeman)
		}
		// Indeks 6: Budowa Chaty Mieszkalnej
		if actionIndex == 6 {
			cmd.ActionType = cmdBPlaceConstruction
			cmd.InteractionTargetID = uint(buildingBarracks)
		}

	case buildingTemple:
		// Indeks 4: Wytwarzanie Kapłana
		if actionIndex == 4 {
			cmd.ActionType = cmdBProduce
			cmd.CreateType = uint8(unitPriest)
		}
		// Indeks 5: Wytwarzanie Kapłanki
		if actionIndex == 5 {
			cmd.ActionType = cmdBProduce
			cmd.CreateType = uint8(unitPriestess)
		}
		// Indeks 6: Tutaj był stary CMD_PRODUCE bez typu, zakładam, że to błąd starego kodu
		// lub puste miejsce. Zostawiamy IDLE.

	case buildingBarracks2:
		// Indeks 4: Wytwarzanie Włócznika
		if actionIndex == 4 {
			cmd.ActionType = cmdBProduce
			cmd.CreateType = uint8(unitSpearman)
		}
		// Indeks 5: Wytwarzanie Miecznika
		if actionIndex == 5 {
			cmd.ActionType = cmdBProduce
			cmd.CreateType = uint8(unitSwordsman)
		}
		// Indeks 6: Pusty w starym kodzie (zwracał gołe CMD_PRODUCE)
		// Indeks 7: Budowa Palisady (stare CMD_BUILD_FENCE)
		if actionIndex == 7 {
			cmd.ActionType = cmdBPlaceConstruction
			cmd.InteractionTargetID = uint(buildingPalisade)
		}

	case buildingAcademy:
		// Indeks 4: Wytwarzanie Kusznika
		if actionIndex == 4 {
			cmd.ActionType = cmdBProduce
			cmd.CreateType = uint8(unitCrossbowman)
		}
		// Indeks 5: Wytwarzanie Dowódcy
		if actionIndex == 5 {
			cmd.ActionType = cmdBProduce
			cmd.CreateType = uint8(unitCommander)
		}
		// Indeks 6: Pusty
	default:
		panic("unhandled default case")
	}

	return cmd
}

func (bld *building) bounds() bounds {
	minX, minY, maxX, maxY := bld.occupiedTilesBounds()

	widthTiles := maxX - minX + 1
	heightTiles := maxY - minY + 1

	return bounds{
		X:        int32(minX) * int32(tileWidth),
		Y:        int32(minY) * int32(tileHeight),
		Width:    int32(widthTiles) * int32(tileWidth),
		Height:   int32(heightTiles) * int32(tileHeight),
		WidthPx:  float32(widthTiles) * float32(tileWidth),
		HeightPx: float32(heightTiles) * float32(tileHeight),
	}
}

func (bld *building) isRepairable(playerID uint8) bool {
	if bld == nil || !bld.Exists || bld.HP >= bld.MaxHP {
		return false
	}

	return bld.Type == buildingPalisade || bld.Type == buildingBridge || bld.Owner == playerID
}

func cleanupDestroyedBuildings(bState *battleState) {
	if bState.GlobalFrameCounter%6000 != 0 {
		return
	}

	if len(bState.Buildings) < int(maxBuildingsPerPlayer)*4 {
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

	if removedCount > 0 {
		log.Printf("INFO: Wyczyszczono %d zniszczonych budynków.", removedCount)
	}
}

func cleanupConvertedBuildings(bState *battleState) {
	newBuildingList := make([]*building, 0, len(bState.Buildings))

	for _, bld := range bState.Buildings {
		if bld.IsPendingRemoval {
			// @reminder: Nie da się zasadzić budowy poza granicami
			// planszy więc nie trzeba sprawdzać przy usuwaniu
			// czy zajmowane kafelki mieszczą się w planszy
			currentTile := &bState.Board.Tiles[bld.OccupiedTiles[0].X][bld.OccupiedTiles[0].Y]
			currentTile.Building = nil
			currentTile.IsWalkable = true
		} else {
			// zachowujemy budynek
			newBuildingList = append(newBuildingList, bld)
		}
	}
	// lista budynków staje się listą zachowanych budowli
	bState.Buildings = newBuildingList
}

func applyPhase2Graphics(bld *building, board *boardData) {
	if bld.Type == buildingPalisade {
		return
	}

	template, ok := constructionTemplatesPhase02[bld.Type]
	if !ok {
		return
	}

	minX, minY := bld.OccupiedTiles[0].X, bld.OccupiedTiles[0].Y

	for dy, row := range template {
		for dx, texID := range row {
			tx, ty := uint16(minX)+uint16(dx), uint16(minY)+uint16(dy)
			board.Tiles[tx][ty].TextureID = texID
		}
	}
}

func completeConstruction(bld *building, board *boardData, playerID uint8) bool {
	bld.IsUnderConstruction = false
	bld.HP = bld.MaxHP

	applyFinishedGraphics(bld, board)

	if bld.Type == buildingBridge {
		bld.IsPendingRemoval = true
	}

	return bld.Owner == playerID
}

func applyFinishedGraphics(bld *building, board *boardData) {
	switch bld.Type {
	case buildingPalisade:
		if bld.Type == buildingPalisade {
			pt := bld.OccupiedTiles[0]
			joinPalisade(pt.X, pt.Y, board)
		}

		return

	case buildingBridge:
		board.Tiles[bld.OccupiedTiles[0].X][bld.OccupiedTiles[0].Y].IsWalkable = true
		board.Tiles[bld.OccupiedTiles[0].X][bld.OccupiedTiles[0].Y].TextureID = spriteBridge01

		return

	default:
		template, ok := buildingTemplates[bld.Type]
		if !ok {
			fmt.Println("Bład przy próbie zastosowania grafiki ukończonej budowy.")
		}

		minX, minY := bld.OccupiedTiles[0].X, bld.OccupiedTiles[0].Y

		for dy, row := range template {
			for dx, texID := range row {
				tx, ty := minX+uint8(dx), minY+uint8(dy)
				if tx < boardMaxX && ty < boardMaxY {
					board.Tiles[tx][ty].TextureID = uint16(texID)
				}
			}
		}
	}
}

func placeConstructionOnBoard(bld *building, x, y uint8, board *boardData) {
	// 1. Sprawdzamy, czy dany rodzaj budynku był określony wcześniej
	stats, ok := buildingDefs[bld.Type]
	if !ok {
		log.Printf("BŁĄD KRYTYCZNY: Brak określenia dla budynku rodzaju %d!", bld.Type)

		return
	}
	// 3. Zajmowanie kafelków
	bld.OccupiedTiles = make([]point, 0, stats.Width*stats.Height)

	startX := x
	startY := y

	for ox := range stats.Width {
		for oy := range stats.Height {
			px, py := startX+ox, startY+oy

			if px < boardMaxX && py < boardMaxY {
				bld.OccupiedTiles = append(bld.OccupiedTiles, point{X: px, Y: py})

				occupieTile := &board.Tiles[px][py]
				occupieTile.Building = bld
				occupieTile.IsWalkable = false
			}
		}
	}
}
