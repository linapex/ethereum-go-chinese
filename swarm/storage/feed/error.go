
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450118757650432>


package feed

import (
	"fmt"
)

const (
	ErrInit = iota
	ErrNotFound
	ErrIO
	ErrUnauthorized
	ErrInvalidValue
	ErrDataOverflow
	ErrNothingToReturn
	ErrCorruptData
	ErrInvalidSignature
	ErrNotSynced
	ErrPeriodDepth
	ErrCnt
)

//错误是用于Swarm提要的类型化错误对象
type Error struct {
	code int
	err  string
}

//错误实现错误接口
func (e *Error) Error() string {
	return e.err
}

//代码返回错误代码
//错误代码在feeds包中的error.go文件中枚举。
func (e *Error) Code() int {
	return e.code
}

//newError用指定的代码和自定义错误消息创建一个新的swarm feeds错误对象
func NewError(code int, s string) error {
	if code < 0 || code >= ErrCnt {
		panic("no such error code!")
	}
	r := &Error{
		err: s,
	}
	switch code {
	case ErrNotFound, ErrIO, ErrUnauthorized, ErrInvalidValue, ErrDataOverflow, ErrNothingToReturn, ErrInvalidSignature, ErrNotSynced, ErrPeriodDepth, ErrCorruptData:
		r.code = code
	}
	return r
}

//newerrorf是newerror的一个方便版本，它包含了printf样式的格式。
func NewErrorf(code int, format string, args ...interface{}) error {
	return NewError(code, fmt.Sprintf(format, args...))
}

