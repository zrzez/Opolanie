package main

import (
	"fmt"
)

// ============
// Próba rozplątania units.go, tutaj powinny trafiać funkcje związane z
// przetwarzaniem rozkazów przez jednoski.

func (u *unit) calculateApproachTile(intention *point, targetID ObjectID, bState *battleState) (*point, error) {
	if u.CurrentSpell != spellNone {
		approachTile, err := u.findApproachTileForSpell(intention, bState)
		if err != nil {
			return nil, err
		}

		return approachTile, nil
	}

	// Budynki, jednostki i drzewa jako cel
	return u.findApproachTileForTarget(intention, targetID, bState)
}

func (u *unit) findApproachTileForSpell(targetPosition *point, bState *battleState) (*point, error) {
	switch u.CurrentSpell {
	case spellMagicShower:

		validCoords, ok := findTileForAttacking(u, nil, nil, targetPosition, bState.Board)
		if !ok {
			return nil, fmt.Errorf("nie ma podejścia do celu")
		}

		return findBestReachableTile(u, validCoords, bState.Board)

	// ↓↓↓↓↓ Poniższe przypadki nie muszą korzystać z A*
	case spellMagicShield, spellMagicSight:
		// Czary, które przyjmują rzucającego jako swój cel.
		return &point{X: u.X, Y: u.Y}, nil
	case spellNone:
		// To nigdy nie powinno mieć miejsca, bo warunkiem wywołania
		// jest u.CurrentSpell != spellNone.
		return &point{X: u.X, Y: u.Y}, fmt.Errorf("próba rzucenia spellNone")
	default:
		// To nigdy nie powinno mieć miejsca, bo wszystkie czary są obsłużone
		return &point{X: u.X, Y: u.Y}, fmt.Errorf("nieznany rodzaj czaru")
	}
}

// @reminder: nazwa dla kafelka z drzewem „intention” jest bardzo kiepska, ale nie mam teraz do tego głowy.
func (u *unit) findApproachTileForTarget(intention *point, targetID ObjectID, bState *battleState) (*point, error) {
	targetUnit, targetBuilding := bState.getObjectByID(targetID)

	var targetU *unit

	var targetBld *building

	var targetTree *point

	switch {
	case targetBuilding != nil && (targetBuilding.Exists || targetBuilding.Type == buildingBridge):
		targetBld = targetBuilding
	case targetUnit != nil && targetUnit.Exists:
		targetU = targetUnit
	case intention != nil && bState.Board.Tiles[intention.X][intention.Y].isTree():
		targetTree = intention
	}

	validCoords, ok := findTileForAttacking(u, targetU, targetBld, targetTree, bState.Board)
	if !ok {
		return nil, fmt.Errorf("nie ma podejścia do celu: %t", ok)
	}

	return findBestReachableTile(u, validCoords, bState.Board)
}
