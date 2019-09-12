package gamedb

import (
	"errors"

	"golang.org/x/net/context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/itfantasy/gonode/components"
	"github.com/itfantasy/gonode/components/mongodb"
)

const (
	LOBBY_COLLECT string = "_lobbys"
	ROOM_COLLECT         = "_rooms"
)

var _context context.Context
var _mongo *mongodb.MongoDB

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

func RoomCol() *mongo.Collection {
	return _mongo.Collect(ROOM_COLLECT)
}

func CreateLobby(entity interface{}) error {
	if _, err := LobbyCol().InsertOne(_context, entity); err != nil {
		return err
	}
	return nil
}

func UpdateLobby(filter map[string]interface{}, data map[string]interface{}) error {
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

func CreateRoom(entity interface{}) error {
	if _, err := RoomCol().InsertOne(_context, entity); err != nil {
		return err
	}
	return nil
}

func UpdateRoom(entity interface{}, filter map[string]interface{}) error {
	if _, err := RoomCol().UpdateOne(_context, filter, entity); err != nil {
		return err
	}
	return nil
}

func DeleteRoom(filter map[string]interface{}) error {
	if _, err := RoomCol().DeleteOne(_context, filter); err != nil {
		return err
	}
	return nil
}

func DeleteRooms(filter map[string]interface{}) error {
	if _, err := RoomCol().DeleteMany(_context, filter); err != nil {
		return err
	}
	return nil
}

func FindRoom(filter map[string]interface{}, entity interface{}) error {
	ret := RoomCol().FindOne(_context, filter)
	if ret.Err() != nil {
		return ret.Err()
	}
	if err := ret.Decode(entity); err != nil {
		return err
	}
	return nil
}

func FindRooms(filter map[string]interface{}, entities interface{}) error {
	cursor, err := RoomCol().Find(_context, filter)
	if err != nil {
		return err
	}
	if cursor.Err() != nil {
		return cursor.Err()
	}
	if err := cursor.All(_context, entities); err != nil {
		return err
	}
	return nil
}
