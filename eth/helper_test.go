
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450089301053440>


//此文件包含一些共享测试功能，多个共享测试功能
//正在测试的不同文件和模块。

package eth

import (
	"crypto/ecdsa"
	"crypto/rand"
	"math/big"
	"sort"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
)

var (
	testBankKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testBank       = crypto.PubkeyToAddress(testBankKey.PublicKey)
)

//NewTestProtocolManager为测试目的创建了一个新的协议管理器，
//已知已知的块数和潜在的通知
//不同事件的频道。
func newTestProtocolManager(mode downloader.SyncMode, blocks int, generator func(int, *core.BlockGen), newtx chan<- []*types.Transaction) (*ProtocolManager, *ethdb.MemDatabase, error) {
	var (
		evmux  = new(event.TypeMux)
		engine = ethash.NewFaker()
		db     = ethdb.NewMemDatabase()
		gspec  = &core.Genesis{
			Config: params.TestChainConfig,
			Alloc:  core.GenesisAlloc{testBank: {Balance: big.NewInt(1000000)}},
		}
		genesis       = gspec.MustCommit(db)
		blockchain, _ = core.NewBlockChain(db, nil, gspec.Config, engine, vm.Config{}, nil)
	)
	chain, _ := core.GenerateChain(gspec.Config, genesis, ethash.NewFaker(), db, blocks, generator)
	if _, err := blockchain.InsertChain(chain); err != nil {
		panic(err)
	}

	pm, err := NewProtocolManager(gspec.Config, mode, DefaultConfig.NetworkId, evmux, &testTxPool{added: newtx}, engine, blockchain, db, nil)
	if err != nil {
		return nil, nil, err
	}
	pm.Start(1000)
	return pm, db, nil
}

//newTestProtocolManager必须为测试目的创建新的协议管理器，
//已知已知的块数和潜在的通知
//不同事件的频道。如果出现错误，构造函数强制-
//测试失败。
func newTestProtocolManagerMust(t *testing.T, mode downloader.SyncMode, blocks int, generator func(int, *core.BlockGen), newtx chan<- []*types.Transaction) (*ProtocolManager, *ethdb.MemDatabase) {
	pm, db, err := newTestProtocolManager(mode, blocks, generator, newtx)
	if err != nil {
		t.Fatalf("Failed to create protocol manager: %v", err)
	}
	return pm, db
}

//testxtpool是一个用于测试的假助手事务池
type testTxPool struct {
	txFeed event.Feed
pool   []*types.Transaction        //收集所有交易
added  chan<- []*types.Transaction //新事务的通知通道

lock sync.RWMutex //保护事务池
}

//AddRemotes向池追加一批事务，并通知
//如果添加通道为非零，则侦听器
func (p *testTxPool) AddRemotes(txs []*types.Transaction) []error {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.pool = append(p.pool, txs...)
	if p.added != nil {
		p.added <- txs
	}
	return make([]error, len(txs))
}

//挂起返回池已知的所有事务
func (p *testTxPool) Pending() (map[common.Address]types.Transactions, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	batches := make(map[common.Address]types.Transactions)
	for _, tx := range p.pool {
		from, _ := types.Sender(types.HomesteadSigner{}, tx)
		batches[from] = append(batches[from], tx)
	}
	for _, batch := range batches {
		sort.Sort(types.TxByNonce(batch))
	}
	return batches, nil
}

func (p *testTxPool) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return p.txFeed.Subscribe(ch)
}

//NewTestTransaction创建新的虚拟事务。
func newTestTransaction(from *ecdsa.PrivateKey, nonce uint64, datasize int) *types.Transaction {
	tx := types.NewTransaction(nonce, common.Address{}, big.NewInt(0), 100000, big.NewInt(0), make([]byte, datasize))
	tx, _ = types.SignTx(tx, types.HomesteadSigner{}, from)
	return tx
}

//testpeer是允许测试直接网络调用的模拟对等机。
type testPeer struct {
net p2p.MsgReadWriter //模拟远程消息传递的网络层读写器
app *p2p.MsgPipeRW    //应用层读写器模拟本地端
	*peer
}

//newtestpeer创建在给定的协议管理器上注册的新对等。
func newTestPeer(name string, version int, pm *ProtocolManager, shake bool) (*testPeer, <-chan error) {
//创建消息管道以通过
	app, net := p2p.MsgPipe()

//生成随机ID并创建对等机
	var id enode.ID
	rand.Read(id[:])

	peer := pm.newPeer(version, p2p.NewPeer(id, name, nil), net)

//在新线程上启动对等机
	errc := make(chan error, 1)
	go func() {
		select {
		case pm.newPeerCh <- peer:
			errc <- pm.handle(peer)
		case <-pm.quitSync:
			errc <- p2p.DiscQuitting
		}
	}()
	tp := &testPeer{app: app, net: net, peer: peer}
//执行任何隐式请求的握手并返回
	if shake {
		var (
			genesis = pm.blockchain.Genesis()
			head    = pm.blockchain.CurrentHeader()
			td      = pm.blockchain.GetTd(head.Hash(), head.Number.Uint64())
		)
		tp.handshake(nil, td, head.Hash(), genesis.Hash())
	}
	return tp, errc
}

//握手模拟一个简单的握手，它期望
//我们在本地模拟的远程端。
func (p *testPeer) handshake(t *testing.T, td *big.Int, head common.Hash, genesis common.Hash) {
	msg := &statusData{
		ProtocolVersion: uint32(p.version),
		NetworkId:       DefaultConfig.NetworkId,
		TD:              td,
		CurrentBlock:    head,
		GenesisBlock:    genesis,
	}
	if err := p2p.ExpectMsg(p.app, StatusMsg, msg); err != nil {
		t.Fatalf("status recv: %v", err)
	}
	if err := p2p.Send(p.app, StatusMsg, msg); err != nil {
		t.Fatalf("status send: %v", err)
	}
}

//CLOSE终止对等端的本地端，通知远程协议
//终止经理。
func (p *testPeer) close() {
	p.app.Close()
}

