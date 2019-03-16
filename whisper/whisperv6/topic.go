
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:46</date>
//</624450125749555200>


//包含耳语协议主题元素。

package whisperv6

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//TopicType表示密码安全的概率部分
//消息的分类，确定为
//sha3消息原始作者给出的某些任意数据的散列。
type TopicType [TopicLength]byte

//BytesToTopic从主题的字节数组表示形式转换
//进入TopicType类型。
func BytesToTopic(b []byte) (t TopicType) {
	sz := TopicLength
	if x := len(b); x < TopicLength {
		sz = x
	}
	for i := 0; i < sz; i++ {
		t[i] = b[i]
	}
	return t
}

//字符串将主题字节数组转换为字符串表示形式。
func (t *TopicType) String() string {
	return common.ToHex(t[:])
}

//marshalText返回t的十六进制表示形式。
func (t TopicType) MarshalText() ([]byte, error) {
	return hexutil.Bytes(t[:]).MarshalText()
}

//UnmarshalText解析主题的十六进制表示。
func (t *TopicType) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Topic", input, t[:])
}

