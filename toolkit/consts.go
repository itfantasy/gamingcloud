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
	Err_InvalidRequestParameters error = gonode.CustomError(-6, "InvalidRequestParameters")
	Err_ArgumentOutOfRange             = gonode.CustomError(-4, "ArgumentOutOfRange")
	Err_OperationDenied                = gonode.CustomError(-3, "OperationDenied")
	Err_OperationInvalid               = gonode.CustomError(-2, "OperationInvalid")
	Err_InternalServerError            = gonode.CustomError(-1, "InternalServerError")
	//Ok
	Err_InvalidAuthentication          = gonode.CustomError(32767, "InvalidAuthentication")
	Err_RoomIdAlreadyExists            = gonode.CustomError(32766, "RoomIdAlreadyExists")
	Err_RoomFull                       = gonode.CustomError(32765, "RoomFull")
	Err_RoomClosed                     = gonode.CustomError(32764, "RoomClosed")
	Err_AlreadyMatched                 = gonode.CustomError(32763, "AlreadyMatched")
	Err_ServerFull                     = gonode.CustomError(32762, "ServerFull")
	Err_UserBlocked                    = gonode.CustomError(32761, "UserBlocked")
	Err_NoMatchFound                   = gonode.CustomError(32760, "NoMatchFound")
	Err_RedirectRepeat                 = gonode.CustomError(32759, "RedirectRepeat")
	Err_RoomIdNotExists                = gonode.CustomError(32758, "RoomIdNotExists")
	Err_MaxCcuReached                  = gonode.CustomError(32757, "MaxCcuReached")
	Err_InvalidRegion                  = gonode.CustomError(32756, "InvalidRegion")
	Err_CustomAuthenticationFailed     = gonode.CustomError(32755, "CustomAuthenticationFailed")
	Err_AuthenticationTokenExpired     = gonode.CustomError(32753, "AuthenticationTokenExpired")
	Err_PluginReportedError            = gonode.CustomError(32752, "PluginReportedError")
	Err_PluginMismatch                 = gonode.CustomError(32751, "PluginMismatch")
	Err_JoinFailedPeerAlreadyJoined    = gonode.CustomError(32750, "JoinFailedPeerAlreadyJoined")
	Err_JoinFailedFoundInactiveJoiner  = gonode.CustomError(32749, "JoinFailedFoundInactiveJoiner")
	Err_JoinFailedWithRejoinerNotFound = gonode.CustomError(32748, "JoinFailedWithRejoinerNotFound")
	Err_JoinFailedFoundExcludedUserId  = gonode.CustomError(32747, "JoinFailedFoundExcludedUserId")
	Err_JoinFailedFoundActiveJoiner    = gonode.CustomError(32746, "JoinFailedFoundActiveJoiner")
	Err_HttpLimitReached               = gonode.CustomError(32745, "HttpLimitReached")
	Err_ExternalHttpCallFailed         = gonode.CustomError(32744, "ExternalHttpCallFailed")
)
