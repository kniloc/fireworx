package sim

import (
	vm "fireworx/vec"
	"math"
	"math/rand/v2"
)

type Rand struct {
	r *rand.Rand
}

func NewRand(seed uint64) *Rand {
	return &Rand{r: rand.New(rand.NewPCG(seed, seed^0x9E3779B97F4A7C15))}
}

func (r *Rand) Float32() float32 {
	return float32(r.r.Float64())
}

func (r *Rand) Uint64() uint64 {
	return r.r.Uint64()
}

func (r *Rand) IntRange(lo, hi uint16) uint16 {
	if hi <= lo {
		return lo
	}

	return lo + uint16(r.r.IntN(int(hi-lo+1)))
}

func mixSeed(seed, k uint64) uint64 {
	x := seed ^ (k*0x9E3779B97F4A7C15 + 0x9E3779B97F4A7C15)
	x ^= x >> 30
	x *= 0xBF58476D1CE4E5B9
	x ^= x >> 27
	x *= 0x94D049BB133111EB
	x ^= x >> 31
	return x
}

func (r *Rand) UnitCap(cosMin float32) vm.Vec3 {
	z := cosMin + (1-cosMin)*r.Float32()
	phi := r.Float32() * float32(math.Pi) * 2
	rad := float32(math.Sqrt(math.Max(0, float64(1-z*z))))

	return vm.Vec3{
		X: rad * float32(math.Cos(float64(phi))),
		Y: rad * float32(math.Sin(float64(phi))),
		Z: z,
	}
}

func (r *Rand) UnitSphere() vm.Vec3 {
	return r.UnitCap(-1.0)
}
