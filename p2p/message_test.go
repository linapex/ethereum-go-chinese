
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450104878698496>


package p2p

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"
	"time"
)

func ExampleMsgPipe() {
	rw1, rw2 := MsgPipe()
	go func() {
		Send(rw1, 8, [][]byte{{0, 0}})
		Send(rw1, 5, [][]byte{{1, 1}})
		rw1.Close()
	}()

	for {
		msg, err := rw2.ReadMsg()
		if err != nil {
			break
		}
		var data [][]byte
		msg.Decode(&data)
		fmt.Printf("msg: %d, %x\n", msg.Code, data[0])
	}
//输出：
//味精：8, 0000
//味精：5, 0101
}

func TestMsgPipeUnblockWrite(t *testing.T) {
loop:
	for i := 0; i < 100; i++ {
		rw1, rw2 := MsgPipe()
		done := make(chan struct{})
		go func() {
			if err := SendItems(rw1, 1); err == nil {
				t.Error("EncodeMsg returned nil error")
			} else if err != ErrPipeClosed {
				t.Errorf("EncodeMsg returned wrong error: got %v, want %v", err, ErrPipeClosed)
			}
			close(done)
		}()

//此调用应确保EncodeMsg正在等待
//有时交付。如果不这样做，关闭很可能
//在encodemsg启动之前执行，然后我们将不测试
//所有病例。
		runtime.Gosched()

		rw2.Close()
		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
			t.Errorf("write didn't unblock")
			break loop
		}
	}
}

//如果未正确执行并发关闭操作，则此测试应该会死机。
func TestMsgPipeConcurrentClose(t *testing.T) {
	rw1, _ := MsgPipe()
	for i := 0; i < 10; i++ {
		go rw1.Close()
	}
}

func TestEOFSignal(t *testing.T) {
	rb := make([]byte, 10)

//空读者
	eof := make(chan struct{}, 1)
	sig := &eofSignal{new(bytes.Buffer), 0, eof}
	if n, err := sig.Read(rb); n != 0 || err != io.EOF {
		t.Errorf("Read returned unexpected values: (%v, %v)", n, err)
	}
	select {
	case <-eof:
	default:
		t.Error("EOF chan not signaled")
	}

//错误前计数
	eof = make(chan struct{}, 1)
	sig = &eofSignal{bytes.NewBufferString("aaaaaaaa"), 4, eof}
	if n, err := sig.Read(rb); n != 4 || err != nil {
		t.Errorf("Read returned unexpected values: (%v, %v)", n, err)
	}
	select {
	case <-eof:
	default:
		t.Error("EOF chan not signaled")
	}

//计数前出错
	eof = make(chan struct{}, 1)
	sig = &eofSignal{bytes.NewBufferString("aaaa"), 999, eof}
	if n, err := sig.Read(rb); n != 4 || err != nil {
		t.Errorf("Read returned unexpected values: (%v, %v)", n, err)
	}
	if n, err := sig.Read(rb); n != 0 || err != io.EOF {
		t.Errorf("Read returned unexpected values: (%v, %v)", n, err)
	}
	select {
	case <-eof:
	default:
		t.Error("EOF chan not signaled")
	}

//如果两者都不发生，则无信号
	eof = make(chan struct{}, 1)
	sig = &eofSignal{bytes.NewBufferString("aaaaaaaaaaaaaaaaaaaaa"), 999, eof}
	if n, err := sig.Read(rb); n != 10 || err != nil {
		t.Errorf("Read returned unexpected values: (%v, %v)", n, err)
	}
	select {
	case <-eof:
		t.Error("unexpected EOF signal")
	default:
	}
}

func unhex(str string) []byte {
	r := strings.NewReplacer("\t", "", " ", "", "\n", "")
	b, err := hex.DecodeString(r.Replace(str))
	if err != nil {
		panic(fmt.Sprintf("invalid hex string: %q", str))
	}
	return b
}

