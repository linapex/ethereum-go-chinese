
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450067952046080>


package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/tests"

	cli "gopkg.in/urfave/cli.v1"
)

var stateTestCommand = cli.Command{
	Action:    stateTestCmd,
	Name:      "statetest",
	Usage:     "executes the given state tests",
	ArgsUsage: "<file>",
}

//StateTestResult包含运行状态测试后的执行状态，任何
//可能发生的错误和最终状态的转储（如果请求）。
type StatetestResult struct {
	Name  string      `json:"name"`
	Pass  bool        `json:"pass"`
	Fork  string      `json:"fork"`
	Error string      `json:"error,omitempty"`
	State *state.Dump `json:"state,omitempty"`
}

func stateTestCmd(ctx *cli.Context) error {
	if len(ctx.Args().First()) == 0 {
		return errors.New("path-to-test argument required")
	}
//配置Go以太坊记录器
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(ctx.GlobalInt(VerbosityFlag.Name)))
	log.Root().SetHandler(glogger)

//配置EVM记录器
	config := &vm.LogConfig{
		DisableMemory: ctx.GlobalBool(DisableMemoryFlag.Name),
		DisableStack:  ctx.GlobalBool(DisableStackFlag.Name),
	}
	var (
		tracer   vm.Tracer
		debugger *vm.StructLogger
	)
	switch {
	case ctx.GlobalBool(MachineFlag.Name):
		tracer = vm.NewJSONLogger(config, os.Stderr)

	case ctx.GlobalBool(DebugFlag.Name):
		debugger = vm.NewStructLogger(config)
		tracer = debugger

	default:
		debugger = vm.NewStructLogger(config)
	}
//从输入文件加载测试内容
	src, err := ioutil.ReadFile(ctx.Args().First())
	if err != nil {
		return err
	}
	var tests map[string]tests.StateTest
	if err = json.Unmarshal(src, &tests); err != nil {
		return err
	}
//迭代所有测试，运行它们并聚合结果
	cfg := vm.Config{
		Tracer: tracer,
		Debug:  ctx.GlobalBool(DebugFlag.Name) || ctx.GlobalBool(MachineFlag.Name),
	}
	results := make([]StatetestResult, 0, len(tests))
	for key, test := range tests {
		for _, st := range test.Subtests() {
//运行测试并汇总结果
			result := &StatetestResult{Name: key, Fork: st.Fork, Pass: true}
			state, err := test.Run(st, cfg)
//为evmlab跟踪打印状态根
			if ctx.GlobalBool(MachineFlag.Name) && state != nil {
				fmt.Fprintf(os.Stderr, "{\"stateRoot\": \"%x\"}\n", state.IntermediateRoot(false))
			}
			if err != nil {
//测试失败，标记为“是”，并转储任何状态以帮助调试
				result.Pass, result.Error = false, err.Error()
				if ctx.GlobalBool(DumpFlag.Name) && state != nil {
					dump := state.RawDump()
					result.State = &dump
				}
			}

			results = append(results, *result)

//打印收集的任何结构化日志
			if ctx.GlobalBool(DebugFlag.Name) {
				if debugger != nil {
					fmt.Fprintln(os.Stderr, "#### TRACE ####")
					vm.WriteTrace(os.Stderr, debugger.StructLogs())
				}
			}
		}
	}
	out, _ := json.MarshalIndent(results, "", "  ")
	fmt.Println(string(out))
	return nil
}

