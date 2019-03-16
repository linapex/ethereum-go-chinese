
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450073274617856>


package common

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

//报告发出警告，要求用户向Github跟踪器提交问题。
func Report(extra ...interface{}) {
fmt.Fprintln(os.Stderr, "You've encountered a sought after, hard to reproduce bug. Please report this to the developers <3 https://github.com/ethereum/go-ethereum/issues“）
	fmt.Fprintln(os.Stderr, extra...)

	_, file, line, _ := runtime.Caller(1)
	fmt.Fprintf(os.Stderr, "%v:%v\n", file, line)

	debug.PrintStack()

	fmt.Fprintln(os.Stderr, "#### BUG! PLEASE REPORT ####")
}

//使用fmt.println在框中打印出给定的字符串。
func PrintDepricationWarning(str string) {
	line := strings.Repeat("#", len(str)+4)
	emptyLine := strings.Repeat(" ", len(str))
	fmt.Printf(`
%s
# %s #
# %s #
# %s #
%s

`, line, emptyLine, str, emptyLine, line)
}

