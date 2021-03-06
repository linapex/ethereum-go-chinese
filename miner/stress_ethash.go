
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450100525010944>


//+不建

//此文件包含基于ethash共识引擎的矿工压力测试。
package main

import (
	"crypto/ecdsa"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/fdlimit"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
)

func main() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
	fdlimit.Raise(2048)

//生成一批封存资金的账户
	faucets := make([]*ecdsa.PrivateKey, 128)
	for i := 0; i < len(faucets); i++ {
		faucets[i], _ = crypto.GenerateKey()
	}
//预先生成ethash挖掘数据，这样我们就不会比赛了
	ethash.MakeDataset(1, filepath.Join(os.Getenv("HOME"), ".ethash"))

//基于Ropsten配置创建ethash网络
	genesis := makeGenesis(faucets)

	var (
		nodes  []*node.Node
		enodes []*enode.Node
	)
	for i := 0; i < 4; i++ {
//启动节点并等待它启动
		node, err := makeMiner(genesis)
		if err != nil {
			panic(err)
		}
		defer node.Stop()

		for node.Server().NodeInfo().Ports.Listener == 0 {
			time.Sleep(250 * time.Millisecond)
		}
//将节点连接到前面的所有节点
		for _, n := range enodes {
			node.Server().AddPeer(n)
		}
//开始跟踪节点，它是enode
		nodes = append(nodes, node)
		enodes = append(enodes, node.Server().Self())

//注入签名者密钥并用它开始密封
		store := node.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
		if _, err := store.NewAccount(""); err != nil {
			panic(err)
		}
	}
//遍历所有节点并开始使用它们进行签名
	time.Sleep(3 * time.Second)

	for _, node := range nodes {
		var ethereum *eth.Ethereum
		if err := node.Service(&ethereum); err != nil {
			panic(err)
		}
		if err := ethereum.StartMining(1); err != nil {
			panic(err)
		}
	}
	time.Sleep(3 * time.Second)

//开始疯狂地从水龙头注入交易
	nonces := make([]uint64, len(faucets))
	for {
		index := rand.Intn(len(faucets))

//获取相关签名者的访问器
		var ethereum *eth.Ethereum
		if err := nodes[index%len(nodes)].Service(&ethereum); err != nil {
			panic(err)
		}
//创建一个自事务并注入池
		tx, err := types.SignTx(types.NewTransaction(nonces[index], crypto.PubkeyToAddress(faucets[index].PublicKey), new(big.Int), 21000, big.NewInt(100000000000+rand.Int63n(65536)), nil), types.HomesteadSigner{}, faucets[index])
		if err != nil {
			panic(err)
		}
		if err := ethereum.TxPool().AddLocal(tx); err != nil {
			panic(err)
		}
		nonces[index]++

//等一下，如果我们太饱和了
		if pend, _ := ethereum.TxPool().Stats(); pend > 2048 {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

//MakeGenesis基于一些预定义的
//水龙头账户。
func makeGenesis(faucets []*ecdsa.PrivateKey) *core.Genesis {
	genesis := core.DefaultTestnetGenesisBlock()
	genesis.Difficulty = params.MinimumDifficulty
	genesis.GasLimit = 25000000

	genesis.Config.ChainID = big.NewInt(18)
	genesis.Config.EIP150Hash = common.Hash{}

	genesis.Alloc = core.GenesisAlloc{}
	for _, faucet := range faucets {
		genesis.Alloc[crypto.PubkeyToAddress(faucet.PublicKey)] = core.GenesisAccount{
			Balance: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
		}
	}
	return genesis
}

func makeMiner(genesis *core.Genesis) (*node.Node, error) {
//定义以太坊节点的基本配置
	datadir, _ := ioutil.TempDir("", "")

	config := &node.Config{
		Name:    "geth",
		Version: params.Version,
		DataDir: datadir,
		P2P: p2p.Config{
			ListenAddr:  "0.0.0.0:0",
			NoDiscovery: true,
			MaxPeers:    25,
		},
		NoUSB:             true,
		UseLightweightKDF: true,
	}
//启动节点并在其上配置完整的以太坊节点
	stack, err := node.New(config)
	if err != nil {
		return nil, err
	}
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return eth.New(ctx, &eth.Config{
			Genesis:         genesis,
			NetworkId:       genesis.Config.ChainID.Uint64(),
			SyncMode:        downloader.FullSync,
			DatabaseCache:   256,
			DatabaseHandles: 256,
			TxPool:          core.DefaultTxPoolConfig,
			GPO:             eth.DefaultConfig.GPO,
			Ethash:          eth.DefaultConfig.Ethash,
			MinerGasFloor:   genesis.GasLimit * 9 / 10,
			MinerGasCeil:    genesis.GasLimit * 11 / 10,
			MinerGasPrice:   big.NewInt(1),
			MinerRecommit:   time.Second,
		})
	}); err != nil {
		return nil, err
	}
//启动节点，成功返回
	return stack, stack.Start()
}

