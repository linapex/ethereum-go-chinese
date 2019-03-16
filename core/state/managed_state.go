
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450080115527680>


package state

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type account struct {
	stateObject *stateObject
	nstart      uint64
	nonces      []bool
}

type ManagedState struct {
	*StateDB

	mu sync.RWMutex

	accounts map[common.Address]*account
}

//managedState返回一个新的托管状态，statedb作为它的支持层。
func ManageState(statedb *StateDB) *ManagedState {
	return &ManagedState{
		StateDB:  statedb.Copy(),
		accounts: make(map[common.Address]*account),
	}
}

//SetState sets the backing layer of the managed state
func (ms *ManagedState) SetState(statedb *StateDB) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.StateDB = statedb
}

//removenonce从托管状态中删除了nonce以及所有将来挂起的nonce
func (ms *ManagedState) RemoveNonce(addr common.Address, n uint64) {
	if ms.hasAccount(addr) {
		ms.mu.Lock()
		defer ms.mu.Unlock()

		account := ms.getAccount(addr)
		if n-account.nstart <= uint64(len(account.nonces)) {
			reslice := make([]bool, n-account.nstart)
			copy(reslice, account.nonces[:n-account.nstart])
			account.nonces = reslice
		}
	}
}

//new nonce返回托管帐户的新规范nonce
func (ms *ManagedState) NewNonce(addr common.Address) uint64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	account := ms.getAccount(addr)
	for i, nonce := range account.nonces {
		if !nonce {
			return account.nstart + uint64(i)
		}
	}
	account.nonces = append(account.nonces, true)

	return uint64(len(account.nonces)-1) + account.nstart
}

//GETNONCE返回托管或非托管帐户的规范NoCE。
//
//因为getnonce改变了db，所以我们必须获得一个写锁。
func (ms *ManagedState) GetNonce(addr common.Address) uint64 {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.hasAccount(addr) {
		account := ms.getAccount(addr)
		return uint64(len(account.nonces)) + account.nstart
	} else {
		return ms.StateDB.GetNonce(addr)
	}
}

//setnonce为托管状态设置新的规范nonce
func (ms *ManagedState) SetNonce(addr common.Address, nonce uint64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	so := ms.GetOrNewStateObject(addr)
	so.SetNonce(nonce)

	ms.accounts[addr] = newAccount(so)
}

//hasAccount返回给定地址是否被管理
func (ms *ManagedState) HasAccount(addr common.Address) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.hasAccount(addr)
}

func (ms *ManagedState) hasAccount(addr common.Address) bool {
	_, ok := ms.accounts[addr]
	return ok
}

//填充托管状态
func (ms *ManagedState) getAccount(addr common.Address) *account {
	if account, ok := ms.accounts[addr]; !ok {
		so := ms.GetOrNewStateObject(addr)
		ms.accounts[addr] = newAccount(so)
	} else {
//Always make sure the state account nonce isn't actually higher
//比跟踪的那个。
		so := ms.StateDB.getStateObject(addr)
		if so != nil && uint64(len(account.nonces))+account.nstart < so.Nonce() {
			ms.accounts[addr] = newAccount(so)
		}

	}

	return ms.accounts[addr]
}

func newAccount(so *stateObject) *account {
	return &account{so, so.Nonce(), nil}
}

