
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450067436146688>


package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/urfave/cli.v1"
)

type outputInspect struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

var commandInspect = cli.Command{
	Name:      "inspect",
	Usage:     "inspect a keyfile",
	ArgsUsage: "<keyfile>",
	Description: `
Print various information about the keyfile.

Private key information can be printed by using the --private flag;
make sure to use this feature with great caution!`,
	Flags: []cli.Flag{
		passphraseFlag,
		jsonFlag,
		cli.BoolFlag{
			Name:  "private",
			Usage: "include the private key in the output",
		},
	},
	Action: func(ctx *cli.Context) error {
		keyfilepath := ctx.Args().First()

//从文件中读取密钥。
		keyjson, err := ioutil.ReadFile(keyfilepath)
		if err != nil {
			utils.Fatalf("Failed to read the keyfile at '%s': %v", keyfilepath, err)
		}

//用密码短语解密密钥。
		passphrase := getPassphrase(ctx)
		key, err := keystore.DecryptKey(keyjson, passphrase)
		if err != nil {
			utils.Fatalf("Error decrypting key: %v", err)
		}

//输出我们能检索到的所有相关信息。
		showPrivate := ctx.Bool("private")
		out := outputInspect{
			Address: key.Address.Hex(),
			PublicKey: hex.EncodeToString(
				crypto.FromECDSAPub(&key.PrivateKey.PublicKey)),
		}
		if showPrivate {
			out.PrivateKey = hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))
		}

		if ctx.Bool(jsonFlag.Name) {
			mustPrintJSON(out)
		} else {
			fmt.Println("Address:       ", out.Address)
			fmt.Println("Public key:    ", out.PublicKey)
			if showPrivate {
				fmt.Println("Private key:   ", out.PrivateKey)
			}
		}
		return nil
	},
}

