package gen_mmo

import (
	"github.com/itfantasy/gonode/core/actors"
)

type Radar struct {
	actionQueue       *actors.Executor
	channel           *MessageChannel
	executor          *actors.Executor
	itemPositions     *DictMmoItemVector
	itemSubscriptions *DictMmoItemIDisposer
}

func NewRadar() *Radar {
	r := new(Radar)
	r.executor = actors.Spawn(1024)
	r.channel = NewMessageChannel()
	r.itemPositions = NewDictMmoItemVector()
	r.itemSubscriptions = NewDictMmoItemIDisposer()
	r.actionQueue = r.executor
	return r
}

func (r *Radar) Channel() *MessageChannel {
	return r.channel
}

func (r *Radar) AddItem(item *MmoItem, position *Vector) {
	r.actionQueue.Execute(func() {
		r.itemPositions.Add(item, position)
		positionUpdates := item.positionUpdateChannel.Subscribe(r.executor, r.UpdatePosition)
		disposeMessage := item.disposeChannel.Subscribe(r.executor, r.RemoveItem)
		r.itemSubscriptions.Add(item, NewUnsubscriberCollection(positionUpdates, disposeMessage))
		r.PublishUpdate(item, position, false)
	})
}

func (r *Radar) SendContentToPeer(peer *MmoPeer) {
	r.actionQueue.Execute(func() {
		r.PublishAll(peer)
	})
}

func (r *Radar) Dispose() {
	r.executor.Dispose()
	r.channel.ClearSubscribers()
	for _, unsubscriber := range r.itemSubscriptions.KeyValuePairs() {
		unsubscriber.Dispose()
	}
	r.itemSubscriptions.Clear()
	r.itemPositions.Clear()
}

func (r *Radar) GetUpdateEvent(item *MmoItem, position *Vector, remove bool) *RadarUpdate {
	return &RadarUpdate{
		ItemId:   item.Id(),
		ItemType: item.Type(),
		Position: position,
		Remove:   remove,
	}
}

func (r *Radar) PublishAll(peer *MmoPeer) {
	kvs := r.itemPositions.KeyValuePairs()
	for item, position := range kvs {
		peer.MmoEventer().OnRadarUpdate(peer, r.GetUpdateEvent(item, position, false))
	}
}

func (r *Radar) PublishUpdate(item *MmoItem, position *Vector, remove bool) {
	updateEvent := r.GetUpdateEvent(item, position, remove)
	message := NewItemEventMessage(item, Event_RadarUpdate, updateEvent)
	r.channel.Publish(message)
}

func (r *Radar) RemoveItem(msg interface{}) {
	message := msg.(*ItemDisposedMessage)
	item := message.Source()
	r.itemPositions.Remove(item)
	if sub, exists := r.itemSubscriptions.Get(item); exists {
		sub.Dispose()
		r.itemSubscriptions.Remove(item)
	}
}

func (r *Radar) UpdatePosition(msg interface{}) {
	message := msg.(*ItemPositionMessage)
	item := message.Source()
	if _, exists := r.itemPositions.Get(item); exists {
		r.itemPositions.Set(item, message.Position())
		r.PublishUpdate(item, message.Position(), false)
	}
}
