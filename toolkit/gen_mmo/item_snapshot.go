package gen_mmo

type ItemSnapshot struct {
	source             *MmoItem
	position           *Vector
	rotation           *Vector
	propertiesRevision int
}

func NewItemSnapshot(source *MmoItem, position *Vector, rotation *Vector, worldRegion *Region, propertiesRevision int) *ItemSnapshot {
	i := new(ItemSnapshot)
	i.source = source
	i.position = position
	i.rotation = rotation
	i.propertiesRevision = propertiesRevision
	return i
}

func (i *ItemSnapshot) Source() *MmoItem {
	return i.source
}

func (i *ItemSnapshot) Position() *Vector {
	return i.position
}

func (i *ItemSnapshot) Rotation() *Vector {
	return i.rotation
}

func (i *ItemSnapshot) PropertiesRevision() int {
	return i.propertiesRevision
}
