
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450089598849024>


package eth

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

//一旦一个真正的块成功，快速同步就被禁用的测试
//导入到区块链中。
func TestFastSyncDisabling(t *testing.T) {
//创建一个原始协议管理器，检查是否启用了快速同步
	pmEmpty, _ := newTestProtocolManagerMust(t, downloader.FastSync, 0, nil, nil)
	if atomic.LoadUint32(&pmEmpty.fastSync) == 0 {
		t.Fatalf("fast sync disabled on pristine blockchain")
	}
//创建完整的协议管理器，检查是否禁用快速同步
	pmFull, _ := newTestProtocolManagerMust(t, downloader.FastSync, 1024, nil, nil)
	if atomic.LoadUint32(&pmFull.fastSync) == 1 {
		t.Fatalf("fast sync not disabled on non-empty blockchain")
	}
//同步两个对等机
	io1, io2 := p2p.MsgPipe()

	go pmFull.handle(pmFull.newPeer(63, p2p.NewPeer(enode.ID{}, "empty", nil), io2))
	go pmEmpty.handle(pmEmpty.newPeer(63, p2p.NewPeer(enode.ID{}, "full", nil), io1))

	time.Sleep(250 * time.Millisecond)
	pmEmpty.synchronise(pmEmpty.peers.BestPeer())

//检查是否禁用了快速同步
	if atomic.LoadUint32(&pmEmpty.fastSync) == 1 {
		t.Fatalf("fast sync not disabled after successful synchronisation")
	}
}

