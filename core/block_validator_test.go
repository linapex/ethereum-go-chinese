
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450077959655424>


package core

import (
	"runtime"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

//测试简单的头验证对好的和坏的块都有效。
func TestHeaderVerification(t *testing.T) {
//创建一个简单的链来验证
	var (
		testdb    = ethdb.NewMemDatabase()
		gspec     = &Genesis{Config: params.TestChainConfig}
		genesis   = gspec.MustCommit(testdb)
		blocks, _ = GenerateChain(params.TestChainConfig, genesis, ethash.NewFaker(), testdb, 8, nil)
	)
	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
//逐个运行块的头检查程序，检查有效和无效的nonce
	chain, _ := NewBlockChain(testdb, nil, params.TestChainConfig, ethash.NewFaker(), vm.Config{}, nil)
	defer chain.Stop()

	for i := 0; i < len(blocks); i++ {
		for j, valid := range []bool{true, false} {
			var results <-chan error

			if valid {
				engine := ethash.NewFaker()
				_, results = engine.VerifyHeaders(chain, []*types.Header{headers[i]}, []bool{true})
			} else {
				engine := ethash.NewFakeFailer(headers[i].Number.Uint64())
				_, results = engine.VerifyHeaders(chain, []*types.Header{headers[i]}, []bool{true})
			}
//等待验证结果
			select {
			case result := <-results:
				if (result == nil) != valid {
					t.Errorf("test %d.%d: validity mismatch: have %v, want %v", i, j, result, valid)
				}
			case <-time.After(time.Second):
				t.Fatalf("test %d.%d: verification timeout", i, j)
			}
//确保不再返回数据
			select {
			case result := <-results:
				t.Fatalf("test %d.%d: unexpected result returned: %v", i, j, result)
			case <-time.After(25 * time.Millisecond):
			}
		}
		chain.InsertChain(blocks[i : i+1])
	}
}

//测试并发头验证对好块和坏块都有效。
func TestHeaderConcurrentVerification2(t *testing.T)  { testHeaderConcurrentVerification(t, 2) }
func TestHeaderConcurrentVerification8(t *testing.T)  { testHeaderConcurrentVerification(t, 8) }
func TestHeaderConcurrentVerification32(t *testing.T) { testHeaderConcurrentVerification(t, 32) }

func testHeaderConcurrentVerification(t *testing.T, threads int) {
//创建一个简单的链来验证
	var (
		testdb    = ethdb.NewMemDatabase()
		gspec     = &Genesis{Config: params.TestChainConfig}
		genesis   = gspec.MustCommit(testdb)
		blocks, _ = GenerateChain(params.TestChainConfig, genesis, ethash.NewFaker(), testdb, 8, nil)
	)
	headers := make([]*types.Header, len(blocks))
	seals := make([]bool, len(blocks))

	for i, block := range blocks {
		headers[i] = block.Header()
		seals[i] = true
	}
//设置要验证的线程数
	old := runtime.GOMAXPROCS(threads)
	defer runtime.GOMAXPROCS(old)

//同时为整个区块链运行标题检查器
//还有一个无效的链（如果一个任意块无效，就足够了）。
	for i, valid := range []bool{true, false} {
		var results <-chan error

		if valid {
			chain, _ := NewBlockChain(testdb, nil, params.TestChainConfig, ethash.NewFaker(), vm.Config{}, nil)
			_, results = chain.engine.VerifyHeaders(chain, headers, seals)
			chain.Stop()
		} else {
			chain, _ := NewBlockChain(testdb, nil, params.TestChainConfig, ethash.NewFakeFailer(uint64(len(headers)-1)), vm.Config{}, nil)
			_, results = chain.engine.VerifyHeaders(chain, headers, seals)
			chain.Stop()
		}
//等待所有验证结果
		checks := make(map[int]error)
		for j := 0; j < len(blocks); j++ {
			select {
			case result := <-results:
				checks[j] = result

			case <-time.After(time.Second):
				t.Fatalf("test %d.%d: verification timeout", i, j)
			}
		}
//检查当前检查有效性
		for j := 0; j < len(blocks); j++ {
want := valid || (j < len(blocks)-2) //我们选择了链中最后一个，但只有一个失败。
			if (checks[j] == nil) != want {
				t.Errorf("test %d.%d: validity mismatch: have %v, want %v", i, j, checks[j], want)
			}
			if !want {
//第一个错误后的几个块可能由于并发而通过验证
//工人。我们不关心这个测试中的那些，只关心正确的模块
//出错。
				break
			}
		}
//确保不再返回数据
		select {
		case result := <-results:
			t.Fatalf("test %d: unexpected result returned: %v", i, result)
		case <-time.After(25 * time.Millisecond):
		}
	}
}

//中止头验证的测试确实阻止了进一步的检查
//运行，以及检查是否没有遗留的Goroutines泄漏。
func TestHeaderConcurrentAbortion2(t *testing.T)  { testHeaderConcurrentAbortion(t, 2) }
func TestHeaderConcurrentAbortion8(t *testing.T)  { testHeaderConcurrentAbortion(t, 8) }
func TestHeaderConcurrentAbortion32(t *testing.T) { testHeaderConcurrentAbortion(t, 32) }

func testHeaderConcurrentAbortion(t *testing.T, threads int) {
//创建一个简单的链来验证
	var (
		testdb    = ethdb.NewMemDatabase()
		gspec     = &Genesis{Config: params.TestChainConfig}
		genesis   = gspec.MustCommit(testdb)
		blocks, _ = GenerateChain(params.TestChainConfig, genesis, ethash.NewFaker(), testdb, 1024, nil)
	)
	headers := make([]*types.Header, len(blocks))
	seals := make([]bool, len(blocks))

	for i, block := range blocks {
		headers[i] = block.Header()
		seals[i] = true
	}
//设置要验证的线程数
	old := runtime.GOMAXPROCS(threads)
	defer runtime.GOMAXPROCS(old)

//启动验证并立即中止
	chain, _ := NewBlockChain(testdb, nil, params.TestChainConfig, ethash.NewFakeDelayer(time.Millisecond), vm.Config{}, nil)
	defer chain.Stop()

	abort, results := chain.engine.VerifyHeaders(chain, headers, seals)
	close(abort)

//耗尽结果通道
	verified := 0
	for depleted := false; !depleted; {
		select {
		case result := <-results:
			if result != nil {
				t.Errorf("header %d: validation failed: %v", verified, result)
			}
			verified++
		case <-time.After(50 * time.Millisecond):
			depleted = true
		}
	}
//检查流产是否因为没有处理过多的战俘而受到尊重。
	if verified > 2*threads {
		t.Errorf("verification count too large: have %d, want below %d", verified, 2*threads)
	}
}

