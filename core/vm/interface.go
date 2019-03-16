
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450082535641088>


package vm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//StateDB是用于完整状态查询的EVM数据库。
type StateDB interface {
	CreateAccount(common.Address)

	SubBalance(common.Address, *big.Int)
	AddBalance(common.Address, *big.Int)
	GetBalance(common.Address) *big.Int

	GetNonce(common.Address) uint64
	SetNonce(common.Address, uint64)

	GetCodeHash(common.Address) common.Hash
	GetCode(common.Address) []byte
	SetCode(common.Address, []byte)
	GetCodeSize(common.Address) int

	AddRefund(uint64)
	SubRefund(uint64)
	GetRefund() uint64

	GetCommittedState(common.Address, common.Hash) common.Hash
	GetState(common.Address, common.Hash) common.Hash
	SetState(common.Address, common.Hash, common.Hash)

	Suicide(common.Address) bool
	HasSuicided(common.Address) bool

//exist报告给定帐户是否处于状态。
//值得注意的是，对于自杀账户来说，这也应该是正确的。
	Exist(common.Address) bool
//empty返回给定帐户是否为空。空的
//根据EIP161定义（balance=nonce=code=0）。
	Empty(common.Address) bool

	RevertToSnapshot(int)
	Snapshot() int

	AddLog(*types.Log)
	AddPreimage(common.Hash, []byte)

	ForEachStorage(common.Address, func(common.Hash, common.Hash) bool)
}

//CallContext为EVM调用约定提供基本接口。EVM
//取决于为执行子调用和初始化新的EVM合同而实现的上下文。
type CallContext interface {
//调用另一个合同
	Call(env *EVM, me ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error)
//采用他人的合同代码并在我们自己的上下文中执行
	CallCode(env *EVM, me ContractRef, addr common.Address, data []byte, gas, value *big.Int) ([]byte, error)
//与callcode相同，但发送方和值从父作用域传播到子作用域
	DelegateCall(env *EVM, me ContractRef, addr common.Address, data []byte, gas *big.Int) ([]byte, error)
//创建新合同
	Create(env *EVM, me ContractRef, data []byte, gas, value *big.Int) ([]byte, common.Address, error)
}

