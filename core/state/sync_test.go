
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450080560123904>


package state

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/trie"
)

//TestAccount是与状态测试使用的帐户关联的数据。
type testAccount struct {
	address common.Address
	balance *big.Int
	nonce   uint64
	code    []byte
}

//maketeststate创建一个样本测试状态来测试节点重构。
func makeTestState() (Database, common.Hash, []*testAccount) {
//创建空状态
	db := NewDatabase(ethdb.NewMemDatabase())
	state, _ := New(common.Hash{}, db)

//用任意数据填充它
	accounts := []*testAccount{}
	for i := byte(0); i < 96; i++ {
		obj := state.GetOrNewStateObject(common.BytesToAddress([]byte{i}))
		acc := &testAccount{address: common.BytesToAddress([]byte{i})}

		obj.AddBalance(big.NewInt(int64(11 * i)))
		acc.balance = big.NewInt(int64(11 * i))

		obj.SetNonce(uint64(42 * i))
		acc.nonce = uint64(42 * i)

		if i%3 == 0 {
			obj.SetCode(crypto.Keccak256Hash([]byte{i, i, i, i, i}), []byte{i, i, i, i, i})
			acc.code = []byte{i, i, i, i, i}
		}
		state.updateStateObject(obj)
		accounts = append(accounts, acc)
	}
	root, _ := state.Commit(false)

//返回生成的状态
	return db, root, accounts
}

//checkStateAccounts交叉引用一个重构的状态，该状态应为
//帐户数组。
func checkStateAccounts(t *testing.T, db ethdb.Database, root common.Hash, accounts []*testAccount) {
//检查根可用性和状态内容
	state, err := New(root, NewDatabase(db))
	if err != nil {
		t.Fatalf("failed to create state trie at %x: %v", root, err)
	}
	if err := checkStateConsistency(db, root); err != nil {
		t.Fatalf("inconsistent state trie at %x: %v", root, err)
	}
	for i, acc := range accounts {
		if balance := state.GetBalance(acc.address); balance.Cmp(acc.balance) != 0 {
			t.Errorf("account %d: balance mismatch: have %v, want %v", i, balance, acc.balance)
		}
		if nonce := state.GetNonce(acc.address); nonce != acc.nonce {
			t.Errorf("account %d: nonce mismatch: have %v, want %v", i, nonce, acc.nonce)
		}
		if code := state.GetCode(acc.address); !bytes.Equal(code, acc.code) {
			t.Errorf("account %d: code mismatch: have %x, want %x", i, code, acc.code)
		}
	}
}

//checktrie一致性检查（sub-）trie中的所有节点是否确实存在。
func checkTrieConsistency(db ethdb.Database, root common.Hash) error {
	if v, _ := db.Get(root[:]); v == nil {
return nil //认为不存在的状态是一致的。
	}
	trie, err := trie.New(root, trie.NewDatabase(db))
	if err != nil {
		return err
	}
	it := trie.NodeIterator(nil)
	for it.Next(true) {
	}
	return it.Error()
}

//checkstateconsistency检查状态根的所有数据是否存在。
func checkStateConsistency(db ethdb.Database, root common.Hash) error {
//创建并迭代子节点中的状态trie
	if _, err := db.Get(root.Bytes()); err != nil {
return nil //认为不存在的状态是一致的。
	}
	state, err := New(root, NewDatabase(db))
	if err != nil {
		return err
	}
	it := NewNodeIterator(state)
	for it.Next() {
	}
	return it.Error
}

//测试是否未计划空状态进行同步。
func TestEmptyStateSync(t *testing.T) {
	empty := common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	if req := NewStateSync(empty, ethdb.NewMemDatabase()).Missing(1); len(req) != 0 {
		t.Errorf("content requested for empty state: %v", req)
	}
}

//测试给定根哈希，状态可以在单个线程上迭代同步，
//请求检索任务并一次性返回所有任务。
func TestIterativeStateSyncIndividual(t *testing.T) { testIterativeStateSync(t, 1) }
func TestIterativeStateSyncBatched(t *testing.T)    { testIterativeStateSync(t, 100) }

func testIterativeStateSync(t *testing.T, batch int) {
//创建要复制的随机状态
	srcDb, srcRoot, srcAccounts := makeTestState()

//Create a destination state and sync with the scheduler
	dstDb := ethdb.NewMemDatabase()
	sched := NewStateSync(srcRoot, dstDb)

	queue := append([]common.Hash{}, sched.Missing(batch)...)
	for len(queue) > 0 {
		results := make([]trie.SyncResult, len(queue))
		for i, hash := range queue {
			data, err := srcDb.TrieDB().Node(hash)
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results[i] = trie.SyncResult{Hash: hash, Data: data}
		}
		if _, index, err := sched.Process(results); err != nil {
			t.Fatalf("failed to process result #%d: %v", index, err)
		}
		if index, err := sched.Commit(dstDb); err != nil {
			t.Fatalf("failed to commit data #%d: %v", index, err)
		}
		queue = append(queue[:0], sched.Missing(batch)...)
	}
//交叉检查两个状态是否同步
	checkStateAccounts(t, dstDb, srcRoot, srcAccounts)
}

//测试trie调度程序是否可以正确地重建状态，即使只有
//partial results are returned, and the others sent only later.
func TestIterativeDelayedStateSync(t *testing.T) {
//创建要复制的随机状态
	srcDb, srcRoot, srcAccounts := makeTestState()

//创建目标状态并与调度程序同步
	dstDb := ethdb.NewMemDatabase()
	sched := NewStateSync(srcRoot, dstDb)

	queue := append([]common.Hash{}, sched.Missing(0)...)
	for len(queue) > 0 {
//只同步一半的计划节点
		results := make([]trie.SyncResult, len(queue)/2+1)
		for i, hash := range queue[:len(results)] {
			data, err := srcDb.TrieDB().Node(hash)
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results[i] = trie.SyncResult{Hash: hash, Data: data}
		}
		if _, index, err := sched.Process(results); err != nil {
			t.Fatalf("failed to process result #%d: %v", index, err)
		}
		if index, err := sched.Commit(dstDb); err != nil {
			t.Fatalf("failed to commit data #%d: %v", index, err)
		}
		queue = append(queue[len(results):], sched.Missing(0)...)
	}
//交叉检查两个状态是否同步
	checkStateAccounts(t, dstDb, srcRoot, srcAccounts)
}

//测试给定根哈希，trie可以在单个线程上迭代同步，
//请求检索任务并一次性返回所有任务，但是在
//随机顺序。
func TestIterativeRandomStateSyncIndividual(t *testing.T) { testIterativeRandomStateSync(t, 1) }
func TestIterativeRandomStateSyncBatched(t *testing.T)    { testIterativeRandomStateSync(t, 100) }

func testIterativeRandomStateSync(t *testing.T, batch int) {
//创建要复制的随机状态
	srcDb, srcRoot, srcAccounts := makeTestState()

//创建目标状态并与调度程序同步
	dstDb := ethdb.NewMemDatabase()
	sched := NewStateSync(srcRoot, dstDb)

	queue := make(map[common.Hash]struct{})
	for _, hash := range sched.Missing(batch) {
		queue[hash] = struct{}{}
	}
	for len(queue) > 0 {
//以随机顺序获取所有排队的节点
		results := make([]trie.SyncResult, 0, len(queue))
		for hash := range queue {
			data, err := srcDb.TrieDB().Node(hash)
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results = append(results, trie.SyncResult{Hash: hash, Data: data})
		}
//将检索到的结果反馈并将新任务排队
		if _, index, err := sched.Process(results); err != nil {
			t.Fatalf("failed to process result #%d: %v", index, err)
		}
		if index, err := sched.Commit(dstDb); err != nil {
			t.Fatalf("failed to commit data #%d: %v", index, err)
		}
		queue = make(map[common.Hash]struct{})
		for _, hash := range sched.Missing(batch) {
			queue[hash] = struct{}{}
		}
	}
//交叉检查两个状态是否同步
	checkStateAccounts(t, dstDb, srcRoot, srcAccounts)
}

//测试trie调度程序是否可以正确地重建状态，即使只有
//部分结果会被返回（甚至是随机返回的结果），其他结果只会在稍后发送。
func TestIterativeRandomDelayedStateSync(t *testing.T) {
//创建要复制的随机状态
	srcDb, srcRoot, srcAccounts := makeTestState()

//创建目标状态并与调度程序同步
	dstDb := ethdb.NewMemDatabase()
	sched := NewStateSync(srcRoot, dstDb)

	queue := make(map[common.Hash]struct{})
	for _, hash := range sched.Missing(0) {
		queue[hash] = struct{}{}
	}
	for len(queue) > 0 {
//只同步一半的计划节点，甚至是随机顺序的节点
		results := make([]trie.SyncResult, 0, len(queue)/2+1)
		for hash := range queue {
			delete(queue, hash)

			data, err := srcDb.TrieDB().Node(hash)
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results = append(results, trie.SyncResult{Hash: hash, Data: data})

			if len(results) >= cap(results) {
				break
			}
		}
//将检索到的结果反馈并将新任务排队
		if _, index, err := sched.Process(results); err != nil {
			t.Fatalf("failed to process result #%d: %v", index, err)
		}
		if index, err := sched.Commit(dstDb); err != nil {
			t.Fatalf("failed to commit data #%d: %v", index, err)
		}
		for _, hash := range sched.Missing(0) {
			queue[hash] = struct{}{}
		}
	}
//交叉检查两个状态是否同步
	checkStateAccounts(t, dstDb, srcRoot, srcAccounts)
}

//测试在同步过程中的任何时间点，只有完整的子尝试处于
//数据库。
func TestIncompleteStateSync(t *testing.T) {
//创建要复制的随机状态
	srcDb, srcRoot, srcAccounts := makeTestState()

	checkTrieConsistency(srcDb.TrieDB().DiskDB().(ethdb.Database), srcRoot)

//创建目标状态并与调度程序同步
	dstDb := ethdb.NewMemDatabase()
	sched := NewStateSync(srcRoot, dstDb)

	added := []common.Hash{}
	queue := append([]common.Hash{}, sched.Missing(1)...)
	for len(queue) > 0 {
//获取一批状态节点
		results := make([]trie.SyncResult, len(queue))
		for i, hash := range queue {
			data, err := srcDb.TrieDB().Node(hash)
			if err != nil {
				t.Fatalf("failed to retrieve node data for %x", hash)
			}
			results[i] = trie.SyncResult{Hash: hash, Data: data}
		}
//处理每个状态节点
		if _, index, err := sched.Process(results); err != nil {
			t.Fatalf("failed to process result #%d: %v", index, err)
		}
		if index, err := sched.Commit(dstDb); err != nil {
			t.Fatalf("failed to commit data #%d: %v", index, err)
		}
		for _, result := range results {
			added = append(added, result.Hash)
		}
//检查到目前为止添加的所有已知子尝试是否完全完成或丢失。
	checkSubtries:
		for _, hash := range added {
			for _, acc := range srcAccounts {
				if hash == crypto.Keccak256Hash(acc.code) {
continue checkSubtries //跳过代码节点的三检。
				}
			}
//无法在此处使用CheckStateConsistency，因为子RIE键可能具有奇数
//长度和撞击力。
			if err := checkTrieConsistency(dstDb, hash); err != nil {
				t.Fatalf("state inconsistent: %v", err)
			}
		}
//获取要检索的下一批
		queue = append(queue[:0], sched.Missing(1)...)
	}
//健全性检查是否检测到从数据库中删除任何节点
	for _, node := range added[1:] {
		key := node.Bytes()
		value, _ := dstDb.Get(key)

		dstDb.Delete(key)
		if err := checkStateConsistency(dstDb, added[0]); err == nil {
			t.Fatalf("trie inconsistency not caught, missing: %x", key)
		}
		dstDb.Put(key, value)
	}
}

