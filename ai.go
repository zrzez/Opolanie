package main

// ai.go

import (
	"log"
	"math"
	"math/rand"
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
func aiMakeDecision(bs *battleState) {
	ai := &bs.AI

	// Inicjalizacja - tylko raz
	if ai.EconomyManager == nil {
		ai.EconomyManager = newEconomyAI()
		ai.MilitaryManager = newMilitaryAI(bs)
		ai.CurrentGoals = []aiGoal{
			{Type: goalProduceCows, Priority: 0.8},
			{Type: goalBuildArmy, Priority: 0.6},
		}
		log.Println("AI: Zainicjowano nowy system")
	}

	// Generuj mleko
	bs.AIEnemyState.Milk += ai.MilkGenerationRate

	// Aktualizuj cele co 100 pyknięć
	if ai.ActionDelay%100 == 0 {
		ai.updateGoals(bs)
	}

	// Wykonuj cele w kolejności priorytetu
	for _, goal := range ai.CurrentGoals {
		if ai.executeGoal(goal, bs) {
			break // Jeden cel na tick
		}
	}

	ai.ActionDelay++
}

func (ai *aiState) updateGoals(bs *battleState) {
	// Dostosowanie celów
	cowCount := ai.countCows(bs)
	armyStrength := ai.countArmyStrength(bs)

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

func (econ *economyAI) produceCows(bs *battleState) bool {
	// Znajdź oborę z miejscem
	barn := econ.findAvailableBarn(bs)
	if barn == nil {
		return false
	}

	// Sprawdź czy stać nas na krowę
	cowCost := unitDefs[unitCow].Cost
	if bs.AIEnemyState.Milk < cowCost {
		return false
	}

	// Wytwarzamy krowę
	cmd := command{
		CommandCategory:     0,
		ActionType:          cmdProduce,
		InteractionTargetID: barn.ID,
		ProduceType:         unitCow,
	}
	bs.CurrentCommands[1] = cmd
	bs.AIEnemyState.Milk -= cowCost
	barn.Food++

	log.Printf("AI: Produkuję krowę w oborze %d", barn.ID)
	return true
}

// TODO: nie bardzo rozumiem czemu mamy takie dziwne poszukiwania zamiast liniowego
func (econ *economyAI) findAvailableBarn(bs *battleState) *building {
	for i := 0; i < len(bs.Buildings); i++ {
		building := bs.Buildings[i]
		if building != nil && building.Exists &&
			building.Owner == bs.AIEnemyState.PlayerID &&
			building.Type == buildingBarn &&
			building.Food < building.MaxFood {
			return building
		}
	}
	return nil
}

func (ai *aiState) countCows(bs *battleState) int {
	count := 0
	for _, unit := range bs.Units {
		if unit != nil && unit.Exists && unit.Owner == bs.AIEnemyState.PlayerID && unit.Type == unitCow {
			count++
		}
	}
	return count
}

func (ai *aiState) countArmyStrength(bs *battleState) int {
	count := 0
	for _, unit := range bs.Units {
		if unit != nil && unit.Exists && unit.Owner == bs.AIEnemyState.PlayerID && unit.Type != unitCow &&
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

func (ai *aiState) executeGoal(goal aiGoal, bs *battleState) bool {
	switch goal.Type {
	case goalProduceCows:
		return ai.EconomyManager.produceCows(bs)
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

func newMilitaryAI(bs *battleState) *militaryAI {
	return &militaryAI{}
}

func (mil *militaryAI) buildUnits(bs *battleState) bool {
	// Na razie nic nie rób
	return false
}

func (mil *militaryAI) attackEnemy(bs *battleState) bool {
	// Na razie nic nie rób
	return false
}

// SI tworzące nowe jednostki w dynamicznych bitwach
func aiMakeNewUnits(bs *battleState) {
	if !bs.CampaignData.GeneratorActive {
		return
	}

	// --- 1. Obsługa ogólnego licznika tworzenia nowych jednostek ---
	// Próg generowania jednostek, dostosowany do poziomu trudności.
	generationThreshold := 150 - int(bs.DifficultyLevel)*50
	if bs.CampaignData.GeneratorTimer < generationThreshold {
		bs.CampaignData.GeneratorTimer++
	} else {
		// Timer osiągnął próg, więc próbujemy stworzyć jednostkę i zerujemy licznik.
		bs.CampaignData.GeneratorTimer = 0

		spawnX, spawnY := uint8(1), uint8(5) // Przypisuję miejsce stworzenia na sztywno @todo: sprawdź, czy przy tworzeniu przez gracza jest już układ doboru miejsca

		tile := &bs.Board.Tiles[spawnX][spawnY]
		isSpawnPointFree := tile.Unit == nil && tile.Building == nil && tile.IsWalkable

		if isSpawnPointFree {
			// Sprawdź, czy "generująca" jednostka (ID 2) już istnieje i jest żywa.
			var existingSpecialUnitID2 *unit
			for _, unit := range bs.Units {
				if unit != nil && unit.Exists && unit.Owner == bs.AIEnemyState.PlayerID && unit.ID == 2 {
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

				newUnit.initUnit(randomUnitType, spawnX, spawnY, cmdIdle, bs)
				newUnit.Owner = bs.AIEnemyState.PlayerID
				newUnit.Experience = uint(20 + int(bs.DifficultyLevel)*20)

				// Rejestracja w nowej strukturze
				newUnit.show(bs)

				bs.Units = append(bs.Units, newUnit)

				// Aktualizacja licznika populacji SI (skoro zmieniliśmy strukturę, warto to trzymać)
				bs.AIEnemyState.CurrentPopulation++

				// Natychmiastowe wydanie rozkazu nowej jednostce
				var targetX, targetY uint8
				var targetUnit *unit
				var targetBuilding *building
				var found bool
				var ok bool
				var tempX uint8
				var tempY uint8
				targetUnit, targetBuilding, found = findNearestEnemyExtended(newUnit, bs) // Szukaj celu dla NOWEJ jednostki

				if found {
					if targetUnit != nil {
						targetX, targetY = targetUnit.X, targetUnit.Y
					} else if targetBuilding != nil {
						// Ważne: dla budynków znajdź dostępny kafelek wokół nich
						tempX, tempY, ok = targetBuilding.getClosestWalkableTile(bs)
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

						newUnit.addUnitCommand(cmdAttack, targetX, targetY, targetIDForCommand, bs)
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
	if bs.CampaignData.GeneratorTimer%25 == 0 && bs.CampaignData.GeneratorTimer != 0 {
		var generatorUnit *unit
		for _, unit := range bs.Units {
			if unit != nil && unit.Exists && unit.Owner == bs.AIEnemyState.PlayerID && unit.ID == 2 {
				generatorUnit = unit
				break
			}
		}

		if generatorUnit != nil {
			var targetX, targetY uint8
			var targetUnit *unit
			var targetBuilding *building
			var found bool
			targetUnit, targetBuilding, found = findNearestEnemyExtended(generatorUnit, bs)

			if found {
				if targetUnit != nil {
					targetX, targetY = targetUnit.X, targetUnit.Y
				} else if targetBuilding != nil {
					tempX, tempY, ok := targetBuilding.getClosestWalkableTile(bs)
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
				generatorUnit.addUnitCommand(cmdAttack, targetX, targetY, targetIDForCommand, bs)
			} else {
				generatorUnit.setIdle()
			}
		}
	}
}

// @todo: ogarnij czemu szuka po całej planszy a nie tylko w zasięgu wzroku
// findNearestEnemyExtended znajduje najbliższego wroga
func findNearestEnemyExtended(seeker *unit, bs *battleState) (*unit, *building, bool) {
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

	for _, unit := range bs.Units {
		if unit == nil || !unit.Exists || unit.Owner == seeker.Owner {
			continue
		}
		distance := float64(seeker.calculateDistanceToTarget(&combatTarget{Unit: unit}))

		if distance > float64(seeker.SightRange) || distance == float64(math.MaxInt32) {
			continue
		}

		currentScore := float64(PriorityUnit) + distance
		if currentScore < bestScore {
			bestScore = currentScore
			chosenUnit = unit
			chosenBuilding = nil
		}
	}

	if chosenUnit != nil {
		return chosenUnit, nil, true
	}

	// Faza 2: Budynki wroga (inne niż palisady)
	for _, bld := range bs.Buildings {
		if bld == nil || !bld.Exists || bld.Owner == seeker.Owner || bld.Type == buildingPalisade {
			continue
		}
		distance := float64(bld.getDistanceToUnit(seeker.X, seeker.Y))

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
	for _, bld := range bs.Buildings {
		if bld == nil || !bld.Exists || bld.Owner == seeker.Owner || bld.Type != buildingPalisade {
			continue
		}

		if !seeker.isImportantPalisade(bld, bs) {
			continue
		}

		distance := float64(bld.getDistanceToUnit(seeker.X, seeker.Y))

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

func findNearestEnemyCached(seeker *unit, bs *battleState) (*unit, *building, bool) {
	if entry, exists := bs.EnemyCache[seeker.ID]; exists {
		if (bs.GlobalFrameCounter - entry.LastUpdateTick) < entry.CacheValidFor {
			return entry.NearestEnemyUnit, entry.NearestEnemyBuilding, entry.Found
		}
	}

	if bs.GlobalFrameCounter != bs.enemyCacheUpdateTick {
		bs.enemyCacheUpdateTick = bs.GlobalFrameCounter
		bs.enemyCacheUpdatesThisTick = 0
	}

	if bs.enemyCacheUpdatesThisTick >= 3 {
		if entry, exists := bs.EnemyCache[seeker.ID]; exists {
			return entry.NearestEnemyUnit, entry.NearestEnemyBuilding, entry.Found
		}
		return nil, nil, false
	}

	unit, building, found := findNearestEnemyExtended(seeker, bs)

	var cacheTime uint16 = 15
	if !found {
		cacheTime = 30
	}

	bs.EnemyCache[seeker.ID] = &enemyCacheEntry{
		NearestEnemyUnit:     unit,
		NearestEnemyBuilding: building,
		Found:                found,
		LastUpdateTick:       bs.GlobalFrameCounter,
		CacheValidFor:        cacheTime,
	}

	bs.enemyCacheUpdatesThisTick++
	return unit, building, found
}

// findGrass znajduje odpowiedni kafelek trawy do wypasu krów.
// @todo: czemu to jest w ai.go a nie cow.go?
func findGrass(xp, yp uint8, xe, ye *uint8, bs *battleState) {
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

				tile := &bs.Board.Tiles[i][j]
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

func sendCowToBarn(cow *unit, bs *battleState) bool {
	for _, building := range bs.Buildings {
		if building == nil || !building.Exists {
			continue
		}

		if building.Type == buildingBarn {
			for _, unitID := range building.AssignedUnits {
				if unitID == cow.ID {
					targetX, targetY, ok := building.getClosestWalkableTile(bs)

					if ok {
						cow.addUnitCommand(cmdMove, targetX, targetY, 0, bs)
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

// who sprawdza do kogo przynależy dana jednostka
func who(objectID uint, bs *battleState) uint8 {
	if objectID == 0 {
		return 0
	}
	// Ta funkcja iteruje po listach, więc jest niezależna od zmian w BoardData
	for _, unit := range bs.Units {
		if unit != nil && unit.Exists && unit.ID == objectID {
			return unit.Owner
		}
	}
	for _, bld := range bs.Buildings {
		if bld != nil && bld.Exists && bld.ID == objectID {
			return bld.Owner
		}
	}
	return 0
}
