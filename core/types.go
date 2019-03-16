
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450081759694848>


package core

import (
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

//validator是定义块验证标准的接口。它
//只负责验证块内容，因为头验证是
//由特定的共识引擎完成。
//
type Validator interface {
//validateBody验证给定块的内容。
	ValidateBody(block *types.Block) error

//validateState验证给定的statedb，以及可选的收据和
//使用的气体。
	ValidateState(block, parent *types.Block, state *state.StateDB, receipts types.Receipts, usedGas uint64) error
}

//处理器是使用给定初始状态处理块的接口。
//
//process接受要处理的块和statedb，在该块上
//初始状态是基于的。它应该返回生成的收据，金额
//过程中使用的气体，如果有任何内部规则，则返回错误
//失败。
type Processor interface {
	Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, error)
}

