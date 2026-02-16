package main

import "log"

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
	"SPRITE_TREE_STUMP_06":    spriteTreeStump06,

	// Specjalne
	"SPRITE_PALISADE_START":      spritePalisadeStart,
	"SPRITE_EFFECT_HEAL_00":      spriteEffectHeal00,
	"SPRITE_EFFECT_TRANSFORM_00": spriteEffectTransform00,
	"SPRITE_VICTORY_POINT":       spriteVictoryPoint,
	"SPRITE_EFFECT_SKELETON_00":  spriteEffectSkeleton_00,
	"SPRITE_EFFECT_SKELETON_01":  spriteEffectSkeleton_01,
	"SPRITE_EFFECT_SKELETON_02":  spriteEffectSkeleton_02,

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
	"SPRITE_FIRE_12": spriteFire12,
	"SPRITE_FIRE_13": spriteFire13,
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
func getID(name string) uint16 {
	if val, ok := idRegistry[name]; ok {
		return val
	}
	// Bezpiecznik dla pustych łańcuchów, które mogły zostać z przenosin jako "" dla starych ID lub błędnych nazw.
	if name != "" {
		log.Printf("OSTRZEŻENIE: Brak ID dla '%s'. Używam trawy.", name)
	}
	return spriteGrassStart
}

func init() {
	initTerrainSprites()
	initUISprites()
	initUnitSprites()
	initBuildingSprites()
}

func initTerrainSprites() {
	// Pomagier do nakładki (16x14)
	setUI := func(id, x, y uint16) {
		if id < maxSpriteID {
			spriteRegistry[id] = spriteDef{
				atlasID: atlasUI,
				x:       x, y: y,
				w: 16, h: 14,
				offX: 0, offY: 0,
			}
		}
	}

	// Pomagier do Units1 (16x14)
	setUnit1 := func(id, x, y uint16) {
		if id < maxSpriteID {
			spriteRegistry[id] = spriteDef{
				atlasID: atlasUnits1,
				x:       x, y: y,
				w: 16, h: 14,
				offX: 0, offY: 0,
			}
		}
	}

	// Pomagier do gabarytów
	setSpecial := func(id, x, y, w, h uint16) {
		if id < maxSpriteID {
			spriteRegistry[id] = spriteDef{
				atlasID: atlasUI,
				x:       x, y: y,
				w: w, h: h,
				offX: 0, offY: 0,
			}
		}
	}

	// 1. Trawa
	setUI(spriteGrassStubbed, 11, 36)
	setUI(spriteGrassGrazed, 27, 36)

	// A. Tekstury z atlasu UI
	setUI(spriteGrass00, 43, 36)
	setUI(spriteGrass01, 139, 134)

	// B. Tekstury z atlasu Units1
	setUnit1(spriteEffectHeal00, 303, 0)
	setUnit1(spriteEffectHeal01, 303, 14)
	setUnit1(spriteGrass02, 303, 28)
	setUnit1(spriteGrass03, 303, 42)
	setUnit1(spriteGrass04, 303, 56)
	setUnit1(spriteGrass05, 303, 70)
	setUnit1(spriteGrass06, 303, 84)
	setUnit1(spriteGrass07, 303, 98)

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
	setUI(spriteRockEnd, 11+(8*16), 92)
	currentRockID := spriteRockStart
	for i := 0; i <= 12; i++ {
		if i == 8 {
			continue
		}
		setUI(currentRockID, uint16(11+(i*16)), 92)
		currentRockID++
	}

	// Sucha ziemia
	setDryEarth := func(id, x, y uint16) {
		if id < maxSpriteID {
			spriteRegistry[id] = spriteDef{
				atlasID: atlasUI,
				x:       x, y: y,
				w: 22, h: 18, // Szersza i wyższa
				offX: -3, offY: -2, // Przesunięcie w lewo i do góry, żeby wyśrodkować
			}
		}
	}

	setDryEarth(spriteDryEarth01, 235, 134)
	setDryEarth(spriteDryEarth02, 235, 153)
	setDryEarth(spriteDryEarth03, 235, 172)

	// Drogi
	for i := uint16(0); i < 5; i++ {
		setUI(spriteRoadStart+i, 59+(i*16), 134)
	}

	setUI(spriteRoadStart+5, 139, 134)

	for i := uint16(0); i < 4; i++ {
		setUI(spriteRoadStart+6+i, 155+(i*16), 134)
	}

	for i := uint16(0); i < 11; i++ {
		setUI(spriteRoadStart+10+i, 11+(i*16), 148)
	}

	// Mosty
	for i := uint16(0); i <= (spriteBridgeEnd - spriteBridgeStart); i++ {
		setUI(spriteBridgeStart+i, 11+(i*16), 162)
	}

	// Woda
	for i := uint16(0); i <= 12; i++ {
		setUI(spriteWaterStart+i, 11+(i*16), 50)
	}
	for i := uint16(0); i <= 12; i++ {
		setUI(spriteWaterStart+13+i, 11+(i*16), 64)
	}
	for i := uint16(0); i <= 12; i++ {
		setUI(spriteWaterStart+26+i, 11+(i*16), 78)
	}

	// Drzewa
	for i := uint16(0); i <= 6; i++ {
		setSpecial(spriteTreeStumpStart+i, 11+(i*32), 120, 32, 14)
	}
	for i := uint16(0); i <= 6; i++ {
		setSpecial(spriteTreeTopStart+i, 11+(i*32), 106, 32, 14)
	}

	// Palisady
	for i := uint16(0); i <= (spritePalisadeEnd - spritePalisadeStart); i++ {
		id := spritePalisadeStart + i
		spriteRegistry[id] = spriteDef{
			atlasID: atlasUnits1,
			x:       287,
			y:       i * 14,
			w:       16, h: 14,
			offX: 0, offY: 0,
		}
	}

	// Pierdoły

	for i := uint16(0); i <= 10; i++ {
		setUI(spriteGadgetStart+i, 59+(i*16), 36)
	}

	for i := uint16(0); i <= 2; i++ {
		setUI(spriteGadgetStart+11+i, 187+(i*16), 148)
	}

	for i := uint16(0); i <= 5; i++ {
		setUI(spriteGadgetStart+14+i, 139+(i*16), 162)
	}

	// Rany
	setUI(spriteEffectHit00, 203, 8)
	setUI(spriteEffectHit01, 219, 8)
}

// Nakładka.
func initUISprites() {
	setUI := func(id, x, y, w, h uint16) {
		if id < maxSpriteID {
			spriteRegistry[id] = spriteDef{
				atlasID: atlasUI,
				x:       x, y: y,
				w: w, h: h,
				offX: 0, offY: 0,
			}
		}
	}

	setCenterUI := func(id, x, y, w, h uint16) {
		if id < maxSpriteID {
			spriteRegistry[id] = spriteDef{
				atlasID: atlasUI,
				x:       x, y: y, w: w, h: h,
				offX: -8, offY: -7,
			}
		}
	}

	// Kursory
	setUI(spriteCursorDefaultBig, 11, 8, 16, 14)
	setCenterUI(spriteCursorCrossWhite, 27, 8, 16, 14)
	setCenterUI(spriteCursorCrossRed, 43, 8, 16, 14)
	setCenterUI(spriteCursorSmallWhite, 59, 8, 16, 14)
	setCenterUI(spriteCursorFrameRed, 75, 8, 16, 14)
	setCenterUI(spriteCursorCrossMedRed, 91, 8, 16, 14)
	setCenterUI(spriteCursorCrossMedWhite, 107, 8, 16, 14)
	setCenterUI(spriteCursorArrowUp, 123, 8, 16, 14)
	setCenterUI(spriteCursorArrowDown, 139, 8, 16, 14)
	setCenterUI(spriteCursorArrowLeft, 155, 8, 16, 14)
	setCenterUI(spriteCursorArrowRight, 171, 8, 16, 14)
	setCenterUI(spriteCursorStop, 187, 8, 16, 14)
	setUI(spriteCursorDefaultSmall, 91, 22, 16, 14)
	setUI(spriteCursorPointer, 107, 22, 16, 14)
	setCenterUI(spriteCursorFrameWhite, 75, 8, 16, 14)

	// Przyciski
	setUI(spriteBtnBuildPalisade, 251, 8, 16, 14)
	setUI(spriteBtnShield, 11, 22, 16, 14)
	// setUI(SPRITE_BTN_, 123, 22, 16, 14) // btn_map @todo: pewnie można usunąć, bo nie używam
	setUI(spriteBtnRepair, 139, 22, 16, 13)
	setUI(spriteBtnBuildBarn, 155, 21, 16, 14)
	setUI(spriteBtnBuildBarracks, 171, 21, 16, 14)
	setUI(spriteBtnBuildTemple, 187, 21, 16, 14)
	setUI(spriteBtnBuildBarracks2, 203, 21, 16, 14)
	setUI(spriteBtnBuildAcademy, 219, 21, 16, 14)

	// Czary @todo: ogarnij potrójne ikonki dla przycisku. Gdzieś jest już gotowa funkcja do tego
	setUI(spriteBtnSpellVision, 235, 8, 16, 14) // spell_vision

	spriteRegistry[609] = spriteDef{atlasID: atlasUnits1, x: 303, y: 112, w: 16, h: 14}

	spriteRegistry[610] = spriteDef{atlasID: atlasUnits2, x: 258, y: 69, w: 16, h: 14}

	spriteRegistry[611] = spriteDef{atlasID: atlasUnits2, x: 255, y: 111, w: 16, h: 14}

	// Zwłoki
	// @todo: kompletnie porąbane nazwy!
	setUI(spriteEffectSkeleton_00, 219, 50, 16, 14) // dead_0
	setUI(spriteEffectSkeleton_01, 219, 64, 16, 14) // dead_1
	setUI(spriteEffectSkeleton_02, 219, 78, 16, 14) // dead_2

	// Rany
	setUI(615, 219, 8, 16, 14) // hit_0
	setUI(616, 203, 8, 16, 14) // hit_1

	// Ogień
	for i := uint16(0); i <= 13; i++ {
		id := spriteFireStart + i
		if id < maxSpriteID {
			spriteRegistry[id] = spriteDef{
				atlasID: atlasUI,
				x:       11 + (i * 16),
				y:       176,
				w:       16, h: 14,
			}
		}
	}
}

// Jednostki
// ID: 700 + (unitType * 200) + (Frame * 8) + Direction
// Frame 0: Idle, Frame 1: Walk1, Frame 2: Walk2, Frame 3: Attack1, Frame 4: Attack2.
// Dodatkowo: offset 40 = świeżo zabity, offset 41 = początek rozkładu
func initUnitSprites() {
	type unitCfg struct {
		Atlas battleAtlasID
		BaseX uint16
		BaseY uint16
		Melee bool
	}
	configs := map[int]unitCfg{
		0:  {Atlas: atlasUnits1, BaseX: 0, BaseY: 0, Melee: false},      // COW
		1:  {Atlas: atlasUnits1, BaseX: 0, BaseY: 42, Melee: true},      // AXEMAN
		2:  {Atlas: atlasUnits1, BaseX: 0, BaseY: 84, Melee: false},     // ARCHER
		3:  {Atlas: atlasUnits1, BaseX: 0, BaseY: 126, Melee: false},    // PRIESTESS
		4:  {Atlas: atlasUnits2, BaseX: 0, BaseY: 0, Melee: false},      // PRIEST
		5:  {Atlas: atlasUnits2, BaseX: 0, BaseY: 42, Melee: true},      // SWORDSMAN
		6:  {Atlas: atlasUnits2, BaseX: 0, BaseY: 84, Melee: false},     // SPEARMAN
		7:  {Atlas: atlasUnits2, BaseX: 0, BaseY: 126, Melee: true},     // COMMANDER
		8:  {Atlas: atlasBuildings, BaseX: 0, BaseY: 0, Melee: true},    // BEAR
		9:  {Atlas: atlasBuildings, BaseX: 0, BaseY: 42, Melee: true},   // STRZYGA
		10: {Atlas: atlasUnits1, BaseX: 160, BaseY: 0, Melee: false},    // SHEPHERD
		11: {Atlas: atlasUnits1, BaseX: 160, BaseY: 126, Melee: false},  // MAGE
		12: {Atlas: atlasBuildings, BaseX: 224, BaseY: 0, Melee: false}, // CROSSBOW
	}

	const StartID = 700
	const StepID = 200

	for uType := 0; uType <= 12; uType++ {
		cfg, exists := configs[uType]
		if !exists {
			continue
		}

		baseID := StartID + (uType * StepID)

		for frame := range 5 {
			if cfg.Melee && (frame == 3 || frame == 4) {
				generateMeleeAttackFrames(baseID, frame, cfg.Atlas, cfg.BaseY)
				continue
			}

			colX := cfg.BaseX + uint16(frame*32)

			if uType == 12 && frame > 0 {
				switch frame {
				case 1:
					colX = 240
				case 2:
					colX = 256
				case 3:
					colX = 272
				}
			}

			for dir := 0; dir < 8; dir++ {
				var x, y uint16
				var w, h uint16 = 16, 14
				var flip bool

				y0 := cfg.BaseY
				y1 := cfg.BaseY + 14
				y2 := cfg.BaseY + 28

				switch dir {
				case 0:
					y = y0
					x = colX + 16
					flip = false
				case 1:
					y = y0
					x = colX
					flip = true
				case 2:
					y = y1
					x = colX
					flip = true
				case 3:
					y = y2
					x = colX
					flip = true
				case 4:
					y = y2
					x = colX + 16
					flip = false
				case 5:
					y = y2
					x = colX
					flip = false
				case 6:
					y = y1
					x = colX
					flip = false
				case 7:
					y = y0
					x = colX + 16
					flip = false
				}

				finalID := baseID + (frame * 8) + dir
				if finalID < maxSpriteID {
					spriteRegistry[finalID] = spriteDef{
						atlasID: cfg.Atlas,
						x:       x,
						y:       y,
						w:       w,
						h:       h,
						flipX:   flip,
						offX:    0,
						offY:    0,
					}
				}
			}
		}

		var x40, y40, x41, y41 uint16
		switch uType {
		case 0: // Cow
			x40, y40, x41, y41 = 112, 14, 144, 14
		case 1: // Axeman
			x40, y40, x41, y41 = 120, 70, 200, 70
		case 2: // Archer
			x40, y40, x41, y41 = 112, 98, 144, 98
		case 3: // Priestess
			x40, y40, x41, y41 = 111, 140, 143, 140
		case 4: // Priest
			x40, y40, x41, y41 = 112, 14, 144, 14
		case 5: // Swordsman
			x40, y40, x41, y41 = 120, 70, 200, 70
		case 6: // Spearman
			x40, y40, x41, y41 = 112, 98, 144, 98
		case 7: // Commander
			x40, y40, x41, y41 = 120, 154, 200, 154
		case 8: // Bear
			x40, y40, x41, y41 = 120, 28, 200, 28
		case 9: // Strzyga
			x40, y40, x41, y41 = 120, 70, 200, 70
		case 10: // Shepherd
			x40, y40, x41, y41 = 160, 98, 192, 98
		case 11: // Mage
			x40, y40, x41, y41 = 224, 98, 256, 98
		case 12: // Crossbow
			x40, y40, x41, y41 = 239, 56, 271, 56
		default:
			continue
		}

		id40 := baseID + 40
		if id40 < maxSpriteID {
			spriteRegistry[id40] = spriteDef{
				atlasID: cfg.Atlas,
				x:       x40,
				y:       y40,
				w:       16,
				h:       14,
				flipX:   false,
				offX:    0,
				offY:    0,
			}
		}

		id41 := baseID + 41
		if id41 < maxSpriteID {
			spriteRegistry[id41] = spriteDef{
				atlasID: cfg.Atlas,
				x:       x41,
				y:       y41,
				w:       16,
				h:       14,
				flipX:   false,
				offX:    0,
				offY:    0,
			}
		}
	}
}

// Specjalna logika dla walczących wręcz.
// Odwzorowuje logikę z TEMP_bridge.go, walczące wręcz powinny działać.
func generateMeleeAttackFrames(baseUnitID int, frame int, atlas battleAtlasID, unitBaseY uint16) {
	// Współrzędne x w atlasie dla poszczególnych faz (hardcoded z bridge)
	// Bridge Col 0 (Diag), Col 1 (Up), Col 2 (Side), Col 3 (DownDiag), Col 4 (Down)

	for dir := range 8 {
		var x, y uint16
		var w, h uint16

		var offX, offY int8

		var flip bool

		// Domyślne wartości (żeby nie było zer)
		w = 16
		h = 14

		if frame == 3 { // === FAZA 1 ATAKU ===
			switch dir {
			case 0: // Góra (Wąska 16x21)
				x = 120
				y = unitBaseY + 7
				w = 16
				h = 21
				flip = false
				offY = -7 // Korekta wysokości
			case 1: // Góra-Prawo (Szeroka 24x21)
				x = 96
				y = unitBaseY + 7
				w = 24
				h = 21
				flip = true
				offX = -4
				offY = -7
			case 2: // Prawo (Wąska 24x14? Bridge mówi 24x14)
				x = 96
				y = unitBaseY + 28
				w = 24
				h = 14
				flip = true
				offX = -4 // Poszerzenie
			case 3: // Dół-Prawo (Szeroka 24x21) - Bridge: Col 3 (AbsX 152, RelY 0 ??)
				// Bridge Row 3 Col 3: AbsX 152, RelY 0.
				x = 152
				y = unitBaseY
				w = 24
				h = 21
				flip = true
				offX = -4
				offY = -7
			case 4: // Dół (Wąska 16x21) - Bridge: Col 4 (AbsX 136, RelY 0)
				x = 136
				y = unitBaseY
				w = 16
				h = 21
				flip = false
				offY = -7
			case 5: // Dół-Lewo (Flip DP)
				x = 152
				y = unitBaseY
				w = 24
				h = 21
				flip = false
				offX = -4
				offY = -7
			case 6: // Lewo (Flip P)
				x = 96
				y = unitBaseY + 28
				w = 24
				h = 14
				flip = false
				offX = -4
			case 7: // Góra-Lewo (Flip GP)
				x = 96
				y = unitBaseY + 7
				w = 24
				h = 21
				flip = false
				offX = -4
				offY = -7
			}
		} else { // === FAZA 2 ATAKU (Frame 4) ===
			switch dir {
			case 0: // Góra
				x = 200
				y = unitBaseY + 7
				w = 16
				h = 21
				flip = false
				offY = -7
			case 1: // GP
				x = 176
				y = unitBaseY + 7
				w = 24
				h = 21
				flip = true
				offX = -4
				offY = -7
			case 2: // P
				x = 200
				y = unitBaseY + 28
				w = 16
				h = 14
				flip = true
			case 3: // DP
				x = 152
				y = unitBaseY + 21
				w = 24
				h = 21
				flip = true
				offX = -4
				offY = -7
			case 4: // D
				x = 136
				y = unitBaseY + 21
				w = 16
				h = 21
				flip = false
				offY = -7
			case 5: // DL
				x = 152
				y = unitBaseY + 21
				w = 24
				h = 21
				flip = false
				offX = -4
				offY = -7
			case 6: // L
				x = 200
				y = unitBaseY + 28
				w = 16
				h = 14
				flip = false
			case 7: // GL
				x = 176
				y = unitBaseY + 7
				w = 24
				h = 21
				flip = false
				offX = -4
				offY = -7
			}
		}

		finalID := baseUnitID + (frame * 8) + dir
		if finalID < maxSpriteID {
			spriteRegistry[finalID] = spriteDef{
				atlasID: atlas,
				x:       x, y: y, w: w, h: h,
				flipX: flip,
				offX:  offX, offY: offY,
			}
		}
	}
}

func initBuildingSprites() {
	assetID := atlasBuildings

	setBuilding := func(id int, x, y, w, h uint16) {
		if id < maxSpriteID {
			spriteRegistry[id] = spriteDef{
				atlasID: assetID,
				x:       x, y: y,
				w: w, h: h,
				offX: 0, offY: 0,
			}
		}
	}
	// Budowa → y:168
	for i := 127; i <= 135; i++ {
		setBuilding(i, uint16((i-127)*16), 168, 16, 14)
	}

	// Zgliszcza → y:182
	for i := 257; i <= 265; i++ {
		setBuilding(i, uint16((i-257)*16), 182, 16, 14)
	}

	// Budynek główny
	// ID 137-155 → y:84
	for i := 137; i <= 155; i++ {
		setBuilding(i, uint16((i-137)*16), 84, 16, 14)
	}
	// ID 156 Most
	setBuilding(156, 304, 84, 16, 14)

	// Obora
	// ID 157-175 → y:98
	for i := 157; i <= 175; i++ {
		setBuilding(i, uint16((i-157)*16), 98, 16, 14)
	}
	// Chata drwali
	// ID 177-195 → y:112
	for i := 177; i <= 195; i++ {
		setBuilding(i, uint16((i-177)*16), 112, 16, 14)
	}
	// Świątynia
	// ID 197-215 → y:126
	for i := 197; i <= 215; i++ {
		setBuilding(i, uint16((i-197)*16), 126, 16, 14)
	}
	// Chata wojów
	// ID 217-235 → y:140
	for i := 217; i <= 235; i++ {
		setBuilding(i, uint16((i-217)*16), 140, 16, 14)
	}
	// Dwór
	// ID 237-255 → y:154
	for i := 237; i <= 255; i++ {
		setBuilding(i, uint16((i-237)*16), 154, 16, 14)
	}
}

// Mapowanie battleAtlasID → rawAssetDef {TopChunk, BotChunk, PaletteID}.
var atlasDefinitions = map[battleAtlasID]rawAssetDef{
	atlasUI:        {3, 18, 3}, // UI
	atlasUnits1:    {4, 19, 4}, // Jednostki
	atlasUnits2:    {5, 20, 3}, // Jednostki i pociski
	atlasBuildings: {7, 22, 3}, // Budynki
}
