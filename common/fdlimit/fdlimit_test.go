
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450073375281152>


package fdlimit

import (
	"fmt"
	"testing"
)

//testfiledescriptor限制只测试文件描述符是否允许
//根据此过程可以检索。
func TestFileDescriptorLimits(t *testing.T) {
	target := 4096
	hardlimit, err := Maximum()
	if err != nil {
		t.Fatal(err)
	}
	if hardlimit < target {
		t.Skip(fmt.Sprintf("system limit is less than desired test target: %d < %d", hardlimit, target))
	}

	if limit, err := Current(); err != nil || limit <= 0 {
		t.Fatalf("failed to retrieve file descriptor limit (%d): %v", limit, err)
	}
	if err := Raise(uint64(target)); err != nil {
		t.Fatalf("failed to raise file allowance")
	}
	if limit, err := Current(); err != nil || limit < target {
		t.Fatalf("failed to retrieve raised descriptor limit (have %v, want %v): %v", limit, target, err)
	}
}

