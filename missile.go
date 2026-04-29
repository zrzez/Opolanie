package main

import (
	"math"
)

// missile.go

// Rodzaje pocisków
const (
	missileArrow uint8 = iota
	missileBolt
	missileFire
	missileLightning
	missileSpear
	missileGhost
)

type projectile struct {
	ID     uint
	Kind   uint8
	Sprite uint16
	Owner  uint8

	// Pozycja
	X, Y float32

	// Wektor ruchu na jedną klatkę (20 TPS)
	DX, DY float32

	// Cel
	TargetX, TargetY uint16

	// Licznik czasu lotu
	Lifetime uint

	// Obrażenia
	Damage uint16 // @reminder: to nie jest używane przez pociski unitPriest, uniteMage. Nie wiem, jak kapłanka.

	// Stan
	Exists bool
}

func (p *projectile) initProjectile(kind, owner uint8, startX, startY, targetX, targetY uint16, damage uint16) {
	p.Kind = kind
	p.Owner = owner
	p.TargetX = targetX
	p.TargetY = targetY
	p.Damage = damage
	p.Exists = true
	// p.IsImpact = false

	// Przeliczenie kafelków na piksele (środek kafelka)
	p.X = float32(startX*uint16(tileWidth)) + float32(tileWidth)/2
	p.Y = float32(startY*uint16(tileHeight)) + float32(tileHeight)/2

	destPixelX := float32(targetX*uint16(tileWidth)) + float32(tileWidth)/2
	destPixelY := float32(targetY*uint16(tileHeight)) + float32(tileHeight)/2

	// Obliczenie wektora różnicy
	diffX := destPixelX - p.X
	diffY := destPixelY - p.Y

	// Ustawienie prędkości
	// kusznik pruł, jak szalony
	var speed float32 = 4.0
	if kind == missileBolt {
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

	p.Sprite = resolveProjectileSprite(kind, p.DX, p.DY)
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
	}
}

// hit zadaje obrażenia w punkcie docelowym.
func (p *projectile) hit(bs *battleState) {
	if p.TargetX >= uint16(boardMaxX) || p.TargetY >= uint16(boardMaxY) {
		p.Exists = false

		return
	}

	targetTile := &bs.Board.Tiles[p.TargetX][p.TargetY]

	// 1. Trafienie jednostki
	if targetTile.Unit != nil && targetTile.Unit.Exists && targetTile.Unit.Owner != p.Owner {
		targetTile.Unit.takeDamage(p.Damage, bs) // @todo: czemu jednostka ma argument bs, a budynek nie?
	}

	// 2. Trafienie budynku
	if targetTile.Building != nil && targetTile.Building.Exists && targetTile.Building.Owner != p.Owner {
		targetTile.Building.takeDamage(p.Damage)
	}

	// 3. Efekty dla specjalnych przypadków
	// @todo: sprawdź, czy dobrze rozumiem, że wywołując efekt specjalny po zadaniu obrażeń
	// uniemożliwiam prawidłowe zadanie obrażeń odpryskami ognia?
	// @todo: dodaj przekazywanie obrażeń, jako argumentu ponieważ potrzebuję tego do ducha.
	p.specialProjectiles(targetTile, bs)

	// @reminder: to się gryzie z duchem, który musi przetrwać uderzenie i dusić cel
	p.Exists = false
}

func (p *projectile) specialProjectiles(targetTile *tile, bs *battleState) {
	// 1. Jeśli to duch (pocisk maga) to zostań widoczny na jednostce
	// @todo: teraz się tym zajmuję 28.04.2026
	// @todo: potrzebuję przekazać tutaj doświadczenie i właściciela maga, żeby poprawnie wyliczyć
	// dodatek do obrażeń.
	if p.Kind == missileGhost {
		p.Sprite = spriteMissileGhostAttack
		p.mageGhost(targetTile, p.Damage, bs)

		return
	}

	// 2. Ogień musi palić się przez jakiś czas
	// @todo: czy można uwspólnić logikę ognia i ducha?
	if p.Kind == missileFire {
		p.priestFireball(targetTile, bs)
	}
}

func (p *projectile) spriteToDirection() (dirX int16, dirY int16) {
	switch p.Sprite {
	case spriteMissileArrowUp, spriteMissileBoltUp, spriteMissileFireUp,
		spriteMissileLightningUp, spriteMissileSpearUp, spriteMissileGhostUp:
		return 0, -1
	case spriteMissileArrowUpLeft, spriteMissileBoltUpLeft, spriteMissileFireUpLeft,
		spriteMissileLightningUpLeft, spriteMissileSpearUpLeft, spriteMissileGhostUpLeft:
		return -1, -1
	case spriteMissileArrowLeft, spriteMissileBoltLeft, spriteMissileFireLeft,
		spriteMissileLightningLeft, spriteMissileSpearLeft, spriteMissileGhostLeft:
		return -1, 0
	case spriteMissileArrowDownLeft, spriteMissileBoltDownLeft, spriteMissileFireDownLeft,
		spriteMissileLightningDownLeft, spriteMissileSpearDownLeft, spriteMissileGhostDownLeft:
		return -1, 1
	case spriteMissileArrowDown, spriteMissileBoltDown, spriteMissileFireDown,
		spriteMissileLightningDown, spriteMissileSpearDown, spriteMissileGhostDown:
		return 0, 1
	case spriteMissileArrowDownRight, spriteMissileBoltDownRight, spriteMissileFireDownRight,
		spriteMissileLightningDownRight, spriteMissileSpearDownRight, spriteMissileGhostDownRight:
		return 1, 1
	case spriteMissileArrowRight, spriteMissileBoltRight, spriteMissileFireRight,
		spriteMissileLightningRight, spriteMissileSpearRight, spriteMissileGhostRight:
		return 1, 0
	case spriteMissileArrowUpRight, spriteMissileBoltUpRight, spriteMissileFireUpRight,
		spriteMissileLightningUpRight, spriteMissileSpearUpRight, spriteMissileGhostUpRight:
		return 1, -1
	}

	// nie powinno się wydarzyć
	return 0, 0
}

func unitTypeToMissileType(unitType unitType) uint8 {
	switch unitType {
	case unitArcher:
		return missileArrow
	case unitSpearman:
		return missileSpear
	case unitPriestess:
		return missileLightning
	case unitPriest:
		return missileFire
	case unitMage:
		return missileGhost
	case unitCrossbowman:
		return missileBolt
	default:
		return missileArrow
	}
}

func resolveProjectileSprite(kind uint8, dx, dy float32) uint16 {
	// Bazowy duszek dla pocisku
	var baseSprite uint16

	switch kind {
	case missileArrow:
		baseSprite = spriteMissileArrowUp
	case missileSpear:
		baseSprite = spriteMissileSpearUp
	case missileLightning:
		baseSprite = spriteMissileLightningUp
	case missileFire:
		baseSprite = spriteMissileFireUp
	case missileGhost:
		baseSprite = spriteMissileGhostUp
	case missileBolt:
		baseSprite = spriteMissileBoltUp
	default:
		baseSprite = spriteMissileArrowUp
	}

	// Określenie kierunku
	dirX := 0

	if dx > 0.5 {
		dirX = 1
	} else if dx < -0.5 {
		dirX = -1
	}

	dirY := 0

	if dy > 0.5 {
		dirY = 1
	} else if dy < -0.5 {
		dirY = -1
	}

	var offset uint16

	switch {
	case dirX == 0 && dirY == -1:
		offset = 0 // góra 1
	case dirX == -1 && dirY == -1:
		offset = 1 // góra lewo 1
	case dirX == -1 && dirY == 0:
		offset = 2 // lewo 1
	case dirX == -1 && dirY == 1:
		offset = 3 // lewy dół 1
	case dirX == 0 && dirY == 1:
		offset = 4 // dół 1
	case dirX == 1 && dirY == -1:
		offset = 5 // góra prawo X
	case dirX == 1 && dirY == 0:
		offset = 6 // prawo X
	case dirX == 1 && dirY == 1:
		offset = 7 // prawy dół 1
	default:
		offset = 2 // na wszelki wypadek
	}

	return baseSprite + offset
}

// @todo: sprawdź, czy odpryski rzeczywiście zadają obrażenia, które powinny
func (p *projectile) priestFireball(affedtedTile *tile, bs *battleState) {
	dx, dy := p.spriteToDirection()

	// @reminder: szkoda, że nie zapisałem dokładnie czemu wyciszyłem tutaj lintera
	splash1X := int16(p.TargetX) + dx      //nolint:gosec
	splash1Y := int16(p.TargetY) + dy      //nolint:gosec
	splash2X := int16(p.TargetX) + dx + dx //nolint:gosec
	splash2Y := int16(p.TargetY) + dy + dy //nolint:gosec

	var splash1 *tile

	var splash2 *tile

	affedtedTile.setOnFire(bigBurn, bs)

	if splash1X >= 0 && splash1X < int16(boardMaxX) && splash1Y >= 0 && splash1Y < int16(boardMaxY) {
		splash1 = &bs.Board.Tiles[splash1X][splash1Y]
		splash1.setOnFire(midBurn, bs)
	}

	if splash2X >= 0 && splash2X < int16(boardMaxX) && splash2Y >= 0 && splash2Y < int16(boardMaxY) {
		splash2 = &bs.Board.Tiles[splash2X][splash2Y]
		splash2.setOnFire(minBurn, bs)
	}

	// 0. efekt to
	// a. obrażenia
	// unit/building.takedamage(p.Damage)
	// b. zapalenie kafelka na którym był cel
	// burningTileEffect()

	// 1. efekt100 w kafelku, damage
	// 2. efekt90 w kafelku+1, damage-10
	// 3. efekt80 w kafelku+2, damage-20
}

// @reminder: funkcja da efekt każdemu zaatakowanemu kafelkowi, ale bezpiecznik musi być
// w wydawaniu rozkazu, a nie tutaj. Tak aby unitMage nie mógł dostać rozkazu atakuj budynek
// reszta logiki w effects.go będzie działać tylko na jednostkach.
func (p *projectile) mageGhost(targetTile *tile, damage uint16, bs *battleState) {
	// @reminder: obrażenia zadajemy przed wejściem do tej funkcji, trzeba to zmienić!
	// KOSZT 20 many
	ownerBonus := uint16(0)
	if p.Owner == bs.AIPlayerID {
		ownerBonus += 20
	}

	totalDamage := damage + ownerBonus
	targetTile.ghost(p.Sprite, totalDamage, bs)
}
