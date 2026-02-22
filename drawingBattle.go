package main

// drawingBattle.go

// @todo: OPTYMALIZACJA PAMIĘCI. Obecnie w pętli rysowania tworzone są nowe obiekty.
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

// Pomagierzy

func isSpecialTile(id uint16) bool {
	return id == spriteEffectHeal00 || id == spriteEffectHeal01 || id == spriteVictoryPoint
}

func isRoad(id uint16) bool {
	return id >= spriteRoadStart && id <= spriteRoadEnd
}

func isRuin(id uint16) bool {
	return id >= spriteRuinStart && id <= spriteRuinEnd
}

func isPalisade(id uint16) bool {
	return id >= spritePalisadeStart && id <= spritePalisadeEnd
}

func isWaterTileOnly(id uint16) bool {
	return id >= spriteWaterStart && id <= spriteWaterEnd
}

func isBridge(id uint16) bool {
	return id >= spriteBridgeStart && id <= spriteBridgeEnd
}

func isWaterOrBridgeForMasking(id uint16) bool {
	return isWaterTileOnly(id) || isBridge(id) || id == spriteBridgeConstruction
}

func isLandOrOther(id uint16) bool {
	return !isWaterTileOnly(id) && !isBridge(id) && id != spriteBridgeConstruction
}

func isBuildingTerrain(id uint16) bool {
	return (id >= spriteConstructionStart && id <= spriteConstructionEnd) ||
		(id >= spriteBuildingMainBase && id <= spriteBuildingEnd)
}

func isHealingShire(id uint16) bool {
	return id == spriteEffectHeal00 || id == spriteEffectHeal01
}

func isDryEarth(id uint16) bool {
	return id >= spriteDryEarth01 && id <= spriteDryEarth03
}

func isGrass(tileID uint16) bool {
	return tileID >= spriteGrassStart && tileID <= spriteGrassEnd
}

func isGadget(tileID uint16) bool {
	return tileID >= spriteGadgetStart && tileID <= spriteGadgetEnd
}

func isGrassOrGadget(tileID uint16) bool {
	return isGrass(tileID) || isGadget(tileID)
}

func isTreeStump(tileID uint16) bool {
	return tileID >= spriteTreeStumpStart && tileID <= spriteTreeStumpEnd
}

func isTree(tileID uint16) bool {
	return tileID >= spriteTreeStumpStart && tileID <= spriteTreeTopEnd ||
		tileID >= spriteTreeBurntStump00 && tileID <= spriteTreeFallingBurnt02
}

func isRockNonWalkable(tileID uint16) bool {
	return tileID >= spriteRockStart && tileID < spriteRockEnd
}

func isPath(tileID uint16) bool {
	return isRoad(tileID) || isBridge(tileID)
}

// === PRZETWARZANIE MAPY ===

func applyRoadProcessing(x, y uint8, board *boardData) {
	var mask uint8 = 0

	if y > 0 && isRoad(board.Tiles[x][y-1].TextureID) {
		mask |= 1
	}

	if x < boardMaxX-1 && isRoad(board.Tiles[x+1][y].TextureID) {
		mask |= 2
	}

	if y < boardMaxY-1 && isRoad(board.Tiles[x][y+1].TextureID) {
		mask |= 4
	}

	if x > 0 && isRoad(board.Tiles[x-1][y].TextureID) {
		mask |= 8
	}

	offset := uint16(13) // Domyślnie środek (mask 0, 7, 15)

	switch mask {
	case 1:
		offset = 9
	case 2:
		offset = 6
	case 3:
		offset = 10
	case 4:
		offset = 8
	case 5:
		offset = 11
	case 6:
		offset = 12
	case 0, 7, 15:
		offset = 13
	case 8:
		offset = 7
	case 9:
		offset = 14
	case 10:
		offset = 15
	case 11:
		offset = 16
	case 12:
		offset = 17
	case 13:
		offset = 18
	case 14:
		offset = 19
	}

	board.Tiles[x][y].TextureID = spriteRoadStart + offset
}

func applyPalisadeProcessing(x, y uint8, board *boardData) {
	var mask uint8 = 0

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

func processMapTiles(bs *battleState) {
	// Tworzymy snapshot TextureID, aby uniknąć efektów ubocznych podczas iteracji
	var snapshot [boardMaxX][boardMaxY]uint16

	for x := uint8(0); x < boardMaxX; x++ {
		for y := uint8(0); y < boardMaxY; y++ {
			snapshot[x][y] = bs.Board.Tiles[x][y].TextureID
		}
	}

	for y := uint8(0); y < boardMaxY; y++ {
		for x := uint8(0); x < boardMaxX; x++ {
			id := bs.Board.Tiles[x][y].TextureID

			switch {
			case isRoad(id):
				applyRoadProcessing(x, y, bs.Board)
			case isPalisade(id):
				applyPalisadeProcessing(x, y, bs.Board)
			case isWaterTileOnly(id):
				applyWaterProcessing(x, y, bs.Board, snapshot)
			case isHealingShire(id):
				bs.HealingShrines = append(bs.HealingShrines, point{X: x, Y: y})
			}
		}
	}

	makeGrassVariations(bs)
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
	var mask uint8 = 0

	if y > 0 && isWaterOrBridgeForMasking(board.Tiles[x][y-1].TextureID) {
		mask |= 1
	}

	if x < boardMaxX-1 && isWaterOrBridgeForMasking(board.Tiles[x+1][y].TextureID) {
		mask |= 2
	}

	if y < boardMaxY-1 && isWaterOrBridgeForMasking(board.Tiles[x][y+1].TextureID) {
		mask |= 4
	}

	if x > 0 && isWaterOrBridgeForMasking(board.Tiles[x-1][y].TextureID) {
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
	if def.w == 0 {
		return // Pusty wpis (np. ID 0)
	}

	// 3. Pobranie Tekstury
	tex := assets.getAtlas(def.atlasID, ownerColor)
	if tex.ID == 0 {
		return // Zabezpieczenie przed brakiem atlasu
	}

	// 4. Obliczenie Prostokąta Źródłowego
	srcW := float32(def.w)

	if def.flipX {
		srcW = -srcW // Raylib obsługuje zwierciadlane odbicie przez ujemną szerokość
	}

	srcRect := rl.NewRectangle(float32(def.x), float32(def.y), srcW, float32(def.h))

	// 5. Obliczenie Pozycji Docelowej z uwzględnieniem offsetów
	finalX := destX + float32(def.offX)
	finalY := destY + float32(def.offY)

	destRect := rl.NewRectangle(finalX, finalY, float32(def.w), float32(def.h))

	// 6. Rysowanie
	rl.DrawTexturePro(tex, srcRect, destRect, rl.NewVector2(0, 0), 0, rl.White)
}

func drawSpriteEx(id uint16, destX, destY float32, ownerColor uint8, tint rl.Color, ps *programState) {
	if id >= maxSpriteID {
		return
	}

	def := spriteRegistry[id]

	if def.w == 0 {
		return // Pusty wpis
	}

	// Pobieramy atlas z managera
	tex := ps.Assets.getAtlas(def.atlasID, ownerColor)

	// Logika flipX i Offsetów z definicji
	srcW := float32(def.w)

	if def.flipX {
		srcW = -srcW
	}

	srcRect := rl.NewRectangle(float32(def.x), float32(def.y), srcW, float32(def.h))

	// Stosujemy offsety z SpriteDef
	finalX := destX + float32(def.offX)
	finalY := destY + float32(def.offY)

	destRect := rl.NewRectangle(finalX, finalY, float32(def.w), float32(def.h))

	rl.DrawTexturePro(tex, srcRect, destRect, rl.NewVector2(0, 0), 0, tint)
}

func drawFrameCorners(x, y, width, height, cLen float32) {
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

func drawBuildingSelectionFrame(bld *building, bounds bounds) {
	var cLen float32

	if bld.Type != buildingPalisade {
		cLen = cornerLenBuilding
	} else {
		cLen = cornerLenPalisade
	}

	drawFrameCorners(float32(bounds.X), float32(bounds.Y), bounds.WidthPx, bounds.HeightPx, cLen)
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

func drawBuildingInterface(bld *building, bs *battleState) {
	bounds := bld.bounds()
	isBuildingSelected := bs.CurrentSelection.BuildingID == bld.ID

	if isBuildingSelected {
		drawBuildingSelectionFrame(bld, bounds)

		if bld.Owner == bs.PlayerID && !bld.IsUnderConstruction {
			drawBuildingCapacity(bld, bounds)
		}
	}

	// Rysowanie paska życia
	if (bld.HP < bld.MaxHP || isBuildingSelected) && bld.HP != 0 {
		drawBuildingHealthBar(bld, bounds)
	}
}

func drawBuildingsInterfaces(bs *battleState) {
	for _, building := range bs.Buildings {
		if !building.Exists || len(building.OccupiedTiles) == 0 {
			continue
		}

		drawBuildingInterface(building, bs)
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

func getRenderDirection(u *unit, bs *battleState) (int, int) {
	if u.AnimationType == "fight" {
		var targetX, targetY uint8
		foundTarget := false

		// 1. Sprawdź, czy mamy cel ataku (Unit lub Building)
		if u.TargetID != 0 {
			targetUnit, targetBld := getObjectByID(u.TargetID, bs)

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

//func DrawProjectiles(bs *battleState, ps *programState) {
//	for _, p := range bs.Projectiles {
//		if p.Exists {
//			DrawSingleProjectile(p)
//		}
//	}
//}

//func DrawSingleProjectile(p *Projectile, ps *programState) {
//	tex, ok := ps.ProjectileSprites[p.Type]
//	if !ok {
//		return
//	}
//
//	dirX := 0
//	if p.DX > 0.5 {
//		dirX = 1
//	} else if p.DX < -0.5 {
//		dirX = -1
//	}
//
//	dirY := 0
//	if p.DY > 0.5 {
//		dirY = 1
//	} else if p.DY < -0.5 {
//		dirY = -1
//	}
//
//	frameW := float32(tex.Width) / 2.0
//	frameH := float32(tex.Height) / 3.0
//
//	col := 0
//	if dirX == 0 {
//		col = 1 // Pionowo (Góra/Dół)
//	} else {
//		col = 0 // Poziomo/Skos
//	}
//
//	row := dirY + 1
//	if row < 0 {
//		row = 0
//	}
//	if row > 2 {
//		row = 2
//	}
//
//	flipX := (dirX == 1)
//
//	sourceRec := rl.NewRectangle(float32(col)*frameW, float32(row)*frameH, frameW, frameH)
//	if flipX {
//		sourceRec.Width = -sourceRec.Width
//	}
//
//	drawX := p.x - (frameW / 2.0)
//	drawY := p.y - (frameH / 2.0)
//
//	rl.DrawTexturePro(tex, sourceRec, rl.NewRectangle(drawX, drawY, frameW, frameH), rl.NewVector2(0, 0), 0, rl.White)
//}

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

func drawUnitSelectionFrame(unit *unit) {
	x := float32(unit.X) * float32(tileWidth)
	y := float32(unit.Y) * float32(tileHeight)

	drawFrameCorners(x, y, float32(tileWidth), float32(tileHeight), cornerLenPalisade)
}

func drawUnitInterface(unit *unit) {
	screenX := int32(unit.X) * int32(tileWidth)
	screenY := int32(unit.Y) * int32(tileHeight)

	if unit.IsSelected {
		drawUnitSelectionFrame(unit)
		drawUnitHealthBar(screenX, screenY, unit)

		if unit.MaxMana > 0 {
			drawManaBar(screenX, screenY, unit)
		}

		if unit.Type == unitCow {
			drawMilkBar(screenX, screenY, unit)
		}
	}
}

// === LOGIKA KAMERY I EKRANU ===

func setupGameCamera(bs *battleState, ps *programState) {
	bs.GameCamera = rl.NewCamera2D(
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

func drawBattleScene(bs *battleState) {
	rl.BeginScissorMode(0, 0, int32(gameViewVirtualWidth), int32(virtualScreenHeight))
	rl.BeginMode2D(bs.GameCamera)
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

func drawSoil(startX, startY, endX, endY uint8, bs *battleState, ps *programState) {
	for yAxis := startY; yAxis < endY; yAxis++ {
		for xAxis := startX; xAxis < endX; xAxis++ {
			currentTile := &bs.Board.Tiles[xAxis][yAxis]
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

			if isSpecialTile(texID) || isTreeStump(texID) {
				drawSprite(ps.Assets, spriteGrass00, xPos, yPos, colorNone)
			}

			if isBridge(texID) {
				waterBaseID := spriteWaterMiddle
				animationOffset := bs.WaterAnimationCounter * 13
				drawSprite(ps.Assets, waterBaseID+animationOffset, xPos, yPos, colorNone)
				drawSprite(ps.Assets, texID, xPos, yPos, colorNone)
			}

			if texID == spriteBridgeConstruction {
				waterBaseID := calculateWaterTileID(xAxis, yAxis, bs.Board)

				if waterBaseID == 999 {
					waterBaseID = spriteWaterMiddle
				}

				animationOffset := bs.WaterAnimationCounter * 13
				drawSprite(ps.Assets, waterBaseID+animationOffset, xPos, yPos, colorNone)
				drawSprite(ps.Assets, texID, xPos, yPos, colorNone)
			}

			if isWaterTileOnly(texID) {
				animationOffset := bs.WaterAnimationCounter * 13
				drawSprite(ps.Assets, texID+animationOffset, xPos, yPos, colorNone)
			}

			if isLandOrOther(texID) && !isTree(texID) {
				drawSprite(ps.Assets, texID, xPos, yPos, colorNone)
			}
		}
	}
}

func drawBuilding(startY, endY uint8, bs *battleState, ps *programState) {
	for yAxis := startY; yAxis < endY; yAxis++ {
		for _, bld := range bs.RenderBuildingRows[yAxis] {
			for _, tile := range bld.OccupiedTiles {
				if tile.Y == yAxis {
					id := bs.Board.Tiles[tile.X][yAxis].TextureID
					if isBuildingTerrain(id) {
						finalID := id
						if flagID, ok := flagAnimationMap[uint8(id)]; ok && (bs.FireAnimationCounter+uint16(tile.X)+uint16(yAxis))%2 == 1 {
							finalID = uint16(flagID)
						}

						drawSprite(ps.Assets, finalID, float32(tile.X)*float32(tileWidth), float32(yAxis)*float32(tileHeight), bld.Owner)
					}
				}
			}
		}
	}
}

func drawPalisade(startX, startY, endX, endY uint8, bs *battleState, ps *programState) {
	for yAxis := startY; yAxis < endY; yAxis++ {
		for xAxis := startX; xAxis < endX; xAxis++ {
			id := bs.Board.Tiles[xAxis][yAxis].TextureID
			if isPalisade(id) {
				drawSprite(ps.Assets, id, float32(xAxis)*float32(tileWidth), float32(yAxis)*float32(tileHeight), colorNone)
			}
		}
	}
}

func drawDryEarth(startX, startY, endX, endY uint8, bs *battleState, ps *programState) {
	for y := startY; y < endY; y++ {
		for x := startX; x < endX; x++ {
			id := bs.Board.Tiles[x][y].TextureID
			if isDryEarth(id) {
				drawSprite(ps.Assets, id, float32(x)*float32(tileWidth), float32(y)*float32(tileHeight), colorNone)
			}
		}
	}
}

func drawBuildings(startX, startY, endX, endY uint8, bs *battleState, ps *programState) {
	drawBuilding(startY, endY, bs, ps)
	drawPalisade(startX, startY, endX, endY, bs, ps)
	drawDryEarth(startX, startY, endX, endY, bs, ps)
}

func drawCorpsesUnitsTrees(startX, startY, endX, endY uint8, bs *battleState, ps *programState) {
	for y := startY; y < endY; y++ {
		drawCorpses(y, startX, endX, bs, ps)
	}

	for y := startY; y < endY; y++ {
		drawUnits(y, bs, ps)
		drawTrees(y, startX, endX, bs, ps)
	}
}

func drawUnits(y uint8, bs *battleState, ps *programState) {
	rowUnits := bs.RenderUnitRows[y]
	if len(rowUnits) > 0 {
		sort.Slice(rowUnits, func(i, j int) bool { return rowUnits[i].X < rowUnits[j].X })

		for _, unit := range rowUnits {
			if unit.Owner == bs.PlayerID {
				drawUnit(unit, bs, ps)
			}
		}

		for _, unit := range rowUnits {
			if unit.Owner != bs.PlayerID {
				drawUnit(unit, bs, ps)
			}
		}
	}
}

func drawCorpses(y, startX, endX uint8, bs *battleState, ps *programState) {
	for _, corpse := range bs.Corpses {
		if corpse.Y == y && corpse.X >= startX && corpse.X < endX {
			posX := float32(corpse.X) * float32(tileWidth)
			posY := float32(corpse.Y) * float32(tileHeight)

			if corpse.Phase < corpsesPhase2 {
				var offsetIndex uint16

				if corpse.Phase == corpsesPhase1 {
					offsetIndex = 1
				}

				baseID := spriteUnitBaseID + (uint16(corpse.UnitType) * spriteUnitStep)
				finalID := baseID + corpsesFrameIndexOffset + offsetIndex
				drawSprite(ps.Assets, finalID, posX, posY, corpse.Owner)
			} else {
				cid := spriteEffectSkeleton_00 + uint16(corpse.SkeletonType)
				tint := rl.Fade(rl.White, float32(corpse.Alpha)/corpsesMaxAlpha)
				drawSpriteEx(cid, posX, posY, colorNone, tint, ps)
			}
		}
	}
}

func drawTrees(y, startX, endX uint8, bs *battleState, ps *programState) {
	for x := startX; x < endX; x++ {
		texID := bs.Board.Tiles[x][y].TextureID
		if texID >= spriteTreeStumpStart && texID <= spriteTreeStumpEnd {
			drawX := float32(x)*float32(tileWidth) - treeOffsetX
			drawY := float32(y) * float32(tileHeight)
			drawSprite(ps.Assets, texID, drawX, drawY, colorNone)

			if y > 0 {
				crownY := float32(y-1) * float32(tileHeight)
				drawSprite(ps.Assets, texID+7, drawX, crownY, colorNone)
			}
		}
	}
}

// @todo: porozbijać na podfunkcje! Inaczej nie rozprawię się z tym potworem.
func drawWorldAndUnits(bs *battleState, ps *programState) {
	// 1. Aktualizacja podręcznych
	bs.updateRenderCache()
	cam := bs.GameCamera
	worldLeft := cam.Target.X - (cam.Offset.X / cam.Zoom)
	worldTop := cam.Target.Y - (cam.Offset.Y / cam.Zoom)
	worldRight := worldLeft + (ps.GameViewWidth / cam.Zoom)
	worldBottom := worldTop + (ps.VirtualHeight / cam.Zoom)
	startXInt := int(worldLeft/float32(tileWidth)) - 1
	startYInt := int(worldTop/float32(tileHeight)) - 1
	endXInt := int(worldRight/float32(tileWidth)) + 2
	endYInt := int(worldBottom/float32(tileHeight)) + 2

	// Zabezpieczenie dolne (teraz działa poprawnie, bo int może być ujemny)
	if startXInt < 0 {
		startXInt = 0
	}

	if startYInt < 0 {
		startYInt = 0
	}

	// Zabezpieczenie górne
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
	drawSoil(startX, startY, endX, endY, bs, ps)
	// Przebieg 2: budynki, sucha ziemia, palisady
	drawBuildings(startX, startY, endX, endY, bs, ps)
	// Przebieg 3: zwłoki, jednostki, drzewa
	drawCorpsesUnitsTrees(startX, startY, endX, endY, bs, ps)

	// Nakładki i efekty.
	drawBuildingsInterfaces(bs)
	drawEffects(bs, ps)
	// @reminder: tymczasowy hack
	// @reminder: 12.02.2026 szkoda, że nie napisałem o co chodziło, bo już nie pamiętam.
	for y := startY; y < endY; y++ {
		rowUnits := bs.RenderUnitRows[y]
		if len(rowUnits) > 0 {
			for _, unit := range rowUnits {
				// Rysujemy tylko nakładkę dla jednostek, żeby była ponad drzewami itd.
				drawUnitInterface(unit)
			}
		}
	}
	// DrawProjectiles(bs)
	drawSelectionBox(bs, ps)
}

// @todo: rozdziel logikę od rysowania, bo to zwyczajne pomieszanie z poplątaniem.
func drawUnit(u *unit, bs *battleState, ps *programState) {
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
	renderDx, renderDy := getRenderDirection(u, bs)
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
		drawUnitWounds(u, ps, screenX, screenY)
	}
}

func drawUnitWounds(u *unit, ps *programState, screenX, screenY float32) {
	// Stałe offsety, aby krew była na środku kafelka 16x14.
	const centerOffsetX = 8.0
	const centerOffsetY = 7.0

	for _, wound := range u.Wounds {
		// Wybór tekstury
		effectID := spriteEffectHit00
		if wound.IsSevere && wound.Timer <= 10 {
			effectID = spriteEffectHit01
		}

		def := spriteRegistry[effectID]
		if def.w == 0 {
			continue
		}

		tex := ps.Assets.getAtlas(def.atlasID, colorNone)

		sourceRec := rl.NewRectangle(
			float32(def.x), float32(def.y),
			float32(def.w), float32(def.h),
		)

		// Skalowanie rany
		destW := float32(def.w) * wound.Scale
		destH := float32(def.h) * wound.Scale

		// Pozycja finalna: pozycja jednostki + środek kafelka + losowy offset rany
		destRect := rl.NewRectangle(
			screenX+centerOffsetX+wound.OffsetX,
			screenY+centerOffsetY+wound.OffsetY,
			destW,
			destH,
		)

		origin := rl.NewVector2(destW/2, destH/2)

		rl.DrawTexturePro(tex, sourceRec, destRect, origin, wound.Rotation, rl.White)
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

func (bs *battleState) updateRenderCache() {
	// @reminder Być może kiedyś okaże się, iż czyszczenie całości nie jest właściwe.
	// Wtedy będę musiał zastanowić się, jak wywalać poszczególne jednostki i budynki.
	// 1. Czyścimy jednostki oraz budynki po poprzedniej klatce
	for y := range boardMaxY {
		bs.RenderUnitRows[y] = bs.RenderUnitRows[y][:0]
		bs.RenderBuildingRows[y] = bs.RenderBuildingRows[y][:0]
	}
	// 2. Dodajemy jednostki
	for _, unit := range bs.Units {
		// Uwzględniamy tylko żywe jednostki
		// Nie da się mieć jednostki poza planszą, więc nie muszę sprawdzać
		// MAX_X i MAX_Y, pewnie kiedyś tego pożałuję
		// @check nie sprawdzam MAX_Y oraz MAX_Y, w razie problemów z rysowaniem
		// tutaj może się kryć przyczyna. Na przyszłość: nie <= bo pierwszy index to 0
		if unit.Exists { // && unit.y >= 0 && unit.y < MAX_Y
			bs.RenderUnitRows[unit.Y] = append(bs.RenderUnitRows[unit.Y], unit)
		}
	}
	// 3. Dodajemy budynki
	for _, bld := range bs.Buildings {
		// Pomijamy palisady i zniszczone
		if !bld.Exists || bld.Type == buildingPalisade {
			continue
		}

		// Ponieważ budynki trzymają informacje o 9 kafelkach, to trzeba się zabezpieczyć
		lastY := uint8(0)
		for _, tile := range bld.OccupiedTiles {
			if tile.Y != lastY {
				// if tile.y >= 0 && tile.y < MAX_Y // @reminder Nie sprawdzam legalności!
				bs.RenderBuildingRows[tile.Y] = append(bs.RenderBuildingRows[tile.Y], bld)
			}

			lastY = tile.Y
		}
	}
}

// drawSelectionBox odpowiada za rysowanie prostokąta do zaznaczania jednostek.
func drawSelectionBox(bs *battleState, ps *programState) {
	if !bs.IsSelectingBox {
		return
	}

	// 1. Pobieramy obecną mysz i zamieniamy ją na WIRTUALNĄ (bo takiej używa logika gry)
	// Bez tego, na dużym monitorze kursor "ucieknie" od ramki.
	currentScreenMouse := rl.GetMousePosition()
	currentVirtualMouse := screenToVirtualCoords(ps, currentScreenMouse)

	// 2. Konwertujemy OBA punkty (Start i Koniec) z Wirtualnego Ekranu na Świat Gry
	// Dzięki temu, że jesteśmy w BeginMode2D, rysujemy wprost na mapie.

	// Start zaznaczania (zapamiętany w input.go jako wirtualny ekran) → Świat
	worldStart := rl.GetScreenToWorld2D(bs.SelectionStart, bs.GameCamera)

	// Obecna mysz (wirtualny ekran) → Świat
	worldEnd := rl.GetScreenToWorld2D(currentVirtualMouse, bs.GameCamera)

	// 3. Obliczamy wymiary w Świecie Gry
	x := float32(math.Min(float64(worldStart.X), float64(worldEnd.X)))
	y := float32(math.Min(float64(worldStart.Y), float64(worldEnd.Y)))
	width := float32(math.Abs(float64(worldEnd.X - worldStart.X)))
	height := float32(math.Abs(float64(worldEnd.Y - worldStart.Y)))

	rect := rl.NewRectangle(x, y, width, height)

	// Linia musi mieć grubość niezależną od przybliżenia, inaczej może stać się niewidoczna
	const minScreenThicknessPixels = 1.5
	requiredWorldThickness := minScreenThicknessPixels / bs.GameCamera.Zoom

	finalThickness := float32(math.Max(1.0, float64(requiredWorldThickness)))

	// Rysowanie
	rl.DrawRectangleLinesEx(rect, finalThickness, rl.White) // Zielona ramka, klasyk RTS
}

func drawEffects(bs *battleState, ps *programState) {
	for y := uint8(0); y < boardMaxY; y++ {
		for x := uint8(0); x < boardMaxX; x++ {
			tile := &bs.Board.Tiles[x][y]
			id := tile.TextureID

			xPos := float32(x) * float32(tileWidth)
			yPos := float32(y) * float32(tileHeight)

			// --- 1. EFEKTY WYNIKAJĄCE Z RODZAJU PODŁOŻA (TextureID) ---
			switch {
			// A. Kapliczka leczenia (stare 282/283)
			case id == spriteEffectHeal00 || id == spriteEffectHeal01:
				frame := (bs.WaterAnimationCounter + uint16(x+y)) % 2
				animID := spriteEffectHeal00 + frame

				// Logika przeźroczystości (Duch)
				var tint rl.Color
				if tile.Unit != nil {
					tint = rl.White // Pełna widoczność, gdy ktoś stoi
				} else {
					tint = rl.Fade(rl.White, 0.4) // Duch, gdy pusto
				}
				drawSpriteEx(animID, xPos, yPos, colorNone, tint, ps)

			// B. Małe ognisko (stare 68 -> SPRITE_GADGET_14)
			case id == spriteGadget14:
				frame := (bs.FireAnimationCounter + uint16(x+y)) % 4
				drawSprite(ps.Assets, spriteGadget14, xPos, yPos, colorNone)
				drawSprite(ps.Assets, spriteFire08+frame, xPos, yPos, colorNone)

			// C. Duże ognisko (stare 69 -> SPRITE_GADGET_15)
			case id == spriteGadget15:
				frame := (bs.FireAnimationCounter + uint16(x+y)) % 4
				drawSprite(ps.Assets, spriteGadget14, xPos, yPos, colorNone)
				drawSprite(ps.Assets, spriteFire04+frame, xPos, yPos, colorNone)

			// D. Punkt Zwycięstwa (stare 301)
			case id == spriteVictoryPoint:
				// Animacja punktu zwycięstwa (zakładamy 4 klatki animacji w atlasie)
				frame := (bs.FireAnimationCounter) % 4
				drawSprite(ps.Assets, spriteVictoryPoint+frame, xPos, yPos, colorNone)
			}

			//// --- 2. EFEKTY NAKŁADANE NA WIERZCH (EffectID / dawne AttackValue) ---
			//// Np. pożar lasu, wybuchy, obrażenia
			//effectVal := tile.EffectID
			//if effectVal > 0 && tile.Visibility != visibilityUnexplored {
			//	switch {
			//	// Ogień bitewny (np. podpalony las)
			//	case effectVal > 70 && effectVal <= 100:
			//		fireFrame := min(bs.FireAnimationCounter%14, 13)
			//		// Używamy bazowego ID ognia (np. SPRITE_FIRE_00)
			//		drawSprite(ps.Assets, spriteFire00+fireFrame, xPos, yPos, colorNone)
			//
			//	// Efekt uderzenia / krwi (hit)
			//	case effectVal > 100 && effectVal < 127:
			//		// drawSprite(SPRITE_EFFECT_HIT, xPos, yPos, COLOR_NONE)
			//
			//	// Zgliszcza / Śmierć obiektu
			//	case effectVal > 127:
			//		// drawSprite(SPRITE_EFFECT_DESTROYED, xPos, yPos, COLOR_NONE)
			//	}
			//}
		}
	}
}

// === UI, MINIMAPA, KURSOR ===
// @reminder: muszę się zastanowić, czy nie przenieść tego poniżej do osobnego pliku, bo
// tutaj już jest strasznie zagracone.

func drawGameCursorOnRealScreen(bs *battleState, ps *programState, scale float32) {
	rl.HideCursor()

	realMousePos := rl.GetMousePosition()

	// 1. Wywołujemy otulinę, która zdecyduje jaki ID kursora zwrócić
	cursorID := getCursorIDFromContext(bs, ps, realMousePos, scale)

	// 2. Animacja i rysowanie (to dzieje się niezależnie od tego, skąd wzięliśmy ID)
	cursorID = animateCursorID(cursorID)
	drawCursorSprite(ps, cursorID, realMousePos, 3*bs.GameCamera.Zoom/scale)
}

// getCursorIDFromContext - otulina
// Sprawdza stan gry. Jeśli gameScreen → deleguje do determineCursorState.
// Jeśli inny stan → zwraca domyślny kursor.
func getCursorIDFromContext(bs *battleState, ps *programState, realMousePos rl.Vector2, scale float32) uint16 {
	if ps.CurrentState == gameScreen {
		virtualMouseX := realMousePos.X / scale
		virtualMouseY := realMousePos.Y / scale
		virtualMousePos := rl.NewVector2(virtualMouseX, virtualMouseY)

		return determineCursorState(bs, virtualMousePos, ps.GameViewWidth, ps.VirtualWidth, ps.VirtualHeight)
	}

	return spriteCursorPointer
}

// determineCursorState - powinno się wypierdzielić to do osobnego pliku, bo to logika, a nie rysowanie.
func determineCursorState(bs *battleState, mousePos rl.Vector2, viewW, totalW, viewH float32) uint16 {
	// A. Sprawdzamy "Scroll w prawo" na samej krawędzi okna
	if mousePos.X > totalW-scrollZoneXThreshold {
		return spriteCursorArrowRight
	}

	// B. Sprawdzamy czy mysz jest nad UI
	if mousePos.X >= viewW {
		return spriteCursorDefaultBig
	}

	// C. Pozostałe przewijanie
	if mousePos.X < scrollZoneXThreshold {
		return spriteCursorArrowLeft // Lewo
	}

	if mousePos.Y < scrollZoneYThreshold {
		return spriteCursorArrowUp // Góra
	}

	if mousePos.Y > viewH-scrollZoneYThreshold {
		return spriteCursorArrowDown // Dół
	}

	// D. Musimy przeliczyć pozycję myszy na świat gry używając kamery
	worldMousePos := rl.GetScreenToWorld2D(mousePos, bs.GameCamera)
	tileX := uint8(worldMousePos.X / float32(tileWidth))
	tileY := uint8(worldMousePos.Y / float32(tileHeight))

	if tileX < 0 || tileX >= boardMaxX || tileY < 0 || tileY >= boardMaxY {
		return spriteCursorStop
	}

	tile := &bs.Board.Tiles[tileX][tileY]
	targetOwner := -1

	var targetBuilding *building // Potrzebujemy wiedzieć, czy to budynek

	if tile.Unit != nil && tile.Unit.Exists {
		targetOwner = int(tile.Unit.Owner)
	} else if tile.Building != nil && tile.Building.Exists {
		targetOwner = int(tile.Building.Owner)
		targetBuilding = tile.Building
	}

	hasSelectedOwnUnit := bs.CurrentSelection.IsUnit && bs.CurrentSelection.OwnerID == bs.PlayerID
	canBeRepaired := targetBuilding != nil && targetBuilding.HP < targetBuilding.MaxHP &&
		(targetBuilding.Type == buildingPalisade || targetBuilding.Type == buildingBridge || targetBuilding.Owner == bs.PlayerID)

	if hasSelectedOwnUnit {
		// 0. Chcemy naprawić budynek
		if bs.MouseCommandMode == cmdRepairStructure {
			if canBeRepaired {
				return spriteBtnRepair
			}

			return spriteCursorStop
		}
		// 1. Wróg
		if targetOwner != -1 && targetOwner != int(bs.PlayerID) {
			// Specjalny przypadek: Palisada
			if targetBuilding != nil && targetBuilding.Type == buildingPalisade {
				selectedUnit, ok := getUnitByID(bs.CurrentSelection.UnitID, bs)
				if ok && !canDamagePalisades(selectedUnit) {
					return spriteCursorStop
				}
			}

			return spriteCursorCrossRed
		}

		// 2. Sojusznik (lub my sami)
		if targetOwner == int(bs.PlayerID) {
			return spriteCursorFrameWhite
		}

		// 3. Puste pole / Teren
		if targetOwner == -1 {
			if !tile.IsWalkable {
				return spriteCursorStop
			}

			return spriteCursorCrossWhite
		}
	} else {
		if targetOwner != -1 && targetOwner != int(bs.PlayerID) {
			return spriteCursorFrameRed
		}

		if targetOwner == int(bs.PlayerID) {
			return spriteCursorFrameWhite
		}
	}

	return spriteCursorDefaultBig
}

func animateCursorID(cursorID uint16) uint16 {
	const cursorAnimSpeed = 4.0
	animPhase := int(rl.GetTime()*cursorAnimSpeed) % 2

	if cursorID == spriteCursorCrossWhite {
		if rl.IsMouseButtonDown(rl.MouseRightButton) {
			return spriteCursorSmallWhite
		}

		if animPhase == 1 {
			return spriteCursorCrossMedWhite
		}
	}

	if cursorID == spriteCursorCrossRed {
		if rl.IsMouseButtonDown(rl.MouseRightButton) {
			return spriteCursorCrossMedRed
		}

		if animPhase == 1 {
			return spriteCursorCrossMedRed
		}
	}

	return cursorID
}

func drawCursorSprite(ps *programState, cursorID uint16, pos rl.Vector2, scale float32) {
	def := spriteRegistry[cursorID]
	if def.w == 0 {
		return
	}

	var tex rl.Texture2D

	var srcRect rl.Rectangle

	if cursorID == spriteCursorFrameWhite {
		tex = ps.Assets.CursorWhiteFrame
		if tex.ID == 0 {
			return
		}

		srcRect = rl.NewRectangle(0, 0, float32(def.w), float32(def.h))
	} else {
		tex = ps.Assets.getAtlas(def.atlasID, colorNone)
		if tex.ID == 0 {
			return
		}

		srcRect = rl.NewRectangle(float32(def.x), float32(def.y), float32(def.w), float32(def.h))
	}

	finalX := pos.X + (float32(def.offX) * scale)
	finalY := pos.Y + (float32(def.offY) * scale)

	destW := float32(def.w) * scale
	destH := float32(def.h) * scale

	ps.RenderDestRect = rl.NewRectangle(finalX, finalY, destW, destH)
	ps.RenderOrigin = rl.NewVector2(0, 0)
	rl.DrawTexturePro(tex, srcRect, ps.RenderDestRect, ps.RenderOrigin, 0, rl.White)
}

func drawMilkBarVisualizer(bs *battleState, ps *programState) {
	anchorX := ps.GameViewWidth
	barX := anchorX + milkBarOffsetX
	barY := milkBarY
	barWidth := milkBarWidth
	barHeight := milkBarHeight

	currentMilk := float32(bs.HumanPlayerState.Milk)
	fillRatio := currentMilk / maxMilk
	if fillRatio > 1.0 {
		fillRatio = 1.0
	}
	filledHeight := barHeight * fillRatio
	fillY := barY + barHeight - filledHeight

	rl.DrawRectangle(int32(barX), int32(fillY), int32(barWidth), int32(filledHeight), rl.White)

	capacityLineRation := float32(bs.HumanPlayerState.MaxMilk) / maxMilk
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

func drawGameUI(bs *battleState, ps *programState) {
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
	drawMilkBarVisualizer(bs, ps)
	drawButtons(bs, ps)

	// 3. Mapa
	minimapX := anchorX + minimapOffsetX
	minimapY := float32(0) + minimapOffsetY
	drawMinimapUnits(bs, minimapX, minimapY, minimapDisplayWidth, minimapDisplayHeight, ps.GameViewWidth)
}

// @todo: funkcja łamie zasadę rozdzielenia Dane-Logika-Rysowanie!
func drawButtons(bs *battleState, ps *programState) {
	// Punkt zaczepienia panelu (prawa strona ekranu)
	anchorX := ps.GameViewWidth

	// Pętla po 5 slotach przycisków
	for actionIndex := range uiActionMaxButtons {
		// 1. Pobieramy dane z battleState
		action := bs.UI.CurrentActions[actionIndex]

		// Jak nieaktywny, to nie rysujemy i nie klikamy
		if !action.IsActive {
			continue
		}

		// 2. Obliczamy pozycję (Layout)
		buttonY := uiAnchorOffsetY + (float32(actionIndex) * (btnHeight + btnMarginY))
		buttonX := anchorX + uiAnchorOffsetX
		rect := rl.NewRectangle(buttonX, buttonY, btnWidth, btnHeight)

		// Zapisujemy prostokąt do battleState (dla input.go)
		bs.UI.ActionButtons[actionIndex] = rect

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
			tex = ps.Assets.getAtlas(def.atlasID, bs.PlayerID)
			iconScale = 0.8
		case cmdBuildStructure, cmdStop, cmdMagicShield, cmdMagicSight, cmdRepairStructure:
			tex = ps.Assets.getAtlas(def.atlasID, colorNone)
			iconScale = 1

		case cmdMagicLightning, cmdMagicFire:
			// @todo: potrójna ikonka do wywołania tutaj
		default:
			// 4. Napis (Label)
			fontSize := int32(10)
			textWidth := rl.MeasureText(action.Label, fontSize)
			textX := rect.X + (rect.Width/2 - float32(textWidth)/2)
			textY := rect.Y + rect.Height - 12 // Na dole przycisku
			rl.DrawText(action.Label, int32(textX), int32(textY), fontSize, rl.Black)
		}

		// Wyznaczamy wycinek z atlasu (Source)
		srcRect := rl.NewRectangle(float32(def.x), float32(def.y), float32(def.w), float32(def.h))
		if def.flipX {
			srcRect.Width = -srcRect.Width
		}
		// Centrujemy ikonę na przycisku
		// Skalujemy ją, żeby zajmowała np. 80% wysokości przycisku
		scale := (rect.Height * iconScale) / float32(def.h)
		destW := float32(def.w) * scale
		destH := float32(def.h) * scale
		destX := rect.X + (rect.Width-destW)/2
		destY := rect.Y + (rect.Height-destH)/2
		rl.DrawTexturePro(tex, srcRect, rl.NewRectangle(destX, destY, destW, destH), rl.Vector2{}, 0, rl.White)
	}
}

// @todo: będzie wykorzystane przy rysowaniu ikonek już z ostatecznego układu
func drawTripleIcon(tex rl.Texture2D, btnRect rl.Rectangle) {
	// 1. Źródło (zakładamy standardową ikonkę 16x14)
	sourceRect := rl.NewRectangle(0, 0, 16, 14)

	// 2. ROZMIAR: 85% wysokości przycisku
	destH := btnRect.Height * 0.85
	destW := destH * (16.0 / 14.0) // Zachowanie proporcji

	// 3. NAKŁADANIE (overlap)
	// Krok przesunięcia to tylko 35% szerokości (ciasne upakowanie)
	stepX := destW * 0.35

	// 4. CENTROWANIE GRUPY w POZIOMIE
	totalGroupWidth := destW + (2 * stepX)
	groupStartX := btnRect.X + (btnRect.Width-totalGroupWidth)/2.0

	// 5. POZYCJE y (GÓRA / DÓŁ)
	paddingY := btnRect.Height - destH

	// Skrajne przy górze, środkowa przy dole
	yHigh := btnRect.Y
	yLow := btnRect.Y + paddingY

	// 6. RYSOWANIE
	// Lewa (Góra)
	rl.DrawTexturePro(tex, sourceRect,
		rl.NewRectangle(groupStartX, yHigh, destW, destH),
		rl.Vector2{}, 0, rl.White)

	// Środkowa (Dół) - Rysujemy drugą, żeby była pod spodem prawej, ale nad lewą (lub na odwrót zależnie od gustu)
	rl.DrawTexturePro(tex, sourceRect,
		rl.NewRectangle(groupStartX+stepX, yLow, destW, destH),
		rl.Vector2{}, 0, rl.White)

	// Prawa (Góra)
	rl.DrawTexturePro(tex, sourceRect,
		rl.NewRectangle(groupStartX+(2*stepX), yHigh, destW, destH),
		rl.Vector2{}, 0, rl.White)
}

func drawMinimapUnits(bs *battleState, minimapX, minimapY, minimapWidth, minimapHeight, actualGameViewWidth float32) {
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
			tileID := bs.Board.Tiles[x][y].TextureID
			col := getMapColor(tileID)
			w := xGridLines[x+1] - xGridLines[x]
			h := yGridLines[y+1] - yGridLines[y]
			rl.DrawRectangle(int32(xGridLines[x]), int32(yGridLines[y]), int32(w), int32(h), col)
		}
	}

	// 2. Budynki
	for _, bld := range bs.Buildings {
		if !bld.Exists {
			continue
		}

		// pomijamy palisady
		if bld.Type == buildingPalisade {
			continue
		}

		var bldColor rl.Color

		if bld.Owner == bs.PlayerID {
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
	for _, unit := range bs.Units {
		if !unit.Exists {
			continue
		}
		var unitColor rl.Color
		if unit.Owner == bs.PlayerID {
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
	camWorldView := getCameraWorldViewRect(bs.GameCamera, actualGameViewWidth, float32(virtualScreenHeight))
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

func drawConstructionDebugBox(bs *battleState, ps *programState) {
	if bs.MouseCommandMode != cmdBuildStructure {
		return
	}

	screenMouse := rl.GetMousePosition()
	virtualMouse := screenToVirtualCoords(ps, screenMouse)

	if isMouseOverUI(ps, virtualMouse) {
		return
	}

	worldMousePos := rl.GetScreenToWorld2D(virtualMouse, bs.GameCamera)
	tileX := uint8(worldMousePos.X / float32(tileWidth))
	tileY := uint8(worldMousePos.Y / float32(tileHeight))

	drawValidationBox(tileX, tileY, bs)

	if bs.PendingBuildingType != buildingPalisade && bs.PendingBuildingType != buildingBridge {
		drawValidationBox(tileX+1, tileY, bs)
		drawValidationBox(tileX+2, tileY, bs)
		drawValidationBox(tileX, tileY+1, bs)
		drawValidationBox(tileX+1, tileY+1, bs)
		drawValidationBox(tileX+2, tileY+1, bs)
		drawValidationBox(tileX, tileY+2, bs)
		drawValidationBox(tileX+1, tileY+2, bs)
		drawValidationBox(tileX+2, tileY+2, bs)
	}
}

func drawValidationBox(tileX, tileY uint8, bs *battleState) {
	// Walidacja
	var isValid bool

	if bs.PendingBuildingType != buildingBridge {
		isValid = isValidConstructionSite(tileX, tileY, 1, 1, bs)
	} else {
		isValid = isWithinBoard(tileX, tileY, bs) && isWaterTileOnly(bs.Board.Tiles[tileX][tileY].TextureID)
	}

	posX := float32(tileX) * float32(tileWidth)
	posY := float32(tileY) * float32(tileHeight)

	widthPx := float32(1) * float32(tileWidth)
	heightPx := float32(1) * float32(tileHeight)

	var color rl.Color
	if isValid {
		color = rl.Fade(rl.Green, validationAlpha)
	} else {
		color = rl.Fade(rl.Red, validationAlpha)
	}

	rl.DrawRectangle(int32(posX), int32(posY), int32(widthPx), int32(heightPx), color)
	// rl.DrawRectangleLines(int32(posX), int32(posY), int32(widthPx), int32(heightPx), rl.White)
}
