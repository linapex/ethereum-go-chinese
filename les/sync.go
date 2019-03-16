
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450096259403776>


package les

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/light"
)

//同步器负责定期与网络同步，两者都是
//下载哈希和块以及处理公告处理程序。
func (pm *ProtocolManager) syncer() {
//启动并确保清除同步机制
//pm.fetcher.start（）。
//延迟pm.fetcher.stop（）
	defer pm.downloader.Terminate()

//等待不同事件触发同步操作
//forceSync：=time.tick（forceSyncCycle）
	for {
		select {
		case <-pm.newPeerCh:
   /*//确保要从中选择对等点，然后同步
      如果pm.peers.len（）<mindesiredpeerCount_
       打破
      }
      转到pm.synchronize（pm.peers.bestpeer（））
   **/

  /*ASE<-强制同步：
  //即使没有足够的对等点，也强制同步
  转到pm.synchronize（pm.peers.bestpeer（））
  **/

		case <-pm.noMorePeers:
			return
		}
	}
}

func (pm *ProtocolManager) needToSync(peerHead blockInfo) bool {
	head := pm.blockchain.CurrentHeader()
	currentTd := rawdb.ReadTd(pm.chainDb, head.Hash(), head.Number.Uint64())
	return currentTd != nil && peerHead.Td.Cmp(currentTd) > 0
}

//同步尝试同步我们的本地块链与远程对等。
func (pm *ProtocolManager) synchronise(peer *peer) {
//如果没有对等点，则短路
	if peer == nil {
		return
	}

//确保同行的TD高于我们自己的TD。
	if !pm.needToSync(peer.headBlockInfo()) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pm.blockchain.(*light.LightChain).SyncCht(ctx)
	pm.downloader.Synchronise(peer.id, peer.Head(), peer.Td(), downloader.LightSync)
}

