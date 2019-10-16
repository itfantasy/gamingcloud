package gen_lobby

import (
	"github.com/itfantasy/gonode-toolkit/toolkit/gamedb"
	"github.com/itfantasy/gonode/components/mongodb"
)

type LobbyEntity struct {
	LobbyId string `bson:"lobbyid"`
	Nick    string `bson:"nick"`
}

func NewLobbyEntity(lobbyId string) *LobbyEntity {
	l := new(LobbyEntity)
	l.LobbyId = lobbyId
	l.Nick = lobbyId
	return l
}

func (l *LobbyEntity) SetNick(nick string) error {
	l.Nick = nick
	fb := mongodb.NewFilter().Equal("lobbyid", l.LobbyId).Serialize()
	op := mongodb.NewOption().Set("nick", l.Nick).Serialize()
	return gamedb.UpdateLobby(fb, op)
}

func (l *LobbyEntity) RoomCount() (int, error) {
	fb := mongodb.NewFilter().Equal("lobbyid", l.LobbyId).Serialize()
	count, err := gamedb.LobbyCol().CountDocuments(gamedb.Cxt(), fb)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (l *LobbyEntity) PeerCount() (int, error) {
	return gamedb.FindLobbyRoomsPeerCount(l.LobbyId)
}

func (l *LobbyEntity) MasterCount() (int, error) {
	// TODO
	return 0, nil
}

func (l *LobbyEntity) Rooms(startIndex int, endIndex int) ([]*LiteRoomEntity, error) {
	arr := make([]*LiteRoomEntity, 0, endIndex-startIndex+1)
	fb := mongodb.NewFilter().Equal("lobbyid", l.LobbyId).Serialize()
	if err := gamedb.FindRooms(fb, arr, l.LobbyId); err != nil {
		return nil, err
	}
	return arr, nil
}

func (l *LobbyEntity) CreateRoom(roomId string) (*LiteRoomEntity, error) {
	nodeId, err := gamedb.FindBalanceNode(l.LobbyId)
	if err != nil {
		return nil, err
	}
	lr := NewLiteRoomEntity(roomId, l.LobbyId, nodeId)
	if err := gamedb.CreateRoom(lr, l.LobbyId); err != nil {
		return nil, err
	}
	return lr, nil
}

func (l *LobbyEntity) FindRoom(roomId string) (*LiteRoomEntity, error) {
	lr := NewLiteRoomEntity("", "", "")
	fb := mongodb.NewFilter().Equal("roomid", roomId).Serialize()
	if err := gamedb.FindRoom(fb, lr, l.LobbyId); err != nil {
		return nil, err
	}
	return lr, nil
}

func (l *LobbyEntity) RandomRoom() (*LiteRoomEntity, error) {
	lr := NewLiteRoomEntity("", "", "")
	if err := gamedb.FindBalanceRoom(lr, l.LobbyId); err != nil {
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
	filter := mongodb.NewFilter().Equal("lobbyid", lobbyId).Serialize()
	return gamedb.DeleteLobby(filter)
}

func (l *LobbyManager) FindLobby(lobbyId string) (*LobbyEntity, error) {
	lobby := NewLobbyEntity("")
	filter := mongodb.NewFilter().Equal("lobbyid", lobbyId).Serialize()
	if err := gamedb.FindLobby(filter, lobby); err != nil {
		return nil, err
	}
	return lobby, nil
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

func NewLiteRoomEntity(roomId string, lobbyId string, nodeId string) *LiteRoomEntity {
	lr := new(LiteRoomEntity)
	lr.RoomId = roomId
	lr.Nick = roomId
	lr.LobbyId = lobbyId
	lr.NodeId = nodeId
	lr.PeerCount = 0
	lr.MaxPeers = 0
	lr.UsrDatas = make(map[string]string)
	return lr
}

var _lobbyManager *LobbyManager

func lobbyManager() *LobbyManager {
	if _lobbyManager == nil {
		_lobbyManager = NewLobbyManager()
		_lobbyManager.CreateLobby("__default")
	}
	return _lobbyManager
}
