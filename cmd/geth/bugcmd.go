
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450068287590400>


package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/cmd/internal/browser"
	"github.com/ethereum/go-ethereum/params"

	"github.com/ethereum/go-ethereum/cmd/utils"
	cli "gopkg.in/urfave/cli.v1"
)

var bugCommand = cli.Command{
	Action:    utils.MigrateFlags(reportBug),
	Name:      "bug",
	Usage:     "opens a window to report a bug on the geth repo",
	ArgsUsage: " ",
	Category:  "MISCELLANEOUS COMMANDS",
}

const issueURL = "https://github.com/ethereum/go-ethereum/issues/new“

//ReportBug通过打开Go-Ethereum GH问题的新URL来报告Bug
//跟踪并将默认值设置为问题主体。
func reportBug(ctx *cli.Context) error {
//执行模板并将内容写入buff
	var buff bytes.Buffer

	fmt.Fprintln(&buff, "#### System information")
	fmt.Fprintln(&buff)
	fmt.Fprintln(&buff, "Version:", params.VersionWithMeta)
	fmt.Fprintln(&buff, "Go Version:", runtime.Version())
	fmt.Fprintln(&buff, "OS:", runtime.GOOS)
	printOSDetails(&buff)
	fmt.Fprintln(&buff, header)

//打开新的GH问题
	if !browser.Open(issueURL + "?body=" + url.QueryEscape(buff.String())) {
		fmt.Printf("Please file a new issue at %s using this template:\n\n%s", issueURL, buff.String())
	}
	return nil
}

//从Go源复制。版权所有2017 The Go作者
func printOSDetails(w io.Writer) {
	switch runtime.GOOS {
	case "darwin":
		printCmdOut(w, "uname -v: ", "uname", "-v")
		printCmdOut(w, "", "sw_vers")
	case "linux":
		printCmdOut(w, "uname -sr: ", "uname", "-sr")
		printCmdOut(w, "", "lsb_release", "-a")
	case "openbsd", "netbsd", "freebsd", "dragonfly":
		printCmdOut(w, "uname -v: ", "uname", "-v")
	case "solaris":
		out, err := ioutil.ReadFile("/etc/release")
		if err == nil {
			fmt.Fprintf(w, "/etc/release: %s\n", out)
		} else {
			fmt.Printf("failed to read /etc/release: %v\n", err)
		}
	}
}

//printcmdout打印运行给定命令的输出。
//它忽略失败；“去错误”是最好的努力。
//
//从Go源复制。版权所有2017 The Go作者
func printCmdOut(w io.Writer, prefix, path string, args ...string) {
	cmd := exec.Command(path, args...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s %s: %v\n", path, strings.Join(args, " "), err)
		return
	}
	fmt.Fprintf(w, "%s%s\n", prefix, bytes.TrimSpace(out))
}

const header = `
#### Expected behaviour


#### Actual behaviour


#### Steps to reproduce the behaviour


#### Backtrace
`

