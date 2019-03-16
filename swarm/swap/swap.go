
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450121551056896>


package swap

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/protocols"
	"github.com/ethereum/go-ethereum/swarm/log"
	"github.com/ethereum/go-ethereum/swarm/state"
)

//swap swarm会计协议
//点对点小额支付系统
//一个节点与每一个对等节点保持单个平衡。
//只有有价格的消息才会被计入
type Swap struct {
stateStore state.Store        //需要Statestore才能在会话之间保持平衡
lock       sync.RWMutex       //锁定余额
balances   map[enode.ID]int64 //每个对等点的平衡图
}

//新建-交换构造函数
func New(stateStore state.Store) (swap *Swap) {
	swap = &Swap{
		stateStore: stateStore,
		balances:   make(map[enode.ID]int64),
	}
	return
}

//swap实现协议。平衡接口
//添加是（唯一）会计功能
func (s *Swap) Add(amount int64, peer *protocols.Peer) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

//从状态存储加载现有余额
	err = s.loadState(peer)
	if err != nil && err != state.ErrNotFound {
		return
	}
//调整平衡
//如果金额为负数，则会减少，否则会增加
	s.balances[peer.ID()] += amount
//将新余额保存到状态存储
	peerBalance := s.balances[peer.ID()]
	err = s.stateStore.Put(peer.ID().String(), &peerBalance)

	log.Debug(fmt.Sprintf("balance for peer %s: %s", peer.ID().String(), strconv.FormatInt(peerBalance, 10)))
	return err
}

//GetPeerBalance返回给定对等机的余额
func (swap *Swap) GetPeerBalance(peer enode.ID) (int64, error) {
	swap.lock.RLock()
	defer swap.lock.RUnlock()
	if p, ok := swap.balances[peer]; ok {
		return p, nil
	}
	return 0, errors.New("Peer not found")
}

//状态存储的负载平衡（持久）
func (s *Swap) loadState(peer *protocols.Peer) (err error) {
	var peerBalance int64
	peerID := peer.ID()
//仅当当前实例没有此对等方的
//内存平衡
	if _, ok := s.balances[peerID]; !ok {
		err = s.stateStore.Get(peerID.String(), &peerBalance)
		s.balances[peerID] = peerBalance
	}
	return
}

//清理交换
func (swap *Swap) Close() {
	swap.stateStore.Close()
}

