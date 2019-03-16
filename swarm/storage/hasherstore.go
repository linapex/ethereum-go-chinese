
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450119906889728>


package storage

import (
	"context"
	"fmt"
	"sync/atomic"

	ch "github.com/ethereum/go-ethereum/swarm/chunk"
	"github.com/ethereum/go-ethereum/swarm/storage/encryption"
	"golang.org/x/crypto/sha3"
)

type hasherStore struct {
	store     ChunkStore
	toEncrypt bool
	hashFunc  SwarmHasher
hashSize  int           //内容哈希大小
refSize   int64         //引用大小（内容哈希+可能的加密密钥）
errC      chan error    //全局错误通道
doneC     chan struct{} //通过close（）调用关闭，以指示count是最后的块数
quitC     chan struct{} //关闭以退出未终止的例程
//nrchunks与原子函数一起使用
//它必须位于结构的末尾，以确保ARM体系结构的64位对齐。
//参见：https://golang.org/pkg/sync/atomic/pkg note bug
nrChunks uint64 //要存储的块数
}

//NewHasherStore创建了一个HasherStore对象，它实现了推杆和getter接口。
//使用hasherstore，您可以将区块数据（仅为[]字节）放入chunkstore中并获取它们。
//如果需要，哈希存储将获取数据加密/解密的核心。
func NewHasherStore(store ChunkStore, hashFunc SwarmHasher, toEncrypt bool) *hasherStore {
	hashSize := hashFunc().Size()
	refSize := int64(hashSize)
	if toEncrypt {
		refSize += encryption.KeyLength
	}

	h := &hasherStore{
		store:     store,
		toEncrypt: toEncrypt,
		hashFunc:  hashFunc,
		hashSize:  hashSize,
		refSize:   refSize,
		errC:      make(chan error),
		doneC:     make(chan struct{}),
		quitC:     make(chan struct{}),
	}

	return h
}

//将chunkdata存储到哈希存储的chunkstore中并返回引用。
//如果hasherstore具有chunkEncryption对象，则数据将被加密。
//异步函数，返回时不必存储数据。
func (h *hasherStore) Put(ctx context.Context, chunkData ChunkData) (Reference, error) {
	c := chunkData
	var encryptionKey encryption.Key
	if h.toEncrypt {
		var err error
		c, encryptionKey, err = h.encryptChunkData(chunkData)
		if err != nil {
			return nil, err
		}
	}
	chunk := h.createChunk(c)
	h.storeChunk(ctx, chunk)

	return Reference(append(chunk.Address(), encryptionKey...)), nil
}

//get返回具有给定引用的块的数据（从HasherStore的ChunkStore中检索）。
//如果数据已加密，并且引用包含加密密钥，则将在
//返回。
func (h *hasherStore) Get(ctx context.Context, ref Reference) (ChunkData, error) {
	addr, encryptionKey, err := parseReference(ref, h.hashSize)
	if err != nil {
		return nil, err
	}

	chunk, err := h.store.Get(ctx, addr)
	if err != nil {
		return nil, err
	}

	chunkData := ChunkData(chunk.Data())
	toDecrypt := (encryptionKey != nil)
	if toDecrypt {
		var err error
		chunkData, err = h.decryptChunkData(chunkData, encryptionKey)
		if err != nil {
			return nil, err
		}
	}
	return chunkData, nil
}

//CLOSE表示不再将块与哈希存储放在一起，因此等待
//函数可以在存储所有以前放置的块后返回。
func (h *hasherStore) Close() {
	close(h.doneC)
}

//等待返回时间
//1）已调用close（）函数，并且
//2）已放置的所有块已存储
func (h *hasherStore) Wait(ctx context.Context) error {
	defer close(h.quitC)
var nrStoredChunks uint64 //存储的块数
	var done bool
	doneC := h.doneC
	for {
		select {
//如果上下文在前面完成，只需返回并返回错误
		case <-ctx.Done():
			return ctx.Err()
//如果所有块都已提交，则DONEC将关闭，从那时起，我们只需等待所有块也被存储。
		case <-doneC:
			done = true
			doneC = nil
//已存储块，如果err为nil，则成功，因此请增加存储块计数器
		case err := <-h.errC:
			if err != nil {
				return err
			}
			nrStoredChunks++
		}
//如果所有的块都已提交，并且所有的块都已存储，那么我们可以返回
		if done {
			if nrStoredChunks >= atomic.LoadUint64(&h.nrChunks) {
				return nil
			}
		}
	}
}

func (h *hasherStore) createHash(chunkData ChunkData) Address {
	hasher := h.hashFunc()
hasher.ResetWithLength(chunkData[:8]) //8字节长度
hasher.Write(chunkData[8:])           //减去8[]字节长度
	return hasher.Sum(nil)
}

func (h *hasherStore) createChunk(chunkData ChunkData) *chunk {
	hash := h.createHash(chunkData)
	chunk := NewChunk(hash, chunkData)
	return chunk
}

func (h *hasherStore) encryptChunkData(chunkData ChunkData) (ChunkData, encryption.Key, error) {
	if len(chunkData) < 8 {
		return nil, nil, fmt.Errorf("Invalid ChunkData, min length 8 got %v", len(chunkData))
	}

	key, encryptedSpan, encryptedData, err := h.encrypt(chunkData)
	if err != nil {
		return nil, nil, err
	}
	c := make(ChunkData, len(encryptedSpan)+len(encryptedData))
	copy(c[:8], encryptedSpan)
	copy(c[8:], encryptedData)
	return c, key, nil
}

func (h *hasherStore) decryptChunkData(chunkData ChunkData, encryptionKey encryption.Key) (ChunkData, error) {
	if len(chunkData) < 8 {
		return nil, fmt.Errorf("Invalid ChunkData, min length 8 got %v", len(chunkData))
	}

	decryptedSpan, decryptedData, err := h.decrypt(chunkData, encryptionKey)
	if err != nil {
		return nil, err
	}

//删除刚添加用于填充的多余字节
	length := ChunkData(decryptedSpan).Size()
	for length > ch.DefaultSize {
		length = length + (ch.DefaultSize - 1)
		length = length / ch.DefaultSize
		length *= uint64(h.refSize)
	}

	c := make(ChunkData, length+8)
	copy(c[:8], decryptedSpan)
	copy(c[8:], decryptedData[:length])

	return c, nil
}

func (h *hasherStore) RefSize() int64 {
	return h.refSize
}

func (h *hasherStore) encrypt(chunkData ChunkData) (encryption.Key, []byte, []byte, error) {
	key := encryption.GenerateRandomKey(encryption.KeyLength)
	encryptedSpan, err := h.newSpanEncryption(key).Encrypt(chunkData[:8])
	if err != nil {
		return nil, nil, nil, err
	}
	encryptedData, err := h.newDataEncryption(key).Encrypt(chunkData[8:])
	if err != nil {
		return nil, nil, nil, err
	}
	return key, encryptedSpan, encryptedData, nil
}

func (h *hasherStore) decrypt(chunkData ChunkData, key encryption.Key) ([]byte, []byte, error) {
	encryptedSpan, err := h.newSpanEncryption(key).Encrypt(chunkData[:8])
	if err != nil {
		return nil, nil, err
	}
	encryptedData, err := h.newDataEncryption(key).Encrypt(chunkData[8:])
	if err != nil {
		return nil, nil, err
	}
	return encryptedSpan, encryptedData, nil
}

func (h *hasherStore) newSpanEncryption(key encryption.Key) encryption.Encryption {
	return encryption.New(key, 0, uint32(ch.DefaultSize/h.refSize), sha3.NewLegacyKeccak256)
}

func (h *hasherStore) newDataEncryption(key encryption.Key) encryption.Encryption {
	return encryption.New(key, int(ch.DefaultSize), 0, sha3.NewLegacyKeccak256)
}

func (h *hasherStore) storeChunk(ctx context.Context, chunk *chunk) {
	atomic.AddUint64(&h.nrChunks, 1)
	go func() {
		select {
		case h.errC <- h.store.Put(ctx, chunk):
		case <-h.quitC:
		}
	}()
}

func parseReference(ref Reference, hashSize int) (Address, encryption.Key, error) {
	encryptedRefLength := hashSize + encryption.KeyLength
	switch len(ref) {
	case AddressLength:
		return Address(ref), nil, nil
	case encryptedRefLength:
		encKeyIdx := len(ref) - encryption.KeyLength
		return Address(ref[:encKeyIdx]), encryption.Key(ref[encKeyIdx:]), nil
	default:
		return nil, nil, fmt.Errorf("Invalid reference length, expected %v or %v got %v", hashSize, encryptedRefLength, len(ref))
	}
}

