
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450114387185664>


package simulation

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/simulations"
)

//包默认值。
var (
	DefaultHTTPSimAddr = ":8888"
)

//WithServer实现用于模拟的生成器模式构造函数
//从HTTP服务器开始
func (s *Simulation) WithServer(addr string) *Simulation {
//如果未提供，则分配默认地址
	if addr == "" {
		addr = DefaultHTTPSimAddr
	}
	log.Info(fmt.Sprintf("Initializing simulation server on %s...", addr))
//初始化HTTP服务器
	s.handler = simulations.NewServer(s.Net)
	s.runC = make(chan struct{})
//向HTTP服务器添加特定于Swarm的路由
	s.addSimulationRoutes()
	s.httpSrv = &http.Server{
		Addr:    addr,
		Handler: s.handler,
	}
	go func() {
		err := s.httpSrv.ListenAndServe()
		if err != nil {
			log.Error("Error starting the HTTP server", "error", err)
		}
	}()
	return s
}

//注册其他HTTP路由
func (s *Simulation) addSimulationRoutes() {
	s.handler.POST("/runsim", s.RunSimulation)
}

//运行模拟是实际的端点后运行程序
func (s *Simulation) RunSimulation(w http.ResponseWriter, req *http.Request) {
	log.Debug("RunSimulation endpoint running")
	s.runC <- struct{}{}
	w.WriteHeader(http.StatusOK)
}

