
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450095663812608>


package les

import (
	"bytes"
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

type odrTestFn func(ctx context.Context, db ethdb.Database, config *params.ChainConfig, bc *core.BlockChain, lc *light.LightChain, bhash common.Hash) []byte

func TestOdrGetBlockLes1(t *testing.T) { testOdr(t, 1, 1, odrGetBlock) }

func TestOdrGetBlockLes2(t *testing.T) { testOdr(t, 2, 1, odrGetBlock) }

func odrGetBlock(ctx context.Context, db ethdb.Database, config *params.ChainConfig, bc *core.BlockChain, lc *light.LightChain, bhash common.Hash) []byte {
	var block *types.Block
	if bc != nil {
		block = bc.GetBlockByHash(bhash)
	} else {
		block, _ = lc.GetBlockByHash(ctx, bhash)
	}
	if block == nil {
		return nil
	}
	rlp, _ := rlp.EncodeToBytes(block)
	return rlp
}

func TestOdrGetReceiptsLes1(t *testing.T) { testOdr(t, 1, 1, odrGetReceipts) }

func TestOdrGetReceiptsLes2(t *testing.T) { testOdr(t, 2, 1, odrGetReceipts) }

func odrGetReceipts(ctx context.Context, db ethdb.Database, config *params.ChainConfig, bc *core.BlockChain, lc *light.LightChain, bhash common.Hash) []byte {
	var receipts types.Receipts
	if bc != nil {
		if number := rawdb.ReadHeaderNumber(db, bhash); number != nil {
			receipts = rawdb.ReadReceipts(db, bhash, *number)
		}
	} else {
		if number := rawdb.ReadHeaderNumber(db, bhash); number != nil {
			receipts, _ = light.GetBlockReceipts(ctx, lc.Odr(), bhash, *number)
		}
	}
	if receipts == nil {
		return nil
	}
	rlp, _ := rlp.EncodeToBytes(receipts)
	return rlp
}

func TestOdrAccountsLes1(t *testing.T) { testOdr(t, 1, 1, odrAccounts) }

func TestOdrAccountsLes2(t *testing.T) { testOdr(t, 2, 1, odrAccounts) }

func odrAccounts(ctx context.Context, db ethdb.Database, config *params.ChainConfig, bc *core.BlockChain, lc *light.LightChain, bhash common.Hash) []byte {
	dummyAddr := common.HexToAddress("1234567812345678123456781234567812345678")
	acc := []common.Address{testBankAddress, acc1Addr, acc2Addr, dummyAddr}

	var (
		res []byte
		st  *state.StateDB
		err error
	)
	for _, addr := range acc {
		if bc != nil {
			header := bc.GetHeaderByHash(bhash)
			st, err = state.New(header.Root, state.NewDatabase(db))
		} else {
			header := lc.GetHeaderByHash(bhash)
			st = light.NewState(ctx, header, lc.Odr())
		}
		if err == nil {
			bal := st.GetBalance(addr)
			rlp, _ := rlp.EncodeToBytes(bal)
			res = append(res, rlp...)
		}
	}
	return res
}

func TestOdrContractCallLes1(t *testing.T) { testOdr(t, 1, 2, odrContractCall) }

func TestOdrContractCallLes2(t *testing.T) { testOdr(t, 2, 2, odrContractCall) }

type callmsg struct {
	types.Message
}

func (callmsg) CheckNonce() bool { return false }

func odrContractCall(ctx context.Context, db ethdb.Database, config *params.ChainConfig, bc *core.BlockChain, lc *light.LightChain, bhash common.Hash) []byte {
	data := common.Hex2Bytes("60CD26850000000000000000000000000000000000000000000000000000000000000000")

	var res []byte
	for i := 0; i < 3; i++ {
		data[35] = byte(i)
		if bc != nil {
			header := bc.GetHeaderByHash(bhash)
			statedb, err := state.New(header.Root, state.NewDatabase(db))

			if err == nil {
				from := statedb.GetOrNewStateObject(testBankAddress)
				from.SetBalance(math.MaxBig256)

				msg := callmsg{types.NewMessage(from.Address(), &testContractAddr, 0, new(big.Int), 100000, new(big.Int), data, false)}

				context := core.NewEVMContext(msg, header, bc, nil)
				vmenv := vm.NewEVM(context, statedb, config, vm.Config{})

//vmenv：=core.newenv（statedb，config，bc，msg，header，vm.config）
				gp := new(core.GasPool).AddGas(math.MaxUint64)
				ret, _, _, _ := core.ApplyMessage(vmenv, msg, gp)
				res = append(res, ret...)
			}
		} else {
			header := lc.GetHeaderByHash(bhash)
			state := light.NewState(ctx, header, lc.Odr())
			state.SetBalance(testBankAddress, math.MaxBig256)
			msg := callmsg{types.NewMessage(testBankAddress, &testContractAddr, 0, new(big.Int), 100000, new(big.Int), data, false)}
			context := core.NewEVMContext(msg, header, lc, nil)
			vmenv := vm.NewEVM(context, state, config, vm.Config{})
			gp := new(core.GasPool).AddGas(math.MaxUint64)
			ret, _, _, _ := core.ApplyMessage(vmenv, msg, gp)
			if state.Error() == nil {
				res = append(res, ret...)
			}
		}
	}
	return res
}

//testOdr tests odr requests whose validation guaranteed by block headers.
func testOdr(t *testing.T, protocol int, expFail uint64, fn odrTestFn) {
//组装测试环境
	server, client, tearDown := newClientServerEnv(t, 4, protocol, nil, true)
	defer tearDown()
	client.pm.synchronise(client.rPeer)

	test := func(expFail uint64) {
		for i := uint64(0); i <= server.pm.blockchain.CurrentHeader().Number.Uint64(); i++ {
			bhash := rawdb.ReadCanonicalHash(server.db, i)
			b1 := fn(light.NoOdr, server.db, server.pm.chainConfig, server.pm.blockchain.(*core.BlockChain), nil, bhash)

			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()
			b2 := fn(ctx, client.db, client.pm.chainConfig, nil, client.pm.blockchain.(*light.LightChain), bhash)

			eq := bytes.Equal(b1, b2)
			exp := i < expFail
			if exp && !eq {
				t.Errorf("odr mismatch")
			}
			if !exp && eq {
				t.Errorf("unexpected odr match")
			}
		}
	}
//暂时删除对等测试ODR失败
//预计在没有LES对等机的情况下检索失败（Genesis块除外）
	client.peers.Unregister(client.rPeer.id)
time.Sleep(time.Millisecond * 10) //ensure that all peerSetNotify callbacks are executed
	test(expFail)
//希望所有检索都通过
	client.peers.Register(client.rPeer)
time.Sleep(time.Millisecond * 10) //确保执行所有peersetnotify回调
	client.peers.lock.Lock()
	client.rPeer.hasBlock = func(common.Hash, uint64, bool) bool { return true }
	client.peers.lock.Unlock()
	test(5)
//仍然希望通过所有检索，现在应该在本地缓存数据
	client.peers.Unregister(client.rPeer.id)
time.Sleep(time.Millisecond * 10) //确保执行所有peersetnotify回调
	test(5)
}

