
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:31</date>
//</624450062050660352>


package abi

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
)

//方法表示给定“name”的可调用项，以及该方法是否为常量。
//如果方法为“const”，则不需要为此创建任何事务
//特定方法调用。它可以很容易地使用本地虚拟机进行模拟。
//例如，“balance（）”方法只需要检索
//从存储器中，因此不需要发送到
//网络。像“transact”这样的方法需要一个tx，因此
//标记为“真”。
//输入指定此给定方法所需的输入参数。
type Method struct {
	Name    string
	Const   bool
	Inputs  Arguments
	Outputs Arguments
}

//SIG根据ABI规范返回方法字符串签名。
//
//例子
//
//函数foo（uint32 a，int b）=“foo（uint32，int256）”
//
//请注意，“int”代替其规范表示“int256”
func (method Method) Sig() string {
	types := make([]string, len(method.Inputs))
	for i, input := range method.Inputs {
		types[i] = input.Type.String()
	}
	return fmt.Sprintf("%v(%v)", method.Name, strings.Join(types, ","))
}

func (method Method) String() string {
	inputs := make([]string, len(method.Inputs))
	for i, input := range method.Inputs {
		inputs[i] = fmt.Sprintf("%v %v", input.Type, input.Name)
	}
	outputs := make([]string, len(method.Outputs))
	for i, output := range method.Outputs {
		outputs[i] = output.Type.String()
		if len(output.Name) > 0 {
			outputs[i] += fmt.Sprintf(" %v", output.Name)
		}
	}
	constant := ""
	if method.Const {
		constant = "constant "
	}
	return fmt.Sprintf("function %v(%v) %sreturns(%v)", method.Name, strings.Join(inputs, ", "), constant, strings.Join(outputs, ", "))
}

func (method Method) Id() []byte {
	return crypto.Keccak256([]byte(method.Sig()))[:4]
}

