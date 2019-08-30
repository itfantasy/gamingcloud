package gen_lobby

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/strs"

	"github.com/itfantasy/gonode-toolkit/toolkit"
)

type LobbyServerInfo struct {
	Id       string
	Url      string
	LogLevel string
	LogComp  string
	RegComp  string
}

func ExpandToNodeInfo(serverInfo LobbyServerInfo) *gen_server.NodeInfo {
	info := gen_server.NewNodeInfo()
	if strs.StartsWith(serverInfo.Id, toolkit.PREFIX_LOBBY) {
		info.Id = serverInfo.Id
	} else {
		info.Id = toolkit.PREFIX_LOBBY + serverInfo.Id
	}
	info.Url = serverInfo.Url
	info.Pub = true
	info.BackEnds = ""
	info.LogLevel = serverInfo.LogLevel
	info.LogComp = serverInfo.LogComp
	info.RegComp = serverInfo.RegComp
	return info
}
