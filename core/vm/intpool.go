
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450082627915776>


package vm

import (
	"math/big"
	"sync"
)

var checkVal = big.NewInt(-42)

const poolLimit = 256

//intpool是一个包含大整数的池，
//可用于所有大的.int操作。
type intPool struct {
	pool *Stack
}

func newIntPool() *intPool {
	return &intPool{pool: newstack()}
}

//get从池中检索一个大整数，如果池为空，则分配一个整数。
//注意，返回的int值是任意的，不会归零！
func (p *intPool) get() *big.Int {
	if p.pool.len() > 0 {
		return p.pool.pop()
	}
	return new(big.Int)
}

//GetZero从池中检索一个大整数，将其设置为零或分配
//如果池是空的，就换一个新的。
func (p *intPool) getZero() *big.Int {
	if p.pool.len() > 0 {
		return p.pool.pop().SetUint64(0)
	}
	return new(big.Int)
}

//Put返回一个分配给池的大int，以便稍后由get调用重用。
//注意，保存为原样的值；既不放置也不获取零，整数都不存在！
func (p *intPool) put(is ...*big.Int) {
	if len(p.pool.data) > poolLimit {
		return
	}
	for _, i := range is {
//VerifyPool是一个生成标志。池验证确保完整性
//通过将值与默认值进行比较来获得整数池的值。
		if verifyPool {
			i.Set(checkVal)
		}
		p.pool.push(i)
	}
}

//Intpool池的默认容量
const poolDefaultCap = 25

//IntpoolPool管理Intpools池。
type intPoolPool struct {
	pools []*intPool
	lock  sync.Mutex
}

var poolOfIntPools = &intPoolPool{
	pools: make([]*intPool, 0, poolDefaultCap),
}

//GET正在寻找可返回的可用池。
func (ipp *intPoolPool) get() *intPool {
	ipp.lock.Lock()
	defer ipp.lock.Unlock()

	if len(poolOfIntPools.pools) > 0 {
		ip := ipp.pools[len(ipp.pools)-1]
		ipp.pools = ipp.pools[:len(ipp.pools)-1]
		return ip
	}
	return newIntPool()
}

//放置已分配GET的池。
func (ipp *intPoolPool) put(ip *intPool) {
	ipp.lock.Lock()
	defer ipp.lock.Unlock()

	if len(ipp.pools) < cap(ipp.pools) {
		ipp.pools = append(ipp.pools, ip)
	}
}

