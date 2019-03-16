
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450120867385344>


//package rpc实现一个连接到集中模拟存储的rpc客户机。
//中心化模拟存储可以是任何其他模拟存储实现，即
//以mockstore名称注册到以太坊RPC服务器。定义的方法
//mock.globalStore与rpc使用的相同。例子：
//
//服务器：=rpc.newserver（）
//server.registername（“mockstore”，mem.newGlobalStore（））
package rpc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/swarm/log"
	"github.com/ethereum/go-ethereum/swarm/storage/mock"
)

//GlobalStore是一个连接到中央模拟商店的rpc.client。
//关闭GlobalStore实例需要释放RPC客户端资源。
type GlobalStore struct {
	client *rpc.Client
}

//NewGlobalStore创建了一个新的GlobalStore实例。
func NewGlobalStore(client *rpc.Client) *GlobalStore {
	return &GlobalStore{
		client: client,
	}
}

//关闭关闭RPC客户端。
func (s *GlobalStore) Close() error {
	s.client.Close()
	return nil
}

//new nodestore返回一个新的nodestore实例，用于检索和存储
//仅对地址为的节点进行数据块处理。
func (s *GlobalStore) NewNodeStore(addr common.Address) *mock.NodeStore {
	return mock.NewNodeStore(addr, s)
}

//get调用rpc服务器的get方法。
func (s *GlobalStore) Get(addr common.Address, key []byte) (data []byte, err error) {
	err = s.client.Call(&data, "mockStore_get", addr, key)
	if err != nil && err.Error() == "not found" {
//传递模拟包的错误值，而不是一个RPC错误
		return data, mock.ErrNotFound
	}
	return data, err
}

//将一个Put方法调用到RPC服务器。
func (s *GlobalStore) Put(addr common.Address, key []byte, data []byte) error {
	err := s.client.Call(nil, "mockStore_put", addr, key, data)
	return err
}

//delete向rpc服务器调用delete方法。
func (s *GlobalStore) Delete(addr common.Address, key []byte) error {
	err := s.client.Call(nil, "mockStore_delete", addr, key)
	return err
}

//haskey向RPC服务器调用haskey方法。
func (s *GlobalStore) HasKey(addr common.Address, key []byte) bool {
	var has bool
	if err := s.client.Call(&has, "mockStore_hasKey", addr, key); err != nil {
		log.Error(fmt.Sprintf("mock store HasKey: addr %s, key %064x: %v", addr, key, err))
		return false
	}
	return has
}

