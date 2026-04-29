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

// Sprawdza, czy tekstura jest drzewem. Drzewa, które nie zostały obalone. Rozróżnienie ważne
// ze względu na właściwość t.isWalkable.
func isTree(tileTexID uint16) bool {
	return tileTexID >= spriteTreeStumpStart && tileTexID <= spriteTreeTopEnd || // drzewa żywe oraz suche
		tileTexID >= spriteTreeBurntStump00 && tileTexID <= spriteTreeBurntTop01 // spalone drzewa
}

// =============
/// ↓↓↓METODY↓↓↓
// =============

// @reminder: dodaję metody do sprawdzenia, czy drzewa. Wydaje się to być potrzebne do
// rozdzielenia stanu od wyglądu ORAZ tego co używamy w grze od tego co potrzebne przy ładowaniu
// map w level.go - 27.04.2026

func (t *tile) isTree() bool {
	return t.treeState != noTree
}

func (t *tile) isStandingTree() bool {
	return t.treeState == treeStraight
}

func (t *tile) isFallingTree() bool {
	return t.treeState == treeFalling
}

func (t *tile) isFallenTree() bool {
	return t.treeState == treeFell
}

func (t *tile) isBurntTree() bool {
	return t.IsBurnt && t.isTree()
}

// Płomienie

// @todo: odsiej kafelki, które nie powinny się palić: np. woda
func (t *tile) setOnFire(fireSize uint16, bs *battleState) {
	t.IsBurning = true
	t.BurnElapsed = fireSize - bigBurn

	if !isTreeStump(t.TextureID) {
		t.hasAsh = true
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
	if !t.hasAsh || t.AshIntensity < 0.01 {
		return
	}

	switch t.AshProcessState {
	case ashDecaying:
		t.AshIntensity *= 1.0 - ashDecayRate

		if t.AshAge >= totalAshLifetime {
			t.hasAsh = false
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

// @todo: z moich obserwacji wynika, że drzewo podpalone odpryskiem szybciej się spala.
// jest to niewłaściwe zachowanie, powinno albo spalić się w tym samym czasie lub później
// ze względu na mniejszy początkowy ogień!
func (t *tile) processTreeFire(bs *battleState) {
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
		// czy mogę tego jeszcze nie zmieniać na FAŁSZ?
		t.IsBurning = false
		t.processBurntTree(bs)

		return
	}

	t.BurnOverlayID = currentFireSprite
}

func (t *tile) processBurntTree(bs *battleState) {
	// Ustalamy tekstury odpowiadające spalonym drzewom.
	if t.TextureID < spriteTreeStump03 {
		t.TextureID = spriteTreeBurntStump00
	} else {
		t.TextureID = spriteTreeBurntStump01
	}

	// Obalamy spalone drzewo
	t.IsBurnt = true
	t.treeFall(bs)
}

// Odpowiada za zadanie obrażeń jednostce lub budynkowi, który się znajduje na danym kafelku.
func (t *tile) applyFireDamage(bs *battleState) {
	if !t.IsBurning {
		return
	}

	if t.Unit != nil && t.Unit.Exists {
		t.Unit.takeDamage(burnDamage, bs)
	}

	if t.Building != nil && t.Building.Exists {
		t.Building.takeDamage(burnDamage)
	}
}

// Odpowiada za zadanie obrażeń jednostce lub budynkowi, który się znajduje na danym kafelku.
func (t *tile) applyFallingTreeDamage(bs *battleState) {
	if t.Unit != nil && t.Unit.Exists {
		t.Unit.takeDamage(fallingTreeDamage, bs)
	}

	if t.Building != nil && t.Building.Exists {
		t.Building.takeDamage(fallingTreeDamage)
	}
}

func (t *tile) accumulateTreeCuts(bs *battleState) {
	t.treeCuts++

	if t.treeCuts >= strikesToCutTree {
		t.treeFall(bs)
	}
}

// Odpowiada za rozpoczęcie całego procesu upadania drzewa.
func (t *tile) treeFall(bs *battleState) {
	// Drzewa, które już upadają nie są obsługiwane!
	if t.treeState != treeStraight {
		return
	}

	// Nie pozwalamy na duplikaty!
	for _, existing := range bs.FallingTreesList {
		if existing == t {
			return
		}
	}

	// 1. Ustawiamy stopień upadku drzewa
	// Za zarządzanie teksturami odpowiada fallingTreeEffect()
	t.treeState = treeStraight
	t.treeCuts = 0

	// 2. Dodajemy kafelkek do listy obsługiwanej centralnie
	bs.FallingTreesList = append(bs.FallingTreesList, t)
}

func (t *tile) ghost(ghostSprite uint16, ghostDamage uint16, bs *battleState) {
	// 1. Ustawiamy wszystkie parametry dla kafelka z duchem
	t.GhostEffect = true
	t.GhostEffectCounter = 40
	t.GhostOverlayID = ghostSprite
	t.GhostDamage = ghostDamage

	// 2. Dodajemy ducha do listy
	// w effects.go będzie osobna funkcja z logiką efektu ducha
	bs.GhostsList = append(bs.GhostsList, t)
}
