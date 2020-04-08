package gen_mmo

type IMmoEventer interface {
	OnItemGenericEvent(peer *MmoPeer, event *ItemGeneric)
	OnItemDestroyed(peer *MmoPeer, itemId string)
	OnItemMoved(peer *MmoPeer, event *ItemMoved)
	OnItemPropertiesSet(peer *MmoPeer, event *ItemPropertiesSet)
	OnWorldExited(peer *MmoPeer, worldName string)
	OnItemSubscribed(peer *MmoPeer, event *ItemSubscribed)
	OnItemUnsubscribed(peer *MmoPeer, event *ItemUnsubscribed)
	OnRadarUpdate(peer *MmoPeer, event *RadarUpdate)
}

const (
	Event_ItemDestroyed byte = iota
	Event_ItemMoved
	Event_ItemPropertiesSet
	Event_ItemGeneric
	Event_RadarUpdate
)

const (
	EventReceiver_ItemSubscriber byte = 1
	EventReceiver_ItemOwner           = 2
)

type ItemGeneric struct {
	ItemId          string
	CustomEventCode byte
	EventData       []byte
}

type ItemDestroyed struct {
	ItemId string
}

type ItemMoved struct {
	ItemId      string
	OldPosition *Vector
	Position    *Vector
	OldRotation *Vector
	Rotation    *Vector
}

type ItemProperties struct {
	ItemId             string
	Source             *MmoItem
	PropertiesRevision int
	PropertiesSet      map[interface{}]interface{}
	Updated            bool
}

type ItemPropertiesSet struct {
	ItemId             string
	PropertiesRevision int
	PropertiesSet      map[interface{}]interface{}
	PropertiesUnset    []interface{}
}

type ItemSubscribed struct {
	InterestAreaId     byte
	ItemId             string
	ItemType           byte
	Position           *Vector
	Rotation           *Vector
	PropertiesRevision int
}

type ItemUnsubscribed struct {
	InterestAreaId byte
	ItemId         string
}

type RadarUpdate struct {
	ItemId   string
	ItemType byte
	Position *Vector
	Remove   bool
}
