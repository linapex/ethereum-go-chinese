
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450093705072640>


//包流控制实现客户端流控制机制
package flowcontrol

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
)

const fcTimeConst = time.Millisecond

type ServerParams struct {
	BufLimit, MinRecharge uint64
}

type ClientNode struct {
	params   *ServerParams
	bufValue uint64
	lastTime mclock.AbsTime
	lock     sync.Mutex
	cm       *ClientManager
	cmNode   *cmNode
}

func NewClientNode(cm *ClientManager, params *ServerParams) *ClientNode {
	node := &ClientNode{
		cm:       cm,
		params:   params,
		bufValue: params.BufLimit,
		lastTime: mclock.Now(),
	}
	node.cmNode = cm.addNode(node)
	return node
}

func (peer *ClientNode) Remove(cm *ClientManager) {
	cm.removeNode(peer.cmNode)
}

func (peer *ClientNode) recalcBV(time mclock.AbsTime) {
	dt := uint64(time - peer.lastTime)
	if time < peer.lastTime {
		dt = 0
	}
	peer.bufValue += peer.params.MinRecharge * dt / uint64(fcTimeConst)
	if peer.bufValue > peer.params.BufLimit {
		peer.bufValue = peer.params.BufLimit
	}
	peer.lastTime = time
}

func (peer *ClientNode) AcceptRequest() (uint64, bool) {
	peer.lock.Lock()
	defer peer.lock.Unlock()

	time := mclock.Now()
	peer.recalcBV(time)
	return peer.bufValue, peer.cm.accept(peer.cmNode, time)
}

func (peer *ClientNode) RequestProcessed(cost uint64) (bv, realCost uint64) {
	peer.lock.Lock()
	defer peer.lock.Unlock()

	time := mclock.Now()
	peer.recalcBV(time)
	peer.bufValue -= cost
	rcValue, rcost := peer.cm.processed(peer.cmNode, time)
	if rcValue < peer.params.BufLimit {
		bv := peer.params.BufLimit - rcValue
		if bv > peer.bufValue {
			peer.bufValue = bv
		}
	}
	return peer.bufValue, rcost
}

type ServerNode struct {
	bufEstimate uint64
	lastTime    mclock.AbsTime
	params      *ServerParams
sumCost     uint64            //发送到此服务器的请求成本总和
pending     map[uint64]uint64 //值=发送给定需求后的总成本
	lock        sync.RWMutex
}

func NewServerNode(params *ServerParams) *ServerNode {
	return &ServerNode{
		bufEstimate: params.BufLimit,
		lastTime:    mclock.Now(),
		params:      params,
		pending:     make(map[uint64]uint64),
	}
}

func (peer *ServerNode) recalcBLE(time mclock.AbsTime) {
	dt := uint64(time - peer.lastTime)
	if time < peer.lastTime {
		dt = 0
	}
	peer.bufEstimate += peer.params.MinRecharge * dt / uint64(fcTimeConst)
	if peer.bufEstimate > peer.params.BufLimit {
		peer.bufEstimate = peer.params.BufLimit
	}
	peer.lastTime = time
}

//当估计的缓冲区值较低时，将安全裕度添加到流控制等待时间中。
const safetyMargin = time.Millisecond

func (peer *ServerNode) canSend(maxCost uint64) (time.Duration, float64) {
	peer.recalcBLE(mclock.Now())
	maxCost += uint64(safetyMargin) * peer.params.MinRecharge / uint64(fcTimeConst)
	if maxCost > peer.params.BufLimit {
		maxCost = peer.params.BufLimit
	}
	if peer.bufEstimate >= maxCost {
		return 0, float64(peer.bufEstimate-maxCost) / float64(peer.params.BufLimit)
	}
	return time.Duration((maxCost - peer.bufEstimate) * uint64(fcTimeConst) / peer.params.MinRecharge), 0
}

//cansend返回发送请求前所需的最小等待时间
//以给定的最大估计成本。第二个返回值是相对值
//发送请求后的估计缓冲区级别（除以buflimit）。
func (peer *ServerNode) CanSend(maxCost uint64) (time.Duration, float64) {
	peer.lock.RLock()
	defer peer.lock.RUnlock()

	return peer.canSend(maxCost)
}

//当请求已分配给给定的
//服务器节点，然后将其放入发送队列。必须要求
//以与发出QueuerRequest调用相同的顺序发送。
func (peer *ServerNode) QueueRequest(reqID, maxCost uint64) {
	peer.lock.Lock()
	defer peer.lock.Unlock()

	peer.bufEstimate -= maxCost
	peer.sumCost += maxCost
	peer.pending[reqID] = peer.sumCost
}

//gotmreply根据包含在
//最新的请求回复。
func (peer *ServerNode) GotReply(reqID, bv uint64) {

	peer.lock.Lock()
	defer peer.lock.Unlock()

	if bv > peer.params.BufLimit {
		bv = peer.params.BufLimit
	}
	sc, ok := peer.pending[reqID]
	if !ok {
		return
	}
	delete(peer.pending, reqID)
	cc := peer.sumCost - sc
	peer.bufEstimate = 0
	if bv > cc {
		peer.bufEstimate = bv - cc
	}
	peer.lastTime = mclock.Now()
}

