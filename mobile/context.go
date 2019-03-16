
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450101049298944>


//包含golang.org/x/net/context包中要支持的所有包装器
//移动平台上的客户端上下文管理。

package geth

import (
	"context"
	"time"
)

//上下文在API中包含截止日期、取消信号和其他值。
//边界。
type Context struct {
	context context.Context
	cancel  context.CancelFunc
}

//NewContext返回非零的空上下文。从来没有取消过，没有
//价值观，没有最后期限。它通常由主功能使用，
//初始化和测试，以及作为传入请求的顶级上下文。
func NewContext() *Context {
	return &Context{
		context: context.Background(),
	}
}

//withcancel返回具有取消机制的原始上下文的副本
//包括。
//
//取消此上下文将释放与其关联的资源，因此代码应该
//在此上下文中运行的操作完成后立即调用Cancel。
func (c *Context) WithCancel() *Context {
	child, cancel := context.WithCancel(c.context)
	return &Context{
		context: child,
		cancel:  cancel,
	}
}

//WithDeadline返回原始上下文的副本，调整了截止时间
//不迟于规定时间。
//
//取消此上下文将释放与其关联的资源，因此代码应该
//在此上下文中运行的操作完成后立即调用Cancel。
func (c *Context) WithDeadline(sec int64, nsec int64) *Context {
	child, cancel := context.WithDeadline(c.context, time.Unix(sec, nsec))
	return &Context{
		context: child,
		cancel:  cancel,
	}
}

//WithTimeout返回原始上下文的副本，并调整最后期限
//不迟于现在+指定的持续时间。
//
//取消此上下文将释放与其关联的资源，因此代码应该
//在此上下文中运行的操作完成后立即调用Cancel。
func (c *Context) WithTimeout(nsec int64) *Context {
	child, cancel := context.WithTimeout(c.context, time.Duration(nsec))
	return &Context{
		context: child,
		cancel:  cancel,
	}
}

