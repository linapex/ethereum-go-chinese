
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450103729459200>


package discv5

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func getnacl() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		_, err := exec.LookPath("sel_ldr_x86_64")
		return "amd64p32", err
	case "i386":
		_, err := exec.LookPath("sel_ldr_i386")
		return "i386", err
	default:
		return "", errors.New("nacl is not supported on " + runtime.GOARCH)
	}
}

//runwithplaygroundtime执行调用方
//在启用faketime的nacl沙盒中。
//
//必须从test*函数调用此函数
//当ishost为true时，调用方必须跳过实际测试。
func runWithPlaygroundTime(t *testing.T) (isHost bool) {
	if runtime.GOOS == "nacl" {
		return false
	}

//打电话给对方。
	callerPC, _, _, ok := runtime.Caller(1)
	if !ok {
		panic("can't get caller")
	}
	callerFunc := runtime.FuncForPC(callerPC)
	if callerFunc == nil {
		panic("can't get caller")
	}
	callerName := callerFunc.Name()[strings.LastIndexByte(callerFunc.Name(), '.')+1:]
	if !strings.HasPrefix(callerName, "Test") {
		panic("must be called from witin a Test* function")
	}
	testPattern := "^" + callerName + "$"

//不幸的是，runtime.faketime（操场时间模式）只在nacl上工作。氯化钠
//必须安装sdk并将其链接到路径才能使其工作。
	arch, err := getnacl()
	if err != nil {
		t.Skip(err)
	}

//使用nacl编译并运行调用测试。
//额外的标签确保使用了sim_main_test.go中的test main功能。
	cmd := exec.Command("go", "test", "-v", "-tags", "faketime_simulation", "-timeout", "100h", "-run", testPattern, ".")
	cmd.Env = append([]string{"GOOS=nacl", "GOARCH=" + arch}, os.Environ()...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	go skipPlaygroundOutputHeaders(os.Stdout, stdout)
	go skipPlaygroundOutputHeaders(os.Stderr, stderr)
	if err := cmd.Run(); err != nil {
		t.Error(err)
	}

//确保测试功能不会在（非Nacl）主机进程中运行。
	return true
}

func skipPlaygroundOutputHeaders(out io.Writer, in io.Reader) {
//附加输出可以不打印标题
//在Nacl二进制文件开始运行之前（例如编译器错误消息）。
	bufin := bufio.NewReader(in)
	output, err := bufin.ReadBytes(0)
	output = bytes.TrimSuffix(output, []byte{0})
	if len(output) > 0 {
		out.Write(output)
	}
	if err != nil {
		return
	}
	bufin.UnreadByte()

//回放头：0 0 p b<8字节时间><4字节数据长度>
	head := make([]byte, 4+8+4)
	for {
		if _, err := io.ReadFull(bufin, head); err != nil {
			if err != io.EOF {
				fmt.Fprintln(out, "read error:", err)
			}
			return
		}
		if !bytes.HasPrefix(head, []byte{0x00, 0x00, 'P', 'B'}) {
			fmt.Fprintf(out, "expected playback header, got %q\n", head)
			io.Copy(out, bufin)
			return
		}
//将数据复制到下一个标题。
		size := binary.BigEndian.Uint32(head[12:])
		io.CopyN(out, bufin, int64(size))
	}
}

