package gen_lobby

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

type LobbyPeer struct {
	peerId  string
	lobbyId string
}

func NewLobbyPeer(peerId string) *LobbyPeer {
	this := new(LobbyPeer)
	this.peerId = peerId
	this.lobbyId = toolkit.DEFAULT_LOBBY
	return this
}

func (l *LobbyPeer) PeerId() string {
	return l.peerId
}

func (l *LobbyPeer) LobbyId() string {
	return l.lobbyId
}

func (l *LobbyPeer) SetLobbyId(lobbyId string) {
	l.lobbyId = lobbyId
}

func (l *LobbyPeer) SetDefaultLobby() {
	l.lobbyId = toolkit.DEFAULT_LOBBY
}

func getLobbyPeer(peerId string) (*LobbyPeer, bool) {
	peer, ok := peerManager().GetPeer(peerId)
	if !ok {
		return nil, false
	}
	cntpeer, ok := peer.(*LobbyPeer)
	if !ok {
		return nil, false
	}
	return cntpeer, true
}
