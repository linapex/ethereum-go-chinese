
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450119462293504>


package feed

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

const (
	testDbDirName = "feeds"
)

type TestHandler struct {
	*Handler
}

func (t *TestHandler) Close() {
	t.chunkStore.Close()
}

type mockNetFetcher struct{}

func (m *mockNetFetcher) Request(ctx context.Context, hopCount uint8) {
}
func (m *mockNetFetcher) Offer(ctx context.Context, source *enode.ID) {
}

func newFakeNetFetcher(context.Context, storage.Address, *sync.Map) storage.NetFetcher {
	return &mockNetFetcher{}
}

//NewTestHandler创建用于测试的处理程序对象。
func NewTestHandler(datadir string, params *HandlerParams) (*TestHandler, error) {
	path := filepath.Join(datadir, testDbDirName)
	fh := NewHandler(params)
	localstoreparams := storage.NewDefaultLocalStoreParams()
	localstoreparams.Init(path)
	localStore, err := storage.NewLocalStore(localstoreparams, nil)
	if err != nil {
		return nil, fmt.Errorf("localstore create fail, path %s: %v", path, err)
	}
	localStore.Validators = append(localStore.Validators, storage.NewContentAddressValidator(storage.MakeHashFunc(feedsHashAlgorithm)))
	localStore.Validators = append(localStore.Validators, fh)
	netStore, err := storage.NewNetStore(localStore, nil)
	if err != nil {
		return nil, err
	}
	netStore.NewNetFetcherFunc = newFakeNetFetcher
	fh.SetStore(netStore)
	return &TestHandler{fh}, nil
}

