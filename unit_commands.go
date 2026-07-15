package main

import (
	"fmt"
	"math"
)

// ============
// Próba rozplątania units.go, tutaj powinny trafiać funkcje związane z
// przetwarzaniem rozkazów przez jednoski.

// Zastanów się, jak wygląda przepływ w tej metodzie, bo coś mi tutaj nie pasuje.
// Czary nie są w przełączniku, potem „z kapelusza” przechodzimy do zwykłych
// rodzajów ataku.
// @reminder:  pracując nad tą metodą zauważyłem, że psuje się podchodzenie do mostów i palisad.
//    jednostka „zatyka się” w pewnej odległości (3) nie przechodzi dalej. Naprawa działa, budowa i atak nie.
func (u *unit) calculateApproachTile(intention *point, targetID ObjectID, bState *battleState) (*point, error) {
	// @reminder: Jeszcze nie wiem czemu miałbym tak rozdzielać odnajdywanie celu, ale niech będzie
	if u.CurrentSpell != spellNone {
		// @reminder: nie korzysta z A*
		approachTile, err := u.findApproachTileForSpell(intention, bState)
		if err != nil {
			return &point{X: 0, Y: 0}, err
		}

		return approachTile, nil
	}

	// Budynki, jednostki i drzewa jako cel
	// @reminder: częściowo korzysta z A*
	return u.findApproachTileForTarget(intention, targetID, bState)
}

func (u *unit) findApproachTileForSpell(targetPosition *point, bState *battleState) (*point, error) {
	switch u.CurrentSpell {
	case spellMagicShower:
		// Tutaj ustalamy, gdzie kapłan/-ka mają stanąć, ażeby rzucić czar.
		// @reminder: nie korzysta z A*
		finalPoint, ok := u.findBestPositionAroundTile(targetPosition, bState.Board)

		if !ok {
			return &point{}, fmt.Errorf("brak miejsca do rzucenia czaru")
		}

		return finalPoint, nil

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

// findBestPostitionAroundTile ma niewłaściwą nazwę. Zwraca współrzędne kafelka z którego
// można napaść na cel. Nie gwarantuje najkrótszej drogi.
// Używany jedynie do czarodziejskich opadów.
// Chyba faworyzuje lewy-górny róg.
// @reminder: nie korzysta z A*.
// @todo: sprawdź, czy rzeczywiście rzucając czar wybieram jeden konkretny kafelek.
//   jeśli tak, to mogę albo połączyć logikę z szukaniem kafelka o najkrótszej
//   drodze ze strzelaniem albo wykorzystać znajomość zasięgu unitPriest i unitPriestess.
//   wtedy mogę brać przykład z boardData.bldNeighborCoords
func (u *unit) findBestPositionAroundTile(targetTile *point, board *boardData) (*point, bool) {
	bestX, bestY := targetTile.X, targetTile.Y
	minDist := 100
	found := false

	for dx := -int(u.AttackRange); dx <= int(u.AttackRange); dx++ {
		checkX := int(targetTile.X) + dx

		if checkX < 0 || checkX >= int(boardMaxX) {
			continue
		}

		for dy := -int(u.AttackRange); dy <= int(u.AttackRange); dy++ {
			checkY := int(targetTile.Y) + dy

			if checkY < 0 || checkY >= int(boardMaxY) {
				continue
			}

			if !isWalkable(board, uint8(checkX), uint8(checkY)) {
				continue
			}

			// Pomijamy kafelki zajęte przez inne jednostki
			// @reminder: to nie jest objęte sprawdzeniem isWalkable
			currentTile := &board.Tiles[checkX][checkY]
			if currentTile.Unit != nil && currentTile.Unit.ID != u.ID {
				continue
			}

			distX := abs(int(u.X) - checkX)
			distY := abs(int(u.Y) - checkY)
			dist := distX

			dist = max(dist, distY)

			if dist < minDist {
				minDist = dist
				bestX, bestY = uint8(checkX), uint8(checkY)
				found = true
			}
		}
	}

	return &point{X: bestX, Y: bestY}, found
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
		return &point{X: 0, Y: 0}, fmt.Errorf("nie ma podejścia do celu: %t", ok)
	}

	// Przeliczamy długość drogi do odsianych kafelków.
	var bestX, bestY uint8

	minPathLen := math.MaxInt32
	found := false

	for _, coord := range validCoords {
		path := findPath(bState.Board, u, u.X, u.Y, coord.X, coord.Y)

		if path != nil && len(path) < minPathLen {
			minPathLen = len(path)
			bestX, bestY = coord.X, coord.Y
			found = true
		}
	}

	// Zwracamy kafelek dający najkrótszą drogę.
	if found {
		return &point{X: bestX, Y: bestY}, nil
	}

	return &point{X: 0, Y: 0}, fmt.Errorf("nie ma podejścia do celu: %t", found)
}
