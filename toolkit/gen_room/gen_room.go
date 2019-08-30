package gen_room

import (
	"errors"
)

type GenRoomCallbacks interface {
	OnJoin(actor *Actor)
	OnLeave(actor *Actor)
	OnDisconn(actor *Actor)
	OnCustomEvent(actor *Actor, data []byte)
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
	return room, actor, nil
}

func JoinRoom(peerId string, roomId string) (*RoomEntity, *Actor, error) {
	room, err := roomManager().GetRoom(roomId)
	if err != nil {
		return nil, nil, err
	}
	actor, err := room.ActorsManager().AddNewActor(peerId)
	if err != nil {
		return nil, nil, err
	}
	return room, actor, nil
}

func LeaveRoom(peerId string, roomId string) error {
	room, err := roomManager().GetRoom(roomId)
	if err != nil {
		return err
	}
	exist := room.ActorsManager().RemoveActorByPeer()
	if !exist {
		return errors.New("can not find the act in the room:" + roomId)
	}
	return nil
}

func DisposeRoom(roomId string) error {
	return nil
}

func RaiseEvent(peerId string, roomId string) error {
	room, err := roomManager().GetRoom(roomId)
	if err != nil {
		return err
	}
	_, exist := room.ActorsManager().GetActorByPeerId(peerId)
	if !exist {
		return errors.New("can not find the act in the room:" + roomId)
	}
	// TODO :

	return nil
}
