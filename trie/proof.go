
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450123149086720>


package trie

import (
	"bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

//prove为key构造了一个merkle证明。结果包含所有编码节点
//在键处的值的路径上。值本身也包含在最后一个
//节点，可以通过验证证据来检索。
//
//如果trie不包含键的值，则返回的证明包含所有
//键的最长现有前缀的节点（至少是根节点），结束
//证明没有键的节点。
func (t *Trie) Prove(key []byte, fromLevel uint, proofDb ethdb.Putter) error {
//收集键路径上的所有节点。
	key = keybytesToHex(key)
	nodes := []node{}
	tn := t.root
	for len(key) > 0 && tn != nil {
		switch n := tn.(type) {
		case *shortNode:
			if len(key) < len(n.Key) || !bytes.Equal(n.Key, key[:len(n.Key)]) {
//trie不包含密钥。
				tn = nil
			} else {
				tn = n.Val
				key = key[len(n.Key):]
			}
			nodes = append(nodes, n)
		case *fullNode:
			tn = n.Children[key[0]]
			key = key[1:]
			nodes = append(nodes, n)
		case hashNode:
			var err error
			tn, err = t.resolveHash(n, nil)
			if err != nil {
				log.Error(fmt.Sprintf("Unhandled trie error: %v", err))
				return err
			}
		default:
			panic(fmt.Sprintf("%T: invalid node: %v", tn, tn))
		}
	}
	hasher := newHasher(0, 0, nil)
	defer returnHasherToPool(hasher)

	for i, n := range nodes {
//别费心检查这里的错误，因为哈瑟惊慌失措
//如果编码不起作用，而且我们没有写入任何数据库。
		n, _, _ = hasher.hashChildren(n, nil)
		hn, _ := hasher.store(n, nil, false)
		if hash, ok := hn.(hashNode); ok || i == 0 {
//如果节点的数据库编码是哈希（或
//根节点），它成为一个证明元素。
			if fromLevel > 0 {
				fromLevel--
			} else {
				enc, _ := rlp.EncodeToBytes(n)
				if !ok {
					hash = crypto.Keccak256(enc)
				}
				proofDb.Put(hash, enc)
			}
		}
	}
	return nil
}

//prove为key构造了一个merkle证明。结果包含所有编码节点
//在键处的值的路径上。值本身也包含在最后一个
//节点，可以通过验证证据来检索。
//
//如果trie不包含键的值，则返回的证明包含所有
//键的最长现有前缀的节点（至少是根节点），结束
//证明没有键的节点。
func (t *SecureTrie) Prove(key []byte, fromLevel uint, proofDb ethdb.Putter) error {
	return t.trie.Prove(key, fromLevel, proofDb)
}

//验证校对检查Merkle校对。给定的证明必须包含
//用给定的根散列值键入trie。VerifyProof返回一个错误，如果
//证明包含无效的trie节点或错误的值。
func VerifyProof(rootHash common.Hash, key []byte, proofDb DatabaseReader) (value []byte, nodes int, err error) {
	key = keybytesToHex(key)
	wantHash := rootHash
	for i := 0; ; i++ {
		buf, _ := proofDb.Get(wantHash[:])
		if buf == nil {
			return nil, i, fmt.Errorf("proof node %d (hash %064x) missing", i, wantHash)
		}
		n, err := decodeNode(wantHash[:], buf, 0)
		if err != nil {
			return nil, i, fmt.Errorf("bad proof node %d: %v", i, err)
		}
		keyrest, cld := get(n, key)
		switch cld := cld.(type) {
		case nil:
//trie不包含密钥。
			return nil, i, nil
		case hashNode:
			key = keyrest
			copy(wantHash[:], cld)
		case valueNode:
			return cld, i + 1, nil
		}
	}
}

func get(tn node, key []byte) ([]byte, node) {
	for {
		switch n := tn.(type) {
		case *shortNode:
			if len(key) < len(n.Key) || !bytes.Equal(n.Key, key[:len(n.Key)]) {
				return nil, nil
			}
			tn = n.Val
			key = key[len(n.Key):]
		case *fullNode:
			tn = n.Children[key[0]]
			key = key[1:]
		case hashNode:
			return key, n
		case nil:
			return key, nil
		case valueNode:
			return nil, n
		default:
			panic(fmt.Sprintf("%T: invalid node: %v", tn, tn))
		}
	}
}

