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
	initUnitSprites()
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

type unitFrame uint8

// Jednostki
// ID: 700 + (unitType * 200) + (Frame * 8) + Direction.
// Frame 0: Idle, Frame 1: Walk1, Frame 2: Walk2, Frame 3: Attack1, Frame 4: Attack2.
func initUnitSprites() {
	type unitCfg struct {
		Atlas battleAtlasID
		BaseX uint16
		BaseY uint16
		Melee bool
	}

	//goland:noinspection GoLinter
	configs := map[unitType]unitCfg{
		unitCow:         {Atlas: atlasUnits1, BaseX: 0, BaseY: 0, Melee: false},
		unitAxeman:      {Atlas: atlasUnits1, BaseX: 0, BaseY: 42, Melee: true},   //nolint:mnd
		unitArcher:      {Atlas: atlasUnits1, BaseX: 0, BaseY: 84, Melee: false},  //nolint:mnd
		unitPriestess:   {Atlas: atlasUnits1, BaseX: 0, BaseY: 126, Melee: false}, //nolint:mnd
		unitPriest:      {Atlas: atlasUnits2, BaseX: 0, BaseY: 0, Melee: false},
		unitSwordsman:   {Atlas: atlasUnits2, BaseX: 0, BaseY: 42, Melee: true},  //nolint:mnd
		unitSpearman:    {Atlas: atlasUnits2, BaseX: 0, BaseY: 84, Melee: false}, //nolint:mnd
		unitCommander:   {Atlas: atlasUnits2, BaseX: 0, BaseY: 126, Melee: true}, //nolint:mnd
		unitBear:        {Atlas: atlasBuildings, BaseX: 0, BaseY: 0, Melee: true},
		unitUnknown:     {Atlas: atlasBuildings, BaseX: 0, BaseY: 42, Melee: true},   //nolint:mnd
		unitShepherd:    {Atlas: atlasUnits1, BaseX: 160, BaseY: 0, Melee: false},    //nolint:mnd
		unitMage:        {Atlas: atlasUnits1, BaseX: 160, BaseY: 126, Melee: false},  //nolint:mnd
		unitCrossbowman: {Atlas: atlasBuildings, BaseX: 224, BaseY: 0, Melee: false}, //nolint:mnd
	}

	for unitIndex := range uint16(unitTypeCount) {
		currentUnitType := unitType(unitIndex)

		cfg, exists := configs[currentUnitType]
		if !exists {
			continue
		}

		spriteBaseID := spriteUnitBaseID + (unitIndex * spriteUnitStep)

		for frame := range frameCount {
			if cfg.Melee && (frame == frameAttack1 || frame == frameAttack2) {
				generateMeleeAttackFrames(spriteBaseID, frame, cfg.Atlas, cfg.BaseY)

				continue
			}

			colX := cfg.BaseX + uint16(frame*32) //nolint:mnd

			if currentUnitType == unitCrossbowman && frame > 0 {
				// @reminder kusznik jest jednym z najgorszych przypadków, który wymaga osobnego
				// podejścia. Tutaj jest tylko częściowo ogarnięty.
				switch frame {
				case frameWalk1:
					colX = 240
				case frameWalk2:
					colX = 256
				case frameAttack1:
					colX = 272
				}
			}

			for direction := range directionCount {
				var sourceX, sourceY uint16

				var flip bool

				var cropWidth, cropHeight uint8 = 16, 14

				atlasY0 := cfg.BaseY
				atlasY1 := cfg.BaseY + 14 //nolint:mnd
				atlasY2 := cfg.BaseY + 28 //nolint:mnd

				// @reminder: tutaj obsługujemy wszystkie możliwości.
				// „brakująca” to directionCount.
				switch direction {
				case directionUp:
					sourceY = atlasY0
					sourceX = colX + 16 //nolint:mnd
					flip = false
				case directionUpRight:
					sourceY = atlasY0
					sourceX = colX
					flip = true
				case directionRight:
					sourceY = atlasY1
					sourceX = colX
					flip = true
				case directionDownRight:
					sourceY = atlasY2
					sourceX = colX
					flip = true
				case directionDown:
					sourceY = atlasY2
					sourceX = colX + 16 //nolint:mnd
					flip = false
				case directionDownLeft:
					sourceY = atlasY2
					sourceX = colX
					flip = false
				case directionLeft:
					sourceY = atlasY1
					sourceX = colX
					flip = false
				case directionUpLeft:
					sourceY = atlasY0
					sourceX = colX + 16 //nolint:mnd
					flip = false
				}

				finalID := spriteBaseID + (uint16(frame) * 8) + uint16(direction) //nolint:mnd

				if finalID < maxSpriteID {
					spriteRegistry[finalID] = spriteDef{
						atlasID:    cfg.Atlas,
						cropX:      sourceX,
						cropY:      sourceY,
						cropWidth:  cropWidth,
						cropHeight: cropHeight,
						flipX:      flip,
						offX:       0,
						offY:       0,
					}
				}
			}
		}

		var freshDeadX, freshDeadY, decayStartX, decayStartY uint16

		// unitNone nie istnieje, jest chwilowym zapychaczem. Dlatego go tutaj nie ma.
		switch currentUnitType { //nolint:exhaustive
		case unitCow: // Cow
			freshDeadX, freshDeadY, decayStartX, decayStartY = 112, 14, 144, 14
		case unitAxeman: // Axeman
			freshDeadX, freshDeadY, decayStartX, decayStartY = 120, 70, 200, 70
		case unitArcher: // Archer
			freshDeadX, freshDeadY, decayStartX, decayStartY = 112, 98, 144, 98
		case unitPriestess: // Priestess
			freshDeadX, freshDeadY, decayStartX, decayStartY = 111, 140, 143, 140
		case unitPriest: // Priest
			freshDeadX, freshDeadY, decayStartX, decayStartY = 112, 14, 144, 14
		case unitSwordsman: // Swordsman
			freshDeadX, freshDeadY, decayStartX, decayStartY = 120, 70, 200, 70
		case unitSpearman: // Spearman
			freshDeadX, freshDeadY, decayStartX, decayStartY = 112, 98, 144, 98
		case unitCommander: // Commander
			freshDeadX, freshDeadY, decayStartX, decayStartY = 120, 154, 200, 154
		case unitBear: // Bear
			freshDeadX, freshDeadY, decayStartX, decayStartY = 120, 28, 200, 28
		case unitUnknown: // Strzyga
			freshDeadX, freshDeadY, decayStartX, decayStartY = 120, 70, 200, 70
		case unitShepherd: // Shepherd
			freshDeadX, freshDeadY, decayStartX, decayStartY = 160, 98, 192, 98
		case unitMage: // Mage
			freshDeadX, freshDeadY, decayStartX, decayStartY = 224, 98, 256, 98
		case unitCrossbowman: // Crossbow
			freshDeadX, freshDeadY, decayStartX, decayStartY = 239, 56, 271, 56
		default:
			continue
		}

		freshDeadSpriteID := spriteBaseID + freshlyDeadSpriteOffset
		if freshDeadSpriteID < maxSpriteID {
			spriteRegistry[freshDeadSpriteID] = spriteDef{
				atlasID:    cfg.Atlas,
				cropX:      freshDeadX,
				cropY:      freshDeadY,
				cropWidth:  tileWidth,
				cropHeight: tileHeight,
				flipX:      false,
				offX:       0,
				offY:       0,
			}
		}

		decayStartSpriteID := spriteBaseID + decayStartSpriteOffset
		if decayStartSpriteID < maxSpriteID {
			spriteRegistry[decayStartSpriteID] = spriteDef{
				atlasID:    cfg.Atlas,
				cropX:      decayStartX,
				cropY:      decayStartY,
				cropWidth:  tileWidth,
				cropHeight: tileHeight,
				flipX:      false,
				offX:       0,
				offY:       0,
			}
		}
	}
}

func generateMeleeAttackFrames(unitBaseID uint16, frame unitFrame, atlas battleAtlasID, unitBaseY uint16) {
	// Współrzędne x w atlasie dla poszczególnych faz (hardcoded z bridge)
	// Bridge Col 0 (Diag), Col 1 (Up), Col 2 (Side), Col 3 (DownDiag), Col 4 (Down)

	for direction := range directionCount { //nolint:nolintlint,mnd
		var cropX, cropY uint16

		var cropWidth, cropHeight uint8

		var offX, offY int8

		var flip bool

		// Domyślne wartości (żeby nie było zer)
		cropWidth = tileWidth
		cropHeight = tileHeight

		// === FAZA 1 ATAKU ===
		if frame == frameAttack1 {
			// @reminder: tutaj obsługujemy wszystkie możliwości.
			// „brakująca” to directionCount.
			switch direction {
			case directionUp: // Góra (Wąska 16x21)
				cropX = 120
				cropY = unitBaseY + 7 //nolint:mnd
				cropWidth = 16
				cropHeight = 21
				flip = false
				// Korekta wysokości
				offY = -7 //nolint:mnd
			case directionUpRight: // Góra-Prawo (Szeroka 24x21)
				cropX = 96
				cropY = unitBaseY + 7 //nolint:mnd
				cropWidth = 24
				cropHeight = 21
				flip = true
				offX = -4 //nolint:mnd
				offY = -7 //nolint:mnd
			case directionRight: // Prawo (Wąska 24x14? Bridge mówi 24x14)
				cropX = 96
				cropY = unitBaseY + 28 //nolint:mnd
				cropWidth = 24
				cropHeight = 14
				flip = true
				offX = -4 //nolint:mnd //nolint:mnd
			case directionDownRight: // Dół-Prawo (Szeroka 24x21) - Bridge: Col 3 (AbsX 152, RelY 0 ??)
				// Bridge Row 3 Col 3: AbsX 152, RelY 0.
				cropX = 152
				cropY = unitBaseY
				cropWidth = 24
				cropHeight = 21
				flip = true
				offX = -4 //nolint:mnd
				offY = -7 //nolint:mnd
			case directionDown: // Dół (Wąska 16x21) - Bridge: Col 4 (AbsX 136, RelY 0)
				cropX = 136
				cropY = unitBaseY
				cropWidth = 16
				cropHeight = 21
				flip = false
				offY = -7 //nolint:mnd
			case directionDownLeft: // Dół-Lewo (Flip DP)
				cropX = 152
				cropY = unitBaseY
				cropWidth = 24
				cropHeight = 21
				flip = false
				offX = -4 //nolint:mnd
				offY = -7 //nolint:mnd
			case directionLeft: // Lewo (Flip P)
				cropX = 96
				cropY = unitBaseY + 28 //nolint:mnd
				cropWidth = 24
				cropHeight = 14
				flip = false
				offX = -4 //nolint:mnd
			case directionUpLeft: // Góra-Lewo (Flip GP)
				cropX = 96
				cropY = unitBaseY + 7 //nolint:mnd
				cropWidth = 24
				cropHeight = 21
				flip = false
				offX = -4 //nolint:mnd
				offY = -7 //nolint:mnd
			}
		} else { // === FAZA 2 ATAKU (Frame 4) ===
			// frameAttack2
			switch direction {
			// @reminder: tutaj obsługujemy wszystkie możliwości.
			// „brakująca” to directionCount.
			case directionUp: // Góra
				cropX = 200
				cropY = unitBaseY + 7 //nolint:mnd
				cropWidth = 16
				cropHeight = 21
				flip = false
				offY = -7 //nolint:mnd
			case directionUpRight: // GP
				cropX = 176
				cropY = unitBaseY + 7 //nolint:mnd
				cropWidth = 24
				cropHeight = 21
				flip = true
				offX = -4 //nolint:mnd
				offY = -7 //nolint:mnd
			case directionRight: // P
				cropX = 200
				cropY = unitBaseY + 28 //nolint:mnd
				cropWidth = 16
				cropHeight = 14
				flip = true
			case directionDownRight: // DP
				cropX = 152
				cropY = unitBaseY + 21 //nolint:mnd
				cropWidth = 24
				cropHeight = 21
				flip = true
				offX = -4 //nolint:mnd
				offY = -7 //nolint:mnd
			case directionDown: // D
				cropX = 136
				cropY = unitBaseY + 21 //nolint:mnd
				cropWidth = 16
				cropHeight = 21
				flip = false
				offY = -7 //nolint:mnd
			case directionDownLeft: // DL
				cropX = 152
				cropY = unitBaseY + 21 //nolint:mnd
				cropWidth = 24
				cropHeight = 21
				flip = false
				offX = -4 //nolint:mnd
				offY = -7 //nolint:mnd
			case directionLeft: // L
				cropX = 200
				cropY = unitBaseY + 28 //nolint:mnd
				cropWidth = 16
				cropHeight = 14
				flip = false
			case directionUpLeft: // GL
				cropX = 176
				cropY = unitBaseY + 7 //nolint:mnd
				cropWidth = 24
				cropHeight = 21
				flip = false
				offX = -4 //nolint:mnd
				offY = -7 //nolint:mnd
			}
		}

		finalID := unitBaseID + (uint16(frame) * 8) + uint16(direction) //nolint:mnd

		if finalID < maxSpriteID {
			spriteRegistry[finalID] = spriteDef{
				atlasID: atlas,
				cropX:   cropX, cropY: cropY, cropWidth: cropWidth, cropHeight: cropHeight,
				flipX: flip,
				offX:  offX, offY: offY,
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
