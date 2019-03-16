
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450101527449600>


//包含扭曲的包装，以允许跨越空接口。

package geth

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//接口表示Go接口的包装版本，具有
//存储任意数据类型。
//
//因为在Go和Mobile之间无法转换任意性
//平台，我们使用显式getter和setter进行转换。那里
//当然，没有必要列举所有的东西，只是足以支持
//合同绑定需要客户端生成的代码。
type Interface struct {
	object interface{}
}

//NewInterface创建一个新的空接口，可用于传递
//泛型类型。
func NewInterface() *Interface {
	return new(Interface)
}

func (i *Interface) SetBool(b bool)                { i.object = &b }
func (i *Interface) SetBools(bs []bool)            { i.object = &bs }
func (i *Interface) SetString(str string)          { i.object = &str }
func (i *Interface) SetStrings(strs *Strings)      { i.object = &strs.strs }
func (i *Interface) SetBinary(binary []byte)       { b := common.CopyBytes(binary); i.object = &b }
func (i *Interface) SetBinaries(binaries [][]byte) { i.object = &binaries }
func (i *Interface) SetAddress(address *Address)   { i.object = &address.address }
func (i *Interface) SetAddresses(addrs *Addresses) { i.object = &addrs.addresses }
func (i *Interface) SetHash(hash *Hash)            { i.object = &hash.hash }
func (i *Interface) SetHashes(hashes *Hashes)      { i.object = &hashes.hashes }
func (i *Interface) SetInt8(n int8)                { i.object = &n }
func (i *Interface) SetInt16(n int16)              { i.object = &n }
func (i *Interface) SetInt32(n int32)              { i.object = &n }
func (i *Interface) SetInt64(n int64)              { i.object = &n }
func (i *Interface) SetUint8(bigint *BigInt)       { n := uint8(bigint.bigint.Uint64()); i.object = &n }
func (i *Interface) SetUint16(bigint *BigInt)      { n := uint16(bigint.bigint.Uint64()); i.object = &n }
func (i *Interface) SetUint32(bigint *BigInt)      { n := uint32(bigint.bigint.Uint64()); i.object = &n }
func (i *Interface) SetUint64(bigint *BigInt)      { n := bigint.bigint.Uint64(); i.object = &n }
func (i *Interface) SetBigInt(bigint *BigInt)      { i.object = &bigint.bigint }
func (i *Interface) SetBigInts(bigints *BigInts)   { i.object = &bigints.bigints }

func (i *Interface) SetDefaultBool()      { i.object = new(bool) }
func (i *Interface) SetDefaultBools()     { i.object = new([]bool) }
func (i *Interface) SetDefaultString()    { i.object = new(string) }
func (i *Interface) SetDefaultStrings()   { i.object = new([]string) }
func (i *Interface) SetDefaultBinary()    { i.object = new([]byte) }
func (i *Interface) SetDefaultBinaries()  { i.object = new([][]byte) }
func (i *Interface) SetDefaultAddress()   { i.object = new(common.Address) }
func (i *Interface) SetDefaultAddresses() { i.object = new([]common.Address) }
func (i *Interface) SetDefaultHash()      { i.object = new(common.Hash) }
func (i *Interface) SetDefaultHashes()    { i.object = new([]common.Hash) }
func (i *Interface) SetDefaultInt8()      { i.object = new(int8) }
func (i *Interface) SetDefaultInt16()     { i.object = new(int16) }
func (i *Interface) SetDefaultInt32()     { i.object = new(int32) }
func (i *Interface) SetDefaultInt64()     { i.object = new(int64) }
func (i *Interface) SetDefaultUint8()     { i.object = new(uint8) }
func (i *Interface) SetDefaultUint16()    { i.object = new(uint16) }
func (i *Interface) SetDefaultUint32()    { i.object = new(uint32) }
func (i *Interface) SetDefaultUint64()    { i.object = new(uint64) }
func (i *Interface) SetDefaultBigInt()    { i.object = new(*big.Int) }
func (i *Interface) SetDefaultBigInts()   { i.object = new([]*big.Int) }

func (i *Interface) GetBool() bool            { return *i.object.(*bool) }
func (i *Interface) GetBools() []bool         { return *i.object.(*[]bool) }
func (i *Interface) GetString() string        { return *i.object.(*string) }
func (i *Interface) GetStrings() *Strings     { return &Strings{*i.object.(*[]string)} }
func (i *Interface) GetBinary() []byte        { return *i.object.(*[]byte) }
func (i *Interface) GetBinaries() [][]byte    { return *i.object.(*[][]byte) }
func (i *Interface) GetAddress() *Address     { return &Address{*i.object.(*common.Address)} }
func (i *Interface) GetAddresses() *Addresses { return &Addresses{*i.object.(*[]common.Address)} }
func (i *Interface) GetHash() *Hash           { return &Hash{*i.object.(*common.Hash)} }
func (i *Interface) GetHashes() *Hashes       { return &Hashes{*i.object.(*[]common.Hash)} }
func (i *Interface) GetInt8() int8            { return *i.object.(*int8) }
func (i *Interface) GetInt16() int16          { return *i.object.(*int16) }
func (i *Interface) GetInt32() int32          { return *i.object.(*int32) }
func (i *Interface) GetInt64() int64          { return *i.object.(*int64) }
func (i *Interface) GetUint8() *BigInt {
	return &BigInt{new(big.Int).SetUint64(uint64(*i.object.(*uint8)))}
}
func (i *Interface) GetUint16() *BigInt {
	return &BigInt{new(big.Int).SetUint64(uint64(*i.object.(*uint16)))}
}
func (i *Interface) GetUint32() *BigInt {
	return &BigInt{new(big.Int).SetUint64(uint64(*i.object.(*uint32)))}
}
func (i *Interface) GetUint64() *BigInt {
	return &BigInt{new(big.Int).SetUint64(*i.object.(*uint64))}
}
func (i *Interface) GetBigInt() *BigInt   { return &BigInt{*i.object.(**big.Int)} }
func (i *Interface) GetBigInts() *BigInts { return &BigInts{*i.object.(*[]*big.Int)} }

//接口是被包装的一般对象的切片。
type Interfaces struct {
	objects []interface{}
}

//NewInterfaces创建一个未初始化的接口切片。
func NewInterfaces(size int) *Interfaces {
	return &Interfaces{
		objects: make([]interface{}, size),
	}
}

//SIZE返回切片中的接口数。
func (i *Interfaces) Size() int {
	return len(i.objects)
}

//get从切片返回给定索引处的bigint。
func (i *Interfaces) Get(index int) (iface *Interface, _ error) {
	if index < 0 || index >= len(i.objects) {
		return nil, errors.New("index out of bounds")
	}
	return &Interface{i.objects[index]}, nil
}

//set在切片中的给定索引处设置big int。
func (i *Interfaces) Set(index int, object *Interface) error {
	if index < 0 || index >= len(i.objects) {
		return errors.New("index out of bounds")
	}
	i.objects[index] = object.object
	return nil
}

