package main

import (
	"testing"
)

// UWAGA!!!! Próby wytworzone przez sztuczną inteligencję. Nie zostały jeszcze sprawdzone.
// Mogą zawierać błędy i niedociągnięcia! - 16.10.2026

// units_test.go — testy units.go

func TestUnitType_isCaster_isMagical(t *testing.T) {
	tests := []struct {
		ut          unitType
		wantCaster  bool
		wantMagical bool
	}{
		{unitMage, true, true},
		{unitPriest, true, true},
		{unitPriestess, true, true},
		{unitUnknown, false, true},
		{unitArcher, false, false},
		{unitCow, false, false},
		{unitAxeman, false, false},
	}

	for _, tt := range tests {
		if got := tt.ut.isCaster(); got != tt.wantCaster {
			t.Errorf("unitType(%d).isCaster() = %v, want %v", tt.ut, got, tt.wantCaster)
		}

		if got := tt.ut.isMagical(); got != tt.wantMagical {
			t.Errorf("unitType(%d).isMagical() = %v, want %v", tt.ut, got, tt.wantMagical)
		}
	}
}

func TestUnit_increaseManaUnit(t *testing.T) {
	tests := []struct {
		name     string
		exists   bool
		mana     uint16
		maxMana  uint16
		amount   uint16
		wantMana uint16
	}{
		{"adds within cap", true, 10, 50, 20, 30},
		{"caps at max", true, 40, 50, 20, 50},
		{"zero amount is no-op", true, 10, 50, 0, 10},
		{"dead unit unchanged", false, 10, 50, 20, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{Exists: tt.exists, Mana: tt.mana, MaxMana: tt.maxMana}
			u.increaseManaUnit(tt.amount)

			if u.Mana != tt.wantMana {
				t.Errorf("Mana = %d, want %d", u.Mana, tt.wantMana)
			}
		})
	}
}

func TestUnit_tryToDecreaseMana(t *testing.T) {
	tests := []struct {
		name     string
		exists   bool
		mana     uint16
		amount   uint16
		wantOK   bool
		wantMana uint16
	}{
		{"enough mana", true, 30, 10, true, 20},
		{"exact mana", true, 10, 10, true, 0},
		{"not enough mana", true, 5, 10, false, 5},
		{"dead unit", false, 30, 10, false, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{Exists: tt.exists, Mana: tt.mana}

			if got := u.tryToDecreaseMana(tt.amount); got != tt.wantOK {
				t.Errorf("ok = %v, want %v", got, tt.wantOK)
			}

			if u.Mana != tt.wantMana {
				t.Errorf("Mana = %d, want %d", u.Mana, tt.wantMana)
			}
		})
	}
}

func TestUnit_increaseHPUnit(t *testing.T) {
	tests := []struct {
		name   string
		exists bool
		hp     uint16
		amount uint16
		wantHP uint16
	}{
		{"heals within cap", true, 10, 5, 15},
		{"caps at max", true, 18, 5, 20},
		{"dead unit unchanged", false, 10, 5, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{Exists: tt.exists, HP: tt.hp, MaxHP: 20}
			u.increaseHPUnit(tt.amount)

			if u.HP != tt.wantHP {
				t.Errorf("HP = %d, want %d", u.HP, tt.wantHP)
			}
		})
	}
}

func TestUnit_decreaseHPUnit(t *testing.T) {
	tests := []struct {
		name       string
		exists     bool
		hp         uint16
		amount     uint16
		wantHP     uint16
		wantExists bool
	}{
		{"partial damage", true, 10, 4, 6, true},
		{"exact kill", true, 10, 10, 0, false},
		{"dead unit unchanged", false, 10, 4, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{Exists: tt.exists, HP: tt.hp, MaxHP: 20}
			u.decreaseHPUnit(tt.amount)

			if u.HP != tt.wantHP || u.Exists != tt.wantExists {
				t.Errorf("HP/Exists = %d/%v, want %d/%v", u.HP, u.Exists, tt.wantHP, tt.wantExists)
			}
		})
	}
}

func TestUnit_handleAttackCooldown(t *testing.T) {
	u := &unit{AttackCooldown: 2}

	u.handleAttackCooldown(0)
	u.handleAttackCooldown(1)
	u.handleAttackCooldown(2)

	if u.AttackCooldown != 0 {
		t.Errorf("AttackCooldown = %d, want 0 (never below zero)", u.AttackCooldown)
	}
}

func TestUnit_handleDelay(t *testing.T) {
	tests := []struct {
		name      string
		state     unitState
		delay     uint16
		wantBusy  bool
		wantDelay uint16
	}{
		{"idle resets delay to max", stateIdle, 1, false, 5},
		{"no delay left", stateMoving, 0, false, 0},
		{"delay ticks down but still busy", stateMoving, 2, true, 1},
		{"last tick frees the unit", stateMoving, 1, false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{State: tt.state, Delay: tt.delay, MaxDelay: 5}

			if got := u.handleDelay(0); got != tt.wantBusy {
				t.Errorf("busy = %v, want %v", got, tt.wantBusy)
			}

			if u.Delay != tt.wantDelay {
				t.Errorf("Delay = %d, want %d", u.Delay, tt.wantDelay)
			}
		})
	}
}

func TestUnit_determineActiveStateFromCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      commandType
		atTarget bool
		want     unitState
	}{
		{"not at target always moves", cmdUBuild, false, stateMoving},
		{"move", cmdUMove, true, stateMoving},
		{"attack", cmdUAttack, true, stateMoving},
		{"graze", cmdUGraze, true, stateGrazing},
		{"build", cmdUBuild, true, stateBuilding},
		{"repair", cmdURepair, true, stateRepairing},
		{"cast", cmdUCastSpell, true, stateCastingSpell},
		{"idle", cmdUIdle, true, stateIdle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{X: 3, Y: 3, Command: tt.cmd, Approach: point{X: 3, Y: 3}}
			if !tt.atTarget {
				u.Approach = point{X: 9, Y: 9}
			}

			if got := u.determineActiveStateFromCommand(); got != tt.want {
				t.Errorf("got %d, want %d", got, tt.want)
			}
		})
	}
}

func TestUnit_handleWaitingToActiveTransition(t *testing.T) {
	u := &unit{X: 1, Y: 1, Approach: point{X: 1, Y: 1}, State: stateWaiting, Command: cmdUGraze}
	u.handleWaitingToActiveTransition()

	if u.State != stateGrazing || u.AnimationType != animationFight {
		t.Errorf("State/Anim = %d/%d, want %d/%d", u.State, u.AnimationType, stateGrazing, animationFight)
	}

	moving := &unit{State: stateMoving, Command: cmdUGraze}
	moving.handleWaitingToActiveTransition()

	if moving.State != stateMoving {
		t.Errorf("non-waiting unit changed state to %d", moving.State)
	}
}

func TestUnit_ensureDelayIsSet(t *testing.T) {
	tests := []struct {
		name  string
		state unitState
		want  uint16
	}{
		{"moving with zero delay gets max", stateMoving, 7},
		{"attacking left alone", stateAttacking, 0},
		{"idle left alone", stateIdle, 0},
		{"waiting left alone", stateWaiting, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{State: tt.state, MaxDelay: 7}
			u.ensureDelayIsSet()

			if u.Delay != tt.want {
				t.Errorf("Delay = %d, want %d", u.Delay, tt.want)
			}
		})
	}
}

func TestUnitType_canDamagePalisades(t *testing.T) {
	for _, ut := range []unitType{unitAxeman, unitPriest} {
		if !ut.canDamagePalisades() {
			t.Errorf("unitType(%d) should damage palisades", ut)
		}
	}

	for _, ut := range []unitType{unitArcher, unitMage, unitSwordsman, unitCow} {
		if ut.canDamagePalisades() {
			t.Errorf("unitType(%d) should not damage palisades", ut)
		}
	}
}

func TestCommandType_isInteraction(t *testing.T) {
	interactions := []commandType{cmdUAttack, cmdUBuild, cmdURepair, cmdBPlaceConstruction, cmdUCastSpell}
	for _, ct := range interactions {
		if !ct.isInteraction() {
			t.Errorf("commandType(%d) should be an interaction", ct)
		}
	}

	nonInteractions := []commandType{cmdUIdle, cmdUMove, cmdUFlee, cmdUGraze, cmdUStop, cmdBMilking}
	for _, ct := range nonInteractions {
		if ct.isInteraction() {
			t.Errorf("commandType(%d) should not be an interaction", ct)
		}
	}
}

func TestUnit_canDamageTree(t *testing.T) {
	tests := []struct {
		name     string
		uType    unitType
		standing bool
		dry      bool
		burning  bool
		want     bool
	}{
		{"priest cuts any standing tree", unitPriest, true, false, false, true},
		{"axeman cuts dry tree", unitAxeman, true, true, false, true},
		{"axeman cannot cut green tree", unitAxeman, true, false, false, false},
		{"archer cannot cut", unitArcher, true, true, false, false},
		{"burning tree untouchable", unitPriest, true, false, true, false},
		{"no tree", unitPriest, false, false, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := &tile{IsBurning: tt.burning}
			if tt.standing {
				makeStandingTree(tl, tt.dry)
			}

			u := &unit{Type: tt.uType}
			if got := u.canDamageTree(tl); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnit_prepareForNewCommand(t *testing.T) {
	u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
	u.Path = []point{{X: 1, Y: 1}}
	u.PathIndex = 1
	u.Delay = 4

	target := tileTargetRef(9, 9)
	approach := point{X: 8, Y: 5}

	u.prepareForNewCommand(cmdUAttack, target, approach)

	if u.Command != cmdUAttack || u.Target != target || u.Approach != approach {
		t.Errorf("command/target/approach not stored: %+v", u)
	}

	if len(u.Path) != 0 || u.PathIndex != 0 {
		t.Errorf("path not cleared: len=%d idx=%d", len(u.Path), u.PathIndex)
	}

	if u.Delay != 0 {
		t.Errorf("Delay = %d, want 0", u.Delay)
	}

	if dir, ok := u.rememberedDirection(); !ok || dir != uint16(directionRight) {
		t.Errorf("unit should face right (approach is to the east), got %d ok=%v", dir, ok)
	}
}

func TestUnit_setIdle(t *testing.T) {
	u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)
	u.State = stateMoving
	u.Command = cmdUAttack
	u.Path = []point{{X: 1, Y: 1}}
	u.RetryAttempts = 3
	u.AllowFriendlyFire = true
	u.IsInQueue = true

	u.setIdle()

	assertState(t, u, stateIdle, cmdUIdle)

	if u.AnimationType != animationIdle {
		t.Errorf("AnimationType = %d, want idle", u.AnimationType)
	}

	if len(u.Path) != 0 || u.RetryAttempts != 0 || u.AllowFriendlyFire || u.IsInQueue {
		t.Errorf("setIdle did not reset auxiliary fields: %+v", u)
	}
}

func TestUnit_updateMovementAnimation(t *testing.T) {
	u := &unit{X: 6, Y: 5}

	for range animationSpeed {
		u.updateMovementAnimation(5, 5)
	}

	if u.AnimationFrame != 1 || u.AnimationCounter != 0 {
		t.Errorf("after %d ticks frame/counter = %d/%d, want 1/0", animationSpeed, u.AnimationFrame, u.AnimationCounter)
	}

	if dir, _ := u.rememberedDirection(); dir != uint16(directionRight) {
		t.Errorf("direction = %d, want right", dir)
	}

	u.AnimationFrame = walkAnimationFrames - 1
	u.AnimationCounter = animationSpeed - 1
	u.updateMovementAnimation(6, 5)

	if u.AnimationFrame != 0 {
		t.Errorf("frame should wrap to 0, got %d", u.AnimationFrame)
	}
}

func TestUnit_faceTarget(t *testing.T) {
	bld := newTestBuilding(1, buildingBarn, point{X: 10, Y: 10}, testAIPlayer)

	tests := []struct {
		name   string
		target *combatTarget
		want   uint8
	}{
		{"unit to the south-east", &combatTarget{Unit: &unit{X: 7, Y: 7}}, directionDownRight},
		{"tile to the north", &combatTarget{Tile: &tile{X: 5, Y: 1}}, directionUp},
		{"building closest tile", &combatTarget{Building: bld}, directionDownRight},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{X: 5, Y: 5}
			u.faceTarget(tt.target)

			if dir, ok := u.rememberedDirection(); !ok || dir != uint16(tt.want) {
				t.Errorf("direction = %d ok=%v, want %d", dir, ok, tt.want)
			}
		})
	}

	t.Run("empty target keeps facing", func(t *testing.T) {
		u := &unit{X: 5, Y: 5}
		u.faceTarget(&combatTarget{})

		if _, ok := u.rememberedDirection(); ok {
			t.Error("direction should stay unset for an empty target")
		}
	})

	t.Run("same tile keeps facing", func(t *testing.T) {
		u := &unit{X: 5, Y: 5}
		u.faceTarget(&combatTarget{Unit: &unit{X: 5, Y: 5}})

		if _, ok := u.rememberedDirection(); ok {
			t.Error("direction should stay unset when target shares the tile")
		}
	})
}

func TestUnit_getProjectileManaCost(t *testing.T) {
	tests := []struct {
		ut   unitType
		want uint16
	}{
		{unitPriestess, magicThunderManaCost},
		{unitPriest, magicFireManaCost},
		{unitMage, magicGhostManaCost},
		{unitArcher, 0},
		{unitCrossbowman, 0},
	}

	for _, tt := range tests {
		u := &unit{Type: tt.ut}
		if got := u.getProjectileManaCost(); got != tt.want {
			t.Errorf("unitType(%d) cost = %d, want %d", tt.ut, got, tt.want)
		}
	}
}

func TestUnit_setAttackTimings(t *testing.T) {
	tests := []struct {
		name         string
		ut           unitType
		attackRange  uint8
		wantDelay    uint16
		wantCooldown uint
	}{
		{"melee default", unitSwordsman, 1, 10, 25},
		{"melee bear", unitBear, 1, 10, 35},
		{"melee unknown", unitUnknown, 1, 10, 15},
		{"ranged default", unitArcher, 5, 12, 65},
		{"ranged crossbowman", unitCrossbowman, 6, 12, 70},
		{"ranged mage", unitMage, 5, 12, 100},
		{"ranged priest", unitPriest, 5, 12, 100},
		{"ranged priestess", unitPriestess, 5, 12, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{Type: tt.ut, AttackRange: tt.attackRange}
			u.setAttackTimings()

			if u.Delay != tt.wantDelay || u.AttackCooldown != tt.wantCooldown {
				t.Errorf("Delay/Cooldown = %d/%d, want %d/%d", u.Delay, u.AttackCooldown, tt.wantDelay, tt.wantCooldown)
			}
		})
	}
}

func TestUnit_resolveMagicShowerStats(t *testing.T) {
	tests := []struct {
		ut         unitType
		wantDamage uint16
		wantKind   uint8
		wantOK     bool
	}{
		{unitPriest, spellDamageFireShower, missileFireRain, true},
		{unitPriestess, spellDamageLightningShower, missileLightning, true},
		{unitMage, 0, 0, false},
		{unitArcher, 0, 0, false},
	}

	for _, tt := range tests {
		u := &unit{Type: tt.ut}
		damage, kind, ok := u.resolveMagicShowerStats()

		if damage != tt.wantDamage || kind != tt.wantKind || ok != tt.wantOK {
			t.Errorf("unitType(%d) = (%d,%d,%v), want (%d,%d,%v)", tt.ut, damage, kind, ok, tt.wantDamage, tt.wantKind, tt.wantOK)
		}
	}
}

func TestUnit_isLightType(t *testing.T) {
	heavy := []unitType{unitCommander, unitSpearman, unitSwordsman, unitBear, unitUnknown, unitCrossbowman}
	for _, ut := range heavy {
		if (&unit{Type: ut}).isLightType() {
			t.Errorf("unitType(%d) should be heavy", ut)
		}
	}

	light := []unitType{unitCow, unitAxeman, unitArcher, unitPriest, unitPriestess, unitMage, unitShepherd}
	for _, ut := range light {
		if !(&unit{Type: ut}).isLightType() {
			t.Errorf("unitType(%d) should be light", ut)
		}
	}
}

func TestUnit_takeDamage(t *testing.T) {
	tests := []struct {
		name       string
		hp         uint16
		armor      uint16
		shield     bool
		damage     uint16
		wantHP     uint16
		wantWounds int
	}{
		{"armor reduces damage", 20, 3, false, 10, 13, 1},
		{"damage equal to armor still deals 1", 20, 5, false, 5, 19, 1},
		{"damage below armor still deals 1", 20, 5, false, 2, 19, 1},
		{"lethal damage clamps at zero", 5, 0, false, 50, 0, 1},
		{"magic shield blocks everything", 20, 0, true, 10, 20, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{Exists: true, Type: unitArcher, HP: tt.hp, MaxHP: 20, Armor: tt.armor, hasMagicShield: tt.shield}
			u.takeDamage(tt.damage)

			if u.HP != tt.wantHP {
				t.Errorf("HP = %d, want %d", u.HP, tt.wantHP)
			}

			if len(u.Wounds) != tt.wantWounds {
				t.Errorf("wounds = %d, want %d", len(u.Wounds), tt.wantWounds)
			}

			if u.wasAttacked == tt.shield {
				t.Errorf("wasAttacked = %v, shield = %v", u.wasAttacked, tt.shield)
			}
		})
	}
}

func TestUnit_unregisterFromBuilding(t *testing.T) {
	bld := newTestBuilding(1, buildingBarracks, point{X: 1, Y: 1}, testHumanPlayer)
	bld.AssignedUnits = []unitID{7, 8}
	bld.Food = 2

	u := &unit{ID: 7, BelongsTo: bld}
	u.unregisterFromBuilding()

	if u.BelongsTo != nil {
		t.Error("BelongsTo should be cleared")
	}

	if len(bld.AssignedUnits) != 1 || bld.AssignedUnits[0] != 8 || bld.Food != 1 {
		t.Errorf("building not updated: %+v", bld)
	}

	// Brak budynku nie może wywołać paniki.
	(&unit{}).unregisterFromBuilding()
}

func TestUnit_getAnimationType(t *testing.T) {
	tests := []struct {
		state unitState
		want  animationType
	}{
		{stateIdle, animationIdle},
		{stateWaiting, animationIdle},
		{stateMilking, animationIdle},
		{stateMoving, animationWalk},
		{stateAttacking, animationFight},
		{stateRepairing, animationFight},
		{stateBuilding, animationFight},
		{stateCastingSpell, animationFight},
		{stateGrazing, animationFight},
		{unitState(99), animationIdle},
	}

	for _, tt := range tests {
		u := &unit{State: tt.state}
		u.setAnimationType()

		if u.AnimationType != tt.want {
			t.Errorf("state %d -> animation %d, want %d", tt.state, u.AnimationType, tt.want)
		}
	}
}

func TestUnit_updateArcherAttackAnimation(t *testing.T) {
	tests := []struct {
		delay    uint16
		wantStep uint8
	}{
		{12, 0},
		{5, 0},
		{4, 1},
		{3, 1},
		{2, 2},
		{0, 2},
	}

	for _, tt := range tests {
		u := &unit{Type: unitArcher, Delay: tt.delay, AnimTick: 9}
		u.updateArcherAttackAnimation()

		if u.AnimStep != tt.wantStep || u.AnimTick != 0 || u.AnimationType != animationFight {
			t.Errorf("delay %d -> step/tick/anim = %d/%d/%d, want %d/0/%d", tt.delay, u.AnimStep, u.AnimTick, u.AnimationType, tt.wantStep, animationFight)
		}
	}
}

func TestUnit_animationDirection(t *testing.T) {
	t.Run("explicit direction wins", func(t *testing.T) {
		u := &unit{}
		u.applyDirection(-1, 0)

		if got := u.animationDirection(); got != uint16(directionLeft) {
			t.Errorf("got %d, want left", got)
		}
	})

	t.Run("idle without direction is none", func(t *testing.T) {
		u := &unit{AnimationType: animationIdle}
		u.rememberDirection(uint16(directionUp))

		if got := u.animationDirection(); got != uint16(directionNone) {
			t.Errorf("got %d, want none", got)
		}
	})

	t.Run("walking without direction uses remembered", func(t *testing.T) {
		u := &unit{AnimationType: animationWalk}
		u.rememberDirection(uint16(directionUpLeft))

		if got := u.animationDirection(); got != uint16(directionUpLeft) {
			t.Errorf("got %d, want up-left", got)
		}
	})

	t.Run("fighting without any direction defaults to down", func(t *testing.T) {
		u := &unit{AnimationType: animationFight}

		if got := u.animationDirection(); got != uint16(directionDown) {
			t.Errorf("got %d, want down", got)
		}
	})
}

func TestUnit_rememberDirection_roundTrip(t *testing.T) {
	u := &unit{}

	if _, ok := u.rememberedDirection(); ok {
		t.Fatal("fresh unit should have no remembered direction")
	}

	for dir := uint16(directionUp); dir <= uint16(directionUpLeft); dir++ {
		u.rememberDirection(dir)

		if got, ok := u.rememberedDirection(); !ok || got != dir {
			t.Errorf("round trip for %d gave %d ok=%v", dir, got, ok)
		}
	}
}

func TestUnit_applyDirection_zeroIsIgnored(t *testing.T) {
	u := &unit{}
	u.applyDirection(1, 1)
	u.applyDirection(0, 0)

	if u.Direction.X != 1 || u.Direction.Y != 1 {
		t.Errorf("Direction = %+v, want (1,1)", u.Direction)
	}
}

func TestUnit_facePoint(t *testing.T) {
	tests := []struct {
		name string
		p    point
		want uint8
	}{
		{"far east normalises to one step", point{X: 60, Y: 10}, directionRight},
		{"far north-west", point{X: 0, Y: 0}, directionUpLeft},
		{"south", point{X: 10, Y: 30}, directionDown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &unit{X: 10, Y: 10}
			u.facePoint(tt.p)

			if dir, ok := u.rememberedDirection(); !ok || dir != uint16(tt.want) {
				t.Errorf("direction = %d ok=%v, want %d", dir, ok, tt.want)
			}
		})
	}

	t.Run("own tile leaves direction untouched", func(t *testing.T) {
		u := &unit{X: 10, Y: 10}
		u.facePoint(point{X: 10, Y: 10})

		if _, ok := u.rememberedDirection(); ok {
			t.Error("direction should stay unset")
		}
	})
}
