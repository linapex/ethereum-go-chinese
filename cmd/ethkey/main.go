
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450067482284032>


package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"gopkg.in/urfave/cli.v1"
)

const (
	defaultKeyfileName = "keyfile.json"
)

//git sha1提交发布的哈希（通过链接器标志设置）
var gitCommit = ""

var app *cli.App

func init() {
	app = utils.NewApp(gitCommit, "an Ethereum key manager")
	app.Commands = []cli.Command{
		commandGenerate,
		commandInspect,
		commandChangePassphrase,
		commandSignMessage,
		commandVerifyMessage,
	}
}

//常用命令行标志。
var (
	passphraseFlag = cli.StringFlag{
		Name:  "passwordfile",
		Usage: "the file that contains the passphrase for the keyfile",
	}
	jsonFlag = cli.BoolFlag{
		Name:  "json",
		Usage: "output JSON instead of human-readable format",
	}
)

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

