
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450120955465728>


package rpc

import (
	"testing"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/ethereum/go-ethereum/swarm/storage/mock/mem"
	"github.com/ethereum/go-ethereum/swarm/storage/mock/test"
)

//testdbstore正在为GlobalStore运行测试
//使用test.mockstore函数。
func TestRPCStore(t *testing.T) {
	serverStore := mem.NewGlobalStore()

	server := rpc.NewServer()
	if err := server.RegisterName("mockStore", serverStore); err != nil {
		t.Fatal(err)
	}

	store := NewGlobalStore(rpc.DialInProc(server))
	defer store.Close()

	test.MockStore(t, store, 30)
}

