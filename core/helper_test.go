
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450079155032064>


package core

import (
	"container/list"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
)

//实施我们的ethtest管理器
type TestManager struct {
//StateManager*状态管理器
	eventMux *event.TypeMux

	db         ethdb.Database
	txPool     *TxPool
	blockChain *BlockChain
	Blocks     []*types.Block
}

func (tm *TestManager) IsListening() bool {
	return false
}

func (tm *TestManager) IsMining() bool {
	return false
}

func (tm *TestManager) PeerCount() int {
	return 0
}

func (tm *TestManager) Peers() *list.List {
	return list.New()
}

func (tm *TestManager) BlockChain() *BlockChain {
	return tm.blockChain
}

func (tm *TestManager) TxPool() *TxPool {
	return tm.txPool
}

//func（tm*testmanager）statemanager（）*statemanager_
//返回tm.statemanager
//}

func (tm *TestManager) EventMux() *event.TypeMux {
	return tm.eventMux
}

//func（tm*testmanager）keymanager（）*crypto.keymanager_
//返回零
//}

func (tm *TestManager) Db() ethdb.Database {
	return tm.db
}

func NewTestManager() *TestManager {
	testManager := &TestManager{}
	testManager.eventMux = new(event.TypeMux)
	testManager.db = ethdb.NewMemDatabase()
//testmanager.txpool=newtxpool（testmanager）
//testmanager.blockchain=newblockchain（testmanager）
//testmanager.statemanager=新状态管理器（testmanager）
	return testManager
}

