package main

import (
	"fmt"
	"log"
	"math"
	"slices"
)

// registerBuilding zapisuje na planszy, które kafelki są zajmowane przez budynek.
// startX, startY to lewy górny róg budynku.
func (board *boardData) registerBuilding(bld *building, startX, startY uint8) error {
	stats, ok := buildingDefs[bld.Type]
	if !ok {
		return fmt.Errorf("nieznany rodzaj budynku %d", bld.Type)
	}

	bld.OccupiedTiles = make([]point, 0, stats.Width*stats.Height)

	for offsetX := range stats.Width {
		for offsetY := range stats.Height {
			finalX, finalY := startX+offsetX, startY+offsetY
			currentTile := &board.Tiles[finalX][finalY]
			bld.OccupiedTiles = append(bld.OccupiedTiles, point{X: finalX, Y: finalY})

			currentTile.Building = bld
			currentTile.IsWalkable = false
		}
	}

	return nil
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

// applyFinishedGraphics nakłada tekstury na kafelki zajmowane przez ukończone budynki.
func (board *boardData) applyFinishedGraphics(bld *building) {
	switch bld.Type {
	case buildingPalisade:
		occupiedTile := bld.OccupiedTiles[0]
		// Palisady muszą tworzyc spójną sieć, dlatego trzeba odświeżyć też sąsiednie kafelki
		joinPalisade(occupiedTile.X, occupiedTile.Y, board)

		return

	case buildingBridge:
		occupiedTile := bld.OccupiedTiles[0]
		board.Tiles[occupiedTile.X][occupiedTile.Y].IsWalkable = true
		// Z tego, co rozumiem to mosty budowane przez graczy zawsze mają ten sam wygląd
		board.Tiles[occupiedTile.X][occupiedTile.Y].TextureID = spriteBridge01

		return

	default:
		// Pobieramy tekstury właściwe dla rodzaju budowy
		template, ok := buildingTemplates[bld.Type]
		if !ok {
			fmt.Printf("Błąd przy próbie zastosowania grafiki ukończonej budowy rodzaju %d!\n", bld.Type)

			return
		}

		// Pobieramy lewy górny kafelek
		occupiedTile := bld.OccupiedTiles[0]

		for offsetY, row := range template {
			for offsetX, texID := range row {
				targetX := occupiedTile.X + uint8(offsetX)
				targetY := occupiedTile.Y + uint8(offsetY)

				// Podmieniamy teksturę na ostateczną
				board.Tiles[targetX][targetY].TextureID = texID
			}
		}
	}
}

// applyConstructionGraphics nakłada tekstury budowy na kafelki zajmowane przez budowę.
// Metoda ta zakłada, że bld.OccupiedTiles są poprawne dzięki zrobione przez boardData.RegisterBuilding.
func (board *boardData) applyConstructionGraphics(bld *building) {
	// Jest to „mechanicznie” wywoływane przez bState.placeConstructionSite
	// przy każdym zasadzeniu budowy.
	switch bld.Type {
	case buildingBridge:
		occupiedTile := bld.OccupiedTiles[0]
		board.Tiles[occupiedTile.X][occupiedTile.Y].TextureID = spriteBridgeConstruction
	case buildingPalisade:
		occupiedTile := bld.OccupiedTiles[0]
		board.Tiles[occupiedTile.X][occupiedTile.Y].TextureID = spritePalisadeDestroyed
		board.Tiles[occupiedTile.X][occupiedTile.Y].IsWalkable = true
	default:
		occupiedTile := bld.OccupiedTiles[0]

		for offsetY, row := range constructionTemplatePhase01 {
			for offsetX, texID := range row {
				targetX := occupiedTile.X + uint8(offsetX)
				targetY := occupiedTile.Y + uint8(offsetY)
				board.Tiles[targetX][targetY].TextureID = texID
			}
		}
	}
}

// Nakładamy połowiczne tekstury na budynek, który jest w połowie skończony
func (board *boardData) applyPhase2Graphics(bld *building) {
	// Palisady oraz mosty nie mają wyglądu połowicznego
	if bld.Type == buildingPalisade || bld.Type == buildingBridge {
		return
	}

	// Pobieramy tekstury odpowiednie dla danego rodzaju budynku w budowie
	template, ok := constructionTemplatesPhase02[bld.Type]
	if !ok {
		fmt.Printf("Ten rodzaj budynku nie ma tekstur połowicznych!\n")

		return
	}

	// Pobieramy lewy górny róg placu budowy
	occupiedTile := bld.OccupiedTiles[0]

	// Przechodzimy przez wszystkie wiersze
	for offsetY, row := range template {
		// Przechodzimy przez wszystkie kolumny
		for offsetX, texID := range row {
			// Ustalamy któremu kafelekowi zmienić teksturę
			targetX := occupiedTile.X + uint8(offsetX)
			targetY := occupiedTile.Y + uint8(offsetY)
			// Przypisujemy nową teksturę
			board.Tiles[targetX][targetY].TextureID = texID
		}
	}
}

// Palisady mają wyjątkowy cykl życia: zawsze można je odbudować.
// Dlatego potrzebują osobnej metody. Jedynym sposobem pozbycia się ich jest
// zbudowanie budynku na nich.
func (board *boardData) handlePalisadeDestruction(occupiedTile *tile) {
	occupiedTile.TextureID = spritePalisadeDestroyed
	occupiedTile.IsWalkable = true
}

// Usuwa budynek z planszy i wstawia ruiny na to miejsce.
// Nie mylić z usuwaniem z bState.Buildings, bo to osobna sprawa.
func (board *boardData) placeRuins(bld *building) {
	occupiedTile := bld.OccupiedTiles[0]

	for offsetY, row := range ruinsTemplate {
		for offsetX, texID := range row {
			targetX := occupiedTile.X + uint8(offsetX)
			targetY := occupiedTile.Y + uint8(offsetY)

			targetTile := &board.Tiles[targetX][targetY]

			targetTile.TextureID = texID
			targetTile.Building = nil
			targetTile.IsWalkable = false
		}
	}
}

func (board *boardData) isValidWalkableTile(x, y int8) bool {
	// Ponieważ chodzi o znalezienie wolnego kafelka
	// nie wykluczam jeszcze, że możemy dostać coś >= boardMax
	if x < 0 || x >= int8(boardMaxX) || y < 0 || y >= int8(boardMaxY) {
		return false
	}

	currentTile := &board.Tiles[x][y]

	// Sprawdzamy kafelek jest przechodni i nie ma na nim żadnego obiektu
	return currentTile.IsWalkable && currentTile.Unit == nil && currentTile.Building == nil
}

func (board *boardData) bldNeighborCoords(bld *building) []point {
	// 1. Bierzemy lewy górny róg budynku
	occupiedTileX := int(bld.OccupiedTiles[0].X)
	occupiedTileY := int(bld.OccupiedTiles[0].Y)

	// 2. Ponieważ budynki są zawsze 3na3 to
	// z góry znamy potencjalne współrzędne dla sąsiadów
	offsets := [16][2]int{
		{-1, -1},
		{0, -1},
		{1, -1},
		{2, -1},
		{3, -1},
		{3, 0},
		{3, 1},
		{3, 2},
		{3, 3},
		{2, 3},
		{1, 3},
		{0, 3},
		{-1, 3},
		{-1, 2},
		{-1, 1},
		{-1, 0},
	}

	// 3. Tworzymy listę współrzędnych mieszczących się
	//    w planszy
	var electedCoords []point

	// 4. Sprawdzamy, czy potencjalne wspołrzędne w rzeczywistości
	//    mieszczą się w planszy
	for _, offset := range offsets {
		// Tworzymy potencjalne współrzędne sąsiada
		electedTileX := occupiedTileX + offset[0]
		electedTileY := occupiedTileY + offset[1]

		// czy mamy poprawny X
		if electedTileX >= 0 && electedTileX < int(boardMaxX) &&
			// oraz Y
			electedTileY >= 0 && electedTileY < int(boardMaxY) {
			// Obie współrzędne są poprawnymi współrzędnymi kafelków więc
			// dodajemy do listy prawidłowych współrzędnych
			electedCoords = append(electedCoords, point{
				X: uint8(electedTileX),
				Y: uint8(electedTileY),
			})
		}
	}

	// Stworzyliśmy listę wszystkich sąsiadów budynku
	return electedCoords
}

func (board *boardData) hasFreeTileInList(electedTiles []point) bool {
	for _, electedTile := range electedTiles {
		if board.Tiles[electedTile.X][electedTile.Y].Unit == nil &&
			board.Tiles[electedTile.X][electedTile.Y].IsWalkable {
			return true
		}
	}

	return false
}

func (board *boardData) hasSpaceAroundBuilding(bld *building) bool {
	coords := board.bldNeighborCoords(bld)

	return board.hasFreeTileInList(coords)
}

func (board *boardData) getFreeTileInList(electedCoords []point) (point, bool) {
	for _, checkTile := range electedCoords {
		if board.Tiles[checkTile.X][checkTile.Y].Unit == nil && board.Tiles[checkTile.X][checkTile.Y].IsWalkable {
			return point{X: checkTile.X, Y: checkTile.Y}, true
		}
	}

	return point{X: 0, Y: 0}, false
}

func (board *boardData) electSpawnTile(bld *building) (point, bool) {
	coords := board.bldNeighborCoords(bld)

	return board.getFreeTileInList(coords)
}

func getDistanceToUnit(bldType buildingType, bldTopLeft point, unitX, unitY uint8) uint8 {
	var minX, minY, maxX, maxY uint8

	// Na podstawie rodzaju budynku wybieramy skrajne współrzędne
	switch bldType {
	case buildingBridge, buildingPalisade:
		minX = bldTopLeft.X
		maxX = minX
		minY = bldTopLeft.Y
		maxY = minY
	// Zakładamy, że są tylko budynki jedno- i dziewięciokafelkowe
	// Dlatego w domyślnym przypadku nie wywalam błędu, a obsługuję
	// resztę rodzajów budynku.
	default:
		// Ponieważ zwyczajne budynki mają 9 kafelków wykorzystajmy
		// fakt, że [0] to zawsze najmniejsze X i Y
		// a [8] to zawsze największe X i Y
		minX = bldTopLeft.X
		maxX = minX + 2
		minY = bldTopLeft.Y
		maxY = minY + 2
	}

	var differenceX, differenceY uint8

	if unitX < minX {
		differenceX = minX - unitX
	} else if unitX > maxX {
		differenceX = unitX - maxX
	}
	// Jeśli unitX = minX, to zostawiamy
	// differenceX jako zero

	if unitY < minY {
		differenceY = minY - unitY
	} else if unitY > maxY {
		differenceY = unitY - maxY
	}
	// Jeśli unitY = minY, to zostawiamy
	// differenceY jako zero

	// Zwracamy największą różnicę
	if differenceX > differenceY {
		return differenceX
	}

	return differenceY
}

// @reminder: funkcja przyjmuje aż 5 argumentów, więc można by przekazać strukturę z celami zamiast każdy osobno.
//    Nie wiem, czy tak byłoby lepiej dlatego tak nie robię. Może w przyszłości się zdecyduję na zmianę.
func findTileForAttacking(attacker *unit, targetU *unit, targetBld *building, targetTile *point, board *boardData) ([]point, bool) {
	var validCoords []point // wykaz prawidłowych kafelków, które można odwiedzić.

	var rangeAdjustment uint8

	var targetX, targetY uint8

	// W zależności od tego co jest celem musimy się inaczej przygotować.
	switch {
	case targetBld != nil:
		if targetBld.Type != buildingPalisade && targetBld.Type != buildingBridge {
			rangeAdjustment = 1
		}

		var ok bool

		targetX, targetY, ok = targetBld.getCenter()
		if !ok {
			return nil, ok
		}

	case targetU != nil:
		targetX, targetY = targetU.X, targetU.Y
	case targetTile != nil:
		targetX, targetY = targetTile.X, targetTile.Y

	default:
		// To nigdy nie powinno mieć miejsca!
		return nil, false
	}

	attackRange := attacker.AttackRange + rangeAdjustment

	// Wszelkie możliwe X. Nigdy nie przekroczymy +-120 więc zmiana na int8 jest bezpieczna.
	for coordX := int8(targetX - attackRange); coordX <= int8(targetX+attackRange); coordX++ { //nolint:gosec
		// Wszelkie możliwe Y
		for coordY := int8(targetY - attackRange); coordY <= int8(targetY+attackRange); coordY++ { //nolint:gosec
			// Tutaj sprawdzamy, czy to prawidłowe współrzędne kafelka.
			if board.isValidWalkableTile(coordX, coordY) {
				// board.isValidWalkableTile gwarantuje, że 0 <= attackX/Y <= 65. Dlatego zmiana na uint8 jest
				// bezpieczna.                             ↓↓↓↓↓             ↓↓↓↓↓
				validCoords = append(validCoords, point{X: uint8(coordX), Y: uint8(coordY)}) //nolint:gosec
			}
		}
	}

	// Jeśli targetTree != nil to musimy wywalić kafelek na lewo od drzewa, inaczej jednostka zginie.
	if targetTile != nil && isTree(board.Tiles[targetTile.X][targetTile.Y].TextureID) && targetTile.X > 0 {
		toRemove := point{X: targetTile.X - 1, Y: targetTile.Y}

		indexToRemove := slices.Index(validCoords, toRemove)
		if indexToRemove != -1 {
			validCoords = slices.Delete(validCoords, indexToRemove, indexToRemove+1)
		}
	}

	return validCoords, true
}

// Odpowiada za wybranie kafelka o najkrótszej drodze.
func findBestReachableTile(u *unit, validCoords []point, board *boardData) (*point, error) {
	var bestX, bestY uint8

	minPathLen := math.MaxInt32
	found := false

	for _, coord := range validCoords {
		path := findPath(board, u, u.X, u.Y, coord.X, coord.Y)

		if path != nil && len(path) < minPathLen {
			minPathLen = len(path)
			bestX, bestY = coord.X, coord.Y
			found = true
		}
	}

	if found {
		return &point{X: bestX, Y: bestY}, nil
	}

	return nil, fmt.Errorf("nie ma prawidłowego kafelka")
}

// @todo: ani to units, ani nie przystaje do obecnej architektury.
func (bld *building) getClosestOccupiedTile(fromX, fromY uint8) (uint8, uint8, bool) {
	if len(bld.OccupiedTiles) == 0 {
		return 0, 0, false
	}

	closestX, closestY := uint8(0), uint8(0)
	minDistSq := math.MaxFloat64

	for _, occupiedTile := range bld.OccupiedTiles {
		dx := float64(occupiedTile.X - fromX)
		dy := float64(occupiedTile.Y - fromY)
		distSq := dx*dx + dy*dy

		if distSq < minDistSq {
			minDistSq = distSq
			closestX = occupiedTile.X
			closestY = occupiedTile.Y
		}
	}

	return closestX, closestY, true
}
