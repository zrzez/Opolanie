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

// @reminder: ogień przygasa w powiązaniu z licznikiem odrastania trawy
// @todo: W tej chwili nad tym pracuję 16.04.2026
func burningTileEffect(bs *battleState) {
	// tymczasowo muszę iterować po całej planszy

	for y := range boardMaxY {
		for x := range boardMaxX {
			affectedTile := &bs.Board.Tiles[x][y]

			processAshDecay(affectedTile)

			// Róznego rodzaju kafelki różnie się palą
			// @todo: przemyśl, czy nie lepiej to złączyć w jedną funkcję
			if isTreeStump(affectedTile.TextureID) {
				// @reminder: Z tego miejsca mogę przekazać sąsiedni kafelek aby przetworzyć poprawnie
				// upadające drzewo.
				processTreeFire(affectedTile)
			} else {
				processNormalFire(affectedTile)
			}

			if bs.GlobalFrameCounter%8 != 0 && affectedTile.IsBurning {
				if affectedTile.Unit != nil && affectedTile.Unit.Exists {
					affectedTile.Unit.takeDamage(3, bs)
				}

				if affectedTile.Building != nil && affectedTile.Building.Exists {
					affectedTile.Building.takeDamage(3)
				}
			}
			// 0. Kafelek płonie
			// tile.IsBurning = true
			// jeśli płonie to jednostka chce zejść z niego - to zapewne musi być gdzieś w units.go

			// 1. Co osiem tików zadaje 3 obrażenia i zmniejsza intensywność o 1
			// tile.FireCounter = powiązany z GrassGrowthCounter
			// tile.FireID =

			// 2. Jeśli kafelek miał drzewo, to kończy się palić i podmienia teksturę na spalone drzewo
			// 2a. spalone drzewo powinno się obalić
		}
	}
}

func processNormalFire(affectedTile *tile) {
	if !affectedTile.IsBurning {
		return
	}

	// Gromadzenie się popiołu
	affectedTile.AshIntensity += ashAccumulationRate

	if affectedTile.AshIntensity > 1.0 {
		affectedTile.AshIntensity = 1.0
	}

	// Właściwe płonięcie
	affectedTile.BurnElapsed++

	var currentFireSprite uint16

	switch {
	case affectedTile.BurnElapsed < bigBurn:
		{
			currentFireSprite = spriteFire00
		}
	case affectedTile.BurnElapsed < midBurn:
		{
			currentFireSprite = spriteFire04
		}
	case affectedTile.BurnElapsed < minBurn:
		{
			currentFireSprite = spriteFire08
		}
	default:
		affectedTile.IsBurning = false
		affectedTile.AshAge = 0
		affectedTile.AshProcessState = ashDecaying

		return
	}

	affectedTile.BurnOverlayID = currentFireSprite
}

func processAshDecay(affectedTile *tile) {
	// Płonące kafelki nie tracą popiołu
	if affectedTile.IsBurning {
		return
	}

	// Jeśli nie ma popiołu lub są śladowe ilości to wychodzimy
	if !affectedTile.IsAsh || affectedTile.AshIntensity < 0.01 {
		return
	}

	switch affectedTile.AshProcessState {
	case ashDecaying:
		affectedTile.AshIntensity *= (1.0 - ashDecayRate)

		if affectedTile.AshAge >= totalAshLifetime {
			affectedTile.IsAsh = false
			affectedTile.AshIntensity = 0.0
			affectedTile.AshProcessState = ashFinished
		}
	case ashFinished:
		return

	// domyślne niepowinno nigdy wystąpić
	default:
		return
	}

	affectedTile.CurrentAshAlpha = affectedTile.AshIntensity
	affectedTile.AshAge++
}

func (t *tile) setOnFire(fireSize uint16) {
	t.IsBurning = true
	t.BurnElapsed = fireSize - bigBurn

	if !isTreeStump(t.TextureID) {
		t.IsAsh = true
	}
}

func processTreeFire(affectedTile *tile) {
	if !affectedTile.IsBurning {
		return
	}

	// Właściwe płonięcie
	affectedTile.BurnElapsed++

	var currentFireSprite uint16

	switch {
	case affectedTile.BurnElapsed < bigBurn:
		{
			currentFireSprite = spriteFire00
		}
	case affectedTile.BurnElapsed < midBurn:
		{
			currentFireSprite = spriteFire04
		}
	case affectedTile.BurnElapsed < minBurn:
		{
			currentFireSprite = spriteFire08
		}
	default:
		affectedTile.IsBurning = false
		processBurntTree(affectedTile)

		return
	}

	affectedTile.BurnOverlayID = currentFireSprite
}

func processBurntTree(affectedTile *tile) {
	// Ustalamy nowe tekstury na spalone drzewa
	if affectedTile.TextureID < spriteTreeStump03 {
		affectedTile.TextureID = spriteTreeBurntStump00
	} else {
		affectedTile.TextureID = spriteTreeBurntStump01
	}

	// Obalamy spalone drzewo
	affectedTile.treeFall()
}

func (t *tile) treeFall() {
	// Po chwili czekania dajemy spriteTreeFalling.
	// Później spriteTreeFallen.
	// @todo: problem z tym, że nie wiem, czy to pojedyncze duszki
	// @todo: potrzebuję dodać obrażenia od spadającej korony
	// @todo: jeśli na lewo od drzewa jest suche drzewo to ono też ma być obalone
	// @reminder: potrzebuję dostępu do współrzędnych! bez tego nie mam wpływu na sąsiedni kafelek
	// kafelek staje się przechodni.
}
