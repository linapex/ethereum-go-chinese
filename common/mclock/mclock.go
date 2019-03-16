
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450074025398272>


//包mclock是单调时钟源的包装器
package mclock

import (
	"time"

	"github.com/aristanetworks/goarista/monotime"
)

//Abstime代表绝对单调时间。
type AbsTime time.Duration

//现在返回当前绝对单调时间。
func Now() AbsTime {
	return AbsTime(monotime.Now())
}

//添加返回t+d。
func (t AbsTime) Add(d time.Duration) AbsTime {
	return t + AbsTime(d)
}

//时钟接口使得用
//模拟时钟。
type Clock interface {
	Now() AbsTime
	Sleep(time.Duration)
	After(time.Duration) <-chan time.Time
}

//系统使用系统时钟实现时钟。
type System struct{}

//现在实现时钟。
func (System) Now() AbsTime {
	return AbsTime(monotime.Now())
}

//睡眠实现时钟。
func (System) Sleep(d time.Duration) {
	time.Sleep(d)
}

//在执行时钟之后。
func (System) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

