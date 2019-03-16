
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450120234045440>


package storage

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	ch "github.com/ethereum/go-ethereum/swarm/chunk"
)

var (
	hashfunc = MakeHashFunc(DefaultHash)
)

//测试内容地址验证器是否正确检查数据
//通过内容地址验证器传递源更新块的测试
//检查资源更新验证器内部正确性的测试在storage/feeds/handler_test.go中找到。
func TestValidator(t *testing.T) {
//设置本地存储
	datadir, err := ioutil.TempDir("", "storage-testvalidator")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(datadir)

	params := NewDefaultLocalStoreParams()
	params.Init(datadir)
	store, err := NewLocalStore(params, nil)
	if err != nil {
		t.Fatal(err)
	}

//不带验证器的检验结果，均成功
	chunks := GenerateRandomChunks(259, 2)
	goodChunk := chunks[0]
	badChunk := chunks[1]
	copy(badChunk.Data(), goodChunk.Data())

	errs := putChunks(store, goodChunk, badChunk)
	if errs[0] != nil {
		t.Fatalf("expected no error on good content address chunk in spite of no validation, but got: %s", err)
	}
	if errs[1] != nil {
		t.Fatalf("expected no error on bad content address chunk in spite of no validation, but got: %s", err)
	}

//添加内容地址验证程序并检查Puts
//坏的应该失败，好的应该通过。
	store.Validators = append(store.Validators, NewContentAddressValidator(hashfunc))
	chunks = GenerateRandomChunks(ch.DefaultSize, 2)
	goodChunk = chunks[0]
	badChunk = chunks[1]
	copy(badChunk.Data(), goodChunk.Data())

	errs = putChunks(store, goodChunk, badChunk)
	if errs[0] != nil {
		t.Fatalf("expected no error on good content address chunk with content address validator only, but got: %s", err)
	}
	if errs[1] == nil {
		t.Fatal("expected error on bad content address chunk with content address validator only, but got nil")
	}

//附加一个始终拒绝的验证器
//坏的应该失败，好的应该通过，
	var negV boolTestValidator
	store.Validators = append(store.Validators, negV)

	chunks = GenerateRandomChunks(ch.DefaultSize, 2)
	goodChunk = chunks[0]
	badChunk = chunks[1]
	copy(badChunk.Data(), goodChunk.Data())

	errs = putChunks(store, goodChunk, badChunk)
	if errs[0] != nil {
		t.Fatalf("expected no error on good content address chunk with content address validator only, but got: %s", err)
	}
	if errs[1] == nil {
		t.Fatal("expected error on bad content address chunk with content address validator only, but got nil")
	}

//附加一个始终批准的验证器
//一切都将通过
	var posV boolTestValidator = true
	store.Validators = append(store.Validators, posV)

	chunks = GenerateRandomChunks(ch.DefaultSize, 2)
	goodChunk = chunks[0]
	badChunk = chunks[1]
	copy(badChunk.Data(), goodChunk.Data())

	errs = putChunks(store, goodChunk, badChunk)
	if errs[0] != nil {
		t.Fatalf("expected no error on good content address chunk with content address validator only, but got: %s", err)
	}
	if errs[1] != nil {
		t.Fatalf("expected no error on bad content address chunk in spite of no validation, but got: %s", err)
	}

}

type boolTestValidator bool

func (self boolTestValidator) Validate(chunk Chunk) bool {
	return bool(self)
}

//PutChunks将块添加到LocalStore
//它等待存储通道上的接收
//它记录但在传递错误时不会失败
func putChunks(store *LocalStore, chunks ...Chunk) []error {
	i := 0
	f := func(n int64) Chunk {
		chunk := chunks[i]
		i++
		return chunk
	}
	_, errs := put(store, len(chunks), f)
	return errs
}

func put(store *LocalStore, n int, f func(i int64) Chunk) (hs []Address, errs []error) {
	for i := int64(0); i < int64(n); i++ {
		chunk := f(ch.DefaultSize)
		err := store.Put(context.TODO(), chunk)
		errs = append(errs, err)
		hs = append(hs, chunk.Address())
	}
	return hs, errs
}

//testgetfrequentlyaccessedchunkwontgetgarbag收集的测试
//频繁访问的块不是从ldbstore收集的垃圾，即，
//当我们达到容量并且垃圾收集器运行时，从磁盘开始。为此
//我们开始将随机块放入数据库，同时不断地访问
//我们关心的块，然后检查我们是否仍然可以从磁盘中检索到它。
func TestGetFrequentlyAccessedChunkWontGetGarbageCollected(t *testing.T) {
	ldbCap := defaultGCRatio
	store, cleanup := setupLocalStore(t, ldbCap)
	defer cleanup()

	var chunks []Chunk
	for i := 0; i < ldbCap; i++ {
		chunks = append(chunks, GenerateRandomChunk(ch.DefaultSize))
	}

	mostAccessed := chunks[0].Address()
	for _, chunk := range chunks {
		if err := store.Put(context.Background(), chunk); err != nil {
			t.Fatal(err)
		}

		if _, err := store.Get(context.Background(), mostAccessed); err != nil {
			t.Fatal(err)
		}
//添加markaccessed（）在单独的goroutine中完成的时间
		time.Sleep(1 * time.Millisecond)
	}

	store.DbStore.collectGarbage()
	if _, err := store.DbStore.Get(context.Background(), mostAccessed); err != nil {
		t.Logf("most frequntly accessed chunk not found on disk (key: %v)", mostAccessed)
		t.Fatal(err)
	}

}

func setupLocalStore(t *testing.T, ldbCap int) (ls *LocalStore, cleanup func()) {
	t.Helper()

	var err error
	datadir, err := ioutil.TempDir("", "storage")
	if err != nil {
		t.Fatal(err)
	}

	params := &LocalStoreParams{
		StoreParams: NewStoreParams(uint64(ldbCap), uint(ldbCap), nil, nil),
	}
	params.Init(datadir)

	store, err := NewLocalStore(params, nil)
	if err != nil {
		_ = os.RemoveAll(datadir)
		t.Fatal(err)
	}

	cleanup = func() {
		store.Close()
		_ = os.RemoveAll(datadir)
	}

	return store, cleanup
}

