package gen_mmo

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
)

type MmoServerInfo struct {
	Id       string
	Url      string
	LogLevel string
	LogComp  string
	RegComp  string
}

func (m *MmoServerInfo) ExpandToNodeInfo() *gen_server.NodeInfo {
	info := gen_server.NewNodeInfo()
	info.Url = m.Url
	info.Pub = true
	info.BackEnds = ""
	info.LogLevel = m.LogLevel
	info.LogComp = m.LogComp
	info.RegComp = m.RegComp
	return info
}
