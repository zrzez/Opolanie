package main

import (
	"log"
	"math"
)

// cow.go

// Konfiguracja zachowania krowy.
const (
	grazingRadius         = 12  // Promień pastwiska od obory
	eatingDelay           = 12  // Czas trwania jedzenia jednej kępki
	fullUdderAmount       = 100 // Pojemność wymienia
	milkingSpeed    uint8 = 3   // celuję w 1,6 sekundy dojenia
)

func (u *unit) handleCowBehavior(resolver objectResolver, board *boardData, pathfindingBudget *int, bState *battleState) {
	// 0. Obsługa techniczna, opóźnienia i blokady
	if u.Delay > 0 {
		u.Delay--

		return
	}

	if u.BlockedCounter > 0 {
		return
	}

	// 0. Ucieczka
	if u.Command == cmdUFlee {
		u.handleCowFlee(bState)
	}

	// 1. Tryb: "Stać bezmyślnie", rozkaz od gracza
	// @todo: muszę to wypróbować i zobaczyć, jak to wyłączyć, krowa nie może w nieskończoność stać bezczynnie
	if u.Command == cmdUStop {
		u.idleCow()

		return
	}

	// 2. Faza: Obora, pełne wymiona lub ucieczka
	if u.Udder >= fullUdderAmount || u.Command == cmdBMilking {
		u.milkCowPhase(resolver, board, pathfindingBudget, bState)

		return
	}

	// 3. Faza: Wypasanie
	// Obejmuje też tryb CMD_GRAZE i CMD_IDLE (szukanie jedzenia) to może być problematyczne dla 1.
	u.grazeCowPhase(resolver, board, pathfindingBudget, bState)
}

// Fazy zachowania

// Tryb 0: Krowa stoi i nic nie robi.
func (u *unit) idleCow() {
	if u.State != stateIdle {
		u.setIdle()
	}
}

// Faza 1: Logika związana z oborą.
func (u *unit) milkCowPhase(resolver objectResolver, board *boardData, pathfindingBudget *int, bState *battleState) {
	barnX, barnY, isSpotFree := findNearestBarnMilkingSpot(u, bState)

	// A. Nie znaleziono żadnej obory
	if !isSpotFree {
		u.setIdleWithReason("brak obory")

		return
	}

	// B. Krowa jest na miejscu dojenia
	if u.X == barnX && u.Y == barnY {
		u.performMilkingAction(resolver, board, pathfindingBudget, bState)

		return
	}

	// C. Krowa musi dojść do obory (lub do kolejki)
	u.moveToBarnOrQueue(resolver, board, pathfindingBudget, bState, barnX, barnY, isSpotFree)
}

// Faza 2 i 3: Logika wypasania
func (u *unit) grazeCowPhase(resolver objectResolver, board *boardData, pathfindingBudget *int, bState *battleState) {
	// Krok A: Jeśli stoimy na trawie → jedz
	if u.tryEatGrass(bState) {
		return
	}

	// Jeśli mamy rozkaz ruchu (np. gracz kliknął), a nie doszliśmy → idź
	if u.Command == cmdUMove && !u.isAtTarget() {
		u.move(resolver, board, pathfindingBudget, bState)

		return
	}

	// Krok B: Jeśli nie jemy, a jesteśmy w ruchu → idź dalej
	if u.State == stateMoving && !u.isAtTarget() {
		u.move(resolver, board, pathfindingBudget, bState)

		return
	}

	// Krok C: Znajdź nowy kawałek trawy
	u.findNewPasture(resolver, board, pathfindingBudget, bState)
}

var (
	isHealthyGrass bool
	isGrazedGrass  bool
)

func (u *unit) tryEatGrass(bState *battleState) bool {
	currentTile := &bState.Board.Tiles[u.X][u.Y]
	texID := currentTile.TextureID
	isHealthyGrass = isGrass(texID) && !currentTile.IsGrazed
	isGrazedGrass = isGrass(texID) && currentTile.IsGrazed && uint16(currentTile.GrazedOverlayID) != spriteGrassStubbed

	if isHealthyGrass {
		// Krok 1: Cała → nadgryziona
		currentTile.IsGrazed = true
		currentTile.GrazedOverlayID = uint8(spriteGrassGrazed)

		u.State = stateGrazing
		u.AnimationType = "idle"
		u.Udder += 25

		if u.Udder > fullUdderAmount {
			u.Udder = fullUdderAmount
		}

		u.Delay = eatingDelay

		return true
	}

	if isGrazedGrass {
		// Krok 2: Nadgryziona → wyżarta
		currentTile.GrazedOverlayID = uint8(spriteGrassStubbed)

		u.State = stateGrazing
		u.AnimationType = "idle"
		u.Udder += 13

		if u.Udder > fullUdderAmount {
			u.Udder = fullUdderAmount
		}

		u.Delay = eatingDelay

		return true
	}

	// Jeśli to ściernisko (STUBBED) lub inny teren → nie da się zjeść
	return false
}

// Szuka nowej trawy i wydaje rozkaz ruchu.
func (u *unit) findNewPasture(resolver objectResolver, board *boardData, pathfindingBudget *int, bState *battleState) {
	originX, originY, ok := u.getGrazingAnchorPoint()
	if !ok {
		return
	}

	// Użycie deterministycznego algorytmu BFS
	// Szuka najbliższej trawy, do której DA SIĘ dojść

	grassX, grassY, found := findReachableGrass(u, bState, originX, originY, grazingRadius)

	if found {
		u.addAndMove(resolver, board, pathfindingBudget, cmdUGraze, grassX, grassY, 0, bState, "Znalazłam trawę.")
	} else {
		// Nie ma trawy -> Wróć pod oborę i czekaj
		u.returnToBarnArea(resolver, board, pathfindingBudget, bState)
	}
}

// Decyzja: idź prosto do wejścia CZY dołącz do kolejki
func (u *unit) moveToBarnOrQueue(resolver objectResolver, board *boardData, pathfindingBudget *int, bState *battleState, barnX, barnY uint8, spotAvailable bool) {
	if spotAvailable {
		// Obora wolna -> Idź prosto do środka
		u.addAndMove(resolver, board, pathfindingBudget, cmdUMove, barnX, barnY, 0, bState, "Idę do obory (wolne).")
		u.IsInQueue = false
	} else {
		// Zajęte -> Dołącz do kolejki
		u.joinMilkingQueue(resolver, board, pathfindingBudget, bState, barnX, barnY)
	}
}

// Logika dołączania do kolejki (Oryginalna funkcjonalność).
func (u *unit) joinMilkingQueue(resolver objectResolver, board *boardData, pathfindingBudget *int, bState *battleState, barnX, barnY uint8) {
	homeBarn := u.BelongsTo

	if homeBarn == nil {
		// Jeśli krowa nie ma domu, a chce się doić, idzie pod najbliższą (awaryjnie)
		// Ale nie może wejść do kolejki "bezdomnej".
		// w takim wypadku po prostu idzie w pobliże znalezionego barnX/y.
		u.addAndMove(resolver, board, pathfindingBudget, cmdUMove, barnX, barnY, 0, bState, "Idę pod oborę (bezdomna).")

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

	if isAtWaitingSpot || (distToWaitingSpot <= 1 && (u.BlockedCounter > 0 || u.Command == cmdUIdle)) {
		// Stoi grzecznie i czeka
		u.State = stateIdle
		u.IsInQueue = true

		return
	}

	// 4. Idź do miejsca w kolejce
	hasCorrectMoveCommand := u.Command == cmdUMove && u.TargetX == waitX && u.TargetY == waitY
	if hasCorrectMoveCommand {
		u.move(resolver, board, pathfindingBudget, bState)
	} else {
		u.addAndMove(resolver, board, pathfindingBudget, cmdUMove, waitX, waitY, 0, bState, "Idę do poczekalni.")
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
func (u *unit) returnToBarnArea(resolver objectResolver, board *boardData, pathfindingBudget *int, bState *battleState) {
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
				tempPath := findPath(bState.Board, u, u.X, u.Y, coord.X, coord.Y)

				if tempPath != nil && len(tempPath) < minPathLen {
					minPathLen = len(tempPath)
					bestX, bestY = coord.X, coord.Y
					found = true
				}
			}
		}
		if found {
			u.addAndMove(resolver, board, pathfindingBudget, cmdUMove, bestX, bestY, 0, bState, "Wracam pod oborę (brak paszy).")
		}
	}
}

// Pomocnik do wykonania ruchu
func (u *unit) addAndMove(resolver objectResolver, board *boardData, pathfindingBudget *int, cmdType commandType, x, y uint8, id ObjectID, bState *battleState, logMsg string) {
	if u.Command != cmdType || u.TargetX != x || u.TargetY != y {
		if logMsg != "" {
			log.Printf("unit %d (COW): %s", u.ID, logMsg)
		}

		cmd := &command{
			ActionType:          cmdType,
			TargetX:             x,
			TargetY:             y,
			InteractionTargetID: id,
		}
		u.addUnitCommand(cmd, bState)
	}

	u.move(resolver, board, pathfindingBudget, bState)
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
func (u *unit) performMilkingAction(resolver objectResolver, board *boardData, pathfindingBudget *int, bState *battleState) {
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
			u.grazeCowPhase(resolver, board, pathfindingBudget, bState)
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
			u.grazeCowPhase(resolver, board, pathfindingBudget, bState)
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
func findNearestAvailableWaitingSpot(targetX, targetY uint8, bState *battleState, myUnitID UnitID) (uint8, uint8, bool) {
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

func findReachableGrass(u *unit, bState *battleState, originX, originY, radius uint8) (uint8, uint8, bool) {
	// Zakładam, że u.X i u.Y są uint8
	queue := []point{{u.X, u.Y}}

	// Hashmapa odwiedzonych
	visited := make(map[int]bool) // Klucz int jest bezpieczniejszy dla hasha
	visited[int(u.Y)*int(boardMaxX)+int(u.X)] = true

	const maxSearchSteps = 400
	steps := 0

	// Definicja kierunków musi być na intach!
	// Nie używaj tu struktury `point` jeśli ona ma pola uint.
	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	for len(queue) > 0 {
		steps++
		if steps > maxSearchSteps {
			break
		}

		curr := queue[0]
		queue = queue[1:]

		for _, d := range dirs {
			// KLUCZOWE: Rzutowanie na int przed dodaniem offsetu
			calcX := int(curr.X) + d[0]
			calcY := int(curr.Y) + d[1]

			// Sprawdzenie granic na intach
			if calcX < 0 || calcX >= int(boardMaxX) || calcY < 0 || calcY >= int(boardMaxY) {
				continue
			}

			// Bezpieczny powrót na uint8
			nextX, nextY := uint8(calcX), uint8(calcY)

			hash := int(nextY)*int(boardMaxX) + int(nextX)
			if visited[hash] {
				continue
			}
			visited[hash] = true

			// 1. Sprawdź zasięg od "Domu".
			// Musimy rzutować na int LUB float, żeby odjąć bezpiecznie.
			// math.Abs wymaga float64, więc najpierw robimy diff na intach.
			distX := abs(int(nextX) - int(originX))
			distY := abs(int(nextY) - int(originY))

			if distX+distY > int(radius) {
				continue
			}

			// 2. Sprawdź cel (Trawa)
			currentTile := &bState.Board.Tiles[nextX][nextY]
			texID := currentTile.TextureID

			if isGrass(texID) && !currentTile.IsGrazed && currentTile.Unit == nil && currentTile.Building == nil {
				return nextX, nextY, true
			}

			// 3. Sprawdź drogę
			if isWalkable(bState.Board, nextX, nextY) {
				queue = append(queue, point{nextX, nextY})
			}
		}
	}

	return 0, 0, false
}

func (u *unit) handleCowFlee(bState *battleState) {
	if u.Type == unitCow && u.Exists {
		if u.Udder < fullUdderAmount && u.Command != cmdUFlee {
			barnX, barnY, foundBarn := findNearestBarnMilkingSpot(u, bState)
			if foundBarn {
				cmd := &command{
					ActionType:          cmdUFlee,
					TargetX:             barnX,
					TargetY:             barnY,
					InteractionTargetID: 0,
				}
				u.addUnitCommand(cmd, bState)
				log.Printf("unit %d (COW): Otrzymała obrażenia, ucieka do obory na (%d,%d).", u.ID, barnX, barnY)
			} else {
				log.Printf("unit %d (COW): Otrzymała obrażenia, ale nie znalazła obory do ucieczki. ", u.ID)
			}
		}
	}
}
