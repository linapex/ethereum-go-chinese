
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450073677271040>


package hexutil_test

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type MyType [5]byte

func (v *MyType) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("MyType", input, v[:])
}

func (v MyType) String() string {
	return hexutil.Bytes(v[:]).String()
}

func ExampleUnmarshalFixedText() {
	var v1, v2 MyType
	fmt.Println("v1 error:", json.Unmarshal([]byte(`"0x01"`), &v1))
	fmt.Println("v2 error:", json.Unmarshal([]byte(`"0x0101010101"`), &v2))
	fmt.Println("v2:", v2)
//输出：
//v1错误：十六进制字符串的长度为2，MyType需要10
//V2错误：<nIL>
//版本2:0x0101010101
}

