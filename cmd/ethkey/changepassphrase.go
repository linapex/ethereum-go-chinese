
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450067331289088>

package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"gopkg.in/urfave/cli.v1"
)

var newPassphraseFlag = cli.StringFlag{
	Name:  "newpasswordfile",
	Usage: "the file that contains the new passphrase for the keyfile",
}

var commandChangePassphrase = cli.Command{
	Name:      "changepassphrase",
	Usage:     "change the passphrase on a keyfile",
	ArgsUsage: "<keyfile>",
	Description: `
Change the passphrase of a keyfile.`,
	Flags: []cli.Flag{
		passphraseFlag,
		newPassphraseFlag,
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

//获取新密码。
		fmt.Println("Please provide a new passphrase")
		var newPhrase string
		if passFile := ctx.String(newPassphraseFlag.Name); passFile != "" {
			content, err := ioutil.ReadFile(passFile)
			if err != nil {
				utils.Fatalf("Failed to read new passphrase file '%s': %v", passFile, err)
			}
			newPhrase = strings.TrimRight(string(content), "\r\n")
		} else {
			newPhrase = promptPassphrase(true)
		}

//用新密码短语加密密钥。
		newJson, err := keystore.EncryptKey(key, newPhrase, keystore.StandardScryptN, keystore.StandardScryptP)
		if err != nil {
			utils.Fatalf("Error encrypting with new passphrase: %v", err)
		}

//然后写新的关键文件代替旧的关键文件。
		if err := ioutil.WriteFile(keyfilepath, newJson, 600); err != nil {
			utils.Fatalf("Error writing new keyfile to disk: %v", err)
		}

//不要打印任何内容。成功返回，
//生成一个正的退出代码。
		return nil
	},
}

