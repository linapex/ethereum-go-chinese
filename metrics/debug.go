
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450098079731712>

package metrics

import (
	"runtime/debug"
	"time"
)

var (
	debugMetrics struct {
		GCStats struct {
			LastGC Gauge
			NumGC  Gauge
			Pause  Histogram
//暂停序列柱状图
			PauseTotal Gauge
		}
		ReadGCStats Timer
	}
	gcStats debug.GCStats
)

//捕获在中导出的Go垃圾收集器统计信息的新值
//调试.gcstats。这被称为Goroutine。
func CaptureDebugGCStats(r Registry, d time.Duration) {
	for range time.Tick(d) {
		CaptureDebugGCStatsOnce(r)
	}
}

//捕获在中导出的Go垃圾收集器统计信息的新值
//调试.gcstats。这被设计为在后台goroutine中调用。
//提供尚未提供给RegisterDebuggCStats的注册表将
//恐慌。
//
//注意（但更不用说），因为debug.readgcstats调用
//C函数runtime·lock（runtime·mheap），它虽然不能阻止世界
//手术，不是你一直想做的事情。
func CaptureDebugGCStatsOnce(r Registry) {
	lastGC := gcStats.LastGC
	t := time.Now()
	debug.ReadGCStats(&gcStats)
	debugMetrics.ReadGCStats.UpdateSince(t)

	debugMetrics.GCStats.LastGC.Update(gcStats.LastGC.UnixNano())
	debugMetrics.GCStats.NumGC.Update(gcStats.NumGC)
	if lastGC != gcStats.LastGC && 0 < len(gcStats.Pause) {
		debugMetrics.GCStats.Pause.Update(int64(gcStats.Pause[0]))
	}
//debugmetrics.gcstats.pauseQuantiles.update（gcstats.pauseQuantiles）
	debugMetrics.GCStats.PauseTotal.Update(int64(gcStats.PauseTotal))
}

//注册中导出的Go垃圾收集器统计信息的度量
//调试.gcstats。这些指标是通过其完全限定的GO符号命名的，
//即debug.gcstats.pausetotal。
func RegisterDebugGCStats(r Registry) {
	debugMetrics.GCStats.LastGC = NewGauge()
	debugMetrics.GCStats.NumGC = NewGauge()
	debugMetrics.GCStats.Pause = NewHistogram(NewExpDecaySample(1028, 0.015))
//debugmetrics.gcstats.pausequeantiles=newhistogram（newexpdecaysample（1028，0.015））。
	debugMetrics.GCStats.PauseTotal = NewGauge()
	debugMetrics.ReadGCStats = NewTimer()

	r.Register("debug.GCStats.LastGC", debugMetrics.GCStats.LastGC)
	r.Register("debug.GCStats.NumGC", debugMetrics.GCStats.NumGC)
	r.Register("debug.GCStats.Pause", debugMetrics.GCStats.Pause)
//r.register（“debug.gcstats.pausequeantiles”，debugmetrics.gcstats.pausequeantiles）
	r.Register("debug.GCStats.PauseTotal", debugMetrics.GCStats.PauseTotal)
	r.Register("debug.ReadGCStats", debugMetrics.ReadGCStats)
}

//为gcstats分配初始切片。暂停以避免在
//正常运行。
func init() {
	gcStats.Pause = make([]time.Duration, 11)
}

