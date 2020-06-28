package gen_room

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/strs"

	"github.com/itfantasy/gonode-toolkit/toolkit"
)

type RoomServerInfo struct {
	RegDC     string
	NameSpace string
	NodeId    string
	EndPoints []string
	GameDB    string
	PubDomain string
}

func (serverInfo *RoomServerInfo) ExpandToNodeInfo() *gen_server.NodeInfo {
	info := gen_server.NewNodeInfo()
	info.RegDC = serverInfo.RegDC
	info.NameSpace = serverInfo.NameSpace
	if strs.StartsWith(serverInfo.NodeId, toolkit.LABEL_ROOM) {
		info.NodeId = serverInfo.NodeId
	} else {
		info.NodeId = toolkit.PREFIX_ROOM + serverInfo.NodeId
	}
	info.EndPoints = serverInfo.EndPoints
	info.IsPub = true
	info.BackEnds = toolkit.LABEL_LOBBY
	info.UsrDatas[toolkit.USRDATA_PUBDOMAIN] = serverInfo.PubDomain
	return info
}
