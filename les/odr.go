
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450095571537920>


package les

import (
	"context"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/log"
)

//lesodr实现light.odrbackend
type LesOdr struct {
	db                                         ethdb.Database
	indexerConfig                              *light.IndexerConfig
	chtIndexer, bloomTrieIndexer, bloomIndexer *core.ChainIndexer
	retriever                                  *retrieveManager
	stop                                       chan struct{}
}

func NewLesOdr(db ethdb.Database, config *light.IndexerConfig, retriever *retrieveManager) *LesOdr {
	return &LesOdr{
		db:            db,
		indexerConfig: config,
		retriever:     retriever,
		stop:          make(chan struct{}),
	}
}

//停止取消所有挂起的检索
func (odr *LesOdr) Stop() {
	close(odr.stop)
}

//数据库返回后备数据库
func (odr *LesOdr) Database() ethdb.Database {
	return odr.db
}

//setindexers将必要的链索引器添加到ODR后端
func (odr *LesOdr) SetIndexers(chtIndexer, bloomTrieIndexer, bloomIndexer *core.ChainIndexer) {
	odr.chtIndexer = chtIndexer
	odr.bloomTrieIndexer = bloomTrieIndexer
	odr.bloomIndexer = bloomIndexer
}

//返回CHT链索引器
func (odr *LesOdr) ChtIndexer() *core.ChainIndexer {
	return odr.chtIndexer
}

//BloomTrieIndexer返回BloomTrie链索引器
func (odr *LesOdr) BloomTrieIndexer() *core.ChainIndexer {
	return odr.bloomTrieIndexer
}

//BloomIndexer返回BloomBits链索引器
func (odr *LesOdr) BloomIndexer() *core.ChainIndexer {
	return odr.bloomIndexer
}

//indexer config返回索引器配置。
func (odr *LesOdr) IndexerConfig() *light.IndexerConfig {
	return odr.indexerConfig
}

const (
	MsgBlockBodies = iota
	MsgCode
	MsgReceipts
	MsgProofsV1
	MsgProofsV2
	MsgHeaderProofs
	MsgHelperTrieProofs
)

//msg对为请求传递答复数据的les消息进行编码
type Msg struct {
	MsgType int
	ReqID   uint64
	Obj     interface{}
}

//retrieve尝试从les网络获取对象。
//如果网络检索成功，它将对象存储在本地数据库中。
func (odr *LesOdr) Retrieve(ctx context.Context, req light.OdrRequest) (err error) {
	lreq := LesRequest(req)

	reqID := genReqID()
	rq := &distReq{
		getCost: func(dp distPeer) uint64 {
			return lreq.GetCost(dp.(*peer))
		},
		canSend: func(dp distPeer) bool {
			p := dp.(*peer)
			return lreq.CanSend(p)
		},
		request: func(dp distPeer) func() {
			p := dp.(*peer)
			cost := lreq.GetCost(p)
			p.fcServer.QueueRequest(reqID, cost)
			return func() { lreq.Request(reqID, p) }
		},
	}

	if err = odr.retriever.retrieve(ctx, reqID, rq, func(p distPeer, msg *Msg) error { return lreq.Validate(odr.db, msg) }, odr.stop); err == nil {
//retrieved from network, store in db
		req.StoreResult(odr.db)
	} else {
		log.Debug("Failed to retrieve data from network", "err", err)
	}
	return
}

