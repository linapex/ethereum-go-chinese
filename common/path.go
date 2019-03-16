
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450074113478656>


package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//makename创建一个遵循以太坊约定的节点名
//对于这样的名字。它添加了操作系统名称和go运行时版本
//这个名字。
func MakeName(name, version string) string {
	return fmt.Sprintf("%s/v%s/%s/%s", name, version, runtime.GOOS, runtime.Version())
}

//fileexist检查文件路径中是否存在文件。
func FileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

//absolutePath返回datadir+filename，如果是绝对的，则返回filename。
func AbsolutePath(datadir string, filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return filepath.Join(datadir, filename)
}

