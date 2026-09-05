package sim

import (
	vm "fireworx/vec"
	"math"
)

type PatternKind uint8

const (
	PatternCone PatternKind = iota
	PatternCrossette
	PatternRing
	PatternSphere
	PatternSpokes
)

type Pattern struct {
	Kind PatternKind

	// Cone
	Angle, Jitter float32

	// Crossette
	Arms    uint16
	Forward float32

	// Spokes
	Cone, Spread float32
	Spokes       uint16
}

type EmitCtx struct {
	Count, Index uint16
}

type Basis struct {
	U, V, W vm.Vec3
}

func cos(x float32) float32 { return float32(math.Cos(float64(x))) }
func sin(x float32) float32 { return float32(math.Sin(float64(x))) }

func jitter(dir vm.Vec3, amount float32, rng *Rand) vm.Vec3 {
	if amount <= 0 {
		return dir
	}
	return dir.Add(rng.UnitSphere().Scale(amount)).NormalizeOr(dir)
}

func BasisFromAxis(axis vm.Vec3) Basis {
	w := axis.NormalizeOr(vm.Y)
	u, v := w.AnyOrthonormalPair()
	return Basis{u, v, w}
}

func (b *Basis) Apply(local vm.Vec3) vm.Vec3 {
	return b.U.Scale(local.X).Add(b.V.Scale(local.Y)).Add(b.W.Scale(local.Z))
}

func (b *Basis) Azimuth(angle float32) vm.Vec3 {
	return b.U.Scale(cos(angle)).Add(b.V.Scale(sin(angle)))
}

func (p Pattern) Direction(rng *Rand, ctx EmitCtx, basis Basis) vm.Vec3 {
	switch p.Kind {
	case PatternCone:
		d := basis.Apply(rng.UnitCap(cos(p.Angle)))
		return jitter(d, p.Jitter, rng)

	case PatternCrossette:
		arms := p.Arms
		if arms == 0 {
			arms = 1
		}

		a := float32(ctx.Index) / float32(arms) * (math.Pi * 2)
		flat := basis.Azimuth(a)
		d := flat.Lerp(basis.W, p.Forward).NormalizeOr(flat)
		return jitter(d, p.Jitter, rng)

	case PatternRing:
		count := ctx.Count
		if count == 0 {
			count = 1
		}

		a := float32(ctx.Index) / float32(count) * (math.Pi * 2)
		return jitter(basis.Azimuth(a), p.Jitter, rng)

	case PatternSphere:
		return rng.UnitSphere()
	case PatternSpokes:
		n := p.Spokes
		if n == 0 {
			n = 1
		}

		k := ctx.Index % n
		az := float32(k) / float32(n) * (math.Pi * 2)
		sc, cc := sin(p.Cone), cos(p.Cone)
		d := basis.Apply(vm.New(sc*cos(az), sc*sin(az), cc))
		return jitter(d, p.Spread, rng)
	}
	return vm.Zero
}
