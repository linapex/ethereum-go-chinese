
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450114299105280>


package simulation

import (
	"context"
	"sync"
	"testing"
	"time"
)

//testpeerEvents创建模拟，添加两个节点，
//注册对等事件，连接链中的节点
//并等待连接事件的数量
//被接受。
func TestPeerEvents(t *testing.T) {
	sim := New(noopServiceFuncMap)
	defer sim.Close()

	_, err := sim.AddNodes(2)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	events := sim.PeerEvents(ctx, sim.NodeIDs())

//两个节点->两个连接事件
	expectedEventCount := 2

	var wg sync.WaitGroup
	wg.Add(expectedEventCount)

	go func() {
		for e := range events {
			if e.Error != nil {
				if e.Error == context.Canceled {
					return
				}
				t.Error(e.Error)
				continue
			}
			wg.Done()
		}
	}()

	err = sim.Net.ConnectNodesChain(sim.NodeIDs())
	if err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}

func TestPeerEventsTimeout(t *testing.T) {
	sim := New(noopServiceFuncMap)
	defer sim.Close()

	_, err := sim.AddNodes(2)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	events := sim.PeerEvents(ctx, sim.NodeIDs())

	done := make(chan struct{})
	errC := make(chan error)
	go func() {
		for e := range events {
			if e.Error == context.Canceled {
				return
			}
			if e.Error == context.DeadlineExceeded {
				close(done)
				return
			} else {
				errC <- e.Error
			}
		}
	}()

	select {
	case <-time.After(time.Second):
		t.Fatal("no context deadline received")
	case err := <-errC:
		t.Fatal(err)
	case <-done:
//一切正常，检测到上下文截止时间
	}
}

