
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450090815197184>


package event_test

import (
	"fmt"

	"github.com/ethereum/go-ethereum/event"
)

func ExampleFeed_acknowledgedEvents() {
//此示例显示如何将send的返回值用于请求/答复
//活动消费者和生产者之间的互动。
	var feed event.Feed
	type ackedEvent struct {
		i   int
		ack chan<- struct{}
	}

//消费者等待feed上的事件并确认处理。
	done := make(chan struct{})
	defer close(done)
	for i := 0; i < 3; i++ {
		ch := make(chan ackedEvent, 100)
		sub := feed.Subscribe(ch)
		go func() {
			defer sub.Unsubscribe()
			for {
				select {
				case ev := <-ch:
fmt.Println(ev.i) //“处理”事件
					ev.ack <- struct{}{}
				case <-done:
					return
				}
			}
		}()
	}

//生产者发送AckedEvent类型的值，增加i的值。
//它在发送下一个事件之前等待所有消费者确认。
	for i := 0; i < 3; i++ {
		acksignal := make(chan struct{})
		n := feed.Send(ackedEvent{i, acksignal})
		for ack := 0; ack < n; ack++ {
			<-acksignal
		}
	}
//输出：
//零
//零
//零
//一
//一
//一
//二
//二
//二
}

