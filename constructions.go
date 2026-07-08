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

	bld.increaseHP(amount)
}

func (bld *building) build(amount uint16) {
	if !bld.Exists || !bld.IsUnderConstruction {
		return
	}

	bld.increaseHP(amount)
}

// increaseHP dla każdej istniejącej budowli zwiększa PŻ o amount
// Pilnuje, aby bld.HP <= bld.MaxHP; Służy do naprawy budynków
func (bld *building) increaseHP(amount uint16) {
	if !bld.Exists {
		return
	}

	bld.HP += amount
	if bld.HP >= bld.MaxHP {
		bld.HP = bld.MaxHP
	}
}

// Obniżamy bezpośrednio HP budynku. Nie mylić processingBuildingDamage, które
// służa sa sito dla nieskutecznych ataków.
func (bld *building) applyBuildingDamage(finalDamage uint16) {
	// Bez tego bld.HP przekręca się na 65 tys.
	if bld.HP >= finalDamage {
		bld.HP -= finalDamage
	} else {
		bld.HP = 0
	}

	log.Printf("Budynek %s przyjął obrażenia w wysokości: %d. HP: %d/%d",
		buildingDefs[bld.Type].Name, finalDamage, bld.HP, bld.MaxHP)
}

func (bldType buildingType) isRegularBuilding() bool {
	return bldType != buildingRoad && bldType != buildingBridge && bldType != buildingPalisade
}

// Odpowiada za dodanie stworzonej jednostki do mieszkańców budynku.
// Sprawdzenie zostało wykonane w canProduceUnit,
// pojemność za pomocą hasRoom
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

func (bld *building) hasRoom() bool {
	return len(bld.AssignedUnits) < int(bld.MaxFood)
}

// @reminer: Zupełnie nie rozumiem po co obecnie miałbym mieć taką metodę.
// Do czasu aż nie ogarnę units.go zostawię, ale czuję, że jest zbędna.
func (bld *building) getCenter() (uint8, uint8, bool) {
	switch bld.Type {
	case buildingPalisade, buildingBridge:
		// Te rodzaje budynków, zawsze mają dokładnie jeden kafelek
		return bld.OccupiedTiles[0].X, bld.OccupiedTiles[0].X, true
	default:
		// Zwyczajne budowle zawsze są 3na3 więc środek jest z góry znany
		return bld.OccupiedTiles[1].X, bld.OccupiedTiles[1].Y, true
	}
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
// @reminder: funkcja wydaje się całkowicie zbyteczna
// bo teraz budynek ma .OccupiedTiles. Nie usuwam jeszcze.
func (bld *building) getBounds() (int, int, int, int) {
	switch bld.Type {

	case buildingBridge, buildingPalisade:
		return int(bld.OccupiedTiles[0].X), int(bld.OccupiedTiles[0].Y), int(bld.OccupiedTiles[0].X), int(bld.OccupiedTiles[0].Y)
	default:
		return int(bld.OccupiedTiles[0].X), int(bld.OccupiedTiles[0].Y), int(bld.OccupiedTiles[2].X), int(bld.OccupiedTiles[2].Y)
	}
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

// @reminder: jak to powinno działać -- budynek gromadzi obrażenia w bld.AccumnulatedDamage
// i jeśli w danym tiku przekroczyły one próg, to są zadawane „zbiorowo”.
// Jeśli się tego progu nie przekroczyło, to budynek zostaje nienaruszony.
// Tutaj mamy samo zbieranie, wewnątrz updateBuildings jest logika
// „rozliczania tiku”.
func (bld *building) takeDamage(damage uint16) {
	// Od sprawdzenia poprawności rozkazu do jego wykonania
	// budynek mógł już zostać zniszczony, dlatego zostawiam bezpiecznik.
	// W tej chwili nie jestem wstanie udowodnić, że jest zbędny.
	// Jednakże w bState.updateBuildings() najpierw ruszają się jednostki,
	// potem przetwarzamy obrażenia zadane budynkom.
	// Dopiero po tym zaczynamy sprzątać zniszczone budynki i budowy.
	if !bld.Exists {
		return
	}

	bld.AccumulatedDamage += damage
	log.Printf("building %d received %d damage (accumulated: %d)", bld.ID, damage, bld.AccumulatedDamage)
}

// @reminder: jeśli to jest sprawdzenie poprawności rozkazu to powinno mieć miejsce w castle.go
// przed samym przekazaniem go do wykonania.
// Dodatkowo jeśli ta metoda rzeczywiści odpowiada za sprawdzenie możności, to nie powinna
// być metodą budynku i znajdować się w constructions. ponieważ potrzeba informacji o stanie
// danego gracza oraz planszy.
func canProduceUnit(unitType unitType, bld *building, bState *battleState) (bool, uint8) {
	// 1. Czy jest wolne miejsce w budynku?
	// spoko 08.07.2026
	if !bld.hasRoom() {
		return false, produceErrNoRoom
	}

	// Ustalamy kto chce wykonać działanie
	// spoko 08.07.2026
	ownerState := bState.getPlayerState(bld.Owner)

	if ownerState == nil {
		return false, produceErrInvalidOwner
	}

	// 2. Czy nie przekraczamy odgórnego ograniczenia?
	// spoko 08.07.2026
	if ownerState.CurrentPopulation > maxUnitsPerPlayer {
		return false, produceErrPopulationLimit
	}

	// Pobieramy dane o jednostce do stworzenia
	// spoko 08.07.2026
	uStats, ok := unitDefs[unitType]

	if !ok {
		return false, produceErrInvalidType
	}

	// 3. Czy stać gracza?
	// spoko 08.07.2026
	if ownerState.Milk < uStats.Cost {
		return false, produceErrMilk
	}

	// 4. Czy jednostka może wyjść z budynku?
	// nie spoko 08.07.2026
	// zamieniam bld.getClosestWalkableTile, które jest pomieszaniem z poplątaniem
	// na metodę boardData, która zajmuje się jedynie sprawdzeniem sąsiadów budynku
	// oraz lewego dolnego rogu - to jest specjalny warunek, który przeoczyłem wcześniej.
	// Tworzenie budynków też tego nie wspiera, ale naprostuję
	if !bState.Board.hasSpaceAroundBuilding(bld) {
		fmt.Println("hasSpaceAroundBuilding FALSE")

		return false, produceErrNoSpace
	}

	// Udało się, można tworzyć
	fmt.Println("canProduceUnit DAJE TRUE")
	return true, produceErrNone
}

// To nie powinna być metoda budynku tylko bState lub board ponieważ w tej chwili
// budynek zna szczegóły tworzenia jednostek, a nie powinien.
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
func (bState *battleState) produceUnit(newUnitType unitType, bld *building) {
	// Bierzemy dane jednostki, nie sprawdzamy, bo już to zrobiliśmy
	// @reminder: jeśli się przebuduje sygnaturę canProduceUnit to handleProductionComman
	// mogłoby dać już staty i oszczędzić, ale też zaciemnić obraz
	uStats := unitDefs[newUnitType]

	// 3. Ustalamy właściciela
	// mamy bld, a on owner więc bez sensu!
	ownerState := bState.getPlayerState(bld.Owner)
	if ownerState == nil {
		// Wypadałoby coś tutaj zrobić, bo to krytyczny błąd
		return
	}

	// 4. Tworzymy jednostkę
	coords, _ := bState.Board.electSpawnTile(bld)

	bld.spawnUnit(newUnitType, coords.X, coords.Y, bState)

	// 5. Pobieramy mleko za jednostkę
	ownerState.Milk -= uStats.Cost
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

func (bld *building) unassignUnitsfromBuilding(resolver unitResolver) {
	// Trzeba dać znać jednostkom, że nie mają już domu
	for _, uID := range bld.AssignedUnits {
		u, ok := resolver.getUnitByID(uID)

		if ok && u.Exists {
			u.BelongsTo = nil
		}
	}
}
