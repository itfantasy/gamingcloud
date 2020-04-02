package gen_mmo

type UnsubscriberCollection struct {
	unsubscriber []IDisposer
}

func NewUnsubscriberCollection(unsubscriber ...IDisposer) *UnsubscriberCollection {
	u := new(UnsubscriberCollection)
	u.unsubscriber = unsubscriber
	return u
}

func (u *UnsubscriberCollection) Dispose() {
	for _, item := range u.unsubscriber {
		item.Dispose()
	}
}
