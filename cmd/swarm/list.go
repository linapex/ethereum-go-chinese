
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450071060025344>


package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/ethereum/go-ethereum/cmd/utils"
	swarm "github.com/ethereum/go-ethereum/swarm/api/client"
	"gopkg.in/urfave/cli.v1"
)

var listCommand = cli.Command{
	Action:             list,
	CustomHelpTemplate: helpTemplate,
	Name:               "ls",
	Usage:              "list files and directories contained in a manifest",
	ArgsUsage:          "<manifest> [<prefix>]",
	Description:        "Lists files and directories contained in a manifest",
}

func list(ctx *cli.Context) {
	args := ctx.Args()

	if len(args) < 1 {
		utils.Fatalf("Please supply a manifest reference as the first argument")
	} else if len(args) > 2 {
		utils.Fatalf("Too many arguments - usage 'swarm ls manifest [prefix]'")
	}
	manifest := args[0]

	var prefix string
	if len(args) == 2 {
		prefix = args[1]
	}

	bzzapi := strings.TrimRight(ctx.GlobalString(SwarmApiFlag.Name), "/")
	client := swarm.NewClient(bzzapi)
	list, err := client.List(manifest, prefix, "")
	if err != nil {
		utils.Fatalf("Failed to generate file and directory list: %s", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 2, 2, ' ', 0)
	defer w.Flush()
	fmt.Fprintln(w, "HASH\tCONTENT TYPE\tPATH")
	for _, prefix := range list.CommonPrefixes {
		fmt.Fprintf(w, "%s\t%s\t%s\n", "", "DIR", prefix)
	}
	for _, entry := range list.Entries {
		fmt.Fprintf(w, "%s\t%s\t%s\n", entry.Hash, entry.ContentType, entry.Path)
	}
}

