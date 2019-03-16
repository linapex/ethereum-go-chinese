
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450106870992896>


package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/simulations"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
	"github.com/ethereum/go-ethereum/rpc"
)

var adapterType = flag.String("adapter", "sim", `node adapter to use (one of "sim", "exec" or "docker")`)

//main（）启动一个模拟网络，其中包含运行简单
//乒乓协议
func main() {
	flag.Parse()

//将日志级别设置为跟踪
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, log.StreamHandler(os.Stderr, log.TerminalFormat(false))))

//注册单个乒乓球服务
	services := map[string]adapters.ServiceFunc{
		"ping-pong": func(ctx *adapters.ServiceContext) (node.Service, error) {
			return newPingPongService(ctx.Config.ID), nil
		},
	}
	adapters.RegisterServices(services)

//创建节点适配器
	var adapter adapters.NodeAdapter

	switch *adapterType {

	case "sim":
		log.Info("using sim adapter")
		adapter = adapters.NewSimAdapter(services)

	case "exec":
		tmpdir, err := ioutil.TempDir("", "p2p-example")
		if err != nil {
			log.Crit("error creating temp dir", "err", err)
		}
		defer os.RemoveAll(tmpdir)
		log.Info("using exec adapter", "tmpdir", tmpdir)
		adapter = adapters.NewExecAdapter(tmpdir)

	default:
		log.Crit(fmt.Sprintf("unknown node adapter %q", *adapterType))
	}

//启动HTTP API
	log.Info("starting simulation server on 0.0.0.0:8888...")
	network := simulations.NewNetwork(adapter, &simulations.NetworkConfig{
		DefaultService: "ping-pong",
	})
	if err := http.ListenAndServe(":8888", simulations.NewServer(network)); err != nil {
		log.Crit("error starting simulation server", "err", err)
	}
}

//PingpongService在每个节点之间运行一个乒乓协议
//每隔10秒向所有连接的对等机发送一次ping，并接收一个pong-in
//返回
type pingPongService struct {
	id       enode.ID
	log      log.Logger
	received int64
}

func newPingPongService(id enode.ID) *pingPongService {
	return &pingPongService{
		id:  id,
		log: log.New("node.id", id),
	}
}

func (p *pingPongService) Protocols() []p2p.Protocol {
	return []p2p.Protocol{{
		Name:     "ping-pong",
		Version:  1,
		Length:   2,
		Run:      p.Run,
		NodeInfo: p.Info,
	}}
}

func (p *pingPongService) APIs() []rpc.API {
	return nil
}

func (p *pingPongService) Start(server *p2p.Server) error {
	p.log.Info("ping-pong service starting")
	return nil
}

func (p *pingPongService) Stop() error {
	p.log.Info("ping-pong service stopping")
	return nil
}

func (p *pingPongService) Info() interface{} {
	return struct {
		Received int64 `json:"received"`
	}{
		atomic.LoadInt64(&p.received),
	}
}

const (
	pingMsgCode = iota
	pongMsgCode
)

//run实现向对等端发送ping消息的乒乓协议
//每隔10秒，并用pong消息响应ping。
func (p *pingPongService) Run(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	log := p.log.New("peer.id", peer.ID())

	errC := make(chan error)
	go func() {
		for range time.Tick(10 * time.Second) {
			log.Info("sending ping")
			if err := p2p.Send(rw, pingMsgCode, "PING"); err != nil {
				errC <- err
				return
			}
		}
	}()
	go func() {
		for {
			msg, err := rw.ReadMsg()
			if err != nil {
				errC <- err
				return
			}
			payload, err := ioutil.ReadAll(msg.Payload)
			if err != nil {
				errC <- err
				return
			}
			log.Info("received message", "msg.code", msg.Code, "msg.payload", string(payload))
			atomic.AddInt64(&p.received, 1)
			if msg.Code == pingMsgCode {
				log.Info("sending pong")
				go p2p.Send(rw, pongMsgCode, "PONG")
			}
		}
	}()
	return <-errC
}

