package gen_mmo

import (
	"github.com/itfantasy/gonode/core/actors"
)

type InterestItems struct {
	peer                    *MmoPeer
	manualItemSubscriptions *DictMmoItemIDisposer
	itemEventExecutor       *actors.Executor
	subManagementExecutor   *actors.Executor
}

func NewInterestItems(peer *MmoPeer) *InterestItems {
	i := new(InterestItems)
	i.peer = peer
	i.manualItemSubscriptions = NewDictMmoItemIDisposer()
	i.itemEventExecutor = actors.Spawn(1024)
	i.subManagementExecutor = actors.Spawn(1024)
	return i
}

func (i *InterestItems) SubscribeItem(item *MmoItem) bool {
	if i.manualItemSubscriptions.ContainsKey(item) {
		return false
	}
	messagesListener := item.EventChannel().Subscribe(i.itemEventExecutor, i.SubscribedItem_OnItemEvent)
	managementListener := item.DisposeChannel().Subscribe(i.subManagementExecutor, i.SubscribedItem_OnItemDisposed)
	i.manualItemSubscriptions.Add(item, NewUnsubscriberCollection(messagesListener, managementListener))
	return true
}

func (i *InterestItems) UnsubscribeItem(item *MmoItem) bool {
	if subscription, exists := i.manualItemSubscriptions.Get(item); exists {
		subscription.Dispose()
		i.manualItemSubscriptions.Remove(item)
		return true
	}
	return false
}

func (i *InterestItems) SubscribedItem_OnItemDisposed(msg interface{}) {
	itemDisposeMessage := msg.(ItemDisposedMessage)
	i.UnsubscribeItem(itemDisposeMessage.Source())
}

func (i *InterestItems) SubscribedItem_OnItemEvent(msg interface{}) {
	m := msg.(ItemEventMessage)
	MmoEventCallback(i.peer, Event_ItemEvent, m)
}

func (i *InterestItems) ClearManualSubscriptions() {
	i.manualItemSubscriptions.ForEach(func(k *MmoItem, v IDisposer) {
		v.Dispose()
	})
	i.manualItemSubscriptions.Clear()
}

func (i *InterestItems) Dispose() {
	i.itemEventExecutor.Dispose()
	i.subManagementExecutor.Dispose()
	i.ClearManualSubscriptions()
}
