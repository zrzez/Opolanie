package main

// tiles.go

// Pomagierzy do sprawdzania rodzaju tekstury kafelka.

/*
Każda zmiana w tym pliku wymaga ode mnie odpowiedzenia sobie na pytanie: metoda, czy funkcja?
Muszę utrzymać ten rozdział, inaczej kod w innych miejscach będzie mniej czytelny.

Funkcja jeśli: sprawdzenie wartości tekstury, nie zmienia niczego w kafelku. Przykład: isDirtRoad()
Metoda jeśli: zmienia stan wewnętrzeny kafelka lub wchodzimy z nim w interakcję. Przykład:
*/

// =============
// ↓↓↓FUNKCJE↓↓↓
// =============

// Sprawdza, czy tekstura jest kapliczką leczącą (świętym miejscem).
func isHealingShire(tileTexID uint16) bool {
	return tileTexID == spriteEffectHeal00 || tileTexID == spriteEffectHeal01
}

// Sprawdza, czy tekstura jest punktem zwycięstwa lub kapliczką.
func isSpecialTile(tileTexID uint16) bool {
	return tileTexID == spriteEffectHeal00 || tileTexID == spriteEffectHeal01 || tileTexID == spriteVictoryPoint
}

// Sprawdza, czy tekstura jest drogą.
func isDirtRoad(tileTexID uint16) bool {
	return tileTexID >= spriteRoadStart && tileTexID <= spriteRoadEnd
}

// Sprawdza, czy tekstura jest ukończonym mostem. Kierunek nie ma znaczenia.
// Nie ma znaczenia, czy grasz postawił.
// @todo: sprawdź, czemu nieukończony most jest traktowany inaczej - 21.04.2026.
func isCompletedBridge(tileTexID uint16) bool {
	return tileTexID >= spriteBridgeStart && tileTexID <= spriteBridgeEnd
}

// Sprawdza, czy tekstura jest drogą lub mostem.
func isPath(tileTexID uint16) bool {
	return isDirtRoad(tileTexID) || isCompletedBridge(tileTexID)
}

// Sprawdza, czy tekstura to zgliszcza budynku.
func isRuin(tileTexID uint16) bool {
	return tileTexID >= spriteRuinStart && tileTexID <= spriteRuinEnd
}

// Sprawdza, czy tekstura jest palisadą. Stan ukończenia nie ma znaczenia.
func isPalisade(tileTexID uint16) bool {
	return tileTexID >= spritePalisadeStart && tileTexID <= spritePalisadeEnd
}

// Sprawdza, czy tekstura jest wodą. Stan animacji nie ma znaczenia.
func isWaterTileOnly(tileTexID uint16) bool {
	return tileTexID >= spriteWaterStart && tileTexID <= spriteWaterEnd
}

// Sprawdza, czy tekstura jest wodą lub ukończonym mostem.
func isWaterOrBridgeForMasking(tileTexID uint16) bool {
	return isWaterTileOnly(tileTexID) || isCompletedBridge(tileTexID) || tileTexID == spriteBridgeConstruction
}

// Sprawdza, czy tekstura jest „lądowa”, czyli nie woda, nie mosty.
func isLandOrOther(tileTexID uint16) bool {
	return !isWaterTileOnly(tileTexID) && !isCompletedBridge(tileTexID) && tileTexID != spriteBridgeConstruction
}

// Sprawdza, czy tekstura jest to teren pod budynkami lub placem budowy.
func isBuildingTerrain(tileTexID uint16) bool {
	return (tileTexID >= spriteConstructionStart && tileTexID <= spriteConstructionEnd) ||
		(tileTexID >= spriteBuildingMainBase && tileTexID <= spriteBuildingEnd)
}

// Sprawdza, czy tekstura jest to sucha ziemia, często stosowana przy mostach.
func isDryEarth(tileTexID uint16) bool {
	return tileTexID >= spriteDryEarth01 && tileTexID <= spriteDryEarth03
}

// Sprawdza, czy tekstura jest to trawa. Stopień wyżarcia przez krowy, czy spopielenia nie ma znaczenia.
func isGrass(tileTexID uint16) bool {
	return tileTexID >= spriteGrassStart && tileTexID <= spriteGrassEnd
}

// Sprawdza, czy tekstura jest ozdoba terenu, jak ołtarz, drogowskaz, studnia.
func isGadget(tileTexID uint16) bool {
	return tileTexID >= spriteGadgetStart && tileTexID <= spriteGadgetEnd
}

// Sprawdza, czy tekstura jest trawa lub ozdoba.
func isGrassOrGadget(tileTexID uint16) bool {
	return isGrass(tileTexID) || isGadget(tileTexID)
}

// Sprawdza, czy tekstura jest kamieniem nieprzechodnim. W ozdobach (gadget) są przechodnie.
func isRockNonWalkable(tileTexID uint16) bool {
	return tileTexID >= spriteRockStart && tileTexID < spriteRockEnd
}

// Sprawdza, czy tekstura jest pień drzewa, w tym suchego. Nie mylić z pieńkiem, który jest ozdobą.
func isTreeStump(tileTexID uint16) bool {
	return tileTexID >= spriteTreeStumpStart && tileTexID <= spriteTreeStumpEnd
}

// Sprawdza, czy tekstura jest drzewem. W tym obalonym.
func isTree(tileTexID uint16) bool {
	return tileTexID >= spriteTreeStumpStart && tileTexID <= spriteTreeTopEnd ||
		tileTexID >= spriteTreeBurntStump00 && tileTexID <= spriteTreeFallingBurnt02
}

// =============
/// ↓↓↓METODY↓↓↓
// =============

// Płomienie

func (t *tile) setOnFire(fireSize uint16, bs *battleState) {
	t.IsBurning = true
	t.BurnElapsed = fireSize - bigBurn

	if !isTreeStump(t.TextureID) {
		t.IsAsh = true
	}

	bs.BurningTilesList = append(bs.BurningTilesList, t)
}

func (t *tile) processNormalFire() {
	// Gromadzenie się popiołu
	t.AshIntensity += ashAccumulationRate

	if t.AshIntensity > 1.0 {
		t.AshIntensity = 1.0
	}

	// Właściwe płonięcie
	t.BurnElapsed++

	var currentFireSprite uint16

	switch {
	case t.BurnElapsed < bigBurn:
		{
			currentFireSprite = spriteFire00
		}
	case t.BurnElapsed < midBurn:
		{
			currentFireSprite = spriteFire04
		}
	case t.BurnElapsed < minBurn:
		{
			currentFireSprite = spriteFire08
		}
	default:
		t.IsBurning = false
		t.AshAge = 0
		t.AshProcessState = ashDecaying

		return
	}

	t.BurnOverlayID = currentFireSprite
}

func (t *tile) processAshDecay() {
	// Płonące kafelki nie tracą popiołu
	if t.IsBurning {
		return
	}

	// Jeśli nie ma popiołu lub są śladowe ilości to wychodzimy
	if !t.IsAsh || t.AshIntensity < 0.01 {
		return
	}

	switch t.AshProcessState {
	case ashDecaying:
		t.AshIntensity *= 1.0 - ashDecayRate

		if t.AshAge >= totalAshLifetime {
			t.IsAsh = false
			t.AshIntensity = 0.0
			t.AshProcessState = ashFinished
		}
	case ashFinished:
		return

	// domyślne niepowinno nigdy wystąpić
	default:
		return
	}

	t.CurrentAshAlpha = t.AshIntensity
	t.AshAge++
}

func (t *tile) processTreeFire() {
	// Właściwe płonięcie
	t.BurnElapsed++

	var currentFireSprite uint16

	switch {
	case t.BurnElapsed < bigBurn:
		{
			currentFireSprite = spriteFire00
		}
	case t.BurnElapsed < midBurn:
		{
			currentFireSprite = spriteFire04
		}
	case t.BurnElapsed < minBurn:
		{
			currentFireSprite = spriteFire08
		}
	default:
		t.IsBurning = false
		t.processBurntTree()

		return
	}

	t.BurnOverlayID = currentFireSprite
}

func (t *tile) processBurntTree() {
	// Ustalamy tekstury odpowiadające spalonym drzewom.
	if t.TextureID < spriteTreeStump03 {
		t.TextureID = spriteTreeBurntStump00
	} else {
		t.TextureID = spriteTreeBurntStump01
	}

	// Obalamy spalone drzewo
	t.treeFall()
}

func (t *tile) applyFireDamage(bs *battleState) {
	if !t.IsBurning {
		return
	}

	damage := burnDamage

	if t.Unit != nil && t.Unit.Exists {
		t.Unit.takeDamage(damage, bs)
	}

	if t.Building != nil && t.Building.Exists {
		t.Building.takeDamage(damage)
	}
}

func (t *tile) treeFall() {
	// Po chwili czekania dajemy spriteTreeFalling.
	// Później spriteTreeFallen.
	// @todo: problem z tym, że nie wiem, czy to pojedyncze duszki
	// @todo: potrzebuję dodać obrażenia od spadającej korony
	// @todo: jeśli na lewo od drzewa jest suche drzewo to ono też ma być obalone
	// @reminder: potrzebuję dostępu do współrzędnych! bez tego nie mam wpływu na sąsiedni kafelek
	// kafelek staje się przechodni.
}
