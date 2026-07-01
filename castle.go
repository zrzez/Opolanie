package main

import (
	"fmt"
	"log"
)

// castle.go

// ============================================================================
// SYSTEM ZAMKÓW - METODY CASTLE
// ============================================================================

// Inicjuje stan zamku gracza/SI.
func (playerS *playerState) init(factionID uint8, maxMilk uint16) {
	playerS.PlayerID = factionID
	playerS.MaxMilk = maxMilk
	playerS.Milk = maxMilk // Początkowe mleko ustawione na MaxMilk
	log.Printf("playerState (FactionID %d) zaczął z MaxMilk: %d", playerS.PlayerID, playerS.MaxMilk)
}

// setCommand przetwarza rozkaz (np. wytwarzania, ruchu, ataku) dla gracza.
// @todo: wróć i posprzątaj
func (playerS *playerState) setCommand(cmd *command, bState *battleState) {
	if cmd == nil {
		return
	}

	// @reminder: staram się tutaj przechwycić rozkazy dla wszystkich zaznaczonych jednsotek
	if cmd.ActionType < cmdDelimiter && cmd.ExecutorID == 0 {
		selectedUnits := getSelectedUnits(bState)
		for _, u := range selectedUnits {
			unitCmd := *cmd
			unitCmd.ExecutorID = u.ID

			playerS.handleUnitCommand(&unitCmd, bState)
		}

		return
	}

	// @reminder: tutaj obsługujemy rozkazy dla pojedynczych jednostek i budynków
	if cmd.ActionType > cmdDelimiter { // Rozkaz dla budynku
		playerS.handleBuildingCommand(cmd, bState)
	} else {
		// @todo: czemu niby miałbym tutaj znowu walidować? Nie rozumiem, do tego ponowne sprawdzanie FF?
		targetUnit, ok := getUnitByID(cmd.ExecutorID, bState)
		if ok {
			targetUnit.AllowFriendlyFire = cmd.FriendlyFire
		}

		playerS.handleUnitCommand(cmd, bState)
	}
}

// handleBuildingCommand przetwarza rozkazy dotyczące budynków.
// Obsługuje np. produkcję jednostek.
func (playerS *playerState) handleBuildingCommand(cmd *command, bState *battleState) {
	// 1. Dla zasadzenia nowej budowy
	if cmd.ActionType == cmdBPlaceConstruction {
		playerS.handleConstructionCommand(cmd, bState)
	}

	// 2. Wytwarzanie jednostek
	targetBuilding, ok := getBuildingByID(cmd.ExecutorID, bState)

	if !ok || !targetBuilding.Exists || targetBuilding.IsUnderConstruction {
		log.Printf("handleBuildingCommand: Nie znaleziono ID %d, nie istnieje lub w budowie.", cmd.InteractionTargetID)

		return
	}

	if targetBuilding.Owner != playerS.PlayerID {
		log.Printf("handleBuildingCommand: Próba wydania komendy dla budynku %d, który nie należy do frakcji %d.",
			targetBuilding.ID, playerS.PlayerID)

		return
	}

	switch cmd.ActionType {
	case cmdBProduce:
		targetBuilding.produceUnit(unitType(cmd.CreateType), bState)
	default:
		log.Printf("handleBuildingCommand: Niezaimplementowany ActionType %d dla budynku %d.",
			cmd.ActionType, targetBuilding.ID)
	}
}

func (playerS *playerState) handleConstructionCommand(cmd *command, bState *battleState) {
	// 0. Poprawność
	bType := buildingType(cmd.CreateType)

	if bType == 0 {
		bState.CurrentMessage.Text = "Błąd: nie określono typu budynku. " // nie powinno się wydarzyć nigdy…
		bState.CurrentMessage.Duration = 60

		return
	}

	// Walidacja kontekstowa
	if ok, errCode := validateConstructionContext(bType, playerS.PlayerID, bState); !ok {
		switch errCode {
		case buildErrLimit:
			bState.CurrentMessage.Text = "Limit budynków!"
			bState.CurrentMessage.Duration = 60
		case buildErrMilk:
			stats := buildingDefs[bType]
			bState.CurrentMessage.Text = fmt.Sprintf("Niedobór mleka! (%d)", stats.Cost)
			bState.CurrentMessage.Duration = 60
		default:
			bState.CurrentMessage.Text = "COŚ POSZŁO PIERUŃSKO NIE TAK Z KONTEKSTEM!"
			bState.CurrentMessage.Duration = 60
		}

		return
	}

	// Walidacja środowiskowa
	if ok, errCode := validateConstructionSite(bType, cmd.TargetX, cmd.TargetY, bState); !ok {
		switch errCode {
		case buildErrOutofBounds:
			bState.CurrentMessage.Text = "Poza mapą!"
			bState.CurrentMessage.Duration = 40
		case buildErrOccupiedUnit:
			bState.CurrentMessage.Text = "Miejsce zajęte przez jednostkę!"
			bState.CurrentMessage.Duration = 40
		case buildErrOccupiedBuilding:
			bState.CurrentMessage.Text = "Miejsce zajęte przez budynek!"
			bState.CurrentMessage.Duration = 40
		case buildErrObstacle:
			bState.CurrentMessage.Text = "Nie można na tym!"
			bState.CurrentMessage.Duration = 40
		case buildErrWater:
			bState.CurrentMessage.Text = "Nie można na wodzie!"
			bState.CurrentMessage.Duration = 40
		case buildErrNoWater:
			bState.CurrentMessage.Text = "Tu nie ma wody!"
			bState.CurrentMessage.Duration = 40
		case buildErrNoRoadAccess:
			bState.CurrentMessage.Text = "Brak dostępu do drogi!"
			bState.CurrentMessage.Duration = 40
		case buildErrAlreadyPath:
			bState.CurrentMessage.Text = "Tu już jest droga!"
			bState.CurrentMessage.Duration = 40
		case buildErrNone:
		default:
			bState.CurrentMessage.Text = "COŚ POSZŁO PIERUŃSKO NIE TAK ZE ŚRODOWISKIEM!"
			bState.CurrentMessage.Duration = 60
		}
		return
	}

	// 1. Wykonanie
	tryBuildStructure(bType, cmd.TargetX, cmd.TargetY, playerS.PlayerID, bState)

	if bType == buildingRoad || bType == buildingPalisade || bType == buildingBridge {
		return
	}

	// 2. Zakończenie, czyścimy
	log.Printf("[castle.go] Przyjęto rozkaz budowy: %d (%d,%d)", bType, cmd.TargetX, cmd.TargetY)
	bState.PendingCommand = nil
	bState.MouseState = mouseStateNormal
}

// handleUnitCommand przetwarza rozkazy dotyczące jednostek.
// Obsługuje np. ruch, atak, stop, magię.
func (playerS *playerState) handleUnitCommand(cmd *command, bState *battleState) {
	targetUnit, ok := getUnitByID(cmd.ExecutorID, bState)
	if !ok || !targetUnit.Exists {
		log.Printf("handleUnitCommand: Nie znaleziono jednostki ID %d lub nie istnieje.", cmd.ExecutorID)

		return
	}

	if targetUnit.Owner != playerS.PlayerID {
		log.Printf("handleUnitCommand: Próba wydania komendy jednostce ID %d, która nie należy do frakcji %d.",
			targetUnit.ID, playerS.PlayerID)

		return
	}

	switch cmd.ActionType {
	case cmdUMove:
		log.Printf("INFO: castle.go wydano cmdMove.")
		// playerS.handleMoveCommand(cmd, targetUnit, bState) nie działa
		targetUnit.addUnitCommand(cmd, bState)
	case cmdUAttack:
		log.Printf("INFO: castle.go wydano cmdAttack.")
		targetUnit.addUnitCommand(cmd, bState)
	case cmdUStop:
		targetUnit.addUnitCommand(cmd, bState) // czemu targetID = 0? może nil?
	case cmdUBuild:
		targetBuilding, ok := getBuildingByID(cmd.InteractionTargetID, bState)
		if !ok {
			return
		}

		if valid, errCode := validateBuildingContext(targetUnit, targetBuilding); !valid {
			switch errCode {
			case workErrWrongWorkerType:
				bState.CurrentMessage.Text = "Ten rodzaj jednostki nie może naprawiać!"
				bState.CurrentMessage.Duration = 60
			case workErrInvalidTarget:
				bState.CurrentMessage.Text = "Tego nie da się budować!"
				bState.CurrentMessage.Duration = 60
			case workErrNotUnderConstruction:
				bState.CurrentMessage.Text = "Nie wymaga budowy!"
				bState.CurrentMessage.Duration = 60
			case workErrNone:
			default:
				bState.CurrentMessage.Text = "COŚ POSZŁO PIERUŃSKO NIE TAK ZE SPRAWDZANIEM BUDOWY!"
				bState.CurrentMessage.Duration = 60
			}
			return
		}

		targetUnit.addUnitCommand(cmd, bState)
		log.Printf("handleUnitCommand: Jednostka %d otrzymała rozkaz BUDOWY budynku %d.",
			targetUnit.ID, cmd.InteractionTargetID)
	case cmdURepair:
		targetBuilding, ok := getBuildingByID(cmd.InteractionTargetID, bState)
		if !ok {
			return
		}
		if valid, errCode := validateRepairContext(targetUnit, targetBuilding); !valid {
			switch errCode {
			case workErrWrongWorkerType:
				bState.CurrentMessage.Text = "Ten rodzaj jednostki nie może naprawiać!"
				bState.CurrentMessage.Duration = 40
			case workErrInvalidTarget:
				bState.CurrentMessage.Text = "Tego nie da się naprawić!"
				bState.CurrentMessage.Duration = 40
			case workErrNotRepairable:
				bState.CurrentMessage.Text = "Nie wymaga naprawy!"
				bState.CurrentMessage.Duration = 40
			case workErrNone:
			default:
				bState.CurrentMessage.Text = "COŚ POSZŁO PIERUŃSKO NIE TAK ZE SPRAWDZANIEM NAPRAWY!"
				bState.CurrentMessage.Duration = 60
			}
			return
		}

		targetUnit.addUnitCommand(cmd, bState)
		log.Printf("handleUnitCommand: Jednostka %d otrzymała rozkaz NAPRAWY budynku %d.",
			targetUnit.ID, cmd.InteractionTargetID)
	case cmdUCastSpell:
		// @todo: dlaczego zakomentowanie tego psuje rzucanie czarów?
		// Przecież to powinno tylko przepchnąć rozkaz we właściwe miejsce do sprawdzenia poprawności!
		// targetUnit.CurrentSpell = cmd.Spell
		targetUnit.addUnitCommand(cmd, bState)
	default:
		log.Printf("handleUnitCommand: Nieznany ActionType %d dla jednostki %d.",
			cmd.ActionType, targetUnit.ID)
	}
}

// handleMoveCommand obsługuje logikę rozkazu ruchu dla jednostki.
// Sprawdza dostępność celu i wyznacza ścieżkę.
func (playerS *playerState) handleMoveCommand(cmd *command, unit *unit, bState *battleState) {
	log.Printf("DEBUG: handleMoveCommand: Rozkaz ruchu dla jednostki ID %d do (%d,%d).",
		unit.ID, cmd.TargetX, cmd.TargetY)

	// 1. Sprawdzenie czy kafelek jest przechodni (używamy isWalkable)
	if !isWalkable(bState, cmd.TargetX, cmd.TargetY) {
		// Pobieramy ID tekstury z nowej struktury Tiles
		var terrainID uint16
		if cmd.TargetX < boardMaxX && cmd.TargetY < boardMaxY {
			terrainID = bState.Board.Tiles[cmd.TargetX][cmd.TargetY].TextureID
		}

		log.Printf(
			"handleMoveCommand: ODRZUCONO ROZKAZ: Cel (%d,%d) jest nieprzechodni (TextureID: %d). Jednostka ID %d.",
			cmd.TargetX, cmd.TargetY, terrainID, unit.ID,
		)

		return
	}

	path := findPath(
		bState,
		unit.ID,
		unit.X,
		unit.Y,
		cmd.TargetX,
		cmd.TargetY,
	)

	if len(path) == 0 {
		log.Printf(
			"DEBUG: handleMoveCommand: Pathfinding nie znalazł ścieżki dla jednostki %d do (%d,%d). Komenda odrzucona.",
			unit.ID, cmd.TargetX, cmd.TargetY,
		)

		return
	}

	unit.Path = path
	unit.PathIndex = 0
	log.Printf(
		"handleMoveCommand: Ścieżka znaleziona, ustawiona dla jednostki %d. Długość: %d. Rozkaz MOVE do (%d,%d).",
		unit.ID,
		len(path),
		cmd.TargetX,
		cmd.TargetY,
	)
}

// szuka budynku w battleState.Buildings.
func getBuildingByID(bldID uint, bState *battleState) (*building, bool) {
	for _, bld := range bState.Buildings {
		if bld.ID == bldID {
			return bld, true
		}
	}

	return nil, false
}

// szuka jednostki w battleState.Units.
func getUnitByID(unitID uint, bState *battleState) (*unit, bool) {
	for _, currentUnit := range bState.Units {
		if currentUnit.ID == unitID {
			return currentUnit, true
		}
	}

	return nil, false
}

func getObjectByID(objectID uint, bState *battleState) (*unit, *building) {
	if currentUnit, ok := getUnitByID(objectID, bState); ok {
		return currentUnit, nil
	}

	if currentBuilding, ok := getBuildingByID(objectID, bState); ok {
		return nil, currentBuilding
	}

	return nil, nil
}
