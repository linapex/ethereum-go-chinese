
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450114424934400>


package simulation

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
)

func TestSimulationWithHTTPServer(t *testing.T) {
	log.Debug("Init simulation")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	sim := New(
		map[string]ServiceFunc{
			"noop": func(_ *adapters.ServiceContext, b *sync.Map) (node.Service, func(), error) {
				return newNoopService(), nil, nil
			},
		}).WithServer(DefaultHTTPSimAddr)
	defer sim.Close()
	log.Debug("Done.")

	_, err := sim.AddNode()
	if err != nil {
		t.Fatal(err)
	}

	log.Debug("Starting sim round and let it time out...")
//第一个不发送到通道的运行测试
//阻止模拟，让它超时
	result := sim.Run(ctx, func(ctx context.Context, sim *Simulation) error {
		log.Debug("Just start the sim without any action and wait for the timeout")
//使用睡眠确保模拟不会在超时之前终止。
		time.Sleep(2 * time.Second)
		return nil
	})

	if result.Error != nil {
		if result.Error.Error() == "context deadline exceeded" {
			log.Debug("Expected timeout error received")
		} else {
			t.Fatal(result.Error)
		}
	}

//现在再次运行它并在等待通道上发送预期信号，
//然后关闭模拟
	log.Debug("Starting sim round and wait for frontend signal...")
//这一次超时时间应该足够长，这样它就不会过早启动。
	ctx, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	errC := make(chan error, 1)
	go triggerSimulationRun(t, errC)
	result = sim.Run(ctx, func(ctx context.Context, sim *Simulation) error {
		log.Debug("This run waits for the run signal from `frontend`...")
//确保睡眠状态下，模拟不会在收到信号之前终止。
		time.Sleep(2 * time.Second)
		return nil
	})
	if result.Error != nil {
		t.Fatal(result.Error)
	}
	if err := <-errC; err != nil {
		t.Fatal(err)
	}
	log.Debug("Test terminated successfully")
}

func triggerSimulationRun(t *testing.T, errC chan error) {
//我们需要首先等待SIM HTTP服务器开始运行…
	time.Sleep(2 * time.Second)
//然后我们可以发送信号

	log.Debug("Sending run signal to simulation: POST /runsim...")
resp, err := http.Post(fmt.Sprintf("http://本地主机%s/runsim“，defaulthttpsimaddr”，“application/json”，nil）
	if err != nil {
		errC <- fmt.Errorf("Request failed: %v", err)
		return
	}
	log.Debug("Signal sent")
	if resp.StatusCode != http.StatusOK {
		errC <- fmt.Errorf("err %s", resp.Status)
		return
	}
	errC <- resp.Body.Close()
}

