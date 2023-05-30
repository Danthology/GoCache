package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashMap struct {
	hash     Hash
	replicas int
	keys     []uint32
	nodeMap  map[uint32]string
}

type Hash func([]byte) uint32

func NewHash(replicas int, fh Hash) *HashMap {
	m := &HashMap{
		hash:     fh,
		replicas: replicas,
		keys:     make([]uint32, 0),
		nodeMap:  make(map[uint32]string),
	}

	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// 添加节点
func (this *HashMap) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < this.replicas; i++ {
			hk := this.hash([]byte(strconv.Itoa(i) + key))
			this.keys = append(this.keys, hk)
			this.nodeMap[hk] = key
		}
	}

	sort.Slice(this.keys, func(i int, j int) bool {
		return this.keys[i] < this.keys[j]
	})
}

// 获取节点
func (this *HashMap) Get(key string) string {
	hv := this.hash([]byte(key))
	idx := sort.Search(len(this.keys), func(i int) bool {
		return this.keys[i] >= hv
	})

	return this.nodeMap[this.keys[idx%len(this.keys)]]
}
