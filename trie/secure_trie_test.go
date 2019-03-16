
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450123396550656>


package trie

import (
	"bytes"
	"runtime"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
)

func newEmptySecure() *SecureTrie {
	trie, _ := NewSecure(common.Hash{}, NewDatabase(ethdb.NewMemDatabase()), 0)
	return trie
}

//maketestsecuretrie为测试创建足够大的安全trie。
func makeTestSecureTrie() (*Database, *SecureTrie, map[string][]byte) {
//创建空的trie
	triedb := NewDatabase(ethdb.NewMemDatabase())

	trie, _ := NewSecure(common.Hash{}, triedb, 0)

//用任意数据填充它
	content := make(map[string][]byte)
	for i := byte(0); i < 255; i++ {
//在多个键下映射相同的数据
		key, val := common.LeftPadBytes([]byte{1, i}, 32), []byte{i}
		content[string(key)] = val
		trie.Update(key, val)

		key, val = common.LeftPadBytes([]byte{2, i}, 32), []byte{i}
		content[string(key)] = val
		trie.Update(key, val)

//添加一些其他数据来填充trie
		for j := byte(3); j < 13; j++ {
			key, val = common.LeftPadBytes([]byte{j, i}, 32), []byte{j, i}
			content[string(key)] = val
			trie.Update(key, val)
		}
	}
	trie.Commit(nil)

//返回生成的trie
	return triedb, trie, content
}

func TestSecureDelete(t *testing.T) {
	trie := newEmptySecure()
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"ether", ""},
		{"dog", "puppy"},
		{"shaman", ""},
	}
	for _, val := range vals {
		if val.v != "" {
			trie.Update([]byte(val.k), []byte(val.v))
		} else {
			trie.Delete([]byte(val.k))
		}
	}
	hash := trie.Hash()
	exp := common.HexToHash("29b235a58c3c25ab83010c327d5932bcf05324b7d6b1185e650798034783ca9d")
	if hash != exp {
		t.Errorf("expected %x got %x", exp, hash)
	}
}

func TestSecureGetKey(t *testing.T) {
	trie := newEmptySecure()
	trie.Update([]byte("foo"), []byte("bar"))

	key := []byte("foo")
	value := []byte("bar")
	seckey := crypto.Keccak256(key)

	if !bytes.Equal(trie.Get(key), value) {
		t.Errorf("Get did not return bar")
	}
	if k := trie.GetKey(seckey); !bytes.Equal(k, key) {
		t.Errorf("GetKey returned %q, want %q", k, key)
	}
}

func TestSecureTrieConcurrency(t *testing.T) {
//为并发访问创建初始trie和copy if
	_, trie, _ := makeTestSecureTrie()

	threads := runtime.NumCPU()
	tries := make([]*SecureTrie, threads)
	for i := 0; i < threads; i++ {
		cpy := *trie
		tries[i] = &cpy
	}
//启动一批Goroutines与Trie交互
	pend := new(sync.WaitGroup)
	pend.Add(threads)
	for i := 0; i < threads; i++ {
		go func(index int) {
			defer pend.Done()

			for j := byte(0); j < 255; j++ {
//在多个键下映射相同的数据
				key, val := common.LeftPadBytes([]byte{byte(index), 1, j}, 32), []byte{j}
				tries[index].Update(key, val)

				key, val = common.LeftPadBytes([]byte{byte(index), 2, j}, 32), []byte{j}
				tries[index].Update(key, val)

//添加一些其他数据来填充trie
				for k := byte(3); k < 13; k++ {
					key, val = common.LeftPadBytes([]byte{byte(index), k, j}, 32), []byte{k, j}
					tries[index].Update(key, val)
				}
			}
			tries[index].Commit(nil)
		}(i)
	}
//等待所有线程完成
	pend.Wait()
}

