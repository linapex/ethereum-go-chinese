
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450088449609728>


package downloader

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

//测试链参数。
var (
	testKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddress = crypto.PubkeyToAddress(testKey.PublicKey)
	testDB      = ethdb.NewMemDatabase()
	testGenesis = core.GenesisBlockForTesting(testDB, testAddress, big.NewInt(1000000000))
)

//所有测试链的通用前缀：
var testChainBase = newTestChain(blockCacheItems+200, testGenesis)

//基链顶部的不同叉：
var testChainForkLightA, testChainForkLightB, testChainForkHeavy *testChain

func init() {
	var forkLen = int(MaxForkAncestry + 50)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() { testChainForkLightA = testChainBase.makeFork(forkLen, false, 1); wg.Done() }()
	go func() { testChainForkLightB = testChainBase.makeFork(forkLen, false, 2); wg.Done() }()
	go func() { testChainForkHeavy = testChainBase.makeFork(forkLen, true, 3); wg.Done() }()
	wg.Wait()
}

type testChain struct {
	genesis  *types.Block
	chain    []common.Hash
	headerm  map[common.Hash]*types.Header
	blockm   map[common.Hash]*types.Block
	receiptm map[common.Hash][]*types.Receipt
	tdm      map[common.Hash]*big.Int
}

//NexTestCu链创建给定长度的块链。
func newTestChain(length int, genesis *types.Block) *testChain {
	tc := new(testChain).copy(length)
	tc.genesis = genesis
	tc.chain = append(tc.chain, genesis.Hash())
	tc.headerm[tc.genesis.Hash()] = tc.genesis.Header()
	tc.tdm[tc.genesis.Hash()] = tc.genesis.Difficulty()
	tc.blockm[tc.genesis.Hash()] = tc.genesis
	tc.generate(length-1, 0, genesis, false)
	return tc
}

//makefork在测试链的顶部创建一个fork。
func (tc *testChain) makeFork(length int, heavy bool, seed byte) *testChain {
	fork := tc.copy(tc.len() + length)
	fork.generate(length, seed, tc.headBlock(), heavy)
	return fork
}

//缩短创建一个给定长度的链副本。如果它是恐慌的
//长度大于可用块的数量。
func (tc *testChain) shorten(length int) *testChain {
	if length > tc.len() {
		panic(fmt.Errorf("can't shorten test chain to %d blocks, it's only %d blocks long", length, tc.len()))
	}
	return tc.copy(length)
}

func (tc *testChain) copy(newlen int) *testChain {
	cpy := &testChain{
		genesis:  tc.genesis,
		headerm:  make(map[common.Hash]*types.Header, newlen),
		blockm:   make(map[common.Hash]*types.Block, newlen),
		receiptm: make(map[common.Hash][]*types.Receipt, newlen),
		tdm:      make(map[common.Hash]*big.Int, newlen),
	}
	for i := 0; i < len(tc.chain) && i < newlen; i++ {
		hash := tc.chain[i]
		cpy.chain = append(cpy.chain, tc.chain[i])
		cpy.tdm[hash] = tc.tdm[hash]
		cpy.blockm[hash] = tc.blockm[hash]
		cpy.headerm[hash] = tc.headerm[hash]
		cpy.receiptm[hash] = tc.receiptm[hash]
	}
	return cpy
}

//生成创建一个由n个块组成的链，从父块开始并包含父块。
//返回的哈希链是有序的head->parent。此外，每22个街区
//包含一个事务，每隔5分钟一个叔叔以允许测试正确的块
//重新组装。
func (tc *testChain) generate(n int, seed byte, parent *types.Block, heavy bool) {
//开始：=时间。现在（））
//defer func（）fmt.printf（“测试链在%v中生成”，time.since（start））（）

	blocks, receipts := core.GenerateChain(params.TestChainConfig, parent, ethash.NewFaker(), testDB, n, func(i int, block *core.BlockGen) {
		block.SetCoinbase(common.Address{seed})
//If a heavy chain is requested, delay blocks to raise difficulty
		if heavy {
			block.OffsetTime(-1)
		}
//将事务包括到矿工中，以使块更有趣。
		if parent == tc.genesis && i%22 == 0 {
			signer := types.MakeSigner(params.TestChainConfig, block.Number())
			tx, err := types.SignTx(types.NewTransaction(block.TxNonce(testAddress), common.Address{seed}, big.NewInt(1000), params.TxGas, nil, nil), signer, testKey)
			if err != nil {
				panic(err)
			}
			block.AddTx(tx)
		}
//如果区块编号是5的倍数，请在区块中添加一个奖金叔叔。
		if i > 0 && i%5 == 0 {
			block.AddUncle(&types.Header{
				ParentHash: block.PrevBlock(i - 1).Hash(),
				Number:     big.NewInt(block.Number().Int64() - 1),
			})
		}
	})

//将块链转换为哈希链和头/块映射
	td := new(big.Int).Set(tc.td(parent.Hash()))
	for i, b := range blocks {
		td := td.Add(td, b.Difficulty())
		hash := b.Hash()
		tc.chain = append(tc.chain, hash)
		tc.blockm[hash] = b
		tc.headerm[hash] = b.Header()
		tc.receiptm[hash] = receipts[i]
		tc.tdm[hash] = new(big.Int).Set(td)
	}
}

//len返回链中的块总数。
func (tc *testChain) len() int {
	return len(tc.chain)
}

//上架返回链条头部。
func (tc *testChain) headBlock() *types.Block {
	return tc.blockm[tc.chain[len(tc.chain)-1]]
}

//td返回给定块的总难度。
func (tc *testChain) td(hash common.Hash) *big.Int {
	return tc.tdm[hash]
}

//HeadersByHash从给定哈希按升序返回头。
func (tc *testChain) headersByHash(origin common.Hash, amount int, skip int) []*types.Header {
	num, _ := tc.hashToNumber(origin)
	return tc.headersByNumber(num, amount, skip)
}

//HeadersByNumber从给定的数字以升序返回标题。
func (tc *testChain) headersByNumber(origin uint64, amount int, skip int) []*types.Header {
	result := make([]*types.Header, 0, amount)
	for num := origin; num < uint64(len(tc.chain)) && len(result) < amount; num += uint64(skip) + 1 {
		if header, ok := tc.headerm[tc.chain[int(num)]]; ok {
			result = append(result, header)
		}
	}
	return result
}

//收据返回给定块散列的收据。
func (tc *testChain) receipts(hashes []common.Hash) [][]*types.Receipt {
	results := make([][]*types.Receipt, 0, len(hashes))
	for _, hash := range hashes {
		if receipt, ok := tc.receiptm[hash]; ok {
			results = append(results, receipt)
		}
	}
	return results
}

//bodies返回给定块散列的块体。
func (tc *testChain) bodies(hashes []common.Hash) ([][]*types.Transaction, [][]*types.Header) {
	transactions := make([][]*types.Transaction, 0, len(hashes))
	uncles := make([][]*types.Header, 0, len(hashes))
	for _, hash := range hashes {
		if block, ok := tc.blockm[hash]; ok {
			transactions = append(transactions, block.Transactions())
			uncles = append(uncles, block.Uncles())
		}
	}
	return transactions, uncles
}

func (tc *testChain) hashToNumber(target common.Hash) (uint64, bool) {
	for num, hash := range tc.chain {
		if hash == target {
			return uint64(num), true
		}
	}
	return 0, false
}

