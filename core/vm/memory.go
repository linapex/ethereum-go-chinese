
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450083005403136>


package vm

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
)

//内存为以太坊虚拟机实现了一个简单的内存模型。
type Memory struct {
	store       []byte
	lastGasCost uint64
}

//new memory返回新的内存模型。
func NewMemory() *Memory {
	return &Memory{}
}

//将“设置偏移量+大小”设置为“值”
func (m *Memory) Set(offset, size uint64, value []byte) {
//偏移量可能大于0，大小等于0。这是因为
//当大小为零（no-op）时，CalcMemsize（common.go）可能返回0。
	if size > 0 {
//存储长度不能小于偏移量+大小。
//在设置内存之前，应调整存储的大小
		if offset+size > uint64(len(m.store)) {
			panic("invalid memory: store empty")
		}
		copy(m.store[offset:offset+size], value)
	}
}

//set32将从偏移量开始的32个字节设置为val值，用零左填充到
//32字节。
func (m *Memory) Set32(offset uint64, val *big.Int) {
//存储长度不能小于偏移量+大小。
//在设置内存之前，应调整存储的大小
	if offset+32 > uint64(len(m.store)) {
		panic("invalid memory: store empty")
	}
//把记忆区归零
	copy(m.store[offset:offset+32], []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
//填写相关位
	math.ReadBits(val, m.store[offset:offset+32])
}

//调整大小将内存大小调整为
func (m *Memory) Resize(size uint64) {
	if uint64(m.Len()) < size {
		m.store = append(m.store, make([]byte, size-uint64(m.Len()))...)
	}
}

//get返回偏移量+作为新切片的大小
func (m *Memory) Get(offset, size int64) (cpy []byte) {
	if size == 0 {
		return nil
	}

	if len(m.store) > int(offset) {
		cpy = make([]byte, size)
		copy(cpy, m.store[offset:offset+size])

		return
	}

	return
}

//getptr返回偏移量+大小
func (m *Memory) GetPtr(offset, size int64) []byte {
	if size == 0 {
		return nil
	}

	if len(m.store) > int(offset) {
		return m.store[offset : offset+size]
	}

	return nil
}

//len返回背衬片的长度
func (m *Memory) Len() int {
	return len(m.store)
}

//数据返回备份切片
func (m *Memory) Data() []byte {
	return m.store
}

//打印转储内存的内容。
func (m *Memory) Print() {
	fmt.Printf("### mem %d bytes ###\n", len(m.store))
	if len(m.store) > 0 {
		addr := 0
		for i := 0; i+32 <= len(m.store); i += 32 {
			fmt.Printf("%03d: % x\n", addr, m.store[i:i+32])
			addr++
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("####################")
}

