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
		// @reminder: nie korzysta z A*
		approachTile, err := u.findApproachTileForSpell(intention, bState)
		if err != nil {
			return point{X: 0, Y: 0}, err
		}

		return approachTile, nil
	}

	// Budynki, jednostki i drzewa jako cel
	// @reminder: częściowo korzysta z A*
	return u.findApproachTileForTarget(intention, targetID, bState)
}

func (u *unit) findApproachTileForSpell(targetPosition point, bState *battleState) (point, error) {
	switch u.CurrentSpell {
	case spellMagicShower:
		// Tutaj ustalamy, gdzie kapłan/-ka mają stanąć, ażeby rzucić czar.
		// @reminder: nie korzysta z A*
		finalPoint, ok := u.findBestPositionAroundTile(targetPosition, bState.Board)

		if !ok {
			return point{}, fmt.Errorf("brak miejsca do rzucenia czaru")
		}

		return finalPoint, nil

	// ↓↓↓↓↓ Poniższe przypadki nie muszą korzystać z A*
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
// @reminder: nie korzysta z A*.
// @todo: sprawdź, czy rzeczywiście rzucając czar wybieram jeden konkretny kafelek.
//   jeśli tak, to mogę albo połączyć logikę z szukaniem kafelka o najkrótszej
//   drodze ze strzelaniem albo wykorzystać znajomość zasięgu unitPriest i unitPriestess.
//   wtedy mogę brać przykład z boardData.bldNeighborCoords
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
		// @reminder: korzysta z A*
		return findTileForBuilding(u, targetBuilding, bState.Board)
	}

	if targetUnit != nil && targetUnit.Exists {
		// @reminder: nie korzysta z A*
		// @reminder: powinno współdzielić logikę z szukaniem
		//    kafelków przy drzewie.
		return findTileForUnit(u, targetUnit, bState.Board)
	}

	if bState.Board.Tiles[intention.X][intention.Y].isTree() {
		// @reminder: nie korzysta z A*
		// @reminder: powinno wspołdzielić logikę z szukaniem
		//    kafelka przy jednostce, ale ogranicać
		//    ten na lewo od drzewa, bo jednostka umrze.
		return findTileForTree(u, intention, bState.Board)
	}

	return point{X: 0, Y: 0}, fmt.Errorf("cel ataku ID %d nie istnieje", targetID)
}

func findTileForBuilding(attacker *unit, bld *building, board *boardData) (point, error) {
	// 0. Tworzymy listę współrzędnych z których można zaatakować budynek.
	// 1. Odsiewamy te, które nie mieszczą się w planszy oraz te na których nie może stanąć napastnik.
	validCoords, err := findTileForAttackingBuilding(attacker, bld, board)
	if err != nil {
		return point{X: 0, Y: 0}, fmt.Errorf("nie ma podejścia do celu: %w", err)
	}

	// 2. Przeliczamy długość drogi do odsianych kafelków.
	var bestX, bestY uint8

	minPathLen := math.MaxInt32
	found := false

	for _, coord := range validCoords {
		path := findPath(board, attacker, attacker.X, attacker.Y, coord.X, coord.Y)

		if path != nil && len(path) < minPathLen {
			minPathLen = len(path)
			bestX, bestY = coord.X, coord.Y
			found = true
		}
	}

	// 3. Zwracamy kafelek dający najkrótszą drogę.
	if found {
		return point{X: bestX, Y: bestY}, nil
	}

	return point{X: 0, Y: 0}, fmt.Errorf("nie ma podejścia do celu: %w", err)
}

func findTileForUnit(u *unit, uTarget *unit, board *boardData) (point, error) {
	bestX, bestY := u.findBestPositionAroundUnit(uTarget, board)

	if bestX == uTarget.X && bestY == uTarget.Y {
		if board.Tiles[bestX][bestY].Unit == uTarget {
			return point{X: 0, Y: 0}, fmt.Errorf("brak wolnego kafelka wokół jednostki ID %d", uTarget.ID)
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

func (u *unit) findBestPositionAroundUnit(targetUnit *unit, board *boardData) (uint8, uint8) {
	bestX, bestY := int(targetUnit.X), int(targetUnit.Y)
	minDist := math.MaxFloat64
	foundFreeSpot := false

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			checkX := int(targetUnit.X) + dx
			checkY := int(targetUnit.Y) + dy

			// @reminder: do wywalenia, bo sprawdza, czy w planszy oraz .Unit i .Building = nil
			//    oraz, czy walkable, które jest obecnie jak tile.IsWalkable.
			if u.isValidMoveTarget(uint8(checkX), uint8(checkY), board) {
				// log.Println("Funkcja findBestPositionAroundUnit isValidMoveTarget = true, szukam freeSpot")
				dist := math.Abs(float64(int(u.X)-checkX)) + math.Abs(float64(int(u.Y)-checkY))
				if dist < minDist {
					minDist = dist
					bestX, bestY = checkX, checkY
					foundFreeSpot = true
				}
			}
		}
	}

	if !foundFreeSpot {
		// log.Println("Funkcja findBestPositionAroundUnit !foundFreeSpot")

		return targetUnit.X, targetUnit.Y // Fallback
	}

	return uint8(bestX), uint8(bestY)
}

func (u *unit) findOptimalAttackTileAroundTree(treeX, treeY uint8, board *boardData) (uint8, uint8, bool) {
	var bestX, bestY uint8

	minDistance := math.MaxFloat64

	for column := int(treeX) - int(u.AttackRange); column <= int(treeX)+int(u.AttackRange); column++ {
		if column < 0 || column >= int(boardMaxX) {
			continue // kolumny poza planszą
		}

		for row := int(treeY) - int(u.AttackRange); row <= int(treeY)+int(u.AttackRange); row++ {
			if row < 0 || row >= int(boardMaxY) {
				continue // wiersz poza planszą
			}

			if !isWalkable(board, uint8(column), uint8(row)) {
				continue // kafelek nieprzechodni
			}

			if column+1 == int(treeX) && row == int(treeY) {
				continue // pomijamy miejsce na które spadnie drzewo
			}

			if column == int(treeX) && row == int(treeY) {
				continue // pomijamy samo drzewo
			}

			electedTile := &board.Tiles[uint8(column)][uint8(row)]

			if electedTile.Unit != nil && electedTile.Unit.ID != u.ID {
				continue // ktoś już tam stoi
			}

			dx := int(u.X) - column
			dy := int(u.Y) - row

			distance := math.Abs(float64(dx) + math.Abs(float64(dy)))

			if distance < minDistance {
				minDistance = distance
				bestX, bestY = uint8(column), uint8(row)
			}
		}
	}

	if minDistance == math.MaxFloat64 {
		return 0, 0, false
	}

	return bestX, bestY, true
}
