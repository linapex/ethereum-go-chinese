
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450099526766592>

package metrics

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"testing"
	"time"
)

const FANOUT = 128

//停止编译器在调试期间抱怨。
var (
	_ = ioutil.Discard
	_ = log.LstdFlags
)

func BenchmarkMetrics(b *testing.B) {
	r := NewRegistry()
	c := NewRegisteredCounter("counter", r)
	g := NewRegisteredGauge("gauge", r)
	gf := NewRegisteredGaugeFloat64("gaugefloat64", r)
	h := NewRegisteredHistogram("histogram", r, NewUniformSample(100))
	m := NewRegisteredMeter("meter", r)
	t := NewRegisteredTimer("timer", r)
	RegisterDebugGCStats(r)
	RegisterRuntimeMemStats(r)
	b.ResetTimer()
	ch := make(chan bool)

	wgD := &sync.WaitGroup{}
 /*
  添加（1）
  转到函数（）
   推迟wgd.done（）
   //log.println（“go capturedebuggstats”）。
   对于{
    选择{
    案例<CH：
     //log.println（“完成capturedebuggcstats”）
     返回
    违约：
     CaptureDebuggCStatsonce（右）
    }
   }
  }（）
 /*/


	wgR := &sync.WaitGroup{}
 /*
 添加（1）
 转到函数（）
  推迟wgr.done（）
  //log.println（“go captureruntimemstats”）。
  对于{
   选择{
   案例<CH：
    //log.println（“done captureruntimemstats”）。
    返回
   违约：
    捕获者untimemstattsone（r）
   }
  }
 }（）
 /*/


	wgW := &sync.WaitGroup{}
 /*
  添加（1）
  转到函数（）
   推迟wgw.done（）
   //日志.println（“go write”）
   对于{
    选择{
    案例<CH：
     //log.println（“完成写入”）
     返回
    违约：
     一次写入（r，ioutil.discard）
    }
   }
  }（）
 /*/


	wg := &sync.WaitGroup{}
	wg.Add(FANOUT)
	for i := 0; i < FANOUT; i++ {
		go func(i int) {
			defer wg.Done()
//log.println（“开始”，i）
			for i := 0; i < b.N; i++ {
				c.Inc(1)
				g.Update(int64(i))
				gf.Update(float64(i))
				h.Update(int64(i))
				m.Mark(1)
				t.Update(1)
			}
//log.println（“完成”，i）
		}(i)
	}
	wg.Wait()
	close(ch)
	wgD.Wait()
	wgR.Wait()
	wgW.Wait()
}

func Example() {
	c := NewCounter()
	Register("money", c)
	c.Inc(17)

//螺纹安全注册
	t := GetOrRegisterTimer("db.get.latency", nil)
	t.Time(func() { time.Sleep(10 * time.Millisecond) })
	t.Update(1)

	fmt.Println(c.Count())
	fmt.Println(t.Min())
//产量：17
//一
}

