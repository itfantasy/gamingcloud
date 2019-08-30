package gen_lobby

func CreateLobby(lobbyId string) error {
	return nil
}

func DisposeLobby(lobbyId string, force bool) error {
	// 如果强制释放，则位于该大厅的房间将被悉数销毁
	// 否则，等待该大厅所有房间均释放时再自动销毁（标志位）
	return nil
}

func JoinLobby(lobbyId string) error {
	return nil
}

func LeaveLobby(lobbyId string) error {
	return nil
}

func LobbyStats(lobbyId string) error {
	return nil
}

func RoomList(lobbyId string) error {
	return nil
}

func CreateRoom(roomId string) error {
	// 自动找寻负载最低的roomserver，并向其请求创建房间
	return nil
}

func JoinRoom(roomId string) error {
	return nil
}

func JoinRandomRoom() error {
	return nil
}
