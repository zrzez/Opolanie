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
		// Informację o tym, dany kafelek zawiera drzewo mamy przechowywaną tylko w teksturze!
		// Drzewo pali się inaczej niż trawa.
		case isTreeStump(burningTile.TextureID):
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

// Odpowiada za zarządzanie logiką upadania drzewa, które spłonęło lub zostało ścięte.
func fallingTreeEffect(bs *battleState) {
	// @todo: jeśli na lewo od drzewa jest suche drzewo to ono też ma być obalone
	for _, currentTile := range bs.FallingTreesList {
		// Drzewo zmienia fazę co 10 logicznych tyknięć.
		if bs.GlobalFrameCounter%10 == 0 {
			switch currentTile.treeFallPhase {
			case treeStraight:
				// czeka i się przechyla bardziej
				// @todo: ustaw docelowe po sprawdzeniu, czy mechanizm działa
				currentTile.TextureID = spriteDryTreeFallingStump02_0
				currentTile.treeFallPhase = treeLeaning

			case treeLeaning:
				// czeka i się przechyla bardziej
				// @todo: ustaw docelowe po sprawdzeniu, czy mechanizm działa
				currentTile.TextureID = spriteDryTreeFallingStump01_1
				currentTile.treeFallPhase = treeImpact

			case treeImpact:
				// @todo: ustaw docelowe tekstury po sprawdzeniu, czy mechanizm działa
				currentTile.TextureID = spriteDryTreeFallingStump00_2

				if currentTile.X-1 > 0 {
					adjacentTile := &bs.Board.Tiles[currentTile.X-1][currentTile.Y]

					if adjacentTile.TextureID == spriteDryTreeStump00 {
						adjacentTile.treeFall(bs)
					} else {
						adjacentTile.applyFallingTreeDamage(bs)
					}
				}

				currentTile.IsWalkable = true
				currentTile.treeFallPhase = treeFell

			case treeFell:
				return
			}
		}
	}
	// „Odczekać chwilkę” powinno trawć jedną pętlę falowania wody.

	// 1.a Odczekać chwilkę
	// 1.b Ustawiamy teksturę jako pierwszy stopień przechylenia

	// 2.a Odczekać chwilkę
	// 2.b Ustawiamy teksturę jako drugi stopień przechylenia

	// 3.a Odczekać chwilkę
	// 3.b Ustawiamy teksturę jako trzeci, ostatni stopień przechylenia
	// 3.c jednostki/budynki na lewo od kafelka powinni otrzymać obrażenia
	// @reminder: przy atakowaniu kafleka jednostka przechowuje współrzędne celu.
	// jeśli przekażę/dam dostęp do battlestate.Board to wtedy mogę dostać się do sąsiada
}
