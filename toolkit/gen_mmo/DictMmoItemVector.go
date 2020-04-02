package gen_mmo

import (
	"sync"
)

type DictMmoItemVector struct {
	_map map[*MmoItem]*Vector
	sync.RWMutex
}

func NewDictMmoItemVector() *DictMmoItemVector {
	dict := DictMmoItemVector{}
	dict._map = make(map[*MmoItem]*Vector)
	return &dict
}

func NewDictMmoItemVectorRaw(raw map[*MmoItem]*Vector) *DictMmoItemVector {
	dict := DictMmoItemVector{}
	dict._map = raw
	return &dict
}

func (d *DictMmoItemVector) Add(key *MmoItem, value *Vector) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		return false
	}
	d._map[key] = value
	return true
}

func (d *DictMmoItemVector) Remove(key *MmoItem) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		delete(d._map, key)
		return true
	}
	return false
}

func (d *DictMmoItemVector) Set(key *MmoItem, value *Vector) {
	d.Lock()
	defer d.Unlock()

	d._map[key] = value
}

func (d *DictMmoItemVector) Get(key *MmoItem) (*Vector, bool) {
	d.RLock()
	defer d.RUnlock()

	v, exist := d._map[key]
	return v, exist
}

func (d *DictMmoItemVector) Len() int {
	d.RLock()
	defer d.RUnlock()

	return len(d._map)
}

func (d *DictMmoItemVector) ContainsKey(key *MmoItem) bool {
	d.RLock()
	defer d.RUnlock()

	_, exist := d._map[key]
	return exist
}

func (d *DictMmoItemVector) ContainsValue(value *Vector) bool {
	d.RLock()
	defer d.RUnlock()

	for _, v := range d._map {
		if v == value {
			return true
		}
	}
	return false
}

func (d *DictMmoItemVector) ForEach(fun func(*MmoItem, *Vector)) {
	d.RLock()
	defer d.RUnlock()

	for k, v := range d._map {
		fun(k, v)
	}
}

func (d *DictMmoItemVector) KeyValuePairs() map[*MmoItem]*Vector {
	d.RLock()
	defer d.RUnlock()

	ret := make(map[*MmoItem]*Vector)
	for k, v := range d._map {
		ret[k] = v
	}
	return ret
}

func (d *DictMmoItemVector) Clear() {
	d.Lock()
	defer d.Unlock()

	for k, _ := range d._map {
		delete(d._map, k)
	}
}
