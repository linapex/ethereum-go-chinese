
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450091163324416>


package event

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/mclock"
)

//订阅表示事件流。事件的载体通常是
//但不是接口的一部分。
//
//建立订阅时可能失败。通过错误报告失败
//通道。如果订阅出现问题（例如
//network connection delivering the events has been closed). Only one value will ever be
//发送。
//
//当订阅成功结束时（即当
//source of events is closed). It is also closed when Unsubscribe is called.
//
//Unsubscribe方法取消发送事件。您必须呼叫取消订阅
//案例以确保与订阅相关的资源被释放。它可以
//调用任意次数。
type Subscription interface {
Err() <-chan error //返回错误通道
Unsubscribe()      //取消发送事件，关闭错误通道
}

//new subscription在新的goroutine中作为订阅运行producer函数。这个
//当调用UNSUBSCRIBE时，将关闭提供给生产者的频道。如果fn返回
//错误，它在订阅的错误通道上发送。
func NewSubscription(producer func(<-chan struct{}) error) Subscription {
	s := &funcSub{unsub: make(chan struct{}), err: make(chan error, 1)}
	go func() {
		defer close(s.err)
		err := producer(s.unsub)
		s.mu.Lock()
		defer s.mu.Unlock()
		if !s.unsubscribed {
			if err != nil {
				s.err <- err
			}
			s.unsubscribed = true
		}
	}()
	return s
}

type funcSub struct {
	unsub        chan struct{}
	err          chan error
	mu           sync.Mutex
	unsubscribed bool
}

func (s *funcSub) Unsubscribe() {
	s.mu.Lock()
	if s.unsubscribed {
		s.mu.Unlock()
		return
	}
	s.unsubscribed = true
	close(s.unsub)
	s.mu.Unlock()
//等待生产商关闭。
	<-s.err
}

func (s *funcSub) Err() <-chan error {
	return s.err
}

//重复重新订阅呼叫fn以保持已建立的订阅。当
//订阅已建立，重新订阅等待失败，然后再次调用fn。这个
//过程重复，直到调用取消订阅或活动订阅结束
//成功地。
//
//resubscribe在对fn的调用之间应用回退。调整通话间隔时间
//基于错误率，但不会超过backoffmax。
func Resubscribe(backoffMax time.Duration, fn ResubscribeFunc) Subscription {
	s := &resubscribeSub{
		waitTime:   backoffMax / 10,
		backoffMax: backoffMax,
		fn:         fn,
		err:        make(chan error),
		unsub:      make(chan struct{}),
	}
	go s.loop()
	return s
}

//A ResubscribeFunc attempts to establish a subscription.
type ResubscribeFunc func(context.Context) (Subscription, error)

type resubscribeSub struct {
	fn                   ResubscribeFunc
	err                  chan error
	unsub                chan struct{}
	unsubOnce            sync.Once
	lastTry              mclock.AbsTime
	waitTime, backoffMax time.Duration
}

func (s *resubscribeSub) Unsubscribe() {
	s.unsubOnce.Do(func() {
		s.unsub <- struct{}{}
		<-s.err
	})
}

func (s *resubscribeSub) Err() <-chan error {
	return s.err
}

func (s *resubscribeSub) loop() {
	defer close(s.err)
	var done bool
	for !done {
		sub := s.subscribe()
		if sub == nil {
			break
		}
		done = s.waitForError(sub)
		sub.Unsubscribe()
	}
}

func (s *resubscribeSub) subscribe() Subscription {
	subscribed := make(chan error)
	var sub Subscription
retry:
	for {
		s.lastTry = mclock.Now()
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			rsub, err := s.fn(ctx)
			sub = rsub
			subscribed <- err
		}()
		select {
		case err := <-subscribed:
			cancel()
			if err != nil {
//订阅失败，请等待，然后启动下一次尝试。
				if s.backoffWait() {
					return nil
				}
				continue retry
			}
			if sub == nil {
				panic("event: ResubscribeFunc returned nil subscription and no error")
			}
			return sub
		case <-s.unsub:
			cancel()
			return nil
		}
	}
}

func (s *resubscribeSub) waitForError(sub Subscription) bool {
	defer sub.Unsubscribe()
	select {
	case err := <-sub.Err():
		return err == nil
	case <-s.unsub:
		return true
	}
}

func (s *resubscribeSub) backoffWait() bool {
	if time.Duration(mclock.Now()-s.lastTry) > s.backoffMax {
		s.waitTime = s.backoffMax / 10
	} else {
		s.waitTime *= 2
		if s.waitTime > s.backoffMax {
			s.waitTime = s.backoffMax
		}
	}

	t := time.NewTimer(s.waitTime)
	defer t.Stop()
	select {
	case <-t.C:
		return false
	case <-s.unsub:
		return true
	}
}

//subscriptionScope提供了一种功能，可以一次取消订阅多个订阅。
//
//对于处理多个订阅的代码，可以方便地使用范围
//只需一个电话就可以取消所有订阅。该示例演示了
//更大的程序。
//
//零值已准备好使用。
type SubscriptionScope struct {
	mu     sync.Mutex
	subs   map[*scopeSub]struct{}
	closed bool
}

type scopeSub struct {
	sc *SubscriptionScope
	s  Subscription
}

//跟踪开始跟踪订阅。如果范围已关闭，则track返回nil。这个
//returned subscription is a wrapper. Unsubscribing the wrapper removes it from the
//范围。
func (sc *SubscriptionScope) Track(s Subscription) Subscription {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.closed {
		return nil
	}
	if sc.subs == nil {
		sc.subs = make(map[*scopeSub]struct{})
	}
	ss := &scopeSub{sc, s}
	sc.subs[ss] = struct{}{}
	return ss
}

//关闭对所有跟踪订阅的取消订阅的呼叫，并阻止进一步添加到
//the tracked set. Calls to Track after Close return nil.
func (sc *SubscriptionScope) Close() {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	if sc.closed {
		return
	}
	sc.closed = true
	for s := range sc.subs {
		s.s.Unsubscribe()
	}
	sc.subs = nil
}

//count返回跟踪订阅的数目。
//它是用来调试的。
func (sc *SubscriptionScope) Count() int {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return len(sc.subs)
}

func (s *scopeSub) Unsubscribe() {
	s.s.Unsubscribe()
	s.sc.mu.Lock()
	defer s.sc.mu.Unlock()
	delete(s.sc.subs, s)
}

func (s *scopeSub) Err() <-chan error {
	return s.s.Err()
}

