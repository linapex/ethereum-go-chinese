
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450070669955072>

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/swarm/api"
	swarm "github.com/ethereum/go-ethereum/swarm/api/client"
	"gopkg.in/urfave/cli.v1"
)

var downloadCommand = cli.Command{
	Action:      download,
	Name:        "down",
	Flags:       []cli.Flag{SwarmRecursiveFlag, SwarmAccessPasswordFlag},
	Usage:       "downloads a swarm manifest or a file inside a manifest",
	ArgsUsage:   " <uri> [<dir>]",
	Description: `Downloads a swarm bzz uri to the given dir. When no dir is provided, working directory is assumed. --recursive flag is expected when downloading a manifest with multiple entries.`,
}

func download(ctx *cli.Context) {
	log.Debug("downloading content using swarm down")
	args := ctx.Args()
	dest := "."

	switch len(args) {
	case 0:
		utils.Fatalf("Usage: swarm down [options] <bzz locator> [<destination path>]")
	case 1:
		log.Trace(fmt.Sprintf("swarm down: no destination path - assuming working dir"))
	default:
		log.Trace(fmt.Sprintf("destination path arg: %s", args[1]))
		if absDest, err := filepath.Abs(args[1]); err == nil {
			dest = absDest
		} else {
			utils.Fatalf("could not get download path: %v", err)
		}
	}

	var (
		bzzapi      = strings.TrimRight(ctx.GlobalString(SwarmApiFlag.Name), "/")
		isRecursive = ctx.Bool(SwarmRecursiveFlag.Name)
		client      = swarm.NewClient(bzzapi)
	)

	if fi, err := os.Stat(dest); err == nil {
		if isRecursive && !fi.Mode().IsDir() {
			utils.Fatalf("destination path is not a directory!")
		}
	} else {
		if !os.IsNotExist(err) {
			utils.Fatalf("could not stat path: %v", err)
		}
	}

	uri, err := api.Parse(args[0])
	if err != nil {
		utils.Fatalf("could not parse uri argument: %v", err)
	}

	dl := func(credentials string) error {
//
		if isRecursive {
			if err := client.DownloadDirectory(uri.Addr, uri.Path, dest, credentials); err != nil {
				if err == swarm.ErrUnauthorized {
					return err
				}
				return fmt.Errorf("directory %s: %v", uri.Path, err)
			}
		} else {
//
			log.Debug("downloading file/path from a manifest", "uri.Addr", uri.Addr, "uri.Path", uri.Path)

			err := client.DownloadFile(uri.Addr, uri.Path, dest, credentials)
			if err != nil {
				if err == swarm.ErrUnauthorized {
					return err
				}
				return fmt.Errorf("file %s from address: %s: %v", uri.Path, uri.Addr, err)
			}
		}
		return nil
	}
	if passwords := makePasswordList(ctx); passwords != nil {
		password := getPassPhrase(fmt.Sprintf("Downloading %s is restricted", uri), 0, passwords)
		err = dl(password)
	} else {
		err = dl("")
	}
	if err != nil {
		utils.Fatalf("download: %v", err)
	}
}

