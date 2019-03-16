
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450067645861888>


package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/pkg/reexec"
	"github.com/ethereum/go-ethereum/internal/cmdtest"
)

type testEthkey struct {
	*cmdtest.TestCmd
}

//使用给定的命令行参数生成ethkey。
func runEthkey(t *testing.T, args ...string) *testEthkey {
	tt := new(testEthkey)
	tt.TestCmd = cmdtest.NewTestCmd(t, tt)
	tt.Run("ethkey-test", args...)
	return tt
}

func TestMain(m *testing.M) {
//如果我们在run ethkey中被执行为“ethkey测试”，请运行该应用程序。
	reexec.Register("ethkey-test", func() {
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	})
//检查我们是否被重新执行了
	if reexec.Init() {
		return
	}
	os.Exit(m.Run())
}

