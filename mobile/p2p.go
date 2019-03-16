
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450101623918592>


//包含p2p包的包装。

package geth

import (
	"errors"

	"github.com/ethereum/go-ethereum/p2p"
)

//nodeinfo表示已知主机信息的pi简短摘要。
type NodeInfo struct {
	info *p2p.NodeInfo
}

func (ni *NodeInfo) GetID() string              { return ni.info.ID }
func (ni *NodeInfo) GetName() string            { return ni.info.Name }
func (ni *NodeInfo) GetEnode() string           { return ni.info.Enode }
func (ni *NodeInfo) GetIP() string              { return ni.info.IP }
func (ni *NodeInfo) GetDiscoveryPort() int      { return ni.info.Ports.Discovery }
func (ni *NodeInfo) GetListenerPort() int       { return ni.info.Ports.Listener }
func (ni *NodeInfo) GetListenerAddress() string { return ni.info.ListenAddr }
func (ni *NodeInfo) GetProtocols() *Strings {
	protos := []string{}
	for proto := range ni.info.Protocols {
		protos = append(protos, proto)
	}
	return &Strings{protos}
}

//peerinfo表示关于pi-connected peer已知信息的pi简短摘要。
type PeerInfo struct {
	info *p2p.PeerInfo
}

func (pi *PeerInfo) GetID() string            { return pi.info.ID }
func (pi *PeerInfo) GetName() string          { return pi.info.Name }
func (pi *PeerInfo) GetCaps() *Strings        { return &Strings{pi.info.Caps} }
func (pi *PeerInfo) GetLocalAddress() string  { return pi.info.Network.LocalAddress }
func (pi *PeerInfo) GetRemoteAddress() string { return pi.info.Network.RemoteAddress }

//PeerInfos表示关于远程对等机的一部分信息。
type PeerInfos struct {
	infos []*p2p.PeerInfo
}

//SIZE返回切片中的对等信息条目数。
func (pi *PeerInfos) Size() int {
	return len(pi.infos)
}

//GET从切片返回给定索引处的对等信息。
func (pi *PeerInfos) Get(index int) (info *PeerInfo, _ error) {
	if index < 0 || index >= len(pi.infos) {
		return nil, errors.New("index out of bounds")
	}
	return &PeerInfo{pi.infos[index]}, nil
}

