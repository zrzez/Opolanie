package main

// effects.go

// Funkcja do łączenia dróg/palisad
func joinRoadsPalisade(x, y uint8, bld *building, bs *battleState) {
	if x >= boardMaxX || y >= boardMaxY {
		return
	}

	// tile := &bs.Board.Tiles[x][y]

	switch bld.Type {
	// case buildingRoad: // droga
	// 	if tile.TextureID >= spriteRoadStart && tile.TextureID <= spriteRoadEnd {
	// 		tile.TextureID = spriteRoadStart + 4
	// 	}
	case buildingPalisade:
		applyPalisadeProcessing(x, y, bs.Board)
		updateAdjacentPalisades(x, y, bs.Board)
	default:
		return
	}
}

func updateAdjacentPalisades(x, y uint8, board *boardData) {
	// Góra
	if y > 0 && isPalisade(board.Tiles[x][y-1].TextureID) && board.Tiles[x][y-1].Building != nil {
		if !board.Tiles[x][y-1].Building.IsUnderConstruction {
			applyPalisadeProcessing(x, y-1, board)
		}
	}

	// Dół
	if y < boardMaxY-1 && isPalisade(board.Tiles[x][y+1].TextureID) && board.Tiles[x][y+1].Building != nil {
		if !board.Tiles[x][y+1].Building.IsUnderConstruction {
			applyPalisadeProcessing(x, y+1, board)
		}
	}

	// Prawo
	if x < boardMaxX-1 && isPalisade(board.Tiles[x+1][y].TextureID) && board.Tiles[x+1][y].Building != nil {
		if !board.Tiles[x+1][y].Building.IsUnderConstruction {
			applyPalisadeProcessing(x+1, y, board)
		}
	}

	// Lewo
	if x > 0 && isPalisade(board.Tiles[x-1][y].TextureID) && board.Tiles[x-1][y].Building != nil {
		if !board.Tiles[x-1][y].Building.IsUnderConstruction {
			applyPalisadeProcessing(x-1, y, board)
		}
	}
}

func joinBridges(x, y uint8, board *boardData) {
	applyBridgeProcessing(x, y, board)
	updateAdjacentBridges(x, y, board)
}

func updateAdjacentBridges(x, y uint8, board *boardData) {
	// Góra
	if y > 0 && isBridge(board.Tiles[x][y-1].TextureID) {
		b := board.Tiles[x][y-1].Building
		if b == nil || !b.IsUnderConstruction {
			applyBridgeProcessing(x, y-1, board)
		}
	}

	// Dół
	if y < boardMaxY-1 && isBridge(board.Tiles[x][y+1].TextureID) {
		b := board.Tiles[x][y+1].Building
		if b == nil || !b.IsUnderConstruction {
			applyBridgeProcessing(x, y+1, board)
		}
	}

	// Prawo
	if x < boardMaxX-1 && isBridge(board.Tiles[x+1][y].TextureID) {
		b := board.Tiles[x+1][y].Building
		if b == nil || !b.IsUnderConstruction {
			applyBridgeProcessing(x+1, y, board)
		}
	}

	// Lewo
	if x > 0 && isBridge(board.Tiles[x-1][y].TextureID) {
		b := board.Tiles[x-1][y].Building
		if b == nil || !b.IsUnderConstruction {
			applyBridgeProcessing(x-1, y, board)
		}
	}
}

func applyBridgeProcessing(x, y uint8, board *boardData) {
	newID := spriteBridgeStart

	// Czy sąsiad poziomo jest mostem?
	hasHorizontalBridge := (x > 0 && isBridge(board.Tiles[x-1][y].TextureID)) ||
		(x < boardMaxX-1 && isBridge(board.Tiles[x+1][y].TextureID))

	// Czy sąsiad pionowo jest mostem?
	hasVerticalBridge := (y > 0 && isBridge(board.Tiles[x][y-1].TextureID)) ||
		(y < boardMaxY-1 && isBridge(board.Tiles[x][y+1].TextureID))

	// Jeśli most ma sąsiada-mostu, dopasuj kierunek
	if hasHorizontalBridge {
		newID = spriteBridgeEnd
	} else if hasVerticalBridge {
		newID = spriteBridgeEnd - 1
	}

	board.Tiles[x][y].TextureID = newID
}

func updateAnimationCounters(bs *battleState) {
	bs.WaterAnimationCounter = (bs.WaterAnimationCounter + 1) % 3
	bs.FireAnimationCounter = (bs.FireAnimationCounter + 1) % 14
	bs.GrassGrowthCounter++
}

// obsługuje efekty na mapie, które zmieniają się w czasie.
func updateWorldTimers(bs *battleState) {
	if bs.GrassGrowthCounter > maxGrassGrowthCounter {
		bs.GrassGrowthCounter = 0

		for y := uint8(1); y < boardMaxY-1; y++ {
			for x := uint8(1); x < boardMaxX-1; x++ {
				currentTile := &bs.Board.Tiles[x][y]

				switch currentTile.GrazedOverlayID {
				case uint8(spriteGrassStubbed):
					currentTile.GrazedOverlayID = uint8(spriteGrassGrazed)
				case uint8(spriteGrassGrazed):
					currentTile.IsGrazed = false
				}
			}
		}
	}
}

// healingShrine powinno leczyć i przywracać manę.
// @todo: Należy powiązać czynność leczenia z animacją
func healingShrine(bs *battleState) {
	if bs.GlobalFrameCounter%4 != 0 {
		return
	}
	// jeżeli kapliczka może leczyć
	for _, point := range bs.HealingShrines {
		// Wybieramy kapliczkę
		tile := &bs.Board.Tiles[point.X][point.Y]
		// Sprawdzamy, czy w kapliczce ktoś jest
		if tile.Unit != nil {
			// bierzemy znalezioną jednostkę
			u := tile.Unit
			// Sprawdzamy, czy potrzebuje leczenia
			if u.HP != u.MaxHP {
				switch u.Type {
				case unitCow:
					u.increaseHPUnit(1)
				case unitUnknown:
					// Specjalny przypadek ponieważ strzyga
					// zamienia manę na życie
					// @todo: ogarnij, czy w pierwowzorze też tak było
					u.increaseHPUnit(4)
					u.increaseManaUnit(5)

				default:
					u.increaseHPUnit(4)
				}
			}

			if u.Type != unitUnknown && u.MaxMana > 0 && u.Mana < u.MaxMana {
				u.increaseManaUnit(5)
			}
		}
	}
}

func manaRegen(bs *battleState) {
	// 15 Hz dla 60 klatek na sekundę, aby zachować zgodność z pierwowzorem
	// @todo: to się wyłoży jeśli będzie można zmienić szybkość gry
	if bs.GlobalFrameCounter%4 != 0 {
		return
	}

	// Każda żywa jednostka powinna odnowić część many
	for _, unit := range bs.Units {
		if unit.Type.hasMana() {
			unit.increaseManaUnit(1)

			// Dodatek dla maga
			if unit.Type == unitMage {
				unit.increaseManaUnit(1)
			}
		}
	}
}
