
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450075292078080>


package misc

import (
	"bytes"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

var (
//如果头不支持
//Pro Fork客户端。
	ErrBadProDAOExtra = errors.New("bad DAO pro-fork extra-data")

//如果头确实支持no-
//分支客户机。
	ErrBadNoDAOExtra = errors.New("bad DAO no-fork extra-data")
)

//verifydaoHeaderextradata验证块头的额外数据字段
//确保符合刀硬叉规则。
//
//DAO硬分叉扩展到头的有效性：
//
//使用fork特定的额外数据集
//b）如果节点是pro fork，则需要特定范围内的块具有
//唯一的额外数据集。
func VerifyDAOHeaderExtraData(config *params.ChainConfig, header *types.Header) error {
//如果节点不关心DAO分叉，则进行短路验证
	if config.DAOForkBlock == nil {
		return nil
	}
//确保块在fork修改的额外数据范围内
	limit := new(big.Int).Add(config.DAOForkBlock, params.DAOForkExtraRange)
	if header.Number.Cmp(config.DAOForkBlock) < 0 || header.Number.Cmp(limit) >= 0 {
		return nil
	}
//根据我们支持还是反对fork，验证额外的数据内容
	if config.DAOForkSupport {
		if !bytes.Equal(header.Extra, params.DAOForkBlockExtra) {
			return ErrBadProDAOExtra
		}
	} else {
		if bytes.Equal(header.Extra, params.DAOForkBlockExtra) {
			return ErrBadNoDAOExtra
		}
	}
//好吧，header有我们期望的额外数据
	return nil
}

//ApplyDaoHardFork根据DAO Hard Fork修改状态数据库
//规则，将一组DAO帐户的所有余额转移到单个退款
//合同。
func ApplyDAOHardFork(statedb *state.StateDB) {
//检索要将余额退款到的合同
	if !statedb.Exist(params.DAORefundContract) {
		statedb.CreateAccount(params.DAORefundContract)
	}

//将每个DAO帐户和额外的余额帐户资金转移到退款合同中
	for _, addr := range params.DAODrainList() {
		statedb.AddBalance(params.DAORefundContract, statedb.GetBalance(addr))
		statedb.SetBalance(addr, new(big.Int))
	}
}

