
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450121802715136>

package tracing

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	cli "gopkg.in/urfave/cli.v1"
)

var Enabled bool = false

//tracingabledflag是用于启用跟踪集合的CLI标志名。
const TracingEnabledFlag = "tracing"

var (
	Closer io.Closer
)

var (
	TracingFlag = cli.BoolFlag{
		Name:  TracingEnabledFlag,
		Usage: "Enable tracing",
	}
	TracingEndpointFlag = cli.StringFlag{
		Name:  "tracing.endpoint",
		Usage: "Tracing endpoint",
		Value: "0.0.0.0:6831",
	}
	TracingSvcFlag = cli.StringFlag{
		Name:  "tracing.svc",
		Usage: "Tracing service name",
		Value: "swarm",
	}
)

//标志保存跟踪集合所需的所有命令行标志。
var Flags = []cli.Flag{
	TracingFlag,
	TracingEndpointFlag,
	TracingSvcFlag,
}

//init启用或禁用打开的跟踪系统。
func init() {
	for _, arg := range os.Args {
		if flag := strings.TrimLeft(arg, "-"); flag == TracingEnabledFlag {
			Enabled = true
		}
	}
}

func Setup(ctx *cli.Context) {
	if Enabled {
		log.Info("Enabling opentracing")
		var (
			endpoint = ctx.GlobalString(TracingEndpointFlag.Name)
			svc      = ctx.GlobalString(TracingSvcFlag.Name)
		)

		Closer = initTracer(endpoint, svc)
	}
}

func initTracer(endpoint, svc string) (closer io.Closer) {
//测试配置示例。用恒定采样法对每条记录道进行采样
//并使log span通过配置的logger记录每个跨度。
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  endpoint,
		},
	}

//示例记录器和度量工厂。使用github.com/uber/jaeger-client-go/log
//和github.com/uber/jaeger-lib/metrics分别绑定到实际日志和指标
//框架。
//jLogger：=jaegerLog.stdLogger
//jmetricsFactory：=metrics.nullFactory

//使用记录器和度量工厂初始化跟踪程序
	closer, err := cfg.InitGlobalTracer(
		svc,
//jaegercfg.记录器（jlogger）
//Jaegercfg.度量（jmetricsFactory）
//jaegercfg.observer（rpcmetrics.newobserver（jmetricsfactory，rpcmetrics.defaultnamenormalizer）），
	)
	if err != nil {
		log.Error("Could not initialize Jaeger tracer", "err", err)
	}

	return closer
}

