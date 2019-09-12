package gen_room

import (
	"errors"
	"fmt"

	"github.com/itfantasy/gonode"
)

type GenRoomCallbacks interface {
	OnJoin(actor *Actor, room *RoomEntity)
	OnLeave(actor *Actor, room *RoomEntity)
	OnCustomEvent(actor *Actor, room *RoomEntity, data []byte)
}

var _callbacks GenRoomCallbacks

func BindCallbacks(callbacks GenRoomCallbacks) {
	_callbacks = callbacks
}

func CreateRoom(peerId string, roomId string) (*RoomEntity, *Actor, error) {
	room, err := roomManager().CreateRoom(roomId)
	if err != nil {
		return nil, nil, err
	}
	actor, err := room.ActorsManager().AddNewActor(peerId)
	if err != nil {
		return nil, nil, err
	}
	room.SetMasterId(actor.ActorNr())
	if _callbacks != nil {
		_callbacks.OnJoin(actor, room)
	}
	return room, actor, nil
}

func JoinRoom(peerId string, roomId string) (*RoomEntity, *Actor, error) {
	//room, err := roomManager().GetRoom(roomId)
	//if err != nil {
	//	return nil, nil, err
	//}
	room := roomManager().FetchRoom(roomId)
	room.SetMasterId(1)

	actor, err := room.ActorsManager().AddNewActor(peerId)
	if err != nil {
		return nil, nil, err
	}
	if _callbacks != nil {
		_callbacks.OnJoin(actor, room)
	}
	return room, actor, nil
}

func LeaveRoom(peerId string, roomId string) (*RoomEntity, *Actor, error) {
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
	if _callbacks != nil {
		_callbacks.OnLeave(actor, room)
	}
	return room, actor, nil
}

func GetActorInRoom(peerId string, roomId string) (*RoomEntity, *Actor, error) {
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
	if _callbacks != nil {
		_callbacks.OnCustomEvent(actor, room, data)
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

func AddPeer(peer Peer) error {
	return peerManager().AddPeer(peer)
}

func RemovePeer(peerId string) error {
	return peerManager().RemovePeer(peerId)
}

func GetPeer(peerId string) (Peer, bool) {
	return peerManager().GetPeer(peerId)
}

func GetClientPeer(peerId string) (*ClientPeer, bool) {
	return peerManager().GetClientPeer(peerId)
}
