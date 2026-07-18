package main

// experience.go

/*
Tutaj trzymam wszystko związane z doświadczeniem i zdobywaniem doświadczenia
przez jednostki.
*/

// @reminder: O ile dobrze rozumiem, to ograniczenie doświadczenia było
// ustawione na 240, ale 224 było ostatnim „poziomem. Pomiędzy 224 a 240
// jest martwa przestrzeń, która nic już nie zmienia. Dlatego ustawiam
// ogranicznik na 224.
var (
	experienceCap         uint8 = 224
	experienceCasterBonus uint8 = 2
	experiencePerLevel    uint8 = 16
)

// @reminder: w pierwowzorze było 1, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 6, 7, 8, 9.
// Zmieniłem logikę z dodatku do każdego uderzenia na jednorazowe podniesienie
// statystyk. Dzięki temu upraszczam, ale wymaga zmian wartości w tablicy.
var levelUpBonusDamage = [15]uint16{0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1}

// @reminder: w pierwowzorze było 0, 1, 1, 2, 2, 2, 2, 3, 3, 3, 4, 4, 4, 5, 6
// Zmieniłem logikę z dodatku do zbroi na jednorazowe podniesienie
// statystyk. Dzięki temu upraszczam, ale wymaga zmian wartości w tablicy.
var levelUpBonusArmor = [15]uint16{0, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 1}

// @reminder: w pierwowzorze było 60, 80, 85, 90, 120, 140, 150, 160, 170, 180, 190, 200, 220, 240, 280.
// Górna granica many była nadpisywana wartością z tablicy, teraz jest podobnie, ale pierwsza
// jest już ustawiona przy tworzeniu jednostek, dlatego różni się od tego, co w komentarzu.
var levelUpBonusMana = [15]uint16{0, 80, 85, 90, 120, 140, 150, 160, 170, 180, 190, 200, 220, 240, 280}

func (u *unit) checkLevelUp() {
	// 0. Sprawdzamy poziom doświadczenia
	currentLevel := u.Experience / experiencePerLevel

	// 1. Sprawdzamy, czy został przekroczony próg, zwany też „tier”
	// Poziom doświadczenia się nie zmienił.
	// @reminder: nie da się obniżyć poziomu, to upraszcza logikę.
	if currentLevel == u.Level {
		return
	}

	// Poziom się podniósł, możemy zmieniać statystyki.
	for currentLevel > u.Level {
		u.levelUp()
	}
}

func (u *unit) levelUp() {
	u.Level++
	u.levelUpBonusDamage(u.Level)
	u.levelUpBonusArmor(u.Level)
	u.levelUpBonusMana(u.Level)
}

func (u *unit) levelUpBonusDamage(level uint8) {
	// Uwzględniamy dodatek za podniesienie poziomu
	u.Damage += levelUpBonusDamage[level]
}

func (u *unit) levelUpBonusArmor(level uint8) {
	// Uwzględniamy dodatek za podniesienie poziomu
	u.Armor += levelUpBonusArmor[level]
}

// Jest to najprostszy przypadek ponieważ tylko podmieniamy u.MaxMana na nową wartość.
func (u *unit) levelUpBonusMana(level uint8) {
	if !u.Type.isMagical() {
		return
	}

	// Uwzględniamy dodatek za podniesienie poziomu
	u.MaxMana = levelUpBonusMana[level]
}

// @reminder: 19.05.2026. Na podstawie pierwowzoru wydaje mi się, że rzeczywistą górną granicą doświadczenia jest
// 224. Dlatego dodałem stałą experienceCap, która służy za bezpiecznik.
// Doświadczenie jest zdobywane w chwili wyprowadzenia ataku.
func handleGainExperience(attacker, target *unit, PlayerID, AIPlayerID PlayerID) {
	// 1. Ustalamy, czy atakowana jest wroga jednostka lub wrogi budynek
	isEnemyUnit := target != nil && target.Owner != attacker.Owner
	// 1a. Ustalamy, czy jednostka może dostać doświadczenie
	//     ↓SI zawsze        ↓gracz tylko przy ataku wrogich jednostek
	canGainExp := attacker.Owner == AIPlayerID || (attacker.Owner == PlayerID && isEnemyUnit)
	// 1b. Jeśli nie to wracamy
	if !canGainExp {
		return
	}

	// 2a. Podstawowy przyrost doświadczenia. Dzięki sprawdzeniu w 0. zawsze możliwy.
	var expAmount uint8 = 1
	// 2b. Dodatek dla jednostek czarujących.
	if attacker.Type.isCaster() {
		expAmount += experienceCasterBonus
	}

	attacker.addExperience(expAmount)
}

func (u *unit) addExperience(amount uint8) {
	// 0. Sprawdzamy, czy osiągnęliśmy górną granicę.
	// TAK: nie obsługujemy zdobywania doświadczenia
	if u.Experience >= experienceCap {
		return
	}
	// NIE: jednostka zdobywa doświadczenie
	u.Experience += amount

	// 1. Bezpiecznik, nie można mieć więcej niż
	// górna granica doświadczenia.
	if u.Experience > experienceCap {
		u.Experience = experienceCap
	}

	u.checkLevelUp()
}
