package gen_room

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

var _peerManager *PeerManager

func init() {
	_peerManager = NewPeerManager()
}

func peerManager() *PeerManager {
	return _peerManager
}

type ClientPeer struct {
	peerId string
	roomId string
}

func NewClientPeer(peerId string) *ClientPeer {
	this := new(ClientPeer)
	this.peerId = peerId
	return this
}

func (this *ClientPeer) PeerId() string {
	return this.peerId
}

func (this *ClientPeer) RoomId() string {
	return this.roomId
}

func (this *ClientPeer) SetRoomId(roomId string) {
	this.roomId = roomId
}

func (this *PeerManager) GetClientPeer(peerId string) (*ClientPeer, bool) {
	peer, ok := this.GetPeer(peerId)
	if !ok {
		return nil, false
	}
	cntpeer, ok := peer.(*ClientPeer)
	if !ok {
		return nil, false
	}
	return cntpeer, true
}
