package main

import (
	"testing"
)

// UWAGA!!!! Próby wytworzone przez sztuczną inteligencję. Nie zostały jeszcze sprawdzone.
// Mogą zawierać błędy i niedociągnięcia! - 16.10.2026

// unit_testing_helpers_test.go
//
// Wspólne środowisko testowe (test fixtures & helpers) dla pakietu testów jednostek
// (units_test.go, unit_commands_test.go, unit_actions_test.go, unit_movement_test.go).
//
// IZOLACJA I OGRANICZENIA WSPÓŁBIEŻNOŚCI:
// Testy korzystające z planszy (boardData) lub symulacji bitwy (battleState) NIE mogą
// być uruchamiane równolegle za pomocą `t.Parallel()`.
// Wynika to z globalnych buforów współdzielonych w celach optymalizacji alokacji w silniku gry:
// - `sharedNodePool` oraz `validCoordsBuffer` w `pathfinding.go` (wyszukiwanie ścieżek A*),
// - `magicShowerProjsBuffer` w `unit_actions.go` (tworzenie pocisków czarów obszarowych).
// Każdy test powinien alokować świeżą instancję przez `newTestBattleState()` lub `newTestBoard()`,
// co gwarantuje pełną izolację stanu planszy i jednostek pomiędzy kolejno wykonywanymi testami.

const (
	testHumanPlayer PlayerID = colorRed
	testAIPlayer    PlayerID = colorYellow
)

// newTestBoard alokuje nową, czystą planszę testową o wymiarach boardMaxX x boardMaxY.
// Każdy kafelek jest zainicjalizowany jako w pełni przechodni (IsWalkable = true),
// z domyślnym kosztem ruchu 1, bez przypisanych jednostek, budynków, ognia czy zniszczeń.
func newTestBoard() *boardData {
	board := &boardData{}

	for x := range boardMaxX {
		for y := range boardMaxY {
			t := &board.Tiles[x][y]
			t.X = x
			t.Y = y
			t.IsWalkable = true
			t.MovementCost = 1
		}
	}

	return board
}

// newTestBattleState tworzy świeży, minimalny stan bitwy (battleState) z nową planszą testową
// oraz zainicjalizowanymi stanami gracza ludzkiego (Red) i przeciwnika SI (Yellow).
// Zapewnia czyste środowisko do testowania logiki rozkazów, akcji i ruchu jednostek.
func newTestBattleState() *battleState {
	return &battleState{
		PlayerID:         testHumanPlayer,
		AIPlayerID:       testAIPlayer,
		HumanPlayerState: &playerState{PlayerID: testHumanPlayer},
		AIEnemyState:     &playerState{PlayerID: testAIPlayer},
		Board:            newTestBoard(),
	}
}

// newTestUnit tworzy instancję jednostki o zadanym identyfikatorze, typie, pozycji startowej i właścicielu.
// Nadaje jednostce pełne zdrowie (HP=MaxHP=20), stan spoczynku (stateIdle, cmdUIdle),
// domyślny zasięg widzenia (8), zasięg ataku (1) oraz domyślny punkt podejścia (Approach) na własnej pozycji.
func newTestUnit(id unitID, uType unitType, x, y uint8, owner PlayerID) *unit {
	return &unit{
		ID:          id,
		Exists:      true,
		X:           x,
		Y:           y,
		Owner:       owner,
		Type:        uType,
		HP:          20,
		MaxHP:       20,
		Command:     cmdUIdle,
		State:       stateIdle,
		SightRange:  8,
		AttackRange: 1,
		Damage:      5,
		Armor:       1,
		MaxDelay:    3,
		Approach:    point{X: x, Y: y},
	}
}

// placeUnit rejestruje jednostkę w kolekcji jednostek stanu bitwy (bState.Units)
// oraz ustawia wskaźnik na nią na odpowiednim kafelku planszy (Tiles[u.X][u.Y].Unit).
// Używane do przygotowywania geometrii testowej.
func placeUnit(tb testing.TB, bState *battleState, u *unit) {
	tb.Helper()

	bState.Units = append(bState.Units, u)
	bState.Board.Tiles[u.X][u.Y].Unit = u
}

// buildingFootprint wylicza listę kafelków zajmowanych przez budynek danego typu
// o zadanym lewym górnym rogu (topLeft).
// Dla obiektów jednokafelkowych (palisada, most) zwraca pojedynczy punkt.
// Dla standardowych budynków 3x3 generuje 9 kafelków w kolejności kolumnowej.
func buildingFootprint(bType buildingType, topLeft point) []point {
	if bType == buildingPalisade || bType == buildingBridge {
		return []point{topLeft}
	}

	tiles := make([]point, 0, 9)

	for dx := range uint8(3) {
		for dy := range uint8(3) {
			tiles = append(tiles, point{X: topLeft.X + dx, Y: topLeft.Y + dy})
		}
	}

	return tiles
}

// newTestBuilding tworzy instancję budynku o zadanym typie, pozycji lewego górnego rogu i właścicielu.
// Budynek ma pełne punkty życia (HP=MaxHP=100) oraz wyliczoną listę zajmowanych kafelków (OccupiedTiles).
func newTestBuilding(id buildingID, bType buildingType, topLeft point, owner PlayerID) *building {
	return &building{
		ID:            id,
		Exists:        true,
		Owner:         owner,
		Type:          bType,
		HP:            100,
		MaxHP:         100,
		OccupiedTiles: buildingFootprint(bType, topLeft),
	}
}

// placeBuilding rejestruje budynek w stanie bitwy (bState.Buildings) oraz oznacza wszystkie
// zajmowane przez niego kafelki planszy jako nieprzechodnie (IsWalkable = false) i wskazuje budynek.
func placeBuilding(tb testing.TB, bState *battleState, bld *building) {
	tb.Helper()

	bState.Buildings = append(bState.Buildings, bld)

	for _, p := range bld.OccupiedTiles {
		t := &bState.Board.Tiles[p.X][p.Y]
		t.Building = bld
		t.IsWalkable = false
	}
}

// unitTargetRef tworzy referencję do celu jednostkowego (targetUnit) ze wskazaną pozycją i ID jednostki.
func unitTargetRef(u *unit) targetReference {
	return targetReference{
		Kind:     targetUnit,
		ID:       uint(u.ID),
		Position: point{X: u.X, Y: u.Y},
	}
}

// buildingTargetRef tworzy referencję do celu budynkowego (targetBuilding) z wyznaczonym środkiem budynku.
func buildingTargetRef(bld *building) targetReference {
	cx, cy, _ := bld.getCenter()

	return targetReference{
		Kind:     targetBuilding,
		ID:       uint(bld.ID),
		Position: point{X: cx, Y: cy},
	}
}

// tileTargetRef tworzy referencję do celu terenowego / kafelka (targetTile) na podanych współrzędnych.
func tileTargetRef(x, y uint8) targetReference {
	return targetReference{
		Kind:     targetTile,
		Position: point{X: x, Y: y},
	}
}

// makeStandingTree konfiguruje kafelek jako stojące drzewo (nieprzechodnie, ze stanem treeStraight).
// Parametr dry określa, czy drzewo ma teksturę suchego (ścięcie przez topornika), czy zielonego.
func makeStandingTree(t *tile, dry bool) {
	t.treeState = treeStraight
	t.IsWalkable = false

	if dry {
		t.TextureID = spriteDryTreeStump00
	} else {
		t.TextureID = spriteDryTreeStump00 + 1
	}
}

// chebyshev oblicza odległość Czebyszewa (odległość w siatce z ruchem po skosie: max(|dx|, |dy|))
// pomiędzy dwoma punktami na planszy.
func chebyshev(a, b point) uint8 {
	dx := max(int(a.X), int(b.X)) - min(int(a.X), int(b.X))
	dy := max(int(a.Y), int(b.Y)) - min(int(a.Y), int(b.Y))

	return uint8(max(dx, dy))
}

// assertState weryfikuje, czy jednostka znajduje się w oczekiwanym stanie (unitState)
// oraz posiada oczekiwany rozkaz (commandType). W razie niezgodności raportuje błąd w teście.
func assertState(tb testing.TB, u *unit, wantState unitState, wantCommand commandType) {
	tb.Helper()

	if u.State != wantState {
		tb.Errorf("State = %d, want %d", u.State, wantState)
	}

	if u.Command != wantCommand {
		tb.Errorf("Command = %d, want %d", u.Command, wantCommand)
	}
}

// mustPanic sprawdza, czy przekazana funkcja zgłasza panikę (panic).
// Jeśli funkcja wykona się bez paniki, asercja testowa zgłasza błąd.
func mustPanic(tb testing.TB, fn func()) {
	tb.Helper()

	defer func() {
		if recover() == nil {
			tb.Errorf("expected panic, got none")
		}
	}()

	fn()
}
