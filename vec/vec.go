package vm

import "math"

type Vec3 struct {
	X, Y, Z float32
}

func New(x, y, z float32) Vec3 {
	return Vec3{x, y, z}
}

func (a Vec3) Add(b Vec3) Vec3 {
	return Vec3{
		a.X + b.X,
		a.Y + b.Y,
		a.Z + b.Z,
	}
}

func (a Vec3) Sub(b Vec3) Vec3 {
	return Vec3{
		a.X - b.X,
		a.Y - b.Y,
		a.Z - b.Z,
	}
}

func (a Vec3) Scale(s float32) Vec3 {
	return Vec3{
		a.X * s,
		a.Y * s,
		a.Z * s,
	}
}

func (a Vec3) Length() float32 {
	return float32(math.Sqrt(float64(a.Dot(a))))
}

func (a Vec3) Dot(b Vec3) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vec3) Cross(b Vec3) Vec3 {
	return Vec3{
		a.Y*b.Z - a.Z*b.Y,
		a.Z*b.X - a.X*b.Z,
		a.X*b.Y - a.Y*b.X,
	}
}

func (a Vec3) NormalizeOr(fallback Vec3) Vec3 {
	l := a.Length()
	if l < 1e-6 {
		return fallback
	}
	return a.Scale(1 / l)
}

func (a Vec3) Lerp(b Vec3, t float32) Vec3 {
	return a.Add(b.Sub(a).Scale(t))
}

func (a Vec3) AnyOrthonormalPair() (u, v Vec3) {
	sign := float32(1.0)
	if a.Z < 0 {
		sign = -1.0
	}
	s := -1.0 / (sign + a.Z)
	b := a.X * a.Y * s
	u = Vec3{1.0 + sign*a.X*a.X*s, sign * b, -sign * a.X}
	v = Vec3{b, sign + a.Y*a.Y*s, -a.Y}
	return
}

var Y = Vec3{0, 1, 0}
var Z = Vec3{0, 0, 1}
var Zero = Vec3{}
