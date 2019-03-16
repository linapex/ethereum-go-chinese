
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450109622456320>


package rpc

import (
	"context"
	"errors"
	"net"
	"os"
	"time"
)

//dialstdio在stdin/stdout上创建客户端。
func DialStdIO(ctx context.Context) (*Client, error) {
	return newClient(ctx, func(_ context.Context) (net.Conn, error) {
		return stdioConn{}, nil
	})
}

type stdioConn struct{}

func (io stdioConn) Read(b []byte) (n int, err error) {
	return os.Stdin.Read(b)
}

func (io stdioConn) Write(b []byte) (n int, err error) {
	return os.Stdout.Write(b)
}

func (io stdioConn) Close() error {
	return nil
}

func (io stdioConn) LocalAddr() net.Addr {
	return &net.UnixAddr{Name: "stdio", Net: "stdio"}
}

func (io stdioConn) RemoteAddr() net.Addr {
	return &net.UnixAddr{Name: "stdio", Net: "stdio"}
}

func (io stdioConn) SetDeadline(t time.Time) error {
	return &net.OpError{Op: "set", Net: "stdio", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

func (io stdioConn) SetReadDeadline(t time.Time) error {
	return &net.OpError{Op: "set", Net: "stdio", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

func (io stdioConn) SetWriteDeadline(t time.Time) error {
	return &net.OpError{Op: "set", Net: "stdio", Source: nil, Addr: nil, Err: errors.New("deadline not supported")}
}

