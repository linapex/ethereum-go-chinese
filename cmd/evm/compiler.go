
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450067738136576>


package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/cmd/evm/internal/compiler"

	cli "gopkg.in/urfave/cli.v1"
)

var compileCommand = cli.Command{
	Action:    compileCmd,
	Name:      "compile",
	Usage:     "compiles easm source to evm binary",
	ArgsUsage: "<file>",
}

func compileCmd(ctx *cli.Context) error {
	debug := ctx.GlobalBool(DebugFlag.Name)

	if len(ctx.Args().First()) == 0 {
		return errors.New("filename required")
	}

	fn := ctx.Args().First()
	src, err := ioutil.ReadFile(fn)
	if err != nil {
		return err
	}

	bin, err := compiler.Compile(fn, src, debug)
	if err != nil {
		return err
	}
	fmt.Println(bin)
	return nil
}

