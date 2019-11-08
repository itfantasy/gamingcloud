package toolkit

import (
	"github.com/itfantasy/gonode"
)

const (
	USRDATA_PUBDOMAIN string = "PubDomain"
)

const (
	LABEL_LOBBY   string = "lobby"
	LABEL_ROOM           = "room"
	PREFIX_LOBBY         = "lobby_"
	PREFIX_ROOM          = "room_"
	DEFAULT_LOBBY        = "__default_lobby"
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
	Group_Others string = "Others"
	Group_All    string = "All"
	Group_Master string = "Master"
)

var (
	Err_InvalidRequestParameters error = gonode.Error(-6, "InvalidRequestParameters")
	Err_ArgumentOutOfRange             = gonode.Error(-4, "ArgumentOutOfRange")
	Err_OperationDenied                = gonode.Error(-3, "OperationDenied")
	Err_OperationInvalid               = gonode.Error(-2, "OperationInvalid")
	Err_InternalServerError            = gonode.Error(-1, "InternalServerError")
	//Ok
	Err_InvalidAuthentication          = gonode.Error(32767, "InvalidAuthentication")
	Err_RoomIdAlreadyExists            = gonode.Error(32766, "RoomIdAlreadyExists")
	Err_RoomFull                       = gonode.Error(32765, "RoomFull")
	Err_RoomClosed                     = gonode.Error(32764, "RoomClosed")
	Err_AlreadyMatched                 = gonode.Error(32763, "AlreadyMatched")
	Err_ServerFull                     = gonode.Error(32762, "ServerFull")
	Err_UserBlocked                    = gonode.Error(32761, "UserBlocked")
	Err_NoMatchFound                   = gonode.Error(32760, "NoMatchFound")
	Err_RedirectRepeat                 = gonode.Error(32759, "RedirectRepeat")
	Err_RoomIdNotExists                = gonode.Error(32758, "RoomIdNotExists")
	Err_MaxCcuReached                  = gonode.Error(32757, "MaxCcuReached")
	Err_InvalidRegion                  = gonode.Error(32756, "InvalidRegion")
	Err_CustomAuthenticationFailed     = gonode.Error(32755, "CustomAuthenticationFailed")
	Err_AuthenticationTokenExpired     = gonode.Error(32753, "AuthenticationTokenExpired")
	Err_PluginReportedError            = gonode.Error(32752, "PluginReportedError")
	Err_PluginMismatch                 = gonode.Error(32751, "PluginMismatch")
	Err_JoinFailedPeerAlreadyJoined    = gonode.Error(32750, "JoinFailedPeerAlreadyJoined")
	Err_JoinFailedFoundInactiveJoiner  = gonode.Error(32749, "JoinFailedFoundInactiveJoiner")
	Err_JoinFailedWithRejoinerNotFound = gonode.Error(32748, "JoinFailedWithRejoinerNotFound")
	Err_JoinFailedFoundExcludedUserId  = gonode.Error(32747, "JoinFailedFoundExcludedUserId")
	Err_JoinFailedFoundActiveJoiner    = gonode.Error(32746, "JoinFailedFoundActiveJoiner")
	Err_HttpLimitReached               = gonode.Error(32745, "HttpLimitReached")
	Err_ExternalHttpCallFailed         = gonode.Error(32744, "ExternalHttpCallFailed")
)
