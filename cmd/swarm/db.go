
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450070623817728>


package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"gopkg.in/urfave/cli.v1"
)

var dbCommand = cli.Command{
	Name:               "db",
	CustomHelpTemplate: helpTemplate,
	Usage:              "manage the local chunk database",
	ArgsUsage:          "db COMMAND",
	Description:        "Manage the local chunk database",
	Subcommands: []cli.Command{
		{
			Action:             dbExport,
			CustomHelpTemplate: helpTemplate,
			Name:               "export",
			Usage:              "export a local chunk database as a tar archive (use - to send to stdout)",
			ArgsUsage:          "<chunkdb> <file>",
			Description: `
Export a local chunk database as a tar archive (use - to send to stdout).

    swarm db export ~/.ethereum/swarm/bzz-KEY/chunks chunks.tar

The export may be quite large, consider piping the output through the Unix
pv(1) tool to get a progress bar:

    swarm db export ~/.ethereum/swarm/bzz-KEY/chunks - | pv > chunks.tar
`,
		},
		{
			Action:             dbImport,
			CustomHelpTemplate: helpTemplate,
			Name:               "import",
			Usage:              "import chunks from a tar archive into a local chunk database (use - to read from stdin)",
			ArgsUsage:          "<chunkdb> <file>",
			Description: `Import chunks from a tar archive into a local chunk database (use - to read from stdin).

    swarm db import ~/.ethereum/swarm/bzz-KEY/chunks chunks.tar

The import may be quite large, consider piping the input through the Unix
pv(1) tool to get a progress bar:

    pv chunks.tar | swarm db import ~/.ethereum/swarm/bzz-KEY/chunks -`,
		},
	},
}

func dbExport(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) != 3 {
		utils.Fatalf("invalid arguments, please specify both <chunkdb> (path to a local chunk database), <file> (path to write the tar archive to, - for stdout) and the base key")
	}

	store, err := openLDBStore(args[0], common.Hex2Bytes(args[2]))
	if err != nil {
		utils.Fatalf("error opening local chunk database: %s", err)
	}
	defer store.Close()

	var out io.Writer
	if args[1] == "-" {
		out = os.Stdout
	} else {
		f, err := os.Create(args[1])
		if err != nil {
			utils.Fatalf("error opening output file: %s", err)
		}
		defer f.Close()
		out = f
	}

	count, err := store.Export(out)
	if err != nil {
		utils.Fatalf("error exporting local chunk database: %s", err)
	}

	log.Info(fmt.Sprintf("successfully exported %d chunks", count))
}

func dbImport(ctx *cli.Context) {
	args := ctx.Args()
	if len(args) != 3 {
		utils.Fatalf("invalid arguments, please specify both <chunkdb> (path to a local chunk database), <file> (path to read the tar archive from, - for stdin) and the base key")
	}

	store, err := openLDBStore(args[0], common.Hex2Bytes(args[2]))
	if err != nil {
		utils.Fatalf("error opening local chunk database: %s", err)
	}
	defer store.Close()

	var in io.Reader
	if args[1] == "-" {
		in = os.Stdin
	} else {
		f, err := os.Open(args[1])
		if err != nil {
			utils.Fatalf("error opening input file: %s", err)
		}
		defer f.Close()
		in = f
	}

	count, err := store.Import(in)
	if err != nil {
		utils.Fatalf("error importing local chunk database: %s", err)
	}

	log.Info(fmt.Sprintf("successfully imported %d chunks", count))
}

func openLDBStore(path string, basekey []byte) (*storage.LDBStore, error) {
	if _, err := os.Stat(filepath.Join(path, "CURRENT")); err != nil {
		return nil, fmt.Errorf("invalid chunkdb path: %s", err)
	}

	storeparams := storage.NewDefaultStoreParams()
	ldbparams := storage.NewLDBStoreParams(storeparams, path)
	ldbparams.BaseKey = basekey
	return storage.NewLDBStore(ldbparams)
}

