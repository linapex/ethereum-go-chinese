
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450080845336576>


package core

import (
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

//测试可以将事务添加到严格的列表和列表内容以及
//正确维护当前边界。
func TestStrictTxListAdd(t *testing.T) {
//生成要插入的事务列表
	key, _ := crypto.GenerateKey()

	txs := make(types.Transactions, 1024)
	for i := 0; i < len(txs); i++ {
		txs[i] = transaction(uint64(i), 0, key)
	}
//按随机顺序插入交易记录
	list := newTxList(true)
	for _, v := range rand.Perm(len(txs)) {
		list.Add(txs[v], DefaultTxPoolConfig.PriceBump)
	}
//验证内部状态
	if len(list.txs.items) != len(txs) {
		t.Errorf("transaction count mismatch: have %d, want %d", len(list.txs.items), len(txs))
	}
	for i, tx := range txs {
		if list.txs.items[tx.Nonce()] != tx {
			t.Errorf("item %d: transaction mismatch: have %v, want %v", i, list.txs.items[tx.Nonce()], tx)
		}
	}
}

