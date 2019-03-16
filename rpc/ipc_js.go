
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450108993310720>


//+构建JS

package rpc

import (
	"context"
	"errors"
	"net"
)

var errNotSupported = errors.New("rpc: not supported")

//ipclisten将在给定的端点上创建命名管道。
func ipcListen(endpoint string) (net.Listener, error) {
	return nil, errNotSupported
}

//NewIPCConnection将连接到具有给定端点作为名称的命名管道。
func newIPCConnection(ctx context.Context, endpoint string) (net.Conn, error) {
	return nil, errNotSupported
}

