
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450103825928192>


//+构建Go1.4、Nacl、Faketime_模拟

package discv5

import (
	"os"
	"runtime"
	"testing"
	"unsafe"
)

//在运行时启用假时间模式，如在游乐场上。
//这有一点可能不起作用，因为有些go代码
//可能在设置变量之前执行。

//转到：linkname faketime runtime.faketime
var faketime = 1

func TestMain(m *testing.M) {
//为了获得go:linkname的访问权限，我们需要以某种方式使用unsafe。
	_ = unsafe.Sizeof(0)

//运行实际测试。runwithplaygroundtime确保
//这就是所谓的跑步。
	runtime.GOMAXPROCS(8)
	os.Exit(m.Run())
}

