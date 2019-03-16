
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450113762234368>


package network

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	p2ptest "github.com/ethereum/go-ethereum/p2p/testing"
	"github.com/ethereum/go-ethereum/swarm/state"
)

func newHiveTester(t *testing.T, params *HiveParams, n int, store state.Store) (*bzzTester, *Hive) {
//设置
addr := RandomAddr() //测试的对等地址
	to := NewKademlia(addr.OAddr, NewKadParams())
pp := NewHive(params, to, store) //蜂箱

	return newBzzBaseTester(t, n, addr, DiscoverySpec, pp.Run), pp
}

func TestRegisterAndConnect(t *testing.T) {
	params := NewHiveParams()
	s, pp := newHiveTester(t, params, 1, nil)

	node := s.Nodes[0]
	raddr := NewAddr(node)
	pp.Register(raddr)

//启动配置单元并等待连接
	err := pp.Start(s.Server)
	if err != nil {
		t.Fatal(err)
	}
	defer pp.Stop()
//检索和广播
	err = s.TestDisconnected(&p2ptest.Disconnect{
		Peer:  s.Nodes[0].ID(),
		Error: nil,
	})

	if err == nil || err.Error() != "timed out waiting for peers to disconnect" {
		t.Fatalf("expected peer to connect")
	}
}

func TestHiveStatePersistance(t *testing.T) {
	log.SetOutput(os.Stdout)

	dir, err := ioutil.TempDir("", "hive_test_store")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

store, err := state.NewDBStore(dir) //用空的dbstore启动配置单元
	if err != nil {
		t.Fatal(err)
	}

	params := NewHiveParams()
	s, pp := newHiveTester(t, params, 5, store)

	peers := make(map[string]bool)
	for _, node := range s.Nodes {
		raddr := NewAddr(node)
		pp.Register(raddr)
		peers[raddr.String()] = true
	}

//启动配置单元并等待连接
	err = pp.Start(s.Server)
	if err != nil {
		t.Fatal(err)
	}
	pp.Stop()
	store.Close()

persistedStore, err := state.NewDBStore(dir) //用空的dbstore启动配置单元
	if err != nil {
		t.Fatal(err)
	}

	s1, pp := newHiveTester(t, params, 1, persistedStore)

//启动配置单元并等待连接

	pp.Start(s1.Server)
	i := 0
	pp.Kademlia.EachAddr(nil, 256, func(addr *BzzAddr, po int) bool {
		delete(peers, addr.String())
		i++
		return true
	})
	if i != 5 {
		t.Errorf("invalid number of entries: got %v, want %v", i, 5)
	}
	if len(peers) != 0 {
		t.Fatalf("%d peers left over: %v", len(peers), peers)
	}
}

