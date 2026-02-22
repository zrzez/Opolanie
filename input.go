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

func handleGameUIClicks(input inputState, bs *battleState, ps *programState) bool {
	virtualMouse := input.MousePosition

	// 1. Najpierw sprawdzamy minimapę (to kluczowe, by nie klikać "przez" nią)
	if handleMinimapInteraction(input, bs, ps) {
		return true
	}

	// 2. Obsługa kliknięć w przyciski akcji (panel boczny)
	if input.IsLeftMouseButtonPressed {
		// Sprawdzamy każdy z 5 przycisków
		for btnIndex := range uiActionMaxButtons {
			action := bs.UI.CurrentActions[btnIndex]

			// Pomijamy wyłączone przyciski
			if !action.IsActive {
				continue
			}

			rect := bs.UI.ActionButtons[btnIndex]
			if rl.CheckCollisionPointRec(virtualMouse, rect) {
				log.Printf("UI: Wybrano akcję przycisku %d: %s (Typ: %d)", btnIndex, action.Label, action.Cmd.ActionType)

				switch action.Cmd.ActionType {
				case cmdBuildStructure:
					// To jest tryb myszy. Nie wysyłamy rozkazu, lecz zmieniamy stan kursora.
					bs.MouseCommandMode = cmdBuildStructure
					// Zapisujemy rodzaj budynku do "plecaka" w battleState
					bs.PendingBuildingType = buildingType(action.Cmd.InteractionTargetID)

					bs.CurrentMessage.Text = "Wskaż miejsce"
					bs.CurrentMessage.Duration = 60

					// Opcjonalnie: czyścimy zaznaczenie jednostek, by nie przeszkadzały
					clearSelection(bs)
				case cmdRepairStructure:
					bs.MouseCommandMode = cmdRepairStructure
					bs.CurrentMessage.Text = "Wskaż budynek do naprawy"
					bs.CurrentMessage.Duration = 60

				default:
					// Rozkazy natychmiastowe (CMD_PRODUCE, CMD_MILK, itp.)
					// Wysyłamy je bezpośrednio do kolejki rozkazów.
					bs.CurrentCommands[0] = action.Cmd
				}

				return true
			}
		}

		// Logowanie kliknięcia w tło UI (pomocne przy debugowaniu wymiarów)
		relativeX := virtualMouse.X - ps.GameViewWidth
		log.Printf("Klik w UI (tło) na: %.1f, %.1f", relativeX, virtualMouse.Y)

		return true
	}

	return false
}

func handleGameShortcuts(bs *battleState) bool {
	if rl.IsKeyPressed(rl.KeyKpSubtract) || rl.IsKeyPressed(rl.KeyMinus) {
		if bs.GameSpeed < 5 {
			bs.GameSpeed++
			log.Printf("SKRÓT: Prędkość gry zmniejszona do %d", bs.GameSpeed)
			return true
		}
	}

	if rl.IsKeyPressed(rl.KeyKpAdd) || rl.IsKeyPressed(rl.KeyEqual) {
		if bs.GameSpeed > 0 {
			bs.GameSpeed--
			log.Printf("SKRÓT: Prędkość gry zwiększona do %d", bs.GameSpeed)
			return true
		}
	}

	isShiftDown := rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)

	for i := 0; i <= 9; i++ {
		key := rl.KeyZero + int32(i)

		if rl.IsKeyPressed(key) {
			if isShiftDown {
				if bs.CurrentSelection.OwnerID != bs.PlayerID {
					return true
				}

				log.Printf("SKRÓT: Zapamiętywanie grupy %d", i)
				var selectedUnitIDs []uint
				for _, unit := range bs.Units {
					if unit.Exists && unit.Owner == bs.PlayerID && unit.IsSelected {
						selectedUnitIDs = append(selectedUnitIDs, unit.ID)
					}
				}
				bs.ControlGroups[i] = controlGroup{}
				for _, unitID := range selectedUnitIDs {
					bs.ControlGroups[i].Units = append(bs.ControlGroups[i].Units, controlGroupUnit{UnitID: unitID})
				}
				log.Printf("SKRÓT: Grupa %d utworzona z %d jednostek.", i, len(bs.ControlGroups[i].Units))
				return true

			}

			if len(bs.ControlGroups[i].Units) == 0 {
				return true
			}

			log.Printf("SKRÓT: Przywoływanie grupy %d", i)
			clearSelection(bs)

			firstUnitInGroup := true
			for _, groupUnit := range bs.ControlGroups[i].Units {
				unit, ok := getUnitByID(groupUnit.UnitID, bs)
				if ok && unit.Exists && unit.Owner == bs.PlayerID {
					unit.IsSelected = true
					if firstUnitInGroup {
						bs.CurrentSelection = selectionState{
							OwnerID:    unit.Owner,
							IsUnit:     true,
							UnitID:     unit.ID,
							BuildingID: 0,
						}
						firstUnitInGroup = false
					}
				}
			}
			if firstUnitInGroup {
				clearSelection(bs)
			}
			bs.MouseCommandMode = 1
			return true
		}
	}

	if bs.CurrentSelection.IsUnit && bs.CurrentSelection.OwnerID == bs.PlayerID {
		selectedUnit, ok := getUnitByID(bs.CurrentSelection.UnitID, bs)
		if !ok || !selectedUnit.Exists {
			clearSelection(bs)
			return false
		}

		if rl.IsKeyPressed(rl.KeyS) {
			log.Println("SKRÓT: Komenda STOP dla jednostki")
			selectedUnit.addUnitCommand(cmdStop, selectedUnit.X, selectedUnit.Y, 0, bs)
			return true
		}
		if rl.IsKeyPressed(rl.KeyC) {
			if (selectedUnit.Type == unitPriestess || selectedUnit.Type == unitPriest ||
				selectedUnit.Type == unitMage) && selectedUnit.Mana >= 79 {
				log.Println("SKRÓT: Wejście w tryb rzucania czaru bojowego")
				bs.MouseCommandMode = cmdCastSpell
				return true
			}
		}
	}

	return false
}

// handleCameraScroll obsługuje przewijanie kamery.
// handleCameraScroll obsługuje przewijanie kamery.
func handleCameraScroll(input inputState, bs *battleState, ps *programState) bool {
	scrollSpeed := 200.0 * rl.GetFrameTime()
	moved := false

	// --- ZMIANA: Obsługa przybliżania skokowego (Integer/Step Scaling) ---
	wheelMove := rl.GetMouseWheelMove()
	if wheelMove != 0 {
		// Krok 0.5 jest bezpieczny dla kafelków 16x14 (bo 16*1.5=24, 14*1.5=21 - liczby całkowite)
		const zoomStep = 0.5

		currentZoom := bs.GameCamera.Zoom
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
			bs.GameCamera.Zoom = newZoom
			log.Printf("Przybliżenie kamery (skokowe): %.1f", bs.GameCamera.Zoom)
			moved = true
		}
	}

	// --- 1. Obsługa Klawiatury ---
	if rl.IsKeyDown(rl.KeyRight) {
		bs.GameCamera.Target.X += scrollSpeed
		moved = true
	}
	// ... reszta funkcji bez zmian ...
	if rl.IsKeyDown(rl.KeyLeft) {
		bs.GameCamera.Target.X -= scrollSpeed
		moved = true
	}
	if rl.IsKeyDown(rl.KeyUp) {
		bs.GameCamera.Target.Y -= scrollSpeed
		moved = true
	}
	if rl.IsKeyDown(rl.KeyDown) {
		bs.GameCamera.Target.Y += scrollSpeed
		moved = true
	}

	// --- 2. Obsługa Myszki (Krawędzie Ekranu) ---
	// ... reszta kodu bez zmian ...

	// (Poniżej wklejam resztę dla kontekstu, żebyś wiedział gdzie jesteśmy)
	if input.MousePosition.X < scrollZoneXThreshold && input.MousePosition.X >= 0 {
		bs.GameCamera.Target.X -= scrollSpeed
		moved = true
	}
	if input.MousePosition.X > ps.VirtualWidth-scrollZoneXThreshold {
		bs.GameCamera.Target.X += scrollSpeed
		moved = true
	}
	if input.MousePosition.Y < scrollZoneYThreshold {
		bs.GameCamera.Target.Y -= scrollSpeed
		moved = true
	}
	if input.MousePosition.Y > ps.VirtualHeight-scrollZoneYThreshold {
		bs.GameCamera.Target.Y += scrollSpeed
		moved = true
	}

	// --- 3. Ograniczenie Kamery (Clamping) ---
	fullMapWidth := float32(uint16(boardMaxX) * uint16(tileWidth))
	fullMapHeight := float32(uint16(boardMaxY) * uint16(tileHeight))

	clampCameraTarget(&bs.GameCamera, fullMapWidth, fullMapHeight, ps.GameViewWidth, ps.VirtualHeight)

	return moved
}

func handleBoardInitialChecks(input inputState, bs *battleState, ps *programState) (bool, uint8, uint8) {
	// Sprawdzamy dynamiczną szerokość widoku
	if input.MousePosition.X < 0 || input.MousePosition.X >= ps.GameViewWidth {
		return true, 0, 0
	}

	worldPos := rl.GetScreenToWorld2D(input.MousePosition, bs.GameCamera)
	tileX := uint8(worldPos.X / float32(tileWidth))
	tileY := uint8(worldPos.Y / float32(tileHeight))

	if tileX < 0 || tileX >= boardMaxX || tileY < 0 || tileY >= boardMaxY {
		if input.IsLeftMouseButtonPressed && !bs.IsSelectingBox {
			clearSelection(bs)
		}
		return true, tileX, tileY
	}
	return false, tileX, tileY
}

func handleBoardRightClick(input inputState, bs *battleState, tileX, tileY uint8) bool {
	if input.IsRightMouseButtonPressed {
		if bs.IsSelectingBox {
			bs.IsSelectingBox = false
			bs.SelectionStart = rl.NewVector2(0, 0)
			bs.InitialClickPos = rl.NewVector2(0, 0)
			bs.MouseCommandMode = 1
			return true
		}

		// Jeżeli jesteśmy w trybie celowania (budowa, naprawa, czar),
		// to prawy przycisk służy jako "ANULUJ", a nie jako "IDŹ TAM".
		if bs.MouseCommandMode > 1 {
			log.Println("INPUT: Anulowano tryb celowania prawym przyciskiem.")
			bs.MouseCommandMode = cmdIdle
			bs.PendingBuildingType = 0 // Na wypadek gdybyśmy anulowali budowanie
			bs.CurrentMessage.Text = "Anulowano"
			bs.CurrentMessage.Duration = 30
			return true // Zjadamy kliknięcie, nie wysyłamy ruchu!
		}

		selectedUnits := getSelectedUnits(bs)
		if len(selectedUnits) > 0 {
			tile := &bs.Board.Tiles[tileX][tileY]
			targetID := uint(0)
			var targetOwner uint8 = 0

			if tile.Unit != nil {
				targetID = tile.Unit.ID
				targetOwner = tile.Unit.Owner
			} else if tile.Building != nil {
				targetID = tile.Building.ID
				targetOwner = tile.Building.Owner
			}

			commandType := cmdMove

			if targetID != 0 && targetOwner != bs.PlayerID {
				commandType = cmdAttack
			} else {
				if !isWalkable(bs, tileX, tileY) {
					bs.CurrentMessage.Text = "Nieprzechodnie!"
					bs.CurrentMessage.Duration = 60
					return true
				}
			}

			sendUnitCommand(bs, selectedUnits, commandType, tileX, tileY, targetID, input.IsCtrlKeyDown)
			return true
		}

		if bs.MouseCommandMode != 1 {
			bs.MouseCommandMode = 1
			return true
		}
	}
	return false
}

const dragThresholdPixels float32 = 3.0

// @todo: tymczasowe ogarnianie drużynowych rozkazów. Muszę wrócić i poprawić!
// @todo: jak poprawnie obsługiwać całe drużyny? Jak dobierać, które przyciski dozwolone?
// w W3 była możliwość „tab” pomiędzy rodzajami jednostek w drużynie i dostawania dostępu
// do przycisków „rodzajowych”. Chyba muszę podobnie zrobić, bo mam miejsce tylko na
// pięć przycisków: atak(0), stop(1), czar1(2), czar2(3),naprawa(4) jeżeli coś innego będzie
// dodane to mam problem. Dodatkowo jest problem mieszania kontekstu bojowego z gospodarczym.
func handleBoardLeftClick(input inputState, bs *battleState, tileX, tileY uint8) bool {
	bs.InitialClickPos = input.MousePosition
	log.Printf("DBG_LCLICK: Kliknięto kafelek (%d, %d). Tryb myszy: %d", tileX, tileY, bs.MouseCommandMode)

	switch bs.MouseCommandMode {

	// === 1. TRYB BUDOWANIA ===
	case cmdBuildStructure:
		log.Printf("DBG_LCLICK: Tryb budowy. Typ z pamięci: %d", bs.PendingBuildingType)

		tryBuildStructure(bs, tileX, tileY)

		bs.MouseCommandMode = cmdIdle
		bs.PendingBuildingType = 0

		return true

	case cmdRepairStructure:
		// === NAPRAWA BĄDŹ NOWA BUDOWA ===
		// Tutaj jest pierwsza okazja, aby dowiedzieć się, czy naprawiamy, czy budujemy
		// wcześniej nie znaliśmy celu. Dlatego nie dało się wybrać. Od tej chwili rozdzielamy.
		log.Println("DBG_LCLICK: Tryb naprawy. Typ z pamięci")
		// 1. Ponieważ nacisnęliśmy lewy przycisk myszy, to bierzemy współrzędne z planszy
		tile := &bs.Board.Tiles[tileX][tileY]
		targetBld := tile.Building

		// Idziemy dalej tylko jeżeli jako cel obraliśmy budynek
		if targetBld == nil {
			bs.CurrentMessage.Text = "Wskaż budynek!"
			bs.CurrentMessage.Duration = 40
			// Nie powinno być drugiej szansy, ale musimy wyjść z tej funkcji
			// O ile dobrze pamiętam, to musi być true, inaczej jest przeciąganie
			return true
		}

		cmd := cmdRepairStructure

		// 2. Możemy naprawiać tylko palisady oraz swoje budynki, które są uszkodzone
		canRepair := ((targetBld.Owner == bs.PlayerID) || (targetBld.Type == buildingPalisade) || targetBld.Type == buildingBridge) &&
			targetBld.HP < targetBld.MaxHP

		if !canRepair {
			bs.CurrentMessage.Text = "Nie możesz naprawiać wrogich budynków!"
			bs.CurrentMessage.Duration = 60
			bs.MouseCommandMode = cmdIdle
			// Niech stoją bezczynnie
			return true
		}

		// 3. Odsianie jednostek, które nie są UNIT_AXEMAN z całej zaznaczonej drużyny
		selectedUnits := getSelectedUnits(bs)

		var repairCrew []*unit

		for _, u := range selectedUnits {
			if u.Type == unitAxeman {
				repairCrew = append(repairCrew, u)
			}
		}

		if len(repairCrew) == 0 {
			bs.CurrentMessage.Text = "Brak Toporników w zaznaczeniu!"
			bs.CurrentMessage.Duration = 60
			bs.MouseCommandMode = cmdIdle

			return true
		}
		// 4. Rozkaz gotowy, wiadomo kto, co, można przekazać dalej
		sendUnitCommand(bs, repairCrew, cmd, tileX, tileY, targetBld.ID, input.IsCtrlKeyDown)
		log.Printf("INPUT: Wysłano %d Toporników do naprawy budynku ID %d.", len(repairCrew), targetBld.ID)

		// Zmieniamy stan myszki i wracamy
		bs.MouseCommandMode = cmdIdle

		return true

	// === 2. RZUCANIE CZARÓW ===
	case cmdCastSpell:
		log.Println("DBG_LCLICK: Tryb rzucania czaru.")
		selectedUnit, ok := getUnitByID(bs.CurrentSelection.UnitID, bs)
		if !ok && selectedUnit.Exists {
			// TODO: Tutaj CMD_MAGIC_FIRE jest na sztywno, docelowo powinno zależeć od wybranego czaru w UI
			selectedUnit.addUnitCommand(cmdMagicFire, tileX, tileY, 0, bs)
			log.Printf("DBG_LCLICK: Wydano komendę czaru dla jednostki %d na (%d,%d).", selectedUnit.ID, tileX, tileY)
		}
		bs.MouseCommandMode = 1
		return true

	// === 3. DOMYŚLNY TRYB (SELEKCJA I RUCH) ===
	default:
		// Sprawdzamy, czy kliknięto w obiekt (Jednostkę lub Budynek)
		tile := &bs.Board.Tiles[tileX][tileY]
		targetID := uint(0)

		if tile.Unit != nil {
			targetID = tile.Unit.ID
		} else if tile.Building != nil {
			targetID = tile.Building.ID
		}

		if targetID != 0 {
			log.Println("DBG_LCLICK: Kliknięto na OBIEKT. Wywołuję selectObjectByClick.")
			selectObjectByClick(tileX, tileY, bs)
			return true
		}

		// Kliknięto w puste pole -> Początek rysowania prostokąta zaznaczenia (Drag Selection)
		log.Println("DBG_LCLICK: Kliknięto na puste pole. Początek zaznaczania.")
		if !rl.IsKeyDown(rl.KeyLeftShift) && !rl.IsKeyDown(rl.KeyRightShift) {
			// Jeśli nie trzymamy Shift, czyścimy poprzednie zaznaczenie
			clearSelection(bs)
		}
		// Zwracamy false, aby pozwolić funkcji nadrzędnej obsłużyć ciągnięcie myszy (drag) w handleBoardDrag
		return false
	}
}

func handleBoardDrag(input inputState, bs *battleState) bool {
	if !input.IsLeftMouseButtonDown {
		return false
	}

	if bs.IsSelectingBox || bs.MouseCommandMode != 1 || bs.CurrentSelection.IsUnit || bs.CurrentSelection.BuildingID != 0 {
		return false
	}

	distance := rl.Vector2Distance(bs.InitialClickPos, input.MousePosition)

	if distance > dragThresholdPixels {
		bs.IsSelectingBox = true
		bs.SelectionStart = bs.InitialClickPos
		return true
	}
	return false
}

// @todo: ogarnij czemu to nie działa jako przekazanie STOP do wszystkich
// zaznaczonych jednostek!
func sendUnitCommand(bs *battleState, units []*unit, command uint16, x, y uint8, targetID uint, ctrlDown bool) {
	for _, u := range units {
		u.AllowFriendlyFire = ctrlDown
	}

	if len(units) > 1 {
		bs.assignGroupCommand(command, x, y, targetID, units)
	} else {
		units[0].addUnitCommand(command, x, y, targetID, bs)
	}

	bs.MouseCommandMode = 1
}

func handleBoardInteraction(input inputState, bs *battleState, ps *programState) {
	// Przekazujemy ps do checks
	handledInitial, tileX, tileY := handleBoardInitialChecks(input, bs, ps)
	if handledInitial {
		return
	}

	if bs.IsSelectingBox && input.IsLeftMouseButtonReleased {
		bs.IsSelectingBox = false
		performBoxSelection(bs, bs.SelectionStart, input.MousePosition)
		bs.SelectionStart = rl.NewVector2(0, 0)
		bs.InitialClickPos = rl.NewVector2(0, 0)
		return
	}

	if handleBoardRightClick(input, bs, tileX, tileY) {
		return
	}

	if input.IsLeftMouseButtonPressed && handleBoardLeftClick(input, bs, tileX, tileY) {
		return
	}

	if input.IsLeftMouseButtonDown {
		handleBoardDrag(input, bs)
	}
}

func handleGameInput(bs *battleState, ps *programState) {
	screenMouse := rl.GetMousePosition()
	virtualMouse := screenToVirtualCoords(ps, screenMouse)

	input := inputState{
		MousePosition:              virtualMouse,
		IsLeftMouseButtonDown:      rl.IsMouseButtonDown(rl.MouseLeftButton),
		IsLeftMouseButtonPressed:   rl.IsMouseButtonPressed(rl.MouseLeftButton),
		IsLeftMouseButtonReleased:  rl.IsMouseButtonReleased(rl.MouseLeftButton),
		IsRightMouseButtonDown:     rl.IsMouseButtonDown(rl.MouseRightButton),
		IsRightMouseButtonPressed:  rl.IsMouseButtonPressed(rl.MouseRightButton),
		IsRightMouseButtonReleased: rl.IsMouseButtonReleased(rl.MouseRightButton),
		IsCtrlKeyDown:              rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl),
	}

	//  if handleCheats(bs) {
	//	return
	//  }

	handleCameraScroll(input, bs, ps)

	if isMouseOverUI(ps, virtualMouse) {
		if input.IsLeftMouseButtonPressed || input.IsRightMouseButtonPressed || input.IsLeftMouseButtonDown ||
			input.IsLeftMouseButtonReleased {
			if handleGameUIClicks(input, bs, ps) {
				return
			}
		}
	} else {
		handleGameShortcuts(bs)
		handleBoardInteraction(input, bs, ps)
	}
}

// OBSŁUGA ZAZNACZANIA

func clearSelection(bs *battleState) {
	log.Println("SELEKCJA: Rozpoczynam clearSelection.")

	for _, unit := range bs.Units {
		if unit.Exists && unit.IsSelected {
			unit.IsSelected = false
		}
	}

	if bs.CurrentSelection.IsUnit || bs.CurrentSelection.BuildingID != 0 {
		bs.CurrentSelection = selectionState{}
	}
}

func selectObjectByClick(tileX, tileY uint8, bs *battleState) {
	tile := &bs.Board.Tiles[tileX][tileY]
	unit := tile.Unit
	bld := tile.Building

	if unit == nil && bld == nil {
		found := false
		originalTileX, originalTileY := tileX, tileY
		log.Printf("DBG_SELECTOBJECT: Na (%d,%d) nie ma bezpośredniego obiektu. Szukam w sąsiedztwie...", tileX, tileY)

		for j := originalTileY - 1; j <= originalTileY+1; j++ {
			for i := originalTileX - 1; i <= originalTileX+1; i++ {
				if i < boardMaxX && j < boardMaxY {
					nt := &bs.Board.Tiles[i][j]
					if nt.Unit != nil || nt.Building != nil {
						unit = nt.Unit
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
			clearSelection(bs)
			bs.MouseCommandMode = 1
			return
		}
	}

	isShiftDown := rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)

	if unit != nil && unit.Exists {
		log.Printf("DBG_SELECTOBJECT: Znaleziono jednostkę ID %d.", unit.ID)

		if unit.Owner != bs.PlayerID {
			clearSelection(bs)
			bs.CurrentSelection = selectionState{
				OwnerID:    unit.Owner,
				IsUnit:     true,
				UnitID:     unit.ID,
				BuildingID: 0,
			}
			bs.CurrentMessage.Text = fmt.Sprintf("Wroga jednostka: %v", unit.Type)
			bs.CurrentMessage.Duration = 20
			bs.MouseCommandMode = 1
			return
		}

		if isShiftDown {
			unit.IsSelected = !unit.IsSelected
			if !unit.IsSelected && bs.CurrentSelection.UnitID == unit.ID {
				foundNewPrimary := false
				for _, u := range bs.Units {
					if u.Exists && u.IsSelected && u.Owner == bs.PlayerID {
						bs.CurrentSelection = selectionState{OwnerID: u.Owner, IsUnit: true, UnitID: u.ID}
						foundNewPrimary = true
						break
					}
				}
				if !foundNewPrimary {
					bs.CurrentSelection = selectionState{}
				}
			} else if unit.IsSelected && !bs.CurrentSelection.IsUnit {
				bs.CurrentSelection = selectionState{OwnerID: unit.Owner, IsUnit: true, UnitID: unit.ID}
			}

		} else {
			clearSelection(bs)
			unit.IsSelected = true
			bs.CurrentSelection = selectionState{
				OwnerID:    unit.Owner,
				IsUnit:     true,
				UnitID:     unit.ID,
				BuildingID: 0,
			}
		}

		if unit.IsSelected && unit.Owner == bs.PlayerID {
			switch unit.Type {
			case unitCow:
				bs.CurrentMessage.Text = "Muuu ?"
			case unitAxeman:
				bs.CurrentMessage.Text = "Tak ?"
			default:
				bs.CurrentMessage.Text = "Rozkaz?"
			}
			bs.CurrentMessage.Duration = 20
		}
		bs.MouseCommandMode = 1

	} else if bld != nil && bld.Exists {
		log.Printf("DBG_SELECTOBJECT: Znaleziono budynek ID %d.", bld.ID)
		if !isShiftDown {
			clearSelection(bs)
		}

		bs.CurrentSelection = selectionState{
			OwnerID:    bld.Owner,
			IsUnit:     false,
			UnitID:     0,
			BuildingID: bld.ID,
		}
		bs.MouseCommandMode = 1

		if bld.Owner == bs.PlayerID {
			bs.CurrentMessage.Text = fmt.Sprintf("Moja budowla: %v", bld.Type)
		} else {
			bs.CurrentMessage.Text = fmt.Sprintf("Wroga budowla: %v", bld.Type)
		}
		bs.CurrentMessage.Duration = 20
	} else {
		clearSelection(bs)
		bs.MouseCommandMode = 1
	}
}

func performBoxSelection(bs *battleState, startPos, endPos rl.Vector2) {
	worldStart := rl.GetScreenToWorld2D(startPos, bs.GameCamera)
	worldEnd := rl.GetScreenToWorld2D(endPos, bs.GameCamera)

	minX := uint8(min(worldStart.X, worldEnd.X) / float32(tileWidth))
	maxX := uint8(max(worldStart.X, worldEnd.X) / float32(tileWidth))
	minY := uint8(min(worldStart.Y, worldEnd.Y) / float32(tileHeight))
	maxY := uint8(max(worldStart.Y, worldEnd.Y) / float32(tileHeight))

	isShiftDown := rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)
	if !isShiftDown {
		clearSelection(bs)
	}

	var selectedCount int
	var firstSelectedUnit *unit

	for _, unit := range bs.Units {
		if unit.Exists && unit.Owner == bs.PlayerID {
			if unit.X >= minX && unit.X <= maxX && unit.Y >= minY && unit.Y <= maxY {
				if isShiftDown {
					unit.IsSelected = !unit.IsSelected
				} else {
					unit.IsSelected = true
				}
			} else {
				if !isShiftDown && unit.IsSelected {
					unit.IsSelected = false
				}
			}

			if unit.IsSelected {
				selectedCount++

				if firstSelectedUnit == nil {
					firstSelectedUnit = unit
				}
			}
		}
	}

	if selectedCount > 0 {
		bs.CurrentSelection = selectionState{
			OwnerID:    firstSelectedUnit.Owner,
			IsUnit:     true,
			UnitID:     firstSelectedUnit.ID,
			BuildingID: 0,
		}
	} else {
		clearSelection(bs)
	}
}

func handleMinimapInteraction(input inputState, bs *battleState, ps *programState) bool {
	// Prostokąt minimapy obliczany dynamicznie!
	// ps.GameViewWidth to początek panelu UI.
	minimapRect := rl.NewRectangle(
		ps.GameViewWidth+minimapOffsetX,
		float32(0)+minimapOffsetY,
		minimapDisplayWidth,
		minimapDisplayHeight,
	)

	isMouseOverMinimap := rl.CheckCollisionPointRec(
		input.MousePosition,
		minimapRect,
	)

	if input.IsLeftMouseButtonReleased && bs.MapInitialClickPos.X != 0.0 {
		bs.IsMapDragging = false
		bs.MapInitialClickPos = rl.NewVector2(0.0, 0.0)
		bs.CameraTargetOnDragStart = rl.NewVector2(0.0, 0.0)
		return true
	}

	if !isMouseOverMinimap && !bs.IsMapDragging {
		return false
	}

	// Przekazujemy ps do clamping
	if handleMinimapLeftMouse(input, bs, minimapRect, ps) {
		return true
	}

	if isMouseOverMinimap {
		if handleMinimapRightMouse(input, bs, minimapRect) {
			return true
		}
	}

	return false
}

func handleMinimapLeftMouse(input inputState, bs *battleState, minimapRect rl.Rectangle, ps *programState) bool {
	fullMapPixelWidth := float32(uint16(boardMaxX) * uint16(tileWidth))
	fullMapPixelHeight := float32(uint16(boardMaxY) * uint16(tileHeight))
	scaleX := fullMapPixelWidth / minimapDisplayWidth
	scaleY := fullMapPixelHeight / minimapDisplayHeight

	if input.IsLeftMouseButtonPressed {
		clickedX := (input.MousePosition.X - minimapRect.X) * scaleX
		clickedY := (input.MousePosition.Y - minimapRect.Y) * scaleY
		bs.GameCamera.Target = rl.NewVector2(clickedX, clickedY)

		// Clamping używa dynamicznych wymiarów
		clampCameraTarget(&bs.GameCamera, fullMapPixelWidth, fullMapPixelHeight,
			ps.GameViewWidth, ps.VirtualHeight)

		bs.MapInitialClickPos = input.MousePosition
		bs.CameraTargetOnDragStart = bs.GameCamera.Target
		bs.IsMapDragging = false
		return true
	}

	if input.IsLeftMouseButtonDown && bs.MapInitialClickPos.X != 0.0 {
		if !bs.IsMapDragging &&
			rl.Vector2Distance(bs.MapInitialClickPos, input.MousePosition) > minimapClickDragThreshold {
			bs.IsMapDragging = true
		}

		if bs.IsMapDragging {
			deltaX := (input.MousePosition.X - bs.MapInitialClickPos.X) * scaleX
			deltaY := (input.MousePosition.Y - bs.MapInitialClickPos.Y) * scaleY
			bs.GameCamera.Target.X = bs.CameraTargetOnDragStart.X + deltaX
			bs.GameCamera.Target.Y = bs.CameraTargetOnDragStart.Y + deltaY

			// Clamping używa dynamicznych wymiarów
			clampCameraTarget(&bs.GameCamera, fullMapPixelWidth, fullMapPixelHeight,
				ps.GameViewWidth, ps.VirtualHeight)
			return true
		}
	}

	return false
}

func handleMinimapRightMouse(
	input inputState,
	bs *battleState,
	minimapRect rl.Rectangle,
) bool {
	if !input.IsRightMouseButtonPressed {
		return false
	}

	selectedUnits := getSelectedUnits(bs)
	if len(selectedUnits) == 0 {
		return true
	}

	scaleX := float32(uint16(boardMaxX)*uint16(tileWidth)) / minimapDisplayWidth
	scaleY := float32(uint16(boardMaxY)*uint16(tileHeight)) / minimapDisplayHeight

	worldX := float64(input.MousePosition.X-minimapRect.X) * float64(scaleX)
	worldY := float64(input.MousePosition.Y-minimapRect.Y) * float64(scaleY)
	tileX := uint8(math.Min(math.Max(worldX/float64(tileWidth), 0), float64(boardMaxX-1)))
	tileY := uint8(math.Min(math.Max(worldY/float64(tileHeight), 0), float64(boardMaxY-1)))

	tile := &bs.Board.Tiles[tileX][tileY]
	var targetID uint = 0
	var targetOwner uint8 = 0
	if tile.Unit != nil {
		targetID = tile.Unit.ID
		targetOwner = tile.Unit.Owner
	} else if tile.Building != nil {
		targetID = tile.Building.ID
		targetOwner = tile.Building.Owner
	}

	cmd := cmdMove

	if targetID != 0 && targetOwner != bs.PlayerID {
		cmd = cmdAttack
	} else if !isWalkable(bs, tileX, tileY) {
		bs.CurrentMessage.Text = "Nieprzechodnie!"
		bs.CurrentMessage.Duration = 60
		return true
	}

	sendUnitCommand(bs, selectedUnits, cmd, tileX, tileY, targetID, input.IsCtrlKeyDown)
	bs.MouseCommandMode = 1
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
