
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450080744673280>


package core

import (
	"errors"
	"io"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

//如果试图插入事务，则返回errnoaactivejournal
//进入日志，但当前没有打开该文件。
var errNoActiveJournal = errors.New("no active journal")

//devnull是一个只丢弃写入其中的任何内容的writecloser。它的
//目标是允许事务日记帐在以下情况下写入假日记帐：
//由于没有文件，在启动时加载事务，但不打印警告
//为写而读。
type devNull struct{}

func (*devNull) Write(p []byte) (n int, err error) { return len(p), nil }
func (*devNull) Close() error                      { return nil }

//txjournal是一个旋转的事务日志，目的是在本地存储
//创建的事务允许未执行的事务在节点重新启动后继续存在。
type txJournal struct {
path   string         //存储事务的文件系统路径
writer io.WriteCloser //将新事务写入的输出流
}

//newtxjournal创建新的交易日记帐到
func newTxJournal(path string) *txJournal {
	return &txJournal{
		path: path,
	}
}

//LOAD分析事务日志从磁盘转储，将其内容加载到
//指定的池。
func (journal *txJournal) load(add func([]*types.Transaction) []error) error {
//如果日志文件根本不存在，则跳过分析
	if _, err := os.Stat(journal.path); os.IsNotExist(err) {
		return nil
	}
//打开日记帐以加载任何过去的交易记录
	input, err := os.Open(journal.path)
	if err != nil {
		return err
	}
	defer input.Close()

//暂时丢弃任何日志添加（加载时不要重复添加）
	journal.writer = new(devNull)
	defer func() { journal.writer = nil }()

//将日记帐中的所有交易记录插入池
	stream := rlp.NewStream(input, 0)
	total, dropped := 0, 0

//创建一个方法来加载有限批事务并
//适当的进度计数器。然后使用此方法加载
//以小批量记录交易。
	loadBatch := func(txs types.Transactions) {
		for _, err := range add(txs) {
			if err != nil {
				log.Debug("Failed to add journaled transaction", "err", err)
				dropped++
			}
		}
	}
	var (
		failure error
		batch   types.Transactions
	)
	for {
//分析下一个事务并在出错时终止
		tx := new(types.Transaction)
		if err = stream.Decode(tx); err != nil {
			if err != io.EOF {
				failure = err
			}
			if batch.Len() > 0 {
				loadBatch(batch)
			}
			break
		}
//已分析新事务，排队等待稍后，达到阈值时导入
		total++

		if batch = append(batch, tx); batch.Len() > 1024 {
			loadBatch(batch)
			batch = batch[:0]
		}
	}
	log.Info("Loaded local transaction journal", "transactions", total, "dropped", dropped)

	return failure
}

//insert将指定的事务添加到本地磁盘日志。
func (journal *txJournal) insert(tx *types.Transaction) error {
	if journal.writer == nil {
		return errNoActiveJournal
	}
	if err := rlp.Encode(journal.writer, tx); err != nil {
		return err
	}
	return nil
}

//Rotate根据当前的内容重新生成事务日记帐
//事务池。
func (journal *txJournal) rotate(all map[common.Address]types.Transactions) error {
//关闭当前日记帐（如果有打开的日记帐）
	if journal.writer != nil {
		if err := journal.writer.Close(); err != nil {
			return err
		}
		journal.writer = nil
	}
//生成包含当前池内容的新日记
	replacement, err := os.OpenFile(journal.path+".new", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	journaled := 0
	for _, txs := range all {
		for _, tx := range txs {
			if err = rlp.Encode(replacement, tx); err != nil {
				replacement.Close()
				return err
			}
		}
		journaled += len(txs)
	}
	replacement.Close()

//用新生成的日志替换活日志
	if err = os.Rename(journal.path+".new", journal.path); err != nil {
		return err
	}
	sink, err := os.OpenFile(journal.path, os.O_WRONLY|os.O_APPEND, 0755)
	if err != nil {
		return err
	}
	journal.writer = sink
	log.Info("Regenerated local transaction journal", "transactions", journaled, "accounts", len(all))

	return nil
}

//close将事务日志内容刷新到磁盘并关闭文件。
func (journal *txJournal) close() error {
	var err error

	if journal.writer != nil {
		err = journal.writer.Close()
		journal.writer = nil
	}
	return err
}

