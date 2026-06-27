package main

// cursor.go

import rl "github.com/gen2brain/raylib-go/raylib"

// getCursorIDFromContext - otulina
// Sprawdza stan gry. Jeśli gameScreen → deleguje do determineCursorState.
// Jeśli inny stan → zwraca domyślny kursor.
func getCursorIDFromContext(bState *battleState, pState *programState, realMousePos rl.Vector2, scale float32, iState inputState) uint16 {
	if pState.CurrentState == gameScreen {
		virtualMouseX := realMousePos.X / scale
		virtualMouseY := realMousePos.Y / scale
		virtualMousePos := rl.NewVector2(virtualMouseX, virtualMouseY)

		return determineCursorState(bState, virtualMousePos, pState.GameViewWidth, pState.VirtualWidth, pState.VirtualHeight, iState)
	}

	return spriteCursorPointer
}

func determineCursorState(bState *battleState, mousePos rl.Vector2, viewW, totalW, viewH float32, iState inputState) uint16 {
	// 1. Czy poza planszą
	if screenCursor := checkScreenCursor(mousePos, viewW, totalW, viewH); screenCursor != 0 {
		return screenCursor
	}

	// 2. Na planszy
	// Musimy przeliczyć pozycję myszy na świat gry używając kamery
	worldMousePos := rl.GetScreenToWorld2D(mousePos, bState.GameCamera)
	tileX := uint8(worldMousePos.X / float32(tileWidth))
	tileY := uint8(worldMousePos.Y / float32(tileHeight))

	if tileX >= boardMaxX || tileY >= boardMaxY {
		return spriteCursorStop
	}

	tileUnderCursor := &bState.Board.Tiles[tileX][tileY]
	targetOwner := -1

	var targetBuilding *building // Potrzebujemy wiedzieć, czy to budynek

	if tileUnderCursor.Unit != nil && tileUnderCursor.Unit.Exists {
		targetOwner = int(tileUnderCursor.Unit.Owner)
	} else if tileUnderCursor.Building != nil && tileUnderCursor.Building.Exists {
		targetOwner = int(tileUnderCursor.Building.Owner)
		targetBuilding = tileUnderCursor.Building
	}

	hasSelection := bState.CurrentSelection.IsUnit && bState.CurrentSelection.OwnerID == bState.PlayerID

	// @todo: przemyśl nazwy poniższych funkcji, bo są bardzo mylące.
	if hasSelection {
		return cursorForSelection(bState, tileUnderCursor, targetOwner, targetBuilding, iState)
	}

	return cursorForNoSelection(targetOwner, bState.PlayerID)
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

func cursorForSelection(bState *battleState, tileUnderCursor *tile, targetOwner int, targetBuilding *building, iState inputState) uint16 {
	// Naprawa
	if bState.MouseState == mouseStateRepairing {
		if canRepair(targetBuilding, bState.PlayerID) {
			return spriteBtnRepair
		}

		return spriteCursorStop
	}

	// Gromobicie/deszcz ognia
	// @todo: sprawdź, czy w pierwowzorze można było tym zaatakować wszystko - 22.06.2026
	if bState.MouseState == mouseStateCasting {
		return spriteCursorCrossRed
	}

	// Wróg
	if targetOwner != -1 && targetOwner != int(bState.PlayerID) {
		return cursorForEnemy(bState, tileUnderCursor)
	}

	// Swój
	if targetOwner == int(bState.PlayerID) {
		// chyba tutaj powinienem dodać warunek dla „szybkiej budowy”
		selectedUnit, ok := getUnitByID(bState.CurrentSelection.UnitID, bState)
		if ok && selectedUnit.Type == unitAxeman && tileUnderCursor.Building != nil && tileUnderCursor.Building.IsUnderConstruction {
			// co jeśli nie jest to prawidłowa zwrotka?!
			return spriteBtnRepair
		}

		if !iState.IsCtrlKeyDown {
			return spriteCursorFrameWhite
		}

		isSelfAttack := tileUnderCursor.Unit != nil && tileUnderCursor.Unit.IsSelected

		if isSelfAttack {
			return spriteCursorStop
		}

		return spriteCursorCrossRed
	}

	// Drzewa
	selectedUnit, unitOK := getUnitByID(bState.CurrentSelection.UnitID, bState)

	if unitOK && tileUnderCursor.isStandingTree() && !tileUnderCursor.IsBurning {
		if selectedUnit.canDamageTree(tileUnderCursor.X, tileUnderCursor.Y, bState) {
			return spriteCursorCrossRed
		}

		return spriteCursorStop
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

// @todo: do przebudowy, bo sprawdzanie legalności w kursorach jest oderwane od legalności rozkazów.
func cursorForEnemy(bState *battleState, tileUnderCursor *tile) uint16 {
	targetBuilding := tileUnderCursor.Building

	if targetBuilding != nil && targetBuilding.Exists {
		// Mag nie może atakować żadnych budynków
		selectedUnit, ok := getUnitByID(bState.CurrentSelection.UnitID, bState)

		if ok && selectedUnit.Type == unitMage {
			return spriteCursorStop
		}

		if targetBuilding.Type == buildingPalisade {
			if !targetBuilding.IsUnderConstruction {
				if ok && !canDamagePalisades(selectedUnit) {
					return spriteCursorStop
				}
			} else {
				return spriteCursorCrossWhite
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
