package gen_room

import (
	"github.com/itfantasy/gonode/utils/stl"
)

type EventData struct {
	actorNr int32
	code    byte
	data    []byte
}

func NewEventData(actor int32, eventCode byte, data []byte) *EventData {
	e := new(EventData)
	e.actorNr = actor
	e.code = eventCode
	lenth := len(data)
	e.data = make([]byte, lenth, lenth)
	copy(e.data, data)
	return e
}

func (e *EventData) ActorNr() int32 {
	return e.actorNr
}

func (e *EventData) Code() byte {
	return e.code
}

func (e *EventData) Data() []byte {
	return e.data
}

type EventCacheManager struct {
	events *stl.List
}

func NewEventCacheManager() *EventCacheManager {
	e := new(EventCacheManager)
	e.events = stl.NewList(50)
	return e
}

func (e *EventCacheManager) AddEvent(actor int32, eventCode byte, data []byte) {
	e.events.Add(NewEventData(actor, eventCode, data))
}

func (e *EventCacheManager) RemoveEventsByActor(actor int32) int {
	dirtyList := stl.NewList(10)
	for _, item := range e.events.Values() {
		customeEvent := item.(*EventData)
		if customeEvent.actorNr == actor {
			dirtyList.Add(customeEvent)
		}
	}
	for _, item := range dirtyList.Values() {
		e.events.Remove(item)
	}
	return dirtyList.Count()
}

func (e *EventCacheManager) Events() []interface{} {
	return e.events.Values()
}

func (e *EventCacheManager) ClearCache() {
	e.events.Clear()
}
