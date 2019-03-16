
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450118795399168>


package feed

import (
	"hash"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

//源表示特定用户对主题的更新流
type Feed struct {
	Topic Topic          `json:"topic"`
	User  common.Address `json:"user"`
}

//饲料布局：
//TopicLength字节
//useraddr common.addresslength字节
const feedLength = TopicLength + common.AddressLength

//mapkey计算此源的唯一ID。由“handler”中的缓存映射使用
func (f *Feed) mapKey() uint64 {
	serializedData := make([]byte, feedLength)
	f.binaryPut(serializedData)
	hasher := hashPool.Get().(hash.Hash)
	defer hashPool.Put(hasher)
	hasher.Reset()
	hasher.Write(serializedData)
	hash := hasher.Sum(nil)
	return *(*uint64)(unsafe.Pointer(&hash[0]))
}

//BinaryPut将此源实例序列化到提供的切片中
func (f *Feed) binaryPut(serializedData []byte) error {
	if len(serializedData) != feedLength {
		return NewErrorf(ErrInvalidValue, "Incorrect slice size to serialize feed. Expected %d, got %d", feedLength, len(serializedData))
	}
	var cursor int
	copy(serializedData[cursor:cursor+TopicLength], f.Topic[:TopicLength])
	cursor += TopicLength

	copy(serializedData[cursor:cursor+common.AddressLength], f.User[:])
	cursor += common.AddressLength

	return nil
}

//BinaryLength返回序列化时此结构的预期大小
func (f *Feed) binaryLength() int {
	return feedLength
}

//binaryget从传递的切片中包含的信息还原当前实例
func (f *Feed) binaryGet(serializedData []byte) error {
	if len(serializedData) != feedLength {
		return NewErrorf(ErrInvalidValue, "Incorrect slice size to read feed. Expected %d, got %d", feedLength, len(serializedData))
	}

	var cursor int
	copy(f.Topic[:], serializedData[cursor:cursor+TopicLength])
	cursor += TopicLength

	copy(f.User[:], serializedData[cursor:cursor+common.AddressLength])
	cursor += common.AddressLength

	return nil
}

//十六进制将提要序列化为十六进制字符串
func (f *Feed) Hex() string {
	serializedData := make([]byte, feedLength)
	f.binaryPut(serializedData)
	return hexutil.Encode(serializedData)
}

//FromValues从字符串键值存储中反序列化此实例
//用于分析查询字符串
func (f *Feed) FromValues(values Values) (err error) {
	topic := values.Get("topic")
	if topic != "" {
		if err := f.Topic.FromHex(values.Get("topic")); err != nil {
			return err
		}
} else { //查看用户集名称和相关内容
		name := values.Get("name")
		relatedContent, _ := hexutil.Decode(values.Get("relatedcontent"))
		if len(relatedContent) > 0 {
			if len(relatedContent) < storage.AddressLength {
				return NewErrorf(ErrInvalidValue, "relatedcontent field must be a hex-encoded byte array exactly %d bytes long", storage.AddressLength)
			}
			relatedContent = relatedContent[:storage.AddressLength]
		}
		f.Topic, err = NewTopic(name, relatedContent)
		if err != nil {
			return err
		}
	}
	f.User = common.HexToAddress(values.Get("user"))
	return nil
}

//AppendValues将此结构序列化到提供的字符串键值存储区中
//用于生成查询字符串
func (f *Feed) AppendValues(values Values) {
	values.Set("topic", f.Topic.Hex())
	values.Set("user", f.User.Hex())
}

