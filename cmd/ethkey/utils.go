
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450067687804928>


package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/crypto"
	"gopkg.in/urfave/cli.v1"
)

//promptppassphrase提示用户输入密码短语。将确认设置为真
//要求用户确认密码短语。
func promptPassphrase(confirmation bool) string {
	passphrase, err := console.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		utils.Fatalf("Failed to read passphrase: %v", err)
	}

	if confirmation {
		confirm, err := console.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			utils.Fatalf("Failed to read passphrase confirmation: %v", err)
		}
		if passphrase != confirm {
			utils.Fatalf("Passphrases do not match")
		}
	}

	return passphrase
}

//getpassphrase获取用户给定的密码。它首先检查
//--passfile命令行标志，并最终提示用户
//口令。
func getPassphrase(ctx *cli.Context) string {
//查找--passwordfile标志。
	passphraseFile := ctx.String(passphraseFlag.Name)
	if passphraseFile != "" {
		content, err := ioutil.ReadFile(passphraseFile)
		if err != nil {
			utils.Fatalf("Failed to read passphrase file '%s': %v",
				passphraseFile, err)
		}
		return strings.TrimRight(string(content), "\r\n")
	}

//否则提示用户输入密码短语。
	return promptPassphrase(false)
}

//signhash是一个帮助函数，用于计算给定消息的哈希值
//可以安全地用于计算签名。
//
//哈希被计算为
//keccak256（“\x19ethereum签名消息：\n”$消息长度$消息）。
//
//这将为已签名的消息提供上下文，并防止对事务进行签名。
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

//mustprintjson打印给定对象的json编码，以及
//在封送失败时退出程序并显示错误消息。
func mustPrintJSON(jsonObject interface{}) {
	str, err := json.MarshalIndent(jsonObject, "", "  ")
	if err != nil {
		utils.Fatalf("Failed to marshal JSON object: %v", err)
	}
	fmt.Println(string(str))
}

