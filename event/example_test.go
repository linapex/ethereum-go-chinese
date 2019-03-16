
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450090999746560>


package event

import "fmt"

func ExampleTypeMux() {
	type someEvent struct{ I int }
	type otherEvent struct{ S string }
	type yetAnotherEvent struct{ X, Y int }

	var mux TypeMux

//Start a subscriber.
	done := make(chan struct{})
	sub := mux.Subscribe(someEvent{}, otherEvent{})
	go func() {
		for event := range sub.Chan() {
			fmt.Printf("Received: %#v\n", event.Data)
		}
		fmt.Println("done")
		close(done)
	}()

//发布一些事件。
	mux.Post(someEvent{5})
	mux.Post(yetAnotherEvent{X: 3, Y: 4})
	mux.Post(someEvent{6})
	mux.Post(otherEvent{"whoa"})

//stop关闭所有订阅通道。
//订户GOUDOTIN将打印“完成”
//然后退出。
	mux.Stop()

//等待订阅服务器返回。
	<-done

//输出：
//接收：事件。某些事件i:5
//接收：事件。某些事件i:6
//事件：其他事件{s：“哇”}
//完成
}

