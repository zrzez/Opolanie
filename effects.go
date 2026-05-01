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

func updateAnimationCounters(bs *battleState) {
	bs.WaterAnimationFrame = (bs.WaterAnimationFrame + 1) % 3
	bs.FireAnimationFrame = (bs.FireAnimationFrame + 1) % 14
	bs.GrassGrowthCycle++
}

// obsługuje efekty na mapie, które zmieniają się w czasie.
func updateWorldTimers(bs *battleState) {
	if bs.GrassGrowthCycle > maxGrassGrowthCounter {
		bs.GrassGrowthCycle = 0

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
	for _, healingSpot := range bs.HealingShrines {
		// Wybieramy kapliczkę
		healingTile := &bs.Board.Tiles[healingSpot.X][healingSpot.Y]
		// Sprawdzamy, czy w kapliczce ktoś jest
		if healingTile.Unit != nil {
			// bierzemy znalezioną jednostkę
			u := healingTile.Unit
			// Sprawdzamy, czy potrzebuje leczenia
			// takie sprawdzenie, to raczej powinno być w tryToHeal, czy cos
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

// Odpowiada za zarządzanie logiką płonięcia kafelka.
func burningTileEffect(bs *battleState) {
	for _, burningTile := range bs.BurningTilesList {
		switch {
		// Drzewo pali się inaczej niż trawa.
		case burningTile.isTree():
			burningTile.processTreeFire(bs)

		// Zwyczajny przypadek zapalenia się kafelka.
		case burningTile.IsBurning:
			burningTile.processNormalFire()

			// Co osiem klatek dodatkowe obrażenia od ognia.
			if bs.GlobalFrameCounter%8 == 0 {
				burningTile.applyFireDamage(bs)
			}

		// Jeśli kafelek się już nie pali, to można przejść do zarządzania popiołem
		case burningTile.hasAsh:
			burningTile.processAshDecay()
		}
	}
}

func ghostEffect(bs *battleState) {
	for _, ghostTile := range bs.GhostsList {
		if ghostTile.GhostEffectCounter > 0 {

			ghostTile.GhostEffectCounter--
			// 2. Zadaję obrażenia
			if ghostTile.Unit != nil && ghostTile.Unit.Exists {
				ghostTile.Unit.takeDamage(ghostTile.GhostDamage, bs)
			}
		} else {
			ghostTile.GhostEffect = false
		}
	}
}

// Odpowiada za zarządzanie logiką upadania drzewa, które spłonęło lub zostało ścięte.
func fallingTreeEffect(bs *battleState) {
	for _, currentTile := range bs.FallingTreesList {
		if bs.GlobalFrameCounter%10 == 0 {
			switch currentTile.treeState {
			case treeStraight:
				currentTile.treeState = treeLeaning

			case treeLeaning:
				// 1. Lekko przechylone, 3 kafelki do narysowania
				if !currentTile.IsBurnt {
					currentTile.TextureID = spriteDryFallingStump
				} else {
					currentTile.TextureID = spriteBurntFallingTreeStump
				}

				currentTile.treeState = treeFalling

			case treeFalling:
				// 2. przechylone 45 stopni, dwa kafelki do narysowania
				if !currentTile.IsBurnt {
					currentTile.TextureID = spriteDryFallenTreeStump
				} else {
					currentTile.TextureID = spriteBurntFallenTreeStump
				}

				currentTile.treeState = treeImpact

			case treeImpact:
				// 3. Upadło, wywołuje określone efekty na sąsiednich kafelkach jeśli
				// Drzewo upada w granicach planszy, żeby móc bezpiecznie wywołać efekty
				// na sąsiednim kafelku
				if currentTile.X-1 > 0 {
					adjacentTile := &bs.Board.Tiles[currentTile.X-1][currentTile.Y]

					// Obalamy sąsiednie suche drzewo
					if adjacentTile.TextureID == spriteDryTreeStump00 {
						adjacentTile.treeFall(bs)
					} else {
						// Lub zadajemy obrażenia
						adjacentTile.applyFallingTreeDamage(bs)
					}
				}

				currentTile.IsWalkable = true
				currentTile.treeState = treeFell

			case treeFell:
				// 4. Leży i już nic nie robi
				return

			case noTree:
				panic("noTree nie może być obsługiwane przez fallingTreeEffect")

			default:
				// Upadanie drzew nie może obslugiwać noTree ponieważ jest to stan braku drzewa.
				panic("noTree nie może być obsługiwane przez fallingTreeEffect")
			}
		}
	}
}
