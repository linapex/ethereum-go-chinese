
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450106082463744>


package protocols

import (
	"encoding/binary"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"

	"github.com/syndtr/goleveldb/leveldb"
)

//AccountMetrics提取了度量数据库，并
//报告者坚持标准
type AccountingMetrics struct {
	reporter *reporter
}

//关闭节点时将调用Close。
//为了优雅的清理
func (am *AccountingMetrics) Close() {
	close(am.reporter.quit)
	am.reporter.db.Close()
}

//Reporter是一种内部结构，用于编写与P2P会计相关的
//指标达到了B级。它将定期向数据库写入应计指标。
type reporter struct {
reg      metrics.Registry //这些度量的注册表（独立于其他度量）
interval time.Duration    //报告者将保持度量的持续时间
db       *leveldb.DB      //实际数据库
quit     chan struct{}    //退出Reporter循环
}

//NewMetricsDB创建一个新的LevelDB实例，用于持久化定义的度量
//在p2p/protocols/accounting.go中
func NewAccountingMetrics(r metrics.Registry, d time.Duration, path string) *AccountingMetrics {
	var val = make([]byte, 8)
	var err error

//创建级别数据库
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

//检查数据库中是否存在值的所有已定义度量
//如果存在，请将其分配给度量。这意味着节点
//以前一直在运行，并且该度量值已被持久化。
	metricsMap := map[string]metrics.Counter{
		"account.balance.credit": mBalanceCredit,
		"account.balance.debit":  mBalanceDebit,
		"account.bytes.credit":   mBytesCredit,
		"account.bytes.debit":    mBytesDebit,
		"account.msg.credit":     mMsgCredit,
		"account.msg.debit":      mMsgDebit,
		"account.peerdrops":      mPeerDrops,
		"account.selfdrops":      mSelfDrops,
	}
//迭代映射并获取值
	for key, metric := range metricsMap {
		val, err = db.Get([]byte(key), nil)
//直到第一次写入值，
//这将返回一个错误。
//尽管以后记录错误是有益的，
//但这需要一个不同的逻辑
		if err == nil {
			metric.Inc(int64(binary.BigEndian.Uint64(val)))
		}
	}

//创建报告人
	rep := &reporter{
		reg:      r,
		interval: d,
		db:       db,
		quit:     make(chan struct{}),
	}

//执行执行例行程序
	go rep.run()

	m := &AccountingMetrics{
		reporter: rep,
	}

	return m
}

//运行是一个goroutine，它定期将度量发送到配置的级别db
func (r *reporter) run() {
	intervalTicker := time.NewTicker(r.interval)

	for {
		select {
		case <-intervalTicker.C:
//在每个勾选处发送指标
			if err := r.save(); err != nil {
				log.Error("unable to send metrics to LevelDB", "err", err)
//如果在写入过程中出现错误，请退出该例程；我们在此假定错误为
//严重，不要再尝试写入。
//此外，这应该可以防止节点停止时发生泄漏。
				return
			}
		case <-r.quit:
//正常关机
			return
		}
	}
}

//将指标发送到数据库
func (r *reporter) save() error {
//创建一个级别数据库批处理
	batch := leveldb.Batch{}
//对于注册表中的每个指标（独立的）
	r.reg.Each(func(name string, i interface{}) {
		metric, ok := i.(metrics.Counter)
		if ok {
//假设这里的每个度量都是一个计数器（单独的注册表）
//…创建快照…
			ms := metric.Snapshot()
			byteVal := make([]byte, 8)
			binary.BigEndian.PutUint64(byteVal, uint64(ms.Count()))
//…并将值保存到数据库
			batch.Put([]byte(name), byteVal)
		}
	})
	return r.db.Write(&batch, nil)
}

