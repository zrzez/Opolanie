package main

// tileHelpers.go

// === Pomagierzy do rysowania

// Sprawcza, czy tekstura jest kapliczką leczącą (świętym miejscem).
func isHealingShire(tileTexID uint16) bool {
	return tileTexID == spriteEffectHeal00 || tileTexID == spriteEffectHeal01
}

// Sprawdza, czy tekstura jest punktem zwycięstwa lub kapliczką.
func isSpecialTile(tileTexID uint16) bool {
	return tileTexID == spriteEffectHeal00 || tileTexID == spriteEffectHeal01 || tileTexID == spriteVictoryPoint
}

// Sprawdza, czy tekstura jest drogą.
func isDirtRoad(tileTexID uint16) bool {
	return tileTexID >= spriteRoadStart && tileTexID <= spriteRoadEnd
}

// Sprawdza, czy tekstura jest ukończonym mostem. Kierunek nie ma znaczenia.
// Nie ma znaczenia, czy grasz postawił.
// @todo: sprawdź, czemu nieukończony most jest traktowany inaczej - 21.04.2026.
func isCompletedBridge(tileTexID uint16) bool {
	return tileTexID >= spriteBridgeStart && tileTexID <= spriteBridgeEnd
}

// Sprawdza, czy tekstura jest drogą lub mostem.
func isPath(tileTexID uint16) bool {
	return isDirtRoad(tileTexID) || isCompletedBridge(tileTexID)
}

// Sprawdza, czy tekstura to zgliszcza budynku.
func isRuin(tileTexID uint16) bool {
	return tileTexID >= spriteRuinStart && tileTexID <= spriteRuinEnd
}

// Sprawdza, czy tekstura jest palisadą. Stan ukończenia nie ma znaczenia.
func isPalisade(tileTexID uint16) bool {
	return tileTexID >= spritePalisadeStart && tileTexID <= spritePalisadeEnd
}

// Sprawdza, czy tekstura jest wodą. Stan animacji nie ma znaczenia.
func isWaterTileOnly(tileTexID uint16) bool {
	return tileTexID >= spriteWaterStart && tileTexID <= spriteWaterEnd
}

// Sprawdza, czy tekstura jest wodą lub ukończonym mostem.
func isWaterOrBridgeForMasking(tileTexID uint16) bool {
	return isWaterTileOnly(tileTexID) || isCompletedBridge(tileTexID) || tileTexID == spriteBridgeConstruction
}

// Sprawdza, czy tekstura jest „lądowa”, czyli nie woda, nie mosty.
func isLandOrOther(tileTexID uint16) bool {
	return !isWaterTileOnly(tileTexID) && !isCompletedBridge(tileTexID) && tileTexID != spriteBridgeConstruction
}

// Sprawdza, czy tekstura jest to teren pod budynkami lub placem budowy.
func isBuildingTerrain(tileTexID uint16) bool {
	return (tileTexID >= spriteConstructionStart && tileTexID <= spriteConstructionEnd) ||
		(tileTexID >= spriteBuildingMainBase && tileTexID <= spriteBuildingEnd)
}

// Sprawdza, czy tekstura jest to sucha ziemia, często stosowana przy mostach.
func isDryEarth(tileTexID uint16) bool {
	return tileTexID >= spriteDryEarth01 && tileTexID <= spriteDryEarth03
}

// Sprawdza, czy tekstura jest to trawa. Stopień wyżarcia przez krowy, czy spopielenia nie ma znaczenia.
func isGrass(tileTexID uint16) bool {
	return tileTexID >= spriteGrassStart && tileTexID <= spriteGrassEnd
}

// Sprawdza, czy tekstura jest ozdoba terenu, jak ołtarz, drogowskaz, studnia.
func isGadget(tileTexID uint16) bool {
	return tileTexID >= spriteGadgetStart && tileTexID <= spriteGadgetEnd
}

// Sprawdza, czy tekstura jest trawa lub ozdoba/
func isGrassOrGadget(tileTexID uint16) bool {
	return isGrass(tileTexID) || isGadget(tileTexID)
}

// Sprawdza, czy tekstura jest
func isRockNonWalkable(tileTexID uint16) bool {
	return tileTexID >= spriteRockStart && tileTexID < spriteRockEnd
}

// Sprawdza, czy tekstura jest pień drzewa, w tym suchego. Nie mylić z pieńkiem, który jest ozdobą.
func isTreeStump(tileTexID uint16) bool {
	return tileTexID >= spriteTreeStumpStart && tileTexID <= spriteTreeStumpEnd
}

// Sprawdza, czy tekstura jest drzewem. W tym obalonym.
func isTree(tileTexID uint16) bool {
	return tileTexID >= spriteTreeStumpStart && tileTexID <= spriteTreeTopEnd ||
		tileTexID >= spriteTreeBurntStump00 && tileTexID <= spriteTreeFallingBurnt02
}
