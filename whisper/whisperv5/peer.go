
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:46</date>
//</624450124763893760>


package whisperv5

import (
	"fmt"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/rlp"
)

//Peer表示一个耳语协议对等连接。
type Peer struct {
	host    *Whisper
	peer    *p2p.Peer
	ws      p2p.MsgReadWriter
	trusted bool

known mapset.Set //对等方已经知道的消息，以避免浪费带宽

	quit chan struct{}
}

//new peer创建一个新的whisper peer对象，但不运行握手本身。
func newPeer(host *Whisper, remote *p2p.Peer, rw p2p.MsgReadWriter) *Peer {
	return &Peer{
		host:    host,
		peer:    remote,
		ws:      rw,
		trusted: false,
		known:   mapset.NewSet(),
		quit:    make(chan struct{}),
	}
}

//启动启动对等更新程序，定期广播低语数据包
//进入网络。
func (peer *Peer) start() {
	go peer.update()
	log.Trace("start", "peer", peer.ID())
}

//stop终止对等更新程序，停止向其转发消息。
func (peer *Peer) stop() {
	close(peer.quit)
	log.Trace("stop", "peer", peer.ID())
}

//握手向远程对等端发送协议启动状态消息，并且
//也验证远程状态。
func (peer *Peer) handshake() error {
//异步发送握手状态消息
	errc := make(chan error, 1)
	go func() {
		errc <- p2p.Send(peer.ws, statusCode, ProtocolVersion)
	}()
//获取远程状态包并验证协议匹配
	packet, err := peer.ws.ReadMsg()
	if err != nil {
		return err
	}
	if packet.Code != statusCode {
		return fmt.Errorf("peer [%x] sent packet %x before status packet", peer.ID(), packet.Code)
	}
	s := rlp.NewStream(packet.Payload, uint64(packet.Size))
	peerVersion, err := s.Uint()
	if err != nil {
		return fmt.Errorf("peer [%x] sent bad status message: %v", peer.ID(), err)
	}
	if peerVersion != ProtocolVersion {
		return fmt.Errorf("peer [%x]: protocol version mismatch %d != %d", peer.ID(), peerVersion, ProtocolVersion)
	}
//等待直到消耗掉自己的状态
	if err := <-errc; err != nil {
		return fmt.Errorf("peer [%x] failed to send status packet: %v", peer.ID(), err)
	}
	return nil
}

//更新在对等机上执行定期操作，包括消息传输
//和呼气。
func (peer *Peer) update() {
//启动更新的滚动条
	expire := time.NewTicker(expirationCycle)
	transmit := time.NewTicker(transmissionCycle)

//循环并发送直到请求终止
	for {
		select {
		case <-expire.C:
			peer.expire()

		case <-transmit.C:
			if err := peer.broadcast(); err != nil {
				log.Trace("broadcast failed", "reason", err, "peer", peer.ID())
				return
			}

		case <-peer.quit:
			return
		}
	}
}

//马克标记了一个同伴知道的信封，这样它就不会被送回。
func (peer *Peer) mark(envelope *Envelope) {
	peer.known.Add(envelope.Hash())
}

//标记检查远程对等机是否已经知道信封。
func (peer *Peer) marked(envelope *Envelope) bool {
	return peer.known.Contains(envelope.Hash())
}

//Expire迭代主机中的所有已知信封，并删除所有
//已知列表中过期（未知）的。
func (peer *Peer) expire() {
	unmark := make(map[common.Hash]struct{})
	peer.known.Each(func(v interface{}) bool {
		if !peer.host.isEnvelopeCached(v.(common.Hash)) {
			unmark[v.(common.Hash)] = struct{}{}
		}
		return true
	})
//转储所有已知但不再缓存的内容
	for hash := range unmark {
		peer.known.Remove(hash)
	}
}

//广播在信封集合上迭代，传输未知信息
//在网络上。
func (peer *Peer) broadcast() error {
	var cnt int
	envelopes := peer.host.Envelopes()
	for _, envelope := range envelopes {
		if !peer.marked(envelope) {
			err := p2p.Send(peer.ws, messagesCode, envelope)
			if err != nil {
				return err
			} else {
				peer.mark(envelope)
				cnt++
			}
		}
	}
	if cnt > 0 {
		log.Trace("broadcast", "num. messages", cnt)
	}
	return nil
}

func (peer *Peer) ID() []byte {
	id := peer.peer.ID()
	return id[:]
}

