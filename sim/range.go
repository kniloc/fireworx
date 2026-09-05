package sim

type Range struct {
	Min, Max float32
}

func NewRange(min, max float32) Range {
	return Range{min, max}
}

func RangeAt(v float32) Range {
	return Range{v, v}
}

func (r Range) Sample(rng *Rand) float32 {
	return r.Min + (r.Max-r.Min)*rng.Float32()
}
