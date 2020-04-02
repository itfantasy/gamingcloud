package gen_mmo

type IMmoEventer interface {
	OnMmoEvent(*MmoPeer, byte, interface{})
}

var _eventer IMmoEventer

func BindEventCallback(eventer IMmoEventer) {
	_eventer = eventer
}

func MmoEventCallback(peer *MmoPeer, evnCode byte, evnData interface{}) {
	if _eventer != nil {
		_eventer.OnMmoEvent(peer, evnCode, evnData)
	}
}

const (
	Event_ItemEvent byte = iota
	Event_ItemDestroyed
	Event_ItemMoved
	Event_ItemPropertiesSet
	Event_WorldExited
	Event_ItemSubscribed
	Event_ItemUnsubscribed
	Event_ItemProperties
	Event_RadarUpdate
	Event_CounterData
	Event_ItemGeneric
)

const (
	EventReceiver_ItemSubscriber byte = 1
	EventReceiver_ItemOwner           = 2
)

type ItemDestroyed struct {
	ItemId string
}

type ItemGeneric struct {
	CustomEventCode byte
	EventData       interface{}
	ItemId          string
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

type WorldExited struct {
	WorldName string
}

type EventData struct {
	Code  byte
	Datas interface{}
}
