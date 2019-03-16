
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450114336854016>


package simulation_test

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
	"github.com/ethereum/go-ethereum/swarm/network"
	"github.com/ethereum/go-ethereum/swarm/network/simulation"
)

//每个节点都可以使用下面的节点bucket关联一个kademlia。
//BucketkeyKademlia键。这允许使用waittillhealthy在
//所有节点的花环都是健康的。
func ExampleSimulation_WaitTillHealthy() {

	sim := simulation.New(map[string]simulation.ServiceFunc{
		"bzz": func(ctx *adapters.ServiceContext, b *sync.Map) (node.Service, func(), error) {
			addr := network.NewAddr(ctx.Config.Node())
			hp := network.NewHiveParams()
			hp.Discovery = false
			config := &network.BzzConfig{
				OverlayAddr:  addr.Over(),
				UnderlayAddr: addr.Under(),
				HiveParams:   hp,
			}
			kad := network.NewKademlia(addr.Over(), network.NewKadParams())
//将kademlia存储在bucketkeykademlia下的节点桶中
//这样就可以通过waittillhealthy方法找到它。
			b.Store(simulation.BucketKeyKademlia, kad)
			return network.NewBzz(config, kad, nil, nil, nil), nil, nil
		},
	})
	defer sim.Close()

	_, err := sim.AddNodesAndConnectRing(10)
	if err != nil {
//正确处理错误…
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	ill, err := sim.WaitTillHealthy(ctx)
	if err != nil {
//检查最新检测到的不健康花环
		for id, kad := range ill {
			fmt.Println("Node", id)
			fmt.Println(kad.String())
		}
//处理错误…
	}

//继续测试

}

//观看模拟网络中的所有对等事件，从一个频道购买接收。
func ExampleSimulation_PeerEvents() {
	sim := simulation.New(nil)
	defer sim.Close()

	events := sim.PeerEvents(context.Background(), sim.NodeIDs())

	go func() {
		for e := range events {
			if e.Error != nil {
				log.Error("peer event", "err", e.Error)
				continue
			}
			log.Info("peer event", "node", e.NodeID, "peer", e.PeerID, "type", e.Event.Type)
		}
	}()
}

//检测节点何时丢弃对等节点。
func ExampleSimulation_PeerEvents_disconnections() {
	sim := simulation.New(nil)
	defer sim.Close()

	disconnections := sim.PeerEvents(
		context.Background(),
		sim.NodeIDs(),
		simulation.NewPeerEventsFilter().Drop(),
	)

	go func() {
		for d := range disconnections {
			if d.Error != nil {
				log.Error("peer drop", "err", d.Error)
				continue
			}
			log.Warn("peer drop", "node", d.NodeID, "peer", d.PeerID)
		}
	}()
}

//观看多种类型的事件或消息。在这种情况下，它们只是不同的
//通过msgcode，但是过滤器也可以为不同的类型或协议设置。
func ExampleSimulation_PeerEvents_multipleFilters() {
	sim := simulation.New(nil)
	defer sim.Close()

	msgs := sim.PeerEvents(
		context.Background(),
		sim.NodeIDs(),
//当收到BZZ消息1和4时，请注意。
		simulation.NewPeerEventsFilter().ReceivedMessages().Protocol("bzz").MsgCode(1),
		simulation.NewPeerEventsFilter().ReceivedMessages().Protocol("bzz").MsgCode(4),
	)

	go func() {
		for m := range msgs {
			if m.Error != nil {
				log.Error("bzz message", "err", m.Error)
				continue
			}
			log.Info("bzz message", "node", m.NodeID, "peer", m.PeerID)
		}
	}()
}

