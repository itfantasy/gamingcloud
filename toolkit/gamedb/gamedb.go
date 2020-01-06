package gamedb

import (
	"errors"

	"golang.org/x/net/context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/itfantasy/gonode/components"
	"github.com/itfantasy/gonode/components/mongodb"

	"github.com/itfantasy/gonode"

	"github.com/itfantasy/gonode-toolkit/toolkit"
)

const (
	LOBBY_COLLECT string = "__lobbys"
	ROOM_COLLECT         = "__rooms"
)

var _context context.Context
var _mongo *mongodb.MongoDB

func Cxt() context.Context {
	return _context
}

func InitMongo(mongoConf string) error {
	comp, err := components.NewComponent(mongoConf)
	if err != nil {
		return err
	}
	mongoComp, ok := comp.(*mongodb.MongoDB)
	if !ok {
		return errors.New("the gamedb only support mongodb compontent to init itself!")
	}
	_mongo = mongoComp
	_context = context.Background()
	return nil
}

func LobbyCol() *mongo.Collection {
	return _mongo.Collect(LOBBY_COLLECT)
}

func RoomCol(lobbyid string) *mongo.Collection {
	return _mongo.Collect(ROOM_COLLECT + "@" + lobbyid)
}

func CreateLobby(entity interface{}) error {
	if _, err := LobbyCol().InsertOne(_context, entity); err != nil {
		return err
	}
	return nil
}

func UpdateLobby(filter map[string]interface{}, data map[string]map[string]interface{}) error {
	if _, err := LobbyCol().UpdateOne(_context, filter, data); err != nil {
		return err
	}
	return nil
}

func DeleteLobby(filter map[string]interface{}) error {
	if _, err := LobbyCol().DeleteOne(_context, filter); err != nil {
		return err
	}
	return nil
}

func FindLobby(filter map[string]interface{}, entity interface{}) error {
	ret := LobbyCol().FindOne(_context, filter)
	if ret.Err() != nil {
		return ret.Err()
	}
	if err := ret.Decode(entity); err != nil {
		return err
	}
	return nil
}

func FindBalanceNode(lobbyid string) (string, error) {
	gb := mongodb.NewGroupBy().Sum("peercount", "nodeid", "num").Serialize()
	ret, err := RoomCol(lobbyid).Aggregate(_context, gb)
	if err != nil {
		return "", err
	}
	arr := make([]BalanceResult, 0, 3)
	if err := ret.All(_context, &arr); err != nil {
		return "", err
	}
	_map := make(map[string]int)
	roomNodes := gonode.Nodes(toolkit.LABEL_ROOM)
	for _, node := range roomNodes {
		_map[node] = 0
	}
	for _, result := range arr {
		_map[result.Id] = result.Num
	}
	var minnode = ""
	var minnum = 999999
	for id, num := range _map {
		if num < minnum {
			minnode = id
			minnum = num
			if minnum <= 0 {
				break
			}
		}
	}
	return minnode, nil
}

func CreateRoom(entity interface{}, lobbyid string) error {
	if _, err := RoomCol(lobbyid).InsertOne(_context, entity); err != nil {
		return err
	}
	return nil
}

func UpdateRoom(entity interface{}, filter map[string]interface{}, lobbyid string) error {
	if _, err := RoomCol(lobbyid).UpdateOne(_context, filter, entity); err != nil {
		return err
	}
	return nil
}

func DeleteRoom(filter map[string]interface{}, lobbyid string) error {
	if _, err := RoomCol(lobbyid).DeleteOne(_context, filter); err != nil {
		return err
	}
	return nil
}

func DeleteRooms(filter map[string]interface{}, lobbyid string) error {
	if _, err := RoomCol(lobbyid).DeleteMany(_context, filter); err != nil {
		return err
	}
	return nil
}

func FindRoom(filter map[string]interface{}, entity interface{}, lobbyid string) error {
	ret := RoomCol(lobbyid).FindOne(_context, filter)
	if ret.Err() != nil {
		return ret.Err()
	}
	if err := ret.Decode(entity); err != nil {
		return err
	}
	return nil
}

func FindRooms(filter map[string]interface{}, entities interface{}, lobbyid string) error {
	cursor, err := RoomCol(lobbyid).Find(_context, filter)
	if err != nil {
		return err
	}
	if cursor.Err() != nil {
		return cursor.Err()
	}
	if err := cursor.All(_context, &entities); err != nil {
		return err
	}
	return nil
}

func FindLobbyRoomsPeerCount(lobbyid string) (int, error) {
	gb := mongodb.NewGroupBy().Min("peercount", "lobbyid", "num").Serialize()
	ret, err := RoomCol(lobbyid).Aggregate(_context, gb)
	if err != nil {
		return 0, err
	}
	arr := make([]BalanceResult, 0, 3)
	if err := ret.All(_context, &arr); err != nil {
		return 0, err
	}
	for _, result := range arr {
		if result.Id == lobbyid {
			return result.Num, nil
		}
	}
	return 0, errors.New("cannot find the lobby:" + lobbyid)
}

func FindBalanceRoom(entity interface{}, lobbyid string) error {
	gb := mongodb.NewGroupBy().Min("peercount", "nodeid", "num").Serialize()
	ret, err := RoomCol(lobbyid).Aggregate(_context, gb)
	if err != nil {
		return err
	}
	arr := make([]BalanceResult, 0, 3)
	if err := ret.All(_context, &arr); err != nil {
		return err
	}
	var minnode = ""
	var minnum = 999999
	for _, result := range arr {
		if result.Num < minnum {
			minnode = result.Id
			minnum = result.Num
		}
	}
	fb := mongodb.NewFilter().Equal("nodeid", minnode).LessEqual("peercount", minnum)
	bs := RoomCol(lobbyid).FindOne(_context, fb)
	if bs.Err() != nil {
		return bs.Err()
	}
	if err := bs.Decode(entity); err != nil {
		return err
	}
	return nil
}

type BalanceResult struct {
	Id  string `bson:"_id"`
	Num int    `bson:"num"`
}
