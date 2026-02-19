package main

// battle.go

import (
	"fmt"
	"log"
)

// ============================================================================
// WARUNKI KOŃCA POZIOMU
// ============================================================================

// checkEndConditions sprawdza warunki zakończenia gry.
// @todo: jeszcze w ogóle nie sprawdzane i nie testowane!
func checkEndConditions(bs *battleState) {
	// Jeżeli już ustalono, że gra zakończona, to nie ma co dalej liczyć
	if bs.BattleOutcome != 0 {
		return
	}

	switch bs.CampaignData.EndCondition {
	case endKillAll, endNothing:
		// Zniszczyć i zabić wszystko LUB?
		// Nie wiem, czy jest end_nothing
		// @todo sprawdź w plikach CPP, co to jest
		allEnemiesDead := true
		playerStillAlive := false

		// Sprawdź jednostki gracza
		for _, unit := range bs.Units {
			if unit.Exists && unit.Owner == bs.PlayerID {
				playerStillAlive = true
				break
			}
		}
		// Sprawdź budynki gracza
		if !playerStillAlive {
			for _, bld := range bs.Buildings {
				if bld.Exists && bld.Owner == bs.PlayerID && bld.Type == buildingMain { //&& !bld.IsUnderConstruction {
					playerStillAlive = true
					break
				}
			}
		}

		if bs.CampaignData.EndCondition == endKillAll {
			// Sprawdź jednostki wroga
			for _, unit := range bs.Units {
				if unit.Exists && unit.Owner == bs.AIPlayerID {
					allEnemiesDead = false
					break
				}
			}
			// Sprawdź budynki wroga
			if allEnemiesDead {
				for _, bld := range bs.Buildings {
					if bld.Exists && bld.Owner == bs.AIPlayerID && bld.Type == buildingMain {
						allEnemiesDead = false
						break
					}
				}
			}
		} else {
			allEnemiesDead = false
		}

		if allEnemiesDead {
			bs.BattleOutcome = outcomeVictory
			// @todo: wróć do poprawienia tego warunku
			log.Println("WARUNEK ZAKOŃCZENIA: Zniszczono wszystkich wrogów. Zwycięstwo!")
			bs.CurrentMessage.Text = "Zwycięstwo!"
			bs.CurrentMessage.Duration = 120
		} else if !playerStillAlive {
			// @todo: wróć do poprawienia tego warunku
			bs.BattleOutcome = outcomeDefeat
			log.Println("WARUNEK ZAKOŃCZENIA: Gracz przegrał (brak jednostek/głównego budynku). Porażka!")
			bs.CurrentMessage.Text = "Porażka!"
			bs.CurrentMessage.Duration = 120
		}
		return

	case endRescue:
		// Cel do uratowania
		targetKilled := false
		rx, ry := bs.CampaignData.RescueTargetX, bs.CampaignData.RescueTargetY

		if rx != 0 || ry != 0 {
			// ZMIANA: Używamy nowej struktury Tiles
			tile := &bs.Board.Tiles[rx][ry]

			// Sprawdzamy widoczność i obecność jednostki bezpośrednio
			if tile.Visibility != visibilityUnexplored { // Dawne PlcFogOfWar != 0
				// Jeśli na polu nie ma jednostki (nil) lub jednostka nie należy do AI (czyli zginęła/została przejęta?)
				// w oryginale sprawdzano czy ID==0 lub Owner!=AI
				if tile.Unit == nil || tile.Unit.Owner != bs.AIPlayerID {
					targetKilled = true
				}
			}
		}

		// Czy dotarto do punktu zwycięstwa
		rescueAchieved := false
		for _, unit := range bs.Units {
			if unit.Exists && unit.Owner == bs.PlayerID {
				if unit.X == bs.CampaignData.VictoryPointX && unit.Y == bs.CampaignData.VictoryPointY {
					rescueAchieved = true
					break
				}
			}
		}

		// Sprawdzenie czy gracz żyje (kopiowane z góry dla bezpieczeństwa)
		playerStillAlive := false
		for _, unit := range bs.Units {
			if unit.Exists && unit.Owner == bs.PlayerID {
				playerStillAlive = true
				break
			}
		}
		if !playerStillAlive {
			for _, bld := range bs.Buildings {
				if bld.Exists && bld.Owner == bs.PlayerID && bld.Type == buildingMain {
					playerStillAlive = true
					break
				}
			}
		}

		if rescueAchieved {
			bs.BattleOutcome = outcomeVictory
			// @todo: wróć do poprawienia tego warunku
			log.Println("WARUNEK ZAKOŃCZENIA: Cel uratowany. Zwycięstwo!")
			bs.CurrentMessage.Text = "Uratowano!"
			bs.CurrentMessage.Duration = 120
		} else if targetKilled {
			bs.BattleOutcome = outcomeDefeat
			// @todo: wróć do poprawienia tego warunku
			log.Println("WARUNEK ZAKOŃCZENIA: Cel ratunkowy został zabity. Porażka!")
			bs.CurrentMessage.Text = "Cel zabity!"
			bs.CurrentMessage.Duration = 120
		} else if !playerStillAlive {
			bs.BattleOutcome = outcomeDefeat
			// @todo: wróć do poprawienia tego warunku
			log.Println("WARUNEK ZAKOŃCZENIA: Gracz przegrał. Porażka!")
			bs.CurrentMessage.Text = "Porażka!"
			bs.CurrentMessage.Duration = 120
		}
		return

	case endKillOne:
		targetKilled := true
		// Szukamy konkretnego dowódcy po ID (to akurat zostaje, bo szukamy w liście Units)
		commanderUnit, ok := getUnitByID(1, bs)
		if ok && commanderUnit.Exists {
			targetKilled = false
		}

		playerStillAlive := false
		for _, unit := range bs.Units {
			if unit.Exists && unit.Owner == bs.PlayerID {
				playerStillAlive = true
				break
			}
		}
		if !playerStillAlive {
			for _, bld := range bs.Buildings {
				if bld.Exists && bld.Owner == bs.PlayerID && bld.Type == buildingMain {
					playerStillAlive = true
					break
				}
			}
		}

		if targetKilled {
			bs.BattleOutcome = outcomeVictory

			log.Println("WARUNEK ZAKOŃCZENIA: Cel zabity. Zwycięstwo!")
			// @todo: wróć do poprawienia tego warunku
			bs.CurrentMessage.Text = "Cel zabity!"
			bs.CurrentMessage.Duration = 120
		} else if !playerStillAlive {
			bs.BattleOutcome = outcomeDefeat
			// @todo: wróć do poprawienia tego warunku
			log.Println("WARUNEK ZAKOŃCZENIA: Gracz przegrał. Porażka!")
			bs.CurrentMessage.Text = "Porażka!"
			bs.CurrentMessage.Duration = 120
		}
		return

	case endBuild:
		// Bitwa w rodzaju „wybuduj x budynków, aby wygrać”
		// w plikach JSON jest to „endType” 2; „targetType” określa
		// ile budynków należy posiadać aby wygrać
		// Przykładowo level_15.json wymaga łącznie 7 budowli
		// To 1xMAIN + 2xBarn + 1xBARRACKS+ 3 nowe budowle żeby wygrać
		buildingCount := uint8(0)
		for _, bld := range bs.Buildings {
			// Nie zniszczone budynki gracza, nieukończone budowle się nie wliczają!
			if bld.Exists && bld.Owner == bs.PlayerID && !bld.IsUnderConstruction {
				buildingCount++
			}
		}

		// Sprawdzamy, czy gracz jeszcze żyje
		// Jest możliwe, że „żyje”, ale nie może wygrać
		// w tej chwili się tym nie zajmuję.
		// @todo: sama jedna krowa bez budynków nie pozwoli wygrać
		// z pastuchem to samo!
		playerStillAlive := false
		for _, unit := range bs.Units {
			if unit.Exists && unit.Owner == bs.PlayerID {
				playerStillAlive = true
				// Wystarczy jedna żywa jednostka, aby gracz „wciąż żył”
				break
			}
		}
		// Jeżeli nie natrafiliśmy nawet na jedną jednostkę
		// to sprawdzamy, czy jest chociaż jeden budynek
		// @todo: wróć do poprawienia tego warunku
		// sama jedna obora bez innych budynków nie pozwala na wygraną!
		if !playerStillAlive {
			for _, bld := range bs.Buildings {
				if bld.Exists && bld.Owner == bs.PlayerID && bld.Type == buildingMain {
					playerStillAlive = true
					break
				}
			}
		}

		// Czy mamy dość ukończonych budynków, aby zakończyć bitwę?
		if buildingCount >= bs.CampaignData.TargetType {
			bs.BattleOutcome = outcomeVictory
			// @todo: usuń po zweryfikowaniu, czy działa poprawnie
			log.Println("WARUNEK ZAKOŃCZENIA: Wymagana liczba budynków zbudowana. Zwycięstwo!")
			bs.CurrentMessage.Text = "Zbudowano!"
			bs.CurrentMessage.Duration = 30
		} else if !playerStillAlive {
			bs.BattleOutcome = outcomeDefeat
			// @todo: usuń po zweryfikowaniu, czy działa poprawnie
			log.Println("WARUNEK ZAKOŃCZENIA: Gracz przegrał. Porażka!")
			bs.CurrentMessage.Text = "Porażka!"
			bs.CurrentMessage.Duration = 30
		}
		return
	}
}

// GŁÓWNA PĘTLA GRY NA POZIOMIE KADRU

// @todo: czy nie lepiej używać ticker time.NewTicker?
func updateGame(bs *battleState) {
	// 1. Odświeżanie logiki co klatkowej
	updatePerFrameLogic(bs)

	// 2. Tick Logiki
	logicTickInterval := bs.GameSpeed + uint16(1)
	if bs.GlobalFrameCounter%logicTickInterval != 0 {
		updateMessageTimer(bs)
		return
	}

	// === LOGIKA TURA ===

	// 3. Czyszczenie pamięci
	if bs.GlobalFrameCounter%100 == 0 {
		performPeriodicCleanup(bs)
	}

	// 4. Sprawdzanie warunki zakończenia bitwy
	// @todo: ogarnij, bo jeszcze nie ruszone w ogóle
	checkEndConditions(bs)

	if bs.BattleOutcome != 0 {
		bs.QuitLevel = true
		log.Printf("Gra zakończona! Wynik: %d. Poziom do zamknięcia: %v", bs.BattleOutcome, bs.QuitLevel)

		return
	}

	// 5. SI
	// @todo: ogarnij, bo jeszcze nie ruszone w ogóle
	processAI(bs)

	// 6. Komendy
	// @todo: ogarnij, bo to CHYBA jeszcze nie ruszone było
	processCommands(bs)

	// 7. Odświeżenie jednostek
	updateUnits(bs)

	// 8. Odświeżenie pocisków
	updateProjectiles(bs)

	// 9. Budynki i niszczenie ich
	// @todo: to brzmi, jak coś co powinno się rozdzielić
	// ogarnij, czy da się tego potwora uprościć, ulepszyć
	updateBuildings(bs)

	// 10. Efekty globalne
	applyGlobalEffects(bs)
	updateCorpses(bs) // Pasuje tutaj, ale nie chcę jeszcze wrzucać w applyGlobalEffects

	// 11. Skrypty poziomu (przemiana w niedźwiedzia itd)
	// @todo: do zrobienia, nie sprawdzałem jeszcze czy w ogóle działają
	handleLevelEvents(bs)

	// 12. Wiadomość
	// @todo: nie wiem nawet po co to jest!
	updateMessageTimer(bs)
}

func updatePerFrameLogic(bs *battleState) {
	bs.pathfindingUnitsThisTick = 0
	bs.enemyCacheUpdateTick = 0

	if bs.CheatsEnabled {
		log.Printf("DEBUG GAME: Level %d, Difficulty %d, GameState %d", bs.CurrentLevel, bs.DifficultyLevel, bs.BattleOutcome)
	}

	// Animacje
	if bs.GlobalFrameCounter%envAnimSpeedDivisor == 0 {
		updateAnimationCounters(bs)
	}

	updateWorldTimers(bs)
}

// Nie pamiętam po co to w ogóle jest potrzebne
func updateMessageTimer(bs *battleState) {
	if bs.CurrentMessage.Duration > 0 {
		bs.CurrentMessage.Duration--
		if bs.CurrentMessage.Duration == 0 {
			bs.CurrentMessage.Text = ""
		}
	}
}

// performPeriodicCleanup sprawdza, czy w budynku znajdują się
// zabite jednostki i zwalnia im miejsce
func performPeriodicCleanup(bs *battleState) {
	for _, building := range bs.Buildings {
		if building != nil && building.Exists {
			building.cleanupDeadUnits(bs)
		}
	}
}

func updateCorpses(bs *battleState) {
	// nextFreeIndex wskaźnik do miejsca dla nowych zwłok
	nextFreeIndex := 0

	for scanIndex := range bs.Corpses {
		corpse := &bs.Corpses[scanIndex]

		corpse.DecayTimer--

		timePassed := corpseDecayTime - corpse.DecayTimer
		currentPhase := uint8(timePassed / corpsesPhaseDuration)

		if currentPhase > corpsesPhase2 {
			currentPhase = corpsesPhase2
		}

		corpse.Phase = currentPhase

		const corpseFadeDuration = 400

		if corpse.DecayTimer > corpseFadeDuration {
			corpse.Alpha = 255
		} else {
			corpse.Alpha = uint8(float32(corpse.DecayTimer) / float32(corpseFadeDuration) * corpsesMaxAlpha)
		}

		// Zwłoki jeszcze się nie rozłożyły całkowicie
		if corpse.DecayTimer > 0 {
			// Jeśli zwłoki są pod adresem większym niż wolny adres, to przenosimy.
			if scanIndex != nextFreeIndex {
				bs.Corpses[nextFreeIndex] = *corpse
			}

			nextFreeIndex++
		}
	}

	bs.Corpses = bs.Corpses[:nextFreeIndex]
}

// processAI zarządza sztuczną inteligencją przeciwnika
func processAI(bs *battleState) {
	switch bs.CampaignData.DecisionType {
	case boardVillage:
		aiMakeDecision(bs)
	case boardBattleDyn:
		aiMakeDecision(bs)
	case boardBattleStat:
		aiMakeDecision(bs)
	case boardNothing:
	default:
		log.Printf("OSTRZEŻENIE: Nieznany typ decyzji planszy: %d", bs.CampaignData.DecisionType)
	}
}

// processCommands przetwarza rozkazy
func processCommands(bs *battleState) {
	if bs.CurrentCommands[0].ActionType != cmdIdle {
		log.Printf("ProcComm: Gracz wykonuje komendę. TargetObject: %d, ActionType: %d",
			bs.CurrentCommands[0].CommandCategory, bs.CurrentCommands[0].ActionType)
		bs.HumanPlayerState.setCommand(&bs.CurrentCommands[0], bs)
		bs.CurrentCommands[0].ActionType = cmdIdle
	}
	if bs.CurrentCommands[1].ActionType != cmdIdle {
		log.Printf("ProcComm: SI wykonuje komendę. TargetObject: %d, ActionType: %d",
			bs.CurrentCommands[1].CommandCategory, bs.CurrentCommands[1].ActionType)
		bs.AIEnemyState.setCommand(&bs.CurrentCommands[1], bs)
		bs.CurrentCommands[1].ActionType = cmdIdle
	}
}

// odświeża jednostki, sprawdza, która została zabita itd
func updateUnits(bs *battleState) {
	for _, unit := range bs.Units {
		if unit.Exists {
			unit.updateUnit(bs)
		}
	}

	cleanupDeadUnits(bs)
}

// updateProjectiles odświeża pociski, sprawdza, czy dodać nowe
func updateProjectiles(bs *battleState) {
	activeProjectiles := bs.Projectiles[:0]
	for _, p := range bs.Projectiles {
		if p.Exists {
			p.updateProjectile(bs)
			if p.Exists {
				activeProjectiles = append(activeProjectiles, p)
			}
		}
	}
	bs.Projectiles = activeProjectiles
}

// updateBuildings TODO: oczyść, ogarnij!
func updateBuildings(bs *battleState) {
	for _, bld := range bs.Buildings {
		if !bld.Exists {
			continue
		}

		if bld.AccumulatedDamage > 0 {
			finalDamage := bld.AccumulatedDamage - uint16(bld.Armor)
			if finalDamage > 0 {
				// Bez tego bld.HP przekręca się na 65 tys.
				if bld.HP >= finalDamage {
					bld.HP -= finalDamage
				} else {
					bld.HP = 0
				}

				log.Printf("building %d took %d final damage. HP: %d/%d",
					bld.ID, finalDamage, bld.HP, bld.MaxHP)

				if bld.HP <= 0 {
					bld.HP = 0
					bld.Exists = false
					log.Printf("building %d destroyed!", bld.ID)

					switch bld.Owner {
					case bs.HumanPlayerState.PlayerID:
						bs.HumanPlayerState.CurrentBuildings--
					case bs.AIEnemyState.PlayerID:
						bs.AIEnemyState.CurrentBuildings--
					}

					placeRuins(bs, bld)

					for _, tile := range bld.OccupiedTiles {
						if tile.X < boardMaxX && tile.Y < boardMaxY {
							// Usuwamy odnośnik do budynku z kafelka
							if bs.Board.Tiles[tile.X][tile.Y].Building == bld {
								bs.Board.Tiles[tile.X][tile.Y].Building = nil
								bs.Board.Tiles[tile.X][tile.Y].IsWalkable = true
							}
						}
					}
				}
			}

			bld.AccumulatedDamage = 0
		}
	}
}

func applyGlobalEffects(bs *battleState) {
	// clearAttack(bs)
	healingShrine(bs)
	manaRegen(bs)
}

// handleLevelEvents przemiana w niedźwiedzia, odprowadzenie jednostki do punktu ucieczki
// @todo: nie ogarnięte w ogóle! przemiana powinna być jak healingShires
// czemu sprawdzenie miejsca ucieczki jest na w tym miejscu a nie przy
// sprawdzaniu warunków końca bitwy?!
func handleLevelEvents(bs *battleState) {
	// Przemiana w niedźwiedzia
	if bs.CampaignData.TransformationSiteX != 0 || bs.CampaignData.TransformationSiteY != 0 {
		tx, ty := bs.CampaignData.TransformationSiteX, bs.CampaignData.TransformationSiteY

		tile := &bs.Board.Tiles[tx][ty]

		if tile.TextureID == spriteEffectTransform00 && tile.Unit != nil {
			unit := tile.Unit // Bezpośredni wskaźnik!
			if unit.Exists && unit.Owner == bs.PlayerID {
				// @todo: Ten warunek kompletnie nie ma sensu! Każdy kto wejdzie w
				// miejsce przemiany zamienia się w niedźwiedzia, a nie tylko AXEMAN LVL50
				// KONIECZNIE SPRAWDŹ TO w PIERWOTNYM KODZIE I OGARNIJ O CO CHODZI!!!
				if unit.Type == unitAxeman && unit.Experience >= 50 {
					log.Printf("GAME: Transformacja jednostki %d!", unit.ID)
					unit.Type = unitBear
					// unit.SightRange = UNIT_TYPE_DATA[UNIT_BEAR][0]
					// unit.AttackRange = UNIT_TYPE_DATA[UNIT_BEAR][1]
					// unit.Damage = UNIT_TYPE_DATA[UNIT_BEAR][2]
					// unit.Armor = UNIT_TYPE_DATA[UNIT_BEAR][3]
					// unit.MaxHP = UNIT_TYPE_DATA[UNIT_BEAR][4]
					// unit.MaxDelay = UNIT_TYPE_DATA[UNIT_BEAR][5]
					stats, ok := unitDefs[unitBear]
					if !ok {
						panic(fmt.Sprintf("BŁĄD KRYTYCZNY: Nie udało się przemienić jednostki ID%d w UNIT_BEAR", unit.ID))
					}
					unit.SightRange = stats.SightRange
					unit.AttackRange = stats.AttackRange
					unit.Damage = stats.BaseDamage
					unit.Armor = stats.BaseArmor
					unit.MaxHP = stats.MaxHP
					unit.MaxDelay = stats.MoveDelay
					unit.Mana = stats.MaxMana
					if unit.HP > stats.MaxHP {
						unit.HP = stats.MaxHP
					}
				}
			}
		}
	}

	// RATUNEK
	if bs.CampaignData.EndCondition == endRescue &&
		bs.CampaignData.RescueTargetX != 0 && bs.CampaignData.RescueTargetY != 0 {

		rx, ry := bs.CampaignData.RescueTargetX, bs.CampaignData.RescueTargetY
		tile := &bs.Board.Tiles[rx][ry]

		// ZMIANA: Visibility i unit pointer
		if tile.Visibility != visibilityUnexplored {
			if tile.Unit != nil {
				rescuedUnit := tile.Unit
				if rescuedUnit.Exists &&
					rescuedUnit.Owner == bs.PlayerID &&
					rescuedUnit.Type == unitType(bs.CampaignData.TargetType) {
					log.Printf("RESCUE: cel jest w punkcie!")
				}
			}
		}
	}
}

// =======================
func placeRuins(bs *battleState, bld *building) {
	if len(bld.OccupiedTiles) == 0 {
		return
	}

	const graphicsRuinsID = 257

	minX, minY := bld.OccupiedTiles[0].X, bld.OccupiedTiles[0].Y
	maxX, maxY := bld.OccupiedTiles[0].X, bld.OccupiedTiles[0].Y

	for _, pt := range bld.OccupiedTiles {
		if pt.X < minX {
			minX = pt.X
		}

		if pt.Y < minY {
			minY = pt.Y
		}

		if pt.X > maxX {
			maxX = pt.X
		}

		if pt.Y > maxY {
			maxY = pt.Y
		}
	}

	width := maxX - minX + 1

	for _, pt := range bld.OccupiedTiles {
		x, y := pt.X, pt.Y
		if x >= boardMaxX || y >= boardMaxY {
			continue
		}

		tile := &bs.Board.Tiles[x][y]

		// Czyścimy wskaźnik na budynek
		tile.Building = nil

		// Ustawiamy grafikę ruin
		if bld.Type == buildingPalisade {
			tile.TextureID = spritePalisadeDestroyed
			tile.IsWalkable = true
		} else {
			dx := pt.X - minX
			dy := pt.Y - minY
			idx := dy*width + dx
			tile.TextureID = graphicsRuinsID + uint16(idx)
			// Zgliszcza uniemożliwiają ruch
			tile.IsWalkable = false
		}
	}
}

// Usuwa uśmiercone jednostki z bs.
// Nie mylić z logiką rozkładu zwłok updateCorpses
func cleanupDeadUnits(bs *battleState) {
	if bs.GlobalFrameCounter%6000 != 0 {
		return
	}
	if len(bs.Units) < 50 {
		return
	}

	log.Println("INFO: Rozpoczynam czyszczenie pamięci...")

	newUnitsList := make([]*unit, 0, len(bs.Units))
	for _, u := range bs.Units {
		if u.Exists {
			newUnitsList = append(newUnitsList, u)
		}
	}

	removedCount := len(bs.Units) - len(newUnitsList)
	bs.Units = newUnitsList

	if removedCount > 0 {
		log.Printf("INFO: Wyczyszczono %d martwych jednostek.", removedCount)
	}
}
