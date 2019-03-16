
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450119634259968>


package feed

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/swarm/chunk"
)

//ProtocolVersion定义将包含在每个更新消息中的协议的当前版本
const ProtocolVersion uint8 = 0

const headerLength = 8

//header定义包含协议版本字节的更新消息头
type Header struct {
Version uint8                   //协议版本
Padding [headerLength - 1]uint8 //留作将来使用
}

//更新将作为源更新的一部分发送的信息封装起来。
type Update struct {
Header Header //
ID            //源更新识别信息
data   []byte //实际数据负载
}

const minimumUpdateDataLength = idLength + headerLength + 1

//maxupdatedatalength指示源更新的最大负载大小
const MaxUpdateDataLength = chunk.DefaultSize - signatureLength - idLength - headerLength

//BinaryPut将源更新信息序列化到给定切片中
func (r *Update) binaryPut(serializedData []byte) error {
	datalength := len(r.data)
	if datalength == 0 {
		return NewError(ErrInvalidValue, "a feed update must contain data")
	}

	if datalength > MaxUpdateDataLength {
		return NewErrorf(ErrInvalidValue, "feed update data is too big (length=%d). Max length=%d", datalength, MaxUpdateDataLength)
	}

	if len(serializedData) != r.binaryLength() {
		return NewErrorf(ErrInvalidValue, "slice passed to putBinary must be of exact size. Expected %d bytes", r.binaryLength())
	}

	var cursor int
//序列化头
	serializedData[cursor] = r.Header.Version
	copy(serializedData[cursor+1:headerLength], r.Header.Padding[:headerLength-1])
	cursor += headerLength

//序列化ID
	if err := r.ID.binaryPut(serializedData[cursor : cursor+idLength]); err != nil {
		return err
	}
	cursor += idLength

//添加数据
	copy(serializedData[cursor:], r.data)
	cursor += datalength

	return nil
}

//BinaryLength返回此结构编码所需的字节数。
func (r *Update) binaryLength() int {
	return idLength + headerLength + len(r.data)
}

//binaryget从传递的字节片中包含的信息填充此实例
func (r *Update) binaryGet(serializedData []byte) error {
	if len(serializedData) < minimumUpdateDataLength {
		return NewErrorf(ErrNothingToReturn, "chunk less than %d bytes cannot be a feed update chunk", minimumUpdateDataLength)
	}
	dataLength := len(serializedData) - idLength - headerLength
//此时，我们可以确信我们有正确的数据长度来读取

	var cursor int

//反序列化头
r.Header.Version = serializedData[cursor]                                      //提取协议版本
copy(r.Header.Padding[:headerLength-1], serializedData[cursor+1:headerLength]) //提取填充物
	cursor += headerLength

	if err := r.ID.binaryGet(serializedData[cursor : cursor+idLength]); err != nil {
		return err
	}
	cursor += idLength

	data := serializedData[cursor : cursor+dataLength]
	cursor += dataLength

//既然所有检查都通过了，那么就将数据复制到结构中。
	r.data = make([]byte, dataLength)
	copy(r.data, data)

	return nil

}

//FromValues从字符串键值存储中反序列化此实例
//用于分析查询字符串
func (r *Update) FromValues(values Values, data []byte) error {
	r.data = data
	version, _ := strconv.ParseUint(values.Get("protocolVersion"), 10, 32)
	r.Header.Version = uint8(version)
	return r.ID.FromValues(values)
}

//AppendValues将此结构序列化到提供的字符串键值存储区中
//用于生成查询字符串
func (r *Update) AppendValues(values Values) []byte {
	r.ID.AppendValues(values)
	values.Set("protocolVersion", fmt.Sprintf("%d", r.Header.Version))
	return r.data
}

