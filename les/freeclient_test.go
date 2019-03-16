
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450093935759360>


//package light实现可按需检索的状态和链对象
//对于以太坊Light客户端。
package les

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
	"github.com/ethereum/go-ethereum/ethdb"
)

func TestFreeClientPoolL10C100(t *testing.T) {
	testFreeClientPool(t, 10, 100)
}

func TestFreeClientPoolL40C200(t *testing.T) {
	testFreeClientPool(t, 40, 200)
}

func TestFreeClientPoolL100C300(t *testing.T) {
	testFreeClientPool(t, 100, 300)
}

const testFreeClientPoolTicks = 500000

func testFreeClientPool(t *testing.T, connLimit, clientCount int) {
	var (
		clock     mclock.Simulated
		db        = ethdb.NewMemDatabase()
		pool      = newFreeClientPool(db, connLimit, 10000, &clock)
		connected = make([]bool, clientCount)
		connTicks = make([]int, clientCount)
		disconnCh = make(chan int, clientCount)
	)
	peerId := func(i int) string {
		return fmt.Sprintf("test peer #%d", i)
	}
	disconnFn := func(i int) func() {
		return func() {
			disconnCh <- i
		}
	}

//池应接受达到其连接限制的新对等方
	for i := 0; i < connLimit; i++ {
		if pool.connect(peerId(i), disconnFn(i)) {
			connected[i] = true
		} else {
			t.Fatalf("Test peer #%d rejected", i)
		}
	}
//因为所有被接受的同龄人都是新的，不应该被淘汰，所以下一个同龄人应该被拒绝。
	if pool.connect(peerId(connLimit), disconnFn(connLimit)) {
		connected[connLimit] = true
		t.Fatalf("Peer accepted over connected limit")
	}

//随机连接和断开对等端，希望在端部具有相似的总连接时间
	for tickCounter := 0; tickCounter < testFreeClientPoolTicks; tickCounter++ {
		clock.Run(1 * time.Second)

		i := rand.Intn(clientCount)
		if connected[i] {
			pool.disconnect(peerId(i))
			connected[i] = false
			connTicks[i] += tickCounter
		} else {
			if pool.connect(peerId(i), disconnFn(i)) {
				connected[i] = true
				connTicks[i] -= tickCounter
			}
		}
	pollDisconnects:
		for {
			select {
			case i := <-disconnCh:
				pool.disconnect(peerId(i))
				if connected[i] {
					connTicks[i] += tickCounter
					connected[i] = false
				}
			default:
				break pollDisconnects
			}
		}
	}

	expTicks := testFreeClientPoolTicks * connLimit / clientCount
	expMin := expTicks - expTicks/10
	expMax := expTicks + expTicks/10

//检查对等节点的总连接时间是否在预期范围内
	for i, c := range connected {
		if c {
			connTicks[i] += testFreeClientPoolTicks
		}
		if connTicks[i] < expMin || connTicks[i] > expMax {
			t.Errorf("Total connected time of test node #%d (%d) outside expected range (%d to %d)", i, connTicks[i], expMin, expMax)
		}
	}

//现在应接受以前未知的对等机
	if !pool.connect("newPeer", func() {}) {
		t.Fatalf("Previously unknown peer rejected")
	}

//关闭并重新启动池
	pool.stop()
	pool = newFreeClientPool(db, connLimit, 10000, &clock)

//尝试连接所有已知对等端（应填写connlimit）
	for i := 0; i < clientCount; i++ {
		pool.connect(peerId(i), func() {})
	}
//期望池记住已知节点，并将其中一个节点踢出以接受新节点
	if !pool.connect("newPeer2", func() {}) {
		t.Errorf("Previously unknown peer rejected after restarting pool")
	}
	pool.stop()
}

