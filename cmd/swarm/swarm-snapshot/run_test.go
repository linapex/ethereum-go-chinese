
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450071995355136>


package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/docker/docker/pkg/reexec"
	"github.com/ethereum/go-ethereum/internal/cmdtest"
)

func init() {
	reexec.Register("swarm-snapshot", func() {
		if err := newApp().Run(os.Args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	})
}

func runSnapshot(t *testing.T, args ...string) *cmdtest.TestCmd {
	tt := cmdtest.NewTestCmd(t, nil)
	tt.Run("swarm-snapshot", args...)
	return tt
}

func TestMain(m *testing.M) {
	if reexec.Init() {
		return
	}
	os.Exit(m.Run())
}

