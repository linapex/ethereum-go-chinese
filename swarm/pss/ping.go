
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450116849242112>


//+建设！NOPSS协议，！诺普斯平

package pss

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/protocols"
	"github.com/ethereum/go-ethereum/swarm/log"
)

//的通用ping协议实现
//PSS devp2p协议仿真
type PingMsg struct {
	Created time.Time
Pong    bool //设置消息是否为pong reply
}

type Ping struct {
Pong bool      //ping接收时切换pong回复
OutC chan bool //触发平
InC  chan bool //可选，返回呼叫代码
}

func (p *Ping) pingHandler(ctx context.Context, msg interface{}) error {
	var pingmsg *PingMsg
	var ok bool
	if pingmsg, ok = msg.(*PingMsg); !ok {
		return errors.New("invalid msg")
	}
	log.Debug("ping handler", "msg", pingmsg, "outc", p.OutC)
	if p.InC != nil {
		p.InC <- pingmsg.Pong
	}
	if p.Pong && !pingmsg.Pong {
		p.OutC <- true
	}
	return nil
}

var PingProtocol = &protocols.Spec{
	Name:       "psstest",
	Version:    1,
	MaxMsgSize: 1024,
	Messages: []interface{}{
		PingMsg{},
	},
}

var PingTopic = ProtocolTopic(PingProtocol)

func NewPingProtocol(ping *Ping) *p2p.Protocol {
	return &p2p.Protocol{
		Name:    PingProtocol.Name,
		Version: PingProtocol.Version,
		Length:  uint64(PingProtocol.MaxMsgSize),
		Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
			quitC := make(chan struct{})
			pp := protocols.NewPeer(p, rw, PingProtocol)
			log.Trace("running pss vprotocol", "peer", p, "outc", ping.OutC)
			go func() {
				for {
					select {
					case ispong := <-ping.OutC:
						pp.Send(context.TODO(), &PingMsg{
							Created: time.Now(),
							Pong:    ispong,
						})
					case <-quitC:
					}
				}
			}()
			err := pp.Run(ping.pingHandler)
			quitC <- struct{}{}
			return err
		},
	}
}

