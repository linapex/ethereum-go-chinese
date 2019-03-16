
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450080002281472>


package state

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//journalEntry is a modification entry in the state change journal that can be
//按需归还。
type journalEntry interface {
//REVERT撤消此日记条目引入的更改。
	revert(*StateDB)

//dirtied returns the Ethereum address modified by this journal entry.
	dirtied() *common.Address
}

//日志包含自上次状态以来应用的状态修改列表
//承诺。它们被跟踪，以便在执行时能够还原。
//exception or revertal request.
type journal struct {
entries []journalEntry         //日记帐跟踪的当前更改
dirties map[common.Address]int //脏帐户和更改数
}

//NewJournal创建新的初始化日志。
func newJournal() *journal {
	return &journal{
		dirties: make(map[common.Address]int),
	}
}

//append在变更日志的末尾插入一个新的修改条目。
func (j *journal) append(entry journalEntry) {
	j.entries = append(j.entries, entry)
	if addr := entry.dirtied(); addr != nil {
		j.dirties[*addr]++
	}
}

//REVERT撤消一批日志修改以及任何已还原的
//处理也很脏。
func (j *journal) revert(statedb *StateDB, snapshot int) {
	for i := len(j.entries) - 1; i >= snapshot; i-- {
//撤消操作所做的更改
		j.entries[i].revert(statedb)

//删除由更改引起的任何脏跟踪
		if addr := j.entries[i].dirtied(); addr != nil {
			if j.dirties[*addr]--; j.dirties[*addr] == 0 {
				delete(j.dirties, *addr)
			}
		}
	}
	j.entries = j.entries[:snapshot]
}

//dirty explicitly sets an address to dirty, even if the change entries would
//否则建议清洁。这个方法是处理ripemd的丑陋的黑客
//precompile consensus exception.
func (j *journal) dirty(addr common.Address) {
	j.dirties[addr]++
}

//length返回日记帐中的当前条目数。
func (j *journal) length() int {
	return len(j.entries)
}

type (
//对帐户trie的更改。
	createObjectChange struct {
		account *common.Address
	}
	resetObjectChange struct {
		prev *stateObject
	}
	suicideChange struct {
		account     *common.Address
prev        bool //账户是否已经自杀
		prevbalance *big.Int
	}

//个人账户变更。
	balanceChange struct {
		account *common.Address
		prev    *big.Int
	}
	nonceChange struct {
		account *common.Address
		prev    uint64
	}
	storageChange struct {
		account       *common.Address
		key, prevalue common.Hash
	}
	codeChange struct {
		account            *common.Address
		prevcode, prevhash []byte
	}

//其他状态值的更改。
	refundChange struct {
		prev uint64
	}
	addLogChange struct {
		txhash common.Hash
	}
	addPreimageChange struct {
		hash common.Hash
	}
	touchChange struct {
		account   *common.Address
		prev      bool
		prevDirty bool
	}
)

func (ch createObjectChange) revert(s *StateDB) {
	delete(s.stateObjects, *ch.account)
	delete(s.stateObjectsDirty, *ch.account)
}

func (ch createObjectChange) dirtied() *common.Address {
	return ch.account
}

func (ch resetObjectChange) revert(s *StateDB) {
	s.setStateObject(ch.prev)
}

func (ch resetObjectChange) dirtied() *common.Address {
	return nil
}

func (ch suicideChange) revert(s *StateDB) {
	obj := s.getStateObject(*ch.account)
	if obj != nil {
		obj.suicided = ch.prev
		obj.setBalance(ch.prevbalance)
	}
}

func (ch suicideChange) dirtied() *common.Address {
	return ch.account
}

var ripemd = common.HexToAddress("0000000000000000000000000000000000000003")

func (ch touchChange) revert(s *StateDB) {
}

func (ch touchChange) dirtied() *common.Address {
	return ch.account
}

func (ch balanceChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setBalance(ch.prev)
}

func (ch balanceChange) dirtied() *common.Address {
	return ch.account
}

func (ch nonceChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setNonce(ch.prev)
}

func (ch nonceChange) dirtied() *common.Address {
	return ch.account
}

func (ch codeChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setCode(common.BytesToHash(ch.prevhash), ch.prevcode)
}

func (ch codeChange) dirtied() *common.Address {
	return ch.account
}

func (ch storageChange) revert(s *StateDB) {
	s.getStateObject(*ch.account).setState(ch.key, ch.prevalue)
}

func (ch storageChange) dirtied() *common.Address {
	return ch.account
}

func (ch refundChange) revert(s *StateDB) {
	s.refund = ch.prev
}

func (ch refundChange) dirtied() *common.Address {
	return nil
}

func (ch addLogChange) revert(s *StateDB) {
	logs := s.logs[ch.txhash]
	if len(logs) == 1 {
		delete(s.logs, ch.txhash)
	} else {
		s.logs[ch.txhash] = logs[:len(logs)-1]
	}
	s.logSize--
}

func (ch addLogChange) dirtied() *common.Address {
	return nil
}

func (ch addPreimageChange) revert(s *StateDB) {
	delete(s.preimages, ch.hash)
}

func (ch addPreimageChange) dirtied() *common.Address {
	return nil
}

