package main

import rl "github.com/gen2brain/raylib-go/raylib"

// constants.go

// ============================================================================
// STAŁE I KONFIGURACJA
// ============================================================================

/*## 1. Odchudzanie struktur danych (Data Locality & Cache Friendliness)
**Cel:** Zmniejszenie rozmiaru pojedynczych obiektów, aby więcej mieściło się w Cache L1 procesora (64 bajty).
	Mniej skakania po RAM = szybsza gra.

### A. Struktura `Tile` (Obecnie: ~56 bajtów)
* **Problem:** Używamy `int` (8 bajtów) i wskaźników tam, gdzie wystarczą małe liczby. Każdy odczyt mapy zapycha cache.
* **Rozwiązanie:**
    * `TextureID`, `EffectID`: Zmień `int` -> `uint16` (0-65535 wystarczy).
    * `MovementCost`: Zmień `float64` -> `uint8` (koszt 1-255, np. 10=droga, 20=trawa).
    * `Flags`: Zamiast wielu `bool` (każdy bierze 1 bajt + padding), użyj jednego `uint8` i masek bitowych
	(np. bit 1: Walkable, bit 2: Water).
* **Zysk:** Rozmiar spadnie do ~12-16 bajtów. Mapa będzie ładowana 4x szybciej.

### B. Struktura `unit` (Obecnie: God Object)
* **Problem:** Struktura jest ogromna. Zawiera dane potrzebne co klatkę (HP, x, y) obok danych rzadkich
	(Inventory, History, Wymiona).
* **Rozwiązanie:** Podział na dane "Gorące" i "Zimne".
    * Wyrzuć `string` (np. `AnimationType`) -> zamień na `enum` (stałe `int`/`byte`). Porównywanie liczb jest 100x
	szybsze niż napisów.
    * Wyrzuć wskaźniki `*pathNode` (slice wskaźników) -> użyj płaskiej tablicy punktów.

## 2. Zarządzanie Pamięcią i Garbage Collector (GC)
**Cel:** Przestać męczyć GC skanowaniem tysięcy małych obiektów.

* **Problem:** `Tile` trzyma wskaźniki `*unit` i `*building`. GC musi skanować całą planszę (4356 pól), żeby sprawdzić,
	czy coś nie zniknęło.
* **Rozwiązanie (Index-based approach):**
    * Zamiast wskaźników (`*unit`), trzymaj w `Tile` numer ID (`UnitID int`).
    * Trzymaj wszystkie jednostki w jednej, wielkiej, alokowanej na starcie tablicy `[]unit`.
    * Dostęp: `GlobalUnits[tile.UnitID]`.
* **Efekt:** Mapa staje się "niewidzialna" dla GC (jeśli nie ma wskaźników), a procesor kocha iterować po ciągłych
	tablicach.

## 3. Logika "Raz a dobrze" (Lekcja z trawy)
**Cel:** Nie liczyć w kółko tego, co się nie zmienia.

* **Zasada:** Jeśli coś jest obliczane w pętli `Draw` lub `updateProjectile` (60 razy na sekundę), zadaj pytanie:
	"Czy wynik zmienił się od ostatniej klatki?".
   * Przykład (Twój sukces): Generowanie wariantów trawy przeniesione z pętli rysowania do inicjalizacji mapy.
   * Kandydaci: Pathfinding (nie szukaj ścieżki co klatkę, jeśli cel się nie ruszył), sortowanie obiektów do rysowania
	(Z-index).
## 4. Typy danych
* Unikaj `int` (64-bit) dla małych wartości w dużych tablicach. w pętlach i obliczeniach `int` jest OK
	(rejestry CPU są duże), ale w **pamięci** (struktury) każdy bajt się liczy.*/

const (
	// @todo: @reminder: Będę musiał tutaj jeszcze wrócić i się ogarnąć prędkość gry. Obecnie jest ona zbyt szybka.
	logicSpeedDivisor   = 1
	envAnimSpeedDivisor = 15
	targetFPS           = 60.0
	tickRate            = 1.0 / targetFPS
)

const (
	boardMaxX             uint8 = 66 // Górna granica szerokości mapy
	boardMaxY             uint8 = 66 // Górna granica wysokości mapy
	maxUnitsPerPlayer     uint8 = 40 // Górna granica ilości wojsk na gracza
	maxBuildingsPerPlayer uint8 = 20 // Górna granica ilości budynków na gracza

	// Wymiary kafelka mapy.
	tileWidth  uint8 = 16 // Szerokość kafelka
	tileHeight uint8 = 14 // Wysokość kafelka

	// Progi poziomów dla wyjątkowych jednostek, skazują kiedy dane jednostki są dostępne.
	shepherdLevel    uint8 = 26
	mageLevel        uint8 = 40
	crossbowmanLevel uint8 = 32

	// Rodzaj bitwy
	// O ile dobrze rozumiem, to wpływa na usposobienie SI w czasie bitwy.
	boardVillage    uint8 = 0 // SI ma osadę
	boardBattleDyn  uint8 = 1 // SI jest napastliwe
	boardBattleStat uint8 = 2 // SI walczy niemrawo
	boardNothing    uint8 = 3 // SI jest bezczynna

	// Warunki zakończenia bitwy.
	endKillAll uint8 = 0 // Ubij wszystkich
	endRescue  uint8 = 1 // Uratuj kogoś
	endBuild   uint8 = 2 // Wybuduj (określoną liczbę budowli)
	endNothing uint8 = 3 // NIC (NIE WIEM CO TO ZNACZY)
	endKillOne uint8 = 4 // Ubij wybraną wrogą jednostkę

	// Wynik bitwy.
	outcomePending uint8 = 0
	outcomeVictory uint8 = 1
	outcomeDefeat  uint8 = 2

	// Barwy graczy.
	colorNone     uint8 = 0 // Niczyje
	colorRed      uint8 = 1 // Czerwień
	colorGreen    uint8 = 2 // Zieleń
	colorYellow   uint8 = 3 // Żółć
	colorBlue     uint8 = 4 // Niebieski
	colorGray     uint8 = 5 // Szary
	maxGameColors uint8 = 6 // @reminder: jeżeli chciałbym dodać nową barwę, to muszę tutaj to uwzględnić

	// Rozmiary wiadomości.
	msgFontColor = 255

	// Przesunięcie drzew.
	treeOffsetX                 float32 = 8.0
	treeCrownTextureOffset      uint16  = 7
	burntTreeCrownTextureOffset uint16  = 2

	// Do uruchomienia obrazów.
	walkAnimationFrames = 5 // Ile klatek ma chodzenie @todo: ustaw prawdziwą wartość!
	animationSpeed      = 4 // Zmień klatkę co 4 cykle logiki gry (mniej = szybciej) @todo: dopasuj do oryginału
)

const (
	// Minimapa.
	minimapOffsetX            float32 = 5
	minimapOffsetY            float32 = 5
	minimapDisplayWidth       float32 = 110
	minimapDisplayHeight      float32 = 100
	minimapClickDragThreshold float32 = 5.0
)

const (
	// Rozkazy dla jednostek.
	cmdIdle   = uint16(iota) // Bezczynność
	cmdMove                  // Ruch
	cmdGoto                  // Idź do (czy to jest teleport maga?!)
	cmdFlee                  // ucieczka
	cmdGraze                 // Wypasaj
	cmdAttack                // Napad
	cmdStop                  // Zatrzymaj się

	// Czary.
	cmdCastSpell
	cmdMagicLightning // Rzuć czar
	cmdMagicShield    // Rzuć magiczną tarczę
	cmdMagicFire      // Rzuć deszcz ognia (NIE JESTEM PEWIEN)
	cmdMagicSight     // Rzuć dalekowidztwo

	// Gospodarka.
	cmdProduce // Wytwarzaj
	cmdMilking // dojenie

	// Budowa.
	cmdRepairStructure // Naprawa uszkodzonego budynku
	cmdBuildStructure  // Budowa nowego budynku
)

const (
	// Ile kosztują czary.
	spellCostMagicShield uint16 = 50
	spellCostMagicSight  uint16 = 25
	// Zasięg ?
	spellRangeMagicSight uint8 = 14
)

const (
	// Ile trwają czary.
	spellDurationMagicShield uint = 100
)

const (
	// Rany.
	maxWoundsCount int    = 6
	severeDamage   uint16 = 4
)

const (
	// Ile naprawia drwal jednym uderzeniem.
	repairAmountPlayer uint16 = 2
	repairAmountAI     uint16 = 5
)

const (
	// Opóźnienia dla LegacyPhase.
	maxPhaseDelay uint16 = 16
	minPhaseDelay uint16 = 0
)

const (
	// rzeczy związane ze zwłokami.
	corpsesFrameIndexOffset uint16  = 40
	corpsesPhase2           uint8   = 2
	corpsesPhase1           uint8   = 1
	corpsesMaxAlpha         float32 = 255.0
)

// ↓↓↓Płonięcie kafelków↓↓↓.
const (
	// rzeczy związane z płomieniami i popiołem.
	bigBurn uint16 = 80
	midBurn uint16 = 224
	minBurn uint16 = 240

	totalAshLifetime    uint16  = 600
	ashDecayRate        float32 = 0.005
	ashAccumulationRate float32 = 0.01

	burnDamage uint16 = 3 // Tyle zadają płomienie co 8 tików
)

const (
	ashDecaying uint8 = iota
	ashFinished
)

// ↓↓↓Upadające drzewo↓↓↓
// Liczba uderzeń unitAxeman potrzebna do ścięcia suchego drzewa.
// @todo: @reminder: zmień na docelową liczbę po skończeniu wdrażania mechanizmu ścinania drzew! 23.04.2026
// @reminder to powinien być odpowiednik 150 obrażeń!
const strikesToCutTree uint8 = 5

type treeState uint8

const (
	noTree treeState = iota
	treeStraight
	treeLeaning
	treeFalling
	treeImpact
	treeFell
)

// @todo: sprawdź ile to było w oryginale.
const fallingTreeDamage uint16 = 10_000

// @todo: jeszcze nie zrobione w drawingBattle, ale niezbędne!
var victoryPointColors = []rl.Color{
	{R: 252, G: 252, B: 188, A: 255}, //nolint:mnd    // LightYellow
	rl.Yellow,                        // Yellow
	rl.Gold,                          // DarkYellow
	rl.Yellow,                        // Yellow
}

// @reminder: O ile dobrze rozumiem, to ograniczenie doświadczenia było
// ustawione na 240, ale 224 było ostatnim „poziomem. Pomiędzy 224 a 240
// jest martwa przestrzeń, która nic już nie zmienia. Dlatego ustawiam
// ogranicznik na 224.
var (
	experienceCap         uint = 224
	experienceCasterBonus uint = 2
)

// @todo: to chyba powiązane ze zdobywanym doświadczeniem przez jednostkę.
var dDamage = [15]uint8{1, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 6, 7, 8, 9}

// @todo: powiązane ze zdobywanym doświadczeniem przez jednostkę.
var dArmor = [15]uint8{0, 1, 1, 2, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 6}

// @todo: powiązane ze zdobywanym doświadczeniem przez jednostkę czarującą.
var dMana = [15]uint16{60, 80, 85, 90, 120, 140, 150, 160, 170, 180, 190, 200, 220, 240, 280}

// Statystyki jednostek.
// @reminder: wyłączyłem lintera ponieważ tutaj powinny znajdować się te zaczarodziejskie liczby!

var unitDefs = map[unitType]unitStats{
	//nolint:mnd
	unitArcher: {
		Name:        "Archer",
		MaxHP:       100,
		BaseDamage:  5,
		BaseArmor:   0,
		SightRange:  4,
		AttackRange: 3,
		Cost:        0, // @todo: zamień na 650
		MoveDelay:   8,
		MaxMana:     0,
	},
	//nolint:mnd
	unitAxeman: {
		Name:        "Axeman",
		MaxHP:       100,
		BaseDamage:  10,
		BaseArmor:   0,
		SightRange:  3,
		AttackRange: 1,
		Cost:        0, // @todo: 450
		MoveDelay:   10,
		MaxMana:     0,
	},
	//nolint:mnd
	unitBear: {
		Name:        "Bear",
		MaxHP:       300,
		BaseDamage:  25,
		BaseArmor:   3,
		SightRange:  4,
		AttackRange: 1,
		Cost:        0, // @todo: 1050
		MoveDelay:   16,
		MaxMana:     0,
	},
	//nolint:mnd
	unitCommander: {
		Name:        "Commander",
		MaxHP:       150,
		BaseDamage:  30,
		BaseArmor:   5,
		SightRange:  4,
		AttackRange: 1,
		Cost:        0, // @todo: 1050
		MoveDelay:   10,
		MaxMana:     0,
	},
	//nolint:mnd
	unitCow: {
		Name:        "Cow",
		MaxHP:       100,
		BaseDamage:  0,
		BaseArmor:   0,
		SightRange:  2,
		AttackRange: 0,
		Cost:        0, // @todo: 450
		MoveDelay:   12,
		MaxMana:     0,
	},
	//nolint:mnd
	unitCrossbowman: {
		Name:        "Crossbowman",
		MaxHP:       130,
		BaseDamage:  25,
		BaseArmor:   5,
		SightRange:  6,
		AttackRange: 5,
		Cost:        0, // @todo: zmień na 1250 po próbach
		MoveDelay:   8,
		MaxMana:     0,
	},
	//nolint:mnd
	unitMage: {
		Name:        "Mage",
		MaxHP:       50,
		BaseDamage:  10,
		BaseArmor:   0,
		SightRange:  6,
		AttackRange: 6,
		Cost:        0, // @todo: zmień na 1250 po próbach
		MoveDelay:   12,
		MaxMana:     60,
	},
	//nolint:mnd
	unitPriest: {
		Name:        "Priest",
		MaxHP:       80,
		BaseDamage:  50,
		BaseArmor:   0,
		SightRange:  5,
		AttackRange: 5,
		Cost:        0, // @todo 1050
		MoveDelay:   16,
		MaxMana:     60,
	},
	//nolint:mnd
	unitPriestess: {
		Name:        "Priestess",
		MaxHP:       70,
		BaseDamage:  35,
		BaseArmor:   0,
		SightRange:  4,
		AttackRange: 3,
		Cost:        0, // @todo: 650
		MoveDelay:   14,
		MaxMana:     60,
	},
	//nolint:mnd
	unitShepherd: {
		Name:        "Shepherd",
		MaxHP:       40,
		BaseDamage:  0,
		BaseArmor:   0,
		SightRange:  3,
		AttackRange: 1,
		Cost:        0, // @todo popraw na 1250 po próbach
		MoveDelay:   12,
		MaxMana:     0,
	},
	//nolint:mnd
	unitSpearman: {
		Name:        "Spearman",
		MaxHP:       120,
		BaseDamage:  15,
		BaseArmor:   3,
		SightRange:  5,
		AttackRange: 4,
		Cost:        0, // @todo: 850
		MoveDelay:   10,
		MaxMana:     0,
	},
	//nolint:mnd
	unitSwordsman: {
		Name:        "Swordsman",
		MaxHP:       120,
		BaseDamage:  20,
		BaseArmor:   3,
		SightRange:  3,
		AttackRange: 1,
		Cost:        0, // @todo: 650
		MoveDelay:   10,
		MaxMana:     0,
	},
	//nolint:mnd
	unitUnknown: {
		Name:        "Unknown",
		MaxHP:       120,
		BaseDamage:  20,
		BaseArmor:   3,
		SightRange:  4,
		AttackRange: 1,
		Cost:        0, // @todo: 1050
		MoveDelay:   6,
		MaxMana:     60,
	},
}

const (
	// Wymiary budowli.
	normalBuildingSize uint8 = 3
	smallBuildingSize  uint8 = 1

	// === Nakładka dla budynków
	// Wymiary ramek.
	cornerLenPalisade float32 = 5.0
	cornerLenBuilding float32 = 10.0
	cornerLenUnit             = 5.0
	cornerThickness           = 1.5

	// Pasek życia.
	healthBarHeight          int32 = 2
	buildingHealthBarMarginY int32 = 2

	// Pojemność budynków.
	capacityRectW   int32 = 3
	capacityReactH  int32 = 2
	capacitySpacing int32 = 2
	capacityMarginX int32 = 2
	capacityMarginY int32 = 2
)

// @reminder: to nie mogą być stałe.
var (
	// Z góry określone barwy dla ramek.
	friendlyFrameColor = rl.NewColor(138, 132, 129, 255) //nolint:mnd
	enemyFrameColor    = rl.NewColor(112, 0, 0, 255)     //nolint:mnd
)

// @todo: Do zrobienia:
// Co z mostami?! cena 80.
// Co z drogami?! cena 45.
var buildingDefs = map[buildingType]buildingStats{
	buildingMain:
	//nolint:mnd
	{
		Name:  "Budynek Główny",
		Width: normalBuildingSize, Height: normalBuildingSize, Cost: 0, MaxHP: 400, MaxFood: 0,
		BaseTextureID: spriteConstructionStart, IsPalisade: false,
	},
	//nolint:mnd
	buildingBarn: {
		Name:  "Obora",
		Width: normalBuildingSize, Height: normalBuildingSize, Cost: 0, MaxHP: 350, MaxFood: 3,
		BaseTextureID: spriteConstructionStart, IsPalisade: false,
	},
	// @todo: cost 650
	//nolint:mnd
	buildingBarracks: {
		Name:  "Chata mieszkalna",
		Width: normalBuildingSize, Height: normalBuildingSize, Cost: 0, MaxHP: 350, MaxFood: 6,
		BaseTextureID: spriteConstructionStart, IsPalisade: false,
	}, // @todo: 850
	//nolint:mnd
	buildingTemple: {
		Name:  "Dwór mocy",
		Width: normalBuildingSize, Height: normalBuildingSize, Cost: 0, MaxHP: 350, MaxFood: 3,
		BaseTextureID: spriteConstructionStart, IsPalisade: false,
	}, // @todo: 1050
	//nolint:mnd
	buildingBarracks2: {
		Name:  "Chata wojów",
		Width: normalBuildingSize, Height: normalBuildingSize, Cost: 0, MaxHP: 350, MaxFood: 4,
		BaseTextureID: spriteConstructionStart, IsPalisade: false,
	}, // @todo ILE TO KOSZTOWAŁO?!
	//nolint:mnd
	buildingAcademy: {
		Name:  "Dwór rycerza",
		Width: normalBuildingSize, Height: normalBuildingSize, Cost: 0, MaxHP: 400, MaxFood: 1,
		BaseTextureID: spriteConstructionStart, IsPalisade: false,
	}, // @todo: 1050
	//nolint:mnd
	buildingPalisade: {
		Name:  "Palisada",
		Width: smallBuildingSize, Height: smallBuildingSize, Cost: 0, MaxHP: 120, MaxFood: 0,
		BaseTextureID: spriteConstructionStart, IsPalisade: true,
	}, // @todo: 60
	//nolint:mnd
	buildingBridge: {
		Name:  "Most", // todo: wszystkie staty! te są tymczasowe!
		Width: smallBuildingSize, Height: smallBuildingSize, Cost: 0, MaxHP: 120, MaxFood: 0,
		BaseTextureID: spriteConstructionStart, IsPalisade: true,
	}, // @todo: 60
	// @todo: dodaj drogę!
	buildingRoad: {
		Name:  "Droga", // todo: wszystkie staty! te są tymczasowe!
		Width: smallBuildingSize, Height: smallBuildingSize, Cost: 0, MaxHP: 0, MaxFood: 0,
		BaseTextureID: spriteRoadButton, IsPalisade: true,
	}, // @todo: 45
}

// Ważne: jednostki czarujące rozpoczynają grę z połową many i maxmana = 60, strzyga 0 many.
var maxManaData = [15]int{60, 80, 85, 90, 120, 140, 150, 160, 170, 180, 190, 200, 220, 240, 280}

// @todo: Ogarnij, czy nie da się tego lepiej załatwić.
// @reminder: Wydaje mi się, że po ósmym najeździe SI jest bezczynna.
var attackTime = [8]int{400, 0, 0, 0, 200, 400, 600, 700}

// Opisuje do kogo przynależą poszczególne ziemie.
// @todo: widok krain do zrobienia.
var provinceInit = [25]uint8{
	colorGreen, colorGreen, colorGreen, colorYellow, colorGreen, colorGreen, colorBlue,
	colorBlue, colorYellow, colorYellow, colorGreen, colorBlue, colorYellow, colorYellow, colorGray, colorBlue,
	colorBlue, colorYellow, colorGray, colorGray, colorGray, colorGray, colorGray, colorGray, colorGray,
}

// Ścieżki dźwiękowe.
// O ile dobrze rozumiem, to całość jest podzielona na dwie części.
// Pierwotna muzyka z wersji na dyskietki oraz dodatkowe wyprawy z CD.
// @todo: dźwięk jeszcze nie zrobiony.
var musicTrack = [52]uint8{
	9, 6, 7, 8, 9, 6, 7, 8, 9, 6, // powrót mirka
	7, 8, 9, 6, 7, 8, 9, 6, 7, 8,
	9, 6, 7, 8, 9,
	6, 11, 13, 10, 8, // przyjaciele
	9, 11, 6, 12, 14,
	7, 13, 12, 6, 10, 8,
	11, 14, 10, 8, 13,
	12, 6, 7, 14, 8, 10,
}

// @reminder @todo: zastanów się, czy nie można tego lepiej rozwiązać.
var legacyShiftX = [13][17]uint8{
	{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 13, 14, 16, 15, 15, 15, 15},    // 0: Krowa
	{0, 1, 2, 3, 4, 6, 8, 10, 12, 14, 16, 25, 15, 15, 15, 15, 15},   // 1: Drwal
	{0, 2, 4, 6, 8, 10, 12, 14, 16, 15, 15, 15, 15, 15, 15, 15, 15}, // 2: Łucznik
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 10, 11, 12, 13, 14, 16, 15, 15},     // 3: Kapłanka
	{0, 1, 2, 3, 4, 5, 6, 7, 7, 8, 9, 10, 11, 12, 13, 14, 15},       // 4: Kapłan
	{0, 1, 3, 5, 6, 8, 9, 11, 12, 14, 16, 15, 15, 15, 15, 15, 15},   // 5: Miecznik
	{0, 1, 3, 5, 6, 8, 9, 11, 12, 14, 16, 15, 15, 15, 15, 15, 15},   // 6: Włócznik
	{0, 1, 3, 5, 6, 8, 9, 11, 12, 14, 15, 15, 15, 15, 15, 15, 15},   // 7: Dowódca
	{0, 1, 2, 3, 4, 5, 6, 7, 7, 8, 9, 10, 11, 12, 13, 14, 15},       // 8: Niedźwiedź
	{0, 3, 5, 7, 9, 12, 14, 15, 15, 8, 9, 10, 11, 12, 13, 14, 15},   // 9: Strzyga
	{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 13, 14, 16, 15, 15, 15, 15},    // 10: Pastuch
	{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 13, 14, 16, 15, 15, 15, 15},    // 11: Mag
	{0, 2, 4, 6, 8, 10, 12, 14, 16, 15, 15, 15, 15, 15, 15, 15, 15}, // 12: Kusznik
}

// @reminder @todo: zastanów się, czy nie można tego lepiej rozwiązać.
var legacyShiftY = [13][17]uint8{
	{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 13, 14, 15, 15, 15, 15, 15},   // Krowa
	{0, 1, 2, 5, 6, 7, 8, 10, 11, 12, 14, 25, 15, 15, 15, 15, 15},  // Drwal
	{0, 1, 3, 5, 7, 9, 11, 13, 14, 15, 15, 15, 15, 15, 15, 15, 15}, // Łucznik
	{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 15},     // Kapłanka
	{0, 1, 2, 3, 4, 5, 6, 7, 7, 8, 9, 10, 11, 12, 13, 14, 14},      // Kapłan
	{0, 1, 2, 5, 6, 7, 8, 10, 11, 12, 14, 25, 15, 15, 15, 15, 15},  // Miecznik
	{0, 1, 2, 5, 6, 7, 8, 10, 11, 12, 14, 25, 15, 15, 15, 15, 15},  // Włócznik
	{0, 1, 2, 5, 6, 7, 8, 10, 11, 12, 14, 25, 15, 15, 15, 15, 15},  // Dowódca
	{0, 1, 2, 3, 4, 5, 6, 7, 7, 8, 9, 10, 11, 12, 13, 14, 14},      // Niedźwiedź
	{0, 3, 5, 7, 9, 12, 14, 14, 14, 8, 9, 10, 11, 12, 13, 14, 14},  // Strzyga
	{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 13, 14, 15, 15, 15, 15, 15},   // Pastuch
	{0, 1, 2, 3, 5, 6, 8, 9, 11, 12, 13, 14, 15, 15, 15, 15, 15},   // Mag
	{0, 1, 3, 5, 7, 9, 11, 13, 14, 15, 15, 15, 15, 15, 15, 15, 15}, // Kusznik
}

// @reminder @todo: zastanów się, czy nie można tego lepiej rozwiązać.
var legacyPhase = [13][19]uint8{
	{0, 0, 0, 1, 1, 1, 0, 0, 0, 2, 2, 2, 0, 6, 1, 1, 0, 8, 4},  // 0: Krowa
	{0, 0, 1, 1, 1, 0, 0, 2, 2, 2, 0, 6, 0, 0, 0, 0, 0, 6, 3},  // 1: Drwal
	{0, 0, 2, 2, 0, 0, 1, 1, 1, 6, 0, 0, 0, 0, 0, 0, 0, 4, 2},  // 2: Łucznik
	{0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 2, 2, 2, 2, 0, 6, 0, 11, 8}, // 3: Kapłanka
	{0, 0, 2, 2, 2, 2, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 13, 9}, // 4: Kapłan
	{0, 0, 2, 2, 2, 0, 0, 0, 1, 1, 1, 6, 0, 0, 0, 0, 0, 6, 3},  // 5: Miecznik
	{0, 0, 2, 2, 2, 0, 0, 0, 1, 1, 1, 6, 0, 0, 0, 0, 0, 6, 3},  // 6: Włócznik
	{0, 0, 2, 2, 2, 0, 0, 1, 1, 1, 1, 6, 0, 0, 1, 1, 0, 5, 2},  // 7: Dowódca
	{0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 2, 2, 2, 2, 0, 0, 8, 2},  // 8: Niedźwiedź
	{0, 1, 1, 0, 2, 2, 0, 6, 0, 0, 0, 2, 2, 2, 2, 0, 0, 4, 2},  // 9: Strzyga
	{0, 0, 0, 1, 1, 1, 0, 0, 0, 2, 2, 2, 0, 6, 1, 1, 0, 8, 4},  // 10: Pastuch
	{0, 0, 0, 1, 1, 1, 0, 0, 0, 2, 2, 2, 0, 6, 1, 1, 0, 8, 4},  // 11: Mag
	{0, 0, 2, 2, 0, 0, 1, 1, 1, 6, 0, 0, 0, 0, 0, 0, 0, 4, 2},  // 12: Kusznik
}

// UI, przyciski, wskaźnik mleka.
const (
	milkBarOffsetX float32 = 71.0
	milkBarY       float32 = 118.0
	milkBarWidth   float32 = 30.0
	milkBarHeight  float32 = 235.0
	maxMilk        float32 = 1850.0

	// Przyciski.
	uiAnchorOffsetX float32 = 18.2
	uiAnchorOffsetY float32 = 112.5
	btnWidth        float32 = 38.5
	btnHeight       float32 = 35.0
	// Odstępy pomiędzy przyciskami.
	btnMarginY float32 = 12.0

	// Liczba przycisków.
	uiActionMaxButtons int = 5
)

// Zestawienie nazw budynków z ich rodzajami.
// @reminder @todo: korzystanie z łańcuchów jest kosztowne, być może trzeba to zmienić.
// Mapa to duża przesada, może zmienię.
var buildingTypeMap = map[string]buildingType{
	"MAIN":      buildingMain,      // Budynek główny
	"BARN":      buildingBarn,      // Obora
	"BARRACKS":  buildingBarracks,  // Chata mieszkalna
	"TEMPLE":    buildingTemple,    // Dwór mocy
	"BARRACKS2": buildingBarracks2, // Chata wojów
	"ACADEMY":   buildingAcademy,   // Dwór rycerza (w tamtych czasach nie było jeszcze rycerzy!)
	"PALISADE":  buildingPalisade,  // Palisada
}

// @reminder @todo: korzystanie z łańcuchów jest kosztowne, być może trzeba to zmienić.
// Mapa to duża przesada, może zmienię.
var unitTypeMap = map[string]unitType{
	"COW":         unitCow,         // Krowa
	"AXEMAN":      unitAxeman,      // Drwal
	"ARCHER":      unitArcher,      // Łucznik
	"PRIESTESS":   unitPriestess,   // Kapłanka
	"PRIEST":      unitPriest,      // Kapłan
	"SWORDSMAN":   unitSwordsman,   // Miecznik
	"SPEARMAN":    unitSpearman,    // Włócznik
	"COMMANDER":   unitCommander,   // Rycerz (w tamtych czasach nie było jeszcze rycerzy!)
	"BEAR":        unitBear,        // Niedźwiedź
	"UNKNOWN":     unitUnknown,     // Strzyga
	"SHEPHERD":    unitShepherd,    // Pastuch
	"MAGE":        unitMage,        // Mag
	"CROSSBOWMAN": unitCrossbowman, // Kusznik
}

const (
	// Przypisanie wartości liczbowych każdemu z rodzajów jednostek.
	unitCow         unitType          = iota // Krowa
	unitAxeman                               // Drwal
	unitArcher                               // Łucznik
	unitPriestess                            // Kapłanka
	unitPriest                               // Kapłan
	unitSwordsman                            // Miecznik
	unitSpearman                             // Włócznik
	unitCommander                            // Rycerz (w tamtych czasach nie było jeszcze rycerzy!)
	unitBear                                 // Niedźwiedź
	unitUnknown                              // Strzyga
	unitShepherd                             // Pastuch
	unitMage                                 // Mag
	unitCrossbowman                          // Kusznik
	unitNone                                 // Zapychacz jakiś @todo: usuń później.
	unitTypeCount   = unitCrossbowman        // @reminder: nie wiem na ile to potrzebne;zawsze ostatnie
)

const (
	// Przypisanie wartości liczbowych każdemu z budynków.
	buildingMain      buildingType = iota + 1 // Budynek główny
	buildingBarn                              // Obora
	buildingBarracks                          // Chata mieszkalna
	buildingTemple                            // Dwór mocy
	buildingBarracks2                         // Chata wojów
	buildingAcademy                           // Dwór
	buildingPalisade                          // Palisada
	buildingBridge                            // Most
	buildingRoad                              // Droga @reminder: nie wiem, czy tak zostanie na dłużej
)

// Rzeczy związane z budynkami.
const (
	initialConstructionHP = 30
	buildingArmor         = 10
)

// Jest to lewy górny róg tekstury, 135 prawy dolny.
// Dodatkowo każdy budynek ma swoją indywidualną „w budowie”.
// @reminder: chyba mogę to przepisać korzystając ze stałych w assets_*.go 1.05.2026.
var constructionTemplatePhase01 = [][]uint16{
	{127, 128, 129},
	{130, 131, 132},
	{133, 134, 135},
}

var constructionTemplatesPhase02 = map[buildingType][][]uint16{
	buildingAcademy: {
		{247, 248, 249},
		{250, 251, 252},
		{253, 254, 255},
	},
	buildingBarn: {
		{167, 168, 169},
		{170, 171, 172},
		{173, 174, 175},
	},
	buildingBarracks: {
		{187, 188, 189},
		{190, 191, 192},
		{193, 194, 195},
	},
	buildingBarracks2: {
		{227, 228, 229},
		{230, 231, 232},
		{233, 234, 235},
	},
	buildingMain: {
		{147, 148, 149},
		{150, 151, 152},
		{153, 154, 155},
	},
	buildingTemple: {
		{207, 208, 209},
		{210, 211, 212},
		{213, 214, 215},
	},
}

// Zestawienie rodzajów budynku z teksturami, które się na niego składają.
var buildingTemplates = map[buildingType][][]uint8{
	buildingMain: { // Budynek główny
		{137, 138, 139},
		{140, 141, 142},
		{143, 144, 145},
	},
	buildingBarn: { // Obora
		{157, 158, 159},
		{160, 161, 162},
		{163, 164, 165},
	},
	buildingBarracks: { // Chata mieszkalna
		{177, 178, 179},
		{180, 181, 182},
		{183, 184, 185},
	},
	buildingTemple: { // Dom mocy
		{197, 198, 199},
		{200, 201, 202},
		{203, 204, 205},
	},
	buildingBarracks2: { // Chata wojów
		{217, 218, 219},
		{220, 221, 222},
		{223, 224, 225},
	},
	buildingAcademy: { // Dwór rycerza (w tamtych czasach nie było jeszcze rycerzy!)
		{237, 238, 239},
		{240, 241, 242},
		{243, 244, 245},
	},
}

// Przypisane tekstur do animacji na budynkach.
var flagAnimationMap = map[uint8]uint8{
	138: 146, 146: 138, //nolint:mnd
	162: 166, 166: 162, //nolint:mnd
	179: 186, 186: 179, //nolint:mnd
	201: 206, 206: 201, //nolint:mnd
	221: 226, 226: 221, //nolint:mnd
	239: 246, 246: 239, //nolint:mnd
}

var (
	// Barwy gracza (czerwone) do podmiany.
	playerColors = []rl.Color{
		// Ciemny czerwony.
		{R: 112, G: 0, B: 0, A: 255}, //nolint:mnd
		// Średni czerwony.
		{R: 252, G: 88, B: 88, A: 255}, //nolint:mnd
		// Jasny czerwony.
		{R: 252, G: 152, B: 152, A: 255}, //nolint:mnd
	}

	// Mapowanie kolorów wrogich plemion.
	// @reminder: są to wartości dobrane w oparciu o oryginalne
	// grafiki. Dodałem nolint, żeby ograniczyć szum.
	enemyColors = map[uint8][]rl.Color{
		colorGreen: {
			{R: 0, G: 128, B: 0, A: 255},     //nolint:mnd
			{R: 0, G: 180, B: 0, A: 255},     //nolint:mnd
			{R: 144, G: 238, B: 144, A: 255}, //nolint:mnd
		},
		colorBlue: {
			{R: 0, G: 0, B: 172, A: 255},     //nolint:mnd
			{R: 88, G: 88, B: 252, A: 255},   //nolint:mnd
			{R: 152, G: 152, B: 252, A: 255}, //nolint:mnd
		},
		colorYellow: {
			{R: 172, G: 172, B: 0, A: 255},   //nolint:mnd
			{R: 252, G: 252, B: 88, A: 255},  //nolint:mnd
			{R: 252, G: 252, B: 188, A: 255}, //nolint:mnd
		},
		colorGray: {
			{R: 88, G: 88, B: 88, A: 255},    //nolint:mnd
			{R: 172, G: 172, B: 172, A: 255}, //nolint:mnd
			{R: 220, G: 220, B: 220, A: 255}, //nolint:mnd
		},
	}
)

const (
	corpseDecayTime = 2100
	// ↓↓↓ corpseDecayTime/3.
	corpsesPhaseDuration = 700
)

const (
	firstCampaignFirstLevel uint8 = 15 // Pierwsza mapa pierwszej kampanii
	// @reminder: @todo: Dodaj kolejne pierwsze poziomy.
)

// ===== UI
// Opis co, kiedy, jak może wytworzyć dany budynek
// przecież to aż krzyczy, że coś zrąbałem dla budynków!
// @todo: Zastąpić obecne rozwiązanie z unitNone typem sumowanym.
// Zamknięty interfejs z typami BuildAction oraz ProduceAction.
// Pozwoli się to pozbyć zbędnego unitNone oraz uciszania lintera.
var buildingRecipes = map[buildingType][]buildingAction{
	buildingMain: {
		// Przycisk 1/5
		{
			BuildingType: buildingTemple,
			Label:        buildingDefs[buildingTemple].Name,
			MinLevel:     0,
			IconID:       spriteBtnBuildTemple,
			UnitType:     unitNone,
		},

		// Przycisk 2/5
		{
			BuildingType: buildingBarracks,
			Label:        buildingDefs[buildingBarracks].Name,
			MinLevel:     0,
			IconID:       spriteBtnBuildBarracks,
			UnitType:     unitNone,
		},
		// Przycisk 3/5
		{
			BuildingType: buildingBarn,
			Label:        buildingDefs[buildingBarn].Name,
			MinLevel:     0,
			IconID:       spriteBtnBuildBarn,
			UnitType:     unitNone,
		},
		// Przycisk 4/5
		{
			BuildingType: buildingRoad,
			Label:        buildingDefs[buildingRoad].Name,
			MinLevel:     0,
			IconID:       spriteRoadButton,
			UnitType:     unitNone,
		},
		// Przycisk 5/5
		{
			BuildingType: buildingBridge,
			Label:        buildingDefs[buildingBridge].Name,
			MinLevel:     0,
			IconID:       spriteBridgeEnd,
			UnitType:     unitNone,
		},
		// N/D główny nie ma piątego przycisku o ile dobrze pamiętam

	},
	buildingBarn: {
		// Przycisk 1/5
		//nolint:exhaustruct
		{
			// @todo: czy pastuj zabierał miejsce w oborze?
			UnitType: unitShepherd,
			Label:    "Pastuch",
			MinLevel: 0, // @todo: SHEPHERD_LEVEL powinno być
			IconID:   0,
		},
		// Przycisk 2/5
		//nolint:exhaustruct
		{
			UnitType: unitCow,
			Label:    "Krowa",
			MinLevel: 0,
			IconID:   0,
		},
		// Przycisk 3/5
		// Przycisk 4/5
		// Przycisk 5/5

	},
	buildingBarracks: {
		// Przycisk 1/5
		{
			BuildingType: buildingBarracks2,
			Label:        buildingDefs[buildingBarracks2].Name,
			MinLevel:     0,
			IconID:       spriteBtnBuildBarracks2,
			UnitType:     unitNone,
		},
		// Przycisk 2/5
		//nolint:exhaustruct
		{
			UnitType: unitArcher,
			Label:    "Łucznik",
			MinLevel: 0,
			IconID:   0,
		},
		// Przycisk 3/5
		//nolint:exhaustruct
		{
			UnitType: unitAxeman,
			Label:    "Drwal",
			MinLevel: 0,
			IconID:   0,
		},
		// Przycisk 4/5
		// Przycisk 5/5

	},
	buildingBarracks2: {
		// Przycisk 1/5
		{
			BuildingType: buildingAcademy,
			Label:        buildingDefs[buildingAcademy].Name,
			MinLevel:     0,
			IconID:       spriteBtnBuildAcademy,
			UnitType:     unitNone,
		},
		// Przycisk 2/5
		//nolint:exhaustruct
		{
			UnitType: unitSpearman, //nolint:nolintlint,nolintlint,exhaustruct
			Label:    "Włócznik",
			MinLevel: 0,
			IconID:   0,
		},
		// Przycisk 3/5
		//nolint:exhaustruct
		{
			UnitType: unitSwordsman,
			Label:    "Miecznik",
			MinLevel: 0,
			IconID:   0,
		},

		// Przycisk 4/5
		// Przycisk 5/5
		{
			BuildingType: buildingPalisade,
			Label:        buildingDefs[buildingPalisade].Name,
			MinLevel:     0,
			IconID:       spriteBtnBuildPalisade,
			UnitType:     unitNone,
		},
	},
	buildingTemple: {
		// Przycisk 1/5
		//nolint:exhaustruct
		{
			UnitType: unitMage,
			Label:    "Mag",
			MinLevel: 0,
			IconID:   0,
		},
		// Przycisk 2/5
		//nolint:exhaustruct
		{
			UnitType: unitPriest,
			Label:    "Kapłan",
			MinLevel: 0,
			IconID:   0,
		},
		// Przycisk 3/5
		//nolint:exhaustruct
		{
			UnitType: unitPriestess,
			Label:    "Kapłanka",
			MinLevel: 0,
			IconID:   0,
		},
		// Przycisk 4/5
		// N/D
		// Przycisk 5/5
		// N/D
	},
	// Tutaj nowy budynek
	buildingAcademy: {
		// Przycisk 1/5
		//nolint:exhaustruct
		{
			UnitType: unitCrossbowman,
			Label:    "Kusznik",
			MinLevel: 0,
			IconID:   0,
		},
		// Przycisk 2/5
		//nolint:exhaustruct
		{
			UnitType: unitCommander,
			Label:    "Rycerz",
			MinLevel: 0,
			IconID:   0,
		},
		// Przycisk 3/5
		// Przycisk 4/5
		// N/D
		// Przycisk 5/5
		// N/D
	},
}

const maxGrassGrowthCounter uint16 = 300

const (
	// Koszty magicznych pocisków.
	magicThunderManaCost uint16 = 10
	magicFireManaCost    uint16 = 20
	magicGhostManaCost   uint16 = 20
)

const (
	// Stałe związane z duszkami zmarłych jednostek.
	freshlyDeadSpriteOffset = 40
	decayStartSpriteOffset  = 41
)

// == Rzeczy dla assets_db.go.
const (
	directionUp uint8 = iota
	directionUpRight
	directionRight
	directionDownRight
	directionDown
	directionDownLeft
	directionLeft
	directionUpLeft
	directionCount
)

const (
	frameIdle unitFrame = iota
	frameWalk1
	frameWalk2
	frameAttack1
	frameAttack2
	frameCount
)
