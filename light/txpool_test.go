
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450096926298112>


package light

import (
	"context"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

type testTxRelay struct {
	send, discard, mined chan int
}

func (self *testTxRelay) Send(txs types.Transactions) {
	self.send <- len(txs)
}

func (self *testTxRelay) NewHead(head common.Hash, mined []common.Hash, rollback []common.Hash) {
	m := len(mined)
	if m != 0 {
		self.mined <- m
	}
}

func (self *testTxRelay) Discard(hashes []common.Hash) {
	self.discard <- len(hashes)
}

const poolTestTxs = 1000
const poolTestBlocks = 100

//测试TX 0…N-1
var testTx [poolTestTxs]*types.Transaction

//一区前发送的TXS
func sentTx(i int) int {
	return int(math.Pow(float64(i)/float64(poolTestBlocks), 0.9) * poolTestTxs)
}

//一区或之前包含的Txs（minedtx（i）<=sentx（i））。
func minedTx(i int) int {
	return int(math.Pow(float64(i)/float64(poolTestBlocks), 1.1) * poolTestTxs)
}

func txPoolTestChainGen(i int, block *core.BlockGen) {
	s := minedTx(i)
	e := minedTx(i + 1)
	for i := s; i < e; i++ {
		block.AddTx(testTx[i])
	}
}

func TestTxPool(t *testing.T) {
	for i := range testTx {
		testTx[i], _ = types.SignTx(types.NewTransaction(uint64(i), acc1Addr, big.NewInt(10000), params.TxGas, nil, nil), types.HomesteadSigner{}, testBankKey)
	}

	var (
		sdb     = ethdb.NewMemDatabase()
		ldb     = ethdb.NewMemDatabase()
		gspec   = core.Genesis{Alloc: core.GenesisAlloc{testBankAddress: {Balance: testBankFunds}}}
		genesis = gspec.MustCommit(sdb)
	)
	gspec.MustCommit(ldb)
//组装测试环境
	blockchain, _ := core.NewBlockChain(sdb, nil, params.TestChainConfig, ethash.NewFullFaker(), vm.Config{}, nil)
	gchain, _ := core.GenerateChain(params.TestChainConfig, genesis, ethash.NewFaker(), sdb, poolTestBlocks, txPoolTestChainGen)
	if _, err := blockchain.InsertChain(gchain); err != nil {
		panic(err)
	}

	odr := &testOdr{sdb: sdb, ldb: ldb, indexerConfig: TestClientIndexerConfig}
	relay := &testTxRelay{
		send:    make(chan int, 1),
		discard: make(chan int, 1),
		mined:   make(chan int, 1),
	}
	lightchain, _ := NewLightChain(odr, params.TestChainConfig, ethash.NewFullFaker())
	txPermanent = 50
	pool := NewTxPool(params.TestChainConfig, lightchain, relay)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for ii, block := range gchain {
		i := ii + 1
		s := sentTx(i - 1)
		e := sentTx(i)
		for i := s; i < e; i++ {
			pool.Add(ctx, testTx[i])
			got := <-relay.send
			exp := 1
			if got != exp {
				t.Errorf("relay.Send expected len = %d, got %d", exp, got)
			}
		}

		if _, err := lightchain.InsertHeaderChain([]*types.Header{block.Header()}, 1); err != nil {
			panic(err)
		}

		got := <-relay.mined
		exp := minedTx(i) - minedTx(i-1)
		if got != exp {
			t.Errorf("relay.NewHead expected len(mined) = %d, got %d", exp, got)
		}

		exp = 0
		if i > int(txPermanent)+1 {
			exp = minedTx(i-int(txPermanent)-1) - minedTx(i-int(txPermanent)-2)
		}
		if exp != 0 {
			got = <-relay.discard
			if got != exp {
				t.Errorf("relay.Discard expected len = %d, got %d", exp, got)
			}
		}
	}
}

