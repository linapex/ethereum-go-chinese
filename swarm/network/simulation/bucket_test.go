
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450114211024896>


package simulation

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
)

//TestServiceBucket使用子测试测试所有Bucket功能。
//它通过向两个节点的bucket中添加项来构造两个节点的模拟。
//在servicefunc构造函数中，然后通过setnodeitem。测试upnodesitems
//通过停止一个节点并验证其项的可用性来完成。
func TestServiceBucket(t *testing.T) {
	testKey := "Key"
	testValue := "Value"

	sim := New(map[string]ServiceFunc{
		"noop": func(ctx *adapters.ServiceContext, b *sync.Map) (node.Service, func(), error) {
			b.Store(testKey, testValue+ctx.Config.ID.String())
			return newNoopService(), nil, nil
		},
	})
	defer sim.Close()

	id1, err := sim.AddNode()
	if err != nil {
		t.Fatal(err)
	}

	id2, err := sim.AddNode()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("ServiceFunc bucket Store", func(t *testing.T) {
		v, ok := sim.NodeItem(id1, testKey)
		if !ok {
			t.Fatal("bucket item not found")
		}
		s, ok := v.(string)
		if !ok {
			t.Fatal("bucket item value is not string")
		}
		if s != testValue+id1.String() {
			t.Fatalf("expected %q, got %q", testValue+id1.String(), s)
		}

		v, ok = sim.NodeItem(id2, testKey)
		if !ok {
			t.Fatal("bucket item not found")
		}
		s, ok = v.(string)
		if !ok {
			t.Fatal("bucket item value is not string")
		}
		if s != testValue+id2.String() {
			t.Fatalf("expected %q, got %q", testValue+id2.String(), s)
		}
	})

	customKey := "anotherKey"
	customValue := "anotherValue"

	t.Run("SetNodeItem", func(t *testing.T) {
		sim.SetNodeItem(id1, customKey, customValue)

		v, ok := sim.NodeItem(id1, customKey)
		if !ok {
			t.Fatal("bucket item not found")
		}
		s, ok := v.(string)
		if !ok {
			t.Fatal("bucket item value is not string")
		}
		if s != customValue {
			t.Fatalf("expected %q, got %q", customValue, s)
		}

		_, ok = sim.NodeItem(id2, customKey)
		if ok {
			t.Fatal("bucket item should not be found")
		}
	})

	if err := sim.StopNode(id2); err != nil {
		t.Fatal(err)
	}

	t.Run("UpNodesItems", func(t *testing.T) {
		items := sim.UpNodesItems(testKey)

		v, ok := items[id1]
		if !ok {
			t.Errorf("node 1 item not found")
		}
		s, ok := v.(string)
		if !ok {
			t.Fatal("node 1 item value is not string")
		}
		if s != testValue+id1.String() {
			t.Fatalf("expected %q, got %q", testValue+id1.String(), s)
		}

		_, ok = items[id2]
		if ok {
			t.Errorf("node 2 item should not be found")
		}
	})

	t.Run("NodeItems", func(t *testing.T) {
		items := sim.NodesItems(testKey)

		v, ok := items[id1]
		if !ok {
			t.Errorf("node 1 item not found")
		}
		s, ok := v.(string)
		if !ok {
			t.Fatal("node 1 item value is not string")
		}
		if s != testValue+id1.String() {
			t.Fatalf("expected %q, got %q", testValue+id1.String(), s)
		}

		v, ok = items[id2]
		if !ok {
			t.Errorf("node 2 item not found")
		}
		s, ok = v.(string)
		if !ok {
			t.Fatal("node 1 item value is not string")
		}
		if s != testValue+id2.String() {
			t.Fatalf("expected %q, got %q", testValue+id2.String(), s)
		}
	})
}

