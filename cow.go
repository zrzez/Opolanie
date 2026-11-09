package main

import (
	"fmt"
	"log"
	"math"
)

// cow.go

// Konfiguracja zachowania krowy.
const (
	grazingRadius         = 12  // Promień pastwiska od obory
	eatingDelay           = 12  // Czas trwania jedzenia jednej kępki
	fullUdderAmount       = 100 // Pojemność wymion
	milkingSpeed    uint8 = 3   // celuję w 1,6 sekundy dojenia
)

// @todo: @reminder: zwinę to później w MSS dla wszystkich jednostek, bo tutaj nie bardzo
//   jest sens trzymać osobno. O ile dobrze rozumiem - 29.10.2026
func (u *unit) handleCowBehavior(pathfindingBudget int, bState *battleState) {
	// 0. Obsługa techniczna, opóźnienia i blokady
	if u.Delay > 0 {
		u.Delay--

		return
	}

	// 1. Ucieczka
	if u.Command == cmdUFlee {
		if u.isAtTarget() {
			u.setIdle()
		} else {
			u.move(pathfindingBudget, bState)
		}

		return
	}

	// 2. Tryb: "Stać bezmyślnie", rozkaz od gracza
	// @todo: muszę to wypróbować i zobaczyć, jak to wyłączyć, krowa nie może w nieskończoność stać bezczynnie
	if u.Command == cmdUStop {
		// @reminder: jedyne użycie idleCow. Pewnie do zastąpienia zwykłym idle
		u.idleCow()

		return
	}

	// 3. Faza: Obora, pełne wymiona lub ucieczka
	if u.Udder >= fullUdderAmount || u.Command == cmdBMilking {
		u.milkCowPhase(pathfindingBudget, bState)

		return
	}

	// 4. Faza: Wypasanie
	// Obejmuje też tryb CMD_GRAZE i CMD_IDLE (szukanie jedzenia) to może być problematyczne dla 1.
	u.grazeCowPhase(pathfindingBudget, bState)
}

// Fazy zachowania

// Tryb 0: Krowa stoi i nic nie robi.
func (u *unit) idleCow() {
	if u.State != stateIdle {
		u.setIdle()
	}
}

// Faza 1: Logika związana z oborą.
func (u *unit) milkCowPhase(pathfindingBudget int, bState *battleState) {
	barnX, barnY, isSpotFree := findNearestBarnMilkingSpot(u, bState)

	// A. Nie znaleziono żadnej obory
	if !isSpotFree {
		u.setIdleWithReason("brak obory")

		return
	}

	// B. Krowa jest na miejscu dojenia
	if u.X == barnX && u.Y == barnY {
		u.performMilkingAction(pathfindingBudget, bState)

		return
	}

	// C. Krowa musi dojść do obory (lub do kolejki)
	u.moveToBarnOrQueue(pathfindingBudget, bState, barnX, barnY, isSpotFree)
}

// Faza 2 i 3: Logika wypasania
func (u *unit) grazeCowPhase(pathfindingBudget int, bState *battleState) {
	// Krok A: Jeśli stoimy na trawie → jedz - 12.09.2026 sprawdzone i posprzątane
	if u.tryEatGrass(bState) {
		return
	}

	// @reminder: wygląda jak coś co powinno być ogarnięte w MSS jak u innych jednostek
	// Jeśli mamy rozkaz ruchu (np. gracz kliknął), a nie doszliśmy → idź
	if u.Command == cmdUMove && !u.isAtTarget() {
		u.move(pathfindingBudget, bState)

		return
	}

	// @reminder: wygląda jak coś co powinno być ogarnięte w MSS jak u innych jednostek
	// Krok B: Jeśli nie jemy, a jesteśmy w ruchu → idź dalej
	if u.State == stateMoving && !u.isAtTarget() {
		u.move(pathfindingBudget, bState)

		return
	}

	// Krok C: Znajdź nowy kawałek trawy - @todo: ogarnij
	u.findNewPasture(pathfindingBudget, bState)
}

const (
	milkGainHealthyGrass = 25
	milkGainGrazedGrass  = 13
)

// Odpowiada za obsłużenie wszystkich przypadków przy wypasaniu krowy.
func (u *unit) tryEatGrass(bState *battleState) bool {
	currentTile := &bState.Board.Tiles[u.X][u.Y]

	// Jeśli nie trawa, to wychodzimy
	if !isGrass(currentTile.TextureID) {
		return false
	}

	var milkGain uint8

	switch {
	// Pełna trawa
	case !currentTile.IsGrazed:
		currentTile.IsGrazed = true
		currentTile.GrazedOverlayID = uint8(spriteGrassGrazed)
		milkGain = milkGainHealthyGrass
	// Nadgryziona trawa
	case uint16(currentTile.GrazedOverlayID) == spriteGrassGrazed:
		currentTile.GrazedOverlayID = uint8(spriteGrassStubbed)
		milkGain = milkGainGrazedGrass
	// Wyżarta w pełni lub coś innego czego nie przewidziałem
	default:
		fmt.Println("tryEatGrass(), przełącznik, przypadek domyślny zwraca false")

		return false
	}

	u.State = stateGrazing
	u.setAnimationType()
	u.Udder = min(u.Udder+milkGain, fullUdderAmount)
	u.Delay = eatingDelay

	return true
}

// Szuka nowej trawy i wydaje rozkaz ruchu.
func (u *unit) findNewPasture(pathfindingBudget int, bState *battleState) {
	/*originX, originY, ok := u.getGrazingAnchorPoint()
	if !ok {
		return
	}

	// Użycie deterministycznego algorytmu BFS
	// Szuka najbliższej trawy, do której DA SIĘ dojść

	grassX, grassY, found := findReachableGrass(u, bState, originX, originY, grazingRadius)*/

	grassTile, found := findReachableGrass(u, bState)
	fmt.Printf("Znalazłem trawę %d,%d", grassTile.X, grassTile.Y)

	if found {
		u.addAndMove(pathfindingBudget, cmdUGraze, grassTile.X, grassTile.Y, 0, bState, "Znalazłam trawę.")
	} else {
		// Nie ma trawy -> Wróć pod oborę i czekaj
		u.returnToBarnArea(pathfindingBudget, bState)
	}
}

// Decyzja: idź prosto do wejścia CZY dołącz do kolejki
func (u *unit) moveToBarnOrQueue(pathfindingBudget int, bState *battleState, barnX, barnY uint8, spotAvailable bool) {
	if spotAvailable {
		// Obora wolna -> Idź prosto do środka
		u.addAndMove(pathfindingBudget, cmdUMove, barnX, barnY, 0, bState, "Idę do obory (wolne).")
		u.IsInQueue = false
	} else {
		// Zajęte -> Dołącz do kolejki
		u.joinMilkingQueue(pathfindingBudget, bState, barnX, barnY)
	}
}

// Logika dołączania do kolejki (Oryginalna funkcjonalność).
func (u *unit) joinMilkingQueue(pathfindingBudget int, bState *battleState, barnX, barnY uint8) {
	homeBarn := u.BelongsTo

	if homeBarn == nil {
		// Jeśli krowa nie ma domu, a chce się doić, idzie pod najbliższą (awaryjnie)
		// Ale nie może wejść do kolejki "bezdomnej".
		// w takim wypadku po prostu idzie w pobliże znalezionego barnX/y.
		u.addAndMove(pathfindingBudget, cmdUMove, barnX, barnY, 0, bState, "Idę pod oborę (bezdomna).")

		return
	}

	// 1. Dodaj się do listy w budynku (jeśli jeszcze nie ma)
	isAlreadyInQueue := false

	for _, cowID := range homeBarn.MilkingQueue {
		if cowID == u.ID {
			isAlreadyInQueue = true

			break
		}
	}

	if !isAlreadyInQueue {
		homeBarn.MilkingQueue = append(homeBarn.MilkingQueue, u.ID)
		u.IsInQueue = true
		log.Printf("unit %d: Dodałem się do kolejki obory %d.", u.ID, homeBarn.ID)
	}

	u.State = stateWaiting

	// 2. Znajdź miejsce fizyczne w poczekalni (wokół obory)
	waitX, waitY, foundWaitingSpot := findNearestAvailableWaitingSpot(barnX, barnY, bState, u.ID)

	if !foundWaitingSpot {
		u.setIdleWithReason("brak miejsca w poczekalni")
		// Mimo braku miejsca fizycznego, logicznie jest w kolejce
		u.IsInQueue = true

		return
	}

	// 3. Sprawdź czy już tam stoi
	distToWaitingSpot := int(math.Abs(float64(u.X-waitX)) + math.Abs(float64(u.Y-waitY)))
	isAtWaitingSpot := u.X == waitX && u.Y == waitY

	if isAtWaitingSpot || (distToWaitingSpot <= 1 && u.Command == cmdUIdle) {
		// Stoi grzecznie i czeka
		u.State = stateIdle
		u.IsInQueue = true

		return
	}

	// 4. Idź do miejsca w kolejce
	hasCorrectMoveCommand := u.Command == cmdUMove && u.Target.Position.X == waitX && u.Target.Position.Y == waitY
	if hasCorrectMoveCommand {
		u.move(pathfindingBudget, bState)
	} else {
		u.addAndMove(pathfindingBudget, cmdUMove, waitX, waitY, 0, bState, "Idę do poczekalni.")
	}
}

// =============================================================================
// Pomocnicze (Helpers) - Oryginalne i Nowe
// =============================================================================

// Zwraca punkt centralny, wokół którego krowa ma się kręcić (Pastwisko)
func (u *unit) getGrazingAnchorPoint() (uint8, uint8, bool) {
	if u.BelongsTo != nil && u.BelongsTo.Exists {
		return u.BelongsTo.getCenter()
	}
	// Krowa bezdomna trzyma się miejsca, gdzie stoi (lub gdzie ostatnio jadła)
	return u.X, u.Y, false
}

// Awaryjny powrót w pobliże obory, gdy nie ma trawy
func (u *unit) returnToBarnArea(pathfindingBudget int, bState *battleState) {
	u.setIdleWithReason("brak trawy w zasięgu")
	if u.BelongsTo != nil {

		// bx, by, ok := u.BelongsTo.getClosestWalkableTile(bState)
		// TUTAJ WPROWADZAM ZMIANĘ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓
		coords := bState.Board.bldNeighborCoords(u.BelongsTo)
		var bestX, bestY uint8
		minPathLen := math.MaxInt32
		var found bool

		for _, coord := range coords {
			if bState.Board.Tiles[coord.X][coord.Y].IsWalkable && bState.Board.Tiles[coord.X][coord.Y].Unit == nil {
				tempPath := findPath(bState.Board, u, u.X, u.Y, coord.X, coord.Y, u.Path)

				if tempPath != nil && len(tempPath) < minPathLen {
					minPathLen = len(tempPath)
					bestX, bestY = coord.X, coord.Y
					found = true
				}
			}
		}
		if found {
			u.addAndMove(pathfindingBudget, cmdUMove, bestX, bestY, 0, bState, "Wracam pod oborę (brak paszy).")
		}
	}
}

// Pomocnik do wykonania ruchu
func (u *unit) addAndMove(pathfindingBudget int, cmdType commandType, x, y uint8, id objectID, bState *battleState, logMsg string) {
	if u.Command != cmdType || u.Target.Position.X != x || u.Target.Position.Y != y {
		if logMsg != "" {
			log.Printf("unit %d (COW): %s", u.ID, logMsg)
		}

		cmd := &command{
			ActionType: cmdType,
			Target:     targetReference{Kind: targetTile, Position: point{X: x, Y: y}},
		}
		u.addUnitCommand(cmd, bState)
	}

	u.move(pathfindingBudget, bState)
}

// Przekazuje mleko do gracza.
func (u *unit) performMilking(bState *battleState) uint16 {
	if u.Udder == 0 {
		return 0
	}

	ownerState := bState.HumanPlayerState
	if u.Owner == bState.AIPlayerID {
		ownerState = bState.AIEnemyState
	}

	// 1. Ustal porcję przekazu
	amountToTransfer := uint16(min(u.Udder, milkingSpeed))

	// 2. Sprawdź miejsce w magazynie
	spaceInBarn := max(ownerState.MaxMilk-ownerState.Milk, 0)

	// 3. Ustal ile można przenieść za jednym razem
	amountToTransfer = min(amountToTransfer, spaceInBarn)

	// 4. Wykonaj transfer
	ownerState.Milk += amountToTransfer
	u.Udder -= uint8(amountToTransfer)

	return amountToTransfer
}

// Wykonuje czynność oddania mleka (gdy stoi w punkcie dojenia).
func (u *unit) performMilkingAction(pathfindingBudget int, bState *battleState) {
	transferredMilk := u.performMilking(bState)

	u.IsInQueue = false

	if transferredMilk > 0 {
		u.State = stateMilking
		// Logujemy postęp co jakiś czas lub przy zmianie
		log.Printf("Przed dojeniem MILK %d, MAX_MILK %d", bState.HumanPlayerState.Milk, bState.HumanPlayerState.MaxMilk)
		log.Printf("Krowa %d: oddaje mleko... (%d/%d)", u.ID, u.Udder, fullUdderAmount)
		log.Printf("Po dojeniu MILK %d, MAX_MILK %d", bState.HumanPlayerState.Milk, bState.HumanPlayerState.MaxMilk)

		if u.Udder <= 0 {
			// Pusta → Koniec dojenia → Idź na pastwisko
			// To nadpisze CMD_FLEE na CMD_GRAZE/CMD_MOVE
			log.Printf("Krowa %d nie ma już mleka %d. Idzie się paść.", u.ID, u.Udder)
			u.grazeCowPhase(pathfindingBudget, bState)
		} else {
			// Nadal ma mleko → MUSI ZOSTAĆ w OBORZE
			// Ustawiamy CMD_FLEE, aby handleCowBehavior w następnej klatce
			// skierowało nas ponownie do milkCowPhase zamiast na pastwisko.
			log.Printf("Krowa %d wciąż ma mleko %d. Zostaje na dalsze dojenie.", u.ID, u.Udder)
			u.Command = cmdBMilking
		}
	} else {
		if u.Udder > 0 {
			// Ma mleko, ale nie może oddać → Czekaj w trybie "Obora"
			u.Command = cmdBMilking
			log.Printf("Krowa %d ma jeszcze mleko %d, ale nie ma miejsca na mleko %d", u.ID, u.Udder, bState.HumanPlayerState.Milk)
		} else {
			// Pusta → Pastwisko
			u.grazeCowPhase(pathfindingBudget, bState)
		}
	}
}

// Szuka najbliższej obory i sprawdza czy wejście jest wolne
func findNearestBarnMilkingSpot(u *unit, bState *battleState) (uint8, uint8, bool) {
	// 1. Sprawdź przypisaną oborę (Dom)
	if u.BelongsTo != nil {
		homeBarn := u.BelongsTo
		if homeBarn.Exists && homeBarn.Type == buildingBarn {
			x, y := homeBarn.getMilkingSpotCoords()

			// Sprawdzamy zajętość w Tiles
			currentTile := &bState.Board.Tiles[x][y]
			occupant := currentTile.Unit
			// Wolne jeśli nikogo nie ma LUB stoi tam ta sama krowa
			isFree := occupant == nil || occupant.ID == u.ID
			return x, y, isFree
		}
	}
	// 2. Jeśli nie ma domu (lub zniszczony), szukaj globalnie
	return findClosestBarnGlobally(u, bState)
}

// Szuka dowolnej obory gracza
func findClosestBarnGlobally(u *unit, bState *battleState) (uint8, uint8, bool) {
	// @todo: współrzędne powinny być uint, jeżeli nie ma miejsca, to powinienem jakoś inaczej to ogarnąć
	closestX, closestY := uint8(0), uint8(0)
	isSpotFree := false
	minDist := math.MaxInt32

	for _, bld := range bState.Buildings {
		if bld == nil || !bld.Exists || bld.Type != buildingBarn || bld.Owner != u.Owner || bld.IsUnderConstruction {
			continue
		}
		x, y := bld.getMilkingSpotCoords()
		dist := int(math.Abs(float64(u.X-x)) + math.Abs(float64(u.Y-y)))

		if dist < minDist {
			minDist = dist
			closestX, closestY = x, y

			currentTile := &bState.Board.Tiles[x][y]
			occupant := currentTile.Unit
			isSpotFree = occupant == nil || occupant.ID == u.ID
		}
	}
	return closestX, closestY, isSpotFree
}

// Zwraca koordynaty wejścia do obory (Naprawione, aby wskazywały punkt PRZED, jeśli potrzeba,
// ale przy nowym odnajdywaniu drogi wchodzimy do środka, więc używamy logiki minX/maxY).
func (bld *building) getMilkingSpotCoords() (uint8, uint8) {
	minX, maxY := uint8(math.MaxUint8), uint8(0)
	for _, currentTile := range bld.OccupiedTiles {
		if currentTile.X < minX {
			minX = currentTile.X
		}

		if currentTile.Y > maxY {
			maxY = currentTile.Y
		}
	}
	// Zwracamy lewy dolny róg samej konstrukcji.
	// Dzięki isWalkableUnit w pathfinding.go krowa tam wejdzie.
	return minX, maxY
}

// Szuka wolnego miejsca w pobliżu obory (spiralnie), żeby krowy nie stały na sobie
func findNearestAvailableWaitingSpot(targetX, targetY uint8, bState *battleState, myUnitID unitID) (uint8, uint8, bool) {
	// Radius w pętli niech będzie int, łatwiej się operuje
	for radius := 1; radius <= 5; radius++ {
		for dx := -radius; dx <= radius; dx++ {
			for dy := -radius; dy <= radius; dy++ {

				// Twoja optymalizacja (tylko ramka) - bez zmian, ale używamy int-math zamiast float
				if abs(dx) != radius && abs(dy) != radius {
					continue
				}

				// KLUCZOWE: Rzutujemy na int przed dodaniem ujemnego dx/dy
				calcX := int(targetX) + dx
				calcY := int(targetY) + dy

				// Sprawdzamy granice na intach (zabezpieczenie przed ujemnymi i wyjściem poza planszę)
				if calcX < 0 || calcX >= int(boardMaxX) || calcY < 0 || calcY >= int(boardMaxY) {
					continue
				}

				// Teraz bezpiecznie wracamy na uint8
				x, y := uint8(calcX), uint8(calcY)

				tile := &bState.Board.Tiles[x][y]

				// Reszta logiki bez zmian...
				isOccupied := tile.Unit != nil && tile.Unit.ID != myUnitID
				if isWalkable(bState.Board, x, y) && !isOccupied && tile.Building == nil {
					return x, y, true
				}
			}
		}
	}

	return 0, 0, false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

// Pierścien wewnętrzny. Bezpośrednio sąsiadujące z krową.
var grazingRing1 = [8][2]int8{
	{-1, -1},
	{0, -1},
	{1, -1},
	{-1, 0},
	{1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
}

// Pierścień zewnętrzny do wypasu krowy.
var grazingRing2 = [16][2]int8{
	{-2, -2},
	{-1, -2},
	{0, -2},
	{1, -2},
	{2, -2},
	{-2, -1},
	{2, -1},
	{-2, 0},
	{2, 0},
	{-2, 1},
	{2, 1},
	{-2, 2},
	{-1, 2},
	{0, 2},
	{1, 2},
	{2, 2},
}

// @reminder: krowa błądzi bez względu na odległość do trawy, z pewnością się to zemści przy rzecze
func findReachableGrass(u *unit, bState *battleState) (point, bool) {
	// 1. Sprawdzamy pierścień wewnętrzny
	full1, grazed1, fallback1 := scanRing(u, bState, grazingRing1[:])

	// Jeśli tuż obok krowy jest soczysta trawa, to przerywamy poszukiwania, zwracamy losowy kafelek.
	if len(full1) > 0 {
		return full1[rng.Intn(len(full1))], true
	}

	// Jeżeli nie było w nienaruszonej trawy, to zwracamy już lekko wyżartą
	if len(grazed1) > 0 {
		// @todo: później się zastanów, czy losowa wyżarta to dobre rozwiązanie. Pomijamy odległość!
		//  w sumie ta z wewnętrznego powinna mieć pierwszeństwo
		return grazed1[rng.Intn(len(grazed1))], true
	}

	// 2. Sprawdzamy pierścień zewnętrzny
	full2, grazed2, fallback2 := scanRing(u, bState, grazingRing2[:])

	// Jeśli jest soczysta trawa, to przerywamy poszukiwania, zwracamy losowy kafelek.
	if len(full2) > 0 {
		return full2[rng.Intn(len(full2))], true
	}

	// Jeżeli nie było w obu pierścieniach nienaruszonej trawy, to zwracamy już lekko wyżartą
	if len(grazed2) > 0 {
		// @todo: później się zastanów, czy losowa wyżarta to dobre rozwiązanie. Pomijamy odległość!
		//  w sumie ta z wewnętrznego powinna mieć pierwszeństwo
		return grazed2[rng.Intn(len(grazed2))], true
	}

	// Przy braku jakiejkolwiek trawy do wypasania, idziemy w losowe wolne pole w zasięgu.
	// W zamyśle to umożliwi przeszukanie innego obszaru do wypasania krów.
	// @todo: sprawdź, czy działa jak myślę
	allFallbacks := append(fallback1, fallback2...)
	if len(allFallbacks) > 0 {
		return allFallbacks[rng.Intn(len(allFallbacks))], true
	}

	return point{}, false
}

func scanRing(u *unit, bState *battleState, ring [][2]int8) (full, grazed, fallback []point) {
	for _, offset := range ring {
		checkX := int8(u.X) + offset[0]
		checkY := int8(u.Y) + offset[1]

		// Jeśli nie można stanąć na danym kafelku to pomijamy.
		if !bState.Board.isValidWalkableTile(checkX, checkY) || bState.Board.Tiles[checkX][checkY].IsBurning {
			continue
		}

		currentPoint := point{X: uint8(checkX), Y: uint8(checkY)}
		// Dodajemy do listy awaryjnej
		fallback = append(fallback, currentPoint)

		currentTile := &bState.Board.Tiles[checkX][checkY]

		if isGrass(currentTile.TextureID) {
			if !currentTile.IsGrazed {
				full = append(full, currentPoint)
			} else if uint16(currentTile.GrazedOverlayID) == spriteGrassGrazed {
				grazed = append(grazed, currentPoint)
			}
		}
	}

	return full, grazed, fallback
}

// @reminder: używam tej funkcji w miejscu, które już gwarantuje, że u.Exists.
func (bState *battleState) handleCowFlee(u *unit) {
	if !u.wasAttacked || u.Type != unitCow {
		return
	}

	u.wasAttacked = false

	fleeCoords, found := bState.findCowFleeDestination(u)
	if !found {
		u.fleeBurningTile(bState)
		u.Command = cmdUFlee
		return
	}

	cmd := &command{
		ActionType: cmdUFlee,
		Target:     targetReference{Kind: targetBuilding, Position: fleeCoords},
	}

	u.addUnitCommand(cmd, bState)

	log.Printf("Krowa otrzymała obrażenia, ucieka do obory na (%d,%d).", fleeCoords.X, fleeCoords.Y)
}

func (bState *battleState) findCowFleeDestination(u *unit) (point, bool) {
	// 1. Sprawdzenie przynależnej obory
	if u.BelongsTo != nil && u.BelongsTo.Exists {
		if fleeCoords, ok := bState.Board.electSpawnTile(u.BelongsTo); ok {
			return fleeCoords, true
		}
	}

	// 2. Szukanie innego budynku lub budowli
	return bState.findFallbackFleeSpot(u)
}

func (bState *battleState) findFallbackFleeSpot(u *unit) (point, bool) {
	var bestCoords point
	priority := 0 // 0 - brak, 1 - zwykły budynek, 2 - budynek główny

	for _, bld := range bState.Buildings {
		if bld.Owner != u.Owner || !bld.Exists || bld.Type == buildingPalisade || bld.Type == buildingBridge {
			continue
		}

		fleeCoords, ok := bState.Board.electSpawnTile(bld)
		if !ok {
			continue
		}

		switch bld.Type {
		case buildingBarn:
			// Byle obora kończy poszukiwania
			return fleeCoords, true

		case buildingMain:
			if priority < 2 {
				bestCoords = fleeCoords
				priority = 2
			}

		default:
			// Dowolny inny budynek lub budowla
			if priority < 1 {
				bestCoords = fleeCoords
				priority = 1
			}
		}
	}

	return bestCoords, priority > 0
}
