
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450107223314432>


package pipes

import (
	"net"
)

//net pipe在返回错误的签名中包装net.pipe
func NetPipe() (net.Conn, net.Conn, error) {
	p1, p2 := net.Pipe()
	return p1, p2, nil
}

//tcp pipe基于本地主机tcp套接字创建进程内全双工管道
func TCPPipe() (net.Conn, net.Conn, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, nil, err
	}
	defer l.Close()

	var aconn net.Conn
	aerr := make(chan error, 1)
	go func() {
		var err error
		aconn, err = l.Accept()
		aerr <- err
	}()

	dconn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		<-aerr
		return nil, nil, err
	}
	if err := <-aerr; err != nil {
		dconn.Close()
		return nil, nil, err
	}
	return aconn, dconn, nil
}

