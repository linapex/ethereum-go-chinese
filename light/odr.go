
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450096498479104>


//package light实现可按需检索的状态和链对象
//对于以太坊Light客户端。
package light

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
)

//noodr是当odr
//不需要服务。
var NoOdr = context.Background()

//如果没有能够为排队请求提供服务的对等方可用，则返回errnopeers。
var ErrNoPeers = errors.New("no suitable peers available")

//odr backend是后端服务的接口，用于处理odr检索类型
type OdrBackend interface {
	Database() ethdb.Database
	ChtIndexer() *core.ChainIndexer
	BloomTrieIndexer() *core.ChainIndexer
	BloomIndexer() *core.ChainIndexer
	Retrieve(ctx context.Context, req OdrRequest) error
	IndexerConfig() *IndexerConfig
}

//ODRRequest是一个用于检索请求的接口
type OdrRequest interface {
	StoreResult(db ethdb.Database)
}

//trieid标识状态或帐户存储trie
type TrieID struct {
	BlockHash, Root common.Hash
	BlockNumber     uint64
	AccKey          []byte
}

//state trieid返回属于某个块的state trie的trieid
//标题。
func StateTrieID(header *types.Header) *TrieID {
	return &TrieID{
		BlockHash:   header.Hash(),
		BlockNumber: header.Number.Uint64(),
		AccKey:      nil,
		Root:        header.Root,
	}
}

//storage trieid返回给定帐户上合同存储trie的trieid
//一个给定的国家的。它还需要trie的根散列
//检查Merkle校样。
func StorageTrieID(state *TrieID, addrHash, root common.Hash) *TrieID {
	return &TrieID{
		BlockHash:   state.BlockHash,
		BlockNumber: state.BlockNumber,
		AccKey:      addrHash[:],
		Root:        root,
	}
}

//trieRequest是状态/存储trie项的ODR请求类型
type TrieRequest struct {
	OdrRequest
	Id    *TrieID
	Key   []byte
	Proof *NodeSet
}

//storeresult将检索到的数据存储在本地数据库中
func (req *TrieRequest) StoreResult(db ethdb.Database) {
	req.Proof.Store(db)
}

//code request是用于检索合同代码的ODR请求类型
type CodeRequest struct {
	OdrRequest
Id   *TrieID //账户参考存储检索
	Hash common.Hash
	Data []byte
}

//storeresult将检索到的数据存储在本地数据库中
func (req *CodeRequest) StoreResult(db ethdb.Database) {
	db.Put(req.Hash[:], req.Data)
}

//BlockRequest是用于检索块体的ODR请求类型
type BlockRequest struct {
	OdrRequest
	Hash   common.Hash
	Number uint64
	Rlp    []byte
}

//storeresult将检索到的数据存储在本地数据库中
func (req *BlockRequest) StoreResult(db ethdb.Database) {
	rawdb.WriteBodyRLP(db, req.Hash, req.Number, req.Rlp)
}

//ReceiptsRequest是用于检索块体的ODR请求类型
type ReceiptsRequest struct {
	OdrRequest
	Hash     common.Hash
	Number   uint64
	Receipts types.Receipts
}

//storeresult将检索到的数据存储在本地数据库中
func (req *ReceiptsRequest) StoreResult(db ethdb.Database) {
	rawdb.WriteReceipts(db, req.Hash, req.Number, req.Receipts)
}

//chtRequest是状态/存储trie项的odr请求类型
type ChtRequest struct {
	OdrRequest
	Config           *IndexerConfig
	ChtNum, BlockNum uint64
	ChtRoot          common.Hash
	Header           *types.Header
	Td               *big.Int
	Proof            *NodeSet
}

//storeresult将检索到的数据存储在本地数据库中
func (req *ChtRequest) StoreResult(db ethdb.Database) {
	hash, num := req.Header.Hash(), req.Header.Number.Uint64()

	rawdb.WriteHeader(db, req.Header)
	rawdb.WriteTd(db, hash, num, req.Td)
	rawdb.WriteCanonicalHash(db, hash, num)
}

//BloomRequest是用于从CHT结构检索Bloom筛选器的ODR请求类型。
type BloomRequest struct {
	OdrRequest
	Config           *IndexerConfig
	BloomTrieNum     uint64
	BitIdx           uint
	SectionIndexList []uint64
	BloomTrieRoot    common.Hash
	BloomBits        [][]byte
	Proofs           *NodeSet
}

//storeresult将检索到的数据存储在本地数据库中
func (req *BloomRequest) StoreResult(db ethdb.Database) {
	for i, sectionIdx := range req.SectionIndexList {
		sectionHead := rawdb.ReadCanonicalHash(db, (sectionIdx+1)*req.Config.BloomTrieSize-1)
//如果没有为此节头编号存储规范散列，我们仍然将其存储在
//一个零分区头的键。如果我们仍然没有规范的
//搞砸。在不太可能的情况下，我们从那以后就检索到了段头散列，我们只检索
//再次从网络中得到位矢量。
		rawdb.WriteBloomBits(db, req.BitIdx, sectionIdx, sectionHead, req.BloomBits[i])
	}
}

