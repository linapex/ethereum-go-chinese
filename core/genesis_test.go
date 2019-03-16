
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450078928539648>


package core

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

func TestDefaultGenesisBlock(t *testing.T) {
	block := DefaultGenesisBlock().ToBlock(nil)
	if block.Hash() != params.MainnetGenesisHash {
		t.Errorf("wrong mainnet genesis hash, got %v, want %v", block.Hash(), params.MainnetGenesisHash)
	}
	block = DefaultTestnetGenesisBlock().ToBlock(nil)
	if block.Hash() != params.TestnetGenesisHash {
		t.Errorf("wrong testnet genesis hash, got %v, want %v", block.Hash(), params.TestnetGenesisHash)
	}
}

func TestSetupGenesis(t *testing.T) {
	var (
		customghash = common.HexToHash("0x89c99d90b79719238d2645c7642f2c9295246e80775b38cfd162b696817fbd50")
		customg     = Genesis{
			Config: &params.ChainConfig{HomesteadBlock: big.NewInt(3)},
			Alloc: GenesisAlloc{
				{1}: {Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}}},
			},
		}
		oldcustomg = customg
	)
	oldcustomg.Config = &params.ChainConfig{HomesteadBlock: big.NewInt(2)}
	tests := []struct {
		name       string
		fn         func(ethdb.Database) (*params.ChainConfig, common.Hash, error)
		wantConfig *params.ChainConfig
		wantHash   common.Hash
		wantErr    error
	}{
		{
			name: "genesis without ChainConfig",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, new(Genesis))
			},
			wantErr:    errGenesisNoConfig,
			wantConfig: params.AllEthashProtocolChanges,
		},
		{
			name: "no block in DB, genesis == nil",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, nil)
			},
			wantHash:   params.MainnetGenesisHash,
			wantConfig: params.MainnetChainConfig,
		},
		{
			name: "mainnet block in DB, genesis == nil",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				DefaultGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, nil)
			},
			wantHash:   params.MainnetGenesisHash,
			wantConfig: params.MainnetChainConfig,
		},
		{
			name: "custom block in DB, genesis == nil",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				customg.MustCommit(db)
				return SetupGenesisBlock(db, nil)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
		},
		{
			name: "custom block in DB, genesis == testnet",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				customg.MustCommit(db)
				return SetupGenesisBlock(db, DefaultTestnetGenesisBlock())
			},
			wantErr:    &GenesisMismatchError{Stored: customghash, New: params.TestnetGenesisHash},
			wantHash:   params.TestnetGenesisHash,
			wantConfig: params.TestnetChainConfig,
		},
		{
			name: "compatible config in DB",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
				oldcustomg.MustCommit(db)
				return SetupGenesisBlock(db, &customg)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
		},
		{
			name: "incompatible config in DB",
			fn: func(db ethdb.Database) (*params.ChainConfig, common.Hash, error) {
//在2进行“古老”的创世地块和宅基地过渡。
//前进到4号街区，经过海关的宅基地过渡街区。
				genesis := oldcustomg.MustCommit(db)

				bc, _ := NewBlockChain(db, nil, oldcustomg.Config, ethash.NewFullFaker(), vm.Config{}, nil)
				defer bc.Stop()

				blocks, _ := GenerateChain(oldcustomg.Config, genesis, ethash.NewFaker(), db, 4, nil)
				bc.InsertChain(blocks)
				bc.CurrentBlock()
//这将返回一个兼容性错误。
				return SetupGenesisBlock(db, &customg)
			},
			wantHash:   customghash,
			wantConfig: customg.Config,
			wantErr: &params.ConfigCompatError{
				What:         "Homestead fork block",
				StoredConfig: big.NewInt(2),
				NewConfig:    big.NewInt(3),
				RewindTo:     1,
			},
		},
	}

	for _, test := range tests {
		db := ethdb.NewMemDatabase()
		config, hash, err := test.fn(db)
//检查返回值。
		if !reflect.DeepEqual(err, test.wantErr) {
			spew := spew.ConfigState{DisablePointerAddresses: true, DisableCapacities: true}
			t.Errorf("%s: returned error %#v, want %#v", test.name, spew.NewFormatter(err), spew.NewFormatter(test.wantErr))
		}
		if !reflect.DeepEqual(config, test.wantConfig) {
			t.Errorf("%s:\nreturned %v\nwant     %v", test.name, config, test.wantConfig)
		}
		if hash != test.wantHash {
			t.Errorf("%s: returned hash %s, want %s", test.name, hash.Hex(), test.wantHash.Hex())
		} else if err == nil {
//检查数据库内容。
			stored := rawdb.ReadBlock(db, test.wantHash, 0)
			if stored.Hash() != test.wantHash {
				t.Errorf("%s: block in DB has hash %s, want %s", test.name, stored.Hash(), test.wantHash)
			}
		}
	}
}

