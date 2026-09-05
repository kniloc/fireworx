package sim

import vm "fireworx/vec"

type Star struct {
	Age      float32
	Effect   EffectID
	Life     float32
	Pos      vm.Vec3
	Seed     uint64
	Stage    uint8
	TrailAcc float32
	Vel      vm.Vec3
}

func (s *Star) T() float32 {
	life := s.Life

	if life < 1e-6 {
		life = 1e-6
	}

	t := s.Age / life
	if t > 1.0 {
		t = 1.0
	}

	return t
}
