
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450079364747264>


package rawdb

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

//Read TxLoopUpTeAccess检索与事务关联的位置元数据
//哈希以允许按哈希检索事务或收据。
func ReadTxLookupEntry(db DatabaseReader, hash common.Hash) (common.Hash, uint64, uint64) {
	data, _ := db.Get(txLookupKey(hash))
	if len(data) == 0 {
		return common.Hash{}, 0, 0
	}
	var entry TxLookupEntry
	if err := rlp.DecodeBytes(data, &entry); err != nil {
		log.Error("Invalid transaction lookup entry RLP", "hash", hash, "err", err)
		return common.Hash{}, 0, 0
	}
	return entry.BlockHash, entry.BlockIndex, entry.Index
}

//WraveTxLoopUpStand为每个事务存储位置元数据。
//一个块，启用基于哈希的事务和收据查找。
func WriteTxLookupEntries(db DatabaseWriter, block *types.Block) {
	for i, tx := range block.Transactions() {
		entry := TxLookupEntry{
			BlockHash:  block.Hash(),
			BlockIndex: block.NumberU64(),
			Index:      uint64(i),
		}
		data, err := rlp.EncodeToBytes(entry)
		if err != nil {
			log.Crit("Failed to encode transaction lookup entry", "err", err)
		}
		if err := db.Put(txLookupKey(tx.Hash()), data); err != nil {
			log.Crit("Failed to store transaction lookup entry", "err", err)
		}
	}
}

//DeleteTxLookupEntry删除与哈希关联的所有事务数据。
func DeleteTxLookupEntry(db DatabaseDeleter, hash common.Hash) {
	db.Delete(txLookupKey(hash))
}

//readTransaction从数据库中检索特定事务，以及
//它添加了位置元数据。
func ReadTransaction(db DatabaseReader, hash common.Hash) (*types.Transaction, common.Hash, uint64, uint64) {
	blockHash, blockNumber, txIndex := ReadTxLookupEntry(db, hash)
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0
	}
	body := ReadBody(db, blockHash, blockNumber)
	if body == nil || len(body.Transactions) <= int(txIndex) {
		log.Error("Transaction referenced missing", "number", blockNumber, "hash", blockHash, "index", txIndex)
		return nil, common.Hash{}, 0, 0
	}
	return body.Transactions[txIndex], blockHash, blockNumber, txIndex
}

//readreceipt从数据库中检索特定的事务回执，以及
//它添加了位置元数据。
func ReadReceipt(db DatabaseReader, hash common.Hash) (*types.Receipt, common.Hash, uint64, uint64) {
	blockHash, blockNumber, receiptIndex := ReadTxLookupEntry(db, hash)
	if blockHash == (common.Hash{}) {
		return nil, common.Hash{}, 0, 0
	}
	receipts := ReadReceipts(db, blockHash, blockNumber)
	if len(receipts) <= int(receiptIndex) {
		log.Error("Receipt refereced missing", "number", blockNumber, "hash", blockHash, "index", receiptIndex)
		return nil, common.Hash{}, 0, 0
	}
	return receipts[receiptIndex], blockHash, blockNumber, receiptIndex
}

//readBloomBits检索属于给定的
//中的节和位索引。
func ReadBloomBits(db DatabaseReader, bit uint, section uint64, head common.Hash) ([]byte, error) {
	return db.Get(bloomBitsKey(bit, section, head))
}

//WriteBloomBits存储属于给定的
//节和位索引。
func WriteBloomBits(db DatabaseWriter, bit uint, section uint64, head common.Hash, bits []byte) {
	if err := db.Put(bloomBitsKey(bit, section, head), bits); err != nil {
		log.Crit("Failed to store bloom bits", "err", err)
	}
}

