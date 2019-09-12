package gen_room

import (
	"errors"

	"github.com/itfantasy/gonode/utils/stl"
)

type RoomEntity struct {
	roomId        string
	nick          string
	masterId      int32
	actorsManager *ActorsManager
	eventCache    *EventCacheManager
}

func NewRoomEntity(roomId string) *RoomEntity {
	r := new(RoomEntity)
	r.roomId = roomId
	r.nick = roomId
	r.masterId = 0
	r.actorsManager = NewActorsManager()
	r.eventCache = NewEventCacheManager()
	return r
}

func (r *RoomEntity) RoomId() string {
	return r.roomId
}

func (r *RoomEntity) Nick() string {
	return r.nick
}

func (r *RoomEntity) SetNick(nick string) {
	r.nick = nick
}

func (r *RoomEntity) MasterId() int32 {
	return r.masterId
}

func (r *RoomEntity) SetMasterId(masterId int32) {
	r.masterId = masterId
}

func (r *RoomEntity) ActorsManager() *ActorsManager {
	return r.actorsManager
}

func (r *RoomEntity) EventCache() *EventCacheManager {
	return r.eventCache
}

func (r *RoomEntity) IsEmpty() bool {
	return r.actorsManager.ActorsCount() <= 0
}

type RoomManager struct {
	dict *stl.Dictionary
}

func NewRoomManager() *RoomManager {
	r := new(RoomManager)
	r.dict = stl.NewDictionary()
	return r
}

func (r *RoomManager) CreateRoom(roomId string) (*RoomEntity, error) {
	if r.dict.ContainsKey(roomId) {
		return nil, errors.New("the roommanager has contained a room with the same roomId:" + roomId)
	}
	room := NewRoomEntity(roomId)
	r.dict.Set(roomId, room)
	return room, nil
}

func (r *RoomManager) FindRoom(roomId string) (*RoomEntity, error) {
	item, exist := r.dict.Get(roomId)
	if exist {
		return item.(*RoomEntity), nil
	} else {
		return nil, errors.New("can not find a room with the roomId:" + roomId)
	}
}

func (r *RoomManager) FetchRoom(roomId string) *RoomEntity {
	item, exist := r.dict.Get(roomId)
	if exist {
		return item.(*RoomEntity)
	} else {
		room := NewRoomEntity(roomId)
		r.dict.Set(roomId, room)
		return room
	}
}

func (r *RoomManager) DisposeRoom(roomId string) {
	// need dispose the actorsManager and the eventCache
	for _, val := range r.dict.KeyValuePairs() {
		room := val.(*RoomEntity)
		room.actorsManager.ClearAll()
		room.eventCache.ClearCache()
	}
}

var _roomManager *RoomManager

func init() {
	_roomManager = NewRoomManager()
}

func roomManager() *RoomManager {
	return _roomManager
}
