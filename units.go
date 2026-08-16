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
	if u.State == stateAttacking {
		return
	}
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

	u.applyDirection(dx, dy)

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
		u.applyDirection(dx, dy)
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

	// Nowy atak = start animacji od początku.
	if u.Type == unitArcher {
		u.updateArcherAttackAnimation()
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

	if u.Type == unitArcher {
		u.AnimStep = 0
		u.AnimTick = 0

		if u.AnimationType == animationFight {
			u.updateArcherAttackAnimation()
		} else {
			u.updateArcherAnimation()
		}

		return
	}
}

func (u *unit) updateArcherAnimation() {
	anim := u.getAnimationType()

	if anim != u.AnimationType {
		u.AnimationType = anim
		u.AnimStep = 0
		u.AnimTick = 0
	}

	// Walka jest obsługiwana osobno przez updateArcherAttackAnimation.
	if u.AnimationType == animationFight {
		return
	}

	u.AnimTick++
	if u.AnimTick >= 2 {
		u.AnimTick = 0
		u.AnimStep++
		if u.AnimStep >= animFramesCount {
			u.AnimStep = 0
		}
	}

	dir := u.animationDirection()
	u.SpriteID = spriteTable[u.Type][u.AnimationType][dir][u.AnimStep]
}

func (u *unit) updateArcherAttackAnimation() {
	if u.AnimationType != animationFight {
		u.AnimationType = animationFight
	}

	// Fazy ataku łucznika:
	// Delay > 4  -> Attack1 / zamach
	// Delay 3..4 -> Attack2 / wystrzał
	// Delay <= 2 -> Idle / złożenie i czekanie
	switch {
	case u.Delay > 4:
		u.AnimStep = 0
	case u.Delay > 2:
		u.AnimStep = 1
	default:
		u.AnimStep = 2
	}

	u.AnimTick = 0

	dir := u.animationDirection()
	u.SpriteID = spriteTable[u.Type][animationFight][dir][u.AnimStep]
}

func (u *unit) animationDirection() uint16 {
	dir := vectorToDirectionIndex(int(u.Direction.X), int(u.Direction.Y))

	if dir != uint16(directionNone) {
		return dir
	}

	if u.AnimationType == animationIdle {
		return uint16(directionNone)
	}

	if last, ok := u.rememberedDirection(); ok {
		return last
	}

	return uint16(directionDown)
}

func (u *unit) rememberDirection(dir uint16) {
	u.Facing = dir + 1
}

func (u *unit) rememberedDirection() (uint16, bool) {
	if u.Facing == 0 {
		return 0, false
	}

	return u.Facing - 1, true
}

func (u *unit) applyDirection(dx, dy int) {
	if dx == 0 && dy == 0 {
		return
	}

	u.Direction = rl.NewVector2(float32(dx), float32(dy))
	u.rememberDirection(vectorToDirectionIndex(dx, dy))
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
		u.applyDirection(dx, dy)
	}
}

func drawUnitArcher(u *unit, bState *battleState, pState *programState) {
	screenX := float32(u.X) * float32(tileWidth)
	screenY := float32(u.Y) * float32(tileHeight)

	if u.State == stateMoving {
		unitIndex := u.Type.getLegacyUnitIndex()
		delayIndex := u.Delay

		if delayIndex > maxPhaseDelay {
			delayIndex = 16
		}
		if delayIndex < minPhaseDelay {
			delayIndex = 0
		}

		rawShiftX := float32(spriteXOffsetByUnitTypeAndDelay[unitIndex][delayIndex])
		rawShiftY := float32(spriteYOffsetByUnitTypeAndDelay[unitIndex][delayIndex])

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

	finalID := getSpriteID(u)

	fmt.Printf("będę rysować finalID: %v, bo u.Delay to: %v\n", finalID, u.Delay)

	drawSpriteEx(finalID, screenX, screenY, u.Owner, rl.White, pState)

	if len(u.Wounds) > 0 {
		drawUnitWounds(u, screenX, screenY, pState)
	}

	drawActiveMagicShield(u, screenX, screenY, bState, pState)
}

// Pobiera jawne ID duszka z tablicy spriteTable na podstawie stanu jednostki.
// Zastępuje matematyczne obliczenia ze starego układu.
func getSpriteID(u *unit) uint16 {
	anim := u.getAnimationType()
	dir := u.animationDirection()

	delayIdx := int(u.Delay)
	if delayIdx > int(maxPhaseDelay) {
		delayIdx = int(maxPhaseDelay)
	}
	if delayIdx < int(minPhaseDelay) {
		delayIdx = int(minPhaseDelay)
	}

	return spriteTable[u.Type][anim][dir][delayIdx]
}

const (
	animFramesCount = 4
	directionsCount = 9
	animTypesCount  = 3

	delayStatesCount = 17
)

// @reminder: zastanów się, jak przechowywać ubitą jednostkę.
var spriteTable [unitTypeCount][animTypesCount][directionsCount][delayStatesCount]uint16

func initSpriteTable() {
	spriteTable[unitArcher][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
	}

	spriteTable[unitArcher][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
	}

	spriteTable[unitArcher][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
	}

	spriteTable[unitArcher][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
	}

	spriteTable[unitArcher][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
	}

	spriteTable[unitArcher][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
	}

	spriteTable[unitArcher][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
	}

	spriteTable[unitArcher][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
	}

	spriteTable[unitArcher][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
	}

	spriteTable[unitArcher][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove2Center,
		spriteManualArcherMove2Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Down,
		spriteManualArcherMove1Down,
		spriteManualArcherMove1Down,
		spriteManualArcherMove2Down,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
	}

	spriteTable[unitArcher][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove2Center,
		spriteManualArcherMove2Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Down,
		spriteManualArcherMove1Down,
		spriteManualArcherMove1Down,
		spriteManualArcherMove2Down,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
		spriteManualArcherMove1Center,
	}

	spriteTable[unitArcher][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherMove1Left,
		spriteManualArcherMove1Left,
		spriteManualArcherMove1Left,
		spriteManualArcherMove2Left,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
	}

	spriteTable[unitArcher][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherMove1Up,
		spriteManualArcherMove1Up,
		spriteManualArcherMove1Up,
		spriteManualArcherMove2Up,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
	}

	spriteTable[unitArcher][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherMove1Right,
		spriteManualArcherMove1Right,
		spriteManualArcherMove1Right,
		spriteManualArcherMove2Right,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
	}

	spriteTable[unitArcher][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherMove1UpLeft,
		spriteManualArcherMove1UpLeft,
		spriteManualArcherMove1UpLeft,
		spriteManualArcherMove2UpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
	}

	spriteTable[unitArcher][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherMove1UpRight,
		spriteManualArcherMove1UpRight,
		spriteManualArcherMove1UpRight,
		spriteManualArcherMove2UpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
	}

	spriteTable[unitArcher][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherMove1DownLeft,
		spriteManualArcherMove1DownLeft,
		spriteManualArcherMove1DownLeft,
		spriteManualArcherMove2DownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
	}

	spriteTable[unitArcher][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherMove1DownRight,
		spriteManualArcherMove1DownRight,
		spriteManualArcherMove1DownRight,
		spriteManualArcherMove2DownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
	}

	spriteTable[unitArcher][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherAttack2Down,
		spriteManualArcherAttack2Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
	}

	spriteTable[unitArcher][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherIdleDown,
		spriteManualArcherAttack2Down,
		spriteManualArcherAttack2Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
		spriteManualArcherAttack1Down,
	}

	spriteTable[unitArcher][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherIdleLeft,
		spriteManualArcherAttack2Left,
		spriteManualArcherAttack2Left,
		spriteManualArcherAttack1Left,
		spriteManualArcherAttack1Left,
		spriteManualArcherAttack1Left,
		spriteManualArcherAttack1Left,
		spriteManualArcherAttack1Left,
		spriteManualArcherAttack1Left,
		spriteManualArcherAttack1Left,
		spriteManualArcherAttack1Left,
		spriteManualArcherAttack1Left,
		spriteManualArcherAttack1Left,
		spriteManualArcherAttack1Left,
		spriteManualArcherAttack1Left,
	}

	spriteTable[unitArcher][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherIdleUp,
		spriteManualArcherAttack2Up,
		spriteManualArcherAttack2Up,
		spriteManualArcherAttack1Up,
		spriteManualArcherAttack1Up,
		spriteManualArcherAttack1Up,
		spriteManualArcherAttack1Up,
		spriteManualArcherAttack1Up,
		spriteManualArcherAttack1Up,
		spriteManualArcherAttack1Up,
		spriteManualArcherAttack1Up,
		spriteManualArcherAttack1Up,
		spriteManualArcherAttack1Up,
		spriteManualArcherAttack1Up,
		spriteManualArcherAttack1Up,
	}

	spriteTable[unitArcher][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherIdleRight,
		spriteManualArcherAttack2Right,
		spriteManualArcherAttack2Right,
		spriteManualArcherAttack1Right,
		spriteManualArcherAttack1Right,
		spriteManualArcherAttack1Right,
		spriteManualArcherAttack1Right,
		spriteManualArcherAttack1Right,
		spriteManualArcherAttack1Right,
		spriteManualArcherAttack1Right,
		spriteManualArcherAttack1Right,
		spriteManualArcherAttack1Right,
		spriteManualArcherAttack1Right,
		spriteManualArcherAttack1Right,
		spriteManualArcherAttack1Right,
	}

	spriteTable[unitArcher][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherIdleUpLeft,
		spriteManualArcherAttack2UpLeft,
		spriteManualArcherAttack2UpLeft,
		spriteManualArcherAttack1UpLeft,
		spriteManualArcherAttack1UpLeft,
		spriteManualArcherAttack1UpLeft,
		spriteManualArcherAttack1UpLeft,
		spriteManualArcherAttack1UpLeft,
		spriteManualArcherAttack1UpLeft,
		spriteManualArcherAttack1UpLeft,
		spriteManualArcherAttack1UpLeft,
		spriteManualArcherAttack1UpLeft,
		spriteManualArcherAttack1UpLeft,
		spriteManualArcherAttack1UpLeft,
		spriteManualArcherAttack1UpLeft,
	}

	spriteTable[unitArcher][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherIdleUpRight,
		spriteManualArcherAttack2UpRight,
		spriteManualArcherAttack2UpRight,
		spriteManualArcherAttack1UpRight,
		spriteManualArcherAttack1UpRight,
		spriteManualArcherAttack1UpRight,
		spriteManualArcherAttack1UpRight,
		spriteManualArcherAttack1UpRight,
		spriteManualArcherAttack1UpRight,
		spriteManualArcherAttack1UpRight,
		spriteManualArcherAttack1UpRight,
		spriteManualArcherAttack1UpRight,
		spriteManualArcherAttack1UpRight,
		spriteManualArcherAttack1UpRight,
		spriteManualArcherAttack1UpRight,
	}

	spriteTable[unitArcher][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherIdleDownLeft,
		spriteManualArcherAttack2DownLeft,
		spriteManualArcherAttack2DownLeft,
		spriteManualArcherAttack1DownLeft,
		spriteManualArcherAttack1DownLeft,
		spriteManualArcherAttack1DownLeft,
		spriteManualArcherAttack1DownLeft,
		spriteManualArcherAttack1DownLeft,
		spriteManualArcherAttack1DownLeft,
		spriteManualArcherAttack1DownLeft,
		spriteManualArcherAttack1DownLeft,
		spriteManualArcherAttack1DownLeft,
		spriteManualArcherAttack1DownLeft,
		spriteManualArcherAttack1DownLeft,
		spriteManualArcherAttack1DownLeft,
	}

	spriteTable[unitArcher][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherIdleDownRight,
		spriteManualArcherAttack2DownRight,
		spriteManualArcherAttack2DownRight,
		spriteManualArcherAttack1DownRight,
		spriteManualArcherAttack1DownRight,
		spriteManualArcherAttack1DownRight,
		spriteManualArcherAttack1DownRight,
		spriteManualArcherAttack1DownRight,
		spriteManualArcherAttack1DownRight,
		spriteManualArcherAttack1DownRight,
		spriteManualArcherAttack1DownRight,
		spriteManualArcherAttack1DownRight,
		spriteManualArcherAttack1DownRight,
		spriteManualArcherAttack1DownRight,
		spriteManualArcherAttack1DownRight,
	}

	// Koniec Łucznika
	spriteTable[unitCow][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown,
	}
	spriteTable[unitCow][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown,
	}
	spriteTable[unitCow][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft,
	}
	spriteTable[unitCow][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp,
	}
	spriteTable[unitCow][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight,
	}
	spriteTable[unitCow][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft,
	}
	spriteTable[unitCow][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight,
	}
	spriteTable[unitCow][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft,
	}
	spriteTable[unitCow][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight,
	}

	spriteTable[unitCow][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowMove1Down, spriteManualCowMove1Down, spriteManualCowMove1Down, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowMove2Down, spriteManualCowMove2Down, spriteManualCowMove2Down, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowMove1Down, spriteManualCowMove1Down, spriteManualCowIdleDown,
	}
	spriteTable[unitCow][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowMove1Down, spriteManualCowMove1Down, spriteManualCowMove1Down, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowMove2Down, spriteManualCowMove2Down, spriteManualCowMove2Down, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowMove1Down, spriteManualCowMove1Down, spriteManualCowIdleDown,
	}
	spriteTable[unitCow][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowMove1Left, spriteManualCowMove1Left, spriteManualCowMove1Left, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowMove2Left, spriteManualCowMove2Left, spriteManualCowMove2Left, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowMove1Left, spriteManualCowMove1Left, spriteManualCowIdleLeft,
	}
	spriteTable[unitCow][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowMove1Up, spriteManualCowMove1Up, spriteManualCowMove1Up, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowMove2Up, spriteManualCowMove2Up, spriteManualCowMove2Up, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowMove1Up, spriteManualCowMove1Up, spriteManualCowIdleUp,
	}
	spriteTable[unitCow][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowMove1Right, spriteManualCowMove1Right, spriteManualCowMove1Right, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowMove2Right, spriteManualCowMove2Right, spriteManualCowMove2Right, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowMove1Right, spriteManualCowMove1Right, spriteManualCowIdleRight,
	}
	spriteTable[unitCow][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowMove1UpLeft, spriteManualCowMove1UpLeft, spriteManualCowMove1UpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowMove2UpLeft, spriteManualCowMove2UpLeft, spriteManualCowMove2UpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowMove1UpLeft, spriteManualCowMove1UpLeft, spriteManualCowIdleUpLeft,
	}
	spriteTable[unitCow][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowMove1UpRight, spriteManualCowMove1UpRight, spriteManualCowMove1UpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowMove2UpRight, spriteManualCowMove2UpRight, spriteManualCowMove2UpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowMove1UpRight, spriteManualCowMove1UpRight, spriteManualCowIdleUpRight,
	}
	spriteTable[unitCow][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowMove1DownLeft, spriteManualCowMove1DownLeft, spriteManualCowMove1DownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowMove2DownLeft, spriteManualCowMove2DownLeft, spriteManualCowMove2DownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowMove1DownLeft, spriteManualCowMove1DownLeft, spriteManualCowIdleDownLeft,
	}
	spriteTable[unitCow][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowMove1DownRight, spriteManualCowMove1DownRight, spriteManualCowMove1DownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowMove2DownRight, spriteManualCowMove2DownRight, spriteManualCowMove2DownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowMove1DownRight, spriteManualCowMove1DownRight, spriteManualCowIdleDownRight,
	}

	spriteTable[unitCow][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowAttack2Down, spriteManualCowAttack2Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down,
	}
	spriteTable[unitCow][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowIdleDown, spriteManualCowAttack2Down, spriteManualCowAttack2Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down, spriteManualCowAttack1Down,
	}
	spriteTable[unitCow][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowIdleLeft, spriteManualCowAttack2Left, spriteManualCowAttack2Left, spriteManualCowAttack1Left, spriteManualCowAttack1Left, spriteManualCowAttack1Left, spriteManualCowAttack1Left, spriteManualCowAttack1Left, spriteManualCowAttack1Left, spriteManualCowAttack1Left, spriteManualCowAttack1Left, spriteManualCowAttack1Left, spriteManualCowAttack1Left, spriteManualCowAttack1Left, spriteManualCowAttack1Left,
	}
	spriteTable[unitCow][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowIdleUp, spriteManualCowAttack2Up, spriteManualCowAttack2Up, spriteManualCowAttack1Up, spriteManualCowAttack1Up, spriteManualCowAttack1Up, spriteManualCowAttack1Up, spriteManualCowAttack1Up, spriteManualCowAttack1Up, spriteManualCowAttack1Up, spriteManualCowAttack1Up, spriteManualCowAttack1Up, spriteManualCowAttack1Up, spriteManualCowAttack1Up, spriteManualCowAttack1Up,
	}
	spriteTable[unitCow][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowIdleRight, spriteManualCowAttack2Right, spriteManualCowAttack2Right, spriteManualCowAttack1Right, spriteManualCowAttack1Right, spriteManualCowAttack1Right, spriteManualCowAttack1Right, spriteManualCowAttack1Right, spriteManualCowAttack1Right, spriteManualCowAttack1Right, spriteManualCowAttack1Right, spriteManualCowAttack1Right, spriteManualCowAttack1Right, spriteManualCowAttack1Right, spriteManualCowAttack1Right,
	}
	spriteTable[unitCow][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowIdleUpLeft, spriteManualCowAttack2UpLeft, spriteManualCowAttack2UpLeft, spriteManualCowAttack1UpLeft, spriteManualCowAttack1UpLeft, spriteManualCowAttack1UpLeft, spriteManualCowAttack1UpLeft, spriteManualCowAttack1UpLeft, spriteManualCowAttack1UpLeft, spriteManualCowAttack1UpLeft, spriteManualCowAttack1UpLeft, spriteManualCowAttack1UpLeft, spriteManualCowAttack1UpLeft, spriteManualCowAttack1UpLeft, spriteManualCowAttack1UpLeft,
	}
	spriteTable[unitCow][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowIdleUpRight, spriteManualCowAttack2UpRight, spriteManualCowAttack2UpRight, spriteManualCowAttack1UpRight, spriteManualCowAttack1UpRight, spriteManualCowAttack1UpRight, spriteManualCowAttack1UpRight, spriteManualCowAttack1UpRight, spriteManualCowAttack1UpRight, spriteManualCowAttack1UpRight, spriteManualCowAttack1UpRight, spriteManualCowAttack1UpRight, spriteManualCowAttack1UpRight, spriteManualCowAttack1UpRight, spriteManualCowAttack1UpRight,
	}
	spriteTable[unitCow][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowIdleDownLeft, spriteManualCowAttack2DownLeft, spriteManualCowAttack2DownLeft, spriteManualCowAttack1DownLeft, spriteManualCowAttack1DownLeft, spriteManualCowAttack1DownLeft, spriteManualCowAttack1DownLeft, spriteManualCowAttack1DownLeft, spriteManualCowAttack1DownLeft, spriteManualCowAttack1DownLeft, spriteManualCowAttack1DownLeft, spriteManualCowAttack1DownLeft, spriteManualCowAttack1DownLeft, spriteManualCowAttack1DownLeft, spriteManualCowAttack1DownLeft,
	}
	spriteTable[unitCow][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowIdleDownRight, spriteManualCowAttack2DownRight, spriteManualCowAttack2DownRight, spriteManualCowAttack1DownRight, spriteManualCowAttack1DownRight, spriteManualCowAttack1DownRight, spriteManualCowAttack1DownRight, spriteManualCowAttack1DownRight, spriteManualCowAttack1DownRight, spriteManualCowAttack1DownRight, spriteManualCowAttack1DownRight, spriteManualCowAttack1DownRight, spriteManualCowAttack1DownRight, spriteManualCowAttack1DownRight, spriteManualCowAttack1DownRight,
	}
}

type corpseSpriteSet struct {
	Fresh uint16
	Decay uint16
}

var corpseSprites [unitTypeCount]corpseSpriteSet

func initCorpseSpriteTable() {
	corpseSprites[unitArcher] = corpseSpriteSet{
		Fresh: spriteManualArcherCorpseFresh,
		Decay: spriteManualArcherCorpseDecay,
	}

	corpseSprites[unitAxeman] = corpseSpriteSet{
		Fresh: spriteManualAxemanCorpseFresh,
		Decay: spriteManualAxemanCorpseDecay,
	}

	corpseSprites[unitBear] = corpseSpriteSet{
		Fresh: spriteManualBearCorpseFresh,
		Decay: spriteManualBearCorpseDecay,
	}

	corpseSprites[unitCommander] = corpseSpriteSet{
		Fresh: spriteManualCommanderCorpseFresh,
		Decay: spriteManualCommanderCorpseDecay,
	}

	corpseSprites[unitCow] = corpseSpriteSet{
		Fresh: spriteManualCowCorpseFresh,
		Decay: spriteManualCowCorpseDecay,
	}

	corpseSprites[unitCrossbowman] = corpseSpriteSet{
		Fresh: spriteManualCrossbowmanCorpseFresh,
		Decay: spriteManualCrossbowmanCorpseDecay,
	}

	corpseSprites[unitMage] = corpseSpriteSet{
		Fresh: spriteManualMageCorpseFresh,
		Decay: spriteManualMageCorpseDecay,
	}

	corpseSprites[unitPriest] = corpseSpriteSet{
		Fresh: spriteManualPriestCorpseFresh,
		Decay: spriteManualPriestCorpseDecay,
	}

	corpseSprites[unitPriestess] = corpseSpriteSet{
		Fresh: spriteManualPriestressCorpseFresh,
		Decay: spriteManualPriestressCorpseDecay,
	}

	corpseSprites[unitShepherd] = corpseSpriteSet{
		Fresh: spriteManualShepherdCorpseFresh,
		Decay: spriteManualShepherdCorpseDecay,
	}

	corpseSprites[unitSpearman] = corpseSpriteSet{
		Fresh: spriteManualSpearmanCorpseFresh,
		Decay: spriteManualSpearmanCorpseDecay,
	}

	corpseSprites[unitSwordsman] = corpseSpriteSet{
		Fresh: spriteManualSwordsmanCorpseFresh,
		Decay: spriteManualSwordsmanCorpseDecay,
	}

	corpseSprites[unitUnknown] = corpseSpriteSet{
		Fresh: spriteManualUnknownCorpseFresh,
		Decay: spriteManualUnknownCorpseDecay,
	}
}
