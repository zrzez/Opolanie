package main

import (
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
// Dzieli logikę na podfunkcje dla lepszej czytelności i utrzymania zasady jednej odpowiedzialności.
func (playerS *playerState) setCommand(cmd *command, bState *battleState) {
	if cmd == nil {
		return
	}

	switch cmd.CommandCategory {
	case 0: // Rozkaz dla budynku
		playerS.handleBuildingCommand(cmd, bState)
	case 1: // Komenda dla jednostki
		playerS.handleUnitCommand(cmd, bState)
	default:
		log.Printf("setCommand: Nieznany TargetObject w komendzie: %d", cmd.CommandCategory)
	}
}

// handleBuildingCommand przetwarza rozkazy dotyczące budynków.
// Obsługuje np. produkcję jednostek.
func (playerS *playerState) handleBuildingCommand(cmd *command, bState *battleState) {
	targetBuilding, ok := getBuildingByID(cmd.InteractionTargetID, bState)

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
	case cmdProduce:
		targetBuilding.produceUnit(cmd.ProduceType, bState)
	default:
		log.Printf("handleBuildingCommand: Niezaimplementowany ActionType %d dla budynku %d.",
			cmd.ActionType, targetBuilding.ID)
	}
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
	case cmdMove:
		log.Printf("INFO: castle.go wydano cmdMove.")
		playerS.handleMoveCommand(cmd, targetUnit, bState)
	case cmdAttack, cmdStop, cmdMagicShower, cmdMagicSight:
		// @todo: tutaj chyba jest poważny błąd w logice, bo jednostka zdaje się atakować samą siebie, ale…
		// w grze działa prawidłowo muszę się przyjrzeć i ogarnąć, co się odjaniepawla.
		log.Printf("INFO: castle.go wydano cmdAttack.")

		targetUnit.addUnitCommand(
			cmd.ActionType,
			cmd.TargetX,
			cmd.TargetY,
			cmd.ExecutorID,
			bState,
		)
		log.Printf(
			"handleUnitCommand: Jednostka %d otrzymała komendę %d do (%d,%d).",
			targetUnit.ID,
			cmd.ActionType,
			cmd.TargetX,
			cmd.TargetY,
		)
	case cmdRepairStructure:
		targetUnit.addUnitCommand(
			cmd.ActionType,
			cmd.TargetX,
			cmd.TargetY,
			cmd.InteractionTargetID,
			bState,
		)
		log.Printf("handleUnitCommand: Jednostka %d otrzymała rozkaz NAPRAWY budynku %d.",
			targetUnit.ID, cmd.InteractionTargetID)
	case cmdMagicShield:
		// @todo: jeżeli jednostka ma rzucić czar ochronny na samą siebie, to chyba
		// muszę tutaj inaczej podać targetX, targetY?
		// w ogóle nie czuję po choinkę jest mi ten castle.go skoro sprowadza się do
		// wywołania jednej komendy. Mogę to zrobić w logic_ui bezpośrednio!!!
		targetUnit.addUnitCommand(
			cmd.ActionType,
			cmd.TargetX,
			cmd.TargetY,
			cmd.ExecutorID,
			bState,
		)
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
