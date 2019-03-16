
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450068493111296>


package main

import (
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/params"
)

const (
	ipcAPIs  = "admin:1.0 debug:1.0 eth:1.0 ethash:1.0 miner:1.0 net:1.0 personal:1.0 rpc:1.0 shh:1.0 txpool:1.0 web3:1.0"
	httpAPIs = "eth:1.0 net:1.0 rpc:1.0 web3:1.0"
)

//测试控制台中嵌入的节点是否可以正确启动，以及
//然后通过关闭输入流终止。
func TestConsoleWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"

//启动一个geth控制台，确保它已清理干净并终止控制台
	geth := runGeth(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--etherbase", coinbase, "--shh",
		"console")

//收集欢迎信息需要包含的所有信息
	geth.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	geth.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	geth.SetTemplateFunc("gover", runtime.Version)
	geth.SetTemplateFunc("gethver", func() string { return params.VersionWithMeta })
	geth.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	geth.SetTemplateFunc("apis", func() string { return ipcAPIs })

//验证所需模板的实际欢迎消息
	geth.Expect(`
Welcome to the Geth JavaScript console!

instance: Geth/v{{gethver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{.Etherbase}}
at block: 0 ({{niltime}})
 datadir: {{.Datadir}}
 modules: {{apis}}

> {{.InputLine "exit"}}
`)
	geth.ExpectExit()
}

//测试控制台是否可以通过各种方式连接到正在运行的节点。
func TestIPCAttachWelcome(t *testing.T) {
//为IPC附件配置实例
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	var ipc string
	if runtime.GOOS == "windows" {
		ipc = `\\.\pipe\geth` + strconv.Itoa(trulyRandInt(100000, 999999))
	} else {
		ws := tmpdir(t)
		defer os.RemoveAll(ws)
		ipc = filepath.Join(ws, "geth.ipc")
	}
//注意：我们需要--shh，因为tesattachwelcome检查默认值
//其中包括IPC模块和SHH列表。
	geth := runGeth(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--etherbase", coinbase, "--shh", "--ipcpath", ipc)

time.Sleep(2 * time.Second) //等待RPC端点打开的简单方法
	testAttachWelcome(t, geth, "ipc:"+ipc, ipcAPIs)

	geth.Interrupt()
	geth.ExpectExit()
}

func TestHTTPAttachWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
port := strconv.Itoa(trulyRandInt(1024, 65536)) //是的，有时会失败，对不起：P
	geth := runGeth(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--etherbase", coinbase, "--rpc", "--rpcport", port)

time.Sleep(2 * time.Second) //等待RPC端点打开的简单方法
testAttachWelcome(t, geth, "http://本地主机：“+端口，httpapis）

	geth.Interrupt()
	geth.ExpectExit()
}

func TestWSAttachWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
port := strconv.Itoa(trulyRandInt(1024, 65536)) //是的，有时会失败，对不起：P

	geth := runGeth(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--etherbase", coinbase, "--ws", "--wsport", port)

time.Sleep(2 * time.Second) //等待RPC端点打开的简单方法
testAttachWelcome(t, geth, "ws://本地主机：“+端口，httpapis）

	geth.Interrupt()
	geth.ExpectExit()
}

func testAttachWelcome(t *testing.T, geth *testgeth, endpoint, apis string) {
//附加到正在运行的geth note并立即终止
	attach := runGeth(t, "attach", endpoint)
	defer attach.ExpectExit()
	attach.CloseStdin()

//收集欢迎信息需要包含的所有信息
	attach.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	attach.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	attach.SetTemplateFunc("gover", runtime.Version)
	attach.SetTemplateFunc("gethver", func() string { return params.VersionWithMeta })
	attach.SetTemplateFunc("etherbase", func() string { return geth.Etherbase })
	attach.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	attach.SetTemplateFunc("ipc", func() bool { return strings.HasPrefix(endpoint, "ipc") })
	attach.SetTemplateFunc("datadir", func() string { return geth.Datadir })
	attach.SetTemplateFunc("apis", func() string { return apis })

//验证所需模板的实际欢迎消息
	attach.Expect(`
Welcome to the Geth JavaScript console!

instance: Geth/v{{gethver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{etherbase}}
at block: 0 ({{niltime}}){{if ipc}}
 datadir: {{datadir}}{{end}}
 modules: {{apis}}

> {{.InputLine "exit" }}
`)
	attach.ExpectExit()
}

//Trulyrandit生成控制台测试使用的加密随机整数
//不会与并行运行的其他测试冲突网络端口。
func trulyRandInt(lo, hi int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(hi-lo)))
	return int(num.Int64()) + lo
}

