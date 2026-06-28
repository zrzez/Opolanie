package main

// input.go

import (
	"fmt"
	"log"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// handleCheats obsługuje wpisywanie i uruchomienie oszustw.
//func handleCheats(bs *battleState) bool {
//	if !bs.IsSinglePlayerGame {
//		return false
//	}
//
//	keyPressed := rl.GetKeyPressed()
//	if keyPressed == 0 {
//		return false
//	}
//
//	// Sprawdzamy, czy odblokować oszustwa
//	if !bs.CheatsEnabled {
//		switch keyPressed {
//		case rl.KeyK:
//			if bs.CheatSequenceProgress == 0 {
//				bs.CheatSequenceProgress = 1
//				return true
//			}
//		case rl.KeyO:
//			if bs.CheatSequenceProgress == 1 {
//				bs.CheatSequenceProgress = 2
//				return true
//			}
//		case rl.KeyD:
//			if bs.CheatSequenceProgress == 2 {
//				bs.CheatSequenceProgress = 3
//				return true
//			}
//		case rl.KeyY:
//			if bs.CheatSequenceProgress == 3 {
//				bs.CheatsEnabled = true
//				bs.CheatSequenceProgress = 0
//				log.Println("OSZUSTWA ODBLOKOWANE")
//				return true
//			}
//		}
//		bs.CheatSequenceProgress = 0
//		return false
//	}
//
//	// Jeżeli oszukiwanie jest dozwolone
//	switch bs.CheatSequenceProgress {
//	case 0: // Oczekiwanie na pierwszy znak
//		switch keyPressed {
//		case rl.KeyD:
//			bs.CheatSequenceProgress = 10
//			return true // DOSW
//		case rl.KeyM:
//			bs.CheatSequenceProgress = 20
//			return true // MAGIC lub MILK
//		case rl.KeyK:
//			bs.CheatSequenceProgress = 40
//			return true // KILL
//		case rl.KeyT:
//			bs.CheatSequenceProgress = 50
//			return true // TREE
//		case rl.KeyE:
//			bs.CheatSequenceProgress = 60
//			return true // ENDV/ENDL
//		case rl.KeyC:
//			bs.CheatSequenceProgress = 70
//			return true // COUNT
//		case rl.KeyS:
//			bs.CheatSequenceProgress = 80
//			return true // SHOW
//		}
//
//	// Ciąg "DOSW"
//	case 10:
//		if keyPressed == rl.KeyO {
//			bs.CheatSequenceProgress = 11
//			return true
//		}
//	case 11:
//		if keyPressed == rl.KeyS {
//			bs.CheatSequenceProgress = 12
//			return true
//		}
//	case 12:
//		if keyPressed == rl.KeyW {
//			selectedUnit := getUnitByID(bs.CurrentSelection.UnitID, bs)
//			if selectedUnit != nil && bs.CurrentSelection.IsUnit {
//				selectedUnit.Experience = 100
//				log.Println("OSZUSTWO: Doświadczenie dla jednostki ustawione na 100")
//			}
//			bs.CheatSequenceProgress = 0
//			return true
//		}
//
//	// Ciąg zaczynające się od 'M' (MAGIC i MILK)
//	case 20:
//		if keyPressed == rl.KeyA {
//			bs.CheatSequenceProgress = 21
//			return true
//		}
//		if keyPressed == rl.KeyI {
//			bs.CheatSequenceProgress = 30
//			return true
//		}
//
//	// Gałąź "MAGIC"
//	case 21:
//		if keyPressed == rl.KeyG {
//			bs.CheatSequenceProgress = 22
//			return true
//		}
//	case 22:
//		if keyPressed == rl.KeyI {
//			bs.CheatSequenceProgress = 23
//			return true
//		}
//	case 23:
//		if keyPressed == rl.KeyC {
//			selectedUnit := getUnitByID(bs.CurrentSelection.UnitID, bs)
//			if selectedUnit != nil && bs.CurrentSelection.IsUnit {
//				t := selectedUnit.Type
//				if t == unitPriestess || t == unitPriest || t == unitUnknown || t == unitMage {
//					selectedUnit.Mana = 300
//					log.Println("OSZUSTWO: Magia dla jednostki ustawiona na 300")
//				}
//			}
//			bs.CheatSequenceProgress = 0
//			return true
//		}
//
//	// Gałąź "MILK"
//	case 30:
//		if keyPressed == rl.KeyL {
//			bs.CheatSequenceProgress = 31
//			return true
//		}
//	case 31:
//		if keyPressed == rl.KeyK {
//			bs.HumanPlayerState.Milk = bs.HumanPlayerState.MaxMilk
//			log.Println("OSZUSTWO: Mleko ustawione na max")
//			bs.CheatSequenceProgress = 0
//			return true
//		}
//
//	// Ciąg "KILL"
//	case 40:
//		if keyPressed == rl.KeyI {
//			bs.CheatSequenceProgress = 41
//			return true
//		}
//	case 41:
//		if keyPressed == rl.KeyL {
//			bs.CheatSequenceProgress = 42
//			return true
//		}
//	case 42:
//		if keyPressed == rl.KeyL {
//			selectedUnit := getUnitByID(bs.CurrentSelection.UnitID, bs)
//			if selectedUnit != nil && bs.CurrentSelection.IsUnit {
//				selectedUnit.HP = 0 // Zabicie jednostki
//				selectedUnit.Exists = false
//
//				tile := &bs.Board.Tiles[selectedUnit.X][selectedUnit.Y]
//				tile.Unit = nil
//				tile.EffectID = 127 // Efekt śmierci
//
//				log.Println("OSZUSTWO: KILL na zaznaczonej jednostce")
//			}
//			bs.CheatSequenceProgress = 0
//			return true
//		}
//
//	// Ciąg "TREE"
//	case 50:
//		if keyPressed == rl.KeyR {
//			bs.CheatSequenceProgress = 51
//			return true
//		}
//	case 51:
//		if keyPressed == rl.KeyE {
//			bs.CheatSequenceProgress = 52
//			return true
//		}
//	case 52:
//		// Powinno się zamienić zaznaczoną jednostkę w drzewo
//		if keyPressed == rl.KeyE {
//			selectedUnit := getUnitByID(bs.CurrentSelection.UnitID, bs)
//			if selectedUnit != nil && bs.CurrentSelection.IsUnit {
//				tx, ty := selectedUnit.X, selectedUnit.Y
//				tile := &bs.Board.Tiles[tx][ty]
//
//				if tile.TextureID <= spriteGrassEnd {
//					selectedUnit.Exists = false
//					tile.EffectID = 100 // @todo: jakaś czarodziejska liczba. Do wywalenia!
//					tile.Unit = nil     // Usuwamy jednostkę z mapy
//
//					tile.TextureID = spriteTreeStump00 + (tx & 3) + (ty & 3)
//					tile.IsWalkable = false
//
//					log.Println("OSZUSTWO: Jednostka zamieniona w drzewo")
//				}
//			}
//			bs.CheatSequenceProgress = 0
//			return true
//		}
//
//	// Ciąg "ENDV" / "ENDL"
//	case 60:
//		if keyPressed == rl.KeyN {
//			bs.CheatSequenceProgress = 61
//			return true
//		}
//	case 61:
//		if keyPressed == rl.KeyD {
//			bs.CheatSequenceProgress = 62
//			return true
//		}
//	case 62:
//		if keyPressed == rl.KeyV {
//			bs.BattleOutcome = 1
//			log.Println("OSZUSTWO: Zakończ poziom - Wygrana")
//			bs.CheatSequenceProgress = 0
//			return true
//		}
//		if keyPressed == rl.KeyL {
//			bs.BattleOutcome = 2
//			log.Println("OSZUSTWO: Zakończ poziom - Przegrana")
//			bs.CheatSequenceProgress = 0
//			return true
//		}
//
//	// Ciąg "COUNT"
//	case 70:
//		if keyPressed == rl.KeyO {
//			bs.CheatSequenceProgress = 71
//			return true
//		}
//	case 71:
//		if keyPressed == rl.KeyU {
//			bs.CheatSequenceProgress = 72
//			return true
//		}
//	case 72:
//		if keyPressed == rl.KeyN {
//			bs.CheatSequenceProgress = 73
//			return true
//		}
//	case 73:
//		if keyPressed == rl.KeyT {
//			bs.CheatsEnabled = !bs.CheatsEnabled
//			log.Printf("OSZUSTWO: Tryb debugowania przełączony na: %v", bs.CheatsEnabled)
//			bs.CheatSequenceProgress = 0
//			return true
//		}
//
//	// Ciąg "SHOW"
//	case 80:
//		if keyPressed == rl.KeyH {
//			bs.CheatSequenceProgress = 81
//			return true
//		}
//	case 81:
//		if keyPressed == rl.KeyO {
//			bs.CheatSequenceProgress = 82
//			return true
//		}
//	case 82:
//		if keyPressed == rl.KeyW {
//			for i := 0; i < boardMaxX; i++ {
//				for j := 0; j < boardMaxY; j++ {
//					bs.Board.Tiles[i][j].Visibility = visibilityVisible
//				}
//			}
//			log.Println("OSZUSTWO: Odkryj całą mapę")
//			bs.CheatSequenceProgress = 0
//			return true
//		}
//	}
//
//	bs.CheatSequenceProgress = 0
//	return true
//}

// OBSŁUGA WEJŚCIA - UI I INTERFEJS

func pollInput(pState *programState) inputState {
	screenMouse := rl.GetMousePosition()
	virtualMouse := screenToVirtualCoords(pState, screenMouse)

	return inputState{
		MousePosition:              virtualMouse,
		IsLeftMouseButtonDown:      rl.IsMouseButtonDown(rl.MouseButtonLeft),
		IsLeftMouseButtonPressed:   rl.IsMouseButtonPressed(rl.MouseButtonLeft),
		IsLeftMouseButtonReleased:  rl.IsMouseButtonReleased(rl.MouseButtonLeft),
		IsRightMouseButtonDown:     rl.IsMouseButtonDown(rl.MouseButtonRight),
		IsRightMouseButtonPressed:  rl.IsMouseButtonPressed(rl.MouseButtonRight),
		IsRightMouseButtonReleased: rl.IsMouseButtonReleased(rl.MouseButtonRight),
		IsCtrlKeyDown:              rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl),
		IsShiftKeyDown:             rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift),
	}
}

func handleGameUIClicks(iState inputState, bState *battleState, pState *programState) bool {
	virtualMouse := iState.MousePosition

	// 1. Najpierw sprawdzamy minimapę (to kluczowe, by nie klikać "przez" nią)
	if handleMinimapInteraction(iState, bState, pState) {
		return true
	}

	// 2. Obsługa kliknięć w przyciski akcji (panel boczny)
	// @todo: wydaje mi się, ze odwrócenie warunków powinno wypłaszyć tego ifoludka.
	if iState.IsLeftMouseButtonPressed {
		// Sprawdzamy każdy z 5 przycisków
		for btnIndex := range uiActionMaxButtons {
			action := bState.UI.CurrentActions[btnIndex]

			// Pomijamy wyłączone przyciski
			if !action.IsActive {
				continue
			}

			rect := bState.UI.ActionButtons[btnIndex]
			if rl.CheckCollisionPointRec(virtualMouse, rect) {
				log.Printf("UI: Wybrano akcję przycisku %d: %s (Typ: %d)", btnIndex, action.Label, action.Cmd.ActionType)

				bState.MouseState = action.State // czy potrzebuję do tego ifa?
				// Na podstawie stanu kursora ustalamy, jak rozumieć kliknięcie
				switch action.State {
				case mouseStatePlaceConstruction:
					// Zapisujemy rodzaj budynku do "plecaka" w battleState
					bState.PendingBuildingType = buildingType(action.Cmd.InteractionTargetID)

					bState.CurrentMessage.Text = "Wskaż miejsce"
					bState.CurrentMessage.Duration = 60

					// Opcjonalnie: czyścimy zaznaczenie jednostek, by nie przeszkadzały
					clearSelection(bState)
				case mouseStateWorking:
					bState.CurrentMessage.Text = "Wskaż budynek do naprawy"
					bState.CurrentMessage.Duration = 60

				case mouseStateCasting:
					bState.CurrentCommands[0] = action.Cmd
					bState.CurrentMessage.Text = "Wskaż cel czaru"
					bState.CurrentMessage.Duration = 60

				default:
					// Rozkazy natychmiastowe (CMD_PRODUCE, CMD_MILK, itp.)
					// Wysyłamy je bezpośrednio do kolejki rozkazów.
					bState.CurrentCommands[0] = action.Cmd

					// bezpiecznik na wypadek gdybym coś przeoczył
					bState.MouseState = mouseStateNormal
				}

				return true
			}
		}

		// Logowanie kliknięcia w tło UI (pomocne przy debugowaniu wymiarów)
		relativeX := virtualMouse.X - pState.GameViewWidth
		log.Printf("Klik w UI (tło) na: %.1f, %.1f", relativeX, virtualMouse.Y)

		return true
	}

	return false
}

func handleGameShortcuts(bState *battleState) bool {
	if rl.IsKeyPressed(rl.KeyKpSubtract) || rl.IsKeyPressed(rl.KeyMinus) {
		if bState.GameSpeed < 5 {
			bState.GameSpeed--
			log.Printf("SKRÓT: Prędkość gry zmniejszona do %d", bState.GameSpeed)
			return true
		}
	}

	if rl.IsKeyPressed(rl.KeyKpAdd) || rl.IsKeyPressed(rl.KeyEqual) {
		if bState.GameSpeed > 0 {
			bState.GameSpeed++
			log.Printf("SKRÓT: Prędkość gry zwiększona do %d", bState.GameSpeed)
			return true
		}
	}

	isShiftDown := rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)

	for i := 0; i <= 9; i++ {
		key := rl.KeyZero + int32(i)

		if rl.IsKeyPressed(key) {
			if isShiftDown {
				if bState.CurrentSelection.OwnerID != bState.PlayerID {
					return true
				}

				log.Printf("SKRÓT: Zapamiętywanie grupy %d", i)
				var selectedUnitIDs []uint
				for _, currentUnit := range bState.Units {
					if currentUnit.Exists && currentUnit.Owner == bState.PlayerID && currentUnit.IsSelected {
						selectedUnitIDs = append(selectedUnitIDs, currentUnit.ID)
					}
				}
				bState.ControlGroups[i] = controlGroup{}
				for _, unitID := range selectedUnitIDs {
					bState.ControlGroups[i].Units = append(bState.ControlGroups[i].Units, controlGroupUnit{UnitID: unitID})
				}
				log.Printf("SKRÓT: Grupa %d utworzona z %d jednostek.", i, len(bState.ControlGroups[i].Units))
				return true

			}

			if len(bState.ControlGroups[i].Units) == 0 {
				return true
			}

			log.Printf("SKRÓT: Przywoływanie grupy %d", i)
			clearSelection(bState)

			firstUnitInGroup := true
			for _, groupUnit := range bState.ControlGroups[i].Units {
				currentUnit, ok := getUnitByID(groupUnit.UnitID, bState)
				if ok && currentUnit.Exists && currentUnit.Owner == bState.PlayerID {
					currentUnit.IsSelected = true
					if firstUnitInGroup {
						bState.CurrentSelection = selectionState{
							OwnerID:    currentUnit.Owner,
							IsUnit:     true,
							UnitID:     currentUnit.ID,
							BuildingID: 0,
						}
						firstUnitInGroup = false
					}
				}
			}
			if firstUnitInGroup {
				clearSelection(bState)
			}
			bState.MouseState = mouseStateNormal
			return true
		}
	}

	if bState.CurrentSelection.IsUnit && bState.CurrentSelection.OwnerID == bState.PlayerID {
		selectedUnit, ok := getUnitByID(bState.CurrentSelection.UnitID, bState)
		if !ok || !selectedUnit.Exists {
			clearSelection(bState)
			return false
		}

		if rl.IsKeyPressed(rl.KeyS) {
			log.Println("SKRÓT: Komenda STOP dla jednostki")
			selectedUnit.addUnitCommand(cmdUStop, selectedUnit.X, selectedUnit.Y, 0, bState)
			return true
		}
		if rl.IsKeyPressed(rl.KeyC) {
			if (selectedUnit.Type == unitPriestess || selectedUnit.Type == unitPriest ||
				selectedUnit.Type == unitMage) && selectedUnit.Mana >= spellBufferMagicShower {
				log.Println("SKRÓT: Wejście w tryb rzucania czaru bojowego")
				bState.MouseState = mouseStateCasting
				return true
			}
		}
	}

	return false
}

// handleCameraScroll obsługuje przewijanie kamery.
func handleCameraScroll(iState inputState, bState *battleState, pState *programState) bool {
	scrollSpeed := 200.0 * rl.GetFrameTime()
	moved := false

	// --- ZMIANA: Obsługa przybliżania skokowego (Integer/Step Scaling) ---
	wheelMove := rl.GetMouseWheelMove()
	if wheelMove != 0 {
		// Krok 0.5 jest bezpieczny dla kafelków 16x14 (bo 16*1.5=24, 14*1.5=21 - liczby całkowite)
		const zoomStep = 0.5

		currentZoom := bState.GameCamera.Zoom
		newZoom := currentZoom

		if wheelMove > 0 {
			// Przybliżanie: Znajdź następny "poziom" w górę
			// math.Floor zapewnia, że jak masz np. 1.1, to "zresetuje" do bazy i skoczy do 1.5
			newZoom = float32(math.Floor(float64(currentZoom/zoomStep))+1) * zoomStep
		} else {
			// Oddalanie: Znajdź poprzedni "poziom" w dół
			newZoom = float32(math.Ceil(float64(currentZoom/zoomStep))-1) * zoomStep
		}

		// Bezpieczniki (Clamp)
		if newZoom < 0.5 {
			newZoom = 0.5
		}
		if newZoom > 4.0 { // Możesz dać 3.0 lub 4.0, zależy jak bardzo chcesz zbliżyć
			newZoom = 4.0
		}

		// Przypisujemy tylko jeśli wartość faktycznie się zmieniła
		if newZoom != currentZoom {
			bState.GameCamera.Zoom = newZoom
			log.Printf("Przybliżenie kamery (skokowe): %.1f", bState.GameCamera.Zoom)
			moved = true
		}
	}

	// --- 1. Obsługa Klawiatury ---
	if rl.IsKeyDown(rl.KeyRight) {
		bState.GameCamera.Target.X += scrollSpeed
		moved = true
	}
	// ... reszta funkcji bez zmian ...
	if rl.IsKeyDown(rl.KeyLeft) {
		bState.GameCamera.Target.X -= scrollSpeed
		moved = true
	}
	if rl.IsKeyDown(rl.KeyUp) {
		bState.GameCamera.Target.Y -= scrollSpeed
		moved = true
	}
	if rl.IsKeyDown(rl.KeyDown) {
		bState.GameCamera.Target.Y += scrollSpeed
		moved = true
	}

	// --- 2. Obsługa Myszki (Krawędzie Ekranu) ---
	if iState.MousePosition.X < scrollZoneXThreshold && iState.MousePosition.X >= 0 {
		bState.GameCamera.Target.X -= scrollSpeed
		moved = true
	}
	if iState.MousePosition.X > pState.VirtualWidth-scrollZoneXThreshold {
		bState.GameCamera.Target.X += scrollSpeed
		moved = true
	}
	if iState.MousePosition.Y < scrollZoneYThreshold {
		bState.GameCamera.Target.Y -= scrollSpeed
		moved = true
	}
	if iState.MousePosition.Y > pState.VirtualHeight-scrollZoneYThreshold {
		bState.GameCamera.Target.Y += scrollSpeed
		moved = true
	}

	// --- 3. Ograniczenie Kamery (Clamping) ---
	fullMapWidth := float32(uint16(boardMaxX) * uint16(tileWidth))
	fullMapHeight := float32(uint16(boardMaxY) * uint16(tileHeight))

	clampCameraTarget(&bState.GameCamera, fullMapWidth, fullMapHeight, pState.GameViewWidth, pState.VirtualHeight)

	return moved
}

func handleBoardInitialChecks(iState inputState, bState *battleState, pState *programState) (bool, uint8, uint8) {
	// Sprawdzamy dynamiczną szerokość widoku
	if iState.MousePosition.X < 0 || iState.MousePosition.X >= pState.GameViewWidth {
		return true, 0, 0
	}

	worldPos := rl.GetScreenToWorld2D(iState.MousePosition, bState.GameCamera)
	tileX := uint8(worldPos.X / float32(tileWidth))
	tileY := uint8(worldPos.Y / float32(tileHeight))

	if tileX >= boardMaxX || tileY >= boardMaxY {
		if iState.IsLeftMouseButtonPressed && !bState.DragContext.IsActive {
			clearSelection(bState)
		}

		return true, tileX, tileY
	}

	return false, tileX, tileY
}

// Odpowiada za ustawienie odpowiedniego rozkazu dla jednostki jeśli użyto prawego przycisku myszy.
func handleBoardRightClick(iState inputState, bState *battleState, tileX, tileY uint8) (clickHandled bool) {
	if !iState.IsRightMouseButtonPressed {
		return false
	}

	if bState.MouseState > mouseStateNormal {
		log.Println("INPUT: Anulowano tryb celowania prawym przyciskiem.")

		bState.MouseState = mouseStateNormal
		bState.PendingBuildingType = 0
		bState.DragContext.IsActive = false
		bState.CurrentMessage.Text = "Anulowano"
		bState.CurrentMessage.Duration = 30

		return true
	}

	selectedUnits := getSelectedUnits(bState)

	if len(selectedUnits) == 0 {
		return false
	}

	targetTile := &bState.Board.Tiles[tileX][tileY]
	targetID, targetOwner := targetTile.getTargetFromTile()

	cmdType, isCommandValid := resolveRightClickCommandType(targetTile, targetID, targetOwner, selectedUnits, bState, iState)

	if !isCommandValid {
		return true
	}

	sendUnitCommand(bState, selectedUnits, cmdType, tileX, tileY, targetID, iState.IsCtrlKeyDown)

	return true
}

// Odpowiada za dopasowanie rozkazu dla jednostki do sytuacji na planszy.
func resolveRightClickCommandType(
	targetTile *tile, targetID uint, targetOwner uint8,
	selectedUnits []*unit, bState *battleState, iState inputState,
) (cmdType commandType, isCommandValid bool) {
	cmdType = cmdUMove
	isCommandValid = true

	switch {
	// 0. Chodzenie po zniszczonej/nie wybudowanej palisadzie
	case targetTile.Building != nil && targetTile.Building.Type == buildingPalisade &&
		targetTile.Building.IsUnderConstruction:
		cmdType = cmdUMove

	// 1. Atak na wrogie jednostki/budynki
	case targetID != 0 && (targetOwner != bState.PlayerID || iState.IsCtrlKeyDown):
		cmdType = cmdUAttack
	// 2. Ścinanie drzewa
	case isTreeStump(targetTile.TextureID):
		canAttackTree := false

		for _, u := range selectedUnits {
			if u.canDamageTree(targetTile.X, targetTile.Y, bState) {
				canAttackTree = true

				break
			}
		}

		if canAttackTree {
			cmdType = cmdUAttack
			// targetID pozostaje 0; koordynaty ataku są przekazywane przez tileX, tileY
		} else {
			bState.CurrentMessage.Text = "Zaznaczone jednostki nie mogą atakować drzew!"
			bState.CurrentMessage.Duration = 60

			isCommandValid = false
		}
	// 3. Nasz drwal najeżdża na naszą budowę
	// a) w zaznaczonych jednostkach mamy drwala
	// b) najechaliśmy na naszą budowę
	// Najprostszy przypadek: jedna jednostka zaznaczona i jest to drwal
	case selectedUnits[0].Type == unitAxeman && targetTile.Building != nil &&
		targetTile.Building.IsUnderConstruction && targetOwner == bState.PlayerID:
		// drwal dostaje rozkaz budowy
		// @reminder: pamiętaj o rysowaniu kursora
		// cmdType = cmdBuild
		// Nie wiem czemu cmdBuild nie działa, ale naprawa już tak
		cmdType = cmdUBuild
	case selectedUnits[0].Type == unitAxeman && targetTile.Building != nil &&
		!targetTile.Building.IsUnderConstruction && targetTile.Building.HP < targetTile.Building.MaxHP &&
		targetOwner == bState.PlayerID:
		cmdType = cmdURepair

	// 4. Nieudane wydanie rozaku pójścia w jakieś miejsce
	case !targetTile.IsWalkable:
		bState.CurrentMessage.Text = "Nieprzechodnie!"
		bState.CurrentMessage.Duration = 60

		isCommandValid = false
	}

	return cmdType, isCommandValid
}

const dragThresholdPixels float32 = 3.0

func handleMouseStatePlacingConstruction(tileX, tileY uint8, bState *battleState) {
	log.Printf("DBG_LCLICK: Tryb budowy. Typ z pamięci: %d", bState.PendingBuildingType)

	tryBuildStructure(bState, tileX, tileY)

	switch bState.PendingBuildingType {
	case buildingRoad, buildingPalisade, buildingBridge:
		// nic, bo chcemy móc dalej budować.
	default:
		// Postawiliśmy „prawdziwy budynek” i kończymy, bo pewnie nie ma więcej mleka.
		bState.MouseState = mouseStateNormal
		bState.PendingBuildingType = 0
	}
}

func handleMouseStateWorking(tileX, tileY uint8, bState *battleState, iState inputState) {
	// === NAPRAWA BĄDŹ NOWA BUDOWA ===
	// Tutaj jest pierwsza okazja, aby dowiedzieć się, czy naprawiamy, czy budujemy
	// wcześniej nie znaliśmy celu. Dlatego nie dało się wybrać. Od tej chwili rozdzielamy.
	log.Println("DBG_LCLICK: Tryb naprawy. Typ z pamięci")
	// 1. Ponieważ nacisnęliśmy lewy przycisk myszy, to bierzemy współrzędne z planszy
	currentTile := &bState.Board.Tiles[tileX][tileY]
	targetBld := currentTile.Building

	// Idziemy dalej tylko jeżeli jako cel obraliśmy budynek
	if targetBld == nil {
		bState.CurrentMessage.Text = "Wskaż budynek!"
		bState.CurrentMessage.Duration = 40
		// Nie powinno być drugiej szansy, ale musimy wyjść z tej funkcji
		// O ile dobrze pamiętam, to musi być true, inaczej jest przeciąganie
		return
	}

	var cmd commandType

	if targetBld.IsUnderConstruction {
		// plac budowy, należy ją ukończyć.
		cmd = cmdUBuild
	} else if targetBld.isRepairable(bState.PlayerID) { // @reminder: wydaje mi się, że metoda zamiast gołego HP
		// uszkodzony budynek, należy go naprawić.
		cmd = cmdURepair
	} else {
		// nic z powyższych, nie ma co robić.
		bState.CurrentMessage.Text = "Ten budynek nie wymaga naprawy!"
		bState.CurrentMessage.Duration = 60
		bState.MouseState = mouseStateNormal

		return
	}

	// 2. Możemy naprawiać tylko palisady oraz swoje budynki, które są uszkodzone
	// @reminder: w tej chwili to metoda bld.IsRepairable spełnia tę samą funkcję!
	// @todo: sprawdź, czy @reminder ma rację i można uprościć
	canBeRepaired := ((targetBld.Owner == bState.PlayerID) || (targetBld.Type == buildingPalisade) ||
		targetBld.Type == buildingBridge) && targetBld.HP < targetBld.MaxHP

	if !canBeRepaired {
		bState.CurrentMessage.Text = "Nie możesz przacować przy wrogich budynkach!"
		bState.CurrentMessage.Duration = 60
		bState.MouseState = mouseStateNormal
		// Niech stoją bezczynnie
		return
	}

	// 3. Odsianie jednostek, które nie są UNIT_AXEMAN z całej zaznaczonej drużyny
	// @todo: sprawdź, czy to w ogóle działa, raczej tak, i wyodrębnij osobną funkcję
	// do tego. Przyda się to w kilku miejscach oraz być może przy zbiorowym rzucaniu czarów.
	selectedUnits := getSelectedUnits(bState)

	var repairCrew []*unit

	for _, u := range selectedUnits {
		if u.Type == unitAxeman {
			repairCrew = append(repairCrew, u)
		}
	}

	if len(repairCrew) == 0 {
		bState.CurrentMessage.Text = "Brak Toporników w zaznaczeniu!"
		bState.CurrentMessage.Duration = 60
		bState.MouseState = mouseStateNormal

		return
	}
	// 4. Rozkaz gotowy, wiadomo kto, co, można przekazać dalej
	sendUnitCommand(bState, repairCrew, cmd, tileX, tileY, targetBld.ID, iState.IsCtrlKeyDown)
	log.Printf("INPUT: Wysłano %d Toporników do naprawy budynku ID %d.", len(repairCrew), targetBld.ID)

	// Zmieniamy stan myszki i wracamy
	bState.MouseState = mouseStateNormal
}

func handleMouseStateCasting(tileX, tileY uint8, bState *battleState) {
	log.Println("DBG_LCLICK: Tryb rzucania czaru ofensywnego.")

	selectedUnit, ok := getUnitByID(bState.CurrentSelection.UnitID, bState)
	if !ok || !selectedUnit.Exists {
		bState.MouseState = mouseStateNormal

		return
	}

	spellActionType := cmdUCastSpell

	// Dodajemy rozkaz do kolejki
	selectedUnit.addUnitCommand(spellActionType, tileX, tileY, 0, bState)
	log.Printf("DBG_LCLICK: Wydano rozkaz czaru %d na (%d,%d).", spellActionType, tileX, tileY)

	bState.MouseState = mouseStateNormal

	return
}

// @todo: tymczasowe ogarnianie drużynowych rozkazów. Muszę wrócić i poprawić!
// @todo: jak poprawnie obsługiwać całe drużyny? Jak dobierać, które przyciski dozwolone?
// w W3 była możliwość „tab” pomiędzy rodzajami jednostek w drużynie i dostawania dostępu
// do przycisków „rodzajowych”. Chyba muszę podobnie zrobić, bo mam miejsce tylko na
// pięć przycisków: atak(0), stop(1), czar1(2), czar2(3),naprawa(4) jeżeli coś innego będzie
// dodane to mam problem. Dodatkowo jest problem mieszania kontekstu bojowego z gospodarczym.
func handleMouseStateNormalPressed(tileX, tileY uint8, bState *battleState, iState inputState, pState *programState) {
	log.Printf("DBG_LCLICK: Naciśnięto kafelek (%d,%d). Tryb myszy: Normal", tileX, tileY)

	tileUnderCursor := &bState.Board.Tiles[tileX][tileY]
	targetID, _ := tileUnderCursor.getTargetFromTile()

	if targetID != 0 {
		selectObjectByClick(tileX, tileY, bState)
	} else {
		if !iState.IsShiftKeyDown {
			clearSelection(bState)
		}

		bState.DragContext.IsActive = true
		bState.DragContext.AnchorPos = clampDragPosition(iState.MousePosition, pState)
		// ↓↓↓ bez tego poprzednie współrzędne zaśmiecają rysowanie ramki.
		bState.DragContext.CurrentPos = bState.DragContext.AnchorPos
	}
}

// Obsługuje chwilę zwolnienia przycisku. Zatwierdzamy ramkę zaznaczania jeśli była stworzona.
func handleMouseStateNormalReleased(bState *battleState) {
	// 0. Gracz nie przeciągał, nie ma ramki do domknięcia.
	if !bState.DragContext.IsActive {
		return
	}

	// @reminder: trzeba wywalić bezpośrednie odwołanie do rl!
	distance := rl.Vector2Distance(bState.DragContext.AnchorPos, bState.DragContext.CurrentPos)

	if distance > dragThresholdPixels {
		log.Println("DBG_LCLICK: Kończę ramkę zaznaczania.")
		performBoxSelection(bState, bState.DragContext.AnchorPos, bState.DragContext.CurrentPos)
	} else {
		log.Println("DGB_LCLICK: Zbyt mały ruch, nie jest to zaznaczanie ramką.")
	}

	// Zwolniony przycisk myszy, obsłużyliśmy udane i nieudane zaznaczanie, kończymy.
	bState.DragContext.IsActive = false
}

// @todo: ogarnij czemu to nie działa jako przekazanie STOP do wszystkich
// zaznaczonych jednostek!
func sendUnitCommand(bState *battleState, units []*unit, command commandType, x, y uint8, targetID uint, ctrlDown bool) {
	log.Printf("INFO: input.go wysłano rozkaz.")

	for _, u := range units {
		u.AllowFriendlyFire = ctrlDown
	}

	if len(units) > 1 {
		bState.assignGroupCommand(command, x, y, targetID, units)
	} else {
		units[0].addUnitCommand(command, x, y, targetID, bState)
	}

	bState.MouseState = mouseStateNormal
}

func handleBoardInteraction(iState inputState, bState *battleState, pState *programState) {
	handledInitial, tileX, tileY := handleBoardInitialChecks(iState, bState, pState)

	if handledInitial {
		return
	}

	if handleBoardRightClick(iState, bState, tileX, tileY) {
		return
	}

	if iState.IsLeftMouseButtonPressed {
		switch bState.MouseState {
		case mouseStatePlaceConstruction:
			handleMouseStatePlacingConstruction(tileX, tileY, bState)
		case mouseStateWorking:
			handleMouseStateWorking(tileX, tileY, bState, iState)
		case mouseStateCasting:
			handleMouseStateCasting(tileX, tileY, bState)
		case mouseStateNormal:
			handleMouseStateNormalPressed(tileX, tileY, bState, iState, pState) // Dodano ps
		default:
			log.Printf("BŁĄD KRYTYCZNY: Nieobsługiwany stan myszy: %d", bState.MouseState)
		}
		return
	}
}

// clampDragPosition ogranicza pozycję ramki do obszaru planszy
func clampDragPosition(position rl.Vector2, pState *programState) rl.Vector2 {
	clamped := position

	// Ograniczenie do lewej/górnej krawędzi
	if clamped.X < 0 {
		clamped.X = 0
	}
	if clamped.Y < 0 {
		clamped.Y = 0
	}

	// Ograniczenie do prawej krawędzi (granica z UI)
	if clamped.X > pState.GameViewWidth {
		clamped.X = pState.GameViewWidth
	}

	// Ograniczenie do dolnej krawędzi
	if clamped.Y > pState.VirtualHeight {
		clamped.Y = pState.VirtualHeight
	}

	return clamped
}

func handleGameInput(bState *battleState, pState *programState, iState inputState) {
	virtualMouse := iState.MousePosition

	handleCameraScroll(iState, bState, pState)

	// Obsługa przeciągania ramki zaznaczania
	if bState.DragContext.IsActive {
		// Jeśli chcemy anulować to PPMem.
		if iState.IsRightMouseButtonPressed {
			bState.DragContext.IsActive = false

			return
		}
		// Zawsze aktualizuj pozycję końca ramki
		bState.DragContext.CurrentPos = clampDragPosition(virtualMouse, pState)

		// Obsługa zwolnienia przycisku myszy
		if iState.IsLeftMouseButtonReleased {
			handleMouseStateNormalReleased(bState)

			return
		}

		// Blokuj inne zdarzenia podczas przeciągania
		return
	}

	// Obsługa UI
	if isMouseOverUI(pState, virtualMouse) {
		if iState.IsLeftMouseButtonPressed || iState.IsRightMouseButtonPressed ||
			iState.IsLeftMouseButtonDown || iState.IsLeftMouseButtonReleased {
			if handleGameUIClicks(iState, bState, pState) {
				return
			}
		}
	} else {
		// Obsługa planszy
		handleGameShortcuts(bState)
		handleBoardInteraction(iState, bState, pState)
	}
}

// OBSŁUGA ZAZNACZANIA

func clearSelection(bState *battleState) {
	log.Println("SELEKCJA: Rozpoczynam clearSelection.")

	for _, currentUnit := range bState.Units {
		if currentUnit.Exists && currentUnit.IsSelected {
			currentUnit.IsSelected = false
		}
	}

	if bState.CurrentSelection.IsUnit || bState.CurrentSelection.BuildingID != 0 {
		bState.CurrentSelection = selectionState{}
	}
}

func selectObjectByClick(tileX, tileY uint8, bState *battleState) {
	currentTile := &bState.Board.Tiles[tileX][tileY]
	currentUnit := currentTile.Unit
	bld := currentTile.Building

	if currentUnit == nil && bld == nil {
		found := false
		originalTileX, originalTileY := tileX, tileY
		log.Printf("DBG_SELECTOBJECT: Na (%d,%d) nie ma bezpośredniego obiektu. Szukam w sąsiedztwie...", tileX, tileY)

		for j := originalTileY - 1; j <= originalTileY+1; j++ {
			for i := originalTileX - 1; i <= originalTileX+1; i++ {
				if i < boardMaxX && j < boardMaxY {
					nt := &bState.Board.Tiles[i][j]
					if nt.Unit != nil || nt.Building != nil {
						currentUnit = nt.Unit
						bld = nt.Building
						found = true

						break
					}
				}
			}

			if found {
				break
			}
		}

		if !found {
			clearSelection(bState)
			bState.MouseState = mouseStateNormal

			return
		}
	}

	isShiftDown := rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)

	if currentUnit != nil && currentUnit.Exists {
		log.Printf("DBG_SELECTOBJECT: Znaleziono jednostkę ID %d.", currentUnit.ID)

		if currentUnit.Owner != bState.PlayerID {
			clearSelection(bState)
			bState.CurrentSelection = selectionState{
				OwnerID:    currentUnit.Owner,
				IsUnit:     true,
				UnitID:     currentUnit.ID,
				BuildingID: 0,
			}
			bState.CurrentMessage.Text = fmt.Sprintf("Wroga jednostka: %v", currentUnit.Type)
			bState.CurrentMessage.Duration = 20
			bState.MouseState = mouseStateNormal

			return
		}

		if isShiftDown {
			currentUnit.IsSelected = !currentUnit.IsSelected
			if !currentUnit.IsSelected && bState.CurrentSelection.UnitID == currentUnit.ID {
				foundNewPrimary := false

				for _, u := range bState.Units {
					if u.Exists && u.IsSelected && u.Owner == bState.PlayerID {
						bState.CurrentSelection = selectionState{OwnerID: u.Owner, IsUnit: true, UnitID: u.ID}
						foundNewPrimary = true

						break
					}
				}

				if !foundNewPrimary {
					bState.CurrentSelection = selectionState{}
				}
			} else if currentUnit.IsSelected && !bState.CurrentSelection.IsUnit {
				bState.CurrentSelection = selectionState{OwnerID: currentUnit.Owner, IsUnit: true, UnitID: currentUnit.ID}
			}

		} else {
			clearSelection(bState)
			currentUnit.IsSelected = true
			bState.CurrentSelection = selectionState{
				OwnerID:    currentUnit.Owner,
				IsUnit:     true,
				UnitID:     currentUnit.ID,
				BuildingID: 0,
			}
		}

		if currentUnit.IsSelected && currentUnit.Owner == bState.PlayerID {
			switch currentUnit.Type {
			case unitCow:
				bState.CurrentMessage.Text = "Muuu ?"
			case unitAxeman:
				bState.CurrentMessage.Text = "Tak ?"
			default:
				bState.CurrentMessage.Text = "Rozkaz?"
			}
			bState.CurrentMessage.Duration = 20
		}
		bState.MouseState = mouseStateNormal

	} else if bld != nil && bld.Exists {
		log.Printf("DBG_SELECTOBJECT: Znaleziono budynek ID %d.", bld.ID)
		if !isShiftDown {
			clearSelection(bState)
		}

		bState.CurrentSelection = selectionState{
			OwnerID:    bld.Owner,
			IsUnit:     false,
			UnitID:     0,
			BuildingID: bld.ID,
		}
		bState.MouseState = mouseStateNormal

		if bld.Owner == bState.PlayerID {
			bState.CurrentMessage.Text = fmt.Sprintf("Moja budowla: %v", bld.Type)
		} else {
			bState.CurrentMessage.Text = fmt.Sprintf("Wroga budowla: %v", bld.Type)
		}
		bState.CurrentMessage.Duration = 20
	} else {
		clearSelection(bState)
		bState.MouseState = mouseStateNormal
	}
}

func performBoxSelection(bState *battleState, startPosition, endPosistion rl.Vector2) {
	worldStart := rl.GetScreenToWorld2D(startPosition, bState.GameCamera)
	worldEnd := rl.GetScreenToWorld2D(endPosistion, bState.GameCamera)

	minX := uint8(min(worldStart.X, worldEnd.X) / float32(tileWidth))
	maxX := uint8(max(worldStart.X, worldEnd.X) / float32(tileWidth))
	minY := uint8(min(worldStart.Y, worldEnd.Y) / float32(tileHeight))
	maxY := uint8(max(worldStart.Y, worldEnd.Y) / float32(tileHeight))

	isShiftDown := rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)
	if !isShiftDown {
		clearSelection(bState)
	}

	var selectedCount int

	var firstSelectedUnit *unit

	for _, currentUnit := range bState.Units {
		if currentUnit.Exists && currentUnit.Owner == bState.PlayerID {
			if currentUnit.X >= minX && currentUnit.X <= maxX && currentUnit.Y >= minY && currentUnit.Y <= maxY {
				if isShiftDown {
					currentUnit.IsSelected = !currentUnit.IsSelected
				} else {
					currentUnit.IsSelected = true
				}
			} else {
				if !isShiftDown && currentUnit.IsSelected {
					currentUnit.IsSelected = false
				}
			}

			if currentUnit.IsSelected {
				selectedCount++

				if firstSelectedUnit == nil {
					firstSelectedUnit = currentUnit
				}
			}
		}
	}

	if selectedCount > 0 {
		bState.CurrentSelection = selectionState{
			OwnerID:    firstSelectedUnit.Owner,
			IsUnit:     true,
			UnitID:     firstSelectedUnit.ID,
			BuildingID: 0,
		}
	} else {
		clearSelection(bState)
	}
}

func handleMinimapInteraction(iState inputState, bState *battleState, pState *programState) bool {
	// Prostokąt minimapy obliczany dynamicznie!
	// ps.GameViewWidth to początek panelu UI.
	minimapRect := rl.NewRectangle(
		pState.GameViewWidth+minimapOffsetX,
		float32(0)+minimapOffsetY,
		minimapDisplayWidth,
		minimapDisplayHeight,
	)

	isMouseOverMinimap := rl.CheckCollisionPointRec(
		iState.MousePosition,
		minimapRect,
	)

	if iState.IsLeftMouseButtonReleased && bState.MapInitialClickPos.X != 0.0 {
		bState.IsMapDragging = false
		bState.MapInitialClickPos = rl.NewVector2(0.0, 0.0)
		bState.CameraTargetOnDragStart = rl.NewVector2(0.0, 0.0)

		return true
	}

	if !isMouseOverMinimap && !bState.IsMapDragging {
		return false
	}

	// Przekazujemy ps do clamping
	if handleMinimapLeftMouse(iState, bState, minimapRect, pState) {
		return true
	}

	if isMouseOverMinimap {
		if handleMinimapRightMouse(iState, bState, minimapRect) {
			return true
		}
	}

	return false
}

func handleMinimapLeftMouse(iState inputState, bState *battleState, minimapRect rl.Rectangle, pState *programState) bool {
	fullMapPixelWidth := float32(uint16(boardMaxX) * uint16(tileWidth))
	fullMapPixelHeight := float32(uint16(boardMaxY) * uint16(tileHeight))
	scaleX := fullMapPixelWidth / minimapDisplayWidth
	scaleY := fullMapPixelHeight / minimapDisplayHeight

	if iState.IsLeftMouseButtonPressed {
		clickedX := (iState.MousePosition.X - minimapRect.X) * scaleX
		clickedY := (iState.MousePosition.Y - minimapRect.Y) * scaleY
		bState.GameCamera.Target = rl.NewVector2(clickedX, clickedY)

		// Clamping używa dynamicznych wymiarów
		clampCameraTarget(&bState.GameCamera, fullMapPixelWidth, fullMapPixelHeight,
			pState.GameViewWidth, pState.VirtualHeight)

		bState.MapInitialClickPos = iState.MousePosition
		bState.CameraTargetOnDragStart = bState.GameCamera.Target
		bState.IsMapDragging = false
		return true
	}

	if iState.IsLeftMouseButtonDown && bState.MapInitialClickPos.X != 0.0 {
		if !bState.IsMapDragging &&
			rl.Vector2Distance(bState.MapInitialClickPos, iState.MousePosition) > minimapClickDragThreshold {
			bState.IsMapDragging = true
		}

		if bState.IsMapDragging {
			deltaX := (iState.MousePosition.X - bState.MapInitialClickPos.X) * scaleX
			deltaY := (iState.MousePosition.Y - bState.MapInitialClickPos.Y) * scaleY
			bState.GameCamera.Target.X = bState.CameraTargetOnDragStart.X + deltaX
			bState.GameCamera.Target.Y = bState.CameraTargetOnDragStart.Y + deltaY

			// Clamping używa dynamicznych wymiarów
			clampCameraTarget(&bState.GameCamera, fullMapPixelWidth, fullMapPixelHeight,
				pState.GameViewWidth, pState.VirtualHeight)

			return true
		}
	}

	return false
}

func handleMinimapRightMouse(
	iState inputState,
	bState *battleState,
	minimapRect rl.Rectangle,
) bool {
	if !iState.IsRightMouseButtonPressed {
		return false
	}

	selectedUnits := getSelectedUnits(bState)
	if len(selectedUnits) == 0 {
		return true
	}

	scaleX := float32(uint16(boardMaxX)*uint16(tileWidth)) / minimapDisplayWidth
	scaleY := float32(uint16(boardMaxY)*uint16(tileHeight)) / minimapDisplayHeight

	worldX := float64(iState.MousePosition.X-minimapRect.X) * float64(scaleX)
	worldY := float64(iState.MousePosition.Y-minimapRect.Y) * float64(scaleY)
	tileX := uint8(math.Min(math.Max(worldX/float64(tileWidth), 0), float64(boardMaxX-1)))
	tileY := uint8(math.Min(math.Max(worldY/float64(tileHeight), 0), float64(boardMaxY-1)))

	currentTile := &bState.Board.Tiles[tileX][tileY]
	targetID, targetOwner := currentTile.getTargetFromTile()

	cmd := cmdUMove

	if targetID != 0 && (targetOwner != bState.PlayerID || iState.IsCtrlKeyDown) {
		cmd = cmdUAttack
	} else if !isWalkable(bState, tileX, tileY) {
		bState.CurrentMessage.Text = "Nieprzechodnie!"
		bState.CurrentMessage.Duration = 60

		return true
	}

	sendUnitCommand(bState, selectedUnits, cmd, tileX, tileY, targetID, iState.IsCtrlKeyDown)
	bState.MouseState = mouseStateNormal

	return true
}

func clampCameraTarget(camera *rl.Camera2D, mapPixelWidth, mapPixelHeight, viewPixelWidth, viewPixelHeight float32) {
	halfViewWidth := (viewPixelWidth / 2.0) / camera.Zoom
	halfViewHeight := (viewPixelHeight / 2.0) / camera.Zoom

	minTargetX := halfViewWidth
	maxTargetX := mapPixelWidth - halfViewWidth
	minTargetY := halfViewHeight
	maxTargetY := mapPixelHeight - halfViewHeight

	if mapPixelWidth < viewPixelWidth/camera.Zoom {
		camera.Target.X = mapPixelWidth / 2.0
	} else {
		camera.Target.X = float32(math.Max(float64(camera.Target.X), float64(minTargetX)))
		camera.Target.X = float32(math.Min(float64(camera.Target.X), float64(maxTargetX)))
	}

	if mapPixelHeight < viewPixelHeight/camera.Zoom {
		camera.Target.Y = mapPixelHeight / 2.0
	} else {
		camera.Target.Y = float32(math.Max(float64(camera.Target.Y), float64(minTargetY)))
		camera.Target.Y = float32(math.Min(float64(camera.Target.Y), float64(maxTargetY)))
	}
}
