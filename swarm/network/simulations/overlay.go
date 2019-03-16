
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450115016331264>


//您可以使用
//
//开始运行/swarm/network/simulations/overlay.go
package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
	"github.com/ethereum/go-ethereum/swarm/network"
	"github.com/ethereum/go-ethereum/swarm/state"
	colorable "github.com/mattn/go-colorable"
)

var (
	noDiscovery = flag.Bool("no-discovery", false, "disable discovery (useful if you want to load a snapshot)")
	vmodule     = flag.String("vmodule", "", "log filters for logger via Vmodule")
	verbosity   = flag.Int("verbosity", 0, "log filters for logger via Vmodule")
	httpSimPort = 8888
)

func init() {
	flag.Parse()
//初始化记录器
//这是有关如何使用vmodule筛选日志的演示
//提供-vmodule作为参数和逗号分隔值，例如：
//-vmodule overlay_test.go=4，simulations=3
//以上示例将overlay-test.go日志设置为级别4，而以“模拟”结尾的包设置为3
	if *vmodule != "" {
//仅当已提供标志时才启用模式匹配处理程序
		glogger := log.NewGlogHandler(log.StreamHandler(colorable.NewColorableStderr(), log.TerminalFormat(true)))
		if *verbosity > 0 {
			glogger.Verbosity(log.Lvl(*verbosity))
		}
		glogger.Vmodule(*vmodule)
		log.Root().SetHandler(glogger)
	}
}

type Simulation struct {
	mtx    sync.Mutex
	stores map[enode.ID]state.Store
}

func NewSimulation() *Simulation {
	return &Simulation{
		stores: make(map[enode.ID]state.Store),
	}
}

func (s *Simulation) NewService(ctx *adapters.ServiceContext) (node.Service, error) {
	node := ctx.Config.Node()
	s.mtx.Lock()
	store, ok := s.stores[node.ID()]
	if !ok {
		store = state.NewInmemoryStore()
		s.stores[node.ID()] = store
	}
	s.mtx.Unlock()

	addr := network.NewAddr(node)

	kp := network.NewKadParams()
	kp.NeighbourhoodSize = 2
	kp.MaxBinSize = 4
	kp.MinBinSize = 1
	kp.MaxRetries = 1000
	kp.RetryExponent = 2
	kp.RetryInterval = 1000000
	kad := network.NewKademlia(addr.Over(), kp)
	hp := network.NewHiveParams()
	hp.Discovery = !*noDiscovery
	hp.KeepAliveInterval = 300 * time.Millisecond

	config := &network.BzzConfig{
		OverlayAddr:  addr.Over(),
		UnderlayAddr: addr.Under(),
		HiveParams:   hp,
	}

	return network.NewBzz(config, kad, store, nil, nil), nil
}

//创建模拟网络
func newSimulationNetwork() *simulations.Network {

	s := NewSimulation()
	services := adapters.Services{
		"overlay": s.NewService,
	}
	adapter := adapters.NewSimAdapter(services)
	simNetwork := simulations.NewNetwork(adapter, &simulations.NetworkConfig{
		DefaultService: "overlay",
	})
	return simNetwork
}

//返回新的HTTP服务器
func newOverlaySim(sim *simulations.Network) *simulations.Server {
	return simulations.NewServer(sim)
}

//无功服务器
func main() {
//CPU优化
	runtime.GOMAXPROCS(runtime.NumCPU())
//运行SIM
	runOverlaySim()
}

func runOverlaySim() {
//创建模拟网络
	net := newSimulationNetwork()
//用它创建一个HTTP服务器
	sim := newOverlaySim(net)
	log.Info(fmt.Sprintf("starting simulation server on 0.0.0.0:%d...", httpSimPort))
//启动HTTP服务器
	http.ListenAndServe(fmt.Sprintf(":%d", httpSimPort), sim)
}

