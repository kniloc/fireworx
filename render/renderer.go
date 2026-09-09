package render

import (
	"fireworx/sim"
	vm "fireworx/vec"
)

const (
	Width          = 960
	Height         = 640
	skylineH       = 0
	focalLength    = 430.0
	horizonY       = Height - skylineH
	persistence    = 205
	pixelsPerMeter = 9.0
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
	// Per-pixel skyline height, built by stamping buildings left-to-right.
	profile := [Width]int{}

	type building struct {
		w, h int
	}

	seq := []building{
		{22, 80}, {8, 55}, {14, 115}, {30, 90}, {10, 140}, {18, 65},
		{26, 128}, {6, 42}, {12, 95}, {20, 148}, {16, 118}, {28, 72},
		{10, 152}, {24, 108}, {40, 62}, {12, 135}, {18, 80},
		{14, 145}, {26, 122}, {10, 100}, {22, 155}, {8, 50}, {16, 70},
		{30, 132}, {12, 112}, {20, 88}, {6, 60}, {28, 140}, {10, 148},
		{22, 120}, {18, 92}, {26, 135}, {6, 48}, {12, 62},
		{32, 128}, {16, 105}, {20, 150}, {10, 85}, {24, 118},
		{14, 95}, {6, 55}, {28, 142}, {18, 78}, {12, 125}, {22, 100},
		{30, 82}, {10, 132}, {16, 108}, {8, 45}, {26, 92}, {14, 120},
		{20, 68}, {18, 138}, {24, 88}, {12, 112}, {28, 75},
	}

	x := 0
	for _, b := range seq {
		if x >= Width {
			break
		}
		xEnd := min(x+b.w, Width)
		for px := x; px < xEnd; px++ {
			profile[px] = b.h
		}
		x = xEnd
	}
	for ; x < Width; x++ {
		profile[x] = 55
	}

	// Draw silhouette from profile
	for x := range Width {
		top := max(Height-profile[x], 0)
		for y := top; y < Height; y++ {
			setRGB(frame, (y*Width+x)*4, 20, 22, 46)
		}
	}

	// Draw windows per building
	x = 0
	for _, b := range seq {
		if x >= Width {
			break
		}
		xEnd := min(x+b.w, Width)
		top := max(Height-b.h, 0)
		for wy := top + 4; wy+3 < Height-3; wy += 7 {
			for wx := x + 3; wx+2 < xEnd-2; wx += 6 {
				if (wx/6+wy/7)%4 != 0 {
					for dy := range 3 {
						for dx := range 2 {
							setRGB(frame, ((wy+dy)*Width+(wx+dx))*4, 58, 46, 18)
						}
					}
				}
			}
		}
		x += b.w
	}
}

func setRGB(frame []byte, offset int, r, g, b uint8) {
	frame[offset], frame[offset+1], frame[offset+2], frame[offset+3] = r, g, b, 0xff
}
