
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450079624794112>


//包RAWDB包含低级别数据库访问器的集合。
package rawdb

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/metrics"
)

//下面的字段定义低级数据库模式前缀。
var (
//databaseverisionkey跟踪当前数据库版本。
	databaseVerisionKey = []byte("DatabaseVersion")

//HeadHeaderKey跟踪最新的已知头散列。
	headHeaderKey = []byte("LastHeader")

//headblockkey跟踪最新的已知完整块哈希。
	headBlockKey = []byte("LastBlock")

//HeadFastBlockKey在快速同步期间跟踪最新的已知不完整块的哈希。
	headFastBlockKey = []byte("LastFast")

//FastTrieProgressKey跟踪在快速同步期间导入的Trie条目数。
	fastTrieProgressKey = []byte("TrieSync")

//数据项前缀（使用单字节避免混合数据类型，避免使用“i”，用于索引）。
headerPrefix       = []byte("h") //headerPrefix+num（uint64 big endian）+hash->header
headerTDSuffix     = []byte("t") //headerPrefix+num（uint64 big endian）+hash+headerTsuffix->td
headerHashSuffix   = []byte("n") //headerPrefix+num（uint64 big endian）+headerHashSuffix->hash
headerNumberPrefix = []byte("H") //headerNumberPrefix+hash->num（uint64 big endian）

blockBodyPrefix     = []byte("b") //blockbodyprefix+num（uint64 big endian）+hash->block body
blockReceiptsPrefix = []byte("r") //blockReceiptsPrefix+num（uint64 big endian）+hash->block receipts

txLookupPrefix  = []byte("l") //txlookupprefix+hash->交易/收据查找元数据
bloomBitsPrefix = []byte("B") //bloombitsprefix+bit（uint16 big endian）+section（uint64 big endian）+hash->bloom位

preimagePrefix = []byte("secure-key-")      //preimagePrefix + hash -> preimage
configPrefix   = []byte("ethereum-config-") //数据库的配置前缀

//链索引前缀（使用'i`+单字节以避免混合数据类型）。
BloomBitsIndexPrefix = []byte("iB") //BloomBitsIndexPrefix是跟踪其进展的链表索引器的数据表。

	preimageCounter    = metrics.NewRegisteredCounter("db/preimage/total", nil)
	preimageHitCounter = metrics.NewRegisteredCounter("db/preimage/hits", nil)
)

//txLookupEntry是一个位置元数据，用于帮助查找
//只给出散列值的交易或收据。
type TxLookupEntry struct {
	BlockHash  common.Hash
	BlockIndex uint64
	Index      uint64
}

//encodeBlockNumber将块编号编码为big endian uint64
func encodeBlockNumber(number uint64) []byte {
	enc := make([]byte, 8)
	binary.BigEndian.PutUint64(enc, number)
	return enc
}

//headerkey=headerprefix+num（uint64 big endian）+哈希
func headerKey(number uint64, hash common.Hash) []byte {
	return append(append(headerPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

//headertdkey=headerprefix+num（uint64 big endian）+hash+headertdsuffix
func headerTDKey(number uint64, hash common.Hash) []byte {
	return append(headerKey(number, hash), headerTDSuffix...)
}

//headerhashkey=headerprefix+num（uint64 big endian）+headerhashsuffix
func headerHashKey(number uint64) []byte {
	return append(append(headerPrefix, encodeBlockNumber(number)...), headerHashSuffix...)
}

//headerNumberKey=headerNumberPrefix+hash
func headerNumberKey(hash common.Hash) []byte {
	return append(headerNumberPrefix, hash.Bytes()...)
}

//blockBodyKey = blockBodyPrefix + num (uint64 big endian) + hash
func blockBodyKey(number uint64, hash common.Hash) []byte {
	return append(append(blockBodyPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

//blockReceiptskey=blockReceiptsPrefix+num（uint64 big endian）+哈希
func blockReceiptsKey(number uint64, hash common.Hash) []byte {
	return append(append(blockReceiptsPrefix, encodeBlockNumber(number)...), hash.Bytes()...)
}

//txLookupKey=txLookupPrefix+哈希
func txLookupKey(hash common.Hash) []byte {
	return append(txLookupPrefix, hash.Bytes()...)
}

//bloomBitsKey = bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash
func bloomBitsKey(bit uint, section uint64, hash common.Hash) []byte {
	key := append(append(bloomBitsPrefix, make([]byte, 10)...), hash.Bytes()...)

	binary.BigEndian.PutUint16(key[1:], uint16(bit))
	binary.BigEndian.PutUint64(key[3:], section)

	return key
}

//preImageKey=preImagePrefix+哈希
func preimageKey(hash common.Hash) []byte {
	return append(preimagePrefix, hash.Bytes()...)
}

//configkey=configPrefix+哈希
func configKey(hash common.Hash) []byte {
	return append(configPrefix, hash.Bytes()...)
}

