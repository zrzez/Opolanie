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

	spriteMagicShield00 = 685
	spriteMagicShield01 = 686
	spriteMagicShield02 = 687
	spriteMagicShield03 = 688

	// 4. JEDNOSTKI (700+)
	// ID uint16 = 700 + (unitType * 200) + (Frame * 8) + Direction.
	spriteUnitBaseID uint16 = 700
	spriteUnitStep   uint16 = 200
)

const (
	// ======
	// Ręczne ustawianie duszków

	// Przycisk na UI do stworzenia jednostki!
	spriteManualArcherBtn = iota + 10_000
	// Bezczynność
	spriteManualArcherIdleUpLeft
	spriteManualArcherIdleLeft
	spriteManualArcherIdleDownLeft
	spriteManualArcherIdleDown
	spriteManualArcherIdleDownRight
	spriteManualArcherIdleRight
	spriteManualArcherIdleUpRight
	spriteManualArcherIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualArcherMove1Center
	spriteManualArcherMove1UpLeft
	spriteManualArcherMove1Left
	spriteManualArcherMove1DownLeft
	spriteManualArcherMove1Down
	spriteManualArcherMove1Up
	spriteManualArcherMove1UpRight
	spriteManualArcherMove1Right
	spriteManualArcherMove1DownRight

	// chodzenie 2
	spriteManualArcherMove2Center
	spriteManualArcherMove2UpLeft
	spriteManualArcherMove2Left
	spriteManualArcherMove2DownLeft
	spriteManualArcherMove2Down
	spriteManualArcherMove2Up
	spriteManualArcherMove2UpRight
	spriteManualArcherMove2Right
	spriteManualArcherMove2DownRight

	// Walka

	// walka 1
	spriteManualArcherAttack1UpLeft
	spriteManualArcherAttack1Left
	spriteManualArcherAttack1DownLeft
	spriteManualArcherAttack1Down
	spriteManualArcherAttack1Up
	spriteManualArcherAttack1UpRight
	spriteManualArcherAttack1Right
	spriteManualArcherAttack1DownRight

	// walka 2
	spriteManualArcherAttack2UpLeft
	spriteManualArcherAttack2Left
	spriteManualArcherAttack2DownLeft
	spriteManualArcherAttack2Down
	spriteManualArcherAttack2Up
	spriteManualArcherAttack2UpRight
	spriteManualArcherAttack2Right
	spriteManualArcherAttack2DownRight

	// KONIEC ŁUCZNIKA

	// Przycisk na UI do stworzenia jednostki!
	spriteManualCowBtn
	// Bezczynność
	spriteManualCowIdleUpLeft
	spriteManualCowIdleLeft
	spriteManualCowIdleDownLeft
	spriteManualCowIdleDown
	spriteManualCowIdleDownRight
	spriteManualCowIdleRight
	spriteManualCowIdleUpRight
	spriteManualCowIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualCowMove1Center
	spriteManualCowMove1UpLeft
	spriteManualCowMove1Left
	spriteManualCowMove1DownLeft
	spriteManualCowMove1Down
	spriteManualCowMove1Up
	spriteManualCowMove1UpRight
	spriteManualCowMove1Right
	spriteManualCowMove1DownRight

	// chodzenie 2
	spriteManualCowMove2Center
	spriteManualCowMove2UpLeft
	spriteManualCowMove2Left
	spriteManualCowMove2DownLeft
	spriteManualCowMove2Down
	spriteManualCowMove2Up
	spriteManualCowMove2UpRight
	spriteManualCowMove2Right
	spriteManualCowMove2DownRight

	// Walka

	// walka 1
	spriteManualCowAttack1UpLeft
	spriteManualCowAttack1Left
	spriteManualCowAttack1DownLeft
	spriteManualCowAttack1Down
	spriteManualCowAttack1Up
	spriteManualCowAttack1UpRight
	spriteManualCowAttack1Right
	spriteManualCowAttack1DownRight

	// walka 2
	spriteManualCowAttack2UpLeft
	spriteManualCowAttack2Left
	spriteManualCowAttack2DownLeft
	spriteManualCowAttack2Down
	spriteManualCowAttack2Up
	spriteManualCowAttack2UpRight
	spriteManualCowAttack2Right
	spriteManualCowAttack2DownRight

	// Koniec krowy.
	// Przycisk na UI do stworzenia jednostki!
	spriteManualPriestessBtn
	// Bezczynność
	spriteManualPriestessIdleUpLeft
	spriteManualPriestessIdleLeft
	spriteManualPriestessIdleDownLeft
	spriteManualPriestessIdleDown
	spriteManualPriestessIdleDownRight
	spriteManualPriestessIdleRight
	spriteManualPriestessIdleUpRight
	spriteManualPriestessIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualPriestessMove1Center
	spriteManualPriestessMove1UpLeft
	spriteManualPriestessMove1Left
	spriteManualPriestessMove1DownLeft
	spriteManualPriestessMove1Down
	spriteManualPriestessMove1Up
	spriteManualPriestessMove1UpRight
	spriteManualPriestessMove1Right
	spriteManualPriestessMove1DownRight

	// chodzenie 2
	spriteManualPriestessMove2Center
	spriteManualPriestessMove2UpLeft
	spriteManualPriestessMove2Left
	spriteManualPriestessMove2DownLeft
	spriteManualPriestessMove2Down
	spriteManualPriestessMove2Up
	spriteManualPriestessMove2UpRight
	spriteManualPriestessMove2Right
	spriteManualPriestessMove2DownRight

	// Walka

	// walka 1
	spriteManualPriestessAttack1UpLeft
	spriteManualPriestessAttack1Left
	spriteManualPriestessAttack1DownLeft
	spriteManualPriestessAttack1Down
	spriteManualPriestessAttack1Up
	spriteManualPriestessAttack1UpRight
	spriteManualPriestessAttack1Right
	spriteManualPriestessAttack1DownRight

	// walka 2
	spriteManualPriestessAttack2UpLeft
	spriteManualPriestessAttack2Left
	spriteManualPriestessAttack2DownLeft
	spriteManualPriestessAttack2Down
	spriteManualPriestessAttack2Up
	spriteManualPriestessAttack2UpRight
	spriteManualPriestessAttack2Right
	spriteManualPriestessAttack2DownRight

	// Koniec kapłanki
	// Przycisk na UI do stworzenia jednostki!
	spriteManualShepherdBtn
	// Bezczynność
	spriteManualShepherdIdleUpLeft
	spriteManualShepherdIdleLeft
	spriteManualShepherdIdleDownLeft
	spriteManualShepherdIdleDown
	spriteManualShepherdIdleDownRight
	spriteManualShepherdIdleRight
	spriteManualShepherdIdleUpRight
	spriteManualShepherdIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualShepherdMove1Center
	spriteManualShepherdMove1UpLeft
	spriteManualShepherdMove1Left
	spriteManualShepherdMove1DownLeft
	spriteManualShepherdMove1Down
	spriteManualShepherdMove1Up
	spriteManualShepherdMove1UpRight
	spriteManualShepherdMove1Right
	spriteManualShepherdMove1DownRight

	// chodzenie 2
	spriteManualShepherdMove2Center
	spriteManualShepherdMove2UpLeft
	spriteManualShepherdMove2Left
	spriteManualShepherdMove2DownLeft
	spriteManualShepherdMove2Down
	spriteManualShepherdMove2Up
	spriteManualShepherdMove2UpRight
	spriteManualShepherdMove2Right
	spriteManualShepherdMove2DownRight

	// Walka

	// walka 1
	spriteManualShepherdAttack1UpLeft
	spriteManualShepherdAttack1Left
	spriteManualShepherdAttack1DownLeft
	spriteManualShepherdAttack1Down
	spriteManualShepherdAttack1Up
	spriteManualShepherdAttack1UpRight
	spriteManualShepherdAttack1Right
	spriteManualShepherdAttack1DownRight

	// walka 2
	spriteManualShepherdAttack2UpLeft
	spriteManualShepherdAttack2Left
	spriteManualShepherdAttack2DownLeft
	spriteManualShepherdAttack2Down
	spriteManualShepherdAttack2Up
	spriteManualShepherdAttack2UpRight
	spriteManualShepherdAttack2Right
	spriteManualShepherdAttack2DownRight

	// Koniec pastucha
	// Przycisk na UI do stworzenia jednostki!
	spriteManualMageBtn
	// Bezczynność
	spriteManualMageIdleUpLeft
	spriteManualMageIdleLeft
	spriteManualMageIdleDownLeft
	spriteManualMageIdleDown
	spriteManualMageIdleDownRight
	spriteManualMageIdleRight
	spriteManualMageIdleUpRight
	spriteManualMageIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualMageMove1Center
	spriteManualMageMove1UpLeft
	spriteManualMageMove1Left
	spriteManualMageMove1DownLeft
	spriteManualMageMove1Down
	spriteManualMageMove1Up
	spriteManualMageMove1UpRight
	spriteManualMageMove1Right
	spriteManualMageMove1DownRight

	// chodzenie 2
	spriteManualMageMove2Center
	spriteManualMageMove2UpLeft
	spriteManualMageMove2Left
	spriteManualMageMove2DownLeft
	spriteManualMageMove2Down
	spriteManualMageMove2Up
	spriteManualMageMove2UpRight
	spriteManualMageMove2Right
	spriteManualMageMove2DownRight

	// Walka

	// walka 1
	spriteManualMageAttack1UpLeft
	spriteManualMageAttack1Left
	spriteManualMageAttack1DownLeft
	spriteManualMageAttack1Down
	spriteManualMageAttack1Up
	spriteManualMageAttack1UpRight
	spriteManualMageAttack1Right
	spriteManualMageAttack1DownRight

	// walka 2
	spriteManualMageAttack2UpLeft
	spriteManualMageAttack2Left
	spriteManualMageAttack2DownLeft
	spriteManualMageAttack2Down
	spriteManualMageAttack2Up
	spriteManualMageAttack2UpRight
	spriteManualMageAttack2Right
	spriteManualMageAttack2DownRight

	// Koniec pastucha
	// Przycisk na UI do stworzenia jednostki!
	spriteManualPriestBtn
	// Bezczynność
	spriteManualPriestIdleUpLeft
	spriteManualPriestIdleLeft
	spriteManualPriestIdleDownLeft
	spriteManualPriestIdleDown
	spriteManualPriestIdleDownRight
	spriteManualPriestIdleRight
	spriteManualPriestIdleUpRight
	spriteManualPriestIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualPriestMove1Center
	spriteManualPriestMove1UpLeft
	spriteManualPriestMove1Left
	spriteManualPriestMove1DownLeft
	spriteManualPriestMove1Down
	spriteManualPriestMove1Up
	spriteManualPriestMove1UpRight
	spriteManualPriestMove1Right
	spriteManualPriestMove1DownRight

	// chodzenie 2
	spriteManualPriestMove2Center
	spriteManualPriestMove2UpLeft
	spriteManualPriestMove2Left
	spriteManualPriestMove2DownLeft
	spriteManualPriestMove2Down
	spriteManualPriestMove2Up
	spriteManualPriestMove2UpRight
	spriteManualPriestMove2Right
	spriteManualPriestMove2DownRight

	// Walka

	// walka 1
	spriteManualPriestAttack1UpLeft
	spriteManualPriestAttack1Left
	spriteManualPriestAttack1DownLeft
	spriteManualPriestAttack1Down
	spriteManualPriestAttack1Up
	spriteManualPriestAttack1UpRight
	spriteManualPriestAttack1Right
	spriteManualPriestAttack1DownRight

	// walka 2
	spriteManualPriestAttack2UpLeft
	spriteManualPriestAttack2Left
	spriteManualPriestAttack2DownLeft
	spriteManualPriestAttack2Down
	spriteManualPriestAttack2Up
	spriteManualPriestAttack2UpRight
	spriteManualPriestAttack2Right
	spriteManualPriestAttack2DownRight

	// Koniec kapłana
	// Przycisk na UI do stworzenia jednostki!
	spriteManualSpearmanBtn
	// Bezczynność
	spriteManualSpearmanIdleUpLeft
	spriteManualSpearmanIdleLeft
	spriteManualSpearmanIdleDownLeft
	spriteManualSpearmanIdleDown
	spriteManualSpearmanIdleDownRight
	spriteManualSpearmanIdleRight
	spriteManualSpearmanIdleUpRight
	spriteManualSpearmanIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualSpearmanMove1Center
	spriteManualSpearmanMove1UpLeft
	spriteManualSpearmanMove1Left
	spriteManualSpearmanMove1DownLeft
	spriteManualSpearmanMove1Down
	spriteManualSpearmanMove1Up
	spriteManualSpearmanMove1UpRight
	spriteManualSpearmanMove1Right
	spriteManualSpearmanMove1DownRight

	// chodzenie 2
	spriteManualSpearmanMove2Center
	spriteManualSpearmanMove2UpLeft
	spriteManualSpearmanMove2Left
	spriteManualSpearmanMove2DownLeft
	spriteManualSpearmanMove2Down
	spriteManualSpearmanMove2Up
	spriteManualSpearmanMove2UpRight
	spriteManualSpearmanMove2Right
	spriteManualSpearmanMove2DownRight

	// Walka

	// walka 1
	spriteManualSpearmanAttack1UpLeft
	spriteManualSpearmanAttack1Left
	spriteManualSpearmanAttack1DownLeft
	spriteManualSpearmanAttack1Down
	spriteManualSpearmanAttack1Up
	spriteManualSpearmanAttack1UpRight
	spriteManualSpearmanAttack1Right
	spriteManualSpearmanAttack1DownRight

	// walka 2
	spriteManualSpearmanAttack2UpLeft
	spriteManualSpearmanAttack2Left
	spriteManualSpearmanAttack2DownLeft
	spriteManualSpearmanAttack2Down
	spriteManualSpearmanAttack2Up
	spriteManualSpearmanAttack2UpRight
	spriteManualSpearmanAttack2Right
	spriteManualSpearmanAttack2DownRight

	// Koniec włócznika

	// Przycisk na UI do stworzenia jednostki!
	spriteManualCrossbowmanBtn
	// Bezczynność
	spriteManualCrossbowmanIdleUpLeft
	spriteManualCrossbowmanIdleLeft
	spriteManualCrossbowmanIdleDownLeft
	spriteManualCrossbowmanIdleDown
	spriteManualCrossbowmanIdleDownRight
	spriteManualCrossbowmanIdleRight
	spriteManualCrossbowmanIdleUpRight
	spriteManualCrossbowmanIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualCrossbowmanMove1Center
	spriteManualCrossbowmanMove1UpLeft
	spriteManualCrossbowmanMove1Left
	spriteManualCrossbowmanMove1DownLeft
	spriteManualCrossbowmanMove1Down
	spriteManualCrossbowmanMove1Up
	spriteManualCrossbowmanMove1UpRight
	spriteManualCrossbowmanMove1Right
	spriteManualCrossbowmanMove1DownRight

	// chodzenie 2
	spriteManualCrossbowmanMove2Center
	spriteManualCrossbowmanMove2UpLeft
	spriteManualCrossbowmanMove2Left
	spriteManualCrossbowmanMove2DownLeft
	spriteManualCrossbowmanMove2Down
	spriteManualCrossbowmanMove2Up
	spriteManualCrossbowmanMove2UpRight
	spriteManualCrossbowmanMove2Right
	spriteManualCrossbowmanMove2DownRight

	// Walka

	// walka 1
	spriteManualCrossbowmanAttack1UpLeft
	spriteManualCrossbowmanAttack1Left
	spriteManualCrossbowmanAttack1DownLeft
	spriteManualCrossbowmanAttack1Down
	spriteManualCrossbowmanAttack1Up
	spriteManualCrossbowmanAttack1UpRight
	spriteManualCrossbowmanAttack1Right
	spriteManualCrossbowmanAttack1DownRight

	// walka 2
	spriteManualCrossbowmanAttack2UpLeft
	spriteManualCrossbowmanAttack2Left
	spriteManualCrossbowmanAttack2DownLeft
	spriteManualCrossbowmanAttack2Down
	spriteManualCrossbowmanAttack2Up
	spriteManualCrossbowmanAttack2UpRight
	spriteManualCrossbowmanAttack2Right
	spriteManualCrossbowmanAttack2DownRight

	// Koniec kusznika

	// Przycisk na UI do stworzenia jednostki!
	spriteManualAxemanBtn
	// Bezczynność
	spriteManualAxemanIdleUpLeft
	spriteManualAxemanIdleLeft
	spriteManualAxemanIdleDownLeft
	spriteManualAxemanIdleDown
	spriteManualAxemanIdleDownRight
	spriteManualAxemanIdleRight
	spriteManualAxemanIdleUpRight
	spriteManualAxemanIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualAxemanMove1Center
	spriteManualAxemanMove1UpLeft
	spriteManualAxemanMove1Left
	spriteManualAxemanMove1DownLeft
	spriteManualAxemanMove1Down
	spriteManualAxemanMove1Up
	spriteManualAxemanMove1UpRight
	spriteManualAxemanMove1Right
	spriteManualAxemanMove1DownRight

	// chodzenie 2
	spriteManualAxemanMove2Center
	spriteManualAxemanMove2UpLeft
	spriteManualAxemanMove2Left
	spriteManualAxemanMove2DownLeft
	spriteManualAxemanMove2Down
	spriteManualAxemanMove2Up
	spriteManualAxemanMove2UpRight
	spriteManualAxemanMove2Right
	spriteManualAxemanMove2DownRight

	// Walka

	// walka 1
	spriteManualAxemanAttack1UpLeft
	spriteManualAxemanAttack1Left
	spriteManualAxemanAttack1DownLeft
	spriteManualAxemanAttack1Down
	spriteManualAxemanAttack1Up
	spriteManualAxemanAttack1UpRight
	spriteManualAxemanAttack1Right
	spriteManualAxemanAttack1DownRight

	// walka 2
	spriteManualAxemanAttack2UpLeft
	spriteManualAxemanAttack2Left
	spriteManualAxemanAttack2DownLeft
	spriteManualAxemanAttack2Down
	spriteManualAxemanAttack2Up
	spriteManualAxemanAttack2UpRight
	spriteManualAxemanAttack2Right
	spriteManualAxemanAttack2DownRight

	// Koniec drwala
	// Przycisk na UI do stworzenia jednostki!
	spriteManualSwordsmanBtn
	// Bezczynność
	spriteManualSwordsmanIdleUpLeft
	spriteManualSwordsmanIdleLeft
	spriteManualSwordsmanIdleDownLeft
	spriteManualSwordsmanIdleDown
	spriteManualSwordsmanIdleDownRight
	spriteManualSwordsmanIdleRight
	spriteManualSwordsmanIdleUpRight
	spriteManualSwordsmanIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualSwordsmanMove1Center
	spriteManualSwordsmanMove1UpLeft
	spriteManualSwordsmanMove1Left
	spriteManualSwordsmanMove1DownLeft
	spriteManualSwordsmanMove1Down
	spriteManualSwordsmanMove1Up
	spriteManualSwordsmanMove1UpRight
	spriteManualSwordsmanMove1Right
	spriteManualSwordsmanMove1DownRight

	// chodzenie 2
	spriteManualSwordsmanMove2Center
	spriteManualSwordsmanMove2UpLeft
	spriteManualSwordsmanMove2Left
	spriteManualSwordsmanMove2DownLeft
	spriteManualSwordsmanMove2Down
	spriteManualSwordsmanMove2Up
	spriteManualSwordsmanMove2UpRight
	spriteManualSwordsmanMove2Right
	spriteManualSwordsmanMove2DownRight

	// Walka

	// walka 1
	spriteManualSwordsmanAttack1UpLeft
	spriteManualSwordsmanAttack1Left
	spriteManualSwordsmanAttack1DownLeft
	spriteManualSwordsmanAttack1Down
	spriteManualSwordsmanAttack1Up
	spriteManualSwordsmanAttack1UpRight
	spriteManualSwordsmanAttack1Right
	spriteManualSwordsmanAttack1DownRight

	// walka 2
	spriteManualSwordsmanAttack2UpLeft
	spriteManualSwordsmanAttack2Left
	spriteManualSwordsmanAttack2DownLeft
	spriteManualSwordsmanAttack2Down
	spriteManualSwordsmanAttack2Up
	spriteManualSwordsmanAttack2UpRight
	spriteManualSwordsmanAttack2Right
	spriteManualSwordsmanAttack2DownRight
	// koniec miecznika
	// Przycisk na UI do stworzenia jednostki!
	spriteManualCommanderBtn
	// Bezczynność
	spriteManualCommanderIdleUpLeft
	spriteManualCommanderIdleLeft
	spriteManualCommanderIdleDownLeft
	spriteManualCommanderIdleDown
	spriteManualCommanderIdleDownRight
	spriteManualCommanderIdleRight
	spriteManualCommanderIdleUpRight
	spriteManualCommanderIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualCommanderMove1Center
	spriteManualCommanderMove1UpLeft
	spriteManualCommanderMove1Left
	spriteManualCommanderMove1DownLeft
	spriteManualCommanderMove1Down
	spriteManualCommanderMove1Up
	spriteManualCommanderMove1UpRight
	spriteManualCommanderMove1Right
	spriteManualCommanderMove1DownRight

	// chodzenie 2
	spriteManualCommanderMove2Center
	spriteManualCommanderMove2UpLeft
	spriteManualCommanderMove2Left
	spriteManualCommanderMove2DownLeft
	spriteManualCommanderMove2Down
	spriteManualCommanderMove2Up
	spriteManualCommanderMove2UpRight
	spriteManualCommanderMove2Right
	spriteManualCommanderMove2DownRight

	// Walka

	// walka 1
	spriteManualCommanderAttack1UpLeft
	spriteManualCommanderAttack1Left
	spriteManualCommanderAttack1DownLeft
	spriteManualCommanderAttack1Down
	spriteManualCommanderAttack1Up
	spriteManualCommanderAttack1UpRight
	spriteManualCommanderAttack1Right
	spriteManualCommanderAttack1DownRight

	// walka 2
	spriteManualCommanderAttack2UpLeft
	spriteManualCommanderAttack2Left
	spriteManualCommanderAttack2DownLeft
	spriteManualCommanderAttack2Down
	spriteManualCommanderAttack2Up
	spriteManualCommanderAttack2UpRight
	spriteManualCommanderAttack2Right
	spriteManualCommanderAttack2DownRight
	// koniec dowódcy
	// Przycisk na UI do stworzenia jednostki!
	spriteManualBearBtn
	// Bezczynność
	spriteManualBearIdleUpLeft
	spriteManualBearIdleLeft
	spriteManualBearIdleDownLeft
	spriteManualBearIdleDown
	spriteManualBearIdleDownRight
	spriteManualBearIdleRight
	spriteManualBearIdleUpRight
	spriteManualBearIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualBearMove1Center
	spriteManualBearMove1UpLeft
	spriteManualBearMove1Left
	spriteManualBearMove1DownLeft
	spriteManualBearMove1Down
	spriteManualBearMove1Up
	spriteManualBearMove1UpRight
	spriteManualBearMove1Right
	spriteManualBearMove1DownRight

	// chodzenie 2
	spriteManualBearMove2Center
	spriteManualBearMove2UpLeft
	spriteManualBearMove2Left
	spriteManualBearMove2DownLeft
	spriteManualBearMove2Down
	spriteManualBearMove2Up
	spriteManualBearMove2UpRight
	spriteManualBearMove2Right
	spriteManualBearMove2DownRight

	// Walka

	// walka 1
	spriteManualBearAttack1UpLeft
	spriteManualBearAttack1Left
	spriteManualBearAttack1DownLeft
	spriteManualBearAttack1Down
	spriteManualBearAttack1Up
	spriteManualBearAttack1UpRight
	spriteManualBearAttack1Right
	spriteManualBearAttack1DownRight

	// walka 2
	spriteManualBearAttack2UpLeft
	spriteManualBearAttack2Left
	spriteManualBearAttack2DownLeft
	spriteManualBearAttack2Down
	spriteManualBearAttack2Up
	spriteManualBearAttack2UpRight
	spriteManualBearAttack2Right
	spriteManualBearAttack2DownRight
	// koneic niedźwiedzia
	// Przycisk na UI do stworzenia jednostki!
	spriteManualUnknownBtn
	// Bezczynność
	spriteManualUnknownIdleUpLeft
	spriteManualUnknownIdleLeft
	spriteManualUnknownIdleDownLeft
	spriteManualUnknownIdleDown
	spriteManualUnknownIdleDownRight
	spriteManualUnknownIdleRight
	spriteManualUnknownIdleUpRight
	spriteManualUnknownIdleUp

	// Chodzenie
	// chodzenie 1
	spriteManualUnknownMove1Center
	spriteManualUnknownMove1UpLeft
	spriteManualUnknownMove1Left
	spriteManualUnknownMove1DownLeft
	spriteManualUnknownMove1Down
	spriteManualUnknownMove1Up
	spriteManualUnknownMove1UpRight
	spriteManualUnknownMove1Right
	spriteManualUnknownMove1DownRight

	// chodzenie 2
	spriteManualUnknownMove2Center
	spriteManualUnknownMove2UpLeft
	spriteManualUnknownMove2Left
	spriteManualUnknownMove2DownLeft
	spriteManualUnknownMove2Down
	spriteManualUnknownMove2Up
	spriteManualUnknownMove2UpRight
	spriteManualUnknownMove2Right
	spriteManualUnknownMove2DownRight

	// Walka

	// walka 1
	spriteManualUnknownAttack1UpLeft
	spriteManualUnknownAttack1Left
	spriteManualUnknownAttack1DownLeft
	spriteManualUnknownAttack1Down
	spriteManualUnknownAttack1Up
	spriteManualUnknownAttack1UpRight
	spriteManualUnknownAttack1Right
	spriteManualUnknownAttack1DownRight

	// walka 2
	spriteManualUnknownAttack2UpLeft
	spriteManualUnknownAttack2Left
	spriteManualUnknownAttack2DownLeft
	spriteManualUnknownAttack2Down
	spriteManualUnknownAttack2Up
	spriteManualUnknownAttack2UpRight
	spriteManualUnknownAttack2Right
	spriteManualUnknownAttack2DownRight
	// koniec strzygi

	// Zwłoki
	spriteManualArcherCorpseFresh
	spriteManualArcherCorpseDecay
	spriteManualAxemanCorpseFresh
	spriteManualAxemanCorpseDecay
	spriteManualBearCorpseFresh
	spriteManualBearCorpseDecay
	spriteManualCommanderCorpseFresh
	spriteManualCommanderCorpseDecay
	spriteManualCowCorpseFresh
	spriteManualCowCorpseDecay
	spriteManualCrossbowmanCorpseFresh
	spriteManualCrossbowmanCorpseDecay
	spriteManualMageCorpseFresh
	spriteManualMageCorpseDecay
	spriteManualPriestCorpseFresh
	spriteManualPriestCorpseDecay
	spriteManualPriestessCorpseFresh
	spriteManualPriestessCorpseDecay
	spriteManualShepherdCorpseFresh
	spriteManualShepherdCorpseDecay
	spriteManualSpearmanCorpseFresh
	spriteManualSpearmanCorpseDecay
	spriteManualSwordsmanCorpseFresh
	spriteManualSwordsmanCorpseDecay
	spriteManualUnknownCorpseFresh
	spriteManualUnknownCorpseDecay
)
