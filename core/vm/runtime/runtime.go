
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450083227701248>


package runtime

import (
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

//config是一种基本类型，指定运行的某些配置标志
//EVM。
type Config struct {
	ChainConfig *params.ChainConfig
	Difficulty  *big.Int
	Origin      common.Address
	Coinbase    common.Address
	BlockNumber *big.Int
	Time        *big.Int
	GasLimit    uint64
	GasPrice    *big.Int
	Value       *big.Int
	Debug       bool
	EVMConfig   vm.Config

	State     *state.StateDB
	GetHashFn func(n uint64) common.Hash
}

//设置配置的默认值
func setDefaults(cfg *Config) {
	if cfg.ChainConfig == nil {
		cfg.ChainConfig = &params.ChainConfig{
			ChainID:        big.NewInt(1),
			HomesteadBlock: new(big.Int),
			DAOForkBlock:   new(big.Int),
			DAOForkSupport: false,
			EIP150Block:    new(big.Int),
			EIP155Block:    new(big.Int),
			EIP158Block:    new(big.Int),
		}
	}

	if cfg.Difficulty == nil {
		cfg.Difficulty = new(big.Int)
	}
	if cfg.Time == nil {
		cfg.Time = big.NewInt(time.Now().Unix())
	}
	if cfg.GasLimit == 0 {
		cfg.GasLimit = math.MaxUint64
	}
	if cfg.GasPrice == nil {
		cfg.GasPrice = new(big.Int)
	}
	if cfg.Value == nil {
		cfg.Value = new(big.Int)
	}
	if cfg.BlockNumber == nil {
		cfg.BlockNumber = new(big.Int)
	}
	if cfg.GetHashFn == nil {
		cfg.GetHashFn = func(n uint64) common.Hash {
			return common.BytesToHash(crypto.Keccak256([]byte(new(big.Int).SetUint64(n).String())))
		}
	}
}

//在执行期间，Execute使用输入作为调用数据来执行代码。
//它返回EVM的返回值、新状态以及失败时的错误。
//
//执行为执行设置内存中的临时环境
//给定的代码。它确保之后它恢复到原来的状态。
func Execute(code, input []byte, cfg *Config) ([]byte, *state.StateDB, error) {
	if cfg == nil {
		cfg = new(Config)
	}
	setDefaults(cfg)

	if cfg.State == nil {
		cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	}
	var (
		address = common.BytesToAddress([]byte("contract"))
		vmenv   = NewEnv(cfg)
		sender  = vm.AccountRef(cfg.Origin)
	)
	cfg.State.CreateAccount(address)
//设置要执行的接收方（执行合同）代码。
	cfg.State.SetCode(address, code)
//使用给定的配置调用代码。
	ret, _, err := vmenv.Call(
		sender,
		common.BytesToAddress([]byte("contract")),
		input,
		cfg.GasLimit,
		cfg.Value,
	)

	return ret, cfg.State, err
}

//create使用evm create方法执行代码
func Create(input []byte, cfg *Config) ([]byte, common.Address, uint64, error) {
	if cfg == nil {
		cfg = new(Config)
	}
	setDefaults(cfg)

	if cfg.State == nil {
		cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(ethdb.NewMemDatabase()))
	}
	var (
		vmenv  = NewEnv(cfg)
		sender = vm.AccountRef(cfg.Origin)
	)

//使用给定的配置调用代码。
	code, address, leftOverGas, err := vmenv.Create(
		sender,
		input,
		cfg.GasLimit,
		cfg.Value,
	)
	return code, address, leftOverGas, err
}

//调用执行由合同地址给定的代码。它将返回
//EVM的返回值或失败时的错误。
//
//与execute不同，call需要config，还需要state字段
//被设定。
func Call(address common.Address, input []byte, cfg *Config) ([]byte, uint64, error) {
	setDefaults(cfg)

	vmenv := NewEnv(cfg)

	sender := cfg.State.GetOrNewStateObject(cfg.Origin)
//使用给定的配置调用代码。
	ret, leftOverGas, err := vmenv.Call(
		sender,
		address,
		input,
		cfg.GasLimit,
		cfg.Value,
	)

	return ret, leftOverGas, err
}

