package sim

import (
	vm "fireworx/vec"
	"math"
)

type EffectID = uint16

type Target struct {
	Effect EffectID
	Stage  uint8
}

func NewTarget(effect EffectID, stage uint8) Target {
	return Target{effect, stage}
}

type AxisKind uint8

const (
	AxisVelocity AxisKind = iota
	AxisWorld
)

type Axis struct {
	Kind  AxisKind
	World vm.Vec3
}
type Trail struct {
	Drag, Gravity, Inherit, Rate, Spread float32
	Life                                 Range
	Ramp                                 Ramp
}

type Effect struct {
	LiftSpeed Range
	Stages    []Stage
}

type Strobe struct {
	Hz, LitFrac float32
}

func (s *Strobe) Lit(age float32, seed uint64) bool {
	phase := hashPhase(seed, 0x5b)
	v := phase + age*s.Hz
	frac := v - float32(math.Floor(float64(v)))
	return frac < s.LitFrac
}

type Motion struct {
	Kind  MotionKind
	Accel float32
	Hz    float32
}

type MotionKind uint8

const (
	MotionBallistic MotionKind = iota
	MotionHelix
	MotionThrust
	MotionWander
)

func (m *Motion) Apply(vel vm.Vec3, age float32, seed uint64, dt float32) vm.Vec3 {
	switch m.Kind {
	case MotionBallistic:
		return vm.Zero

	case MotionHelix:
		basis := BasisFromAxis(vel)
		ph := (hashPhase(seed, 0x51) + age*m.Hz) * (math.Pi * 2)
		return basis.Azimuth(ph).Scale(m.Accel * dt)

	case MotionThrust:
		return vel.NormalizeOr(vm.Y).Scale(m.Accel * dt)

	case MotionWander:
		t := age * m.Hz
		k := uint64(t)
		frac := t - float32(k)
		frac = frac * frac * (3.0 - 2.0*frac)
		a := hashDir(seed, k)
		b := hashDir(seed, k+1)
		return a.Lerp(b, frac).NormalizeOr(a).Scale(m.Accel * dt)
	}
	return vm.Zero
}

func hashPhase(seed, k uint64) float32 {
	return float32(mixSeed(seed, k)>>40) * (1.0 / 16777216.0)
}

func hashDir(seed, k uint64) vm.Vec3 {
	rng := NewRand(mixSeed(seed, k))
	return rng.UnitSphere()
}

type Burst struct {
	Axis    Axis
	Child   Target
	Count   [2]uint16
	Inherit float32
	Offset  float32
	Pattern Pattern
	Speed   Range
}

func SphereBurst(count [2]uint16, speed Range, child Target) Burst {
	return Burst{
		Axis:    Axis{Kind: AxisVelocity},
		Child:   child,
		Count:   count,
		Inherit: 0.15,
		Offset:  0.2,
		Pattern: Pattern{Kind: PatternSphere},
		Speed:   speed,
	}
}

type Stage struct {
	Drag     float32
	Gravity  float32
	Life     Range
	Motion   Motion
	Ramp     Ramp
	Strobe   *Strobe
	Terminal []Burst
	Trail    *Trail
}

func LiftStage(terminal []Burst) Stage {
	return StarStage(LiftRamp).
		Burn(1.5, 1.7).
		WithDrag(0.35).
		WithTrail(LiftTrail).
		WithTerminal(terminal)
}

func MineStage(terminal []Burst) Stage {
	return StarStage(LiftRamp).Burn(0.04, 0.06).WithTerminal(terminal)
}

func StarStage(ramp Ramp) Stage {
	return Stage{
		Drag:     0,
		Gravity:  1.0,
		Life:     RangeAt(1.0),
		Motion:   Motion{Kind: MotionBallistic},
		Ramp:     ramp,
		Terminal: nil,
	}
}

func (s Stage) Burn(min, max float32) Stage {
	s.Life = NewRange(min, max)
	return s
}

func (s Stage) WithDrag(v float32) Stage {
	s.Drag = v
	return s
}

func (s Stage) Fuse(min, max float32) Stage {
	return s.Burn(min, max)
}

func (s Stage) WithGravity(v float32) Stage {
	s.Gravity = v
	return s
}

func (s Stage) WithMotion(m Motion) Stage {
	s.Motion = m
	return s
}

func (s Stage) WithStrobe(hz, lit float32) Stage {
	st := Strobe{Hz: hz, LitFrac: lit}
	s.Strobe = &st
	return s
}

func (s Stage) WithTerminal(b []Burst) Stage {
	s.Terminal = b
	return s
}

func (s Stage) WithTrail(t Trail) Stage {
	s.Trail = &t
	return s
}
