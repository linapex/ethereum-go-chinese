
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450105306517504>


package netutil

import (
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
)

//IPtracker预测本地主机的外部端点，即IP地址和端口
//基于其他主机的语句。
type IPTracker struct {
	window          time.Duration
	contactWindow   time.Duration
	minStatements   int
	clock           mclock.Clock
	statements      map[string]ipStatement
	contact         map[string]mclock.AbsTime
	lastStatementGC mclock.AbsTime
	lastContactGC   mclock.AbsTime
}

type ipStatement struct {
	endpoint string
	time     mclock.AbsTime
}

//newiptracker创建一个IP跟踪器。
//
//窗口参数配置保留的过去网络事件的数量。这个
//minStatements参数强制执行必须记录的最小语句数
//在做出任何预测之前。这些参数的值越高，则“拍打”越小。
//网络条件变化时的预测。窗口持续时间值通常应在
//分钟的范围。
func NewIPTracker(window, contactWindow time.Duration, minStatements int) *IPTracker {
	return &IPTracker{
		window:        window,
		contactWindow: contactWindow,
		statements:    make(map[string]ipStatement),
		minStatements: minStatements,
		contact:       make(map[string]mclock.AbsTime),
		clock:         mclock.System{},
	}
}

//PredictTfullconenat检查本地主机是否位于全锥NAT之后。它预测
//正在检查是否已从以前未联系的节点接收到任何语句
//声明已作出。
func (it *IPTracker) PredictFullConeNAT() bool {
	now := it.clock.Now()
	it.gcContact(now)
	it.gcStatements(now)
	for host, st := range it.statements {
		if c, ok := it.contact[host]; !ok || c > st.time {
			return true
		}
	}
	return false
}

//PredictEndPoint返回外部终结点的当前预测。
func (it *IPTracker) PredictEndpoint() string {
	it.gcStatements(it.clock.Now())

//当前的策略很简单：使用大多数语句查找端点。
	counts := make(map[string]int)
	maxcount, max := 0, ""
	for _, s := range it.statements {
		c := counts[s.endpoint] + 1
		counts[s.endpoint] = c
		if c > maxcount && c >= it.minStatements {
			maxcount, max = c, s.endpoint
		}
	}
	return max
}

//addStatement记录某个主机认为我们的外部端点是给定的端点。
func (it *IPTracker) AddStatement(host, endpoint string) {
	now := it.clock.Now()
	it.statements[host] = ipStatement{endpoint, now}
	if time.Duration(now-it.lastStatementGC) >= it.window {
		it.gcStatements(now)
	}
}

//addcontact记录包含端点信息的数据包已发送到
//某些宿主。
func (it *IPTracker) AddContact(host string) {
	now := it.clock.Now()
	it.contact[host] = now
	if time.Duration(now-it.lastContactGC) >= it.contactWindow {
		it.gcContact(now)
	}
}

func (it *IPTracker) gcStatements(now mclock.AbsTime) {
	it.lastStatementGC = now
	cutoff := now.Add(-it.window)
	for host, s := range it.statements {
		if s.time < cutoff {
			delete(it.statements, host)
		}
	}
}

func (it *IPTracker) gcContact(now mclock.AbsTime) {
	it.lastContactGC = now
	cutoff := now.Add(-it.contactWindow)
	for host, ct := range it.contact {
		if ct < cutoff {
			delete(it.contact, host)
		}
	}
}

