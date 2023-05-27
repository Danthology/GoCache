package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestFind(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Put("key1", String("1234"))
	if v, ok := lru.Find("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := lru.Find("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestRemove(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := New(int64(cap), nil)
	lru.Put(k1, String(v1))
	lru.Put(k2, String(v2))
	lru.Put(k3, String(v3))

	if _, ok := lru.Find("key1"); ok || lru.Len() != 2 {
		t.Fatalf("Removeoldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Put("key1", String("123456"))
	lru.Put("k2", String("k2"))
	lru.Put("k3", String("k3"))
	lru.Put("k4", String("k4"))

	expect := []string{"key1", "k2"}

	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}
}

func TestPut(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Put("key", String("1"))
	lru.Put("key", String("111"))

	if lru.nowBytes != int64(len("key")+len("111")) {
		t.Fatal("expected 6 but got", lru.nowBytes)
	}
}