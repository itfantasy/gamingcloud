package gen_mmo

import (
	"sync"
)

type DictByteInterestArea struct {
	_map map[byte]*InterestArea
	sync.RWMutex
}

func NewDictByteInterestArea() *DictByteInterestArea {
	dict := DictByteInterestArea{}
	dict._map = make(map[byte]*InterestArea)
	return &dict
}

func NewDictByteInterestAreaRaw(raw map[byte]*InterestArea) *DictByteInterestArea {
	dict := DictByteInterestArea{}
	dict._map = raw
	return &dict
}

func (d *DictByteInterestArea) Add(key byte, value *InterestArea) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		return false
	}
	d._map[key] = value
	return true
}

func (d *DictByteInterestArea) Remove(key byte) bool {
	d.Lock()
	defer d.Unlock()

	_, exist := d._map[key]
	if exist {
		delete(d._map, key)
		return true
	}
	return false
}

func (d *DictByteInterestArea) Set(key byte, value *InterestArea) {
	d.Lock()
	defer d.Unlock()

	d._map[key] = value
}

func (d *DictByteInterestArea) Get(key byte) (*InterestArea, bool) {
	d.RLock()
	defer d.RUnlock()

	v, exist := d._map[key]
	return v, exist
}

func (d *DictByteInterestArea) Len() int {
	d.RLock()
	defer d.RUnlock()

	return len(d._map)
}

func (d *DictByteInterestArea) ContainsKey(key byte) bool {
	d.RLock()
	defer d.RUnlock()

	_, exist := d._map[key]
	return exist
}

func (d *DictByteInterestArea) ContainsValue(value *InterestArea) bool {
	d.RLock()
	defer d.RUnlock()

	for _, v := range d._map {
		if v == value {
			return true
		}
	}
	return false
}

func (d *DictByteInterestArea) ForEach(fun func(byte, *InterestArea)) {
	d.RLock()
	defer d.RUnlock()

	for k, v := range d._map {
		fun(k, v)
	}
}

func (d *DictByteInterestArea) KeyValuePairs() map[byte]*InterestArea {
	d.RLock()
	defer d.RUnlock()

	ret := make(map[byte]*InterestArea)
	for k, v := range d._map {
		ret[k] = v
	}
	return ret
}

func (d *DictByteInterestArea) Clear() {
	d.Lock()
	defer d.Unlock()

	for k, _ := range d._map {
		delete(d._map, k)
	}
}
