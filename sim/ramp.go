package sim

type Ramp struct {
	Base, Len uint8
}

func NewRamp(hue, length uint8) Ramp {
	return Ramp{
		Base: hue << 4,
		Len:  length,
	}
}

func (r Ramp) Sample(t float32) uint8 {
	if r.Len == 0 {
		return r.Base
	}

	i := min(uint16(t*float32(r.Len)), uint16(r.Len)-1)

	return r.Base + uint8(i)
}
