package gen_mmo

import (
	"math"

	"github.com/itfantasy/gonode/utils/stl"
)

type World struct {
	name           string
	worldRegions   [][]*Region
	area           *BoundingBox
	tileDimensions *Vector
	tileX          int
	tileY          int
	itemCatch      *ItemManager
	radar          *Radar
}

func NewWorld(name string, boundingBox *BoundingBox, tileDimensions *Vector) *World {
	w := NewGridWorld(boundingBox, tileDimensions)
	w.name = name
	w.itemCatch = NewItemManager()
	w.radar = NewRadar()
	return w
}

func NewGridWorld(area *BoundingBox, tileDimensions *Vector) *World {
	w := new(World)
	w.area = area
	w.tileDimensions = tileDimensions
	w.tileX = int(area.Size().X() / tileDimensions.X())
	w.tileY = int(area.Size().Y() / tileDimensions.Y())

	regions := make([][]*Region, 0, w.tileX)
	for x := 0; x < w.tileX; x++ {
		sub := make([]*Region, 0, w.tileY)
		for y := 0; y < w.tileY; y++ {
			sub = append(sub, NewRegion(x, y))
		}
		regions = append(regions, sub)
	}

	w.worldRegions = regions
	return w
}

func (w *World) Name() string {
	return w.name
}

func (w *World) Area() *BoundingBox {
	return w.area
}

func (w *World) TileDimensions() *Vector {
	return w.tileDimensions
}

func (w *World) TileX() int {
	return w.tileX
}

func (w *World) TileY() int {
	return w.tileY
}

func (w *World) ItemCatch() *ItemManager {
	return w.itemCatch
}

func (w *World) Radar() *Radar {
	return w.radar
}

func (w *World) GetRegion(position *Vector) (*Region, bool) {
	p := VSubtract(position, w.Area().Min())
	if p.X() >= 0 && p.X() < w.Area().Size().X() &&
		p.Y() >= 0 && p.Y() < w.Area().Size().Y() {
		x := int(p.X() / w.TileDimensions().X())
		y := int(p.Y() / w.TileDimensions().Y())
		return w.worldRegions[x][y], true
	} else {
		return nil, false
	}
}

func (w *World) GetRegions(area *BoundingBox) *stl.HashSet {
	overlap := w.Area().IntersectWith(area)
	min := VSubtract(overlap.Min(), w.Area().Min())
	max := VSubtract(overlap.Max(), w.Area().Min())
	x0 := int(math.Max(min.X()/w.TileDimensions().X(), 0))
	x1 := int(math.Min(float64(math.Ceil(max.X()/w.TileDimensions().X())), float64(w.TileX())))
	y0 := int(math.Max(min.Y()/w.TileDimensions().Y(), 0))
	y1 := int(math.Min(float64(math.Ceil(max.Y()/w.TileDimensions().Y())), float64(w.TileY())))

	regions := stl.NewHashSet()
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			regions.Add(w.worldRegions[x][y])
		}
	}
	return regions
}

func (w *World) Dispose() {
	for _, regions := range w.worldRegions {
		for _, region := range regions {
			region.Dispose()
		}
	}
}
