
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450101254819840>


//包含节点包中支持客户端节点的所有包装器
//移动平台管理。

package geth

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethstats"
	"github.com/ethereum/go-ethereum/internal/debug"
	"github.com/ethereum/go-ethereum/les"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/params"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"
)

//nodeconfig表示对geth进行微调的配置值集合。
//嵌入到移动进程中的节点。可用值是
//Go Ethereum提供的整个API，以减少维护面和开发
//复杂性。
type NodeConfig struct {
//引导节点用于建立与网络其余部分的连接。
	BootstrapNodes *Enodes

//MaxPeers是可以连接的最大对等数。如果这是
//设置为零，则只有配置的静态和受信任对等方可以连接。
	MaxPeers int

//ethereumenabled指定节点是否应运行ethereum协议。
	EthereumEnabled bool

//ethereumnetworkid是以太坊协议用于
//决定是否接受远程对等。
EthereumNetworkID int64 //UIT64实际上是，但Java不能处理…

//以太坊是用来为区块链播种的Genesis JSON。安
//空的genesis状态相当于使用mainnet的状态。
	EthereumGenesis string

//ethereumdatabasecache是以MB为单位分配给数据库缓存的系统内存。
//始终保留至少16MB。
	EthereumDatabaseCache int

//ethereumnetstats是一个netstats连接字符串，用于报告各种
//链、事务和节点状态到监控服务器。
//
//格式为“nodename:secret@host:port”
	EthereumNetStats string

//WhisperEnabled指定节点是否应运行Whisper协议。
	WhisperEnabled bool

//pprof服务器的侦听地址。
	PprofAddress string
}

//defaultnodeconfig包含默认节点配置值（如果全部使用）
//或者用户指定列表中缺少某些字段。
var defaultNodeConfig = &NodeConfig{
	BootstrapNodes:        FoundationBootnodes(),
	MaxPeers:              25,
	EthereumEnabled:       true,
	EthereumNetworkID:     1,
	EthereumDatabaseCache: 16,
}

//newnodeconfig创建一个新的节点选项集，初始化为默认值。
func NewNodeConfig() *NodeConfig {
	config := *defaultNodeConfig
	return &config
}

//节点表示一个geth-ethereum节点实例。
type Node struct {
	node *node.Node
}

//new node创建和配置新的geth节点。
func NewNode(datadir string, config *NodeConfig) (stack *Node, _ error) {
//如果未指定或未指定部分配置，请使用默认值
	if config == nil {
		config = NewNodeConfig()
	}
	if config.MaxPeers == 0 {
		config.MaxPeers = defaultNodeConfig.MaxPeers
	}
	if config.BootstrapNodes == nil || config.BootstrapNodes.Size() == 0 {
		config.BootstrapNodes = defaultNodeConfig.BootstrapNodes
	}

	if config.PprofAddress != "" {
		debug.StartPProf(config.PprofAddress)
	}

//创建空的网络堆栈
	nodeConf := &node.Config{
		Name:        clientIdentifier,
		Version:     params.VersionWithMeta,
		DataDir:     datadir,
KeyStoreDir: filepath.Join(datadir, "keystore"), //手机不应使用内部密钥库！
		P2P: p2p.Config{
			NoDiscovery:      true,
			DiscoveryV5:      true,
			BootstrapNodesV5: config.BootstrapNodes.nodes,
			ListenAddr:       ":0",
			NAT:              nat.Any(),
			MaxPeers:         config.MaxPeers,
		},
	}
	rawStack, err := node.New(nodeConf)
	if err != nil {
		return nil, err
	}

	debug.Memsize.Add("node", rawStack)

	var genesis *core.Genesis
	if config.EthereumGenesis != "" {
//分析用户提供的Genesis规范（如果不是Mainnet）
		genesis = new(core.Genesis)
		if err := json.Unmarshal([]byte(config.EthereumGenesis), genesis); err != nil {
			return nil, fmt.Errorf("invalid genesis spec: %v", err)
		}
//如果我们有了测试网，那么也要对链配置进行硬编码。
		if config.EthereumGenesis == TestnetGenesis() {
			genesis.Config = params.TestnetChainConfig
			if config.EthereumNetworkID == 1 {
				config.EthereumNetworkID = 3
			}
		}
	}
//如果需要，注册以太坊协议
	if config.EthereumEnabled {
		ethConf := eth.DefaultConfig
		ethConf.Genesis = genesis
		ethConf.SyncMode = downloader.LightSync
		ethConf.NetworkId = uint64(config.EthereumNetworkID)
		ethConf.DatabaseCache = config.EthereumDatabaseCache
		if err := rawStack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
			return les.New(ctx, &ethConf)
		}); err != nil {
			return nil, fmt.Errorf("ethereum init: %v", err)
		}
//如果请求Netstats报告，请执行此操作
		if config.EthereumNetStats != "" {
			if err := rawStack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
				var lesServ *les.LightEthereum
				ctx.Service(&lesServ)

				return ethstats.New(config.EthereumNetStats, nil, lesServ)
			}); err != nil {
				return nil, fmt.Errorf("netstats init: %v", err)
			}
		}
	}
//如有要求，注册窃听协议
	if config.WhisperEnabled {
		if err := rawStack.Register(func(*node.ServiceContext) (node.Service, error) {
			return whisper.New(&whisper.DefaultConfig), nil
		}); err != nil {
			return nil, fmt.Errorf("whisper init: %v", err)
		}
	}
	return &Node{rawStack}, nil
}

//start创建一个活动的p2p节点并开始运行它。
func (n *Node) Start() error {
	return n.node.Start()
}

//stop终止正在运行的节点及其所有服务。如果节点是
//未启动，返回错误。
func (n *Node) Stop() error {
	return n.node.Stop()
}

//getethereumclient检索客户端以访问ethereum子系统。
func (n *Node) GetEthereumClient() (client *EthereumClient, _ error) {
	rpc, err := n.node.Attach()
	if err != nil {
		return nil, err
	}
	return &EthereumClient{ethclient.NewClient(rpc)}, nil
}

//getnodeinfo收集并返回有关主机的已知元数据集合。
func (n *Node) GetNodeInfo() *NodeInfo {
	return &NodeInfo{n.node.Server().NodeInfo()}
}

//GetPeerSinfo返回描述已连接对等端的元数据对象数组。
func (n *Node) GetPeersInfo() *PeerInfos {
	return &PeerInfos{n.node.Server().PeersInfo()}
}

