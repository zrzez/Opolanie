package main

// buttons.go

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func setupMainMenuButtons(ps *programState, bState *battleState) {
	ps.ActiveMenuButtons = []button{
		{
			Rectangle:  rl.NewRectangle(215, 60, 222, 36),
			DebugLabel: "Nowa gra",
			OnClick: func() {
				ps.CurrentState = newCampaignMenuScreen
				setupCampaignMenuButtons(ps, bState)
				startNewCampaign(ps) // @todo: pytanie, czy musimy to tak trzymać, bo wydaje się, że to zbędne
				// można od razy wywołać drawSelectCampaignMenuLegacy
			},
		},
		{
			Rectangle:  rl.NewRectangle(215, 114, 222, 36),
			DebugLabel: "Wczytaj grę",
			OnClick: func() {
				// @todo do wdrożenia!
				log.Println("@todo wczytaj grę")
			},
		},
		{
			Rectangle:  rl.NewRectangle(215, 285, 222, 36),
			DebugLabel: "Wyjście",
			OnClick:    func() { rl.CloseWindow() },
		},
	}
}

func handleMenuInput(ps *programState) {
	screenMouse := rl.GetMousePosition()
	virtualMouse := screenToVirtualCoords(ps, screenMouse)

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		for _, btn := range ps.ActiveMenuButtons {
			if rl.CheckCollisionPointRec(virtualMouse, btn.Rectangle) {
				log.Printf("MENU: Kliknięto przycisk: %s", btn.DebugLabel)

				if btn.OnClick != nil {
					btn.OnClick()
				}

				return
			}
		}
	}
}

func setupCampaignMenuButtons(ps *programState, bState *battleState) {
	// Pomagier do szybkiego ustawiania przycisków w tym menu
	addCampaignButton := func(y float32, label string, campaignFunc func()) button {
		return button{
			Rectangle:  rl.NewRectangle(90, y, 331, 46),
			DebugLabel: label,
			OnClick:    campaignFunc,
		}
	}

	addDifficultyButton := func(y float32, label string, selectedLevel difficultyLevel) button {
		return button{
			Rectangle:  rl.NewRectangle(505, y, 60, 70),
			DebugLabel: label,
			OnClick:    func() { ps.SelectedDifficulty = selectedLevel },
		}
	}

	ps.ActiveMenuButtons = []button{
		addCampaignButton(22, "Powrót Mirka", func() {
			startFirstCampaign(ps, bState)
		}),
		addCampaignButton(82, "Przyjaciele", func() { log.Println("Kampania 2, @todo") }),
		addCampaignButton(142, "Porwanie", func() { log.Println("Kampania 3, @todo") }),
		addCampaignButton(203, "Wojna Magów", func() { log.Println("Kampania 4, @todo") }),
		addCampaignButton(363, "Wschodnia Pożoga", func() { log.Println("Kampania 5, @todo") }),
		addCampaignButton(323, "Południe w ogniu", func() { log.Println("Kampania 6, @todo") }),

		addDifficultyButton(40, "Łatwy", difficultyEasy),
		addDifficultyButton(125, "Normalny", difficultyNormal),
		addDifficultyButton(210, "Trudny", difficultyHard),
		{
			Rectangle:  rl.NewRectangle(468, 322, 130, 45),
			DebugLabel: "Menu",
			OnClick: func() {
				ps.CurrentState = mainMenuScreen
				setupMainMenuButtons(ps, bState)
			},
		},
	}
}

// Działanie przycisków w menu
func startFirstCampaign(ps *programState, bState *battleState) {
	log.Println("Naciśnięto pierwszą kampanię. Rozpoczynam sekwencję startową.")

	// KROK 1: Ustalenie parametrów (Sztywno dla 1. kampanii)
	startLevel := firstCampaignFirstLevel // (15)
	enemyColor := getEnemyColor(startLevel)

	log.Printf("START KAMPANII: Poziom %d. Przeciwnik: %d", startLevel, enemyColor)

	// KROK 2: Reset Stanu Bitwy
	bState.IsSinglePlayerGame = true
	bState.DifficultyLevel = ps.SelectedDifficulty

	bState.HumanPlayerState.init(colorRed, 0)
	bState.AIEnemyState.init(enemyColor, 0)

	bState.PlayerID = colorRed
	bState.AIPlayerID = enemyColor

	// Zerowanie liczników
	bState.CheatsEnabled = false
	bState.CheatSequenceProgress = 0
	bState.GrassGrowthCycle = 0
	bState.WaterAnimationFrame = 0
	bState.FireAnimationFrame = 0
	bState.GlobalFrameCounter = 0
	bState.NextUnitID = 1

	// Nowe, czyste wycinki
	bState.Units = make([]*unit, 0)
	bState.Buildings = make([]*building, 0)
	bState.Projectiles = make([]*projectile, 0)
	bState.Board = &boardData{}

	// KROK 3: Grafika (Assets)
	ps.Assets.unloadBattleAssets()

	activePlayers := []PlayerID{colorRed, enemyColor}
	err := ps.Assets.loadBattleAssets(activePlayers)
	if err != nil {
		log.Printf("KRYTYCZNY BŁĄD GRAFIKI: %v", err)

		ps.CurrentState = mainMenuScreen
		setupMainMenuButtons(ps, bState)

		return
	}

	// KROK 4: Ładowanie Poziomu (JSON)
	loader := newLevelLoader("./")

	err = loader.initBattle(startLevel, bState)
	if err != nil {
		log.Printf("KRYTYCZNY BŁĄD BITWY: Nie udało się załadować poziomu %d: %v", startLevel, err)

		ps.CurrentState = mainMenuScreen
		setupMainMenuButtons(ps, bState)

		return
	}

	// KROK 5: Wstępniaki

	// playIntroSegment(ps, "s003_frames", "i003.wav")

	// KROK 6: Start
	log.Println("Stan programu zmieniony na GameScreen. Rozpoczyna się bitwa.")
	ps.changeState(gameScreen, bState)
}
