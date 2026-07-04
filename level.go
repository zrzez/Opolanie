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

func initBoard(bState *battleState) {
	for boardColumn := uint8(0); boardColumn < boardMaxX; boardColumn++ {
		for boardRow := uint8(0); boardRow < boardMaxY; boardRow++ {
			// Dostęp do kafelka przez wskaźnik
			currentTile := &bState.Board.Tiles[boardColumn][boardRow]

			// Wartości domyślne
			currentTile.X = boardColumn
			currentTile.Y = boardRow
			currentTile.TextureID = spriteGrassStart
			currentTile.Unit = nil
			currentTile.Building = nil
			currentTile.Visibility = visibilityUnexplored // Na start wszystko czarne (gdy wdrożymy mgłę)

			// Ustawianie przejezdności
			// Granice mapy są nieprzejezdne
			if boardColumn == 0 || boardColumn == boardMaxX-1 || boardRow == 0 || boardRow == boardMaxY-1 {
				currentTile.IsWalkable = false
			} else {
				currentTile.IsWalkable = true
			}

			currentTile.MovementCost = 1.0
		}
	}

	bState.NextUniqueObjectID = 1

	log.Println("INFO: Zaczyn planszy zrobiony.")
}

// Stosuje dane z JSON do stanu gry.
func (l *jsonLevelLoader) applyJSONLevel(jsonLevel *jsonLevel, bState *battleState) {
	l.clearGameState(bState)

	// Metadane poziomu
	// @todo: Nigdy nie dodałem CampaignData do battleState, być może trzeba to zmienić!
	bState.CampaignData.DecisionType = jsonLevel.Metadata.DecisionType
	bState.CampaignData.EndCondition = jsonLevel.Metadata.EndType
	bState.CampaignData.TargetType = jsonLevel.Metadata.TargetType
	bState.CampaignData.LevelsMilkLimit = jsonLevel.Metadata.MaxMilk
	bState.CampaignData.GeneratorActive = jsonLevel.Metadata.Generator
	bState.CampaignData.NextLevel = jsonLevel.Metadata.NextLevel
	bState.CampaignData.Name = jsonLevel.Metadata.Name

	// Ustawienie początkowej pozycji kamery
	bState.GameCamera.Target.X = float32(jsonLevel.Metadata.StartPos.X * tileWidth)
	bState.GameCamera.Target.Y = float32(jsonLevel.Metadata.StartPos.Y * tileHeight)
	log.Printf("INFO: Początkowa pozycja kamery ustawiona na (%f, %f) w kafelkach.",
		bState.GameCamera.Target.X, bState.GameCamera.Target.Y)

	// Ustawienia SI
	bState.AI.GatherPointX = jsonLevel.AISettings.GatherPoint.X
	bState.AI.GatherPointY = jsonLevel.AISettings.GatherPoint.Y
	bState.AI.PastureX = jsonLevel.AISettings.Pasture.X
	bState.AI.PastureY = jsonLevel.AISettings.Pasture.Y

	// Specjalne lokacje
	if jsonLevel.SpecialLocations.TransformationPoint != nil {
		bState.CampaignData.TransformationSiteX = jsonLevel.SpecialLocations.TransformationPoint.X
		bState.CampaignData.TransformationSiteY = jsonLevel.SpecialLocations.TransformationPoint.Y
	}

	if jsonLevel.SpecialLocations.VictoryPoint != nil {
		bState.CampaignData.VictoryPointX = jsonLevel.SpecialLocations.VictoryPoint.X
		bState.CampaignData.VictoryPointY = jsonLevel.SpecialLocations.VictoryPoint.Y
	}

	if jsonLevel.SpecialLocations.RescueTarget != nil {
		bState.CampaignData.RescueTargetX = jsonLevel.SpecialLocations.RescueTarget.X
		bState.CampaignData.RescueTargetY = jsonLevel.SpecialLocations.RescueTarget.Y
	}

	// Zastosuj surowy teren
	l.applyTerrain(&jsonLevel.Terrain, bState)

	// Zastosuj budynki
	l.applyBuildings(jsonLevel.Buildings, bState)

	// Zastosuj jednostki
	l.applyUnits(jsonLevel.Units, bState)

	// Zasoby graczy
	bState.HumanPlayerState.MaxMilk = jsonLevel.Metadata.MaxMilk
	bState.HumanPlayerState.Milk = jsonLevel.Metadata.MaxMilk
	bState.AIEnemyState.MaxMilk = 1800
	bState.AIEnemyState.Milk = 1800

	log.Printf("Zastosowano poziom: %s (%dx%d)",
		jsonLevel.Metadata.Name, jsonLevel.Terrain.Width, jsonLevel.Terrain.Height)
}

// Wyczyść stan gry
func (l *jsonLevelLoader) clearGameState(bState *battleState) {
	initBoard(bState)

	// POPRAWKA: Resetujemy długość do 0, ale zachowujemy pojemność (Capacity).
	// Zapobiega to utracie "rezerwacji" 40 miejsc i zbędnym alokacjom.
	bState.Units = bState.Units[:0]
	bState.Buildings = bState.Buildings[:0]
	// Jeśli projectiles są zainicjowane, też je czyścimy
	if bState.Projectiles != nil {
		bState.Projectiles = bState.Projectiles[:0]
	}

	// Zerujemy graczy
	bState.HumanPlayerState.init(bState.PlayerID, 0)
	bState.AIEnemyState.init(bState.AIPlayerID, 0)

	// ZEROWANIE LICZNIKÓW POPULACJI I BUDYNKÓW
	bState.HumanPlayerState.CurrentPopulation = 0
	bState.AIEnemyState.CurrentPopulation = 0
	bState.HumanPlayerState.CurrentBuildings = 0
	bState.AIEnemyState.CurrentBuildings = 0

	bState.CurrentSelection = selectionState{}
	bState.MouseState = mouseStateNormal
	bState.DragContext.IsActive = false

	// Zerowanie liczników animacji
	bState.GrassGrowthCycle = 0
	bState.WaterAnimationFrame = 0
	bState.FireAnimationFrame = 0
	bState.GlobalFrameCounter = 0

	bState.AI = aiState{}
	bState.AI.MilkGenerationRate = uint16(bState.DifficultyLevel)
}

func configureTile(currentTile *tile, graphicID uint16) {
	currentTile.IsWalkable = true
	currentTile.MovementCost = 1.0
	// ↓↓↓ Powinienem się zastanowić, co otwarcie deklarować
	// ↓↓↓ bo niektóre wartości są domyślnie ustawiane prawidłowo
	// currentTile.treeState = noTree

	if isWater(graphicID) {
		currentTile.IsWalkable = false
	}

	if isRockNonWalkable(graphicID) {
		currentTile.IsWalkable = false
	}

	if isDirtRoad(graphicID) {
		currentTile.MovementCost = 0.5
	}

	if isTreeStump(graphicID) || isRuin(graphicID) {
		currentTile.IsWalkable = false
	}

	if isGadget(graphicID) {
		switch graphicID {
		case spriteGadget02, spriteGadget03, spriteGadget13:
			currentTile.IsWalkable = true
		default:
			currentTile.IsWalkable = false
		}
	}

	if state, isBurnt, isWalkableOverride := classifyTreeFromTexture(graphicID); state != noTree {
		currentTile.treeState = state
		currentTile.IsBurnt = isBurnt
		currentTile.IsWalkable = isWalkableOverride
	}
}

// Przekłada tekstury używane w JSON-ie map na stan logiczny kafelka.
// Zwraca następujące informacje: treeState, isBurnt, isWalkable
// Dzięki temu wiemy jednoznacznie co powinno się dziać z drzewem.
func classifyTreeFromTexture(textureID uint16) (treeState, bool, bool) {
	switch textureID {
	// stojące żywe i suche
	case spriteTreeStump00, spriteTreeStump01, spriteTreeStump02, spriteTreeStump03,
		spriteTreeStump04, spriteTreeStump05, spriteDryTreeStump00:
		return treeStraight, false, false
	// stojące spalone
	case spriteTreeBurntStump00, spriteTreeBurntStump01:
		return treeStraight, true, false
	// suche upadające
	case spriteDryLeaningTreeStump:
		return treeLeaning, false, false
	case spriteDryFallingStump:
		return treeFalling, false, false
	case spriteDryFallenTreeStump:
		return treeFell, false, true
	// spalone upadające
	case spriteBurntLeaningTreeStump:
		return treeLeaning, true, false
	case spriteBurntFallingTreeStump:
		return treeFalling, true, false
	case spriteBurntFallenTreeStump:
		return treeFell, true, true
	// jeśli coś nie jest drzewem
	default:
		return noTree, false, true
	}
}

func (l *jsonLevelLoader) spawnPalisade(tileX, tileY uint8, graphicID uint16, bState *battleState) {
	if graphicID == spritePalisadeNE {
		newPalisade := &building{}
		newPalisade.initConstruction(buildingPalisade, colorNone, BuildingID(bState.NextUniqueObjectID))
		bState.NextUniqueObjectID++
		placeConstructionOnBoard(newPalisade, tileX, tileY, bState.Board)

		bState.Buildings = append(bState.Buildings, newPalisade)
	}
}

func (l *jsonLevelLoader) applyTerrain(terrain *jsonTerrainData, bState *battleState) {
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
			currentTile := &bState.Board.Tiles[tX][tY]

			// 3. Przypisanie grafiki
			currentTile.TextureID = graphicID

			// 4. Przetworzenie
			configureTile(currentTile, graphicID)
			l.spawnPalisade(tX, tY, graphicID, bState)
		}
	}

	log.Println("INFO: Teren nałożony pomyślnie.")
}

func (l *jsonLevelLoader) applyBuildings(buildingsData []jsonBuildingData, bState *battleState) {
	log.Printf("INFO: Ładowanie %d budynków do battleState...", len(buildingsData))

	for _, data := range buildingsData {
		var ownerID uint8

		switch data.Owner {
		case "ENEMY":
			ownerID = bState.AIPlayerID
		case "PLAYER":
			ownerID = bState.PlayerID
		default:
			ownerID = colorNone
		}

		bldType, exists := buildingTypeMap[data.Type]
		if !exists {
			continue
		}

		// Pobieramy definicję, aby znać wymiary (Width/Height)
		stats, ok := buildingDefs[bldType]
		if !ok {
			log.Printf("BŁĄD: Brak definicji statystyk dla budynku typu %d", bldType)

			continue
		}

		// Przeliczenie współrzędnych
		// JSON zawiera współrzędne Prawego-Dolnego rogu (stary system).
		// Funkcja init() wymaga teraz Lewego-Górnego rogu (nowy system).
		topLeftX := data.Position.X - stats.Width + 1
		topLeftY := data.Position.Y - stats.Height + 1

		newBuilding := &building{}

		// Wywołujemy init z przeliczonymi współrzędnymi Top-Left
		newBuilding.initConstruction(bldType, ownerID, BuildingID(bState.NextUniqueObjectID))
		bState.NextUniqueObjectID++

		placeConstructionOnBoard(newBuilding, topLeftX, topLeftY, bState.Board)

		bState.Buildings = append(bState.Buildings, newBuilding)

		// Liczniki
		switch newBuilding.Owner {
		case bState.HumanPlayerState.PlayerID:
			bState.HumanPlayerState.CurrentBuildings++
		case bState.AIEnemyState.PlayerID:
			bState.AIEnemyState.CurrentBuildings++
		}

		// Nakładanie grafiki przy użyciu wyliczonego Top-Left
		template, templateExists := buildingTemplates[bldType]
		if templateExists {
			for dy, row := range template {
				for dx, graphicID := range row {
					tileX := topLeftX + uint8(dx)
					tileY := topLeftY + uint8(dy)

					if tileX < boardMaxX && tileY < boardMaxY {
						bState.Board.Tiles[tileX][tileY].TextureID = uint16(graphicID)
					}
				}
			}
		}
	}
}

// levelLoader.go - applyUnits
func (l *jsonLevelLoader) applyUnits(units []jsonUnitData, bState *battleState) {
	log.Printf("INFO: Ładowanie %d jednostek do battleState...", len(units))

	for _, unitData := range units {
		var ownerID uint8
		if unitData.Owner == "ENEMY" {
			ownerID = bState.AIPlayerID
		} else {
			ownerID = bState.PlayerID
		}

		uType, exists := unitTypeMap[unitData.Type]
		if !exists {
			log.Printf("OSTRZEŻENIE: Pomięto nieznany typ jednostki: %s", unitData.Type)
			continue
		}

		newUnit := &unit{}
		// Inicjalizacja nowej jednostki
		newUnit.initUnit(uType, unitData.Position.X, unitData.Position.Y, cmdUIdle, bState)
		newUnit.Owner = ownerID

		// Wstawienie na mapę
		newUnit.show(bState.Board)

		// Dodanie do głównej listy
		bState.Units = append(bState.Units, newUnit)

		// Zliczamy jednostki startowe, aby limit 40 działał poprawnie od początku gry.
		switch newUnit.Owner {
		case bState.HumanPlayerState.PlayerID:
			bState.HumanPlayerState.CurrentPopulation++
		case bState.AIEnemyState.PlayerID:
			bState.AIEnemyState.CurrentPopulation++
		}
	}
	log.Printf("INFO: Załadowano %d jednostek. Populacja Gracza: %d, AI: %d.",
		len(bState.Units), bState.HumanPlayerState.CurrentPopulation, bState.AIEnemyState.CurrentPopulation)
}

// Załaduj poziom
func (l *jsonLevelLoader) loadLevel(levelNum uint8, bState *battleState) error {
	if err := l.loadJSONLevel(levelNum, bState); err == nil {
		return nil
	}
	return fmt.Errorf("nie znaleziono poziomu JSON %d", levelNum)
}

func (l *jsonLevelLoader) loadJSONLevel(levelNum uint8, bState *battleState) error {
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
	l.applyJSONLevel(&jsonLevel, bState)

	return nil
}

// Waliduj pozycję kamery
func (l *jsonLevelLoader) validateScreenPosition(bState *battleState) {
	// Ile kafelków mieści się na ekranie (wirtualnym)
	tilesVisibleX := virtualScreenWidth / uint16(tileWidth)   // 640/16 = 40
	tilesVisibleY := virtualScreenHeight / uint16(tileHeight) // 400/14 = 28

	// Konwersja celu kamery z pikseli na kafelki
	cameraTargetTileX := uint8(bState.GameCamera.Target.X / float32(tileWidth))
	cameraTargetTileY := uint8(bState.GameCamera.Target.Y / float32(tileHeight))

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
	bState.GameCamera.Target.X = float32(cameraTargetTileX * tileWidth)
	bState.GameCamera.Target.Y = float32(cameraTargetTileY * tileHeight)

	log.Printf("DEBUG: Ostateczna pozycja celu kamery (skorygowana): (%d, %d) kafelków = (%.1f, %.1f) pikseli",
		cameraTargetTileX, cameraTargetTileY, bState.GameCamera.Target.X, bState.GameCamera.Target.Y)
}

// initBattle - metoda LevelLoader która inicjuje bitwę z ProgramState.
func (l *jsonLevelLoader) initBattle(levelNumber uint8, bState *battleState) error {
	log.Printf("Inicjalizacja bitwy poziom %d", levelNumber)
	// @todo: przekaż poziom trudności do „battlestate”!!

	initBoard(bState)

	if err := l.loadLevel(levelNumber, bState); err != nil {
		return fmt.Errorf("nie można załadować poziomu %d: %w", levelNumber, err)
	}

	processMapTiles(bState)

	l.validateScreenPosition(bState)

	log.Println("INFO: Bitwa rozpoczęta pomyślnie.")

	return nil
}
