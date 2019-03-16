
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450122679324672>


package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

//vmtest检查没有块或事务上下文的evm执行。
//有关测试格式规范，请参阅https://github.com/ethereum/tests/wiki/vm-tests。
type VMTest struct {
	json vmJSON
}

func (t *VMTest) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.json)
}

type vmJSON struct {
	Env           stEnv                 `json:"env"`
	Exec          vmExec                `json:"exec"`
	Logs          common.UnprefixedHash `json:"logs"`
	GasRemaining  *math.HexOrDecimal64  `json:"gas"`
	Out           hexutil.Bytes         `json:"out"`
	Pre           core.GenesisAlloc     `json:"pre"`
	Post          core.GenesisAlloc     `json:"post"`
	PostStateRoot common.Hash           `json:"postStateRoot"`
}

//go:generate gencodec-type vmexec-field override vmexecmarshaling-out gen_vmexec.go

type vmExec struct {
	Address  common.Address `json:"address"  gencodec:"required"`
	Caller   common.Address `json:"caller"   gencodec:"required"`
	Origin   common.Address `json:"origin"   gencodec:"required"`
	Code     []byte         `json:"code"     gencodec:"required"`
	Data     []byte         `json:"data"     gencodec:"required"`
	Value    *big.Int       `json:"value"    gencodec:"required"`
	GasLimit uint64         `json:"gas"      gencodec:"required"`
	GasPrice *big.Int       `json:"gasPrice" gencodec:"required"`
}

type vmExecMarshaling struct {
	Address  common.UnprefixedAddress
	Caller   common.UnprefixedAddress
	Origin   common.UnprefixedAddress
	Code     hexutil.Bytes
	Data     hexutil.Bytes
	Value    *math.HexOrDecimal256
	GasLimit math.HexOrDecimal64
	GasPrice *math.HexOrDecimal256
}

func (t *VMTest) Run(vmconfig vm.Config) error {
	statedb := MakePreState(ethdb.NewMemDatabase(), t.json.Pre)
	ret, gasRemaining, err := t.exec(statedb, vmconfig)

	if t.json.GasRemaining == nil {
		if err == nil {
			return fmt.Errorf("gas unspecified (indicating an error), but VM returned no error")
		}
		if gasRemaining > 0 {
			return fmt.Errorf("gas unspecified (indicating an error), but VM returned gas remaining > 0")
		}
		return nil
	}
//测试声明气体，期望输出匹配。
	if !bytes.Equal(ret, t.json.Out) {
		return fmt.Errorf("return data mismatch: got %x, want %x", ret, t.json.Out)
	}
	if gasRemaining != uint64(*t.json.GasRemaining) {
		return fmt.Errorf("remaining gas %v, want %v", gasRemaining, *t.json.GasRemaining)
	}
	for addr, account := range t.json.Post {
		for k, wantV := range account.Storage {
			if haveV := statedb.GetState(addr, k); haveV != wantV {
				return fmt.Errorf("wrong storage value at %x:\n  got  %x\n  want %x", k, haveV, wantV)
			}
		}
	}
//如果根：=statedb.intermediateroot（false）；根！=t.json.poststateroot
//返回fmt.errorf（“状态后根不匹配，得到了%x，需要%x”，根，t.json.post state root）
//}
	if logs := rlpHash(statedb.Logs()); logs != common.Hash(t.json.Logs) {
		return fmt.Errorf("post state logs hash mismatch: got %x, want %x", logs, t.json.Logs)
	}
	return nil
}

func (t *VMTest) exec(statedb *state.StateDB, vmconfig vm.Config) ([]byte, uint64, error) {
	evm := t.newEVM(statedb, vmconfig)
	e := t.json.Exec
	return evm.Call(vm.AccountRef(e.Caller), e.Address, e.Data, e.GasLimit, e.Value)
}

func (t *VMTest) newEVM(statedb *state.StateDB, vmconfig vm.Config) *vm.EVM {
	initialCall := true
	canTransfer := func(db vm.StateDB, address common.Address, amount *big.Int) bool {
		if initialCall {
			initialCall = false
			return true
		}
		return core.CanTransfer(db, address, amount)
	}
	transfer := func(db vm.StateDB, sender, recipient common.Address, amount *big.Int) {}
	context := vm.Context{
		CanTransfer: canTransfer,
		Transfer:    transfer,
		GetHash:     vmTestBlockHash,
		Origin:      t.json.Exec.Origin,
		Coinbase:    t.json.Env.Coinbase,
		BlockNumber: new(big.Int).SetUint64(t.json.Env.Number),
		Time:        new(big.Int).SetUint64(t.json.Env.Timestamp),
		GasLimit:    t.json.Env.GasLimit,
		Difficulty:  t.json.Env.Difficulty,
		GasPrice:    t.json.Exec.GasPrice,
	}
	vmconfig.NoRecursion = true
	return vm.NewEVM(context, statedb, params.MainnetChainConfig, vmconfig)
}

func vmTestBlockHash(n uint64) common.Hash {
	return common.BytesToHash(crypto.Keccak256([]byte(big.NewInt(int64(n)).String())))
}

