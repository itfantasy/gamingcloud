package gen_mmo

import (
	"sync"
)

type DictMmoItemIDisposer struct {
	_map map[*MmoItem]IDisposer
	sync.RWMutex
}

func NewDictMmoItemIDisposer() *DictMmoItemIDisposer {
	dict := DictMmoItemIDisposer{}
	dict._map = make(map[*MmoItem]IDisposer)
	return &dict
}

func NewDictMmoItemIDisposerRaw(raw map[*MmoItem]IDisposer) *DictMmoItemIDisposer {
	dict := DictMmoItemIDisposer{}
	dict._map = raw
	return &dict
}

func (d *DictMmoItemIDisposer) Add(key *MmoItem, value IDisposer) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		return false
	}
	d._map[key] = value
	return true
}

func (d *DictMmoItemIDisposer) Remove(key *MmoItem) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		delete(d._map, key)
		return true
	}
	return false
}

func (d *DictMmoItemIDisposer) Set(key *MmoItem, value IDisposer) {
	d.Lock()
	defer d.Unlock()

	d._map[key] = value
}

func (d *DictMmoItemIDisposer) Get(key *MmoItem) (IDisposer, bool) {
	d.RLock()
	defer d.RUnlock()

	v, exist := d._map[key]
	return v, exist
}

func (d *DictMmoItemIDisposer) Len() int {
	d.RLock()
	defer d.RUnlock()

	return len(d._map)
}

func (d *DictMmoItemIDisposer) ContainsKey(key *MmoItem) bool {
	d.RLock()
	defer d.RUnlock()

	_, exist := d._map[key]
	return exist
}

func (d *DictMmoItemIDisposer) ContainsValue(value IDisposer) bool {
	d.RLock()
	defer d.RUnlock()

	for _, v := range d._map {
		if v == value {
			return true
		}
	}
	return false
}

func (d *DictMmoItemIDisposer) ForEach(fun func(*MmoItem, IDisposer)) {
	d.RLock()
	defer d.RUnlock()

	for k, v := range d._map {
		fun(k, v)
	}
}

func (d *DictMmoItemIDisposer) KeyValuePairs() map[*MmoItem]IDisposer {
	d.RLock()
	defer d.RUnlock()

	ret := make(map[*MmoItem]IDisposer)
	for k, v := range d._map {
		ret[k] = v
	}
	return ret
}

func (d *DictMmoItemIDisposer) Clear() {
	d.Lock()
	defer d.Unlock()

	for k, _ := range d._map {
		delete(d._map, k)
	}
}
