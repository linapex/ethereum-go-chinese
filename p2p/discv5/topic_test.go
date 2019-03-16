
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450104085975040>


package discv5

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/mclock"
)

func TestTopicRadius(t *testing.T) {
	now := mclock.Now()
	topic := Topic("qwerty")
	rad := newTopicRadius(topic)
	targetRad := (^uint64(0)) / 100

	waitFn := func(addr common.Hash) time.Duration {
		prefix := binary.BigEndian.Uint64(addr[0:8])
		dist := prefix ^ rad.topicHashPrefix
		relDist := float64(dist) / float64(targetRad)
		relTime := (1 - relDist/2) * 2
		if relTime < 0 {
			relTime = 0
		}
		return time.Duration(float64(targetWaitTime) * relTime)
	}

	bcnt := 0
	cnt := 0
	var sum float64
	for cnt < 100 {
		addr := rad.nextTarget(false).target
		wait := waitFn(addr)
		ticket := &ticket{
			topics:  []Topic{topic},
			regTime: []mclock.AbsTime{mclock.AbsTime(wait)},
			node:    &Node{nodeNetGuts: nodeNetGuts{sha: addr}},
		}
		rad.adjustWithTicket(now, addr, ticketRef{ticket, 0})
		if rad.radius != maxRadius {
			cnt++
			sum += float64(rad.radius)
		} else {
			bcnt++
			if bcnt > 500 {
				t.Errorf("Radius did not converge in 500 iterations")
			}
		}
	}
	avgRel := sum / float64(cnt) / float64(targetRad)
	if avgRel > 1.05 || avgRel < 0.95 {
		t.Errorf("Average/target ratio is too far from 1 (%v)", avgRel)
	}
}

