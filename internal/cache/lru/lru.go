package lru

import (
	"container/list"
	"sync"
)

type Lru struct {
	mu       *sync.RWMutex
	itemMap  map[string]*list.Element
	queue    *list.List
	capacity int
}

type item struct {
	key   string
	value interface{}
}

func New(capacity int) *Lru {
	return &Lru{
		mu:       new(sync.RWMutex),
		itemMap:  make(map[string]*list.Element),
		queue:    list.New(),
		capacity: capacity,
	}
}

func (l *Lru) Set(key string, value interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if elem, ok := l.itemMap[key]; ok {
		item := elem.Value.(*item)
		item.value = value
		l.queue.MoveToFront(elem)
		return
	}

	if l.queue.Len() == l.capacity {
		l.removeLeastRecentlyUsed()
	}

	item := &item{key: key, value: value}
	elem := l.queue.PushFront(item)
	l.itemMap[key] = elem
}

func (l *Lru) Get(key string) (val interface{}, ok bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if elem, ok := l.itemMap[key]; ok {
		return elem.Value.(*item).value, true
	}
	return struct{}{}, false
}

func (l *Lru) Delete(key string) (ok bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	elem, ok := l.itemMap[key]
	if !ok {
		return false
	}

	l.queue.Remove(elem)
	delete(l.itemMap, key)
	return true
}

func (l *Lru) removeLeastRecentlyUsed() {
	lruElem := l.queue.Back()
	l.queue.Remove(lruElem)
	delete(l.itemMap, lruElem.Value.(*item).key)
}
