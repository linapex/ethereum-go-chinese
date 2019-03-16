
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:31</date>
//</624450061908054016>


package abi

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//事件是可能由EVM的日志机制触发的事件。事件
//保存有关生成的输出的类型信息（输入）。匿名事件
//不要将签名规范表示作为第一个日志主题。
type Event struct {
	Name      string
	Anonymous bool
	Inputs    Arguments
}

func (e Event) String() string {
	inputs := make([]string, len(e.Inputs))
	for i, input := range e.Inputs {
		inputs[i] = fmt.Sprintf("%v %v", input.Type, input.Name)
		if input.Indexed {
			inputs[i] = fmt.Sprintf("%v indexed %v", input.Type, input.Name)
		}
	}
	return fmt.Sprintf("event %v(%v)", e.Name, strings.Join(inputs, ", "))
}

//ID返回由
//用于标识事件名称和类型的ABI定义。
func (e Event) Id() common.Hash {
	types := make([]string, len(e.Inputs))
	i := 0
	for _, input := range e.Inputs {
		types[i] = input.Type.String()
		i++
	}
	return common.BytesToHash(crypto.Keccak256([]byte(fmt.Sprintf("%v(%v)", e.Name, strings.Join(types, ",")))))
}

