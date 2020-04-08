package gen_mmo

import (
	"errors"
	"sync"

	"github.com/itfantasy/gonode/core/actors"
	"github.com/itfantasy/gonode/utils/stl"
)

type InterestArea struct {
	id                       byte
	requestItemEnterMessage  *RequestItemEnterMessage
	requestItemExitMessage   *RequestItemExitMessage
	regions                  *stl.HashSet
	subManagementExecutor    *actors.Executor
	world                    *World
	currentInnerFocus        *BoundingBox
	currentOuterFocus        *BoundingBox
	itemMovementSubscription IDisposer
	regionSubscriptions      *DictRegionIDisposer

	attachedItem      *MmoItem
	position          *Vector
	viewDistanceEnter *Vector
	viewDistanceExit  *Vector

	eventChannelSubscriptions *DictRegionIDisposer
	peer                      *MmoPeer

	sync.RWMutex
}

func NewInterestArea(peer *MmoPeer, id byte, world *World) *InterestArea {
	i := new(InterestArea)
	i.id = id
	i.world = world
	i.requestItemEnterMessage = NewRequestItemEnterMessage(i)
	i.requestItemExitMessage = NewRequestItemExitMessage(i)
	i.regions = stl.NewHashSet()
	i.subManagementExecutor = actors.Spawn(1024)
	i.regionSubscriptions = NewDictRegionIDisposer()

	i.peer = peer
	i.eventChannelSubscriptions = NewDictRegionIDisposer()
	return i
}

func (i *InterestArea) Sync(fun func()) {
	i.Lock()
	defer i.Unlock()

	fun()
}

func (i *InterestArea) AttachedItem() *MmoItem {
	return i.attachedItem
}

func (i *InterestArea) Id() byte {
	return i.id
}

func (i *InterestArea) Position() *Vector {
	return i.position
}

func (i *InterestArea) SetPosition(position *Vector) {
	i.position = position
}

func (i *InterestArea) ViewDistanceEnter() *Vector {
	return i.viewDistanceEnter
}

func (i *InterestArea) SetViewDistanceEnter(v *Vector) {
	i.viewDistanceEnter = v
}

func (i *InterestArea) ViewDistanceExit() *Vector {
	return i.viewDistanceExit
}

func (i *InterestArea) SetViewDistanceExit(v *Vector) {
	i.viewDistanceExit = v
}

func (i *InterestArea) AttachToItem(item *MmoItem) error {
	if i.attachedItem != nil {
		return errors.New("invalid operation!! there has been a attached item!")
	}
	i.attachedItem = item
	i.position = item.position
	disposeSubscription := item.DisposeChannel().Subscribe(i.subManagementExecutor, i.AttachedItem_OnItemDisposed)
	positionSubscription := item.PositionUpdateChannel().Subscribe(i.subManagementExecutor, i.AttachedItem_OnItemPosition)
	i.itemMovementSubscription = NewUnsubscriberCollection(disposeSubscription, positionSubscription)
	return nil
}

func (i *InterestArea) Detach() {
	if i.attachedItem != nil {
		i.itemMovementSubscription.Dispose()
		i.itemMovementSubscription = nil

		i.attachedItem = nil
	}
}

func (i *InterestArea) UpdateInterestManagement() {
	focus, _ := NewBoundingBoxFromPoints(VSubtract(i.position, i.viewDistanceExit), VAdd(i.position, i.viewDistanceExit))
	i.currentOuterFocus = focus.IntersectWith(i.world.Area())

	focus2 := NewBoundingBox(VSubtract(i.position, i.viewDistanceEnter), VAdd(i.position, i.viewDistanceEnter))
	i.currentInnerFocus = focus2.IntersectWith(i.world.Area())

	i.SubscribeRegions(i.world.GetRegions(i.currentInnerFocus))

}

func (i *InterestArea) OnRegionEnter(region *Region) {
	subscription := region.ItemRegionChangedChannel().Subscribe(i.subManagementExecutor, i.OnItemRegionChange)
	i.regionSubscriptions.Add(region, subscription)

	messageReceiver := region.ItemEventChannel().Subscribe(i.peer.RequestExecutor(), i.Region_OnItemEvent)
	i.eventChannelSubscriptions.Set(region, messageReceiver)
}

func (i *InterestArea) OnRegionExit(region *Region) {
	if subscription, exists := i.regionSubscriptions.Get(region); exists {
		subscription.Dispose()
		i.regionSubscriptions.Remove(region)
	}

	if messageReceiver, exists := i.eventChannelSubscriptions.Get(region); exists {
		i.eventChannelSubscriptions.Remove(region)
		messageReceiver.Dispose()
	}
}

func (i *InterestArea) OnItemRegionChange(msg interface{}) {
	message := msg.(*ItemRegionChangedMessage)
	r0 := i.regions.Contains(message.Region0())
	r1 := i.regions.Contains(message.Region1())
	if r0 && r1 {
		// nothing to do
	} else if r0 {
		i.OnItemExit(message.ItemSnapshot().Source())
	} else if r1 {
		i.OnItemEnter(message.ItemSnapshot())
	}
}

func (i *InterestArea) OnItemEnter(snapshot *ItemSnapshot) {
	item := snapshot.Source()
	i.peer.MmoEventer().OnItemSubscribed(i.peer, &ItemSubscribed{
		InterestAreaId:     i.Id(),
		ItemId:             item.Id(),
		ItemType:           item.Type(),
		Position:           snapshot.Position(),
		Rotation:           snapshot.Rotation(),
		PropertiesRevision: snapshot.PropertiesRevision(),
	})
}

func (i *InterestArea) OnItemExit(item *MmoItem) {
	i.peer.MmoEventer().OnItemUnsubscribed(i.peer, &ItemUnsubscribed{
		InterestAreaId: i.Id(),
		ItemId:         item.Id(),
	})
}

func (i *InterestArea) AttachedItem_OnItemDisposed(msg interface{}) {
	message := msg.(*ItemDisposedMessage)
	i.Sync(func() {
		if message.source == i.attachedItem {
			i.Detach()
		}
	})
}

func (i *InterestArea) AttachedItem_OnItemPosition(msg interface{}) {
	message := msg.(*ItemPositionMessage)
	i.Sync(func() {
		if message.source == i.attachedItem {
			i.position = message.position
			i.UpdateInterestManagement()
		}
	})
}

func (i *InterestArea) Region_OnItemEvent(msg interface{}) {
	m := msg.(ItemEventMessage)
	switch m.Code() {
	case Event_ItemDestroyed:
		message := m.Data().(*ItemDestroyed)
		i.peer.MmoEventer().OnItemDestroyed(i.peer, message.ItemId)
		break
	case Event_ItemMoved:
		message := m.Data().(*ItemMoved)
		i.peer.MmoEventer().OnItemMoved(i.peer, message)
		break
	case Event_ItemPropertiesSet:
		message := m.Data().(*ItemPropertiesSet)
		i.peer.MmoEventer().OnItemPropertiesSet(i.peer, message)
		break
	case Event_ItemGeneric:
		message := m.Data().(*ItemGeneric)
		i.peer.MmoEventer().OnItemGenericEvent(i.peer, message)
		break
	}
}

func (i *InterestArea) SubscribeRegions(newRegions *stl.HashSet) {
	newRegions.ForEach(func(item interface{}) {
		r := item.(*Region)
		if !i.regions.Contains(r) {
			i.regions.Add(r)
			i.OnRegionEnter(r)
			r.RequestItemEnterChannel().Publish(i.requestItemEnterMessage)
		}
	})
	newRegions.Clear()
	newRegions = nil
}

func (i *InterestArea) UnsubscribeRegionsNotIn(regionsToSurvive *stl.HashSet) {
	toUnsubscribeSet := i.regions.Except(regionsToSurvive)
	toUnsubscribeSet.ForEach(func(item interface{}) {
		r := item.(*Region)
		i.regions.Remove(r)
		i.OnRegionExit(r)
		r.RequestItemExitChannel().Publish(i.requestItemExitMessage)
	})
	regionsToSurvive.Clear()
	regionsToSurvive = nil
}

func (i *InterestArea) Dispose() {
	i.subManagementExecutor.Dispose()
	if i.attachedItem != nil {
		i.itemMovementSubscription.Dispose()
		i.itemMovementSubscription = nil
		i.attachedItem = nil
	}
	i.regions.Clear()
	i.regionSubscriptions.ForEach(func(r *Region, s IDisposer) {
		s.Dispose()
	})
	i.regionSubscriptions.Clear()
	i.regionSubscriptions = nil

	i.eventChannelSubscriptions.ForEach(func(r *Region, s IDisposer) {
		s.Dispose()
	})
	i.eventChannelSubscriptions.Clear()
	i.eventChannelSubscriptions = nil
}
