
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450111509893120>


package api

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestConfig(t *testing.T) {

	var hexprvkey = "65138b2aa745041b372153550584587da326ab440576b2a1191dd95cee30039c"

	prvkey, err := crypto.HexToECDSA(hexprvkey)
	if err != nil {
		t.Fatalf("failed to load private key: %v", err)
	}

	one := NewConfig()
	two := NewConfig()

	one.LocalStoreParams = two.LocalStoreParams
	if equal := reflect.DeepEqual(one, two); !equal {
		t.Fatal("Two default configs are not equal")
	}

	one.Init(prvkey)

//init函数应设置以下字段
	if one.BzzKey == "" {
		t.Fatal("Expected BzzKey to be set")
	}
	if one.PublicKey == "" {
		t.Fatal("Expected PublicKey to be set")
	}
	if one.Swap.PayProfile.Beneficiary == (common.Address{}) && one.SwapEnabled {
		t.Fatal("Failed to correctly initialize SwapParams")
	}
	if one.ChunkDbPath == one.Path {
		t.Fatal("Failed to correctly initialize StoreParams")
	}
}

