package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type assetManager struct {
	//              ↓ barwa gracza   ↓ ID atlasu np. UI, units1, units2
	Atlases          [maxGameColors][atlasCount]rl.Texture2D
	SpecialTextures  [specialCount]rl.Texture2D // Mapa Polski, Menu, etc.
	WoodPanel        rl.Texture2D               // Nakładka widoczna w czasie bitwy po prawej stronie
	CursorWhiteFrame rl.Texture2D               // Biała ramka wskazująca, która z naszych jednostek jest zaznaczona
	loader           *assetLoader               // Wczytywacz
}

func newAssetManager(loader *assetLoader) *assetManager {
	return &assetManager{
		loader: loader,
	}
}

func (am *assetManager) loadGlobalAssets(basePath string) error {
	log.Println("Zasoby: Ładowanie zasobów globalnych (nakładka)...")

	// 1. Atlas UI
	// Pobieramy definicję kawałków z mapy w assets_db.go
	uiDef, ok := atlasDefinitions[atlasUI]
	if !ok {
		return fmt.Errorf("błąd: brak definicji dla AtlasUI w assets_db.go")
	}

	if err := am.loadAndProcess(atlasUI, uiDef, colorNone, nil); err != nil {
		return fmt.Errorf("błąd ładowania AtlasUI: %v", err)
	}

	// Biała ramka jednostek
	img := rl.LoadImageFromTexture(am.Atlases[colorNone][atlasUI])
	rl.ImageCrop(img, rl.NewRectangle(75, 8, 16, 14))
	for y := int32(0); y < img.Height; y++ {
		for x := int32(0); x < img.Width; x++ {
			px := rl.GetImageColor(*img, x, y)
			if px.A != 0 {
				rl.ImageDrawPixel(img, x, y, rl.White)
			}
		}
	}

	am.CursorWhiteFrame = rl.LoadTextureFromImage(img)
	rl.SetTextureFilter(am.CursorWhiteFrame, rl.FilterPoint)
	rl.UnloadImage(img)

	// 2. Nakładka UI (new_wood)
	if err := am.loadExternalAssets(basePath); err != nil {
		return fmt.Errorf("błąd ładowania new_wood: %v", err)
	}

	// 3. Ekrany specjalne (Menu główne, Mapa Polski itp.)
	for index := 0; index < specialCount; index++ {
		def := specialAssetsDB[index]
		if err := am.loadRawAssetScaled(index, def); err != nil {
			return fmt.Errorf("błąd ładowania ekranu ID %d: %v", index, err)
		}
	}

	return nil
}

// Potrzebna do zrobienia grafik o wymiarach 640na400 (skalowanie 2x w pamięci).
func (am *assetManager) loadRawAssetScaled(specialID int, def rawAssetDef) error {
	// Używamy nowej metody LoadRawImage z assets_loader.go
	srcPixels, err := am.loader.loadRawImage(def)
	if err != nil {
		return err
	}

	const srcW, srcH = 320, 200
	const dstW, dstH = 640, 400

	dstDataSize := dstW * dstH * 4
	dstPixels := make([]byte, dstDataSize)

	// Ręczne skalowanie, każdy piksel zamienia się w blok 2x2
	for sy := 0; sy < srcH; sy++ {
		for sx := 0; sx < srcW; sx++ {
			// Pobierz piksel źródłowy (4 bajty)
			sIdx := (sy*srcW + sx) * 4
			r, g, b, a := srcPixels[sIdx], srcPixels[sIdx+1], srcPixels[sIdx+2], srcPixels[sIdx+3]

			// Wylicz koordynaty w celu (2x2 blok)
			dy := sy * 2
			dx := sx * 2

			// Wypełnij 4 piksele w docelowym buforze
			// Wiersz 1 (dy)
			dIdx1 := (dy*dstW + dx) * 4
			dstPixels[dIdx1], dstPixels[dIdx1+1], dstPixels[dIdx1+2], dstPixels[dIdx1+3] = r, g, b, a   // Lewy
			dstPixels[dIdx1+4], dstPixels[dIdx1+5], dstPixels[dIdx1+6], dstPixels[dIdx1+7] = r, g, b, a // Prawy

			// Wiersz 2 (dy+1)
			dIdx2 := ((dy+1)*dstW + dx) * 4
			dstPixels[dIdx2], dstPixels[dIdx2+1], dstPixels[dIdx2+2], dstPixels[dIdx2+3] = r, g, b, a   // Lewy
			dstPixels[dIdx2+4], dstPixels[dIdx2+5], dstPixels[dIdx2+6], dstPixels[dIdx2+7] = r, g, b, a // Prawy
		}
	}

	finalImg := rl.NewImage(dstPixels, int32(dstW), int32(dstH), 1, rl.UncompressedR8g8b8a8)
	tex := rl.LoadTextureFromImage(finalImg)
	rl.SetTextureFilter(tex, rl.FilterPoint)

	am.SpecialTextures[specialID] = tex

	return nil
}

func (am *assetManager) loadBattleAssets(activePlayers []uint8) error {
	log.Println("Zasoby: Ładowanie zasobów (jednostki, budynki)...")

	// Wykaz atlasów do przetworzenia dla każdego gracza
	gameplayAtlases := []battleAtlasID{atlasUnits1, atlasUnits2, atlasBuildings}

	for _, playerColor := range activePlayers {
		// Ustalenie palety docelowej
		var targetPalette []rl.Color
		if playerColor != colorRed && playerColor != colorNone {
			targetPalette = enemyColors[playerColor]
		}

		// Ładowanie i kolorowanie
		for _, atlasID := range gameplayAtlases {
			// Pobieramy definicję z assets_db.go
			def, ok := atlasDefinitions[atlasID]
			if !ok {
				return fmt.Errorf("brak definicji kawałka dla atlasu ID %d", atlasID)
			}

			if err := am.loadAndProcess(atlasID, def, playerColor, targetPalette); err != nil {
				return err
			}
		}
	}

	return nil
}

func (am *assetManager) loadAndProcess(atlasID battleAtlasID, def rawAssetDef, ownerColor uint8, palette []rl.Color) error {
	// Jeśli mamy już atlasy dla czerwonego, to nie ma sensu znowu go ładować
	if ownerColor < maxGameColors && am.Atlases[ownerColor][atlasID].ID != 0 {
		return nil
	}

	// 1. Wczytaj do RAM
	pixels, err := am.loader.loadRawImage(def)
	if err != nil {
		return fmt.Errorf("błąd ładowania obrazu atlasu %d: %v", atlasID, err)
	}

	// 2. Przebarwianie
	if len(palette) > 0 {
		replaceColorsInMemory(pixels, palette)
	}

	// 3. Wysyłka do VRAM
	img := rl.NewImage(pixels, 320, 200, 1, rl.UncompressedR8g8b8a8)
	texture := rl.LoadTextureFromImage(img)
	rl.SetTextureFilter(texture, rl.FilterPoint)

	// 4. Zapis do tablicy
	if ownerColor < maxGameColors {
		am.Atlases[ownerColor][atlasID] = texture
	}

	return nil
}

// Podmienia kolory operując bezpośrednio na bajtach.
func replaceColorsInMemory(pixels []byte, targetPalette []rl.Color) {
	dataSize := len(pixels)

	// Przygotowanie mapy kolorów
	count := len(targetPalette)
	if count > len(playerColors) {
		count = len(playerColors)
	}

	// Iteracja po wszystkich pikselach
	for i := 0; i < dataSize; i += 4 {
		// Czy piksel jest przezroczysty?
		if pixels[i+3] == 0 {
			continue
		}

		r := pixels[i]
		g := pixels[i+1]
		b := pixels[i+2]
		// a := pixels[i+3]

		// Sprawdź czy piksel pasuje do któregoś z kolorów bazowych gracza
		for c := 0; c < count; c++ {
			base := playerColors[c]

			if r == base.R && g == base.G && b == base.B {
				// Podmieniamy na kolor z palety celu
				target := targetPalette[c]
				pixels[i] = target.R
				pixels[i+1] = target.G
				pixels[i+2] = target.B

				break
			}
		}
	}
}

func (am *assetManager) unloadBattleAssets() {
	// Czyścimy atlasy graczy (Units1, Units2, Buildings)
	// Przechodzimy po wszystkich kolorach
	for i := 0; i < int(maxGameColors); i++ {
		if i == int(colorRed) {
			continue
		}
		// Przechodzimy po atlasach, ale pomijamy nakładkę UI
		for j := 0; j < int(atlasCount); j++ {
			atlasID := battleAtlasID(j)

			if atlasID == atlasUI {
				continue
			}

			if am.Atlases[i][atlasID].ID != 0 {
				rl.UnloadTexture(am.Atlases[i][atlasID])
				am.Atlases[i][atlasID] = rl.Texture2D{}
			}
		}
	}
	log.Println("Zasoby: Wyczyszczono zasoby bitewne (VRAM zwolniony).")
}

func (am *assetManager) unload() {
	for i := 0; i < int(maxGameColors); i++ {
		for j := 0; j < int(atlasCount); j++ {
			if am.Atlases[i][j].ID != 0 {
				rl.UnloadTexture(am.Atlases[i][j])
				am.Atlases[i][j] = rl.Texture2D{}
			}
		}
	}

	for i := range specialCount {
		if am.SpecialTextures[i].ID != 0 {
			rl.UnloadTexture(am.SpecialTextures[i])
			am.SpecialTextures[i] = rl.Texture2D{}
		}
	}
}

func (am *assetManager) getAtlas(id battleAtlasID, ownerColor uint8) rl.Texture2D {
	if ownerColor >= maxGameColors {
		return am.Atlases[colorRed][id]
	}

	tex := am.Atlases[ownerColor][id]

	if tex.ID == 0 {
		return am.Atlases[colorRed][id]
	}

	return tex
}

func (am *assetManager) getSpecial(id int) rl.Texture2D {
	if id < 0 || id >= specialCount {
		return rl.Texture2D{}
	}

	return am.SpecialTextures[id]
}

func (am *assetManager) loadExternalAssets(basePath string) error {
	woodPath := filepath.Join(basePath, "new_wood.png")

	// Bezpiecznik na wypadek braku pliku w oczekiwanym miejscu
	if _, err := os.Stat(woodPath); os.IsNotExist(err) {
		return fmt.Errorf("BŁAD KRYTYCZNY: brak pliku new_wood.png w oczekiwanym miejscu")
	}

	tex := rl.LoadTexture(woodPath)
	if tex.ID == 0 {
		return fmt.Errorf("BŁĄD KRYTYCZNY: nie udało się załadować new_wood")
	}

	rl.SetTextureFilter(tex, rl.FilterPoint)

	am.WoodPanel = tex
	return nil
}
