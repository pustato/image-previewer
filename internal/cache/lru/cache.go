package lru

import (
	"container/list"
	"sync"
)

var _ Cache = (*CacheLRU)(nil)

type Cache interface {
	Get(key string) (*Item, bool)
	Set(key string, item *Item) bool
}

type Item struct {
	key      string
	FileName string
	Size     uint64
}

type RemoveItemCallback func(item *Item)

type CacheLRU struct {
	mu           sync.RWMutex
	list         *list.List
	items        map[string]*list.Element
	limit        uint64
	size         uint64
	onRemoveFunc RemoveItemCallback
}

func NewCache(limit uint64, onRemove RemoveItemCallback) *CacheLRU {
	return &CacheLRU{
		list:         list.New(),
		items:        make(map[string]*list.Element),
		limit:        limit,
		onRemoveFunc: onRemove,
	}
}

func (c *CacheLRU) Get(key string) (*Item, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	element, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.list.MoveToFront(element)

	return element.Value.(*Item), true
}

func (c *CacheLRU) Set(key string, item *Item) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item.key = key

	if element, ok := c.items[key]; ok {
		element.Value = item
		c.list.MoveToFront(element)

		return true
	}

	element := c.list.PushFront(item)
	c.items[key] = element
	c.size += item.Size

	c.gc()

	return false
}

func (c *CacheLRU) remove(key string) {
	element, ok := c.items[key]
	if !ok {
		return
	}

	delete(c.items, key)
	item := element.Value.(*Item)
	c.list.Remove(element)
	c.size -= item.Size

	c.onRemoveFunc(item)
}

func (c *CacheLRU) gc() {
	for c.size > c.limit {
		element := c.list.Back()
		item := element.Value.(*Item)
		c.remove(item.key)
	}
}
