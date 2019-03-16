
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450071856943104>


package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/p2p/simulations"
)

//testsnapshotcreate是一个高级别的e2e测试，用于测试快照生成。
//它运行一些带有不同标志值和生成的加载的“创建”命令
//快照文件以验证其内容。
func TestSnapshotCreate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip()
	}

	for _, v := range []struct {
		name     string
		nodes    int
		services string
	}{
		{
			name: "defaults",
		},
		{
			name:  "more nodes",
			nodes: defaultNodes + 5,
		},
		{
			name:     "services",
			services: "stream,pss,zorglub",
		},
		{
			name:     "services with bzz",
			services: "bzz,pss",
		},
	} {
		t.Run(v.name, func(t *testing.T) {
			t.Parallel()

			file, err := ioutil.TempFile("", "swarm-snapshot")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(file.Name())

			if err = file.Close(); err != nil {
				t.Error(err)
			}

			args := []string{"create"}
			if v.nodes > 0 {
				args = append(args, "--nodes", strconv.Itoa(v.nodes))
			}
			if v.services != "" {
				args = append(args, "--services", v.services)
			}
			testCmd := runSnapshot(t, append(args, file.Name())...)

			testCmd.ExpectExit()
			if code := testCmd.ExitStatus(); code != 0 {
				t.Fatalf("command exit code %v, expected 0", code)
			}

			f, err := os.Open(file.Name())
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				err := f.Close()
				if err != nil {
					t.Error("closing snapshot file", "err", err)
				}
			}()

			b, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}
			var snap simulations.Snapshot
			err = json.Unmarshal(b, &snap)
			if err != nil {
				t.Fatal(err)
			}

			wantNodes := v.nodes
			if wantNodes == 0 {
				wantNodes = defaultNodes
			}
			gotNodes := len(snap.Nodes)
			if gotNodes != wantNodes {
				t.Errorf("got %v nodes, want %v", gotNodes, wantNodes)
			}

			if len(snap.Conns) == 0 {
				t.Error("no connections in a snapshot")
			}

			var wantServices []string
			if v.services != "" {
				wantServices = strings.Split(v.services, ",")
			} else {
				wantServices = []string{"bzz"}
			}
//对服务名进行排序，以便进行比较
//作为每个节点排序服务的字符串
			sort.Strings(wantServices)

			for i, n := range snap.Nodes {
				gotServices := n.Node.Config.Services
				sort.Strings(gotServices)
				if fmt.Sprint(gotServices) != fmt.Sprint(wantServices) {
					t.Errorf("got services %v for node %v, want %v", gotServices, i, wantServices)
				}
			}

		})
	}
}

