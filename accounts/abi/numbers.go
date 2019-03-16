
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:31</date>
//</624450062205849600>


package abi

import (
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var (
	bigT      = reflect.TypeOf(&big.Int{})
	derefbigT = reflect.TypeOf(big.Int{})
	uint8T    = reflect.TypeOf(uint8(0))
	uint16T   = reflect.TypeOf(uint16(0))
	uint32T   = reflect.TypeOf(uint32(0))
	uint64T   = reflect.TypeOf(uint64(0))
	int8T     = reflect.TypeOf(int8(0))
	int16T    = reflect.TypeOf(int16(0))
	int32T    = reflect.TypeOf(int32(0))
	int64T    = reflect.TypeOf(int64(0))
	addressT  = reflect.TypeOf(common.Address{})
)

//u256将一个大整数转换为256bit evm数。
func U256(n *big.Int) []byte {
	return math.PaddedBigBytes(math.U256(n), 32)
}

