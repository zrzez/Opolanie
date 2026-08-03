package main

// units.go

import (
	"fmt"

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
func (u *unit) handleWaitingToActiveTransition() {
	if u.State == stateWaiting {
		u.State = u.determineActiveStateFromCommand()
		u.setAnimationType()
	}
}

func (u *unit) ensureDelayIsSet() {
	if u.State != stateWaiting && u.State != stateIdle && u.Delay == 0 {
		fmt.Print("DUPA ZAPOMŁES USTAWIĆ DELAY!!!!\n")
		u.Delay = u.MaxDelay
	}
}

func (u *unit) determineActiveStateFromCommand() unitState {
	switch u.Command {
	case cmdUMove:
		return stateMoving
	case cmdUAttack:
		return stateMoving
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

func (u *unit) prepareForNewCommand(cmdType commandType, target targetReference, approach point) {
	u.clearPath()
	u.Command = cmdType
	u.Target = target
	u.Approach = approach

	u.Delay = 0
	u.facePoint(approach)
}

func (u *unit) setIdle() {
	u.State = stateIdle
	u.setAnimationType()
	u.Command = cmdUIdle
	u.clearPath()
	u.RetryAttempts = 0
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
		targetCoords, _ := getClosestOccupiedTile(&point{X: u.X, Y: u.Y}, &target.Building.OccupiedTiles)
		tx, ty = targetCoords.X, targetCoords.Y
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

func (u *unit) isLightType() bool {
	return u.Type != unitCommander && u.Type != unitSpearman && u.Type != unitSwordsman &&
		u.Type != unitBear && u.Type != unitUnknown && u.Type != unitCrossbowman
}

func (u *unit) takeDamage(damage uint16) {
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

	u.wasAttacked = true

	createWound(finalDamage, u)
}

func (u *unit) unregisterFromBuilding() {
	if u.BelongsTo != nil {
		// @reminder: metoda unregisterUnit zwraca bool, który ignorujemy
		u.BelongsTo.unregisterUnit(u.ID)
		u.BelongsTo = nil
	}
}

// Na podstawie rodzaju jednostki zwraca indeks potrzebny do
// doboru odpowiedniego przesunięcia duszka w czasie wykonywania działań.
func (ut unitType) getLegacyUnitIndex() int {
	return int(ut)
}

func (u *unit) getAnimationType() animationType {
	switch u.State {
	case stateIdle, stateWaiting, stateMilking:
		return animationIdle
	case stateMoving:
		return animationWalk
	case stateAttacking, stateRepairing, stateBuilding, stateCastingSpell:
		return animationFight
	case stateGrazing: // @reminder: 02.08.2026 nie jestem pewien w tej chwili więc dam walkę
		return animationFight
	default:
		return animationIdle
	}
}

func (u *unit) setAnimationType() {
	u.AnimationType = u.getAnimationType()
}

// @reminder: próba napisania nowego układu rysowania jednostek 03.08.2026
// Układ zmieniający wygląd łucznika w zależności od tego, co robi.
func (u *unit) updateArcherAnimation() {
	anim := u.getAnimationType()

	// Jeśli zmienił się stan, to zaczynamy animację od zera
	if anim != u.AnimationType {
		u.AnimationType = anim
		u.AnimStep = 0
		u.AnimTick = 0
	}

	// Co tyknięcie podmieniamy grafikę animacji
	// @reminder: dzieje się to zbyt szybko, chyba zapomniałem „dzielić przez szybkość logiki”!
	// bezczynność przy animationSpeed = 24 wygląda spoko
	// atak przy animationSpeed = ? wygląda spoko
	// chodzenie przy animationSpeed = ? wygląda spoko
	u.AnimTick++

	if u.AnimTick >= animationSpeed {
		u.AnimTick = 0
		u.AnimStep++

		if u.AnimStep >= animFramesCount {
			u.AnimStep = 0
		}
	}

	dir := vectorToDirectionIndex(int(u.Direction.X), int(u.Direction.Y))

	u.SpriteID = spriteTable[u.Type][u.AnimationType][dir][u.AnimStep]
}

func (u *unit) facePoint(p point) {
	dx := int(p.X) - int(u.X)
	dy := int(p.Y) - int(u.Y)

	if dx > 0 {
		dx = 1
	} else if dx < 0 {
		dx = -1
	}

	if dy > 0 {
		dy = 1
	} else if dy < 0 {
		dy = -1
	}

	if dx != 0 || dy != 0 {
		u.Direction = rl.NewVector2(float32(dx), float32(dy))
	}
}

func drawUnitArcher(u *unit, bState *battleState, pState *programState) {
	// 1. Gdzie stoi jednostka na ekranie
	screenX := float32(u.X) * float32(tileWidth)
	screenY := float32(u.Y) * float32(tileHeight)

	// 2. Poprawki ułożenia grafiki na ekranie
	if u.State == stateMoving {
		// Każda jednostka ma swóją tablicę, musimy ustalić która to
		unitIndex := u.Type.getLegacyUnitIndex()
		// Z tablicy musimy pobrać odpowiednią wartość zależną od delay
		delayIndex := u.Delay

		// Delay = 16 jest dla bezczynności, dzięki temu
		// z tablicy dobieramy
		if delayIndex > maxPhaseDelay {
			delayIndex = 16
		}

		// Nie pozwalamy delay nie mniejsze niż zero, bo nie ma
		// takiego indeksu w tablicy
		if delayIndex < minPhaseDelay {
			delayIndex = 0
		}

		// Dzięki tym pokręconym obliczeniom mamy dokładne wartości dla
		// przesunięć X i Y sprawiające, że jednostka się płynnie przesuwa
		// po ekranie zamiast „przeskakiwać” z kafelka na kafelek.
		rawShiftX := float32(spriteXOffsetByUnitTypeAndDelay[unitIndex][delayIndex])
		rawShiftY := float32(spriteYOffsetByUnitTypeAndDelay[unitIndex][delayIndex])

		// Poprawki na kierunek patrzenia jednostki
		if int(u.Direction.X) > 0 {
			screenX -= rawShiftX
		} else if int(u.Direction.X) < 0 {
			screenX += rawShiftX
		}

		if int(u.Direction.Y) > 0 {
			screenY -= rawShiftY
		} else if int(u.Direction.Y) < 0 {
			screenY += rawShiftY
		}
	}

	// 3. Rysujemy
	drawSpriteEx(u.SpriteID, screenX, screenY, u.Owner, rl.White, pState)

	// 4. Dodatki
	if len(u.Wounds) > 0 {
		drawUnitWounds(u, screenX, screenY, pState)
	}

	drawActiveMagicShield(u, screenX, screenY, bState, pState)
}

const (
	animFramesCount = 4
	directionsCount = 9
	animTypesCount  = 3
)

var spriteTable [unitTypeCount][animTypesCount][directionsCount][animFramesCount]uint16

func initSpriteTable() {
	// Tutaj przechowuję zestawienie wszystkich grafik (tfu, duszków) używanych przez jednostki.

	// Bezczynność
	// Nie ma kierunku
	spriteTable[unitArcher][animationIdle][directionNone][0] = spriteManualArcherIdleDown  // DOBRZE
	spriteTable[unitArcher][animationIdle][directionNone][1] = spriteManualArcherIdleLeft  // DOBRZE
	spriteTable[unitArcher][animationIdle][directionNone][2] = spriteManualArcherIdleUp    // DOBRZE
	spriteTable[unitArcher][animationIdle][directionNone][3] = spriteManualArcherIdleRight // DOBRZE

	// Chodzenie
	// w dół
	spriteTable[unitArcher][animationWalk][directionDown][0] = spriteManualArcherMove1Center // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionDown][1] = spriteManualArcherMove1Down   // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionDown][2] = spriteManualArcherMove2Center // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionDown][3] = spriteManualArcherMove2Down   // DOBRZE, ALE ZA WOLNO ZMIENIA!

	// w lewo
	spriteTable[unitArcher][animationWalk][directionLeft][0] = spriteManualArcherIdleLeft  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionLeft][1] = spriteManualArcherMove1Left // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionLeft][2] = spriteManualArcherIdleLeft  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionLeft][3] = spriteManualArcherMove2Left // DOBRZE, ALE ZA WOLNO ZMIENIA!

	// w górę
	spriteTable[unitArcher][animationWalk][directionUp][0] = spriteManualArcherIdleUp  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionUp][1] = spriteManualArcherMove1Up // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionUp][2] = spriteManualArcherIdleUp  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionUp][3] = spriteManualArcherMove2Up // DOBRZE, ALE ZA WOLNO ZMIENIA!

	// w prawo (w lewo + zwierciadlane odbicie)
	spriteTable[unitArcher][animationWalk][directionRight][0] = spriteManualArcherIdleRight  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionRight][1] = spriteManualArcherMove1Right // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionRight][2] = spriteManualArcherIdleRight  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionRight][3] = spriteManualArcherMove2Right // DOBRZE, ALE ZA WOLNO ZMIENIA!

	// w lewy górny róg
	spriteTable[unitArcher][animationWalk][directionUpLeft][0] = spriteManualArcherIdleUpLeft  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionUpLeft][1] = spriteManualArcherMove1UpLeft // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionUpLeft][2] = spriteManualArcherIdleUpLeft  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionUpLeft][3] = spriteManualArcherMove2UpLeft // DOBRZE, ALE ZA WOLNO ZMIENIA!

	// w prawy górny róg (zwierciadlane odbicie)
	spriteTable[unitArcher][animationWalk][directionUpRight][0] = spriteManualArcherIdleUpRight  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionUpRight][1] = spriteManualArcherMove1UpRight // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionUpRight][2] = spriteManualArcherIdleUpRight  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionUpRight][3] = spriteManualArcherMove2UpRight // DOBRZE, ALE ZA WOLNO ZMIENIA!

	// w lewy dolny róg
	spriteTable[unitArcher][animationWalk][directionDownLeft][0] = spriteManualArcherIdleDownLeft  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionDownLeft][1] = spriteManualArcherMove1DownLeft // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionDownLeft][2] = spriteManualArcherIdleDownLeft  // DOBRZE, ALE ZA WOLNO ZMIENIA!
	spriteTable[unitArcher][animationWalk][directionDownLeft][3] = spriteManualArcherMove2DownLeft // DOBRZE, ALE ZA WOLNO ZMIENIA!

	// w prawy dolny róg (zwierdlane odbicie)
	spriteTable[unitArcher][animationWalk][directionDownRight][0] = spriteManualArcherIdleDownRight  // cośtam
	spriteTable[unitArcher][animationWalk][directionDownRight][1] = spriteManualArcherMove1DownRight // cośtam
	spriteTable[unitArcher][animationWalk][directionDownRight][2] = spriteManualArcherIdleDownRight  // cośtam
	spriteTable[unitArcher][animationWalk][directionDownRight][3] = spriteManualArcherMove2DownRight // cośtam
}
