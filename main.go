package main

// main.go

import (
	"fmt"
	"log"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var rng *rand.Rand

// @todo: sprawdź później, czy nadal jest to potrzebne przy odnajdywaniu drogi przez jednostki
func init() {
	source := rand.NewSource(time.Now().UnixNano())
	rng = rand.New(source)
}

// Stałe potrzebne do rysowania urojonego płótna.
const (
	virtualScreenWidth  uint16 = 760
	virtualScreenHeight uint16 = 400
)

// Stałe potrzebne do rysowania widoku bitwy
const (
	gameViewVirtualWidth = 640                  // Szerokość widoku bitwy
	uiPanelVirtualWidth  = 120                  // Szerokość nakładki
	uiPanelVirtualStartX = gameViewVirtualWidth // Początek nakładki
	scrollZoneXThreshold = 11
	scrollZoneYThreshold = 8
)

// Rodzaj ekranów, które mogą wystąpić podczas działania programu
type screenState uint8

const (
	gameScreen            screenState = iota // Eran bitwy
	mainMenuScreen                           // Ekran menu głównego
	newCampaignMenuScreen                    // Ekran wyboru wyprawy
	// SELECT_MISSION_MENU_SCREEN                    // Ekran wyboru misji w oryginalnej wyprawie
	// SAVE_MENU_SCREEN                              // Ekran zapisu gry
	// LOAD_MENU_SCREEN                              // Ekran wczytywania gry
	// SETTINGS_MENU_SCREEN                          // Ekran ustawień gry
	// INTRO_SCREEN                                  // Ekran odtwarzania wstępniaków
)

// Określnik poziomu trudności rozgrywki
type difficultyLevel int

// Poziomy trudności rozgrywki
const (
	difficultyEasy   = iota // Łatwy
	difficultyNormal        // Zwyczajny
	difficultyHard          // Trudny
)

// Określenie stanu programu
type programState struct {
	// === OGÓLNE ===
	BaseAssetsPath string      // Ścieżka dostępu do zasobów gry
	ShouldClose    bool        // Wskaźnik, czy należy zamknąć grę
	CurrentState   screenState // Określenie, który ekran jest obecny używany
	// programContext                     // Zbiór rzeczy określających kontekst programu
	SelectedDifficulty difficultyLevel // Wskaźnik wybranego poziomu trudności

	// === DŹWIĘK ===
	MusicEnabled       bool     // Wskaźnik, czy włączono odtwarzanie muzyki
	CurrentMusic       rl.Music // Bieżąca muzyka
	CurrentTruckNumber int      // Numer ścieżki odtwarzanej muzyki
	SFXVolume          float32  // Głośność dźwięków
	MusicVolume        float32  // Głośność muzyki
	GameMusicFiles     []string // Pliki z muzyką TODO: ogarnij czemu tak jest to zrobione

	// === EKRAN ===
	VirtualWidth  float32 // Całkowita szerokość urojonego płótna
	VirtualHeight float32 // Całkowita wysokość urojonego płótna. Stała 400.0
	GameViewWidth float32 // Szerokość widoku gry (VirtualWidth - UI_WIDTH)
	// -- Zmienne robocze, renderer
	RenderSrcRect  rl.Rectangle // Do wycinania z tekstury źródłowej
	RenderDestRect rl.Rectangle // Do wklejania na ekran
	RenderOrigin   rl.Vector2   // Punk obrotu, zazwyczaj mam 0,0

	ShowMainMenu       bool               // Wskaźnik, czy pokazać menu główne
	ActualScreenWidth  int32              // Bieżąca szerokość ekranu
	ActualScreenHeight int32              // Bieżąca wysokość ekranu
	ScreenTarget       rl.RenderTexture2D // NIE PAMIĘTAM @todo: ogarnij, co to jest
	MenuCamera         rl.Camera2D        // Kamera wykorzystywana w menu głównym
	GameCamera         rl.Camera2D        // Kamera wykorzystywana w bitwie

	Assets *assetManager // zbiór atlasów z pliku .dat; Przechodzę na ten układ korzystania z tekstur
	// NOWE PODEJŚCIE DO PRZYCISKÓW 31.12.2025 próba ostatecznego rozwiązania sprawy
	ActiveMenuButtons []button
}

func newProgramState() *programState {
	// musicFiles := []string{} // @todo: samoczynne ładowanie z BaseAssetsPath

	ps := &programState{
		CurrentState:       mainMenuScreen,
		ActualScreenWidth:  1280,
		ActualScreenHeight: 720,
		VirtualHeight:      400.0,
		VirtualWidth:       640.0 + 120.0,
		ShouldClose:        false,
		MusicEnabled:       true,
		SFXVolume:          0.03,
		MusicVolume:        0.03,
		// GameMusicFiles:     musicFiles,
		ShowMainMenu:       true,
		SelectedDifficulty: 1,
		RenderSrcRect:      rl.NewRectangle(0, 0, 0, 0),
		RenderDestRect:     rl.NewRectangle(0, 0, 0, 0),
		RenderOrigin:       rl.NewVector2(0, 0),
	}

	return ps
}

// Przelicza wymiary urojonego płótna na podstawie fizycznych wymiarów okna
// Zachowuje proporcje
func (ps *programState) recalculateVirtualResolution() {
	fmt.Print("DEBUG: zmiana wielkości okna. To nie powinno być ciągle wywoływane")
	if ps.ActualScreenHeight == 0 {
		ps.ActualScreenHeight = 1
	}

	if ps.CurrentState == gameScreen {

		targetScale := float32(ps.ActualScreenHeight) / 400.0

		// Ustawiamy wymiary wirtualne
		ps.VirtualHeight = 400.0
		ps.VirtualWidth = float32(ps.ActualScreenWidth) / targetScale

		// Zabezpieczenie: minimalna szerokość panelu
		if ps.VirtualWidth < 120.0 {
			ps.VirtualWidth = 120.0
		}

		// GameViewWidth to szerokość gry BEZ panelu UI
		ps.GameViewWidth = ps.VirtualWidth - 120.0
		if ps.GameViewWidth < 0 {
			ps.GameViewWidth = 0
		}
	} else {
		// Dla menu
		ps.VirtualWidth = 640.0
		ps.VirtualHeight = 400.0
		ps.GameViewWidth = 640.0
	}

	// Aktualizacja tekstury renderowania
	if ps.ScreenTarget.ID != 0 {
		rl.UnloadRenderTexture(ps.ScreenTarget)
	}
	// Tworzymy teksturę o nowych wymiarach
	ps.ScreenTarget = rl.LoadRenderTexture(int32(ps.VirtualWidth), int32(ps.VirtualHeight))

	rl.SetTextureFilter(ps.ScreenTarget.Texture, rl.FilterPoint)

	// Reset kamery
	setupMenuCamera(ps)
}

// Odtwarza animację bez dźwięku
func playFrameAnimation(framesDirPath string, targetFps int) bool {
	if _, err := os.Stat(framesDirPath); os.IsNotExist(err) {
		log.Printf("OSTRZEŻENIE: Katalog animacji nie istnieje: %s\n", framesDirPath)
		return true // błąd nie jest krytyczny, można działać dalej
		// brak folderu z animacjami nie powinien zawieszać/zamykać gry
	}

	files, err := os.ReadDir(framesDirPath)
	if err != nil {
		log.Printf("OSTRZEŻENIE: Nie można odczytać katalogu animacji: %s\n", framesDirPath)
		return true // błąd nie jest krytyczny, można działać dalej
		// brak klatek z animacji nie powinien zawieszać/zamykać gry
	}

	var frameFiles []string
	for _, file := range files {
		if !file.IsDir() { // tylko pliki, które nie są folderami
			frameFiles = append(frameFiles, filepath.Join(framesDirPath, file.Name()))
		}
	}

	if len(frameFiles) == 0 {
		log.Printf("OSTRZEŻENIE: Brak klatek w katalogu: %s\n", framesDirPath)
		return true // brak animacji nie jest krytycznym błędem, można działać dalej
	}

	// TODO: później sprawdź, czy układanie jest w ogóle potrzebne skoro są odpowiednio nazwane
	sort.Strings(frameFiles) // klatki powinny być w odpowiedniej kolejności

	textures := make([]rl.Texture2D, len(frameFiles))
	for i, filePath := range frameFiles {
		tex := rl.LoadTexture(filePath)
		if tex.ID == 0 { // sprawdzenie, czy tekstura została załadowana poprawnie
			log.Printf("OSTRZEŻENIE: Nie udało się załadować klatki: %s\n", filePath)
			for j := 0; j < i; j++ {
				// czyścimy załadowane tekstury
				rl.UnloadTexture(textures[j])
			}
			return true // brak animacji nie jest krytycznym błędem, można działać dalej
		}
		textures[i] = tex
	}
	defer func() {
		for _, tex := range textures {
			if tex.ID != 0 {
				rl.UnloadTexture(tex)
			}
		}
	}()
	animShortName := filepath.Base(framesDirPath)
	log.Printf("INFO: Rozpoczęcie odtwarzania animacji klatkowej: %s (%d klatek @ %d FPS)\n", animShortName,
		len(textures), targetFps)

	frameDuration := 1.0 / float32(targetFps)
	currentFrameIndex := 0
	frameTimer := float32(0.0)

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeySpace) {
			log.Println("INFO: Animacja pominięta przez użytkownika.")
			return false // pomijamy
		}

		frameTimer += rl.GetFrameTime()
		if frameTimer >= frameDuration {
			frameTimer -= frameDuration
			currentFrameIndex++
			if currentFrameIndex >= len(textures) {
				break
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black) // tło animacji

		// wyśrodkowanie, być może się nie przyda
		currentTexture := textures[currentFrameIndex]
		posX := (rl.GetScreenWidth() - int(currentTexture.Width)) / 2
		posY := (rl.GetScreenHeight() - int(currentTexture.Height)) / 2
		rl.DrawTexture(currentTexture, int32(posX), int32(posY), rl.White)

		rl.DrawText(fmt.Sprintf("Animacja: %s [%d/%d]", animShortName, currentFrameIndex+1, len(textures)),
			10, 10, 10, rl.LightGray)
		rl.DrawText("Naciśnij SPACJĘ, aby pominąć wstęp", 10, int32(rl.GetScreenHeight()-20), 10, rl.LightGray)
		rl.EndDrawing()
	}

	return true
}

// Odtwarza dźwięk dla wstępniaka
// prawdopodobnie tylko .wav wchodzi w grę
func playSoundEffect(ps *programState, soundPath string) rl.Sound {
	if _, err := os.Stat(soundPath); os.IsNotExist(err) {
		log.Printf("OSTRZEŻENIE: Plik ze ścieżką dźwiękową nie istnieje: %s\n", soundPath)
		return rl.Sound{}
	}

	sound := rl.LoadSound(soundPath)
	if !rl.IsSoundValid(sound) {
		log.Printf("OSTRZEŻENIE: Nie udało się załadować dźwięku: %s\n", soundPath)
		return rl.Sound{}
	}

	rl.SetSoundVolume(sound, ps.SFXVolume)
	rl.PlaySound(sound)
	log.Printf("INFO: Odtwarzanie dźwięku: %s\n", filepath.Base(soundPath))
	return sound
}

// Zatrzymuje i zwalnia dźwięk
func stopAndUnloadSound(sound rl.Sound) {
	if rl.IsSoundValid(sound) {
		log.Printf("INFO: zatrzymywanie dźwięku i zwalnianie kanału z dźwiękiem.\n")
		rl.StopSound(sound)
		rl.UnloadSound(sound)
	}
}

// @todo: co to jest właściwie? Film przy uruchomieniu gry, czy już dla poszczególnych bitew?
func playIntroSegment(ps *programState, framesDirName, soundFileName string) bool {
	dataPath := filepath.Join(ps.BaseAssetsPath, "data")
	framesPath := filepath.Join(dataPath, framesDirName)

	var currentSound rl.Sound
	if soundFileName != "" {
		soundPath := filepath.Join(dataPath, soundFileName)
		currentSound = playSoundEffect(ps, soundPath)
	}

	const defaultIntroFPS = 60
	continueSequence := playFrameAnimation(framesPath, defaultIntroFPS)
	if !continueSequence {
		log.Println("INFO: Wstęp przerwany przez gracza.")
		stopAndUnloadSound(currentSound)
		return false
	}

	stopAndUnloadSound(currentSound)

	return continueSequence
}

// Odtwarza cały wstęp.
// Zwraca false, jeśli użytkownik je pominie
func playIntro(ps *programState) bool {
	log.Println("INFO: Rozpoczęcie sekwencji intro")

	if ps.BaseAssetsPath == "" {
		log.Println("OSTRZEŻENIE: programState.BaseAssetsPath nie jest ustawione! Nie można odtworzyć wstępu")
		return true
	}

	// USER
	playIntroSegment(ps, "s000_frames", "")

	// DRAHMA
	playIntroSegment(ps, "s001_frames", "i001.wav")

	// WSTĘP
	playIntroSegment(ps, "s002_frames", "i002.wav")

	log.Println("INFO: Wstęp odtworzony pomyślnie")
	return true
}

// Powinien zainicjować odtwarzacz pieśni
func initMusicPlayer(ps *programState) {
	if !rl.IsAudioDeviceReady() {
		log.Println("OSTRZEŻENIE: urządzenie audio nie jest gotowe.")
	}

	if ps.CurrentMusic.Stream.Buffer != nil {
		rl.SetMusicVolume(ps.CurrentMusic, ps.MusicVolume)
	}
}

// Odtworzenie pieśni
func playMusicByNumber(ps *programState, trackNumber int) error {
	if !ps.MusicEnabled {
		return nil
	}
	if trackNumber < 1 || trackNumber > len(ps.GameMusicFiles) {
		return fmt.Errorf("nieprawidłowy numer pieśni: %d", trackNumber)
	}
	if ps.CurrentMusic.Stream.Buffer != nil {
		rl.StopMusicStream(ps.CurrentMusic)
		rl.UnloadMusicStream(ps.CurrentMusic)
	}

	filePath := filepath.Join(ps.BaseAssetsPath, ps.GameMusicFiles[trackNumber-1])

	ps.CurrentMusic = rl.LoadMusicStream(filePath)
	rl.PlayMusicStream(ps.CurrentMusic)
	rl.SetMusicVolume(ps.CurrentMusic, ps.MusicVolume)
	ps.CurrentTruckNumber = trackNumber
	log.Printf("INFO: Odtwarzanie pieśni: %s (ścieżka %d)", ps.GameMusicFiles[trackNumber-1], trackNumber)
	return nil
}

// Uaktualnienie pieśni
func updateMusic(ps *programState) {
	if ps.CurrentMusic.Stream.Buffer != nil && ps.MusicEnabled {
		rl.UpdateMusicStream(ps.CurrentMusic)
	}
}

// Wyłącz odtwarzanie pieśni
func shutDownMusicPlayer(ps *programState) {
	if ps.CurrentMusic.Stream.Buffer != nil {
		rl.StopMusicStream(ps.CurrentMusic)
		rl.UnloadMusicStream(ps.CurrentMusic)
		ps.MusicEnabled = false
	}
	log.Println("INFO: Układ odtwarzania pieśni wyłączony.")
}

func setupMenuCamera(ps *programState) {
	ps.MenuCamera = rl.NewCamera2D(
		rl.NewVector2(ps.VirtualWidth/2.0, ps.VirtualHeight/2.0),
		rl.NewVector2(ps.VirtualWidth/2.0, ps.VirtualHeight/2.0), // Cel też na środek
		0.0,
		1.0,
	)
}

func drawMainMenu(ps *programState) {
	rl.BeginMode2D(ps.MenuCamera)
	rl.ClearBackground(rl.Black) // czemu robię to tylko tutaj, a np. w select campaign już nie?

	tex := ps.Assets.getSpecial(specialMainMenu)
	rl.DrawTexture(tex, 0, 0, rl.White)

	for _, btn := range ps.ActiveMenuButtons {
		rl.DrawRectangleLinesEx(btn.Rectangle, 1, rl.Green)
	}
	rl.EndMode2D()
}

// Odpowiada za tłumaczenie współrzędnych z płótna rzeczywistego
// do płótna urojonego, które jest płótnem podstawowym w programie
func screenToVirtualCoords(ps *programState, screenPos rl.Vector2) rl.Vector2 {
	if ps.CurrentState == gameScreen {
		// Mysz w bitwie
		// Proste skalowanie, bo gra wypełnia ekran
		scale := float32(ps.ActualScreenHeight) / ps.VirtualHeight
		return rl.NewVector2(screenPos.X/scale, screenPos.Y/scale)

	}

	// Mysz w menu
	scaleX := float32(ps.ActualScreenWidth) / 640.0
	scaleY := float32(ps.ActualScreenHeight) / 400.0

	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	// Obliczamy przesunięcie (czarne pasy)
	offsetX := (float32(ps.ActualScreenWidth) - (640.0 * scale)) / 2.0
	offsetY := (float32(ps.ActualScreenHeight) - (400.0 * scale)) / 2.0

	// Odejmujemy pasy i dzielimy przez skalę
	virtualX := (screenPos.X - offsetX) / scale
	virtualY := (screenPos.Y - offsetY) / scale

	return rl.NewVector2(virtualX, virtualY)
}

// Na podstawie ps.currentState decyduje, który ekran rysować.
func (ps *programState) changeState(newState screenState, bs *battleState) {
	ps.CurrentState = newState

	switch newState {
	case gameScreen:
		// Logika wejścia do gry
		log.Println("SYSTEM: Transition to GAME_SCREEN → Recalculating View.")
		ps.recalculateVirtualResolution() // Naprawia UI i skalowanie
		setupGameCamera(bs, ps)           // Centruje kamerę na nowym wymiarze

	case mainMenuScreen, newCampaignMenuScreen:
		// Logika wejścia do menu
		ps.recalculateVirtualResolution()
		setupMenuCamera(ps)
	}
}

// Odpowiada za rysowanie urojonego płótna, które później ma być skalowane do docelowej wielkości
func drawSceneToVirtualScreen(bs *battleState, ps *programState) {
	rl.BeginTextureMode(ps.ScreenTarget)
	rl.ClearBackground(rl.Black)

	// @todo: czemu ten przełącznik powtarza się w renderCurrentScene?! Co przeoczyłem?!
	switch ps.CurrentState {
	case gameScreen:
		// 1. Rysuj GRĘ (Po lewej)
		// Nożyczki ucinają wszystko, co wchodzi na teren panelu UI
		if ps.GameViewWidth > 0 {
			rl.BeginScissorMode(0, 0, int32(ps.GameViewWidth), int32(ps.VirtualHeight))
			rl.BeginMode2D(bs.GameCamera)
			drawWorldAndUnits(bs, ps)
			drawConstructionDebugBox(bs, ps)
			rl.EndMode2D()
			rl.EndScissorMode()
		}

		drawGameUI(bs, ps)

	default:
		// Rysuj menu, intro, itp.
		renderCurrentScene(bs, ps)
	}

	rl.EndTextureMode()
}

func drawSceneToActualScreen(bs *battleState, ps *programState) {
	rl.ClearBackground(rl.Black) // Czarne tło (pasy)

	// 1. USTAWIENIE ŹRÓDŁA (Source)
	ps.RenderSrcRect.X = 0
	ps.RenderSrcRect.Y = 0
	ps.RenderSrcRect.Width = float32(ps.ScreenTarget.Texture.Width)
	ps.RenderSrcRect.Height = float32(-ps.ScreenTarget.Texture.Height)

	// 2. USTAWIENIE CELU (Destination) i SKALI
	var scale float32

	if ps.CurrentState == gameScreen {
		// Gra na cały ekran
		ps.RenderDestRect.X = 0
		ps.RenderDestRect.Y = 0
		ps.RenderDestRect.Width = float32(ps.ActualScreenWidth)
		ps.RenderDestRect.Height = float32(ps.ActualScreenHeight)

		scale = float32(ps.ActualScreenHeight) / ps.VirtualHeight
	} else {
		// Menu z zachowaniem proporcji
		scaleX := float32(ps.ActualScreenWidth) / 640.0
		scaleY := float32(ps.ActualScreenHeight) / 400.0

		scale = scaleX
		if scaleY < scaleX {
			scale = scaleY
		}

		finalWidth := 640.0 * scale
		finalHeight := 400.0 * scale
		posX := (float32(ps.ActualScreenWidth) - finalWidth) / 2.0
		posY := (float32(ps.ActualScreenHeight) - finalHeight) / 2.0

		ps.RenderDestRect.X = posX
		ps.RenderDestRect.Y = posY
		ps.RenderDestRect.Width = finalWidth
		ps.RenderDestRect.Height = finalHeight
	}

	// 3. RYSOWANIE
	rl.DrawTexturePro(
		ps.ScreenTarget.Texture,
		ps.RenderSrcRect,
		ps.RenderDestRect,
		ps.RenderOrigin,
		0.0,
		rl.White,
	)

	drawGameCursorOnRealScreen(bs, ps, scale)
}

// Odpowiada za dopasowanie rozmiaru okna w ProgramState do rzeczywistego
// Pozwala prawidłowo dopasować tekstury do okna
func updateWindowSize(bs *battleState, ps *programState) {
	newWidth := int32(rl.GetScreenWidth())
	newHeight := int32(rl.GetScreenHeight())

	// Sprawdzamy, czy zmienił się rozmiar okna
	if ps.ActualScreenWidth != newWidth || ps.ActualScreenHeight != newHeight {
		ps.ActualScreenWidth = newWidth
		ps.ActualScreenHeight = newHeight

		// 1. Przeliczamy wymiary wirtualne (GameViewWidth, VirtualHeight)
		ps.recalculateVirtualResolution()

		// 2. Aktualizujemy kamerę MENU (centrowanie)
		setupMenuCamera(ps)

		// 3. Aktualizujemy kamerę GRY (Offset musi być nowym środkiem ekranu!)
		if bs != nil {
			// Offset to punkt na ekranie, gdzie rysowany jest cel (Target) kamery.
			// Musi być na środku dynamicznego widoku gry.
			bs.GameCamera.Offset = rl.NewVector2(
				ps.GameViewWidth/2.0,
				ps.VirtualHeight/2.0,
			)

			// Opcjonalnie: Po zmianie rozmiaru warto upewnić się, że kamera nie patrzy w pustkę
			// (wywołujemy clamping na obecnej pozycji)
			fullMapWidth := float32(uint16(boardMaxX) * uint16(tileWidth))
			fullMapHeight := float32(uint16(boardMaxY) * uint16(tileHeight))
			clampCameraTarget(&bs.GameCamera, fullMapWidth, fullMapHeight, ps.GameViewWidth, ps.VirtualHeight)
		}
	}
}

// Zwraca współrzędne myszy. Używać do nastawiania przycisków
func logVirtualMouseCoordinates(ps *programState) {
	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		screenMousePos := rl.GetMousePosition()
		// Proste skalowanie, bez zbędnych offsetów, bo tekstura wypełnia ekran
		scale := float32(ps.ActualScreenHeight) / ps.VirtualHeight
		virtualX := screenMousePos.X / scale
		virtualY := screenMousePos.Y / scale

		log.Printf("DEBUG: Mysz Virtual: %.1f, %.1f (Screen: %.1f, %.1f)",
			virtualX, virtualY, screenMousePos.X, screenMousePos.Y)
	}
}

func startNewCampaign(ps *programState) {
	drawSelectCampaignMenuLegacy(ps)
}

// Odpowiada za rysowanie menu wyboru wyprawy.
func drawSelectCampaignMenuLegacy(ps *programState) {
	rl.BeginMode2D(ps.MenuCamera)

	tex := ps.Assets.getSpecial(specialCampaignSelect)
	rl.DrawTexture(tex, 0, 0, rl.White)

	drawEyes(ps)

	rl.EndMode2D()
}

// @todo: wywal zaczarowane liczby!
func drawEyes(ps *programState) {
	drawEye := func(x, y int32, color rl.Color) {
		rl.DrawLine(x, y, x+5, y, color)
		rl.DrawLine(x, y+1, x+5, y+1, color)
	}

	// Współrzędne oczu (bazowe dla 640×400)
	// @todo: to nie powinny być „liczby z palca”. Muszę przenieść do constants.go
	var yBase int32

	switch ps.SelectedDifficulty {
	case 0: // Łatwy
		yBase = 70
	case 1: // Normalny
		yBase = 156
	case 2: // Trudny
		yBase = 242
	default:
		return
	}

	// Lewe oko
	// @todo: to nie powinny być „liczby z palca”. Muszę przenieść do constants.go
	drawEye(519, yBase+1, rl.Gold)
	drawEye(519, yBase, rl.Gold)
	// Prawe oko
	// @todo: to nie powinny być „liczby z palca”. Muszę przenieść do constants.go
	drawEye(537, yBase+1, rl.Gold)
	drawEye(537, yBase, rl.Gold)
}

// Odpowiada za wybór funkcji do rysowania właściwego widoku
// menu głównego, wybór wyprawy itd
func renderCurrentScene(bs *battleState, ps *programState) {
	switch ps.CurrentState {
	case mainMenuScreen:
		drawMainMenu(ps)
	case newCampaignMenuScreen:
		drawSelectCampaignMenuLegacy(ps)
	case gameScreen:
		drawBattleScene(bs)
	default:
		log.Printf("OSTRZEŻENIE: Nieznany CurrentState (%d) do narysowania. Rysuję menu główne.", ps.CurrentState)
		// @todo: ogarnij, czy dobrze jest wrócić do głównego jeżeli nie jesteśmy pewni, co narysować
		// Chyba program powinien się wywalić i zostawić coś w logach
		drawMainMenu(ps)
	}
}

// Dobiera logikę obsługi wejścia na postawie programState.CurrentState
func handleCurrentScreenInput(bs *battleState, ps *programState) {
	// logVirtualMouseCoordinates(ps) // @todo: zakomentuj później, bo zapycha konsolę niepotrzebnie
	switch ps.CurrentState {
	case mainMenuScreen:
		handleMenuInput(ps)
	case newCampaignMenuScreen:
		handleMenuInput(ps)
	case gameScreen:
		handleGameInput(bs, ps)
	default:
		log.Printf("OSTRZEŻENIE: Nieznany stan (%d) dla obsługi wejścia", ps.CurrentState)
	}
}

// Odpowiada za wywołanie odświeżania obrazu
func updateCurrentScreen(bs *battleState, ps *programState) {
	switch ps.CurrentState {
	case mainMenuScreen:
		// można zostawić puste, bo nie powinno być logiki zmiany ekranu co klatkę
	case newCampaignMenuScreen:
		// można zostawić puste, bo nie powinno być logiki zmiany ekranu co klatkę
	case gameScreen:
		updateGame(bs)
		updateActionButtons(bs)
	}
}

// Odpowiada za wczytanie odpowiedniej rodziny tekstur
// @todo: to opiera się o ułomny i tymczasowy sposób ładowania tekstur.
// Do przebudowania wraz z wdrożeniem obsługi .dat!
func (ps *programState) loadTextureCategory(texMap map[int]rl.Texture2D, dirPath string,
	re *regexp.Regexp, categoryName string,
) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("OSTRZEŻENIE: Nie można odczytać katalogu %s ('%s'): %v", categoryName, dirPath, err)
		return
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		matches := re.FindStringSubmatch(file.Name())
		if len(matches) > 1 {
			id, err := strconv.Atoi(matches[len(matches)-1])
			if err != nil {
				log.Printf("BŁĄD: %s - Nie można sparsować ID '%s' z pliku '%s'",
					categoryName, matches[len(matches)-1], file.Name())
				continue
			}

			if _, found := texMap[id]; found {
				log.Printf("OSTRZEŻENIE: %s - Konflikt ID %d. Plik '%s' zignorowany", categoryName, id, file.Name())
				continue
			}

			texPath := filepath.Join(dirPath, file.Name())
			texture := rl.LoadTexture(texPath)
			if texture.ID == 0 {
				log.Printf("BŁĄD: %s - Nie udało się załadować: %s", categoryName, texPath)
			} else {
				texMap[id] = texture
			}
		}
	}
	log.Printf("INFO: %s - Załadowano %d tekstur", categoryName, len(texMap))
}

// @todo: sprawdź, czy do zmiany po ogarnięciu ładowania zasobów z .dat
func (ps *programState) loadSingleTexture(path string) (rl.Texture2D, error) {
	texture := rl.LoadTexture(path)
	if texture.ID == 0 {
		return rl.Texture2D{}, fmt.Errorf("nie udało się załadować tekstury: %s", path)
	}
	return texture, nil
}

// newBattleState @todo: w kampanii gracz zawsze czerwony, ale przeciwnicy to już różnie
// zmień to zakodowane na sztywno na dynamiczne, w zależności od prowincji
// te dane przechowujemy nawet w „prowintionInit”
// w ogóle powinno się przenieść określenie battleState w miejsce, gdzie odpalamy bitwy
func newBattleState(ps *programState) *battleState {
	return &battleState{
		// Ustawiamy na NONE (0) - stan jest "nieuzbrojony"
		PlayerID:   colorNone,
		AIPlayerID: colorNone,

		// Alokacja pamięci dla obiektów dynamicznych
		Units:       make([]*unit, 0, maxUnitsPerPlayer*2),
		Buildings:   make([]*building, 0, maxBuildingsPerPlayer*2),
		Projectiles: make([]*projectile, 0, 100),
		Corpses:     make([]corpse, 0, 100),

		// Inicjalizacja struktur graczy wartościami zerowymi/pustymi
		// Handler nadpisze PlayerID i Milk właściwymi danymi przed startem
		HumanPlayerState: &playerState{PlayerID: colorNone, Milk: 0, MaxMilk: 0},
		AIEnemyState:     &playerState{PlayerID: colorNone, Milk: 0, MaxMilk: 0},

		AI:              aiState{},
		CurrentCommands: [2]command{}, // Puste komendy

		// Konfiguracja Kamery (To jest uniwersalne, może zostać)
		GameCamera: rl.NewCamera2D(
			rl.NewVector2(float32(virtualScreenWidth/2), float32(virtualScreenHeight/2)),
			rl.NewVector2(float32(virtualScreenWidth/2), float32(virtualScreenHeight/2)),
			0.0,
			1.0,
		),

		// Reszta domyślnych ustawień technicznych
		NextUniqueObjectID: 1,
		Board:              &boardData{},
		CurrentSelection:   selectionState{},
		MouseCommandMode:   1,
		IsSelectingBox:     false,
		InitialClickPos:    rl.Vector2{X: 0, Y: 0},
		SelectionStart:     rl.Vector2{X: 0, Y: 0},
		Map:                false,

		CurrentLevel:          0,
		DifficultyLevel:       ps.SelectedDifficulty,
		GameSpeed:             1,
		IsSinglePlayerGame:    true, // To zazwyczaj prawda, ale Handler może zmienić
		CheatsEnabled:         false,
		CheatSequenceProgress: 0,
		BattleOutcome:         0,
		QuitLevel:             false,

		// Liczniki
		GrassGrowthCounter:    0,
		WaterAnimationCounter: 0,
		FireAnimationCounter:  0,
		GlobalFrameCounter:    0,

		CampaignData: campaignData{},

		// Cache dla AI
		EnemyCache:                make(map[uint]*enemyCacheEntry),
		enemyCacheUpdateTick:      0,
		pathfindingUnitsThisTick:  0,
		enemyCacheUpdatesThisTick: 0,

		IsMapDragging:           false,
		MapInitialClickPos:      rl.NewVector2(0, 0),
		CameraTargetOnDragStart: rl.NewVector2(0, 0),
	}
}

// @todo: przenieś to do innego pliku, bo tutaj nie pasuje chyba level.go będzie lepsze, bo tam są różne „apply”.
func makeGrassVariations(bs *battleState) {
	for y := range boardMaxY {
		y5 := y * 5
		for x := range boardMaxX {
			tile := &bs.Board.Tiles[x][y]
			id := tile.TextureID

			if id >= spriteGrassStart && id <= spriteGrassEnd {
				hash := uint16((x + y5) & 0x0F)
				finalID := spriteGrassStart + hash

				tile.TextureID = finalID
			}
		}
	}
}

func getEnemyColor(levelNumber uint8) uint8 {
	provinceIndex := levelNumber - 1

	if provinceIndex < uint8(len(provinceInit)) {
		return provinceInit[provinceIndex]
	}

	// bezpiecznik
	// @todo: ogarnij coś sensowniejszego dla nowych kampanii
	return colorGray
}

func main() {
	// Debugger (opcjonalny)
	//go func() {
	//	log.Println("Uruchamiam serwer pprof na localhost:6060")
	//	log.Println(http.ListenAndServe("localhost:6060", nil))
	//}()

	// 1. Inicjalizacja zmiennych stanu (bez grafiki!)
	programState := newProgramState()

	// 2. Ścieżki
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Krytyczny błąd: Nie można ustalić ścieżki pliku wykonywalnego: %v", err)
	}
	executableDir := filepath.Dir(executablePath)
	programState.BaseAssetsPath = filepath.Join(executableDir, "assets")
	log.Printf("Bazowa ścieżka zasobów ustawiona na: %s", programState.BaseAssetsPath)

	// 3. Okno
	rl.SetConfigFlags(rl.FlagWindowResizable | rl.FlagWindowAlwaysRun)
	rl.InitWindow(programState.ActualScreenWidth, programState.ActualScreenHeight, "Opolanie")
	rl.SetWindowMinSize(640, 400)
	defer rl.CloseWindow()

	// 4. Audio
	rl.InitAudioDevice()
	defer func() {
		shutDownMusicPlayer(programState)
		rl.CloseAudioDevice()
	}()

	// 5. System Plików (.DAT)
	datPath := filepath.Join(programState.BaseAssetsPath, "graf.dat")
	palPath := filepath.Join(programState.BaseAssetsPath, "pal.dat")

	loader, err := newAssetLoader(datPath, palPath)
	if err != nil {
		log.Printf("BŁĄD KRYTYCZNY: Nie udało się otworzyć .dat: %v", err)
	} else {
		log.Println("INFO: assetLoader załadowany pomyślnie.")
		defer loader.close()
	}

	// 6. Inicjalizacja AssetManagera
	if loader != nil {
		programState.Assets = newAssetManager(loader)
		log.Println("INFO: AssetManager utworzony i gotowy do pracy.")
	}

	if err := programState.Assets.loadGlobalAssets(programState.BaseAssetsPath); err != nil {
		log.Fatalf("BŁAD KRYTYCZNY: nie udało się załadować UI %v", err)
	}
	// 7. Stan bitwy (Pusty)
	battleState := newBattleState(programState)

	setupMainMenuButtons(programState, battleState)

	// 8. Skalowanie ekranu
	programState.ActualScreenWidth = int32(rl.GetScreenWidth())
	programState.ActualScreenHeight = int32(rl.GetScreenHeight())
	programState.recalculateVirtualResolution()

	if programState.ScreenTarget.ID == 0 {
		rl.CloseWindow()
		log.Fatalln("Nie udało się załadować tekstury renderującej (Błąd FBO)")
	}
	rl.SetTextureFilter(programState.ScreenTarget.Texture, rl.FilterPoint)
	defer func() {
		if programState.ScreenTarget.ID != 0 {
			rl.UnloadRenderTexture(programState.ScreenTarget)
		}
	}()

	rl.SetTargetFPS(60)

	if rl.IsWindowResized() {
		updateWindowSize(battleState, programState)
	}

	setupMenuCamera(programState)
	setupGameCamera(battleState, programState)

	// Pętla stałego kroku czasowego
	timeAccumulator := 0.0

	// === PĘTLA GŁÓWNA ===
	for !programState.ShouldClose && !rl.WindowShouldClose() {
		if programState.CurrentState == gameScreen && rl.GetScreenWidth() < 640+120 {
			rl.SetWindowMinSize(640+120, 400)
			rl.SetWindowSize(640+120, 400)
		}

		dt := float64(rl.GetFrameTime())
		if dt > 0.25 {
			dt = 0.25
		}

		timeAccumulator += dt

		for timeAccumulator >= tickRate {

			updateWindowSize(battleState, programState)
			handleCurrentScreenInput(battleState, programState)
			updateCurrentScreen(battleState, programState)

			battleState.GlobalFrameCounter++
			timeAccumulator -= tickRate
		}

		rl.BeginDrawing()
		logVirtualMouseCoordinates(programState)
		drawSceneToVirtualScreen(battleState, programState)
		drawSceneToActualScreen(battleState, programState)
		rl.EndDrawing()
	}
}
