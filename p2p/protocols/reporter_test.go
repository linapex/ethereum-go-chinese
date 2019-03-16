
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450106153766912>


package protocols

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

//TestReporter测试为P2P会计收集的度量
//在重新启动节点后将被持久化并可用。
//它通过重新创建数据库模拟重新启动，就像节点重新启动一样。
func TestReporter(t *testing.T) {
//创建测试目录
	dir, err := ioutil.TempDir("", "reporter-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

//设置指标
	log.Debug("Setting up metrics first time")
	reportInterval := 5 * time.Millisecond
	metrics := SetupAccountingMetrics(reportInterval, filepath.Join(dir, "test.db"))
	log.Debug("Done.")

//做一些度量
	mBalanceCredit.Inc(12)
	mBytesCredit.Inc(34)
	mMsgDebit.Inc(9)

//给报告者时间将指标写入数据库
	time.Sleep(20 * time.Millisecond)

//将度量值设置为零-这有效地模拟了关闭的节点…
	mBalanceCredit = nil
	mBytesCredit = nil
	mMsgDebit = nil
//同时关闭数据库，否则无法创建新数据库
	metrics.Close()

//再次设置指标
	log.Debug("Setting up metrics second time")
	metrics = SetupAccountingMetrics(reportInterval, filepath.Join(dir, "test.db"))
	defer metrics.Close()
	log.Debug("Done.")

//现在检查度量，它们应该与“关闭”之前的值相同。
	if mBalanceCredit.Count() != 12 {
		t.Fatalf("Expected counter to be %d, but is %d", 12, mBalanceCredit.Count())
	}
	if mBytesCredit.Count() != 34 {
		t.Fatalf("Expected counter to be %d, but is %d", 23, mBytesCredit.Count())
	}
	if mMsgDebit.Count() != 9 {
		t.Fatalf("Expected counter to be %d, but is %d", 9, mMsgDebit.Count())
	}
}

