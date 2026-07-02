package main

// logic_ui.go

// updateActionButtons przelicza stan przycisków na podstawie obecnego zaznaczenia.
// Powinna być wywoływana w każdej klatce logicznej (updateGame).
func updateActionButtons(bState *battleState) {
	// 1. Czyścimy obecne działanie przycisków
	for btnIndex := range uiActionMaxButtons {
		bState.UI.CurrentActions[btnIndex] = uiAction{IsActive: false}
	}

	// 2. Sprawdzamy, czy jesteśmy właścicielami zaznaczonego obiektu
	// jeżeli nie, to nie ma nic do pokazania
	selected := bState.CurrentSelection

	if selected.OwnerID != bState.PlayerID {
		return
	}

	// 3. Przekazujemy dobór logiki do odpowiednich pomocników
	if selected.IsUnit {
		fillUnitActions(bState, selected.UnitID)
	} else if selected.BuildingID != 0 {
		fillBuildingActions(bState, selected.BuildingID)
	}
}

// Wypełnia przyciski na podstawie przepisu budynku.
func fillBuildingActions(bState *battleState, buildingID uint) {
	bld, ok := bState.getBuildingByID(buildingID)
	// Jeśli budynek został zniszczony lub jest w trakcie budowy, to nie może działać.
	if !ok || !bld.Exists || bld.IsUnderConstruction {
		return
	}

	// Sprawdzamy, czy zdefiniowano dla budynku przyciski.
	recipes, exists := buildingRecipes[bld.Type]
	// jeśli nie, to niestety nie ma czym wypełnić.
	if !exists {
		return
	}

	// Przechodzimy przez wszystkie przepisy dla danego rodzaju budynków.
	for rIndex := 0; rIndex < len(recipes) && rIndex < uiActionMaxButtons; rIndex++ {
		recipe := recipes[rIndex]

		// Warunek poziomu. Np. Pastuch jest dostępny dopiero od 26 „misji”.
		if recipe.MinLevel <= bState.CurrentLevel {
			// var cmd command
			action := uiAction{
				IsActive: true,
				Label:    recipe.Label,
				IconID:   recipe.IconID,
			}

			// === ROZGAŁĘZIENIE LOGIKI ===
			// Sprawdzamy, czy przepis dotyczy budowy struktury, czy produkcji jednostki.
			// @todo: czemu wykluczam nieistniejący buildingType? Nie pamiętam, ale źle to wygląda.
			if recipe.BuildingType != 0 {
				// PRZYPADEK 1: Zasadzenie nowej budowy
				// Jest to rozkaz wymagający dodatkowej informacji z planszy.
				action.State = mouseStatePlaceConstruction
				action.Cmd = command{
					ActionType: cmdBPlaceConstruction,
					CreateType: uint8(recipe.BuildingType),
				}
			} else {
				// PRZYPADEK 2: PRODUKCJA (np. Krowa, Drwal)
				// Jest to rozkaz możliwy do wykonania natychmiastowo.
				// Niczego z planszy nie potrzebujemy, dlatego stan myszy jest zwyczajny
				action.State = mouseStateNormal
				action.Cmd = command{
					ActionType: cmdBProduce, // cmdB oznacza, że to „budynkowy rozkaz”
					// na wypadek gdyby categoryBuilding nie było widoczne w kodzie i powstała wątpliwość
					ExecutorID: bld.ID,                 // tenże budynek ma wykonać rozkaz
					CreateType: uint8(recipe.UnitType), // rodzaj jednostki do wytworzenia
				}
			}
			// Przypisanie gotowego rozkazu do UI
			bState.UI.CurrentActions[rIndex] = action
		}
	}
}

// Wypełnia przyciski na podstawie rodzaju jednostki.
func fillUnitActions(bState *battleState, unitID uint) {
	currentUnit, ok := bState.getUnitByID(unitID)

	if !ok || !currentUnit.Exists || currentUnit.Owner != bState.PlayerID {
		return
	}

	// @todo: jest to nieprawidłowy rozkaz, powinno być „broń się”, czy coś podobnego
	// @reminder sprawdź, jak to wyglądało w pierwowzorze!
	bState.UI.CurrentActions[0] = uiAction{
		IsActive: true,
		Label:    "Stop",
		IconID:   spriteBtnShield,
		State:    mouseStateNormal, // rozkaz „natychmiastowy”
		Cmd: command{
			ActionType: cmdUStop,
		},
	}
	// @todo: podejrzewam, że zamiast if-ów będzie potrzebny switch później
	if currentUnit.Type == unitAxeman {
		bState.UI.CurrentActions[1] = uiAction{
			IsActive: true,
			Label:    "Napraw",
			IconID:   spriteBtnRepair,
			State:    mouseStateWorking, // rozkaz „złożony”
			Cmd: command{
				ActionType: cmdUWork, // cmdUWork, ponieważ tutaj nie możemy wiedzieć czy build czy repair
			},
		}
	}

	if currentUnit.Type == unitPriestess {
		bState.UI.CurrentActions[1] = uiAction{
			IsActive: true, // @todo: czy powinien być widoczny również, gdy nie ma many?
			Label:    "Magiczna tarcza",
			IconID:   spriteBtnSpellMagicShield,
			State:    mouseStateNormal, // rozkaz „natychmiastowy”
			Cmd: command{
				ActionType: cmdUCastSpell,
				Spell:      spellMagicShield,
				ExecutorID: unitID,
			},
		}
	}

	if currentUnit.Type == unitPriestess {
		bState.UI.CurrentActions[2] = uiAction{
			IsActive: true, // @todo: czy powinien być widoczny również, gdy nie ma many?
			Label:    "Gromobicie",
			IconID:   spriteBtnSpellMagicLighting,
			State:    mouseStateCasting,
			Cmd: command{
				ActionType: cmdUCastSpell,
				Spell:      spellMagicShower,
				ExecutorID: unitID,
			},
		}
	}

	if currentUnit.Type == unitPriest {
		bState.UI.CurrentActions[1] = uiAction{
			IsActive: true, // @todo: czy powinien być widoczny również, gdy nie ma many?
			Label:    "Dalekie widzenie",
			IconID:   spriteBtnSpellVision,
			State:    mouseStateCasting,
			Cmd: command{
				ActionType: cmdUCastSpell,
				Spell:      spellMagicSight,
				ExecutorID: unitID,
			},
		}
	}

	if currentUnit.Type == unitPriest {
		bState.UI.CurrentActions[2] = uiAction{
			IsActive: true, // @todo: czy powinien być widoczny również, gdy nie ma many?
			Label:    "Deszcz ognia",
			IconID:   spriteBtnSpellMagicFire,
			State:    mouseStateCasting,
			Cmd: command{
				ActionType: cmdUCastSpell,
				Spell:      spellMagicShower,
				ExecutorID: unitID,
			},
		}
	}
	// @todo: dodaj działania dla czarujących jednostek czy coś
}
