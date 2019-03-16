
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450066081386496>


//此文件包含用于与Trezor硬件交互的实现
//钱包。有线协议规范可在Satoshilabs网站上找到：
//https://doc.satoshilabs.com/trezor-tech/api-protobuf.html网站

//go:generate protoc--go_out=import_path=trezor:。types.proto消息.proto

//包trezor在go中包含有线协议包装器。
package trezor

import (
	"reflect"

	"github.com/golang/protobuf/proto"
)

//类型返回特定消息的协议缓冲区类型号。如果
//消息为零，此方法会恐慌！
func Type(msg proto.Message) uint16 {
	return uint16(MessageType_value["MessageType_"+reflect.TypeOf(msg).Elem().Name()])
}

//name返回特定协议缓冲区的友好消息类型名称
//类型号。
func Name(kind uint16) string {
	name := MessageType_name[int32(kind)]
	if len(name) < 12 {
		return name
	}
	return name[12:]
}

