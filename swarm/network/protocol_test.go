
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450114131333120>


package network

import (
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/protocols"
	p2ptest "github.com/ethereum/go-ethereum/p2p/testing"
)

const (
	TestProtocolVersion   = 8
	TestProtocolNetworkID = 3
)

var (
	loglevel = flag.Int("loglevel", 2, "verbosity of logs")
)

func init() {
	flag.Parse()
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(*loglevel), log.StreamHandler(os.Stderr, log.TerminalFormat(true))))
}

func HandshakeMsgExchange(lhs, rhs *HandshakeMsg, id enode.ID) []p2ptest.Exchange {
	return []p2ptest.Exchange{
		{
			Expects: []p2ptest.Expect{
				{
					Code: 0,
					Msg:  lhs,
					Peer: id,
				},
			},
		},
		{
			Triggers: []p2ptest.Trigger{
				{
					Code: 0,
					Msg:  rhs,
					Peer: id,
				},
			},
		},
	}
}

func newBzzBaseTester(t *testing.T, n int, addr *BzzAddr, spec *protocols.Spec, run func(*BzzPeer) error) *bzzTester {
	cs := make(map[string]chan bool)

	srv := func(p *BzzPeer) error {
		defer func() {
			if cs[p.ID().String()] != nil {
				close(cs[p.ID().String()])
			}
		}()
		return run(p)
	}

	protocol := func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
		return srv(&BzzPeer{Peer: protocols.NewPeer(p, rw, spec), BzzAddr: NewAddr(p.Node())})
	}

	s := p2ptest.NewProtocolTester(t, addr.ID(), n, protocol)

	for _, node := range s.Nodes {
		cs[node.ID().String()] = make(chan bool)
	}

	return &bzzTester{
		addr:           addr,
		ProtocolTester: s,
		cs:             cs,
	}
}

type bzzTester struct {
	*p2ptest.ProtocolTester
	addr *BzzAddr
	cs   map[string]chan bool
	bzz  *Bzz
}

func newBzz(addr *BzzAddr, lightNode bool) *Bzz {
	config := &BzzConfig{
		OverlayAddr:  addr.Over(),
		UnderlayAddr: addr.Under(),
		HiveParams:   NewHiveParams(),
		NetworkID:    DefaultNetworkID,
		LightNode:    lightNode,
	}
	kad := NewKademlia(addr.OAddr, NewKadParams())
	bzz := NewBzz(config, kad, nil, nil, nil)
	return bzz
}

func newBzzHandshakeTester(t *testing.T, n int, addr *BzzAddr, lightNode bool) *bzzTester {
	bzz := newBzz(addr, lightNode)
	pt := p2ptest.NewProtocolTester(t, addr.ID(), n, bzz.runBzz)

	return &bzzTester{
		addr:           addr,
		ProtocolTester: pt,
		bzz:            bzz,
	}
}

//应该在一次交换中测试握手吗？并行化
func (s *bzzTester) testHandshake(lhs, rhs *HandshakeMsg, disconnects ...*p2ptest.Disconnect) error {
	if err := s.TestExchanges(HandshakeMsgExchange(lhs, rhs, rhs.Addr.ID())...); err != nil {
		return err
	}

	if len(disconnects) > 0 {
		return s.TestDisconnected(disconnects...)
	}

//如果我们不期望断开连接，请确保对等端保持连接
	err := s.TestDisconnected(&p2ptest.Disconnect{
		Peer:  s.Nodes[0].ID(),
		Error: nil,
	})

	if err == nil {
		return fmt.Errorf("Unexpected peer disconnect")
	}

	if err.Error() != "timed out waiting for peers to disconnect" {
		return err
	}

	return nil
}

func correctBzzHandshake(addr *BzzAddr, lightNode bool) *HandshakeMsg {
	return &HandshakeMsg{
		Version:   TestProtocolVersion,
		NetworkID: TestProtocolNetworkID,
		Addr:      addr,
		LightNode: lightNode,
	}
}

func TestBzzHandshakeNetworkIDMismatch(t *testing.T) {
	lightNode := false
	addr := RandomAddr()
	s := newBzzHandshakeTester(t, 1, addr, lightNode)
	node := s.Nodes[0]

	err := s.testHandshake(
		correctBzzHandshake(addr, lightNode),
		&HandshakeMsg{Version: TestProtocolVersion, NetworkID: 321, Addr: NewAddr(node)},
		&p2ptest.Disconnect{Peer: node.ID(), Error: fmt.Errorf("Handshake error: Message handler error: (msg code 0): network id mismatch 321 (!= 3)")},
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestBzzHandshakeVersionMismatch(t *testing.T) {
	lightNode := false
	addr := RandomAddr()
	s := newBzzHandshakeTester(t, 1, addr, lightNode)
	node := s.Nodes[0]

	err := s.testHandshake(
		correctBzzHandshake(addr, lightNode),
		&HandshakeMsg{Version: 0, NetworkID: TestProtocolNetworkID, Addr: NewAddr(node)},
		&p2ptest.Disconnect{Peer: node.ID(), Error: fmt.Errorf("Handshake error: Message handler error: (msg code 0): version mismatch 0 (!= %d)", TestProtocolVersion)},
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestBzzHandshakeSuccess(t *testing.T) {
	lightNode := false
	addr := RandomAddr()
	s := newBzzHandshakeTester(t, 1, addr, lightNode)
	node := s.Nodes[0]

	err := s.testHandshake(
		correctBzzHandshake(addr, lightNode),
		&HandshakeMsg{Version: TestProtocolVersion, NetworkID: TestProtocolNetworkID, Addr: NewAddr(node)},
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestBzzHandshakeLightNode(t *testing.T) {
	var lightNodeTests = []struct {
		name      string
		lightNode bool
	}{
		{"on", true},
		{"off", false},
	}

	for _, test := range lightNodeTests {
		t.Run(test.name, func(t *testing.T) {
			randomAddr := RandomAddr()
pt := newBzzHandshakeTester(nil, 1, randomAddr, false) //TODO更改签名-t未在任何位置使用
			node := pt.Nodes[0]
			addr := NewAddr(node)

			err := pt.testHandshake(
				correctBzzHandshake(randomAddr, false),
				&HandshakeMsg{Version: TestProtocolVersion, NetworkID: TestProtocolNetworkID, Addr: addr, LightNode: test.lightNode},
			)

			if err != nil {
				t.Fatal(err)
			}

			select {

			case <-pt.bzz.handshakes[node.ID()].done:
				if pt.bzz.handshakes[node.ID()].LightNode != test.lightNode {
					t.Fatalf("peer LightNode flag is %v, should be %v", pt.bzz.handshakes[node.ID()].LightNode, test.lightNode)
				}
			case <-time.After(10 * time.Second):
				t.Fatal("test timeout")
			}
		})
	}
}

