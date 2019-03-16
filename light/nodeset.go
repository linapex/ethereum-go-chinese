
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450096456536064>


package light

import (
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
)

//nodeset存储一组trie节点。它实现了trie.database，还可以
//作为另一个trie.database的缓存。
type NodeSet struct {
	nodes map[string][]byte
	order []string

	dataSize int
	lock     sync.RWMutex
}

//newnodeset创建空节点集
func NewNodeSet() *NodeSet {
	return &NodeSet{
		nodes: make(map[string][]byte),
	}
}

//在集合中放置存储新节点
func (db *NodeSet) Put(key []byte, value []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if _, ok := db.nodes[string(key)]; ok {
		return nil
	}
	keystr := string(key)

	db.nodes[keystr] = common.CopyBytes(value)
	db.order = append(db.order, keystr)
	db.dataSize += len(value)

	return nil
}

//get返回存储节点
func (db *NodeSet) Get(key []byte) ([]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if entry, ok := db.nodes[string(key)]; ok {
		return entry, nil
	}
	return nil, errors.New("not found")
}

//如果节点集包含给定的键，则返回true
func (db *NodeSet) Has(key []byte) (bool, error) {
	_, err := db.Get(key)
	return err == nil, nil
}

//keycount返回集合中的节点数
func (db *NodeSet) KeyCount() int {
	db.lock.RLock()
	defer db.lock.RUnlock()

	return len(db.nodes)
}

//data size返回集合中节点的聚合数据大小
func (db *NodeSet) DataSize() int {
	db.lock.RLock()
	defer db.lock.RUnlock()

	return db.dataSize
}

//节点列表将节点集转换为节点列表
func (db *NodeSet) NodeList() NodeList {
	db.lock.RLock()
	defer db.lock.RUnlock()

	var values NodeList
	for _, key := range db.order {
		values = append(values, db.nodes[key])
	}
	return values
}

//存储将集的内容写入给定的数据库
func (db *NodeSet) Store(target ethdb.Putter) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	for key, value := range db.nodes {
		target.Put([]byte(key), value)
	}
}

//nodelist存储trie节点的有序列表。它实现ethdb.putter。
type NodeList []rlp.RawValue

//存储将列表的内容写入给定的数据库
func (n NodeList) Store(db ethdb.Putter) {
	for _, node := range n {
		db.Put(crypto.Keccak256(node), node)
	}
}

//节点集将节点列表转换为节点集
func (n NodeList) NodeSet() *NodeSet {
	db := NewNodeSet()
	n.Store(db)
	return db
}

//在列表末尾放置一个新节点
func (n *NodeList) Put(key []byte, value []byte) error {
	*n = append(*n, value)
	return nil
}

//data size返回列表中节点的聚合数据大小
func (n NodeList) DataSize() int {
	var size int
	for _, node := range n {
		size += len(node)
	}
	return size
}

