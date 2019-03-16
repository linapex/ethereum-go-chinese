
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450103121285120>


package discover

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"sync"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
)

var nullNode *enode.Node

func init() {
	var r enr.Record
	r.Set(enr.IP{0, 0, 0, 0})
	nullNode = enode.SignNull(&r, enode.ID{})
}

func newTestTable(t transport) (*Table, *enode.DB) {
	db, _ := enode.OpenDB("")
	tab, _ := newTable(t, db, nil)
	return tab, db
}

//nodeAtDistance为其创建一个enode.logDist（base，n.id）=ld的节点。
func nodeAtDistance(base enode.ID, ld int, ip net.IP) *node {
	var r enr.Record
	r.Set(enr.IP(ip))
	return wrapNode(enode.SignNull(&r, idAtDistance(base, ld)))
}

//IDATInstance返回一个随机哈希，使enode.logdist（a，b）==n
func idAtDistance(a enode.ID, n int) (b enode.ID) {
	if n == 0 {
		return a
	}
//在n位置翻转钻头，其余部分用随机钻头填满。
	b = a
	pos := len(a) - n/8 - 1
	bit := byte(0x01) << (byte(n%8) - 1)
	if bit == 0 {
		pos++
		bit = 0x80
	}
b[pos] = a[pos]&^bit | ^a[pos]&bit //TODO:随机结束位
	for i := pos + 1; i < len(a); i++ {
		b[i] = byte(rand.Intn(255))
	}
	return b
}

func intIP(i int) net.IP {
	return net.IP{byte(i), 0, 2, byte(i)}
}

//FillBucket将节点插入给定的bucket，直到它满为止。
func fillBucket(tab *Table, n *node) (last *node) {
	ld := enode.LogDist(tab.self().ID(), n.ID())
	b := tab.bucket(n.ID())
	for len(b.entries) < bucketSize {
		b.entries = append(b.entries, nodeAtDistance(tab.self().ID(), ld, intIP(ld)))
	}
	return b.entries[bucketSize-1]
}

type pingRecorder struct {
	mu           sync.Mutex
	dead, pinged map[enode.ID]bool
	n            *enode.Node
}

func newPingRecorder() *pingRecorder {
	var r enr.Record
	r.Set(enr.IP{0, 0, 0, 0})
	n := enode.SignNull(&r, enode.ID{})

	return &pingRecorder{
		dead:   make(map[enode.ID]bool),
		pinged: make(map[enode.ID]bool),
		n:      n,
	}
}

func (t *pingRecorder) self() *enode.Node {
	return nullNode
}

func (t *pingRecorder) findnode(toid enode.ID, toaddr *net.UDPAddr, target encPubkey) ([]*node, error) {
	return nil, nil
}

func (t *pingRecorder) waitping(from enode.ID) error {
return nil //远程总是ping
}

func (t *pingRecorder) ping(toid enode.ID, toaddr *net.UDPAddr) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.pinged[toid] = true
	if t.dead[toid] {
		return errTimeout
	} else {
		return nil
	}
}

func (t *pingRecorder) close() {}

func hasDuplicates(slice []*node) bool {
	seen := make(map[enode.ID]bool)
	for i, e := range slice {
		if e == nil {
			panic(fmt.Sprintf("nil *Node at %d", i))
		}
		if seen[e.ID()] {
			return true
		}
		seen[e.ID()] = true
	}
	return false
}

func contains(ns []*node, id enode.ID) bool {
	for _, n := range ns {
		if n.ID() == id {
			return true
		}
	}
	return false
}

func sortedByDistanceTo(distbase enode.ID, slice []*node) bool {
	var last enode.ID
	for i, e := range slice {
		if i > 0 && enode.DistCmp(distbase, e.ID(), last) < 0 {
			return false
		}
		last = e.ID()
	}
	return true
}

func hexEncPubkey(h string) (ret encPubkey) {
	b, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	if len(b) != len(ret) {
		panic("invalid length")
	}
	copy(ret[:], b)
	return ret
}

func hexPubkey(h string) *ecdsa.PublicKey {
	k, err := decodePubkey(hexEncPubkey(h))
	if err != nil {
		panic(err)
	}
	return k
}

