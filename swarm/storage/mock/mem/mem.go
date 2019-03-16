
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450120540229632>


//package mem实现了一个模拟存储，将所有块数据保存在内存中。
//虽然它可以用于小规模的测试，但其主要目的是
//包是提供模拟存储的最简单的参考实现。
package mem

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage/mock"
)

//GlobalStore存储所有块数据以及键和节点地址关系。
//它实现mock.globalStore接口。
type GlobalStore struct {
	nodes map[string]map[common.Address]struct{}
	data  map[string][]byte
	mu    sync.Mutex
}

//NewGlobalStore创建了一个新的GlobalStore实例。
func NewGlobalStore() *GlobalStore {
	return &GlobalStore{
		nodes: make(map[string]map[common.Address]struct{}),
		data:  make(map[string][]byte),
	}
}

//new nodestore返回一个新的nodestore实例，用于检索和存储
//仅对地址为的节点进行数据块处理。
func (s *GlobalStore) NewNodeStore(addr common.Address) *mock.NodeStore {
	return mock.NewNodeStore(addr, s)
}

//如果节点存在键为的块，则get返回块数据
//地址地址。
func (s *GlobalStore) Get(addr common.Address, key []byte) (data []byte, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.nodes[string(key)][addr]; !ok {
		return nil, mock.ErrNotFound
	}

	data, ok := s.data[string(key)]
	if !ok {
		return nil, mock.ErrNotFound
	}
	return data, nil
}

//Put保存带有地址addr的节点的块数据。
func (s *GlobalStore) Put(addr common.Address, key []byte, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.nodes[string(key)]; !ok {
		s.nodes[string(key)] = make(map[common.Address]struct{})
	}
	s.nodes[string(key)][addr] = struct{}{}
	s.data[string(key)] = data
	return nil
}

//删除删除地址为的节点的块数据。
func (s *GlobalStore) Delete(addr common.Address, key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var count int
	if _, ok := s.nodes[string(key)]; ok {
		delete(s.nodes[string(key)], addr)
		count = len(s.nodes[string(key)])
	}
	if count == 0 {
		delete(s.data, string(key))
	}
	return nil
}

//haskey返回带有addr的节点是否包含键。
func (s *GlobalStore) HasKey(addr common.Address, key []byte) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.nodes[string(key)][addr]
	return ok
}

//import从包含导出块数据的读卡器读取tar存档。
//它返回导入的块的数量和错误。
func (s *GlobalStore) Import(r io.Reader) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

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

		addrs := make(map[common.Address]struct{})
		for _, a := range c.Addrs {
			addrs[a] = struct{}{}
		}

		key := string(common.Hex2Bytes(hdr.Name))
		s.nodes[key] = addrs
		s.data[key] = c.Data
		n++
	}
	return n, err
}

//将包含所有块数据的tar存档导出到写入程序
//商店。它返回导出的块数和错误。
func (s *GlobalStore) Export(w io.Writer) (n int, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tw := tar.NewWriter(w)
	defer tw.Close()

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	encoder := json.NewEncoder(buf)
	for key, addrs := range s.nodes {
		al := make([]common.Address, 0, len(addrs))
		for a := range addrs {
			al = append(al, a)
		}

		buf.Reset()
		if err = encoder.Encode(mock.ExportedChunk{
			Addrs: al,
			Data:  s.data[key],
		}); err != nil {
			return n, err
		}

		data := buf.Bytes()
		hdr := &tar.Header{
			Name: common.Bytes2Hex([]byte(key)),
			Mode: 0644,
			Size: int64(len(data)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return n, err
		}
		if _, err := tw.Write(data); err != nil {
			return n, err
		}
		n++
	}
	return n, err
}

