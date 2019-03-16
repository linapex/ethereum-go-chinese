
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450118602461184>


package feed

import "github.com/ethereum/go-ethereum/common/hexutil"

type binarySerializer interface {
	binaryPut(serializedData []byte) error
	binaryLength() int
	binaryGet(serializedData []byte) error
}

//值接口表示字符串键值存储
//用于生成查询字符串
type Values interface {
	Get(key string) string
	Set(key, value string)
}

type valueSerializer interface {
	FromValues(values Values) error
	AppendValues(values Values)
}

//十六进制序列化结构并将其转换为十六进制字符串
func Hex(bin binarySerializer) string {
	b := make([]byte, bin.binaryLength())
	bin.binaryPut(b)
	return hexutil.Encode(b)
}

