
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450106829049856>


package simulations

import (
	"fmt"
	"time"
)

//EventType是模拟网络发出的事件类型
type EventType string

const (
//EventTypeNode是当节点为
//创建、启动或停止
	EventTypeNode EventType = "node"

//EventTypeConn是连接时发出的事件类型
//在两个节点之间建立或删除
	EventTypeConn EventType = "conn"

//eventtypmsg是p2p消息时发出的事件类型。
//在两个节点之间发送
	EventTypeMsg EventType = "msg"
)

//事件是模拟网络发出的事件
type Event struct {
//类型是事件的类型
	Type EventType `json:"type"`

//时间是事件发生的时间
	Time time.Time `json:"time"`

//控件指示事件是否是受控件的结果
//网络中的操作
	Control bool `json:"control"`

//如果类型为EventTypeNode，则设置节点
	Node *Node `json:"node,omitempty"`

//如果类型为eventtypconn，则设置conn
	Conn *Conn `json:"conn,omitempty"`

//如果类型为eventtypmsg，则设置msg。
	Msg *Msg `json:"msg,omitempty"`

//可选提供数据（当前仅用于模拟前端）
	Data interface{} `json:"data"`
}

//NewEvent为给定对象创建一个新事件，该事件应为
//节点、连接或消息。
//
//复制对象以便事件表示对象的状态
//调用NewEvent时。
func NewEvent(v interface{}) *Event {
	event := &Event{Time: time.Now()}
	switch v := v.(type) {
	case *Node:
		event.Type = EventTypeNode
		node := *v
		event.Node = &node
	case *Conn:
		event.Type = EventTypeConn
		conn := *v
		event.Conn = &conn
	case *Msg:
		event.Type = EventTypeMsg
		msg := *v
		event.Msg = &msg
	default:
		panic(fmt.Sprintf("invalid event type: %T", v))
	}
	return event
}

//ControlEvent创建新的控件事件
func ControlEvent(v interface{}) *Event {
	event := NewEvent(v)
	event.Control = true
	return event
}

//字符串返回事件的字符串表示形式
func (e *Event) String() string {
	switch e.Type {
	case EventTypeNode:
		return fmt.Sprintf("<node-event> id: %s up: %t", e.Node.ID().TerminalString(), e.Node.Up)
	case EventTypeConn:
		return fmt.Sprintf("<conn-event> nodes: %s->%s up: %t", e.Conn.One.TerminalString(), e.Conn.Other.TerminalString(), e.Conn.Up)
	case EventTypeMsg:
		return fmt.Sprintf("<msg-event> nodes: %s->%s proto: %s, code: %d, received: %t", e.Msg.One.TerminalString(), e.Msg.Other.TerminalString(), e.Msg.Protocol, e.Msg.Code, e.Msg.Received)
	default:
		return ""
	}
}

