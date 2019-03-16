
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450101581975552>


package geth

import (
	"os"

	"github.com/ethereum/go-ethereum/log"
)

//setverbosity设置全局详细级别（介于0和6之间-请参阅logger/verbosity.go）。
func SetVerbosity(level int) {
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(level), log.StreamHandler(os.Stderr, log.TerminalFormat(false))))
}

