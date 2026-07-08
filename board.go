package main

import (
	"fmt"
	"log"
)

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

// @todo zweryfikuj, czy to technicznie możliwe aby dostać x,y >= boardMaxX.
func (board *boardData) isValidWalkableTile(x, y uint8) bool {
	// Ponieważ chodzi o znalezienie wolnego kafelka
	// nie wykluczam jeszcze, że możemy dostać coś >= boardMax
	if x >= boardMaxX || y >= boardMaxY {
		return false
	}

	currentTile := &board.Tiles[x][y]

	// Sprawdzamy kafelek jest przechodni i nie ma na nim żadnego obiektu
	return currentTile.IsWalkable && currentTile.Unit == nil && currentTile.Building == nil
}

func (board *boardData) neighborCoords(bld *building) []point {
	// LOGOWANIE
	fmt.Println("=== neighborCoords: budynek ===")

	for i, p := range bld.OccupiedTiles {
		fmt.Printf("kafelek %d: (%d, %d)\n", i, p.X, p.Y)
	}

	fmt.Println("=== neighborCoords: budynek ===")

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

		// czy mamy legalny X
		if electedTileX >= 0 && electedTileX < int(boardMaxX) &&
			// czy mamy Y
			electedTileY >= 0 && electedTileY < int(boardMaxY) {
			// Obie współrzędne są poprawnymi współrzędnymi kafelków więc
			// dodajemy do listy prawidłowych współrzędnych
			electedCoords = append(electedCoords, point{
				X: uint8(electedTileX),
				Y: uint8(electedTileY),
			})
		}
	}

	fmt.Println("=== neighborCoords: sąsiedzi w planszy ===")

	for i, p := range electedCoords {
		fmt.Printf("sąsiad %d: (%d, %d)\n", i, p.X, p.Y)
	}

	fmt.Println("=== neighborCoords: sąsiedzi w planszy ===")

	// Stworzyliśmy listę wszystkich sąsiadów budynku
	return electedCoords
}

func (board *boardData) hasFreeTileInList(electedTiles []point) bool {
	fmt.Println("=== hasFreeTileInList: sprawdzanie kafelków ===")

	for _, electedTile := range electedTiles {
		unitNil := board.Tiles[electedTile.X][electedTile.Y].Unit == nil
		walkable := board.Tiles[electedTile.X][electedTile.Y].IsWalkable
		fmt.Printf("   (%d, %d): Unit==nil = %v, IsWalkable = %v\n", electedTile.X, electedTile.Y, unitNil, walkable)
		if unitNil && walkable {
			return true
		}
	}

	fmt.Println("=== hasFreeTileInList: sprawdzanie kafelków ===")

	/*for _, electedTile := range electedTiles {
		if board.Tiles[electedTile.X][electedTile.Y].Unit == nil &&
			board.Tiles[electedTile.X][electedTile.Y].IsWalkable {
			return true
		}
	}
	*/
	return false
}

func (board *boardData) hasSpaceAroundBuilding(bld *building) bool {
	coords := board.neighborCoords(bld)

	return board.hasFreeTileInList(coords)
}

func (board *boardData) getFreeTileInList(electedCoords []point) (point, bool) {
	for _, checkTile := range electedCoords {
		if board.Tiles[checkTile.X][checkTile.Y].Unit == nil && board.Tiles[checkTile.X][checkTile.Y].IsWalkable {
			return point{X: checkTile.X, Y: checkTile.Y}, true
		}
	}

	// @todo: nie można przekazywać 0,0 jako „poprawnego”
	return point{}, false
}

func (board *boardData) electSpawnTile(bld *building) (point, bool) {
	coords := board.neighborCoords(bld)

	return board.getFreeTileInList(coords)
}
