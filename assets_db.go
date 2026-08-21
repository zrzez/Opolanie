package main

import "log"

/*
Plik jest pełen zaczarodziejskich liczb i jest tak powinno być. Wyciszam linter, żeby nie zawracał mi głowy o to.
*/

// Zestawienie wszystkich składowych.
var spriteRegistry [maxSpriteID]spriteDef

var idRegistry = map[string]uint16{
	// Trawa
	"SPRITE_GRASS_START": spriteGrassStart, // 2
	"SPRITE_GRASS_END":   spriteGrassEnd,   // 17
	"SPRITE_GRASS_00":    spriteGrass00,
	"SPRITE_GRASS_06":    spriteGrass06, // Stara 8-ka

	// Woda
	"SPRITE_WATER_START":  spriteWaterStart,  // 74
	"SPRITE_WATER_MIDDLE": spriteWaterMiddle, // 82
	"SPRITE_WATER_END":    spriteWaterEnd,    // 112

	// Skały
	"SPRITE_ROCK_START": spriteRockStart,
	"SPRITE_ROCK_00":    spriteRock00,
	"SPRITE_ROCK_01":    spriteRock01,
	"SPRITE_ROCK_02":    spriteRock02,
	"SPRITE_ROCK_03":    spriteRock03,
	"SPRITE_ROCK_04":    spriteRock04,
	"SPRITE_ROCK_05":    spriteRock05,
	"SPRITE_ROCK_06":    spriteRock06,
	"SPRITE_ROCK_07":    spriteRock07,
	"SPRITE_ROCK_08":    spriteRock08,
	"SPRITE_ROCK_09":    spriteRock09,
	"SPRITE_ROCK_10":    spriteRock10,
	"SPRITE_ROCK_11":    spriteRock11,
	"SPRITE_ROCK_12":    spriteRock12, // Przechodnie

	// Drogi
	"SPRITE_ROAD_START": spriteRoadStart,

	// Mosty
	"SPRITE_BRIDGE_START": spriteBridgeStart,
	"SPRITE_BRIDGE_00":    spriteBridge00,
	"SPRITE_BRIDGE_01":    spriteBridge01,
	"SPRITE_BRIDGE_02":    spriteBridge02,
	"SPRITE_BRIDGE_03":    spriteBridge03,
	"SPRITE_BRIDGE_04":    spriteBridge04,
	"SPRITE_BRIDGE_05":    spriteBridge05,
	"SPRITE_BRIDGE_06":    spriteBridge06,
	"SPRITE_BRIDGE_07":    spriteBridge07,

	// Sucha ziemia
	"SPRITE_DRY_EARTH_1": spriteDryEarth01,
	"SPRITE_DRY_EARTH_2": spriteDryEarth02,
	"SPRITE_DRY_EARTH_3": spriteDryEarth03,

	// Pierdoły (gadgets)
	"SPRITE_GADGET_START": spriteGadgetStart,
	"SPRITE_GADGET_00":    spriteGadget00,
	"SPRITE_GADGET_01":    spriteGadget01,
	"SPRITE_GADGET_02":    spriteGadget02,
	"SPRITE_GADGET_03":    spriteGadget03,
	"SPRITE_GADGET_04":    spriteGadget04,
	"SPRITE_GADGET_05":    spriteGadget05,
	"SPRITE_GADGET_06":    spriteGadget06,
	"SPRITE_GADGET_07":    spriteGadget07,
	"SPRITE_GADGET_08":    spriteGadget08,
	"SPRITE_GADGET_09":    spriteGadget09,
	"SPRITE_GADGET_10":    spriteGadget10,
	"SPRITE_GADGET_11":    spriteGadget11,
	"SPRITE_GADGET_12":    spriteGadget12,
	"SPRITE_GADGET_13":    spriteGadget13,
	"SPRITE_GADGET_14":    spriteGadget14,
	"SPRITE_GADGET_15":    spriteGadget15,

	// Drzewa
	"SPRITE_TREE_STUMP_START": spriteTreeStumpStart,
	"SPRITE_TREE_STUMP_00":    spriteTreeStump00,
	"SPRITE_TREE_STUMP_01":    spriteTreeStump01,
	"SPRITE_TREE_STUMP_02":    spriteTreeStump02,
	"SPRITE_TREE_STUMP_03":    spriteTreeStump03,
	"SPRITE_TREE_STUMP_04":    spriteTreeStump04,
	"SPRITE_TREE_STUMP_05":    spriteTreeStump05,
	"SPRITE_TREE_STUMP_06":    spriteDryTreeStump00,

	// Upadające drzewa
	"SPRITE_DRY_TREE_FALLEN_STUMP_02":     spriteDryFallenTreeStump,
	"SPRITE_DRY_TREE_FALLEN_TOP_02":       spriteDryFallenTreeTop,
	"SPRITE_DRY_TREE_FALLING_STUMP_01":    spriteDryFallingStump,
	"SPRITE_DRY_TREE_FALLING_TOP_01":      spriteDryFallingTreeTop,
	"SPRITE_DRY_TREE_LEANING_TOP_LEFT_00": spriteDryLeaningTreeCrownLeft,
	"SPRITE_DRY_TREE_LEANING_TOP_00":      spriteDryLeaningTreeTop,
	"SPRITE_DRY_TREE_LEANING_STUMP_00":    spriteDryLeaningTreeStump,

	// Specjalne
	"SPRITE_PALISADE_START":      spritePalisadeStart,
	"SPRITE_EFFECT_HEAL_00":      spriteEffectHeal00,
	"SPRITE_EFFECT_TRANSFORM_00": spriteEffectTransform00,
	"SPRITE_VICTORY_POINT":       spriteVictoryPoint,
	"SPRITE_EFFECT_SKELETON_00":  spriteEffectskeleton00,
	"SPRITE_EFFECT_SKELETON_01":  spriteEffectskeleton01,
	"SPRITE_EFFECT_SKELETON_02":  spriteEffectskeleton02,

	// Ogień
	"SPRITE_FIRE_00": spriteFire00,
	"SPRITE_FIRE_01": spriteFire01,
	"SPRITE_FIRE_02": spriteFire02,
	"SPRITE_FIRE_03": spriteFire03,
	"SPRITE_FIRE_04": spriteFire04,
	"SPRITE_FIRE_05": spriteFire05,
	"SPRITE_FIRE_06": spriteFire06,
	"SPRITE_FIRE_07": spriteFire07,
	"SPRITE_FIRE_08": spriteFire08,
	"SPRITE_FIRE_09": spriteFire09,
	"SPRITE_FIRE_10": spriteFire10,
	"SPRITE_FIRE_11": spriteFire11,
	"SPRITE_FIRE_12": spriteAsh00,
	"SPRITE_FIRE_13": spriteAsh01,

	"SPRITE_ARROW":     spriteMissileArrowUp,
	"SPRITE_BOLT":      spriteMissileBoltUp,
	"SPRITE_FIRE":      spriteMissileFireUp,
	"SPRITE_LIGHTNING": spriteMissileLightningUp,
	"SPRITE_SPEAR":     spriteMissileSpearUp,
	"SPRITE_GHOST":     spriteMissileGhostUp,
}

const (
	specialMainMenu = iota
	specialCampaignSelect
	specialMapPoland
	specialProvinces
	specialBriefing
	specialCount // ← TO MUSI BYĆ OSTATNIE!
)

var specialAssetsDB = [specialCount]rawAssetDef{
	specialMainMenu:       {1, 16, 1},
	specialCampaignSelect: {10, 24, 10},
	specialMapPoland:      {6, 21, 6},
	specialProvinces:      {12, 27, 0},
	specialBriefing:       {8, 23, 8},
}

// Zwraca ID dla nazwy. Domyślnie trawa.
func getID(spriteName string) uint16 {
	if val, ok := idRegistry[spriteName]; ok {
		return val
	}
	// Bezpiecznik dla pustych łańcuchów, które mogły zostać z przenosin jako "" dla starych ID lub błędnych nazw.
	if spriteName != "" {
		log.Printf("OSTRZEŻENIE: Brak ID dla '%s'. Używam trawy.", spriteName)
	}

	return spriteGrassStart
}

func init() {
	initTerrainSprites()
	initUISprites()
	initManualUnitSprites()
	initManualCorpseSprites()
	initBuildingSprites()
	initProjectileSprites()
}

/*
Kod poniżej wycina i przygotowuje do rysowania „duszki” (ang. sprite).
Całość jest bardzo prosta:
- mamy duży zbiór obrazków, zwany atlasem.
- wycinamy z niego mały prostokąt („duszek”) z wyglądem tego co chcemy (np. studnia)

Aby móc to zrobić potrzebujemy kilku informacji:
- atlasu z którego chcemy coś wyciąć; dlatego potrzebujemy atlasUI.
- nazwę wycinanego duszka; dlatego potrzebujemy spriteID.
- punkt początkowy do wycięcia duszka; dlatego potrzebujemy cropX, cropY.
- wysokość i szerokość duszka; dlatego potrzebujemy cropWidth, cropHeight, które zazwyczaj mają wymiary 16x14.
- poprawki przy rysowaniu, przydaje się np. przy drzewach; dlatego potrzebujemy drawOffsetX, drawOffsetY.
- czy chcemy lustrzane odbicie duszka. Szczególnie przydatne przy duszkach ponieważ dla oszczędności miejsca wiele
duszków to właśnie lustrzane odbicie; dlatego potrzebujemy flipX.

W przypadku jednostek sprawa się komplikuje z kilku powodów:
- niektóre jednostki, np. drwal (unitAxeman) walczą wręcz, a to oznacza, że duszki mają różne wymiary. Jeśli stoją, to
są to zwyczajne 16x14. Jednak przy ataku, który wizualnie zachodzi na atakowany kafelkek, wymiary te ulegają zmianie.
Pole na planszy zajmowane przez atakowaną wrogą jednostkę lub budynek też musi mieć częściowo duszka atakującego.
- jednostki są w różnych atlasach
- jednostki dodane w wersji na CD, tj. pastuch, kusznik, mag mają swój własny rozkład duszków.
Przez to potrzebują osobnego podejścia do załadowania wszystkich potrzebnych duszków.
- wszystkie duszki są przygotowane w barwach gracza (czerwień) i przy ładowaniu mapy przebarwiane. Odbywa się to
w innym miejscu.
*/

func initTerrainSprites() {
	// Pomagier do nakładki (16x14)
	setUI := func(spriteID, cropX, cropY uint16) {
		if spriteID < maxSpriteID {
			spriteRegistry[spriteID] = spriteDef{
				atlasID: atlasUI,
				cropX:   cropX, cropY: cropY,
				cropWidth: tileWidth, cropHeight: tileHeight,
				offX: 0, offY: 0,
				flipX: false,
			}
		}
	}

	// Pomagier do Units1 (16x14)
	setUnit1 := func(spriteID, cropY uint16) {
		if spriteID < maxSpriteID {
			spriteRegistry[spriteID] = spriteDef{
				atlasID: atlasUnits1,
				cropX:   303, cropY: cropY, //nolint:mnd
				cropWidth: tileWidth, cropHeight: tileHeight,
				offX: 0, offY: 0,
				flipX: false,
			}
		}
	}

	// Pomagier do gabarytów
	setSpecial := func(spriteID, cropX, cropY uint16, cropWidth, cropHeight uint8) {
		if spriteID < maxSpriteID {
			spriteRegistry[spriteID] = spriteDef{
				atlasID: atlasUI,
				cropX:   cropX, cropY: cropY,
				cropWidth: cropWidth, cropHeight: cropHeight,
				offX: 0, offY: 0,
				flipX: false,
			}
		}
	}

	setFallingTree := func(spriteID, cropX, cropY uint16, offX, offY int8) {
		if spriteID < maxSpriteID {
			spriteRegistry[spriteID] = spriteDef{
				atlasID: atlasUI,
				cropX:   cropX, cropY: cropY,
				cropWidth: tileWidth, cropHeight: tileHeight,
				offX: offX, offY: offY,
				flipX: false,
			}
		}
	}

	// 1. Trawa
	setUI(spriteGrassStubbed, 11, 36) //nolint:mnd
	setUI(spriteGrassGrazed, 27, 36)  //nolint:mnd

	// A. Tekstury z atlasu UI
	setUI(spriteGrass00, 43, 36)   //nolint:mnd
	setUI(spriteGrass01, 139, 134) //nolint:mnd

	// B. Tekstury z atlasu Units1
	setUnit1(spriteEffectHeal00, 0)  //nolint:nolintlint,mnd
	setUnit1(spriteEffectHeal01, 14) //nolint:mnd
	setUnit1(spriteGrass02, 28)      //nolint:mnd
	setUnit1(spriteGrass03, 42)      //nolint:mnd
	setUnit1(spriteGrass04, 56)      //nolint:mnd
	setUnit1(spriteGrass05, 70)      //nolint:mnd
	setUnit1(spriteGrass06, 84)      //nolint:mnd
	setUnit1(spriteGrass07, 98)      //nolint:mnd

	// Ba. Magiczna tarcza, próba
	// 303, 112 pierwsza magiczna tarcza, później miarowo w dół
	setUnit1(spriteMagicShield00, 112) //nolint:mnd
	setUnit1(spriteMagicShield01, 126) //nolint:mnd
	setUnit1(spriteMagicShield02, 140) //nolint:mnd
	setUnit1(spriteMagicShield03, 154) //nolint:mnd

	// C. Mieszanie
	spriteRegistry[spriteGrass08] = spriteRegistry[spriteGrass02]
	spriteRegistry[spriteGrass09] = spriteRegistry[spriteGrass00]
	spriteRegistry[spriteGrass10] = spriteRegistry[spriteGrass03]
	spriteRegistry[spriteGrass11] = spriteRegistry[spriteGrass01]
	spriteRegistry[spriteGrass12] = spriteRegistry[spriteGrass04]
	spriteRegistry[spriteGrass13] = spriteRegistry[spriteGrass05]
	spriteRegistry[spriteGrass14] = spriteRegistry[spriteGrass06]
	spriteRegistry[spriteGrass15] = spriteRegistry[spriteGrass07]

	// 2. Reszta terenu

	// Skały
	setUI(spriteRockEnd, 11+(8*16), 92) //nolint:mnd

	currentRockID := spriteRockStart

	for i := range 13 {
		if i == 8 { //nolint:mnd
			continue
		}

		setUI(currentRockID, uint16(11+(i*16)), 92) //nolint:mnd
		currentRockID++
	}

	// Sucha ziemia
	setDryEarth := func(spriteID, cropX, cropY uint16) {
		if spriteID < maxSpriteID {
			spriteRegistry[spriteID] = spriteDef{
				atlasID: atlasUI,
				cropX:   cropX, cropY: cropY,
				cropWidth: 22, cropHeight: 18, //nolint:mnd
				offX: -3, offY: -2, //nolint:nolintlint
				flipX: false,
			}
		}
	}

	setDryEarth(spriteDryEarth01, 235, 134) //nolint:mnd
	setDryEarth(spriteDryEarth02, 235, 153) //nolint:mnd
	setDryEarth(spriteDryEarth03, 235, 172) //nolint:mnd

	// Drogi
	for i := range uint16(5) { //nolint:mnd
		setUI(spriteRoadStart+i, 59+(i*16), 134) //nolint:mnd
	}

	setUI(spriteRoadStart+5, 139, 134) //nolint:mnd

	for i := range uint16(4) { //nolint:mnd
		setUI(spriteRoadStart+6+i, 155+(i*16), 134) //nolint:mnd
	}

	for i := range uint16(11) { //nolint:mnd
		setUI(spriteRoadStart+10+i, 11+(i*16), 148) //nolint:mnd
	}

	// Mosty
	for i := uint16(0); i <= (spriteBridgeEnd - spriteBridgeStart); i++ {
		setUI(spriteBridgeStart+i, 11+(i*16), 162) //nolint:mnd
	}

	// Woda
	for i := range uint16(13) { //nolint:mnd
		setUI(spriteWaterStart+i, 11+(i*16), 50) //nolint:mnd
	}

	for i := range uint16(13) { //nolint:mnd
		setUI(spriteWaterStart+13+i, 11+(i*16), 64) //nolint:mnd
	}

	for i := range uint16(13) { //nolint:mnd
		setUI(spriteWaterStart+26+i, 11+(i*16), 78) //nolint:mnd
	}

	// Drzewa
	for i := range uint16(7) { //nolint:mnd
		setSpecial(spriteTreeStumpStart+i, 11+(i*32), 120, 32, 14) //nolint:mnd
	}

	for i := range uint16(7) { //nolint:mnd
		setSpecial(spriteTreeTopStart+i, 11+(i*32), 106, 32, 14) //nolint:mnd
	}

	// Palisady
	for i := uint16(0); i <= (spritePalisadeEnd - spritePalisadeStart); i++ {
		spriteID := spritePalisadeStart + i
		spriteRegistry[spriteID] = spriteDef{
			atlasID:   atlasUnits1,
			cropX:     287,    //nolint:mnd
			cropY:     i * 14, //nolint:mnd
			cropWidth: tileWidth, cropHeight: tileHeight,
			offX: 0, offY: 0,
			flipX: false,
		}
	}

	// Pierdoły

	for i := range uint16(11) { //nolint:mnd
		setUI(spriteGadgetStart+i, 59+(i*16), 36) //nolint:mnd
	}

	for i := range uint16(3) { //nolint:mnd
		setUI(spriteGadgetStart+11+i, 187+(i*16), 148) //nolint:mnd
	}

	for i := range uint16(6) { //nolint:mnd
		setUI(spriteGadgetStart+14+i, 139+(i*16), 162) //nolint:mnd
	}

	// Rany
	setUI(spriteEffectHit00, 203, 8) //nolint:mnd
	setUI(spriteEffectHit01, 219, 8) //nolint:mnd

	// Spalone drzewa
	// @reminder: to może się nie zapisywać w przyszłości przez level.go, bo nie jest dodane w idRegistry 22.04.2026
	// przy rysowaniu działa poprawnie.
	setUI(spriteTreeBurntStump00, 187, 162) //nolint:mnd
	setUI(spriteTreeBurntStump01, 219, 162) //nolint:mnd
	setUI(spriteTreeBurntTop00, 171, 162)   //nolint:mnd
	setUI(spriteTreeBurntTop01, 203, 162)   //nolint:mnd

	// Upadające suche drzewa
	// Obalone
	setFallingTree(spriteDryFallenTreeTop, 235, 22, -16, 0) //nolint:mnd
	setFallingTree(spriteDryFallenTreeStump, 251, 22, 0, 0) //nolint:mnd
	// Upadające
	setFallingTree(spriteDryFallingTreeTop, 235, 36, -16, 0) //nolint:mnd
	setFallingTree(spriteDryFallingStump, 251, 36, 0, 0)     //nolint:mnd
	// Przechylające się
	setFallingTree(spriteDryLeaningTreeCrownLeft, 235, 50, -16, -14) //nolint:mnd
	setFallingTree(spriteDryLeaningTreeTop, 251, 50, 0, -14)         //nolint:mnd
	setFallingTree(spriteDryLeaningTreeStump, 251, 64, 0, 0)         //nolint:mnd

	// Upadające zwęglone drzewa
	// Obalone
	setFallingTree(spriteBurntFallenTreeCrown, 235, 78, -16, 0) //nolint:mnd
	setFallingTree(spriteBurntFallenTreeStump, 251, 78, 0, 0)   //nolint:mnd
	// Upadające
	setFallingTree(spriteBurntFallingTreeCrown, 235, 92, -16, 0) //nolint:mnd
	setFallingTree(spriteBurntFallingTreeStump, 251, 92, 0, 0)   //nolint:mnd
	// Przechylające się
	setFallingTree(spriteBurntLeaningTreeCrownLeft, 235, 106, -16, -14) //nolint:mnd
	setFallingTree(spriteBurntLeaningTreeCrown, 251, 106, 0, -14)       //nolint:mnd
	setFallingTree(spriteBurntLeaningTreeStump, 251, 120, 0, 0)         //nolint:mnd
}

// Nakładka.
func initUISprites() {
	setUI := func(spriteID, cropX, cropY uint16) {
		if spriteID < maxSpriteID {
			spriteRegistry[spriteID] = spriteDef{
				atlasID: atlasUI,
				cropX:   cropX, cropY: cropY,
				cropWidth: tileWidth, cropHeight: tileHeight,
				offX: 0, offY: 0,
				flipX: false,
			}
		}
	}

	setCenterUI := func(spriteID, cropX uint16) {
		if spriteID < maxSpriteID {
			spriteRegistry[spriteID] = spriteDef{
				atlasID: atlasUI,
				cropX:   cropX, cropY: 8, //nolint:mnd
				cropWidth: tileWidth, cropHeight: tileHeight,
				offX: -8, offY: -7,
				flipX: false,
			}
		}
	}

	setUIRepairBtn := func(spriteID, cropX, cropY uint16, cropHeight uint8) {
		if spriteID < maxSpriteID {
			spriteRegistry[spriteID] = spriteDef{
				atlasID: atlasUI,
				cropX:   cropX, cropY: cropY,
				cropWidth: tileWidth, cropHeight: cropHeight,
				offX: 0, offY: 0,
				flipX: false,
			}
		}
	}

	// Kursory
	setCenterUI(spriteCursorCrossWhite, 27)     //nolint:mnd
	setCenterUI(spriteCursorCrossRed, 43)       //nolint:mnd
	setCenterUI(spriteCursorSmallWhite, 59)     //nolint:mnd
	setCenterUI(spriteCursorFrameRed, 75)       //nolint:mnd
	setCenterUI(spriteCursorCrossMedRed, 91)    //nolint:mnd
	setCenterUI(spriteCursorCrossMedWhite, 107) //nolint:mnd
	setCenterUI(spriteCursorArrowUp, 123)       //nolint:mnd
	setCenterUI(spriteCursorArrowDown, 139)     //nolint:mnd
	setCenterUI(spriteCursorArrowLeft, 155)     //nolint:mnd
	setCenterUI(spriteCursorArrowRight, 171)    //nolint:mnd
	setCenterUI(spriteCursorStop, 187)          //nolint:mnd
	setCenterUI(spriteCursorFrameWhite, 75)     //nolint:mnd

	setUI(spriteCursorDefaultBig, 11, 8)    //nolint:mnd
	setUI(spriteCursorDefaultSmall, 91, 22) //nolint:mnd
	setUI(spriteCursorPointer, 107, 22)     //nolint:mnd

	// Przyciski
	setUI(spriteBtnBuildPalisade, 251, 8)   //nolint:mnd
	setUI(spriteBtnShield, 11, 22)          //nolint:mnd
	setUI(spriteBtnBuildBarn, 155, 21)      //nolint:mnd
	setUI(spriteBtnBuildBarracks, 171, 21)  //nolint:mnd
	setUI(spriteBtnBuildTemple, 187, 21)    //nolint:mnd
	setUI(spriteBtnBuildBarracks2, 203, 21) //nolint:mnd
	setUI(spriteBtnBuildAcademy, 219, 21)   //nolint:mnd
	// setUI(SPRITE_BTN_, 123, 22, 16, 14) // btn_map @todo: pewnie można usunąć, bo nie używam
	setUIRepairBtn(spriteBtnRepair, 139, 22, 13) //nolint:mnd

	// Czary @todo: ogarnij potrójne ikonki dla przycisku. Gdzieś jest już gotowa funkcja do tego
	setUI(spriteBtnSpellVision, 235, 8) //nolint:mnd

	//goland:noinspection GoLinter
	spriteRegistry[spriteBtnSpellMagicShield] = spriteDef{
		atlasID: atlasUnits1, cropX: 303, cropY: 112, //nolint:mnd
		cropWidth: 16, cropHeight: 14, offX: 0, offY: 0, flipX: false, //nolint:mnd
	}

	// @todo: @reminder: ikonki te wymagają specjalnej funkcji do „podrójnego” rysowania.
	// Gromobicie.
	//goland:noinspection GoLinter
	spriteRegistry[spriteBtnSpellMagicLighting] = spriteDef{
		atlasID: atlasUnits2, cropX: 258, cropY: 69, //nolint:mnd
		cropWidth: 16, cropHeight: 14, offX: 0, offY: 0, flipX: false, //nolint:mnd
	}
	// Deszcz ognia.
	//goland:noinspection GoLinter
	spriteRegistry[spriteBtnSpellMagicFire] = spriteDef{
		atlasID: atlasUnits2, cropX: 255, cropY: 111, //nolint:mnd
		cropWidth: 16, cropHeight: 14, offX: 0, offY: 0, flipX: false,
	}

	// Zwłoki
	// @todo: kompletnie porąbane nazwy!
	setUI(spriteEffectskeleton00, 219, 50) //nolint:mnd
	setUI(spriteEffectskeleton01, 219, 64) //nolint:mnd
	setUI(spriteEffectskeleton02, 219, 78) //nolint:mnd

	// Rany
	setUI(spriteEffectHit00, 219, 8) //nolint:mnd
	setUI(spriteEffectHit01, 203, 8) //nolint:mnd

	// Ogień
	//nolint:mnd
	for i := range uint16(14) {
		spriteID := spriteFireStart + i
		if spriteID < maxSpriteID {
			spriteRegistry[spriteID] = spriteDef{
				atlasID:   atlasUI,
				cropX:     11 + (i * 16), //nolint:mnd
				cropY:     176,           //nolint:mnd
				cropWidth: tileWidth, cropHeight: tileHeight,
				offX: 0, offY: 0,
				flipX: false,
			}
		}
	}
}

// @todo: wiele tych liczb można zamienić stałymi sprite….
func initBuildingSprites() {
	assetID := atlasBuildings

	setBuilding := func(spriteID uint16, cropX, cropY uint16) {
		if spriteID < maxSpriteID {
			spriteRegistry[spriteID] = spriteDef{
				atlasID: assetID,
				cropX:   cropX, cropY: cropY,
				cropWidth: tileWidth, cropHeight: tileHeight,
				offX: 0, offY: 0,
				flipX: false,
			}
		}
	}

	for i := spriteConstructionStart; i <= spriteConstructionEnd; i++ {
		setBuilding(i, (i-127)*16, 168) //nolint:mnd
	}

	for i := spriteRuinStart; i <= spriteRuinEnd; i++ {
		setBuilding(i, (i-257)*16, 182) //nolint:mnd
	}

	// Budynek główny
	for i := spriteBuildingMainBase; i <= 155; i++ {
		setBuilding(i, (i-137)*16, 84) //nolint:mnd
	}
	// ID 156 Most //@todo: czemu to nie ma zdefiniowanego spriteID, ale działa?!
	setBuilding(156, 304, 84) //nolint:mnd

	// Obora
	for i := spriteBuildingBarnBase; i <= 175; i++ {
		setBuilding(i, (i-157)*16, 98) //nolint:mnd
	}
	// Chata drwali
	for i := spriteBuildingBarracksBase; i <= 195; i++ {
		setBuilding(i, (i-177)*16, 112) //nolint:mnd
	}
	// Świątynia
	for i := spriteBuildingTempleBase; i <= 215; i++ {
		setBuilding(i, (i-197)*16, 126) //nolint:mnd
	}
	// Chata wojów
	for i := spriteBuildingBarracks2Base; i <= 235; i++ {
		setBuilding(i, (i-217)*16, 140) //nolint:mnd
	}
	// Dwór
	for i := spriteBuildingAcademyBase; i <= spriteBuildingEnd; i++ {
		setBuilding(i, (i-237)*16, 154) //nolint:mnd
	}

	setBuilding(spriteBridgeConstruction, 304, 84) //nolint:mnd
}

func initProjectileSprites() {
	setProjectile := func(spiteID uint16, cropX, cropY uint16, flip bool) {
		if spiteID < maxSpriteID {
			spriteRegistry[spiteID] = spriteDef{
				atlasID: atlasUnits2,
				cropX:   cropX, cropY: cropY,
				cropWidth: tileWidth, cropHeight: tileHeight,
				offX: -8, offY: -7,
				flipX: flip,
			}
		}
	}

	// === STRZAŁA ===
	setProjectile(spriteMissileArrowUp, 256, 0, false)        //nolint:mnd
	setProjectile(spriteMissileArrowUpLeft, 240, 0, false)    //nolint:mnd
	setProjectile(spriteMissileArrowLeft, 240, 14, false)     //nolint:mnd
	setProjectile(spriteMissileArrowDownLeft, 240, 28, false) //nolint:mnd
	setProjectile(spriteMissileArrowDown, 256, 28, false)     //nolint:mnd
	setProjectile(spriteMissileArrowUpRight, 240, 0, true)    //nolint:mnd
	setProjectile(spriteMissileArrowRight, 240, 14, true)     //nolint:mnd
	setProjectile(spriteMissileArrowDownRight, 240, 28, true) //nolint:mnd

	// === PIORUN ===
	setProjectile(spriteMissileLightningUp, 258, 42, false)       //nolint:mnd
	setProjectile(spriteMissileLightningUpLeft, 242, 40, false)   //nolint:mnd
	setProjectile(spriteMissileLightningLeft, 242, 54, false)     //nolint:mnd
	setProjectile(spriteMissileLightningDownLeft, 242, 68, false) //nolint:mnd
	setProjectile(spriteMissileLightningDown, 258, 69, false)     //nolint:mnd
	setProjectile(spriteMissileLightningUpRight, 242, 40, true)   //nolint:mnd
	setProjectile(spriteMissileLightningRight, 242, 54, true)     //nolint:mnd
	setProjectile(spriteMissileLightningDownRight, 242, 68, true) //nolint:mnd

	// === OGIEŃ ===
	setProjectile(spriteMissileFireUp, 255, 84, false)        //nolint:mnd
	setProjectile(spriteMissileFireUpLeft, 241, 83, false)    //nolint:mnd
	setProjectile(spriteMissileFireLeft, 241, 97, false)      //nolint:mnd
	setProjectile(spriteMissileFireDownLeft, 241, 110, false) //nolint:mnd
	setProjectile(spriteMissileFireDown, 255, 111, false)     //nolint:mnd
	setProjectile(spriteMissileFireUpRight, 241, 83, true)    //nolint:mnd
	setProjectile(spriteMissileFireRight, 241, 97, true)      //nolint:mnd
	setProjectile(spriteMissileFireDownRight, 241, 110, true) //nolint:mnd

	// === WŁÓCZNIA ===
	setProjectile(spriteMissileSpearUp, 255, 126, false)       //nolint:mnd
	setProjectile(spriteMissileSpearUpLeft, 239, 125, false)   //nolint:mnd
	setProjectile(spriteMissileSpearLeft, 239, 139, false)     //nolint:mnd
	setProjectile(spriteMissileSpearDownLeft, 239, 154, false) //nolint:mnd
	setProjectile(spriteMissileSpearDown, 255, 154, false)     //nolint:mnd
	setProjectile(spriteMissileSpearUpRight, 239, 125, true)   //nolint:mnd
	setProjectile(spriteMissileSpearRight, 239, 139, true)     //nolint:mnd
	setProjectile(spriteMissileSpearDownRight, 239, 154, true) //nolint:mnd

	// === DUCH ===
	setProjectile(spriteMissileGhostUp, 288, 0, false)        //nolint:mnd
	setProjectile(spriteMissileGhostUpLeft, 272, 0, false)    //nolint:mnd
	setProjectile(spriteMissileGhostLeft, 272, 14, false)     //nolint:mnd
	setProjectile(spriteMissileGhostDownLeft, 272, 28, false) //nolint:mnd
	setProjectile(spriteMissileGhostDown, 288, 28, false)     //nolint:mnd
	setProjectile(spriteMissileGhostUpRight, 272, 0, true)    //nolint:mnd
	setProjectile(spriteMissileGhostRight, 272, 14, true)     //nolint:mnd
	setProjectile(spriteMissileGhostDownRight, 272, 28, true) //nolint:mnd

	// === BEŁT ===
	setProjectile(spriteMissileBoltUp, 288, 42, false)       //nolint:mnd
	setProjectile(spriteMissileBoltUpLeft, 272, 42, false)   //nolint:mnd
	setProjectile(spriteMissileBoltLeft, 272, 56, false)     //nolint:mnd
	setProjectile(spriteMissileBoltDownLeft, 272, 70, false) //nolint:mnd
	setProjectile(spriteMissileBoltDown, 288, 70, false)     //nolint:mnd
	setProjectile(spriteMissileBoltUpRight, 272, 42, true)   //nolint:mnd
	setProjectile(spriteMissileBoltRight, 272, 56, true)     //nolint:mnd
	setProjectile(spriteMissileBoltDownRight, 272, 70, true) //nolint:mnd

	// Efekt ducha. Musi być wydzielony ponieważ nie jest to zwykły pocisk.
	setGhost := func(spiteID uint16, cropX, cropY uint16) {
		if spiteID < maxSpriteID {
			spriteRegistry[spiteID] = spriteDef{
				atlasID: atlasUnits2,
				cropX:   cropX, cropY: cropY,
				cropWidth: tileWidth, cropHeight: tileHeight,
				offX: 0, offY: 0,
				flipX: false,
			}
		}
	}

	setGhost(spriteMissileGhostAttack, 288, 14) //nolint:mnd
}

// Mapowanie battleAtlasID → rawAssetDef {TopChunk, BotChunk, PaletteID}.
var atlasDefinitions = map[battleAtlasID]rawAssetDef{
	atlasUI:        {3, 18, 3}, // UI
	atlasUnits1:    {4, 19, 4}, // Jednostki
	atlasUnits2:    {5, 20, 3}, // Jednostki i pociski
	atlasBuildings: {7, 22, 3}, // Budynki
}

func initManualUnitSprites() {
	// Przygotowanie funkcji do załadowania „duszków” do atlasu units1.
	setManualUnits1 := func(spriteID, cropX, cropY uint16, cropWidth, cropHeight uint8, offX, offY int8, flipX bool) {
		if spriteID >= maxSpriteID {
			log.Printf("OSTRZEŻENIE: ręcznie wpisany duszek nie mieści się!")

			return
		}

		spriteRegistry[spriteID] = spriteDef{
			atlasID:    atlasUnits1,
			cropX:      cropX,
			cropY:      cropY,
			cropWidth:  cropWidth,
			cropHeight: cropHeight,
			offX:       offX,
			offY:       offY,
			flipX:      flipX,
		}
	}

	setManualUnits2 := func(spriteID, cropX, cropY uint16, cropWidth, cropHeight uint8, offX, offY int8, flipX bool) {
		if spriteID >= maxSpriteID {
			log.Printf("OSTRZEŻENIE: ręcznie wpisany duszek nie mieści się!")

			return
		}

		spriteRegistry[spriteID] = spriteDef{
			atlasID:    atlasUnits2,
			cropX:      cropX,
			cropY:      cropY,
			cropWidth:  cropWidth,
			cropHeight: cropHeight,
			offX:       offX,
			offY:       offY,
			flipX:      flipX,
		}
	}

	setManualUnits3 := func(spriteID, cropX, cropY uint16, cropWidth, cropHeight uint8, offX, offY int8, flipX bool) {
		if spriteID >= maxSpriteID {
			log.Printf("OSTRZEŻENIE: ręcznie wpisany duszek nie mieści się!")

			return
		}

		spriteRegistry[spriteID] = spriteDef{
			atlasID:    atlasBuildings,
			cropX:      cropX,
			cropY:      cropY,
			cropWidth:  cropWidth,
			cropHeight: cropHeight,
			offX:       offX,
			offY:       offY,
			flipX:      flipX,
		}
	}

	// Tutaj deklaruję dokładnie co i z jakiego wycinka jest duszkiem
	setManualUnits1(spriteManualArcherBtn, 16, 98, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualArcherIdleUpLeft, 0, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherIdleLeft, 0, 98, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherIdleDownLeft, 0, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherIdleDown, 16, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherIdleDownRight, 0, 112, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualArcherIdleRight, 0, 98, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualArcherIdleUpRight, 0, 84, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualArcherIdleUp, 16, 84, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualArcherMove1Center, 48, 98, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove1UpLeft, 32, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove1Left, 32, 98, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove1DownLeft, 32, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove1Down, 48, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove1Up, 48, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove1UpRight, 32, 84, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualArcherMove1Right, 32, 98, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualArcherMove1DownRight, 32, 112, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualArcherMove2Center, 80, 98, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove2UpLeft, 64, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove2Left, 64, 98, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove2DownLeft, 64, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove2Down, 80, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove2Up, 80, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherMove2UpRight, 64, 84, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualArcherMove2Right, 64, 98, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualArcherMove2DownRight, 64, 112, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualArcherAttack1UpLeft, 96, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherAttack1Left, 96, 98, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherAttack1DownLeft, 96, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherAttack1Down, 112, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherAttack1Up, 112, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherAttack1UpRight, 96, 84, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualArcherAttack1Right, 96, 98, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualArcherAttack1DownRight, 96, 112, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualArcherAttack2UpLeft, 128, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherAttack2Left, 128, 98, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherAttack2DownLeft, 128, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherAttack2Down, 144, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherAttack2Up, 144, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualArcherAttack2UpRight, 128, 84, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualArcherAttack2Right, 128, 98, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualArcherAttack2DownRight, 128, 112, 16, 14, 0, 0, true)
	// Koniec łucznika

	setManualUnits1(spriteManualCowBtn, 16, 14, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualCowIdleUpLeft, 0, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowIdleLeft, 0, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowIdleDownLeft, 0, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowIdleDown, 16, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowIdleDownRight, 0, 28, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualCowIdleRight, 0, 14, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualCowIdleUpRight, 0, 0, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualCowIdleUp, 16, 0, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualCowMove1Center, 48, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove1UpLeft, 32, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove1Left, 32, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove1DownLeft, 32, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove1Down, 48, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove1Up, 48, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove1UpRight, 32, 0, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualCowMove1Right, 32, 14, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualCowMove1DownRight, 32, 28, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualCowMove2Center, 80, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove2UpLeft, 64, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove2Left, 64, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove2DownLeft, 64, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove2Down, 80, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove2Up, 80, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowMove2UpRight, 64, 0, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualCowMove2Right, 64, 14, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualCowMove2DownRight, 64, 28, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualCowAttack1UpLeft, 96, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowAttack1Left, 96, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowAttack1DownLeft, 96, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowAttack1Down, 112, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowAttack1Up, 112, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowAttack1UpRight, 96, 0, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualCowAttack1Right, 96, 14, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualCowAttack1DownRight, 96, 28, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualCowAttack2UpLeft, 128, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowAttack2Left, 128, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowAttack2DownLeft, 128, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowAttack2Down, 144, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowAttack2Up, 144, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualCowAttack2UpRight, 128, 0, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualCowAttack2Right, 128, 14, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualCowAttack2DownRight, 128, 28, 16, 14, 0, 0, true)
	// Koniec krowy

	setManualUnits1(spriteManualPriestessBtn, 16, 140, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualPriestessIdleUpLeft, 0, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessIdleLeft, 0, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessIdleDownLeft, 0, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessIdleDown, 16, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessIdleDownRight, 0, 154, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualPriestessIdleRight, 0, 140, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualPriestessIdleUpRight, 0, 126, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualPriestessIdleUp, 16, 126, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualPriestessMove1Center, 48, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove1UpLeft, 32, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove1Left, 32, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove1DownLeft, 32, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove1Down, 48, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove1Up, 48, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove1UpRight, 32, 126, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualPriestessMove1Right, 32, 140, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualPriestessMove1DownRight, 32, 154, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualPriestessMove2Center, 80, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove2UpLeft, 64, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove2Left, 64, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove2DownLeft, 64, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove2Down, 80, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove2Up, 80, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessMove2UpRight, 64, 126, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualPriestessMove2Right, 64, 140, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualPriestessMove2DownRight, 64, 154, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualPriestessAttack1UpLeft, 96, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessAttack1Left, 96, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessAttack1DownLeft, 96, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessAttack1Down, 112, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessAttack1Up, 112, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessAttack1UpRight, 96, 126, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualPriestessAttack1Right, 96, 140, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualPriestessAttack1DownRight, 96, 154, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualPriestessAttack2UpLeft, 128, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessAttack2Left, 128, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessAttack2DownLeft, 128, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessAttack2Down, 144, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessAttack2Up, 144, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualPriestessAttack2UpRight, 128, 126, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualPriestessAttack2Right, 128, 140, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualPriestessAttack2DownRight, 128, 154, 16, 14, 0, 0, true)
	// koniec kapłanki

	setManualUnits1(spriteManualShepherdBtn, 176, 14, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualShepherdIdleUpLeft, 160, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdIdleLeft, 160, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdIdleDownLeft, 160, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdIdleDown, 176, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdIdleDownRight, 160, 28, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualShepherdIdleRight, 160, 14, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualShepherdIdleUpRight, 160, 0, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualShepherdIdleUp, 176, 0, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualShepherdMove1Center, 208, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove1UpLeft, 192, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove1Left, 192, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove1DownLeft, 192, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove1Down, 208, 29, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove1Up, 208, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove1UpRight, 192, 0, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualShepherdMove1Right, 192, 14, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualShepherdMove1DownRight, 192, 28, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualShepherdMove2Center, 240, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove2UpLeft, 224, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove2Left, 224, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove2DownLeft, 224, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove2Down, 240, 29, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove2Up, 240, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdMove2UpRight, 224, 0, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualShepherdMove2Right, 224, 14, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualShepherdMove2DownRight, 224, 28, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualShepherdAttack1UpLeft, 257, 0, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdAttack1Left, 257, 14, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdAttack1DownLeft, 257, 28, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdAttack1Down, 160, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdAttack1Up, 160, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdAttack1UpRight, 257, 0, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualShepherdAttack1Right, 257, 14, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualShepherdAttack1DownRight, 257, 28, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualShepherdAttack2UpLeft, 176, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdAttack2Left, 176, 98, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdAttack2DownLeft, 176, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdAttack2Down, 192, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdAttack2Up, 192, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualShepherdAttack2UpRight, 176, 84, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualShepherdAttack2Right, 176, 98, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualShepherdAttack2DownRight, 176, 112, 16, 14, 0, 0, true)
	// koniec pastucha
	setManualUnits1(spriteManualMageBtn, 176, 140, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualMageIdleUpLeft, 160, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageIdleLeft, 160, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageIdleDownLeft, 160, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageIdleDown, 176, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageIdleDownRight, 160, 154, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualMageIdleRight, 160, 140, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualMageIdleUpRight, 160, 126, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualMageIdleUp, 176, 126, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualMageMove1Center, 208, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove1UpLeft, 192, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove1Left, 192, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove1DownLeft, 192, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove1Down, 208, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove1Up, 208, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove1UpRight, 192, 126, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualMageMove1Right, 192, 140, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualMageMove1DownRight, 192, 154, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualMageMove2Center, 240, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove2UpLeft, 224, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove2Left, 224, 140, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove2DownLeft, 224, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove2Down, 240, 154, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove2Up, 240, 126, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageMove2UpRight, 224, 126, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualMageMove2Right, 224, 140, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualMageMove2DownRight, 224, 154, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualMageAttack1UpLeft, 208, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageAttack1Left, 208, 98, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageAttack1DownLeft, 208, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageAttack1Down, 224, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageAttack1Up, 224, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageAttack1UpRight, 208, 84, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualMageAttack1Right, 208, 98, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualMageAttack1DownRight, 208, 112, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualMageAttack2UpLeft, 240, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageAttack2Left, 240, 98, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageAttack2DownLeft, 240, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageAttack2Down, 256, 112, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageAttack2Up, 256, 84, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualMageAttack2UpRight, 240, 84, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualMageAttack2Right, 240, 98, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualMageAttack2DownRight, 240, 112, 16, 14, 0, 0, true)
	// koniec maga
	setManualUnits2(spriteManualPriestBtn, 16, 14, 16, 14, 0, 0, false)

	setManualUnits2(spriteManualPriestIdleUpLeft, 0, 0, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestIdleLeft, 0, 14, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestIdleDownLeft, 0, 28, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestIdleDown, 16, 28, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestIdleDownRight, 0, 28, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualPriestIdleRight, 0, 14, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualPriestIdleUpRight, 0, 0, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualPriestIdleUp, 16, 0, 16, 14, 0, 0, false)

	setManualUnits2(spriteManualPriestMove1Center, 48, 14, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove1UpLeft, 32, 0, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove1Left, 32, 14, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove1DownLeft, 32, 28, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove1Down, 48, 28, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove1Up, 48, 0, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove1UpRight, 32, 0, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualPriestMove1Right, 32, 14, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualPriestMove1DownRight, 32, 28, 16, 14, 0, 0, true)

	setManualUnits2(spriteManualPriestMove2Center, 80, 14, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove2UpLeft, 64, 0, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove2Left, 64, 14, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove2DownLeft, 64, 28, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove2Down, 80, 28, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove2Up, 80, 0, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestMove2UpRight, 64, 0, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualPriestMove2Right, 64, 14, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualPriestMove2DownRight, 64, 28, 16, 14, 0, 0, true)

	setManualUnits2(spriteManualPriestAttack1UpLeft, 96, 0, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestAttack1Left, 96, 14, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestAttack1DownLeft, 96, 28, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestAttack1Down, 112, 28, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestAttack1Up, 112, 0, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestAttack1UpRight, 96, 0, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualPriestAttack1Right, 96, 14, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualPriestAttack1DownRight, 96, 28, 16, 14, 0, 0, true)

	setManualUnits2(spriteManualPriestAttack2UpLeft, 128, 0, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestAttack2Left, 128, 14, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestAttack2DownLeft, 128, 28, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestAttack2Down, 144, 28, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestAttack2Up, 144, 0, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualPriestAttack2UpRight, 128, 0, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualPriestAttack2Right, 128, 14, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualPriestAttack2DownRight, 128, 28, 16, 14, 0, 0, true)
	// koniec kapłana
	setManualUnits2(spriteManualSpearmanBtn, 16, 98, 16, 14, 0, 0, false)

	setManualUnits2(spriteManualSpearmanIdleUpLeft, 0, 84, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanIdleLeft, 0, 98, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanIdleDownLeft, 0, 112, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanIdleDown, 16, 112, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanIdleDownRight, 0, 112, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSpearmanIdleRight, 0, 98, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSpearmanIdleUpRight, 0, 84, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSpearmanIdleUp, 16, 84, 16, 14, 0, 0, false)

	setManualUnits2(spriteManualSpearmanMove1Center, 48, 98, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove1UpLeft, 32, 84, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove1Left, 32, 98, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove1DownLeft, 32, 112, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove1Down, 48, 112, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove1Up, 48, 84, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove1UpRight, 32, 84, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSpearmanMove1Right, 32, 98, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSpearmanMove1DownRight, 32, 112, 16, 14, 0, 0, true)

	setManualUnits2(spriteManualSpearmanMove2Center, 80, 98, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove2UpLeft, 64, 84, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove2Left, 64, 98, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove2DownLeft, 64, 112, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove2Down, 80, 112, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove2Up, 80, 84, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanMove2UpRight, 64, 84, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSpearmanMove2Right, 64, 98, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSpearmanMove2DownRight, 64, 112, 16, 14, 0, 0, true)

	setManualUnits2(spriteManualSpearmanAttack1UpLeft, 96, 84, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanAttack1Left, 96, 98, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanAttack1DownLeft, 96, 112, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanAttack1Down, 112, 112, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanAttack1Up, 112, 84, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanAttack1UpRight, 96, 84, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSpearmanAttack1Right, 96, 98, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSpearmanAttack1DownRight, 96, 112, 16, 14, 0, 0, true)

	setManualUnits2(spriteManualSpearmanAttack2UpLeft, 128, 84, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanAttack2Left, 128, 98, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanAttack2DownLeft, 128, 112, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanAttack2Down, 144, 112, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanAttack2Up, 144, 84, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSpearmanAttack2UpRight, 128, 84, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSpearmanAttack2Right, 128, 98, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSpearmanAttack2DownRight, 128, 112, 16, 14, 0, 0, true)
	// koniec włócznika
	setManualUnits3(spriteManualCrossbowmanBtn, 240, 14, 16, 14, 0, 0, false)

	setManualUnits3(spriteManualCrossbowmanIdleUpLeft, 224, 0, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanIdleLeft, 224, 14, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanIdleDownLeft, 224, 28, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanIdleDown, 240, 28, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanIdleDownRight, 224, 28, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualCrossbowmanIdleRight, 224, 14, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualCrossbowmanIdleUpRight, 224, 0, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualCrossbowmanIdleUp, 240, 0, 16, 14, 0, 0, false)

	setManualUnits3(spriteManualCrossbowmanMove1Center, 272, 14, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove1UpLeft, 256, 0, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove1Left, 256, 14, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove1DownLeft, 256, 28, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove1Down, 272, 28, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove1Up, 272, 0, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove1UpRight, 256, 0, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualCrossbowmanMove1Right, 256, 14, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualCrossbowmanMove1DownRight, 256, 28, 16, 14, 0, 0, true)

	setManualUnits3(spriteManualCrossbowmanMove2Center, 224, 56, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove2UpLeft, 256, 0, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove2Left, 256, 14, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove2DownLeft, 256, 28, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove2Down, 224, 70, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove2Up, 224, 42, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanMove2UpRight, 256, 0, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualCrossbowmanMove2Right, 256, 14, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualCrossbowmanMove2DownRight, 256, 28, 16, 14, 0, 0, true)

	setManualUnits3(spriteManualCrossbowmanAttack1UpLeft, 240, 42, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanAttack1Left, 240, 56, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanAttack1DownLeft, 240, 70, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanAttack1Down, 256, 70, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanAttack1Up, 256, 42, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanAttack1UpRight, 240, 42, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualCrossbowmanAttack1Right, 240, 56, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualCrossbowmanAttack1DownRight, 240, 70, 16, 14, 0, 0, true)

	setManualUnits3(spriteManualCrossbowmanAttack2UpLeft, 272, 42, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanAttack2Left, 272, 56, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanAttack2DownLeft, 272, 70, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanAttack2Down, 288, 70, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanAttack2Up, 288, 42, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualCrossbowmanAttack2UpRight, 272, 42, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualCrossbowmanAttack2Right, 272, 56, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualCrossbowmanAttack2DownRight, 272, 70, 16, 14, 0, 0, true)
	// koniec kusznika

	// JEDNOSTKI WRĘCZNE!!!
	setManualUnits1(spriteManualAxemanBtn, 16, 56, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualAxemanIdleUpLeft, 0, 42, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanIdleLeft, 0, 56, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanIdleDownLeft, 0, 70, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanIdleDown, 16, 70, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanIdleDownRight, 0, 70, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualAxemanIdleRight, 0, 56, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualAxemanIdleUpRight, 0, 42, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualAxemanIdleUp, 16, 42, 16, 14, 0, 0, false)

	setManualUnits1(spriteManualAxemanMove1Center, 48, 56, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanMove1UpLeft, 32, 42, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanMove1Left, 32, 56, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanMove1DownLeft, 32, 70, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanMove1Down, 48, 70, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanMove1Up, 48, 42, 16, 15, 0, -1, false) // nietypowy rozmiar!
	setManualUnits1(spriteManualAxemanMove1UpRight, 32, 42, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualAxemanMove1Right, 32, 56, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualAxemanMove1DownRight, 32, 70, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualAxemanMove2Center, 80, 56, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanMove2UpLeft, 64, 42, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanMove2Left, 64, 56, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanMove2DownLeft, 64, 70, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanMove2Down, 80, 70, 16, 14, 0, 0, false)
	setManualUnits1(spriteManualAxemanMove2Up, 80, 42, 16, 15, 0, -1, false) // nietypowy rozmiar!
	setManualUnits1(spriteManualAxemanMove2UpRight, 64, 42, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualAxemanMove2Right, 64, 56, 16, 14, 0, 0, true)
	setManualUnits1(spriteManualAxemanMove2DownRight, 64, 70, 16, 14, 0, 0, true)

	setManualUnits1(spriteManualAxemanAttack1UpLeft, 96, 49, 24, 21, -8, -7, false)
	setManualUnits1(spriteManualAxemanAttack1Left, 96, 70, 24, 14, -8, 0, false)
	setManualUnits1(spriteManualAxemanAttack1DownLeft, 152, 42, 24, 21, -8, 0, false)
	setManualUnits1(spriteManualAxemanAttack1Down, 136, 42, 16, 21, 0, 0, false)
	setManualUnits1(spriteManualAxemanAttack1Up, 120, 49, 16, 21, 0, -7, false)
	setManualUnits1(spriteManualAxemanAttack1UpRight, 96, 49, 24, 21, 0, -7, true)
	setManualUnits1(spriteManualAxemanAttack1Right, 96, 70, 24, 14, 0, 0, true)
	setManualUnits1(spriteManualAxemanAttack1DownRight, 152, 42, 24, 21, 0, 0, true)

	setManualUnits1(spriteManualAxemanAttack2UpLeft, 176, 49, 24, 21, -8, -7, false)
	setManualUnits1(spriteManualAxemanAttack2Left, 176, 70, 24, 14, -8, 0, false)
	setManualUnits1(spriteManualAxemanAttack2DownLeft, 152, 63, 24, 21, -8, 0, false)
	setManualUnits1(spriteManualAxemanAttack2Down, 136, 63, 16, 21, 0, 0, false)
	setManualUnits1(spriteManualAxemanAttack2Up, 200, 48, 16, 21, 0, -7, false)
	setManualUnits1(spriteManualAxemanAttack2UpRight, 176, 49, 24, 21, 0, -7, true)
	setManualUnits1(spriteManualAxemanAttack2Right, 176, 70, 24, 14, 0, 0, true)
	setManualUnits1(spriteManualAxemanAttack2DownRight, 152, 63, 24, 21, 0, 0, true)
	// koniec drwala

	setManualUnits2(spriteManualSwordsmanBtn, 16, 56, 16, 14, 0, 0, false)

	setManualUnits2(spriteManualSwordsmanIdleUpLeft, 0, 42, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanIdleLeft, 0, 56, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanIdleDownLeft, 0, 70, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanIdleDown, 16, 70, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanIdleDownRight, 0, 70, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSwordsmanIdleRight, 0, 56, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSwordsmanIdleUpRight, 0, 42, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSwordsmanIdleUp, 16, 42, 16, 14, 0, 0, false)

	setManualUnits2(spriteManualSwordsmanMove1Center, 48, 56, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove1UpLeft, 32, 42, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove1Left, 32, 56, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove1DownLeft, 32, 70, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove1Down, 48, 70, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove1Up, 48, 42, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove1UpRight, 32, 42, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSwordsmanMove1Right, 32, 56, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSwordsmanMove1DownRight, 32, 70, 16, 14, 0, 0, true)

	setManualUnits2(spriteManualSwordsmanMove2Center, 80, 56, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove2UpLeft, 64, 42, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove2Left, 64, 56, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove2DownLeft, 64, 70, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove2Down, 80, 70, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove2Up, 80, 42, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanMove2UpRight, 64, 42, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSwordsmanMove2Right, 64, 56, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualSwordsmanMove2DownRight, 64, 70, 16, 14, 0, 0, true)

	setManualUnits2(spriteManualSwordsmanAttack1UpLeft, 96, 49, 24, 21, -8, -7, false)
	setManualUnits2(spriteManualSwordsmanAttack1Left, 96, 70, 24, 14, -8, 0, false)
	setManualUnits2(spriteManualSwordsmanAttack1DownLeft, 152, 42, 24, 21, -8, 0, false)
	setManualUnits2(spriteManualSwordsmanAttack1Down, 136, 42, 16, 21, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanAttack1Up, 120, 49, 16, 21, 0, -7, false)
	setManualUnits2(spriteManualSwordsmanAttack1UpRight, 96, 49, 24, 21, 0, -7, true)
	setManualUnits2(spriteManualSwordsmanAttack1Right, 96, 70, 24, 14, 0, 0, true)
	setManualUnits2(spriteManualSwordsmanAttack1DownRight, 152, 42, 24, 21, 0, 0, true)

	setManualUnits2(spriteManualSwordsmanAttack2UpLeft, 176, 49, 24, 21, -8, -7, false)
	setManualUnits2(spriteManualSwordsmanAttack2Left, 176, 70, 24, 14, -8, 0, false)
	setManualUnits2(spriteManualSwordsmanAttack2DownLeft, 152, 63, 24, 21, -8, 0, false)
	setManualUnits2(spriteManualSwordsmanAttack2Down, 136, 63, 16, 21, 0, 0, false)
	setManualUnits2(spriteManualSwordsmanAttack2Up, 200, 49, 16, 21, 0, -7, false)
	setManualUnits2(spriteManualSwordsmanAttack2UpRight, 176, 49, 24, 21, 0, -7, true)
	setManualUnits2(spriteManualSwordsmanAttack2Right, 176, 70, 24, 14, 0, 0, true)
	setManualUnits2(spriteManualSwordsmanAttack2DownRight, 152, 63, 24, 21, 0, 0, true)
	// Koniec miecznika

	setManualUnits2(spriteManualCommanderBtn, 16, 140, 16, 14, 0, 0, false)

	setManualUnits2(spriteManualCommanderIdleUpLeft, 0, 126, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderIdleLeft, 0, 140, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderIdleDownLeft, 0, 154, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderIdleDown, 16, 154, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderIdleDownRight, 0, 154, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualCommanderIdleRight, 0, 140, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualCommanderIdleUpRight, 0, 126, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualCommanderIdleUp, 16, 126, 16, 14, 0, 0, false)

	setManualUnits2(spriteManualCommanderMove1Center, 48, 140, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove1UpLeft, 32, 126, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove1Left, 32, 140, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove1DownLeft, 32, 154, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove1Down, 48, 154, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove1Up, 48, 126, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove1UpRight, 32, 126, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualCommanderMove1Right, 32, 140, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualCommanderMove1DownRight, 32, 154, 16, 14, 0, 0, true)

	setManualUnits2(spriteManualCommanderMove2Center, 80, 140, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove2UpLeft, 64, 126, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove2Left, 64, 140, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove2DownLeft, 64, 154, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove2Down, 80, 154, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove2Up, 80, 126, 16, 14, 0, 0, false)
	setManualUnits2(spriteManualCommanderMove2UpRight, 64, 126, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualCommanderMove2Right, 64, 140, 16, 14, 0, 0, true)
	setManualUnits2(spriteManualCommanderMove2DownRight, 64, 154, 16, 14, 0, 0, true)

	setManualUnits2(spriteManualCommanderAttack1UpLeft, 96, 133, 24, 21, -8, -7, false)
	setManualUnits2(spriteManualCommanderAttack1Left, 96, 154, 24, 14, -8, 0, false)
	setManualUnits2(spriteManualCommanderAttack1DownLeft, 152, 126, 24, 21, -8, 0, false)
	setManualUnits2(spriteManualCommanderAttack1Down, 136, 126, 16, 21, 0, 0, false)
	setManualUnits2(spriteManualCommanderAttack1Up, 120, 133, 16, 21, 0, -7, false)
	setManualUnits2(spriteManualCommanderAttack1UpRight, 96, 133, 24, 21, 0, -7, true)
	setManualUnits2(spriteManualCommanderAttack1Right, 96, 154, 24, 14, 0, 0, true)
	setManualUnits2(spriteManualCommanderAttack1DownRight, 152, 126, 24, 21, 0, 0, true)

	setManualUnits2(spriteManualCommanderAttack2UpLeft, 176, 133, 24, 21, -8, -7, false)
	setManualUnits2(spriteManualCommanderAttack2Left, 176, 154, 24, 14, -8, 0, false)
	setManualUnits2(spriteManualCommanderAttack2DownLeft, 152, 147, 24, 21, -8, 0, false)
	setManualUnits2(spriteManualCommanderAttack2Down, 136, 147, 16, 21, 0, 0, false)
	setManualUnits2(spriteManualCommanderAttack2Up, 200, 133, 16, 21, 0, -7, false)
	setManualUnits2(spriteManualCommanderAttack2UpRight, 176, 133, 24, 21, 0, -7, true)
	setManualUnits2(spriteManualCommanderAttack2Right, 176, 154, 24, 14, 0, 0, true)
	setManualUnits2(spriteManualCommanderAttack2DownRight, 152, 147, 24, 21, 0, 0, true)
	// Koniec dowódcy

	setManualUnits3(spriteManualBearBtn, 16, 14, 16, 14, 0, 0, false)

	setManualUnits3(spriteManualBearIdleUpLeft, 0, 0, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearIdleLeft, 0, 14, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearIdleDownLeft, 0, 28, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearIdleDown, 16, 28, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearIdleDownRight, 0, 28, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualBearIdleRight, 0, 14, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualBearIdleUpRight, 0, 0, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualBearIdleUp, 16, 0, 16, 14, 0, 0, false)

	setManualUnits3(spriteManualBearMove1Center, 48, 14, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove1UpLeft, 32, 0, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove1Left, 32, 14, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove1DownLeft, 32, 28, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove1Down, 48, 28, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove1Up, 48, 0, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove1UpRight, 32, 0, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualBearMove1Right, 32, 14, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualBearMove1DownRight, 32, 28, 16, 14, 0, 0, true)

	setManualUnits3(spriteManualBearMove2Center, 80, 14, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove2UpLeft, 64, 0, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove2Left, 64, 14, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove2DownLeft, 64, 28, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove2Down, 80, 28, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove2Up, 80, 0, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualBearMove2UpRight, 64, 0, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualBearMove2Right, 64, 14, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualBearMove2DownRight, 64, 28, 16, 14, 0, 0, true)

	setManualUnits3(spriteManualBearAttack1UpLeft, 96, 7, 24, 21, -8, -7, false)
	setManualUnits3(spriteManualBearAttack1Left, 96, 28, 24, 14, -8, 0, false)
	setManualUnits3(spriteManualBearAttack1DownLeft, 152, 0, 24, 21, -8, 0, false)
	setManualUnits3(spriteManualBearAttack1Down, 136, 0, 16, 21, 0, 0, false)
	setManualUnits3(spriteManualBearAttack1Up, 120, 7, 16, 21, 0, -7, false)
	setManualUnits3(spriteManualBearAttack1UpRight, 96, 7, 24, 21, 0, -7, true)
	setManualUnits3(spriteManualBearAttack1Right, 96, 28, 24, 14, 0, 0, true)
	setManualUnits3(spriteManualBearAttack1DownRight, 152, 0, 24, 21, 0, 0, true)

	setManualUnits3(spriteManualBearAttack2UpLeft, 176, 7, 24, 21, -8, -7, false)
	setManualUnits3(spriteManualBearAttack2Left, 176, 28, 24, 14, -8, 0, false)
	setManualUnits3(spriteManualBearAttack2DownLeft, 152, 21, 24, 21, -8, 0, false)
	setManualUnits3(spriteManualBearAttack2Down, 136, 21, 16, 21, 0, 0, false)
	setManualUnits3(spriteManualBearAttack2Up, 200, 7, 16, 21, 0, -7, false)
	setManualUnits3(spriteManualBearAttack2UpRight, 176, 7, 24, 21, 0, -7, true)
	setManualUnits3(spriteManualBearAttack2Right, 176, 28, 24, 14, 0, 0, true)
	setManualUnits3(spriteManualBearAttack2DownRight, 152, 21, 24, 21, 0, 0, true)
	// Koniec niedźwiedzia

	setManualUnits3(spriteManualUnknownBtn, 16, 56, 16, 14, 0, 0, false)

	setManualUnits3(spriteManualUnknownIdleUpLeft, 0, 42, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownIdleLeft, 0, 56, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownIdleDownLeft, 0, 70, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownIdleDown, 16, 70, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownIdleDownRight, 0, 70, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualUnknownIdleRight, 0, 56, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualUnknownIdleUpRight, 0, 42, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualUnknownIdleUp, 16, 42, 16, 14, 0, 0, false)

	setManualUnits3(spriteManualUnknownMove1Center, 48, 56, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove1UpLeft, 32, 42, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove1Left, 32, 56, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove1DownLeft, 32, 70, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove1Down, 48, 70, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove1Up, 48, 42, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove1UpRight, 32, 42, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualUnknownMove1Right, 32, 56, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualUnknownMove1DownRight, 32, 70, 16, 14, 0, 0, true)

	setManualUnits3(spriteManualUnknownMove2Center, 80, 56, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove2UpLeft, 64, 42, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove2Left, 64, 56, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove2DownLeft, 64, 70, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove2Down, 80, 70, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove2Up, 80, 42, 16, 14, 0, 0, false)
	setManualUnits3(spriteManualUnknownMove2UpRight, 64, 42, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualUnknownMove2Right, 64, 56, 16, 14, 0, 0, true)
	setManualUnits3(spriteManualUnknownMove2DownRight, 64, 70, 16, 14, 0, 0, true)

	setManualUnits3(spriteManualUnknownAttack1UpLeft, 96, 49, 24, 21, -8, -7, false)
	setManualUnits3(spriteManualUnknownAttack1Left, 96, 70, 24, 14, -8, 0, false)
	setManualUnits3(spriteManualUnknownAttack1DownLeft, 152, 42, 24, 21, -8, 0, false)
	setManualUnits3(spriteManualUnknownAttack1Down, 136, 42, 16, 21, 0, 0, false)
	setManualUnits3(spriteManualUnknownAttack1Up, 120, 49, 16, 21, 0, -7, false)
	setManualUnits3(spriteManualUnknownAttack1UpRight, 96, 49, 24, 21, 0, -7, true)
	setManualUnits3(spriteManualUnknownAttack1Right, 96, 70, 24, 14, 0, 0, true)
	setManualUnits3(spriteManualUnknownAttack1DownRight, 152, 42, 24, 21, 0, 0, true)

	setManualUnits3(spriteManualUnknownAttack2UpLeft, 176, 49, 24, 21, -8, -7, false)
	setManualUnits3(spriteManualUnknownAttack2Left, 176, 70, 24, 14, -8, 0, false)
	setManualUnits3(spriteManualUnknownAttack2DownLeft, 152, 63, 24, 21, -8, 0, false)
	setManualUnits3(spriteManualUnknownAttack2Down, 136, 63, 16, 21, 0, 0, false)
	setManualUnits3(spriteManualUnknownAttack2Up, 200, 49, 16, 21, 0, -7, false)
	setManualUnits3(spriteManualUnknownAttack2UpRight, 176, 49, 24, 21, 0, -7, true)
	setManualUnits3(spriteManualUnknownAttack2Right, 176, 70, 24, 14, 0, 0, true)
	setManualUnits3(spriteManualUnknownAttack2DownRight, 152, 63, 24, 21, 0, 0, true)
	// Koniec strzygi
}

func initManualCorpseSprites() {
	setManualCorpse := func(spriteID uint16, atlas battleAtlasID, cropX, cropY uint16) {
		if spriteID >= maxSpriteID {
			log.Printf("OSTRZEŻENIE: ręcznie wpisany duszek zwłok nie mieści się!")
			return
		}

		spriteRegistry[spriteID] = spriteDef{
			atlasID:    atlas,
			cropX:      cropX,
			cropY:      cropY,
			cropWidth:  tileWidth,
			cropHeight: tileHeight,
			offX:       0,
			offY:       0,
			flipX:      false,
		}
	}

	setManualCorpse(spriteManualArcherCorpseFresh, atlasUnits1, 112, 98)
	setManualCorpse(spriteManualArcherCorpseDecay, atlasUnits1, 144, 98)

	setManualCorpse(spriteManualAxemanCorpseFresh, atlasUnits1, 120, 70)
	setManualCorpse(spriteManualAxemanCorpseDecay, atlasUnits1, 200, 70)

	setManualCorpse(spriteManualBearCorpseFresh, atlasBuildings, 120, 28)
	setManualCorpse(spriteManualBearCorpseDecay, atlasBuildings, 200, 28)

	setManualCorpse(spriteManualCommanderCorpseFresh, atlasUnits2, 120, 154)
	setManualCorpse(spriteManualCommanderCorpseDecay, atlasUnits2, 200, 154)

	setManualCorpse(spriteManualCowCorpseFresh, atlasUnits1, 112, 14)
	setManualCorpse(spriteManualCowCorpseDecay, atlasUnits1, 144, 14)

	setManualCorpse(spriteManualCrossbowmanCorpseFresh, atlasBuildings, 255, 56)
	setManualCorpse(spriteManualCrossbowmanCorpseDecay, atlasBuildings, 287, 56)

	setManualCorpse(spriteManualMageCorpseFresh, atlasUnits1, 224, 98)
	setManualCorpse(spriteManualMageCorpseDecay, atlasUnits1, 256, 98)

	setManualCorpse(spriteManualPriestCorpseFresh, atlasUnits2, 112, 14)
	setManualCorpse(spriteManualPriestCorpseDecay, atlasUnits2, 144, 14)

	setManualCorpse(spriteManualPriestessCorpseFresh, atlasUnits1, 111, 140)
	setManualCorpse(spriteManualPriestessCorpseDecay, atlasUnits1, 143, 140)

	setManualCorpse(spriteManualShepherdCorpseFresh, atlasUnits1, 160, 98)
	setManualCorpse(spriteManualShepherdCorpseDecay, atlasUnits1, 192, 98)

	setManualCorpse(spriteManualSpearmanCorpseFresh, atlasUnits2, 112, 98)
	setManualCorpse(spriteManualSpearmanCorpseDecay, atlasUnits2, 144, 98)

	setManualCorpse(spriteManualSwordsmanCorpseFresh, atlasUnits2, 120, 70)
	setManualCorpse(spriteManualSwordsmanCorpseDecay, atlasUnits2, 200, 70)

	setManualCorpse(spriteManualUnknownCorpseFresh, atlasBuildings, 120, 70)
	setManualCorpse(spriteManualUnknownCorpseDecay, atlasBuildings, 200, 70)
}

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
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
	}
	spriteTable[unitCow][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
	}
	spriteTable[unitCow][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
	}
	spriteTable[unitCow][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
	}
	spriteTable[unitCow][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
	}
	spriteTable[unitCow][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
	}
	spriteTable[unitCow][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
	}
	spriteTable[unitCow][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
	}
	spriteTable[unitCow][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
	}

	spriteTable[unitCow][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowMove1Down,
		spriteManualCowMove1Down,
		spriteManualCowMove1Down,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowMove2Down,
		spriteManualCowMove2Down,
		spriteManualCowMove2Down,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowMove1Down,
		spriteManualCowMove1Down,
		spriteManualCowIdleDown,
	}
	spriteTable[unitCow][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowMove1Down,
		spriteManualCowMove1Down,
		spriteManualCowMove1Down,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowMove2Down,
		spriteManualCowMove2Down,
		spriteManualCowMove2Down,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowMove1Down,
		spriteManualCowMove1Down,
		spriteManualCowIdleDown,
	}
	spriteTable[unitCow][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowMove1Left,
		spriteManualCowMove1Left,
		spriteManualCowMove1Left,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowMove2Left,
		spriteManualCowMove2Left,
		spriteManualCowMove2Left,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowMove1Left,
		spriteManualCowMove1Left,
		spriteManualCowIdleLeft,
	}
	spriteTable[unitCow][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowMove1Up,
		spriteManualCowMove1Up,
		spriteManualCowMove1Up,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowMove2Up,
		spriteManualCowMove2Up,
		spriteManualCowMove2Up,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowMove1Up,
		spriteManualCowMove1Up,
		spriteManualCowIdleUp,
	}
	spriteTable[unitCow][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowMove1Right,
		spriteManualCowMove1Right,
		spriteManualCowMove1Right,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowMove2Right,
		spriteManualCowMove2Right,
		spriteManualCowMove2Right,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowMove1Right,
		spriteManualCowMove1Right,
		spriteManualCowIdleRight,
	}
	spriteTable[unitCow][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowMove1UpLeft,
		spriteManualCowMove1UpLeft,
		spriteManualCowMove1UpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowMove2UpLeft,
		spriteManualCowMove2UpLeft,
		spriteManualCowMove2UpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowMove1UpLeft,
		spriteManualCowMove1UpLeft,
		spriteManualCowIdleUpLeft,
	}
	spriteTable[unitCow][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowMove1UpRight,
		spriteManualCowMove1UpRight,
		spriteManualCowMove1UpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowMove2UpRight,
		spriteManualCowMove2UpRight,
		spriteManualCowMove2UpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowMove1UpRight,
		spriteManualCowMove1UpRight,
		spriteManualCowIdleUpRight,
	}
	spriteTable[unitCow][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowMove1DownLeft,
		spriteManualCowMove1DownLeft,
		spriteManualCowMove1DownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowMove2DownLeft,
		spriteManualCowMove2DownLeft,
		spriteManualCowMove2DownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowMove1DownLeft,
		spriteManualCowMove1DownLeft,
		spriteManualCowIdleDownLeft,
	}
	spriteTable[unitCow][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowMove1DownRight,
		spriteManualCowMove1DownRight,
		spriteManualCowMove1DownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowMove2DownRight,
		spriteManualCowMove2DownRight,
		spriteManualCowMove2DownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowMove1DownRight,
		spriteManualCowMove1DownRight,
		spriteManualCowIdleDownRight,
	}

	spriteTable[unitCow][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowAttack2Down,
		spriteManualCowAttack2Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
	}
	spriteTable[unitCow][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowIdleDown,
		spriteManualCowAttack2Down,
		spriteManualCowAttack2Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
		spriteManualCowAttack1Down,
	}
	spriteTable[unitCow][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowIdleLeft,
		spriteManualCowAttack2Left,
		spriteManualCowAttack2Left,
		spriteManualCowAttack1Left,
		spriteManualCowAttack1Left,
		spriteManualCowAttack1Left,
		spriteManualCowAttack1Left,
		spriteManualCowAttack1Left,
		spriteManualCowAttack1Left,
		spriteManualCowAttack1Left,
		spriteManualCowAttack1Left,
		spriteManualCowAttack1Left,
		spriteManualCowAttack1Left,
		spriteManualCowAttack1Left,
		spriteManualCowAttack1Left,
	}
	spriteTable[unitCow][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowIdleUp,
		spriteManualCowAttack2Up,
		spriteManualCowAttack2Up,
		spriteManualCowAttack1Up,
		spriteManualCowAttack1Up,
		spriteManualCowAttack1Up,
		spriteManualCowAttack1Up,
		spriteManualCowAttack1Up,
		spriteManualCowAttack1Up,
		spriteManualCowAttack1Up,
		spriteManualCowAttack1Up,
		spriteManualCowAttack1Up,
		spriteManualCowAttack1Up,
		spriteManualCowAttack1Up,
		spriteManualCowAttack1Up,
	}
	spriteTable[unitCow][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowIdleRight,
		spriteManualCowAttack2Right,
		spriteManualCowAttack2Right,
		spriteManualCowAttack1Right,
		spriteManualCowAttack1Right,
		spriteManualCowAttack1Right,
		spriteManualCowAttack1Right,
		spriteManualCowAttack1Right,
		spriteManualCowAttack1Right,
		spriteManualCowAttack1Right,
		spriteManualCowAttack1Right,
		spriteManualCowAttack1Right,
		spriteManualCowAttack1Right,
		spriteManualCowAttack1Right,
		spriteManualCowAttack1Right,
	}
	spriteTable[unitCow][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowIdleUpLeft,
		spriteManualCowAttack2UpLeft,
		spriteManualCowAttack2UpLeft,
		spriteManualCowAttack1UpLeft,
		spriteManualCowAttack1UpLeft,
		spriteManualCowAttack1UpLeft,
		spriteManualCowAttack1UpLeft,
		spriteManualCowAttack1UpLeft,
		spriteManualCowAttack1UpLeft,
		spriteManualCowAttack1UpLeft,
		spriteManualCowAttack1UpLeft,
		spriteManualCowAttack1UpLeft,
		spriteManualCowAttack1UpLeft,
		spriteManualCowAttack1UpLeft,
		spriteManualCowAttack1UpLeft,
	}
	spriteTable[unitCow][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowIdleUpRight,
		spriteManualCowAttack2UpRight,
		spriteManualCowAttack2UpRight,
		spriteManualCowAttack1UpRight,
		spriteManualCowAttack1UpRight,
		spriteManualCowAttack1UpRight,
		spriteManualCowAttack1UpRight,
		spriteManualCowAttack1UpRight,
		spriteManualCowAttack1UpRight,
		spriteManualCowAttack1UpRight,
		spriteManualCowAttack1UpRight,
		spriteManualCowAttack1UpRight,
		spriteManualCowAttack1UpRight,
		spriteManualCowAttack1UpRight,
		spriteManualCowAttack1UpRight,
	}
	spriteTable[unitCow][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowIdleDownLeft,
		spriteManualCowAttack2DownLeft,
		spriteManualCowAttack2DownLeft,
		spriteManualCowAttack1DownLeft,
		spriteManualCowAttack1DownLeft,
		spriteManualCowAttack1DownLeft,
		spriteManualCowAttack1DownLeft,
		spriteManualCowAttack1DownLeft,
		spriteManualCowAttack1DownLeft,
		spriteManualCowAttack1DownLeft,
		spriteManualCowAttack1DownLeft,
		spriteManualCowAttack1DownLeft,
		spriteManualCowAttack1DownLeft,
		spriteManualCowAttack1DownLeft,
		spriteManualCowAttack1DownLeft,
	}
	spriteTable[unitCow][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowIdleDownRight,
		spriteManualCowAttack2DownRight,
		spriteManualCowAttack2DownRight,
		spriteManualCowAttack1DownRight,
		spriteManualCowAttack1DownRight,
		spriteManualCowAttack1DownRight,
		spriteManualCowAttack1DownRight,
		spriteManualCowAttack1DownRight,
		spriteManualCowAttack1DownRight,
		spriteManualCowAttack1DownRight,
		spriteManualCowAttack1DownRight,
		spriteManualCowAttack1DownRight,
		spriteManualCowAttack1DownRight,
		spriteManualCowAttack1DownRight,
		spriteManualCowAttack1DownRight,
	}
	// koniec krowy

	spriteTable[unitPriestess][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
	}
	spriteTable[unitPriestess][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
	}
	spriteTable[unitPriestess][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
	}
	spriteTable[unitPriestess][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
	}
	spriteTable[unitPriestess][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
	}
	spriteTable[unitPriestess][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
	}
	spriteTable[unitPriestess][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
	}
	spriteTable[unitPriestess][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
	}
	spriteTable[unitPriestess][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
	}

	spriteTable[unitPriestess][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessMove1Down,
		spriteManualPriestessMove1Down,
		spriteManualPriestessMove1Down,
		spriteManualPriestessMove1Down,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessMove2Down,
		spriteManualPriestessMove2Down,
		spriteManualPriestessMove2Down,
		spriteManualPriestessMove2Down,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
	}
	spriteTable[unitPriestess][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessMove1Down,
		spriteManualPriestessMove1Down,
		spriteManualPriestessMove1Down,
		spriteManualPriestessMove1Down,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessMove2Down,
		spriteManualPriestessMove2Down,
		spriteManualPriestessMove2Down,
		spriteManualPriestessMove2Down,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
	}
	spriteTable[unitPriestess][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessMove1Left,
		spriteManualPriestessMove1Left,
		spriteManualPriestessMove1Left,
		spriteManualPriestessMove1Left,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessMove2Left,
		spriteManualPriestessMove2Left,
		spriteManualPriestessMove2Left,
		spriteManualPriestessMove2Left,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
	}
	spriteTable[unitPriestess][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessMove1Up,
		spriteManualPriestessMove1Up,
		spriteManualPriestessMove1Up,
		spriteManualPriestessMove1Up,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessMove2Up,
		spriteManualPriestessMove2Up,
		spriteManualPriestessMove2Up,
		spriteManualPriestessMove2Up,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
	}
	spriteTable[unitPriestess][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessMove1Right,
		spriteManualPriestessMove1Right,
		spriteManualPriestessMove1Right,
		spriteManualPriestessMove1Right,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessMove2Right,
		spriteManualPriestessMove2Right,
		spriteManualPriestessMove2Right,
		spriteManualPriestessMove2Right,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
	}
	spriteTable[unitPriestess][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessMove1UpLeft,
		spriteManualPriestessMove1UpLeft,
		spriteManualPriestessMove1UpLeft,
		spriteManualPriestessMove1UpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessMove2UpLeft,
		spriteManualPriestessMove2UpLeft,
		spriteManualPriestessMove2UpLeft,
		spriteManualPriestessMove2UpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
	}
	spriteTable[unitPriestess][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessMove1UpRight,
		spriteManualPriestessMove1UpRight,
		spriteManualPriestessMove1UpRight,
		spriteManualPriestessMove1UpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessMove2UpRight,
		spriteManualPriestessMove2UpRight,
		spriteManualPriestessMove2UpRight,
		spriteManualPriestessMove2UpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
	}
	spriteTable[unitPriestess][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessMove1DownLeft,
		spriteManualPriestessMove1DownLeft,
		spriteManualPriestessMove1DownLeft,
		spriteManualPriestessMove1DownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessMove2DownLeft,
		spriteManualPriestessMove2DownLeft,
		spriteManualPriestessMove2DownLeft,
		spriteManualPriestessMove2DownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
	}
	spriteTable[unitPriestess][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessMove1DownRight,
		spriteManualPriestessMove1DownRight,
		spriteManualPriestessMove1DownRight,
		spriteManualPriestessMove1DownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessMove2DownRight,
		spriteManualPriestessMove2DownRight,
		spriteManualPriestessMove2DownRight,
		spriteManualPriestessMove2DownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
	}

	spriteTable[unitPriestess][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessAttack1Down,
		spriteManualPriestessIdleDown,
		spriteManualPriestessAttack2Down,
		spriteManualPriestessAttack2Down,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
	}
	spriteTable[unitPriestess][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
		spriteManualPriestessAttack1Down,
		spriteManualPriestessIdleDown,
		spriteManualPriestessAttack2Down,
		spriteManualPriestessAttack2Down,
		spriteManualPriestessIdleDown,
		spriteManualPriestessIdleDown,
	}
	spriteTable[unitPriestess][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessAttack1Left,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessAttack2Left,
		spriteManualPriestessAttack2Left,
		spriteManualPriestessIdleLeft,
		spriteManualPriestessIdleLeft,
	}
	spriteTable[unitPriestess][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
		spriteManualPriestessAttack1Up,
		spriteManualPriestessIdleUp,
		spriteManualPriestessAttack2Up,
		spriteManualPriestessAttack2Up,
		spriteManualPriestessIdleUp,
		spriteManualPriestessIdleUp,
	}
	spriteTable[unitPriestess][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
		spriteManualPriestessAttack1Right,
		spriteManualPriestessIdleRight,
		spriteManualPriestessAttack2Right,
		spriteManualPriestessAttack2Right,
		spriteManualPriestessIdleRight,
		spriteManualPriestessIdleRight,
	}
	spriteTable[unitPriestess][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessAttack1UpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessAttack2UpLeft,
		spriteManualPriestessAttack2UpLeft,
		spriteManualPriestessIdleUpLeft,
		spriteManualPriestessIdleUpLeft,
	}
	spriteTable[unitPriestess][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessAttack1UpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessAttack2UpRight,
		spriteManualPriestessAttack2UpRight,
		spriteManualPriestessIdleUpRight,
		spriteManualPriestessIdleUpRight,
	}
	spriteTable[unitPriestess][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessAttack1DownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessAttack2DownLeft,
		spriteManualPriestessAttack2DownLeft,
		spriteManualPriestessIdleDownLeft,
		spriteManualPriestessIdleDownLeft,
	}
	spriteTable[unitPriestess][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessAttack1DownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessAttack2DownRight,
		spriteManualPriestessAttack2DownRight,
		spriteManualPriestessIdleDownRight,
		spriteManualPriestessIdleDownRight,
	}
	// Koniec Kapłanki

	spriteTable[unitShepherd][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
	}
	spriteTable[unitShepherd][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
	}
	spriteTable[unitShepherd][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
	}
	spriteTable[unitShepherd][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
	}
	spriteTable[unitShepherd][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
	}
	spriteTable[unitShepherd][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
	}
	spriteTable[unitShepherd][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
	}
	spriteTable[unitShepherd][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
	}
	spriteTable[unitShepherd][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
	}

	spriteTable[unitShepherd][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdMove1Down,
		spriteManualShepherdMove1Down,
		spriteManualShepherdMove1Down,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdMove2Down,
		spriteManualShepherdMove2Down,
		spriteManualShepherdMove2Down,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdMove1Down,
		spriteManualShepherdMove1Down,
		spriteManualShepherdIdleDown,
	}
	spriteTable[unitShepherd][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdMove1Down,
		spriteManualShepherdMove1Down,
		spriteManualShepherdMove1Down,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdMove2Down,
		spriteManualShepherdMove2Down,
		spriteManualShepherdMove2Down,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdMove1Down,
		spriteManualShepherdMove1Down,
		spriteManualShepherdIdleDown,
	}
	spriteTable[unitShepherd][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdMove1Left,
		spriteManualShepherdMove1Left,
		spriteManualShepherdMove1Left,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdMove2Left,
		spriteManualShepherdMove2Left,
		spriteManualShepherdMove2Left,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdMove1Left,
		spriteManualShepherdMove1Left,
		spriteManualShepherdIdleLeft,
	}
	spriteTable[unitShepherd][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdMove1Up,
		spriteManualShepherdMove1Up,
		spriteManualShepherdMove1Up,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdMove2Up,
		spriteManualShepherdMove2Up,
		spriteManualShepherdMove2Up,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdMove1Up,
		spriteManualShepherdMove1Up,
		spriteManualShepherdIdleUp,
	}
	spriteTable[unitShepherd][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdMove1Right,
		spriteManualShepherdMove1Right,
		spriteManualShepherdMove1Right,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdMove2Right,
		spriteManualShepherdMove2Right,
		spriteManualShepherdMove2Right,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdMove1Right,
		spriteManualShepherdMove1Right,
		spriteManualShepherdIdleRight,
	}
	spriteTable[unitShepherd][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdMove1UpLeft,
		spriteManualShepherdMove1UpLeft,
		spriteManualShepherdMove1UpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdMove2UpLeft,
		spriteManualShepherdMove2UpLeft,
		spriteManualShepherdMove2UpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdMove1UpLeft,
		spriteManualShepherdMove1UpLeft,
		spriteManualShepherdIdleUpLeft,
	}
	spriteTable[unitShepherd][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdMove1UpRight,
		spriteManualShepherdMove1UpRight,
		spriteManualShepherdMove1UpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdMove2UpRight,
		spriteManualShepherdMove2UpRight,
		spriteManualShepherdMove2UpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdMove1UpRight,
		spriteManualShepherdMove1UpRight,
		spriteManualShepherdIdleUpRight,
	}
	spriteTable[unitShepherd][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdMove1DownLeft,
		spriteManualShepherdMove1DownLeft,
		spriteManualShepherdMove1DownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdMove2DownLeft,
		spriteManualShepherdMove2DownLeft,
		spriteManualShepherdMove2DownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdMove1DownLeft,
		spriteManualShepherdMove1DownLeft,
		spriteManualShepherdIdleDownLeft,
	}
	spriteTable[unitShepherd][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdMove1DownRight,
		spriteManualShepherdMove1DownRight,
		spriteManualShepherdMove1DownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdMove2DownRight,
		spriteManualShepherdMove2DownRight,
		spriteManualShepherdMove2DownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdMove1DownRight,
		spriteManualShepherdMove1DownRight,
		spriteManualShepherdIdleDownRight,
	}

	spriteTable[unitShepherd][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdAttack2Down,
		spriteManualShepherdAttack2Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
	}
	spriteTable[unitShepherd][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdIdleDown,
		spriteManualShepherdAttack2Down,
		spriteManualShepherdAttack2Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
		spriteManualShepherdAttack1Down,
	}
	spriteTable[unitShepherd][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdIdleLeft,
		spriteManualShepherdAttack2Left,
		spriteManualShepherdAttack2Left,
		spriteManualShepherdAttack1Left,
		spriteManualShepherdAttack1Left,
		spriteManualShepherdAttack1Left,
		spriteManualShepherdAttack1Left,
		spriteManualShepherdAttack1Left,
		spriteManualShepherdAttack1Left,
		spriteManualShepherdAttack1Left,
		spriteManualShepherdAttack1Left,
		spriteManualShepherdAttack1Left,
		spriteManualShepherdAttack1Left,
		spriteManualShepherdAttack1Left,
		spriteManualShepherdAttack1Left,
	}
	spriteTable[unitShepherd][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdIdleUp,
		spriteManualShepherdAttack2Up,
		spriteManualShepherdAttack2Up,
		spriteManualShepherdAttack1Up,
		spriteManualShepherdAttack1Up,
		spriteManualShepherdAttack1Up,
		spriteManualShepherdAttack1Up,
		spriteManualShepherdAttack1Up,
		spriteManualShepherdAttack1Up,
		spriteManualShepherdAttack1Up,
		spriteManualShepherdAttack1Up,
		spriteManualShepherdAttack1Up,
		spriteManualShepherdAttack1Up,
		spriteManualShepherdAttack1Up,
		spriteManualShepherdAttack1Up,
	}
	spriteTable[unitShepherd][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdIdleRight,
		spriteManualShepherdAttack2Right,
		spriteManualShepherdAttack2Right,
		spriteManualShepherdAttack1Right,
		spriteManualShepherdAttack1Right,
		spriteManualShepherdAttack1Right,
		spriteManualShepherdAttack1Right,
		spriteManualShepherdAttack1Right,
		spriteManualShepherdAttack1Right,
		spriteManualShepherdAttack1Right,
		spriteManualShepherdAttack1Right,
		spriteManualShepherdAttack1Right,
		spriteManualShepherdAttack1Right,
		spriteManualShepherdAttack1Right,
		spriteManualShepherdAttack1Right,
	}
	spriteTable[unitShepherd][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdIdleUpLeft,
		spriteManualShepherdAttack2UpLeft,
		spriteManualShepherdAttack2UpLeft,
		spriteManualShepherdAttack1UpLeft,
		spriteManualShepherdAttack1UpLeft,
		spriteManualShepherdAttack1UpLeft,
		spriteManualShepherdAttack1UpLeft,
		spriteManualShepherdAttack1UpLeft,
		spriteManualShepherdAttack1UpLeft,
		spriteManualShepherdAttack1UpLeft,
		spriteManualShepherdAttack1UpLeft,
		spriteManualShepherdAttack1UpLeft,
		spriteManualShepherdAttack1UpLeft,
		spriteManualShepherdAttack1UpLeft,
		spriteManualShepherdAttack1UpLeft,
	}
	spriteTable[unitShepherd][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdIdleUpRight,
		spriteManualShepherdAttack2UpRight,
		spriteManualShepherdAttack2UpRight,
		spriteManualShepherdAttack1UpRight,
		spriteManualShepherdAttack1UpRight,
		spriteManualShepherdAttack1UpRight,
		spriteManualShepherdAttack1UpRight,
		spriteManualShepherdAttack1UpRight,
		spriteManualShepherdAttack1UpRight,
		spriteManualShepherdAttack1UpRight,
		spriteManualShepherdAttack1UpRight,
		spriteManualShepherdAttack1UpRight,
		spriteManualShepherdAttack1UpRight,
		spriteManualShepherdAttack1UpRight,
		spriteManualShepherdAttack1UpRight,
	}
	spriteTable[unitShepherd][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdIdleDownLeft,
		spriteManualShepherdAttack2DownLeft,
		spriteManualShepherdAttack2DownLeft,
		spriteManualShepherdAttack1DownLeft,
		spriteManualShepherdAttack1DownLeft,
		spriteManualShepherdAttack1DownLeft,
		spriteManualShepherdAttack1DownLeft,
		spriteManualShepherdAttack1DownLeft,
		spriteManualShepherdAttack1DownLeft,
		spriteManualShepherdAttack1DownLeft,
		spriteManualShepherdAttack1DownLeft,
		spriteManualShepherdAttack1DownLeft,
		spriteManualShepherdAttack1DownLeft,
		spriteManualShepherdAttack1DownLeft,
		spriteManualShepherdAttack1DownLeft,
	}
	spriteTable[unitShepherd][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdIdleDownRight,
		spriteManualShepherdAttack2DownRight,
		spriteManualShepherdAttack2DownRight,
		spriteManualShepherdAttack1DownRight,
		spriteManualShepherdAttack1DownRight,
		spriteManualShepherdAttack1DownRight,
		spriteManualShepherdAttack1DownRight,
		spriteManualShepherdAttack1DownRight,
		spriteManualShepherdAttack1DownRight,
		spriteManualShepherdAttack1DownRight,
		spriteManualShepherdAttack1DownRight,
		spriteManualShepherdAttack1DownRight,
		spriteManualShepherdAttack1DownRight,
		spriteManualShepherdAttack1DownRight,
		spriteManualShepherdAttack1DownRight,
	}
	// Koniec pastucha

	spriteTable[unitMage][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
	}
	spriteTable[unitMage][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
	}
	spriteTable[unitMage][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
	}
	spriteTable[unitMage][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
	}
	spriteTable[unitMage][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
	}
	spriteTable[unitMage][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
	}
	spriteTable[unitMage][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
	}
	spriteTable[unitMage][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
	}
	spriteTable[unitMage][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
	}

	spriteTable[unitMage][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageMove1Down,
		spriteManualMageMove1Down,
		spriteManualMageMove1Down,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageMove2Down,
		spriteManualMageMove2Down,
		spriteManualMageMove2Down,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageMove1Down,
		spriteManualMageMove1Down,
		spriteManualMageIdleDown,
	}
	spriteTable[unitMage][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageMove1Down,
		spriteManualMageMove1Down,
		spriteManualMageMove1Down,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageMove2Down,
		spriteManualMageMove2Down,
		spriteManualMageMove2Down,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageMove1Down,
		spriteManualMageMove1Down,
		spriteManualMageIdleDown,
	}
	spriteTable[unitMage][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageMove1Left,
		spriteManualMageMove1Left,
		spriteManualMageMove1Left,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageMove2Left,
		spriteManualMageMove2Left,
		spriteManualMageMove2Left,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageMove1Left,
		spriteManualMageMove1Left,
		spriteManualMageIdleLeft,
	}
	spriteTable[unitMage][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageMove1Up,
		spriteManualMageMove1Up,
		spriteManualMageMove1Up,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageMove2Up,
		spriteManualMageMove2Up,
		spriteManualMageMove2Up,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageMove1Up,
		spriteManualMageMove1Up,
		spriteManualMageIdleUp,
	}
	spriteTable[unitMage][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageMove1Right,
		spriteManualMageMove1Right,
		spriteManualMageMove1Right,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageMove2Right,
		spriteManualMageMove2Right,
		spriteManualMageMove2Right,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageMove1Right,
		spriteManualMageMove1Right,
		spriteManualMageIdleRight,
	}
	spriteTable[unitMage][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageMove1UpLeft,
		spriteManualMageMove1UpLeft,
		spriteManualMageMove1UpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageMove2UpLeft,
		spriteManualMageMove2UpLeft,
		spriteManualMageMove2UpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageMove1UpLeft,
		spriteManualMageMove1UpLeft,
		spriteManualMageIdleUpLeft,
	}
	spriteTable[unitMage][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageMove1UpRight,
		spriteManualMageMove1UpRight,
		spriteManualMageMove1UpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageMove2UpRight,
		spriteManualMageMove2UpRight,
		spriteManualMageMove2UpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageMove1UpRight,
		spriteManualMageMove1UpRight,
		spriteManualMageIdleUpRight,
	}
	spriteTable[unitMage][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageMove1DownLeft,
		spriteManualMageMove1DownLeft,
		spriteManualMageMove1DownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageMove2DownLeft,
		spriteManualMageMove2DownLeft,
		spriteManualMageMove2DownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageMove1DownLeft,
		spriteManualMageMove1DownLeft,
		spriteManualMageIdleDownLeft,
	}
	spriteTable[unitMage][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageMove1DownRight,
		spriteManualMageMove1DownRight,
		spriteManualMageMove1DownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageMove2DownRight,
		spriteManualMageMove2DownRight,
		spriteManualMageMove2DownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageMove1DownRight,
		spriteManualMageMove1DownRight,
		spriteManualMageIdleDownRight,
	}

	spriteTable[unitMage][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageAttack2Down,
		spriteManualMageAttack2Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
	}
	spriteTable[unitMage][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageIdleDown,
		spriteManualMageAttack2Down,
		spriteManualMageAttack2Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
		spriteManualMageAttack1Down,
	}
	spriteTable[unitMage][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageIdleLeft,
		spriteManualMageAttack2Left,
		spriteManualMageAttack2Left,
		spriteManualMageAttack1Left,
		spriteManualMageAttack1Left,
		spriteManualMageAttack1Left,
		spriteManualMageAttack1Left,
		spriteManualMageAttack1Left,
		spriteManualMageAttack1Left,
		spriteManualMageAttack1Left,
		spriteManualMageAttack1Left,
		spriteManualMageAttack1Left,
		spriteManualMageAttack1Left,
		spriteManualMageAttack1Left,
		spriteManualMageAttack1Left,
	}
	spriteTable[unitMage][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageIdleUp,
		spriteManualMageAttack2Up,
		spriteManualMageAttack2Up,
		spriteManualMageAttack1Up,
		spriteManualMageAttack1Up,
		spriteManualMageAttack1Up,
		spriteManualMageAttack1Up,
		spriteManualMageAttack1Up,
		spriteManualMageAttack1Up,
		spriteManualMageAttack1Up,
		spriteManualMageAttack1Up,
		spriteManualMageAttack1Up,
		spriteManualMageAttack1Up,
		spriteManualMageAttack1Up,
		spriteManualMageAttack1Up,
	}
	spriteTable[unitMage][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageIdleRight,
		spriteManualMageAttack2Right,
		spriteManualMageAttack2Right,
		spriteManualMageAttack1Right,
		spriteManualMageAttack1Right,
		spriteManualMageAttack1Right,
		spriteManualMageAttack1Right,
		spriteManualMageAttack1Right,
		spriteManualMageAttack1Right,
		spriteManualMageAttack1Right,
		spriteManualMageAttack1Right,
		spriteManualMageAttack1Right,
		spriteManualMageAttack1Right,
		spriteManualMageAttack1Right,
		spriteManualMageAttack1Right,
	}
	spriteTable[unitMage][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageIdleUpLeft,
		spriteManualMageAttack2UpLeft,
		spriteManualMageAttack2UpLeft,
		spriteManualMageAttack1UpLeft,
		spriteManualMageAttack1UpLeft,
		spriteManualMageAttack1UpLeft,
		spriteManualMageAttack1UpLeft,
		spriteManualMageAttack1UpLeft,
		spriteManualMageAttack1UpLeft,
		spriteManualMageAttack1UpLeft,
		spriteManualMageAttack1UpLeft,
		spriteManualMageAttack1UpLeft,
		spriteManualMageAttack1UpLeft,
		spriteManualMageAttack1UpLeft,
		spriteManualMageAttack1UpLeft,
	}
	spriteTable[unitMage][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageIdleUpRight,
		spriteManualMageAttack2UpRight,
		spriteManualMageAttack2UpRight,
		spriteManualMageAttack1UpRight,
		spriteManualMageAttack1UpRight,
		spriteManualMageAttack1UpRight,
		spriteManualMageAttack1UpRight,
		spriteManualMageAttack1UpRight,
		spriteManualMageAttack1UpRight,
		spriteManualMageAttack1UpRight,
		spriteManualMageAttack1UpRight,
		spriteManualMageAttack1UpRight,
		spriteManualMageAttack1UpRight,
		spriteManualMageAttack1UpRight,
		spriteManualMageAttack1UpRight,
	}
	spriteTable[unitMage][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageIdleDownLeft,
		spriteManualMageAttack2DownLeft,
		spriteManualMageAttack2DownLeft,
		spriteManualMageAttack1DownLeft,
		spriteManualMageAttack1DownLeft,
		spriteManualMageAttack1DownLeft,
		spriteManualMageAttack1DownLeft,
		spriteManualMageAttack1DownLeft,
		spriteManualMageAttack1DownLeft,
		spriteManualMageAttack1DownLeft,
		spriteManualMageAttack1DownLeft,
		spriteManualMageAttack1DownLeft,
		spriteManualMageAttack1DownLeft,
		spriteManualMageAttack1DownLeft,
		spriteManualMageAttack1DownLeft,
	}
	spriteTable[unitMage][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageIdleDownRight,
		spriteManualMageAttack2DownRight,
		spriteManualMageAttack2DownRight,
		spriteManualMageAttack1DownRight,
		spriteManualMageAttack1DownRight,
		spriteManualMageAttack1DownRight,
		spriteManualMageAttack1DownRight,
		spriteManualMageAttack1DownRight,
		spriteManualMageAttack1DownRight,
		spriteManualMageAttack1DownRight,
		spriteManualMageAttack1DownRight,
		spriteManualMageAttack1DownRight,
		spriteManualMageAttack1DownRight,
		spriteManualMageAttack1DownRight,
		spriteManualMageAttack1DownRight,
	}
	// Koniec maga
	// === KAPŁAN (Priest) ===
	spriteTable[unitPriest][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
	}
	spriteTable[unitPriest][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
	}
	spriteTable[unitPriest][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
	}
	spriteTable[unitPriest][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
	}
	spriteTable[unitPriest][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
	}
	spriteTable[unitPriest][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
	}
	spriteTable[unitPriest][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
	}
	spriteTable[unitPriest][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
	}
	spriteTable[unitPriest][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
	}

	// Walk: Phase[4] = {0,0,2,2,2,2,0,0,0,0,0,1,1,1,1,0,0}
	spriteTable[unitPriest][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestMove2Down,
		spriteManualPriestMove2Down,
		spriteManualPriestMove2Down,
		spriteManualPriestMove2Down,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestMove1Down,
		spriteManualPriestMove1Down,
		spriteManualPriestMove1Down,
		spriteManualPriestMove1Down,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
	}
	spriteTable[unitPriest][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestMove2Down,
		spriteManualPriestMove2Down,
		spriteManualPriestMove2Down,
		spriteManualPriestMove2Down,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestMove1Down,
		spriteManualPriestMove1Down,
		spriteManualPriestMove1Down,
		spriteManualPriestMove1Down,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
	}
	spriteTable[unitPriest][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestMove2Left,
		spriteManualPriestMove2Left,
		spriteManualPriestMove2Left,
		spriteManualPriestMove2Left,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestMove1Left,
		spriteManualPriestMove1Left,
		spriteManualPriestMove1Left,
		spriteManualPriestMove1Left,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
	}
	spriteTable[unitPriest][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestMove2Up,
		spriteManualPriestMove2Up,
		spriteManualPriestMove2Up,
		spriteManualPriestMove2Up,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestMove1Up,
		spriteManualPriestMove1Up,
		spriteManualPriestMove1Up,
		spriteManualPriestMove1Up,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
	}
	spriteTable[unitPriest][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestMove2Right,
		spriteManualPriestMove2Right,
		spriteManualPriestMove2Right,
		spriteManualPriestMove2Right,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestMove1Right,
		spriteManualPriestMove1Right,
		spriteManualPriestMove1Right,
		spriteManualPriestMove1Right,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
	}
	spriteTable[unitPriest][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestMove2UpLeft,
		spriteManualPriestMove2UpLeft,
		spriteManualPriestMove2UpLeft,
		spriteManualPriestMove2UpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestMove1UpLeft,
		spriteManualPriestMove1UpLeft,
		spriteManualPriestMove1UpLeft,
		spriteManualPriestMove1UpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
	}
	spriteTable[unitPriest][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestMove2UpRight,
		spriteManualPriestMove2UpRight,
		spriteManualPriestMove2UpRight,
		spriteManualPriestMove2UpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestMove1UpRight,
		spriteManualPriestMove1UpRight,
		spriteManualPriestMove1UpRight,
		spriteManualPriestMove1UpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
	}
	spriteTable[unitPriest][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestMove2DownLeft,
		spriteManualPriestMove2DownLeft,
		spriteManualPriestMove2DownLeft,
		spriteManualPriestMove2DownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestMove1DownLeft,
		spriteManualPriestMove1DownLeft,
		spriteManualPriestMove1DownLeft,
		spriteManualPriestMove1DownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
	}
	spriteTable[unitPriest][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestMove2DownRight,
		spriteManualPriestMove2DownRight,
		spriteManualPriestMove2DownRight,
		spriteManualPriestMove2DownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestMove1DownRight,
		spriteManualPriestMove1DownRight,
		spriteManualPriestMove1DownRight,
		spriteManualPriestMove1DownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
	}

	// Fight: triggerHit=13, triggerReturn=9
	// delay 0..9 Idle, 10..13 Attack2, 14..16 Attack1
	spriteTable[unitPriest][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestAttack2Down,
		spriteManualPriestAttack2Down,
		spriteManualPriestAttack2Down,
		spriteManualPriestAttack2Down,
		spriteManualPriestAttack1Down,
		spriteManualPriestAttack1Down,
		spriteManualPriestAttack1Down,
	}
	spriteTable[unitPriest][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestIdleDown,
		spriteManualPriestAttack2Down,
		spriteManualPriestAttack2Down,
		spriteManualPriestAttack2Down,
		spriteManualPriestAttack2Down,
		spriteManualPriestAttack1Down,
		spriteManualPriestAttack1Down,
		spriteManualPriestAttack1Down,
	}
	spriteTable[unitPriest][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestIdleLeft,
		spriteManualPriestAttack2Left,
		spriteManualPriestAttack2Left,
		spriteManualPriestAttack2Left,
		spriteManualPriestAttack2Left,
		spriteManualPriestAttack1Left,
		spriteManualPriestAttack1Left,
		spriteManualPriestAttack1Left,
	}
	spriteTable[unitPriest][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestIdleUp,
		spriteManualPriestAttack2Up,
		spriteManualPriestAttack2Up,
		spriteManualPriestAttack2Up,
		spriteManualPriestAttack2Up,
		spriteManualPriestAttack1Up,
		spriteManualPriestAttack1Up,
		spriteManualPriestAttack1Up,
	}
	spriteTable[unitPriest][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestIdleRight,
		spriteManualPriestAttack2Right,
		spriteManualPriestAttack2Right,
		spriteManualPriestAttack2Right,
		spriteManualPriestAttack2Right,
		spriteManualPriestAttack1Right,
		spriteManualPriestAttack1Right,
		spriteManualPriestAttack1Right,
	}
	spriteTable[unitPriest][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestIdleUpLeft,
		spriteManualPriestAttack2UpLeft,
		spriteManualPriestAttack2UpLeft,
		spriteManualPriestAttack2UpLeft,
		spriteManualPriestAttack2UpLeft,
		spriteManualPriestAttack1UpLeft,
		spriteManualPriestAttack1UpLeft,
		spriteManualPriestAttack1UpLeft,
	}
	spriteTable[unitPriest][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestIdleUpRight,
		spriteManualPriestAttack2UpRight,
		spriteManualPriestAttack2UpRight,
		spriteManualPriestAttack2UpRight,
		spriteManualPriestAttack2UpRight,
		spriteManualPriestAttack1UpRight,
		spriteManualPriestAttack1UpRight,
		spriteManualPriestAttack1UpRight,
	}
	spriteTable[unitPriest][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestIdleDownLeft,
		spriteManualPriestAttack2DownLeft,
		spriteManualPriestAttack2DownLeft,
		spriteManualPriestAttack2DownLeft,
		spriteManualPriestAttack2DownLeft,
		spriteManualPriestAttack1DownLeft,
		spriteManualPriestAttack1DownLeft,
		spriteManualPriestAttack1DownLeft,
	}
	spriteTable[unitPriest][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestIdleDownRight,
		spriteManualPriestAttack2DownRight,
		spriteManualPriestAttack2DownRight,
		spriteManualPriestAttack2DownRight,
		spriteManualPriestAttack2DownRight,
		spriteManualPriestAttack1DownRight,
		spriteManualPriestAttack1DownRight,
		spriteManualPriestAttack1DownRight,
	}
	// Koniec Kapłana
	spriteTable[unitSpearman][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
	}
	spriteTable[unitSpearman][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
	}
	spriteTable[unitSpearman][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
	}
	spriteTable[unitSpearman][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
	}
	spriteTable[unitSpearman][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
	}
	spriteTable[unitSpearman][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
	}
	spriteTable[unitSpearman][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
	}
	spriteTable[unitSpearman][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
	}
	spriteTable[unitSpearman][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
	}

	spriteTable[unitSpearman][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanMove2Down,
		spriteManualSpearmanMove2Down,
		spriteManualSpearmanMove2Down,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanMove1Down,
		spriteManualSpearmanMove1Down,
		spriteManualSpearmanMove1Down,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
	}
	spriteTable[unitSpearman][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanMove2Down,
		spriteManualSpearmanMove2Down,
		spriteManualSpearmanMove2Down,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanMove1Down,
		spriteManualSpearmanMove1Down,
		spriteManualSpearmanMove1Down,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
	}
	spriteTable[unitSpearman][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanMove2Left,
		spriteManualSpearmanMove2Left,
		spriteManualSpearmanMove2Left,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanMove1Left,
		spriteManualSpearmanMove1Left,
		spriteManualSpearmanMove1Left,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
	}
	spriteTable[unitSpearman][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanMove2Up,
		spriteManualSpearmanMove2Up,
		spriteManualSpearmanMove2Up,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanMove1Up,
		spriteManualSpearmanMove1Up,
		spriteManualSpearmanMove1Up,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
	}
	spriteTable[unitSpearman][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanMove2Right,
		spriteManualSpearmanMove2Right,
		spriteManualSpearmanMove2Right,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanMove1Right,
		spriteManualSpearmanMove1Right,
		spriteManualSpearmanMove1Right,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
	}
	spriteTable[unitSpearman][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanMove2UpLeft,
		spriteManualSpearmanMove2UpLeft,
		spriteManualSpearmanMove2UpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanMove1UpLeft,
		spriteManualSpearmanMove1UpLeft,
		spriteManualSpearmanMove1UpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
	}
	spriteTable[unitSpearman][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanMove2UpRight,
		spriteManualSpearmanMove2UpRight,
		spriteManualSpearmanMove2UpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanMove1UpRight,
		spriteManualSpearmanMove1UpRight,
		spriteManualSpearmanMove1UpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
	}
	spriteTable[unitSpearman][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanMove2DownLeft,
		spriteManualSpearmanMove2DownLeft,
		spriteManualSpearmanMove2DownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanMove1DownLeft,
		spriteManualSpearmanMove1DownLeft,
		spriteManualSpearmanMove1DownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
	}
	spriteTable[unitSpearman][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanMove2DownRight,
		spriteManualSpearmanMove2DownRight,
		spriteManualSpearmanMove2DownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanMove1DownRight,
		spriteManualSpearmanMove1DownRight,
		spriteManualSpearmanMove1DownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
	}

	spriteTable[unitSpearman][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanAttack2Down,
		spriteManualSpearmanAttack2Down,
		spriteManualSpearmanAttack2Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
	}
	spriteTable[unitSpearman][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanIdleDown,
		spriteManualSpearmanAttack2Down,
		spriteManualSpearmanAttack2Down,
		spriteManualSpearmanAttack2Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
		spriteManualSpearmanAttack1Down,
	}
	spriteTable[unitSpearman][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanIdleLeft,
		spriteManualSpearmanAttack2Left,
		spriteManualSpearmanAttack2Left,
		spriteManualSpearmanAttack2Left,
		spriteManualSpearmanAttack1Left,
		spriteManualSpearmanAttack1Left,
		spriteManualSpearmanAttack1Left,
		spriteManualSpearmanAttack1Left,
		spriteManualSpearmanAttack1Left,
		spriteManualSpearmanAttack1Left,
		spriteManualSpearmanAttack1Left,
		spriteManualSpearmanAttack1Left,
		spriteManualSpearmanAttack1Left,
		spriteManualSpearmanAttack1Left,
	}
	spriteTable[unitSpearman][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanIdleUp,
		spriteManualSpearmanAttack2Up,
		spriteManualSpearmanAttack2Up,
		spriteManualSpearmanAttack2Up,
		spriteManualSpearmanAttack1Up,
		spriteManualSpearmanAttack1Up,
		spriteManualSpearmanAttack1Up,
		spriteManualSpearmanAttack1Up,
		spriteManualSpearmanAttack1Up,
		spriteManualSpearmanAttack1Up,
		spriteManualSpearmanAttack1Up,
		spriteManualSpearmanAttack1Up,
		spriteManualSpearmanAttack1Up,
		spriteManualSpearmanAttack1Up,
	}
	spriteTable[unitSpearman][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanIdleRight,
		spriteManualSpearmanAttack2Right,
		spriteManualSpearmanAttack2Right,
		spriteManualSpearmanAttack2Right,
		spriteManualSpearmanAttack1Right,
		spriteManualSpearmanAttack1Right,
		spriteManualSpearmanAttack1Right,
		spriteManualSpearmanAttack1Right,
		spriteManualSpearmanAttack1Right,
		spriteManualSpearmanAttack1Right,
		spriteManualSpearmanAttack1Right,
		spriteManualSpearmanAttack1Right,
		spriteManualSpearmanAttack1Right,
		spriteManualSpearmanAttack1Right,
	}
	spriteTable[unitSpearman][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanIdleUpLeft,
		spriteManualSpearmanAttack2UpLeft,
		spriteManualSpearmanAttack2UpLeft,
		spriteManualSpearmanAttack2UpLeft,
		spriteManualSpearmanAttack1UpLeft,
		spriteManualSpearmanAttack1UpLeft,
		spriteManualSpearmanAttack1UpLeft,
		spriteManualSpearmanAttack1UpLeft,
		spriteManualSpearmanAttack1UpLeft,
		spriteManualSpearmanAttack1UpLeft,
		spriteManualSpearmanAttack1UpLeft,
		spriteManualSpearmanAttack1UpLeft,
		spriteManualSpearmanAttack1UpLeft,
		spriteManualSpearmanAttack1UpLeft,
	}
	spriteTable[unitSpearman][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanIdleUpRight,
		spriteManualSpearmanAttack2UpRight,
		spriteManualSpearmanAttack2UpRight,
		spriteManualSpearmanAttack2UpRight,
		spriteManualSpearmanAttack1UpRight,
		spriteManualSpearmanAttack1UpRight,
		spriteManualSpearmanAttack1UpRight,
		spriteManualSpearmanAttack1UpRight,
		spriteManualSpearmanAttack1UpRight,
		spriteManualSpearmanAttack1UpRight,
		spriteManualSpearmanAttack1UpRight,
		spriteManualSpearmanAttack1UpRight,
		spriteManualSpearmanAttack1UpRight,
		spriteManualSpearmanAttack1UpRight,
	}
	spriteTable[unitSpearman][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanIdleDownLeft,
		spriteManualSpearmanAttack2DownLeft,
		spriteManualSpearmanAttack2DownLeft,
		spriteManualSpearmanAttack2DownLeft,
		spriteManualSpearmanAttack1DownLeft,
		spriteManualSpearmanAttack1DownLeft,
		spriteManualSpearmanAttack1DownLeft,
		spriteManualSpearmanAttack1DownLeft,
		spriteManualSpearmanAttack1DownLeft,
		spriteManualSpearmanAttack1DownLeft,
		spriteManualSpearmanAttack1DownLeft,
		spriteManualSpearmanAttack1DownLeft,
		spriteManualSpearmanAttack1DownLeft,
		spriteManualSpearmanAttack1DownLeft,
	}
	spriteTable[unitSpearman][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanIdleDownRight,
		spriteManualSpearmanAttack2DownRight,
		spriteManualSpearmanAttack2DownRight,
		spriteManualSpearmanAttack2DownRight,
		spriteManualSpearmanAttack1DownRight,
		spriteManualSpearmanAttack1DownRight,
		spriteManualSpearmanAttack1DownRight,
		spriteManualSpearmanAttack1DownRight,
		spriteManualSpearmanAttack1DownRight,
		spriteManualSpearmanAttack1DownRight,
		spriteManualSpearmanAttack1DownRight,
		spriteManualSpearmanAttack1DownRight,
		spriteManualSpearmanAttack1DownRight,
		spriteManualSpearmanAttack1DownRight,
	}
	// Koniec Włócznika

	spriteTable[unitCrossbowman][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
	}
	spriteTable[unitCrossbowman][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
	}
	spriteTable[unitCrossbowman][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
	}
	spriteTable[unitCrossbowman][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
	}
	spriteTable[unitCrossbowman][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
	}
	spriteTable[unitCrossbowman][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
	}
	spriteTable[unitCrossbowman][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
	}
	spriteTable[unitCrossbowman][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
	}
	spriteTable[unitCrossbowman][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
	}

	spriteTable[unitCrossbowman][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanMove2Down,
		spriteManualCrossbowmanMove2Down,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanMove1Down,
		spriteManualCrossbowmanMove1Down,
		spriteManualCrossbowmanMove1Down,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
	}
	spriteTable[unitCrossbowman][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanMove2Down,
		spriteManualCrossbowmanMove2Down,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanMove1Down,
		spriteManualCrossbowmanMove1Down,
		spriteManualCrossbowmanMove1Down,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
	}
	spriteTable[unitCrossbowman][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanMove2Left,
		spriteManualCrossbowmanMove2Left,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanMove1Left,
		spriteManualCrossbowmanMove1Left,
		spriteManualCrossbowmanMove1Left,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
	}
	spriteTable[unitCrossbowman][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanMove2Up,
		spriteManualCrossbowmanMove2Up,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanMove1Up,
		spriteManualCrossbowmanMove1Up,
		spriteManualCrossbowmanMove1Up,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
	}
	spriteTable[unitCrossbowman][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanMove2Right,
		spriteManualCrossbowmanMove2Right,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanMove1Right,
		spriteManualCrossbowmanMove1Right,
		spriteManualCrossbowmanMove1Right,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
	}
	spriteTable[unitCrossbowman][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanMove2UpLeft,
		spriteManualCrossbowmanMove2UpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanMove1UpLeft,
		spriteManualCrossbowmanMove1UpLeft,
		spriteManualCrossbowmanMove1UpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
	}
	spriteTable[unitCrossbowman][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanMove2UpRight,
		spriteManualCrossbowmanMove2UpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanMove1UpRight,
		spriteManualCrossbowmanMove1UpRight,
		spriteManualCrossbowmanMove1UpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
	}
	spriteTable[unitCrossbowman][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanMove2DownLeft,
		spriteManualCrossbowmanMove2DownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanMove1DownLeft,
		spriteManualCrossbowmanMove1DownLeft,
		spriteManualCrossbowmanMove1DownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
	}
	spriteTable[unitCrossbowman][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanMove2DownRight,
		spriteManualCrossbowmanMove2DownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanMove1DownRight,
		spriteManualCrossbowmanMove1DownRight,
		spriteManualCrossbowmanMove1DownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
	}

	spriteTable[unitCrossbowman][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanAttack2Down,
		spriteManualCrossbowmanAttack2Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
	}
	spriteTable[unitCrossbowman][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanIdleDown,
		spriteManualCrossbowmanAttack2Down,
		spriteManualCrossbowmanAttack2Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
		spriteManualCrossbowmanAttack1Down,
	}
	spriteTable[unitCrossbowman][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanIdleLeft,
		spriteManualCrossbowmanAttack2Left,
		spriteManualCrossbowmanAttack2Left,
		spriteManualCrossbowmanAttack1Left,
		spriteManualCrossbowmanAttack1Left,
		spriteManualCrossbowmanAttack1Left,
		spriteManualCrossbowmanAttack1Left,
		spriteManualCrossbowmanAttack1Left,
		spriteManualCrossbowmanAttack1Left,
		spriteManualCrossbowmanAttack1Left,
		spriteManualCrossbowmanAttack1Left,
		spriteManualCrossbowmanAttack1Left,
		spriteManualCrossbowmanAttack1Left,
		spriteManualCrossbowmanAttack1Left,
		spriteManualCrossbowmanAttack1Left,
	}
	spriteTable[unitCrossbowman][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanIdleUp,
		spriteManualCrossbowmanAttack2Up,
		spriteManualCrossbowmanAttack2Up,
		spriteManualCrossbowmanAttack1Up,
		spriteManualCrossbowmanAttack1Up,
		spriteManualCrossbowmanAttack1Up,
		spriteManualCrossbowmanAttack1Up,
		spriteManualCrossbowmanAttack1Up,
		spriteManualCrossbowmanAttack1Up,
		spriteManualCrossbowmanAttack1Up,
		spriteManualCrossbowmanAttack1Up,
		spriteManualCrossbowmanAttack1Up,
		spriteManualCrossbowmanAttack1Up,
		spriteManualCrossbowmanAttack1Up,
		spriteManualCrossbowmanAttack1Up,
	}
	spriteTable[unitCrossbowman][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanIdleRight,
		spriteManualCrossbowmanAttack2Right,
		spriteManualCrossbowmanAttack2Right,
		spriteManualCrossbowmanAttack1Right,
		spriteManualCrossbowmanAttack1Right,
		spriteManualCrossbowmanAttack1Right,
		spriteManualCrossbowmanAttack1Right,
		spriteManualCrossbowmanAttack1Right,
		spriteManualCrossbowmanAttack1Right,
		spriteManualCrossbowmanAttack1Right,
		spriteManualCrossbowmanAttack1Right,
		spriteManualCrossbowmanAttack1Right,
		spriteManualCrossbowmanAttack1Right,
		spriteManualCrossbowmanAttack1Right,
		spriteManualCrossbowmanAttack1Right,
	}
	spriteTable[unitCrossbowman][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanIdleUpLeft,
		spriteManualCrossbowmanAttack2UpLeft,
		spriteManualCrossbowmanAttack2UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
		spriteManualCrossbowmanAttack1UpLeft,
	}
	spriteTable[unitCrossbowman][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanIdleUpRight,
		spriteManualCrossbowmanAttack2UpRight,
		spriteManualCrossbowmanAttack2UpRight,
		spriteManualCrossbowmanAttack1UpRight,
		spriteManualCrossbowmanAttack1UpRight,
		spriteManualCrossbowmanAttack1UpRight,
		spriteManualCrossbowmanAttack1UpRight,
		spriteManualCrossbowmanAttack1UpRight,
		spriteManualCrossbowmanAttack1UpRight,
		spriteManualCrossbowmanAttack1UpRight,
		spriteManualCrossbowmanAttack1UpRight,
		spriteManualCrossbowmanAttack1UpRight,
		spriteManualCrossbowmanAttack1UpRight,
		spriteManualCrossbowmanAttack1UpRight,
		spriteManualCrossbowmanAttack1UpRight,
	}
	spriteTable[unitCrossbowman][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanIdleDownLeft,
		spriteManualCrossbowmanAttack2DownLeft,
		spriteManualCrossbowmanAttack2DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
		spriteManualCrossbowmanAttack1DownLeft,
	}
	spriteTable[unitCrossbowman][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanIdleDownRight,
		spriteManualCrossbowmanAttack2DownRight,
		spriteManualCrossbowmanAttack2DownRight,
		spriteManualCrossbowmanAttack1DownRight,
		spriteManualCrossbowmanAttack1DownRight,
		spriteManualCrossbowmanAttack1DownRight,
		spriteManualCrossbowmanAttack1DownRight,
		spriteManualCrossbowmanAttack1DownRight,
		spriteManualCrossbowmanAttack1DownRight,
		spriteManualCrossbowmanAttack1DownRight,
		spriteManualCrossbowmanAttack1DownRight,
		spriteManualCrossbowmanAttack1DownRight,
		spriteManualCrossbowmanAttack1DownRight,
		spriteManualCrossbowmanAttack1DownRight,
		spriteManualCrossbowmanAttack1DownRight,
	}
	// Koniec Kusznika

	spriteTable[unitAxeman][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
	}
	spriteTable[unitAxeman][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
	}
	spriteTable[unitAxeman][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
	}
	spriteTable[unitAxeman][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
	}
	spriteTable[unitAxeman][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
	}
	spriteTable[unitAxeman][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
	}
	spriteTable[unitAxeman][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
	}
	spriteTable[unitAxeman][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
	}
	spriteTable[unitAxeman][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
	}

	spriteTable[unitAxeman][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanMove1Down,
		spriteManualAxemanMove1Down,
		spriteManualAxemanMove1Down,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanMove2Down,
		spriteManualAxemanMove2Down,
		spriteManualAxemanMove2Down,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
	}
	spriteTable[unitAxeman][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanMove1Down,
		spriteManualAxemanMove1Down,
		spriteManualAxemanMove1Down,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanMove2Down,
		spriteManualAxemanMove2Down,
		spriteManualAxemanMove2Down,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
	}
	spriteTable[unitAxeman][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanMove1Left,
		spriteManualAxemanMove1Left,
		spriteManualAxemanMove1Left,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanMove2Left,
		spriteManualAxemanMove2Left,
		spriteManualAxemanMove2Left,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
	}
	spriteTable[unitAxeman][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanMove1Up,
		spriteManualAxemanMove1Up,
		spriteManualAxemanMove1Up,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanMove2Up,
		spriteManualAxemanMove2Up,
		spriteManualAxemanMove2Up,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
	}
	spriteTable[unitAxeman][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanMove1Right,
		spriteManualAxemanMove1Right,
		spriteManualAxemanMove1Right,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanMove2Right,
		spriteManualAxemanMove2Right,
		spriteManualAxemanMove2Right,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
	}
	spriteTable[unitAxeman][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanMove1UpLeft,
		spriteManualAxemanMove1UpLeft,
		spriteManualAxemanMove1UpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanMove2UpLeft,
		spriteManualAxemanMove2UpLeft,
		spriteManualAxemanMove2UpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
	}
	spriteTable[unitAxeman][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanMove1UpRight,
		spriteManualAxemanMove1UpRight,
		spriteManualAxemanMove1UpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanMove2UpRight,
		spriteManualAxemanMove2UpRight,
		spriteManualAxemanMove2UpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
	}
	spriteTable[unitAxeman][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanMove1DownLeft,
		spriteManualAxemanMove1DownLeft,
		spriteManualAxemanMove1DownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanMove2DownLeft,
		spriteManualAxemanMove2DownLeft,
		spriteManualAxemanMove2DownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
	}
	spriteTable[unitAxeman][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanMove1DownRight,
		spriteManualAxemanMove1DownRight,
		spriteManualAxemanMove1DownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanMove2DownRight,
		spriteManualAxemanMove2DownRight,
		spriteManualAxemanMove2DownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
	}

	spriteTable[unitAxeman][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanAttack2Down,
		spriteManualAxemanAttack2Down,
		spriteManualAxemanAttack2Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
	}
	spriteTable[unitAxeman][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanIdleDown,
		spriteManualAxemanAttack2Down,
		spriteManualAxemanAttack2Down,
		spriteManualAxemanAttack2Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
		spriteManualAxemanAttack1Down,
	}
	spriteTable[unitAxeman][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanIdleLeft,
		spriteManualAxemanAttack2Left,
		spriteManualAxemanAttack2Left,
		spriteManualAxemanAttack2Left,
		spriteManualAxemanAttack1Left,
		spriteManualAxemanAttack1Left,
		spriteManualAxemanAttack1Left,
		spriteManualAxemanAttack1Left,
		spriteManualAxemanAttack1Left,
		spriteManualAxemanAttack1Left,
		spriteManualAxemanAttack1Left,
		spriteManualAxemanAttack1Left,
		spriteManualAxemanAttack1Left,
		spriteManualAxemanAttack1Left,
	}
	spriteTable[unitAxeman][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanIdleUp,
		spriteManualAxemanAttack2Up,
		spriteManualAxemanAttack2Up,
		spriteManualAxemanAttack2Up,
		spriteManualAxemanAttack1Up,
		spriteManualAxemanAttack1Up,
		spriteManualAxemanAttack1Up,
		spriteManualAxemanAttack1Up,
		spriteManualAxemanAttack1Up,
		spriteManualAxemanAttack1Up,
		spriteManualAxemanAttack1Up,
		spriteManualAxemanAttack1Up,
		spriteManualAxemanAttack1Up,
		spriteManualAxemanAttack1Up,
	}
	spriteTable[unitAxeman][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanIdleRight,
		spriteManualAxemanAttack2Right,
		spriteManualAxemanAttack2Right,
		spriteManualAxemanAttack2Right,
		spriteManualAxemanAttack1Right,
		spriteManualAxemanAttack1Right,
		spriteManualAxemanAttack1Right,
		spriteManualAxemanAttack1Right,
		spriteManualAxemanAttack1Right,
		spriteManualAxemanAttack1Right,
		spriteManualAxemanAttack1Right,
		spriteManualAxemanAttack1Right,
		spriteManualAxemanAttack1Right,
		spriteManualAxemanAttack1Right,
	}
	spriteTable[unitAxeman][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanIdleUpLeft,
		spriteManualAxemanAttack2UpLeft,
		spriteManualAxemanAttack2UpLeft,
		spriteManualAxemanAttack2UpLeft,
		spriteManualAxemanAttack1UpLeft,
		spriteManualAxemanAttack1UpLeft,
		spriteManualAxemanAttack1UpLeft,
		spriteManualAxemanAttack1UpLeft,
		spriteManualAxemanAttack1UpLeft,
		spriteManualAxemanAttack1UpLeft,
		spriteManualAxemanAttack1UpLeft,
		spriteManualAxemanAttack1UpLeft,
		spriteManualAxemanAttack1UpLeft,
		spriteManualAxemanAttack1UpLeft,
	}
	spriteTable[unitAxeman][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanIdleUpRight,
		spriteManualAxemanAttack2UpRight,
		spriteManualAxemanAttack2UpRight,
		spriteManualAxemanAttack2UpRight,
		spriteManualAxemanAttack1UpRight,
		spriteManualAxemanAttack1UpRight,
		spriteManualAxemanAttack1UpRight,
		spriteManualAxemanAttack1UpRight,
		spriteManualAxemanAttack1UpRight,
		spriteManualAxemanAttack1UpRight,
		spriteManualAxemanAttack1UpRight,
		spriteManualAxemanAttack1UpRight,
		spriteManualAxemanAttack1UpRight,
		spriteManualAxemanAttack1UpRight,
	}
	spriteTable[unitAxeman][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanIdleDownLeft,
		spriteManualAxemanAttack2DownLeft,
		spriteManualAxemanAttack2DownLeft,
		spriteManualAxemanAttack2DownLeft,
		spriteManualAxemanAttack1DownLeft,
		spriteManualAxemanAttack1DownLeft,
		spriteManualAxemanAttack1DownLeft,
		spriteManualAxemanAttack1DownLeft,
		spriteManualAxemanAttack1DownLeft,
		spriteManualAxemanAttack1DownLeft,
		spriteManualAxemanAttack1DownLeft,
		spriteManualAxemanAttack1DownLeft,
		spriteManualAxemanAttack1DownLeft,
		spriteManualAxemanAttack1DownLeft,
	}
	spriteTable[unitAxeman][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanIdleDownRight,
		spriteManualAxemanAttack2DownRight,
		spriteManualAxemanAttack2DownRight,
		spriteManualAxemanAttack2DownRight,
		spriteManualAxemanAttack1DownRight,
		spriteManualAxemanAttack1DownRight,
		spriteManualAxemanAttack1DownRight,
		spriteManualAxemanAttack1DownRight,
		spriteManualAxemanAttack1DownRight,
		spriteManualAxemanAttack1DownRight,
		spriteManualAxemanAttack1DownRight,
		spriteManualAxemanAttack1DownRight,
		spriteManualAxemanAttack1DownRight,
		spriteManualAxemanAttack1DownRight,
	}
	// Koniec Drwala

	spriteTable[unitSwordsman][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
	}
	spriteTable[unitSwordsman][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
	}
	spriteTable[unitSwordsman][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
	}
	spriteTable[unitSwordsman][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
	}
	spriteTable[unitSwordsman][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
	}
	spriteTable[unitSwordsman][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
	}
	spriteTable[unitSwordsman][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
	}
	spriteTable[unitSwordsman][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
	}
	spriteTable[unitSwordsman][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
	}

	spriteTable[unitSwordsman][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanMove2Down,
		spriteManualSwordsmanMove2Down,
		spriteManualSwordsmanMove2Down,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanMove1Down,
		spriteManualSwordsmanMove1Down,
		spriteManualSwordsmanMove1Down,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
	}
	spriteTable[unitSwordsman][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanMove2Down,
		spriteManualSwordsmanMove2Down,
		spriteManualSwordsmanMove2Down,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanMove1Down,
		spriteManualSwordsmanMove1Down,
		spriteManualSwordsmanMove1Down,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
	}
	spriteTable[unitSwordsman][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanMove2Left,
		spriteManualSwordsmanMove2Left,
		spriteManualSwordsmanMove2Left,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanMove1Left,
		spriteManualSwordsmanMove1Left,
		spriteManualSwordsmanMove1Left,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
	}
	spriteTable[unitSwordsman][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanMove2Up,
		spriteManualSwordsmanMove2Up,
		spriteManualSwordsmanMove2Up,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanMove1Up,
		spriteManualSwordsmanMove1Up,
		spriteManualSwordsmanMove1Up,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
	}
	spriteTable[unitSwordsman][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanMove2Right,
		spriteManualSwordsmanMove2Right,
		spriteManualSwordsmanMove2Right,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanMove1Right,
		spriteManualSwordsmanMove1Right,
		spriteManualSwordsmanMove1Right,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
	}
	spriteTable[unitSwordsman][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanMove2UpLeft,
		spriteManualSwordsmanMove2UpLeft,
		spriteManualSwordsmanMove2UpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanMove1UpLeft,
		spriteManualSwordsmanMove1UpLeft,
		spriteManualSwordsmanMove1UpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
	}
	spriteTable[unitSwordsman][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanMove2UpRight,
		spriteManualSwordsmanMove2UpRight,
		spriteManualSwordsmanMove2UpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanMove1UpRight,
		spriteManualSwordsmanMove1UpRight,
		spriteManualSwordsmanMove1UpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
	}
	spriteTable[unitSwordsman][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanMove2DownLeft,
		spriteManualSwordsmanMove2DownLeft,
		spriteManualSwordsmanMove2DownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanMove1DownLeft,
		spriteManualSwordsmanMove1DownLeft,
		spriteManualSwordsmanMove1DownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
	}
	spriteTable[unitSwordsman][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanMove2DownRight,
		spriteManualSwordsmanMove2DownRight,
		spriteManualSwordsmanMove2DownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanMove1DownRight,
		spriteManualSwordsmanMove1DownRight,
		spriteManualSwordsmanMove1DownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
	}

	spriteTable[unitSwordsman][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanAttack2Down,
		spriteManualSwordsmanAttack2Down,
		spriteManualSwordsmanAttack2Down,
		spriteManualSwordsmanAttack2Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
	}
	spriteTable[unitSwordsman][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanIdleDown,
		spriteManualSwordsmanAttack2Down,
		spriteManualSwordsmanAttack2Down,
		spriteManualSwordsmanAttack2Down,
		spriteManualSwordsmanAttack2Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
		spriteManualSwordsmanAttack1Down,
	}
	spriteTable[unitSwordsman][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanIdleLeft,
		spriteManualSwordsmanAttack2Left,
		spriteManualSwordsmanAttack2Left,
		spriteManualSwordsmanAttack2Left,
		spriteManualSwordsmanAttack2Left,
		spriteManualSwordsmanAttack1Left,
		spriteManualSwordsmanAttack1Left,
		spriteManualSwordsmanAttack1Left,
		spriteManualSwordsmanAttack1Left,
		spriteManualSwordsmanAttack1Left,
		spriteManualSwordsmanAttack1Left,
		spriteManualSwordsmanAttack1Left,
		spriteManualSwordsmanAttack1Left,
		spriteManualSwordsmanAttack1Left,
		spriteManualSwordsmanAttack1Left,
	}
	spriteTable[unitSwordsman][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanIdleUp,
		spriteManualSwordsmanAttack2Up,
		spriteManualSwordsmanAttack2Up,
		spriteManualSwordsmanAttack2Up,
		spriteManualSwordsmanAttack2Up,
		spriteManualSwordsmanAttack1Up,
		spriteManualSwordsmanAttack1Up,
		spriteManualSwordsmanAttack1Up,
		spriteManualSwordsmanAttack1Up,
		spriteManualSwordsmanAttack1Up,
		spriteManualSwordsmanAttack1Up,
		spriteManualSwordsmanAttack1Up,
		spriteManualSwordsmanAttack1Up,
		spriteManualSwordsmanAttack1Up,
		spriteManualSwordsmanAttack1Up,
	}
	spriteTable[unitSwordsman][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanIdleRight,
		spriteManualSwordsmanAttack2Right,
		spriteManualSwordsmanAttack2Right,
		spriteManualSwordsmanAttack2Right,
		spriteManualSwordsmanAttack2Right,
		spriteManualSwordsmanAttack1Right,
		spriteManualSwordsmanAttack1Right,
		spriteManualSwordsmanAttack1Right,
		spriteManualSwordsmanAttack1Right,
		spriteManualSwordsmanAttack1Right,
		spriteManualSwordsmanAttack1Right,
		spriteManualSwordsmanAttack1Right,
		spriteManualSwordsmanAttack1Right,
		spriteManualSwordsmanAttack1Right,
		spriteManualSwordsmanAttack1Right,
	}
	spriteTable[unitSwordsman][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanIdleUpLeft,
		spriteManualSwordsmanAttack2UpLeft,
		spriteManualSwordsmanAttack2UpLeft,
		spriteManualSwordsmanAttack2UpLeft,
		spriteManualSwordsmanAttack2UpLeft,
		spriteManualSwordsmanAttack1UpLeft,
		spriteManualSwordsmanAttack1UpLeft,
		spriteManualSwordsmanAttack1UpLeft,
		spriteManualSwordsmanAttack1UpLeft,
		spriteManualSwordsmanAttack1UpLeft,
		spriteManualSwordsmanAttack1UpLeft,
		spriteManualSwordsmanAttack1UpLeft,
		spriteManualSwordsmanAttack1UpLeft,
		spriteManualSwordsmanAttack1UpLeft,
		spriteManualSwordsmanAttack1UpLeft,
	}
	spriteTable[unitSwordsman][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanIdleUpRight,
		spriteManualSwordsmanAttack2UpRight,
		spriteManualSwordsmanAttack2UpRight,
		spriteManualSwordsmanAttack2UpRight,
		spriteManualSwordsmanAttack2UpRight,
		spriteManualSwordsmanAttack1UpRight,
		spriteManualSwordsmanAttack1UpRight,
		spriteManualSwordsmanAttack1UpRight,
		spriteManualSwordsmanAttack1UpRight,
		spriteManualSwordsmanAttack1UpRight,
		spriteManualSwordsmanAttack1UpRight,
		spriteManualSwordsmanAttack1UpRight,
		spriteManualSwordsmanAttack1UpRight,
		spriteManualSwordsmanAttack1UpRight,
		spriteManualSwordsmanAttack1UpRight,
	}
	spriteTable[unitSwordsman][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanIdleDownLeft,
		spriteManualSwordsmanAttack2DownLeft,
		spriteManualSwordsmanAttack2DownLeft,
		spriteManualSwordsmanAttack2DownLeft,
		spriteManualSwordsmanAttack2DownLeft,
		spriteManualSwordsmanAttack1DownLeft,
		spriteManualSwordsmanAttack1DownLeft,
		spriteManualSwordsmanAttack1DownLeft,
		spriteManualSwordsmanAttack1DownLeft,
		spriteManualSwordsmanAttack1DownLeft,
		spriteManualSwordsmanAttack1DownLeft,
		spriteManualSwordsmanAttack1DownLeft,
		spriteManualSwordsmanAttack1DownLeft,
		spriteManualSwordsmanAttack1DownLeft,
		spriteManualSwordsmanAttack1DownLeft,
	}
	spriteTable[unitSwordsman][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanIdleDownRight,
		spriteManualSwordsmanAttack2DownRight,
		spriteManualSwordsmanAttack2DownRight,
		spriteManualSwordsmanAttack2DownRight,
		spriteManualSwordsmanAttack2DownRight,
		spriteManualSwordsmanAttack1DownRight,
		spriteManualSwordsmanAttack1DownRight,
		spriteManualSwordsmanAttack1DownRight,
		spriteManualSwordsmanAttack1DownRight,
		spriteManualSwordsmanAttack1DownRight,
		spriteManualSwordsmanAttack1DownRight,
		spriteManualSwordsmanAttack1DownRight,
		spriteManualSwordsmanAttack1DownRight,
		spriteManualSwordsmanAttack1DownRight,
		spriteManualSwordsmanAttack1DownRight,
	}
	// Koniec miecznika

	spriteTable[unitCommander][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
	}
	spriteTable[unitCommander][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
	}
	spriteTable[unitCommander][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
	}
	spriteTable[unitCommander][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
	}
	spriteTable[unitCommander][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
	}
	spriteTable[unitCommander][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
	}
	spriteTable[unitCommander][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
	}
	spriteTable[unitCommander][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
	}
	spriteTable[unitCommander][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
	}

	spriteTable[unitCommander][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderMove2Down,
		spriteManualCommanderMove2Down,
		spriteManualCommanderMove2Down,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderMove1Down,
		spriteManualCommanderMove1Down,
		spriteManualCommanderMove1Down,
		spriteManualCommanderMove1Down,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderMove1Down,
		spriteManualCommanderMove1Down,
		spriteManualCommanderIdleDown,
	}
	spriteTable[unitCommander][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderMove2Down,
		spriteManualCommanderMove2Down,
		spriteManualCommanderMove2Down,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderMove1Down,
		spriteManualCommanderMove1Down,
		spriteManualCommanderMove1Down,
		spriteManualCommanderMove1Down,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderMove1Down,
		spriteManualCommanderMove1Down,
		spriteManualCommanderIdleDown,
	}
	spriteTable[unitCommander][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderMove2Left,
		spriteManualCommanderMove2Left,
		spriteManualCommanderMove2Left,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderMove1Left,
		spriteManualCommanderMove1Left,
		spriteManualCommanderMove1Left,
		spriteManualCommanderMove1Left,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderMove1Left,
		spriteManualCommanderMove1Left,
		spriteManualCommanderIdleLeft,
	}
	spriteTable[unitCommander][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderMove2Up,
		spriteManualCommanderMove2Up,
		spriteManualCommanderMove2Up,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderMove1Up,
		spriteManualCommanderMove1Up,
		spriteManualCommanderMove1Up,
		spriteManualCommanderMove1Up,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderMove1Up,
		spriteManualCommanderMove1Up,
		spriteManualCommanderIdleUp,
	}
	spriteTable[unitCommander][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderMove2Right,
		spriteManualCommanderMove2Right,
		spriteManualCommanderMove2Right,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderMove1Right,
		spriteManualCommanderMove1Right,
		spriteManualCommanderMove1Right,
		spriteManualCommanderMove1Right,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderMove1Right,
		spriteManualCommanderMove1Right,
		spriteManualCommanderIdleRight,
	}
	spriteTable[unitCommander][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderMove2UpLeft,
		spriteManualCommanderMove2UpLeft,
		spriteManualCommanderMove2UpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderMove1UpLeft,
		spriteManualCommanderMove1UpLeft,
		spriteManualCommanderMove1UpLeft,
		spriteManualCommanderMove1UpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderMove1UpLeft,
		spriteManualCommanderMove1UpLeft,
		spriteManualCommanderIdleUpLeft,
	}
	spriteTable[unitCommander][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderMove2UpRight,
		spriteManualCommanderMove2UpRight,
		spriteManualCommanderMove2UpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderMove1UpRight,
		spriteManualCommanderMove1UpRight,
		spriteManualCommanderMove1UpRight,
		spriteManualCommanderMove1UpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderMove1UpRight,
		spriteManualCommanderMove1UpRight,
		spriteManualCommanderIdleUpRight,
	}
	spriteTable[unitCommander][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderMove2DownLeft,
		spriteManualCommanderMove2DownLeft,
		spriteManualCommanderMove2DownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderMove1DownLeft,
		spriteManualCommanderMove1DownLeft,
		spriteManualCommanderMove1DownLeft,
		spriteManualCommanderMove1DownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderMove1DownLeft,
		spriteManualCommanderMove1DownLeft,
		spriteManualCommanderIdleDownLeft,
	}
	spriteTable[unitCommander][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderMove2DownRight,
		spriteManualCommanderMove2DownRight,
		spriteManualCommanderMove2DownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderMove1DownRight,
		spriteManualCommanderMove1DownRight,
		spriteManualCommanderMove1DownRight,
		spriteManualCommanderMove1DownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderMove1DownRight,
		spriteManualCommanderMove1DownRight,
		spriteManualCommanderIdleDownRight,
	}

	spriteTable[unitCommander][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderAttack2Down,
		spriteManualCommanderAttack2Down,
		spriteManualCommanderAttack2Down,
		spriteManualCommanderAttack2Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
	}
	spriteTable[unitCommander][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDown,
		spriteManualCommanderIdleDown,
		spriteManualCommanderAttack2Down,
		spriteManualCommanderAttack2Down,
		spriteManualCommanderAttack2Down,
		spriteManualCommanderAttack2Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
		spriteManualCommanderAttack1Down,
	}
	spriteTable[unitCommander][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualCommanderIdleLeft,
		spriteManualCommanderIdleLeft,
		spriteManualCommanderAttack2Left,
		spriteManualCommanderAttack2Left,
		spriteManualCommanderAttack2Left,
		spriteManualCommanderAttack2Left,
		spriteManualCommanderAttack1Left,
		spriteManualCommanderAttack1Left,
		spriteManualCommanderAttack1Left,
		spriteManualCommanderAttack1Left,
		spriteManualCommanderAttack1Left,
		spriteManualCommanderAttack1Left,
		spriteManualCommanderAttack1Left,
		spriteManualCommanderAttack1Left,
		spriteManualCommanderAttack1Left,
		spriteManualCommanderAttack1Left,
		spriteManualCommanderAttack1Left,
	}
	spriteTable[unitCommander][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualCommanderIdleUp,
		spriteManualCommanderIdleUp,
		spriteManualCommanderAttack2Up,
		spriteManualCommanderAttack2Up,
		spriteManualCommanderAttack2Up,
		spriteManualCommanderAttack2Up,
		spriteManualCommanderAttack1Up,
		spriteManualCommanderAttack1Up,
		spriteManualCommanderAttack1Up,
		spriteManualCommanderAttack1Up,
		spriteManualCommanderAttack1Up,
		spriteManualCommanderAttack1Up,
		spriteManualCommanderAttack1Up,
		spriteManualCommanderAttack1Up,
		spriteManualCommanderAttack1Up,
		spriteManualCommanderAttack1Up,
		spriteManualCommanderAttack1Up,
	}
	spriteTable[unitCommander][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualCommanderIdleRight,
		spriteManualCommanderIdleRight,
		spriteManualCommanderAttack2Right,
		spriteManualCommanderAttack2Right,
		spriteManualCommanderAttack2Right,
		spriteManualCommanderAttack2Right,
		spriteManualCommanderAttack1Right,
		spriteManualCommanderAttack1Right,
		spriteManualCommanderAttack1Right,
		spriteManualCommanderAttack1Right,
		spriteManualCommanderAttack1Right,
		spriteManualCommanderAttack1Right,
		spriteManualCommanderAttack1Right,
		spriteManualCommanderAttack1Right,
		spriteManualCommanderAttack1Right,
		spriteManualCommanderAttack1Right,
		spriteManualCommanderAttack1Right,
	}
	spriteTable[unitCommander][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderIdleUpLeft,
		spriteManualCommanderAttack2UpLeft,
		spriteManualCommanderAttack2UpLeft,
		spriteManualCommanderAttack2UpLeft,
		spriteManualCommanderAttack2UpLeft,
		spriteManualCommanderAttack1UpLeft,
		spriteManualCommanderAttack1UpLeft,
		spriteManualCommanderAttack1UpLeft,
		spriteManualCommanderAttack1UpLeft,
		spriteManualCommanderAttack1UpLeft,
		spriteManualCommanderAttack1UpLeft,
		spriteManualCommanderAttack1UpLeft,
		spriteManualCommanderAttack1UpLeft,
		spriteManualCommanderAttack1UpLeft,
		spriteManualCommanderAttack1UpLeft,
		spriteManualCommanderAttack1UpLeft,
	}
	spriteTable[unitCommander][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderIdleUpRight,
		spriteManualCommanderAttack2UpRight,
		spriteManualCommanderAttack2UpRight,
		spriteManualCommanderAttack2UpRight,
		spriteManualCommanderAttack2UpRight,
		spriteManualCommanderAttack1UpRight,
		spriteManualCommanderAttack1UpRight,
		spriteManualCommanderAttack1UpRight,
		spriteManualCommanderAttack1UpRight,
		spriteManualCommanderAttack1UpRight,
		spriteManualCommanderAttack1UpRight,
		spriteManualCommanderAttack1UpRight,
		spriteManualCommanderAttack1UpRight,
		spriteManualCommanderAttack1UpRight,
		spriteManualCommanderAttack1UpRight,
		spriteManualCommanderAttack1UpRight,
	}
	spriteTable[unitCommander][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderIdleDownLeft,
		spriteManualCommanderAttack2DownLeft,
		spriteManualCommanderAttack2DownLeft,
		spriteManualCommanderAttack2DownLeft,
		spriteManualCommanderAttack2DownLeft,
		spriteManualCommanderAttack1DownLeft,
		spriteManualCommanderAttack1DownLeft,
		spriteManualCommanderAttack1DownLeft,
		spriteManualCommanderAttack1DownLeft,
		spriteManualCommanderAttack1DownLeft,
		spriteManualCommanderAttack1DownLeft,
		spriteManualCommanderAttack1DownLeft,
		spriteManualCommanderAttack1DownLeft,
		spriteManualCommanderAttack1DownLeft,
		spriteManualCommanderAttack1DownLeft,
		spriteManualCommanderAttack1DownLeft,
	}
	spriteTable[unitCommander][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderIdleDownRight,
		spriteManualCommanderAttack2DownRight,
		spriteManualCommanderAttack2DownRight,
		spriteManualCommanderAttack2DownRight,
		spriteManualCommanderAttack2DownRight,
		spriteManualCommanderAttack1DownRight,
		spriteManualCommanderAttack1DownRight,
		spriteManualCommanderAttack1DownRight,
		spriteManualCommanderAttack1DownRight,
		spriteManualCommanderAttack1DownRight,
		spriteManualCommanderAttack1DownRight,
		spriteManualCommanderAttack1DownRight,
		spriteManualCommanderAttack1DownRight,
		spriteManualCommanderAttack1DownRight,
		spriteManualCommanderAttack1DownRight,
		spriteManualCommanderAttack1DownRight,
	}
	// Koniec dowódcy

	spriteTable[unitBear][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
	}
	spriteTable[unitBear][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
	}
	spriteTable[unitBear][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
	}
	spriteTable[unitBear][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
	}
	spriteTable[unitBear][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
	}
	spriteTable[unitBear][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
	}
	spriteTable[unitBear][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
	}
	spriteTable[unitBear][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
	}
	spriteTable[unitBear][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
	}

	spriteTable[unitBear][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearMove1Down,
		spriteManualBearMove1Down,
		spriteManualBearMove1Down,
		spriteManualBearMove1Down,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearMove2Down,
		spriteManualBearMove2Down,
		spriteManualBearMove2Down,
		spriteManualBearMove2Down,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
	}
	spriteTable[unitBear][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearMove1Down,
		spriteManualBearMove1Down,
		spriteManualBearMove1Down,
		spriteManualBearMove1Down,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearMove2Down,
		spriteManualBearMove2Down,
		spriteManualBearMove2Down,
		spriteManualBearMove2Down,
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
	}
	spriteTable[unitBear][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearMove1Left,
		spriteManualBearMove1Left,
		spriteManualBearMove1Left,
		spriteManualBearMove1Left,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearMove2Left,
		spriteManualBearMove2Left,
		spriteManualBearMove2Left,
		spriteManualBearMove2Left,
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
	}
	spriteTable[unitBear][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearMove1Up,
		spriteManualBearMove1Up,
		spriteManualBearMove1Up,
		spriteManualBearMove1Up,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearMove2Up,
		spriteManualBearMove2Up,
		spriteManualBearMove2Up,
		spriteManualBearMove2Up,
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
	}
	spriteTable[unitBear][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearMove1Right,
		spriteManualBearMove1Right,
		spriteManualBearMove1Right,
		spriteManualBearMove1Right,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearMove2Right,
		spriteManualBearMove2Right,
		spriteManualBearMove2Right,
		spriteManualBearMove2Right,
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
	}
	spriteTable[unitBear][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearMove1UpLeft,
		spriteManualBearMove1UpLeft,
		spriteManualBearMove1UpLeft,
		spriteManualBearMove1UpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearMove2UpLeft,
		spriteManualBearMove2UpLeft,
		spriteManualBearMove2UpLeft,
		spriteManualBearMove2UpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
	}
	spriteTable[unitBear][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearMove1UpRight,
		spriteManualBearMove1UpRight,
		spriteManualBearMove1UpRight,
		spriteManualBearMove1UpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearMove2UpRight,
		spriteManualBearMove2UpRight,
		spriteManualBearMove2UpRight,
		spriteManualBearMove2UpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
	}
	spriteTable[unitBear][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearMove1DownLeft,
		spriteManualBearMove1DownLeft,
		spriteManualBearMove1DownLeft,
		spriteManualBearMove1DownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearMove2DownLeft,
		spriteManualBearMove2DownLeft,
		spriteManualBearMove2DownLeft,
		spriteManualBearMove2DownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
	}
	spriteTable[unitBear][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearMove1DownRight,
		spriteManualBearMove1DownRight,
		spriteManualBearMove1DownRight,
		spriteManualBearMove1DownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearMove2DownRight,
		spriteManualBearMove2DownRight,
		spriteManualBearMove2DownRight,
		spriteManualBearMove2DownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
	}

	spriteTable[unitBear][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
	}
	spriteTable[unitBear][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualBearIdleDown,
		spriteManualBearIdleDown,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack2Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
		spriteManualBearAttack1Down,
	}
	spriteTable[unitBear][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualBearIdleLeft,
		spriteManualBearIdleLeft,
		spriteManualBearAttack2Left,
		spriteManualBearAttack2Left,
		spriteManualBearAttack2Left,
		spriteManualBearAttack2Left,
		spriteManualBearAttack2Left,
		spriteManualBearAttack2Left,
		spriteManualBearAttack2Left,
		spriteManualBearAttack1Left,
		spriteManualBearAttack1Left,
		spriteManualBearAttack1Left,
		spriteManualBearAttack1Left,
		spriteManualBearAttack1Left,
		spriteManualBearAttack1Left,
		spriteManualBearAttack1Left,
		spriteManualBearAttack1Left,
	}
	spriteTable[unitBear][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualBearIdleUp,
		spriteManualBearIdleUp,
		spriteManualBearAttack2Up,
		spriteManualBearAttack2Up,
		spriteManualBearAttack2Up,
		spriteManualBearAttack2Up,
		spriteManualBearAttack2Up,
		spriteManualBearAttack2Up,
		spriteManualBearAttack2Up,
		spriteManualBearAttack1Up,
		spriteManualBearAttack1Up,
		spriteManualBearAttack1Up,
		spriteManualBearAttack1Up,
		spriteManualBearAttack1Up,
		spriteManualBearAttack1Up,
		spriteManualBearAttack1Up,
		spriteManualBearAttack1Up,
	}
	spriteTable[unitBear][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualBearIdleRight,
		spriteManualBearIdleRight,
		spriteManualBearAttack2Right,
		spriteManualBearAttack2Right,
		spriteManualBearAttack2Right,
		spriteManualBearAttack2Right,
		spriteManualBearAttack2Right,
		spriteManualBearAttack2Right,
		spriteManualBearAttack2Right,
		spriteManualBearAttack1Right,
		spriteManualBearAttack1Right,
		spriteManualBearAttack1Right,
		spriteManualBearAttack1Right,
		spriteManualBearAttack1Right,
		spriteManualBearAttack1Right,
		spriteManualBearAttack1Right,
		spriteManualBearAttack1Right,
	}
	spriteTable[unitBear][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualBearIdleUpLeft,
		spriteManualBearIdleUpLeft,
		spriteManualBearAttack2UpLeft,
		spriteManualBearAttack2UpLeft,
		spriteManualBearAttack2UpLeft,
		spriteManualBearAttack2UpLeft,
		spriteManualBearAttack2UpLeft,
		spriteManualBearAttack2UpLeft,
		spriteManualBearAttack2UpLeft,
		spriteManualBearAttack1UpLeft,
		spriteManualBearAttack1UpLeft,
		spriteManualBearAttack1UpLeft,
		spriteManualBearAttack1UpLeft,
		spriteManualBearAttack1UpLeft,
		spriteManualBearAttack1UpLeft,
		spriteManualBearAttack1UpLeft,
		spriteManualBearAttack1UpLeft,
	}
	spriteTable[unitBear][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualBearIdleUpRight,
		spriteManualBearIdleUpRight,
		spriteManualBearAttack2UpRight,
		spriteManualBearAttack2UpRight,
		spriteManualBearAttack2UpRight,
		spriteManualBearAttack2UpRight,
		spriteManualBearAttack2UpRight,
		spriteManualBearAttack2UpRight,
		spriteManualBearAttack2UpRight,
		spriteManualBearAttack1UpRight,
		spriteManualBearAttack1UpRight,
		spriteManualBearAttack1UpRight,
		spriteManualBearAttack1UpRight,
		spriteManualBearAttack1UpRight,
		spriteManualBearAttack1UpRight,
		spriteManualBearAttack1UpRight,
		spriteManualBearAttack1UpRight,
	}
	spriteTable[unitBear][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualBearIdleDownLeft,
		spriteManualBearIdleDownLeft,
		spriteManualBearAttack2DownLeft,
		spriteManualBearAttack2DownLeft,
		spriteManualBearAttack2DownLeft,
		spriteManualBearAttack2DownLeft,
		spriteManualBearAttack2DownLeft,
		spriteManualBearAttack2DownLeft,
		spriteManualBearAttack2DownLeft,
		spriteManualBearAttack1DownLeft,
		spriteManualBearAttack1DownLeft,
		spriteManualBearAttack1DownLeft,
		spriteManualBearAttack1DownLeft,
		spriteManualBearAttack1DownLeft,
		spriteManualBearAttack1DownLeft,
		spriteManualBearAttack1DownLeft,
		spriteManualBearAttack1DownLeft,
	}
	spriteTable[unitBear][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualBearIdleDownRight,
		spriteManualBearIdleDownRight,
		spriteManualBearAttack2DownRight,
		spriteManualBearAttack2DownRight,
		spriteManualBearAttack2DownRight,
		spriteManualBearAttack2DownRight,
		spriteManualBearAttack2DownRight,
		spriteManualBearAttack2DownRight,
		spriteManualBearAttack2DownRight,
		spriteManualBearAttack1DownRight,
		spriteManualBearAttack1DownRight,
		spriteManualBearAttack1DownRight,
		spriteManualBearAttack1DownRight,
		spriteManualBearAttack1DownRight,
		spriteManualBearAttack1DownRight,
		spriteManualBearAttack1DownRight,
		spriteManualBearAttack1DownRight,
	}
	// Koniec niedźwiedzia

	spriteTable[unitUnknown][animationIdle][directionNone] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
	}
	spriteTable[unitUnknown][animationIdle][directionDown] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
	}
	spriteTable[unitUnknown][animationIdle][directionLeft] = [delayStatesCount]uint16{
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
	}
	spriteTable[unitUnknown][animationIdle][directionUp] = [delayStatesCount]uint16{
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
	}
	spriteTable[unitUnknown][animationIdle][directionRight] = [delayStatesCount]uint16{
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
	}
	spriteTable[unitUnknown][animationIdle][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
	}
	spriteTable[unitUnknown][animationIdle][directionUpRight] = [delayStatesCount]uint16{
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
	}
	spriteTable[unitUnknown][animationIdle][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
	}
	spriteTable[unitUnknown][animationIdle][directionDownRight] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
	}

	spriteTable[unitUnknown][animationWalk][directionNone] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDown,
		spriteManualUnknownMove1Down,
		spriteManualUnknownMove1Down,
		spriteManualUnknownIdleDown,
		spriteManualUnknownMove2Down,
		spriteManualUnknownMove2Down,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownMove2Down,
		spriteManualUnknownMove2Down,
		spriteManualUnknownMove2Down,
		spriteManualUnknownMove2Down,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
	}
	spriteTable[unitUnknown][animationWalk][directionDown] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDown,
		spriteManualUnknownMove1Down,
		spriteManualUnknownMove1Down,
		spriteManualUnknownIdleDown,
		spriteManualUnknownMove2Down,
		spriteManualUnknownMove2Down,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownMove2Down,
		spriteManualUnknownMove2Down,
		spriteManualUnknownMove2Down,
		spriteManualUnknownMove2Down,
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
	}
	spriteTable[unitUnknown][animationWalk][directionLeft] = [delayStatesCount]uint16{
		spriteManualUnknownIdleLeft,
		spriteManualUnknownMove1Left,
		spriteManualUnknownMove1Left,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownMove2Left,
		spriteManualUnknownMove2Left,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownMove2Left,
		spriteManualUnknownMove2Left,
		spriteManualUnknownMove2Left,
		spriteManualUnknownMove2Left,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
	}
	spriteTable[unitUnknown][animationWalk][directionUp] = [delayStatesCount]uint16{
		spriteManualUnknownIdleUp,
		spriteManualUnknownMove1Up,
		spriteManualUnknownMove1Up,
		spriteManualUnknownIdleUp,
		spriteManualUnknownMove2Up,
		spriteManualUnknownMove2Up,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownMove2Up,
		spriteManualUnknownMove2Up,
		spriteManualUnknownMove2Up,
		spriteManualUnknownMove2Up,
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
	}
	spriteTable[unitUnknown][animationWalk][directionRight] = [delayStatesCount]uint16{
		spriteManualUnknownIdleRight,
		spriteManualUnknownMove1Right,
		spriteManualUnknownMove1Right,
		spriteManualUnknownIdleRight,
		spriteManualUnknownMove2Right,
		spriteManualUnknownMove2Right,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownMove2Right,
		spriteManualUnknownMove2Right,
		spriteManualUnknownMove2Right,
		spriteManualUnknownMove2Right,
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
	}
	spriteTable[unitUnknown][animationWalk][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownMove1UpLeft,
		spriteManualUnknownMove1UpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownMove2UpLeft,
		spriteManualUnknownMove2UpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownMove2UpLeft,
		spriteManualUnknownMove2UpLeft,
		spriteManualUnknownMove2UpLeft,
		spriteManualUnknownMove2UpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
	}
	spriteTable[unitUnknown][animationWalk][directionUpRight] = [delayStatesCount]uint16{
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownMove1UpRight,
		spriteManualUnknownMove1UpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownMove2UpRight,
		spriteManualUnknownMove2UpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownMove2UpRight,
		spriteManualUnknownMove2UpRight,
		spriteManualUnknownMove2UpRight,
		spriteManualUnknownMove2UpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
	}
	spriteTable[unitUnknown][animationWalk][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownMove1DownLeft,
		spriteManualUnknownMove1DownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownMove2DownLeft,
		spriteManualUnknownMove2DownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownMove2DownLeft,
		spriteManualUnknownMove2DownLeft,
		spriteManualUnknownMove2DownLeft,
		spriteManualUnknownMove2DownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
	}
	spriteTable[unitUnknown][animationWalk][directionDownRight] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownMove1DownRight,
		spriteManualUnknownMove1DownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownMove2DownRight,
		spriteManualUnknownMove2DownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownMove2DownRight,
		spriteManualUnknownMove2DownRight,
		spriteManualUnknownMove2DownRight,
		spriteManualUnknownMove2DownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
	}

	spriteTable[unitUnknown][animationFight][directionNone] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownAttack2Down,
		spriteManualUnknownAttack2Down,
		spriteManualUnknownAttack2Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
	}
	spriteTable[unitUnknown][animationFight][directionDown] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDown,
		spriteManualUnknownIdleDown,
		spriteManualUnknownAttack2Down,
		spriteManualUnknownAttack2Down,
		spriteManualUnknownAttack2Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
		spriteManualUnknownAttack1Down,
	}
	spriteTable[unitUnknown][animationFight][directionLeft] = [delayStatesCount]uint16{
		spriteManualUnknownIdleLeft,
		spriteManualUnknownIdleLeft,
		spriteManualUnknownAttack2Left,
		spriteManualUnknownAttack2Left,
		spriteManualUnknownAttack2Left,
		spriteManualUnknownAttack1Left,
		spriteManualUnknownAttack1Left,
		spriteManualUnknownAttack1Left,
		spriteManualUnknownAttack1Left,
		spriteManualUnknownAttack1Left,
		spriteManualUnknownAttack1Left,
		spriteManualUnknownAttack1Left,
		spriteManualUnknownAttack1Left,
		spriteManualUnknownAttack1Left,
		spriteManualUnknownAttack1Left,
		spriteManualUnknownAttack1Left,
		spriteManualUnknownAttack1Left,
	}
	spriteTable[unitUnknown][animationFight][directionUp] = [delayStatesCount]uint16{
		spriteManualUnknownIdleUp,
		spriteManualUnknownIdleUp,
		spriteManualUnknownAttack2Up,
		spriteManualUnknownAttack2Up,
		spriteManualUnknownAttack2Up,
		spriteManualUnknownAttack1Up,
		spriteManualUnknownAttack1Up,
		spriteManualUnknownAttack1Up,
		spriteManualUnknownAttack1Up,
		spriteManualUnknownAttack1Up,
		spriteManualUnknownAttack1Up,
		spriteManualUnknownAttack1Up,
		spriteManualUnknownAttack1Up,
		spriteManualUnknownAttack1Up,
		spriteManualUnknownAttack1Up,
		spriteManualUnknownAttack1Up,
		spriteManualUnknownAttack1Up,
	}
	spriteTable[unitUnknown][animationFight][directionRight] = [delayStatesCount]uint16{
		spriteManualUnknownIdleRight,
		spriteManualUnknownIdleRight,
		spriteManualUnknownAttack2Right,
		spriteManualUnknownAttack2Right,
		spriteManualUnknownAttack2Right,
		spriteManualUnknownAttack1Right,
		spriteManualUnknownAttack1Right,
		spriteManualUnknownAttack1Right,
		spriteManualUnknownAttack1Right,
		spriteManualUnknownAttack1Right,
		spriteManualUnknownAttack1Right,
		spriteManualUnknownAttack1Right,
		spriteManualUnknownAttack1Right,
		spriteManualUnknownAttack1Right,
		spriteManualUnknownAttack1Right,
		spriteManualUnknownAttack1Right,
		spriteManualUnknownAttack1Right,
	}
	spriteTable[unitUnknown][animationFight][directionUpLeft] = [delayStatesCount]uint16{
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownIdleUpLeft,
		spriteManualUnknownAttack2UpLeft,
		spriteManualUnknownAttack2UpLeft,
		spriteManualUnknownAttack2UpLeft,
		spriteManualUnknownAttack1UpLeft,
		spriteManualUnknownAttack1UpLeft,
		spriteManualUnknownAttack1UpLeft,
		spriteManualUnknownAttack1UpLeft,
		spriteManualUnknownAttack1UpLeft,
		spriteManualUnknownAttack1UpLeft,
		spriteManualUnknownAttack1UpLeft,
		spriteManualUnknownAttack1UpLeft,
		spriteManualUnknownAttack1UpLeft,
		spriteManualUnknownAttack1UpLeft,
		spriteManualUnknownAttack1UpLeft,
		spriteManualUnknownAttack1UpLeft,
	}
	spriteTable[unitUnknown][animationFight][directionUpRight] = [delayStatesCount]uint16{
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownIdleUpRight,
		spriteManualUnknownAttack2UpRight,
		spriteManualUnknownAttack2UpRight,
		spriteManualUnknownAttack2UpRight,
		spriteManualUnknownAttack1UpRight,
		spriteManualUnknownAttack1UpRight,
		spriteManualUnknownAttack1UpRight,
		spriteManualUnknownAttack1UpRight,
		spriteManualUnknownAttack1UpRight,
		spriteManualUnknownAttack1UpRight,
		spriteManualUnknownAttack1UpRight,
		spriteManualUnknownAttack1UpRight,
		spriteManualUnknownAttack1UpRight,
		spriteManualUnknownAttack1UpRight,
		spriteManualUnknownAttack1UpRight,
		spriteManualUnknownAttack1UpRight,
	}
	spriteTable[unitUnknown][animationFight][directionDownLeft] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownIdleDownLeft,
		spriteManualUnknownAttack2DownLeft,
		spriteManualUnknownAttack2DownLeft,
		spriteManualUnknownAttack2DownLeft,
		spriteManualUnknownAttack1DownLeft,
		spriteManualUnknownAttack1DownLeft,
		spriteManualUnknownAttack1DownLeft,
		spriteManualUnknownAttack1DownLeft,
		spriteManualUnknownAttack1DownLeft,
		spriteManualUnknownAttack1DownLeft,
		spriteManualUnknownAttack1DownLeft,
		spriteManualUnknownAttack1DownLeft,
		spriteManualUnknownAttack1DownLeft,
		spriteManualUnknownAttack1DownLeft,
		spriteManualUnknownAttack1DownLeft,
		spriteManualUnknownAttack1DownLeft,
	}
	spriteTable[unitUnknown][animationFight][directionDownRight] = [delayStatesCount]uint16{
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownIdleDownRight,
		spriteManualUnknownAttack2DownRight,
		spriteManualUnknownAttack2DownRight,
		spriteManualUnknownAttack2DownRight,
		spriteManualUnknownAttack1DownRight,
		spriteManualUnknownAttack1DownRight,
		spriteManualUnknownAttack1DownRight,
		spriteManualUnknownAttack1DownRight,
		spriteManualUnknownAttack1DownRight,
		spriteManualUnknownAttack1DownRight,
		spriteManualUnknownAttack1DownRight,
		spriteManualUnknownAttack1DownRight,
		spriteManualUnknownAttack1DownRight,
		spriteManualUnknownAttack1DownRight,
		spriteManualUnknownAttack1DownRight,
		spriteManualUnknownAttack1DownRight,
	}
	// Koniec strzygi
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
		Fresh: spriteManualPriestessCorpseFresh,
		Decay: spriteManualPriestessCorpseDecay,
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
