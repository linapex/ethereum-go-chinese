
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450105902108672>


package protocols

import (
	"testing"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"
	"github.com/ethereum/go-ethereum/rlp"
)

//虚拟平衡实现
type dummyBalance struct {
	amount int64
	peer   *Peer
}

//虚拟价格实施
type dummyPrices struct{}

//需要基于大小的记帐的虚拟消息
//发送者付费
type perBytesMsgSenderPays struct {
	Content string
}

//需要基于大小的记帐的虚拟消息
//接收者付费
type perBytesMsgReceiverPays struct {
	Content string
}

//按单位付费的假消息
//发送者付费
type perUnitMsgSenderPays struct{}

//接收者付费
type perUnitMsgReceiverPays struct{}

//价格为零的假消息
type zeroPriceMsg struct{}

//没有记帐的虚拟消息
type nilPriceMsg struct{}

//返回已定义消息的价格
func (d *dummyPrices) Price(msg interface{}) *Price {
	switch msg.(type) {
//基于大小的消息成本，接收者支付
	case *perBytesMsgReceiverPays:
		return &Price{
			PerByte: true,
			Value:   uint64(100),
			Payer:   Receiver,
		}
//基于大小的邮件成本，发件人支付
	case *perBytesMsgSenderPays:
		return &Price{
			PerByte: true,
			Value:   uint64(100),
			Payer:   Sender,
		}
//单一成本，收款人支付
	case *perUnitMsgReceiverPays:
		return &Price{
			PerByte: false,
			Value:   uint64(99),
			Payer:   Receiver,
		}
//单一成本，寄件人支付
	case *perUnitMsgSenderPays:
		return &Price{
			PerByte: false,
			Value:   uint64(99),
			Payer:   Sender,
		}
	case *zeroPriceMsg:
		return &Price{
			PerByte: false,
			Value:   uint64(0),
			Payer:   Sender,
		}
	case *nilPriceMsg:
		return nil
	}
	return nil
}

//虚拟会计实现，只存储值供以后检查
func (d *dummyBalance) Add(amount int64, peer *Peer) error {
	d.amount = amount
	d.peer = peer
	return nil
}

type testCase struct {
	msg        interface{}
	size       uint32
	sendResult int64
	recvResult int64
}

//最低水平单元测试
func TestBalance(t *testing.T) {
//创建实例
	balance := &dummyBalance{}
	prices := &dummyPrices{}
//创建规格
	spec := createTestSpec()
//为规范创建记帐挂钩
	acc := NewAccounting(balance, prices)
//创建对等体
	id := adapters.RandomNodeConfig().ID
	p := p2p.NewPeer(id, "testPeer", nil)
	peer := NewPeer(p, &dummyRW{}, spec)
//价格取决于尺寸，接收者支付
	msg := &perBytesMsgReceiverPays{Content: "testBalance"}
	size, _ := rlp.EncodeToBytes(msg)

	testCases := []testCase{
		{
			msg,
			uint32(len(size)),
			int64(len(size) * 100),
			int64(len(size) * -100),
		},
		{
			&perBytesMsgSenderPays{Content: "testBalance"},
			uint32(len(size)),
			int64(len(size) * -100),
			int64(len(size) * 100),
		},
		{
			&perUnitMsgSenderPays{},
			0,
			int64(-99),
			int64(99),
		},
		{
			&perUnitMsgReceiverPays{},
			0,
			int64(99),
			int64(-99),
		},
		{
			&zeroPriceMsg{},
			0,
			int64(0),
			int64(0),
		},
		{
			&nilPriceMsg{},
			0,
			int64(0),
			int64(0),
		},
	}
	checkAccountingTestCases(t, testCases, acc, peer, balance, true)
	checkAccountingTestCases(t, testCases, acc, peer, balance, false)
}

func checkAccountingTestCases(t *testing.T, cases []testCase, acc *Accounting, peer *Peer, balance *dummyBalance, send bool) {
	for _, c := range cases {
		var err error
		var expectedResult int64
//每次检查前重置平衡
		balance.amount = 0
		if send {
			err = acc.Send(peer, c.size, c.msg)
			expectedResult = c.sendResult
		} else {
			err = acc.Receive(peer, c.size, c.msg)
			expectedResult = c.recvResult
		}

		checkResults(t, err, balance, peer, expectedResult)
	}
}

func checkResults(t *testing.T, err error, balance *dummyBalance, peer *Peer, result int64) {
	if err != nil {
		t.Fatal(err)
	}
	if balance.peer != peer {
		t.Fatalf("expected Add to be called with peer %v, got %v", peer, balance.peer)
	}
	if balance.amount != result {
		t.Fatalf("Expected balance to be %d but is %d", result, balance.amount)
	}
}

//创建测试规范
func createTestSpec() *Spec {
	spec := &Spec{
		Name:       "test",
		Version:    42,
		MaxMsgSize: 10 * 1024,
		Messages: []interface{}{
			&perBytesMsgReceiverPays{},
			&perBytesMsgSenderPays{},
			&perUnitMsgReceiverPays{},
			&perUnitMsgSenderPays{},
			&zeroPriceMsg{},
			&nilPriceMsg{},
		},
	}
	return spec
}

