
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450077045297152>


package ens

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/ens/contract"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	key, _   = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	name     = "my name on ENS"
	hash     = crypto.Keccak256Hash([]byte("my content"))
	addr     = crypto.PubkeyToAddress(key.PublicKey)
	testAddr = common.HexToAddress("0x1234123412341234123412341234123412341234")
)

func TestENS(t *testing.T) {
	contractBackend := backends.NewSimulatedBackend(core.GenesisAlloc{addr: {Balance: big.NewInt(1000000000)}}, 10000000)
	transactOpts := bind.NewKeyedTransactor(key)

	ensAddr, ens, err := DeployENS(transactOpts, contractBackend)
	if err != nil {
		t.Fatalf("can't deploy root registry: %v", err)
	}
	contractBackend.Commit()

//将自己设置为名称的所有者。
	if _, err := ens.Register(name); err != nil {
		t.Fatalf("can't register: %v", err)
	}
	contractBackend.Commit()

//部署一个冲突解决程序并使其负责该名称。
	resolverAddr, _, _, err := contract.DeployPublicResolver(transactOpts, contractBackend, ensAddr)
	if err != nil {
		t.Fatalf("can't deploy resolver: %v", err)
	}
	if _, err := ens.SetResolver(EnsNode(name), resolverAddr); err != nil {
		t.Fatalf("can't set resolver: %v", err)
	}
	contractBackend.Commit()

//设置名称的内容哈希。
	if _, err = ens.SetContentHash(name, hash); err != nil {
		t.Fatalf("can't set content hash: %v", err)
	}
	contractBackend.Commit()

//尝试解析名称。
	vhost, err := ens.Resolve(name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if vhost != hash {
		t.Fatalf("resolve error, expected %v, got %v", hash.Hex(), vhost.Hex())
	}

//设置名称的地址
	if _, err = ens.SetAddr(name, testAddr); err != nil {
		t.Fatalf("can't set address: %v", err)
	}
	contractBackend.Commit()

//尝试将名称解析为地址
	recoveredAddr, err := ens.Addr(name)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if vhost != hash {
		t.Fatalf("resolve error, expected %v, got %v", testAddr.Hex(), recoveredAddr.Hex())
	}
}

