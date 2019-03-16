
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450121760772096>


package testutil

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"testing"
)

//tempfilewithcontent是一个助手函数，它创建一个包含以下字符串内容的临时文件，然后关闭文件句柄
//它返回完整的文件路径
func TempFileWithContent(t *testing.T, content string) string {
	tempFile, err := ioutil.TempFile("", "swarm-temp-file")
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(tempFile, strings.NewReader(content))
	if err != nil {
		os.RemoveAll(tempFile.Name())
		t.Fatal(err)
	}
	if err = tempFile.Close(); err != nil {
		t.Fatal(err)
	}
	return tempFile.Name()
}

//RandomBytes返回伪随机确定性结果
//因为测试失败必须是可复制的
func RandomBytes(seed, length int) []byte {
	b := make([]byte, length)
	reader := rand.New(rand.NewSource(int64(seed)))
	for n := 0; n < length; {
		read, err := reader.Read(b[n:])
		if err != nil {
			panic(err)
		}
		n += read
	}
	return b
}

func RandomReader(seed, length int) *bytes.Reader {
	return bytes.NewReader(RandomBytes(seed, length))
}

