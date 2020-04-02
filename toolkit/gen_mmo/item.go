package gen_mmo

import (
	"github.com/itfantasy/gonode/core/actors"
	"github.com/itfantasy/gonode/utils/stl"
)

const (
	ItemType_Avatar byte = iota
	ItemType_Bot
)

type MmoItem struct {
	id                    string
	eventChannel          *MessageChannel
	positionUpdateChannel *MessageChannel
	disposeChannel        *MessageChannel
	properties            *stl.HashTable
	itype                 byte
	world                 *World
	disposed              bool
	regionSubscription    IDisposer

	currentWorldRegion *Region
	rotation           *Vector
	position           *Vector
	propertiesRevision int
	owner              *MmoActor
}

func NewMmoItem(position *Vector, rotation *Vector, properties map[interface{}]interface{}, owner *MmoActor, id string, itype byte, world *World) *MmoItem {
	i := new(MmoItem)
	i.position = position
	i.rotation = rotation
	i.owner = owner
	i.eventChannel = NewMessageChannel()
	i.disposeChannel = NewMessageChannel()
	i.positionUpdateChannel = NewMessageChannel()
	if properties == nil {
		i.properties = stl.NewHashTable()
		i.propertiesRevision = 0
	} else {
		i.properties = stl.NewHashTableRaw(properties)
		i.propertiesRevision++
	}
	i.id = id
	i.world = world
	i.itype = itype
	return i
}

func (i *MmoItem) Id() string {
	return i.id
}

func (i *MmoItem) CurrentWorldRegion() *Region {
	return i.currentWorldRegion
}

func (i *MmoItem) DisposeChannel() *MessageChannel {
	return i.disposeChannel
}

func (i *MmoItem) Disposed() bool {
	return i.disposed
}

func (i *MmoItem) EventChannel() *MessageChannel {
	return i.eventChannel
}

func (i *MmoItem) Excutor() *actors.Executor {
	return i.Owner().Peer().RequestExecutor()
}

func (i *MmoItem) Owner() *MmoActor {
	return i.owner
}

func (i *MmoItem) Rotation() *Vector {
	return i.rotation
}

func (i *MmoItem) SetRotation(rotation *Vector) {
	i.rotation = rotation
}

func (i *MmoItem) Position() *Vector {
	return i.position
}

func (i *MmoItem) SetPosition(position *Vector) {
	i.position = position
}

func (i *MmoItem) PositionUpdateChannel() *MessageChannel {
	return i.positionUpdateChannel
}

func (i *MmoItem) Properties() *stl.HashTable {
	return i.properties
}

func (i *MmoItem) PropertiesRevision() int {
	return i.propertiesRevision
}

func (i *MmoItem) Type() byte {
	return i.itype
}

func (i *MmoItem) World() *World {
	return i.world
}

func (i *MmoItem) Destroy() {
	i.OnDestroy()
}

func (i *MmoItem) UpdateInterestManagement() {
	message := i.GetPositionUpdateMessage(i.position)
	i.positionUpdateChannel.Publish(message)

	prevRegion := i.currentWorldRegion
	newRegion, _ := i.world.GetRegion(i.position)

	if newRegion != i.currentWorldRegion {
		i.currentWorldRegion = newRegion
		if i.regionSubscription != nil {
			i.regionSubscription.Dispose()
		}
		snapshot := i.GetItemSnapshot()
		regMessage := NewItemRegionChangedMessage(prevRegion, newRegion, snapshot)
		if prevRegion != nil {
			prevRegion.ItemRegionChangedChannel().Publish(regMessage)
		}
		if newRegion != nil {
			newRegion.ItemRegionChangedChannel().Publish(regMessage)

			i.regionSubscription = NewUnsubscriberCollection(
				i.eventChannel.Subscribe(i.Excutor(), func(m interface{}) {
					newRegion.ItemEventChannel().Publish(m)
				}),
				newRegion.RequestItemEnterChannel().Subscribe(i.Excutor(), func(m interface{}) {
					msg := m.(*RequestItemEnterMessage)
					msg.InterestArea().OnItemEnter(i.GetItemSnapshot())
				}),
				newRegion.RequestItemExitChannel().Subscribe(i.Excutor(), func(m interface{}) {
					msg := m.(*RequestItemExitMessage)
					msg.InterestArea().OnItemExit(i)
				}),
			)
		}

	}
}

func (i *MmoItem) SetProperties(propertiesSet map[interface{}]interface{}, propertiesUnset []interface{}) {
	if propertiesSet != nil {
		for k, v := range propertiesSet {
			i.properties.Set(k, v)
		}
	}
	if propertiesUnset != nil {
		for _, item := range propertiesUnset {
			i.properties.Remove(item)
		}
	}
	i.propertiesRevision++
}

func (i *MmoItem) GetItemSnapshot() *ItemSnapshot {
	return NewItemSnapshot(i, i.position, i.rotation, i.currentWorldRegion, i.propertiesRevision)
}

func (i *MmoItem) GetPositionUpdateMessage(position *Vector) *ItemPositionMessage {
	return NewItemPositionMessage(i, position)
}

func (i *MmoItem) Dispose() {
	if i.regionSubscription != nil {
		i.regionSubscription.Dispose()
	}
	i.currentWorldRegion = nil
	i.disposeChannel.Publish(NewItemDisposedMessage(i))
	i.eventChannel.Dispose()
	i.disposeChannel.Dispose()
	i.positionUpdateChannel.Dispose()
	i.disposed = true
}

func (i *MmoItem) OnDestroy() {
	eventInstance := &ItemDestroyed{
		ItemId: i.Id(),
	}
	message := NewItemEventMessage(i, Event_ItemDestroyed, eventInstance)
	i.eventChannel.Publish(message)
}

func (i *MmoItem) Move(position *Vector) {
	i.position = position
	i.UpdateInterestManagement()
}

func (i *MmoItem) Spawn(position *Vector) {
	i.position = position
	i.UpdateInterestManagement()
}

func (i *MmoItem) GrantWriteAccess(actor *MmoActor) bool {
	return i.owner == actor
}
