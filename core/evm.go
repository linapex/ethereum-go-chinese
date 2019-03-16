
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450078748184576>


package core

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
)

//ChainContext支持从
//交易处理过程中使用的当前区块链。
type ChainContext interface {
//引擎检索链的共识引擎。
	Engine() consensus.Engine

//GetHeader返回与其哈希相对应的哈希。
	GetHeader(common.Hash, uint64) *types.Header
}

//new evm context创建一个新上下文以在evm中使用。
func NewEVMContext(msg Message, header *types.Header, chain ChainContext, author *common.Address) vm.Context {
//如果没有明确的作者（即没有挖掘），则从头中提取
	var beneficiary common.Address
	if author == nil {
beneficiary, _ = chain.Engine().Author(header) //忽略错误，我们已经过了头验证
	} else {
		beneficiary = *author
	}
	return vm.Context{
		CanTransfer: CanTransfer,
		Transfer:    Transfer,
		GetHash:     GetHashFn(header, chain),
		Origin:      msg.From(),
		Coinbase:    beneficiary,
		BlockNumber: new(big.Int).Set(header.Number),
		Time:        new(big.Int).Set(header.Time),
		Difficulty:  new(big.Int).Set(header.Difficulty),
		GasLimit:    header.GasLimit,
		GasPrice:    new(big.Int).Set(msg.GasPrice()),
	}
}

//gethashfn返回gethashfunc，该函数按数字检索头哈希
func GetHashFn(ref *types.Header, chain ChainContext) func(n uint64) common.Hash {
	var cache map[uint64]common.Hash

	return func(n uint64) common.Hash {
//如果还没有哈希缓存，请创建一个
		if cache == nil {
			cache = map[uint64]common.Hash{
				ref.Number.Uint64() - 1: ref.ParentHash,
			}
		}
//尝试完成来自缓存的请求
		if hash, ok := cache[n]; ok {
			return hash
		}
//不缓存，迭代块并缓存哈希
		for header := chain.GetHeader(ref.ParentHash, ref.Number.Uint64()-1); header != nil; header = chain.GetHeader(header.ParentHash, header.Number.Uint64()-1) {
			cache[header.Number.Uint64()-1] = header.ParentHash
			if n == header.Number.Uint64()-1 {
				return header.ParentHash
			}
		}
		return common.Hash{}
	}
}

//CanTransfer检查地址的账户中是否有足够的资金进行转账。
//这不需要考虑必要的气体以使转移有效。
func CanTransfer(db vm.StateDB, addr common.Address, amount *big.Int) bool {
	return db.GetBalance(addr).Cmp(amount) >= 0
}

//转账从发送方减去金额，并使用给定的数据库向接收方添加金额。
func Transfer(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {
	db.SubBalance(sender, amount)
	db.AddBalance(recipient, amount)
}

