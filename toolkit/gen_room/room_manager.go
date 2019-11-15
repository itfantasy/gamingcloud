package gen_room

import (
	"errors"

	"github.com/itfantasy/gonode/utils/stl"

	"github.com/itfantasy/gonode-toolkit/toolkit/gamedb"
	"github.com/itfantasy/gonode/components/mongodb"
)

type RoomEntity struct {
	roomId        string
	nick          string
	lobbyId       string
	masterId      int32
	maxPeers      byte
	actorsManager *ActorsManager
	eventCache    *EventCacheManager
}

func NewRoomEntity(roomId string, lobbyId string, maxPeers byte) *RoomEntity {
	r := new(RoomEntity)
	r.roomId = roomId
	r.nick = roomId
	r.lobbyId = lobbyId
	r.masterId = 0
	r.maxPeers = maxPeers
	r.actorsManager = NewActorsManager()
	r.eventCache = NewEventCacheManager()
	return r
}

type LiteRoomEntity struct {
	RoomId    string            `bson:"roomid"`
	Nick      string            `bson:"nick"`
	LobbyId   string            `bson:"lobbyid"`
	NodeId    string            `bson:"nodeid"`
	PeerCount int               `bson:"peercount"`
	MaxPeers  int               `bson:"maxpeers"`
	UsrDatas  map[string]string `bson:"usrdatas"`
}

func NewRoomEntityFromLite(lite *LiteRoomEntity) *RoomEntity {
	return NewRoomEntity(lite.RoomId, lite.LobbyId, byte(lite.MaxPeers))
}

func (r *RoomEntity) RoomId() string {
	return r.roomId
}

func (r *RoomEntity) Nick() string {
	return r.nick
}

func (r *RoomEntity) LobbyId() string {
	return r.lobbyId
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

func (r *RoomEntity) MaxPeers() byte {
	return r.maxPeers
}

func (r *RoomEntity) PeerCount() byte {
	return byte(r.actorsManager.ActorsCount())
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

func (r *RoomEntity) IsFull() bool {
	return r.actorsManager.ActorsCount() >= int(r.MaxPeers())
}

func (r *RoomEntity) UpdateStatusToGameDB() error {
	fb := mongodb.NewFilter().Equal("roomid", r.RoomId()).Serialize()
	op := mongodb.NewOption().Set("nick", r.Nick()).
		Set("peercount", r.PeerCount()).
		Set("maxpeers", r.MaxPeers()).Serialize()
	_, err := gamedb.RoomCol(r.LobbyId()).UpdateOne(gamedb.Cxt(), fb, op)
	if err != nil {
		return err
	}
	return nil
}

type RoomManager struct {
	dict *stl.Dictionary
}

func NewRoomManager() *RoomManager {
	r := new(RoomManager)
	r.dict = stl.NewDictionary()
	return r
}

func (r *RoomManager) CreateRoom(roomId string, lobbyId string, maxPeers byte) (*RoomEntity, error) {
	if r.dict.ContainsKey(roomId) {
		return nil, errors.New("the roommanager has contained a room with the same roomId:" + roomId)
	}
	room := NewRoomEntity(roomId, lobbyId, maxPeers)
	r.dict.Set(roomId, room)
	return room, nil
}

func (r *RoomManager) FindRoom(roomId string) (*RoomEntity, error) {
	item, exist := r.dict.Get(roomId)
	if exist {
		return item.(*RoomEntity), nil
	} else {
		//lite := new(LiteRoomEntity)
		//fb := mongodb.NewFilter().Equal("roomid", roomId).Serialize()
		//if err := gamedb.FindRoom(fb, lite, lobbyId); err != nil {
		//	return nil, errors.New("can not find a room with the roomId:" + roomId)
		//}
		//room := NewRoomEntityFromLite(lite)
		//r.dict.Set(roomId, room)
		//return room
		return nil, errors.New("can not find a room with the roomId:" + roomId)
	}
}

func (r *RoomManager) FetchRoom(roomId string, lobbyId string, maxPeers byte) *RoomEntity {
	item, exist := r.dict.Get(roomId)
	if exist {
		return item.(*RoomEntity)
	} else {
		room := NewRoomEntity(roomId, lobbyId, maxPeers)
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
