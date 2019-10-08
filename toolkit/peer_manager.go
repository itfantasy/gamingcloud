package toolkit

import (
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
	return this.dict.Add(peer.PeerId(), peer)
}

func (this *PeerManager) RemovePeer(peerId string) error {
	return this.dict.Remove(peerId)
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
