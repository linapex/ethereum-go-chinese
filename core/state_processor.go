
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450080614649856>


package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

//StateProcessor is a basic Processor, which takes care of transitioning
//从一点到另一点的状态。
//
//StateProcessor实现处理器。
type StateProcessor struct {
config *params.ChainConfig //链配置选项
bc     *BlockChain         //规范区块链
engine consensus.Engine    //集体奖励的共识引擎
}

//NewStateProcessor初始化新的StateProcessor。
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

//进程通过运行来根据以太坊规则处理状态更改
//事务消息使用statedb并对两者应用任何奖励
//处理器（coinbase）和任何包括的叔叔。
//
//流程返回流程中累积的收据和日志，以及
//返回过程中使用的气体量。如果有
//由于气体不足，无法执行事务，它将返回错误。
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, error) {
	var (
		receipts types.Receipts
		usedGas  = new(uint64)
		header   = block.Header()
		allLogs  []*types.Log
		gp       = new(GasPool).AddGas(block.GasLimit())
	)
//根据任何硬分叉规格改变块和状态
	if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0 {
		misc.ApplyDAOHardFork(statedb)
	}
//迭代并处理单个事务
	for i, tx := range block.Transactions() {
		statedb.Prepare(tx.Hash(), block.Hash(), i)
		receipt, _, err := ApplyTransaction(p.config, p.bc, nil, gp, statedb, header, tx, usedGas, cfg)
		if err != nil {
			return nil, nil, 0, err
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}
//完成区块，应用任何共识引擎特定的额外项目（例如区块奖励）
	p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles(), receipts)

	return receipts, allLogs, *usedGas, nil
}

//ApplyTransaction尝试将事务应用到给定的状态数据库
//并为其环境使用输入参数。它返回收据
//对于交易、使用的天然气以及交易失败时的错误，
//指示块无效。
func ApplyTransaction(config *params.ChainConfig, bc ChainContext, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*types.Receipt, uint64, error) {
	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		return nil, 0, err
	}
//创建要在EVM环境中使用的新上下文
	context := NewEVMContext(msg, header, bc, author)
//创建一个保存所有相关信息的新环境
//关于事务和调用机制。
	vmenv := vm.NewEVM(context, statedb, config, cfg)
//将事务应用于当前状态（包含在env中）
	_, gas, failed, err := ApplyMessage(vmenv, msg, gp)
	if err != nil {
		return nil, 0, err
	}
//用挂起的更改更新状态
	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	*usedGas += gas

//Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
//基于EIP阶段，我们将传递是否根触摸删除帐户。
	receipt := types.NewReceipt(root, failed, *usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = gas
//如果事务创建了合同，请将创建地址存储在收据中。
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
	}
//设置收据日志并创建一个用于过滤的Bloom
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	return receipt, gas, err
}

