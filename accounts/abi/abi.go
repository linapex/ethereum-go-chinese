
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:30</date>
//</624450059311779840>


package abi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

//ABI保存有关合同上下文的信息，并提供
//可调用方法。它将允许您键入check函数调用和
//相应地打包数据。
type ABI struct {
	Constructor Method
	Methods     map[string]Method
	Events      map[string]Event
}

//JSON返回已分析的ABI接口，如果失败则返回错误。
func JSON(reader io.Reader) (ABI, error) {
	dec := json.NewDecoder(reader)

	var abi ABI
	if err := dec.Decode(&abi); err != nil {
		return ABI{}, err
	}

	return abi, nil
}

//打包给定的方法名以符合ABI。方法调用的数据
//将由方法“id”、“args0”、“arg1…”组成。阿尔金方法ID由
//其中4个字节和参数都是32个字节。
//方法ID是从
//方法字符串签名。（签名=baz（uint32，string32））
func (abi ABI) Pack(name string, args ...interface{}) ([]byte, error) {
//获取所请求方法的ABI
	if name == "" {
//构造函数
		arguments, err := abi.Constructor.Inputs.Pack(args...)
		if err != nil {
			return nil, err
		}
		return arguments, nil
	}
	method, exist := abi.Methods[name]
	if !exist {
		return nil, fmt.Errorf("method '%s' not found", name)
	}
	arguments, err := method.Inputs.Pack(args...)
	if err != nil {
		return nil, err
	}
//如果不是构造函数，也打包方法ID并返回
	return append(method.Id(), arguments...), nil
}

//根据ABI规范以V为单位解包输出
func (abi ABI) Unpack(v interface{}, name string, output []byte) (err error) {
	if len(output) == 0 {
		return fmt.Errorf("abi: unmarshalling empty output")
	}
//既然不能命名与契约和事件的冲突，
//我们需要决定是调用方法还是调用事件
	if method, ok := abi.Methods[name]; ok {
		if len(output)%32 != 0 {
			return fmt.Errorf("abi: improperly formatted output: %s - Bytes: [%+v]", string(output), output)
		}
		return method.Outputs.Unpack(v, output)
	} else if event, ok := abi.Events[name]; ok {
		return event.Inputs.Unpack(v, output)
	}
	return fmt.Errorf("abi: could not locate named method or event")
}

//unmashaljson实现json.unmasheler接口
func (abi *ABI) UnmarshalJSON(data []byte) error {
	var fields []struct {
		Type      string
		Name      string
		Constant  bool
		Anonymous bool
		Inputs    []Argument
		Outputs   []Argument
	}

	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	abi.Methods = make(map[string]Method)
	abi.Events = make(map[string]Event)
	for _, field := range fields {
		switch field.Type {
		case "constructor":
			abi.Constructor = Method{
				Inputs: field.Inputs,
			}
//空默认为根据ABI规范的功能
		case "function", "":
			abi.Methods[field.Name] = Method{
				Name:    field.Name,
				Const:   field.Constant,
				Inputs:  field.Inputs,
				Outputs: field.Outputs,
			}
		case "event":
			abi.Events[field.Name] = Event{
				Name:      field.Name,
				Anonymous: field.Anonymous,
				Inputs:    field.Inputs,
			}
		}
	}

	return nil
}

//method by id按4字节的ID查找方法
//如果未找到，则返回nil
func (abi *ABI) MethodById(sigdata []byte) (*Method, error) {
	if len(sigdata) < 4 {
		return nil, fmt.Errorf("data too short (% bytes) for abi method lookup", len(sigdata))
	}
	for _, method := range abi.Methods {
		if bytes.Equal(method.Id(), sigdata[:4]) {
			return &method, nil
		}
	}
	return nil, fmt.Errorf("no method with id: %#x", sigdata[:4])
}

