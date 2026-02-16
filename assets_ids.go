package main

// assets_ids.go

// Zestawienie tekstura - ID
const (
	// Trawa 0 - 17

	spriteGrassStart uint16 = 2
	spriteGrassEnd   uint16 = 17

	spriteGrass00 uint16 = 2
	spriteGrass01 uint16 = 3
	spriteGrass02 uint16 = 4
	spriteGrass03 uint16 = 5
	spriteGrass04 uint16 = 6
	spriteGrass05 uint16 = 7
	spriteGrass06 uint16 = 8
	spriteGrass07 uint16 = 9
	spriteGrass08 uint16 = 10
	spriteGrass09 uint16 = 11
	spriteGrass10 uint16 = 12
	spriteGrass11 uint16 = 13
	spriteGrass12 uint16 = 14
	spriteGrass13 uint16 = 15
	spriteGrass14 uint16 = 16
	spriteGrass15 uint16 = 17

	spriteGrassStubbed uint16 = 0
	spriteGrassGrazed  uint16 = 1

	// --- SKAŁY (18 - 30) ---
	// były 9-21
	spriteRockStart uint16 = 18
	spriteRockEnd   uint16 = 30

	spriteRock00 uint16 = 18
	spriteRock01 uint16 = 19
	spriteRock02 uint16 = 20
	spriteRock03 uint16 = 21
	spriteRock04 uint16 = 22
	spriteRock05 uint16 = 23
	spriteRock06 uint16 = 24
	spriteRock07 uint16 = 25
	spriteRock08 uint16 = 26
	spriteRock09 uint16 = 27
	spriteRock10 uint16 = 28
	spriteRock11 uint16 = 29
	spriteRock12 uint16 = 30 // przechodnie! dawne rocks_17.png

	// --- DROGI (31 - 51) ---
	// były 25-45
	spriteRoadStart uint16 = 31
	// nie wiem, czy trzeba dodawać kolejne
	spriteRoadEnd uint16 = 51

	// --- MOSTY (55 - 62) ---
	// były 46-53
	spriteBridgeStart uint16 = 52
	spriteBridgeEnd   uint16 = 59

	spriteBridge00 uint16 = 52
	spriteBridge01 uint16 = 53
	spriteBridge02 uint16 = 54
	spriteBridge03 uint16 = 55
	spriteBridge04 uint16 = 56
	spriteBridge05 uint16 = 57
	spriteBridge06 uint16 = 58
	spriteBridge07 uint16 = 59

	// --- SUCHA ZIEMIA () ---
	// były 22-24
	spriteDryEarth01 uint16 = 60
	spriteDryEarth02 uint16 = 61
	spriteDryEarth03 uint16 = 62

	// --- GADŻETY
	spriteGadgetStart uint16 = 363
	spriteGadgetEnd   uint16 = 378

	spriteGadget00 uint16 = 363
	spriteGadget01 uint16 = 364
	spriteGadget02 uint16 = 365
	spriteGadget03 uint16 = 366
	spriteGadget04 uint16 = 367
	spriteGadget05 uint16 = 368
	spriteGadget06 uint16 = 369
	spriteGadget07 uint16 = 370
	spriteGadget08 uint16 = 371
	spriteGadget09 uint16 = 372
	spriteGadget10 uint16 = 373
	spriteGadget11 uint16 = 374
	spriteGadget12 uint16 = 375
	spriteGadget13 uint16 = 376
	spriteGadget14 uint16 = 377
	spriteGadget15 uint16 = 378

	// --- WODA ---
	spriteWaterStart  uint16 = 74 // Początek kafelków wody
	spriteWaterMiddle uint16 = 82
	spriteWaterEnd    uint16 = 112

	// --- DRZEWA ---
	spriteTreeStumpStart uint16 = 113 // Pnie
	spriteTreeStumpEnd   uint16 = 119
	spriteTreeTopStart   uint16 = 120 // Korony
	spriteTreeTopEnd     uint16 = 126

	spriteTreeStump00 uint16 = 113
	spriteTreeStump01 uint16 = 114
	spriteTreeStump02 uint16 = 115
	spriteTreeStump03 uint16 = 116
	spriteTreeStump04 uint16 = 117
	spriteTreeStump05 uint16 = 118
	spriteTreeStump06 uint16 = 119 // Suche

	spriteTreeTop00 uint16 = 120
	spriteTreeTop01 uint16 = 121
	spriteTreeTop02 uint16 = 122
	spriteTreeTop03 uint16 = 123
	spriteTreeTop04 uint16 = 124
	spriteTreeTop05 uint16 = 125
	spriteTreeTop06 uint16 = 126 // Suche

	spriteTreeBurntStump00 uint16 = 380 // Spalone stojące
	spriteTreeBurntStump01 uint16 = 381 // Spalone stojące

	spriteTreeBurntTop00 uint16 = 382 // Spalone stojące
	spriteTreeBurntTop01 uint16 = 383 // Spalone stojące

	spriteTreeFallingDry00 uint16 = 384
	spriteTreeFallingDry01 uint16 = 385
	spriteTreeFallingDry02 uint16 = 386

	spriteTreeFallingBurnt00 uint16 = 387
	spriteTreeFallingBurnt01 uint16 = 388
	spriteTreeFallingBurnt02 uint16 = 389

	spriteTreeFallenDry00 uint16 = 390
	spriteTreeFallenDry01 uint16 = 391

	spriteTreeFallenBurnt00 uint16 = 392
	spriteTreeFallenBurnt01 uint16 = 393

	// 2. BUDYNKI I KONSTRUKCJE

	// --- PLAC BUDOWY (127 - 135) ---
	spriteConstructionStart uint16 = 127
	spriteConstructionEnd   uint16 = 135

	// --- RUINY (257 - 265) ---
	spriteRuinStart uint16 = 257
	spriteRuinEnd   uint16 = 265

	// --- PALISADY (266 - 277) ---
	spritePalisadeStart uint16 = 266
	spritePalisadeEnd   uint16 = 277

	spritePalisadeNE        uint16 = 266 // Niepołączona (pojedynczy słup)
	spritePalisadeV         uint16 = 267 // Pionowa
	spritePalisadeH         uint16 = 271 // Pozioma
	spritePalisadeDestroyed uint16 = 277 // Zniszczona

	// --- GOTOWE BUDYNKI --
	spriteBuildingMainBase      uint16 = 137
	spriteBuildingBarnBase      uint16 = 157
	spriteBuildingBarracksBase  uint16 = 177
	spriteBuildingTempleBase    uint16 = 197
	spriteBuildingBarracks2Base uint16 = 217
	spriteBuildingAcademyBase   uint16 = 237

	spriteBuildingEnd uint16 = 255
	// 3. UI I EFEKTY

	// --- KURSORY ---
	spriteCursorDefaultBig    uint16 = 501
	spriteCursorCrossWhite    uint16 = 502
	spriteCursorCrossRed      uint16 = 503
	spriteCursorSmallWhite    uint16 = 504
	spriteCursorFrameRed      uint16 = 505
	spriteCursorCrossMedRed   uint16 = 506
	spriteCursorCrossMedWhite uint16 = 507
	spriteCursorArrowUp       uint16 = 508
	spriteCursorArrowDown     uint16 = 509
	spriteCursorArrowLeft     uint16 = 510
	spriteCursorArrowRight    uint16 = 511
	spriteCursorStop          uint16 = 512
	spriteCursorDefaultSmall  uint16 = 514
	spriteCursorPointer       uint16 = 515
	spriteCursorFrameWhite    uint16 = 516

	// --- PRZYCISKI (600+) ---
	spriteBtnBuildPalisade      uint16 = 600
	spriteBtnShield             uint16 = 601
	spriteBtnRepair             uint16 = 603
	spriteBtnBuildBarn          uint16 = 604
	spriteBtnBuildBarracks      uint16 = 605
	spriteBtnBuildTemple        uint16 = 606
	spriteBtnBuildBarracks2     uint16 = 607
	spriteBtnBuildAcademy       uint16 = 608
	spriteBtnSpellVision        uint16 = 513
	spriteBtnSpellMagicShield   uint16 = 609
	spriteBtnSpellMagicLighting uint16 = 610
	spriteBtnSpellMagicFire     uint16 = 611

	// --- EFEKTY ---
	spriteEffectSkeleton_00 uint16 = 612
	spriteEffectSkeleton_01 uint16 = 613
	spriteEffectSkeleton_02 uint16 = 614
	spriteEffectHit00       uint16 = 615
	spriteEffectHit01       uint16 = 616

	// --- LECZENIE
	spriteEffectHeal00 uint16 = 617
	spriteEffectHeal01 uint16 = 618

	spriteEffectTransform00 uint16 = 619
	spriteEffectTransform01 uint16 = 620

	spriteVictoryPoint uint16 = 621

	// --- OGIEŃ
	spriteFireStart uint16 = 622
	spriteFireEnd   uint16 = 635

	spriteFire00 uint16 = 622
	spriteFire01 uint16 = 623
	spriteFire02 uint16 = 624
	spriteFire03 uint16 = 625
	spriteFire04 uint16 = 626
	spriteFire05 uint16 = 627
	spriteFire06 uint16 = 628
	spriteFire07 uint16 = 629
	spriteFire08 uint16 = 630
	spriteFire09 uint16 = 631
	spriteFire10 uint16 = 632
	spriteFire11 uint16 = 633
	spriteFire12 uint16 = 634
	spriteFire13 uint16 = 635

	// 4. JEDNOSTKI (700+)
	// ID uint16 = 700 + (Type * 200) + (Frame * 8) + Dir.
	spriteUnitBaseID uint16 = 700
	spriteUnitStep   uint16 = 200
)
