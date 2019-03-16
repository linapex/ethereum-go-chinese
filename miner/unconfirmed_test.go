
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450100638257152>


package miner

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

//noopchainretriever是headerretriever的一个实现，它始终
//对于任何请求的头返回nil。
type noopChainRetriever struct{}

func (r *noopChainRetriever) GetHeaderByNumber(number uint64) *types.Header {
	return nil
}
func (r *noopChainRetriever) GetBlockByNumber(number uint64) *types.Block {
	return nil
}

//将块插入未确认集的测试将累积这些块，直到
//达到所需深度后，开始下降。
func TestUnconfirmedInsertBounds(t *testing.T) {
	limit := uint(10)

	pool := newUnconfirmedBlocks(new(noopChainRetriever), limit)
	for depth := uint64(0); depth < 2*uint64(limit); depth++ {
//为同一级别插入多个块以强调它
		for i := 0; i < int(depth); i++ {
			pool.Insert(depth, common.Hash([32]byte{byte(depth), byte(i)}))
		}
//确认深度允许范围以下没有块留在
		pool.blocks.Do(func(block interface{}) {
			if block := block.(*unconfirmedBlock); block.index+uint64(limit) <= depth {
				t.Errorf("depth %d: block %x not dropped", depth, block.hash)
			}
		})
	}
}

//将块移出未确认集的测试正常工作
//箱，以及角箱，如空箱、空班或满箱
//轮班。
func TestUnconfirmedShifts(t *testing.T) {
//在不同深度上创建带有几个块的池
	limit, start := uint(10), uint64(25)

	pool := newUnconfirmedBlocks(new(noopChainRetriever), limit)
	for depth := start; depth < start+uint64(limit); depth++ {
		pool.Insert(depth, common.Hash([32]byte{byte(depth)}))
	}
//试着移到极限以下，确保没有掉块。
	pool.Shift(start + uint64(limit) - 1)
	if n := pool.blocks.Len(); n != int(limit) {
		t.Errorf("unconfirmed count mismatch: have %d, want %d", n, limit)
	}
//试着把一半的积木移走，并核实剩余部分。
	pool.Shift(start + uint64(limit) - 1 + uint64(limit/2))
	if n := pool.blocks.Len(); n != int(limit)/2 {
		t.Errorf("unconfirmed count mismatch: have %d, want %d", n, limit/2)
	}
//尝试将所有剩余的块移出并验证是否为空。
	pool.Shift(start + 2*uint64(limit))
	if n := pool.blocks.Len(); n != 0 {
		t.Errorf("unconfirmed count mismatch: have %d, want %d", n, 0)
	}
//试着从空的那套换出来，确保它不坏
	pool.Shift(start + 3*uint64(limit))
	if n := pool.blocks.Len(); n != 0 {
		t.Errorf("unconfirmed count mismatch: have %d, want %d", n, 0)
	}
}

