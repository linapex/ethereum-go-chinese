
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450101904936960>


//包含core/types包中的所有包装器。

package geth

import (
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
)

//日志表示合同日志事件。这些事件由日志生成
//操作码并由节点存储/索引。
type Log struct {
	log *types.Log
}

func (l *Log) GetAddress() *Address  { return &Address{l.log.Address} }
func (l *Log) GetTopics() *Hashes    { return &Hashes{l.log.Topics} }
func (l *Log) GetData() []byte       { return l.log.Data }
func (l *Log) GetBlockNumber() int64 { return int64(l.log.BlockNumber) }
func (l *Log) GetTxHash() *Hash      { return &Hash{l.log.TxHash} }
func (l *Log) GetTxIndex() int       { return int(l.log.TxIndex) }
func (l *Log) GetBlockHash() *Hash   { return &Hash{l.log.BlockHash} }
func (l *Log) GetIndex() int         { return int(l.log.Index) }

//日志表示VM日志的一部分。
type Logs struct{ logs []*types.Log }

//SIZE返回切片中的日志数。
func (l *Logs) Size() int {
	return len(l.logs)
}

//get返回切片中给定索引处的日志。
func (l *Logs) Get(index int) (log *Log, _ error) {
	if index < 0 || index >= len(l.logs) {
		return nil, errors.New("index out of bounds")
	}
	return &Log{l.logs[index]}, nil
}

