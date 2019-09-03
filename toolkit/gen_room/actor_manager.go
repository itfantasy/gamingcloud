package gen_room

import (
	"errors"
	"sync"

	"github.com/itfantasy/gonode/utils/stl"
)

type Actor struct {
	peerId  string
	actorNr int32
}

func NewActor(peerId string, actorNr int32) *Actor {
	a := new(Actor)
	a.peerId = peerId
	a.actorNr = actorNr
	return a
}

func (a *Actor) PeerId() string {
	return a.peerId
}

func (a *Actor) ActorNr() int32 {
	return a.actorNr
}

type ActorsManager struct {
	allActors  *stl.List
	lock       sync.RWMutex
	curActorNr int32
}

func NewActorsManager() *ActorsManager {
	a := new(ActorsManager)
	a.allActors = stl.NewList(10)
	a.curActorNr = 0
	return a
}

func (a *ActorsManager) CreateActorNr() int32 {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.curActorNr += 1
	return a.curActorNr
}

func (a *ActorsManager) AddNewActor(peerId string) (*Actor, error) {
	actorNr := a.CreateActorNr()
	if _, exist := a.GetActorByNr(actorNr); exist {
		return nil, errors.New("there has been an actor with the same actorNr!")
	}
	if _, exist := a.GetActorByPeerId(peerId); exist {
		return nil, errors.New("there has been an actor with the same peerId!")
	}
	actor := NewActor(peerId, actorNr)
	a.allActors.Add(actor)
	return actor, nil
}

func (a *ActorsManager) GetActorByNr(actorNr int32) (*Actor, bool) {
	for _, item := range a.allActors.Values() {
		actor := item.(*Actor)
		if actor.actorNr == actorNr {
			return actor, true
		}
	}
	return nil, false
}

func (a *ActorsManager) GetActorByPeerId(peerId string) (*Actor, bool) {
	for _, item := range a.allActors.Values() {
		actor := item.(*Actor)
		if actor.peerId == peerId {
			return actor, true
		}
	}
	return nil, false
}

func (a *ActorsManager) GetActorByIndex(index int) (*Actor, bool) {
	actor, err := a.allActors.Get(index)
	if err != nil {
		return nil, false
	}
	return actor.(*Actor), true
}

func (a *ActorsManager) RemoveActorByNr(actorNr int32) (*Actor, bool) {
	if actor, exist := a.GetActorByNr(actorNr); exist {
		a.allActors.Remove(actor)
		return actor, true
	}
	return nil, false
}

func (a *ActorsManager) RemoveActorByPeer(peerId string) (*Actor, bool) {
	if actor, exist := a.GetActorByPeerId(peerId); exist {
		a.allActors.Remove(actor)
		return actor, true
	}
	return nil, false
}

func (a *ActorsManager) GetAllActorNrs() []int32 {
	list := make([]int32, 0, a.allActors.Count())
	for _, item := range a.allActors.Values() {
		actor := item.(*Actor)
		list = append(list, actor.actorNr)
	}
	return list
}

func (a *ActorsManager) GetAllPeerIds() []string {
	list := make([]string, 0, a.allActors.Count())
	for _, item := range a.allActors.Values() {
		actor := item.(*Actor)
		list = append(list, actor.peerId)
	}
	return list
}

func (a *ActorsManager) ActorsCount() int {
	return a.allActors.Count()
}

func (a *ActorsManager) ClearAll() {
	a.allActors.Clear()
}
