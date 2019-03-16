
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450102521499648>


package node_test

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rpc"
)

//sampleService是一个普通的网络服务，可以连接到
//生命周期管理。
//
//实现node.service需要以下方法：
//-protocols（）[]p2p.protocol-服务可以通信的devp2p协议
//-apis（）[]rpc.api-服务希望在rpc通道上公开的api方法
//-start（）错误-节点准备启动服务时调用的方法
//-stop（）错误-节点终止服务时调用的方法
type SampleService struct{}

func (s *SampleService) Protocols() []p2p.Protocol { return nil }
func (s *SampleService) APIs() []rpc.API           { return nil }
func (s *SampleService) Start(*p2p.Server) error   { fmt.Println("Service starting..."); return nil }
func (s *SampleService) Stop() error               { fmt.Println("Service stopping..."); return nil }

func ExampleService() {
//创建一个网络节点以运行具有默认值的协议。
	stack, err := node.New(&node.Config{})
	if err != nil {
		log.Fatalf("Failed to create network node: %v", err)
	}
//创建并注册一个简单的网络服务。这是通过定义完成的
//将实例化node.service的node.serviceconstructor。原因
//工厂方法的方法是支持服务重新启动，而不依赖
//单个实现对此类操作的支持。
	constructor := func(context *node.ServiceContext) (node.Service, error) {
		return new(SampleService), nil
	}
	if err := stack.Register(constructor); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
//启动整个协议栈，重新启动并终止
	if err := stack.Start(); err != nil {
		log.Fatalf("Failed to start the protocol stack: %v", err)
	}
	if err := stack.Restart(); err != nil {
		log.Fatalf("Failed to restart the protocol stack: %v", err)
	}
	if err := stack.Stop(); err != nil {
		log.Fatalf("Failed to stop the protocol stack: %v", err)
	}
//输出：
//服务正在启动…
//服务正在停止…
//服务正在启动…
//服务正在停止…
}

