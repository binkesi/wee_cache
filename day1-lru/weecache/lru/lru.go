package lru

import (
	"container/list"
)

type Cache struct {
	ll        *list.List
	maxBytes  int64
	nBytes    int64
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}
