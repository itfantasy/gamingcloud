package gen_mmo

import (
	"sync"
)

type ItemManager struct {
	items map[string]*MmoItem
	sync.RWMutex
}

func NewItemManager() *ItemManager {
	i := new(ItemManager)
	i.items = make(map[string]*MmoItem)
	return i
}

func (i *ItemManager) AddItem(item *MmoItem) bool {
	i.Lock()
	defer i.Unlock()

	if _, exists := i.items[item.id]; exists {
		return false
	}
	i.items[item.Id()] = item
	return true
}

func (i *ItemManager) RemoveItem(itemId string) bool {
	i.Lock()
	defer i.Unlock()

	if _, exists := i.items[itemId]; exists {
		delete(i.items, itemId)
		return true
	}
	return false
}

func (i *ItemManager) GetItem(itemId string) (*MmoItem, bool) {
	i.Lock()
	defer i.Unlock()

	item, exists := i.items[itemId]
	return item, exists
}

func (i ItemManager) Dispose() {
	i.items = nil
}
