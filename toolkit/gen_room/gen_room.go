package gen_room

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itfantasy/gonode"
	"github.com/itfantasy/gonode-toolkit/toolkit/gamedb"
	"github.com/itfantasy/gonode/components/mongodb"
)

type IGenRoomEventer interface {
	OnJoinRoom(actor *Actor, room *RoomEntity)
	OnLeaveRoom(actor *Actor, room *RoomEntity)
	OnCustomEvent(actor *Actor, room *RoomEntity, data []byte)
}

var _eventer IGenRoomEventer

func BindEventCallback(eventer IGenRoomEventer) {
	_eventer = eventer
}

func InitGameDB(mongoConf string) error {
	if err := gamedb.InitMongo(mongoConf); err != nil {
		return err
	}
	return nil
}

func CreateRoom(peerId string, roomId string, lobbyId string, maxPeers byte) (*RoomEntity, *Actor, error) {
	if peerId == "" || roomId == "" || lobbyId == "" || maxPeers <= 0 {
		return nil, nil, errors.New("illegal args for room creating!!")
	}
	if !strings.HasSuffix(roomId, lobbyId) {
		roomId = roomId + "@" + lobbyId
	}
	room, err := roomManager().CreateRoom(roomId, lobbyId, maxPeers)
	if err != nil {
		return nil, nil, err
	}
	actor, err := room.ActorsManager().AddNewActor(peerId)
	if err != nil {
		return nil, nil, err
	}
	room.SetMasterId(actor.ActorNr())
	if err := room.UpdateStatusToGameDB(); err != nil {
		room.ActorsManager().RemoveActorByPeer(peerId)
		room.SetMasterId(0)
		return nil, nil, err
	}
	if _eventer != nil {
		_eventer.OnJoinRoom(actor, room)
	}
	return room, actor, nil
}

func JoinRoom(peerId string, roomId string) (*RoomEntity, *Actor, error) {
	if peerId == "" || roomId == "" {
		return nil, nil, errors.New("illegal args for room joining!!")
	}
	room, err := roomManager().FindRoom(roomId)
	if err != nil {
		room, err = reuseLiteRoom(roomId)
		if err != nil {
			return nil, nil, err
		}
	}
	if room.IsFull() {
		return nil, nil, errors.New("the room is full!!" + roomId)
	}
	actor, err := room.ActorsManager().AddNewActor(peerId)
	if err != nil {
		return nil, nil, err
	}
	if room.IsEmpty() {
		room.SetMasterId(actor.ActorNr())
	}
	if err := room.UpdateStatusToGameDB(); err != nil {
		room.ActorsManager().RemoveActorByPeer(peerId)
		return nil, nil, err
	}
	if _eventer != nil {
		_eventer.OnJoinRoom(actor, room)
	}
	return room, actor, nil
}

func reuseLiteRoom(roomId string) (*RoomEntity, error) {
	infos := strings.Split(roomId, "@")
	if len(infos) < 2 {
		return nil, errors.New("illegal roomId:" + roomId)
	}
	lobbyId := infos[1]
	lite := new(LiteRoomEntity)
	fb := mongodb.NewFilter().Equal("roomid", roomId).Serialize()
	if err := gamedb.FindRoom(fb, lite, lobbyId); err != nil {
		return nil, errors.New("can not find a room with the roomId:" + roomId)
	}
	room := NewRoomEntityFromLite(lite)
	return room, nil
}

func LeaveRoom(peerId string, roomId string) (*RoomEntity, *Actor, error) {
	if peerId == "" || roomId == "" {
		return nil, nil, errors.New("illegal args for room leaving!!")
	}
	room, err := roomManager().FindRoom(roomId)
	if err != nil {
		return nil, nil, err
	}
	actor, exist := room.ActorsManager().RemoveActorByPeer(peerId)
	if !exist {
		return nil, nil, errors.New("can not find the act in the room:" + roomId)
	}
	room.EventCache().RemoveEventsByActor(actor.ActorNr())
	if room.MasterId() == actor.actorNr {
		if room.IsEmpty() {
			room.SetMasterId(0)
		} else {
			newMaster, _ := room.ActorsManager().GetActorByIndex(0)
			room.SetMasterId(newMaster.ActorNr())
		}
	}
	if err := room.UpdateStatusToGameDB(); err != nil {
		return nil, nil, err
	}
	if _eventer != nil {
		_eventer.OnLeaveRoom(actor, room)
	}
	return room, actor, nil
}

func GetActorInRoom(peerId string, roomId string) (*RoomEntity, *Actor, error) {
	if peerId == "" || roomId == "" {
		return nil, nil, errors.New("illegal args for getting room actor!!")
	}
	room, err := roomManager().FindRoom(roomId)
	if err != nil {
		return nil, nil, err
	}
	actor, exist := room.ActorsManager().GetActorByPeerId(peerId)
	if !exist {
		return nil, nil, errors.New("can not find the actor [ " + peerId + "] in the room:" + roomId)
	}
	return room, actor, nil
}

func DisposeRoom(roomId string) error {
	if roomId == "" {
		return errors.New("illegal args for room disposing!!")
	}
	room, err := roomManager().FindRoom(roomId)
	if err != nil {
		return err
	}
	peers := room.ActorsManager().GetAllPeerIds()
	for _, peerId := range peers {
		gonode.Close(peerId)
	}
	gonode.LogWarn("room has been closed!!" + roomId)
	return nil
}

const (
	RcvGroup_Others byte = 0
	RcvGroup_All         = 1
	RcvGroup_Master      = 2
)

func RaiseEvent(peerId string, roomId string, data []byte, rcvGroup byte, addToRoomCache bool) error {
	room, actor, err := GetActorInRoom(peerId, roomId)
	if err != nil {
		return err
	}
	if addToRoomCache {
		room.EventCache().AddEvent(actor.ActorNr(), data)
	}
	var theErr error = nil
	switch rcvGroup {
	case RcvGroup_All:
		peerIds := room.ActorsManager().GetAllPeerIds()
		for _, otherId := range peerIds {
			if err := gonode.Send(otherId, data); err != nil {
				theErr = err
			}
		}
	case RcvGroup_Others:
		peerIds := room.ActorsManager().GetAllPeerIds()
		for _, otherId := range peerIds {
			if otherId != peerId {
				if err := gonode.Send(otherId, data); err != nil {
					theErr = err
				}
			}
		}
	case RcvGroup_Master:
		masterId := room.MasterId()
		if masterId <= 0 {
			theErr = errors.New("the room has no master!" + roomId)
		} else {
			master, exist := room.ActorsManager().GetActorByNr(masterId)
			if !exist {
				theErr = errors.New("the room has no master with the masterId! " + fmt.Sprint(masterId))
			} else {
				if err := gonode.Send(master.PeerId(), data); err != nil {
					theErr = err
				}
			}
		}
	default:
		theErr = errors.New("unkown RcvGroup type!!")
	}
	if _eventer != nil {
		_eventer.OnCustomEvent(actor, room, data)
	}
	return theErr
}

func RcvCacheEvent(peerId string, roomId string) error {
	room, _, err := GetActorInRoom(peerId, roomId)
	if err != nil {
		return err
	}
	events := room.EventCache().Events()
	for _, item := range events {
		event := item.(*EventData)
		gonode.Send(peerId, event.Data())
	}
	return nil
}

func ClrEventCache(roomId string, peerId string) error {
	room, actor, err := GetActorInRoom(peerId, roomId)
	if err != nil {
		return err
	}
	if peerId == "" {
		room.EventCache().ClearCache()
	} else {
		room.EventCache().RemoveEventsByActor(actor.ActorNr())
	}
	return nil
}

func AddPeer(peer *RoomPeer) error {
	return peerManager().AddPeer(peer)
}

func RemovePeer(peerId string) error {
	return peerManager().RemovePeer(peerId)
}

func GetPeer(peerId string) (*RoomPeer, bool) {
	return getRoomPeer(peerId)
}
