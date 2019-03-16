
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450092148985856>


package ethapi

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type AddrLocker struct {
	mu    sync.Mutex
	locks map[common.Address]*sync.Mutex
}

//lock returns the lock of the given address.
func (l *AddrLocker) lock(address common.Address) *sync.Mutex {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.locks == nil {
		l.locks = make(map[common.Address]*sync.Mutex)
	}
	if _, ok := l.locks[address]; !ok {
		l.locks[address] = new(sync.Mutex)
	}
	return l.locks[address]
}

//lockaddr锁定帐户的mutex。这用于阻止另一个Tx获取
//直到释放锁。互斥体阻止（相同的nonce）
//在签署第一个事务期间被再次读取。
func (l *AddrLocker) LockAddr(address common.Address) {
	l.lock(address).Lock()
}

//unlockaddr解锁给定帐户的互斥体。
func (l *AddrLocker) UnlockAddr(address common.Address) {
	l.lock(address).Unlock()
}

