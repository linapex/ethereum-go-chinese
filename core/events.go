
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450078714630144>


package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//当一批交易进入交易池时，将过账newtxsevent。
type NewTxsEvent struct{ Txs []*types.Transaction }

//PendingLogSevent在挖掘前发布，并通知挂起的日志。
type PendingLogsEvent struct {
	Logs []*types.Log
}

//当块被导入时，将发布NewMinedBlockEvent。
type NewMinedBlockEvent struct{ Block *types.Block }

//当发生REORG时，会发布REMOVEDLogsevent
type RemovedLogsEvent struct{ Logs []*types.Log }

type ChainEvent struct {
	Block *types.Block
	Hash  common.Hash
	Logs  []*types.Log
}

type ChainSideEvent struct {
	Block *types.Block
}

type ChainHeadEvent struct{ Block *types.Block }

