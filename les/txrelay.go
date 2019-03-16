
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450096301346816>


package les

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ltrInfo struct {
	tx     *types.Transaction
	sentTo map[*peer]struct{}
}

type LesTxRelay struct {
	txSent       map[common.Hash]*ltrInfo
	txPending    map[common.Hash]struct{}
	ps           *peerSet
	peerList     []*peer
	peerStartPos int
	lock         sync.RWMutex

	reqDist *requestDistributor
}

func NewLesTxRelay(ps *peerSet, reqDist *requestDistributor) *LesTxRelay {
	r := &LesTxRelay{
		txSent:    make(map[common.Hash]*ltrInfo),
		txPending: make(map[common.Hash]struct{}),
		ps:        ps,
		reqDist:   reqDist,
	}
	ps.notify(r)
	return r
}

func (self *LesTxRelay) registerPeer(p *peer) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.peerList = self.ps.AllPeers()
}

func (self *LesTxRelay) unregisterPeer(p *peer) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.peerList = self.ps.AllPeers()
}

//send最多向给定数量的对等方发送一个事务列表
//一次，不要将任何特定事务重新发送到同一个对等机两次
func (self *LesTxRelay) send(txs types.Transactions, count int) {
	sendTo := make(map[*peer]types.Transactions)

self.peerStartPos++ //旋转对等列表的起始位置
	if self.peerStartPos >= len(self.peerList) {
		self.peerStartPos = 0
	}

	for _, tx := range txs {
		hash := tx.Hash()
		ltr, ok := self.txSent[hash]
		if !ok {
			ltr = &ltrInfo{
				tx:     tx,
				sentTo: make(map[*peer]struct{}),
			}
			self.txSent[hash] = ltr
			self.txPending[hash] = struct{}{}
		}

		if len(self.peerList) > 0 {
			cnt := count
			pos := self.peerStartPos
			for {
				peer := self.peerList[pos]
				if _, ok := ltr.sentTo[peer]; !ok {
					sendTo[peer] = append(sendTo[peer], tx)
					ltr.sentTo[peer] = struct{}{}
					cnt--
				}
				if cnt == 0 {
break //发送到所需的对等数
				}
				pos++
				if pos == len(self.peerList) {
					pos = 0
				}
				if pos == self.peerStartPos {
break //已尝试所有可用的对等机
				}
			}
		}
	}

	for p, list := range sendTo {
		pp := p
		ll := list

		reqID := genReqID()
		rq := &distReq{
			getCost: func(dp distPeer) uint64 {
				peer := dp.(*peer)
				return peer.GetRequestCost(SendTxMsg, len(ll))
			},
			canSend: func(dp distPeer) bool {
				return dp.(*peer) == pp
			},
			request: func(dp distPeer) func() {
				peer := dp.(*peer)
				cost := peer.GetRequestCost(SendTxMsg, len(ll))
				peer.fcServer.QueueRequest(reqID, cost)
				return func() { peer.SendTxs(reqID, cost, ll) }
			},
		}
		self.reqDist.queue(rq)
	}
}

func (self *LesTxRelay) Send(txs types.Transactions) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.send(txs, 3)
}

func (self *LesTxRelay) NewHead(head common.Hash, mined []common.Hash, rollback []common.Hash) {
	self.lock.Lock()
	defer self.lock.Unlock()

	for _, hash := range mined {
		delete(self.txPending, hash)
	}

	for _, hash := range rollback {
		self.txPending[hash] = struct{}{}
	}

	if len(self.txPending) > 0 {
		txs := make(types.Transactions, len(self.txPending))
		i := 0
		for hash := range self.txPending {
			txs[i] = self.txSent[hash].tx
			i++
		}
		self.send(txs, 1)
	}
}

func (self *LesTxRelay) Discard(hashes []common.Hash) {
	self.lock.Lock()
	defer self.lock.Unlock()

	for _, hash := range hashes {
		delete(self.txSent, hash)
		delete(self.txPending, hash)
	}
}

