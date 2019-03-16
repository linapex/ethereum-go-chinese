
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450093465997312>


package les

import "sync"

//execqueue实现在单个线程中执行函数调用的队列，
//与排队的顺序相同。
type execQueue struct {
	mu        sync.Mutex
	cond      *sync.Cond
	funcs     []func()
	closeWait chan struct{}
}

//NeXeExcel队列创建一个新的执行队列。
func newExecQueue(capacity int) *execQueue {
	q := &execQueue{funcs: make([]func(), 0, capacity)}
	q.cond = sync.NewCond(&q.mu)
	go q.loop()
	return q
}

func (q *execQueue) loop() {
	for f := q.waitNext(false); f != nil; f = q.waitNext(true) {
		f()
	}
	close(q.closeWait)
}

func (q *execQueue) waitNext(drop bool) (f func()) {
	q.mu.Lock()
	if drop {
//删除刚刚执行的函数。我们在这里而不是在什么时候
//出列so len（q.funcs）包括正在运行的函数。
		q.funcs = append(q.funcs[:0], q.funcs[1:]...)
	}
	for !q.isClosed() {
		if len(q.funcs) > 0 {
			f = q.funcs[0]
			break
		}
		q.cond.Wait()
	}
	q.mu.Unlock()
	return f
}

func (q *execQueue) isClosed() bool {
	return q.closeWait != nil
}

//如果可以向执行队列中添加更多函数调用，则can queue返回true。
func (q *execQueue) canQueue() bool {
	q.mu.Lock()
	ok := !q.isClosed() && len(q.funcs) < cap(q.funcs)
	q.mu.Unlock()
	return ok
}

//queue向执行队列添加函数调用。如果成功，则返回true。
func (q *execQueue) queue(f func()) bool {
	q.mu.Lock()
	ok := !q.isClosed() && len(q.funcs) < cap(q.funcs)
	if ok {
		q.funcs = append(q.funcs, f)
		q.cond.Signal()
	}
	q.mu.Unlock()
	return ok
}

//quit停止执行队列。
//quit在返回之前等待当前执行完成。
func (q *execQueue) quit() {
	q.mu.Lock()
	if !q.isClosed() {
		q.closeWait = make(chan struct{})
		q.cond.Signal()
	}
	q.mu.Unlock()
	<-q.closeWait
}

