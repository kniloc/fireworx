package sim

import vm "fireworx/vec"

type Spark struct {
	Age     float32
	Drag    float32
	Gravity float32
	Life    float32
	Pos     vm.Vec3
	Ramp    Ramp
	Vel     vm.Vec3
}

func (s *Spark) T() float32 {
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
