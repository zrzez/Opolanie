package main

// units.go

import (
	"log"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// isCaster zwraca true jeżeli dana jednostka czaruje.
func (ut unitType) isCaster() bool {
	return ut == unitMage || ut == unitPriest || ut == unitPriestess
}

// isMagical zwraca true jeżeli dana jednostka ma MaxMana > 0.
func (ut unitType) isMagical() bool {
	return ut.isCaster() || ut == unitUnknown
}

// increaseManaUnit zwiększa manę jednostki o amount.
// Pilnuje, aby u.Mana <= u.MaxMana.
func (u *unit) increaseManaUnit(amount uint16) {
	// @reminder: sprawdź, czy sprawdzanie !u.Exists ma sens
	if !u.Exists || amount == 0 {
		return
	}

	// @reminder: najwięcej many ma unitMage i unitUnknown 15 poziomu, tj. 60+280=340.
	// Czyli nie ma możliwości, aby przyrost many przekroczył górną granicę
	// uint16 i „przekręcił się” w okolice zera. Dlatego nie ma bezpiecznika.
	u.Mana += amount

	if u.Mana > u.MaxMana {
		u.Mana = u.MaxMana
	}
}

// tryToDecreaseMana zmniejsza manę o amount oraz zwraca prawda. Jeśli u.Mana < amount, to zwraca fałsz.
func (u *unit) tryToDecreaseMana(amount uint16) bool {
	if !u.Exists {
		return false
	}

	if u.Mana < amount {
		return false
	}

	u.Mana -= amount

	return true
}

// increaseHPUnit zwiększa punkty życia o amount. Pilnuje, aby u.HP <= u.MaxHP.
func (u *unit) increaseHPUnit(amount uint16) {
	if !u.Exists {
		return
	}

	u.HP += amount

	if u.HP > u.MaxHP {
		u.HP = u.MaxHP
	}
}

// @reminder: jeszcze nie usuwam, być może się przyda.
// decreaseHPUnit dla każdej istniejącej jednostki zmniejsza PŻ o amount.
// Pilnuje, aby ustawić u.Exists = false.
func (u *unit) decreaseHPUnit(amount uint16) {
	if !u.Exists {
		return
	}

	u.HP -= amount

	if u.HP > 0 {
		return
	}

	u.HP = 0
	u.Exists = false
}

// ============================================================================
// LOGIKA JEDNOSTEK
// ============================================================================

// @reminder: jest to metoda odpowiadająca za aktualizowanie wszystkiego w jednostce.
/*
Podejrzewam, że:
1) branie bState *battleState jako argumentu, czyli „całego świata gry” jest efektem poważnego błędu.
2) są tutaj zbyteczne rzeczy, które już robię w innej części kodu
3) uda się mocno odchudzić kod bez utraty funkcjonalności
4) nie przystaje on do aktualnej architektury
5) utknę na długie tygodnie - 11.07.2026

Chyba będzie najlepiej jeśli zacznę iść od samego spodu. W ten sposób zmiany same wypłyną i będą oczywiste
po spojrzeniu na sygnatury i zmiany w nazwach.


*/
func (u *unit) handleAttackCooldown(currentTick uint16) {
	if currentTick%logicSpeedDivisor == 0 {
		if u.AttackCooldown > 0 {
			u.AttackCooldown--
		}
	}
}

func (u *unit) handleDelay(currentTick uint16) bool {
	if u.State == stateIdle {
		u.Delay = u.MaxDelay

		return false
	}
	if u.Delay <= 0 {
		return false
	}
	if currentTick%logicSpeedDivisor != 0 {
		return true
	}
	u.Delay--

	return u.Delay > 0
}

// @todo: ogarnij co to robi, bo zupełnie nie pamiętam.
func (u *unit) handleBlockedCounter() bool {
	if u.BlockedCounter > 0 {
		u.BlockedCounter--
		u.AnimationType = "idle"
		u.AnimationFrame = 0
		u.Delay = u.MaxDelay

		return true
	}

	return false
}

// @todo: ogarnij co to robi, bo zupełnie nie pamiętam.
func (u *unit) handleWaitingToActiveTransition() {
	if u.State == stateWaiting {
		u.State = u.determineActiveStateFromCommand()
		if u.State == stateAttacking {
			u.AnimationType = "fight"
		} else {
			u.AnimationType = "walk"
		}
	}
}

func (u *unit) resetDelayIfActive() {
	if u.State != stateWaiting {
		if u.Delay < u.MaxDelay {
			u.Delay = u.MaxDelay
		}
	}
}

func (u *unit) determineActiveStateFromCommand() unitState {
	switch u.Command {
	case cmdUMove:
		return stateMoving
	case cmdUAttack:
		return stateAttacking
	case cmdUGraze:
		return stateGrazing
	case cmdUBuild:
		return stateBuilding
	case cmdURepair:
		return stateRepairing
	case cmdUCastSpell:
		return stateCastingSpell
	default:
		return stateIdle
	}
}

func (ut unitType) canDamagePalisades() bool {
	return ut == unitAxeman || ut == unitPriest
}

func (ct commandType) isInteraction() bool {
	switch ct {
	case cmdUAttack, cmdUBuild, cmdURepair, cmdBPlaceConstruction, cmdUCastSpell:
		return true
	default:
		return false
	}
}

// caDamageTree sprawdza, czy jednostka może zaatakować dane drzewo.
func (u *unit) canDamageTree(treeTile *tile) bool {
	if !treeTile.isStandingTree() || treeTile.IsBurning {
		return false
	}

	switch u.Type {
	case unitPriest:
		return true
	case unitAxeman:
		return treeTile.isDryTree()
	default:
		return false
	}
}

func (u *unit) prepareForNewCommand(cmdType commandType, intentionX, intentionY uint8, targetID ObjectID, approachX, approachY uint8) {
	u.clearPath()
	u.History = nil
	u.LoopCount = 0
	u.TicksNoProgress = 0
	u.LastPathIndex = 0
	u.Command = cmdType
	u.TargetX = intentionX
	u.TargetY = intentionY
	u.ApproachX = approachX
	u.ApproachY = approachY
	u.TargetID = targetID
	u.Delay = 0
}

func (u *unit) setIdle() {
	u.State = stateIdle
	u.AnimationType = "idle"
	u.Command = cmdUIdle
	u.clearPath()
	u.BlockedCounter = 0
	u.RetryAttempts = 0
	u.PathfindingCooldown = 0
	u.AllowFriendlyFire = false

	if u.State != stateWaiting {
		u.IsInQueue = false
	}
}

func (u *unit) updateMovementAnimation(prevX, prevY uint8) {
	dx := int(u.X) - int(prevX)
	dy := int(u.Y) - int(prevY)

	if dx != 0 || dy != 0 {
		u.Direction = rl.NewVector2(float32(dx), float32(dy))
	}

	u.AnimationType = "walk"
	u.AnimationCounter++

	if u.AnimationCounter >= animationSpeed {
		u.AnimationCounter = 0
		u.AnimationFrame++

		if u.AnimationFrame >= walkAnimationFrames {
			u.AnimationFrame = 0
		}
	}
}

func (u *unit) faceTarget(target *combatTarget) {
	var tx, ty uint8

	if target.Unit != nil {
		tx, ty = target.Unit.X, target.Unit.Y
	} else if target.Building != nil {
		// Dla budynków celujemy w ich środek lub najbliższy punkt
		tx, ty, _ = target.Building.getClosestOccupiedTile(u.X, u.Y)
	} else if target.Tile != nil {
		tx, ty = target.Tile.X, target.Tile.Y
	} else {
		return
	}

	dx := 0
	dy := 0

	if tx > u.X {
		dx = 1
	} else if tx < u.X {
		dx = -1
	}

	if ty > u.Y {
		dy = 1
	} else if ty < u.Y {
		dy = -1
	}

	// Ustawiamy wektor kierunku
	if dx != 0 || dy != 0 {
		u.Direction = rl.NewVector2(float32(dx), float32(dy))
	}
}

func (u *unit) getProjectileManaCost() uint16 {
	// @reminder: inne typy jednostek nie rzucają magicznych pocisków.
	// Dlatego mam default zwracający zawsze prawdę i nie rozpisałem
	// wszystkich typów.
	switch u.Type { //nolint:exhaustive
	case unitPriestess:
		return magicThunderManaCost
	case unitPriest:
		return magicFireManaCost
	case unitMage:
		return magicGhostManaCost
	default:
		return 0
	}
}

// @todo: @reminder: wróć i zastanów się, czy to rozdzielenie ma sens, bo wygląda
// jak przekombinowane, zwykły przełącznik bez podfunkcji powinien też się sprawdzić.
func (u *unit) setAttackTimings() {
	if u.AttackRange > 1 {
		u.setRangedTimings()
	} else {
		u.setMeleeTimings()
	}
}

func (u *unit) setRangedTimings() {
	u.Delay = 12

	switch u.Type {
	case unitCrossbowman:
		u.AttackCooldown = 70
	case unitMage, unitPriest, unitPriestess:
		u.AttackCooldown = 100
	default:
		u.AttackCooldown = 65
	}
}

func (u *unit) setMeleeTimings() {
	u.Delay = 10

	switch u.Type {
	case unitBear:
		u.AttackCooldown = 35
	case unitUnknown:
		u.AttackCooldown = 15
	default:
		u.AttackCooldown = 25
	}
}

// To już chyba nie powinna być metoda jednostki. Nic jej do szczegółów pocisków.
func (u *unit) resolveMagicShowerStats() (damage uint16, missileKind uint8, ok bool) {
	// @reminder: inne rodzaje jednostek nie mają tego czaru więc ich nie wymieniam tutaj.
	switch u.Type {
	case unitPriest:
		return spellDamageFireShower, missileFireRain, true
	case unitPriestess:
		return spellDamageLightningShower, missileLightning, true
	default:
		return 0, 0, false
	}
}

// @todo: przenieś logikę tworzenia ran, bo jednostka nie musi o tym wiedzieć.
func (u *unit) takeDamage(damage uint16, bState *battleState) {
	// 0. Sprawdzamy, czy jednostka jest chroniona przed obrażeniami.
	if u.hasMagicShield {
		return
	}

	var finalDamage uint16

	if damage <= u.Armor {
		finalDamage = 1
	} else {
		finalDamage = damage - u.Armor
	}

	if finalDamage <= u.HP {
		u.HP -= finalDamage
	} else {
		u.HP = 0
	}

	log.Printf("unit %d otrzymał %d obrażeń. HP: %d/%d", u.ID, finalDamage, u.HP, u.MaxHP)
	// === Zapisywanie informacji o odniesionych ranach

	if len(u.Wounds) < maxWoundsCount {
		offX := float32(rand.Intn(11) - 4)
		offY := float32(rand.Intn(9) - 3)

		isSevere := false
		// @todo zastąp to prostym sprawdzeniem Armor dla ranionej jednostki
		// lub % MaxHP, żeby to było wizualnie czytelniejsze w trakcie bitwy.
		// Za maxhp przemawia, że w finalDamage już jest ukryta poprawka na u.Armor
		isUnitLight := u.Type != unitCommander && u.Type != unitSpearman && u.Type != unitSwordsman &&
			u.Type != unitBear && u.Type != unitUnknown && u.Type != unitCrossbowman

		var baseScale float32 = 1.0

		if finalDamage > severeDamage && isUnitLight {
			isSevere = true
			baseScale = 1.1
		}
		// @todo: skala powinna zależeć od ilości obrażeń i tylko dla isSevere

		rotation := rand.Float32() * 120.0
		newWound := wound{
			Timer:    20,
			OffsetX:  offX,
			OffsetY:  offY,
			IsSevere: isSevere,
			Scale:    baseScale,
			Rotation: rotation,
		}

		u.Wounds = append(u.Wounds, newWound)
	}

	if u.Type == unitCow && u.Exists {
		if u.Udder < fullUdderAmount && u.Command != cmdUFlee {
			barnX, barnY, foundBarn := findNearestBarnMilkingSpot(u, bState)
			if foundBarn {
				cmd := &command{
					ActionType:          cmdUFlee,
					TargetX:             barnX,
					TargetY:             barnY,
					InteractionTargetID: 0,
				}
				u.addUnitCommand(cmd, bState)
				log.Printf("unit %d (COW): Otrzymała obrażenia, ucieka do obory na (%d,%d).", u.ID, barnX, barnY)
			} else {
				log.Printf("unit %d (COW): Otrzymała obrażenia, ale nie znalazła obory do ucieczki. ", u.ID)
			}
		}
	}

	// Sprawdzamy, czy jednostka została zabita
	if u.HP == 0 {
		u.Exists = false

		// @todo: to nie powinno w ogóle być w jednostce, a tam obok aktualizowania
		//  budynków. Bo mamy dostęp do potrzebnych rzeczy.
		// Zabita jednostka nie powinna zliczać się do górnej granicy ludności
		bState.decreasePopulation(u.Owner)

		occupiedTile := &bState.Board.Tiles[u.X][u.Y]
		if occupiedTile.Unit == u {
			occupiedTile.Unit = nil
		}

		bState.createCorpses(u)
		u.unregisterFromBuilding()
		log.Printf("Jednostka %d została zabita!", u.ID)
	}
}

func (u *unit) unregisterFromBuilding() {
	if u.BelongsTo != nil {
		log.Printf("Jednostka %d próbuje się wyrejestrować z budynku %d", u.ID, u.BelongsTo.ID)
		// @reminder: metoda unregisterUnit zwraca bool, który ignorujemy
		u.BelongsTo.unregisterUnit(u.ID)
		u.BelongsTo = nil
	}
}

func (ut unitType) getLegacyUnitIndex() int {
	return int(ut)
}
