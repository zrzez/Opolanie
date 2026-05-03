package main

// assets_ids.go

const (
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

	spriteRoadStart  uint16 = 31
	spriteRoadButton uint16 = 46
	spriteRoadEnd    uint16 = 51

	/*
		R - prawo
		L - lewo
		U - góra
		D - dół
	*/
	spriteRoadR    uint16 = 37
	spriteRoadL    uint16 = 38
	spriteRoadD    uint16 = 39
	spriteRoadU    uint16 = 40
	spriteRoadRU   uint16 = 41
	spriteRoadUD   uint16 = 42
	spriteRoadRD   uint16 = 43
	spriteRoadRUD  uint16 = 44
	spriteRoadLU   uint16 = 45
	spriteRoadLR   uint16 = 46
	spriteRoadLRU  uint16 = 47
	spriteRoadLD   uint16 = 48
	spriteRoadLUD  uint16 = 49
	spriteRoadLRD  uint16 = 50
	spriteRoadLRUD uint16 = 51

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

	spriteBridgeConstruction uint16 = 516

	spriteDryEarth01 uint16 = 60
	spriteDryEarth02 uint16 = 61
	spriteDryEarth03 uint16 = 62

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

	spriteWaterStart  uint16 = 74
	spriteWaterMiddle uint16 = 82
	spriteWaterEnd    uint16 = 112

	spriteTreeStumpStart uint16 = 113
	spriteTreeStumpEnd   uint16 = 119
	spriteTreeTopStart   uint16 = 120
	spriteTreeTopEnd     uint16 = 126

	spriteTreeStump00    uint16 = 113
	spriteTreeStump01    uint16 = 114
	spriteTreeStump02    uint16 = 115
	spriteTreeStump03    uint16 = 116
	spriteTreeStump04    uint16 = 117
	spriteTreeStump05    uint16 = 118
	spriteDryTreeStump00 uint16 = 119

	// Nie są używane bezpośrednio - 03.05.2026.
	// Pośrednio przy rysowaniu pni dodajemy przesunięcie
	// aby uzyskać właściwy numer duszka.
	spriteTreeTop00    uint16 = 120
	spriteTreeTop01    uint16 = 121
	spriteTreeTop02    uint16 = 122
	spriteTreeTop03    uint16 = 123
	spriteTreeTop04    uint16 = 124
	spriteTreeTop05    uint16 = 125
	spriteDryTreeTop00 uint16 = 126

	spriteTreeBurntStump00 uint16 = 380
	spriteTreeBurntStump01 uint16 = 381

	spriteTreeBurntTop00 uint16 = 382
	spriteTreeBurntTop01 uint16 = 383

	spriteDryFallenTreeStump uint16 = 384
	spriteDryFallenTreeTop   uint16 = 385

	spriteDryFallingStump   uint16 = 386
	spriteDryFallingTreeTop uint16 = 387

	spriteDryLeaningTreeStump     uint16 = 388
	spriteDryLeaningTreeTop       uint16 = 389
	spriteDryLeaningTreeCrownLeft uint16 = 390

	spriteBurntFallenTreeStump uint16 = 391
	spriteBurntFallenTreeCrown uint16 = 392

	spriteBurntFallingTreeStump uint16 = 393
	spriteBurntFallingTreeCrown uint16 = 394

	spriteBurntLeaningTreeStump     uint16 = 395
	spriteBurntLeaningTreeCrown     uint16 = 396
	spriteBurntLeaningTreeCrownLeft uint16 = 397

	spriteConstructionStart uint16 = 127
	spriteConstructionEnd   uint16 = 135

	spriteRuinStart uint16 = 257
	spriteRuinEnd   uint16 = 265

	spritePalisadeStart uint16 = 266
	spritePalisadeEnd   uint16 = 277

	spritePalisadeNE        uint16 = 266
	spritePalisadeV         uint16 = 267
	spritePalisadeH         uint16 = 271
	spritePalisadeDestroyed uint16 = 277

	spriteBuildingMainBase      uint16 = 137
	spriteBuildingBarnBase      uint16 = 157
	spriteBuildingBarracksBase  uint16 = 177
	spriteBuildingTempleBase    uint16 = 197
	spriteBuildingBarracks2Base uint16 = 217
	spriteBuildingAcademyBase   uint16 = 237

	spriteBuildingEnd uint16 = 255

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

	spriteEffectskeleton00 uint16 = 612
	spriteEffectskeleton01 uint16 = 613
	spriteEffectskeleton02 uint16 = 614

	spriteEffectHit00 uint16 = 615
	spriteEffectHit01 uint16 = 616

	spriteEffectHeal00 uint16 = 617
	spriteEffectHeal01 uint16 = 618

	// @todo: oba powinny być wykorzystywane. Jakieś niedociągnięcie przy rysowaniu!
	spriteEffectTransform00 uint16 = 619
	spriteEffectTransform01 uint16 = 620

	spriteVictoryPoint uint16 = 621

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

	spriteAsh00 uint16 = 634
	spriteAsh01 uint16 = 635

	spriteMissileArrowUp        = 636
	spriteMissileArrowUpLeft    = 637
	spriteMissileArrowLeft      = 638
	spriteMissileArrowDownLeft  = 639
	spriteMissileArrowDown      = 640
	spriteMissileArrowUpRight   = 641
	spriteMissileArrowRight     = 642
	spriteMissileArrowDownRight = 643

	spriteMissileBoltUp        = 644
	spriteMissileBoltUpLeft    = 645
	spriteMissileBoltLeft      = 646
	spriteMissileBoltDownLeft  = 647
	spriteMissileBoltDown      = 648
	spriteMissileBoltUpRight   = 649
	spriteMissileBoltRight     = 650
	spriteMissileBoltDownRight = 651

	spriteMissileFireUp        = 652
	spriteMissileFireUpLeft    = 653
	spriteMissileFireLeft      = 654
	spriteMissileFireDownLeft  = 655
	spriteMissileFireDown      = 656
	spriteMissileFireUpRight   = 657
	spriteMissileFireRight     = 658
	spriteMissileFireDownRight = 659

	spriteMissileLightningUp        = 660
	spriteMissileLightningUpLeft    = 661
	spriteMissileLightningLeft      = 662
	spriteMissileLightningDownLeft  = 663
	spriteMissileLightningDown      = 664
	spriteMissileLightningUpRight   = 665
	spriteMissileLightningRight     = 666
	spriteMissileLightningDownRight = 667

	spriteMissileSpearUp        = 668
	spriteMissileSpearUpLeft    = 669
	spriteMissileSpearLeft      = 670
	spriteMissileSpearDownLeft  = 671
	spriteMissileSpearDown      = 672
	spriteMissileSpearUpRight   = 673
	spriteMissileSpearRight     = 674
	spriteMissileSpearDownRight = 675

	spriteMissileGhostUp        = 676
	spriteMissileGhostUpLeft    = 677
	spriteMissileGhostLeft      = 678
	spriteMissileGhostDownLeft  = 679
	spriteMissileGhostDown      = 680
	spriteMissileGhostUpRight   = 681
	spriteMissileGhostRight     = 682
	spriteMissileGhostDownRight = 683
	spriteMissileGhostAttack    = 684

	// 4. JEDNOSTKI (700+)
	// ID uint16 = 700 + (unitType * 200) + (Frame * 8) + Direction.
	spriteUnitBaseID uint16 = 700
	spriteUnitStep   uint16 = 200
)
