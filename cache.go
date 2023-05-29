package gocache

import (
	"gocache/lru"
	"sync"
)

type Cache struct {
	lock      sync.Mutex
	lru       *lru.Cache
	cacheByte int64
}

func (this *Cache) Put(key string, value Byteview) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.lru == nil {
		this.lru = lru.New(this.cacheByte, nil)
	}
	this.lru.Put(key, value)
}

func (this *Cache) Get(key string) (value Byteview, ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	if this.lru == nil {
		return
	}

	if res, ok := this.lru.Get(key); ok {
		return res.(Byteview), ok
	}
	return
}
