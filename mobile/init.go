
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450101485506560>


//包含mbile库的初始化代码。

package geth

import (
	"os"
	"runtime"

	"github.com/ethereum/go-ethereum/log"
)

func init() {
//初始化记录器
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlInfo, log.StreamHandler(os.Stderr, log.TerminalFormat(false))))

//初始化goroutine计数
	runtime.GOMAXPROCS(runtime.NumCPU())
}

