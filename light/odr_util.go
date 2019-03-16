
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450096611725312>


package light

import (
	"bytes"
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

var sha3_nil = crypto.Keccak256Hash(nil)

func GetHeaderByNumber(ctx context.Context, odr OdrBackend, number uint64) (*types.Header, error) {
	db := odr.Database()
	hash := rawdb.ReadCanonicalHash(db, number)
	if (hash != common.Hash{}) {
//如果有一个规范散列，也有一个头
		header := rawdb.ReadHeader(db, hash, number)
		if header == nil {
			panic("Canonical hash present but header not found")
		}
		return header, nil
	}

	var (
		chtCount, sectionHeadNum uint64
		sectionHead              common.Hash
	)
	if odr.ChtIndexer() != nil {
		chtCount, sectionHeadNum, sectionHead = odr.ChtIndexer().Sections()
		canonicalHash := rawdb.ReadCanonicalHash(db, sectionHeadNum)
//如果将cht作为可信检查点注入，那么我们还没有规范散列，因此我们也接受零散列。
		for chtCount > 0 && canonicalHash != sectionHead && canonicalHash != (common.Hash{}) {
			chtCount--
			if chtCount > 0 {
				sectionHeadNum = chtCount*odr.IndexerConfig().ChtSize - 1
				sectionHead = odr.ChtIndexer().SectionHead(chtCount - 1)
				canonicalHash = rawdb.ReadCanonicalHash(db, sectionHeadNum)
			}
		}
	}
	if number >= chtCount*odr.IndexerConfig().ChtSize {
		return nil, ErrNoTrustedCht
	}
	r := &ChtRequest{ChtRoot: GetChtRoot(db, chtCount-1, sectionHead), ChtNum: chtCount - 1, BlockNum: number, Config: odr.IndexerConfig()}
	if err := odr.Retrieve(ctx, r); err != nil {
		return nil, err
	}
	return r.Header, nil
}

func GetCanonicalHash(ctx context.Context, odr OdrBackend, number uint64) (common.Hash, error) {
	hash := rawdb.ReadCanonicalHash(odr.Database(), number)
	if (hash != common.Hash{}) {
		return hash, nil
	}
	header, err := GetHeaderByNumber(ctx, odr, number)
	if header != nil {
		return header.Hash(), nil
	}
	return common.Hash{}, err
}

//getBodyrlp在rlp编码中检索块体（事务和uncles）。
func GetBodyRLP(ctx context.Context, odr OdrBackend, hash common.Hash, number uint64) (rlp.RawValue, error) {
	if data := rawdb.ReadBodyRLP(odr.Database(), hash, number); data != nil {
		return data, nil
	}
	r := &BlockRequest{Hash: hash, Number: number}
	if err := odr.Retrieve(ctx, r); err != nil {
		return nil, err
	} else {
		return r.Rlp, nil
	}
}

//getBody检索与
//搞砸。
func GetBody(ctx context.Context, odr OdrBackend, hash common.Hash, number uint64) (*types.Body, error) {
	data, err := GetBodyRLP(ctx, odr, hash, number)
	if err != nil {
		return nil, err
	}
	body := new(types.Body)
	if err := rlp.Decode(bytes.NewReader(data), body); err != nil {
		return nil, err
	}
	return body, nil
}

//GetBlock检索与哈希对应的整个块，并对其进行组装
//从存储的标题和正文返回。
func GetBlock(ctx context.Context, odr OdrBackend, hash common.Hash, number uint64) (*types.Block, error) {
//检索块标题和正文内容
	header := rawdb.ReadHeader(odr.Database(), hash, number)
	if header == nil {
		return nil, ErrNoHeader
	}
	body, err := GetBody(ctx, odr, hash, number)
	if err != nil {
		return nil, err
	}
//重新组装阀块并返回
	return types.NewBlockWithHeader(header).WithBody(body.Transactions, body.Uncles), nil
}

//GetBlockReceipts检索由包含的事务生成的收据
//在由散列给出的块中。
func GetBlockReceipts(ctx context.Context, odr OdrBackend, hash common.Hash, number uint64) (types.Receipts, error) {
//从磁盘或网络检索可能不完整的收据
	receipts := rawdb.ReadReceipts(odr.Database(), hash, number)
	if receipts == nil {
		r := &ReceiptsRequest{Hash: hash, Number: number}
		if err := odr.Retrieve(ctx, r); err != nil {
			return nil, err
		}
		receipts = r.Receipts
	}
//如果收据不完整，请填写派生字段
	if len(receipts) > 0 && receipts[0].TxHash == (common.Hash{}) {
		block, err := GetBlock(ctx, odr, hash, number)
		if err != nil {
			return nil, err
		}
		genesis := rawdb.ReadCanonicalHash(odr.Database(), 0)
		config := rawdb.ReadChainConfig(odr.Database(), genesis)

		if err := core.SetReceiptsData(config, block, receipts); err != nil {
			return nil, err
		}
		rawdb.WriteReceipts(odr.Database(), hash, number, receipts)
	}
	return receipts, nil
}

//GetBlockLogs检索包含在
//由散列给出的块。
func GetBlockLogs(ctx context.Context, odr OdrBackend, hash common.Hash, number uint64) ([][]*types.Log, error) {
//从磁盘或网络检索可能不完整的收据
	receipts := rawdb.ReadReceipts(odr.Database(), hash, number)
	if receipts == nil {
		r := &ReceiptsRequest{Hash: hash, Number: number}
		if err := odr.Retrieve(ctx, r); err != nil {
			return nil, err
		}
		receipts = r.Receipts
	}
//返回日志，而不导出收据上的任何计算字段
	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

//GetBloomBits检索属于给定位索引和节索引的一批压缩BloomBits向量
func GetBloomBits(ctx context.Context, odr OdrBackend, bitIdx uint, sectionIdxList []uint64) ([][]byte, error) {
	var (
		db      = odr.Database()
		result  = make([][]byte, len(sectionIdxList))
		reqList []uint64
		reqIdx  []int
	)

	var (
		bloomTrieCount, sectionHeadNum uint64
		sectionHead                    common.Hash
	)
	if odr.BloomTrieIndexer() != nil {
		bloomTrieCount, sectionHeadNum, sectionHead = odr.BloomTrieIndexer().Sections()
		canonicalHash := rawdb.ReadCanonicalHash(db, sectionHeadNum)
//如果将Bloomtrie作为受信任的检查点注入，那么我们还没有规范的哈希，因此我们也接受零哈希。
		for bloomTrieCount > 0 && canonicalHash != sectionHead && canonicalHash != (common.Hash{}) {
			bloomTrieCount--
			if bloomTrieCount > 0 {
				sectionHeadNum = bloomTrieCount*odr.IndexerConfig().BloomTrieSize - 1
				sectionHead = odr.BloomTrieIndexer().SectionHead(bloomTrieCount - 1)
				canonicalHash = rawdb.ReadCanonicalHash(db, sectionHeadNum)
			}
		}
	}

	for i, sectionIdx := range sectionIdxList {
		sectionHead := rawdb.ReadCanonicalHash(db, (sectionIdx+1)*odr.IndexerConfig().BloomSize-1)
//如果没有为该节头编号存储规范散列，我们仍将查找
//一个零分区头的条目（如果我们不知道，我们也用零分区头存储它
//在检索时）
		bloomBits, err := rawdb.ReadBloomBits(db, bitIdx, sectionIdx, sectionHead)
		if err == nil {
			result[i] = bloomBits
		} else {
//todo（rjl493456442）将sectionindex转换为bloomtrie相对索引
			if sectionIdx >= bloomTrieCount {
				return nil, ErrNoTrustedBloomTrie
			}
			reqList = append(reqList, sectionIdx)
			reqIdx = append(reqIdx, i)
		}
	}
	if reqList == nil {
		return result, nil
	}

	r := &BloomRequest{BloomTrieRoot: GetBloomTrieRoot(db, bloomTrieCount-1, sectionHead), BloomTrieNum: bloomTrieCount - 1,
		BitIdx: bitIdx, SectionIndexList: reqList, Config: odr.IndexerConfig()}
	if err := odr.Retrieve(ctx, r); err != nil {
		return nil, err
	} else {
		for i, idx := range reqIdx {
			result[idx] = r.BloomBits[i]
		}
		return result, nil
	}
}

