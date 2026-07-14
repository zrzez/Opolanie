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
func (u *unit) calculateApproachTile(intention point, targetID ObjectID, bState *battleState) (point, error) {
	// @reminder: Jeszcze nie wiem czemu miałbym tak rozdzielać odnajdywanie celu, ale niech będzie
	if u.CurrentSpell != spellNone {
		approachTile, err := u.findApproachTileForSpell(intention, bState)
		if err != nil {
			return point{X: 0, Y: 0}, err
		}

		return approachTile, nil
	}

	// Budynki, jednostki i drzewa jako cel
	return u.findApproachTileForTarget(intention, targetID, bState)
}

func (u *unit) findApproachTileForSpell(targetPosition point, bState *battleState) (point, error) {
	switch u.CurrentSpell {
	case spellMagicShower:
		// Tutaj ustalamy, gdzie kapłan/-ka mają stanąć, ażeby rzucić czar.
		finalPoint, ok := u.findBestPositionAroundTile(targetPosition, bState.Board)

		if !ok {
			return point{}, fmt.Errorf("brak miejsca do rzucenia czaru")
		}

		return finalPoint, nil

	case spellMagicShield, spellMagicSight:
		// Czary, które przyjmują rzucającego jako swój cel.
		return point{X: u.X, Y: u.Y}, nil
	case spellNone:
		// To nigdy nie powinno mieć miejsca, bo warunkiem wywołania
		// jest u.CurrentSpell != spellNone.
		return point{X: u.X, Y: u.Y}, fmt.Errorf("próba rzucenia spellNone")
	default:
		// To nigdy nie powinno mieć miejsca, bo wszystkie czary są obsłużone
		return point{X: u.X, Y: u.Y}, fmt.Errorf("nieznany rodzaj czaru")
	}
}

// findBestPostitionAroundTile ma niewłaściwą nazwę. Zwraca współrzędne kafelka z którego
// można napaść na cel. Nie gwarantuje najkrótszej drogi.
// Używany jedynie do czarodziejskich opadów.
// Chyba faworyzuje lewy-górny róg.
func (u *unit) findBestPositionAroundTile(targetTile point, board *boardData) (point, bool) {
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

	return point{X: bestX, Y: bestY}, found
}

func (u *unit) findApproachTileForTarget(intention point, targetID ObjectID, bState *battleState) (point, error) {
	targetUnit, targetBuilding := bState.getObjectByID(targetID)

	if targetBuilding != nil && (targetBuilding.Exists || targetBuilding.Type == buildingBridge) {
		return findTileForBuilding(u, targetBuilding, bState.Board)
	}

	if targetUnit != nil && targetUnit.Exists {
		return findTileForUnit(u, targetUnit, bState.Board)
	}

	if bState.Board.Tiles[intention.X][intention.Y].isTree() {
		return findTileForTree(u, intention, bState.Board)
	}

	return point{X: 0, Y: 0}, fmt.Errorf("cel ataku ID %d nie istnieje", targetID)
}

func findTileForBuilding(u *unit, b *building, board *boardData) (point, error) {
	// 1. Próba ataku dystansowego
	if u.AttackRange > 1 {
		x, y, ok := findOptimalRangedAttackTile(u.X, u.Y, u.AttackRange, b, board)
		if ok {
			return point{X: x, Y: y}, nil
		}
	}

	// 2. Szukanie najbliższego wolnego sąsiedniego kafelka
	coords := board.neighborCoords(b)

	var bestX, bestY uint8

	minPathLen := math.MaxInt32
	found := false

	for _, coord := range coords {
		electedTile := &board.Tiles[coord.X][coord.Y]

		if electedTile.IsWalkable && electedTile.Unit == nil {
			path := findPath(board, u, u.X, u.Y, coord.X, coord.Y)

			if path != nil && len(path) < minPathLen {
				minPathLen = len(path)
				bestX, bestY = coord.X, coord.Y
				found = true
			}
		}
	}

	if found {
		return point{X: bestX, Y: bestY}, nil
	}

	return point{X: 0, Y: 0}, fmt.Errorf("cel (budynek) jest nieosiągalny")
}

func findTileForUnit(u *unit, target *unit, board *boardData) (point, error) {
	bestX, bestY := u.findBestPositionAroundUnit(target, board) // zakładam, że zmienisz sygnaturę na przyjmującą *Board

	// Jeśli zwrócono pozycję celu, sprawdź, czy nie jest to fallback (kafelek zajęty przez cel)
	if bestX == target.X && bestY == target.Y {
		if board.Tiles[bestX][bestY].Unit == target {
			return point{X: 0, Y: 0}, fmt.Errorf("brak wolnego kafelka wokół jednostki ID %d", target.ID)
		}
	}

	return point{X: bestX, Y: bestY}, nil
}

func findTileForTree(u *unit, treeTile point, board *boardData) (point, error) {
	bestX, bestY, ok := u.findOptimalAttackTileAroundTree(treeTile.X, treeTile.Y, board)
	if !ok {
		return point{X: 0, Y: 0}, fmt.Errorf("nie ma pozycji do ataku tego drzewa")
	}

	return point{X: bestX, Y: bestY}, nil
}
