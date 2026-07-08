package main

// types.go

import rl "github.com/gen2brain/raylib-go/raylib"

// === Typy podstawowe i struktury ===

// inputState przechowuje, co jest na wejściu
// Obecnie tylko mysz.
type inputState struct {
	MousePosition              rl.Vector2 // Położenie myszy
	IsLeftMouseButtonDown      bool       // Czy lewy przycisk myszy jest wciśnięty
	IsLeftMouseButtonPressed   bool       // Czy lewy przycisk myszy został przyciśnięty
	IsLeftMouseButtonReleased  bool       // Czy lewy przycisk myszy został zwolniony
	IsRightMouseButtonDown     bool       // Czy prawy przycisk myszy jest wciśnięty
	IsRightMouseButtonPressed  bool       // Czy prawy przycisk myszy został przyciśnięty
	IsRightMouseButtonReleased bool       // Czy prawy przycisk myszy został zwolniony
	IsCtrlKeyDown              bool       // Czy ctrl jest przyciśnięty
	IsShiftKeyDown             bool       //
}

// mouseState określa stan kursora
type mouseState uint8

// unitState opisuje usposobienie jednostki.
type unitState int

// unitType opisuje rodzaj jednostki (drwal, krowa itd.)
type unitType uint8

type ObjectID uint

type (
	UnitID     ObjectID
	BuildingID ObjectID
)

// unitStats to definicja statystyk bazowych dla danego rodzaju jednostki.
type unitStats struct {
	Name string

	// Cechy bojowe
	MaxHP      uint16
	BaseDamage uint16
	BaseArmor  uint16

	// Zasięgi
	SightRange  uint8
	AttackRange uint8

	// Ekonomia i czas
	Cost      uint16
	MoveDelay uint16

	// Mana
	MaxMana       uint16
	BaseManaRegen uint16
}

const (
	stateIdle         unitState = iota // bezczynność
	stateMoving                        // ruch
	stateAttacking                     // napaść
	stateRepairing                     // naprawa
	stateBuilding                      // praca na budowie
	stateGrazing                       // wypasanie, dotyczy tylko krów
	stateCastingSpell                  // rzucanie czarów
	stateWaiting                       // oczekiwanie, jeżeli przytkana jest droga
	stateMilking                       // dojenie krowy, ew. oczekiwanie w kolejce na dojenie
)

// unit określa pojedynczą jednostkę podczas bitwy.
type unit struct {
	ID           UnitID      // Unikatowy numer jednostki
	Exists       bool        // Czy jednostka nie została jeszcze zabita
	X, Y         uint8       // Współrzędne jednostki
	Owner        PlayerID    // Kto jest właścicielem. colorRed gracz, inne SI
	Type         unitType    // Rodzaj jednostki (Drwal = unitAxeman itd.)
	HP           uint16      // Bieżący wskaźnik życia
	MaxHP        uint16      // Górna granica wskaźnika życia
	Command      commandType // Bieżący rozkaz (cmdMove, cmdAttack itd.)
	CurrentSpell spellID     // Jaki czar ma rzucić
	TargetID     ObjectID    // Identyfikator przedmiotu (jednostka bądź budynek) rozkazu
	TargetX      uint8       // Intencja: Współrzędna X celu (dla cmdMove. cmdAttack)
	TargetY      uint8       // Intencja: Współrzędna Y celu (dla cmdMove. cmdAttack)
	// ↓↓↓↓↓↓ NOWOŚĆ! TYLKO OD ZRZEZA!
	ApproachX uint8 // PODEJŚCIE: kafalek na którym trzeba stanąć, aby dobrać się do celu
	ApproachY uint8 // PODEJŚCIE: kafalek na którym trzeba stanąć, aby dobrać się do celu
	// ↑↑↑↑↑↑ Powinno zastąpić ułomne interactionTargetX/Y i uzdrowić rozkazy w units.go
	// Do obsługi atakowania drzew, być może zbędne - 20.04.2026
	// @todo:↓↓↓ wypierdziel jeśli ApproachXY siądą 02.07.2026
	interactionTargetX, interactionTargetY uint8

	Experience uint8     // Miara doświadczenia jednostki
	Level      uint8     // Poziom doświadczenia jednostki
	IsSelected bool      // Określa czy dana jednostka jest wybrana przez gracza
	BelongsTo  *building // Określa do którego budynku jest przywiązana dana jednostka
	IsInQueue  bool      // Wskaźnik, by wiedzieć, czy jednostka jest w kolejce
	// Dodawalibyśmy jednostkę do określonego budynku i w ten sposób śledzimy
	// kto do kogo przynależy i np. krowy gdzie idą dać się wydoić!
	AllowFriendlyFire bool // Wskaźnik, czy możemy napaść swoich
	// cechy w grze
	SightRange  uint8  // Pole widzenia
	AttackRange uint8  // Zasięg uderzenia
	Damage      uint16 // Obrażenia
	Armor       uint16 // Obrona
	MaxDelay    uint16 // Najwyższe opóźnienie jednostki. Nie rozumiem, co to znaczy
	Mana        uint16 // Miara many jednostki
	MaxMana     uint16 // Górna granica many
	ManaRegen   uint16 // Miara odnowy many
	// Effects     []activeEffect // Wykaz działających efektów
	// przechowywanie mleka przez jednostki
	Udder  uint8   // Dosłownie wymiona
	Wounds []wound // Wykaz świeżo zadanych ran

	// Magiczna tarcza, tylko dla unitPriestess
	hasMagicShield      bool // Wskazuje, czy jednostka jest chroniona magiczną tarczą
	MagicShieldCooldown uint // Wskazuje ile zostało do końca ochrony magiczną tarczą

	// Pola do rysowania jednostki
	AnimationType    string     // Rodzaj rysunku ("walk", "attack"). "idle" NIE ISTNIEJE!
	AnimationFrame   int        // Bieżąca klatka ruchu
	AnimationCounter int        // Pomocniczy licznik do zarządzania prędkością rysowania ruchu
	Direction        rl.Vector2 // Kierunek jednostki. (0,1) dla góry

	State          unitState // Obecne usposobienie jednostki (bezczynność, ruch, napaść itd.)
	Delay          uint16    // Bieżące opóźnienie jednostki przy chodzeniu
	AttackCooldown uint      // Opóźnienie przy ataku

	// Obieranie ścieżki
	Path                  []*pathNode  `json:"path"`        // Bieżąca ścieżka
	PathIndex             int          `json:"pathIndex"`   // Bieżący wskaźnik w ścieżce
	LastTargetX           uint8        `json:"lastTargetX"` // Ostatnia współrzędna X celu
	LastTargetY           uint8        `json:"lastTargetY"` // Ostatnia współrzędna Y celu
	PathfindingCooldown   int          // Opóźnienie dla ponownego obierania
	RetryAttempts         int          // Liczb prób obejścia przeszkody
	LastPathfindTime      int          // Ostatni czas kiedy przeliczono ścieżkę
	BlockedCounter        int          // Na ile zatkano jednostkę
	History               []rl.Vector2 // ostatnie pozycje jednostki
	LoopCount             int
	PathfindingRetryCount int
	LastPathIndex         int
	TicksNoProgress       int
	TicksNearTarget       int
	LastX, LastY          uint8 // ostatnia pozycja
	NoMoveTicks           int   // liczymy ticki bez zmiany pozycji
}

// Do trzymania wiedzy o odniesionych ranach.
type wound struct {
	Timer    uint    // Jak długo będzie widoczna
	OffsetX  float32 // Przesunięcie w poziomie
	OffsetY  float32 // Przesunięcie w pionie
	IsSevere bool    // Wskaźnik, czy animacja będzie dwuklatkowa
	Scale    float32 // Jak duża rana
	Rotation float32 // kąt pod którym narysujemy ranę
}

// Do trzymania wiedzy o zwłokach.
type corpse struct {
	X, Y         uint8    // gdzie poległa jednostka
	UnitType     unitType // jaka to jednostka, ważne dla 2 pierwszych faz rozkładu
	DecayTimer   uint     // do mierzenia kiedy podmienić teksturę
	Phase        uint8    // faza rozkładu
	SkeletonType uint8    // wskazuje na rodzaj szkieletu w ostatniej fazie rozkładu
	Rotation     float32  // kąd pod którym narysujemy zwłoki
	Alpha        uint8    // do zanikania kości
	Owner        PlayerID
}

// struktura do przedstawiania celu bojowego.
type combatTarget struct {
	Unit     *unit
	Building *building
	Tile     *tile
}

type economyAI struct {
	lastBarnChecked int
}

type militaryAI struct {
	// chwilowo puste
}

type buildingType uint8

type point struct {
	X, Y uint8
}
type buildingStats struct {
	Name          string
	Width, Height uint8  // Wymiary, 3x3, dla palisady 1x1
	Cost          uint16 // Cena w mleku
	MaxHP         uint16 // Górna granica wytrzymałości
	MaxFood       uint8  // Ile jednostek może być przypisanych do budynku
	BaseTextureID uint16 // ID tekstury, 1x1 całość, 3x3 początek ciągu
	IsPalisade    bool   // Wskaźnik, czy budowla jest palisadą, potrzebna do napraw wrogich palisad
}

type constructionProgress uint8

// building określa pojedynczy budynek podczas bitwy.
type building struct {
	ID                BuildingID   // Unikatowy numer budynku
	Exists            bool         // Czy budynek nie został jeszcze zniszczony
	Owner             PlayerID     // Kto jest właścicielem. colorRed gracz, inne SI
	Type              buildingType // Rodzaj budynku (obora = buildingBarn itd.)
	HP                uint16       // Bieżący wskaźnik wytrzymałości
	Armor             uint8        // Obrona budynku, zawsze równa 10
	AccumulatedDamage uint16       // Obrażenia przyjęte w jednej turze
	MaxHP             uint16       // Górna granica wskaźnika wytrzymałości
	Food              uint8        // Bieżący wskaźnik liczby jednostek w budynku
	MaxFood           uint8        // Górna granica wskaźnika liczby jednostek w budynku;
	AssignedUnits     []UnitID     // Identyfikatory jednostek przypisanych do budynku @todo: sprawdź, czy potrzebne
	// czemu nie wskaźnik?↑↑↑
	OccupiedTiles []point  // Współrzędne budynku
	MilkingQueue  []UnitID // Wycinek z ID krowami będącymi w kolejce do dojenia w danej oborze
	// Budowa
	IsUnderConstruction bool                 // Wskazuje, czy budowla jest w trakcie wznoszenia
	ConstructionPhase   constructionProgress // Wskaźnik zaawansowania 0 - początek budowy, 1 - połowa budowy, 2 - zakończona
}

// playerState przedstawia usposobienie gracza.
type playerState struct {
	PlayerID          PlayerID // Identyfikator gracza (colorRed, colorYellow)
	Milk              uint16   // Ilość mleka
	MaxMilk           uint16   // górna granica wskaźnika mleka
	CurrentPopulation uint8    // Do pilnowania liczby jednostek
	CurrentBuildings  uint8    // Do pilnowania liczby budynków
	// Brakuje określenia fazy, potrzebnej do SI
	// Brakuje określenia górnych granic liczby budynków tudzież jednostek
}

type (
	commandType uint8
	spellID     uint8
)

// command przechowuje rozkazy dla jednostki lub budynku.
type command struct {
	// ExecutorID - ID konkretnego obiektu, który MA WYKONAĆ rozkaz.
	// Jeśli Category=0, to jest to ID Budynku. Jeśli Category=1, to ID Jednostki.
	ExecutorID ObjectID

	// === CO? ===
	ActionType commandType // Co ma zostać zrobione (cmdMove, cmdAttack, cmdRepairStructure)

	// == JAKI CZAR @reminder: jeszcze nie wiem, czy tak powinno być!
	Spell spellID

	// === GDZIE / NA KIM? (Cel działania) ===

	// InteractionTargetID - ID obiektu, NA KTÓRYM wykonujemy akcję.
	// Np. ID Budynku do naprawy lub ID Wroga do ataku.
	// WAŻNE: w przyciskach UI to pole jest zazwyczaj puste (0)! Wypełnia je dopiero Input po kliknięciu w mapę.
	InteractionTargetID ObjectID

	TargetX uint8 // Współrzędna x kliknięcia w mapę (dla ruchu/ataku obszarowego)
	TargetY uint8 // Współrzędna y kliknięcia w mapę

	// === INNE ===
	CreateType   uint8 // Rodzaj obiektu, który stworzymy
	FriendlyFire bool  // Określnik, czy bratobójczy napad
}

// aiState przechowuje usposobienie SI.
type aiState struct {
	// CurrentPhase        int // Bieżące usposobienie SI
	// BuildingIndex       int // Identyfikator budynku jako podmiotu podejmowania decyzji
	// AttackBuildingIndex int // Identyfikator budynku do napadnięcia
	GatherPointX uint8 // Współrzędna x miejsca zbornego
	GatherPointY uint8 // Współrzędna y miejsca zbornego
	PastureX     uint8 // Współrzędna x pastwiska
	PastureY     uint8 // Współrzędna y pastwiska
	// DistanceToEnemy     int // Odległość do wroga
	MilkGenerationRate uint16 // Przyrost mleka
	ActionDelay        uint   // Opóźnienie przed kolejnym działaniem
	// AttackWaveIndex     int // Numer napadu
	// TargetEnemyX        int // Współrzędna x przedmiotu napadu SI
	// TargetEnemyY        int // Współrzędna y przedmiotu napadu SI

	// Coś tam do SI.
	CurrentGoals    []aiGoal
	EconomyManager  *economyAI
	MilitaryManager *militaryAI
}

type aiGoal struct {
	Type     aiGoalType
	Priority float32
	Target   any // budynek, jednostka, pozycja
	Progress float32
}

type aiGoalType int

const (
	goalProduceCows aiGoalType = iota
	goalBuildArmy
	goalAttackEnemy
	goalDefendBase
	goalExpandEconomy
)

// message opisuje wiadomości.
type message struct {
	Text      string // Treść wiadomości
	Duration  int    // Licznik czasu wyświetlania wiadomości
	SoundID   int    // Identyfikator dźwięku do otworzenia
	PositionX int    // Współrzędna x wiadomości
	PositionY int    // Współrzędna y wiadomości
}

// Opisuje bieżący przedmiot zaznaczenia.
type selectionState struct {
	OwnerID PlayerID // Identyfikator zaznaczonego przedmiotu (colorRed, colorYellow itd.)
	IsUnit  bool     // Określa, czy zaznaczony przedmiot jest jednostką (prawda), czy budynkiem (fałsz)
	// Nie wiem, co ze zgliszczami
	UnitID     UnitID     // Identyfikator przedmiotu zaznaczenia jeżeli jednostka
	BuildingID BuildingID // Identyfikator przedmiotu zaznaczenia jeżeli budynek
}

// controlGroupUnit określa jednostkę w zarządzanym zespole - 02.01.2026 ciekawe co miałem na myśli pisząc to.
// @todo: czy to powód przekazywania rozkazu tylko jednej jednostce w całej drużynie?
type controlGroupUnit struct {
	UnitID UnitID // Identyfikator jednostki
}

// controlGroup określa zarządzany zespół (przypisany do cyfr od 1 do 0?)
type controlGroup struct {
	Units []controlGroupUnit
}

// visibilityState określa stan mgły wojny dla danego kafelka.
type visibilityState int

const (
	visibilityUnexplored visibilityState = iota // niezbadany
	visibilityFog                               // niewidziany
	visibilityVisible                           // widziany
)

// przechowuje kompletną wiedzę o polu (x,y).
type tile struct {
	X, Y      uint8  // współrznędne kafelka
	TextureID uint16 // ID grafiki np. trawa, droga

	// Trawa
	IsGrazed        bool
	GrazedCounter   uint8
	GrazedOverlayID uint8

	// Płomień
	IsBurning     bool
	BurnElapsed   uint16
	BurnOverlayID uint16
	IsBurnt       bool

	// Ścinanie drzew
	// @reminder: być może będzie trzeba przerobić treeCuts z liczby uderzeń na sumę obrażeń.
	treeCuts uint8 // Ile uderzeń drwala do ścięcia drzewa.

	/*
		Śledzimy stan drzewa

		noTree - nie ma drzewa

		treeStraight - stoi normalnie drzewo

		treeLeaning - drzewo zaczyna się przechylać

		treeFalling - drzewo jest przechylone

		treeImpact - drzewo właśnie uderzyło w sąsiedni kafelek, wywołując efekty

		treeFell - drzewo leży i nie zagradza już drogi
	*/
	treeState treeState

	// Popiół
	hasAsh          bool
	AshIntensity    float32
	AshAge          uint16
	AshProcessState uint8
	CurrentAshAlpha float32

	// duszenie przez pocisk unitMage
	GhostEffect        bool  // do rysowania efektu ducha po trafieniu
	GhostEffectCounter uint8 // do odliczania czasu od pojawienia się efektu`
	GhostDamage        uint16

	// --- WŁAŚCIWOŚCI FIZYCZNE ---
	IsWalkable   bool    // przechodniość kafelka
	MovementCost float64 // drożność kafelka. Dla A*.

	// --- ZAWARTOŚĆ ---
	Unit     *unit
	Building *building

	// --- MGŁA WOJNY ---
	Visibility visibilityState
}

type boardData struct {
	Tiles [boardMaxX][boardMaxY]tile
}

type dragContext struct {
	IsActive   bool       // Czy mamy punkt początkowy
	AnchorPos  rl.Vector2 // Punkt początkowy
	CurrentPos rl.Vector2 // Gdzie jest teraz mysz
}

type PlayerID uint8

// battleState przechowuje opis bitwy
// @todo: czemu nazywamy to „battle” skoro powinien przechowywać też dane o całej wyprawie?
type battleState struct {
	// === GRACZE I ZASOBY === //
	PlayerID         PlayerID     // Identyfikator gracza (domyślnie colorRed)
	AIPlayerID       PlayerID     // Identyfikator SI
	HumanPlayerState *playerState // Opis usposobienia gracza
	AIEnemyState     *playerState // Opis usposobienia wroga
	AI               aiState      // Pamięć i usposobienie SI
	CurrentCommands  [2]command   // Rozkazy SI @todo: być może CommandQueue
	CampaignData     campaignData // Zawiera rzeczy związane wyprawą i obecną bitwą

	// === PRZEDMIOTY I JEDNOSTKI NA MAPIE ===
	Units              []*unit                // Żyjące jednostki
	Buildings          []*building            // Działające budynki
	NextUnitID         UnitID                 // Licznik do tworzenia nowych identyfikatorów dla jednostek
	NextBuildingID     BuildingID             // Licznik do tworzenia nowych identyfikatorów dla budynków
	NextUniqueObjectID ObjectID               // Ogólny licznik identyfikatorów dla wszystkich.
	Board              *boardData             // Wszystko co związane z przechowywaniem współrzędnych na planszy
	Projectiles        []*projectile          // Pociski
	HealingShrines     []point                // Wykaz współrzędnych miejsc leczenia
	RenderUnitRows     [boardMaxY][]*unit     // Służy do wskazywania jednostek wg współrzędnej y i rysowania na planszy
	RenderBuildingRows [boardMaxY][]*building // Służy do wskazywania budynków wg współrzędnych na potrzeby rysowania
	CorpsesList        []corpse               // Służy do wskazywania zwłok wg współrzędnych na potrzeby rysowania
	BurningTilesList   []*tile                // Służy do wskazania kafelków, które zostały podpalone
	FallingTreesList   []*tile                // Służy do wskazania kafelków z upadającymi drzewami
	GhostsList         []*tile                // Służy do wskazywania jednostek z efektem ducha
	// === UI I INTERAKCJA ===
	GameCamera              rl.Camera2D      // Kamera widoku gry
	CurrentSelection        selectionState   // Bieżące zaznaczenie
	ControlGroups           [10]controlGroup // Zarządzane zespoły
	MouseState              mouseState       // Tryb kursora
	DragContext             dragContext      // Przeciąganie
	Map                     bool             // Czy mapa jest widoczna @todo: czy w ogóle tego używam?!
	IsMapDragging           bool
	MapInitialClickPos      rl.Vector2
	CameraTargetOnDragStart rl.Vector2
	// PendingBuildingType     buildingType // Rodzaj budynku do wybudowania - 28.06.2026 raczej się pozbędę tego
	PendingCommand *command // @reminder: jeśli się uda, to ujednolici rozkazy
	UI             uiAssets

	// === Przymioty bitwy ===
	CurrentLevel          uint8           // Numer bieżącego poziomu
	DifficultyLevel       difficultyLevel // Poziom trudności
	GameSpeed             uint16          // Szybkość gry (bo skalar, prędkość jest wektorem)
	IsSinglePlayerGame    bool            // Wskaźnik, czy gra jednoosobowa
	CheatsEnabled         bool            // Wskaźnik, czy oszustwa są dozwolone
	CheatSequenceProgress int             // Do sprawdzania, czy wpisujemy prawidłowe oszustwa
	BattleOutcome         uint8           // Wynik gry ( 0 = trwa, endKillAll, endRescue itd.)
	QuitLevel             bool            // Wskaźnik, czy poziom ma zostać zakończony
	CurrentMessage        message         // Wiadomość wyświetlana na ekranie @todo: częściowo wdrożone, nigdy nie użyte

	// === LICZNIKI I URUCHAMIANIE OBRAZÓW
	GrassGrowthCycle    uint16 // Licznik odrastania trawy oraz innych ogólnych uruchomień
	WaterAnimationFrame uint16 // Wskaźnik falowania wody
	FireAnimationFrame  uint16 // Wskaźnik uruchomienia ognia
	GlobalFrameCounter  uint16 // Ogólny licznik klatek do różnych uruchomień i logiki
	// TODO: może uda się zastąpić czymś z raylib albo Go, bo to raczej przestarzałe

	EnemyCache                map[uint]*enemyCacheEntry
	enemyCacheUpdateTick      uint16
	pathfindingUnitsThisTick  uint16
	enemyCacheUpdatesThisTick uint16
	PathfindingBudget         int
}

type enemyCacheEntry struct {
	NearestEnemyUnit     *unit
	NearestEnemyBuilding *building
	Found                bool
	LastUpdateTick       uint16
	CacheValidFor        uint16
}

// campaignData przechowuje przymioty bieżącej wyprawy
// TODO: czemu nazywamy to „campaign” skoro dotyczy tylko bitwy?
type campaignData struct {
	DecisionType uint8 // NIE JESTEM PEWIEN: chyba rodzaj strategii SI (boardVillage, boardBattleDyn itd.)
	//  TODO: Ogarnij, czy to nie powinno być w battleState
	GeneratorActive bool // Wskaźnik, czy wytwórstwo jednostek działa
	//  TODO: Ogarnij, czy to nie powinno być w battleState
	EndCondition        uint8  // Rodzaj zadania (endKillAll itd.) TODO: Ogarnij, czy to nie powinno być w battleState
	TargetType          uint8  // Rodzaj celu do uratowania TODO: Ogarnij, czy to nie powinno być w battleState
	VictoryPointX       uint8  // Współrzędna x miejsca zbornego  TODO: Ogarnij, czy to nie powinno być w battleState
	VictoryPointY       uint8  // Współrzędna y miejsca zbornego TODO: Ogarnij, czy to nie powinno być w battleState
	RescueTargetX       uint8  // Współrzędna x jednostki do uratowania TODO: Ogarnij, czy to nie powinno być w battleState
	RescueTargetY       uint8  // Współrzędna y jednostki do uratowania TODO: Ogarnij, czy to nie powinno być w battleState
	TransformationSiteX uint8  // Współrzędna x miejsca przemiany TODO: Ogarnij, czy to nie powinno być w battleState
	TransformationSiteY uint8  // Współrzędna y miejsca przemiany TODO: Ogarnij, czy to nie powinno być w battleState
	TransformationType  int    // Rodzaj przemiany TODO: Ogarnij, czy to nie powinno być w battleState
	GeneratorTimer      int    // Licznik wytwarzacza (jednostek?)
	LevelsMilkLimit     uint16 // Wskaźnik Górnej granicy mleka dla wybranej misji
	Name                string // Imię postaci lub nazwa poziomu
	Password            int    // Hasło TODO: powinno się to usunąć, bo nikt już nie zabezpiecza gry hasłem z instrukcji!
	PasswordNumber      int    // Numer hasła
	//  TODO: powinno się to usunąć, bo nikt już nie zabezpiecza gry hasłem z instrukcji!
	LevelNumber    uint8 // Numer poziomu
	CurrentEventID int   // Licznik zdarzeń (co to w ogóle jest?)
	NextLevel      uint8 // Numer następnego poziomu (czemu tak, skoro niektóre udrażniają wiele poziomów)
}

// === UI, PRZYCISKI I INNE TAKIE

// Pozwala opisać każdy przycisk wytwórczy.
type buildingAction struct {
	UnitType     unitType     // jaki rodzaj jednostki można stworzyć unitCow
	BuildingType buildingType // Jaki rodzaj budynku zamierzamy wybudować
	Label        string       // @todo: TYMCZASOWE ZAMIAST TEKSTURY; być może zostawimy dla dymków
	MinLevel     uint8
	IconID       uint16
}

// === RZECZY ZWIĄZANE Z ŁADOWANIEM MAP ===

// JSON - ładowanie mapy.
type jsonLevelLoader struct {
	drivePath string
}

// Struktury JSON.
type jsonLevel struct {
	Metadata         jsonLevelMetadata    `json:"metadata"`
	AISettings       jsonAISettings       `json:"aiSettings"`
	SpecialLocations jsonSpecialLocations `json:"specialLocations"`
	Terrain          jsonTerrainData      `json:"terrain"`
	Buildings        []jsonBuildingData   `json:"buildings"`
	Units            []jsonUnitData       `json:"units"`
}

type jsonLevelMetadata struct {
	Name         string    `json:"name"`
	LevelNumber  uint8     `json:"levelNumber"`
	DecisionType uint8     `json:"decisionType"`
	EndType      uint8     `json:"endType"`
	TargetType   uint8     `json:"targetType"`
	MaxMilk      uint16    `json:"maxMilk"`
	Generator    bool      `json:"generator"`
	NextLevel    uint8     `json:"nextLevel"`
	StartPos     jsonPoint `json:"startPosition"`
}

type jsonAISettings struct {
	GatherPoint jsonPoint `json:"gatherPoint"`
	Pasture     jsonPoint `json:"pasture"`
}

type jsonSpecialLocations struct {
	HealingShrine       *jsonPoint `json:"healingShrine,omitempty"`
	TransformationPoint *jsonPoint `json:"transformationPoint,omitempty"`
	VictoryPoint        *jsonPoint `json:"victoryPoint,omitempty"`
	RescueTarget        *jsonPoint `json:"rescueTarget,omitempty"`
}

type jsonTerrainData struct {
	Width  int        `json:"width"`
	Height int        `json:"height"`
	Tiles  [][]string `json:"tiles"`
}

type jsonBuildingData struct {
	Type     string    `json:"type"`
	Owner    string    `json:"owner"`
	Position jsonPoint `json:"position"`
	HP       int       `json:"hp,omitempty"`
}

type jsonUnitData struct {
	Type     string    `json:"type"`
	Owner    string    `json:"owner"`
	Position jsonPoint `json:"position"`
}

type jsonPoint struct {
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}

// === NAKŁADKA DLA UŻYTKOWNIKA

// uiAssets przechowuje tekstury i przyciski związane z nakładką.
type uiAssets struct {
	// Zdefiniowanie 5 fizycznych przycisków na nakładce
	ActionButtons [5]rl.Rectangle

	// Miejsce w którym przechowujemy informację o przyciskach
	// w zależności od zaznaczonego obiektu
	CurrentActions [5]uiAction
}

// opisuje pojedyńcze działanie dostępne pod przyciskiem z nakładki, pomost pomiędzy rysowaniem a logiką.
type uiAction struct {
	IsActive bool       // widoczność
	Label    string     // @todo: tymczasowy napis
	IconID   uint16     // wskaźnik tekstury docelowej
	Cmd      command    // rozkaz do wykonania przez budynek lub jednostkę
	State    mouseState // tryb kursora (budowa, naprawa, celowanie)
}

type button struct {
	Rectangle  rl.Rectangle // Gdzie rysujemy przycisk
	OnClick    func()       // Co się stanie, gdy naciśniemy
	DebugLabel string       // @reminder tymczasowe rozwiązanie, aby móc debugować!
}

// Nakładka dla budynków.
type bounds struct {
	X, Y              int32 // lewy górny róg w pikselach
	Width, Height     int32
	WidthPx, HeightPx float32
}

type unitDirection uint8

// ============ INTERFACE
type objectResolver interface {
	getObjectByID(id ObjectID) (*unit, *building)
}

type unitResolver interface {
	getUnitByID(id UnitID) (*unit, bool)
}

type buildingResolver interface {
	getBuildingByID(id BuildingID) (*building, bool)
}
