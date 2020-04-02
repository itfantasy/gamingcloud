package gen_mmo

import (
	"sync"
)

type DictStringWorld struct {
	_map map[string]*World
	sync.RWMutex
}

func NewDictStringWorld() *DictStringWorld {
	dict := DictStringWorld{}
	dict._map = make(map[string]*World)
	return &dict
}

func NewDictStringWorldRaw(raw map[string]*World) *DictStringWorld {
	dict := DictStringWorld{}
	dict._map = raw
	return &dict
}

func (d *DictStringWorld) Add(key string, value *World) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		return false
	}
	d._map[key] = value
	return true
}

func (d *DictStringWorld) Remove(key string) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		delete(d._map, key)
		return true
	}
	return false
}

func (d *DictStringWorld) Set(key string, value *World) {
	d.Lock()
	defer d.Unlock()

	d._map[key] = value
}

func (d *DictStringWorld) Get(key string) (*World, bool) {
	d.RLock()
	defer d.RUnlock()

	v, exist := d._map[key]
	return v, exist
}

func (d *DictStringWorld) Len() int {
	d.RLock()
	defer d.RUnlock()

	return len(d._map)
}

func (d *DictStringWorld) ContainsKey(key string) bool {
	d.RLock()
	defer d.RUnlock()

	_, exist := d._map[key]
	return exist
}

func (d *DictStringWorld) ContainsValue(value *World) bool {
	d.RLock()
	defer d.RUnlock()

	for _, v := range d._map {
		if v == value {
			return true
		}
	}
	return false
}

func (d *DictStringWorld) ForEach(fun func(string, *World)) {
	d.RLock()
	defer d.RUnlock()

	for k, v := range d._map {
		fun(k, v)
	}
}

func (d *DictStringWorld) KeyValuePairs() map[string]*World {
	d.RLock()
	defer d.RUnlock()

	ret := make(map[string]*World)
	for k, v := range d._map {
		ret[k] = v
	}
	return ret
}

func (d *DictStringWorld) Clear() {
	d.Lock()
	defer d.Unlock()

	for k, _ := range d._map {
		delete(d._map, k)
	}
}
