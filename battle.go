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
func checkEndConditions(bState *battleState) {
	// Jeżeli już ustalono, że gra zakończona, to nie ma co dalej liczyć
	if bState.BattleOutcome != 0 {
		return
	}

	switch bState.CampaignData.EndCondition {
	case endKillAll, endNothing:
		// Zniszczyć i zabić wszystko LUB?
		// Nie wiem, czy jest end_nothing
		// @todo sprawdź w plikach CPP, co to jest
		allEnemiesDead := true
		playerStillAlive := false

		// Sprawdź jednostki gracza
		for _, u := range bState.Units {
			if u.Exists && u.Owner == bState.PlayerID {
				playerStillAlive = true

				break
			}
		}
		// Sprawdź budynki gracza
		if !playerStillAlive {
			for _, bld := range bState.Buildings {
				if bld.Exists && bld.Owner == bState.PlayerID && bld.Type == buildingMain {
					playerStillAlive = true

					break
				}
			}
		}

		if bState.CampaignData.EndCondition == endKillAll {
			// Sprawdź jednostki wroga
			for _, u := range bState.Units {
				if u.Exists && u.Owner == bState.AIPlayerID {
					allEnemiesDead = false

					break
				}
			}
			// Sprawdź budynki wroga
			if allEnemiesDead {
				for _, bld := range bState.Buildings {
					if bld.Exists && bld.Owner == bState.AIPlayerID && bld.Type == buildingMain {
						allEnemiesDead = false

						break
					}
				}
			}
		} else {
			allEnemiesDead = false
		}

		if allEnemiesDead {
			bState.BattleOutcome = outcomeVictory
			// @todo: wróć do poprawienia tego warunku
			log.Println("WARUNEK ZAKOŃCZENIA: Zniszczono wszystkich wrogów. Zwycięstwo!")
			bState.CurrentMessage.Text = "Zwycięstwo!"
			bState.CurrentMessage.Duration = 120
		} else if !playerStillAlive {
			// @todo: wróć do poprawienia tego warunku
			bState.BattleOutcome = outcomeDefeat
			log.Println("WARUNEK ZAKOŃCZENIA: Gracz przegrał (brak jednostek/głównego budynku). Porażka!")
			bState.CurrentMessage.Text = "Porażka!"
			bState.CurrentMessage.Duration = 120
		}

		return

	case endRescue:
		// Cel do uratowania
		targetKilled := false
		rx, ry := bState.CampaignData.RescueTargetX, bState.CampaignData.RescueTargetY

		// @todo: czemu to do cholery jest 0,0?
		if rx != 0 || ry != 0 {
			t := &bState.Board.Tiles[rx][ry]

			// Sprawdzamy widoczność i obecność jednostki bezpośrednio
			if t.Visibility != visibilityUnexplored { // Dawne PlcFogOfWar != 0
				// Jeśli na polu nie ma jednostki (nil) lub jednostka nie należy do AI (czyli zginęła/została przejęta?)
				// w oryginale sprawdzano czy ID==0 lub Owner!=AI
				if t.Unit == nil || t.Unit.Owner != bState.AIPlayerID {
					targetKilled = true
				}
			}
		}

		// Czy dotarto do punktu zwycięstwa
		rescueAchieved := false
		for _, u := range bState.Units {
			if u.Exists && u.Owner == bState.PlayerID {
				if u.X == bState.CampaignData.VictoryPointX && u.Y == bState.CampaignData.VictoryPointY {
					rescueAchieved = true
					break
				}
			}
		}

		// Sprawdzenie czy gracz żyje (kopiowane z góry dla bezpieczeństwa)
		playerStillAlive := false

		for _, u := range bState.Units {
			if u.Exists && u.Owner == bState.PlayerID {
				playerStillAlive = true

				break
			}
		}
		if !playerStillAlive {
			for _, bld := range bState.Buildings {
				if bld.Exists && bld.Owner == bState.PlayerID && bld.Type == buildingMain {
					playerStillAlive = true

					break
				}
			}
		}

		if rescueAchieved {
			bState.BattleOutcome = outcomeVictory

			// @todo: wróć do poprawienia tego warunku
			log.Println("WARUNEK ZAKOŃCZENIA: Cel uratowany. Zwycięstwo!")

			bState.CurrentMessage.Text = "Uratowano!"
			bState.CurrentMessage.Duration = 120
		} else if targetKilled {
			bState.BattleOutcome = outcomeDefeat

			// @todo: wróć do poprawienia tego warunku
			log.Println("WARUNEK ZAKOŃCZENIA: Cel ratunkowy został zabity. Porażka!")

			bState.CurrentMessage.Text = "Cel zabity!"
			bState.CurrentMessage.Duration = 120
		} else if !playerStillAlive {
			bState.BattleOutcome = outcomeDefeat

			// @todo: wróć do poprawienia tego warunku
			log.Println("WARUNEK ZAKOŃCZENIA: Gracz przegrał. Porażka!")

			bState.CurrentMessage.Text = "Porażka!"
			bState.CurrentMessage.Duration = 120
		}

		return

	case endKillOne:
		targetKilled := true
		// Szukamy konkretnego dowódcy po ID (to akurat zostaje, bo szukamy w liście Units)
		commanderUnit, ok := bState.getUnitByID(1) // @todo: @reminder: to nie może być sztywne id!
		if ok && commanderUnit.Exists {
			targetKilled = false
		}

		playerStillAlive := false
		for _, u := range bState.Units {
			if u.Exists && u.Owner == bState.PlayerID {
				playerStillAlive = true

				break
			}
		}

		if !playerStillAlive {
			for _, bld := range bState.Buildings {
				if bld.Exists && bld.Owner == bState.PlayerID && bld.Type == buildingMain {
					playerStillAlive = true

					break
				}
			}
		}

		if targetKilled {
			bState.BattleOutcome = outcomeVictory

			log.Println("WARUNEK ZAKOŃCZENIA: Cel zabity. Zwycięstwo!")

			// @todo: wróć do poprawienia tego warunku
			bState.CurrentMessage.Text = "Cel zabity!"
			bState.CurrentMessage.Duration = 120
		} else if !playerStillAlive {
			bState.BattleOutcome = outcomeDefeat

			// @todo: wróć do poprawienia tego warunku
			log.Println("WARUNEK ZAKOŃCZENIA: Gracz przegrał. Porażka!")

			bState.CurrentMessage.Text = "Porażka!"
			bState.CurrentMessage.Duration = 120
		}

		return

	case endBuild:
		// Bitwa w rodzaju „wybuduj x budynków, aby wygrać”
		// w plikach JSON jest to „endType” 2; „targetType” określa
		// ile budynków należy posiadać aby wygrać
		// Przykładowo level_15.json wymaga łącznie 7 budowli
		// To 1xMAIN + 2xBarn + 1xBARRACKS+ 3 nowe budowle żeby wygrać
		buildingCount := uint8(0)

		for _, bld := range bState.Buildings {
			// Nie zniszczone budynki gracza, nieukończone budowle się nie wliczają!
			if bld.Exists && bld.Owner == bState.PlayerID && !bld.IsUnderConstruction {
				buildingCount++
			}
		}

		// Sprawdzamy, czy gracz jeszcze żyje
		// Jest możliwe, że „żyje”, ale nie może wygrać
		// w tej chwili się tym nie zajmuję.
		// @todo: sama jedna krowa bez budynków nie pozwoli wygrać
		// z pastuchem to samo!
		playerStillAlive := false

		for _, u := range bState.Units {
			if u.Exists && u.Owner == bState.PlayerID {
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
			for _, bld := range bState.Buildings {
				if bld.Exists && bld.Owner == bState.PlayerID && bld.Type == buildingMain {
					playerStillAlive = true

					break
				}
			}
		}

		// Czy mamy dość ukończonych budynków, aby zakończyć bitwę?
		if buildingCount >= bState.CampaignData.TargetType {
			bState.BattleOutcome = outcomeVictory

			// @todo: usuń po zweryfikowaniu, czy działa poprawnie
			log.Println("WARUNEK ZAKOŃCZENIA: Wymagana liczba budynków zbudowana. Zwycięstwo!")

			bState.CurrentMessage.Text = "Zbudowano!"
			bState.CurrentMessage.Duration = 30
		} else if !playerStillAlive {
			bState.BattleOutcome = outcomeDefeat

			// @todo: usuń po zweryfikowaniu, czy działa poprawnie
			log.Println("WARUNEK ZAKOŃCZENIA: Gracz przegrał. Porażka!")

			bState.CurrentMessage.Text = "Porażka!"
			bState.CurrentMessage.Duration = 30
		}

		return
	}
}

// GŁÓWNA PĘTLA GRY NA POZIOMIE KADRU

// @todo: czy nie lepiej używać ticker time.NewTicker?
func updateGame(bState *battleState) {
	// 1. Odświeżanie logiki co klatkowej
	updatePerFrameLogic(bState)

	// 2. Tick Logiki
	logicTickInterval := bState.GameSpeed + uint16(1)
	if bState.GlobalFrameCounter%logicTickInterval != 0 {
		updateMessageTimer(bState)

		return
	}

	// === LOGIKA TURA ===

	// 3. Czyszczenie pamięci
	// @todo: zastanów się, czy nie dodać czyszczenia innych list, jak budynki, zwłoki, czy płonące kafelki.
	if bState.GlobalFrameCounter%120 == 0 {
		performPeriodicCleanup(bState)
		updateBurningTilesList(bState)
		updateFallingTreesList(bState)
		updateGhostList(bState)
	}

	// 4. Sprawdzanie warunki zakończenia bitwy
	// @todo: ogarnij, bo jeszcze nie ruszone w ogóle
	checkEndConditions(bState)

	if bState.BattleOutcome != 0 {
		bState.QuitLevel = true
		log.Printf("Gra zakończona! Wynik: %d. Poziom do zamknięcia: %v", bState.BattleOutcome, bState.QuitLevel)

		return
	}

	// 5. SI
	// @todo: ogarnij, bo jeszcze nie ruszone w ogóle
	processAI(bState)

	// 6. Komendy
	// @todo: ogarnij, bo to CHYBA jeszcze nie ruszone było
	processCommands(bState)

	// 7. Odświeżenie jednostek
	updateUnits(bState)

	// 8. Odświeżenie pocisków
	updateProjectiles(bState)

	// 9. Budynki i niszczenie ich
	// @todo: to brzmi, jak coś co powinno się rozdzielić
	// ogarnij, czy da się tego potwora uprościć, ulepszyć
	updateBuildings(bState)

	// 10. Efekty globalne
	applyGlobalEffects(bState)

	// 11. Skrypty poziomu (przemiana w niedźwiedzia itd)
	// @todo: do zrobienia, nie sprawdzałem jeszcze czy w ogóle działają
	handleLevelEvents(bState)

	// 12. Wiadomość
	// @todo: nie wiem nawet po co to jest!
	updateMessageTimer(bState)
}

func updatePerFrameLogic(bState *battleState) {
	bState.pathfindingUnitsThisTick = 0
	bState.enemyCacheUpdateTick = 0

	if bState.CheatsEnabled {
		log.Printf("DEBUG GAME: Level %d, Difficulty %d, GameState %d",
			bState.CurrentLevel, bState.DifficultyLevel, bState.BattleOutcome)
	}

	// Animacje
	if bState.GlobalFrameCounter%envAnimSpeedDivisor == 0 {
		updateAnimationCounters(bState)
	}

	updateWorldTimers(bState)
}

// Nie pamiętam po co to w ogóle jest potrzebne.
func updateMessageTimer(bState *battleState) {
	if bState.CurrentMessage.Duration > 0 {
		bState.CurrentMessage.Duration--
		if bState.CurrentMessage.Duration == 0 {
			bState.CurrentMessage.Text = ""
		}
	}
}

// Sprawdza, czy w budynku znajdują się zabite jednostki i zwalnia im miejsce.
func performPeriodicCleanup(bState *battleState) {
	for _, bld := range bState.Buildings {
		if bld != nil && bld.Exists {
			bld.cleanupDeadUnits(bState)
		}
	}
}

// Sprawdza, czy można usunąć kafelki, które zakończyły cykl płonięcia.
func updateBurningTilesList(bState *battleState) {
	var burningTiles []*tile

	for _, currentTile := range bState.BurningTilesList {
		if currentTile.IsBurning || currentTile.hasAsh {
			burningTiles = append(burningTiles, currentTile)
		}
	}

	bState.BurningTilesList = burningTiles
}

// Sprawdza, czy drzewo dopełniło cykl upadku i można je usunąć z listy.
func updateFallingTreesList(bState *battleState) {
	var fallingTrees []*tile

	for _, currentTile := range bState.FallingTreesList {
		if currentTile.treeState != treeFell {
			fallingTrees = append(fallingTrees, currentTile)
		}
	}

	bState.FallingTreesList = fallingTrees
}

func updateGhostList(bState *battleState) {
	var ghostList []*tile

	for _, ghostTile := range bState.GhostsList {
		if ghostTile.GhostEffect {
			ghostList = append(ghostList, ghostTile)
		}
	}

	bState.GhostsList = ghostList
}

func updateCorpses(bState *battleState) {
	// nextFreeIndex wskaźnik do miejsca dla nowych zwłok
	nextFreeIndex := 0

	for scanIndex := range bState.CorpsesList {
		corpseToUpdate := &bState.CorpsesList[scanIndex]

		corpseToUpdate.DecayTimer--

		timePassed := corpseDecayTime - corpseToUpdate.DecayTimer
		currentPhase := uint8(timePassed / corpsesPhaseDuration)

		if currentPhase > corpsesPhase2 {
			currentPhase = corpsesPhase2
		}

		corpseToUpdate.Phase = currentPhase

		const corpseFadeDuration = 400

		if corpseToUpdate.DecayTimer > corpseFadeDuration {
			corpseToUpdate.Alpha = 255
		} else {
			corpseToUpdate.Alpha = uint8(float32(corpseToUpdate.DecayTimer) / float32(corpseFadeDuration) * corpsesMaxAlpha)
		}

		// Zwłoki jeszcze się nie rozłożyły całkowicie
		if corpseToUpdate.DecayTimer > 0 {
			// Jeśli zwłoki są pod adresem większym niż wolny adres, to przenosimy.
			if scanIndex != nextFreeIndex {
				bState.CorpsesList[nextFreeIndex] = *corpseToUpdate
			}

			nextFreeIndex++
		}
	}

	bState.CorpsesList = bState.CorpsesList[:nextFreeIndex]
}

// processAI zarządza sztuczną inteligencją przeciwnika
func processAI(bState *battleState) {
	switch bState.CampaignData.DecisionType {
	case boardVillage:
		aiMakeDecision(bState)
	case boardBattleDyn:
		aiMakeDecision(bState)
	case boardBattleStat:
		aiMakeDecision(bState)
	case boardNothing:
	default:
		log.Printf("OSTRZEŻENIE: Nieznany typ decyzji planszy: %d", bState.CampaignData.DecisionType)
	}
}

// processCommands przetwarza rozkazy
func processCommands(bState *battleState) {
	if bState.CurrentCommands[0].ActionType != cmdUIdle {
		log.Printf("ProcComm: Gracz wykonuje komendę. ActionType: %d",
			bState.CurrentCommands[0].ActionType)
		bState.HumanPlayerState.setCommand(&bState.CurrentCommands[0], bState)
		bState.CurrentCommands[0].ActionType = cmdUIdle
	}

	if bState.CurrentCommands[1].ActionType != cmdUIdle {
		log.Printf("ProcComm: SI wykonuje komendę. ActionType: %d",
			bState.CurrentCommands[1].ActionType)
		bState.AIEnemyState.setCommand(&bState.CurrentCommands[1], bState)
		bState.CurrentCommands[1].ActionType = cmdUIdle
	}
}

// odświeża jednostki, sprawdza, która została zabita itd.
func updateUnits(bState *battleState) {
	bState.PathfindingBudget = 0

	for _, currentUnit := range bState.Units {
		if currentUnit.Exists {
			currentUnit.updateUnit(bState)
		}
	}

	cleanupDeadUnits(bState)
}

// updateProjectiles odświeża pociski, sprawdza, czy dodać nowe.
func updateProjectiles(bState *battleState) {
	activeProjectiles := bState.Projectiles[:0]
	for _, p := range bState.Projectiles {
		if p.Exists {
			p.updateProjectile(bState)

			if p.Exists {
				activeProjectiles = append(activeProjectiles, p)
			}
		}
	}

	bState.Projectiles = activeProjectiles
}

func updateBuildings(bState *battleState) {
	for _, bld := range bState.Buildings {
		if !bld.Exists {
			continue
		}

		if bld.AccumulatedDamage > uint16(bld.Armor) {
			finalDamage := bld.AccumulatedDamage - uint16(bld.Armor)

			if finalDamage > 0 {
				applyBuildingDamage(bld, finalDamage, bState)
			}

			bld.AccumulatedDamage = 0
		}
	}

	cleanupDestroyedBuildings(bState)
}

// @todo: CZEMU TO JEST W battle.go a nie constructions.go?! 26.05.2026
func applyBuildingDamage(bld *building, finalDamage uint16, bState *battleState) {
	// Bez tego bld.HP przekręca się na 65 tys.
	if bld.HP >= finalDamage {
		bld.HP -= finalDamage
	} else {
		bld.HP = 0
	}

	log.Printf("building %d took %d final damage. HP: %d/%d",
		bld.ID, finalDamage, bld.HP, bld.MaxHP)

	if bld.HP <= 0 {
		removeBuilding(bld, bState)
	}
}

func removeBuilding(bld *building, bState *battleState) {
	bld.HP = 0

	if bld.Type != buildingPalisade {
		placeDestroyedBuilding(bld, bState)
	} else {
		placeDestroyedPalisade(bld.OccupiedTiles[0].X, bld.OccupiedTiles[0].Y, bld, bState)
	}
}

func placeDestroyedBuilding(bld *building, bState *battleState) {
	bld.Exists = false
	log.Printf("building %d destroyed!", bld.ID)

	switch bld.Owner {
	case bState.HumanPlayerState.PlayerID:
		bState.HumanPlayerState.CurrentBuildings--
	case bState.AIEnemyState.PlayerID:
		bState.AIEnemyState.CurrentBuildings--
	}

	placeRuins(bState, bld)

	for _, occupiedTile := range bld.OccupiedTiles {
		if occupiedTile.X < boardMaxX && occupiedTile.Y < boardMaxY {
			// Usuwamy odnośnik do budynku z kafelka
			checkTile := bState.Board.Tiles[occupiedTile.X][occupiedTile.Y]

			if checkTile.Building == bld {
				checkTile.Building = nil
				checkTile.IsWalkable = true
			}
		}
	}
}

func applyGlobalEffects(bState *battleState) {
	healingShrine(bState)
	manaRegen(bState)
	updateCorpses(bState)
	burningTileEffect(bState)
	fallingTreeEffect(bState)
	ghostEffect(bState)
	handleMagicShieldEffect(bState)
}

// handleLevelEvents przemiana w niedźwiedzia, odprowadzenie jednostki do punktu ucieczki
// @todo: nie ogarnięte w ogóle! przemiana powinna być jak healingShires
// czemu sprawdzenie miejsca ucieczki jest na w tym miejscu a nie przy
// sprawdzaniu warunków końca bitwy?!
func handleLevelEvents(bState *battleState) {
	// Przemiana w niedźwiedzia
	if bState.CampaignData.TransformationSiteX != 0 || bState.CampaignData.TransformationSiteY != 0 {
		tx, ty := bState.CampaignData.TransformationSiteX, bState.CampaignData.TransformationSiteY

		currentTile := &bState.Board.Tiles[tx][ty]

		if currentTile.TextureID == spriteEffectTransform00 && currentTile.Unit != nil {
			currentUnit := currentTile.Unit // Bezpośredni wskaźnik!
			if currentUnit.Exists && currentUnit.Owner == bState.PlayerID {
				// @todo: Ten warunek kompletnie nie ma sensu! Każdy kto wejdzie w
				// miejsce przemiany zamienia się w niedźwiedzia, a nie tylko AXEMAN LVL50
				// KONIECZNIE SPRAWDŹ TO w PIERWOTNYM KODZIE I OGARNIJ O CO CHODZI!!!
				if currentUnit.Type == unitAxeman && currentUnit.Experience >= 50 {
					log.Printf("GAME: Przemiana jednostki %d!", currentUnit.ID)
					currentUnit.Type = unitBear
					// unit.SightRange = UNIT_TYPE_DATA[UNIT_BEAR][0]
					// unit.AttackRange = UNIT_TYPE_DATA[UNIT_BEAR][1]
					// unit.Damage = UNIT_TYPE_DATA[UNIT_BEAR][2]
					// unit.Armor = UNIT_TYPE_DATA[UNIT_BEAR][3]
					// unit.MaxHP = UNIT_TYPE_DATA[UNIT_BEAR][4]
					// unit.MaxDelay = UNIT_TYPE_DATA[UNIT_BEAR][5]
					stats, ok := unitDefs[unitBear]
					if !ok {
						panic(fmt.Sprintf("BŁĄD KRYTYCZNY: Nie udało się przemienić jednostki ID%d w UNIT_BEAR", currentUnit.ID))
					}

					currentUnit.SightRange = stats.SightRange
					currentUnit.AttackRange = stats.AttackRange
					currentUnit.Damage = stats.BaseDamage
					currentUnit.Armor = stats.BaseArmor
					currentUnit.MaxHP = stats.MaxHP
					currentUnit.MaxDelay = stats.MoveDelay
					currentUnit.Mana = stats.MaxMana

					if currentUnit.HP > stats.MaxHP {
						currentUnit.HP = stats.MaxHP
					}
				}
			}
		}
	}

	// RATUNEK
	if bState.CampaignData.EndCondition == endRescue &&
		bState.CampaignData.RescueTargetX != 0 && bState.CampaignData.RescueTargetY != 0 {

		rx, ry := bState.CampaignData.RescueTargetX, bState.CampaignData.RescueTargetY
		rescueTile := &bState.Board.Tiles[rx][ry]

		// ZMIANA: Visibility i unit pointer
		if rescueTile.Visibility != visibilityUnexplored {
			if rescueTile.Unit != nil {
				rescuedUnit := rescueTile.Unit
				if rescuedUnit.Exists &&
					rescuedUnit.Owner == bState.PlayerID &&
					rescuedUnit.Type == unitType(bState.CampaignData.TargetType) {
					log.Printf("RESCUE: cel jest w punkcie!")
				}
			}
		}
	}
}

func placeDestroyedPalisade(x, y uint8, bld *building, bState *battleState) {
	currentTile := &bState.Board.Tiles[x][y]
	currentTile.TextureID = spritePalisadeDestroyed
	currentTile.IsWalkable = true
	bld.IsUnderConstruction = true
}

func placeRuins(bState *battleState, bld *building) {
	if len(bld.OccupiedTiles) == 0 {
		return
	}

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

		occupedTile := &bState.Board.Tiles[x][y]

		// Czyścimy wskaźnik na budynek
		occupedTile.Building = nil

		// Ustawiamy grafikę ruin
		dx := pt.X - minX
		dy := pt.Y - minY
		idx := dy*width + dx
		occupedTile.TextureID = spriteRuinStart + uint16(idx)
		// Zgliszcza uniemożliwiają ruch
		occupedTile.IsWalkable = false
	}
}

// Usuwa uśmiercone jednostki z bs.
// Nie mylić z logiką rozkładu zwłok updateCorpses().
func cleanupDeadUnits(bState *battleState) {
	if bState.GlobalFrameCounter%6000 != 0 {
		return
	}

	if len(bState.Units) < int(maxUnitsPerPlayer)*4 {
		return
	}

	newUnitsList := make([]*unit, 0, len(bState.Units))
	for _, u := range bState.Units {
		if u.Exists {
			newUnitsList = append(newUnitsList, u)
		}
	}

	bState.Units = newUnitsList
}

func cleanupDestroyedBuildings(bState *battleState) {
	if bState.GlobalFrameCounter%6000 != 0 {
		return
	}

	if len(bState.Buildings) < int(maxBuildingsPerPlayer)*4 {
		return
	}

	log.Println("INFO: Rozpoczynam czyszczenie pamięci z budynków...")

	newBuildingsList := make([]*building, 0, len(bState.Buildings))
	for _, bld := range bState.Buildings {
		if bld.Exists {
			newBuildingsList = append(newBuildingsList, bld)
		}
	}

	removedCount := len(bState.Buildings) - len(newBuildingsList)
	bState.Buildings = newBuildingsList

	if removedCount > 0 {
		log.Printf("INFO: Wyczyszczono %d zniszczonych budynków.", removedCount)
	}
}
