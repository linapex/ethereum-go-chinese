
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450074239307776>


package common

import (
	"fmt"
)

//StorageSize是一个围绕浮点值的包装器，它支持用户友好的
//格式化。
type StorageSize float64

//字符串实现字符串接口。
func (s StorageSize) String() string {
	if s > 1000000 {
		return fmt.Sprintf("%.2f mB", s/1000000)
	} else if s > 1000 {
		return fmt.Sprintf("%.2f kB", s/1000)
	} else {
		return fmt.Sprintf("%.2f B", s)
	}
}

//terminalString实现log.terminalStringer，为控制台格式化字符串
//日志记录期间的输出。
func (s StorageSize) TerminalString() string {
	if s > 1000000 {
		return fmt.Sprintf("%.2fmB", s/1000000)
	} else if s > 1000 {
		return fmt.Sprintf("%.2fkB", s/1000)
	} else {
		return fmt.Sprintf("%.2fB", s)
	}
}

