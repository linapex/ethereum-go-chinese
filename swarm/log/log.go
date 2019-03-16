
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450113250529280>

package log

import (
	l "github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
)

const (
//CallDepth设置为1，以影响
//调用l.output（）时跳过堆栈帧为1的日志消息
	CallDepth = 1
)

//warn是日志的方便别名。warn with stats
func Warn(msg string, ctx ...interface{}) {
	metrics.GetOrRegisterCounter("warn", nil).Inc(1)
	l.Output(msg, l.LvlWarn, CallDepth, ctx...)
}

//错误是日志的方便别名。使用stats时出错
func Error(msg string, ctx ...interface{}) {
	metrics.GetOrRegisterCounter("error", nil).Inc(1)
	l.Output(msg, l.LvlError, CallDepth, ctx...)
}

//crit是log.crit和stats的方便别名
func Crit(msg string, ctx ...interface{}) {
	metrics.GetOrRegisterCounter("crit", nil).Inc(1)
	l.Output(msg, l.LvlCrit, CallDepth, ctx...)
}

//信息是日志的一个方便别名。带有统计信息的信息
func Info(msg string, ctx ...interface{}) {
	metrics.GetOrRegisterCounter("info", nil).Inc(1)
	l.Output(msg, l.LvlInfo, CallDepth, ctx...)
}

//调试是日志的一个方便别名。使用stats进行调试
func Debug(msg string, ctx ...interface{}) {
	metrics.GetOrRegisterCounter("debug", nil).Inc(1)
	l.Output(msg, l.LvlDebug, CallDepth, ctx...)
}

//trace是日志的方便别名。trace with stats
func Trace(msg string, ctx ...interface{}) {
	metrics.GetOrRegisterCounter("trace", nil).Inc(1)
	l.Output(msg, l.LvlTrace, CallDepth, ctx...)
}

