
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450074918785024>


package ethash

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

var errEthashStopped = errors.New("ethash stopped")

//API为RPC接口公开与ethash相关的方法。
type API struct {
ethash *Ethash //确保ethash模式正常。
}

//GetWork返回外部矿工的工作包。
//
//工作包由3个字符串组成：
//结果[0]-32字节十六进制编码的当前块头POW哈希
//结果[1]-用于DAG的32字节十六进制编码种子哈希
//结果[2]-32字节十六进制编码边界条件（“目标”），2^256/难度
//结果[3]-十六进制编码的块号
func (api *API) GetWork() ([4]string, error) {
	if api.ethash.config.PowMode != ModeNormal && api.ethash.config.PowMode != ModeTest {
		return [4]string{}, errors.New("not supported")
	}

	var (
		workCh = make(chan [4]string, 1)
		errc   = make(chan error, 1)
	)

	select {
	case api.ethash.fetchWorkCh <- &sealWork{errc: errc, res: workCh}:
	case <-api.ethash.exitCh:
		return [4]string{}, errEthashStopped
	}

	select {
	case work := <-workCh:
		return work, nil
	case err := <-errc:
		return [4]string{}, err
	}
}

//外部矿工可使用SubNetwork提交其POW解决方案。
//如果工作被接受，它会返回一个指示。
//注意：如果解决方案无效，则过时的工作不存在的工作将返回false。
func (api *API) SubmitWork(nonce types.BlockNonce, hash, digest common.Hash) bool {
	if api.ethash.config.PowMode != ModeNormal && api.ethash.config.PowMode != ModeTest {
		return false
	}

	var errc = make(chan error, 1)

	select {
	case api.ethash.submitWorkCh <- &mineResult{
		nonce:     nonce,
		mixDigest: digest,
		hash:      hash,
		errc:      errc,
	}:
	case <-api.ethash.exitCh:
		return false
	}

	err := <-errc
	return err == nil
}

//submit hash rate可用于远程矿工提交哈希率。
//这使节点能够报告所有矿工的组合哈希率
//通过这个节点提交工作。
//
//它接受矿工哈希率和标识符，该标识符必须是唯一的
//节点之间。
func (api *API) SubmitHashRate(rate hexutil.Uint64, id common.Hash) bool {
	if api.ethash.config.PowMode != ModeNormal && api.ethash.config.PowMode != ModeTest {
		return false
	}

	var done = make(chan struct{}, 1)

	select {
	case api.ethash.submitRateCh <- &hashrate{done: done, rate: uint64(rate), id: id}:
	case <-api.ethash.exitCh:
		return false
	}

//阻止，直到哈希率成功提交。
	<-done

	return true
}

//GetHashrate返回本地CPU矿工和远程矿工的当前哈希率。
func (api *API) GetHashrate() uint64 {
	return uint64(api.ethash.Hashrate())
}

