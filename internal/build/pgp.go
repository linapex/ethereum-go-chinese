
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450091691806720>


//signfile读取输入文件的内容并对其进行签名（装甲格式）
//在提供密钥的情况下，将签名放入输出文件。

package build

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/crypto/openpgp"
)

//pgpsignfile从指定的字符串分析pgp私钥并创建
//signature file into the output parameter of the input file.
//
//注意，此方法假定一个键将是pgpkey参数中的容器，
//此外，它是装甲格式。
func PGPSignFile(input string, output string, pgpkey string) error {
//解析密钥环并确保我们只有一个私钥。
	keys, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(pgpkey))
	if err != nil {
		return err
	}
	if len(keys) != 1 {
		return fmt.Errorf("key count mismatch: have %d, want %d", len(keys), 1)
	}
//创建用于签名的输入和输出流
	in, err := os.Open(input)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(output)
	if err != nil {
		return err
	}
	defer out.Close()

//生成签名并返回
	return openpgp.ArmoredDetachSign(out, keys[0], in, nil)
}

//PGPKeyID parses an armored key and returns the key ID.
func PGPKeyID(pgpkey string) (string, error) {
	keys, err := openpgp.ReadArmoredKeyRing(bytes.NewBufferString(pgpkey))
	if err != nil {
		return "", err
	}
	if len(keys) != 1 {
		return "", fmt.Errorf("key count mismatch: have %d, want %d", len(keys), 1)
	}
	return keys[0].PrimaryKey.KeyIdString(), nil
}

