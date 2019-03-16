
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450107932151808>


package rlp

import (
	"bytes"
	"fmt"
)

type structWithTail struct {
	A, B uint
	C    []uint `rlp:"tail"`
}

func ExampleDecode_structTagTail() {
//在本例中，“tail”结构标记用于解码
//结构中的不同长度。
	var val structWithTail

	err := Decode(bytes.NewReader([]byte{0xC4, 0x01, 0x02, 0x03, 0x04}), &val)
	fmt.Printf("with 4 elements: err=%v val=%v\n", err, val)

	err = Decode(bytes.NewReader([]byte{0xC6, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}), &val)
	fmt.Printf("with 6 elements: err=%v val=%v\n", err, val)

//请注意，必须至少有两个list元素存在于
//填充字段A和B：
	err = Decode(bytes.NewReader([]byte{0xC1, 0x01}), &val)
	fmt.Printf("with 1 element: err=%q\n", err)

//输出：
//有4个元素：err=<nil>val=1 2[3 4]
//有6个元素：err=<nil>val=1 2[3 4 5 6]
//with 1 element:err=“rlp:rlp.structWithTail的元素太少”
}

