
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450105537204224>


//+构建窗口

package netutil

import (
	"net"
	"os"
	"syscall"
)

const _WSAEMSGSIZE = syscall.Errno(10040)

//ispackettoobig报告err是否指示UDP数据包没有
//安装接收缓冲区。在Windows上，wsarecvfrom返回
//对wsaemsgsize进行编码，如果发生这种情况，则没有数据。
func isPacketTooBig(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if scErr, ok := opErr.Err.(*os.SyscallError); ok {
			return scErr.Err == _WSAEMSGSIZE
		}
		return opErr.Err == _WSAEMSGSIZE
	}
	return false
}

