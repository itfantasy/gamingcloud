package gen_lobby

import (
	"github.com/itfantasy/gonode-toolkit/toolkit/gamedb"
	"github.com/itfantasy/gonode/components/mongodb"
)

type LobbyEntity struct {
	LobbyId   string `bson:"lobbyid"`
	Nick      string `bson:"nick"`
	RoomCount int
	PeerCount int
	Rooms     []LiteRoomEntity
}

func NewLobbyEntity(lobbyId string) *LobbyEntity {
	l := new(LobbyEntity)
	l.LobbyId = lobbyId
	l.Nick = lobbyId
	l.RoomCount = 0
	l.PeerCount = 0
	l.Rooms = make([]LiteRoomEntity, 0, 10)
	return l
}

func (l *LobbyEntity) SetNick(nick string) {
	l.Nick = nick
	fb := mongodb.NewFilterBuilder().Equal("lobbyid", l.LobbyId).Serialize()
	op := mongodb.NewOptionBuilder().Set("nick", l.Nick).Serialize()
	gamedb.UpdateLobby(fb, op)
}

func (l *LobbyEntity) CreateRoom(roomId string) (*LiteRoomEntity, error) {
	lr := NewLiteRoomEntity(roomId, l.LobbyId)
	if err := gamedb.CreateRoom(lr); err != nil {
		return nil, err
	}
	return lr, nil
}

type LobbyManager struct {
}

func NewLobbyManager() *LobbyManager {
	l := new(LobbyManager)
	return l
}

func (l *LobbyManager) CreateLobby(lobbyId string) (*LobbyEntity, error) {
	lobby := NewLobbyEntity(lobbyId)
	if err := gamedb.CreateLobby(lobby); err != nil {
		return nil, err
	}
	return lobby, nil
}

func (l *LobbyManager) DisposeLobby(lobbyId string) error {
	filter := mongodb.NewFilterBuilder().Equal("lobbyid", lobbyId).Serialize()
	if err := gamedb.DeleteRooms(filter); err != nil {
		return err
	}
	return gamedb.DeleteLobby(filter)
}

func (l *LobbyManager) FindLobby(lobbyId string) (*LobbyEntity, error) {
	lobby := NewLobbyEntity("")
	filter := mongodb.NewFilterBuilder().Equal("lobbyid", lobbyId).Serialize()
	if err := gamedb.FindLobby(filter, lobby); err != nil {
		return nil, err
	}
	rooms := make([]LiteRoomEntity, 0, 10)
	if err := gamedb.FindRooms(filter, rooms); err != nil {
		return nil, err
	}
	lobby.Rooms = rooms
	for _, room := range rooms {
		lobby.RoomCount = lobby.RoomCount + 1
		lobby.PeerCount = lobby.PeerCount + room.PeerCount
	}
	return lobby, nil
}

type LiteRoomEntity struct {
	RoomId    string `bson:"roomid"`
	Nick      string `bson:"nick"`
	LobbyId   string `bson:"lobbyid"`
	NodeId    string `bson:"nodeid"`
	PeerCount int    `bson:"peercount"`
}

func NewLiteRoomEntity(roomId string, lobbyId string) *LiteRoomEntity {
	lr := new(LiteRoomEntity)
	lr.RoomId = roomId
	lr.Nick = roomId
	lr.LobbyId = lobbyId
	lr.PeerCount = 0
	return lr
}

func (lr *LiteRoomEntity) SetNick(nick string) {
	lr.Nick = nick
}

var _lobbyManager *LobbyManager

func lobbyManager() *LobbyManager {
	if _lobbyManager == nil {
		_lobbyManager = NewLobbyManager()
		_lobbyManager.CreateLobby("__default")
	}
	return _lobbyManager
}
