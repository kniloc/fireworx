package render

import (
	"fireworx/sim"
	vm "fireworx/vec"
)

const (
	Width          = 960
	Height         = 640
	focalLength    = 430.0
	horizonY       = Height - 20.0
	persistence    = 205
	pixelsPerMeter = 9.0
	skylineH       = 18
)

type Renderer struct {
	palette Palette
}

func New() *Renderer {
	return &Renderer{
		palette: NewPalette(),
	}
}

func (r *Renderer) Draw(fw *sim.Fireworks, frame []byte) {
	skyBytes := (Height - skylineH) * Width * 4
	r.fade(frame[:skyBytes])
	fw.Visit(func(pos vm.Vec3, index uint8, isStar bool) {
		r.plot(frame, pos, index, isStar)
	})
	r.drawSkyline(frame)
}

func Unproject(sx, sy float32) vm.Vec3 {
	return vm.New(
		(sx-Width*0.5)/pixelsPerMeter,
		(horizonY-sy)/pixelsPerMeter,
		0,
	)
}

func (r *Renderer) fade(frame []byte) {
	for i := 0; i+3 < len(frame); i += 4 {
		for c := range 3 {
			faded := (uint16(frame[i+c]) * persistence) >> 8
			if faded > 0 {
				faded--
			}
			frame[i+c] = byte(faded)
		}
		frame[i+3] = 0xff
	}
}

func (r *Renderer) plot(frame []byte, pos vm.Vec3, index uint8, bright bool) {
	scale := focalLength / (focalLength + pos.Z*pixelsPerMeter)
	sx := int32(Width*0.5 + pos.X*pixelsPerMeter*scale)
	sy := int32(horizonY - pos.Y*pixelsPerMeter*scale)
	if sx < 1 || sy < 1 || sx >= Width-1 || sy >= Height-1 {
		return
	}

	rr, gg, bb := r.palette.Lookup(index)
	offset := (int(sy)*Width + int(sx)) * 4
	addRGB(frame, offset, rr, gg, bb)
	if bright {
		row := Width * 4
		for _, n := range []int{offset - 4, offset + 4, offset - row, offset + row} {
			addRGB(frame, n, rr>>1, gg>>1, bb>>1)
		}
	}
}

func addRGB(frame []byte, offset int, r, g, b uint8) {
	frame[offset] = saturatingAdd(frame[offset], r)
	frame[offset+1] = saturatingAdd(frame[offset+1], g)
	frame[offset+2] = saturatingAdd(frame[offset+2], b)
}

func saturatingAdd(a, b uint8) uint8 {
	sum := uint16(a) + uint16(b)
	if sum > 255 {
		return 255
	}

	return uint8(sum)
}

func (r *Renderer) drawSkyline(frame []byte) {
	ground := Height - skylineH
	for y := ground; y < Height; y++ {
		for x := range Width {
			setRGB(frame, (y*Width+x)*4, 10, 12, 26)
		}
	}

	for i := range 36 {
		bx := i*27 + (i%3)*4
		bh := 10 + (i*37)%26
		top := max(ground-bh, 0)

		xEnd := min(bx+21, Width)

		for y := top; y < ground; y++ {
			for x := bx; x < xEnd; x++ {
				setRGB(frame, (y*Width+x)*4, 17, 20, 42)
			}
		}
	}
}

func setRGB(frame []byte, offset int, r, g, b uint8) {
	frame[offset], frame[offset+1], frame[offset+2], frame[offset+3] = r, g, b, 0xff
}
