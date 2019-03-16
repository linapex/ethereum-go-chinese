
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450105348460544>


package netutil

import (
	"fmt"
	mrand "math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
)

const (
	opStatement = iota
	opContact
	opPredict
	opCheckFullCone
)

type iptrackTestEvent struct {
	op       int
time     int //绝对值（毫秒）
	ip, from string
}

func TestIPTracker(t *testing.T) {
	tests := map[string][]iptrackTestEvent{
		"minStatements": {
			{opPredict, 0, "", ""},
			{opStatement, 0, "127.0.0.1", "127.0.0.2"},
			{opPredict, 1000, "", ""},
			{opStatement, 1000, "127.0.0.1", "127.0.0.3"},
			{opPredict, 1000, "", ""},
			{opStatement, 1000, "127.0.0.1", "127.0.0.4"},
			{opPredict, 1000, "127.0.0.1", ""},
		},
		"window": {
			{opStatement, 0, "127.0.0.1", "127.0.0.2"},
			{opStatement, 2000, "127.0.0.1", "127.0.0.3"},
			{opStatement, 3000, "127.0.0.1", "127.0.0.4"},
			{opPredict, 10000, "127.0.0.1", ""},
{opPredict, 10001, "", ""}, //第一条语句已过期
			{opStatement, 10100, "127.0.0.1", "127.0.0.2"},
			{opPredict, 10200, "127.0.0.1", ""},
		},
		"fullcone": {
			{opContact, 0, "", "127.0.0.2"},
			{opStatement, 10, "127.0.0.1", "127.0.0.2"},
			{opContact, 2000, "", "127.0.0.3"},
			{opStatement, 2010, "127.0.0.1", "127.0.0.3"},
			{opContact, 3000, "", "127.0.0.4"},
			{opStatement, 3010, "127.0.0.1", "127.0.0.4"},
			{opCheckFullCone, 3500, "false", ""},
		},
		"fullcone_2": {
			{opContact, 0, "", "127.0.0.2"},
			{opStatement, 10, "127.0.0.1", "127.0.0.2"},
			{opContact, 2000, "", "127.0.0.3"},
			{opStatement, 2010, "127.0.0.1", "127.0.0.3"},
			{opStatement, 3000, "127.0.0.1", "127.0.0.4"},
			{opContact, 3010, "", "127.0.0.4"},
			{opCheckFullCone, 3500, "true", ""},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) { runIPTrackerTest(t, test) })
	}
}

func runIPTrackerTest(t *testing.T, evs []iptrackTestEvent) {
	var (
		clock mclock.Simulated
		it    = NewIPTracker(10*time.Second, 10*time.Second, 3)
	)
	it.clock = &clock
	for i, ev := range evs {
		evtime := time.Duration(ev.time) * time.Millisecond
		clock.Run(evtime - time.Duration(clock.Now()))
		switch ev.op {
		case opStatement:
			it.AddStatement(ev.from, ev.ip)
		case opContact:
			it.AddContact(ev.from)
		case opPredict:
			if pred := it.PredictEndpoint(); pred != ev.ip {
				t.Errorf("op %d: wrong prediction %q, want %q", i, pred, ev.ip)
			}
		case opCheckFullCone:
			pred := fmt.Sprintf("%t", it.PredictFullConeNAT())
			if pred != ev.ip {
				t.Errorf("op %d: wrong prediction %s, want %s", i, pred, ev.ip)
			}
		}
	}
}

//这将检查旧的语句和联系人是否已GCED，即使没有调用Predict*。
func TestIPTrackerForceGC(t *testing.T) {
	var (
		clock  mclock.Simulated
		window = 10 * time.Second
		rate   = 50 * time.Millisecond
		max    = int(window/rate) + 1
		it     = NewIPTracker(window, window, 3)
	)
	it.clock = &clock

	for i := 0; i < 5*max; i++ {
		e1 := make([]byte, 4)
		e2 := make([]byte, 4)
		mrand.Read(e1)
		mrand.Read(e2)
		it.AddStatement(string(e1), string(e2))
		it.AddContact(string(e1))
		clock.Run(rate)
	}
	if len(it.contact) > 2*max {
		t.Errorf("contacts not GCed, have %d", len(it.contact))
	}
	if len(it.statements) > 2*max {
		t.Errorf("statements not GCed, have %d", len(it.statements))
	}
}
