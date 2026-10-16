package main

import (
	"testing"
)

// UWAGA!!!! Próby wytworzone przez sztuczną inteligencję. Nie zostały jeszcze sprawdzone.
// Mogą zawierać błędy i niedociągnięcia! - 16.10.2026

// unit_commands_test.go — testy unit_commands.go

func TestUnit_setIdleWithReason(t *testing.T) {
	u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
	u.State = stateAttacking
	u.Command = cmdUAttack
	u.Target = tileTargetRef(9, 9)
	u.Approach = point{X: 8, Y: 8}
	u.Path = []point{{X: 6, Y: 6}}
	u.AllowFriendlyFire = true
	u.IsInQueue = true

	u.setIdleWithReason("test")

	assertState(t, u, stateIdle, cmdUIdle)

	if u.Target != (targetReference{}) || u.Approach != (point{}) {
		t.Errorf("target/approach not cleared: %+v / %+v", u.Target, u.Approach)
	}

	if len(u.Path) != 0 || u.AllowFriendlyFire || u.IsInQueue {
		t.Errorf("auxiliary fields not reset: %+v", u)
	}
}

func TestUnit_applyCommandState(t *testing.T) {
	tests := []struct {
		name          string
		cmd           commandType
		needsMovement bool
		wantState     unitState
		wantAnim      animationType
	}{
		{"interaction needing movement moves", cmdUAttack, true, stateMoving, animationWalk},
		{"attack in place", cmdUAttack, false, stateAttacking, animationFight},
		{"move", cmdUMove, false, stateMoving, animationWalk},
		{"move ignores needsMovement flag", cmdUMove, true, stateMoving, animationWalk},
		{"flee", cmdUFlee, false, stateMoving, animationWalk},
		{"cast in place", cmdUCastSpell, false, stateCastingSpell, animationFight},
		{"graze", cmdUGraze, false, stateGrazing, animationIdle},
		{"build in place", cmdUBuild, false, stateBuilding, animationIdle},
		{"repair in place", cmdURepair, false, stateRepairing, animationIdle},
		{"milking", cmdBMilking, false, stateMilking, animationIdle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{State: stateIdle, Command: tt.cmd, AnimationCounter: 5}
			u.applyCommandState(tt.cmd, tt.needsMovement)

			if u.State != tt.wantState {
				t.Errorf("State = %d, want %d", u.State, tt.wantState)
			}

			if u.AnimationType != tt.wantAnim {
				t.Errorf("AnimationType = %d, want %d", u.AnimationType, tt.wantAnim)
			}
		})
	}

	t.Run("stop resets to idle command", func(t *testing.T) {
		u := &unit{State: stateMoving, Command: cmdUStop}
		u.applyCommandState(cmdUStop, false)
		assertState(t, u, stateIdle, cmdUIdle)
	})

	t.Run("unknown command panics", func(t *testing.T) {
		mustPanic(t, func() { (&unit{}).applyCommandState(cmdUGoto, false) })
	})
}

func TestUnit_calculateDistanceToTarget(t *testing.T) {
	u := &unit{X: 10, Y: 10}
	bld := newTestBuilding(1, buildingBarn, point{X: 20, Y: 8}, testAIPlayer)
	pal := newTestBuilding(2, buildingPalisade, point{X: 13, Y: 13}, testAIPlayer)

	tests := []struct {
		name   string
		target *combatTarget
		want   uint8
	}{
		{"unit chebyshev", &combatTarget{Unit: &unit{X: 13, Y: 12}}, 3},
		{"unit same tile", &combatTarget{Unit: &unit{X: 10, Y: 10}}, 0},
		{"tile chebyshev", &combatTarget{Tile: &tile{X: 8, Y: 15}}, 5},
		{"3x3 building measures to nearest edge", &combatTarget{Building: bld}, 10},
		{"palisade single tile", &combatTarget{Building: pal}, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := u.calculateDistanceToTarget(tt.target); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestCanAttackBuilding(t *testing.T) {
	pal := newTestBuilding(1, buildingPalisade, point{X: 1, Y: 1}, testAIPlayer)
	barn := newTestBuilding(2, buildingBarn, point{X: 5, Y: 5}, testAIPlayer)
	gone := newTestBuilding(3, buildingBarn, point{X: 9, Y: 9}, testAIPlayer)
	gone.Exists = false

	tests := []struct {
		name string
		ut   unitType
		bld  *building
		want bool
	}{
		{"destroyed building", unitSwordsman, gone, false},
		{"mage never attacks buildings", unitMage, barn, false},
		{"archer cannot break palisade", unitArcher, pal, false},
		{"axeman breaks palisade", unitAxeman, pal, true},
		{"priest breaks palisade", unitPriest, pal, true},
		{"swordsman attacks barn", unitSwordsman, barn, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canAttackBuilding(tt.ut, tt.bld); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_canAttack(t *testing.T) {
	self := newTestUnit(1, unitAxeman, 5, 5, testHumanPlayer)
	enemy := newTestUnit(2, unitArcher, 6, 5, testAIPlayer)
	dead := newTestUnit(3, unitArcher, 7, 5, testAIPlayer)
	dead.Exists = false

	dryTree := &tile{}
	makeStandingTree(dryTree, true)
	greenTree := &tile{}
	makeStandingTree(greenTree, false)

	tests := []struct {
		name   string
		target *combatTarget
		want   bool
	}{
		{"enemy unit", &combatTarget{Unit: enemy}, true},
		{"self", &combatTarget{Unit: self}, false},
		{"dead unit", &combatTarget{Unit: dead}, false},
		{"dry tree by axeman", &combatTarget{Tile: dryTree}, true},
		{"green tree by axeman", &combatTarget{Tile: greenTree}, false},
		{"plain tile", &combatTarget{Tile: &tile{}}, false},
		{"barn", &combatTarget{Building: newTestBuilding(9, buildingBarn, point{X: 1, Y: 1}, testAIPlayer)}, true},
		{"empty target", &combatTarget{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := self.canAttack(tt.target); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_validateAttackTarget(t *testing.T) {
	attacker := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
	axeman := newTestUnit(2, unitAxeman, 5, 5, testHumanPlayer)

	friendly := newTestUnit(3, unitArcher, 6, 5, testHumanPlayer)
	enemy := newTestUnit(4, unitArcher, 6, 5, testAIPlayer)
	corpse := newTestUnit(5, unitArcher, 6, 5, testAIPlayer)
	corpse.HP = 0

	enemyPal := newTestBuilding(1, buildingPalisade, point{X: 1, Y: 1}, testAIPlayer)
	enemyBarn := newTestBuilding(2, buildingBarn, point{X: 10, Y: 10}, testAIPlayer)
	ruin := newTestBuilding(3, buildingBarn, point{X: 20, Y: 20}, testAIPlayer)
	ruin.HP = 0
	ownBarn := newTestBuilding(4, buildingBarn, point{X: 30, Y: 30}, testHumanPlayer)

	standing := &tile{}
	makeStandingTree(standing, true)
	stump := &tile{treeState: treeFell}

	tests := []struct {
		name         string
		attacker     *unit
		friendlyFire bool
		target       *combatTarget
		wantErr      bool
	}{
		{"enemy unit ok", attacker, false, &combatTarget{Unit: enemy}, false},
		{"friendly unit forbidden", attacker, false, &combatTarget{Unit: friendly}, true},
		{"friendly unit allowed with friendly fire", attacker, true, &combatTarget{Unit: friendly}, false},
		{"dead unit", attacker, false, &combatTarget{Unit: corpse}, true},
		{"enemy barn ok", attacker, false, &combatTarget{Building: enemyBarn}, false},
		{"own barn forbidden", attacker, false, &combatTarget{Building: ownBarn}, true},
		{"palisade by archer forbidden", attacker, false, &combatTarget{Building: enemyPal}, true},
		{"palisade by axeman ok", axeman, false, &combatTarget{Building: enemyPal}, false},
		{"ruined building", attacker, false, &combatTarget{Building: ruin}, true},
		{"standing tree ok", attacker, false, &combatTarget{Tile: standing}, false},
		{"fallen tree", attacker, false, &combatTarget{Tile: stump}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.attacker.AllowFriendlyFire = tt.friendlyFire
			err := tt.attacker.validateAttackTarget(tt.target)

			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnit_canAttackTarget(t *testing.T) {
	u := &unit{X: 5, Y: 5, AttackRange: 3}

	if !u.canAttackTarget(&combatTarget{Unit: &unit{X: 8, Y: 5}}) {
		t.Error("target at range 3 should be attackable")
	}

	if u.canAttackTarget(&combatTarget{Unit: &unit{X: 9, Y: 5}}) {
		t.Error("target at range 4 should be out of reach")
	}
}

func TestUnit_getRangedTargetCoords(t *testing.T) {
	u := &unit{X: 5, Y: 5}
	bld := newTestBuilding(1, buildingBarn, point{X: 10, Y: 10}, testAIPlayer)

	tests := []struct {
		name   string
		target *combatTarget
		want   point
		wantOK bool
	}{
		{"unit", &combatTarget{Unit: &unit{X: 7, Y: 8}}, point{X: 7, Y: 8}, true},
		{"tile", &combatTarget{Tile: &tile{X: 2, Y: 3}}, point{X: 2, Y: 3}, true},
		{"building closest tile", &combatTarget{Building: bld}, point{X: 10, Y: 10}, true},
		{"empty", &combatTarget{}, point{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := u.getRangedTargetCoords(tt.target)
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}

			if ok && *got != tt.want {
				t.Errorf("got %+v, want %+v", *got, tt.want)
			}
		})
	}
}

func TestUnit_handleTargetPostAttack(t *testing.T) {
	alive := newTestUnit(2, unitArcher, 6, 5, testAIPlayer)
	dead := newTestUnit(3, unitArcher, 6, 5, testAIPlayer)
	dead.HP = 0
	barn := newTestBuilding(1, buildingBarn, point{X: 10, Y: 10}, testAIPlayer)
	ruin := newTestBuilding(2, buildingBarn, point{X: 20, Y: 20}, testAIPlayer)
	ruin.Exists = false
	bridge := newTestBuilding(3, buildingBridge, point{X: 1, Y: 1}, testAIPlayer)
	bridge.HP = 0

	tests := []struct {
		name      string
		target    *combatTarget
		wantState unitState
		wantCmd   commandType
	}{
		{"living unit keeps attacking", &combatTarget{Unit: alive}, stateAttacking, cmdUAttack},
		{"dead unit goes idle", &combatTarget{Unit: dead}, stateIdle, cmdUIdle},
		{"standing building keeps attacking", &combatTarget{Building: barn}, stateAttacking, cmdUAttack},
		{"destroyed building goes idle", &combatTarget{Building: ruin}, stateIdle, cmdUIdle},
		{"bridge is never destroyed", &combatTarget{Building: bridge}, stateAttacking, cmdUAttack},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
			u.Command = cmdUAttack
			u.AnimationFrame = 3
			u.handleTargetPostAttack(tt.target)
			assertState(t, u, tt.wantState, tt.wantCmd)

			if tt.wantState == stateAttacking && u.AnimationFrame != 0 {
				t.Errorf("AnimationFrame = %d, want 0", u.AnimationFrame)
			}
		})
	}
}

func TestUnit_canActOnIdle(t *testing.T) {
	for _, ut := range []unitType{unitCow, unitShepherd} {
		if (&unit{Type: ut}).canActOnIdle() {
			t.Errorf("unitType(%d) should not act on idle", ut)
		}
	}

	if !(&unit{Type: unitArcher}).canActOnIdle() {
		t.Error("archer should act on idle")
	}
}

func TestUnit_isReadyToAct(t *testing.T) {
	bState := newTestBattleState()
	enemy := newTestUnit(2, unitArcher, 8, 8, testAIPlayer)
	placeUnit(t, bState, enemy)

	tests := []struct {
		name   string
		state  unitState
		cmd    commandType
		target targetReference
		want   bool
	}{
		{"idle and no command", stateIdle, cmdUIdle, targetReference{}, true},
		{"attacking a live target", stateAttacking, cmdUAttack, unitTargetRef(enemy), false},
		{"attacking a vanished target", stateAttacking, cmdUAttack, targetReference{Kind: targetUnit, ID: 99}, true},
		{"moving", stateMoving, cmdUMove, targetReference{}, false},
		{"idle state but busy command", stateIdle, cmdUBuild, targetReference{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
			u.State = tt.state
			u.Command = tt.cmd
			u.Target = tt.target

			if got := u.isReadyToAct(bState); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_executeActionByDistance(t *testing.T) {
	tests := []struct {
		name      string
		distance  uint8
		wantState unitState
		wantCmd   commandType
	}{
		{"out of sight goes idle", 9, stateIdle, cmdUIdle},
		{"in attack range attacks", 1, stateAttacking, cmdUAttack},
		{"in sight but out of reach moves", 5, stateMoving, cmdUAttack},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
			u.Approach = point{X: 9, Y: 5}
			u.Path = []point{{X: 6, Y: 5}}
			u.executeActionByDistance(tt.distance)
			assertState(t, u, tt.wantState, tt.wantCmd)

			if tt.wantState == stateMoving && len(u.Path) != 0 {
				t.Error("path should be invalidated for recalculation")
			}
		})
	}
}

func TestUnit_findApproachTileForTarget(t *testing.T) {
	bState := newTestBattleState()
	attacker := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
	enemy := newTestUnit(2, unitArcher, 12, 5, testAIPlayer)
	placeUnit(t, bState, attacker)
	placeUnit(t, bState, enemy)

	t.Run("unit target yields adjacent free tile", func(t *testing.T) {
		approach, err := attacker.findApproachTileForTarget(unitTargetRef(enemy), bState)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if d := chebyshev(approach, point{X: enemy.X, Y: enemy.Y}); d != 1 {
			t.Errorf("approach %+v is %d tiles from the target, want 1", approach, d)
		}

		if bState.Board.Tiles[approach.X][approach.Y].Unit != nil {
			t.Error("approach tile must be empty")
		}
	})

	t.Run("building target yields tile next to footprint", func(t *testing.T) {
		barn := newTestBuilding(1, buildingBarn, point{X: 20, Y: 20}, testAIPlayer)
		placeBuilding(t, bState, barn)

		approach, err := attacker.findApproachTileForTarget(buildingTargetRef(barn), bState)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if d := getDistanceToUnit(barn.Type, barn.OccupiedTiles[0], approach.X, approach.Y); d != 1 {
			t.Errorf("approach %+v is %d tiles from the barn, want 1", approach, d)
		}
	})

	t.Run("vanished target is an error", func(t *testing.T) {
		if _, err := attacker.findApproachTileForTarget(targetReference{Kind: targetUnit, ID: 99}, bState); err == nil {
			t.Error("expected error for a missing target")
		}
	})

	t.Run("target surrounded by walls is unreachable", func(t *testing.T) {
		isolated := newTestUnit(3, unitArcher, 40, 40, testAIPlayer)
		placeUnit(t, bState, isolated)

		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				if dx != 0 || dy != 0 {
					bState.Board.Tiles[40+dx][40+dy].IsWalkable = false
				}
			}
		}

		if _, err := attacker.findApproachTileForTarget(unitTargetRef(isolated), bState); err == nil {
			t.Error("expected error for an unreachable target")
		}
	})
}

func TestUnit_addUnitCommand(t *testing.T) {
	newScenario := func(t *testing.T) (*battleState, *unit, *unit) {
		t.Helper()

		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		enemy := newTestUnit(2, unitArcher, 10, 5, testAIPlayer)
		placeUnit(t, bState, u)
		placeUnit(t, bState, enemy)

		return bState, u, enemy
	}

	t.Run("move goes straight to the destination", func(t *testing.T) {
		bState, u, _ := newScenario(t)
		u.addUnitCommand(&command{ActionType: cmdUMove, Target: tileTargetRef(20, 20)}, bState)

		assertState(t, u, stateMoving, cmdUMove)

		if u.Approach != (point{X: 20, Y: 20}) {
			t.Errorf("Approach = %+v, want (20,20)", u.Approach)
		}
	})

	t.Run("attack on distant enemy starts moving", func(t *testing.T) {
		bState, u, enemy := newScenario(t)
		u.addUnitCommand(&command{ActionType: cmdUAttack, Target: unitTargetRef(enemy), FriendlyFire: true}, bState)

		assertState(t, u, stateMoving, cmdUAttack)

		if !u.AllowFriendlyFire {
			t.Error("friendly fire flag should be copied from the command")
		}

		if d := chebyshev(u.Approach, point{X: enemy.X, Y: enemy.Y}); d != 1 {
			t.Errorf("approach %+v is %d tiles from the enemy, want 1", u.Approach, d)
		}
	})

	// @reminder: findApproachTileForTarget odrzuca kafelek, na którym stoi sama jednostka
	// (tile.Unit != nil), więc nawet przy celu obok rozkaz zaczyna się od stateMoving,
	// a właściwy atak wykonuje executeStandardUnitCommand. Test utrwala to zachowanie.
	t.Run("attack on adjacent enemy keeps attack command with adjacent approach", func(t *testing.T) {
		bState, u, enemy := newScenario(t)
		enemy.X = 6
		bState.Board.Tiles[10][5].Unit = nil
		bState.Board.Tiles[6][5].Unit = enemy

		u.addUnitCommand(&command{ActionType: cmdUAttack, Target: unitTargetRef(enemy)}, bState)

		if u.Command != cmdUAttack {
			t.Errorf("Command = %d, want attack", u.Command)
		}

		if d := chebyshev(u.Approach, point{X: enemy.X, Y: enemy.Y}); d != 1 {
			t.Errorf("approach %+v is %d tiles from the enemy, want 1", u.Approach, d)
		}
	})

	t.Run("attack on self is rejected", func(t *testing.T) {
		bState, u, _ := newScenario(t)
		u.addUnitCommand(&command{ActionType: cmdUAttack, Target: unitTargetRef(u)}, bState)

		assertState(t, u, stateIdle, cmdUIdle)
	})

	t.Run("attack on missing target is rejected", func(t *testing.T) {
		bState, u, _ := newScenario(t)
		u.addUnitCommand(&command{ActionType: cmdUAttack, Target: targetReference{Kind: targetUnit, ID: 99}}, bState)

		assertState(t, u, stateIdle, cmdUIdle)
	})

	t.Run("self-targeting spell casts in place", func(t *testing.T) {
		bState, _, _ := newScenario(t)
		priestess := newTestUnit(3, unitPriestess, 7, 7, testHumanPlayer)
		placeUnit(t, bState, priestess)

		priestess.addUnitCommand(&command{ActionType: cmdUCastSpell, Spell: spellMagicShield, Target: unitTargetRef(priestess)}, bState)

		assertState(t, priestess, stateCastingSpell, cmdUCastSpell)

		if priestess.CurrentSpell != spellMagicShield {
			t.Errorf("CurrentSpell = %d, want magic shield", priestess.CurrentSpell)
		}
	})
}

func TestUnit_canAttackTargetFromCurrentPosition(t *testing.T) {
	t.Run("unit in range", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
		u.AttackRange = 4
		enemy := newTestUnit(2, unitArcher, 9, 5, testAIPlayer)
		placeUnit(t, bState, u)
		placeUnit(t, bState, enemy)
		u.Target = unitTargetRef(enemy)

		if !u.canAttackTargetFromCurrentPosition(bState) {
			t.Error("enemy at range 4 should be attackable")
		}

		enemy.X = 10

		if u.canAttackTargetFromCurrentPosition(bState) {
			t.Error("enemy at range 5 should be out of reach")
		}
	})

	t.Run("burning tile forces relocation", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
		enemy := newTestUnit(2, unitArcher, 6, 5, testAIPlayer)
		placeUnit(t, bState, u)
		placeUnit(t, bState, enemy)
		u.Target = unitTargetRef(enemy)
		u.Path = []point{{X: 5, Y: 5}}
		bState.Board.Tiles[5][5].IsBurning = true

		if u.canAttackTargetFromCurrentPosition(bState) {
			t.Error("standing in fire must never allow attacking")
		}

		if len(u.Path) != 0 {
			t.Error("path should be invalidated")
		}

		if u.Approach == (point{X: 5, Y: 5}) {
			t.Error("approach should be moved off the burning tile")
		}
	})

	t.Run("missing target", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
		placeUnit(t, bState, u)
		u.Target = targetReference{Kind: targetUnit, ID: 99}

		if u.canAttackTargetFromCurrentPosition(bState) {
			t.Error("missing target cannot be attacked")
		}
	})

	t.Run("axeman avoids tile where the tree will fall", func(t *testing.T) {
		bState := newTestBattleState()
		makeStandingTree(&bState.Board.Tiles[10][10], true)
		axe := newTestUnit(1, unitAxeman, 9, 10, testHumanPlayer)
		placeUnit(t, bState, axe)
		axe.Command = cmdUAttack
		axe.Target = tileTargetRef(10, 10)

		if axe.canAttackTargetFromCurrentPosition(bState) {
			t.Error("axeman standing west of the tree must not cut it")
		}

		axe.X, axe.Y = 10, 11

		if !axe.canAttackTargetFromCurrentPosition(bState) {
			t.Error("axeman south of the tree should be able to cut it")
		}
	})
}

func TestUnit_fleeBurningTile(t *testing.T) {
	bState := newTestBattleState()
	u := newTestUnit(1, unitArcher, 20, 20, testHumanPlayer)
	placeUnit(t, bState, u)
	bState.Board.Tiles[20][20].IsBurning = true
	// Bezpośredni sąsiedzi też płoną, ucieczka musi sięgnąć dalej.
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			bState.Board.Tiles[20+dx][20+dy].IsBurning = true
		}
	}

	u.fleeBurningTile(bState)

	assertState(t, u, stateMoving, cmdUMove)

	if bState.Board.Tiles[u.Approach.X][u.Approach.Y].IsBurning {
		t.Errorf("approach %+v is still on fire", u.Approach)
	}

	if u.Approach == (point{X: 20, Y: 20}) {
		t.Error("approach must differ from the burning position")
	}

	if u.Delay != 1 {
		t.Errorf("Delay = %d, want 1", u.Delay)
	}

	t.Run("no safe tile leaves unit untouched", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitArcher, 20, 20, testHumanPlayer)
		placeUnit(t, bState, u)

		for x := 14; x <= 25; x++ {
			for y := 14; y <= 25; y++ {
				bState.Board.Tiles[x][y].IsBurning = true
			}
		}

		u.fleeBurningTile(bState)
		assertState(t, u, stateIdle, cmdUIdle)
	})
}

func TestUnit_actOnIdle(t *testing.T) {
	t.Run("cow never acts", func(t *testing.T) {
		bState := newTestBattleState()
		cow := newTestUnit(1, unitCow, 5, 5, testHumanPlayer)
		placeUnit(t, bState, cow)
		bState.Board.Tiles[5][5].IsBurning = true

		cow.actOnIdle(bState)
		assertState(t, cow, stateIdle, cmdUIdle)
	})

	t.Run("busy unit is left alone", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
		placeUnit(t, bState, u)
		u.State = stateMoving
		u.Command = cmdUMove
		bState.Board.Tiles[5][5].IsBurning = true

		u.actOnIdle(bState)
		assertState(t, u, stateMoving, cmdUMove)
	})

	t.Run("idle unit flees fire", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
		placeUnit(t, bState, u)
		bState.Board.Tiles[5][5].IsBurning = true

		u.actOnIdle(bState)
		assertState(t, u, stateMoving, cmdUMove)
	})

	t.Run("idle unit without enemies stays idle", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
		placeUnit(t, bState, u)

		u.actOnIdle(bState)
		assertState(t, u, stateIdle, cmdUIdle)
	})

	t.Run("idle human unit engages visible enemy", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
		u.AttackRange = 4
		enemy := newTestUnit(2, unitArcher, 8, 5, testAIPlayer)
		placeUnit(t, bState, u)
		placeUnit(t, bState, enemy)

		u.actOnIdle(bState)

		if u.Command != cmdUAttack || u.Target.Kind != targetUnit || u.Target.ID != uint(enemy.ID) {
			t.Errorf("unit should target the enemy: cmd=%d target=%+v", u.Command, u.Target)
		}
	})

	t.Run("idle human unit ignores enemy buildings", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		placeUnit(t, bState, u)
		placeBuilding(t, bState, newTestBuilding(1, buildingBarn, point{X: 7, Y: 5}, testAIPlayer))

		u.actOnIdle(bState)
		assertState(t, u, stateIdle, cmdUIdle)
	})

	t.Run("idle AI unit attacks enemy building", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testAIPlayer)
		placeUnit(t, bState, u)
		barn := newTestBuilding(1, buildingBarn, point{X: 7, Y: 5}, testHumanPlayer)
		placeBuilding(t, bState, barn)

		u.actOnIdle(bState)

		if u.Command != cmdUAttack || u.Target.Kind != targetBuilding || u.Target.ID != uint(barn.ID) {
			t.Errorf("AI unit should target the barn: cmd=%d target=%+v", u.Command, u.Target)
		}
	})
}

func TestUnit_handleTargetReached(t *testing.T) {
	t.Run("attack", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
		u.Command = cmdUAttack
		u.Path = []point{{X: 5, Y: 5}}

		u.handleTargetReached(bState)
		assertState(t, u, stateAttacking, cmdUAttack)

		if u.Delay != 1 || len(u.Path) != 0 {
			t.Errorf("Delay/path = %d/%d, want 1/0", u.Delay, len(u.Path))
		}
	})

	t.Run("cast spell", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriest, 5, 5, testHumanPlayer)
		u.Command = cmdUCastSpell
		u.AnimationFrame = 2

		u.handleTargetReached(bState)
		assertState(t, u, stateCastingSpell, cmdUCastSpell)

		if u.AnimationFrame != 0 {
			t.Errorf("AnimationFrame = %d, want 0", u.AnimationFrame)
		}
	})

	t.Run("build works on the site", func(t *testing.T) {
		bState := newTestBattleState()
		site := newTestBuilding(1, buildingBarn, point{X: 10, Y: 10}, testHumanPlayer)
		site.HP = 10
		site.IsUnderConstruction = true
		placeBuilding(t, bState, site)

		u := newTestUnit(1, unitAxeman, 9, 10, testHumanPlayer)
		u.Command = cmdUBuild
		u.Target = buildingTargetRef(site)

		u.handleTargetReached(bState)
		assertState(t, u, stateBuilding, cmdUBuild)

		if site.HP != 10+repairAmountPlayer {
			t.Errorf("site HP = %d, want %d", site.HP, 10+repairAmountPlayer)
		}
	})

	t.Run("move ends idle", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
		u.Command = cmdUMove
		u.State = stateMoving

		u.handleTargetReached(bState)
		assertState(t, u, stateIdle, cmdUIdle)
	})
}

func TestUnit_work(t *testing.T) {
	t.Run("AI repairs faster than the player", func(t *testing.T) {
		bState := newTestBattleState()
		bld := newTestBuilding(1, buildingBarn, point{X: 10, Y: 10}, testAIPlayer)
		bld.HP = 50
		placeBuilding(t, bState, bld)

		u := newTestUnit(1, unitAxeman, 9, 10, testAIPlayer)
		u.Command = cmdURepair
		u.Target = buildingTargetRef(bld)

		u.work(bState)
		assertState(t, u, stateRepairing, cmdURepair)

		if bld.HP != 50+repairAmountAI {
			t.Errorf("HP = %d, want %d", bld.HP, 50+repairAmountAI)
		}
	})

	t.Run("finished repair goes idle", func(t *testing.T) {
		bState := newTestBattleState()
		bld := newTestBuilding(1, buildingBarn, point{X: 10, Y: 10}, testHumanPlayer)
		placeBuilding(t, bState, bld)

		u := newTestUnit(1, unitAxeman, 9, 10, testHumanPlayer)
		u.Command = cmdURepair
		u.Target = buildingTargetRef(bld)

		u.work(bState)
		assertState(t, u, stateIdle, cmdUIdle)
	})

	t.Run("vanished target goes idle", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitAxeman, 9, 10, testHumanPlayer)
		u.Command = cmdUBuild
		u.Target = targetReference{Kind: targetBuilding, ID: 42}

		u.work(bState)
		assertState(t, u, stateIdle, cmdUIdle)
	})
}

func TestUnit_handleWorkCommand(t *testing.T) {
	t.Run("in range builds and sets delay", func(t *testing.T) {
		bState := newTestBattleState()
		site := newTestBuilding(1, buildingBarn, point{X: 10, Y: 10}, testHumanPlayer)
		site.HP = 10
		site.IsUnderConstruction = true
		placeBuilding(t, bState, site)

		u := newTestUnit(1, unitAxeman, 9, 10, testHumanPlayer)
		placeUnit(t, bState, u)
		u.Command = cmdUBuild
		u.Target = buildingTargetRef(site)
		u.Approach = point{X: 9, Y: 10}

		u.handleWorkCommand(0, bState)
		assertState(t, u, stateBuilding, cmdUBuild)

		if site.HP != 10+repairAmountPlayer || u.Delay != u.MaxDelay {
			t.Errorf("HP/Delay = %d/%d, want %d/%d", site.HP, u.Delay, 10+repairAmountPlayer, u.MaxDelay)
		}
	})

	t.Run("out of range keeps moving", func(t *testing.T) {
		bState := newTestBattleState()
		site := newTestBuilding(1, buildingBarn, point{X: 10, Y: 10}, testHumanPlayer)
		site.HP = 10
		site.IsUnderConstruction = true
		placeBuilding(t, bState, site)

		u := newTestUnit(1, unitAxeman, 2, 2, testHumanPlayer)
		placeUnit(t, bState, u)
		u.Command = cmdUBuild
		u.Target = buildingTargetRef(site)
		u.Approach = point{X: 9, Y: 10}
		u.Delay = 1

		u.handleWorkCommand(0, bState)

		if u.State != stateMoving || site.HP != 10 {
			t.Errorf("State/HP = %d/%d, want moving/10", u.State, site.HP)
		}
	})
}

func TestUnit_executeStandardUnitCommand(t *testing.T) {
	t.Run("adjacent melee attack lands", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		enemy := newTestUnit(2, unitArcher, 6, 5, testAIPlayer)
		placeUnit(t, bState, u)
		placeUnit(t, bState, enemy)
		u.Command = cmdUAttack
		u.Target = unitTargetRef(enemy)

		u.executeStandardUnitCommand(0, bState)

		wantHP := enemy.MaxHP - (u.Damage - enemy.Armor)
		if enemy.HP != wantHP {
			t.Errorf("enemy HP = %d, want %d", enemy.HP, wantHP)
		}

		assertState(t, u, stateAttacking, cmdUAttack)
	})

	t.Run("cooldown blocks the strike", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		enemy := newTestUnit(2, unitArcher, 6, 5, testAIPlayer)
		placeUnit(t, bState, u)
		placeUnit(t, bState, enemy)
		u.Command = cmdUAttack
		u.Target = unitTargetRef(enemy)
		u.AttackCooldown = 5

		u.executeStandardUnitCommand(0, bState)

		if enemy.HP != enemy.MaxHP {
			t.Errorf("enemy HP = %d, want untouched %d", enemy.HP, enemy.MaxHP)
		}
	})

	t.Run("distant target re-plans approach", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		enemy := newTestUnit(2, unitArcher, 15, 5, testAIPlayer)
		placeUnit(t, bState, u)
		placeUnit(t, bState, enemy)
		u.Command = cmdUAttack
		u.Target = unitTargetRef(enemy)
		u.Delay = 1

		u.executeStandardUnitCommand(0, bState)

		if u.State != stateMoving {
			t.Errorf("State = %d, want moving", u.State)
		}

		if d := chebyshev(u.Approach, point{X: enemy.X, Y: enemy.Y}); d != 1 {
			t.Errorf("approach %+v is %d tiles from enemy, want 1", u.Approach, d)
		}
	})

	t.Run("unknown command panics", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		u.Command = cmdUGoto

		mustPanic(t, func() { u.executeStandardUnitCommand(0, bState) })
	})
}

func TestUnit_findNearestPalisade(t *testing.T) {
	bState := newTestBattleState()
	breaker := newTestUnit(1, unitAxeman, 5, 5, testAIPlayer)
	placeUnit(t, bState, breaker)

	if got := breaker.findNearestPalisade(bState, 10); got != nil {
		t.Errorf("no palisades on the board, got %+v", got)
	}

	// Palisada bez strategicznego budynku obok nie jest „ważna”.
	// Scenariusz: palisada wolnostojąca nie chroni żadnego budynku, więc nie powinna być celem strategicznym dla SI.
	lonePal := newTestBuilding(1, buildingPalisade, point{X: 3, Y: 5}, testHumanPlayer)
	placeBuilding(t, bState, lonePal)

	if got := breaker.findNearestPalisade(bState, 10); got != nil {
		t.Errorf("palisade without a guarded building should be ignored, got %+v", got)
	}

	// Palisada otaczająca budynek gracza jest ważna.
	// Tworzymy pierścień palisad wokół stodoły gracza na współrzędnych (20,20).
	barn := newTestBuilding(2, buildingBarn, point{X: 20, Y: 20}, testHumanPlayer)
	placeBuilding(t, bState, barn)

	var ringID buildingID = 10

	for x := 19; x <= 23; x++ {
		for y := 19; y <= 23; y++ {
			if x == 19 || x == 23 || y == 19 || y == 23 {
				placeBuilding(t, bState, newTestBuilding(ringID, buildingPalisade, point{X: uint8(x), Y: uint8(y)}, testHumanPlayer))
				ringID++
			}
		}
	}

	// Scenariusz: jednostka SI (topornik niszczący palisady) podchodząca do bazy gracza z różnych stron świata
	// powinna zawsze poprawnie zidentyfikować najbliższy segment pierścienia palisad chroniącego budynek.
	// Test weryfikuje wszystkie 4 kierunki kardynalne (wschód, zachód, południe, północ),
	// aby jednoznacznie sprawdzić kontrakt bez sztucznego ograniczania pozycji jednostki.
	directions := []struct {
		name        string
		unitPos     point
		wantPalTile point
		desc        string
	}{
		{
			name:        "east of palisade ring",
			unitPos:     point{X: 25, Y: 21},
			wantPalTile: point{X: 23, Y: 21},
			desc:        "Jednostka na wschód od pierścienia palisad (X=25, Y=21) powinna wybrać wschodnią ścianę (X=23, Y=21)",
		},
		{
			name:        "west of palisade ring",
			unitPos:     point{X: 17, Y: 21},
			wantPalTile: point{X: 19, Y: 21},
			desc:        "Jednostka na zachód od pierścienia palisad (X=17, Y=21) powinna wybrać zachodnią ścianę (X=19, Y=21)",
		},
		{
			name:        "south of palisade ring",
			unitPos:     point{X: 21, Y: 25},
			wantPalTile: point{X: 21, Y: 23},
			desc:        "Jednostka na południe od pierścienia palisad (X=21, Y=25) powinna wybrać południową ścianę (X=21, Y=23)",
		},
		{
			name:        "north of palisade ring",
			unitPos:     point{X: 21, Y: 17},
			wantPalTile: point{X: 21, Y: 19},
			desc:        "Jednostka na północ od pierścienia palisad (X=21, Y=17) powinna wybrać północną ścianę (X=21, Y=19)",
		},
	}

	for _, tt := range directions {
		t.Run(tt.name, func(t *testing.T) {
			breaker.X, breaker.Y = tt.unitPos.X, tt.unitPos.Y

			got := breaker.findNearestPalisade(bState, 10)
			if got == nil || got.Type != buildingPalisade {
				t.Fatalf("%s: expected a palisade from the ring, got %+v", tt.desc, got)
			}

			if got.OccupiedTiles[0] != tt.wantPalTile {
				t.Errorf("%s: nearest palisade = %+v, want %+v", tt.desc, got.OccupiedTiles[0], tt.wantPalTile)
			}
		})
	}

	t.Run("human owner never considers palisades important", func(t *testing.T) {
		human := newTestUnit(3, unitAxeman, 25, 21, testHumanPlayer)

		if human.findNearestPalisade(bState, 10) != nil {
			t.Error("human units must not auto-target palisades")
		}
	})

	t.Run("non-breaker ignores palisades", func(t *testing.T) {
		archer := newTestUnit(4, unitArcher, 25, 21, testAIPlayer)

		if archer.findNearestPalisade(bState, 10) != nil {
			t.Error("archers cannot break palisades")
		}
	})
}
