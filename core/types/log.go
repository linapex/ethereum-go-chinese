
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450081424150528>


package types

import (
	"io"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
)

//go:生成gencodec-type log-field override logmarshaling-out gen_log_json.go

//日志表示合同日志事件。这些事件由日志操作码生成，并且
//由节点存储/索引。
type Log struct {
//共识领域：
//产生事件的合同的地址
	Address common.Address `json:"address" gencodec:"required"`
//合同提供的主题列表。
	Topics []common.Hash `json:"topics" gencodec:"required"`
//由合同提供，通常为ABI编码
	Data []byte `json:"data" gencodec:"required"`

//派生字段。这些字段由节点填充
//但没有达成共识。
//包含事务的块
	BlockNumber uint64 `json:"blockNumber"`
//事务的哈希
	TxHash common.Hash `json:"transactionHash" gencodec:"required"`
//块中事务的索引
	TxIndex uint `json:"transactionIndex" gencodec:"required"`
//包含事务的块的哈希
	BlockHash common.Hash `json:"blockHash"`
//块中日志的索引
	Index uint `json:"logIndex" gencodec:"required"`

//如果由于链重组而还原此日志，则删除的字段为真。
//如果通过筛选查询接收日志，则必须注意此字段。
	Removed bool `json:"removed"`
}

type logMarshaling struct {
	Data        hexutil.Bytes
	BlockNumber hexutil.Uint64
	TxIndex     hexutil.Uint
	Index       hexutil.Uint
}

type rlpLog struct {
	Address common.Address
	Topics  []common.Hash
	Data    []byte
}

type rlpStorageLog struct {
	Address     common.Address
	Topics      []common.Hash
	Data        []byte
	BlockNumber uint64
	TxHash      common.Hash
	TxIndex     uint
	BlockHash   common.Hash
	Index       uint
}

//encoderlp实现rlp.encoder。
func (l *Log) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpLog{Address: l.Address, Topics: l.Topics, Data: l.Data})
}

//decoderlp实现rlp.decoder。
func (l *Log) DecodeRLP(s *rlp.Stream) error {
	var dec rlpLog
	err := s.Decode(&dec)
	if err == nil {
		l.Address, l.Topics, l.Data = dec.Address, dec.Topics, dec.Data
	}
	return err
}

//logForStorage是一个围绕一个日志的包装器，它扁平化并解析
//包含非共识字段的日志。
type LogForStorage Log

//encoderlp实现rlp.encoder。
func (l *LogForStorage) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, rlpStorageLog{
		Address:     l.Address,
		Topics:      l.Topics,
		Data:        l.Data,
		BlockNumber: l.BlockNumber,
		TxHash:      l.TxHash,
		TxIndex:     l.TxIndex,
		BlockHash:   l.BlockHash,
		Index:       l.Index,
	})
}

//decoderlp实现rlp.decoder。
func (l *LogForStorage) DecodeRLP(s *rlp.Stream) error {
	var dec rlpStorageLog
	err := s.Decode(&dec)
	if err == nil {
		*l = LogForStorage{
			Address:     dec.Address,
			Topics:      dec.Topics,
			Data:        dec.Data,
			BlockNumber: dec.BlockNumber,
			TxHash:      dec.TxHash,
			TxIndex:     dec.TxIndex,
			BlockHash:   dec.BlockHash,
			Index:       dec.Index,
		}
	}
	return err
}

