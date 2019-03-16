
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450072909713408>


//+构建GouuZZ

package bitutil

import "bytes"

//Fuzz实现了Go-Fuzz引信方法来测试各种编码方法
//调用。
func Fuzz(data []byte) int {
	if len(data) == 0 {
		return -1
	}
	if data[0]%2 == 0 {
		return fuzzEncode(data[1:])
	}
	return fuzzDecode(data[1:])
}

//Fuzzencode实现了一种go-fuzz引信方法来测试位集编码和
//解码算法。
func fuzzEncode(data []byte) int {
	proc, _ := bitsetDecodeBytes(bitsetEncodeBytes(data), len(data))
	if !bytes.Equal(data, proc) {
		panic("content mismatch")
	}
	return 0
}

//fuzzdecode实现了一种go-fuzz引信方法来测试位解码和
//重新编码算法。
func fuzzDecode(data []byte) int {
	blob, err := bitsetDecodeBytes(data, 1024)
	if err != nil {
		return 0
	}
	if comp := bitsetEncodeBytes(blob); !bytes.Equal(comp, data) {
		panic("content mismatch")
	}
	return 0
}

