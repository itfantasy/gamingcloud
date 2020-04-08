package gen_mmo

import (
	"sync"

	"github.com/itfantasy/gonode/core/actors"
	"github.com/itfantasy/gonode/utils/snowflake"
)

type MessageChannel struct {
	subscribers map[int64]*MsgChanSubscriber
	sync.RWMutex
}

func NewMessageChannel() *MessageChannel {
	m := new(MessageChannel)
	m.subscribers = make(map[int64]*MsgChanSubscriber)
	return m
}

func (m *MessageChannel) Publish(msg interface{}) bool {
	m.Lock()
	defer m.Unlock()

	if len(m.subscribers) > 0 {
		for _, subscriber := range m.subscribers {
			if subscriber.executor.Living() {
				subscriber.executor.Execute(func() {
					subscriber.action(msg)
				})
			} else {
				return false
			}
		}
	}

	return true
}

func (m *MessageChannel) Subscribe(executor *actors.Executor, action func(interface{})) IDisposer {
	m.Lock()
	defer m.Unlock()

	sid := snowflake.GenerateRaw()
	sub := NewMsgChanSubscriber(sid, action, executor, m)
	m.subscribers[sid] = sub
	return sub
}

func (m *MessageChannel) NumSubscribers() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.subscribers)
}

func (m *MessageChannel) HasSubscriptions() bool {
	return m.NumSubscribers() > 0
}

func (m *MessageChannel) ClearSubscribers() {
	m.Lock()
	defer m.Unlock()

	for _, subscriber := range m.subscribers {
		subscriber.Dispose()
	}
}

func (m *MessageChannel) Dispose() {
	m.ClearSubscribers()
	m = nil
}

type MsgChanSubscriber struct {
	sid      int64
	action   func(interface{})
	executor *actors.Executor
	owner    *MessageChannel
}

func NewMsgChanSubscriber(sid int64, action func(interface{}), executor *actors.Executor, owner *MessageChannel) *MsgChanSubscriber {
	m := new(MsgChanSubscriber)
	m.sid = sid
	m.action = action
	m.executor = executor
	m.owner = owner
	return m
}

func (m *MsgChanSubscriber) Dispose() {
	m.action = nil
	delete(m.owner.subscribers, m.sid)
}
