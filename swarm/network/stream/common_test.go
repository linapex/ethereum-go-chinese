
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450115121188864>


package stream

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p/enode"
	p2ptest "github.com/ethereum/go-ethereum/p2p/testing"
	"github.com/ethereum/go-ethereum/swarm/network"
	"github.com/ethereum/go-ethereum/swarm/network/simulation"
	"github.com/ethereum/go-ethereum/swarm/state"
	"github.com/ethereum/go-ethereum/swarm/storage"
	"github.com/ethereum/go-ethereum/swarm/testutil"
	colorable "github.com/mattn/go-colorable"
)

var (
	loglevel     = flag.Int("loglevel", 2, "verbosity of logs")
	nodes        = flag.Int("nodes", 0, "number of nodes")
	chunks       = flag.Int("chunks", 0, "number of chunks")
	useMockStore = flag.Bool("mockstore", false, "disabled mock store (default: enabled)")
	longrunning  = flag.Bool("longrunning", false, "do run long-running tests")

	bucketKeyDB        = simulation.BucketKey("db")
	bucketKeyStore     = simulation.BucketKey("store")
	bucketKeyFileStore = simulation.BucketKey("filestore")
	bucketKeyNetStore  = simulation.BucketKey("netstore")
	bucketKeyDelivery  = simulation.BucketKey("delivery")
	bucketKeyRegistry  = simulation.BucketKey("registry")

	chunkSize = 4096
	pof       = network.Pof
)

func init() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	log.PrintOrigins(true)
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(*loglevel), log.StreamHandler(colorable.NewColorableStderr(), log.TerminalFormat(true))))
}

func newStreamerTester(t *testing.T, registryOptions *RegistryOptions) (*p2ptest.ProtocolTester, *Registry, *storage.LocalStore, func(), error) {
//设置
addr := network.RandomAddr() //测试的对等地址
	to := network.NewKademlia(addr.OAddr, network.NewKadParams())

//临时数据
	datadir, err := ioutil.TempDir("", "streamer")
	if err != nil {
		return nil, nil, nil, func() {}, err
	}
	removeDataDir := func() {
		os.RemoveAll(datadir)
	}

	params := storage.NewDefaultLocalStoreParams()
	params.Init(datadir)
	params.BaseKey = addr.Over()

	localStore, err := storage.NewTestLocalStoreForAddr(params)
	if err != nil {
		return nil, nil, nil, removeDataDir, err
	}

	netStore, err := storage.NewNetStore(localStore, nil)
	if err != nil {
		return nil, nil, nil, removeDataDir, err
	}

	delivery := NewDelivery(to, netStore)
	netStore.NewNetFetcherFunc = network.NewFetcherFactory(delivery.RequestFromPeers, true).New
	streamer := NewRegistry(addr.ID(), delivery, netStore, state.NewInmemoryStore(), registryOptions, nil)
	teardown := func() {
		streamer.Close()
		removeDataDir()
	}
	protocolTester := p2ptest.NewProtocolTester(t, addr.ID(), 1, streamer.runProtocol)

	err = waitForPeers(streamer, 1*time.Second, 1)
	if err != nil {
		return nil, nil, nil, nil, errors.New("timeout: peer is not created")
	}

	return protocolTester, streamer, localStore, teardown, nil
}

func waitForPeers(streamer *Registry, timeout time.Duration, expectedPeers int) error {
	ticker := time.NewTicker(10 * time.Millisecond)
	timeoutTimer := time.NewTimer(timeout)
	for {
		select {
		case <-ticker.C:
			if streamer.peersCount() >= expectedPeers {
				return nil
			}
		case <-timeoutTimer.C:
			return errors.New("timeout")
		}
	}
}

type roundRobinStore struct {
	index  uint32
	stores []storage.ChunkStore
}

func newRoundRobinStore(stores ...storage.ChunkStore) *roundRobinStore {
	return &roundRobinStore{
		stores: stores,
	}
}

func (rrs *roundRobinStore) Get(ctx context.Context, addr storage.Address) (storage.Chunk, error) {
	return nil, errors.New("get not well defined on round robin store")
}

func (rrs *roundRobinStore) Put(ctx context.Context, chunk storage.Chunk) error {
	i := atomic.AddUint32(&rrs.index, 1)
	idx := int(i) % len(rrs.stores)
	return rrs.stores[idx].Put(ctx, chunk)
}

func (rrs *roundRobinStore) Close() {
	for _, store := range rrs.stores {
		store.Close()
	}
}

func readAll(fileStore *storage.FileStore, hash []byte) (int64, error) {
	r, _ := fileStore.Retrieve(context.TODO(), hash)
	buf := make([]byte, 1024)
	var n int
	var total int64
	var err error
	for (total == 0 || n > 0) && err == nil {
		n, err = r.ReadAt(buf, total)
		total += int64(n)
	}
	if err != nil && err != io.EOF {
		return total, err
	}
	return total, nil
}

func uploadFilesToNodes(sim *simulation.Simulation) ([]storage.Address, []string, error) {
	nodes := sim.UpNodeIDs()
	nodeCnt := len(nodes)
	log.Debug(fmt.Sprintf("Uploading %d files to nodes", nodeCnt))
//保存生成文件的数组
	rfiles := make([]string, nodeCnt)
//保存文件根哈希的数组
	rootAddrs := make([]storage.Address, nodeCnt)

	var err error
//对于每个节点，生成一个文件并上载
	for i, id := range nodes {
		item, ok := sim.NodeItem(id, bucketKeyFileStore)
		if !ok {
			return nil, nil, fmt.Errorf("Error accessing localstore")
		}
		fileStore := item.(*storage.FileStore)
//生成文件
		rfiles[i], err = generateRandomFile()
		if err != nil {
			return nil, nil, err
		}
//将其存储（上载）到文件存储区
		ctx := context.TODO()
		rk, wait, err := fileStore.Store(ctx, strings.NewReader(rfiles[i]), int64(len(rfiles[i])), false)
		log.Debug("Uploaded random string file to node")
		if err != nil {
			return nil, nil, err
		}
		err = wait(ctx)
		if err != nil {
			return nil, nil, err
		}
		rootAddrs[i] = rk
	}
	return rootAddrs, rfiles, nil
}

//生成随机文件（字符串）
func generateRandomFile() (string, error) {
//在minfileSize和maxfileSize之间生成随机文件大小
	fileSize := rand.Intn(maxFileSize-minFileSize) + minFileSize
	log.Debug(fmt.Sprintf("Generated file with filesize %d kB", fileSize))
	b := testutil.RandomBytes(1, fileSize*1024)
	return string(b), nil
}

//为给定节点创建本地存储
func createTestLocalStorageForID(id enode.ID, addr *network.BzzAddr) (storage.ChunkStore, string, error) {
	var datadir string
	var err error
	datadir, err = ioutil.TempDir("", fmt.Sprintf("syncer-test-%s", id.TerminalString()))
	if err != nil {
		return nil, "", err
	}
	var store storage.ChunkStore
	params := storage.NewDefaultLocalStoreParams()
	params.ChunkDbPath = datadir
	params.BaseKey = addr.Over()
	store, err = storage.NewTestLocalStoreForAddr(params)
	if err != nil {
		os.RemoveAll(datadir)
		return nil, "", err
	}
	return store, datadir, nil
}

