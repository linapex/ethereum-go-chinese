
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450077896740864>


package core

import (
	"fmt"

	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

//blockvalidator负责验证块头、uncles和
//已处理状态。
//
//BlockValidator实现验证程序。
type BlockValidator struct {
config *params.ChainConfig //链配置选项
bc     *BlockChain         //规范区块链
engine consensus.Engine    //用于验证的共识引擎
}

//new block validator返回一个可安全重用的新块验证程序
func NewBlockValidator(config *params.ChainConfig, blockchain *BlockChain, engine consensus.Engine) *BlockValidator {
	validator := &BlockValidator{
		config: config,
		engine: engine,
		bc:     blockchain,
	}
	return validator
}

//validateBody验证给定块的叔叔并验证该块
//头的事务和叔叔根。假定标题已经
//此时已验证。
func (v *BlockValidator) ValidateBody(block *types.Block) error {
//检查块是否已知，如果不知道，它是否可链接
	if v.bc.HasBlockAndState(block.Hash(), block.NumberU64()) {
		return ErrKnownBlock
	}
//此时知道头的有效性，检查叔叔和事务
	header := block.Header()
	if err := v.engine.VerifyUncles(v.bc, block); err != nil {
		return err
	}
	if hash := types.CalcUncleHash(block.Uncles()); hash != header.UncleHash {
		return fmt.Errorf("uncle root hash mismatch: have %x, want %x", hash, header.UncleHash)
	}
	if hash := types.DeriveSha(block.Transactions()); hash != header.TxHash {
		return fmt.Errorf("transaction root hash mismatch: have %x, want %x", hash, header.TxHash)
	}
	if !v.bc.HasBlockAndState(block.ParentHash(), block.NumberU64()-1) {
		if !v.bc.HasBlock(block.ParentHash(), block.NumberU64()-1) {
			return consensus.ErrUnknownAncestor
		}
		return consensus.ErrPrunedAncestor
	}
	return nil
}

//validateState验证状态之后发生的各种更改
//过渡，如已用气体量、接收根和状态根
//本身。如果验证成功，则validateState返回数据库批处理
//否则为零，返回错误。
func (v *BlockValidator) ValidateState(block, parent *types.Block, statedb *state.StateDB, receipts types.Receipts, usedGas uint64) error {
	header := block.Header()
	if block.GasUsed() != usedGas {
		return fmt.Errorf("invalid gas used (remote: %d local: %d)", block.GasUsed(), usedGas)
	}
//使用从生成的收据中派生的块验证接收到的块的Bloom。
//对于有效块，应始终验证为真。
	rbloom := types.CreateBloom(receipts)
	if rbloom != header.Bloom {
		return fmt.Errorf("invalid bloom (remote: %x  local: %x)", header.Bloom, rbloom)
	}
//tre receipt trie's root（r=（tr[[h1，r1），……[HN，R1]）
	receiptSha := types.DeriveSha(receipts)
	if receiptSha != header.ReceiptHash {
		return fmt.Errorf("invalid receipt root hash (remote: %x local: %x)", header.ReceiptHash, receiptSha)
	}
//根据接收到的状态根验证状态根并引发
//如果不匹配则为错误。
	if root := statedb.IntermediateRoot(v.config.IsEIP158(header.Number)); header.Root != root {
		return fmt.Errorf("invalid merkle root (remote: %x local: %x)", header.Root, root)
	}
	return nil
}

//CalcGasLimit计算父块后面下一个块的气体限制。它的目标
//将基线气体保持在提供的地面以上，并向
//如果块已满，则返回CEIL。如果超过CEIL，它将始终降低
//煤气补贴。
func CalcGasLimit(parent *types.Block, gasFloor, gasCeil uint64) uint64 {
//contrib=（parentgasused*3/2）/1024
	contrib := (parent.GasUsed() + parent.GasUsed()/2) / params.GasLimitBoundDivisor

//衰变=父气体极限/1024-1
	decay := parent.GasLimit()/params.GasLimitBoundDivisor - 1

 /*
  策略：区块到矿井的气限是根据母公司的
  气体使用值。如果parentgasused>parentgaslimit*（2/3），那么我们
  增加它，否则降低它（或者如果它是正确的话保持不变
  使用时）增加/减少的数量取决于距离
  来自parentgaslimit*（2/3）parentgasused是。
 **/

	limit := parent.GasLimit() - decay + contrib
	if limit < params.MinGasLimit {
		limit = params.MinGasLimit
	}
//如果我们超出了允许的加油范围，我们会努力向他们靠拢。
	if limit < gasFloor {
		limit = parent.GasLimit() + decay
		if limit > gasFloor {
			limit = gasFloor
		}
	} else if limit > gasCeil {
		limit = parent.GasLimit() - decay
		if limit < gasCeil {
			limit = gasCeil
		}
	}
	return limit
}

