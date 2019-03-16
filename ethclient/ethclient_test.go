
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450090219606016>


package ethclient

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

//验证客户端是否实现了以太坊接口。
var (
	_ = ethereum.ChainReader(&Client{})
	_ = ethereum.TransactionReader(&Client{})
	_ = ethereum.ChainStateReader(&Client{})
	_ = ethereum.ChainSyncReader(&Client{})
	_ = ethereum.ContractCaller(&Client{})
	_ = ethereum.GasEstimator(&Client{})
	_ = ethereum.GasPricer(&Client{})
	_ = ethereum.LogFilterer(&Client{})
	_ = ethereum.PendingStateReader(&Client{})
//_u=ethereum.pendingstateeventer（&client）
	_ = ethereum.PendingContractCaller(&Client{})
)

func TestToFilterArg(t *testing.T) {
	blockHashErr := fmt.Errorf("cannot specify both BlockHash and FromBlock/ToBlock")
	addresses := []common.Address{
		common.HexToAddress("0xD36722ADeC3EdCB29c8e7b5a47f352D701393462"),
	}
	blockHash := common.HexToHash(
		"0xeb94bb7d78b73657a9d7a99792413f50c0a45c51fc62bdcb08a53f18e9a2b4eb",
	)

	for _, testCase := range []struct {
		name   string
		input  ethereum.FilterQuery
		output interface{}
		err    error
	}{
		{
			"without BlockHash",
			ethereum.FilterQuery{
				Addresses: addresses,
				FromBlock: big.NewInt(1),
				ToBlock:   big.NewInt(2),
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"fromBlock": "0x1",
				"toBlock":   "0x2",
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with nil fromBlock and nil toBlock",
			ethereum.FilterQuery{
				Addresses: addresses,
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"fromBlock": "0x0",
				"toBlock":   "latest",
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with blockhash",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				Topics:    [][]common.Hash{},
			},
			map[string]interface{}{
				"address":   addresses,
				"blockHash": blockHash,
				"topics":    [][]common.Hash{},
			},
			nil,
		},
		{
			"with blockhash and from block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				FromBlock: big.NewInt(1),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
		{
			"with blockhash and to block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				ToBlock:   big.NewInt(1),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
		{
			"with blockhash and both from / to block",
			ethereum.FilterQuery{
				Addresses: addresses,
				BlockHash: &blockHash,
				FromBlock: big.NewInt(1),
				ToBlock:   big.NewInt(2),
				Topics:    [][]common.Hash{},
			},
			nil,
			blockHashErr,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			output, err := toFilterArg(testCase.input)
			if (testCase.err == nil) != (err == nil) {
				t.Fatalf("expected error %v but got %v", testCase.err, err)
			}
			if testCase.err != nil {
				if testCase.err.Error() != err.Error() {
					t.Fatalf("expected error %v but got %v", testCase.err, err)
				}
			} else if !reflect.DeepEqual(testCase.output, output) {
				t.Fatalf("expected filter arg %v but got %v", testCase.output, output)
			}
		})
	}
}

