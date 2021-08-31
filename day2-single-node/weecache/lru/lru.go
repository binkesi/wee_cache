package lru

import (
	"container/list"
)

type Cache struct {
	ll        *list.List
	maxBytes  int64
	nBytes    int64
	cache     map[string]*list.Element
	OnEvicted func(key string, value Mvalue)
}

type Mvalue interface {
	Len() int
}

type entry struct {
	key   string
	value Mvalue
}

func New(maxBytes int64, onEvicted func(string, Mvalue)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}

func (c *Cache) Get(key string) (value Mvalue, ok bool) {
	if elm, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elm)
		kv := elm.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Mvalue) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		kv.value = value
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		c.ll.MoveToFront(ele)
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(value.Len()) + int64(len(key))
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}
