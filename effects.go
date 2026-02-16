package main

// effects.go

// ============================================================================
// EFEKTY I ATAK
// ============================================================================

// Funkcja do łączenia dróg/palisad
func joinRoadsPalisade(x, y uint8, baseType uint16, bs *battleState) {
	if x < 0 || x >= boardMaxX || y < 0 || y >= boardMaxY {
		return
	}

	tile := &bs.Board.Tiles[x][y]

	switch baseType {
	case spriteRoadStart: // droga
		if tile.TextureID >= spriteRoadStart && tile.TextureID <= spriteRoadEnd {
			tile.TextureID = spriteRoadStart + 4
		}
	case spritePalisadeStart: // palisada
		if tile.TextureID >= spritePalisadeStart && tile.TextureID <= spritePalisadeEnd {
			tile.TextureID = spritePalisadeStart
		}
	}
}

func updateAnimationCounters(bs *battleState) {
	bs.WaterAnimationCounter = (bs.WaterAnimationCounter + 1) % 3
	bs.FireAnimationCounter = (bs.FireAnimationCounter + 1) % 14
	// ↓↓↓↓↓ Z jakiegoś powodu wzrost tego licznika niszczy trawę
	bs.GrassGrowthCounter++ // @todo: ogarnij, czy mogę to usunąć/odkomentować żeby trawa odrastała
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
					currentTile.IsGrazed = false // Powrót do bazowej zdrowej trawy
				}
				// === NOWA LOGIKA ODRASTANIA (KONIEC) ===

				// @cholera wie co to jest (zostawiamy bez zmian, bo to dotyczy efektów, nie tekstury)
				//  if currentTile.EffectID == 2 || currentTile.EffectID == 3 {
				//  	currentTile.EffectID--
				//  }
			}
		}
	}

	// --- EFEKTY KRÓTKOTERMINOWE (co klatkę) ---
	// for y := range boardMaxY {
	//	for x := range boardMaxX {
	//		tile := &bs.Board.Tiles[x][y]

	// Logika ognia (używano PlcFogOfWar jako licznika?)
	// w oryginale PlcN to 'sąsiedztwo'/flagi, a PlaceAttackEffects to ogień.
	// Tutaj w Pana kodzie: PlcFogOfWar > 70 to ogień?
	// OK, przepisuję logikę 1:1, ale mapuję na EffectID (bo to efekt).
	// Jeśli PlcFogOfWar miało robić za Effect, to użyjmy EffectID.

	//effectVal := &tile.EffectID // Zamiast PlcFogOfWar dla efektów ognia
	//
	//// Logika ognia
	//if *effectVal > 70 && *effectVal < 101 {
	//	if (bs.GrassGrowthCounter & 7) == 0 {
	//		// Tu było PlcAttackEffects += 3 i *PlcNValue--
	//		// Rozbijmy to:
	//		// Jeśli EffectID to licznik czasu ognia, to go zmniejszamy.
	//		// A siłę ognia (dawne PlcAttackEffects) zwiększamy?
	//		// Kod był: bs.Board.PlcAttackEffects[x][y] += 3; *PlcNValue--
	//		// Zakładam, że teraz EffectID przechowuje stan ognia.
	//
	//		// Uproszczenie dla nowej architektury:
	//		// Zmniejszamy czas trwania efektu
	//		*effectVal--
	//	}
	//}
	//
	//// Wypalenie się ognia
	//if *effectVal == 70 {
	//	*effectVal = 3 // Zostaje popiół
	//
	//	if tile.TextureID > spriteTreeStump00 && tile.TextureID <= spriteTreeStump03 {
	//		tile.TextureID = spriteTreeBurntStump01
	//	}
	//
	//	if tile.TextureID >= spriteTreeStump03 && tile.TextureID <= spriteTreeStumpEnd {
	//		tile.TextureID = spriteTreeBurntStump00
	//	}
	//}
	//	}
	//}
}

// @todo: Zupełnie nie pamiętam, co się tutaj odwala. Wiele rzeczy z wymienionych tutaj
// nigdy nie sprawdzałem więc raczej ładnie wyglądające bzdury, a nie działający kod
// clearAttack przetwarza i resetuje stany związane z atakami i zmianami terenu.
//func clearAttack(bs *battleState) {
//	for j := range boardMaxX {
//		for i := range boardMaxX {
//			tile := &bs.Board.Tiles[i][j]
//
//			// Postaw palisadę
//			if tile.TextureID == spritePalisadeDestroyed && tile.EffectID > 225 {
//				tile.EffectID = 0
//				tile.TextureID = spritePalisadeStart
//
//				// Jeśli brak obiektu, stwórz budynek Palisady
//				if tile.Building == nil {
//					newPalisade := &building{}
//					newPalisade.initConstruction(i, j, buildingPalisade, colorNone, bs)
//					bs.Buildings = append(bs.Buildings, newPalisade)
//					// Rejestracja w Tiles (zrobi to init, ale dla pewności)
//					tile.Building = newPalisade
//					tile.IsWalkable = false
//					log.Printf("INFO: Zbudowano palisadę (obiekt) na (%d,%d).", i, j)
//				}
//
//				joinRoadsPalisade(i, j, 266, bs)
//				joinRoadsPalisade(i-1, j, 266, bs)
//				joinRoadsPalisade(i+1, j, 266, bs)
//				joinRoadsPalisade(i, j-1, 266, bs)
//				joinRoadsPalisade(i, j+1, 266, bs)
//
//				// CMD_BUILD_ROAD to 15, ale Droga() oczekuje typu bazowego grafiki (np. 35).
//				// Zostawmy 266.
//
//				tile.Visibility = visibilityVisible // Odkryj teren
//			}
//
//			// Postaw most
//			// Tutaj stary kod sprawdzał Plc[i][j] == 2 (most w budowie?).
//			// w nowej architekturze: sprawdzamy TextureID lub building?
//			// Załóżmy, że TextureID == most w budowie (np. 2).
//			// Jeśli nie mamy ID mostu w budowie, pomijam ten blok lub adaptuję:
//			// if tile.TextureID == ??? && tile.EffectID > 225 ...
//
//			// Wyczyść atak jeśli nie ma jednostki (i nie jest to efekt terenu)
//			// Dawniej: if Plc < 256 ...
//			if tile.Unit == nil && tile.Building == nil && tile.EffectID != 0 && tile.EffectID < 200 {
//				// Jeśli to nie jest spadające drzewo (200), a pole puste -> czyść
//				tile.EffectID = 0
//			}
//
//			// Zniszcz palisadę
//			if tile.TextureID >= spritePalisadeStart && tile.TextureID < spritePalisadeEnd &&
//				(tile.EffectID > 120 || tile.EffectID < 0) {
//				tile.EffectID = 0
//				tile.TextureID = spritePalisadeDestroyed // Zniszczona
//
//				// Usuń obiekt
//				tile.Building = nil
//				tile.IsWalkable = true // Można przejść po zgliszczach?
//				tile.Visibility = visibilityVisible
//				log.Printf("INFO: Palisada zniszczona w (%d,%d)", i, j)
//			}
//
//			// Obsługa ognia (powtórka z updateWorldTimers?)
//			if tile.Visibility != visibilityUnexplored && tile.EffectID > 70 && tile.EffectID < 101 {
//				if (bs.GrassGrowthCounter & 7) == 0 {
//					// Zwiększamy intensywność czy zmniejszamy?
//					// Kod: += 3, --. Czyli +2 netto?
//					tile.EffectID += 2
//				}
//			}
//
//			// Drzewo się wali
//			// Użyjemy EffectID dla stanu walenia się
//			if tile.EffectID > 190 && tile.EffectID < 201 {
//				tile.EffectID--
//			}
//
//			// Ścięcie drzewa
//			if tile.TextureID >= spriteTreeStumpStart && tile.TextureID <= spriteTreeStumpEnd &&
//				tile.EffectID > 150 {
//
//				tile.EffectID = 200 // Animacja
//				log.Printf("INFO: Drzewo ścięte w (%d, %d)", i, j)
//			}
//
//			// Ścięcie spalonego drzewa
//			if tile.EffectID > 75 && (tile.TextureID == spriteTreeBurntStump00 ||
//				tile.TextureID == spriteTreeBurntStump01) {
//
//				tile.EffectID = 200
//				log.Printf("INFO: Spalone drzewo ścięte w (%d, %d)", i, j)
//			}
//
//			// Koniec ognia - popiół
//			if tile.EffectID == 70 {
//				tile.EffectID = 3
//				if tile.TextureID > spriteTreeStumpStart && tile.TextureID < spriteTreeStump03 {
//					tile.TextureID = spriteTreeBurntStump00
//				}
//				if tile.TextureID >= spriteTreeStump03 && tile.TextureID <= spriteTreeStumpEnd {
//					tile.TextureID = spriteTreeBurntStump01
//				}
//			}
//
//			// Drzewo upadło na zwykły teren
//			if tile.EffectID == 190 && tile.TextureID < 112 {
//				tile.EffectID = 1
//				tile.Unit = nil // Drzewo znika jako przeszkoda?
//				tile.Building = nil
//				tile.TextureID = spriteGadget00 // Pieniek, chyba
//
//				// Efekt upadku obok
//				if i > 0 && bs.Board.Tiles[i-1][j].TextureID < spriteWaterEnd {
//					bs.Board.Tiles[i-1][j].EffectID = spriteWaterMiddle
//				}
//			}
//
//			// Drzewo upadło na las
//			// To chyba jest dodanie efektu ścinania drzewa
//			// @todo: całość nigdy nie wypróbowana, podmieniam textureID, ale pewnie nie działa
//			// @todo: effectID korzysta z czarodziejskiej liczby. Muszę to wywalić!
//			if tile.EffectID == 190 && tile.TextureID > spriteTreeStumpStart {
//				tile.EffectID = 1
//				tile.Unit = nil
//				tile.Building = nil
//				tile.TextureID = spriteTreeFallenDry00 // Kłoda @reminder: dałem byle co, później poprawię
//				// @todo: ogarnij, czy prawidłowo, ale było textureid < 113
//				// czyli, co jeśli na lewo od drzewa jest woda
//				if i > 0 && bs.Board.Tiles[i-1][j].TextureID < spriteWaterEnd {
//					// @todo: było 110, dałem środek wody
//					bs.Board.Tiles[i-1][j].EffectID = spriteWaterMiddle
//				}
//			}
//		}
//	}
//}

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
		// Sprawdzamy, czy w kapliczce ktoś jest
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

// manaRegen odpowiada za odnawianie many
func manaRegen(bs *battleState) {
	// 0. 15 Hz dla 60 klatek na sekundę, aby zachować zgodność z pierwowzorem
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
