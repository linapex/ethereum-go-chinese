
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450120351485952>


package storage

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/swarm/log"
)

func newTestMemStore() *MemStore {
	storeparams := NewDefaultStoreParams()
	return NewMemStore(storeparams, nil)
}

func testMemStoreRandom(n int, chunksize int64, t *testing.T) {
	m := newTestMemStore()
	defer m.Close()
	testStoreRandom(m, n, chunksize, t)
}

func testMemStoreCorrect(n int, chunksize int64, t *testing.T) {
	m := newTestMemStore()
	defer m.Close()
	testStoreCorrect(m, n, chunksize, t)
}

func TestMemStoreRandom_1(t *testing.T) {
	testMemStoreRandom(1, 0, t)
}

func TestMemStoreCorrect_1(t *testing.T) {
	testMemStoreCorrect(1, 4104, t)
}

func TestMemStoreRandom_1k(t *testing.T) {
	testMemStoreRandom(1000, 0, t)
}

func TestMemStoreCorrect_1k(t *testing.T) {
	testMemStoreCorrect(100, 4096, t)
}

func TestMemStoreNotFound(t *testing.T) {
	m := newTestMemStore()
	defer m.Close()

	_, err := m.Get(context.TODO(), ZeroAddr)
	if err != ErrChunkNotFound {
		t.Errorf("Expected ErrChunkNotFound, got %v", err)
	}
}

func benchmarkMemStorePut(n int, processors int, chunksize int64, b *testing.B) {
	m := newTestMemStore()
	defer m.Close()
	benchmarkStorePut(m, n, chunksize, b)
}

func benchmarkMemStoreGet(n int, processors int, chunksize int64, b *testing.B) {
	m := newTestMemStore()
	defer m.Close()
	benchmarkStoreGet(m, n, chunksize, b)
}

func BenchmarkMemStorePut_1_500(b *testing.B) {
	benchmarkMemStorePut(500, 1, 4096, b)
}

func BenchmarkMemStorePut_8_500(b *testing.B) {
	benchmarkMemStorePut(500, 8, 4096, b)
}

func BenchmarkMemStoreGet_1_500(b *testing.B) {
	benchmarkMemStoreGet(500, 1, 4096, b)
}

func BenchmarkMemStoreGet_8_500(b *testing.B) {
	benchmarkMemStoreGet(500, 8, 4096, b)
}

func TestMemStoreAndLDBStore(t *testing.T) {
	ldb, cleanup := newLDBStore(t)
	ldb.setCapacity(4000)
	defer cleanup()

	cacheCap := 200
	memStore := NewMemStore(NewStoreParams(4000, 200, nil, nil), nil)

	tests := []struct {
n         int   //要推送到memstore的块数
chunkSize int64 //块的大小（在Swarm-4096中默认）
	}{
		{
			n:         1,
			chunkSize: 4096,
		},
		{
			n:         101,
			chunkSize: 4096,
		},
		{
			n:         501,
			chunkSize: 4096,
		},
		{
			n:         1100,
			chunkSize: 4096,
		},
	}

	for i, tt := range tests {
		log.Info("running test", "idx", i, "tt", tt)
		var chunks []Chunk

		for i := 0; i < tt.n; i++ {
			c := GenerateRandomChunk(tt.chunkSize)
			chunks = append(chunks, c)
		}

		for i := 0; i < tt.n; i++ {
			err := ldb.Put(context.TODO(), chunks[i])
			if err != nil {
				t.Fatal(err)
			}
			err = memStore.Put(context.TODO(), chunks[i])
			if err != nil {
				t.Fatal(err)
			}

			if got := memStore.cache.Len(); got > cacheCap {
				t.Fatalf("expected to get cache capacity less than %v, but got %v", cacheCap, got)
			}

		}

		for i := 0; i < tt.n; i++ {
			_, err := memStore.Get(context.TODO(), chunks[i].Address())
			if err != nil {
				if err == ErrChunkNotFound {
					_, err := ldb.Get(context.TODO(), chunks[i].Address())
					if err != nil {
						t.Fatalf("couldn't get chunk %v from ldb, got error: %v", i, err)
					}
				} else {
					t.Fatalf("got error from memstore: %v", err)
				}
			}
		}
	}
}

