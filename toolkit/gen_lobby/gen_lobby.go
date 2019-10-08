package gen_lobby

import (
	"errors"
)

// --------------------- super admin

func CreateLobby(lobbyId string) (*LobbyEntity, error) {
	return lobbyManager().CreateLobby(lobbyId)
}

func DisposeLobby(lobbyId string, force bool) error {
	// 如果强制释放，则位于该大厅的房间将被悉数销毁
	if force {
		return lobbyManager().DisposeLobby(lobbyId)
	}
	// 否则，等待该大厅所有房间均释放时再自动销毁（标志位）
	return nil
}

func LobbyStats(lobbyId string) error {
	return nil
}

func RoomList(lobbyId string, startIndex int, endIndex int) error {
	return nil
}

// --------------------- guest usr

func peerCannotFind(peerId string) error {
	return errors.New("cannot find the peer:" + peerId)
}

func JoinLobby(peerId string, lobbyId string) error {
	l, exist := getLobbyPeer(peerId)
	if !exist {
		return peerCannotFind(peerId)
	}
	l.SetLobbyId(lobbyId)
	return nil
}

func LeaveLobby(peerId string, lobbyId string) error {
	l, exist := getLobbyPeer(peerId)
	if !exist {
		return peerCannotFind(peerId)
	}
	l.SetDefaultLobby()
	return nil
}

func CreateRoom(peerId string, roomId string) error {
	l, exist := getLobbyPeer(peerId)
	if !exist {
		return peerCannotFind(peerId)
	}
	lobby, err := lobbyManager().FindLobby(l.LobbyId())
	if err != nil {
		return err
	}
	_, err := lobby.CreateRoom(roomId)
	// TODO
	return nil
}

func JoinRoom(peerId string, roomId string) error {
	return nil
}

func JoinRandomRoom(peerId string) error {
	return nil
}
