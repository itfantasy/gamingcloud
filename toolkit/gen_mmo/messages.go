package gen_mmo

type ItemRegionChangedMessage struct {
	region0      *Region
	region1      *Region
	itemSnapshot *ItemSnapshot
}

func NewItemRegionChangedMessage(r0 *Region, r1 *Region, snapshot *ItemSnapshot) *ItemRegionChangedMessage {
	i := new(ItemRegionChangedMessage)
	i.region0 = r0
	i.region1 = r1
	i.itemSnapshot = snapshot
	return i
}

func (i *ItemRegionChangedMessage) Region0() *Region {
	return i.region0
}

func (i *ItemRegionChangedMessage) Region1() *Region {
	return i.region1
}

func (i *ItemRegionChangedMessage) ItemSnapshot() *ItemSnapshot {
	return i.itemSnapshot
}

type RequestItemEnterMessage struct {
	interestArea *InterestArea
}

func NewRequestItemEnterMessage(interestArea *InterestArea) *RequestItemEnterMessage {
	i := new(RequestItemEnterMessage)
	i.interestArea = interestArea
	return i
}

func (i *RequestItemEnterMessage) InterestArea() *InterestArea {
	return i.interestArea
}

type RequestItemExitMessage struct {
	interestArea *InterestArea
}

func NewRequestItemExitMessage(interestArea *InterestArea) *RequestItemExitMessage {
	i := new(RequestItemExitMessage)
	i.interestArea = interestArea
	return i
}

func (i *RequestItemExitMessage) InterestArea() *InterestArea {
	return i.interestArea
}

type ItemDisposedMessage struct {
	source *MmoItem
}

func NewItemDisposedMessage(source *MmoItem) *ItemDisposedMessage {
	i := new(ItemDisposedMessage)
	i.source = source
	return i
}

func (i *ItemDisposedMessage) Source() *MmoItem {
	return i.source
}

type ItemPositionMessage struct {
	source   *MmoItem
	position *Vector
}

func NewItemPositionMessage(source *MmoItem, position *Vector) *ItemPositionMessage {
	i := new(ItemPositionMessage)
	i.source = source
	i.position = position
	return i
}

func (i *ItemPositionMessage) Source() *MmoItem {
	return i.source
}

func (i *ItemPositionMessage) Position() *Vector {
	return i.position
}

type ItemEventMessage struct {
	source    *MmoItem
	eventData *EventData
}

func NewItemEventMessage(source *MmoItem, code byte, datas interface{}) *ItemEventMessage {
	i := new(ItemEventMessage)
	i.source = source
	i.eventData = &EventData{
		Code:  code,
		Datas: datas,
	}
	return i
}

func (i *ItemEventMessage) Source() *MmoItem {
	return i.source
}

func (i ItemEventMessage) EventData() *EventData {
	return i.eventData
}
