
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450114483654656>


package simulation

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/swarm/network"
)

//BucketkeyKademlia是用于存储Kademlia的密钥
//特定节点的实例，通常位于servicefunc函数内部。
var BucketKeyKademlia BucketKey = "kademlia"

//等到所有花环的健康都是真实的时候，Waittillhealthy才会停止。
//如果错误不为零，则返回发现不健康的卡德米利亚地图。
//TODO:检查Kademilia深度计算逻辑更改后的正确性
func (s *Simulation) WaitTillHealthy(ctx context.Context) (ill map[enode.ID]*network.Kademlia, err error) {
//准备Peerpot地图检查卡德米利亚健康
	var ppmap map[string]*network.PeerPot
	kademlias := s.kademlias()
	addrs := make([][]byte, 0, len(kademlias))
//要验证所有花冠是否具有相同的参数
	for _, k := range kademlias {
		addrs = append(addrs, k.BaseAddr())
	}
	ppmap = network.NewPeerPotMap(s.neighbourhoodSize, addrs)

//在检查文件之前，在每个节点上等待健康的kademlia
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	ill = make(map[enode.ID]*network.Kademlia)
	for {
		select {
		case <-ctx.Done():
			return ill, ctx.Err()
		case <-ticker.C:
			for k := range ill {
				delete(ill, k)
			}
			log.Debug("kademlia health check", "addr count", len(addrs))
			for id, k := range kademlias {
//此节点的Peerpot
				addr := common.Bytes2Hex(k.BaseAddr())
				pp := ppmap[addr]
//调用健康的RPC
				h := k.Healthy(pp)
//打印信息
				log.Debug(k.String())
				log.Debug("kademlia", "connectNN", h.ConnectNN, "knowNN", h.KnowNN)
				log.Debug("kademlia", "health", h.ConnectNN && h.KnowNN, "addr", hex.EncodeToString(k.BaseAddr()), "node", id)
				log.Debug("kademlia", "ill condition", !h.ConnectNN, "addr", hex.EncodeToString(k.BaseAddr()), "node", id)
				if !h.ConnectNN {
					ill[id] = k
				}
			}
			if len(ill) == 0 {
				return nil, nil
			}
		}
	}
}

//Kademlias返回设置的所有Kademlia实例
//在模拟桶中。
func (s *Simulation) kademlias() (ks map[enode.ID]*network.Kademlia) {
	items := s.UpNodesItems(BucketKeyKademlia)
	ks = make(map[enode.ID]*network.Kademlia, len(items))
	for id, v := range items {
		k, ok := v.(*network.Kademlia)
		if !ok {
			continue
		}
		ks[id] = k
	}
	return ks
}

