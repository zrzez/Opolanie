package main

import (
	"fmt"
	"log"
	"math"
	"slices"
)

// constructions.go

// initConstruction odpowiada tylko za techniczne utworzenie obiektu w pamięci i na planszy.
// Nie mylić z budowaniem nowych w czasie bitwy. To służy tylko ładowaniu z mapy lub zapisu gry.
func (bld *building) initConstruction(x, y uint8, buildingType buildingType, owner uint8, bs *battleState) {
	// 1. Sprawdzamy, czy dany rodzaj budynku był określony wcześniej
	stats, ok := buildingDefs[buildingType]
	if !ok {
		log.Printf("BŁĄD KRYTYCZNY: Brak określenia dla budynku rodzaju %d!", bld.Type)
		return
	}

	// 2. Podstawowe właściwości
	bld.ID = bs.NextUniqueObjectID
	bs.NextUniqueObjectID++
	bld.Type = buildingType
	bld.Owner = owner
	bld.Exists = true
	bld.Armor = buildingArmor
	bld.MaxHP = stats.MaxHP
	bld.HP = stats.MaxHP
	bld.MaxFood = stats.MaxFood

	bld.AssignedUnits = make([]uint, 0)

	// 3. Zajmowanie kafelków
	bld.OccupiedTiles = make([]point, 0, stats.Width*stats.Height)

	startX := x
	startY := y

	for ox := uint8(0); ox < stats.Width; ox++ {
		for oy := uint8(0); oy < stats.Height; oy++ {
			px, py := startX+ox, startY+oy

			if px < boardMaxX && py < boardMaxY {
				bld.OccupiedTiles = append(bld.OccupiedTiles, point{X: px, Y: py})

				tile := &bs.Board.Tiles[px][py]
				tile.Building = bld
				tile.IsWalkable = false
			}
		}
	}
}

// setAsFinished ustawia budynek jako gotowy (przy wczytywaniu gry, czy nowej bitwy)
// currentHP = -1 dla MaxHP (nowa mapa) lub inna, aby odtworzyć konkretną wartość HP
func (bld *building) setAsFinished(bs *battleState, currentHP uint16) {
	bld.IsUnderConstruction = false

	// Wczytanie obecnego stanu z zapisu gry
	if currentHP > 0 {
		bld.HP = min(currentHP, bld.MaxHP)
	} else {
		bld.HP = bld.MaxHP
	}

	bld.applyFinishedGraphics(bs)
}

// startConstruction odpowiada za postawienie pola budowy
func (bld *building) startConstruction(bs *battleState) {
	// 1. Pobieramy dane i bronimy się przed głupim błędem
	stats, ok := buildingDefs[bld.Type]

	if !ok {
		log.Printf("BŁĄD KRYTYCZNY: Nie ma %s w buildingDefs, nie można rozpocząć budowy", stats.Name)
	}

	bld.IsUnderConstruction = true

	// Nie zaczynamy z zerem HP
	bld.HP = initialConstructionHP

	for index, tilePoint := range bld.OccupiedTiles {
		if tilePoint.X >= boardMaxX || tilePoint.Y >= boardMaxY {
			continue
		}

		row := index / int(normalBuildingSize)
		column := index % int(normalBuildingSize)

		bs.Board.Tiles[tilePoint.X][tilePoint.Y].TextureID = constructionTemplatePhase01[column][row]
		bs.Board.Tiles[tilePoint.X][tilePoint.Y].Building = bld
		bs.Board.Tiles[tilePoint.X][tilePoint.Y].IsWalkable = false
	}
}

// applyWork wywoływana przez UNIT_AXE
// TODO: Budowa powinna dodawać określoną liczbę HP
// Różna dla człowieka i SI!
func (bld *building) applyWork(amount uint16, bs *battleState) bool {
	if !bld.Exists {
		return false
	}

	// SCENARIUSZ 1: BUDOWA
	if bld.IsUnderConstruction {
		// 1. Aktualizacja postępu

		// 2. Aktualizacja HP wg sztywnych zasad
		hpGain := uint16(0)

		// TODO: Chyba powinienem to wrzucić jako stałe
		// Sprawdzamy, czy buduje człowiek czy SI
		if bld.Owner == bs.HumanPlayerState.PlayerID {
			// TODO: zaczarowane liczby!
			hpGain = 2 // Gracz: +2 HP
		} else {
			// TODO: zaczarowane liczby!
			hpGain = 5 // SI: +5 HP
		}

		bld.increaseHPBuilding(hpGain)

		// 3. Sprawdzenie końca budowy
		if bld.HP >= bld.MaxHP/2 {
			bld.applyPhase2Graphics(bs)
		}

		if bld.HP >= bld.MaxHP {
			bld.completeConstruction(bs)
		}

		return true
	}

	// SCENARIUSZ 2: NAPRAWA
	// Naprawa działa doskonale
	if bld.HP < bld.MaxHP {
		bld.increaseHPBuilding(amount)

		return true
	}

	return false
}

// completeConstruction domyka budowę, zmienia flagi grafikę
func (bld *building) completeConstruction(bs *battleState) {
	bld.IsUnderConstruction = false
	bld.HP = bld.MaxHP

	bld.applyFinishedGraphics(bs)

	if bld.Owner == bs.PlayerID {
		bs.CurrentMessage.Text = "Budowa %s zakończona!"
		bs.CurrentMessage.Duration = 60
	}

	log.Printf("INFO: Budynek ID %d (Typ %d) został ukończony.", bld.ID, bld.Type)
}

// applyFinishedGraphics nakłada docelową teksturę na kafelki
// TODO: jeżeli HP == MaxHP/2 to tekstury powinny zmienić się
// z ogólnej budowy na wybraną dla danego rodzaju budowli
// TODO: jeżeli HP == MAXHP to budowa jest zakończona.
// Nie tylko budynek zmienia teksturę, ale też normalnie działa!
// TODO: w oryginale raz zmieniona tekstura z podwalin na połowiczną
// już nie wraca na podwaliny pomimo spadku HP. Zostawić tak?
func (bld *building) applyFinishedGraphics(bs *battleState) {
	template, ok := buildingTemplates[bld.Type]
	if !ok {
		// TODO: ogarnij, czy nie lepiej jest dodać palisady do BUILDING_TEMPLATES
		// Takie obsłużenie palisad i mostów jest ryzykowne!
		if bld.Type == buildingPalisade && len(bld.OccupiedTiles) > 0 {
			pt := bld.OccupiedTiles[0]
			joinRoadsPalisade(pt.X, pt.Y, 266, bs)
		}
		return
	}

	minX, minY := bld.OccupiedTiles[0].X, bld.OccupiedTiles[0].Y

	for dy, row := range template {
		for dx, texID := range row {
			tx, ty := minX+uint8(dx), minY+uint8(dy)
			if tx < boardMaxX && ty < boardMaxY {
				bs.Board.Tiles[tx][ty].TextureID = uint16(texID)
			}
		}
	}
}

func (bld *building) applyPhase2Graphics(bs *battleState) {
	template, ok := constructionTemplatesPhase02[bld.Type]
	if !ok {
		return
	}

	minX, minY := bld.OccupiedTiles[0].X, bld.OccupiedTiles[0].Y

	for dy, row := range template {
		for dx, texID := range row {
			tx, ty := uint16(minX)+uint16(dx), uint16(minY)+uint16(dy)
			if tx < uint16(boardMaxX) && ty < uint16(boardMaxY) {
				bs.Board.Tiles[tx][ty].TextureID = uint16(texID)
			}
		}
	}
}

// sprawdza, czy w danym miejscu (tileX, tileY)
// można postawić budynek, którego rodzaj mamy zapisany w pamięci (bs.PendingBuildingType).
func tryBuildStructure(bs *battleState, tileX, tileY uint8) {
	buildingType := bs.PendingBuildingType
	stats, ok := buildingDefs[buildingType]

	if !ok {
		return
	}

	// Limity i koszty bez zmian...
	if bs.HumanPlayerState.CurrentBuildings >= maxBuildingsPerPlayer {
		bs.CurrentMessage.Text = "Limit budynków!"
		bs.CurrentMessage.Duration = 60

		return
	}

	if bs.HumanPlayerState.Milk < stats.Cost {
		bs.CurrentMessage.Text = fmt.Sprintf("Brak mleka! (%d)", stats.Cost)
		bs.CurrentMessage.Duration = 60

		return
	}

	// Walidacja terenu (Teraz tileX, tileY to lewy górny róg)
	if !isValidConstructionSite(tileX, tileY, stats.Width, stats.Height, bs) {
		return
	}

	// === DECYZJA POZYTYWNA ===
	bs.HumanPlayerState.Milk -= stats.Cost
	bs.HumanPlayerState.CurrentBuildings++

	newBld := &building{}
	// init teraz przyjmie tileX, tileY jako lewy górny róg
	newBld.initConstruction(tileX, tileY, buildingType, bs.PlayerID, bs)

	bs.Buildings = append(bs.Buildings, newBld)
	newBld.startConstruction(bs)

	bs.CurrentMessage.Text = fmt.Sprintf("Wznoszenie: %s", stats.Name)
	bs.CurrentMessage.Duration = 60

	log.Printf("BUDOWA: Rozpoczęto %s (ID: %d) na (%d,%d).", stats.Name, newBld.ID, tileX, tileY)
}

// isValidConstructionSite sprawdza wszelkie warunki, które należy spełnić, aby można było zasadzić budowlę.
func isValidConstructionSite(tileX, tileY, width, height uint8, bs *battleState) bool {
	for dx := range width {
		for dy := range height {
			// Po prostu dodajemy offset do kursora. Czysto i przyjemnie.
			cx, cy := tileX+dx, tileY+dy

			// 1. Poza mapą?
			if cx < 0 || cx >= boardMaxX || cy < 0 || cy >= boardMaxY {
				bs.CurrentMessage.Text = "Poza mapą!"
				bs.CurrentMessage.Duration = 40
				return false
			}

			tile := &bs.Board.Tiles[cx][cy]

			// 2. Czy coś tu stoi?
			if tile.Unit != nil || tile.Building != nil {
				bs.CurrentMessage.Text = "Miejsce zajęte!"
				bs.CurrentMessage.Duration = 40
				return false
			}

			// 3. Czy teren nadaje się pod budowę?
			if isObstacle(tile.TextureID) {
				bs.CurrentMessage.Text = "Nie można na tym!"
				bs.CurrentMessage.Duration = 40
				return false
			}
		}
	}
	// 4. Czy sąsiaduje z drogą
	// Przekazujemy tileX, tileY jako lewy górny róg.
	// Pole budowy to zawsze kwadrat.
	if !isBorderingRoad(tileX, tileY, width, bs) {
		return false
	}

	return true
}

func isObstacle(texID uint16) bool {
	switch {
	case isRockNonWalkable(texID):
		return true
	case isBridge(texID):
		return true
	case isGadget(texID):
		return true
	case isTreeStump(texID):
		return true
	case texID >= spriteTreeBurntStump00 && texID <= spriteTreeBurntStump01:
		return true
	case isSpecialTile(texID):
		return true
	}

	return false
}

func isBorderingRoad(x, y, size uint8, bs *battleState) bool {
	// 1. Sprawdzamy boki LEWY i PRAWY (iterujemy po wysokości budynku)
	for i := uint8(0); i < size; i++ {
		currentY := y + i

		// Zabezpieczenie, gdyby budynek wystawał poza mapę dołem (teoretycznie niemożliwe po walidacji, ale bezpiecznie)
		if currentY >= boardMaxY {
			break
		}

		// --- LEWA KRAWĘDŹ (x-1) ---
		// Sprawdzamy tylko, jeśli nie jesteśmy przytuleni do lewej ściany mapy (x=0)
		// To zapobiega błędowi "uint8 overflow" (0 - 1 = 255)
		if x > 0 {
			if isRoad(bs.Board.Tiles[x-1][currentY].TextureID) {
				return true
			}
		}

		// --- PRAWA KRAWĘDŹ (x+size) ---
		if x+size < boardMaxX {
			if isRoad(bs.Board.Tiles[x+size][currentY].TextureID) {
				return true
			}
		}
	}

	// 2. Sprawdzamy boki GÓRA i DÓŁ (iterujemy po szerokości budynku)
	for i := uint8(0); i < size; i++ {
		currentX := x + i

		if currentX >= boardMaxX {
			break
		}

		// --- GÓRNA KRAWĘDŹ (y-1) ---
		// Sprawdzamy tylko, jeśli nie jesteśmy na samej górze mapy (y=0)
		if y > 0 {
			if isRoad(bs.Board.Tiles[currentX][y-1].TextureID) {
				return true
			}
		}

		// --- DOLNA KRAWĘDŹ (y+size) ---
		if y+size < boardMaxY {
			if isRoad(bs.Board.Tiles[currentX][y+size].TextureID) {
				return true
			}
		}
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

func (bld *building) getAssignedUnits(bs *battleState) []*unit {
	var units []*unit
	for _, unitID := range bld.AssignedUnits {
		unit, ok := getUnitByID(unitID, bs)
		if ok && unit.Exists {
			units = append(units, unit)
		}
	}
	return units
}

func (bld *building) cleanupDeadUnits(bs *battleState) {
	var validUnits []uint
	for _, unitID := range bld.AssignedUnits {
		unit, ok := getUnitByID(unitID, bs)
		if ok && unit.Exists {
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

	for _, tile := range bld.OccupiedTiles {
		tileX, tileY := int(tile.X), int(tile.Y)
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
func (bld *building) isValidSpawnTile(x, y int, bs *battleState) bool {
	// 1. Czy mieści się na mapie?
	if x < 0 || x >= int(boardMaxX) || y < 0 || y >= int(boardMaxY) {
		return false
	}

	tile := &bs.Board.Tiles[x][y]

	// 2. Czy teren jest przechodni?
	if !tile.IsWalkable {
		return false
	}

	// 3. Czy pole jest puste?
	if tile.Unit != nil || tile.Building != nil {
		return false
	}

	return true
}

func (bld *building) getClosestWalkableTile(bs *battleState) (uint8, uint8, bool) {
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
			if bld.isValidSpawnTile(x, y, bs) {
				return uint8(x), uint8(y), true
			}
		}
	}

	return 0, 0, false
}

func (bld *building) getOptimalRangedAttackTile(unitX, unitY, attackRange uint8, bs *battleState) (uint8, uint8, bool) {
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
		if !bld.isValidWalkableTile(candidate.X, candidate.Y, bs) {
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

func (bld *building) containsTile(x, y uint8) bool {
	for _, tile := range bld.OccupiedTiles {
		if tile.X == x && tile.Y == y {
			return true
		}
	}
	return false
}

func (bld *building) isValidWalkableTile(x, y uint8, bs *battleState) bool {
	if x < 0 || x >= boardMaxX || y < 0 || y >= boardMaxY {
		return false
	}

	// ZMIANA: Użycie nowej struktury Tiles zamiast Plc/IsWalkable(old)
	tile := &bs.Board.Tiles[x][y]

	// Jest przejezdne I nie ma tam nikogo
	return tile.IsWalkable && tile.Unit == nil && tile.Building == nil
}

func (bld *building) takeDamage(damage uint16) {
	if !bld.Exists {
		return
	}
	bld.AccumulatedDamage += damage
	log.Printf("building %d received %d damage (accumulated: %d)", bld.ID, damage, bld.AccumulatedDamage)
}

// @todo: trzeba to wpiąć w logikę odblokowywania przycisków dla SI oraz rysowania dla ludzi.
func (bld *building) allowedUnitTypes(unitType unitType, bs *battleState) bool {
	switch bld.Type {
	case buildingBarn:
		if unitType == unitCow {
			return true
		}
		if unitType == unitShepherd {
			return shepherdLevel >= bs.CurrentLevel
		}
	case buildingTemple:
		if unitType == unitMage {
			return bs.CurrentLevel >= mageLevel
		}
	case buildingAcademy:
		if unitType == unitCrossbowman {
			return bs.CurrentLevel >= crossbowmanLevel
		}
	default:
		return true
	}

	return true
}

func (bs *battleState) getPlayerState(ownerID uint8) *playerState {
	if ownerID == bs.HumanPlayerState.PlayerID {
		return bs.HumanPlayerState
	}

	return bs.AIEnemyState
}

func (bld *building) canProduceUnit(unitType unitType, bs *battleState) bool {
	reject := func(reason string) bool {
		if bld.Owner == bs.PlayerID {
			bs.CurrentMessage.Text = reason
			bs.CurrentMessage.Duration = 60
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
	ownerState := bs.getPlayerState(bld.Owner)

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
	if _, _, ok := bld.getClosestWalkableTile(bs); !ok {
		fmt.Println("getClosestWalkableTile FALSE")

		return reject("Wyjście z budynku zastawione!")
	}

	fmt.Println("jest miejsce obok budynku")

	fmt.Println("wszystkie warunki spełnione")

	return true
}

func (bld *building) spawnUnit(unitType unitType, spawnX, spawnY uint8, bs *battleState) {
	newUnit := &unit{}
	newUnit.initUnit(unitType, spawnX, spawnY, cmdIdle, bs)
	newUnit.Owner = bld.Owner
	newUnit.BelongsTo = bld

	newUnit.show(bs)

	bs.Units = append(bs.Units, newUnit)
	bld.registerUnit(newUnit.ID)

	if bld.Owner == bs.HumanPlayerState.PlayerID {
		bs.HumanPlayerState.CurrentPopulation++
	} else {
		bs.AIEnemyState.CurrentPopulation++
	}

	// W sumie, to może należałoby odwrócić logikę, bo != jest dużo częsciej?
	if newUnit.Type == unitCow {
		// Bez tego nowa krowa stoi bezczynnie
		newUnit.grazeCowPhase(bs)
	}

	log.Printf("DEBUG: Stworzono jednostkę. Populacja Gracza: %d, AI: %d",
		bs.HumanPlayerState.CurrentPopulation, bs.AIEnemyState.CurrentPopulation)
}

// produceUnit odpowiada za próbę wytworzenia jednostki.
func (bld *building) produceUnit(unitType unitType, bs *battleState) {
	// 1. Sprawdzamy, czy są jakieś przeszkody w stworzeniu jednostki
	if !bld.canProduceUnit(unitType, bs) {
		return
	}

	// 2. Weryfikujemy, czy taki rodzaj jednostki istnieje
	stats, ok := unitDefs[unitType]
	if !ok {
		panic(fmt.Sprintf("BŁĄD KRYTYCZNY: Brak definicji dla jednostki %d w unitDefs", unitType))
	}

	// 3. Ustalamy właściciela
	playerState := bs.HumanPlayerState
	if bld.Owner == bs.AIPlayerID {
		playerState = bs.AIEnemyState
	}

	// 4. Pobieramy mleko za jednostkę
	playerState.Milk -= stats.Cost

	// 5. Tworzymy jednostkę

	spawnX, spawnY, ok := bld.getClosestWalkableTile(bs)
	if ok {
		bld.spawnUnit(unitType, spawnX, spawnY, bs)
		log.Printf("INFO: Budynek ID %d zrobił jednostkę typu %v. Mleka gracza: %d.", bld.ID, unitType, playerState.Milk)
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
func (bld *building) getButtonCommand(actionIndex int) command {
	// Domyślny, pusty rozkaz
	cmd := command{ActionType: cmdIdle}

	switch bld.Type {
	case buildingMain:
		// Indeks 6: Budowa drogi/palisady (w zależności od kontekstu UI)
		if actionIndex == 6 {
			cmd.ActionType = cmdBuildStructure
			// Traktujemy to jako strukturę. Jeżeli droga nie ma osobnego buildingType,
			// używamy PALISADE lub innego odpowiednika.
			cmd.InteractionTargetID = uint(buildingPalisade)
		}

	case buildingBarn:
		// Indeks 5: Wytwarzanie Krowy
		if actionIndex == 5 {
			cmd.ActionType = cmdProduce
			cmd.ProduceType = unitCow
		}
		// Indeks 6: Budowa nowej Obory
		if actionIndex == 6 {
			cmd.ActionType = cmdBuildStructure
			cmd.InteractionTargetID = uint(buildingBarn)
		}

	case buildingBarracks:
		// Indeks 4: Wytwarzanie Łucznika
		if actionIndex == 4 {
			cmd.ActionType = cmdProduce
			cmd.ProduceType = unitArcher
		}
		// Indeks 5: Wytwarzanie Drwala
		if actionIndex == 5 {
			cmd.ActionType = cmdProduce
			cmd.ProduceType = unitAxeman
		}
		// Indeks 6: Budowa Chaty Mieszkalnej
		if actionIndex == 6 {
			cmd.ActionType = cmdBuildStructure
			cmd.InteractionTargetID = uint(buildingBarracks)
		}

	case buildingTemple:
		// Indeks 4: Wytwarzanie Kapłana
		if actionIndex == 4 {
			cmd.ActionType = cmdProduce
			cmd.ProduceType = unitPriest
		}
		// Indeks 5: Wytwarzanie Kapłanki
		if actionIndex == 5 {
			cmd.ActionType = cmdProduce
			cmd.ProduceType = unitPriestess
		}
		// Indeks 6: Tutaj był stary CMD_PRODUCE bez typu, zakładam, że to błąd starego kodu
		// lub puste miejsce. Zostawiamy IDLE.

	case buildingBarracks2:
		// Indeks 4: Wytwarzanie Włócznika
		if actionIndex == 4 {
			cmd.ActionType = cmdProduce
			cmd.ProduceType = unitSpearman
		}
		// Indeks 5: Wytwarzanie Miecznika
		if actionIndex == 5 {
			cmd.ActionType = cmdProduce
			cmd.ProduceType = unitSwordsman
		}
		// Indeks 6: Pusty w starym kodzie (zwracał gołe CMD_PRODUCE)
		// Indeks 7: Budowa Palisady (stare CMD_BUILD_FENCE)
		if actionIndex == 7 {
			cmd.ActionType = cmdBuildStructure
			cmd.InteractionTargetID = uint(buildingPalisade)
		}

	case buildingAcademy:
		// Indeks 4: Wytwarzanie Kusznika
		if actionIndex == 4 {
			cmd.ActionType = cmdProduce
			cmd.ProduceType = unitCrossbowman
		}
		// Indeks 5: Wytwarzanie Dowódcy
		if actionIndex == 5 {
			cmd.ActionType = cmdProduce
			cmd.ProduceType = unitCommander
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
