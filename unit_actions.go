package main

import (
	"math"
)

// attack zadaje obrażenia celowi lub ustawia ruch w jego kierunku.
func (u *unit) attack(bState *battleState) {
	// 1. Bierzemy wskaźnik do celu
	target, err := bState.resolveTarget(u.Target)
	if err != nil {
		u.setIdleWithReason(err.Error())

		return
	}

	// 2. Sprawdzamy cel
	err = u.validateAttackTarget(target)
	if err != nil {
		u.setIdleWithReason(err.Error())

		return
	}

	// 3. Sprawdzanie przerwy
	if u.AttackCooldown > 0 {
		u.State = stateIdle
		u.AnimationType = "idle"

		// Obracamy jednostkę w stronę celu, żeby nie stała bokiem/tyłem
		u.faceTarget(target)

		// Ustawiamy minimalny Delay (1), aby w następnej klatce
		// znów wejść do tej funkcji i sprawdzić, czy cooldown już minął.
		u.Delay = 1

		return
	}

	// 4. Sprawdzanie zasięgu
	if u.canAttackTarget(target) {
		proj := u.performAttack(target, bState.HumanPlayerState.PlayerID, bState.AIEnemyState.PlayerID,
			bState)

		if proj != nil {
			bState.Projectiles = append(bState.Projectiles, proj)
		}

		return
	}

	// Jeśli cel oddalił się, gonimy go
	// ! trzeba posprzątać
	whereToGo, err := u.findApproachTileForTarget(u.Target, bState)
	if err != nil {
		return
	}

	u.invalidatePathForRecalculation()

	u.State = stateMoving
	u.AnimationType = "walk"
	u.Approach = whereToGo
}

func (u *unit) performAttack(target *combatTarget, hPID, aiPID PlayerID, bState *battleState) *projectile {
	if u.AttackRange > 1 {
		return u.performRangedAttack(target, u.Damage, hPID, aiPID)
	}

	u.performMeleeAttack(target, u.Damage, hPID, aiPID, bState)
	u.setAttackTimings()
	u.handleTargetPostAttack(target)

	return nil
}

func (u *unit) performRangedAttack(target *combatTarget, damage uint16, hPID, aiPID PlayerID) *projectile {
	targetCoords, ok := u.getRangedTargetCoords(target)
	if !ok {
		return nil
	}

	// Mechanizm odejmowania many za rzucenie magicznego pocisku
	if u.Type.isCaster() && !u.tryToDecreaseMana(u.getProjectileManaCost()) {
		return nil
	}

	projParams := projectileParameters{
		owner:        u.Owner,
		spawnX:       uint16(u.X),
		spawnY:       uint16(u.Y),
		targetX:      uint16(targetCoords.X),
		targetY:      uint16(targetCoords.Y),
		missileKind:  unitTypeToMissileType(u.Type),
		damage:       damage,
		friendlyfire: u.AllowFriendlyFire,
	}

	proj := spawnProjectile(projParams)

	// Za stworzenie jakiegokolwiek pocisku jest przyznawane doświadczenie.
	// Muszę dodać logikę rozdziało pomięcy celem jednostką a celem budynkiem.
	// u.gainExperience tutaj!
	handleGainExperience(u, target.Unit, hPID, aiPID)

	return proj
}

// @reminder: zdobywanie doświadczenia jest niezależne od wyniku ataku. Wykonał atak→gainExperience().
func (u *unit) performMeleeAttack(target *combatTarget, damage uint16, hPID, aiPID PlayerID, bState *battleState) {
	switch {
	case target.Unit != nil && target.Unit.Exists:
		target.Unit.takeDamage(damage)
		handleGainExperience(u, target.Unit, hPID, aiPID)
	case target.Building != nil && target.Building.Exists:
		target.Building.takeDamage(damage)
		handleGainExperience(u, nil, hPID, aiPID)
	case target.Tile.isTree():
		target.Tile.accumulateTreeCuts(bState)
	default:
		// @reminder: jak coś to można dodać logi tutaj
	}
}

func (u *unit) repair(targetBuilding *building, amount uint16) {
	if !targetBuilding.repair(amount) {
		u.setIdleWithReason("naprawa ukończona")
	}
}

func (u *unit) build(targetBuilding *building, amount uint16) {
	if !targetBuilding.build(amount) {
		u.setIdleWithReason("budowa ukończona")
	}
}

// @reminder: najprawdopodobniej objectResolver nie jest prawidłowo użyty i będzie wyrzucony.
// @reminder: wydaje mi się, że każde „idle” ustawiane wewnątrz tej metody jest zbyteczne.
// @todo: brakuje ustawienia uruchomienia (animacji) ataku przy rzucaniu czarów.
func (u *unit) castSpell(pathfindingBudget int, bState *battleState) {
	if u.AttackCooldown > 0 {
		u.State = stateIdle
		u.AnimationType = "idle"
		u.Delay = 1

		return
	}

	switch u.CurrentSpell {
	case spellMagicShield:
		u.castMagicShield()

	case spellMagicSight:
		u.castMagicSight(bState.Board)

	case spellMagicShower:
		if u.canAttackTargetFromCurrentPosition(bState) {
			u.State = stateCastingSpell
			u.AnimationType = "fight"
			u.clearPath()
			u.castMagicShower(bState)
		} else {
			u.State = stateMoving
			u.AnimationType = "walk"
			u.move(pathfindingBudget, bState)

			return
		}

	case spellNone:
	// Nigdy nie powinno się przytrafić

	default:
		// Nigdy nie powinno się przytafić
	}
}

// @reminder: przechodzenie w idle powinno być inaczej załatwione.
func (u *unit) castMagicShield() {
	// 0. Jeśli już jest magiczna tarcza, to nie można rzucić nowej
	if u.hasMagicShield {
		return
	}
	// 1. Odejmujemy potrzebną manę
	// @todo: sprawdź ile rzeczywiście kosztowała
	if u.tryToDecreaseMana(spellCostMagicShield) {
		// 2. aktywujemy efekt
		u.hasMagicShield = true
		// 3. Ustalamy ile ma trwać
		// @todo: sprawdź ile rzeczywiście trwała
		u.MagicShieldCooldown = spellDurationMagicShield
	}
}

var magicShowerProjsBuffer = make([]*projectile, 3)

// Metoda odpowiedzialna za gromobicie i deszcz ognia.
func (u *unit) magicShower(target *point, bState *battleState) []*projectile {
	// 0. Koszt czaru
	if u.Mana < spellBufferMagicShower || !u.tryToDecreaseMana(spellCostMagicShower) {
		return nil
	}

	// 1. Tworzymy czarodziejski deszcz
	damage, missileKind, ok := u.resolveMagicShowerStats()
	if !ok {
		return nil
	}

	// 2. Bezpiecznik pozycji początkowej tworzonych pocisków
	spawnY := target.Y

	if spawnY >= 4 {
		spawnY -= 4
	} else {
		spawnY = 0
	}

	projIndex := 0

	// 3. Tworzenie opadów
	for offset := -1; offset <= 1; offset++ {
		spawnX := int(target.X) + offset

		if spawnX < 0 || spawnX >= int(boardMaxX) {
			continue
		}

		projParameters := projectileParameters{
			owner:       u.Owner,
			spawnX:      uint16(spawnX),
			spawnY:      uint16(spawnY),
			targetY:     uint16(target.Y),
			missileKind: missileKind,
			damage:      damage,
		}

		proj := spawnMagicShowerProjectile(projParameters)

		if proj != nil {
			magicShowerProjsBuffer[projIndex] = proj
			projIndex++
		}

		// 4. Przyzanie doświadczenia za zaatakowanie
		targetedTile := &bState.Board.Tiles[spawnX][target.Y]

		switch {
		case targetedTile.Unit != nil && targetedTile.Unit.Exists:
			handleGainExperience(u, targetedTile.Unit, bState.HumanPlayerState.PlayerID, bState.AIEnemyState.PlayerID)
		case targetedTile.Building != nil && targetedTile.Building.Exists:
			handleGainExperience(u, nil, bState.HumanPlayerState.PlayerID, bState.AIEnemyState.PlayerID)
		default:
			// Nie przyznajemy nic doświadczenia za napaść na otoczenie
		}
	}

	return magicShowerProjsBuffer[:projIndex]
}

// @reminder: przechodzenie w idle powinno być inaczej załatwione.
func (u *unit) castMagicShower(bState *battleState) {
	if u.AttackCooldown > 0 {
		u.State = stateIdle
		u.AnimationType = "idle"
		u.Delay = 1

		return
	}

	target := &u.Target.Position
	newProjs := u.magicShower(target, bState)

	if len(newProjs) > 0 {
		u.setRangedTimings()
		u.setIdleWithReason("czar rzucony")

		bState.Projectiles = append(bState.Projectiles, newProjs...)
	} else {
		u.State = stateIdle
		u.AnimationType = "idle"
		u.Command = cmdUIdle
	}
}

// @reminder: napisane na czuja, pewnie nie działa.
// @todo: ogarnij
// @reminder: przechodzenie w idle powinno być inaczej załatwione.
func (u *unit) castMagicSight(board *boardData) {
	if u.Mana >= spellCostMagicSight {
		u.Mana -= spellCostMagicSight

		revealRadius := spellCostRangeMagicSight
		for i := u.X - revealRadius; i <= u.X+revealRadius; i++ {
			for j := u.Y - revealRadius; j <= u.Y+revealRadius; j++ {
				if i <= boardMaxX && j <= boardMaxY {
					// @todo: czemu 18?!
					if math.Abs(float64(u.X-i))+math.Abs(float64(u.Y-j)) < 18 {
						board.Tiles[i][j].Visibility = visibilityVisible
					}
				}
			}
		}
	}
}
