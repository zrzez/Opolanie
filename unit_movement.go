package main

import "math"

func (u *unit) executeAStarMovement(resolver objectResolver, board *boardData, pathfindingBudget *int) {
	if !u.ensureValidPath(resolver, board, pathfindingBudget) {
		return
	}

	u.moveAlongPath(board)
}

func (u *unit) ensureValidPath(resolver objectResolver, board *boardData, pathfindingBudget *int) bool {
	if u.hasValidPath(resolver, board) {
		return true
	}

	if *pathfindingBudget >= maxPathfindingBudget {
		u.waitForPathfindingBudget()

		return false
	}

	*pathfindingBudget++

	return u.calculateNewPath(board)
}

func (u *unit) calculateNewPath(board *boardData) bool {
	newPath := findPath(board, u, u.X, u.Y, u.ApproachX, u.ApproachY)

	if newPath == nil {
		u.handlePathfindingFailure()

		return false
	}

	u.setPathAndState(newPath)

	return true
}

func (u *unit) findLocalDetour(board *boardData, blockedX, blockedY uint8) (uint8, uint8, bool) {
	bestX, bestY := 0, 0
	bestScore := math.MaxFloat64

	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}

			x, y := int(u.X)+dx, int(u.Y)+dy

			if x == int(blockedX) && y == int(blockedY) {
				continue
			}

			if !u.canMoveTo(uint8(x), uint8(y), board) {
				continue
			}

			score := u.calculateDetourScore(uint8(x), uint8(y))
			if score < bestScore {
				bestScore = score
				bestX, bestY = x, y
			}
		}
	}

	return uint8(bestX), uint8(bestY), true
}

func (u *unit) handleMovementBlocked(board *boardData, blockedX, blockedY uint8) {
	detourX, detourY, ok := u.findLocalDetour(board, blockedX, blockedY)

	if !ok {
		u.executeSuccessfulMove(detourX, detourY, board)
		u.invalidatePathForRecalculation()

		return
	}

	u.handlePersistentBlock()
}

func (u *unit) calculateDetourScore(x, y uint8) float64 {
	distToTarget := math.Abs(float64(u.TargetX-x)) + math.Abs(float64(u.TargetY-y))
	mainDirection := u.getMainDirection()
	directionBonus := 0.0

	dx := x - u.X
	dy := y - u.Y

	if (mainDirection.X > 0 && dx > 0) || (mainDirection.X < 0 && dx < 0) {
		directionBonus -= 0.5
	}
	if (mainDirection.Y > 0 && dy > 0) || (mainDirection.Y < 0 && dy < 0) {
		directionBonus -= 0.5
	}

	return distToTarget + directionBonus
}
