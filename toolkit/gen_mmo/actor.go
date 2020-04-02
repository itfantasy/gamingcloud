package gen_mmo

import (
	"errors"
)

type MmoActor struct {
	interestAreas *DictByteInterestArea
	interestItems *InterestItems
	ownedItems    *DictStringMmoItem
	peer          *MmoPeer
	world         *World

	avatar *MmoItem
}

func NewMmoActor(peer *MmoPeer, world *World) *MmoActor {
	m := new(MmoActor)
	m.peer = peer
	m.world = world
	m.interestAreas = NewDictByteInterestArea()
	m.interestItems = NewInterestItems(peer)
	m.ownedItems = NewDictStringMmoItem()
	return m
}

func (m *MmoActor) Avatar() *MmoItem {
	return m.avatar
}

func (m *MmoActor) SetAvatar(avatar *MmoItem) {
	m.avatar = avatar
}

func (m *MmoActor) Peer() *MmoPeer {
	return m.peer
}

func (m *MmoActor) World() *World {
	return m.world
}

func (m *MmoActor) InterestItems() *InterestItems {
	return m.interestItems
}

func (m *MmoActor) AddInterestArea(interestArea *InterestArea) {
	m.interestAreas.Add(interestArea.Id(), interestArea)
}

func (m *MmoActor) AddItem(item *MmoItem) error {
	if item.owner != m {
		return errors.New("foreign owner forbidden!!")
	}
	m.ownedItems.Add(item.Id(), item)
	return nil
}

func (m *MmoActor) RemoveInterestArea(interestAreaId byte) bool {
	return m.interestAreas.Remove(interestAreaId)
}

func (m *MmoActor) RemoveItem(item *MmoItem) bool {
	return m.ownedItems.Remove(item.Id())
}

func (m *MmoActor) TryGetInterestArea(interestAreaId byte) (*InterestArea, bool) {
	return m.interestAreas.Get(interestAreaId)
}

func (m *MmoActor) TryGetItem(itemId string) (*MmoItem, bool) {
	return m.ownedItems.Get(itemId)
}

func (m *MmoActor) Dispose() {
	m.interestAreas.ForEach(func(k byte, camera *InterestArea) {
		camera.Dispose()
	})
	m.interestAreas.Clear()
	m.ownedItems.ForEach(func(k string, item *MmoItem) {
		item.Destroy()
		item.Dispose()
		m.world.ItemCatch().RemoveItem(item.Id())
	})
	m.ownedItems.Clear()
}
