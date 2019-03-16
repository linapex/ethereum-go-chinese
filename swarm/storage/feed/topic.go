
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450119541985280>


package feed

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/bitutil"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

//topic length建立主题字符串的最大长度
const TopicLength = storage.AddressLength

//主题表示提要的内容
type Topic [TopicLength]byte

//创建名称/相关内容太长的主题时返回errtopictolong
var ErrTopicTooLong = fmt.Errorf("Topic is too long. Max length is %d", TopicLength)

//new topic从提供的名称和“相关内容”字节数组创建新主题，
//将两者合并在一起。
//如果RelatedContent或Name长于TopicLength，它们将被截断并返回错误
//名称可以是空字符串
//相关内容可以为零
func NewTopic(name string, relatedContent []byte) (topic Topic, err error) {
	if relatedContent != nil {
		contentLength := len(relatedContent)
		if contentLength > TopicLength {
			contentLength = TopicLength
			err = ErrTopicTooLong
		}
		copy(topic[:], relatedContent[:contentLength])
	}
	nameBytes := []byte(name)
	nameLength := len(nameBytes)
	if nameLength > TopicLength {
		nameLength = TopicLength
		err = ErrTopicTooLong
	}
	bitutil.XORBytes(topic[:], topic[:], nameBytes[:nameLength])
	return topic, err
}

//hex将返回编码为十六进制字符串的主题
func (t *Topic) Hex() string {
	return hexutil.Encode(t[:])
}

//FromHex将把十六进制字符串解析到此主题实例中
func (t *Topic) FromHex(hex string) error {
	bytes, err := hexutil.Decode(hex)
	if err != nil || len(bytes) != len(t) {
		return NewErrorf(ErrInvalidValue, "Cannot decode topic")
	}
	copy(t[:], bytes)
	return nil
}

//name将尝试从主题中提取主题名称
func (t *Topic) Name(relatedContent []byte) string {
	nameBytes := *t
	if relatedContent != nil {
		contentLength := len(relatedContent)
		if contentLength > TopicLength {
			contentLength = TopicLength
		}
		bitutil.XORBytes(nameBytes[:], t[:], relatedContent[:contentLength])
	}
	z := bytes.IndexByte(nameBytes[:], 0)
	if z < 0 {
		z = TopicLength
	}
	return string(nameBytes[:z])

}

//unmashaljson实现json.unmarshaller接口
func (t *Topic) UnmarshalJSON(data []byte) error {
	var hex string
	json.Unmarshal(data, &hex)
	return t.FromHex(hex)
}

//marshaljson实现json.marshaller接口
func (t *Topic) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Hex())
}

