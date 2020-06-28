package gen_mmo

import (
	"github.com/itfantasy/gonode/behaviors/gen_server"
)

type MmoServerInfo struct {
	RegDC     string
	NameSpace string
	NodeId    string
	EndPoints []string
}

func (m *MmoServerInfo) ExpandToNodeInfo() *gen_server.NodeInfo {
	info := gen_server.NewNodeInfo()
	info.RegDC = m.RegDC
	info.NameSpace = m.NameSpace
	info.NodeId = m.NodeId
	info.EndPoints = m.EndPoints
	info.IsPub = true
	info.BackEnds = ""
	return info
}
