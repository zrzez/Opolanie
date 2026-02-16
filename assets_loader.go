package main

// assets_loader.go

import (
	"fmt"
	"io"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Stałe formatu pliku .dat
const (
	loaderChunkSize  = 33000
	loaderHeaderSize = 6
	loaderImageWidth = 320 // Szerokość w pamięci
	loaderFileWidth  = 319 // Szerokość danych w pliku @reminder to MUSI być 319
)

// Odpowiada za niskopoziomowy odczyt z pliku .dat.
type assetLoader struct {
	file     *os.File
	palettes [][]rl.Color
}

func newAssetLoader(grafPath, palPath string) (*assetLoader, error) {
	// 1. Ładowanie palet
	pals, err := loadAllPalettes(palPath)
	if err != nil {
		return nil, fmt.Errorf("błąd palet: %v", err)
	}

	// 2. Otwarcie pliku graficznego
	f, err := os.Open(grafPath)
	if err != nil {
		return nil, fmt.Errorf("błąd pliku graf: %v", err)
	}

	return &assetLoader{
		file:     f,
		palettes: pals,
	}, nil
}

func (al *assetLoader) close() {
	if al.file != nil {
		err := al.file.Close()
		if err != nil {
			return
		}
	}
}

// Wczytuje kawałki zdefiniowane w rawAssetDef i zwraca surowy bufor pikseli (RGBA).
func (al *assetLoader) loadRawImage(def rawAssetDef) ([]byte, error) {
	const linesPerChunk = 100
	const totalHeight = 200

	// Alokacja pamięci na pełny obraz 320x200
	pixels := make([]byte, loaderImageWidth*totalHeight*4)

	// Wybór palety
	palIdx := def.palette
	if palIdx >= len(al.palettes) {
		palIdx = 0
	}
	activePalette := al.palettes[palIdx]

	// Czytanie górnego kawałka (linie 0-99)
	err := al.readChunkToBuffer(def.topChunk, 0, linesPerChunk, activePalette, pixels)
	if err != nil {
		return nil, fmt.Errorf("błąd chunka %d: %v", def.topChunk, err)
	}

	// Czytanie dolnego kawałka (linie 100-199), jeśli jest inny
	if def.botChunk != def.topChunk {
		err = al.readChunkToBuffer(def.botChunk, 100, linesPerChunk, activePalette, pixels)
		if err != nil {
			return nil, fmt.Errorf("błąd chunka %d: %v", def.botChunk, err)
		}
	}

	return pixels, nil
}

// Czyta fragment pliku i mapuje go na piksele w buforze docelowym.
func (al *assetLoader) readChunkToBuffer(chunkID int, destY int, lines int, pal []rl.Color, dest []byte) error {
	offset := int64(chunkID) * loaderChunkSize
	if _, err := al.file.Seek(offset, 0); err != nil {
		return err
	}

	// Bufor na skompresowane dane z dysku
	buf := make([]byte, loaderChunkSize)
	if _, err := io.ReadFull(al.file, buf); err != nil {
		return err
	}

	ptr := loaderHeaderSize

	for i := 0; i < lines; i++ {
		targetY := destY + i

		for x := 0; x < loaderFileWidth; x++ {
			if ptr >= len(buf) {
				break
			}
			colIdx := buf[ptr]
			ptr++

			// Pobranie koloru z palety
			c := pal[colIdx]

			// Zapis do bufora docelowego (RGBA)
			pos := (targetY*loaderImageWidth + x) * 4

			if pos+3 >= len(dest) {
				continue
			}

			dest[pos] = c.R
			dest[pos+1] = c.G
			dest[pos+2] = c.B

			// Przezroczystość
			if colIdx == 0 {
				dest[pos+3] = 0
			} else {
				dest[pos+3] = 255
			}
		}
	}
	return nil
}

func loadAllPalettes(path string) ([][]rl.Color, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)

	st, err := f.Stat()
	if err != nil {
		return nil, err
	}

	// Rozmiar jednej palety to 256 kolorów * 3 bajty (RGB) = 768 bajtów
	count := int(st.Size() / 768)
	res := make([][]rl.Color, 0, count)

	for i := 0; i < count; i++ {
		buf := make([]byte, 768)
		if _, err := io.ReadFull(f, buf); err != nil {
			return nil, err
		}

		p := make([]rl.Color, 256)
		for j := 0; j < 256; j++ {
			r := buf[j*3]
			g := buf[j*3+1]
			b := buf[j*3+2]
			p[j] = rl.NewColor(r, g, b, 255)
		}
		res = append(res, p)
	}
	return res, nil
}
