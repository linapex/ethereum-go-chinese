
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450106636111872>


package adapters

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/p2p/simulations/pipes"
)

func TestTCPPipe(t *testing.T) {
	c1, c2, err := pipes.TCPPipe()
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		msgs := 50
		size := 1024
		for i := 0; i < msgs; i++ {
			msg := make([]byte, size)
			_ = binary.PutUvarint(msg, uint64(i))

			_, err := c1.Write(msg)
			if err != nil {
				t.Fatal(err)
			}
		}

		for i := 0; i < msgs; i++ {
			msg := make([]byte, size)
			_ = binary.PutUvarint(msg, uint64(i))

			out := make([]byte, size)
			_, err := c2.Read(out)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(msg, out) {
				t.Fatalf("expected %#v, got %#v", msg, out)
			}
		}
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("test timeout")
	}
}

func TestTCPPipeBidirections(t *testing.T) {
	c1, c2, err := pipes.TCPPipe()
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		msgs := 50
		size := 7
		for i := 0; i < msgs; i++ {
			msg := []byte(fmt.Sprintf("ping %02d", i))

			_, err := c1.Write(msg)
			if err != nil {
				t.Fatal(err)
			}
		}

		for i := 0; i < msgs; i++ {
			expected := []byte(fmt.Sprintf("ping %02d", i))

			out := make([]byte, size)
			_, err := c2.Read(out)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(expected, out) {
				t.Fatalf("expected %#v, got %#v", out, expected)
			} else {
				msg := []byte(fmt.Sprintf("pong %02d", i))
				_, err := c2.Write(msg)
				if err != nil {
					t.Fatal(err)
				}
			}
		}

		for i := 0; i < msgs; i++ {
			expected := []byte(fmt.Sprintf("pong %02d", i))

			out := make([]byte, size)
			_, err := c1.Read(out)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(expected, out) {
				t.Fatalf("expected %#v, got %#v", out, expected)
			}
		}
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("test timeout")
	}
}

func TestNetPipe(t *testing.T) {
	c1, c2, err := pipes.NetPipe()
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		msgs := 50
		size := 1024
//网管阻塞，因此写操作是异步发出的。
		go func() {
			for i := 0; i < msgs; i++ {
				msg := make([]byte, size)
				_ = binary.PutUvarint(msg, uint64(i))

				_, err := c1.Write(msg)
				if err != nil {
					t.Fatal(err)
				}
			}
		}()

		for i := 0; i < msgs; i++ {
			msg := make([]byte, size)
			_ = binary.PutUvarint(msg, uint64(i))

			out := make([]byte, size)
			_, err := c2.Read(out)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(msg, out) {
				t.Fatalf("expected %#v, got %#v", msg, out)
			}
		}

		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("test timeout")
	}
}

func TestNetPipeBidirections(t *testing.T) {
	c1, c2, err := pipes.NetPipe()
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		msgs := 1000
		size := 8
		pingTemplate := "ping %03d"
		pongTemplate := "pong %03d"

//网管阻塞，因此写操作是异步发出的。
		go func() {
			for i := 0; i < msgs; i++ {
				msg := []byte(fmt.Sprintf(pingTemplate, i))

				_, err := c1.Write(msg)
				if err != nil {
					t.Fatal(err)
				}
			}
		}()

//网管阻塞，因此对pong的读取是异步发出的。
		go func() {
			for i := 0; i < msgs; i++ {
				expected := []byte(fmt.Sprintf(pongTemplate, i))

				out := make([]byte, size)
				_, err := c1.Read(out)
				if err != nil {
					t.Fatal(err)
				}

				if !bytes.Equal(expected, out) {
					t.Fatalf("expected %#v, got %#v", expected, out)
				}
			}

			done <- struct{}{}
		}()

//期望读取ping，并用pong响应备用连接
		for i := 0; i < msgs; i++ {
			expected := []byte(fmt.Sprintf(pingTemplate, i))

			out := make([]byte, size)
			_, err := c2.Read(out)
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(expected, out) {
				t.Fatalf("expected %#v, got %#v", expected, out)
			} else {
				msg := []byte(fmt.Sprintf(pongTemplate, i))

				_, err := c2.Write(msg)
				if err != nil {
					t.Fatal(err)
				}
			}
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("test timeout")
	}
}

