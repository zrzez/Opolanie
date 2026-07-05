package main

import (
	"fmt"
	"log"
	"math"
)

// constructions.go

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

// Czyścimy plac budowy zwykłych budynków ze znisczonych/nieukończonych palisad
func (board *boardData) clearConstructionSite(tileX, tileY uint8, bldStats buildingStats) {
	// Wszystkie kafelki, które będą zajęte przez plac budowy
	for dx := range bldStats.Width {
		for dy := range bldStats.Height {
			cx, cy := tileX+dx, tileY+dy

			// Wybieramy konkretny kafelek do sprawdzenia
			currentTile := &board.Tiles[cx][cy]

			// Przez walidację przepuszczamy palisady w budowie
			// Jeśli nil, to nie mamy co robić
			if currentTile.Building == nil {
				continue
			}

			// Jeśli !nil, to isFreeForConstruction gwarantuje, że
			// trafiliśmy na palisadę w budowie i trzeba się jej pozbyć
			// 1. Ustawiamy budynek do usunięcia z listy []*buildings w battleState
			currentTile.Building.Exists = false
			// 2. Wywalamy tenże budynek z planszy
			currentTile.Building = nil
		}
	}
}

func (board *boardData) placeRoad(tileX, tileY uint8) {
	board.Tiles[tileX][tileY].TextureID = spriteRoadStart
	cx := int(tileX)
	cy := int(tileY)

	refreshRoadTile(cx, cy, board)

	refreshRoadTile(cx+1, cy, board) // prawo
	refreshRoadTile(cx-1, cy, board) // lewo
	refreshRoadTile(cx, cy+1, board) // góra
	refreshRoadTile(cx, cy-1, board) // dół

	log.Printf("BUDOWA: Postawiono drogę na (%d,%d).", tileX, tileY)
}

// Odpowiada za dodanie stworzonej jednostki do mieszkańców budynku.
// Sprawdzenie zostało wykonane w canProduceUnit,
// pojemność za pomocą hasSpace
// Zwracany bool jest całkowicie ignorowany
func (bld *building) registerUnit(uID UnitID) bool {
	// Budynek poszerza listę zameldowanych jednostek
	// To chyba powinny być wskaźniki
	// @todo: sprawdź po cholerę mi w ogole ta lista
	// przecież samo śledzenie licznika domowników można zrobić prościej
	// @reminder: jeśli znajdę sposób na przypisane krowy do obory
	// najprawdopodobniej zniknie jedyna przyczyna dla której ta lista istnieje
	bld.AssignedUnits = append(bld.AssignedUnits, uID)
	// Budynek aktualizuje licznik posiadania
	bld.Food++

	return true
}

// Wywoływana przez u.unregisterFromBuilding gdy jednostka zmarła
// Zwracane bool jest ignorowane
func (bld *building) unregisterUnit(unregisterUnitID UnitID) bool {
	// Przechodzimy przez listę jednostek zamieszkujących
	for index, registeredUnitID := range bld.AssignedUnits {
		if registeredUnitID == unregisterUnitID {
			// po znalezieniu miejsca w którym znajduje się jednostka, pomijamy ją przy
			// odświeżeniu listy
			bld.AssignedUnits = append(bld.AssignedUnits[:index], bld.AssignedUnits[index+1:]...)

			bld.Food--

			return true
		}
	}
	// @reminder: co z jednostkami, które dostaliśmy z początek wyprawy?
	return false
}

func (bld *building) hasSpace() bool {
	return len(bld.AssignedUnits) < int(bld.MaxFood)
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

func (bState *battleState) getPlayerState(ownerID PlayerID) *playerState {
	if ownerID == bState.HumanPlayerState.PlayerID {
		return bState.HumanPlayerState
	}

	return bState.AIEnemyState
}

func (bld *building) canProduceUnit(unitType unitType, bState *battleState) bool {
	// @todo: co to do cholery jest? Zupełnie nie pamiętam - 05.07.2026
	reject := func(reason string) bool {
		if bld.Owner == bState.PlayerID {
			bState.CurrentMessage.Text = reason
			bState.CurrentMessage.Duration = 60
		}

		return false
	}

	fmt.Println("sprawdzam, czy da się stworzyć jednostkę.")

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
			cmd.InteractionTargetID = ObjectID(buildingPalisade)
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
			cmd.InteractionTargetID = ObjectID(buildingBarn)
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
			cmd.InteractionTargetID = ObjectID(buildingBarracks)
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
			cmd.InteractionTargetID = ObjectID(buildingPalisade)
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
		panic("getButtonCommand unhandled default case")
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

func (bld *building) isRepairable(playerID PlayerID) bool {
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

func cleanupConvertedBridges(bState *battleState) {
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

func completeConstruction(bld *building, board *boardData, playerID PlayerID) bool {
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

func (board *boardData) registerBuilding(bld *building, startX, startY uint8) error {
	stats, ok := buildingDefs[bld.Type]
	if !ok {
		return fmt.Errorf("nieznany rodzaj budynku %d", bld.Type)
	}

	bld.OccupiedTiles = make([]point, 0, stats.Width*stats.Height)

	for ox := range stats.Width {
		for oy := range stats.Height {
			px, py := startX+ox, startY+oy

			// Bezpiecznik przed wyjściem poza planszę
			// @todo: sprawdź, czy jest potrzebny
			// @todo: czemu przepuszczamy a nie wypierdzielamy?
			// przecież wyjście poza planszę, to błąd krytyczny!
			if px >= boardMaxX || py >= boardMaxY {
				continue
			}

			currentTile := &board.Tiles[px][py]
			// @reminder: sprawdź, czy poszła walidacja wcześniej!
			// jeśli nie to może się napatoczyć budynek i zablokować
			bld.OccupiedTiles = append(bld.OccupiedTiles, point{X: px, Y: py})

			currentTile.Building = bld
			currentTile.IsWalkable = false
		}
	}

	return nil
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
	newBld.applyConstructionGraphics(bState.Board)

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

// applyConstructionGraphics nakłada tekstury budowy na kafelki zajmowane przez budynek.
// Metoda ta zakłada, że bld.OccupiedTiles jest już poprawnie zainicjalizowane
// (zrobione przez boardData.RegisterBuilding).
func (bld *building) applyConstructionGraphics(board *boardData) {
	for index, tilePoint := range bld.OccupiedTiles {
		currentTile := &board.Tiles[tilePoint.X][tilePoint.Y]

		// Ustalenie tekstury i właściwości w zależności od typu budynku
		switch bld.Type {
		case buildingBridge:
			currentTile.TextureID = spriteBridgeConstruction

		case buildingPalisade:
			currentTile.TextureID = spritePalisadeDestroyed
			// Wyjątek: zniszczona palisada w budowie jest przechodnia
			currentTile.IsWalkable = true

		default:
			// Standardowe budynki - użycie szablonu budowy (faza 01)
			row := index / int(regularBuildingSize)
			column := index % int(regularBuildingSize)
			currentTile.TextureID = constructionTemplatePhase01[column][row]
		}
	}
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
	applyFinishedGraphics(newBld, bState.Board)

	return newBld
}
