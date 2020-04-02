package gen_mmo

import (
	"sync"
)

type DictStringMmoItem struct {
	_map map[string]*MmoItem
	sync.RWMutex
}

func NewDictStringMmoItem() *DictStringMmoItem {
	dict := DictStringMmoItem{}
	dict._map = make(map[string]*MmoItem)
	return &dict
}

func NewDictStringMmoItemRaw(raw map[string]*MmoItem) *DictStringMmoItem {
	dict := DictStringMmoItem{}
	dict._map = raw
	return &dict
}

func (d *DictStringMmoItem) Add(key string, value *MmoItem) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		return false
	}
	d._map[key] = value
	return true
}

func (d *DictStringMmoItem) Remove(key string) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		delete(d._map, key)
		return true
	}
	return false
}

func (d *DictStringMmoItem) Set(key string, value *MmoItem) {
	d.Lock()
	defer d.Unlock()

	d._map[key] = value
}

func (d *DictStringMmoItem) Get(key string) (*MmoItem, bool) {
	d.RLock()
	defer d.RUnlock()

	v, exist := d._map[key]
	return v, exist
}

func (d *DictStringMmoItem) Len() int {
	d.RLock()
	defer d.RUnlock()

	return len(d._map)
}

func (d *DictStringMmoItem) ContainsKey(key string) bool {
	d.RLock()
	defer d.RUnlock()

	_, exist := d._map[key]
	return exist
}

func (d *DictStringMmoItem) ContainsValue(value *MmoItem) bool {
	d.RLock()
	defer d.RUnlock()

	for _, v := range d._map {
		if v == value {
			return true
		}
	}
	return false
}

func (d *DictStringMmoItem) ForEach(fun func(string, *MmoItem)) {
	d.RLock()
	defer d.RUnlock()

	for k, v := range d._map {
		fun(k, v)
	}
}

func (d *DictStringMmoItem) KeyValuePairs() map[string]*MmoItem {
	d.RLock()
	defer d.RUnlock()

	ret := make(map[string]*MmoItem)
	for k, v := range d._map {
		ret[k] = v
	}
	return ret
}

func (d *DictStringMmoItem) Clear() {
	d.Lock()
	defer d.Unlock()

	for k, _ := range d._map {
		delete(d._map, k)
	}
}
