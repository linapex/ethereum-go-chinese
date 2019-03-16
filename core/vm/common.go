
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450081893912576>


package vm

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

//计算步骤所需的内存大小
func calcMemSize(off, l *big.Int) *big.Int {
	if l.Sign() == 0 {
		return common.Big0
	}

	return new(big.Int).Add(off, l)
}

//GetData根据起始点、大小和焊盘从数据返回一个切片。
//最大大小为零。此函数是溢出安全的。
func getData(data []byte, start uint64, size uint64) []byte {
	length := uint64(len(data))
	if start > length {
		start = length
	}
	end := start + size
	if end > length {
		end = length
	}
	return common.RightPadBytes(data[start:end], int(size))
}

//GetDataBig根据开始、大小和焊盘从数据返回一个切片
//最大大小为零。此函数是溢出安全的。
func getDataBig(data []byte, start *big.Int, size *big.Int) []byte {
	dlen := big.NewInt(int64(len(data)))

	s := math.BigMin(start, dlen)
	e := math.BigMin(new(big.Int).Add(s, size), dlen)
	return common.RightPadBytes(data[s.Uint64():e.Uint64()], int(size.Uint64()))
}

//bigunt64返回强制转换为uint64的整数，并返回它是否
//过程中溢出。
func bigUint64(v *big.Int) (uint64, bool) {
	return v.Uint64(), v.BitLen() > 64
}

//TowordSize返回内存扩展所需的接收字大小。
func toWordSize(size uint64) uint64 {
	if size > math.MaxUint64-31 {
		return math.MaxUint64/32 + 1
	}

	return (size + 31) / 32
}

func allZero(b []byte) bool {
	for _, byte := range b {
		if byte != 0 {
			return false
		}
	}
	return true
}

