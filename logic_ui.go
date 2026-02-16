package main

// logic_ui.go

// updateActionButtons przelicza stan przycisków na podstawie obecnego zaznaczenia.
// Powinna być wywoływana w każdej klatce logicznej (updateGame).
func updateActionButtons(bs *battleState) {
	// 1. Czyścimy obecne działanie przycisków
	for btnIndex := range uiActionMaxButtons {
		bs.UI.CurrentActions[btnIndex] = uiAction{IsActive: false}
	}

	// 2. Sprawdzamy, czy jesteśmy właścicielami zaznaczonego obiektu
	// jeżeli nie, to nie ma nic do pokazania
	selected := bs.CurrentSelection

	if selected.OwnerID != bs.PlayerID {
		return
	}

	// 3. Przekazujemy dobór logiki do odpowiednich pomocników
	if selected.IsUnit {
		fillUnitActions(bs, selected.UnitID)
	} else if selected.BuildingID != 0 {
		fillBuildingActions(bs, selected.BuildingID)
	}
}

// Wypełnia przyciski na podstawie przepisu budynku.
func fillBuildingActions(bs *battleState, buildingID uint) {
	bld, ok := getBuildingByID(buildingID, bs)
	if !ok || !bld.Exists || bld.IsUnderConstruction {
		return
	}

	recipes, exists := buildingRecipes[bld.Type]
	if !exists {
		return
	}

	for rIndex := 0; rIndex < len(recipes) && rIndex < uiActionMaxButtons; rIndex++ {
		recipe := recipes[rIndex]

		// Warunek poziomu (np. Pastuch wymaga poziomu 26)
		if recipe.MinLevel <= bs.CurrentLevel {
			var cmd command

			// === ROZGAŁĘZIENIE LOGIKI ===
			// Sprawdzamy, czy przepis dotyczy budowy struktury, czy produkcji jednostki.

			if recipe.BuildingType != 0 {
				// PRZYPADEK 1: BUDOWA (np. Nowa Obora, Droga)
				// TargetBuildingID pełni tu rolę nośnika TYPU budynku (buildingType),
				// który zostanie przekazany do bs.PendingBuildingType w input.go.
				cmd = command{
					ActionType:          cmdBuildStructure,
					InteractionTargetID: uint(recipe.BuildingType),
					CommandCategory:     5, // 5 = brak konkretnego celu (tryb myszy)
				}
			} else {
				// PRZYPADEK 2: PRODUKCJA (np. Krowa, Drwal)
				// TargetBuildingID wskazuje na INSTANCJĘ budynku, który ma produkować (bld.ID).
				cmd = command{
					ActionType:          cmdProduce,
					InteractionTargetID: bld.ID,
					ProduceType:         recipe.UnitType,
					CommandCategory:     0, // 0 = budynek wykonujący rozkaz
				}
			}

			// Przypisanie gotowego rozkazu do UI
			bs.UI.CurrentActions[rIndex] = uiAction{
				IsActive: true,
				Label:    recipe.Label,
				IconID:   recipe.IconID,
				Cmd:      cmd,
			}
		}
	}
}

// Wypełnia przyciski na podstawie rodzaju jednostki.
func fillUnitActions(bs *battleState, unitID uint) {
	unit, ok := getUnitByID(unitID, bs)

	if !ok || !unit.Exists || unit.Owner != bs.PlayerID {
		return
	}

	// @todo: jest to nieprawidłowy rozkaz, powinno być „broń się”, czy coś podobnego
	// @reminder sprawdź, jak to wyglądało w pierwowzorze!
	bs.UI.CurrentActions[0] = uiAction{
		IsActive: true,
		Label:    "Stop",
		IconID:   spriteBtnShield,
		Cmd: command{
			ActionType:      cmdStop,
			ExecutorID:      unit.ID,
			CommandCategory: 1,
		},
	}
	// @todo: podejrzewam, że zamiast if-ów będzie potrzebny switch później
	if unit.Type == unitAxeman {
		bs.UI.CurrentActions[1] = uiAction{
			IsActive: true,
			Label:    "Napraw",
			IconID:   spriteBtnRepair,
			Cmd: command{
				ActionType:      cmdRepairStructure,
				ExecutorID:      unitID,
				CommandCategory: 1,
			},
		}
	}

	if unit.Type == unitPriestess {
		bs.UI.CurrentActions[1] = uiAction{
			IsActive: true, // @todo: czy powinien być widoczny również, gdy nie ma many?
			Label:    "Magiczna tarcza",
			IconID:   spriteBtnSpellMagicShield,
			Cmd: command{
				ActionType:      cmdMagicShield,
				ExecutorID:      unitID,
				CommandCategory: 1,
			},
		}
	}

	if unit.Type == unitPriestess {
		bs.UI.CurrentActions[2] = uiAction{
			IsActive: true, // @todo: czy powinien być widoczny również, gdy nie ma many?
			Label:    "Gromobicie",
			IconID:   spriteBtnSpellMagicLighting,
			Cmd: command{
				ActionType:      cmdMagicLightning,
				ExecutorID:      unitID,
				CommandCategory: 1,
			},
		}
	}

	if unit.Type == unitPriest {
		bs.UI.CurrentActions[1] = uiAction{
			IsActive: true, // @todo: czy powinien być widoczny również, gdy nie ma many?
			Label:    "Dalekie widzenie",
			IconID:   spriteBtnSpellVision,
			Cmd: command{
				ActionType:      cmdMagicSight,
				ExecutorID:      unitID,
				CommandCategory: 1,
			},
		}
	}

	if unit.Type == unitPriest {
		bs.UI.CurrentActions[2] = uiAction{
			IsActive: true, // @todo: czy powinien być widoczny również, gdy nie ma many?
			Label:    "Deszcz ognia",
			IconID:   spriteBtnSpellMagicFire, // @todo powinno być potrójne, a nie jeden!
			Cmd: command{
				ActionType:      cmdMagicFire,
				ExecutorID:      unitID,
				CommandCategory: 1,
			},
		}
	}
	// @todo: dodaj działania dla czarujących jednostek czy coś
}
