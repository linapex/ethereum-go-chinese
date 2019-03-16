
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450090936832000>


package event_test

import (
	"fmt"

	"github.com/ethereum/go-ethereum/event"
)

func ExampleNewSubscription() {
//创建在ch上发送10个整数的订阅。
	ch := make(chan int)
	sub := event.NewSubscription(func(quit <-chan struct{}) error {
		for i := 0; i < 10; i++ {
			select {
			case ch <- i:
			case <-quit:
				fmt.Println("unsubscribed")
				return nil
			}
		}
		return nil
	})

//这是消费者。它读取5个整数，然后中止订阅。
//请注意，取消订阅会一直等到生产者关闭。
	for i := range ch {
		fmt.Println(i)
		if i == 4 {
			sub.Unsubscribe()
			break
		}
	}
//输出：
//零
//一
//二
//三
//四
//退订
}

