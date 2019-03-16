
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450114856947712>


package simulation

import (
	"context"
	"errors"
	"flag"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
	"github.com/mattn/go-colorable"
)

var (
	loglevel = flag.Int("loglevel", 2, "verbosity of logs")
)

func init() {
	flag.Parse()
	log.PrintOrigins(true)
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(*loglevel), log.StreamHandler(colorable.NewColorableStderr(), log.TerminalFormat(true))))
}

//如果run方法调用runfunc并正确处理上下文，则测试run。
func TestRun(t *testing.T) {
	sim := New(noopServiceFuncMap)
	defer sim.Close()

	t.Run("call", func(t *testing.T) {
		expect := "something"
		var got string
		r := sim.Run(context.Background(), func(ctx context.Context, sim *Simulation) error {
			got = expect
			return nil
		})

		if r.Error != nil {
			t.Errorf("unexpected error: %v", r.Error)
		}
		if got != expect {
			t.Errorf("expected %q, got %q", expect, got)
		}
	})

	t.Run("cancellation", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		r := sim.Run(ctx, func(ctx context.Context, sim *Simulation) error {
			time.Sleep(time.Second)
			return nil
		})

		if r.Error != context.DeadlineExceeded {
			t.Errorf("unexpected error: %v", r.Error)
		}
	})

	t.Run("context value and duration", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), "hey", "there")
		sleep := 50 * time.Millisecond

		r := sim.Run(ctx, func(ctx context.Context, sim *Simulation) error {
			if ctx.Value("hey") != "there" {
				return errors.New("expected context value not passed")
			}
			time.Sleep(sleep)
			return nil
		})

		if r.Error != nil {
			t.Errorf("unexpected error: %v", r.Error)
		}
		if r.Duration < sleep {
			t.Errorf("reported run duration less then expected: %s", r.Duration)
		}
	})
}

//testclose测试是close方法，触发所有close函数，所有节点都不再是up。
func TestClose(t *testing.T) {
	var mu sync.Mutex
	var cleanupCount int

	sleep := 50 * time.Millisecond

	sim := New(map[string]ServiceFunc{
		"noop": func(ctx *adapters.ServiceContext, b *sync.Map) (node.Service, func(), error) {
			return newNoopService(), func() {
				time.Sleep(sleep)
				mu.Lock()
				defer mu.Unlock()
				cleanupCount++
			}, nil
		},
	})

	nodeCount := 30

	_, err := sim.AddNodes(nodeCount)
	if err != nil {
		t.Fatal(err)
	}

	var upNodeCount int
	for _, n := range sim.Net.GetNodes() {
		if n.Up {
			upNodeCount++
		}
	}
	if upNodeCount != nodeCount {
		t.Errorf("all nodes should be up, insted only %v are up", upNodeCount)
	}

	sim.Close()

	if cleanupCount != nodeCount {
		t.Errorf("number of cleanups expected %v, got %v", nodeCount, cleanupCount)
	}

	upNodeCount = 0
	for _, n := range sim.Net.GetNodes() {
		if n.Up {
			upNodeCount++
		}
	}
	if upNodeCount != 0 {
		t.Errorf("all nodes should be down, insted %v are up", upNodeCount)
	}
}

//testdone检查close方法是否触发了done通道的关闭。
func TestDone(t *testing.T) {
	sim := New(noopServiceFuncMap)
	sleep := 50 * time.Millisecond
	timeout := 2 * time.Second

	start := time.Now()
	go func() {
		time.Sleep(sleep)
		sim.Close()
	}()

	select {
	case <-time.After(timeout):
		t.Error("done channel closing timed out")
	case <-sim.Done():
		if d := time.Since(start); d < sleep {
			t.Errorf("done channel closed sooner then expected: %s", d)
		}
	}
}

//用于不执行任何操作的常规服务的帮助者映射
var noopServiceFuncMap = map[string]ServiceFunc{
	"noop": noopServiceFunc,
}

//最基本的noop服务的助手函数
func noopServiceFunc(_ *adapters.ServiceContext, _ *sync.Map) (node.Service, func(), error) {
	return newNoopService(), nil, nil
}

func newNoopService() node.Service {
	return &noopService{}
}

//最基本的noop服务的助手函数
//不同类型的，然后NoopService进行测试
//一个节点上有多个服务。
func noopService2Func(_ *adapters.ServiceContext, _ *sync.Map) (node.Service, func(), error) {
	return new(noopService2), nil, nil
}

//NoopService2是不做任何事情的服务
//但实现了node.service接口。
type noopService2 struct {
	simulations.NoopService
}

type noopService struct {
	simulations.NoopService
}

