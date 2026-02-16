package main

// level.go

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// @todo @reminder Wróć do tego pliku, ogarnij te wszystkie „apply” ponieważ nie powinno to się odpalać przy każdej nowej grze a tylko
// przy wczytywaniu nowe mapy. Ponadto wydaje mi się, że można zmodyfikować też mapy z kampanii, aby się tego w ogóle pozbyć i zapisać przetworzone

// Tworzy nowy przedmiot levelLoader.
func newLevelLoader(drivePath string) *jsonLevelLoader {
	return &jsonLevelLoader{drivePath: drivePath}
}

func initBoard(bs *battleState) {
	for x := uint8(0); x < boardMaxX; x++ {
		for y := uint8(0); y < boardMaxY; y++ {
			// Dostęp do kafelka przez wskaźnik (dla wydajności przy edycji)
			tile := &bs.Board.Tiles[x][y]

			// Wartości domyślne
			tile.TextureID = spriteGrassStart
			// tile.EffectID = 0
			tile.Unit = nil
			tile.Building = nil
			tile.Visibility = visibilityUnexplored // Na start wszystko czarne (gdy wdrożymy mgłę)

			// Ustawianie przejezdności
			// Granice mapy są nieprzejezdne
			if x == 0 || x == boardMaxX-1 || y == 0 || y == boardMaxY-1 {
				tile.IsWalkable = false
			} else {
				tile.IsWalkable = true
			}

			tile.MovementCost = 1.0
		}
	}

	bs.NextUniqueObjectID = 1
	log.Println("INFO: Zaczyn planszy skończony.")
}

// Stosuje dane z JSON do stanu gry.
func (l *jsonLevelLoader) applyJSONLevel(jsonLevel *jsonLevel, bs *battleState) {
	l.clearGameState(bs)

	// Metadane poziomu
	// @todo: Nigdy nie dodałem CampaignData do battleState, być może trzeba to zmienić!
	bs.CampaignData.DecisionType = jsonLevel.Metadata.DecisionType
	bs.CampaignData.EndCondition = jsonLevel.Metadata.EndType
	bs.CampaignData.TargetType = jsonLevel.Metadata.TargetType
	bs.CampaignData.LevelsMilkLimit = jsonLevel.Metadata.MaxMilk
	bs.CampaignData.GeneratorActive = jsonLevel.Metadata.Generator
	bs.CampaignData.NextLevel = jsonLevel.Metadata.NextLevel
	bs.CampaignData.Name = jsonLevel.Metadata.Name

	// Ustawienie początkowej pozycji kamery
	bs.GameCamera.Target.X = float32(jsonLevel.Metadata.StartPos.X * tileWidth)
	bs.GameCamera.Target.Y = float32(jsonLevel.Metadata.StartPos.Y * tileHeight)
	log.Printf("INFO: Początkowa pozycja kamery ustawiona na (%f, %f) w kafelkach.",
		bs.GameCamera.Target.X, bs.GameCamera.Target.Y)

	// Ustawienia SI
	bs.AI.GatherPointX = jsonLevel.AISettings.GatherPoint.X
	bs.AI.GatherPointY = jsonLevel.AISettings.GatherPoint.Y
	bs.AI.PastureX = jsonLevel.AISettings.Pasture.X
	bs.AI.PastureY = jsonLevel.AISettings.Pasture.Y

	// Specjalne lokacje
	if jsonLevel.SpecialLocations.TransformationPoint != nil {
		bs.CampaignData.TransformationSiteX = jsonLevel.SpecialLocations.TransformationPoint.X
		bs.CampaignData.TransformationSiteY = jsonLevel.SpecialLocations.TransformationPoint.Y
	}
	if jsonLevel.SpecialLocations.VictoryPoint != nil {
		bs.CampaignData.VictoryPointX = jsonLevel.SpecialLocations.VictoryPoint.X
		bs.CampaignData.VictoryPointY = jsonLevel.SpecialLocations.VictoryPoint.Y
	}
	if jsonLevel.SpecialLocations.RescueTarget != nil {
		bs.CampaignData.RescueTargetX = jsonLevel.SpecialLocations.RescueTarget.X
		bs.CampaignData.RescueTargetY = jsonLevel.SpecialLocations.RescueTarget.Y
	}

	// Zastosuj surowy teren
	l.applyTerrain(&jsonLevel.Terrain, bs)

	// Zastosuj budynki
	l.applyBuildings(jsonLevel.Buildings, bs)

	// Zastosuj jednostki
	l.applyUnits(jsonLevel.Units, bs)

	// Zasoby graczy
	bs.HumanPlayerState.MaxMilk = jsonLevel.Metadata.MaxMilk
	bs.HumanPlayerState.Milk = jsonLevel.Metadata.MaxMilk
	bs.AIEnemyState.MaxMilk = 1800
	bs.AIEnemyState.Milk = 1800

	log.Printf("Zastosowano poziom: %s (%dx%d)",
		jsonLevel.Metadata.Name, jsonLevel.Terrain.Width, jsonLevel.Terrain.Height)
}

// Wyczyść stan gry
func (l *jsonLevelLoader) clearGameState(bs *battleState) {
	initBoard(bs)

	// POPRAWKA: Resetujemy długość do 0, ale zachowujemy pojemność (Capacity).
	// Zapobiega to utracie "rezerwacji" 40 miejsc i zbędnym alokacjom.
	bs.Units = bs.Units[:0]
	bs.Buildings = bs.Buildings[:0]
	// Jeśli projectiles są zainicjowane, też je czyścimy
	if bs.Projectiles != nil {
		bs.Projectiles = bs.Projectiles[:0]
	}

	// Zerujemy graczy
	bs.HumanPlayerState.init(bs.PlayerID, 0)
	bs.AIEnemyState.init(bs.AIPlayerID, 0)

	// ZEROWANIE LICZNIKÓW POPULACJI I BUDYNKÓW
	bs.HumanPlayerState.CurrentPopulation = 0
	bs.AIEnemyState.CurrentPopulation = 0
	bs.HumanPlayerState.CurrentBuildings = 0
	bs.AIEnemyState.CurrentBuildings = 0

	bs.CurrentSelection = selectionState{}
	bs.MouseCommandMode = 1
	bs.IsSelectingBox = false

	// Zerowanie liczników animacji
	bs.GrassGrowthCounter = 0
	bs.WaterAnimationCounter = 0
	bs.FireAnimationCounter = 0
	bs.GlobalFrameCounter = 0

	bs.AI = aiState{}
	bs.AI.MilkGenerationRate = uint16(bs.DifficultyLevel)
}

func configureTile(tile *tile, graphicID uint16) {
	tile.IsWalkable = true
	tile.MovementCost = 1.0

	if isWaterTileOnly(graphicID) {
		tile.IsWalkable = false
	}

	if isRockNonWalkable(graphicID) {
		tile.IsWalkable = false
	}

	if isRoad(graphicID) {
		tile.MovementCost = 0.5
	}

	if isTreeStump(graphicID) || isRuin(graphicID) {
		tile.IsWalkable = false
	}

	if isGadget(graphicID) {
		tile.IsWalkable = false
	}
}

func (l *jsonLevelLoader) spawnPalisade(x, y uint8, graphicID uint16, bs *battleState) {
	if graphicID == spritePalisadeNE {
		newPalisade := &building{}
		newPalisade.initConstruction(x, y, buildingPalisade, colorNone, bs)

		bs.Buildings = append(bs.Buildings, newPalisade)
	}
}

func (l *jsonLevelLoader) applyTerrain(terrain *jsonTerrainData, bs *battleState) {
	log.Println("INFO: Nakładanie terenu...")

	// Bezpieczne granice pętli (zamiast rzutowania uint w każdym obiegu)
	maxY := min(len(terrain.Tiles), int(boardMaxY))

	for yAxis := range maxY {
		maxX := min(len(terrain.Tiles[yAxis]), int(boardMaxX))

		for xAxis := range maxX {
			// 1. Pobranie ID
			tileName := terrain.Tiles[yAxis][xAxis]
			graphicID := getID(tileName)

			// 2. Pobranie wskaźnika na kafelek
			tX, tY := uint8(xAxis), uint8(yAxis)
			tile := &bs.Board.Tiles[tX][tY]

			// 3. Przypisanie grafiki
			tile.TextureID = graphicID

			// 4. Delegacja logiki (Czystość!)
			configureTile(tile, graphicID)
			l.spawnPalisade(tX, tY, graphicID, bs)
		}
	}

	log.Println("INFO: Teren nałożony pomyślnie.")
}

func (l *jsonLevelLoader) applyBuildings(buildingsData []jsonBuildingData, bs *battleState) {
	log.Printf("INFO: Ładowanie %d budynków do battleState...", len(buildingsData))

	for _, data := range buildingsData {
		var ownerID uint8
		switch data.Owner {
		case "ENEMY":
			ownerID = bs.AIPlayerID
		case "PLAYER":
			ownerID = bs.PlayerID
		default:
			ownerID = colorNone
		}

		buildingType, exists := buildingTypeMap[data.Type]
		if !exists {
			continue
		}

		// Pobieramy definicję, aby znać wymiary (Width/Height)
		stats, ok := buildingDefs[buildingType]
		if !ok {
			log.Printf("BŁĄD: Brak definicji statystyk dla budynku typu %d", buildingType)
			continue
		}

		// Przeliczenie współrzędnych
		// JSON zawiera współrzędne Prawego-Dolnego rogu (stary system).
		// Funkcja init() wymaga teraz Lewego-Górnego rogu (nowy system).
		topLeftX := data.Position.X - stats.Width + 1
		topLeftY := data.Position.Y - stats.Height + 1

		newBuilding := &building{}

		// Wywołujemy init z przeliczonymi współrzędnymi Top-Left
		newBuilding.initConstruction(topLeftX, topLeftY, buildingType, ownerID, bs)

		bs.Buildings = append(bs.Buildings, newBuilding)

		// Liczniki
		switch newBuilding.Owner {
		case bs.HumanPlayerState.PlayerID:
			bs.HumanPlayerState.CurrentBuildings++
		case bs.AIEnemyState.PlayerID:
			bs.AIEnemyState.CurrentBuildings++
		}

		// Nakładanie grafiki przy użyciu wyliczonego Top-Left
		template, templateExists := buildingTemplates[buildingType]
		if templateExists {
			for dy, row := range template {
				for dx, graphicID := range row {
					tileX := topLeftX + uint8(dx)
					tileY := topLeftY + uint8(dy)

					if tileX < boardMaxX && tileY < boardMaxY {
						bs.Board.Tiles[tileX][tileY].TextureID = uint16(graphicID)
					}
				}
			}
		}
	}
}

// levelLoader.go - applyUnits
func (l *jsonLevelLoader) applyUnits(units []jsonUnitData, bs *battleState) {
	log.Printf("INFO: Ładowanie %d jednostek do battleState...", len(units))

	for _, unitData := range units {
		var ownerID uint8
		if unitData.Owner == "ENEMY" {
			ownerID = bs.AIPlayerID
		} else {
			ownerID = bs.PlayerID
		}

		unitType, exists := unitTypeMap[unitData.Type]
		if !exists {
			log.Printf("OSTRZEŻENIE: Pomięto nieznany typ jednostki: %s", unitData.Type)
			continue
		}

		newUnit := &unit{}
		// Inicjalizacja nowej jednostki
		newUnit.initUnit(unitType, unitData.Position.X, unitData.Position.Y, cmdIdle, bs)
		newUnit.Owner = ownerID

		// Wstawienie na mapę
		newUnit.show(bs)

		// Dodanie do głównej listy
		bs.Units = append(bs.Units, newUnit)

		// Zliczamy jednostki startowe, aby limit 40 działał poprawnie od początku gry.
		switch newUnit.Owner {
		case bs.HumanPlayerState.PlayerID:
			bs.HumanPlayerState.CurrentPopulation++
		case bs.AIEnemyState.PlayerID:
			bs.AIEnemyState.CurrentPopulation++
		}
	}
	log.Printf("INFO: Załadowano %d jednostek. Populacja Gracza: %d, AI: %d.",
		len(bs.Units), bs.HumanPlayerState.CurrentPopulation, bs.AIEnemyState.CurrentPopulation)
}

// Załaduj poziom
func (l *jsonLevelLoader) loadLevel(levelNum uint8, bs *battleState) error {
	if err := l.loadJSONLevel(levelNum, bs); err == nil {
		return nil
	}
	return fmt.Errorf("nie znaleziono poziomu JSON %d", levelNum)
}

func (l *jsonLevelLoader) loadJSONLevel(levelNum uint8, bs *battleState) error {
	log.Printf("Ładowanie poziomu %d z JSON...", levelNum)

	filename := filepath.Join(l.drivePath, "levels_json", fmt.Sprintf("level_%02d.json", levelNum))

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("nie można wczytać pliku %s: %w", filename, err)
	}

	var jsonLevel jsonLevel
	if err := json.Unmarshal(data, &jsonLevel); err != nil {
		return fmt.Errorf("błąd parsowania JSON: %w", err)
	}

	log.Printf("Załadowano poziom: %s", jsonLevel.Metadata.Name)
	l.applyJSONLevel(&jsonLevel, bs)
	return nil
}

// Waliduj pozycję kamery
func (l *jsonLevelLoader) validateScreenPosition(bs *battleState) {
	// Ile kafelków mieści się na ekranie (wirtualnym)
	tilesVisibleX := virtualScreenWidth / uint16(tileWidth)   // 640/16 = 40
	tilesVisibleY := virtualScreenHeight / uint16(tileHeight) // 400/14 = 28

	// Konwersja celu kamery z pikseli na kafelki
	cameraTargetTileX := uint8(bs.GameCamera.Target.X / float32(tileWidth))
	cameraTargetTileY := uint8(bs.GameCamera.Target.Y / float32(tileHeight))

	// Minimalne i maksymalne pozycje celu kamery (w kafelkach), tak żeby kamera nie wyszła poza mapę
	// Cel kamery to środek ekranu. Offset kamery to środek ekranu wirtualnego (320, 200).
	// Więc min target to środek ekranu, a max target to (rozmiar mapy - połowa ekranu).
	minTargetTileX := uint8(tilesVisibleX / 2)
	minTargetTileY := uint8(tilesVisibleY / 2)
	maxTargetTileX := boardMaxX - uint8(tilesVisibleX/2)
	maxTargetTileY := boardMaxY - uint8(tilesVisibleY/2)

	log.Printf("DEBUG: Limity kamery w kafelkach: x(%d-%d), y(%d-%d)", minTargetTileX,
		maxTargetTileX, minTargetTileY, maxTargetTileY)

	// Korekta pozycji celu kamery
	if cameraTargetTileX < minTargetTileX {
		cameraTargetTileX = minTargetTileX
	}
	if cameraTargetTileY < minTargetTileY {
		cameraTargetTileY = minTargetTileY
	}
	if cameraTargetTileX > maxTargetTileX {
		cameraTargetTileX = maxTargetTileX
	}
	if cameraTargetTileY > maxTargetTileY {
		cameraTargetTileY = maxTargetTileY
	}

	// Zaktualizuj cel kamery w pikselach
	bs.GameCamera.Target.X = float32(cameraTargetTileX * tileWidth)
	bs.GameCamera.Target.Y = float32(cameraTargetTileY * tileHeight)

	log.Printf("DEBUG: Ostateczna pozycja celu kamery (skorygowana): (%d, %d) kafelków = (%.1f, %.1f) pikseli",
		cameraTargetTileX, cameraTargetTileY, bs.GameCamera.Target.X, bs.GameCamera.Target.Y)
}

// initBattle - metoda LevelLoader która inicjuje bitwę z ProgramState
func (l *jsonLevelLoader) initBattle(levelNumber uint8, bs *battleState) error {
	log.Printf("Inicjalizacja bitwy poziom %d", levelNumber)
	// @todo: przekaż poziom trudności do „battlestate”!!

	initBoard(bs)

	if err := l.loadLevel(levelNumber, bs); err != nil {
		return fmt.Errorf("nie można załadować poziomu %d: %w", levelNumber, err)
	}

	processMapTiles(bs)

	l.validateScreenPosition(bs)

	log.Println("INFO: Bitwa rozpoczęta pomyślnie.")
	return nil
}
