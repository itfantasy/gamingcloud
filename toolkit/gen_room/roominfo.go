package gen_room

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
	"github.com/itfantasy/gonode/utils/strs"

	"github.com/itfantasy/gonode-toolkit/toolkit"
)

// 基础房间节点信息
type RoomServerInfo struct {
	Id       string
	Url      string
	LogLevel string
	LogComp  string
	RegComp  string

	PubDomain string
}

func ExpandToNodeInfo(serverInfo *RoomServerInfo) *gen_server.NodeInfo {
	info := gen_server.NewNodeInfo()
	if strs.StartsWith(serverInfo.Id, toolkit.PREFIX_ROOM) {
		info.Id = serverInfo.Id
	} else {
		info.Id = toolkit.PREFIX_ROOM + serverInfo.Id
	}
	info.Url = serverInfo.Url
	info.Pub = true
	info.BackEnds = toolkit.LABEL_LOBBY
	info.LogLevel = serverInfo.LogLevel
	info.LogComp = serverInfo.LogComp
	info.RegComp = serverInfo.RegComp
	info.UserDatas["PubDomain"] = serverInfo.PubDomain
	return info
}
