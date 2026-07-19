package main

import (
	"log"
)

// constructions.go

func (bld *building) repair(amount uint16) {
	if !bld.Exists || bld.HP >= bld.MaxHP {
		return
	}

	bld.increaseHP(amount)
}

func (bld *building) build(amount uint16) {
	if !bld.Exists || !bld.IsUnderConstruction {
		return
	}

	bld.increaseHP(amount)
}

// increaseHP dla każdej istniejącej budowli zwiększa PŻ o amount
// Pilnuje, aby bld.HP <= bld.MaxHP; Służy do naprawy budynków.
func (bld *building) increaseHP(amount uint16) {
	if !bld.Exists {
		return
	}

	bld.HP += amount
	if bld.HP >= bld.MaxHP {
		bld.HP = bld.MaxHP
	}
}

// Obniżamy bezpośrednio HP budynku. Nie mylić processingBuildingDamage, które
// służa sa sito dla nieskutecznych ataków.
func (bld *building) applyBuildingDamage(finalDamage uint16) {
	// Bez tego bld.HP przekręca się na 65 tys.
	if bld.HP >= finalDamage {
		bld.HP -= finalDamage
	} else {
		bld.HP = 0
	}

	log.Printf("Budynek %s przyjął obrażenia w wysokości: %d. HP: %d/%d",
		buildingDefs[bld.Type].Name, finalDamage, bld.HP, bld.MaxHP)
}

func (bldType buildingType) isRegularBuilding() bool {
	return bldType != buildingRoad && bldType != buildingBridge && bldType != buildingPalisade
}

// Odpowiada za dodanie stworzonej jednostki do mieszkańców budynku.
// Sprawdzenie zostało wykonane w canProduceUnit,
// pojemność za pomocą hasRoom
// Zwracany bool jest całkowicie ignorowany.
func (bld *building) registerUnit(uID UnitID) bool {
	// Budynek poszerza listę zameldowanych jednostek
	// To chyba powinny być wskaźniki
	// @todo: sprawdź po cholerę mi w ogole ta lista
	// przecież samo śledzenie licznika domowników można zrobić prościej
	// @reminder: jeśli znajdę sposób na przypisane krowy do obory
	// najprawdopodobniej zniknie jedyna przyczyna dla której ta lista istnieje
	bld.AssignedUnits = append(bld.AssignedUnits, uID)
	// Budynek aktualizuje licznik posiadania
	bld.Food++

	return true
}

// Wywoływana przez u.unregisterFromBuilding gdy jednostka zmarła.
// @todo: co się dzieje z jednostkami, które sotajemy na początku bitwy?
func (bld *building) unregisterUnit(unregisterUnitID UnitID) {
	// Przechodzimy przez listę jednostek zamieszkujących
	for index, registeredUnitID := range bld.AssignedUnits {
		if registeredUnitID == unregisterUnitID {
			// po znalezieniu miejsca w którym znajduje się jednostka, pomijamy ją przy
			// odświeżeniu listy
			bld.AssignedUnits = append(bld.AssignedUnits[:index], bld.AssignedUnits[index+1:]...)

			bld.Food--

			// @reminder: próba sprawdzenia, czy to powoduje błąd
			return
		}
	}
}

func (bld *building) hasRoom() bool {
	return len(bld.AssignedUnits) < int(bld.MaxFood)
}

// @reminer: Zupełnie nie rozumiem po co obecnie miałbym mieć taką metodę.
// Do czasu aż nie ogarnę units.go zostawię, ale czuję, że jest zbędna.
func (bld *building) getCenter() (uint8, uint8, bool) {
	// @reminder: celowo nie daję każdego rodzaju budynku w przłączniku.
	switch bld.Type {
	case buildingPalisade, buildingBridge:
		// Te rodzaje budynków, zawsze mają dokładnie jeden kafelek
		return bld.OccupiedTiles[0].X, bld.OccupiedTiles[0].Y, true
	default:
		// Zwyczajne budowle zawsze są 3na3 więc środek jest z góry znany
		return bld.OccupiedTiles[4].X, bld.OccupiedTiles[4].Y, true
	}
}

// @reminder: jak to powinno działać -- budynek gromadzi obrażenia w bld.AccumnulatedDamage
// i jeśli w danym tiku przekroczyły one próg, to są zadawane „zbiorowo”.
// Jeśli się tego progu nie przekroczyło, to budynek zostaje nienaruszony.
// Tutaj mamy samo zbieranie, wewnątrz updateBuildings jest logika
// „rozliczania tiku”.
func (bld *building) takeDamage(damage uint16) {
	// Od sprawdzenia poprawności rozkazu do jego wykonania
	// budynek mógł już zostać zniszczony, dlatego zostawiam bezpiecznik.
	// W tej chwili nie jestem wstanie udowodnić, że jest zbędny.
	// Jednakże w bState.updateBuildings() najpierw ruszają się jednostki,
	// potem przetwarzamy obrażenia zadane budynkom.
	// Dopiero po tym zaczynamy sprzątać zniszczone budynki i budowy.
	if !bld.Exists {
		return
	}

	bld.AccumulatedDamage += damage
	log.Printf("Budynek %d otrzymał %d obrażen (łącznie: %d)", bld.ID, damage, bld.AccumulatedDamage)
}

func (bld *building) isRepairable(playerID PlayerID) bool {
	if bld == nil || !bld.Exists || bld.HP >= bld.MaxHP {
		return false
	}

	return bld.Type == buildingPalisade || bld.Type == buildingBridge || bld.Owner == playerID
}

func (bld *building) unassignUnitsfromBuilding(resolver unitResolver) {
	// Trzeba dać znać jednostkom, że nie mają już domu
	for _, uID := range bld.AssignedUnits {
		u, ok := resolver.getUnitByID(uID)

		if ok && u.Exists {
			u.BelongsTo = nil
		}
	}
}
