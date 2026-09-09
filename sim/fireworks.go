package sim

import (
	vm "fireworx/vec"
)

type Fireworks struct {
	frame   uint64
	lib     []Effect
	pending []Star
	rng     *Rand
	sparks  []Spark
	stars   []Star
	World   World
}

func NewFireWorks(lib []Effect, w World, seed uint64) *Fireworks {
	return &Fireworks{
		lib:     lib,
		pending: make([]Star, 0, 512),
		rng:     NewRand(seed),
		sparks:  make([]Spark, 0, w.SparkCap),
		stars:   make([]Star, 0, w.StarCap),
		World:   w,
	}
}

func (f *Fireworks) stage(effect EffectID, stage uint8) *Stage {
	return &f.lib[effect].Stages[stage]
}

func (f *Fireworks) LaunchFromGround(effect EffectID, x float32) {
	up := f.lib[effect].LiftSpeed.Sample(f.rng)
	f.Launch(effect, vm.New(x, 0, 0), vm.New(0, up, 0))
}

func (f *Fireworks) Launch(effect EffectID, pos, vel vm.Vec3) {
	stage := f.stage(effect, 0)
	life := stage.Life.Sample(f.rng)
	seed := f.rng.Uint64()
	f.pushStar(Star{
		Effect: effect,
		Life:   life,
		Pos:    pos,
		Seed:   seed,
		Vel:    vel,
	})
}

func (f *Fireworks) Update(dt float32) {
	f.frame++
	f.updateSparks(dt)
	f.updateStars(dt)
	f.stars = append(f.stars, f.pending...)
	f.pending = f.pending[:0]
}

func (f *Fireworks) Visit(fn func(pos vm.Vec3, index uint8, isStar bool)) {
	for _, s := range f.sparks {
		fn(s.Pos, s.Ramp.Sample(s.T()), false)
	}

	for _, s := range f.stars {
		stage := f.stage(s.Effect, s.Stage)
		if stage.Strobe != nil && !stage.Strobe.Lit(s.Age, s.Seed) {
			continue
		}
		fn(s.Pos, stage.Ramp.Sample(s.T()), true)
	}
}

func dragAndAdvance(pos, vel *vm.Vec3, drag, dt float32) {
	*vel = vel.Scale(1.0 / (1.0 + drag*dt))
	*pos = pos.Add(vel.Scale(dt))
}

func (f *Fireworks) emitTrail(trail *Trail, prev vm.Vec3, s *Star, dt float32) {
	s.TrailAcc += trail.Rate * dt
	n := uint32(s.TrailAcc)
	if n == 0 {
		return
	}

	s.TrailAcc -= float32(n)
	rng := NewRand(mixSeed(s.Seed, f.frame))
	step := s.Pos.Sub(prev)
	for k := range n {
		if len(f.sparks) >= f.World.SparkCap {
			break
		}

		frac := (float32(k) + 0.5) / float32(n)
		life := trail.Life.Sample(rng)
		f.sparks = append(f.sparks, Spark{
			Drag: trail.Drag, Gravity: trail.Gravity, Life: life,
			Pos:  prev.Add(step.Scale(frac)),
			Ramp: trail.Ramp,
			Vel:  s.Vel.Scale(trail.Inherit).Add(rng.UnitSphere().Scale(trail.Spread)),
		})
	}
}

func (f *Fireworks) fireBurst(b Burst, parent *Star) {
	rng := NewRand(mixSeed(parent.Seed, 0x42))
	var axis vm.Vec3

	if b.Axis.Kind == AxisVelocity {
		axis = parent.Vel
	} else {
		axis = b.Axis.World
	}

	basis := BasisFromAxis(axis)
	count := rng.IntRange(b.Count[0], b.Count[1])
	childStage := f.stage(b.Child.Effect, b.Child.Stage)

	for index := range count {
		ctx := EmitCtx{Count: count, Index: index}
		dir := b.Pattern.Direction(rng, ctx, basis)
		speed := b.Speed.Sample(rng)
		life := childStage.Life.Sample(rng)
		seed := rng.Uint64()
		f.pushStar(Star{
			Effect: b.Child.Effect, Life: life,
			Pos:      parent.Pos.Add(dir.Scale(b.Offset)),
			Seed:     seed,
			Stage:    b.Child.Stage,
			TrailAcc: rng.Float32(),
			Vel:      parent.Vel.Scale(b.Inherit).Add(dir.Scale(speed)),
		})
	}
}

func (f *Fireworks) pushStar(s Star) {
	if len(f.stars)+len(f.pending) < f.World.StarCap {
		f.pending = append(f.pending, s)
	}
}

func (f *Fireworks) updateSparks(dt float32) {
	g, wind := f.World.Gravity, f.World.Wind
	out := f.sparks[:0]

	for i := range f.sparks {
		s := &f.sparks[i]
		s.Age += dt
		if s.Age >= s.Life {
			continue
		}
		s.Vel = s.Vel.Add(g.Scale(s.Gravity).Add(wind).Scale(dt))
		dragAndAdvance(&s.Pos, &s.Vel, s.Drag, dt)
		out = append(out, *s)
	}

	f.sparks = out
}

func (f *Fireworks) updateStars(dt float32) {
	g, wind := f.World.Gravity, f.World.Wind
	i := 0

	for i < len(f.stars) {
		s := f.stars[i]
		stage := f.stage(s.Effect, s.Stage)
		prev := s.Pos
		s.Age += dt
		s.Vel = s.Vel.Add(g.Scale(stage.Gravity).Add(wind).Scale(dt))
		s.Vel = s.Vel.Add(stage.Motion.Apply(s.Vel, s.Age, s.Seed, dt))
		dragAndAdvance(&s.Pos, &s.Vel, stage.Drag, dt)

		if stage.Trail != nil {
			f.emitTrail(stage.Trail, prev, &s, dt)
		}

		if s.Age >= s.Life {
			for _, burst := range stage.Terminal {
				f.fireBurst(burst, &s)
			}
			// swap-remove
			f.stars[i] = f.stars[len(f.stars)-1]
			f.stars = f.stars[:len(f.stars)-1]
		} else {
			f.stars[i] = s
			i++
		}
	}
}
