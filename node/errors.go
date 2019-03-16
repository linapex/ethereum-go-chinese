
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450102370504704>


package node

import (
	"errors"
	"fmt"
	"reflect"
	"syscall"
)

var (
	ErrDatadirUsed    = errors.New("datadir already used by another process")
	ErrNodeStopped    = errors.New("node not started")
	ErrNodeRunning    = errors.New("node already running")
	ErrServiceUnknown = errors.New("unknown service")

	datadirInUseErrnos = map[uint]bool{11: true, 32: true, 35: true}
)

func convertFileLockError(err error) error {
	if errno, ok := err.(syscall.Errno); ok && datadirInUseErrnos[uint(errno)] {
		return ErrDatadirUsed
	}
	return err
}

//如果注册的服务在节点启动期间返回DuplicateServiceError
//构造函数返回已启动的同一类型的服务。
type DuplicateServiceError struct {
	Kind reflect.Type
}

//错误生成重复服务错误的文本表示形式。
func (e *DuplicateServiceError) Error() string {
	return fmt.Sprintf("duplicate service: %v", e.Kind)
}

//如果节点未能停止其任何注册节点，则返回stopError。
//服务或其本身。
type StopError struct {
	Server   error
	Services map[reflect.Type]error
}

//错误生成停止错误的文本表示形式。
func (e *StopError) Error() string {
	return fmt.Sprintf("server: %v, services: %v", e.Server, e.Services)
}

