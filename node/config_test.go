
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450102127235072>


package node

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p"
)

//可以成功创建datadir的测试，是否可以手动配置
//或自动生成的临时文件。
func TestDatadirCreation(t *testing.T) {
//创建一个临时数据目录，并检查它是否可由节点使用
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create manual data dir: %v", err)
	}
	defer os.RemoveAll(dir)

	if _, err := New(&Config{DataDir: dir}); err != nil {
		t.Fatalf("failed to create stack with existing datadir: %v", err)
	}
//生成一个长的不存在的datadir路径，并检查它是否由节点创建。
	dir = filepath.Join(dir, "a", "b", "c", "d", "e", "f")
	if _, err := New(&Config{DataDir: dir}); err != nil {
		t.Fatalf("failed to create stack with creatable datadir: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("freshly created datadir not accessible: %v", err)
	}
//验证不可能的datadir创建失败
	file, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	dir = filepath.Join(file.Name(), "invalid/path")
	if _, err := New(&Config{DataDir: dir}); err == nil {
		t.Fatalf("protocol stack created with an invalid datadir")
	}
}

//测试IPC路径是否正确解析为不同的有效终结点
//平台。
func TestIPCPathResolution(t *testing.T) {
	var tests = []struct {
		DataDir  string
		IPCPath  string
		Windows  bool
		Endpoint string
	}{
		{"", "", false, ""},
		{"data", "", false, ""},
		{"", "geth.ipc", false, filepath.Join(os.TempDir(), "geth.ipc")},
		{"data", "geth.ipc", false, "data/geth.ipc"},
		{"data", "./geth.ipc", false, "./geth.ipc"},
		{"data", "/geth.ipc", false, "/geth.ipc"},
		{"", "", true, ``},
		{"data", "", true, ``},
		{"", "geth.ipc", true, `\\.\pipe\geth.ipc`},
		{"data", "geth.ipc", true, `\\.\pipe\geth.ipc`},
		{"data", `\\.\pipe\geth.ipc`, true, `\\.\pipe\geth.ipc`},
	}
	for i, test := range tests {
//仅在平台/测试匹配时运行
		if (runtime.GOOS == "windows") == test.Windows {
			if endpoint := (&Config{DataDir: test.DataDir, IPCPath: test.IPCPath}).IPCEndpoint(); endpoint != test.Endpoint {
				t.Errorf("test %d: IPC endpoint mismatch: have %s, want %s", i, endpoint, test.Endpoint)
			}
		}
	}
}

//测试节点键是否可以正确创建、持久化、加载和/或生成
//短暂的
func TestNodeKeyPersistency(t *testing.T) {
//创建临时文件夹并确保没有密钥
	dir, err := ioutil.TempDir("", "node-test")
	if err != nil {
		t.Fatalf("failed to create temporary data directory: %v", err)
	}
	defer os.RemoveAll(dir)

	keyfile := filepath.Join(dir, "unit-test", datadirPrivateKey)

//使用预设键配置节点，并确保它不会持久化。
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate one-shot node key: %v", err)
	}
	config := &Config{Name: "unit-test", DataDir: dir, P2P: p2p.Config{PrivateKey: key}}
	config.NodeKey()
	if _, err := os.Stat(filepath.Join(keyfile)); err == nil {
		t.Fatalf("one-shot node key persisted to data directory")
	}

//配置一个没有预设键的节点，并确保这一次将其持久化。
	config = &Config{Name: "unit-test", DataDir: dir}
	config.NodeKey()
	if _, err := os.Stat(keyfile); err != nil {
		t.Fatalf("node key not persisted to data directory: %v", err)
	}
	if _, err = crypto.LoadECDSA(keyfile); err != nil {
		t.Fatalf("failed to load freshly persisted node key: %v", err)
	}
	blob1, err := ioutil.ReadFile(keyfile)
	if err != nil {
		t.Fatalf("failed to read freshly persisted node key: %v", err)
	}

//配置新节点并确保加载以前保存的密钥
	config = &Config{Name: "unit-test", DataDir: dir}
	config.NodeKey()
	blob2, err := ioutil.ReadFile(filepath.Join(keyfile))
	if err != nil {
		t.Fatalf("failed to read previously persisted node key: %v", err)
	}
	if !bytes.Equal(blob1, blob2) {
		t.Fatalf("persisted node key mismatch: have %x, want %x", blob2, blob1)
	}

//配置临时节点并确保没有密钥在本地转储
	config = &Config{Name: "unit-test", DataDir: ""}
	config.NodeKey()
	if _, err := os.Stat(filepath.Join(".", "unit-test", datadirPrivateKey)); err == nil {
		t.Fatalf("ephemeral node key persisted to disk")
	}
}

