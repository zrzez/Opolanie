package main

type battleAtlasID uint8

const (
	atlasUI battleAtlasID = iota
	atlasUnits1
	atlasUnits2
	atlasBuildings
	// ↓↓↓ZAWSZE NA KOŃCU!↓↓↓.
	atlasCount // Musi być na końcu, zawsze!
)

// Rysopis duszka dla jednostek
type spriteDef struct {
	atlasID               battleAtlasID // cropWidth którym atlasie
	cropX, cropY          uint16        // współrzędne; 303 to największa liczba, wiec uint8 jest zbyt małe
	cropWidth, cropHeight uint8
	offX                  int8 // poprawka cropX
	offY                  int8 // poprawka cropY
	flipX                 bool // Wskaźnik, czy potrzeba zrobić zwierciadlane odbicie
}

type rawAssetDef struct {
	topChunk int
	botChunk int
	palette  int
}

// Liczba wszystkich składowych; mnie wyszło 743, ale dam zapas
// @reminder: @todo: sprawdź ile rzeczywiście się ładuje i później popraw
const maxSpriteID uint16 = 4096
