package gen_room

import (
	"github.com/itfantasy/gonode-toolkit/toolkit"
)

var _peerManager *toolkit.PeerManager

func init() {
	_peerManager = toolkit.NewPeerManager()
}

func peerManager() *toolkit.PeerManager {
	return _peerManager
}

type RoomPeer struct {
	peerId string
	roomId string
}

func NewRoomPeer(peerId string) *RoomPeer {
	this := new(RoomPeer)
	this.peerId = peerId
	return this
}

func (r *RoomPeer) PeerId() string {
	return r.peerId
}

func (r *RoomPeer) RoomId() string {
	return r.roomId
}

func (r *RoomPeer) SetRoomId(roomId string) {
	r.roomId = roomId
}

func getRoomPeer(peerId string) (*RoomPeer, bool) {
	peer, ok := peerManager().GetPeer(peerId)
	if !ok {
		return nil, false
	}
	cntpeer, ok := peer.(*RoomPeer)
	if !ok {
		return nil, false
	}
	return cntpeer, true
}
