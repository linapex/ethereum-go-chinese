
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450083043151872>


package vm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
)

func memorySha3(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(0), stack.Back(1))
}

func memoryCallDataCopy(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(0), stack.Back(2))
}

func memoryReturnDataCopy(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(0), stack.Back(2))
}

func memoryCodeCopy(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(0), stack.Back(2))
}

func memoryExtCodeCopy(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(1), stack.Back(3))
}

func memoryMLoad(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(0), big.NewInt(32))
}

func memoryMStore8(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(0), big.NewInt(1))
}

func memoryMStore(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(0), big.NewInt(32))
}

func memoryCreate(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(1), stack.Back(2))
}

func memoryCreate2(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(1), stack.Back(2))
}

func memoryCall(stack *Stack) *big.Int {
	x := calcMemSize(stack.Back(5), stack.Back(6))
	y := calcMemSize(stack.Back(3), stack.Back(4))

	return math.BigMax(x, y)
}

func memoryDelegateCall(stack *Stack) *big.Int {
	x := calcMemSize(stack.Back(4), stack.Back(5))
	y := calcMemSize(stack.Back(2), stack.Back(3))

	return math.BigMax(x, y)
}

func memoryStaticCall(stack *Stack) *big.Int {
	x := calcMemSize(stack.Back(4), stack.Back(5))
	y := calcMemSize(stack.Back(2), stack.Back(3))

	return math.BigMax(x, y)
}

func memoryReturn(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(0), stack.Back(1))
}

func memoryRevert(stack *Stack) *big.Int {
	return calcMemSize(stack.Back(0), stack.Back(1))
}

func memoryLog(stack *Stack) *big.Int {
	mSize, mStart := stack.Back(1), stack.Back(0)
	return calcMemSize(mStart, mSize)
}

