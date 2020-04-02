package gen_mmo

type WorldManager struct {
	dict *DictStringWorld
}

func NewWorldManager() *WorldManager {
	w := new(WorldManager)
	w.dict = NewDictStringWorld()
	return w
}

func (w *WorldManager) Clear() {
	w.dict.Clear()
}

func (w *WorldManager) TryCreate(name string, boundingBox *BoundingBox, tileDimensions *Vector) (*World, bool) {
	_, exists := w.dict.Get(name)
	if exists {
		return nil, false
	}
	world := NewWorld(name, boundingBox, tileDimensions)
	w.dict.Add(name, world)
	return world, true
}

func (w *WorldManager) TryGet(name string) (*World, bool) {
	return w.dict.Get(name)
}

func (w *WorldManager) Dispose() {
	w.dict.ForEach(func(name string, world *World) {
		world.Dispose()
	})
}
