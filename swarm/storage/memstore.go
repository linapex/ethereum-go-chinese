
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450120301154304>


//包块哈希的内存存储层

package storage

import (
	"context"

	lru "github.com/hashicorp/golang-lru"
)

type MemStore struct {
	cache    *lru.Cache
	disabled bool
}

//newmemstore正在实例化memstore缓存，以保留所有经常请求的缓存
//“cache”lru缓存中的块。
func NewMemStore(params *StoreParams, _ *LDBStore) (m *MemStore) {
	if params.CacheCapacity == 0 {
		return &MemStore{
			disabled: true,
		}
	}

	c, err := lru.New(int(params.CacheCapacity))
	if err != nil {
		panic(err)
	}

	return &MemStore{
		cache: c,
	}
}

func (m *MemStore) Get(_ context.Context, addr Address) (Chunk, error) {
	if m.disabled {
		return nil, ErrChunkNotFound
	}

	c, ok := m.cache.Get(string(addr))
	if !ok {
		return nil, ErrChunkNotFound
	}
	return c.(Chunk), nil
}

func (m *MemStore) Put(_ context.Context, c Chunk) error {
	if m.disabled {
		return nil
	}

	m.cache.Add(string(c.Address()), c)
	return nil
}

func (m *MemStore) setCapacity(n int) {
	if n <= 0 {
		m.disabled = true
	} else {
		c, err := lru.New(n)
		if err != nil {
			panic(err)
		}

		*m = MemStore{
			cache: c,
		}
	}
}

func (s *MemStore) Close() {}

