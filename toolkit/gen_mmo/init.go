package gen_mmo

import (
	"reflect"

	"github.com/itfantasy/gonode-toolkit/toolkit"
	"github.com/itfantasy/gonode/core/binbuf"
)

const (
	CustomType_Vector byte = iota
	CustomType_BoundingBox
)

var _worldCache *WorldManager

var _peerManager *toolkit.PeerManager

func init() {
	_peerManager = toolkit.NewPeerManager()
	_worldCache = NewWorldManager()

	zero := NewVector(0, 0, 0)
	binbuf.ExtendCustomType(reflect.TypeOf(zero), CustomType_Vector, func(b *binbuf.BinBuffer, value interface{}) {
		v := value.(*Vector)
		b.PushFloat(float32(v.X()))
		b.PushFloat(float32(v.Y()))
		b.PushFloat(float32(v.Z()))
	}, func(p *binbuf.BinParser) interface{} {
		v := NewVector(
			float64(p.Float()),
			float64(p.Float()),
			float64(p.Float()))
		return v
	})
	binbuf.ExtendCustomType(reflect.TypeOf(NewBoundingBox(zero, zero)), CustomType_BoundingBox, func(b *binbuf.BinBuffer, value interface{}) {
		v := value.(*BoundingBox)
		b.PushFloat(float32(v.Min().X()))
		b.PushFloat(float32(v.Min().Y()))
		b.PushFloat(float32(v.Min().Z()))
		b.PushFloat(float32(v.Max().X()))
		b.PushFloat(float32(v.Max().Y()))
		b.PushFloat(float32(v.Max().Z()))
	}, func(p *binbuf.BinParser) interface{} {
		b := NewBoundingBox(NewVector(
			float64(p.Float()),
			float64(p.Float()),
			float64(p.Float())), NewVector(
			float64(p.Float()),
			float64(p.Float()),
			float64(p.Float())))
		return b
	})
}

func peerManager() *toolkit.PeerManager {
	return _peerManager
}

func worldCacheIns() *WorldManager {
	return _worldCache
}
