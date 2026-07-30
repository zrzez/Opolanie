package main

import (
	"container/heap"
)

// pathfinding.go
// @todo:
// 1) Wywal wyjątek dla krów po poprawieniu przechodności budynków
// 2) Zastanów się, czy posiadanie gotowej mapy „kosztowności terenu” będzie w praktyce usprawnieniem

var currentGeneration uint32

// pathNode jest węzłem tj. kafelkiem na planszy.
type pathNode struct {
	goCost        float32 // Cena dotarcia od początku do tego węzła
	heuristicCost float32 // Szacowana cena od tego węzła do celu
	finalCost     float32 // Suma goCost i heuristicCost
	parent        int32   // Rodzica do odtworzenia ścieżki
	heapIndex     int32
	generation    uint32 // Do śledzenia, czy węzeł jest z „tego szukania”, czy poprzednich
}

type nodeHeap struct {
	indices []int32
	nodes   *[]pathNode
}

func (h *nodeHeap) Len() int { return len(h.indices) }

func (h *nodeHeap) Less(i, j int) bool {
	ni, nj := h.indices[i], h.indices[j]

	return (*h.nodes)[ni].finalCost < (*h.nodes)[nj].finalCost
}

func (h *nodeHeap) Swap(i, j int) {
	h.indices[i], h.indices[j] = h.indices[j], h.indices[i]

	(*h.nodes)[h.indices[i]].heapIndex = int32(i)
	(*h.nodes)[h.indices[j]].heapIndex = int32(j)
}

func (h *nodeHeap) Push(x any) {
	idx := x.(int32)
	n := len(h.indices)

	(*h.nodes)[idx].heapIndex = int32(n)
	h.indices = append(h.indices, idx)
}

func (h *nodeHeap) Pop() any {
	old := h.indices
	n := len(old)
	idx := old[n-1]
	(*h.nodes)[idx].heapIndex = -1
	h.indices = old[:n-1]

	return idx
}

const (
	maxPathfindingIterations = 10000
	maxNodes                 = 256 * 256
)

var (
	sharedNodePool    = make([]pathNode, maxNodes)
	sharedOpenIndices = make([]int32, 0, 512)
)

// odnajduje ścieżkę do celu używając algo A*.
// @todo: powinno przyjmować point zamiast uint8.
func findPath(board *boardData, mover *unit, startX, startY, endX, endY uint8, buffer []point) []point {
	currentGeneration++

	if currentGeneration == 0 {
		clear(sharedNodePool)

		currentGeneration = 1
	}

	startIndex := int32(startY)<<8 | int32(startX)

	sharedNodePool[startIndex] = pathNode{
		parent:     -1,
		generation: currentGeneration,
		heapIndex:  -1,
	}

	sharedOpenIndices = sharedOpenIndices[:0]

	open := &nodeHeap{
		indices: sharedOpenIndices,
		nodes:   &sharedNodePool,
	}

	heap.Init(open)
	heap.Push(open, startIndex)

	iterations := 0
	for open.Len() > 0 {
		iterations++
		if iterations >= maxPathfindingIterations {
			sharedOpenIndices = open.indices

			return nil
		}

		currentIndex := heap.Pop(open).(int32)
		currentX := uint8(currentIndex & 0xFF)
		currentY := uint8(currentIndex >> 8)

		if currentX == endX && currentY == endY {
			sharedOpenIndices = open.indices

			return reconstructPath(currentIndex, buffer)
		}

		currentNode := &sharedNodePool[currentIndex]

		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 && dy == 0 {
					continue
				}

				checkX, checkY := int(currentX)+dx, int(currentY)+dy

				if checkX < 0 || checkX >= int(boardMaxX) || checkY < 0 || checkY >= int(boardMaxY) {
					continue
				}

				if !isWalkableUnit(board, uint8(checkX), uint8(checkY), mover) {
					continue
				}

				index := int32(checkY)<<8 | int32(checkX)
				existingNode := &sharedNodePool[index]

				newGoCost := currentNode.goCost + calculateMoveCost(checkX, checkY, board)

				if existingNode.generation == currentGeneration {
					if newGoCost < existingNode.goCost {
						// Znaleźliśmy lepszą ścieżkę – aktualizujemy
						existingNode.parent = currentIndex
						existingNode.goCost = newGoCost
						existingNode.heuristicCost = calcHeuristic(checkX, checkY, int(endX), int(endY))
						existingNode.finalCost = newGoCost + existingNode.heuristicCost
						heap.Fix(open, int(existingNode.heapIndex))
					}
				} else {
					hCost := calcHeuristic(checkX, checkY, int(endX), int(endY))
					sharedNodePool[index] = pathNode{
						parent:        currentIndex,
						generation:    currentGeneration,
						goCost:        newGoCost,
						heuristicCost: hCost,
						finalCost:     newGoCost + hCost,
						heapIndex:     -1, // Push ustawi właściwy indeks
					}

					// @todo: tutaj spędzam większość pamięci. Może uda się coś wykombinować?
					heap.Push(open, index)
				}
			}
		}
	}

	sharedOpenIndices = open.indices

	return nil
}

// @todo: ogarnij czemu w ogóle potrzebujemy tej otuliny.
func isWalkable(board *boardData, x, y uint8) bool {
	return isWalkableUnit(board, x, y, nil)
}

// isWalkableUnit - Sprawdza czy dana jednostka może wejść na kafelek.
// Obsługuje wyjątek: Krowa wchodzi do swojej Obory (punkt dojenia).
func isWalkableUnit(board *boardData, x, y uint8, mover *unit) bool {
	currentTile := &board.Tiles[x][y]

	// 3. Flaga z loadera mapy (jeśli loader oznaczył coś jako nieprzechodnie ręcznie)
	if currentTile.IsWalkable {
		return true
	}

	// 1. Sprawdź czy to budynek
	// ! handlePalisadeDestruction gwarantuje, że zniszczona palisada jest przechodnia.
	// Oznacza to, eż mogę pozbyć się sprawdzenia „czy palisada w budowie”.
	// @reminder: Po ogarnięciu obory będę mógł w ogóle pozbyć się zewnętrznego wyrażenia ponieważ
	// każdy budynek prawidłowo ustawia currentTile.IsWalkable.
	if currentTile.Building != nil {
		// a. Krowa + obora (milking spot) - TYLKO jeden kafelek
		// @reminder: w pierwowzorze lewy-dolny kafelek był przechodni w każdym budynku.
		// Zgaduję, że to się wywali niedługo ponieważ lewy dolny róg budynków będzie
		// normalnie przechodni.
		if mover != nil && mover.Type == unitCow &&
			currentTile.Building.Type == buildingBarn &&
			currentTile.Building.Owner == mover.Owner {

			mx, my, ok := currentTile.Building.calculateMilkingSpot()
			if ok && x == mx && y == my {
				return !isWaterOrObstacle(currentTile.TextureID)
			}
		}

		// b. Każdy inny budynek blokuje, ukończony most nie jest budynkiem
		return false
	}

	return false
}

func isWaterOrObstacle(spriteID uint16) bool {
	// Woda
	if spriteID >= spriteWaterStart && spriteID <= spriteWaterEnd {
		return true
	}
	// Skały
	if spriteID >= spriteRockStart && spriteID <= spriteRockEnd {
		return true
	}
	// Dół drzewa (pnie)
	if spriteID >= spriteTreeStumpStart && spriteID <= spriteTreeStumpEnd {
		return true
	}

	return false
}

func calculateMoveCost(toX, toY int, board *boardData) float32 {
	cost := float32(1.0)
	currentTile := &board.Tiles[toX][toY]

	// Drogi ułatwiają ruch
	if currentTile.TextureID >= spriteRoadStart && currentTile.TextureID <= spriteRoadEnd {
		cost *= 0.5
	}

	// Jednostki nie są przeszkodą, ale zniechęcamy do próby wejścia w nich
	if currentTile.Unit != nil {
		cost *= 100
	}

	// Płonące kafelki powinny być unikane
	if currentTile.IsBurning {
		cost *= 90
	}

	// Jeśli kafelek jest nawiedzony, to unikamy
	if currentTile.GhostEffect {
		cost *= 100
	}

	return cost
}

func calcHeuristic(fromX, fromY, toX, toY int) float32 {
	dx := toX - fromX

	if dx < 0 {
		dx = -dx
	}

	dy := toY - fromY

	if dy < 0 {
		dy = -dy
	}

	if dx > dy {
		return float32(dx)
	}

	return float32(dy)
}

// W założeniu wykorzystujemy u.Path aby uniknać tworzenia wycinków.
func reconstructPath(endIndex int32, buffer []point) []point {
	// 0. Ustalamy jak długa jest droga
	length := 0

	for index := endIndex; index != -1; index = sharedNodePool[index].parent {
		length++
	}

	// 1. Oraz, czy zmeiści się w podręcznej liście
	if cap(buffer) < length {
		buffer = make([]point, length)
	}

	buffer = buffer[:length]

	i := length - 1

	// 2. A* zwraca drogę od celu do przemieszczającego się.
	// Dlatego to zmieniam, aby przemieszczająca się jednostka
	// miała gotową ścieżkę.
	for index := endIndex; index != -1; index = sharedNodePool[index].parent {
		buffer[i] = point{X: uint8(index & 0xFF), Y: uint8(index >> 8)}
		i--
	}

	return buffer
}
