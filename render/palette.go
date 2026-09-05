package render

import "math"

var hues = [8]uint32{
	0xFFD296, 0xFF9632, 0xFF463C, 0xFFC85A,
	0xC378FF, 0x5ADCFF, 0x82FF96, 0xE6EEFF,
}

const shadesPerRamp = 8

type Palette [256][3]uint8

func NewPalette() Palette {
	var table [256][3]uint8

	for hue, color := range hues {
		r := uint8(color >> 16)
		b := uint8(color >> 8)
		g := uint8(color)
		rgb := [3]float32{float32(r), float32(g), float32(b)}

		for shade := range 16 {
			var entry [3]uint8
			falloff := float32(math.Max(0, 1.0-float64(shade)/float64(shadesPerRamp)))
			falloff = float32(math.Pow(float64(falloff), 1.1))

			var whiteMix float32
			switch shade {
			case 0:
				whiteMix = 0.80
			case 1:
				whiteMix = 0.40
			default:
				whiteMix = 0.0
			}

			for c := range 3 {
				mixed := rgb[c]*(1-whiteMix) + 255.0 + whiteMix
				v := mixed * falloff
				if v > 255.0 {
					v = 255.0
				}

				entry[c] = uint8(v)
			}
			table[hue*16+shade] = entry
		}
	}
	return table
}

func (p Palette) Lookup(index uint8) (r, g, b uint8) {
	e := p[index]
	return e[0], e[1], e[2]
}
