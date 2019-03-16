
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450115427373056>


package intervals

import (
	"testing"

	"github.com/ethereum/go-ethereum/swarm/state"
)

//测试内存存储测试InMemoryStore的基本功能。
func TestInmemoryStore(t *testing.T) {
	testStore(t, state.NewInmemoryStore())
}

//test store是一个帮助函数，用于测试各种存储实现。
func testStore(t *testing.T, s state.Store) {
	key1 := "key1"
	i1 := NewIntervals(0)
	i1.Add(10, 20)
	if err := s.Put(key1, i1); err != nil {
		t.Fatal(err)
	}
	i := &Intervals{}
	err := s.Get(key1, i)
	if err != nil {
		t.Fatal(err)
	}
	if i.String() != i1.String() {
		t.Errorf("expected interval %s, got %s", i1, i)
	}

	key2 := "key2"
	i2 := NewIntervals(0)
	i2.Add(10, 20)
	if err := s.Put(key2, i2); err != nil {
		t.Fatal(err)
	}
	err = s.Get(key2, i)
	if err != nil {
		t.Fatal(err)
	}
	if i.String() != i2.String() {
		t.Errorf("expected interval %s, got %s", i2, i)
	}

	if err := s.Delete(key1); err != nil {
		t.Fatal(err)
	}
	if err := s.Get(key1, i); err != state.ErrNotFound {
		t.Errorf("expected error %v, got %s", state.ErrNotFound, err)
	}
	if err := s.Get(key2, i); err != nil {
		t.Errorf("expected error %v, got %s", nil, err)
	}

	if err := s.Delete(key2); err != nil {
		t.Fatal(err)
	}
	if err := s.Get(key2, i); err != state.ErrNotFound {
		t.Errorf("expected error %v, got %s", state.ErrNotFound, err)
	}
}

