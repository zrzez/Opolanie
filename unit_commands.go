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
func (u *unit) calculateApproachTile(intentionX, intentionY uint8, targetID ObjectID, bState *battleState) (uint8, uint8, error) {
	// @reminder: Jeszcze nie wiem czemu miałbym tak rozdzielać odnajdywanie celu, ale niech będzie
	if u.CurrentSpell != spellNone {
		approachTile, err := u.findApproachTileForSpell(point{X: intentionX, Y: intentionY}, bState)
		if err != nil {
			return 0, 0, err
		}

		return approachTile.X, approachTile.Y, nil
	}

	// Budynki, jednostki i drzewa jako cel
	return u.findApproachTileForTarget(intentionX, intentionY, targetID, bState)
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

			if distY > dist {
				dist = distY
			}

			if dist < minDist {
				minDist = dist
				bestX, bestY = uint8(checkX), uint8(checkY)
				found = true
			}
		}
	}

	return point{X: bestX, Y: bestY}, found
}

func (u *unit) findApproachTileForTarget(intentionX, intentionY uint8, targetID ObjectID, bState *battleState) (uint8, uint8, error) {
	targetUnit, targetBuilding := bState.getObjectByID(targetID)

	// Cel jest budynkiem
	if targetBuilding != nil && (targetBuilding.Exists || targetBuilding.Type == buildingBridge) {
		if u.AttackRange > 1 {
			x, y, ok := findOptimalRangedAttackTile(u.X, u.Y, u.AttackRange, targetBuilding, bState.Board)
			if ok {
				return x, y, nil
			}
		}

		// TUTAJ WPROWADZAM ZMIANĘ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓
		coords := bState.Board.neighborCoords(targetBuilding)
		var bestX, bestY uint8
		minPathLen := math.MaxInt32
		var found bool

		for _, coord := range coords {
			if bState.Board.Tiles[coord.X][coord.Y].IsWalkable && bState.Board.Tiles[coord.X][coord.Y].Unit == nil {
				tempPath := findPath(bState.Board, u, u.X, u.Y, coord.X, coord.Y)

				if tempPath != nil && len(tempPath) < minPathLen {
					minPathLen = len(tempPath)
					bestX, bestY = coord.X, coord.Y
					found = true
				}
			}
		}

		// x, y, ok := targetBuilding.getClosestWalkableTile(bState)
		if found {
			return bestX, bestY, nil
		}

		return 0, 0, fmt.Errorf("cel (budynek) jest nieosiągalny")
	}

	// Cel jest jednostką
	if targetUnit != nil && targetUnit.Exists {
		bestX, bestY := u.findBestPositionAroundUnit(targetUnit, bState)

		if bestX == targetUnit.X && bestY == targetUnit.Y {
			// Sprawdź, czy to naprawdę fallback (kafel jest zajęty przez cel)
			targetTile := &bState.Board.Tiles[bestX][bestY]
			if targetTile.Unit == targetUnit {
				return 0, 0, fmt.Errorf("brak wolnego kafelka wokół jednostki ID %d", targetID)
			}
		}

		return bestX, bestY, nil
	}

	// Cel jest drzewem
	targetTile := &bState.Board.Tiles[intentionX][intentionY]

	if targetTile.isTree() {
		bestX, bestY, ok := u.findOptimalAttackTileAroundTree(intentionX, intentionY, bState.Board)
		if !ok {
			return 0, 0, fmt.Errorf("nie ma pozycji do ataku tego drzewa")
		}

		return bestX, bestY, nil
	}

	return 0, 0, fmt.Errorf("cel ataku ID %d nie istnieje", targetID)
}
