
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450090852945920>


package event_test

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/event"
)

//此示例演示如何使用subscriptionScope来控制
//订阅。
//
//我们的示例程序由两个服务器组成，每个服务器在
//请求。服务器还允许订阅所有计算的结果。
type divServer struct{ results event.Feed }
type mulServer struct{ results event.Feed }

func (s *divServer) do(a, b int) int {
	r := a / b
	s.results.Send(r)
	return r
}

func (s *mulServer) do(a, b int) int {
	r := a * b
	s.results.Send(r)
	return r
}

//The servers are contained in an App. The app controls the servers and exposes them
//通过它的API。
type App struct {
	divServer
	mulServer
	scope event.SubscriptionScope
}

func (s *App) Calc(op byte, a, b int) int {
	switch op {
	case '/':
		return s.divServer.do(a, b)
	case '*':
		return s.mulServer.do(a, b)
	default:
		panic("invalid op")
	}
}

//应用程序的subscripresults方法开始将计算结果发送给给定的
//通道。通过此方法创建的订阅与应用程序的生存期绑定。
//因为它们是在作用域中注册的。
func (s *App) SubscribeResults(op byte, ch chan<- int) event.Subscription {
	switch op {
	case '/':
		return s.scope.Track(s.divServer.results.Subscribe(ch))
	case '*':
		return s.scope.Track(s.mulServer.results.Subscribe(ch))
	default:
		panic("invalid op")
	}
}

//stop停止应用程序，关闭通过subscripresults创建的所有订阅。
func (s *App) Stop() {
	s.scope.Close()
}

func ExampleSubscriptionScope() {
//创建应用程序。
	var (
		app  App
		wg   sync.WaitGroup
		divs = make(chan int)
		muls = make(chan int)
	)

//在后台运行订阅服务器。
	divsub := app.SubscribeResults('/', divs)
	mulsub := app.SubscribeResults('*', muls)
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer fmt.Println("subscriber exited")
		defer divsub.Unsubscribe()
		defer mulsub.Unsubscribe()
		for {
			select {
			case result := <-divs:
				fmt.Println("division happened:", result)
			case result := <-muls:
				fmt.Println("multiplication happened:", result)
			case <-divsub.Err():
				return
			case <-mulsub.Err():
				return
			}
		}
	}()

//与应用程序交互。
	app.Calc('/', 22, 11)
	app.Calc('*', 3, 4)

//停止应用程序。这将关闭订阅，导致订阅服务器退出。
	app.Stop()
	wg.Wait()

//输出：
//分部发生：2
//发生乘法：12
//已退出订阅服务器
}

