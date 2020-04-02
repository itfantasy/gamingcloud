package gen_mmo

import (
	"sync"
)

type DictRegionIDisposer struct {
	_map map[*Region]IDisposer
	sync.RWMutex
}

func NewDictRegionIDisposer() *DictRegionIDisposer {
	dict := DictRegionIDisposer{}
	dict._map = make(map[*Region]IDisposer)
	return &dict
}

func NewDictRegionIDisposerRaw(raw map[*Region]IDisposer) *DictRegionIDisposer {
	dict := DictRegionIDisposer{}
	dict._map = raw
	return &dict
}

func (d *DictRegionIDisposer) Add(key *Region, value IDisposer) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		return false
	}
	d._map[key] = value
	return true
}

func (d *DictRegionIDisposer) Remove(key *Region) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		delete(d._map, key)
		return true
	}
	return false
}

func (d *DictRegionIDisposer) Set(key *Region, value IDisposer) {
	d.Lock()
	defer d.Unlock()

	d._map[key] = value
}

func (d *DictRegionIDisposer) Get(key *Region) (IDisposer, bool) {
	d.RLock()
	defer d.RUnlock()

	v, exist := d._map[key]
	return v, exist
}

func (d *DictRegionIDisposer) Len() int {
	d.RLock()
	defer d.RUnlock()

	return len(d._map)
}

func (d *DictRegionIDisposer) ContainsKey(key *Region) bool {
	d.RLock()
	defer d.RUnlock()

	_, exist := d._map[key]
	return exist
}

func (d *DictRegionIDisposer) ContainsValue(value IDisposer) bool {
	d.RLock()
	defer d.RUnlock()

	for _, v := range d._map {
		if v == value {
			return true
		}
	}
	return false
}

func (d *DictRegionIDisposer) ForEach(fun func(*Region, IDisposer)) {
	d.RLock()
	defer d.RUnlock()

	for k, v := range d._map {
		fun(k, v)
	}
}

func (d *DictRegionIDisposer) KeyValuePairs() map[*Region]IDisposer {
	d.RLock()
	defer d.RUnlock()

	ret := make(map[*Region]IDisposer)
	for k, v := range d._map {
		ret[k] = v
	}
	return ret
}

func (d *DictRegionIDisposer) Clear() {
	d.Lock()
	defer d.Unlock()

	for k, _ := range d._map {
		delete(d._map, k)
	}
}
