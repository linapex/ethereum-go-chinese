
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450099468046336>

//Coda Hale度量库的Go端口
//
//<https://github.com/rcrowley/go-metrics>
//
//Coda Hale的原创作品：<https://github.com/coda hale/metrics>
package metrics

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

//对于所有的
//标准度量。如果为真，则返回的度量是存根。
//
//这种全局杀伤开关有助于量化观察者效应，并使
//以减少混乱的PPROF配置文件。
var Enabled bool = false

//MetricsEnabledFlag是用于启用度量集合的CLI标志名。
const MetricsEnabledFlag = "metrics"
const DashboardEnabledFlag = "dashboard"

//init启用或禁用度量系统。因为我们以前需要这个
//任何其他代码都可以创建仪表和计时器，实际上我们会做一个丑陋的黑客。
//并查看命令行参数中的度量标志。
func init() {
	for _, arg := range os.Args {
		if flag := strings.TrimLeft(arg, "-"); flag == MetricsEnabledFlag || flag == DashboardEnabledFlag {
			log.Info("Enabling metrics collection")
			Enabled = true
		}
	}
}

//CollectProcessMetrics定期收集有关运行的各种指标
//过程。
func CollectProcessMetrics(refresh time.Duration) {
//如果禁用度量系统，则短路
	if !Enabled {
		return
	}
//创建各种数据收集器
	memstats := make([]*runtime.MemStats, 2)
	diskstats := make([]*DiskStats, 2)
	for i := 0; i < len(memstats); i++ {
		memstats[i] = new(runtime.MemStats)
		diskstats[i] = new(DiskStats)
	}
//定义要收集的各种指标
	memAllocs := GetOrRegisterMeter("system/memory/allocs", DefaultRegistry)
	memFrees := GetOrRegisterMeter("system/memory/frees", DefaultRegistry)
	memInuse := GetOrRegisterMeter("system/memory/inuse", DefaultRegistry)
	memPauses := GetOrRegisterMeter("system/memory/pauses", DefaultRegistry)

	var diskReads, diskReadBytes, diskWrites, diskWriteBytes Meter
	var diskReadBytesCounter, diskWriteBytesCounter Counter
	if err := ReadDiskStats(diskstats[0]); err == nil {
		diskReads = GetOrRegisterMeter("system/disk/readcount", DefaultRegistry)
		diskReadBytes = GetOrRegisterMeter("system/disk/readdata", DefaultRegistry)
		diskReadBytesCounter = GetOrRegisterCounter("system/disk/readbytes", DefaultRegistry)
		diskWrites = GetOrRegisterMeter("system/disk/writecount", DefaultRegistry)
		diskWriteBytes = GetOrRegisterMeter("system/disk/writedata", DefaultRegistry)
		diskWriteBytesCounter = GetOrRegisterCounter("system/disk/writebytes", DefaultRegistry)
	} else {
		log.Debug("Failed to read disk metrics", "err", err)
	}
//重复加载不同的统计数据并更新仪表
	for i := 1; ; i++ {
		location1 := i % 2
		location2 := (i - 1) % 2

		runtime.ReadMemStats(memstats[location1])
		memAllocs.Mark(int64(memstats[location1].Mallocs - memstats[location2].Mallocs))
		memFrees.Mark(int64(memstats[location1].Frees - memstats[location2].Frees))
		memInuse.Mark(int64(memstats[location1].Alloc - memstats[location2].Alloc))
		memPauses.Mark(int64(memstats[location1].PauseTotalNs - memstats[location2].PauseTotalNs))

		if ReadDiskStats(diskstats[location1]) == nil {
			diskReads.Mark(diskstats[location1].ReadCount - diskstats[location2].ReadCount)
			diskReadBytes.Mark(diskstats[location1].ReadBytes - diskstats[location2].ReadBytes)
			diskWrites.Mark(diskstats[location1].WriteCount - diskstats[location2].WriteCount)
			diskWriteBytes.Mark(diskstats[location1].WriteBytes - diskstats[location2].WriteBytes)

			diskReadBytesCounter.Inc(diskstats[location1].ReadBytes - diskstats[location2].ReadBytes)
			diskWriteBytesCounter.Inc(diskstats[location1].WriteBytes - diskstats[location2].WriteBytes)
		}
		time.Sleep(refresh)
	}
}

