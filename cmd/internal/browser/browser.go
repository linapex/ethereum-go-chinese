
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450068925124608>

//版权所有2016 Go作者。版权所有。
//此源代码的使用受BSD样式的控制
//可以在许可文件中找到的许可证。

//包浏览器提供了与用户浏览器交互的实用程序。
package browser

import (
	"os"
	"os/exec"
	"runtime"
)

//命令返回用于打开URL的可能命令的列表。
func Commands() [][]string {
	var cmds [][]string
	if exe := os.Getenv("BROWSER"); exe != "" {
		cmds = append(cmds, []string{exe})
	}
	switch runtime.GOOS {
	case "darwin":
		cmds = append(cmds, []string{"/usr/bin/open"})
	case "windows":
		cmds = append(cmds, []string{"cmd", "/c", "start"})
	default:
		cmds = append(cmds, []string{"xdg-open"})
	}
	cmds = append(cmds,
		[]string{"chrome"},
		[]string{"google-chrome"},
		[]string{"chromium"},
		[]string{"firefox"},
	)
	return cmds
}

//open尝试在浏览器中打开URL并报告是否成功。
func Open(url string) bool {
	for _, args := range Commands() {
		cmd := exec.Command(args[0], append(args[1:], url)...)
		if cmd.Start() == nil {
			return true
		}
	}
	return false
}

