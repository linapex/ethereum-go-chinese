
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450121920155648>


//包测试实现了以太坊JSON测试的执行。
package tests

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

//块测试检查整个块的处理。
type BlockTest struct {
	json btJSON
}

//unmashaljson实现json.unmasheler接口。
func (t *BlockTest) UnmarshalJSON(in []byte) error {
	return json.Unmarshal(in, &t.json)
}

type btJSON struct {
	Blocks     []btBlock             `json:"blocks"`
	Genesis    btHeader              `json:"genesisBlockHeader"`
	Pre        core.GenesisAlloc     `json:"pre"`
	Post       core.GenesisAlloc     `json:"postState"`
	BestBlock  common.UnprefixedHash `json:"lastblockhash"`
	Network    string                `json:"network"`
	SealEngine string                `json:"sealEngine"`
}

type btBlock struct {
	BlockHeader  *btHeader
	Rlp          string
	UncleHeaders []*btHeader
}

//go:生成gencodec-type btheader-field override btheadermarshaling-out gen_btheader.go

type btHeader struct {
	Bloom            types.Bloom
	Coinbase         common.Address
	MixHash          common.Hash
	Nonce            types.BlockNonce
	Number           *big.Int
	Hash             common.Hash
	ParentHash       common.Hash
	ReceiptTrie      common.Hash
	StateRoot        common.Hash
	TransactionsTrie common.Hash
	UncleHash        common.Hash
	ExtraData        []byte
	Difficulty       *big.Int
	GasLimit         uint64
	GasUsed          uint64
	Timestamp        *big.Int
}

type btHeaderMarshaling struct {
	ExtraData  hexutil.Bytes
	Number     *math.HexOrDecimal256
	Difficulty *math.HexOrDecimal256
	GasLimit   math.HexOrDecimal64
	GasUsed    math.HexOrDecimal64
	Timestamp  *math.HexOrDecimal256
}

func (t *BlockTest) Run() error {
	config, ok := Forks[t.json.Network]
	if !ok {
		return UnsupportedForkError{t.json.Network}
	}

//导入预帐户并构建测试Genesis块和状态根
	db := ethdb.NewMemDatabase()
	gblock, err := t.genesis(config).Commit(db)
	if err != nil {
		return err
	}
	if gblock.Hash() != t.json.Genesis.Hash {
		return fmt.Errorf("genesis block hash doesn't match test: computed=%x, test=%x", gblock.Hash().Bytes()[:6], t.json.Genesis.Hash[:6])
	}
	if gblock.Root() != t.json.Genesis.StateRoot {
		return fmt.Errorf("genesis block state root does not match test: computed=%x, test=%x", gblock.Root().Bytes()[:6], t.json.Genesis.StateRoot[:6])
	}
	var engine consensus.Engine
	if t.json.SealEngine == "NoProof" {
		engine = ethash.NewFaker()
	} else {
		engine = ethash.NewShared()
	}
	chain, err := core.NewBlockChain(db, &core.CacheConfig{TrieCleanLimit: 0}, config, engine, vm.Config{}, nil)
	if err != nil {
		return err
	}
	defer chain.Stop()

	validBlocks, err := t.insertBlocks(chain)
	if err != nil {
		return err
	}
	cmlast := chain.CurrentBlock().Hash()
	if common.Hash(t.json.BestBlock) != cmlast {
		return fmt.Errorf("last block hash validation mismatch: want: %x, have: %x", t.json.BestBlock, cmlast)
	}
	newDB, err := chain.State()
	if err != nil {
		return err
	}
	if err = t.validatePostState(newDB); err != nil {
		return fmt.Errorf("post state validation failed: %v", err)
	}
	return t.validateImportedHeaders(chain, validBlocks)
}

func (t *BlockTest) genesis(config *params.ChainConfig) *core.Genesis {
	return &core.Genesis{
		Config:     config,
		Nonce:      t.json.Genesis.Nonce.Uint64(),
		Timestamp:  t.json.Genesis.Timestamp.Uint64(),
		ParentHash: t.json.Genesis.ParentHash,
		ExtraData:  t.json.Genesis.ExtraData,
		GasLimit:   t.json.Genesis.GasLimit,
		GasUsed:    t.json.Genesis.GasUsed,
		Difficulty: t.json.Genesis.Difficulty,
		Mixhash:    t.json.Genesis.MixHash,
		Coinbase:   t.json.Genesis.Coinbase,
		Alloc:      t.json.Pre,
	}
}

/*参见https://github.com/ethereum/tests/wiki/blockback-tests-ii

   一个块是否有效有点微妙，它是通过存在
   BlockHeader、Transactions和UncleHeaders字段。如果他们不见了，这个街区就是
   无效，我们必须确认我们不接受它。

   由于一些测试混合了有效和无效的块，我们需要为每个块检查这一点。

   如果一个块是无效的，它不一定会通过测试，如果它的无效性是
   期望我们忽略它并继续处理，然后验证
   后状态。
**/

func (t *BlockTest) insertBlocks(blockchain *core.BlockChain) ([]btBlock, error) {
	validBlocks := make([]btBlock, 0)
//插入将执行所有事务的测试块
	for _, b := range t.json.Blocks {
		cb, err := b.decode()
		if err != nil {
			if b.BlockHeader == nil {
continue //好-块应该无效，继续下一个块
			} else {
				return nil, fmt.Errorf("Block RLP decoding failed when expected to succeed: %v", err)
			}
		}
//RLP解码成功，尝试插入链：
		blocks := types.Blocks{cb}
		i, err := blockchain.InsertChain(blocks)
		if err != nil {
			if b.BlockHeader == nil {
continue //好-块应该无效，继续下一个块
			} else {
				return nil, fmt.Errorf("Block #%v insertion into chain failed: %v", blocks[i].Number(), err)
			}
		}
		if b.BlockHeader == nil {
			return nil, fmt.Errorf("Block insertion should have failed")
		}

//通过检查测试文件json的所有值来验证rlp解码
		if err = validateHeader(b.BlockHeader, cb.Header()); err != nil {
			return nil, fmt.Errorf("Deserialised block header validation failed: %v", err)
		}
		validBlocks = append(validBlocks, b)
	}
	return validBlocks, nil
}

func validateHeader(h *btHeader, h2 *types.Header) error {
	if h.Bloom != h2.Bloom {
		return fmt.Errorf("Bloom: want: %x have: %x", h.Bloom, h2.Bloom)
	}
	if h.Coinbase != h2.Coinbase {
		return fmt.Errorf("Coinbase: want: %x have: %x", h.Coinbase, h2.Coinbase)
	}
	if h.MixHash != h2.MixDigest {
		return fmt.Errorf("MixHash: want: %x have: %x", h.MixHash, h2.MixDigest)
	}
	if h.Nonce != h2.Nonce {
		return fmt.Errorf("Nonce: want: %x have: %x", h.Nonce, h2.Nonce)
	}
	if h.Number.Cmp(h2.Number) != 0 {
		return fmt.Errorf("Number: want: %v have: %v", h.Number, h2.Number)
	}
	if h.ParentHash != h2.ParentHash {
		return fmt.Errorf("Parent hash: want: %x have: %x", h.ParentHash, h2.ParentHash)
	}
	if h.ReceiptTrie != h2.ReceiptHash {
		return fmt.Errorf("Receipt hash: want: %x have: %x", h.ReceiptTrie, h2.ReceiptHash)
	}
	if h.TransactionsTrie != h2.TxHash {
		return fmt.Errorf("Tx hash: want: %x have: %x", h.TransactionsTrie, h2.TxHash)
	}
	if h.StateRoot != h2.Root {
		return fmt.Errorf("State hash: want: %x have: %x", h.StateRoot, h2.Root)
	}
	if h.UncleHash != h2.UncleHash {
		return fmt.Errorf("Uncle hash: want: %x have: %x", h.UncleHash, h2.UncleHash)
	}
	if !bytes.Equal(h.ExtraData, h2.Extra) {
		return fmt.Errorf("Extra data: want: %x have: %x", h.ExtraData, h2.Extra)
	}
	if h.Difficulty.Cmp(h2.Difficulty) != 0 {
		return fmt.Errorf("Difficulty: want: %v have: %v", h.Difficulty, h2.Difficulty)
	}
	if h.GasLimit != h2.GasLimit {
		return fmt.Errorf("GasLimit: want: %d have: %d", h.GasLimit, h2.GasLimit)
	}
	if h.GasUsed != h2.GasUsed {
		return fmt.Errorf("GasUsed: want: %d have: %d", h.GasUsed, h2.GasUsed)
	}
	if h.Timestamp.Cmp(h2.Time) != 0 {
		return fmt.Errorf("Timestamp: want: %v have: %v", h.Timestamp, h2.Time)
	}
	return nil
}

func (t *BlockTest) validatePostState(statedb *state.StateDB) error {
//根据状态数据库中的内容验证测试文件中的状态后帐户
	for addr, acct := range t.json.Post {
//地址由其他字段间接验证，因为它是db键
		code2 := statedb.GetCode(addr)
		balance2 := statedb.GetBalance(addr)
		nonce2 := statedb.GetNonce(addr)
		if !bytes.Equal(code2, acct.Code) {
			return fmt.Errorf("account code mismatch for addr: %s want: %v have: %s", addr, acct.Code, hex.EncodeToString(code2))
		}
		if balance2.Cmp(acct.Balance) != 0 {
			return fmt.Errorf("account balance mismatch for addr: %s, want: %d, have: %d", addr, acct.Balance, balance2)
		}
		if nonce2 != acct.Nonce {
			return fmt.Errorf("account nonce mismatch for addr: %s want: %d have: %d", addr, acct.Nonce, nonce2)
		}
	}
	return nil
}

func (t *BlockTest) validateImportedHeaders(cm *core.BlockChain, validBlocks []btBlock) error {
//通过哈希验证块头时获取常量查找（某些测试有许多块）
	bmap := make(map[common.Hash]btBlock, len(t.json.Blocks))
	for _, b := range validBlocks {
		bmap[b.BlockHeader.Hash] = b
	}
//从头部向后迭代块并验证导入的
//头与测试文件。有些测试有REORG，我们导入
//逐块，因此我们只能在
//所有块都已由区块链处理，因为它们可能不是
//成为最长链的一部分，直到导入最后一个块。
	for b := cm.CurrentBlock(); b != nil && b.NumberU64() != 0; b = cm.GetBlockByHash(b.Header().ParentHash) {
		if err := validateHeader(bmap[b.Hash()].BlockHeader, b.Header()); err != nil {
			return fmt.Errorf("Imported block header validation failed: %v", err)
		}
	}
	return nil
}

func (bb *btBlock) decode() (*types.Block, error) {
	data, err := hexutil.Decode(bb.Rlp)
	if err != nil {
		return nil, err
	}
	var b types.Block
	err = rlp.DecodeBytes(data, &b)
	return &b, err
}

