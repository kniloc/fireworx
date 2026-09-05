package sim

import vm "fireworx/vec"

type World struct {
	Gravity           vm.Vec3
	SparkCap, StarCap int
	Wind              vm.Vec3
}

func DefaultWorld() World {
	return World{
		Gravity:  vm.New(0, -9.81, 0),
		SparkCap: 32768,
		StarCap:  4096,
		Wind:     vm.Zero,
	}
}
