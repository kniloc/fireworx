package sim

import vm "fireworx/vec"

const (
	FxBrocade EffectID = iota
	FxChrysanthemum
	FxComet
	FxCrossette
	FxDragonsEggs
	FxFish
	FxPalm
	FxPeony
	FxPistil
	FxStrobe
	FxTourbillon
	FxWillow
	FxRing
)

var (
	orangeRamp = NewRamp(1, 8)
	redRamp    = NewRamp(2, 8)
	goldRamp   = NewRamp(3, 8)
	purpleRamp = NewRamp(4, 8)
	cyanRamp   = NewRamp(5, 8)
	greenRamp  = NewRamp(6, 8)
	silverRamp = NewRamp(7, 6)
	whiteRamp  = NewRamp(7, 8)

	LiftRamp  = goldRamp
	shellLift = NewRange(40.0, 44.0)

	LiftTrail = Trail{
		Drag: 2.0, Gravity: 0.3, Inherit: 0.1,
		Life: NewRange(0.15, 0.4), Ramp: LiftRamp, Rate: 120.0, Spread: 0.7,
	}
)

var Brocade = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			SphereBurst([2]uint16{85, 105}, NewRange(7.0, 9.0), NewTarget(FxBrocade, 1)),
		}),
		StarStage(goldRamp).
			Burn(3.4, 4.4).
			WithDrag(0.3).
			WithTrail(Trail{
				Drag: 1.0, Gravity: 0.28, Inherit: 0.04,
				Life: NewRange(1.3, 2.1), Ramp: goldRamp, Rate: 150.0, Spread: 0.2,
			}),
	},
}

var Chrysanthemum = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			SphereBurst([2]uint16{170, 210}, NewRange(9.5, 12.5), NewTarget(FxChrysanthemum, 1)),
		}),
		StarStage(purpleRamp).
			Burn(1.4, 1.9).
			WithDrag(1.1).
			WithTrail(Trail{
				Drag: 1.4, Gravity: 0.6, Inherit: 0.06,
				Life: NewRange(0.3, 0.55), Ramp: purpleRamp, Rate: 80.0, Spread: 0.25,
			}),
	},
}

var Comet = Effect{
	LiftSpeed: RangeAt(4.0),
	Stages: []Stage{
		MineStage([]Burst{
			{
				Axis:    Axis{Kind: AxisVelocity},
				Child:   NewTarget(FxComet, 1),
				Count:   [2]uint16{5, 8},
				Inherit: 0.25,
				Offset:  0.0,
				Pattern: Pattern{Kind: PatternCone, Angle: 0.30, Jitter: 0.02},
				Speed:   NewRange(20.0, 27.0),
			},
		}),
		StarStage(orangeRamp).
			Burn(1.5, 2.1).
			WithDrag(0.55).
			WithMotion(Motion{Kind: MotionThrust, Accel: 3.5}).
			WithTrail(Trail{
				Drag: 1.6, Gravity: 0.45, Inherit: 0.05,
				Life: NewRange(0.5, 0.9), Ramp: orangeRamp, Rate: 220.0, Spread: 0.5,
			}),
	},
}

var Crossette = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			SphereBurst([2]uint16{14, 18}, NewRange(8.0, 9.5), NewTarget(FxCrossette, 1)),
		}),
		StarStage(redRamp).
			Burn(0.7, 0.85).
			WithDrag(0.5).
			WithTrail(Trail{
				Drag: 1.5, Gravity: 0.6, Inherit: 0.08,
				Life: NewRange(0.25, 0.5), Ramp: redRamp, Rate: 90.0, Spread: 0.3,
			}).
			WithTerminal([]Burst{
				{
					Axis:    Axis{Kind: AxisVelocity},
					Child:   NewTarget(FxCrossette, 2),
					Count:   [2]uint16{4, 4},
					Inherit: 0.35,
					Offset:  0.0,
					Pattern: Pattern{Kind: PatternCrossette, Arms: 4, Forward: 0.15, Jitter: 0.05},
					Speed:   NewRange(4.5, 5.5),
				},
			}),
		StarStage(whiteRamp).
			Burn(0.5, 0.7).
			WithDrag(1.0).
			WithTerminal([]Burst{
				{
					Axis:    Axis{Kind: AxisVelocity},
					Child:   NewTarget(FxCrossette, 3),
					Count:   [2]uint16{6, 10},
					Inherit: 0.4,
					Offset:  0.0,
					Pattern: Pattern{Kind: PatternSphere},
					Speed:   NewRange(1.0, 2.2),
				},
			}),
		StarStage(whiteRamp).
			Burn(0.15, 0.3).
			WithDrag(3.0).
			WithGravity(0.8),
	},
}

var DragonsEggs = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			SphereBurst([2]uint16{110, 140}, NewRange(6.5, 9.5), NewTarget(FxDragonsEggs, 1)),
		}),
		StarStage(goldRamp).
			Burn(0.7, 2.0).
			WithDrag(1.5).
			WithTrail(Trail{
				Drag: 2.0, Gravity: 0.7, Inherit: 0.05,
				Life: NewRange(0.15, 0.3), Ramp: goldRamp, Rate: 30.0, Spread: 0.15,
			}).
			WithTerminal([]Burst{
				{
					Axis:    Axis{Kind: AxisVelocity},
					Child:   NewTarget(FxDragonsEggs, 2),
					Count:   [2]uint16{10, 16},
					Inherit: 0.3,
					Offset:  0.0,
					Pattern: Pattern{Kind: PatternSphere},
					Speed:   NewRange(1.4, 3.2),
				},
			}),
		StarStage(whiteRamp).
			Burn(0.06, 0.16).
			WithDrag(4.0).
			WithGravity(0.9),
	},
}

var Fish = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			SphereBurst([2]uint16{45, 65}, NewRange(3.0, 5.0), NewTarget(FxFish, 1)),
		}),
		StarStage(purpleRamp).
			Burn(1.1, 1.7).
			WithDrag(0.9).
			WithGravity(0.5).
			WithMotion(Motion{Kind: MotionWander, Accel: 34.0, Hz: 7.0}).
			WithTrail(Trail{
				Drag: 2.5, Gravity: 0.5, Inherit: 0.05,
				Life: NewRange(0.12, 0.28), Ramp: purpleRamp, Rate: 110.0, Spread: 0.2,
			}),
	},
}

var Palm = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			{
				Axis:    Axis{Kind: AxisVelocity},
				Child:   NewTarget(FxPalm, 1),
				Count:   [2]uint16{56, 56},
				Inherit: 0.2,
				Offset:  0.3,
				Pattern: Pattern{Kind: PatternSpokes, Cone: 1.05, Spokes: 8, Spread: 0.055},
				Speed:   NewRange(9.0, 13.0),
			},
		}),
		StarStage(goldRamp).
			Burn(1.7, 2.2).
			WithDrag(0.55).
			WithMotion(Motion{Kind: MotionThrust, Accel: 2.0}).
			WithTrail(Trail{
				Drag: 1.3, Gravity: 0.4, Inherit: 0.05,
				Life: NewRange(0.55, 1.0), Ramp: goldRamp, Rate: 190.0, Spread: 0.3,
			}),
	},
}

var Peony = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			SphereBurst([2]uint16{150, 190}, NewRange(9.0, 12.5), NewTarget(FxPeony, 1)),
		}),
		StarStage(greenRamp).Burn(1.2, 1.7).WithDrag(1.4),
	},
}

var Pistil = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			SphereBurst([2]uint16{150, 180}, NewRange(10.0, 13.0), NewTarget(FxPistil, 1)),
			SphereBurst([2]uint16{55, 70}, NewRange(3.0, 4.6), NewTarget(FxPistil, 2)),
		}),
		StarStage(silverRamp).Burn(1.3, 1.8).WithDrag(1.3),
		StarStage(redRamp).Burn(1.5, 2.0).WithDrag(1.3),
	},
}

var Ring = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			{
				Axis:    Axis{Kind: AxisWorld, World: vm.Z},
				Child:   NewTarget(FxRing, 1),
				Count:   [2]uint16{64, 64},
				Inherit: 0.1,
				Offset:  0.1,
				Pattern: Pattern{Kind: PatternRing, Jitter: 0.03},
				Speed:   NewRange(10.0, 10.4),
			},
			SphereBurst([2]uint16{30, 40}, NewRange(1.5, 3.0), NewTarget(FxRing, 1)),
		}),
		StarStage(cyanRamp).Burn(1.4, 1.8).WithDrag(1.3),
	},
}

var StrobeFx = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			SphereBurst([2]uint16{90, 120}, NewRange(5.0, 8.0), NewTarget(FxStrobe, 1)),
		}),
		StarStage(whiteRamp).
			Burn(2.6, 3.6).
			WithDrag(2.0).
			WithGravity(0.55).
			WithStrobe(11.0, 0.3),
	},
}

var Tourbillon = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			SphereBurst([2]uint16{9, 14}, NewRange(7.0, 9.0), NewTarget(FxTourbillon, 1)),
		}),
		StarStage(silverRamp).
			Burn(1.3, 1.8).
			WithDrag(0.7).
			WithMotion(Motion{Kind: MotionHelix, Accel: 120.0, Hz: 5.5}).
			WithTrail(Trail{
				Drag: 1.5, Gravity: 0.5, Inherit: 0.03,
				Life: NewRange(0.35, 0.7), Ramp: silverRamp, Rate: 260.0, Spread: 0.1,
			}),
	},
}

var Willow = Effect{
	LiftSpeed: shellLift,
	Stages: []Stage{
		LiftStage([]Burst{
			SphereBurst([2]uint16{70, 90}, NewRange(6.0, 8.0), NewTarget(FxWillow, 1)),
		}).Fuse(1.7, 1.9),
		StarStage(goldRamp).
			Burn(2.6, 3.4).
			WithDrag(0.45).
			WithTrail(Trail{
				Drag: 1.2, Gravity: 0.55, Inherit: 0.05,
				Life: NewRange(0.8, 1.4), Ramp: goldRamp, Rate: 60.0, Spread: 0.25,
			}),
	},
}

var Catalog = []Effect{
	Brocade, Chrysanthemum, Comet, Crossette, DragonsEggs, Fish,
	Palm, Peony, Pistil, StrobeFx, Tourbillon, Willow, Ring,
}
