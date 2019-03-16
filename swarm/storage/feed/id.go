
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450119026085888>


package feed

import (
	"fmt"
	"hash"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage/feed/lookup"

	"github.com/ethereum/go-ethereum/swarm/storage"
)

//ID唯一标识网络上的更新。
type ID struct {
	Feed         `json:"feed"`
	lookup.Epoch `json:"epoch"`
}

//ID布局：
//进纸长度字节
//时代纪元
const idLength = feedLength + lookup.EpochLength

//addr计算与此ID对应的源更新块地址
func (u *ID) Addr() (updateAddr storage.Address) {
	serializedData := make([]byte, idLength)
	var cursor int
	u.Feed.binaryPut(serializedData[cursor : cursor+feedLength])
	cursor += feedLength

	eid := u.Epoch.ID()
	copy(serializedData[cursor:cursor+lookup.EpochLength], eid[:])

	hasher := hashPool.Get().(hash.Hash)
	defer hashPool.Put(hasher)
	hasher.Reset()
	hasher.Write(serializedData)
	return hasher.Sum(nil)
}

//BinaryPut将此实例序列化到提供的切片中
func (u *ID) binaryPut(serializedData []byte) error {
	if len(serializedData) != idLength {
		return NewErrorf(ErrInvalidValue, "Incorrect slice size to serialize ID. Expected %d, got %d", idLength, len(serializedData))
	}
	var cursor int
	if err := u.Feed.binaryPut(serializedData[cursor : cursor+feedLength]); err != nil {
		return err
	}
	cursor += feedLength

	epochBytes, err := u.Epoch.MarshalBinary()
	if err != nil {
		return err
	}
	copy(serializedData[cursor:cursor+lookup.EpochLength], epochBytes[:])
	cursor += lookup.EpochLength

	return nil
}

//BinaryLength返回序列化时此结构的预期大小
func (u *ID) binaryLength() int {
	return idLength
}

//binaryget从传递的切片中包含的信息还原当前实例
func (u *ID) binaryGet(serializedData []byte) error {
	if len(serializedData) != idLength {
		return NewErrorf(ErrInvalidValue, "Incorrect slice size to read ID. Expected %d, got %d", idLength, len(serializedData))
	}

	var cursor int
	if err := u.Feed.binaryGet(serializedData[cursor : cursor+feedLength]); err != nil {
		return err
	}
	cursor += feedLength

	if err := u.Epoch.UnmarshalBinary(serializedData[cursor : cursor+lookup.EpochLength]); err != nil {
		return err
	}
	cursor += lookup.EpochLength

	return nil
}

//FromValues从字符串键值存储中反序列化此实例
//用于分析查询字符串
func (u *ID) FromValues(values Values) error {
	level, _ := strconv.ParseUint(values.Get("level"), 10, 32)
	u.Epoch.Level = uint8(level)
	u.Epoch.Time, _ = strconv.ParseUint(values.Get("time"), 10, 64)

	if u.Feed.User == (common.Address{}) {
		return u.Feed.FromValues(values)
	}
	return nil
}

//AppendValues将此结构序列化到提供的字符串键值存储区中
//用于生成查询字符串
func (u *ID) AppendValues(values Values) {
	values.Set("level", fmt.Sprintf("%d", u.Epoch.Level))
	values.Set("time", fmt.Sprintf("%d", u.Epoch.Time))
	u.Feed.AppendValues(values)
}

