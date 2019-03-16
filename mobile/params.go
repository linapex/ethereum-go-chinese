
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450101716193280>


//包含params包中的所有包装器。

package geth

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/p2p/discv5"
	"github.com/ethereum/go-ethereum/params"
)

//MainNetGenesis返回用于主以太坊网络的JSON规范。它
//实际上是空的，因为它默认为硬编码的二进制genesis块。
func MainnetGenesis() string {
	return ""
}

//TestNetGenesis返回用于以太坊测试网络的JSON规范。
func TestnetGenesis() string {
	enc, err := json.Marshal(core.DefaultTestnetGenesisBlock())
	if err != nil {
		panic(err)
	}
	return string(enc)
}

//RinkebyGenesis返回用于Rinkeby测试网络的JSON规范
func RinkebyGenesis() string {
	enc, err := json.Marshal(core.DefaultRinkebyGenesisBlock())
	if err != nil {
		panic(err)
	}
	return string(enc)
}

//FoundationBootnodes返回所操作的p2p引导节点的enode URL
//通过运行V5发现协议的基础。
func FoundationBootnodes() *Enodes {
	nodes := &Enodes{nodes: make([]*discv5.Node, len(params.DiscoveryV5Bootnodes))}
	for i, url := range params.DiscoveryV5Bootnodes {
		nodes.nodes[i] = discv5.MustParseNode(url)
	}
	return nodes
}

