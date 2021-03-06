
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450121165180928>


//包测试提供用于测试的函数
//GlobalStrer实施。
package test

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"github.com/ethereum/go-ethereum/swarm/storage/mock"
)

//mockstore从提供的globalstorer创建nodestore实例，
//每个都有一个唯一的地址，在上面存储不同的块
//并检查它们是否可在所有节点上检索。
//属性n定义将创建的节点存储数。
func MockStore(t *testing.T, globalStore mock.GlobalStorer, n int) {
	t.Run("GlobalStore", func(t *testing.T) {
		addrs := make([]common.Address, n)
		for i := 0; i < n; i++ {
			addrs[i] = common.HexToAddress(strconv.FormatInt(int64(i)+1, 16))
		}

		for i, addr := range addrs {
			chunkAddr := storage.Address(append(addr[:], []byte(strconv.FormatInt(int64(i)+1, 16))...))
			data := []byte(strconv.FormatInt(int64(i)+1, 16))
			data = append(data, make([]byte, 4096-len(data))...)
			globalStore.Put(addr, chunkAddr, data)

			for _, cAddr := range addrs {
				cData, err := globalStore.Get(cAddr, chunkAddr)
				if cAddr == addr {
					if err != nil {
						t.Fatalf("get data from store %s key %s: %v", cAddr.Hex(), chunkAddr.Hex(), err)
					}
					if !bytes.Equal(data, cData) {
						t.Fatalf("data on store %s: expected %x, got %x", cAddr.Hex(), data, cData)
					}
					if !globalStore.HasKey(cAddr, chunkAddr) {
						t.Fatalf("expected key %s on global store for node %s, but it was not found", chunkAddr.Hex(), cAddr.Hex())
					}
				} else {
					if err != mock.ErrNotFound {
						t.Fatalf("expected error from store %s: %v, got %v", cAddr.Hex(), mock.ErrNotFound, err)
					}
					if len(cData) > 0 {
						t.Fatalf("data on store %s: expected nil, got %x", cAddr.Hex(), cData)
					}
					if globalStore.HasKey(cAddr, chunkAddr) {
						t.Fatalf("not expected key %s on global store for node %s, but it was found", chunkAddr.Hex(), cAddr.Hex())
					}
				}
			}
		}
		t.Run("delete", func(t *testing.T) {
			chunkAddr := storage.Address([]byte("1234567890abcd"))
			for _, addr := range addrs {
				err := globalStore.Put(addr, chunkAddr, []byte("data"))
				if err != nil {
					t.Fatalf("put data to store %s key %s: %v", addr.Hex(), chunkAddr.Hex(), err)
				}
			}
			firstNodeAddr := addrs[0]
			if err := globalStore.Delete(firstNodeAddr, chunkAddr); err != nil {
				t.Fatalf("delete from store %s key %s: %v", firstNodeAddr.Hex(), chunkAddr.Hex(), err)
			}
			for i, addr := range addrs {
				_, err := globalStore.Get(addr, chunkAddr)
				if i == 0 {
					if err != mock.ErrNotFound {
						t.Errorf("get data from store %s key %s: expected mock.ErrNotFound error, got %v", addr.Hex(), chunkAddr.Hex(), err)
					}
				} else {
					if err != nil {
						t.Errorf("get data from store %s key %s: %v", addr.Hex(), chunkAddr.Hex(), err)
					}
				}
			}
		})
	})

	t.Run("NodeStore", func(t *testing.T) {
		nodes := make(map[common.Address]*mock.NodeStore)
		for i := 0; i < n; i++ {
			addr := common.HexToAddress(strconv.FormatInt(int64(i)+1, 16))
			nodes[addr] = globalStore.NewNodeStore(addr)
		}

		i := 0
		for addr, store := range nodes {
			i++
			chunkAddr := storage.Address(append(addr[:], []byte(fmt.Sprintf("%x", i))...))
			data := []byte(strconv.FormatInt(int64(i)+1, 16))
			data = append(data, make([]byte, 4096-len(data))...)
			store.Put(chunkAddr, data)

			for cAddr, cStore := range nodes {
				cData, err := cStore.Get(chunkAddr)
				if cAddr == addr {
					if err != nil {
						t.Fatalf("get data from store %s key %s: %v", cAddr.Hex(), chunkAddr.Hex(), err)
					}
					if !bytes.Equal(data, cData) {
						t.Fatalf("data on store %s: expected %x, got %x", cAddr.Hex(), data, cData)
					}
					if !globalStore.HasKey(cAddr, chunkAddr) {
						t.Fatalf("expected key %s on global store for node %s, but it was not found", chunkAddr.Hex(), cAddr.Hex())
					}
				} else {
					if err != mock.ErrNotFound {
						t.Fatalf("expected error from store %s: %v, got %v", cAddr.Hex(), mock.ErrNotFound, err)
					}
					if len(cData) > 0 {
						t.Fatalf("data on store %s: expected nil, got %x", cAddr.Hex(), cData)
					}
					if globalStore.HasKey(cAddr, chunkAddr) {
						t.Fatalf("not expected key %s on global store for node %s, but it was found", chunkAddr.Hex(), cAddr.Hex())
					}
				}
			}
		}
		t.Run("delete", func(t *testing.T) {
			chunkAddr := storage.Address([]byte("1234567890abcd"))
			var chosenStore *mock.NodeStore
			for addr, store := range nodes {
				if chosenStore == nil {
					chosenStore = store
				}
				err := store.Put(chunkAddr, []byte("data"))
				if err != nil {
					t.Fatalf("put data to store %s key %s: %v", addr.Hex(), chunkAddr.Hex(), err)
				}
			}
			if err := chosenStore.Delete(chunkAddr); err != nil {
				t.Fatalf("delete key %s: %v", chunkAddr.Hex(), err)
			}
			for addr, store := range nodes {
				_, err := store.Get(chunkAddr)
				if store == chosenStore {
					if err != mock.ErrNotFound {
						t.Errorf("get data from store %s key %s: expected mock.ErrNotFound error, got %v", addr.Hex(), chunkAddr.Hex(), err)
					}
				} else {
					if err != nil {
						t.Errorf("get data from store %s key %s: %v", addr.Hex(), chunkAddr.Hex(), err)
					}
				}
			}
		})
	})
}

//importexport将块保存到出口，将它们导出到tar存档，
//将tar存档导入到instore并检查是否正确导入了所有块。
func ImportExport(t *testing.T, outStore, inStore mock.GlobalStorer, n int) {
	exporter, ok := outStore.(mock.Exporter)
	if !ok {
		t.Fatal("outStore does not implement mock.Exporter")
	}
	importer, ok := inStore.(mock.Importer)
	if !ok {
		t.Fatal("inStore does not implement mock.Importer")
	}
	addrs := make([]common.Address, n)
	for i := 0; i < n; i++ {
		addrs[i] = common.HexToAddress(strconv.FormatInt(int64(i)+1, 16))
	}

	for i, addr := range addrs {
		chunkAddr := storage.Address(append(addr[:], []byte(strconv.FormatInt(int64(i)+1, 16))...))
		data := []byte(strconv.FormatInt(int64(i)+1, 16))
		data = append(data, make([]byte, 4096-len(data))...)
		outStore.Put(addr, chunkAddr, data)
	}

	r, w := io.Pipe()
	defer r.Close()

	exportErrChan := make(chan error)
	go func() {
		defer w.Close()

		_, err := exporter.Export(w)
		exportErrChan <- err
	}()

	if _, err := importer.Import(r); err != nil {
		t.Fatalf("import: %v", err)
	}

	if err := <-exportErrChan; err != nil {
		t.Fatalf("export: %v", err)
	}

	for i, addr := range addrs {
		chunkAddr := storage.Address(append(addr[:], []byte(strconv.FormatInt(int64(i)+1, 16))...))
		data := []byte(strconv.FormatInt(int64(i)+1, 16))
		data = append(data, make([]byte, 4096-len(data))...)
		for _, cAddr := range addrs {
			cData, err := inStore.Get(cAddr, chunkAddr)
			if cAddr == addr {
				if err != nil {
					t.Fatalf("get data from store %s key %s: %v", cAddr.Hex(), chunkAddr.Hex(), err)
				}
				if !bytes.Equal(data, cData) {
					t.Fatalf("data on store %s: expected %x, got %x", cAddr.Hex(), data, cData)
				}
				if !inStore.HasKey(cAddr, chunkAddr) {
					t.Fatalf("expected key %s on global store for node %s, but it was not found", chunkAddr.Hex(), cAddr.Hex())
				}
			} else {
				if err != mock.ErrNotFound {
					t.Fatalf("expected error from store %s: %v, got %v", cAddr.Hex(), mock.ErrNotFound, err)
				}
				if len(cData) > 0 {
					t.Fatalf("data on store %s: expected nil, got %x", cAddr.Hex(), cData)
				}
				if inStore.HasKey(cAddr, chunkAddr) {
					t.Fatalf("not expected key %s on global store for node %s, but it was found", chunkAddr.Hex(), cAddr.Hex())
				}
			}
		}
	}
}

