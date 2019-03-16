
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450081172492288>


package types

import (
	"math/big"
	"testing"
)

func TestBloom(t *testing.T) {
	positive := []string{
		"testtest",
		"test",
		"hallo",
		"other",
	}
	negative := []string{
		"tes",
		"lo",
	}

	var bloom Bloom
	for _, data := range positive {
		bloom.Add(new(big.Int).SetBytes([]byte(data)))
	}

	for _, data := range positive {
		if !bloom.TestBytes([]byte(data)) {
			t.Error("expected", data, "to test true")
		}
	}
	for _, data := range negative {
		if bloom.TestBytes([]byte(data)) {
			t.Error("did not expect", data, "to test true")
		}
	}
}

/*
进口（
 “测试”

 “github.com/ethereum/go-ethereum/core/state”
）

func测试bloom9（t*testing.t）
 测试用例：=[]字节（“测试测试”）
 bin：=logsbloom（[]state.log_
  测试用例，[]字节[]字节（“hellohello”），零，
 }。字节（）
 res：=bloomlookup（bin，测试用例）

 如果！RES{
  t.errorf（“Bloom查找失败”）
 }
}


func测试地址（t*testing.t）
 块：=&block
 block.coinbase=common.hex2bytes（“2234AE42D6DD7384BC8584E50419EA3AC75B83F”）。
 fmt.printf（“%x\n”，crypto.keccak256（block.coinbase））。

 bin:=创建Bloom（块）
 fmt.printf（“bin=%x\n”，common.leftpaddbytes（bin，64））。
}
**/


