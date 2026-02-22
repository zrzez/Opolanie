package main

import (
	"container/heap"
	"fmt"
	"math"
)

// pathfinding.go

// pathNode jest węzłem tj. kafelkiem na planszy.
type pathNode struct {
	parent        *pathNode // Wskaźnik na rodzica do odtworzenia ścieżki
	X, Y          uint8     // Współrzędne na mapie
	goCost        float64   // Cena dotarcia od początku do tego węzła
	heuristicCost float64   // Szacowana cena od tego węzła do celu
	finalCost     float64   // Suma goCost i heuristicCost
	index         int       // Wskaźnik miejsca węzła w stosie
}

// pathNodeHeap wdraża heap.Interface i przechowuje węzły
type pathNodeHeap []*pathNode

func (h pathNodeHeap) Len() int           { return len(h) }
func (h pathNodeHeap) Less(i, j int) bool { return h[i].finalCost < h[j].finalCost }
func (h pathNodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *pathNodeHeap) Push(x any) {
	n := len(*h)
	node := x.(*pathNode)
	node.index = n
	*h = append(*h, node)
}

func (h *pathNodeHeap) Pop() any {
	old := *h
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*h = old[0 : n-1]

	return node
}

func (h *pathNodeHeap) update(node, parent *pathNode, goCost, heuristicCost float64) {
	node.goCost = goCost
	node.heuristicCost = heuristicCost
	node.finalCost = goCost + heuristicCost
	node.parent = parent
	heap.Fix(h, node.index)
}

// Ograniczenie prób.
const maxPathfindingIterations = 10000

// odnajduje ścieżkę do celu używając algo A*.
func findPath(bs *battleState, moverID uint, startX, startY, endX, endY uint8) []*pathNode {
	startNode := &pathNode{parent: nil, X: startX, Y: startY}

	// 1. Pobieramy jednostkę, żeby wiedzieć kim jesteśmy (dla wyjątków ruchu - np. Krowa->Obora)
	var mover *unit
	if moverID > 0 {
		mover, _ = getObjectByID(moverID, bs)
	}

	openHeap := &pathNodeHeap{}
	heap.Init(openHeap)
	heap.Push(openHeap, startNode)

	nodeMap := make(map[string]*pathNode)
	nodeMap[fmt.Sprintf("%d.%d", startX, startY)] = startNode

	iterations := 0

	for openHeap.Len() > 0 {
		iterations++
		if iterations >= maxPathfindingIterations {
			return nil
		}

		currentNode := heap.Pop(openHeap).(*pathNode)

		// Jeśli dotarliśmy do celu
		if currentNode.X == endX && currentNode.Y == endY {
			return reconstructPath(currentNode)
		}

		// Sprawdzamy sąsiadów
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}

				checkX, checkY := int(currentNode.X)+dx, int(currentNode.Y)+dy
				coordKey := fmt.Sprintf("%d.%d", checkX, checkY)

				// 2. KLUCZOWE: Używamy isWalkableUnit zamiast zwykłego isWalkable
				// To pozwala krowie wejść w ścianę obory, jeśli to jej cel.
				if !isWalkableUnit(bs, uint8(checkX), uint8(checkY), mover) {
					continue
				}

				newGoCost := currentNode.goCost + calculateMoveCost(currentNode, &pathNode{X: uint8(checkX), Y: uint8(checkY)}, bs, moverID)

				if existingNode, found := nodeMap[coordKey]; found {
					if newGoCost <= existingNode.goCost {
						openHeap.update(existingNode, currentNode, newGoCost, calcHeuristic(existingNode, &pathNode{X: endX, Y: endY}))
					}
				} else {
					newNode := &pathNode{
						parent:        currentNode,
						X:             uint8(checkX),
						Y:             uint8(checkY),
						goCost:        newGoCost,
						heuristicCost: calcHeuristic(&pathNode{X: uint8(checkX), Y: uint8(checkY)}, &pathNode{X: endX, Y: endY}),
					}
					newNode.finalCost = newNode.goCost + newNode.heuristicCost
					heap.Push(openHeap, newNode)
					nodeMap[coordKey] = newNode
				}
			}
		}
	}

	return nil
}

func isWalkable(bs *battleState, x, y uint8) bool {
	return isWalkableUnit(bs, x, y, nil)
}

// isWalkableUnit - Sprawdza czy dana jednostka może wejść na kafelek.
// Obsługuje wyjątek: Krowa wchodzi do swojej Obory (punkt dojenia).
func isWalkableUnit(bs *battleState, x, y uint8, mover *unit) bool {
	if x >= boardMaxX || y >= boardMaxY {
		return false
	}

	currentTile := &bs.Board.Tiles[x][y]

	// 1. Sprawdź czy to budynek
	if currentTile.Building != nil {
		// 1. Krowa + obora (milking spot) - TYLKO jeden kafelek
		if mover != nil && mover.Type == unitCow &&
			currentTile.Building.Type == buildingBarn &&
			currentTile.Building.Owner == mover.Owner {

			mx, my, ok := calculateMilkingSpot(currentTile.Building)
			if ok && x == mx && y == my {
				return !isWaterOrObstacle(currentTile.TextureID)
			}
		}

		// 2. Palisada w budowie - pozwala na naprawę (cała powierzchnia)
		if currentTile.Building.Type == buildingPalisade && currentTile.Building.IsUnderConstruction {
			return true
		}

		// 3. Ukończony most - przechodni na CAŁEJ powierzchni
		if currentTile.Building.Type == buildingBridge && !currentTile.Building.IsUnderConstruction {
			return !isWaterOrObstacle(currentTile.TextureID)
		}

		// 4. Każdy inny budynek blokuje
		return false
	}

	// 2. Standardowa weryfikacja terenu
	if isWaterOrObstacle(currentTile.TextureID) {
		return false
	}

	// 3. Flaga z loadera mapy (jeśli loader oznaczył coś jako nieprzechodnie ręcznie)
	if !currentTile.IsWalkable {
		return false
	}

	return true
}

func isWaterOrObstacle(id uint16) bool {
	// Woda
	if id >= spriteWaterStart && id <= spriteWaterEnd {
		return true
	}
	// Skały
	if id >= spriteRockStart && id <= spriteRockEnd {
		return true
	}
	// Dół drzewa (pnie)
	if id >= spriteTreeStumpStart && id <= spriteTreeStumpEnd {
		return true
	}

	// Gadżety blokujące (wybrane offsety)
	// To wymaga doprecyzowania, które konkretnie gadżety blokują,
	// ale na start zablokujmy cały zakres gadżetów, żeby sprawdzić czy działa,
	// albo użyj switcha na konkretne ID.
	// Stary kod blokował: 54, 58, 60... (wybiórczo).
	// Nowy start gadżetów to 363. Stare 54 to teraz 363.
	// Przykładowo: blokujemy wszystko co wygląda na duży kamień/płot.
	if id >= spriteGadgetStart && id <= spriteGadgetEnd {
		// Tu można dodać wyjątki dla grzybków (przechodnich)
		return true
	}

	// Zgliszcza i palisady (jeśli nie mają obiektu building)
	// UWAGA: Palisady (building) są obsługiwane przez tile.building != nil w isWalkableUnit.
	// Ale jeśli została sama tekstura (bez Logic building), to blokujemy.
	// if id >= spriteRuinStart && id <= spritePalisadeEnd {
	//	return true
	// }

	return false
}

func calculateMoveCost(from, to *pathNode, bs *battleState, moverID uint) float64 {
	cost := 1.0
	// koszt ruchu po przekątnej
	if from.X != to.X && from.Y != to.Y {
		cost = 1.414
	}

	terrainID := bs.Board.Tiles[to.X][to.Y].TextureID

	// Drogi ułatwiają ruch
	if terrainID >= spriteRoadStart && terrainID <= spriteRoadEnd {
		cost *= 0.5
	}

	// Unikanie tłoku (sztuczka A*):
	// Inne jednostki nie są ścianą (isWalkableUnit puszcza), ale są bardzo drogie.
	tile := &bs.Board.Tiles[to.X][to.Y]
	if tile.Unit != nil && uint(tile.Unit.Owner) == moverID {
		cost *= 3
	}

	return cost
}

func calcHeuristic(from, to *pathNode) float64 {
	dx := float64(to.X - from.X)
	dy := float64(to.Y - from.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

func reconstructPath(node *pathNode) []*pathNode {
	var path []*pathNode
	current := node
	for current != nil {
		path = append(path, current)
		current = current.parent
	}
	return path
}
