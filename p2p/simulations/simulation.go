
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450107256868865>


package simulations

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

//模拟为在模拟网络中运行操作提供了一个框架
//然后等待期望得到满足
type Simulation struct {
	network *Network
}

//新模拟返回在给定网络中运行的新模拟
func NewSimulation(network *Network) *Simulation {
	return &Simulation{
		network: network,
	}
}

//运行通过执行步骤的操作执行模拟步骤，并
//然后等待步骤的期望得到满足
func (s *Simulation) Run(ctx context.Context, step *Step) (result *StepResult) {
	result = newStepResult()

	result.StartedAt = time.Now()
	defer func() { result.FinishedAt = time.Now() }()

//在步骤期间监视网络事件
	stop := s.watchNetwork(result)
	defer stop()

//执行操作
	if err := step.Action(ctx); err != nil {
		result.Error = err
		return
	}

//等待所有节点期望通过、错误或超时
	nodes := make(map[enode.ID]struct{}, len(step.Expect.Nodes))
	for _, id := range step.Expect.Nodes {
		nodes[id] = struct{}{}
	}
	for len(result.Passes) < len(nodes) {
		select {
		case id := <-step.Trigger:
//如果不检查节点，则跳过
			if _, ok := nodes[id]; !ok {
				continue
			}

//如果节点已通过，则跳过
			if _, ok := result.Passes[id]; ok {
				continue
			}

//运行节点期望检查
			pass, err := step.Expect.Check(ctx, id)
			if err != nil {
				result.Error = err
				return
			}
			if pass {
				result.Passes[id] = time.Now()
			}
		case <-ctx.Done():
			result.Error = ctx.Err()
			return
		}
	}

	return
}

func (s *Simulation) watchNetwork(result *StepResult) func() {
	stop := make(chan struct{})
	done := make(chan struct{})
	events := make(chan *Event)
	sub := s.network.Events().Subscribe(events)
	go func() {
		defer close(done)
		defer sub.Unsubscribe()
		for {
			select {
			case event := <-events:
				result.NetworkEvents = append(result.NetworkEvents, event)
			case <-stop:
				return
			}
		}
	}()
	return func() {
		close(stop)
		<-done
	}
}

type Step struct {
//操作是为此步骤执行的操作
	Action func(context.Context) error

//触发器是一个接收节点ID并触发
//该节点的预期检查
	Trigger chan enode.ID

//Expect是执行此步骤时等待的期望
	Expect *Expectation
}

type Expectation struct {
//节点是要检查的节点列表
	Nodes []enode.ID

//检查给定节点是否满足预期
	Check func(context.Context, enode.ID) (bool, error)
}

func newStepResult() *StepResult {
	return &StepResult{
		Passes: make(map[enode.ID]time.Time),
	}
}

type StepResult struct {
//错误是运行步骤时遇到的错误
	Error error

//Startedat是步骤开始的时间
	StartedAt time.Time

//FinishedAt是步骤完成的时间。
	FinishedAt time.Time

//传递是成功节点期望的时间戳
	Passes map[enode.ID]time.Time

//NetworkEvents是在步骤中发生的网络事件。
	NetworkEvents []*Event
}

