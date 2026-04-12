package main

// cursor.go

import rl "github.com/gen2brain/raylib-go/raylib"

// getCursorIDFromContext - otulina
// Sprawdza stan gry. Jeśli gameScreen → deleguje do determineCursorState.
// Jeśli inny stan → zwraca domyślny kursor.
func getCursorIDFromContext(bs *battleState, ps *programState, realMousePos rl.Vector2, scale float32) uint16 {
	if ps.CurrentState == gameScreen {
		virtualMouseX := realMousePos.X / scale
		virtualMouseY := realMousePos.Y / scale
		virtualMousePos := rl.NewVector2(virtualMouseX, virtualMouseY)

		return determineCursorState(bs, virtualMousePos, ps.GameViewWidth, ps.VirtualWidth, ps.VirtualHeight)
	}

	return spriteCursorPointer
}

func determineCursorState(bs *battleState, mousePos rl.Vector2, viewW, totalW, viewH float32) uint16 {
	// 1. Czy poza planszą
	if screenCursor := checkScreenCursor(mousePos, viewW, totalW, viewH); screenCursor != 0 {
		return screenCursor
	}

	// 2. Na planszy
	// Musimy przeliczyć pozycję myszy na świat gry używając kamery
	worldMousePos := rl.GetScreenToWorld2D(mousePos, bs.GameCamera)
	tileX := uint8(worldMousePos.X / float32(tileWidth))
	tileY := uint8(worldMousePos.Y / float32(tileHeight))

	if tileX >= boardMaxX || tileY >= boardMaxY {
		return spriteCursorStop
	}

	tileUnderCursor := &bs.Board.Tiles[tileX][tileY]
	targetOwner := -1

	var targetBuilding *building // Potrzebujemy wiedzieć, czy to budynek

	if tileUnderCursor.Unit != nil && tileUnderCursor.Unit.Exists {
		targetOwner = int(tileUnderCursor.Unit.Owner)
	} else if tileUnderCursor.Building != nil && tileUnderCursor.Building.Exists {
		targetOwner = int(tileUnderCursor.Building.Owner)
		targetBuilding = tileUnderCursor.Building
	}

	hasSelection := bs.CurrentSelection.IsUnit && bs.CurrentSelection.OwnerID == bs.PlayerID

	if hasSelection {
		return cursorForSelection(bs, tileUnderCursor, targetOwner, targetBuilding)
	}

	return cursorForNoSelection(targetOwner, bs.PlayerID)
}

// Zwraca zero dla kursora na planszy lub odpowiedni duszek w pozostałych przypadkach.
func checkScreenCursor(mousePos rl.Vector2, viewW, totalW, viewH float32) uint16 {
	// Sprawdzamy przewijanie
	if mousePos.X > totalW-scrollZoneXThreshold {
		return spriteCursorArrowRight
	}

	if mousePos.X < scrollZoneXThreshold {
		return spriteCursorArrowLeft
	}

	if mousePos.Y < scrollZoneYThreshold {
		return spriteCursorArrowUp
	}

	if mousePos.Y > viewH-scrollZoneYThreshold {
		return spriteCursorArrowDown
	}

	// Sprawdzamy, czy mysz jest nad UI
	if mousePos.X >= viewW {
		return spriteCursorDefaultBig
	}

	// kursor jest nad planszą więc zwracam zero
	return 0
}

func cursorForSelection(bs *battleState, tileUnderCursor *tile, targetOwner int, targetBuilding *building) uint16 {
	// Naprawa
	if bs.MouseCommandMode == cmdRepairStructure {
		if canRepair(targetBuilding, bs.PlayerID) {
			return spriteBtnRepair
		}

		return spriteCursorStop
	}

	// Wróg
	if targetOwner != -1 && targetOwner != int(bs.PlayerID) {
		return cursorForEnemy(bs, targetBuilding)
	}

	// Swój
	if targetOwner == int(bs.PlayerID) {
		return spriteCursorFrameWhite
	}

	// Puste pole
	if !tileUnderCursor.IsWalkable {
		return spriteCursorStop
	}

	return spriteCursorCrossWhite
}

func cursorForNoSelection(targetOwner int, playerID uint8) uint16 {
	// Wróg
	if targetOwner != -1 && targetOwner != int(playerID) {
		return spriteCursorFrameRed
	}

	// Swój
	if targetOwner == int(playerID) {
		return spriteCursorFrameWhite
	}

	// Pusto
	return spriteCursorDefaultBig
}

func cursorForEnemy(bs *battleState, targetBuilding *building) uint16 {
	if targetBuilding != nil {
		if targetBuilding.Type == buildingPalisade {
			selectedUnit, ok := getUnitByID(bs.CurrentSelection.UnitID, bs)

			if ok && !canDamagePalisades(selectedUnit) {
				return spriteCursorStop
			}
		}

		if targetBuilding.Type == buildingBridge {
			return spriteCursorStop
		}
	}

	return spriteCursorCrossRed
}

func canRepair(bld *building, playerID uint8) bool {
	if bld == nil || bld.HP >= bld.MaxHP {
		return false
	}

	return bld.Type == buildingPalisade || bld.Type == buildingBridge || bld.Owner == playerID
}

func animateCursorID(cursorID uint16) uint16 {
	const cursorAnimSpeed = 4.0
	animPhase := int(rl.GetTime()*cursorAnimSpeed) % 2

	if cursorID == spriteCursorCrossWhite {
		if rl.IsMouseButtonDown(rl.MouseRightButton) {
			return spriteCursorSmallWhite
		}

		if animPhase == 1 {
			return spriteCursorCrossMedWhite
		}
	}

	if cursorID == spriteCursorCrossRed {
		if rl.IsMouseButtonDown(rl.MouseRightButton) {
			return spriteCursorCrossMedRed
		}

		if animPhase == 1 {
			return spriteCursorCrossMedRed
		}
	}

	return cursorID
}
