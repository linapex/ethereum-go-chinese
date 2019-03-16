
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450082242039808>


package vm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

//天然气成本
const (
	GasQuickStep   uint64 = 2
	GasFastestStep uint64 = 3
	GasFastStep    uint64 = 5
	GasMidStep     uint64 = 8
	GasSlowStep    uint64 = 10
	GasExtStep     uint64 = 20

	GasReturn       uint64 = 0
	GasStop         uint64 = 0
	GasContractByte uint64 = 200
)

//CalcGas返回呼叫的实际气体成本。
//
//天然气成本在宅基地价格变化期间发生了变化。考虑EIP150
//待实施。返回气体为气基*63/64。
func callGas(gasTable params.GasTable, availableGas, base uint64, callCost *big.Int) (uint64, error) {
	if gasTable.CreateBySuicide > 0 {
		availableGas = availableGas - base
		gas := availableGas - availableGas/64
//如果位长度超过64位，我们知道新计算的EIP150“气体”
//小于请求的金额。因此，我们把新的天然气换回来。
//返回错误。
		if callCost.BitLen() > 64 || gas < callCost.Uint64() {
			return gas, nil
		}
	}
	if callCost.BitLen() > 64 {
		return 0, errGasUintOverflow
	}

	return callCost.Uint64(), nil
}

