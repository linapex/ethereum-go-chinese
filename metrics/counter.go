
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450097995845632>

package metrics

import (
	"sync/atomic"
)

//计数器保存一个可以递增和递减的Int64值。
type Counter interface {
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
	Snapshot() Counter
}

//GetOrRegisterCounter返回现有计数器或构造和注册
//一个新的标准计数器。
func GetOrRegisterCounter(name string, r Registry) Counter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewCounter).(Counter)
}

//GetOrRegisterCounterForced返回现有计数器或构造并注册
//新计数器，无论全局开关是否启用。
//一旦计数器对
//允许垃圾收集。
func GetOrRegisterCounterForced(name string, r Registry) Counter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewCounterForced).(Counter)
}

//NewCounter构造一个新的标准计数器。
func NewCounter() Counter {
	if !Enabled {
		return NilCounter{}
	}
	return &StandardCounter{0}
}

//NewCounterForced构造新的StandardCounter并返回它，无论
//全局开关是否启用。
func NewCounterForced() Counter {
	return &StandardCounter{0}
}

//NewRegisteredCounter构造并注册新的标准计数器。
func NewRegisteredCounter(name string, r Registry) Counter {
	c := NewCounter()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

//newregisteredcounterforce构造并注册新的标准计数器
//并启动Goroutine，无论全局开关是否启用。
//一旦计数器对
//允许垃圾收集。
func NewRegisteredCounterForced(name string, r Registry) Counter {
	c := NewCounterForced()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

//CounterSnapshot是另一个计数器的只读副本。
type CounterSnapshot int64

//清晰的恐慌。
func (CounterSnapshot) Clear() {
	panic("Clear called on a CounterSnapshot")
}

//count返回快照拍摄时的计数。
func (c CounterSnapshot) Count() int64 { return int64(c) }

//十足恐慌。
func (CounterSnapshot) Dec(int64) {
	panic("Dec called on a CounterSnapshot")
}

//公司恐慌。
func (CounterSnapshot) Inc(int64) {
	panic("Inc called on a CounterSnapshot")
}

//快照返回快照。
func (c CounterSnapshot) Snapshot() Counter { return c }

//nilcounter是一个禁止操作的计数器。
type NilCounter struct{}

//清除是不可操作的。
func (NilCounter) Clear() {}

//计数是不允许的。
func (NilCounter) Count() int64 { return 0 }

//DEC是NO-OP。
func (NilCounter) Dec(i int64) {}

//公司是一个NO-OP。
func (NilCounter) Inc(i int64) {}

//快照是不可操作的。
func (NilCounter) Snapshot() Counter { return NilCounter{} }

//StandardCounter是计数器的标准实现，它使用
//同步/atomic包以管理单个int64值。
type StandardCounter struct {
	count int64
}

//清除将计数器设置为零。
func (c *StandardCounter) Clear() {
	atomic.StoreInt64(&c.count, 0)
}

//count返回当前计数。
func (c *StandardCounter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

//DEC按给定的数量递减计数器。
func (c *StandardCounter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

//inc将计数器递增给定的数量。
func (c *StandardCounter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}

//快照返回计数器的只读副本。
func (c *StandardCounter) Snapshot() Counter {
	return CounterSnapshot(c.Count())
}

