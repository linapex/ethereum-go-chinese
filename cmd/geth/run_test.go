
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450068807684096>


package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/docker/docker/pkg/reexec"
	"github.com/ethereum/go-ethereum/internal/cmdtest"
)

func tmpdir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "geth-test")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

type testgeth struct {
	*cmdtest.TestCmd

//Expect的模板变量
	Datadir   string
	Etherbase string
}

func init() {
//运行应用程序，如果我们已经被执行为“Geth测试”在朗格思。
	reexec.Register("geth-test", func() {
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	})
}

func TestMain(m *testing.M) {
//检查我们是否被重新执行了
	if reexec.Init() {
		return
	}
	os.Exit(m.Run())
}

//使用给定的命令行参数生成geth。如果参数未设置--datadir，则
//子G获取临时数据目录。
func runGeth(t *testing.T, args ...string) *testgeth {
	tt := &testgeth{}
	tt.TestCmd = cmdtest.NewTestCmd(t, tt)
	for i, arg := range args {
		switch {
		case arg == "-datadir" || arg == "--datadir":
			if i < len(args)-1 {
				tt.Datadir = args[i+1]
			}
		case arg == "-etherbase" || arg == "--etherbase":
			if i < len(args)-1 {
				tt.Etherbase = args[i+1]
			}
		}
	}
	if tt.Datadir == "" {
		tt.Datadir = tmpdir(t)
		tt.Cleanup = func() { os.RemoveAll(tt.Datadir) }
		args = append([]string{"-datadir", tt.Datadir}, args...)
//如果下面发生故障，请删除临时datadir。
		defer func() {
			if t.Failed() {
				tt.Cleanup()
			}
		}()
	}

//启动“GETH”。这实际上运行了测试二进制文件，但是testmain
//函数将阻止任何测试运行。
	tt.Run("geth-test", args...)

	return tt
}

