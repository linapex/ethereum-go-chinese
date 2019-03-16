
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450079884840960>


package state

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

//nodeiterator是遍历整个状态trie post顺序的迭代器，
//including all of the contract code and contract state tries.
type NodeIterator struct {
state *StateDB //正在迭代的状态

stateIt trie.NodeIterator //全局状态trie的主迭代器
dataIt  trie.NodeIterator //Secondary iterator for the data trie of a contract

accountHash common.Hash //包含帐户的节点的哈希
codeHash    common.Hash //合同源代码的哈希
code        []byte      //与合同相关的源代码

Hash   common.Hash //正在迭代的当前条目的哈希（如果不是独立的，则为零）
Parent common.Hash //第一个完整祖先节点的哈希（如果当前是根节点，则为零）

Error error //迭代器中出现内部错误时的故障集
}

//newnodeiterator创建一个后序状态节点迭代器。
func NewNodeIterator(state *StateDB) *NodeIterator {
	return &NodeIterator{
		state: state,
	}
}

//next将迭代器移动到下一个节点，返回是否存在
//进一步的节点。如果出现内部错误，此方法将返回false，并且
//将错误字段设置为遇到的故障。
func (it *NodeIterator) Next() bool {
//如果迭代器以前失败，则不要执行任何操作
	if it.Error != nil {
		return false
	}
//否则，使用迭代器前进并报告任何错误
	if err := it.step(); err != nil {
		it.Error = err
		return false
	}
	return it.retrieve()
}

//步骤将迭代器移动到状态trie的下一个条目。
func (it *NodeIterator) step() error {
//如果到达迭代结束，则中止
	if it.state == nil {
		return nil
	}
//如果我们刚开始初始化迭代器
	if it.stateIt == nil {
		it.stateIt = it.state.trie.NodeIterator(nil)
	}
//如果我们以前有数据节点，那么我们肯定至少有状态节点
	if it.dataIt != nil {
		if cont := it.dataIt.Next(true); !cont {
			if it.dataIt.Error() != nil {
				return it.dataIt.Error()
			}
			it.dataIt = nil
		}
		return nil
	}
//If we had source code previously, discard that
	if it.code != nil {
		it.code = nil
		return nil
	}
//进入下一个状态trie节点，如果节点用完则终止
	if cont := it.stateIt.Next(true); !cont {
		if it.stateIt.Error() != nil {
			return it.stateIt.Error()
		}
		it.state, it.stateIt = nil, nil
		return nil
	}
//如果状态trie节点是内部条目，则保持原样
	if !it.stateIt.Leaf() {
		return nil
	}
//否则，我们将到达一个帐户节点，开始数据迭代
	var account Account
	if err := rlp.Decode(bytes.NewReader(it.stateIt.LeafBlob()), &account); err != nil {
		return err
	}
	dataTrie, err := it.state.db.OpenStorageTrie(common.BytesToHash(it.stateIt.LeafKey()), account.Root)
	if err != nil {
		return err
	}
	it.dataIt = dataTrie.NodeIterator(nil)
	if !it.dataIt.Next(true) {
		it.dataIt = nil
	}
	if !bytes.Equal(account.CodeHash, emptyCodeHash) {
		it.codeHash = common.BytesToHash(account.CodeHash)
		addrHash := common.BytesToHash(it.stateIt.LeafKey())
		it.code, err = it.state.db.ContractCode(addrHash, common.BytesToHash(account.CodeHash))
		if err != nil {
			return fmt.Errorf("code %x: %v", account.CodeHash, err)
		}
	}
	it.accountHash = it.stateIt.Parent()
	return nil
}

//检索拉取和缓存迭代器正在遍历的当前状态条目。
//该方法返回是否还有其他数据要检查。
func (it *NodeIterator) retrieve() bool {
//清除任何预先设置的值
	it.Hash = common.Hash{}

//如果迭代完成，则不返回可用数据
	if it.state == nil {
		return false
	}
//否则检索当前条目
	switch {
	case it.dataIt != nil:
		it.Hash, it.Parent = it.dataIt.Hash(), it.dataIt.Parent()
		if it.Parent == (common.Hash{}) {
			it.Parent = it.accountHash
		}
	case it.code != nil:
		it.Hash, it.Parent = it.codeHash, it.accountHash
	case it.stateIt != nil:
		it.Hash, it.Parent = it.stateIt.Hash(), it.stateIt.Parent()
	}
	return true
}

