package gen_mmo

import (
	"errors"
)

type BoundingBox struct {
	max *Vector
	min *Vector
}

func NewBoundingBox(min, max *Vector) *BoundingBox {
	b := new(BoundingBox)
	b.min = min
	b.max = max
	return b
}

func NewBoundingBoxFromPoints(points ...*Vector) (*BoundingBox, error) {
	if points == nil {
		return nil, errors.New("the points can not be nil!")
	}
	if len(points) <= 0 {
		return nil, errors.New("the points' len can not be zero!")
	}
	min := points[0]
	max := points[1]
	for _, point := range points {
		min = VMin(min, point)
		max = VMax(max, point)
	}
	return NewBoundingBox(min, max), nil
}

func (b *BoundingBox) Max() *Vector {
	return b.max
}

func (b *BoundingBox) Min() *Vector {
	return b.min
}

func (b *BoundingBox) SetMax(max *Vector) {
	b.max = max
}

func (b *BoundingBox) SetMin(min *Vector) {
	b.min = min
}

func (b *BoundingBox) Size() *Vector {
	return VSubtract(b.Max(), b.Min())
}

func (b *BoundingBox) Contains(point *Vector) bool {
	return (point.X() < b.Min().X() || point.X() > b.Max().X() ||
		point.Y() < b.Min().Y() || point.Y() > b.Max().Y() ||
		point.Z() < b.Min().Z() || point.Z() > b.Max().Z()) == false
}

func (b *BoundingBox) Contains2d(point *Vector) bool {
	return (point.X() < b.Min().X() || point.X() > b.Max().X() ||
		point.Y() < b.Min().Y() || point.Y() > b.Max().Y()) == false
}

func (b *BoundingBox) IntersectWith(other *BoundingBox) *BoundingBox {
	return NewBoundingBox(VMax(b.Min(), other.Min()), VMin(b.Max(), other.Max()))
}

func (b *BoundingBox) UnionWith(other *BoundingBox) *BoundingBox {
	return NewBoundingBox(VMin(b.Min(), other.Min()), VMax(b.Max(), other.Max()))
}

func (b *BoundingBox) IsValid() bool {
	return (b.Max().X() < b.Min().X() || b.Max().Y() < b.Min().Y() || b.Max().Z() < b.Min().Z()) == false
}
