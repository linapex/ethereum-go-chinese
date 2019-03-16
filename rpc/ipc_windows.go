
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450109081391104>


//+构建窗口

package rpc

import (
	"context"
	"net"
	"time"

	"gopkg.in/natefinch/npipe.v2"
)

//如果拨号上下文没有截止时间，则使用此选项。它比
//默认拨号超时，因为命名管道是本地的，不需要等待太长时间。
const defaultPipeDialTimeout = 2 * time.Second

//ipclisten将在给定的端点上创建命名管道。
func ipcListen(endpoint string) (net.Listener, error) {
	return npipe.Listen(endpoint)
}

//NewIPCConnection将连接到具有给定端点作为名称的命名管道。
func newIPCConnection(ctx context.Context, endpoint string) (net.Conn, error) {
	timeout := defaultPipeDialTimeout
	if deadline, ok := ctx.Deadline(); ok {
		timeout = deadline.Sub(time.Now())
		if timeout < 0 {
			timeout = 0
		}
	}
	return npipe.DialTimeout(endpoint, timeout)
}

