
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450106745163776>


package simulations

import (
	"errors"
	"strings"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

var (
	ErrNodeNotFound = errors.New("node not found")
)

//ConnectToLastNode将节点与提供的节点ID连接起来
//到上一个节点，并避免连接到自身。
//它在构建链网络拓扑结构时很有用
//当网络动态添加和删除节点时。
func (net *Network) ConnectToLastNode(id enode.ID) (err error) {
	ids := net.getUpNodeIDs()
	l := len(ids)
	if l < 2 {
		return nil
	}
	last := ids[l-1]
	if last == id {
		last = ids[l-2]
	}
	return net.connect(last, id)
}

//connecttorandomnode将节点与提供的nodeid连接起来
//向上的随机节点发送。
func (net *Network) ConnectToRandomNode(id enode.ID) (err error) {
	selected := net.GetRandomUpNode(id)
	if selected == nil {
		return ErrNodeNotFound
	}
	return net.connect(selected.ID(), id)
}

//ConnectNodesFull将所有节点连接到另一个。
//它在网络中提供了完整的连接
//这应该是很少需要的。
func (net *Network) ConnectNodesFull(ids []enode.ID) (err error) {
	if ids == nil {
		ids = net.getUpNodeIDs()
	}
	for i, lid := range ids {
		for _, rid := range ids[i+1:] {
			if err = net.connect(lid, rid); err != nil {
				return err
			}
		}
	}
	return nil
}

//connectnodeschain连接链拓扑中的所有节点。
//如果ids参数为nil，则所有打开的节点都将被连接。
func (net *Network) ConnectNodesChain(ids []enode.ID) (err error) {
	if ids == nil {
		ids = net.getUpNodeIDs()
	}
	l := len(ids)
	for i := 0; i < l-1; i++ {
		if err := net.connect(ids[i], ids[i+1]); err != nil {
			return err
		}
	}
	return nil
}

//ConnectNodesRing连接环拓扑中的所有节点。
//如果ids参数为nil，则所有打开的节点都将被连接。
func (net *Network) ConnectNodesRing(ids []enode.ID) (err error) {
	if ids == nil {
		ids = net.getUpNodeIDs()
	}
	l := len(ids)
	if l < 2 {
		return nil
	}
	if err := net.ConnectNodesChain(ids); err != nil {
		return err
	}
	return net.connect(ids[l-1], ids[0])
}

//connectnodestar将所有节点连接到星形拓扑中
//如果ids参数为nil，则所有打开的节点都将被连接。
func (net *Network) ConnectNodesStar(ids []enode.ID, center enode.ID) (err error) {
	if ids == nil {
		ids = net.getUpNodeIDs()
	}
	for _, id := range ids {
		if center == id {
			continue
		}
		if err := net.connect(center, id); err != nil {
			return err
		}
	}
	return nil
}

//连接连接两个节点，但忽略已连接的错误。
func (net *Network) connect(oneID, otherID enode.ID) error {
	return ignoreAlreadyConnectedErr(net.Connect(oneID, otherID))
}

func ignoreAlreadyConnectedErr(err error) error {
	if err == nil || strings.Contains(err.Error(), "already connected") {
		return nil
	}
	return err
}

