package gocache

import (
	"fmt"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (g GetterFunc) Get(key string) ([]byte, error) {
	return g(key)
}

type Group struct {
	name      string
	getter    Getter
	mainCache Cache
}

var (
	rw     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, capacity int64, getter Getter) *Group {
	if getter == nil {
		panic("nil getter")
	}

	rw.Lock()
	defer rw.Unlock()

	temp := &Group{
		name:   name,
		getter: getter,
		mainCache: Cache{
			cacheByte: capacity,
		},
	}

	groups[name] = temp
	return temp
}

func GetGroup(name string) *Group {
	rw.RLock()
	temp := groups[name]
	rw.RUnlock()
	return temp
}

func (this *Group) Get(key string) (Byteview, error) {
	if key == "" {
		return Byteview{}, fmt.Errorf("no key")
	}

	if v, ok := this.mainCache.Get(key); ok {
		log.Println("Group", this.name, key, "hit")
		return v, nil
	}

	return this.Load(key)
}

func (this *Group) Load(key string) (Byteview, error) {
	return this.GetLocally(key)
}

func (this *Group) GetLocally(key string) (Byteview, error) {
	bytes, err := this.getter.Get(key)
	if err != nil {
		log.Println("Local load failed -", key)
		return Byteview{}, err
	}

	v := make([]byte, len(bytes))
	copy(v, bytes)
	bv := Byteview{
		v: v,
	}

	this.PopulateCache(key, bv)
	return bv, nil
}

func (this *Group) PopulateCache(key string, value Byteview) {
	this.mainCache.Put(key, value)
}
