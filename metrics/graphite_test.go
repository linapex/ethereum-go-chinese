
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450098872455168>

package metrics

import (
	"net"
	"time"
)

func ExampleGraphite() {
	addr, _ := net.ResolveTCPAddr("net", ":2003")
	go Graphite(DefaultRegistry, 1*time.Second, "some.prefix", addr)
}

func ExampleGraphiteWithConfig() {
	addr, _ := net.ResolveTCPAddr("net", ":2003")
	go GraphiteWithConfig(GraphiteConfig{
		Addr:          addr,
		Registry:      DefaultRegistry,
		FlushInterval: 1 * time.Second,
		DurationUnit:  time.Millisecond,
		Percentiles:   []float64{0.5, 0.75, 0.99, 0.999},
	})
}

