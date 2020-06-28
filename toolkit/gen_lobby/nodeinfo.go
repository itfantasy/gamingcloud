package gen_lobby

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/strs"

	"github.com/itfantasy/gonode-toolkit/toolkit"
)

type LobbyServerInfo struct {
	RegDC     string
	NameSpace string
	NodeId    string
	EndPoints []string
	GameDB    string
}

func (serverInfo *LobbyServerInfo) ExpandToNodeInfo() *gen_server.NodeInfo {
	info := gen_server.NewNodeInfo()
	info.RegDC = serverInfo.RegDC
	info.NameSpace = serverInfo.NameSpace
	if strs.StartsWith(serverInfo.NodeId, toolkit.LABEL_LOBBY) {
		info.NodeId = serverInfo.NodeId
	} else {
		info.NodeId = toolkit.PREFIX_LOBBY + serverInfo.NodeId
	}
	info.EndPoints = serverInfo.EndPoints
	info.IsPub = true
	info.BackEnds = ""
	return info
}
