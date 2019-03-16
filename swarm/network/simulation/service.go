
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450114710147072>


package simulation

import (
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
)

//服务在特定节点上按名称返回单个服务
//提供ID。
func (s *Simulation) Service(name string, id enode.ID) node.Service {
	simNode, ok := s.Net.GetNode(id).Node.(*adapters.SimNode)
	if !ok {
		return nil
	}
	services := simNode.ServiceMap()
	if len(services) == 0 {
		return nil
	}
	return services[name]
}

//RandomService按名称返回
//随机选择的向上的节点。
func (s *Simulation) RandomService(name string) node.Service {
	n := s.Net.GetRandomUpNode().Node.(*adapters.SimNode)
	if n == nil {
		return nil
	}
	return n.Service(name)
}

//服务返回具有所提供名称的所有服务
//从向上的节点。
func (s *Simulation) Services(name string) (services map[enode.ID]node.Service) {
	nodes := s.Net.GetNodes()
	services = make(map[enode.ID]node.Service)
	for _, node := range nodes {
		if !node.Up {
			continue
		}
		simNode, ok := node.Node.(*adapters.SimNode)
		if !ok {
			continue
		}
		services[node.ID()] = simNode.Service(name)
	}
	return services
}

