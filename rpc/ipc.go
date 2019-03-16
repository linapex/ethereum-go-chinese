
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450108938784768>


package rpc

import (
	"context"
	"net"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/netutil"
)

//servelistener接受L上的连接，并为它们提供JSON-RPC。
func (srv *Server) ServeListener(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if netutil.IsTemporaryError(err) {
			log.Warn("IPC accept error", "err", err)
			continue
		} else if err != nil {
			return err
		}
		log.Trace("IPC accepted connection")
		go srv.ServeCodec(NewJSONCodec(conn), OptionMethodInvocation|OptionSubscriptions)
	}
}

//DialIPC创建一个新的连接到给定端点的IPC客户端。在Unix上，它假定
//端点是指向Unix套接字的完整路径，而Windows端点是
//命名管道的标识符。
//
//上下文用于建立初始连接。它不
//影响与客户的后续交互。
func DialIPC(ctx context.Context, endpoint string) (*Client, error) {
	return newClient(ctx, func(ctx context.Context) (net.Conn, error) {
		return newIPCConnection(ctx, endpoint)
	})
}

