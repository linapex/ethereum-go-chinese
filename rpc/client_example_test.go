
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450108527742976>


package rpc_test

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

//在本例中，我们的客户希望跟踪最新的“块编号”
//服务器已知。服务器支持两种方法：
//
//eth_getBlockByNumber（“最新”，）
//返回最新的块对象。
//
//ETH订阅（“newblocks”）
//创建在新块到达时激发块对象的订阅。

type Block struct {
	Number *big.Int
}

func ExampleClientSubscription() {
//连接客户端。
client, _ := rpc.Dial("ws://127.0.0.1:8485“）
	subch := make(chan Block)

//确保Subch接收到最新的块。
	go func() {
		for i := 0; ; i++ {
			if i > 0 {
				time.Sleep(2 * time.Second)
			}
			subscribeBlocks(client, subch)
		}
	}()

//到达时打印订阅中的事件。
	for block := range subch {
		fmt.Println("latest block:", block.Number)
	}
}

//subscribeBlocks在自己的goroutine中运行并维护
//新块的订阅。
func subscribeBlocks(client *rpc.Client, subch chan Block) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

//订阅新块。
	sub, err := client.EthSubscribe(ctx, subch, "newHeads")
	if err != nil {
		fmt.Println("subscribe error:", err)
		return
	}

//现在已建立连接。
//用当前块更新频道。
	var lastBlock Block
	if err := client.CallContext(ctx, &lastBlock, "eth_getBlockByNumber", "latest"); err != nil {
		fmt.Println("can't get latest block:", err)
		return
	}
	subch <- lastBlock

//订阅将向通道传递事件。等待
//订阅以任何原因结束，然后循环重新建立
//连接。
	fmt.Println("connection lost: ", <-sub.Err())
}

