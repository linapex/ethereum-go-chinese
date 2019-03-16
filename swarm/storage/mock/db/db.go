
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450120414400512>


//包DB实现了一个模拟存储，它将所有块数据保存在LevelDB数据库中。
package db

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage/mock"
)

//GlobalStore包含正在存储的LevelDB数据库
//所有群节点的块数据。
//使用关闭方法关闭GlobalStore需要
//释放数据库使用的资源。
type GlobalStore struct {
	db *leveldb.DB
}

//NewGlobalStore创建了一个新的GlobalStore实例。
func NewGlobalStore(path string) (s *GlobalStore, err error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &GlobalStore{
		db: db,
	}, nil
}

//close释放基础级别db使用的资源。
func (s *GlobalStore) Close() error {
	return s.db.Close()
}

//new nodestore返回一个新的nodestore实例，用于检索和存储
//仅对地址为的节点进行数据块处理。
func (s *GlobalStore) NewNodeStore(addr common.Address) *mock.NodeStore {
	return mock.NewNodeStore(addr, s)
}

//如果节点存在键为的块，则get返回块数据
//地址地址。
func (s *GlobalStore) Get(addr common.Address, key []byte) (data []byte, err error) {
	has, err := s.db.Has(nodeDBKey(addr, key), nil)
	if err != nil {
		return nil, mock.ErrNotFound
	}
	if !has {
		return nil, mock.ErrNotFound
	}
	data, err = s.db.Get(dataDBKey(key), nil)
	if err == leveldb.ErrNotFound {
		err = mock.ErrNotFound
	}
	return
}

//Put保存带有地址addr的节点的块数据。
func (s *GlobalStore) Put(addr common.Address, key []byte, data []byte) error {
	batch := new(leveldb.Batch)
	batch.Put(nodeDBKey(addr, key), nil)
	batch.Put(dataDBKey(key), data)
	return s.db.Write(batch, nil)
}

//删除删除对地址为addr的节点的块引用。
func (s *GlobalStore) Delete(addr common.Address, key []byte) error {
	batch := new(leveldb.Batch)
	batch.Delete(nodeDBKey(addr, key))
	return s.db.Write(batch, nil)
}

//haskey返回带有addr的节点是否包含键。
func (s *GlobalStore) HasKey(addr common.Address, key []byte) bool {
	has, err := s.db.Has(nodeDBKey(addr, key), nil)
	if err != nil {
		has = false
	}
	return has
}

//import从包含导出块数据的读卡器读取tar存档。
//它返回导入的块的数量和错误。
func (s *GlobalStore) Import(r io.Reader) (n int, err error) {
	tr := tar.NewReader(r)

	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}
			return n, err
		}

		data, err := ioutil.ReadAll(tr)
		if err != nil {
			return n, err
		}

		var c mock.ExportedChunk
		if err = json.Unmarshal(data, &c); err != nil {
			return n, err
		}

		batch := new(leveldb.Batch)
		for _, addr := range c.Addrs {
			batch.Put(nodeDBKeyHex(addr, hdr.Name), nil)
		}

		batch.Put(dataDBKey(common.Hex2Bytes(hdr.Name)), c.Data)
		if err = s.db.Write(batch, nil); err != nil {
			return n, err
		}

		n++
	}
	return n, err
}

//将包含所有块数据的tar存档导出到写入程序
//商店。它返回导出的块的数量和错误。
func (s *GlobalStore) Export(w io.Writer) (n int, err error) {
	tw := tar.NewWriter(w)
	defer tw.Close()

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := json.NewEncoder(buf)

	iter := s.db.NewIterator(util.BytesPrefix(nodeKeyPrefix), nil)
	defer iter.Release()

	var currentKey string
	var addrs []common.Address

	saveChunk := func(hexKey string) error {
		key := common.Hex2Bytes(hexKey)

		data, err := s.db.Get(dataDBKey(key), nil)
		if err != nil {
			return err
		}

		buf.Reset()
		if err = encoder.Encode(mock.ExportedChunk{
			Addrs: addrs,
			Data:  data,
		}); err != nil {
			return err
		}

		d := buf.Bytes()
		hdr := &tar.Header{
			Name: hexKey,
			Mode: 0644,
			Size: int64(len(d)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if _, err := tw.Write(d); err != nil {
			return err
		}
		n++
		return nil
	}

	for iter.Next() {
		k := bytes.TrimPrefix(iter.Key(), nodeKeyPrefix)
		i := bytes.Index(k, []byte("-"))
		if i < 0 {
			continue
		}
		hexKey := string(k[:i])

		if currentKey == "" {
			currentKey = hexKey
		}

		if hexKey != currentKey {
			if err = saveChunk(currentKey); err != nil {
				return n, err
			}

			addrs = addrs[:0]
		}

		currentKey = hexKey
		addrs = append(addrs, common.BytesToAddress(k[i:]))
	}

	if len(addrs) > 0 {
		if err = saveChunk(currentKey); err != nil {
			return n, err
		}
	}

	return n, err
}

var (
	nodeKeyPrefix = []byte("node-")
	dataKeyPrefix = []byte("data-")
)

//nodedbkey为键/节点映射构造数据库键。
func nodeDBKey(addr common.Address, key []byte) []byte {
	return nodeDBKeyHex(addr, common.Bytes2Hex(key))
}

//nodedbkeyhex为键/节点映射构造数据库键
//使用键的十六进制字符串表示形式。
func nodeDBKeyHex(addr common.Address, hexKey string) []byte {
	return append(append(nodeKeyPrefix, []byte(hexKey+"-")...), addr[:]...)
}

//datadbkey为键/数据存储构造数据库键。
func dataDBKey(key []byte) []byte {
	return append(dataKeyPrefix, key...)
}

