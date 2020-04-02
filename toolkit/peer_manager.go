package toolkit

import (
	"errors"

	"github.com/itfantasy/gonode/utils/stl"
)

type Peer interface {
	PeerId() string
}

type PeerManager struct {
	dict *stl.Dictionary
}

func NewPeerManager() *PeerManager {
	this := new(PeerManager)
	this.dict = stl.NewDictionary()
	return this
}

func (this *PeerManager) AddPeer(peer Peer) error {
	if ret := this.dict.Add(peer.PeerId(), peer); !ret {
		return errors.New("the manager has contains the peer!" + peer.PeerId())
	}
	return nil
}

func (this *PeerManager) RemovePeer(peerId string) error {
	if ret := this.dict.Remove(peerId); !ret {
		return errors.New("the manager can not find the peer!" + peerId)
	}
	return nil
}

func (this *PeerManager) GetPeer(peerId string) (Peer, bool) {
	ret, exist := this.dict.Get(peerId)
	if !exist {
		return nil, false
	}
	peer, ok := ret.(Peer)
	if !ok {
		return nil, false
	}
	return peer, true
}
