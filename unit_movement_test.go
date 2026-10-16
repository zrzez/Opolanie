package main

import (
	"testing"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// UWAGA!!!! Próby wytworzone przez sztuczną inteligencję. Nie zostały jeszcze sprawdzone.
// Mogą zawierać błędy i niedociągnięcia! - 16.10.2026

// unit_movement_test.go — testy jednostkowe podsystemu ruchu (unit_movement.go)

func TestUnit_move_BasicProgression(t *testing.T) {
	t.Run("steps along path and updates position, board, and animation", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		placeUnit(t, bState, u)

		u.Command = cmdUMove
		u.State = stateMoving
		u.Target = tileTargetRef(7, 5)
		u.Approach = point{X: 7, Y: 5}
		u.Path = []point{{X: 5, Y: 5}, {X: 6, Y: 5}, {X: 7, Y: 5}}
		u.PathIndex = 1

		u.move(0, bState)

		// Po pierwszym kroku jednostka powinna być na (6,5)
		if u.X != 6 || u.Y != 5 {
			t.Errorf("u position = (%d,%d), want (6,5)", u.X, u.Y)
		}
		if u.PathIndex != 2 {
			t.Errorf("u.PathIndex = %d, want 2", u.PathIndex)
		}
		if bState.Board.Tiles[5][5].Unit != nil {
			t.Error("old tile (5,5) should be cleared")
		}
		if bState.Board.Tiles[6][5].Unit != u {
			t.Error("new tile (6,5) should reference unit")
		}
		if u.Delay != u.MaxDelay {
			t.Errorf("u.Delay = %d, want %d", u.Delay, u.MaxDelay)
		}
		if u.Direction != (rl.Vector2{X: 1, Y: 0}) {
			t.Errorf("u.Direction = %+v, want {1, 0}", u.Direction)
		}

		// Drugi krok do celu (7,5)
		u.move(0, bState)

		if u.X != 7 || u.Y != 5 {
			t.Errorf("u position = (%d,%d), want (7,5)", u.X, u.Y)
		}
		if bState.Board.Tiles[6][5].Unit != nil {
			t.Error("old tile (6,5) should be cleared")
		}
		if bState.Board.Tiles[7][5].Unit != u {
			t.Error("new tile (7,5) should reference unit")
		}
	})

	t.Run("diagonal step updates direction correctly", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 10, 10, testHumanPlayer)
		placeUnit(t, bState, u)

		u.Command = cmdUMove
		u.State = stateMoving
		u.Target = tileTargetRef(11, 11)
		u.Approach = point{X: 11, Y: 11}
		u.Path = []point{{X: 10, Y: 10}, {X: 11, Y: 11}}
		u.PathIndex = 1

		u.move(0, bState)

		if u.X != 11 || u.Y != 11 {
			t.Errorf("u position = (%d,%d), want (11,11)", u.X, u.Y)
		}
		if u.Direction != (rl.Vector2{X: 1, Y: 1}) {
			t.Errorf("u.Direction = %+v, want {1, 1}", u.Direction)
		}
	})
}

func TestUnit_move_ArrivalAndTransitions(t *testing.T) {
	t.Run("arrival on move command transitions to idle", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 10, 10, testHumanPlayer)
		placeUnit(t, bState, u)

		u.Command = cmdUMove
		u.State = stateMoving
		u.Approach = point{X: 10, Y: 10}
		u.Target = tileTargetRef(10, 10)

		u.move(0, bState)

		assertState(t, u, stateIdle, cmdUIdle)
	})

	t.Run("early attack range check stops movement and transitions to attack", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		enemy := newTestUnit(2, unitArcher, 6, 5, testAIPlayer)
		placeUnit(t, bState, u)
		placeUnit(t, bState, enemy)

		u.Command = cmdUAttack
		u.State = stateMoving
		u.Target = unitTargetRef(enemy)
		u.Approach = point{X: 6, Y: 6}
		u.Path = []point{{X: 5, Y: 5}, {X: 5, Y: 6}, {X: 6, Y: 6}}
		u.PathIndex = 1

		// Jednostka Swordsman ma AttackRange=1 i stoi obok wroga (5,5 vs 6,5)
		u.move(0, bState)

		if u.State != stateAttacking {
			t.Errorf("State = %d, want %d (stateAttacking)", u.State, stateAttacking)
		}
		if len(u.Path) != 0 {
			t.Error("path should be cleared when attack range is reached")
		}
	})

	t.Run("vanished attack target stops movement and sets idle", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		placeUnit(t, bState, u)

		u.Command = cmdUAttack
		u.State = stateMoving
		u.Target = targetReference{Kind: targetUnit, ID: 999}
		u.Approach = point{X: 10, Y: 10}

		u.move(0, bState)

		assertState(t, u, stateIdle, cmdUIdle)
	})
}

func TestUnit_move_ObstacleAndRecalculation(t *testing.T) {
	t.Run("blocked waypoint waits and increments NoMoveTicks", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		obstacleUnit := newTestUnit(2, unitArcher, 6, 5, testHumanPlayer)
		placeUnit(t, bState, u)
		placeUnit(t, bState, obstacleUnit)

		u.Command = cmdUMove
		u.State = stateMoving
		u.Approach = point{X: 7, Y: 5}
		u.Target = tileTargetRef(7, 5)
		u.Path = []point{{X: 5, Y: 5}, {X: 6, Y: 5}, {X: 7, Y: 5}}
		u.PathIndex = 1

		u.move(0, bState)

		if u.X != 5 || u.Y != 5 {
			t.Errorf("unit should not have moved, got (%d,%d)", u.X, u.Y)
		}
		if u.State != stateWaiting {
			t.Errorf("State = %d, want %d (stateWaiting)", u.State, stateWaiting)
		}
		if u.NoMoveTicks != 1 {
			t.Errorf("NoMoveTicks = %d, want 1", u.NoMoveTicks)
		}
		if u.Delay != 1 {
			t.Errorf("Delay = %d, want 1", u.Delay)
		}
	})

	t.Run("exceeding maxNoMoveTicks aborts movement and sets idle", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		placeUnit(t, bState, u)

		u.Command = cmdUMove
		u.State = stateMoving
		u.Approach = point{X: 10, Y: 10}
		u.NoMoveTicks = maxNoMoveTicks + 1

		u.move(0, bState)

		assertState(t, u, stateIdle, cmdUIdle)
	})

	t.Run("exhausting pathfinding retries sets idle", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		placeUnit(t, bState, u)

		// Otaczamy jednostkę nieprzejezdnymi kafelkami, żeby findPath zawiódł
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx != 0 || dy != 0 {
					bState.Board.Tiles[5+dx][5+dy].IsWalkable = false
				}
			}
		}

		u.Command = cmdUMove
		u.State = stateMoving
		u.Approach = point{X: 20, Y: 20}
		u.Target = tileTargetRef(20, 20)
		u.RetryAttempts = maxPathfindingRetries - 1

		u.move(0, bState)

		assertState(t, u, stateIdle, cmdUIdle)
		if u.RetryAttempts != 0 {
			t.Errorf("RetryAttempts = %d, want 0 after reset", u.RetryAttempts)
		}
	})

	t.Run("moving target invalidates path and recalculates", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		enemy := newTestUnit(2, unitArcher, 15, 5, testAIPlayer)
		placeUnit(t, bState, u)
		placeUnit(t, bState, enemy)

		u.Command = cmdUAttack
		u.State = stateMoving
		u.Target = unitTargetRef(enemy)
		u.Approach = point{X: 14, Y: 5}
		u.Path = []point{{X: 5, Y: 5}, {X: 6, Y: 5}}
		u.PathIndex = 1

		// Przemieszczamy wroga na (20, 5) i aktualizujemy kafelek planszy
		bState.Board.Tiles[15][5].Unit = nil
		enemy.X, enemy.Y = 20, 5
		bState.Board.Tiles[20][5].Unit = enemy

		// hasValidPath zwróci false, bo Target.Position w u.Target nie zgadza się z pozycją wroga
		if u.hasValidPath(bState) {
			t.Error("hasValidPath should be false when target moved")
		}
	})
}

func TestUnit_move_PathfindingBudget(t *testing.T) {
	bState := newTestBattleState()
	u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
	placeUnit(t, bState, u)

	u.Command = cmdUMove
	u.State = stateMoving
	u.Approach = point{X: 20, Y: 20}
	u.Target = tileTargetRef(20, 20)
	u.Path = nil // brak ścieżki, wymaga przeliczenia

	// Przy przekroczonym budżecie jednostka powinna poczekać
	u.move(maxPathfindingBudget, bState)

	if u.State != stateWaiting {
		t.Errorf("State = %d, want %d (stateWaiting)", u.State, stateWaiting)
	}
}

func TestUnit_canMoveTo(t *testing.T) {
	board := newTestBoard()
	u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
	other := newTestUnit(2, unitArcher, 6, 5, testAIPlayer)
	board.Tiles[6][5].Unit = other

	tests := []struct {
		name string
		x, y uint8
		want bool
	}{
		{"empty walkable tile", 5, 6, true},
		{"same unit tile", 5, 5, true},
		{"occupied by another unit", 6, 5, false},
		{"out of bounds X", boardMaxX, 5, false},
		{"out of bounds Y", 5, boardMaxY, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.canMoveTo(tt.x, tt.y, board); got != tt.want {
				t.Errorf("canMoveTo(%d,%d) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}

	t.Run("unwalkable tile returns false", func(t *testing.T) {
		board.Tiles[7][7].IsWalkable = false
		if u.canMoveTo(7, 7, board) {
			t.Error("unwalkable tile should return false")
		}
	})
}

func TestBuilding_calculateMilkingSpot(t *testing.T) {
	t.Run("standard 3x3 barn returns top-left minX and maxY", func(t *testing.T) {
		barn := newTestBuilding(1, buildingBarn, point{X: 10, Y: 10}, testHumanPlayer)
		x, y, ok := barn.calculateMilkingSpot()

		if !ok {
			t.Fatal("expected ok to be true")
		}
		// buildingFootprint generuje kafelki od (10,10) do (12,12)
		if x != 10 || y != 12 {
			t.Errorf("milking spot = (%d,%d), want (10,12)", x, y)
		}
	})

	t.Run("empty occupied tiles returns false", func(t *testing.T) {
		bld := &building{}
		_, _, ok := bld.calculateMilkingSpot()
		if ok {
			t.Error("expected ok to be false for building without occupied tiles")
		}
	})
}
