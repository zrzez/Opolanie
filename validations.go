package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// ==============
// SPRAWDZANIE POPRAWNOŚCI ROZKAZÓW DLA BUDYNKÓW
// ==============

// Stałe związane z tymczasowym układem zwracania informacji o błędach
const (
	buildErrNone uint8 = iota
	buildErrLimit
	buildErrMilk
	buildErrOutofBounds
	buildErrOccupiedUnit
	buildErrOccupiedBuilding
	buildErrObstacle
	buildErrWater
	buildErrNoWater
	buildErrNoRoadAccess
	buildErrAlreadyPath
	buildErrNoDefinition
	workErrNone
	workErrWrongWorkerType
	workErrInvalidTarget
	workErrNotUnderConstruction
	workErrNotRepairable
)

// ========
// KONTEKSTOWE
// =======
func validateConstructionContext(bType buildingType, owner uint8, bState *battleState) (bool, uint8) {
	stats, ok := buildingDefs[bType]

	if !ok {
		return false, buildErrNoDefinition // nigdy nie powinno się wydarzyć
	}

	regularBuilding := bType != buildingPalisade && bType != buildingBridge && bType != buildingRoad

	ownerState := bState.getPlayerState(owner)

	if regularBuilding && ownerState.CurrentBuildings >= maxBuildingsPerPlayer {
		return false, buildErrLimit
	}

	if ownerState.Milk < stats.Cost {
		return false, buildErrMilk
	}

	return true, buildErrNone
}

func validateConstructionSite(bType buildingType, tileX, tileY uint8, bState *battleState) (bool, uint8) {
	stats, ok := buildingDefs[bType]

	if !ok {
		return false, buildErrNoDefinition // nigdy nie powinno mieć to miejsca
	}

	textureID := bState.Board.Tiles[tileX][tileY].TextureID

	switch bType {
	case buildingRoad:
		if isPath(textureID) {
			return false, buildErrAlreadyPath
		}
		if isWater(textureID) {
			return false, buildErrWater
		}
		if !isWalkable(bState, tileX, tileY) {
			return false, buildErrObstacle
		}
		if !hasRoadAccess(tileX, tileY, smallBuildingSize, bState) {
			return false, buildErrNoRoadAccess
		}
	case buildingBridge:
		if !isWater(textureID) {
			return false, buildErrNoWater
		}
		if !hasRoadAccess(tileX, tileY, smallBuildingSize, bState) {
			return false, buildErrNoRoadAccess
		}
	default:
		valid, err := isValidConstructionSite(tileX, tileY, stats.Width, stats.Height, bType, bState)

		if !valid {
			return false, err
		}
	}

	return true, buildErrNone
}

// isValidConstructionSite sprawdza wszelkie warunki, które należy spełnić, aby można było zasadzić budowlę.
func isValidConstructionSite(tileX, tileY, width, height uint8, bType buildingType, bState *battleState) (bool, uint8) {
	can, err := canFitBuilding(tileX, tileY, width, height, bType, bState)

	if !can {
		return false, err
	}

	if bType != buildingPalisade && !hasRoadAccess(tileX, tileY, width, bState) {
		return false, buildErrNoRoadAccess
	}

	return true, buildErrNone
}

func canFitBuilding(tileX, tileY, width, height uint8, bType buildingType, bState *battleState) (bool, uint8) {
	for dx := range width {
		for dy := range height {
			// Dodajemy przesunięcie do kursora.
			constructionX, constructionY := tileX+dx, tileY+dy

			ok, err := isWithinBoard(constructionX, constructionY)

			if !ok {
				return false, err
			}

			free, err := isFreeForConstruction(constructionX, constructionY, bType, bState)
			if !free {
				return false, err
			}
		}
	}

	return true, buildErrNone
}

func isWithinBoard(constructionX, constructionY uint8) (bool, uint8) {
	// @reminder: usunąłem sprawdzenie constructionX i constructionY < 0
	if constructionX >= boardMaxX || constructionY >= boardMaxY {
		return false, buildErrOutofBounds
	}

	return true, buildErrNone
}

func isFreeForConstruction(constructionX, constructionY uint8, bType buildingType, bState *battleState) (bool, uint8) {
	currentTile := &bState.Board.Tiles[constructionX][constructionY]

	// Czy coś tu stoi?
	if currentTile.Unit != nil {
		return false, buildErrOccupiedUnit
	}

	if currentTile.Building != nil {
		existingIsUnfinishedPalisade := currentTile.Building.Type == buildingPalisade && currentTile.Building.IsUnderConstruction
		buildingNonPalisade := bType != buildingPalisade

		// Zezwalamy na budowę TYLKO w jednym przypadku:
		// istniejąca nieukończona palisada + budujemy coś innego niż palisadę
		if existingIsUnfinishedPalisade && buildingNonPalisade {
			// OK – nie blokujemy, spadamy do sprawdzenia terenu
		} else {
			// We wszystkich innych przypadkach budynek blokuje miejsce
			return false, buildErrOccupiedBuilding
		}
	}

	// Czy teren nadaje się pod budowę?
	if isObstacle(currentTile.TextureID) {
		return false, buildErrObstacle
	}

	return true, buildErrNone
}

func hasRoadAccess(x, y, size uint8, bState *battleState) bool {
	// 1. Sprawdzamy boki LEWY i PRAWY (iterujemy po wysokości budynku)
	for i := uint8(0); i < size; i++ {
		currentY := y + i

		// Zabezpieczenie, gdyby budynek wystawał poza mapę dołem (teoretycznie niemożliwe po walidacji, ale bezpiecznie)
		if currentY >= boardMaxY {
			break
		}

		// --- LEWA KRAWĘDŹ (x-1) ---
		// Sprawdzamy tylko, jeśli nie jesteśmy przytuleni do lewej ściany mapy (x=0)
		// To zapobiega błędowi "uint8 overflow" (0 - 1 = 255)
		if x > 0 {
			if isPath(bState.Board.Tiles[x-1][currentY].TextureID) {
				return true
			}
		}

		// --- PRAWA KRAWĘDŹ (x+size) ---
		if x+size < boardMaxX {
			if isPath(bState.Board.Tiles[x+size][currentY].TextureID) {
				return true
			}
		}
	}

	// 2. Sprawdzamy boki GÓRA i DÓŁ (iterujemy po szerokości budynku)
	for i := uint8(0); i < size; i++ {
		currentX := x + i

		if currentX >= boardMaxX {
			break
		}

		// --- GÓRNA KRAWĘDŹ (y-1) ---
		// Sprawdzamy tylko, jeśli nie jesteśmy na samej górze mapy (y=0)
		if y > 0 {
			if isPath(bState.Board.Tiles[currentX][y-1].TextureID) {
				return true
			}
		}

		// --- DOLNA KRAWĘDŹ (y+size) ---
		if y+size < boardMaxY {
			if isPath(bState.Board.Tiles[currentX][y+size].TextureID) {
				return true
			}
		}
	}

	return false
}

// ====
// GRAFICZNA ZWROTKA WALIDACJI ZASADZANIA BUDOWY
func drawConstructionValidationBox(bState *battleState, pState *programState) {
	if bState.MouseState != mouseStatePlaceConstruction {
		return
	}

	screenMouse := rl.GetMousePosition()
	virtualMouse := screenToVirtualCoords(pState, screenMouse)

	if isMouseOverUI(pState, virtualMouse) {
		return
	}

	worldMousePos := rl.GetScreenToWorld2D(virtualMouse, bState.GameCamera)

	startX := uint8(worldMousePos.X / float32(tileWidth))
	startY := uint8(worldMousePos.Y / float32(tileHeight))

	size := buildingDefs[buildingType(bState.PendingCommand.CreateType)].Width

	var tileStates [3][3]uint8
	var hasRoadAccessAnywhere bool

	for dy := uint8(0); dy < size; dy++ {
		for dx := uint8(0); dx < size; dx++ {
			cx := startX + dx
			cy := startY + dy

			rawColor := validationBoxColor(cx, cy, bState)

			var state uint8

			if rawColor == rl.Red {
				state = 0
			} else if rawColor == rl.Orange {
				state = 1
			} else {
				state = 2
				hasRoadAccessAnywhere = true
			}

			tileStates[dy][dx] = state
		}
	}

	for dy := uint8(0); dy < size; dy++ {
		for dx := uint8(0); dx < size; dx++ {
			state := tileStates[dy][dx]
			finalColor := rl.Red

			switch state {
			case 0:
				finalColor = rl.Red
			case 1:
				if hasRoadAccessAnywhere {
					finalColor = rl.Green
				} else {
					finalColor = rl.Orange
				}
			case 2:
				finalColor = rl.DarkGreen
			}

			finalColor = rl.Fade(finalColor, 0.3)

			posX := float32(startX+dx) * float32(tileWidth)
			posY := float32(startY+dy) * float32(tileHeight)

			rl.DrawRectangle(int32(posX), int32(posY), int32(tileWidth), int32(tileHeight), finalColor)
		}
	}
}

func validationBoxColor(tileX, tileY uint8, bState *battleState) rl.Color {
	// 1. Walidacja specyficzna dla typu budynku
	var isValid bool

	switch buildingType(bState.PendingCommand.CreateType) {
	case buildingBridge:
		isOnBoard, _ := isWithinBoard(tileX, tileY)
		isValid = isOnBoard && isWater(bState.Board.Tiles[tileX][tileY].TextureID)

	case buildingRoad:
		canFit, _ := canFitBuilding(tileX, tileY, smallBuildingSize, smallBuildingSize, buildingType(bState.PendingCommand.CreateType), bState)
		isValid = !isPath(bState.Board.Tiles[tileX][tileY].TextureID) && canFit

	default:
		isValid, _ = canFitBuilding(tileX, tileY, smallBuildingSize, smallBuildingSize, buildingType(bState.PendingCommand.CreateType), bState)
	}

	// 2. Jeśli miejsce jest nieważne, od razu zwracamy czerwień
	if !isValid {
		return rl.Red
	}

	// 3. Palisada nie wymaga dostępu do drogi, więc od razu jest zielona
	if buildingType(bState.PendingCommand.CreateType) == buildingPalisade {
		return rl.DarkGreen
	}

	// 4. Budynki wymagają drogi
	if hasRoadAccess(tileX, tileY, smallBuildingSize, bState) {
		return rl.DarkGreen
	}

	// 5. Ważne, ale bez dostępu do drogi
	return rl.Orange
}

// ====
// Budowanie i naprawianie
// ====
func validateBuildingContext(u *unit, targetBld *building) (bool, uint8) {
	// 1. Czy jednostka ma prawo budować?
	if u.Type != unitAxeman {
		return false, workErrWrongWorkerType
	}

	// 2. Czy właściwa budowa?
	if targetBld == nil || !targetBld.Exists {
		return false, workErrInvalidTarget
	}

	isNeutral := targetBld.Type == buildingBridge || targetBld.Type == buildingPalisade

	if !isNeutral && targetBld.Owner != u.Owner {
		return false, workErrInvalidTarget
	}

	// 3. Czy to plac budowy?
	if !targetBld.IsUnderConstruction {
		return false, workErrNotUnderConstruction
	}

	return true, workErrNone
}

func validateRepairContext(u *unit, targetBld *building) (bool, uint8) {
	// 1. Czy jednostka ma prawo naprawiać?
	if u.Type != unitAxeman {
		return false, workErrWrongWorkerType
	}

	// 2. Czy właściwy budynek?
	if targetBld == nil || !targetBld.Exists {
		return false, workErrInvalidTarget
	}

	isNeutral := targetBld.Type == buildingBridge || targetBld.Type == buildingPalisade

	if !isNeutral && targetBld.Owner != u.Owner {
		return false, workErrInvalidTarget
	}

	// 3. Czy naprawialny?
	if !targetBld.isRepairable(u.Owner) {
		return false, workErrNotRepairable
	}

	return true, workErrNone
}

func isValidWorkTarget(targetBld *building, playerID uint8) bool {
	if targetBld == nil || !targetBld.Exists {
		return false
	}

	isNeutral := targetBld.Type == buildingBridge || targetBld.Type == buildingPalisade

	if !isNeutral && targetBld.Owner != playerID {
		return false
	}

	return targetBld.IsUnderConstruction || targetBld.isRepairable(playerID)
}
