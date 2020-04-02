package gen_mmo

type Region struct {
	x                        int
	y                        int
	itemRegionChangedChannel *MessageChannel
	requestItemEnterChannel  *MessageChannel
	requestItemExitChannel   *MessageChannel
	itemEventChannel         *MessageChannel
}

func NewRegion(x, y int) *Region {
	r := new(Region)
	r.x = x
	r.y = y
	return r
}

func (r *Region) X() int {
	return r.x
}

func (r *Region) Y() int {
	return r.y
}

func (r *Region) ItemRegionChangedChannel() *MessageChannel {
	return r.itemRegionChangedChannel
}

func (r *Region) RequestItemEnterChannel() *MessageChannel {
	return r.requestItemEnterChannel
}

func (r *Region) RequestItemExitChannel() *MessageChannel {
	return r.requestItemExitChannel
}

func (r *Region) ItemEventChannel() *MessageChannel {
	return r.itemEventChannel
}

func (r *Region) Dispose() {
	r.itemRegionChangedChannel.Dispose()
	r.requestItemEnterChannel.Dispose()
	r.requestItemExitChannel.Dispose()
	r.itemEventChannel.Dispose()
}
