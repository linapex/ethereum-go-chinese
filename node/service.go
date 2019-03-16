
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450102638940160>


package node

import (
	"reflect"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

//ServiceContext是从继承的与服务无关的选项的集合
//协议栈，它被传递给所有的构造器以供选择使用；
//以及在服务环境中操作的实用方法。
type ServiceContext struct {
	config         *Config
services       map[reflect.Type]Service //已建服务索引
EventMux       *event.TypeMux           //用于分离通知的事件多路复用器
AccountManager *accounts.Manager        //节点创建的客户经理。
}

//opendatabase打开具有给定名称的现有数据库（或创建一个
//如果在节点的数据目录中找不到上一个）。如果
//节点是短暂的，返回内存数据库。
func (ctx *ServiceContext) OpenDatabase(name string, cache int, handles int) (ethdb.Database, error) {
	if ctx.config.DataDir == "" {
		return ethdb.NewMemDatabase(), nil
	}
	db, err := ethdb.NewLDBDatabase(ctx.config.ResolvePath(name), cache, handles)
	if err != nil {
		return nil, err
	}
	return db, nil
}

//resolvepath将用户路径解析为数据目录（如果该路径是相对的）
//如果用户实际使用持久存储。它将返回空字符串
//对于临时存储和用户自己的绝对路径输入。
func (ctx *ServiceContext) ResolvePath(path string) string {
	return ctx.config.ResolvePath(path)
}

//服务检索当前正在运行的特定类型的注册服务。
func (ctx *ServiceContext) Service(service interface{}) error {
	element := reflect.ValueOf(service).Elem()
	if running, ok := ctx.services[element.Type()]; ok {
		element.Set(reflect.ValueOf(running))
		return nil
	}
	return ErrServiceUnknown
}

//ServiceConstructor是需要
//已注册服务实例化。
type ServiceConstructor func(ctx *ServiceContext) (Service, error)

//服务是可以注册到节点中的单个协议。
//
//笔记：
//
//•将服务生命周期管理委托给节点。允许服务
//创建时初始化自身，但不应在
//启动方法。
//
//•不需要重新启动逻辑，因为节点将创建新实例
//每次启动服务时。
type Service interface {
//协议检索服务希望启动的P2P协议。
	Protocols() []p2p.Protocol

//API检索服务提供的RPC描述符列表
	APIs() []rpc.API

//在构建所有服务和网络之后调用Start
//层也被初始化以生成服务所需的任何goroutine。
	Start(server *p2p.Server) error

//stop终止属于服务的所有goroutine，阻塞直到它们
//全部终止。
	Stop() error
}

