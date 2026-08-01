package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Rzeczy związane z rysowaniem mapy. Nie mylić z planszą.
var (
	mapXGrid           [boardMaxX + 1]int
	mapYGrid           [boardMaxY + 1]int
	mapGridInitialized bool
	mapTerrain         rl.RenderTexture2D
	mapTerrainReady    bool
	mapTerrainInit     = true
)

func initMapGrid(mapWidth, mapHeight float32) {
	if mapGridInitialized {
		return
	}

	totalDisplayWidthPx := int(math.Round(float64(mapWidth)))
	totalDisplayHeightPx := int(math.Round(float64(mapHeight)))

	for index := uint8(0); index <= boardMaxX; index++ {
		currentXFloat := float64(index) / float64(boardMaxX) * float64(totalDisplayWidthPx)
		mapXGrid[index] = int(math.Round(currentXFloat))
	}

	for index := uint8(0); index <= boardMaxY; index++ {
		currentYFloat := float64(index) / float64(boardMaxY) * float64(totalDisplayHeightPx)
		mapYGrid[index] = int(math.Round(currentYFloat))
	}

	mapXGrid[boardMaxX] = totalDisplayWidthPx
	mapYGrid[boardMaxY] = totalDisplayHeightPx

	mapGridInitialized = true
}

func updateMapTerrainCache(bState *battleState, pState *programState) {
	if !mapTerrainReady {
		return
	}

	rl.BeginTextureMode(mapTerrain)
	rl.ClearBackground(rl.Black)

	for x := range boardMaxX {
		for y := range boardMaxY {
			tileID := bState.Board.Tiles[x][y].TextureID
			color := getMapColor(tileID)

			w := mapXGrid[x+1] - mapXGrid[x]
			h := mapYGrid[y+1] - mapYGrid[y]

			rl.DrawRectangle(int32(mapXGrid[x]), int32(mapYGrid[y]), int32(w), int32(h), color)
		}
	}

	rl.EndTextureMode()
	rl.BeginTextureMode(pState.ScreenTarget)
}

func drawMinimapUnits(mapX, mapY, mapWidth, mapHeight, actualGameViewWidth float32, bState *battleState, pState *programState) {
	const worldMapFullPixelWidth = float32(boardMaxX) * float32(tileWidth)
	const worldMapFullPixelHeight = float32(boardMaxY) * float32(tileHeight)

	pxStart := int(math.Round(float64(mapX)))
	pyStart := int(math.Round(float64(mapY)))

	// 1. Inicjalizacja siatki (wykona się tylko raz)
	initMapGrid(mapWidth, mapHeight)

	// 2. Odświeżenie cache'u terenu
	if mapTerrainInit || bState.GlobalFrameCounter%60 == 0 || !mapTerrainReady {
		updateMapTerrainCache(bState, pState)
		mapTerrainInit = false
	}

	// 3. Rysowanie terenu z cache'u
	if mapTerrainReady {
		srcRec := rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  float32(mapTerrain.Texture.Width),
			Height: -float32(mapTerrain.Texture.Height),
		}
		destPos := rl.Vector2{X: float32(pxStart), Y: float32(pyStart)}
		rl.DrawTextureRec(mapTerrain.Texture, srcRec, destPos, rl.White)
	}

	// 4. Budynki
	for _, bld := range bState.Buildings {
		if !bld.Exists || bld.Type == buildingPalisade {
			continue
		}

		var bldColor rl.Color

		if bld.Owner == bState.PlayerID {
			bldColor = rl.White
		} else {
			bldColor = rl.Red
		}

		for _, currentTile := range bld.OccupiedTiles {
			w := mapXGrid[currentTile.X+1] - mapXGrid[currentTile.X]
			h := mapYGrid[currentTile.Y+1] - mapYGrid[currentTile.Y]

			rl.DrawRectangle(int32(pxStart+mapXGrid[currentTile.X]), int32(pyStart+mapYGrid[currentTile.Y]), int32(w), int32(h), bldColor)
		}
	}

	// 5. Jednostki
	for _, currentUnit := range bState.Units {
		if !currentUnit.Exists {
			continue
		}
		var unitColor rl.Color
		if currentUnit.Owner == bState.PlayerID {
			unitColor = rl.White
		} else {
			unitColor = rl.Red
		}

		w := mapXGrid[currentUnit.X+1] - mapXGrid[currentUnit.X]
		h := mapYGrid[currentUnit.Y+1] - mapYGrid[currentUnit.Y]

		rl.DrawRectangle(int32(pxStart+mapXGrid[currentUnit.X]), int32(pyStart+mapYGrid[currentUnit.Y]), int32(w), int32(h), unitColor)
	}

	// 6. Ramka Kamery
	camWorldView := getCameraWorldViewRect(bState.GameCamera, actualGameViewWidth, float32(virtualScreenHeight))
	scaleFactorX := mapWidth / worldMapFullPixelWidth
	scaleFactorY := mapHeight / worldMapFullPixelHeight

	rawCameraRect := rl.NewRectangle(
		mapX+(camWorldView.X*scaleFactorX),
		mapY+(camWorldView.Y*scaleFactorY),
		camWorldView.Width*scaleFactorX,
		camWorldView.Height*scaleFactorY,
	)

	minimapBounds := rl.NewRectangle(mapX, mapY, mapWidth, mapHeight)
	clampedRect := rl.GetCollisionRec(minimapBounds, rawCameraRect)
	rl.DrawRectangleLinesEx(clampedRect, 1.0, rl.Yellow)
}
