
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450114009698304>


//包优先级队列实现基于通道的优先级队列
//在任意类型上。它提供了一个
//一个自动操作循环，将一个函数应用于始终遵守的项
//他们的优先权。结构只是准一致的，即如果
//优先项是自动停止的，保证有一点
//当没有更高优先级的项目时，即不能保证
//有一点低优先级的项目存在
//但更高的不是

package priorityqueue

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/log"
)

var (
	ErrContention = errors.New("contention")

	errBadPriority = errors.New("bad priority")

	wakey = struct{}{}
)

//PriorityQueue是基本结构
type PriorityQueue struct {
	Queues []chan interface{}
	wakeup chan struct{}
}

//New是PriorityQueue的构造函数
func New(n int, l int) *PriorityQueue {
	var queues = make([]chan interface{}, n)
	for i := range queues {
		queues[i] = make(chan interface{}, l)
	}
	return &PriorityQueue{
		Queues: queues,
		wakeup: make(chan struct{}, 1),
	}
}

//运行是从队列中弹出项目的永久循环
func (pq *PriorityQueue) Run(ctx context.Context, f func(interface{})) {
	top := len(pq.Queues) - 1
	p := top
READ:
	for {
		q := pq.Queues[p]
		select {
		case <-ctx.Done():
			return
		case x := <-q:
			log.Trace("priority.queue f(x)", "p", p, "len(Queues[p])", len(pq.Queues[p]))
			f(x)
			p = top
		default:
			if p > 0 {
				p--
				log.Trace("priority.queue p > 0", "p", p)
				continue READ
			}
			p = top
			select {
			case <-ctx.Done():
				return
			case <-pq.wakeup:
				log.Trace("priority.queue wakeup", "p", p)
			}
		}
	}
}

//push将项目推送到priority参数中指定的适当队列
//如果给定了上下文，它将一直等到推送该项或上下文中止为止。
func (pq *PriorityQueue) Push(x interface{}, p int) error {
	if p < 0 || p >= len(pq.Queues) {
		return errBadPriority
	}
	log.Trace("priority.queue push", "p", p, "len(Queues[p])", len(pq.Queues[p]))
	select {
	case pq.Queues[p] <- x:
	default:
		return ErrContention
	}
	select {
	case pq.wakeup <- wakey:
	default:
	}
	return nil
}

