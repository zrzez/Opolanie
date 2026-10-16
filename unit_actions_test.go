package main

import (
	"testing"
)

// UWAGA!!!! Próby wytworzone przez sztuczną inteligencję. Nie zostały jeszcze sprawdzone.
// Mogą zawierać błędy i niedociągnięcia! - 16.10.2026

// unit_actions_test.go — testy unit_actions.go

// newDuel ustawia napastnika i wroga w zadanej odległości w poziomie.
func newDuel(t *testing.T, attackerType unitType, attackRange, gap uint8) (*battleState, *unit, *unit) {
	t.Helper()

	bState := newTestBattleState()
	attacker := newTestUnit(1, attackerType, 10, 10, testHumanPlayer)
	attacker.AttackRange = attackRange
	attacker.Command = cmdUAttack
	enemy := newTestUnit(2, unitArcher, 10+gap, 10, testAIPlayer)
	placeUnit(t, bState, attacker)
	placeUnit(t, bState, enemy)
	attacker.Target = unitTargetRef(enemy)

	return bState, attacker, enemy
}

func TestUnit_attack(t *testing.T) {
	t.Run("melee strike damages target and sets timings", func(t *testing.T) {
		bState, u, enemy := newDuel(t, unitSwordsman, 1, 1)

		u.attack(bState)

		wantHP := enemy.MaxHP - (u.Damage - enemy.Armor)
		if enemy.HP != wantHP {
			t.Errorf("enemy HP = %d, want %d", enemy.HP, wantHP)
		}

		assertState(t, u, stateAttacking, cmdUAttack)

		if u.Delay != 10 || u.AttackCooldown != 25 {
			t.Errorf("Delay/Cooldown = %d/%d, want 10/25", u.Delay, u.AttackCooldown)
		}

		if len(bState.Projectiles) != 0 {
			t.Error("melee attack must not spawn projectiles")
		}
	})

	t.Run("killing blow ends in idle", func(t *testing.T) {
		bState, u, enemy := newDuel(t, unitSwordsman, 1, 1)
		enemy.HP = 1

		u.attack(bState)

		if enemy.HP != 0 {
			t.Errorf("enemy HP = %d, want 0", enemy.HP)
		}

		assertState(t, u, stateIdle, cmdUIdle)
	})

	t.Run("ranged attack spawns a projectile", func(t *testing.T) {
		bState, u, enemy := newDuel(t, unitArcher, 5, 4)

		u.attack(bState)

		if len(bState.Projectiles) != 1 {
			t.Fatalf("projectiles = %d, want 1", len(bState.Projectiles))
		}

		proj := bState.Projectiles[0]
		if proj.Kind != missileArrow || proj.Owner != u.Owner || proj.Damage != u.Damage {
			t.Errorf("projectile = %+v", proj)
		}

		if proj.TargetX != uint16(enemy.X) || proj.TargetY != uint16(enemy.Y) {
			t.Errorf("projectile target = (%d,%d), want (%d,%d)", proj.TargetX, proj.TargetY, enemy.X, enemy.Y)
		}

		if enemy.HP != enemy.MaxHP {
			t.Error("ranged attack must not deal instant damage")
		}

		if u.Delay != 12 || u.AttackCooldown != 65 {
			t.Errorf("Delay/Cooldown = %d/%d, want 12/65", u.Delay, u.AttackCooldown)
		}
	})

	t.Run("cooldown only turns towards target", func(t *testing.T) {
		bState, u, enemy := newDuel(t, unitSwordsman, 1, 1)
		u.AttackCooldown = 3

		u.attack(bState)

		if enemy.HP != enemy.MaxHP {
			t.Error("no damage expected during cooldown")
		}

		assertState(t, u, stateAttacking, cmdUAttack)

		if u.Delay != u.MaxDelay {
			t.Errorf("Delay = %d, want MaxDelay %d", u.Delay, u.MaxDelay)
		}

		if dir, _ := u.rememberedDirection(); dir != uint16(directionRight) {
			t.Errorf("unit should face the enemy to the east, got %d", dir)
		}
	})

	t.Run("out of range chases the target", func(t *testing.T) {
		bState, u, enemy := newDuel(t, unitSwordsman, 1, 6)
		u.Path = []point{{X: 10, Y: 10}}

		u.attack(bState)

		if u.State != stateMoving || len(u.Path) != 0 {
			t.Errorf("State/path = %d/%d, want moving/0", u.State, len(u.Path))
		}

		if d := chebyshev(u.Approach, point{X: enemy.X, Y: enemy.Y}); d != 1 {
			t.Errorf("approach %+v is %d tiles from enemy, want 1", u.Approach, d)
		}

		if enemy.HP != enemy.MaxHP {
			t.Error("no damage expected out of range")
		}
	})

	t.Run("missing target goes idle", func(t *testing.T) {
		bState, u, _ := newDuel(t, unitSwordsman, 1, 1)
		u.Target = targetReference{Kind: targetUnit, ID: 99}

		u.attack(bState)
		assertState(t, u, stateIdle, cmdUIdle)
	})

	t.Run("friendly target goes idle", func(t *testing.T) {
		bState, u, enemy := newDuel(t, unitSwordsman, 1, 1)
		enemy.Owner = u.Owner

		u.attack(bState)
		assertState(t, u, stateIdle, cmdUIdle)

		if enemy.HP != enemy.MaxHP {
			t.Error("friendly unit must not be damaged")
		}
	})
}

func TestUnit_performRangedAttack(t *testing.T) {
	tests := []struct {
		name         string
		ut           unitType
		mana         uint16
		friendlyFire bool
		wantProj     bool
		wantKind     uint8
		wantMana     uint16
	}{
		{"archer needs no mana", unitArcher, 0, false, true, missileArrow, 0},
		{"crossbowman bolt", unitCrossbowman, 0, true, true, missileBolt, 0},
		{"priestess pays thunder cost", unitPriestess, 30, false, true, missileLightning, 30 - magicThunderManaCost},
		{"priest pays fire cost", unitPriest, 30, false, true, missileFire, 30 - magicFireManaCost},
		{"mage pays ghost cost", unitMage, 30, false, true, missileGhost, 30 - magicGhostManaCost},
		{"caster without mana fires nothing", unitMage, 5, false, false, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUnit(1, tt.ut, 5, 5, testHumanPlayer)
			u.AttackRange = 5
			u.Mana = tt.mana
			u.MaxMana = 100
			u.AllowFriendlyFire = tt.friendlyFire
			enemy := newTestUnit(2, unitArcher, 9, 5, testAIPlayer)

			proj := u.performRangedAttack(&combatTarget{Unit: enemy}, u.Damage, testHumanPlayer, testAIPlayer)

			if (proj != nil) != tt.wantProj {
				t.Fatalf("projectile spawned = %v, want %v", proj != nil, tt.wantProj)
			}

			if u.Mana != tt.wantMana {
				t.Errorf("Mana = %d, want %d", u.Mana, tt.wantMana)
			}

			if proj == nil {
				return
			}

			if proj.Kind != tt.wantKind || proj.AllowFriendlyFire != tt.friendlyFire {
				t.Errorf("Kind/FriendlyFire = %d/%v, want %d/%v", proj.Kind, proj.AllowFriendlyFire, tt.wantKind, tt.friendlyFire)
			}

			if u.Experience == 0 {
				t.Error("shooting an enemy should grant experience")
			}
		})
	}

	t.Run("empty target yields nothing", func(t *testing.T) {
		u := newTestUnit(1, unitArcher, 5, 5, testHumanPlayer)

		if u.performRangedAttack(&combatTarget{}, 5, testHumanPlayer, testAIPlayer) != nil {
			t.Error("no projectile expected for an empty target")
		}
	})
}

func TestUnit_performMeleeAttack(t *testing.T) {
	t.Run("unit target takes damage and grants experience", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		enemy := newTestUnit(2, unitArcher, 6, 5, testAIPlayer)

		u.performMeleeAttack(&combatTarget{Unit: enemy}, 7, testHumanPlayer, testAIPlayer, bState)

		if enemy.HP != enemy.MaxHP-(7-enemy.Armor) {
			t.Errorf("enemy HP = %d", enemy.HP)
		}

		if u.Experience != 1 {
			t.Errorf("Experience = %d, want 1", u.Experience)
		}
	})

	t.Run("building target accumulates damage", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testAIPlayer)
		bld := newTestBuilding(1, buildingBarn, point{X: 7, Y: 5}, testHumanPlayer)

		u.performMeleeAttack(&combatTarget{Building: bld}, 7, testHumanPlayer, testAIPlayer, bState)

		if bld.AccumulatedDamage != 7 || bld.HP != bld.MaxHP {
			t.Errorf("AccumulatedDamage/HP = %d/%d, want 7/%d", bld.AccumulatedDamage, bld.HP, bld.MaxHP)
		}

		if u.Experience != 1 {
			t.Errorf("AI attacker should gain experience for buildings, got %d", u.Experience)
		}
	})

	t.Run("human attacker gets no experience for buildings", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		bld := newTestBuilding(1, buildingBarn, point{X: 7, Y: 5}, testAIPlayer)

		u.performMeleeAttack(&combatTarget{Building: bld}, 7, testHumanPlayer, testAIPlayer, bState)

		if u.Experience != 0 {
			t.Errorf("Experience = %d, want 0", u.Experience)
		}
	})

	t.Run("tree target accumulates cuts", func(t *testing.T) {
		bState := newTestBattleState()
		tree := &bState.Board.Tiles[7][5]
		makeStandingTree(tree, true)
		u := newTestUnit(1, unitAxeman, 6, 5, testHumanPlayer)

		u.performMeleeAttack(&combatTarget{Tile: tree}, 7, testHumanPlayer, testAIPlayer, bState)

		if tree.treeCuts != 1 {
			t.Errorf("treeCuts = %d, want 1", tree.treeCuts)
		}
	})

	t.Run("dead unit is ignored", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitSwordsman, 5, 5, testHumanPlayer)
		corpse := newTestUnit(2, unitArcher, 6, 5, testAIPlayer)
		corpse.Exists = false
		corpse.HP = 3

		u.performMeleeAttack(&combatTarget{Unit: corpse, Tile: &tile{}}, 7, testHumanPlayer, testAIPlayer, bState)

		if corpse.HP != 3 {
			t.Error("dead unit must not take damage")
		}
	})
}

func TestUnit_repairAndBuild(t *testing.T) {
	t.Run("repair in progress keeps working", func(t *testing.T) {
		u := newTestUnit(1, unitAxeman, 5, 5, testHumanPlayer)
		u.Command = cmdURepair
		u.State = stateRepairing
		bld := newTestBuilding(1, buildingBarn, point{X: 7, Y: 5}, testHumanPlayer)
		bld.HP = 90

		u.repair(bld, 5)

		if bld.HP != 95 {
			t.Errorf("HP = %d, want 95", bld.HP)
		}

		assertState(t, u, stateRepairing, cmdURepair)
	})

	t.Run("repair completes and unit idles", func(t *testing.T) {
		u := newTestUnit(1, unitAxeman, 5, 5, testHumanPlayer)
		u.Command = cmdURepair
		u.State = stateRepairing
		bld := newTestBuilding(1, buildingBarn, point{X: 7, Y: 5}, testHumanPlayer)

		u.repair(bld, 5)
		assertState(t, u, stateIdle, cmdUIdle)
	})

	t.Run("build progresses on construction site", func(t *testing.T) {
		u := newTestUnit(1, unitAxeman, 5, 5, testHumanPlayer)
		u.Command = cmdUBuild
		u.State = stateBuilding
		site := newTestBuilding(1, buildingBarn, point{X: 7, Y: 5}, testHumanPlayer)
		site.HP = 10
		site.IsUnderConstruction = true

		u.build(site, 5)

		if site.HP != 15 {
			t.Errorf("HP = %d, want 15", site.HP)
		}

		assertState(t, u, stateBuilding, cmdUBuild)
	})

	t.Run("build on finished building idles", func(t *testing.T) {
		u := newTestUnit(1, unitAxeman, 5, 5, testHumanPlayer)
		u.Command = cmdUBuild
		u.State = stateBuilding
		bld := newTestBuilding(1, buildingBarn, point{X: 7, Y: 5}, testHumanPlayer)

		u.build(bld, 5)
		assertState(t, u, stateIdle, cmdUIdle)
	})
}

func TestUnit_castMagicShield(t *testing.T) {
	tests := []struct {
		name         string
		mana         uint16
		alreadyUp    bool
		wantShield   bool
		wantMana     uint16
		wantCooldown uint
	}{
		{"enough mana activates shield", spellCostMagicShield + 10, false, true, 10, spellDurationMagicShield},
		{"not enough mana", spellCostMagicShield - 1, false, false, spellCostMagicShield - 1, 0},
		{"existing shield is not recast", spellCostMagicShield + 10, true, true, spellCostMagicShield + 10, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := newTestUnit(1, unitPriestess, 5, 5, testHumanPlayer)
			u.Mana = tt.mana
			u.MaxMana = 200
			u.hasMagicShield = tt.alreadyUp

			u.castMagicShield()

			if u.hasMagicShield != tt.wantShield || u.Mana != tt.wantMana || u.MagicShieldCooldown != tt.wantCooldown {
				t.Errorf("shield/mana/cooldown = %v/%d/%d, want %v/%d/%d",
					u.hasMagicShield, u.Mana, u.MagicShieldCooldown, tt.wantShield, tt.wantMana, tt.wantCooldown)
			}
		})
	}
}

func TestUnit_magicShower(t *testing.T) {
	t.Run("priest spawns three fire rain projectiles", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriest, 5, 5, testAIPlayer)
		u.Mana = 50
		u.MaxMana = 100
		victim := newTestUnit(2, unitArcher, 20, 20, testHumanPlayer)
		placeUnit(t, bState, victim)

		projs := u.magicShower(&point{X: 20, Y: 20}, bState)

		if len(projs) != 3 {
			t.Fatalf("projectiles = %d, want 3", len(projs))
		}

		for i, p := range projs {
			if p.Kind != missileFireRain || p.Damage != spellDamageFireShower || !p.AllowFriendlyFire {
				t.Errorf("proj[%d] = %+v", i, p)
			}

			if p.TargetY != 20 || p.TargetX != uint16(19+i) {
				t.Errorf("proj[%d] target = (%d,%d), want (%d,20)", i, p.TargetX, p.TargetY, 19+i)
			}
		}

		if u.Mana != 50-spellCostMagicShower {
			t.Errorf("Mana = %d, want %d", u.Mana, 50-spellCostMagicShower)
		}

		if u.Experience == 0 {
			t.Error("hitting a unit should grant experience")
		}
	})

	t.Run("priestess uses lightning", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriestess, 5, 5, testHumanPlayer)
		u.Mana = 50

		projs := u.magicShower(&point{X: 20, Y: 20}, bState)

		if len(projs) != 3 || projs[0].Kind != missileLightning || projs[0].Damage != spellDamageLightningShower {
			t.Errorf("projs = %+v", projs)
		}
	})

	t.Run("left board edge drops the off-board column", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriest, 5, 5, testHumanPlayer)
		u.Mana = 50

		projs := u.magicShower(&point{X: 0, Y: 20}, bState)

		if len(projs) != 2 {
			t.Errorf("projectiles = %d, want 2", len(projs))
		}
	})

	t.Run("right board edge drops the off-board column", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriest, 5, 5, testHumanPlayer)
		u.Mana = 50

		projs := u.magicShower(&point{X: boardMaxX - 1, Y: 20}, bState)

		if len(projs) != 2 {
			t.Errorf("projectiles = %d, want 2", len(projs))
		}
	})

	t.Run("not enough mana", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriest, 5, 5, testHumanPlayer)
		u.Mana = spellCostMagicShower - 1

		if projs := u.magicShower(&point{X: 20, Y: 20}, bState); projs != nil {
			t.Errorf("expected nil, got %d projectiles", len(projs))
		}

		if u.Mana != spellCostMagicShower-1 {
			t.Error("mana must not be spent")
		}
	})

	t.Run("non-caster has no shower", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitMage, 5, 5, testHumanPlayer)
		u.Mana = 50

		if projs := u.magicShower(&point{X: 20, Y: 20}, bState); projs != nil {
			t.Errorf("expected nil, got %d projectiles", len(projs))
		}
	})
}

func TestUnit_castMagicShower(t *testing.T) {
	t.Run("cooldown keeps unit in attack pose", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriest, 5, 5, testHumanPlayer)
		u.Command = cmdUCastSpell
		u.AttackCooldown = 4

		u.castMagicShower(bState)
		assertState(t, u, stateAttacking, cmdUCastSpell)

		if u.Delay != 1 || len(bState.Projectiles) != 0 {
			t.Errorf("Delay/projectiles = %d/%d, want 1/0", u.Delay, len(bState.Projectiles))
		}
	})

	t.Run("successful cast adds projectiles and idles", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriest, 5, 5, testHumanPlayer)
		u.Command = cmdUCastSpell
		u.Mana = 50
		u.Target = tileTargetRef(20, 20)

		u.castMagicShower(bState)

		if len(bState.Projectiles) != 3 {
			t.Errorf("projectiles = %d, want 3", len(bState.Projectiles))
		}

		assertState(t, u, stateIdle, cmdUIdle)

		if u.AttackCooldown != 100 {
			t.Errorf("AttackCooldown = %d, want 100 (ranged timings)", u.AttackCooldown)
		}
	})

	t.Run("failed cast idles without projectiles", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriest, 5, 5, testHumanPlayer)
		u.Command = cmdUCastSpell
		u.Mana = 0
		u.Target = tileTargetRef(20, 20)

		u.castMagicShower(bState)

		if len(bState.Projectiles) != 0 {
			t.Errorf("projectiles = %d, want 0", len(bState.Projectiles))
		}

		assertState(t, u, stateIdle, cmdUIdle)
	})
}

// TestUnit_castMagicSight — weryfikacja czaru MagicSight (Oko Proroka / Magiczne Spojrzenie).
//
// UWAGA / SUGESTIA STRUKTURALNA:
// Mechanika mgły wojny (Fog of War) oraz czar MagicSight nie są obecnie zaimplementowane
// w aktywnym silniku gry (brak integracji z mechaniką widoczności i renderowaniem).
// Poniższy test służy jako propozycja i wzorzec architektoniczny (sugestia), jak testy
// podsystemu widoczności kafelków powinny być ustrukturyzowane w przyszłości po wdrożeniu mgły wojny.
//
// Uwaga o kodzie produkcyjnym: castMagicSight w unit_actions.go zawiera błąd arytmetyki uint8
// przy obliczaniu odległości (math.Abs(float64(u.X-i)) przekręca się przy i > u.X),
// przez co ujawnia niespełnienie asercji widoczności kafelków wokół rzucającego czar.
func TestUnit_castMagicSight(t *testing.T) {
	t.Run("reveals tiles around caster", func(t *testing.T) {
		board := newTestBoard()
		u := newTestUnit(1, unitMage, 30, 30, testHumanPlayer)
		u.Mana = spellCostMagicSight + 5

		u.castMagicSight(board)

		if u.Mana != 5 {
			t.Errorf("Mana = %d, want 5", u.Mana)
		}

		if board.Tiles[30][30].Visibility != visibilityVisible || board.Tiles[35][35].Visibility != visibilityVisible {
			t.Error("tiles near the caster should be visible")
		}

		// Odległość Manhattan 14+14=28 ≥ 18 → poza zasięgiem odsłaniania.
		if board.Tiles[44][44].Visibility == visibilityVisible {
			t.Error("far diagonal tile should stay hidden")
		}

		if board.Tiles[10][10].Visibility == visibilityVisible {
			t.Error("tile outside the radius should stay hidden")
		}
	})

	t.Run("not enough mana changes nothing", func(t *testing.T) {
		board := newTestBoard()
		u := newTestUnit(1, unitMage, 30, 30, testHumanPlayer)
		u.Mana = spellCostMagicSight - 1

		u.castMagicSight(board)

		if u.Mana != spellCostMagicSight-1 || board.Tiles[30][30].Visibility == visibilityVisible {
			t.Error("spell must not fire without mana")
		}
	})
}

func TestUnit_castSpell(t *testing.T) {
	t.Run("cooldown blocks casting", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriestess, 5, 5, testHumanPlayer)
		u.Command = cmdUCastSpell
		u.CurrentSpell = spellMagicShield
		u.Mana = 100
		u.AttackCooldown = 2

		u.castSpell(0, bState)

		if u.hasMagicShield || u.State != stateAttacking || u.Delay != 1 {
			t.Errorf("shield/state/delay = %v/%d/%d", u.hasMagicShield, u.State, u.Delay)
		}
	})

	t.Run("magic shield dispatch", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriestess, 5, 5, testHumanPlayer)
		u.CurrentSpell = spellMagicShield
		u.Mana = 100

		u.castSpell(0, bState)

		if !u.hasMagicShield {
			t.Error("shield should be active")
		}
	})

	t.Run("magic sight dispatch", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitMage, 30, 30, testHumanPlayer)
		u.CurrentSpell = spellMagicSight
		u.Mana = 100

		u.castSpell(0, bState)

		if bState.Board.Tiles[30][30].Visibility != visibilityVisible {
			t.Error("sight should reveal the caster's tile")
		}
	})

	t.Run("magic shower in range fires", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriest, 5, 5, testHumanPlayer)
		u.AttackRange = 5
		u.Command = cmdUCastSpell
		u.CurrentSpell = spellMagicShower
		u.Mana = 50
		u.Target = tileTargetRef(8, 5)
		placeUnit(t, bState, u)

		u.castSpell(0, bState)

		if len(bState.Projectiles) != 3 {
			t.Errorf("projectiles = %d, want 3", len(bState.Projectiles))
		}
	})

	t.Run("magic shower out of range moves instead", func(t *testing.T) {
		bState := newTestBattleState()
		u := newTestUnit(1, unitPriest, 5, 5, testHumanPlayer)
		u.AttackRange = 5
		u.Command = cmdUCastSpell
		u.CurrentSpell = spellMagicShower
		u.Mana = 50
		u.Target = tileTargetRef(30, 5)
		u.Approach = point{X: 25, Y: 5}
		u.Delay = 1
		placeUnit(t, bState, u)

		u.castSpell(0, bState)

		if u.State != stateMoving || len(bState.Projectiles) != 0 {
			t.Errorf("State/projectiles = %d/%d, want moving/0", u.State, len(bState.Projectiles))
		}
	})
}
