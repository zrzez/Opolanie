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
