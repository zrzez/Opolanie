package main

import (
	"math"
)

// missile.go

// Rodzaje pocisków
const (
	missileArrow     = "ARROW"
	missileBolt      = "CROSSBOW"
	missileFire      = "FIRE"
	missileLightning = "LIGHTNING"
	missileSpear     = "SPEAR"
	missileGhost     = "GHOST"
)

type projectile struct {
	ID    uint
	Type  string // Klucz do mapy ProjectileSprites
	Owner uint8  // Kto wystrzelił

	// Pozycja
	X, Y float32

	// Wektor ruchu na jedną klatkę (20 TPS)
	DX, DY float32

	// Cel
	TargetX, TargetY uint8

	// Licznik czasu lotu
	Lifetime uint

	// Obrażenia
	Damage uint16

	// Stan
	Exists bool
}

func (p *projectile) initProjectile(typeName string, owner uint8, startX, startY, targetX, targetY uint8, damage uint16) {
	p.Type = typeName
	p.Owner = owner
	p.TargetX = targetX
	p.TargetY = targetY
	p.Damage = damage
	p.Exists = true

	// Przeliczenie kafelków na piksele (środek kafelka)
	p.X = float32(startX*tileWidth) + float32(tileWidth)/2
	p.Y = float32(startY*tileHeight) + float32(tileHeight)/2

	destPixelX := float32(targetX*tileWidth) + float32(tileWidth)/2
	destPixelY := float32(targetY*tileHeight) + float32(tileHeight)/2

	// Obliczenie wektora różnicy
	diffX := destPixelX - p.X
	diffY := destPixelY - p.Y

	// Normalizacja prędkości
	// ZMIANA: Kusznik (BOLT) ma szybszy pocisk, reszta standardowy.
	// To przywraca dynamikę z oryginału, gdzie bełt leciał 2x szybciej.
	var speed float32 = 4.0
	if typeName == missileBolt {
		speed = 8.0
	}

	distance := float32(math.Sqrt(float64(diffX*diffX + diffY*diffY)))

	if distance > 0 {
		p.DX = (diffX / distance) * speed
		p.DY = (diffY / distance) * speed
		// Lifetime to czas potrzebny na dolot
		p.Lifetime = uint(int(distance / speed))
	} else {
		p.Exists = false // Cel tożsamy ze startem
	}
}

// updateProjectile aktualizuje pozycję pocisku.
func (p *projectile) updateProjectile(bs *battleState) {
	if !p.Exists {
		return
	}

	// Ruch
	p.X += p.DX
	p.Y += p.DY
	p.Lifetime--

	// Sprawdzenie trafienia (gdy czas się skończył)
	if p.Lifetime <= 0 {
		p.hit(bs)
		p.Exists = false
	}
}

// hit zadaje obrażenia w punkcie docelowym.
func (p *projectile) hit(bs *battleState) {
	// Sprawdzamy granice mapy
	if p.TargetX < boardMaxX && p.TargetY < boardMaxY {

		// === ZMIANA ARCHITEKTONICZNA ===
		// Pobieramy wskaźnik do kafelka (Tile) w miejscu trafienia.
		tile := &bs.Board.Tiles[p.TargetX][p.TargetY]

		// Zamiast szukać ID w starej tablicy Plc i potem wołać getObjectByID,
		// sprawdzamy bezpośrednio, czy na kafelku jest wskaźnik do unit lub building.

		// 1. Trafienie w Jednostkę
		if tile.Unit != nil && tile.Unit.Exists && tile.Unit.Owner != p.Owner {
			tile.Unit.takeDamage(p.Damage, bs)
		}

		// 2. Trafienie w Budynek
		if tile.Building != nil && tile.Building.Exists && tile.Building.Owner != p.Owner {
			tile.Building.takeDamage(p.Damage)
		}

		// 3. Efekt wizualny (np. ogień)
		// Dawniej: bs.Board.PlcAttackEffects[x][y] = 100
		// Teraz: tile.EffectID = 100
		if p.Type == missileFire {
			// tile.EffectID = 100
		}
	}
}
