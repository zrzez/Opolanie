package main

// drawingBattle.go

// @todo: OPTYMALIZACJA PAMIĘCI. Obecnie w pętli rysowania tworzone są nowiuśkie obiekty.
// (rl.NewRectangle, rl.NewVector2), co powoduje wzrost zużycia RAMu.
// Docelowo należy wprowadzić pre-alokowane zmienne robocze na wzór programState
// RenderSourceRect rl.Rectangle
// RenderDestRect   rl.Rectangle
// RenderOrigin     rl.Vector2
// zamiast tworzyć je w każdej klatce.

import (
	"math"
	"sort"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// === PRZETWARZANIE MAPY ===

func applyRoadProcessing(x, y uint8, board *boardData) {
	var mask uint8

	if y > 0 && isDirtRoad(board.Tiles[x][y-1].TextureID) {
		mask |= 1
	}

	if x < boardMaxX-1 && isDirtRoad(board.Tiles[x+1][y].TextureID) {
		mask |= 2
	}

	if y < boardMaxY-1 && isDirtRoad(board.Tiles[x][y+1].TextureID) {
		mask |= 4
	}

	if x > 0 && isDirtRoad(board.Tiles[x-1][y].TextureID) {
		mask |= 8
	}

	var finalSprite uint16

	switch mask {
	case 0:
		finalSprite = spriteRoadLU // lewo-góra @reminder: to raczej nie powinno być możliwe
	case 1:
		finalSprite = spriteRoadU // góra
	case 2:
		finalSprite = spriteRoadR // prawo
	case 3:
		finalSprite = spriteRoadRU // góra-prawo
	case 4:
		finalSprite = spriteRoadD // dół
	case 5:
		finalSprite = spriteRoadUD // góra-dół
	case 6:
		finalSprite = spriteRoadRD // prawy-dół
	case 7:
		finalSprite = spriteRoadRUD // prawo-góra-dół
	case 8:
		finalSprite = spriteRoadL // lewo
	case 9:
		finalSprite = spriteRoadLU // lewo-góra
	case 10:
		finalSprite = spriteRoadLR // lewo-prawo
	case 11:
		finalSprite = spriteRoadLRU // lewo-prawo-góra
	case 12:
		finalSprite = spriteRoadLD // lewo-dół
	case 13:
		finalSprite = spriteRoadLUD // lewo-góra-dół
	case 14:
		finalSprite = spriteRoadLRD // lewo-prawo-dół
	case 15:
		finalSprite = spriteRoadLRUD // lewo-prawo-góra-doł
	}

	board.Tiles[x][y].TextureID = finalSprite
}

func refreshRoadTile(x, y int, board *boardData) {
	if x < 0 || x >= int(boardMaxX) || y < 0 || y >= int(boardMaxY) {
		return
	}

	if isDirtRoad(board.Tiles[x][y].TextureID) {
		applyRoadProcessing(uint8(x), uint8(y), board)
	}
}

func applyPalisadeProcessing(x, y uint8, board *boardData) {
	var mask uint8

	if y > 0 && isPalisade(board.Tiles[x][y-1].TextureID) {
		mask |= 1
	}

	if x < boardMaxX-1 && isPalisade(board.Tiles[x+1][y].TextureID) {
		mask |= 2
	}

	if y < boardMaxY-1 && isPalisade(board.Tiles[x][y+1].TextureID) {
		mask |= 4
	}

	if x > 0 && isPalisade(board.Tiles[x-1][y].TextureID) {
		mask |= 8
	}

	// Mapujemy maskę na konkretne ID z assets_ids.go lub offset
	newID := spritePalisadeStart

	switch mask {
	case 0, 3:
		newID = spritePalisadeNE // 266
	case 1, 4, 5:
		newID = spritePalisadeV // 267
	case 2, 8, 10:
		newID = spritePalisadeH // 271 (266+5)
	case 6:
		newID = spritePalisadeStart + 2 // 268
	case 7:
		newID = spritePalisadeStart + 3 // 269
	case 9:
		newID = spritePalisadeStart + 4 // 270
	case 11:
		newID = spritePalisadeStart + 6 // 272
	case 12:
		newID = spritePalisadeStart + 7 // 273
	case 13:
		newID = spritePalisadeStart + 8 // 274
	case 14:
		newID = spritePalisadeStart + 9 // 275
	case 15:
		newID = spritePalisadeStart + 10 // 276
	}

	board.Tiles[x][y].TextureID = newID
}

func processMapTiles(bState *battleState) {
	// Tworzymy snapshot TextureID, aby uniknąć efektów ubocznych podczas iteracji
	var snapshot [boardMaxX][boardMaxY]uint16

	for x := uint8(0); x < boardMaxX; x++ {
		for y := uint8(0); y < boardMaxY; y++ {
			snapshot[x][y] = bState.Board.Tiles[x][y].TextureID
		}
	}

	for y := uint8(0); y < boardMaxY; y++ {
		for x := uint8(0); x < boardMaxX; x++ {
			id := bState.Board.Tiles[x][y].TextureID

			switch {
			case isDirtRoad(id):
				applyRoadProcessing(x, y, bState.Board)
			case isPalisade(id):
				applyPalisadeProcessing(x, y, bState.Board)
			case isWaterTileOnly(id):
				applyWaterProcessing(x, y, bState.Board, snapshot)
			case isHealingShire(id):
				bState.HealingShrines = append(bState.HealingShrines, point{X: x, Y: y})
			}
		}
	}

	makeGrassVariations(bState)
}

func applyWaterProcessing(x, y uint8, board *boardData, snapshot [boardMaxX][boardMaxY]uint16) {
	currentTile := board.Tiles[x][y].TextureID
	if !isWaterTileOnly(currentTile) {
		return
	}

	var mask uint8 = 0

	if y > 0 && isWaterOrBridgeForMasking(snapshot[x][y-1]) {
		mask |= 1
	}

	if x < boardMaxX-1 && isWaterOrBridgeForMasking(snapshot[x+1][y]) {
		mask |= 2
	}

	if y < boardMaxY-1 && isWaterOrBridgeForMasking(snapshot[x][y+1]) {
		mask |= 4
	}

	if x > 0 && isWaterOrBridgeForMasking(snapshot[x-1][y]) {
		mask |= 8
	}

	newID := uint16(0)

	switch mask {
	case 0:
		newID = 60
	case 1:
		newID = spriteWaterStart + uint16(8)
	case 3:
		newID = spriteWaterStart + uint16(0)
	case 4:
		newID = spriteWaterStart + uint16(8)
	case 5:
		newID = spriteWaterStart + uint16(8)
	case 6:
		newID = spriteWaterStart + uint16(1)
	case 7:
		newID = spriteWaterStart + uint16(2)
	case 8:
		newID = spriteWaterStart + uint16(8)
	case 9:
		newID = spriteWaterStart + uint16(3)
	case 10:
		newID = spriteWaterStart + uint16(8)
	case 11:
		newID = spriteWaterStart + uint16(4)
	case 12:
		newID = spriteWaterStart + uint16(5)
	case 13:
		newID = spriteWaterStart + uint16(6)
	case 14:
		newID = spriteWaterStart + uint16(7)
	case 15:
		newID = spriteWaterStart + uint16(8)
		if x > 0 && y > 0 && isLandOrOther(snapshot[x-1][y-1]) {
			newID = spriteWaterStart + uint16(9)
		}

		if x < boardMaxX-1 && y > 0 && isLandOrOther(snapshot[x+1][y-1]) {
			newID = spriteWaterStart + uint16(10)
		}

		if x < boardMaxX-1 && y < boardMaxY-1 && isLandOrOther(snapshot[x+1][y+1]) {
			newID = spriteWaterStart + uint16(11)
		}

		if x > 0 && y < boardMaxY-1 && isLandOrOther(snapshot[x-1][y+1]) {
			newID = spriteWaterStart + uint16(12)
		}
	case 2:
		// Przypadek teoretycznie możliwy, ale nie ma dla niego tektury
		newID = 999 // @todo: ten przypadek powinien być ręcznie usunięty z każdej mapy
	}

	board.Tiles[x][y].TextureID = newID
}

func calculateWaterTileID(x, y uint8, board *boardData) uint16 {
	var mask uint8

	if isWaterOrBridgeForMasking(board.Tiles[x][y-1].TextureID) {
		mask |= 1
	}

	if x < boardMaxX-1 && isWaterOrBridgeForMasking(board.Tiles[x+1][y].TextureID) {
		mask |= 2
	}

	if y < boardMaxY-1 && isWaterOrBridgeForMasking(board.Tiles[x][y+1].TextureID) {
		mask |= 4
	}

	if isWaterOrBridgeForMasking(board.Tiles[x-1][y].TextureID) {
		mask |= 8
	}

	newID := uint16(0)

	switch mask {
	case 0:
		newID = 60
	case 1:
		newID = spriteWaterStart + uint16(8)
	case 3:
		newID = spriteWaterStart + uint16(0)
	case 4:
		newID = spriteWaterStart + uint16(8)
	case 5:
		newID = spriteWaterStart + uint16(8)
	case 6:
		newID = spriteWaterStart + uint16(1)
	case 7:
		newID = spriteWaterStart + uint16(2)
	case 8:
		newID = spriteWaterStart + uint16(8)
	case 9:
		newID = spriteWaterStart + uint16(3)
	case 10:
		newID = spriteWaterStart + uint16(8)
	case 11:
		newID = spriteWaterStart + uint16(4)
	case 12:
		newID = spriteWaterStart + uint16(5)
	case 13:
		newID = spriteWaterStart + uint16(6)
	case 14:
		newID = spriteWaterStart + uint16(7)
	case 15:
		newID = spriteWaterStart + uint16(8)
		// Sprawdzenie narożników (diagonale)
		if x > 0 && y > 0 && isLandOrOther(board.Tiles[x-1][y-1].TextureID) {
			newID = spriteWaterStart + uint16(9)
		}

		if x < boardMaxX-1 && y > 0 && isLandOrOther(board.Tiles[x+1][y-1].TextureID) {
			newID = spriteWaterStart + uint16(10)
		}

		if x < boardMaxX-1 && y < boardMaxY-1 && isLandOrOther(board.Tiles[x+1][y+1].TextureID) {
			newID = spriteWaterStart + uint16(11)
		}

		if x > 0 && y < boardMaxY-1 && isLandOrOther(board.Tiles[x-1][y+1].TextureID) {
			newID = spriteWaterStart + uint16(12)
		}
	case 2:
		// Przypadek teoretycznie możliwy, ale nie ma dla niego tekstury
		newID = 999 // @todo: ten przypadek powinien być ręcznie usunięty z każdej mapy
	}

	return newID
}

func drawSprite(assets *assetManager, id uint16, destX, destY float32, ownerColor uint8) {
	// 1. Walidacja ID
	if id >= maxSpriteID {
		return
	}

	// 2. Pobranie definicji
	def := spriteRegistry[id]
	if def.cropWidth == 0 {
		return // Pusty wpis (np. ID 0)
	}

	// 3. Pobranie Tekstury
	tex := assets.getAtlas(def.atlasID, ownerColor)
	if tex.ID == 0 {
		return // Zabezpieczenie przed brakiem atlasu
	}

	// 4. Obliczenie Prostokąta Źródłowego
	srcW := float32(def.cropWidth)

	if def.flipX {
		srcW = -srcW // Raylib obsługuje zwierciadlane odbicie przez ujemną szerokość
	}

	srcRect := rl.NewRectangle(float32(def.cropX), float32(def.cropY), srcW, float32(def.cropHeight))

	// 5. Obliczenie Pozycji Docelowej z uwzględnieniem offsetów
	finalX := destX + float32(def.offX)
	finalY := destY + float32(def.offY)

	destRect := rl.NewRectangle(finalX, finalY, float32(def.cropWidth), float32(def.cropHeight))

	// 6. Rysowanie
	rl.DrawTexturePro(tex, srcRect, destRect, rl.NewVector2(0, 0), 0, rl.White)
}

func drawSpriteEx(id uint16, destX, destY float32, ownerColor uint8, tint rl.Color, ps *programState) {
	if id >= maxSpriteID {
		return
	}

	def := spriteRegistry[id]

	if def.cropWidth == 0 {
		return // Pusty wpis
	}

	// Pobieramy atlas z managera
	tex := ps.Assets.getAtlas(def.atlasID, ownerColor)

	// Logika flipX i Offsetów z definicji
	srcW := float32(def.cropWidth)

	if def.flipX {
		srcW = -srcW
	}

	srcRect := rl.NewRectangle(float32(def.cropX), float32(def.cropY), srcW, float32(def.cropHeight))

	// Stosujemy offsety z SpriteDef
	finalX := destX + float32(def.offX)
	finalY := destY + float32(def.offY)

	destRect := rl.NewRectangle(finalX, finalY, float32(def.cropWidth), float32(def.cropHeight))

	rl.DrawTexturePro(tex, srcRect, destRect, rl.NewVector2(0, 0), 0, tint)
}

func drawFrameCorners(x, y, width, height, cLen float32, frameColor rl.Color) {
	// lewy góry narożnik
	rl.DrawLineEx(rl.NewVector2(x, y), rl.NewVector2(x+cLen, y), cornerThickness, frameColor)
	rl.DrawLineEx(rl.NewVector2(x, y), rl.NewVector2(x, cLen+y), cornerThickness, frameColor)
	// prawy górny narożnik
	rl.DrawLineEx(rl.NewVector2(x+width, y), rl.NewVector2(x+width-cLen, y), cornerThickness, frameColor)
	rl.DrawLineEx(rl.NewVector2(x+width, y), rl.NewVector2(x+width, cLen+y), cornerThickness, frameColor)
	// lewy dolny narożnik
	rl.DrawLineEx(rl.NewVector2(x, y+height), rl.NewVector2(x+cLen, y+height), cornerThickness, frameColor)
	rl.DrawLineEx(rl.NewVector2(x, y+height), rl.NewVector2(x, y+height-cLen), cornerThickness, frameColor)
	// prawy dolny narożnik
	rl.DrawLineEx(rl.NewVector2(x+width, y+height), rl.NewVector2(x+width-cLen, y+height), cornerThickness, frameColor)
	rl.DrawLineEx(rl.NewVector2(x+width, y+height), rl.NewVector2(x+width, y+height-cLen), cornerThickness, frameColor)
}

func (bld *building) occupiedTilesBounds() (minX, minY, maxX, maxY uint8) {
	minX, minY = bld.OccupiedTiles[0].X, bld.OccupiedTiles[0].Y
	maxX, maxY = minX, minY

	for _, tile := range bld.OccupiedTiles[1:] {
		tileX, tileY := tile.X, tile.Y
		if tileX < minX {
			minX = tileX
		}

		if tileX > maxX {
			maxX = tileX
		}

		if tileY < minY {
			minY = tileY
		}

		if tileY > maxY {
			maxY = tileY
		}
	}

	return minX, minY, maxX, maxY
}

func drawBuildingSelectionFrame(bld *building, bounds bounds, bState *battleState) {
	var cLen float32

	if bld.Type != buildingPalisade {
		cLen = cornerLenBuilding
	} else {
		cLen = cornerLenPalisade
	}

	if bld.Owner == bState.PlayerID {
		drawFrameCorners(float32(bounds.X), float32(bounds.Y), bounds.WidthPx, bounds.HeightPx, cLen, friendlyFrameColor)
	} else {
		drawFrameCorners(float32(bounds.X), float32(bounds.Y), bounds.WidthPx, bounds.HeightPx, cLen, enemyFrameColor)
	}
}

func drawBuildingCapacity(bld *building, bounds bounds) {
	if bld.MaxFood == 0 {
		return
	}

	baseX := bounds.X + capacityMarginX
	baseY := bounds.Y + bounds.Height - capacityReactH - capacityMarginY
	step := capacityRectW + capacityMarginX

	for i := range bld.MaxFood {
		x := baseX + int32(i)*step
		color := rl.DarkGray

		if i < bld.Food {
			color = rl.Gold
		}

		rl.DrawRectangle(x, baseY, capacityRectW, capacityReactH, color)
	}
}

func drawBuildingHealthBar(bld *building, bounds bounds) {
	barY := bounds.Y + buildingHealthBarMarginY
	rl.DrawRectangle(bounds.X, barY, bounds.Width, healthBarHeight, rl.Red)

	fillWidth := int32(float32(bounds.Width) * (float32(bld.HP) / float32(bld.MaxHP)))
	rl.DrawRectangle(bounds.X, barY, fillWidth, healthBarHeight, rl.Green)
}

func drawBuildingInterface(bld *building, bState *battleState) {
	buildingBounds := bld.bounds()
	isBuildingSelected := bState.CurrentSelection.BuildingID == bld.ID

	if isBuildingSelected {
		drawBuildingSelectionFrame(bld, buildingBounds, bState)

		if bld.Owner == bState.PlayerID && !bld.IsUnderConstruction {
			drawBuildingCapacity(bld, buildingBounds)
		}
	}

	// Rysowanie paska życia
	if (bld.HP < bld.MaxHP || isBuildingSelected) && bld.HP != 0 {
		drawBuildingHealthBar(bld, buildingBounds)
	}
}

func drawBuildingsInterfaces(bState *battleState) {
	for _, currentBuilding := range bState.Buildings {
		if !currentBuilding.Exists || len(currentBuilding.OccupiedTiles) == 0 {
			continue
		}

		drawBuildingInterface(currentBuilding, bState)
	}
}

// === RYSOWANIE JEDNOSTEK I POCISKÓW ===

func calculateLegacyPhase(u *unit) uint8 {
	idx := getLegacyUnitIndex(u.Type)
	delay := u.Delay

	delay = max(delay, minPhaseDelay)
	delay = min(delay, maxPhaseDelay)
	phase := legacyPhase[idx][delay]

	if u.AnimationType == "fight" {
		triggerHit := uint16(legacyPhase[idx][17])
		triggerReturn := uint16(legacyPhase[idx][18])

		if delay > triggerHit {
			phase = 3
		}

		if delay <= triggerHit && delay > triggerReturn {
			phase = 4
		}

		if delay <= triggerReturn {
			phase = 1
			if u.Type == unitCow {
				phase = 0
			}
		}
	}

	return phase
}

func getRenderDirection(u *unit, bState *battleState) (int, int) {
	if u.AnimationType == "fight" {
		var targetX, targetY uint8
		foundTarget := false

		// 1. Sprawdź, czy mamy cel ataku (Unit lub Building)
		if u.TargetID != 0 {
			targetUnit, targetBld := getObjectByID(u.TargetID, bState)

			if targetUnit != nil && targetUnit.Exists {
				targetX = targetUnit.X
				targetY = targetUnit.Y
				foundTarget = true // WAŻNE: Musimy to oznaczyć!
			} else if targetBld != nil && targetBld.Exists {
				// Tutaj funkcja zwraca bool, więc używamy go
				var ok bool
				targetX, targetY, ok = targetBld.getClosestOccupiedTile(u.X, u.Y)
				foundTarget = ok
			}
		}

		// 2. Fallback: Jeśli nie mamy celu ataku (lub zniknął), patrzymy na cel ruchu
		if !foundTarget {
			// Sprawdzamy, czy cel ruchu jest różny od obecnej pozycji
			if u.TargetX != u.X || u.TargetY != u.Y {
				targetX = u.TargetX
				targetY = u.TargetY
				foundTarget = true
			}
		}

		// 3. Obliczanie delty (kierunku) - TYLKO jeśli mamy cel
		if foundTarget {
			dx := 0
			dy := 0

			// Bezpieczne porównanie uintów
			if targetX > u.X {
				dx = 1
			} else if targetX < u.X {
				dx = -1
			}

			if targetY > u.Y {
				dy = 1
			} else if targetY < u.Y {
				dy = -1
			}

			// Zwracamy tylko jeśli faktycznie jest jakiś kierunek
			if dx != 0 || dy != 0 {
				return dx, dy
			}
		}
	}

	// Domyślny kierunek (np. z wektora ruchu plynnego)
	// Rzutujemy na int, żeby zachować znak
	dx := int(math.Round(float64(u.Direction.X)))
	dy := int(math.Round(float64(u.Direction.Y)))

	return dx, dy
}

func drawProjectiles(bState *battleState, ps *programState) {
	for _, p := range bState.Projectiles {
		if p.Exists {
			drawSingleProjectile(p, bState.GlobalFrameCounter, ps)
		}
	}
}

func drawGhostlySprite(spriteID uint16, x, y float32, phase1, phase2 float64, frameCounter uint16, ps *programState) {
	// 1. Pobranie definicji duszka z atlasu
	if spriteID >= maxSpriteID {
		return
	}

	def := spriteRegistry[spriteID]

	if def.cropWidth == 0 {
		return
	}

	texture := ps.Assets.getAtlas(def.atlasID, colorNone)

	if texture.ID == 0 {
		return
	}

	// 2. Wodotryski
	t := float64(frameCounter)
	wave := math.Sin(t*0.25+phase1) + 0.5*math.Sin(t*0.33+phase2)
	pulse := float32(0.6 + 0.4*(0.5+0.5*wave))
	wobbleX := float32(2.0 * math.Sin(t*0.4+phase1*1.1))
	wobbleY := float32(2.0 * math.Cos(t*0.4+phase1*1.1))
	alpha := 0.6 + 0.4*float32(math.Sin(t*0.35+(phase1+phase2)*0.5))
	tint := rl.Fade(rl.White, alpha)

	// 3. Zwierciadlane odbicie jeśli potrzeba
	srcW := float32(def.cropWidth)
	if def.flipX {
		srcW = -srcW
	}

	// 4. Wycięcie duszka z atlasu
	srcRect := rl.NewRectangle(float32(def.cropX), float32(def.cropY), srcW, float32(def.cropHeight))

	// 5. Pozycjonowanie
	centerX := x + float32(def.offX) + wobbleX + float32(def.cropWidth)
	centerY := y + float32(def.offY) + wobbleY + float32(def.cropHeight)

	// 6, Umiejscowienie
	destW := float32(def.cropWidth) * pulse
	destH := float32(def.cropHeight) * pulse

	destRect := rl.NewRectangle(centerX-destW*0.5, centerY-destH*0.5, destW, destH)

	origin := rl.NewVector2(destW*0.5, destH*0.5)

	// 7. Rysujemy
	rl.DrawTexturePro(texture, srcRect, destRect, origin, 0, tint)
}

func drawSingleProjectile(p *projectile, frameCounter uint16, ps *programState) {
	if p.Kind != missileGhost {
		drawSpriteEx(p.Sprite, p.X, p.Y, colorRed, rl.White, ps)
	} else {
		drawGhostlySprite(p.Sprite, p.X, p.Y, p.Phase1, p.Phase2, frameCounter, ps)
	}
}

// @todo: być może warto połączyć logikę rysowania paska HP budynków i jednostek?
func drawUnitHealthBar(screenX, screenY int32, unit *unit) {
	// Wskaźnik MaxHP jednostki
	barW := int32(tileWidth)
	barY := screenY - 3
	rl.DrawRectangle(screenX, barY, barW, healthBarHeight, rl.Red)

	// Wskaźnik bieżącego HP jednostki →
	fillWidth := int32(float32(barW) * (float32(unit.HP) / float32(unit.MaxHP)))
	rl.DrawRectangle(screenX, barY, fillWidth, healthBarHeight, rl.Green)
}

func drawManaBar(screenX, screenY int32, unit *unit) {
	barW, barH := int32(1), int32(tileHeight)
	barX := screenX + int32(tileWidth) + 1
	barY := screenY

	// Wypełnienie wskazujące na poziom many
	manaPercent := float32(unit.Mana) / float32(unit.MaxMana)
	fillH := int32(float32(barH) * manaPercent)

	// Rysujemy ↑
	newY := (barY + barH) - fillH
	rl.DrawRectangle(barX, newY, barW, fillH, rl.Blue)
}

func drawExperienceBar(screenX, screenY int32, unit *unit) {
	barW, barH := int32(1), int32(tileHeight)
	barX := screenX - 2 // nie dodaję szerokości kafelka ponieważ chcę mieć pasek po lewej
	barY := screenY

	// Wypełnienie wskazujące na poziom doświadczenia
	experiencePercent := float32(unit.Experience) / float32(experienceCap) // @reminder: sprawdź ile to maxExp!
	fillH := int32(float32(barH) * experiencePercent)

	// Rysujemy ↑
	newY := (barY + barH) - fillH
	rl.DrawRectangle(barX, newY, barW, fillH, rl.Gold)
}

func drawMilkBar(screenX, screenY int32, unit *unit) {
	barW, barH := int32(1), int32(tileHeight)
	barX := screenX + int32(tileWidth) + 1
	barY := screenY

	// Wypełnienie wskazujące na poziom many
	milkPercent := float32(unit.Udder) / float32(fullUdderAmount)
	fillH := int32(float32(barH) * milkPercent)

	// Rysujemy ↑
	newY := (barY + barH) - fillH
	rl.DrawRectangle(barX, newY, barW, fillH, rl.White)
}

func drawUnitSelectionFrame(selectedUnit *unit, bState *battleState) {
	x := float32(selectedUnit.X) * float32(tileWidth)
	y := float32(selectedUnit.Y) * float32(tileHeight)

	if selectedUnit.Owner == bState.PlayerID {
		drawFrameCorners(x, y, float32(tileWidth), float32(tileHeight), cornerLenUnit, friendlyFrameColor)
	} else {
		drawFrameCorners(x, y, float32(tileWidth), float32(tileHeight), cornerLenUnit, enemyFrameColor)
	}
}

// Odpowiada za rysowanie nakładek (ramka, pasek życia itd.) jednostkom widocznym na ekranie.
func drawUnitsInterfaces(startY, endY uint8, bState *battleState) {
	for boardRow := startY; boardRow < endY; boardRow++ {
		unitsInRow := bState.RenderUnitRows[boardRow]

		if len(unitsInRow) > 0 {
			for _, renderUnit := range unitsInRow {
				drawUnitInterface(renderUnit, bState)
			}
		}
	}
}

func drawUnitInterface(renderUnit *unit, bState *battleState) {
	screenX := int32(renderUnit.X) * int32(tileWidth)
	screenY := int32(renderUnit.Y) * int32(tileHeight)

	isUnitSelected := renderUnit.IsSelected ||
		(bState.CurrentSelection.IsUnit && bState.CurrentSelection.UnitID == renderUnit.ID)

	if isUnitSelected {
		drawUnitSelectionFrame(renderUnit, bState)
		drawUnitHealthBar(screenX, screenY, renderUnit)

		// @reminder: pracuję nad tym w tej chwili - 19.05.2026
		drawExperienceBar(screenX, screenY, renderUnit)

		if renderUnit.MaxMana > 0 {
			drawManaBar(screenX, screenY, renderUnit)
		}

		if renderUnit.Type == unitCow {
			drawMilkBar(screenX, screenY, renderUnit)
		}
	}
}

// === LOGIKA KAMERY I EKRANU ===

func setupGameCamera(bState *battleState, ps *programState) {
	bState.GameCamera = rl.NewCamera2D(
		// Offset: środek widoku gry (DYNAMICZNY)
		rl.NewVector2(
			ps.GameViewWidth/2.0,
			ps.VirtualHeight/2.0,
		),
		// Target: środek świata (lub inna pozycja startowa)
		rl.NewVector2(
			ps.GameViewWidth/2.0,
			ps.VirtualHeight/2.0,
		),
		0.0,
		1.0,
	)
}

func drawBattleScene(bState *battleState) {
	rl.BeginScissorMode(0, 0, int32(gameViewVirtualWidth), int32(virtualScreenHeight))
	rl.BeginMode2D(bState.GameCamera)
	rl.EndMode2D()
	rl.EndScissorMode()
}

func isMouseOverUI(ps *programState, virtualPos rl.Vector2) bool {
	return virtualPos.X >= ps.GameViewWidth
}

// Dobiera barwy dla mapy na podstawie zawartości kafelka.
//
//nolint:cyclop
func getMapColor(tileID uint16) rl.Color {
	switch {
	case isGrassOrGadget(tileID):
		return rl.Green
	case isRockNonWalkable(tileID):
		return rl.DarkGray
	case tileID == spriteRock12:
		return rl.LightGray
	case isTree(tileID):
		return rl.DarkGreen
	case isWaterTileOnly(tileID):
		return rl.DarkBlue
	case tileID == spriteEffectHeal00 || tileID == spriteEffectHeal01:
		return rl.Blue
	case isPalisade(tileID):
		return rl.Brown
	case isPath(tileID):
		return rl.DarkBrown
	case isRuin(tileID):
		return rl.Black
	case isDryEarth(tileID):
		return rl.Beige
	default:
		return rl.Pink
	}
}

func getCameraWorldViewRect(camera rl.Camera2D, virtualScreenWidth, virtualScreenHeight float32) rl.Rectangle {
	topLeftScreen := rl.NewVector2(0, 0)
	bottomRightScreen := rl.NewVector2(virtualScreenWidth, virtualScreenHeight)
	topLeftWorld := rl.GetScreenToWorld2D(topLeftScreen, camera)
	bottomRightWorld := rl.GetScreenToWorld2D(bottomRightScreen, camera)
	worldWidth := bottomRightWorld.X - topLeftWorld.X
	worldHeight := bottomRightWorld.Y - topLeftWorld.Y

	return rl.NewRectangle(topLeftWorld.X, topLeftWorld.Y, worldWidth, worldHeight)
}

func drawSoil(startX, startY, endX, endY uint8, bState *battleState, ps *programState) {
	for yAxis := startY; yAxis < endY; yAxis++ {
		for xAxis := startX; xAxis < endX; xAxis++ {
			currentTile := &bState.Board.Tiles[xAxis][yAxis]
			texID := currentTile.TextureID
			xPos := float32(xAxis) * float32(tileWidth)
			yPos := float32(yAxis) * float32(tileHeight)

			if isDryEarth(texID) {
				drawSprite(ps.Assets, spriteGrass00, xPos, yPos, colorNone)

				continue
			}

			if currentTile.IsGrazed {
				drawSprite(ps.Assets, texID, xPos, yPos, colorNone)
				drawSprite(ps.Assets, uint16(currentTile.GrazedOverlayID), xPos, yPos, colorNone)

				continue
			}

			// Rysowanie trawy pod kapliczkami, miejscem zwycięstwa oraz drzewami.
			if isSpecialTile(texID) || currentTile.isTree() || currentTile.isFallingTree() {
				drawSprite(ps.Assets, spriteGrass00, xPos, yPos, colorNone)
			}

			if isCompletedBridge(texID) || texID == spriteBridgeConstruction {
				waterBaseID := calculateWaterTileID(xAxis, yAxis, bState.Board)

				if waterBaseID == 999 {
					waterBaseID = spriteWaterMiddle
				}

				animationOffset := bState.WaterAnimationFrame * 13
				drawSprite(ps.Assets, waterBaseID+animationOffset, xPos, yPos, colorNone)
				drawSprite(ps.Assets, texID, xPos, yPos, colorNone)
			}

			if isWaterTileOnly(texID) {
				animationOffset := bState.WaterAnimationFrame * 13
				drawSprite(ps.Assets, texID+animationOffset, xPos, yPos, colorNone)
			}

			if isLandOrOther(texID) && !currentTile.isTree() {
				drawSprite(ps.Assets, texID, xPos, yPos, colorNone)
			}
		}
	}
}

// @todo: przemyśl, czy musi to być tak być, bo wygląda tragicznie
func drawBuilding(startY, endY uint8, bState *battleState, ps *programState) {
	for yAxis := startY; yAxis < endY; yAxis++ {
		for _, bld := range bState.RenderBuildingRows[yAxis] {
			for _, currentTile := range bld.OccupiedTiles {
				if currentTile.Y == yAxis {
					id := bState.Board.Tiles[currentTile.X][yAxis].TextureID
					if isBuildingTerrain(id) {
						finalID := id
						if flagID, ok := flagAnimationMap[uint8(id)]; ok && (bState.FireAnimationFrame+uint16(currentTile.X)+uint16(yAxis))%2 == 1 {
							finalID = uint16(flagID)
						}

						drawSprite(ps.Assets, finalID, float32(currentTile.X)*float32(tileWidth), float32(yAxis)*float32(tileHeight), bld.Owner)
					}
				}
			}
		}
	}
}

func drawPalisade(startX, startY, endX, endY uint8, bState *battleState, ps *programState) {
	for yAxis := startY; yAxis < endY; yAxis++ {
		for xAxis := startX; xAxis < endX; xAxis++ {
			textureID := bState.Board.Tiles[xAxis][yAxis].TextureID
			if isPalisade(textureID) {
				drawSprite(ps.Assets, textureID, float32(xAxis)*float32(tileWidth), float32(yAxis)*float32(tileHeight), colorNone)
			}
		}
	}
}

func drawDryEarth(startX, startY, endX, endY uint8, bState *battleState, ps *programState) {
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			textureID := bState.Board.Tiles[x][y].TextureID
			if isDryEarth(textureID) {
				drawSprite(ps.Assets, textureID, float32(x)*float32(tileWidth), float32(y)*float32(tileHeight), colorNone)
			}
		}
	}
}

func drawBuildings(startX, startY, endX, endY uint8, bState *battleState, ps *programState) {
	drawBuilding(startY, endY, bState, ps)
	drawPalisade(startX, startY, endX, endY, bState, ps)
	drawDryEarth(startX, startY, endX, endY, bState, ps)
}

func drawCorpsesUnitsTrees(startX, startY, endX, endY uint8, bState *battleState, ps *programState) {
	for boardRow := startY; boardRow < endY; boardRow++ {
		drawCorpses(boardRow, startX, endX, bState, ps)
	}

	for boardRow := startY; boardRow < endY; boardRow++ {
		drawUnits(boardRow, bState, ps)
		drawTrees(boardRow, startX, endX, bState, ps)
	}
}

func drawUnits(boardRow uint8, bState *battleState, ps *programState) {
	rowUnits := bState.RenderUnitRows[boardRow]
	if len(rowUnits) > 0 {
		sort.Slice(rowUnits, func(i, j int) bool { return rowUnits[i].X < rowUnits[j].X })

		for _, currentUnit := range rowUnits {
			if currentUnit.Owner == bState.PlayerID {
				drawUnit(currentUnit, bState, ps)
			}
		}

		for _, currentUnit := range rowUnits {
			if currentUnit.Owner != bState.PlayerID {
				drawUnit(currentUnit, bState, ps)
			}
		}
	}
}

func drawCorpses(y, startX, endX uint8, bState *battleState, ps *programState) {
	for _, currentCorpse := range bState.CorpsesList {
		if currentCorpse.Y == y && currentCorpse.X >= startX && currentCorpse.X < endX {
			posX := float32(currentCorpse.X) * float32(tileWidth)
			posY := float32(currentCorpse.Y) * float32(tileHeight)

			if currentCorpse.Phase < corpsesPhase2 {
				var offsetIndex uint16

				if currentCorpse.Phase == corpsesPhase1 {
					offsetIndex = 1
				}

				baseID := spriteUnitBaseID + (uint16(currentCorpse.UnitType) * spriteUnitStep)
				finalID := baseID + corpsesFrameIndexOffset + offsetIndex
				drawSprite(ps.Assets, finalID, posX, posY, currentCorpse.Owner)
			} else {
				cid := spriteEffectskeleton00 + uint16(currentCorpse.SkeletonType)
				tint := rl.Fade(rl.White, float32(currentCorpse.Alpha)/corpsesMaxAlpha)
				drawSpriteEx(cid, posX, posY, colorNone, tint, ps)
			}
		}
	}
}

func drawTrees(boardRow, startX, endX uint8, bState *battleState, ps *programState) {
	for boardColumn := startX; boardColumn < endX; boardColumn++ {
		currentTile := &bState.Board.Tiles[boardColumn][boardRow]

		// jeśli nie jest drzewem to pomijamy
		if currentTile.treeState == noTree {
			continue
		}

		switch currentTile.treeState {
		case noTree:
			continue
		case treeStraight:
			drawStandingTree(boardColumn, boardRow, currentTile, bState, ps)
		case treeLeaning:
			drawLeaningTree(boardColumn, boardRow, currentTile, ps)
		case treeFalling, treeImpact, treeFell:
			drawFallingOrFallenTree(boardColumn, boardRow, currentTile, ps)
		}
	}
}

func drawStandingTree(boardColumn, boardRow uint8, t *tile, bState *battleState, ps *programState) {
	drawX := float32(boardColumn) * float32(tileWidth)
	drawY := float32(boardRow) * float32(tileHeight)

	if !t.IsBurnt {
		drawX -= treeOffsetX
	}

	// 1. Pień
	drawSprite(ps.Assets, t.TextureID, drawX, drawY, colorNone)

	// 2. Korona
	if boardRow > 0 {
		crownY := float32(boardRow-1) * float32(tileHeight)

		var crownID uint16

		if t.IsBurnt {
			crownID = t.TextureID + burntTreeCrownTextureOffset
		} else {
			crownID = t.TextureID + treeCrownTextureOffset
		}

		drawSprite(ps.Assets, crownID, drawX, crownY, colorNone)

		// 3. Uruchomienie ognia
		if t.IsBurning {
			frame := (bState.FireAnimationFrame + uint16(boardColumn+boardRow)) % 4
			drawSprite(ps.Assets, spriteFire00+frame, drawX+treeOffsetX, crownY, colorNone)
		}
	}
}

func drawLeaningTree(boardColumn, boardRow uint8, t *tile, ps *programState) {
	drawX := float32(boardColumn) * float32(tileWidth)
	drawY := float32(boardRow) * float32(tileHeight)

	var baseStump uint16

	if t.IsBurnt {
		baseStump = spriteBurntLeaningTreeStump
	} else {
		baseStump = spriteDryLeaningTreeStump
	}

	drawSprite(ps.Assets, baseStump, drawX, drawY, colorNone)
	drawSprite(ps.Assets, baseStump+1, drawX, drawY, colorNone)
	drawSprite(ps.Assets, baseStump+2, drawX, drawY, colorNone)
}

func drawFallingOrFallenTree(boardColumn, boardRow uint8, t *tile, ps *programState) {
	drawX := float32(boardColumn) * float32(tileWidth)
	drawY := float32(boardRow) * float32(tileHeight)

	var stumpID uint16

	// Tutaj nie wykorzystujemy wszystkich stanów ponieważ interesują nas tylko
	// i wyłącznie upadające drzewa
	switch t.treeState {
	case treeFalling:
		if t.IsBurnt {
			stumpID = spriteBurntFallingTreeStump
		} else {
			stumpID = spriteDryFallingStump
		}
	case treeImpact, treeFell:
		if t.IsBurnt {
			stumpID = spriteBurntFallenTreeStump
		} else {
			stumpID = spriteDryFallenTreeStump
		}
	}

	drawSprite(ps.Assets, stumpID, drawX, drawY, colorNone)
	drawSprite(ps.Assets, stumpID+1, drawX, drawY, colorNone)
}

// @todo: porozbijać na podfunkcje! Inaczej nie rozprawię się z tym potworem.
func drawWorldAndUnits(bState *battleState, ps *programState) {
	// 1. Aktualizacja podręcznych
	// @todo: Przydałoby się dokładniej opisać, co jest co, bo ciężko się rozeznać - 24.04.2026
	bState.updateRenderCache()
	cam := bState.GameCamera

	// @todo: Sprawdź, czy mam rację: world… to część planszy, którą widzimy
	worldLeft := cam.Target.X - (cam.Offset.X / cam.Zoom)
	worldTop := cam.Target.Y - (cam.Offset.Y / cam.Zoom)
	worldRight := worldLeft + (ps.GameViewWidth / cam.Zoom)
	worldBottom := worldTop + (ps.VirtualHeight / cam.Zoom)

	// @todo: Sprawdź, czy mam rację: start… to część planszy, którą widzimy przeliczona na kafelki
	startXInt := int(worldLeft/float32(tileWidth)) - 1
	startYInt := int(worldTop/float32(tileHeight)) - 1
	endXInt := int(worldRight/float32(tileWidth)) + 2
	endYInt := int(worldBottom/float32(tileHeight)) + 2

	// Zabezpieczenie przed wyjściem poza początek planszy
	if startXInt < 0 {
		startXInt = 0
	}

	if startYInt < 0 {
		startYInt = 0
	}

	// Zabezpieczenie przed wyjściem poza koniec planszy
	if endXInt > int(boardMaxX) {
		endXInt = int(boardMaxX)
	}

	if endYInt > int(boardMaxY) {
		endYInt = int(boardMaxY)
	}

	// Konwersja bezpiecznych zakresów na uint8 dla funkcji rysujących
	startX := uint8(startXInt)
	startY := uint8(startYInt)
	endX := uint8(endXInt)
	endY := uint8(endYInt)

	// Przebieg 1: „płaskie” rzeczy, jak trawa, woda, mosty, sucha ziemia
	drawSoil(startX, startY, endX, endY, bState, ps)
	// Przebieg 2: budynki, sucha ziemia, palisady
	drawBuildings(startX, startY, endX, endY, bState, ps)
	// Przebieg 3: zwłoki, jednostki, drzewa
	drawCorpsesUnitsTrees(startX, startY, endX, endY, bState, ps)

	// Przebieg 4: efekty oraz nakładki
	drawEffectsAndInterfaces(startY, endY, bState, ps)

	// @reminder: wydaje mi się, że pociski powinny być pod nakładkami.
	// inaczej będą przesłaniać pasek życia.
	drawProjectiles(bState, ps)
	drawSelectionBox(bState, ps)
}

// @todo: rozdziel logikę od rysowania, bo to zwyczajne pomieszanie z poplątaniem.
func drawUnit(u *unit, bState *battleState, ps *programState) {
	// 1. Bazowa pozycja to CEL (bo w units.go X,Y zmieniają się na początku kroku)
	screenX := float32(u.X) * float32(tileWidth)
	screenY := float32(u.Y) * float32(tileHeight)

	// 2. Logika ruchu - Cofanie od celu

	if u.State == stateMoving {
		idx := getLegacyUnitIndex(u.Type)
		delayIdx := u.Delay

		// Zabezpieczenie zakresu
		if delayIdx > maxPhaseDelay {
			delayIdx = 16
		}

		if delayIdx < minPhaseDelay {
			delayIdx = 0
		}

		// Pobieramy wartość przesunięcia (0..16)
		rawShiftX := float32(legacyShiftX[idx][delayIdx])
		rawShiftY := float32(legacyShiftY[idx][delayIdx])

		// Fizyczny kierunek ruchu (directionX, directionY)
		directionX := int(u.Direction.X)
		directionY := int(u.Direction.Y)

		// LOGIKA: Jesteśmy na CELU. Musimy cofnąć się o 'rawShift' w stronę STARTU.
		// Start (Delay=Max, Shift=16): Cel - 16 = Start. (Jednostka w miejscu startu)
		// Koniec (Delay=0, Shift=0):   Cel - 0  = Cel.   (Jednostka w miejscu celu)

		if directionX > 0 {
			screenX -= rawShiftX // Idzie w prawo, więc cofamy w lewo
		} else if directionX < 0 {
			screenX += rawShiftX // Idzie w lewo, więc cofamy w prawo
		}

		if directionY > 0 {
			screenY -= rawShiftY // Idzie w dół, więc cofamy w górę
		} else if directionY < 0 {
			screenY += rawShiftY // Idzie w górę, więc cofamy w dół
		}
	}

	// 3. Wyliczenie klatki animacji (Phase)
	frame := calculateLegacyPhase(u)

	// 4. Kierunek duszka
	// W ruchu wymuszamy patrzenie w kierunku idzenia (naprawia "moonwalk")
	renderDx, renderDy := getRenderDirection(u, bState)
	if u.State == stateMoving {
		renderDx = int(u.Direction.X)
		renderDy = int(u.Direction.Y)
	}

	dir := legacyDirToNewDir(renderDx, renderDy)

	// 5. Korekty dla jednostek Melee (zgodnie z MOVER1.CPP)
	isMelee := u.Type == unitAxeman || u.Type == unitSwordsman ||
		u.Type == unitCommander || u.Type == unitBear || u.Type == unitUnknown

	if u.State == stateAttacking && frame > 2 && isMelee {
		if renderDy > 0 {
			screenY -= 7.0
		}

		if renderDx > 0 {
			screenX -= 8.0
		}
	}

	// 6. Rysowanie
	baseID := uint16(700 + (int(u.Type) * 200))
	finalID := baseID + uint16(frame*8) + uint16(dir)

	drawSpriteEx(finalID, screenX, screenY, u.Owner, rl.White, ps)

	if len(u.Wounds) > 0 {
		drawUnitWounds(u, screenX, screenY, ps)
	}

	// 7. Rysowanie magicznej tarczy dla czarodziejek
	drawActiveMagicShield(u, screenX, screenY, bState, ps)
}

func drawActiveMagicShield(u *unit, screenX, screenY float32, bState *battleState, ps *programState) {
	if u.Type == unitPriestess && u.hasMagicShield {
		// 0. Podkładka.
		centerX := int32(screenX) + int32(tileWidth/2)  //nolint:mnd
		centerY := int32(screenY) + int32(tileHeight/2) //nolint:mnd
		radius := float32(6.0)                          //nolint:mnd
		shieldGlow := rl.Fade(rl.Violet, 0.1)           //nolint:mnd

		rl.DrawCircle(centerX, centerY, radius, shieldGlow)

		// 1. Rysowanie właściwej tarczy.
		// @reminder: dzielę przez cztery, bo tyle klatek ma uruchomienie magicznej tarczy.
		frame := bState.FireAnimationFrame % 4 //nolint:mnd
		drawSprite(ps.Assets, spriteMagicShield00+frame, screenX, screenY, colorNone)
	}
}

// @todo: koniecznie wrócić do tej funkcji i ją przepisz, rozbij, bo jest dramat!
func drawUnitWounds(u *unit, screenX, screenY float32, ps *programState) {
	// Stałe offsety, aby krew była na środku kafelka 16x14.
	const centerOffsetX = 8.0

	const centerOffsetY = 7.0

	for _, currentWound := range u.Wounds {
		// Wybór tekstury
		effectID := spriteEffectHit00
		if currentWound.IsSevere && currentWound.Timer <= 10 {
			effectID = spriteEffectHit01
		}

		def := spriteRegistry[effectID]
		if def.cropWidth == 0 {
			continue
		}

		tex := ps.Assets.getAtlas(def.atlasID, colorNone)

		sourceRec := rl.NewRectangle(
			float32(def.cropX), float32(def.cropY),
			float32(def.cropWidth), float32(def.cropHeight),
		)

		// Skalowanie rany
		destW := float32(def.cropWidth) * currentWound.Scale
		destH := float32(def.cropHeight) * currentWound.Scale

		// Pozycja finalna: pozycja jednostki + środek kafelka + losowy offset rany
		destRect := rl.NewRectangle(
			screenX+centerOffsetX+currentWound.OffsetX,
			screenY+centerOffsetY+currentWound.OffsetY,
			destW,
			destH,
		)

		origin := rl.NewVector2(destW/2, destH/2)

		rl.DrawTexturePro(tex, sourceRec, destRect, origin, currentWound.Rotation, rl.White)
	}
}

// Konwerter kierunków (dx, dy) -> 0..7 (zgodnie z assets_db)
// @todo: to powinno się chyba później wywalić!
func legacyDirToNewDir(dx, dy int) int {
	if dx == 0 && dy == -1 {
		return 0
	} // Góra
	if dx == 1 && dy == -1 {
		return 1
	} // GP
	if dx == 1 && dy == 0 {
		return 2
	} // Prawo
	if dx == 1 && dy == 1 {
		return 3
	} // DP
	if dx == 0 && dy == 1 {
		return 4
	} // Dół
	if dx == -1 && dy == 1 {
		return 5
	} // DL
	if dx == -1 && dy == 0 {
		return 6
	} // Lewo
	if dx == -1 && dy == -1 {
		return 7
	} // GL
	return 4 // Domyślnie dół
}

func (bState *battleState) updateRenderCache() {
	// @reminder Być może kiedyś okaże się, iż czyszczenie całości nie jest właściwe.
	// Wtedy będę musiał zastanowić się, jak wywalać poszczególne jednostki i budynki.
	// 1. Czyścimy jednostki oraz budynki po poprzedniej klatce
	for y := range boardMaxY {
		bState.RenderUnitRows[y] = bState.RenderUnitRows[y][:0]
		bState.RenderBuildingRows[y] = bState.RenderBuildingRows[y][:0]
	}
	// 2. Dodajemy jednostki
	for _, currentUnit := range bState.Units {
		// Uwzględniamy tylko żywe jednostki
		// Nie da się mieć jednostki poza planszą, więc nie muszę sprawdzać
		// MAX_X i MAX_Y, pewnie kiedyś tego pożałuję
		// @check nie sprawdzam MAX_Y oraz MAX_Y, w razie problemów z rysowaniem
		// tutaj może się kryć przyczyna. Na przyszłość: nie <= bo pierwszy index to 0
		if currentUnit.Exists { // && unit.y >= 0 && unit.y < MAX_Y
			bState.RenderUnitRows[currentUnit.Y] = append(bState.RenderUnitRows[currentUnit.Y], currentUnit)
		}
	}
	// 3. Dodajemy budynki
	for _, bld := range bState.Buildings {
		// Pomijamy palisady i zniszczone
		if !bld.Exists || bld.Type == buildingPalisade {
			continue
		}

		// Ponieważ budynki trzymają informacje o 9 kafelkach, to trzeba się zabezpieczyć
		lastY := uint8(0)
		for _, currentTile := range bld.OccupiedTiles {
			if currentTile.Y != lastY {
				// if tile.y >= 0 && tile.y < MAX_Y // @reminder Nie sprawdzam legalności!
				bState.RenderBuildingRows[currentTile.Y] = append(bState.RenderBuildingRows[currentTile.Y], bld)
			}

			lastY = currentTile.Y
		}
	}
}

// drawSelectionBox odpowiada za rysowanie prostokąta do zaznaczania jednostek.
func drawSelectionBox(bState *battleState, ps *programState) {
	if !bState.IsSelectingBox {
		return
	}

	// 1. Pobieramy obecną mysz i zamieniamy ją na WIRTUALNĄ (bo takiej używa logika gry)
	// Bez tego, na dużym monitorze kursor "ucieknie" od ramki.
	currentScreenMouse := rl.GetMousePosition()
	currentVirtualMouse := screenToVirtualCoords(ps, currentScreenMouse)

	// 2. Konwertujemy OBA punkty (Start i Koniec) z Wirtualnego Ekranu na Świat Gry
	// Dzięki temu, że jesteśmy w BeginMode2D, rysujemy wprost na mapie.

	// Start zaznaczania (zapamiętany w input.go jako wirtualny ekran) → Świat
	worldStart := rl.GetScreenToWorld2D(bState.SelectionStart, bState.GameCamera)

	// Obecna mysz (wirtualny ekran) → Świat
	worldEnd := rl.GetScreenToWorld2D(currentVirtualMouse, bState.GameCamera)

	// 3. Obliczamy wymiary w Świecie Gry
	x := float32(math.Min(float64(worldStart.X), float64(worldEnd.X)))
	y := float32(math.Min(float64(worldStart.Y), float64(worldEnd.Y)))
	width := float32(math.Abs(float64(worldEnd.X - worldStart.X)))
	height := float32(math.Abs(float64(worldEnd.Y - worldStart.Y)))

	rect := rl.NewRectangle(x, y, width, height)

	// Linia musi mieć grubość niezależną od przybliżenia, inaczej może stać się niewidoczna
	const minScreenThicknessPixels = 1.5

	requiredWorldThickness := minScreenThicknessPixels / bState.GameCamera.Zoom

	finalThickness := float32(math.Max(1.0, float64(requiredWorldThickness)))

	// Rysowanie
	rl.DrawRectangleLinesEx(rect, finalThickness, rl.White) // Zielona ramka, klasyk RTS
}

func drawEffectsAndInterfaces(startY, endY uint8, bState *battleState, ps *programState) {
	drawEffects(bState, ps)
	drawBuildingsInterfaces(bState)
	drawUnitsInterfaces(startY, endY, bState)
}

// @reminder: Chodzenie po całej planszy i sprawdzanie dwukrotnie, czy mamy jakiś efekt do narysowania jest bardzo
// niewydajny! Powinienem mieć jakąś podręczną listę na stałe efekty, jak święte miejsce oraz
// na tymczasowe. Coś, jak ze zwłokami, poręczne szybkie i w ogóle.
// @todo: zmień tak, żeby nie sprawdzać każdego kafelka. Efekty nie są niewiadomą.
func drawEffects(bState *battleState, ps *programState) {
	// Czemu y jest zewnętrzną pętlą?!
	for y := range boardMaxY {
		for x := range boardMaxX {
			affectedTile := &bState.Board.Tiles[x][y]

			xPos := float32(x) * float32(tileWidth)
			yPos := float32(y) * float32(tileHeight)

			// Stałe, u dołu
			permanentEffects(affectedTile, x, y, xPos, yPos, bState, ps)

			// Tymczasowe, na wierzchu
			temporaryEffects(affectedTile, x, y, xPos, yPos, bState, ps)
		}
	}
}

func permanentEffects(affectedTile *tile, x, y uint8, xPos, yPos float32, bState *battleState, ps *programState) {
	textureID := affectedTile.TextureID

	switch {
	// A. Kapliczka leczenia (stare 282/283)
	case textureID == spriteEffectHeal00 || textureID == spriteEffectHeal01:
		frame := (bState.WaterAnimationFrame + uint16(x+y)) % 2
		animID := spriteEffectHeal00 + frame

		// Logika przeźroczystości (Duch)
		var tint rl.Color
		if affectedTile.Unit != nil {
			tint = rl.White // Pełna widoczność, gdy ktoś stoi
		} else {
			tint = rl.Fade(rl.White, 0.4) // Duch, gdy pusto
		}
		drawSpriteEx(animID, xPos, yPos, colorNone, tint, ps)

	// B. Małe ognisko (stare 68 -> SPRITE_GADGET_14)
	case textureID == spriteGadget14:
		frame := (bState.FireAnimationFrame + uint16(x+y)) % 4
		drawSprite(ps.Assets, spriteGadget14, xPos, yPos, colorNone)
		drawSprite(ps.Assets, spriteFire08+frame, xPos, yPos, colorNone)

	// C. Duże ognisko (stare 69 -> SPRITE_GADGET_15)
	case textureID == spriteGadget15:
		frame := (bState.FireAnimationFrame + uint16(x+y)) % 4
		drawSprite(ps.Assets, spriteGadget14, xPos, yPos, colorNone)
		drawSprite(ps.Assets, spriteFire04+frame, xPos, yPos, colorNone)

	// D. Punkt Zwycięstwa (stare 301)
	// @reminder: tekstura spriteVictoryPoint się nie rysuje 22.04.202
	case textureID == spriteVictoryPoint:
		// Animacja punktu zwycięstwa (zakładamy 4 klatki animacji w atlasie)
		frame := (bState.FireAnimationFrame) % 4
		drawSprite(ps.Assets, spriteVictoryPoint+frame, xPos, yPos, colorNone)
	}
}

func temporaryEffects(affectedTile *tile, x, y uint8, xPos, yPos float32, bState *battleState, ps *programState) {
	// 1. Popiół
	if affectedTile.hasAsh && affectedTile.AshIntensity > 0.01 {
		alphaFactor := affectedTile.AshIntensity
		tint := rl.Fade(rl.White, alphaFactor)

		drawSpriteEx(spriteAsh00, xPos, yPos, colorNone, tint, ps)
	}

	// 2. Płonące kafelki
	if affectedTile.IsBurning {
		frame := (bState.FireAnimationFrame + uint16(x+y)) % 4
		drawSprite(ps.Assets, affectedTile.BurnOverlayID+frame, xPos, yPos, colorNone)
	}

	// 3. Duch
	if affectedTile.GhostEffect {
		phase1 := float64(x)*31 + float64(y)*17
		phase2 := float64(x)*43 + float64(y)*59

		drawGhostlySprite(spriteMissileGhostAttack, xPos, yPos, phase1, phase2, bState.GlobalFrameCounter, ps)
	}
}

// === UI, MINIMAPA, KURSOR ===
// @reminder: muszę się zastanowić, czy nie przenieść tego poniżej do osobnego pliku, bo
// tutaj już jest strasznie zagracone.

func drawGameCursorOnRealScreen(bState *battleState, pState *programState, scale float32, iState inputState) {
	rl.HideCursor()

	realMousePos := rl.GetMousePosition()

	// 1. Wywołujemy otulinę, która zdecyduje jaki ID kursora zwrócić
	cursorID := getCursorIDFromContext(bState, pState, realMousePos, scale, iState)

	// 2. Animacja i rysowanie (to dzieje się niezależnie od tego, skąd wzięliśmy ID)
	cursorID = animateCursorID(cursorID)
	drawCursorSprite(pState, cursorID, realMousePos, 3*bState.GameCamera.Zoom/scale)
}

func drawCursorSprite(pState *programState, cursorID uint16, pos rl.Vector2, scale float32) {
	def := spriteRegistry[cursorID]
	if def.cropWidth == 0 {
		return
	}

	var tex rl.Texture2D

	var srcRect rl.Rectangle

	if cursorID == spriteCursorFrameWhite {
		tex = pState.Assets.CursorWhiteFrame
		if tex.ID == 0 {
			return
		}

		srcRect = rl.NewRectangle(0, 0, float32(def.cropWidth), float32(def.cropHeight))
	} else {
		tex = pState.Assets.getAtlas(def.atlasID, colorNone)
		if tex.ID == 0 {
			return
		}

		srcRect = rl.NewRectangle(float32(def.cropX), float32(def.cropY), float32(def.cropWidth), float32(def.cropHeight))
	}

	finalX := pos.X + (float32(def.offX) * scale)
	finalY := pos.Y + (float32(def.offY) * scale)

	destW := float32(def.cropWidth) * scale
	destH := float32(def.cropHeight) * scale

	pState.RenderDestRect = rl.NewRectangle(finalX, finalY, destW, destH)
	pState.RenderOrigin = rl.NewVector2(0, 0)
	rl.DrawTexturePro(tex, srcRect, pState.RenderDestRect, pState.RenderOrigin, 0, rl.White)
}

func drawMilkBarVisualizer(bState *battleState, ps *programState) {
	anchorX := ps.GameViewWidth
	barX := anchorX + milkBarOffsetX
	barY := milkBarY
	barWidth := milkBarWidth
	barHeight := milkBarHeight

	currentMilk := float32(bState.HumanPlayerState.Milk)

	fillRatio := currentMilk / maxMilk

	if fillRatio > 1.0 {
		fillRatio = 1.0
	}

	filledHeight := barHeight * fillRatio
	fillY := barY + barHeight - filledHeight

	rl.DrawRectangle(int32(barX), int32(fillY), int32(barWidth), int32(filledHeight), rl.White)

	capacityLineRation := float32(bState.HumanPlayerState.MaxMilk) / maxMilk
	if capacityLineRation > 1.0 {
		capacityLineRation = 1.0
	}

	redLineY := barY + barHeight - (barHeight * capacityLineRation)

	rl.DrawLineEx(
		rl.NewVector2(barX, redLineY),
		rl.NewVector2(barX+barWidth, redLineY),
		2.0,
		rl.Red,
	)
}

func drawGameUI(bState *battleState, ps *programState) {
	anchorX := ps.GameViewWidth

	// 1. Rysowanie TŁA (Drewniany panel)
	if ps.Assets.WoodPanel.ID != 0 {
		ps.RenderDestRect = rl.NewRectangle(
			anchorX,
			0,
			float32(uiPanelVirtualWidth),
			float32(virtualScreenHeight),
		)
		// Źródło tekstury
		ps.RenderSrcRect = rl.NewRectangle(0, 0, float32(ps.Assets.WoodPanel.Width), float32(ps.Assets.WoodPanel.Height))

		rl.DrawTexturePro(ps.Assets.WoodPanel, ps.RenderSrcRect, ps.RenderDestRect, rl.NewVector2(0, 0), 0, rl.White)
	}

	// 2. Przyciski, pasek z mlekiem
	drawMilkBarVisualizer(bState, ps)
	drawButtons(bState, ps)

	// 3. Mapa
	minimapX := anchorX + minimapOffsetX
	minimapY := float32(0) + minimapOffsetY
	drawMinimapUnits(bState, minimapX, minimapY, minimapDisplayWidth, minimapDisplayHeight, ps.GameViewWidth)
}

// @todo: funkcja łamie zasadę rozdzielenia Dane-Logika-Rysowanie!
func drawButtons(bState *battleState, ps *programState) {
	// Punkt zaczepienia panelu (prawa strona ekranu)
	anchorX := ps.GameViewWidth

	// Pętla po 5 slotach przycisków
	for actionIndex := range uiActionMaxButtons {
		// 1. Pobieramy dane z battleState
		action := bState.UI.CurrentActions[actionIndex]

		// Jak nieaktywny, to nie rysujemy i nie klikamy
		if !action.IsActive {
			continue
		}

		// 2. Obliczamy pozycję (Layout)
		buttonY := uiAnchorOffsetY + (float32(actionIndex) * (btnHeight + btnMarginY))
		buttonX := anchorX + uiAnchorOffsetX
		rect := rl.NewRectangle(buttonX, buttonY, btnWidth, btnHeight)

		// Zapisujemy prostokąt do battleState (dla input.go)
		bState.UI.ActionButtons[actionIndex] = rect

		// @reminder: chwilowo wyłączyłem, ale myślę, że rozjaśnienie ikonki przycisku nad którym
		// się znajduje mysz będzie przydatny, powinno się to też powiązać ze wskazaniem ile miar
		// mleka to będzie kosztować
		// 3. Obsługa myszy (podświetlenie)
		// mousePos := screenToVirtualCoords(ps, rl.GetMousePosition())
		// isHover := rl.CheckCollisionPointRec(mousePos, rect)
		// bgColor := rl.Pink
		//
		// if isHover {
		// 	bgColor = rl.LightGray
		// }
		//
		// // Rysujemy tło przycisku
		// rl.DrawRectangleRec(rect, bgColor)

		iconID := action.IconID

		// jeśli IconID == 0 (nie ustawiono ręcznie),
		// wyliczamy ikonę na podstawie typu jednostki (np. "twarz" krowy)
		// @todo: teraz wszystkie jednostki z tego korzystają, ale muszę to zmienić
		// oby się udało wpisać „na sztywno”, co chcę wyświetlić. Bardzo uprości życie
		if action.Cmd.ActionType == cmdProduce {
			// Base (700) + Typ * 200 + Kierunek Dół (4) = Portret
			iconID = uint16(int(spriteUnitBaseID + (uint16(action.Cmd.ProduceType) * spriteUnitStep) + 4))
		}

		def := spriteRegistry[iconID]
		// Pobieramy odpowiedni atlas z AssetManagera
		// Używamy bs.PlayerID, żeby jednostki miały kolor gracza

		tex := rl.Texture2D{}

		var iconScale float32

		switch action.Cmd.ActionType {
		case cmdProduce:
			tex = ps.Assets.getAtlas(def.atlasID, bState.PlayerID)
			iconScale = 0.8
		case cmdBuildStructure, cmdStop, cmdRepairStructure:
			tex = ps.Assets.getAtlas(def.atlasID, colorNone)
			iconScale = 1
		case cmdCastSpell:
			switch action.Cmd.Spell {
			case spellMagicShower:
				tex = ps.Assets.getAtlas(def.atlasID, colorNone)
				drawTripleIcon(tex, def, rect)
				continue
			case spellMagicShield, spellMagicSight:
				tex = ps.Assets.getAtlas(def.atlasID, colorNone)
				iconScale = 1
			default:
				panic("nieobsługiwany rodzaj czaru. drawingBattle.go drawingButtons się wykrzaczyło")
			}

		default:
			// 4. Napis (Label)
			fontSize := int32(10)
			textWidth := rl.MeasureText(action.Label, fontSize)
			textX := rect.X + (rect.Width/2 - float32(textWidth)/2)
			textY := rect.Y + rect.Height - 12 // Na dole przycisku
			rl.DrawText(action.Label, int32(textX), int32(textY), fontSize, rl.Black)
		}

		// Wyznaczamy wycinek z atlasu (Source)
		srcRect := rl.NewRectangle(float32(def.cropX), float32(def.cropY), float32(def.cropWidth), float32(def.cropHeight))
		if def.flipX {
			srcRect.Width = -srcRect.Width
		}
		// Centrujemy ikonę na przycisku
		// Skalujemy ją, żeby zajmowała np. 80% wysokości przycisku
		scale := (rect.Height * iconScale) / float32(def.cropHeight)
		destW := float32(def.cropWidth) * scale
		destH := float32(def.cropHeight) * scale
		destX := rect.X + (rect.Width-destW)/2
		destY := rect.Y + (rect.Height-destH)/2
		rl.DrawTexturePro(tex, srcRect, rl.NewRectangle(destX, destY, destW, destH), rl.Vector2{}, 0, rl.White)
	}
}

// @todo: będzie wykorzystane przy rysowaniu ikonek już z ostatecznego układu.
func drawTripleIcon(tex rl.Texture2D, def spriteDef, btnRect rl.Rectangle) {
	// 0. Źródło: Jaki duszek będzie wykorzystany do narysowania potrojonej ikonki.
	srcW := float32(def.cropWidth)

	if def.flipX {
		srcW = -srcW
	}

	sourceRect := rl.NewRectangle(float32(def.cropX), float32(def.cropY), srcW, float32(def.cropHeight))
	textureOrigin := rl.NewVector2(0, 0)

	// 1. Wymiary: 85% wysokości przycisku.
	destH := btnRect.Height * tripleIconScale
	destW := destH * tileAspectRatio // Zachowanie proporcji

	// 2. Składanie.
	// Przesunięcie 35% szerokości.
	iconOffsetX := destW * tripleIconOverlapFactor

	// 3. Środkowanie potrojonego duszka.
	totalGroupWidth := destW + (2 * iconOffsetX)
	groupStartX := btnRect.X + (btnRect.Width-totalGroupWidth)/2.0

	// 4. Pozycjonowanie góra-dół.
	paddingY := btnRect.Height - destH

	// Skrajne przy górze, środkowa niżej.
	// aby uzyskać coś na podobieństwo '.'
	yHigh := btnRect.Y
	yLow := btnRect.Y + paddingY

	// 5. Rysowanie.
	// '__
	rl.DrawTexturePro(tex, sourceRect,
		rl.NewRectangle(groupStartX, yHigh, destW, destH),
		textureOrigin, 0, rl.White)

	// _._
	rl.DrawTexturePro(tex, sourceRect,
		rl.NewRectangle(groupStartX+iconOffsetX, yLow, destW, destH),
		textureOrigin, 0, rl.White)

	// __'
	rl.DrawTexturePro(tex, sourceRect,
		rl.NewRectangle(groupStartX+(2*iconOffsetX), yHigh, destW, destH),
		textureOrigin, 0, rl.White)
}

func drawMinimapUnits(bState *battleState, minimapX, minimapY, minimapWidth, minimapHeight, actualGameViewWidth float32) {
	const tileWidthF = float32(tileWidth)

	const tileHeightF = float32(tileHeight)

	const worldMapFullPixelWidth = float32(boardMaxX) * tileWidthF

	const worldMapFullPixelHeight = float32(boardMaxY) * tileHeightF

	pxStart := int(math.Round(float64(minimapX)))
	pyStart := int(math.Round(float64(minimapY)))
	totalDisplayWidthPx := int(math.Round(float64(minimapWidth)))
	totalDisplayHeightPx := int(math.Round(float64(minimapHeight)))

	xGridLines := make([]int, boardMaxX+1)
	yGridLines := make([]int, boardMaxY+1)

	for i := uint8(0); i <= boardMaxX; i++ {
		currentXFloat := float64(i) / float64(boardMaxX) * float64(totalDisplayWidthPx)
		xGridLines[i] = pxStart + int(math.Round(currentXFloat))
	}
	for i := uint8(0); i <= boardMaxY; i++ {
		currentYFloat := float64(i) / float64(boardMaxY) * float64(totalDisplayHeightPx)
		yGridLines[i] = pyStart + int(math.Round(currentYFloat))
	}
	xGridLines[boardMaxX] = pxStart + totalDisplayWidthPx
	yGridLines[boardMaxY] = pyStart + totalDisplayHeightPx

	// 1. Kafelki Terenu
	for x := range boardMaxX {
		for y := range boardMaxY {
			tileID := bState.Board.Tiles[x][y].TextureID
			col := getMapColor(tileID)
			w := xGridLines[x+1] - xGridLines[x]
			h := yGridLines[y+1] - yGridLines[y]
			rl.DrawRectangle(int32(xGridLines[x]), int32(yGridLines[y]), int32(w), int32(h), col)
		}
	}

	// 2. Budynki
	for _, bld := range bState.Buildings {
		if !bld.Exists {
			continue
		}

		// pomijamy palisady
		if bld.Type == buildingPalisade {
			continue
		}

		var bldColor rl.Color

		if bld.Owner == bState.PlayerID {
			bldColor = rl.White
		} else {
			bldColor = rl.Red
		}

		for _, tile := range bld.OccupiedTiles {
			if tile.X < 0 || tile.X >= boardMaxX || tile.Y < 0 || tile.Y >= boardMaxY {
				continue
			}
			w := xGridLines[tile.X+1] - xGridLines[tile.X]
			h := yGridLines[tile.Y+1] - yGridLines[tile.Y]
			rl.DrawRectangle(int32(xGridLines[tile.X]), int32(yGridLines[tile.Y]), int32(w), int32(h), bldColor)
		}
	}

	// 3. Jednostki
	for _, unit := range bState.Units {
		if !unit.Exists {
			continue
		}
		var unitColor rl.Color
		if unit.Owner == bState.PlayerID {
			unitColor = rl.White
		} else {
			unitColor = rl.Red
		}

		if unit.X < 0 || unit.X >= boardMaxX || unit.Y < 0 || unit.Y >= boardMaxY {
			continue
		}

		w := xGridLines[unit.X+1] - xGridLines[unit.X]
		h := yGridLines[unit.Y+1] - yGridLines[unit.Y]
		rl.DrawRectangle(int32(xGridLines[unit.X]), int32(yGridLines[unit.Y]), int32(w), int32(h), unitColor)
	}

	// 4. Ramka Kamery
	camWorldView := getCameraWorldViewRect(bState.GameCamera, actualGameViewWidth, float32(virtualScreenHeight))
	scaleFactorX := minimapWidth / worldMapFullPixelWidth
	scaleFactorY := minimapHeight / worldMapFullPixelHeight

	// Jaki powinien być prostokąt pokazujący
	rawCameraRect := rl.NewRectangle(
		minimapX+(camWorldView.X*scaleFactorX),
		minimapY+(camWorldView.Y*scaleFactorY),
		camWorldView.Width*scaleFactorX,
		camWorldView.Height*scaleFactorY,
	)

	minimapBounds := rl.NewRectangle(
		minimapX,
		minimapY,
		minimapWidth,
		minimapHeight,
	)

	// Obliczamy część wspólną dla rawCameraRect oraz minimapBounds
	// dzięki temu nawet jeżeli widzimy czarny ekran po bokach planszy, to żółty prostokąt
	// nie wyjdzie poza mapę
	clampedRect := rl.GetCollisionRec(minimapBounds, rawCameraRect)

	rl.DrawRectangleLinesEx(clampedRect, 1.0, rl.Yellow)
}

// ========= DEBUGOWANIE ZASADZANIA BUDOWLI

func drawConstructionValidationBox(bState *battleState, ps *programState) {
	if bState.MouseState != mouseStateBuilding {
		return
	}

	screenMouse := rl.GetMousePosition()
	virtualMouse := screenToVirtualCoords(ps, screenMouse)

	if isMouseOverUI(ps, virtualMouse) {
		return
	}

	worldMousePos := rl.GetScreenToWorld2D(virtualMouse, bState.GameCamera)

	startX := uint8(worldMousePos.X / float32(tileWidth))
	startY := uint8(worldMousePos.Y / float32(tileHeight))

	size := buildingDefs[bState.PendingBuildingType].Width

	var tileStates [3][3]uint8
	var hasRoadAccessAnywhere bool

	for dy := uint8(0); dy < size; dy++ {
		for dx := uint8(0); dx < size; dx++ {
			cx := startX + dx
			cy := startY + dy

			rawColor := validationBoxColor(cx, cy, bState)

			var state uint8

			if rawColor == rl.Red {
				state = 0
			} else if rawColor == rl.Orange {
				state = 1
			} else {
				state = 2
				hasRoadAccessAnywhere = true
			}

			tileStates[dy][dx] = state
		}
	}

	for dy := uint8(0); dy < size; dy++ {
		for dx := uint8(0); dx < size; dx++ {
			state := tileStates[dy][dx]
			finalColor := rl.Red

			switch state {
			case 0:
				finalColor = rl.Red
			case 1:
				if hasRoadAccessAnywhere {
					finalColor = rl.Green
				} else {
					finalColor = rl.Orange
				}
			case 2:
				finalColor = rl.DarkGreen
			}

			finalColor = rl.Fade(finalColor, 0.3)

			posX := float32(startX+dx) * float32(tileWidth)
			posY := float32(startY+dy) * float32(tileHeight)

			rl.DrawRectangle(int32(posX), int32(posY), int32(tileWidth), int32(tileHeight), finalColor)
		}
	}
}

func validationBoxColor(tileX, tileY uint8, bState *battleState) rl.Color {
	// 1. Walidacja specyficzna dla typu budynku
	var isValid bool

	switch bState.PendingBuildingType {
	case buildingBridge:
		isValid = isWithinBoard(tileX, tileY, bState) &&
			isWaterTileOnly(bState.Board.Tiles[tileX][tileY].TextureID)

	case buildingRoad:
		isValid = !isPath(bState.Board.Tiles[tileX][tileY].TextureID) &&
			canFitBuilding(tileX, tileY, smallBuildingSize, smallBuildingSize, bState)

	default:
		isValid = canFitBuilding(tileX, tileY, smallBuildingSize, smallBuildingSize, bState)
	}

	// 2. Jeśli miejsce jest nieważne, od razu zwracamy czerwień
	if !isValid {
		return rl.Red
	}

	// 3. Palisada nie wymaga dostępu do drogi, więc od razu jest zielona
	if bState.PendingBuildingType == buildingPalisade {
		return rl.DarkGreen
	}

	// 4. Budynki wymagają drogi
	if hasRoadAccess(tileX, tileY, smallBuildingSize, bState) {
		return rl.DarkGreen
	}

	// 5. Ważne, ale bez dostępu do drogi
	return rl.Orange
}
