package main

// ai.go

import (
	"log"
	"math"
	"sort"
)

// 01.01.2026 Ta część programu nie była jeszcze nigdy sprawdzona. Wygląda spójnie i sensownie, ale to tylko pozory.

// Chyba już nieużywane.
const (
	aiPhaseInit          = 0
	aiPhaseBuildCows     = 1
	aiPhaseBuildUnits    = 11
	aiPhaseDeployUnits   = 2
	aiPhaseCycleComplete = 3
	aiPhaseDefeated      = 5
	aiPhaseStaticDefense = 6
)

// SZTUCZNA INTELIGENCJA

// Główna funkcja decyzyjna SI, operująca na battleState
func aiMakeDecision(bState *battleState) {
	ai := &bState.AI

	// Inicjalizacja - tylko raz
	if ai.EconomyManager == nil {
		ai.EconomyManager = newEconomyAI()
		ai.MilitaryManager = newMilitaryAI(bState)
		ai.CurrentGoals = []aiGoal{
			{Type: goalProduceCows, Priority: 0.8},
			{Type: goalBuildArmy, Priority: 0.6},
		}
		log.Println("AI: Zainicjowano nowy system")
	}

	// Generuj mleko
	bState.AIEnemyState.Milk += ai.MilkGenerationRate

	// Aktualizuj cele co 100 pyknięć
	if ai.ActionDelay%100 == 0 {
		ai.updateGoals(bState)
	}

	// Wykonuj cele w kolejności priorytetu
	for _, goal := range ai.CurrentGoals {
		if ai.executeGoal(goal, bState) {
			break // Jeden cel na tick
		}
	}

	ai.ActionDelay++
}

func (ai *aiState) updateGoals(bState *battleState) {
	// Dostosowanie celów
	cowCount := ai.countCows(bState)
	armyStrength := ai.countArmyStrength(bState)

	// Jeżeli niedobór, krów to zwiększamy wagę
	if cowCount < 5 {
		ai.setGoalPriority(goalProduceCows, 1.0)
	} else {
		ai.setGoalPriority(goalProduceCows, 0.3)
	}

	// Jeżeli mało wojsk, to zwiększamy wagę
	if armyStrength < 10 {
		ai.setGoalPriority(goalBuildArmy, 0.9)
	}

	// Ułóż cele wg wagi
	sort.Slice(ai.CurrentGoals, func(i, j int) bool {
		return ai.CurrentGoals[i].Priority > ai.CurrentGoals[j].Priority
	})
}

func newEconomyAI() *economyAI {
	return &economyAI{lastBarnChecked: 0}
}

func (econ *economyAI) produceCows(bState *battleState) bool {
	// Znajdź oborę z miejscem
	barn := econ.findAvailableBarn(bState)
	if barn == nil {
		return false
	}

	// Sprawdź czy stać nas na krowę
	cowCost := unitDefs[unitCow].Cost
	if bState.AIEnemyState.Milk < cowCost {
		return false
	}

	// Wytwarzamy krowę
	cmd := command{
		ActionType:          cmdBProduce,
		InteractionTargetID: ObjectID(barn.ID),
		CreateType:          uint8(unitCow),
	}
	bState.CurrentCommands[1] = cmd
	bState.AIEnemyState.Milk -= cowCost
	barn.Food++

	log.Printf("AI: Produkuję krowę w oborze %d", barn.ID)
	return true
}

// TODO: nie bardzo rozumiem czemu mamy takie dziwne poszukiwania zamiast liniowego
func (econ *economyAI) findAvailableBarn(bState *battleState) *building {
	for i := 0; i < len(bState.Buildings); i++ {
		bld := bState.Buildings[i]
		if bld != nil && bld.Exists &&
			bld.Owner == bState.AIEnemyState.PlayerID &&
			bld.Type == buildingBarn &&
			bld.Food < bld.MaxFood {
			return bld
		}
	}

	return nil
}

func (ai *aiState) countCows(bState *battleState) int {
	count := 0
	for _, unit := range bState.Units {
		if unit != nil && unit.Exists && unit.Owner == bState.AIEnemyState.PlayerID && unit.Type == unitCow {
			count++
		}
	}
	return count
}

func (ai *aiState) countArmyStrength(bState *battleState) int {
	count := 0
	for _, unit := range bState.Units {
		if unit != nil && unit.Exists && unit.Owner == bState.AIEnemyState.PlayerID && unit.Type != unitCow &&
			unit.Type != unitShepherd {
			count++
		}
	}
	return count
}

func (ai *aiState) setGoalPriority(goalType aiGoalType, priority float32) {
	// Sprawdzamy, czy cel już jest w kolejce
	for i := range ai.CurrentGoals {
		if ai.CurrentGoals[i].Type == goalType {
			ai.CurrentGoals[i].Priority = priority
			return
		}
	}

	// Jeżeli jest nowe, to dodajemy
	ai.CurrentGoals = append(ai.CurrentGoals, aiGoal{
		Type:     goalType,
		Priority: priority,
	})
}

func (ai *aiState) executeGoal(goal aiGoal, bState *battleState) bool {
	switch goal.Type {
	case goalProduceCows:
		return ai.EconomyManager.produceCows(bState)
	case goalBuildArmy:
		// @todo: wytwarzanie jednostek
		return false
	case goalAttackEnemy:
		// @todo: napadanie na przeciwnika
		return false
	default:
		return false
	}
}

func newMilitaryAI(bState *battleState) *militaryAI {
	return &militaryAI{}
}

func (mil *militaryAI) buildUnits(bState *battleState) bool {
	// Na razie nic nie rób
	return false
}

func (mil *militaryAI) attackEnemy(bState *battleState) bool {
	// Na razie nic nie rób
	return false
}

// SI tworzące nowe jednostki w dynamicznych bitwach
/*func aiMakeNewUnits(bState *battleState) {
	if !bState.CampaignData.GeneratorActive {
		return
	}

	// --- 1. Obsługa ogólnego licznika tworzenia nowych jednostek ---
	// Próg generowania jednostek, dostosowany do poziomu trudności.
	generationThreshold := 150 - int(bState.DifficultyLevel)*50
	if bState.CampaignData.GeneratorTimer < generationThreshold {
		bState.CampaignData.GeneratorTimer++
	} else {
		// Timer osiągnął próg, więc próbujemy stworzyć jednostkę i zerujemy licznik.
		bState.CampaignData.GeneratorTimer = 0

		spawnX, spawnY := uint8(1), uint8(5) // Przypisuję miejsce stworzenia na sztywno @todo: sprawdź, czy przy tworzeniu przez gracza jest już układ doboru miejsca

		tile := &bState.Board.Tiles[spawnX][spawnY]
		isSpawnPointFree := tile.Unit == nil && tile.Building == nil && tile.IsWalkable

		if isSpawnPointFree {
			// Sprawdź, czy "generująca" jednostka (ID 2) już istnieje i jest żywa.
			var existingSpecialUnitID2 *unit
			for _, unit := range bState.Units {
				if unit != nil && unit.Exists && unit.Owner == bState.AIEnemyState.PlayerID && unit.ID == 2 {
					existingSpecialUnitID2 = unit
					break
				}
			}

			// Tylko jeśli specjalna jednostka ID 2 nie istnieje, tworzymy nową.
			if existingSpecialUnitID2 == nil {
				newUnit := &unit{}
				unitTypesForGenerator := []unitType{
					unitAxeman, unitArcher, unitSwordsman, unitSpearman,
					unitCommander, unitBear, unitPriest, unitPriestess,
				}
				randomUnitType := unitTypesForGenerator[rand.Intn(len(unitTypesForGenerator))]

				newUnit.initUnit(randomUnitType, spawnX, spawnY, cmdUIdle, bState)
				newUnit.Owner = bState.AIEnemyState.PlayerID
				newUnit.Experience = uint8(20 + int(bState.DifficultyLevel)*20)

				// Rejestracja w nowej strukturze
				newUnit.show(bState)

				bState.Units = append(bState.Units, newUnit)

				// Aktualizacja licznika populacji SI (skoro zmieniliśmy strukturę, warto to trzymać)
				bState.AIEnemyState.CurrentPopulation++

				// Natychmiastowe wydanie rozkazu nowej jednostce
				var targetX, targetY uint8
				var targetUnit *unit
				var targetBuilding *building
				var found bool
				var ok bool
				var tempX uint8
				var tempY uint8
				targetUnit, targetBuilding, found = findNearestEnemyExtended(newUnit, bState) // Szukaj celu dla NOWEJ jednostki

				if found {
					if targetUnit != nil {
						targetX, targetY = targetUnit.X, targetUnit.Y
					} else if targetBuilding != nil {
						// Ważne: dla budynków znajdź dostępny kafelek wokół nich
						tempX, tempY, ok = targetBuilding.getClosestWalkableTile(bState)
						if ok {
							targetX, targetY = tempX, tempY
						} else {
							targetX, targetY, ok = targetBuilding.getCenter()
						}
					}

					if ok {
						targetIDForCommand := newUnit.ID
						if targetUnit != nil {
							targetIDForCommand = targetUnit.ID
						} else if targetBuilding != nil {
							targetIDForCommand = targetBuilding.ID
						}

						newUnit.addUnitCommand(cmdUAttack, targetX, targetY, targetIDForCommand, bState)
						log.Printf("SI (wytwórca): Nowa jednostka %d (%d) wysłana na cel (%d,%d) ID %d.",
							newUnit.ID, newUnit.Type, targetX, targetY, targetIDForCommand)
					} else {
						newUnit.setIdle()
					}
				} else {
					newUnit.setIdle()
				}
			}
		} else {
			log.Printf("SI (wytwórca): Punkt (%d,%d) jest zablokowany. Nie można stworzyć nowej jednostki.", spawnX, spawnY)
		}
	}

	// --- 2. Obsługa komend dla jednostki „wytwórczej” (ID 2), jeśli istnieje ---
	if bState.CampaignData.GeneratorTimer%25 == 0 && bState.CampaignData.GeneratorTimer != 0 {
		var generatorUnit *unit
		for _, unit := range bState.Units {
			if unit != nil && unit.Exists && unit.Owner == bState.AIEnemyState.PlayerID && unit.ID == 2 {
				generatorUnit = unit
				break
			}
		}

		if generatorUnit != nil {
			var targetX, targetY uint8
			var targetUnit *unit
			var targetBuilding *building
			var found bool
			targetUnit, targetBuilding, found = findNearestEnemyExtended(generatorUnit, bState)

			if found {
				if targetUnit != nil {
					targetX, targetY = targetUnit.X, targetUnit.Y
				} else if targetBuilding != nil {
					tempX, tempY, ok := targetBuilding.getClosestWalkableTile(bState)
					if ok {
						targetX, targetY = tempX, tempY
					} else {
						targetX, targetY, ok = targetBuilding.getCenter()
					}
				}

				targetIDForCommand := generatorUnit.ID
				if targetUnit != nil {
					targetIDForCommand = targetUnit.ID
				} else if targetBuilding != nil {
					targetIDForCommand = targetBuilding.ID
				}
				generatorUnit.addUnitCommand(cmdUAttack, targetX, targetY, targetIDForCommand, bState)
			} else {
				generatorUnit.setIdle()
			}
		}
	}
}
*/
// @todo: ogarnij czemu szuka po całej planszy a nie tylko w zasięgu wzroku
// findNearestEnemyExtended znajduje najbliższego wroga
func findNearestEnemyExtended(seeker *unit, bState *battleState) (*unit, *building, bool) {
	const (
		PriorityUnit         = 100
		PriorityBuilding     = 200
		PriorityImportantPal = 300
		MaxScore             = math.MaxFloat64
	)

	bestScore := MaxScore
	var chosenUnit *unit
	var chosenBuilding *building

	// Faza 1: Jednostki wroga

	for _, u := range bState.Units {
		if u == nil || !u.Exists || u.Owner == seeker.Owner {
			continue
		}
		distance := float64(seeker.calculateDistanceToTarget(&combatTarget{Unit: u}, bState.Board))

		if distance > float64(seeker.SightRange) || distance == float64(math.MaxInt32) {
			continue
		}

		currentScore := float64(PriorityUnit) + distance
		if currentScore < bestScore {
			bestScore = currentScore
			chosenUnit = u
			chosenBuilding = nil
		}
	}

	if chosenUnit != nil {
		return chosenUnit, nil, true
	}

	// Faza 2: Budynki wroga (inne niż palisady)
	for _, bld := range bState.Buildings {
		if bld == nil || !bld.Exists || bld.Owner == seeker.Owner || bld.Type == buildingPalisade {
			continue
		}
		distance := float64(getDistanceToUnit(bld.Type, bld.OccupiedTiles[0], seeker.X, seeker.Y))

		if distance > float64(seeker.SightRange) || distance == float64(math.MaxInt32) {
			continue
		}

		currentScore := float64(PriorityBuilding) + distance
		if currentScore < bestScore {
			bestScore = currentScore
			chosenUnit = nil
			chosenBuilding = bld
		}
	}

	if chosenBuilding != nil {
		return nil, chosenBuilding, true
	}

	// Faza 3: Ważne palisady
	for _, bld := range bState.Buildings {
		if bld == nil || !bld.Exists || bld.Owner == seeker.Owner || bld.Type != buildingPalisade {
			continue
		}

		if !seeker.isImportantPalisade(bld, bState) {
			continue
		}

		distance := float64(getDistanceToUnit(bld.Type, bld.OccupiedTiles[0], seeker.X, seeker.Y))

		if distance > float64(seeker.SightRange) || distance == float64(math.MaxInt32) {
			continue
		}

		currentScore := float64(PriorityImportantPal) + distance
		if currentScore < bestScore {
			bestScore = currentScore
			chosenUnit = nil
			chosenBuilding = bld
		}
	}

	if chosenBuilding != nil {
		return nil, chosenBuilding, true
	}

	return nil, nil, false
}

/*func findNearestEnemyCached(seeker *unit, bState *battleState) (*unit, *building, bool) {
	if entry, exists := bState.EnemyCache[seeker.ID]; exists {
		if (bState.GlobalFrameCounter - entry.LastUpdateTick) < entry.CacheValidFor {
			return entry.NearestEnemyUnit, entry.NearestEnemyBuilding, entry.Found
		}
	}

	if bState.GlobalFrameCounter != bState.enemyCacheUpdateTick {
		bState.enemyCacheUpdateTick = bState.GlobalFrameCounter
		bState.enemyCacheUpdatesThisTick = 0
	}

	if bState.enemyCacheUpdatesThisTick >= 3 {
		if entry, exists := bState.EnemyCache[seeker.ID]; exists {
			return entry.NearestEnemyUnit, entry.NearestEnemyBuilding, entry.Found
		}
		return nil, nil, false
	}

	u, bld, found := findNearestEnemyExtended(seeker, bState)

	var cacheTime uint16 = 15
	if !found {
		cacheTime = 30
	}

	bState.EnemyCache[seeker.ID] = &enemyCacheEntry{
		NearestEnemyUnit:     u,
		NearestEnemyBuilding: bld,
		Found:                found,
		LastUpdateTick:       bState.GlobalFrameCounter,
		CacheValidFor:        cacheTime,
	}

	bState.enemyCacheUpdatesThisTick++
	return u, bld, found
}*/

// findGrass znajduje odpowiedni kafelek trawy do wypasu krów.
// @todo: czemu to jest w ai.go a nie cow.go?
func findGrass(xp, yp uint8, xe, ye *uint8, bState *battleState) {
	*xe = uint8(0)
	*ye = uint8(0)

	for k := uint8(2); k <= 15; k++ {
		for i := xp - k; i <= xp+k && i < boardMaxX-1; i++ {
			if i <= 0 {
				continue
			}
			for j := yp - k; j <= yp+k && j < boardMaxY-1; j++ {
				if j <= 0 {
					continue
				}

				tile := &bState.Board.Tiles[i][j]
				if (i == xp-k || i == xp+k || j == yp-k || j == yp+k) &&
					// @todo: @reminder: ogarnij, czy dobry SPRITE_GRASS dodałeś
					tile.TextureID < spriteGrassEnd &&
					tile.Unit == nil && tile.Building == nil {
					*xe = i
					*ye = j
					return
				}
			}
		}
	}

	*xe = xp + 5
	*ye = yp + 3
}

/*func sendCowToBarn(cow *unit, bState *battleState) bool {
	for _, building := range bState.Buildings {
		if building == nil || !building.Exists {
			continue
		}

		if building.Type == buildingBarn {
			for _, unitID := range building.AssignedUnits {
				if unitID == cow.ID {
					targetX, targetY, ok := building.getClosestWalkableTile(bState)

					if ok {
						cow.addUnitCommand(cmdUMove, targetX, targetY, 0, bState)
						log.Printf("Krowa %d wraca do obory %d na pozycję (%d,%d).", cow.ID, building.ID, targetX, targetY)
						return true
					}

					log.Printf("OSTRZEŻENIE: Nie można znaleźć wolnego miejsca obok obory %d dla krowy %d.", building.ID, cow.ID)
					return false
				}
			}
		}
	}
	return false
}
*/
// who sprawdza do kogo przynależy dana jednostka
func who(oID ObjectID, bState *battleState) PlayerID {
	if oID == 0 {
		return 0
	}
	// Ta funkcja iteruje po listach, więc jest niezależna od zmian w BoardData
	for _, u := range bState.Units {
		if u != nil && u.Exists && ObjectID(u.ID) == oID {
			return u.Owner
		}
	}
	for _, bld := range bState.Buildings {
		if bld != nil && bld.Exists && ObjectID(bld.ID) == oID {
			return bld.Owner
		}
	}
	return 0
}
