package lru

import (
	"container/list"
)

type Value interface {
	Len() int
}

type Node struct {
	key   string
	value Value
}

type Cache struct {
	maxBytes  int64
	nowBytes  int64
	ll        *list.List
	hashmap   map[string]*list.Element
	onEvicted func(key string, value Value)
}

func New(capacity int64, OnEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  capacity,
		nowBytes:  0,
		ll:        list.New(),
		hashmap:   make(map[string]*list.Element),
		onEvicted: OnEvicted,
	}
}

func (this *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := this.hashmap[key]; ok {
		this.ll.MoveToFront(ele)
		v := ele.Value.(*Node)
		return v.value, true
	}
	return
}

func (this *Cache) Put(key string, value Value) {
	ncap := int64(value.Len()) + int64(len(key))
	if ele, ok := this.hashmap[key]; ok {
		this.ll.MoveToFront(ele)
		node := ele.Value.(*Node)
		this.nowBytes += ncap - (int64(node.value.Len()) + int64(len(node.key)))
		node.value = value
	} else {
		ele := this.ll.PushFront(&Node{
			key:   key,
			value: value,
		})
		this.hashmap[key] = ele
		this.nowBytes += ncap
	}
	for this.nowBytes > this.maxBytes && this.maxBytes != 0 {
		this.Remove()
	}
}

func (this *Cache) Remove() {
	tail := this.ll.Back()
	if tail != nil {
		node := tail.Value.(*Node)
		this.ll.Remove(tail)
		delete(this.hashmap, node.key)
		this.nowBytes -= int64(node.value.Len()) + int64(len(node.key))

		if this.onEvicted != nil {
			this.onEvicted(node.key, node.value)
		}
	}
}

func (this *Cache) Len() int {
	return this.ll.Len()
}
