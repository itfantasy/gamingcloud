package gen_mmo

import (
	"math"
)

type Vector struct {
	x float64
	y float64
	z float64
}

const (
	TOLERANCE float64 = 0.000001
)

func NewVector2(x, y float64) *Vector {
	v := new(Vector)
	v.x = x
	v.y = y
	v.z = 0
	return v
}

func NewVector(x, y, z float64) *Vector {
	v := new(Vector)
	v.x = x
	v.y = y
	v.z = z
	return v
}

func (v *Vector) X() float64 {
	return v.x
}

func (v *Vector) Y() float64 {
	return v.y
}

func (v *Vector) Z() float64 {
	return v.z
}

func (v *Vector) SetX(x float64) {
	v.x = x
}

func (v *Vector) SetY(y float64) {
	v.y = y
}

func (v *Vector) SetZ(z float64) {
	v.z = z
}

func (v *Vector) AddWith(a *Vector) {
	v.x += a.X()
	v.y += a.Y()
	v.z += a.Z()
}

func (v *Vector) SubtractWith(a *Vector) {
	v.x -= a.X()
	v.y -= a.Y()
	v.z -= a.Z()
}

func (v *Vector) MultiplyWith(a *Vector) {
	v.x *= a.X()
	v.y *= a.Y()
	v.z *= a.Z()
}

func (v *Vector) DivideWith(a *Vector) {
	v.x /= a.X()
	v.y /= a.Y()
	v.z /= a.Z()
}

func (v *Vector) IsZero() bool {
	return math.Abs(v.x) < TOLERANCE && math.Abs(v.y) < TOLERANCE && math.Abs(v.z) < TOLERANCE
}

func (v *Vector) Len2() float64 {
	return v.X()*v.X() + v.Y()*v.Y() + v.Z()*v.Z()
}

func (a *Vector) Add(b *Vector) *Vector {
	return NewVector(a.X()+b.X(), a.Y()+b.Y(), a.Z()+b.Z())
}

func (a *Vector) Subtract(b *Vector) *Vector {
	return NewVector(a.X()-b.X(), a.Y()-b.Y(), a.Z()-b.Z())
}

func (a *Vector) Multiply(b *Vector) *Vector {
	return NewVector(a.X()*b.X(), a.Y()*b.Y(), a.Z()*b.Z())
}

func (a *Vector) Divide(b *Vector) *Vector {
	return NewVector(a.X()/b.X(), a.Y()/b.Y(), a.Z()/b.Z())
}

func (a *Vector) Max(b *Vector) *Vector {
	return NewVector(math.Max(a.X(), b.X()), math.Max(a.Y(), b.Y()), math.Max(a.Z(), b.Z()))
}

func (a *Vector) Min(b *Vector) *Vector {
	return NewVector(math.Min(a.X(), b.X()), math.Min(a.Y(), b.Y()), math.Min(a.Z(), b.Z()))
}

func VAdd(a *Vector, b *Vector) *Vector {
	return a.Add(b)
}

func VSubtract(a *Vector, b *Vector) *Vector {
	return a.Subtract(b)
}

func VMultiply(a *Vector, b *Vector) *Vector {
	return a.Multiply(b)
}

func VDivide(a *Vector, b *Vector) *Vector {
	return a.Divide(b)
}

func VMax(a *Vector, b *Vector) *Vector {
	return a.Max(b)
}

func VMin(a *Vector, b *Vector) *Vector {
	return a.Min(b)
}
