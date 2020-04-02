package gen_mmo

// 基础mmo地图同步接口

func CreateWorld(peerId string, worldName string, boundingBox *BoundingBox, tileDimensions *Vector) error {
	if _, exists := getMmoPeer(peerId); !exists {
		return Err_PeerNotFound
	}
	if _, ok := worldCacheIns().TryCreate(worldName, boundingBox, tileDimensions); !ok {
		return Err_WorldAlreadyExists
	}
	return nil
}

func EnterWorld(peerId string, usrName string, worldName string, interestAreaId byte, position *Vector, rotation *Vector, viewDistanceEnter *Vector, viewDistanceExit *Vector, properties map[interface{}]interface{}) error {
	peer, exists := getMmoPeer(peerId)
	if !exists {
		return Err_PeerNotFound
	}
	world, ok := worldCacheIns().TryGet(worldName)
	if !ok {
		return Err_WorldNotFound
	}
	interestArea := NewInterestArea(peer, interestAreaId, world)
	actor := NewMmoActor(peer, world)
	actor.AddInterestArea(interestArea)
	avatar := NewMmoItem(position, rotation, properties, actor, usrName, ItemType_Avatar, world)
	if !world.ItemCatch().AddItem(avatar) {
		if _, ok := world.ItemCatch().GetItem(avatar.Id()); ok {
			avatar.Dispose()
			actor.Dispose()
			interestArea.Dispose()
			ExitWorld(peerId)
			return Err_ItemAlreadyExists
		}
	}
	actor.AddItem(avatar)
	actor.SetAvatar(avatar)
	peer.SetActorHandler(actor)
	interestArea.AttachToItem(avatar)
	interestArea.UpdateInterestManagement()
	avatar.Spawn(position)
	world.Radar().AddItem(avatar, position)
	return nil
}

func ExitWorld(peerId string) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	worldName := peer.MmoActor().World().Name()
	peer.DisposeActor()
	MmoEventCallback(peer, Event_WorldExited, &WorldExited{
		WorldName: worldName,
	})
	return nil
}

func RadarSubscribe(peerId string, worldName string) error {
	peer, exists := getMmoPeer(peerId)
	if !exists {
		return Err_PeerNotFound
	}
	if peer.RadarSubscription() != nil {
		peer.RadarSubscription().Dispose()
		peer.DisposeRadarSubscription()
	}
	world, ok := worldCacheIns().TryGet(worldName)
	if !ok {
		return Err_WorldNotFound
	}
	peer.SetRadarSubscription(world.Radar().Channel().Subscribe(peer.RequestExecutor(), func(m interface{}) {
		message := m.(*ItemEventMessage)
		MmoEventCallback(peer, Event_ItemEvent, message.EventData())
	}))
	world.Radar().SendContentToPeer(peer)
	return nil
}

func doCheckAccess(actor *MmoActor, item *MmoItem) error {
	if item.Disposed() {
		return Err_ItemNotFound
	}
	if !item.GrantWriteAccess(actor) {
		return Err_ItemAccessDenied
	}
	return nil
}

func doCheckPeerAndActor(peerId string) (*MmoPeer, error) {
	peer, exists := getMmoPeer(peerId)
	if !exists {
		return nil, Err_PeerNotFound
	}
	if peer.MmoActor() == nil {
		return nil, Err_InvalidOperation
	}
	return peer, nil
}

func AddInterestArea(peerId string, interestAreaId byte, itemId string, position *Vector, viewDistanceEnter *Vector, viewDistanceExit *Vector) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	if _, ok := peer.MmoActor().TryGetInterestArea(interestAreaId); ok {
		return Err_InterestAreaAlreadyExists
	}
	interestArea := NewInterestArea(peer, interestAreaId, peer.MmoActor().World())
	peer.MmoActor().AddInterestArea(interestArea)
	if itemId != "" {
		item, exists := peer.MmoActor().TryGetItem(itemId)
		if exists {
			return doItemAddInterestArea(item, interestArea, viewDistanceEnter, viewDistanceExit)
		} else {
			newItem, ok := peer.MmoActor().World().ItemCatch().GetItem(itemId)
			if !ok {
				return Err_ItemNotFound
			} else {
				return doItemAddInterestArea(newItem, interestArea, viewDistanceEnter, viewDistanceExit)
			}
		}
	} else {
		// free floating interestArea
		interestArea.Sync(func() {
			interestArea.SetPosition(position)
			interestArea.SetViewDistanceEnter(viewDistanceEnter)
			interestArea.SetViewDistanceExit(viewDistanceExit)
			interestArea.UpdateInterestManagement()
		})
		return nil
	}
}

func doItemAddInterestArea(item *MmoItem, interestArea *InterestArea, viewDistanceEnter *Vector, viewDistanceExit *Vector) error {
	if item.Disposed() {
		return Err_ItemNotFound
	}
	interestArea.Sync(func() {
		interestArea.AttachToItem(item)
		interestArea.SetViewDistanceEnter(viewDistanceEnter)
		interestArea.SetViewDistanceExit(viewDistanceExit)
		interestArea.UpdateInterestManagement()
	})
	return nil
}

func AttachInterestArea(peerId string, interestAreaId byte, itemId string) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	interestArea, ok := peer.MmoActor().TryGetInterestArea(interestAreaId)
	if !ok {
		return Err_InterestAreaNotFound
	}
	var item *MmoItem = nil
	var actorItem bool = false
	if itemId == "" {
		item = peer.MmoActor().Avatar()
		actorItem = true
		itemId = item.Id()
	} else {
		item, actorItem = peer.MmoActor().TryGetItem(itemId)
	}
	if actorItem {
		return doItemAttachInterestArea(item, interestArea)
	} else {
		item, ok := peer.MmoActor().World().ItemCatch().GetItem(itemId)
		if !ok {
			return Err_ItemNotFound
		} else {
			return doItemAttachInterestArea(item, interestArea)
		}
	}
}

func doItemAttachInterestArea(item *MmoItem, interestArea *InterestArea) error {
	if item.Disposed() {
		return Err_ItemNotFound
	}
	interestArea.Sync(func() {
		interestArea.Detach()
		interestArea.AttachToItem(item)
		interestArea.UpdateInterestManagement()
	})
	return nil
}

func DestroyItem(peerId string, itemId string) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	item, actorItem := peer.MmoActor().TryGetItem(itemId)
	if actorItem {
		return doItemDestroy(peer.MmoActor(), item)
	} else {
		item, ok := peer.MmoActor().World().ItemCatch().GetItem(itemId)
		if !ok {
			return Err_ItemNotFound
		} else {
			return doItemDestroy(peer.MmoActor(), item)
		}
	}
}

func doItemDestroy(actor *MmoActor, item *MmoItem) error {
	if err := doCheckAccess(actor, item); err != nil {
		return err
	}
	item.Destroy()
	item.Dispose()
	actor.RemoveItem(item)
	item.World().ItemCatch().RemoveItem(item.Id())
	MmoEventCallback(actor.Peer(), Event_ItemDestroyed, &ItemDestroyed{
		ItemId: item.Id(),
	})
	return nil
}

func DetachInterestArea(peerId string, interestAreaId byte) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	interestArea, ok := peer.MmoActor().TryGetInterestArea(interestAreaId)
	if !ok {
		return Err_InterestAreaNotFound
	}
	interestArea.Sync(func() {
		interestArea.Detach()
	})
	return nil
}

func GetProperties(peerId string, itemId string, propertiesRevision int) (*ItemProperties, error) {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return nil, err
	}
	var item *MmoItem = nil
	item, actorItem := peer.MmoActor().TryGetItem(itemId)
	if !actorItem {
		item, ok := peer.MmoActor().World().ItemCatch().GetItem(itemId)
		if !ok {
			return nil, Err_ItemNotFound
		}
		return doItemGetProperties(item, propertiesRevision)
	}
	return doItemGetProperties(item, propertiesRevision)
}

func doItemGetProperties(item *MmoItem, propertiesRevision int) (*ItemProperties, error) {
	if item.Disposed() {
		return nil, Err_ItemNotFound
	}
	return &ItemProperties{
		ItemId:             item.Id(),
		PropertiesRevision: item.PropertiesRevision(),
		PropertiesSet:      item.Properties().KeyValuePairs(),
		Updated:            item.PropertiesRevision() != propertiesRevision,
	}, nil
}

func Move(peerId string, itemId string, position *Vector, rotation *Vector) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	var item *MmoItem = nil
	if itemId == "" {
		item = peer.MmoActor().Avatar()
		itemId = item.Id()
	} else {
		item, ok := peer.MmoActor().TryGetItem(itemId)
		if !ok {
			return Err_ItemNotFound
		}
		return doItemMove(peer.MmoActor(), item, position, rotation)
	}
	return doItemMove(peer.MmoActor(), item, position, rotation)
}

func doItemMove(actor *MmoActor, item *MmoItem, position *Vector, rotation *Vector) error {
	if err := doCheckAccess(actor, item); err != nil {
		return err
	}
	oldPosition := item.Position()
	oldRotation := item.Rotation()
	item.SetRotation(rotation)
	item.Move(position)
	eventInstance := &ItemMoved{
		ItemId:      item.Id(),
		OldPosition: oldPosition,
		OldRotation: oldRotation,
		Position:    position,
		Rotation:    rotation,
	}
	message := NewItemEventMessage(item, Event_ItemMoved, eventInstance)
	item.EventChannel().Publish(message)
	return nil
}

func MoveInterestArea(peerId string, interestAreaId byte, position *Vector) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	interestArea, ok := peer.MmoActor().TryGetInterestArea(interestAreaId)
	if !ok {
		return Err_InterestAreaNotFound
	}
	interestArea.Sync(func() {
		interestArea.SetPosition(position)
		interestArea.UpdateInterestManagement()
	})
	return nil
}

func RemoveInterestArea(peerId string, interestAreaId byte) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	interestArea, ok := peer.MmoActor().TryGetInterestArea(interestAreaId)
	if !ok {
		return Err_InterestAreaNotFound
	}
	interestArea.Sync(func() {
		interestArea.Detach()
		interestArea.Dispose()
	})
	peer.MmoActor().RemoveInterestArea(interestAreaId)
	return nil
}

func SetProperties(peerId string, itemId string, propertiesSet map[interface{}]interface{}, propertiesUnset []interface{}) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	var item *MmoItem = nil
	if itemId == "" {
		item = peer.MmoActor().Avatar()
		itemId = item.Id()
		return doItemSetProperties(peer.MmoActor(), item, propertiesSet, propertiesUnset)
	} else {
		item, ok := peer.MmoActor().TryGetItem(itemId)
		if !ok {
			return Err_ItemNotFound
		}
		return doItemSetProperties(peer.MmoActor(), item, propertiesSet, propertiesUnset)
	}
}

func doItemSetProperties(actor *MmoActor, item *MmoItem, propertiesSet map[interface{}]interface{}, propertiesUnset []interface{}) error {
	if err := doCheckAccess(actor, item); err != nil {
		return err
	}
	item.SetProperties(propertiesSet, propertiesUnset)
	eventInstance := &ItemPropertiesSet{
		ItemId:             item.Id(),
		PropertiesRevision: item.PropertiesRevision(),
		PropertiesSet:      propertiesSet,
		PropertiesUnset:    propertiesUnset,
	}
	message := NewItemEventMessage(item, Event_ItemPropertiesSet, eventInstance)
	item.EventChannel().Publish(message)
	return nil
}

func SetViewDistance(peerId string, interestAreaId byte, viewDistanceEnter *Vector, viewDistanceExit *Vector) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	interestArea, ok := peer.MmoActor().TryGetInterestArea(interestAreaId)
	if !ok {
		return Err_InterestAreaNotFound
	}
	interestArea.Sync(func() {
		interestArea.SetViewDistanceEnter(viewDistanceEnter)
		interestArea.SetViewDistanceExit(viewDistanceExit)
		interestArea.UpdateInterestManagement()
	})
	return nil
}

func SpawnItem(peerId string, itemId string, itemType byte, position *Vector, rotation *Vector, properties map[interface{}]interface{}) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	item := NewMmoItem(position, rotation, properties, peer.MmoActor(), itemId, itemType, peer.MmoActor().World())
	if !peer.MmoActor().World().ItemCatch().AddItem(item) {
		item.Dispose()
		return Err_ItemAlreadyExists
	}
	peer.MmoActor().AddItem(item)
	return doItemSpawn(peer.MmoActor(), item, position, rotation)
}

func doItemSpawn(actor *MmoActor, item *MmoItem, position *Vector, rotation *Vector) error {
	if err := doCheckAccess(actor, item); err != nil {
		return err
	}
	item.SetRotation(rotation)
	item.Spawn(position)
	actor.World().Radar().AddItem(item, position)
	return nil
}

func SubscribeItem(peerId string, itemId string, propertiesRevision int) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	var item *MmoItem = nil
	item, actorItem := peer.MmoActor().TryGetItem(itemId)
	if !actorItem {
		item, ok := peer.MmoActor().World().ItemCatch().GetItem(itemId)
		if !ok {
			return Err_ItemNotFound
		}
		return doItemSubscribeItem(peer.MmoActor(), item, propertiesRevision)
	}
	return doItemSubscribeItem(peer.MmoActor(), item, propertiesRevision)
}

func doItemSubscribeItem(actor *MmoActor, item *MmoItem, propertiesRevision int) error {
	if item.Disposed() {
		return Err_ItemNotFound
	}
	actor.InterestItems().SubscribeItem(item)
	subscribeEvent := &ItemSubscribed{
		ItemId:             item.Id(),
		ItemType:           item.Type(),
		Position:           item.Position(),
		Rotation:           item.Rotation(),
		PropertiesRevision: item.PropertiesRevision(),
	}
	MmoEventCallback(actor.Peer(), Event_ItemSubscribed, subscribeEvent)
	if propertiesRevision != item.PropertiesRevision() {
		properties := &ItemPropertiesSet{
			ItemId:             item.Id(),
			PropertiesRevision: propertiesRevision,
			PropertiesSet:      item.Properties().KeyValuePairs(),
		}
		MmoEventCallback(actor.Peer(), Event_ItemPropertiesSet, properties)
	}
	return nil
}

func UnsubscribeItem(peerId string, itemId string) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	var item *MmoItem = nil
	item, actorItem := peer.MmoActor().TryGetItem(itemId)
	if !actorItem {
		item, ok := peer.MmoActor().World().ItemCatch().GetItem(itemId)
		if !ok {
			return Err_ItemNotFound
		}
		return doItemUnsubscribeItem(peer.MmoActor(), item)
	}
	return doItemUnsubscribeItem(peer.MmoActor(), item)
}

func doItemUnsubscribeItem(actor *MmoActor, item *MmoItem) error {
	if item.Disposed() {
		return Err_ItemNotFound
	}
	actor.InterestItems().UnsubscribeItem(item)
	MmoEventCallback(actor.Peer(), Event_ItemUnsubscribed, &ItemSubscribed{
		ItemId: item.Id(),
	})
	return nil
}

func RaiseGenericEvent(peerId string, itemId string, customEventCode byte, eventData interface{}, eventReceiver byte) error {
	peer, err := doCheckPeerAndActor(peerId)
	if err != nil {
		return err
	}
	var item *MmoItem = nil
	item, actorItem := peer.MmoActor().TryGetItem(itemId)
	if !actorItem {
		item, ok := peer.MmoActor().World().ItemCatch().GetItem(itemId)
		if !ok {
			return Err_ItemNotFound
		}
		return doItemRaiseGenericEvent(item, customEventCode, eventData, eventReceiver)
	}
	return doItemRaiseGenericEvent(item, customEventCode, eventData, eventReceiver)
}

func doItemRaiseGenericEvent(item *MmoItem, customEventCode byte, eventData interface{}, eventReceiver byte) error {
	if item.Disposed() {
		return Err_ItemNotFound
	}
	eventInstance := &ItemGeneric{
		ItemId:          item.Id(),
		CustomEventCode: customEventCode,
		EventData:       eventData,
	}
	switch eventReceiver {
	case EventReceiver_ItemOwner:
		MmoEventCallback(item.Owner().Peer(), Event_ItemGeneric, eventInstance)
	case EventReceiver_ItemSubscriber:
		message := NewItemEventMessage(item, Event_ItemGeneric, eventInstance)
		item.EventChannel().Publish(message)
	}
	return nil
}

func AddPeer(peer *MmoPeer) error {
	return peerManager().AddPeer(peer)
}

func RemovePeer(peerId string) error {
	return peerManager().RemovePeer(peerId)
}

func GetPeer(peerId string) (*MmoPeer, bool) {
	return getMmoPeer(peerId)
}
