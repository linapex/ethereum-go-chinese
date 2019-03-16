
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450090722922496>


//包事件处理对实时事件的订阅。
package event

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

//typemuxevent是一个推送到订户的带有时间标签的通知。
type TypeMuxEvent struct {
	Time time.Time
	Data interface{}
}

//typemux将事件发送给已注册的接收者。接收器可以
//注册以处理特定类型的事件。任何操作
//在mux停止后调用将返回errmuxClosed。
//
//The zero value is ready to use.
//
//已弃用：使用源
type TypeMux struct {
	mutex   sync.RWMutex
	subm    map[reflect.Type][]*TypeMuxSubscription
	stopped bool
}

//在已关闭的typemux上发布时返回errmuxClosed。
var ErrMuxClosed = errors.New("event: mux closed")

//订阅为给定类型的事件创建订阅。这个
//订阅的频道在取消订阅时关闭
//或者MUX关闭。
func (mux *TypeMux) Subscribe(types ...interface{}) *TypeMuxSubscription {
	sub := newsub(mux)
	mux.mutex.Lock()
	defer mux.mutex.Unlock()
	if mux.stopped {
//将状态设置为“已关闭”，以便在此之后调用Unsubscribe
//呼叫将短路。
		sub.closed = true
		close(sub.postC)
	} else {
		if mux.subm == nil {
			mux.subm = make(map[reflect.Type][]*TypeMuxSubscription)
		}
		for _, t := range types {
			rtyp := reflect.TypeOf(t)
			oldsubs := mux.subm[rtyp]
			if find(oldsubs, sub) != -1 {
				panic(fmt.Sprintf("event: duplicate type %s in Subscribe", rtyp))
			}
			subs := make([]*TypeMuxSubscription, len(oldsubs)+1)
			copy(subs, oldsubs)
			subs[len(oldsubs)] = sub
			mux.subm[rtyp] = subs
		}
	}
	return sub
}

//Post向为给定类型注册的所有接收器发送事件。
//如果MUX已停止，则返回errmuxClosed。
func (mux *TypeMux) Post(ev interface{}) error {
	event := &TypeMuxEvent{
		Time: time.Now(),
		Data: ev,
	}
	rtyp := reflect.TypeOf(ev)
	mux.mutex.RLock()
	if mux.stopped {
		mux.mutex.RUnlock()
		return ErrMuxClosed
	}
	subs := mux.subm[rtyp]
	mux.mutex.RUnlock()
	for _, sub := range subs {
		sub.deliver(event)
	}
	return nil
}

//停止关闭一个多路复用器。MUX不能再使用了。
//以后的Post调用将失败，并关闭errmuxClose。
//停止块，直到所有当前交货完成。
func (mux *TypeMux) Stop() {
	mux.mutex.Lock()
	for _, subs := range mux.subm {
		for _, sub := range subs {
			sub.closewait()
		}
	}
	mux.subm = nil
	mux.stopped = true
	mux.mutex.Unlock()
}

func (mux *TypeMux) del(s *TypeMuxSubscription) {
	mux.mutex.Lock()
	for typ, subs := range mux.subm {
		if pos := find(subs, s); pos >= 0 {
			if len(subs) == 1 {
				delete(mux.subm, typ)
			} else {
				mux.subm[typ] = posdelete(subs, pos)
			}
		}
	}
	s.mux.mutex.Unlock()
}

func find(slice []*TypeMuxSubscription, item *TypeMuxSubscription) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
}

func posdelete(slice []*TypeMuxSubscription, pos int) []*TypeMuxSubscription {
	news := make([]*TypeMuxSubscription, len(slice)-1)
	copy(news[:pos], slice[:pos])
	copy(news[pos:], slice[pos+1:])
	return news
}

//typemux订阅是通过typemux建立的订阅。
type TypeMuxSubscription struct {
	mux     *TypeMux
	created time.Time
	closeMu sync.Mutex
	closing chan struct{}
	closed  bool

//这两个频道是同一频道。它们分开存放，所以
//post可以设置为nil，而不影响
//陈。
	postMu sync.RWMutex
	readC  <-chan *TypeMuxEvent
	postC  chan<- *TypeMuxEvent
}

func newsub(mux *TypeMux) *TypeMuxSubscription {
	c := make(chan *TypeMuxEvent)
	return &TypeMuxSubscription{
		mux:     mux,
		created: time.Now(),
		readC:   c,
		postC:   c,
		closing: make(chan struct{}),
	}
}

func (s *TypeMuxSubscription) Chan() <-chan *TypeMuxEvent {
	return s.readC
}

func (s *TypeMuxSubscription) Unsubscribe() {
	s.mux.del(s)
	s.closewait()
}

func (s *TypeMuxSubscription) Closed() bool {
	s.closeMu.Lock()
	defer s.closeMu.Unlock()
	return s.closed
}

func (s *TypeMuxSubscription) closewait() {
	s.closeMu.Lock()
	defer s.closeMu.Unlock()
	if s.closed {
		return
	}
	close(s.closing)
	s.closed = true

	s.postMu.Lock()
	close(s.postC)
	s.postC = nil
	s.postMu.Unlock()
}

func (s *TypeMuxSubscription) deliver(event *TypeMuxEvent) {
//失效事件时短路交付
	if s.created.After(event.Time) {
		return
	}
//否则，交付活动
	s.postMu.RLock()
	defer s.postMu.RUnlock()

	select {
	case s.postC <- event:
	case <-s.closing:
	}
}

