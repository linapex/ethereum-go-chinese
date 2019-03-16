
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450080706924544>


package core

import (
	"runtime"

	"github.com/ethereum/go-ethereum/core/types"
)

//SENDECCHACER是并发事务发送者恢复和缓存。
var senderCacher = newTxSenderCacher(runtime.NumCPU())

//txsendercacherRequest是一个用于恢复事务发送方的请求，
//特定的签名方案并将其缓存到事务本身中。
//
//inc字段定义每次恢复后要跳过的事务数，
//它用于向不同的线程提供相同的基础输入数组，但
//确保他们快速处理早期事务。
type txSenderCacherRequest struct {
	signer types.Signer
	txs    []*types.Transaction
	inc    int
}

//txsendercacher是用于并发ecrecover事务的辅助结构
//来自后台线程上数字签名的发件人。
type txSenderCacher struct {
	threads int
	tasks   chan *txSenderCacherRequest
}

//newTxSenderCacher creates a new transaction sender background cacher and starts
//gomaxprocs在构建时允许的处理goroutine的数量。
func newTxSenderCacher(threads int) *txSenderCacher {
	cacher := &txSenderCacher{
		tasks:   make(chan *txSenderCacherRequest, threads),
		threads: threads,
	}
	for i := 0; i < threads; i++ {
		go cacher.cache()
	}
	return cacher
}

//缓存是一个无限循环，缓存来自各种形式的事务发送者
//数据结构。
func (cacher *txSenderCacher) cache() {
	for task := range cacher.tasks {
		for i := 0; i < len(task.txs); i += task.inc {
			types.Sender(task.signer, task.txs[i])
		}
	}
}

//recover从一批事务中恢复发送方并缓存它们
//回到相同的数据结构中。没有进行验证，也没有
//对无效签名的任何反应。这取决于以后调用代码。
func (cacher *txSenderCacher) recover(signer types.Signer, txs []*types.Transaction) {
//如果没有什么可恢复的，中止
	if len(txs) == 0 {
		return
	}
//确保我们拥有有意义的任务规模并计划恢复
	tasks := cacher.threads
	if len(txs) < tasks*4 {
		tasks = (len(txs) + 3) / 4
	}
	for i := 0; i < tasks; i++ {
		cacher.tasks <- &txSenderCacherRequest{
			signer: signer,
			txs:    txs[i:],
			inc:    tasks,
		}
	}
}

//恢复器块从批处理中恢复发件人并缓存它们。
//回到相同的数据结构中。没有进行验证，也没有
//对无效签名的任何反应。这取决于以后调用代码。
func (cacher *txSenderCacher) recoverFromBlocks(signer types.Signer, blocks []*types.Block) {
	count := 0
	for _, block := range blocks {
		count += len(block.Transactions())
	}
	txs := make([]*types.Transaction, 0, count)
	for _, block := range blocks {
		txs = append(txs, block.Transactions()...)
	}
	cacher.recover(signer, txs)
}

