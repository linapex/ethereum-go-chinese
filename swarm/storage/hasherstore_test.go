
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450119965609984>


package storage

import (
	"bytes"
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/swarm/storage/encryption"

	"github.com/ethereum/go-ethereum/common"
)

func TestHasherStore(t *testing.T) {
	var tests = []struct {
		chunkLength int
		toEncrypt   bool
	}{
		{10, false},
		{100, false},
		{1000, false},
		{4096, false},
		{10, true},
		{100, true},
		{1000, true},
		{4096, true},
	}

	for _, tt := range tests {
		chunkStore := NewMapChunkStore()
		hasherStore := NewHasherStore(chunkStore, MakeHashFunc(DefaultHash), tt.toEncrypt)

//将两个随机块放入哈希存储中
		chunkData1 := GenerateRandomChunk(int64(tt.chunkLength)).Data()
		ctx, cancel := context.WithTimeout(context.Background(), getTimeout)
		defer cancel()
		key1, err := hasherStore.Put(ctx, chunkData1)
		if err != nil {
			t.Fatalf("Expected no error got \"%v\"", err)
		}

		chunkData2 := GenerateRandomChunk(int64(tt.chunkLength)).Data()
		key2, err := hasherStore.Put(ctx, chunkData2)
		if err != nil {
			t.Fatalf("Expected no error got \"%v\"", err)
		}

		hasherStore.Close()

//等待块真正存储
		err = hasherStore.Wait(ctx)
		if err != nil {
			t.Fatalf("Expected no error got \"%v\"", err)
		}

//获取第一块
		retrievedChunkData1, err := hasherStore.Get(ctx, key1)
		if err != nil {
			t.Fatalf("Expected no error, got \"%v\"", err)
		}

//检索到的数据应与原始数据相同
		if !bytes.Equal(chunkData1, retrievedChunkData1) {
			t.Fatalf("Expected retrieved chunk data %v, got %v", common.Bytes2Hex(chunkData1), common.Bytes2Hex(retrievedChunkData1))
		}

//获取第二块
		retrievedChunkData2, err := hasherStore.Get(ctx, key2)
		if err != nil {
			t.Fatalf("Expected no error, got \"%v\"", err)
		}

//检索到的数据应与原始数据相同
		if !bytes.Equal(chunkData2, retrievedChunkData2) {
			t.Fatalf("Expected retrieved chunk data %v, got %v", common.Bytes2Hex(chunkData2), common.Bytes2Hex(retrievedChunkData2))
		}

		hash1, encryptionKey1, err := parseReference(key1, hasherStore.hashSize)
		if err != nil {
			t.Fatalf("Expected no error, got \"%v\"", err)
		}

		if tt.toEncrypt {
			if encryptionKey1 == nil {
				t.Fatal("Expected non-nil encryption key, got nil")
			} else if len(encryptionKey1) != encryption.KeyLength {
				t.Fatalf("Expected encryption key length %v, got %v", encryption.KeyLength, len(encryptionKey1))
			}
		}
		if !tt.toEncrypt && encryptionKey1 != nil {
			t.Fatalf("Expected nil encryption key, got key with length %v", len(encryptionKey1))
		}

//检查存储区中的块数据是否加密
		chunkInStore, err := chunkStore.Get(ctx, hash1)
		if err != nil {
			t.Fatalf("Expected no error got \"%v\"", err)
		}

		chunkDataInStore := chunkInStore.Data()

		if tt.toEncrypt && bytes.Equal(chunkData1, chunkDataInStore) {
			t.Fatalf("Chunk expected to be encrypted but it is stored without encryption")
		}
		if !tt.toEncrypt && !bytes.Equal(chunkData1, chunkDataInStore) {
			t.Fatalf("Chunk expected to be not encrypted but stored content is different. Expected %v got %v", common.Bytes2Hex(chunkData1), common.Bytes2Hex(chunkDataInStore))
		}
	}
}

