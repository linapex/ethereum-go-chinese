
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450108720680960>


package rpc

import "fmt"

//请求用于未知服务
type methodNotFoundError struct {
	service string
	method  string
}

func (e *methodNotFoundError) ErrorCode() int { return -32601 }

func (e *methodNotFoundError) Error() string {
	return fmt.Sprintf("The method %s%s%s does not exist/is not available", e.service, serviceMethodSeparator, e.method)
}

//收到的消息不是有效的请求
type invalidRequestError struct{ message string }

func (e *invalidRequestError) ErrorCode() int { return -32600 }

func (e *invalidRequestError) Error() string { return e.message }

//接收到的消息无效
type invalidMessageError struct{ message string }

func (e *invalidMessageError) ErrorCode() int { return -32700 }

func (e *invalidMessageError) Error() string { return e.message }

//无法解码提供的参数，或参数数目无效
type invalidParamsError struct{ message string }

func (e *invalidParamsError) ErrorCode() int { return -32602 }

func (e *invalidParamsError) Error() string { return e.message }

//逻辑错误，回调返回错误
type callbackError struct{ message string }

func (e *callbackError) ErrorCode() int { return -32000 }

func (e *callbackError) Error() string { return e.message }

//在服务器发出停止后收到请求时发出。
type shutdownError struct{}

func (e *shutdownError) ErrorCode() int { return -32000 }

func (e *shutdownError) Error() string { return "server is shutting down" }

