package toolkit

const (
	LABEL_LOBBY  string = "lobby"
	LABEL_ROOM          = "room"
	PREFIX_LOBBY        = "lobby_"
	PREFIX_ROOM         = "room_"
)

const (
	Api_JoinLobby      string = "JoinLobby"
	Api_LeaveLobby            = "LeaveLobby"
	Api_LobbyStats            = "LobbyStats"
	Api_RoomList              = "RoomList"
	Api_CreateRoom            = "CreateRoom"
	Api_JoinRoom              = "JoinRoom"
	Api_JoinRandomRoom        = "JoinRandomRoom"
	Api_LeaveRoom             = "LeaveRoom"
)

const (
	Event_Join    string = "Join"
	Event_Leave   string = "Leave"
	Event_Disconn string = "Disconn"
	Event_Custom  string = "Custom"
)

const (
	Net_RaiseEvent string = "RaiseEvent"
)

const (
	Group_Other  string = "Others"
	Group_All    string = "All"
	Group_Master string = "Master"
)
