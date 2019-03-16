
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450071928246272>


package main

import (
	"os"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	cli "gopkg.in/urfave/cli.v1"
)

var gitCommit string //git sha1提交发布的哈希（通过链接器标志设置）

//“创建”命令的默认值--节点标志
const defaultNodes = 10

func main() {
	err := newApp().Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

//newapp构建了一个新的swarm快照实用程序实例。
//方法run在主函数和测试中对其进行调用。
func newApp() (app *cli.App) {
	app = utils.NewApp(gitCommit, "Swarm Snapshot Utility")

	app.Name = "swarm-snapshot"
	app.Usage = ""

//应用程序标志（用于所有命令）
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "verbosity",
			Value: 1,
			Usage: "verbosity level",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create a swarm snapshot",
			Action:  create,
//仅用于“创建”命令的标志。
//允许在
//命令参数。
			Flags: append(app.Flags,
				cli.IntFlag{
					Name:  "nodes",
					Value: defaultNodes,
					Usage: "number of nodes",
				},
				cli.StringFlag{
					Name:  "services",
					Value: "bzz",
					Usage: "comma separated list of services to boot the nodes with",
				},
			),
		},
	}

	return app
}

