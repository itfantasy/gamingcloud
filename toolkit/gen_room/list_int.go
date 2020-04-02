package gen_room

import (
	"errors"
	"sync"
	//"fmt"
)

type ListInt struct {
	array []int32
	sync.RWMutex
}

func NewListInt(capacity int) *ListInt {
	list := ListInt{}
	list.array = make([]int32, 0, capacity)
	return &list
}

func NewListIntRaw(raw []int32) *ListInt {
	list := ListInt{}
	list.array = raw
	return &list
}

func (l *ListInt) Add(item int32) {
	l.Lock()
	defer l.Unlock()

	l.array = append(l.array, item)
}

func (l *ListInt) Insert(index int, item int32) error {
	l.Lock()
	defer l.Unlock()

	if index > len(l.array) {
		return errors.New("ArgumentOutOfRange")
	}

	temp := make([]int32, 0)
	after := append(temp, l.array[index:]...)
	before := l.array[0:index]
	l.array = append(before, item)
	l.array = append(l.array, after...)
	return nil
}

func (l *ListInt) RemoveAt(index int) error {
	l.Lock()
	defer l.Unlock()

	if index > len(l.array) {
		return errors.New("ArgumentOutOfRange")
	}

	l.array = append(l.array[:index], l.array[index+1:]...)
	return nil
}

func (l *ListInt) Remove(item int32) bool {
	index := l.IndexOf(item)
	if index < 0 {
		return false
	}
	l.RemoveAt(index)
	return true
}

func (l *ListInt) IndexOf(item int32) int {
	l.RLock()
	defer l.RUnlock()

	count := len(l.array)
	for i := 0; i < count; i++ {
		if l.array[i] == item {
			return i
		}
	}
	return -1
}

func (l *ListInt) Contains(item int32) bool {
	return l.IndexOf(item) >= 0
}

func (l *ListInt) Count() int {
	l.RLock()
	defer l.RUnlock()

	return len(l.array)
}

func (l *ListInt) Capacity() int {
	l.RLock()
	defer l.RUnlock()

	return cap(l.array)
}

func (l *ListInt) Items() []int32 {
	l.RLock()
	defer l.RUnlock()

	return l.array
}

func (l *ListInt) Get(index int) (int32, error) {
	l.RLock()
	defer l.RUnlock()

	if index >= len(l.array) {
		return 0, errors.New("ArgumentOutOfRange")
	}
	return l.array[index], nil
}

func (l *ListInt) Set(index int, item int32) error {
	l.Lock()
	defer l.Unlock()

	if index > len(l.array) {
		return errors.New("ArgumentOutOfRange")
	}
	l.array[index] = item
	return nil
}

func (l *ListInt) ForEach(fun func(int32)) {
	l.RLock()
	defer l.RUnlock()

	for _, v := range l.array {
		fun(v)
	}
}

func (l *ListInt) Clear() {
	l.Lock()
	defer l.Unlock()

	l.array = l.array[0:0]
}
