package main

import (
	"log"
	"math"
)

// attack zadaje obrażenia celowi lub ustawia ruch w jego kierunku.
func (u *unit) attack(resolver objectResolver, board *boardData, bState *battleState) {
	log.Printf("units.go attack weszliśmy do metody")

	target, err := u.validateAttackTarget(resolver, board)
	if err != nil {
		u.setIdleWithReason(err.Error())

		return
	}

	// Sprawdzanie przerwy
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

	if u.canAttackTarget(target, board) {
		log.Printf("units.go attack weszliśmy do canAttackTarget")

		u.performDirectAttack(target, bState)
	} else {
		log.Printf("units.go attack weszliśmy do !canAttackTarget.")

		// Jeśli cel oddalił się, gonimy go
		u.startMoveToAttack(bState)
	}
}

func (u *unit) performDirectAttack(target *combatTarget, bState *battleState) {
	if u.AttackRange > 1 {
		u.performRangedAttack(target, u.Damage, bState)
	} else {
		u.performMeleeAttack(target, u.Damage, bState)
	}

	u.setAttackTimings()
	u.handleTargetPostAttack(target.Unit, target.Building)
}

// @todo: do zmiany, jednostka nie może znać szczegółów wdrożenia pocisków.
func (u *unit) performRangedAttack(target *combatTarget, damage uint16, bState *battleState) {
	targetX, targetY, ok := u.getRangedTargetCoords(target)
	if !ok {
		log.Printf("UWAGA: jednostka %d: nie udało się określić koordynatów celu dla pocisku", u.ID)
	}

	// Mechanizm odejmowania many za rzucenie magicznego pocisku
	if u.Type.isCaster() && !u.tryToDecreaseMana(u.getProjectileManaCost()) {
		return
	}

	proj := &projectile{}
	proj.initProjectile(
		unitTypeToMissileType(u.Type),
		u.Owner,
		uint16(u.X), uint16(u.Y),
		uint16(targetX), uint16(targetY),
		damage,
	)

	proj.AllowFriendlyFire = u.AllowFriendlyFire

	bState.Projectiles = append(bState.Projectiles, proj)

	// Za stworzenie jakiegokolwiek pocisku jest przyznawane doświadczenie.
	// Muszę dodać logikę rozdziało pomięcy celem jednostką a celem budynkiem.
	// u.gainExperience tutaj!
	u.gainExperience(target.Unit, bState)

	log.Printf("jednostka %d wystrzeliła pocisk w (%d, %d) z obrażeniami %d", u.ID, targetX, targetY, damage)
}

// @reminder: zdobywanie doświadczenia jest niezależne od wyniku ataku. Wykonał atak→gainExperience()
func (u *unit) performMeleeAttack(target *combatTarget, damage uint16, bState *battleState) {
	switch {
	case target.Unit != nil && target.Unit.Exists:
		target.Unit.takeDamage(damage, bState)
		u.gainExperience(target.Unit, bState)
	case target.Building != nil && target.Building.Exists:
		target.Building.takeDamage(damage)
		u.gainExperience(nil, bState)
	case target.Tile != nil:
		target.Tile.accumulateTreeCuts(bState)
	default:
		log.Printf("UWAGA: jednostka %d: cel ataku wręcz już nie istnieje", u.ID)
	}
}

func (u *unit) repair(bState *battleState) {
	// Bezpiecznik, bo cel mógł się zmienić w międzyczasie
	_, targetBuilding := bState.getObjectByID(u.TargetID)
	validRepairTarget, _ := validateRepairContext(u, targetBuilding)

	if !validRepairTarget {
		u.setIdleWithReason("Nie można naprawić budynku")

		return
	}

	// Szukamy drogi do celu
	distance := getDistanceToUnit(targetBuilding.Type, targetBuilding.OccupiedTiles[0], u.X, u.Y)

	var amount uint16

	switch u.Owner {
	case bState.PlayerID:
		amount = repairAmountPlayer
	case bState.AIPlayerID:
		amount = repairAmountAI
	}

	if distance == 1 {
		targetBuilding.repair(amount)
	}
}

func (u *unit) build(bState *battleState) {
	// Bezpiecznik, bo budowa mogą się zmienić w międzyczasie
	_, targetBuilding := bState.getObjectByID(u.TargetID)
	validBuildTarget, _ := validateBuildingContext(u, targetBuilding)

	if !validBuildTarget {
		u.setIdleWithReason("Nie można budować")

		return
	}

	// Szukamy drogi do celu
	distance := getDistanceToUnit(targetBuilding.Type, targetBuilding.OccupiedTiles[0], u.X, u.Y)

	var amount uint16

	switch u.Owner {
	case bState.PlayerID:
		amount = repairAmountPlayer
	case bState.AIPlayerID:
		amount = repairAmountAI
	}

	if distance == 1 {
		targetBuilding.build(amount)
	}
}

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

	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdUIdle
}

// Metoda odpowiedzialna za gromobicie i deszcz ognia.
func (u *unit) magicShower(targetX, targetY uint8, bState *battleState) bool {
	// 0. Koszt czaru
	if u.Mana < spellBufferMagicShower || !u.tryToDecreaseMana(spellCostMagicShower) {
		log.Printf("INFO: Jednostka %d nie ma wystarczająco many na rzucenie czaru", u.ID)

		return false
	}

	// 1. Tworzymy czarodziejski deszcz
	u.createMagicShower(targetX, targetY, bState)

	// 2. Skończyliśmy czarowanie, stoimy bezczynnie
	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdUIdle

	log.Printf("INFO: Jednostka %d rzuciła czar na (%d, %d)", u.ID, targetX, targetY)

	return true
}

func (u *unit) castSpell(bState *battleState) {
	if u.AttackCooldown > 0 {
		u.State = stateIdle
		u.AnimationType = "idle"
		u.Delay = 1

		return
	}

	targetX, targetY := u.interactionTargetX, u.interactionTargetY

	// @todo: sprawdź, czy rzeczywiście tak jest, bo wydaje mi się, że
	//  zmieniłem to ostatnio. 18.07.2026
	if targetX == 0 && targetY == 0 {
		targetX, targetY = u.TargetX, u.TargetY
	}

	if u.magicShower(targetX, targetY, bState) {
		u.setRangedTimings()
		u.setIdleWithReason("czar rzucony")
	} else {
		u.State = stateIdle
		u.AnimationType = "idle"
		u.Command = cmdUIdle
	}
}

// @reminder: napisane na czuja, pewnie nie działa.
func (u *unit) castMagicSight(board *boardData) {
	if u.Mana >= spellCostMagicSight {
		u.Mana -= spellCostMagicSight
		log.Printf("Jednostka %d rzuca czar widzenia", u.ID)

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

	log.Printf("unit %d: Rzucono Czar Widoczności. Cała mapa odkryta.", u.ID)

	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdUIdle
}

// To już chyba nie powinna być metoda jednostki. Nic jej do szczegółów pocisków.
func (u *unit) spawnMagicShowerProjectiles(targetX, targetY uint8, missileKind uint8, damage uint16, bState *battleState) {
	// 0. Bezpiecznik pozycji początkowej tworzonych pocisków
	spawnY := targetY

	if spawnY >= 4 {
		spawnY -= 4
	} else {
		spawnY = 0
	}

	// 1. Tworzenie opadów
	for offset := -1; offset <= 1; offset++ {
		spawnX := int(targetX) + offset

		if spawnX < 0 || spawnX >= int(boardMaxX) {
			continue
		}

		proj := &projectile{}
		proj.initProjectile(
			missileKind,
			u.Owner,
			uint16(spawnX), uint16(spawnY),
			uint16(spawnX), uint16(targetY),
			damage,
		)

		if proj.Exists {
			bState.Projectiles = append(bState.Projectiles, proj)
		}

		// 2. Przyzanie doświadczenia za zaatakowanie
		currentTile := &bState.Board.Tiles[spawnX][targetY]

		switch {
		case currentTile.Unit != nil && currentTile.Unit.Exists:
			u.gainExperience(currentTile.Unit, bState)
		case currentTile.Building != nil && currentTile.Building.Exists:
			u.gainExperience(nil, bState)
		default:
			// Nie przyznajemy nic doświadczenia za napaść na otoczenie
		}
	}
}

// To już chyba nie powinna być metoda jednostki. Nic jej do szczegółów pocisków.
func (u *unit) createMagicShower(targetX, targetY uint8, bState *battleState) {
	damage, missileKind, ok := u.resolveMagicShowerStats()
	if !ok {
		log.Printf("UWAGA: magicShower wywołany dla jednostki o nieobsługiwanym rodzaju %d!", u.Type)
		return
	}

	u.spawnMagicShowerProjectiles(targetX, targetY, missileKind, damage, bState)
}
