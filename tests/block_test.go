
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450121886601216>


package tests

import (
	"testing"
)

func TestBlockchain(t *testing.T) {
	t.Parallel()

	bt := new(testMatcher)
//一般的状态测试作为区块链测试“导出”，但我们可以在本地运行它们。
	bt.skipLoad(`^GeneralStateTests/`)
//跳过由于自私的挖掘测试而导致的随机失败。
	bt.skipLoad(`^bcForgedTest/bcForkUncle\.json`)
	bt.skipLoad(`^bcMultiChainTest/(ChainAtoChainB_blockorder|CallContractFromNotBestBlock)`)
	bt.skipLoad(`^bcTotalDifficultyTest/(lotsOfLeafs|lotsOfBranches|sideChainWithMoreTransactions)`)
//慢测试
	bt.slow(`^bcExploitTest/DelegateCallSpam.json`)
	bt.slow(`^bcExploitTest/ShanghaiLove.json`)
	bt.slow(`^bcExploitTest/SuicideIssue.json`)
	bt.slow(`^bcForkStressTest/`)
	bt.slow(`^bcGasPricerTest/RPC_API_Test.json`)
	bt.slow(`^bcWalletTest/`)

//仍未能通过我们需要调查的测试
//bt.失败（`^bcstatetests/suicidethecheckbalance.json/suicidethecheckbalance constantinople`，'todo:investive'）
//bt.失败（`^bcstatetests/suicidestoragecheckvcreate2.json/suicidestoragecheckvcreate2_Constantinople`，'todo:investive'）
//bt.失败（`^bcstatetests/suicidestoragecheckvcreate.json/suicidestoragecheckvcreate_Constantinople`，'todo:investive'）
//bt.失败（`^bcstatetests/suicidestoragecheck.json/suicidestoragecheck_constantinople`，'todo:investive'）

	bt.walk(t, blockTestDir, func(t *testing.T, name string, test *BlockTest) {
		if err := bt.checkFailure(t, name, test.Run()); err != nil {
			t.Error(err)
		}
	})
}

