
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450067377426432>


package main

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pborman/uuid"
	"gopkg.in/urfave/cli.v1"
)

type outputGenerate struct {
	Address      string
	AddressEIP55 string
}

var commandGenerate = cli.Command{
	Name:      "generate",
	Usage:     "generate new keyfile",
	ArgsUsage: "[ <keyfile> ]",
	Description: `
Generate a new keyfile.

If you want to encrypt an existing private key, it can be specified by setting
--privatekey with the location of the file containing the private key.
`,
	Flags: []cli.Flag{
		passphraseFlag,
		jsonFlag,
		cli.StringFlag{
			Name:  "privatekey",
			Usage: "file containing a raw private key to encrypt",
		},
	},
	Action: func(ctx *cli.Context) error {
//检查是否指定了关键文件路径，并确保它不存在。
		keyfilepath := ctx.Args().First()
		if keyfilepath == "" {
			keyfilepath = defaultKeyfileName
		}
		if _, err := os.Stat(keyfilepath); err == nil {
			utils.Fatalf("Keyfile already exists at %s.", keyfilepath)
		} else if !os.IsNotExist(err) {
			utils.Fatalf("Error checking if keyfile exists: %v", err)
		}

		var privateKey *ecdsa.PrivateKey
		var err error
		if file := ctx.String("privatekey"); file != "" {
//从文件加载私钥。
			privateKey, err = crypto.LoadECDSA(file)
			if err != nil {
				utils.Fatalf("Can't load private key: %v", err)
			}
		} else {
//如果未加载，则生成随机。
			privateKey, err = crypto.GenerateKey()
			if err != nil {
				utils.Fatalf("Failed to generate random private key: %v", err)
			}
		}

//使用随机UUID创建keyfile对象。
		id := uuid.NewRandom()
		key := &keystore.Key{
			Id:         id,
			Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
			PrivateKey: privateKey,
		}

//用密码短语加密密钥。
		passphrase := promptPassphrase(true)
		keyjson, err := keystore.EncryptKey(key, passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
		if err != nil {
			utils.Fatalf("Error encrypting key: %v", err)
		}

//将文件存储到磁盘。
		if err := os.MkdirAll(filepath.Dir(keyfilepath), 0700); err != nil {
			utils.Fatalf("Could not create directory %s", filepath.Dir(keyfilepath))
		}
		if err := ioutil.WriteFile(keyfilepath, keyjson, 0600); err != nil {
			utils.Fatalf("Failed to write keyfile to %s: %v", keyfilepath, err)
		}

//输出一些信息。
		out := outputGenerate{
			Address: key.Address.Hex(),
		}
		if ctx.Bool(jsonFlag.Name) {
			mustPrintJSON(out)
		} else {
			fmt.Println("Address:", out.Address)
		}
		return nil
	},
}

