
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450075338215424>


package misc

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

//verifyforkhashes验证符合网络硬分叉的块是否具有
//正确的散列，以避免客户在不同的链上离开。这是一个
//可选功能。
func VerifyForkHashes(config *params.ChainConfig, header *types.Header, uncle bool) error {
//我们不关心叔叔
	if uncle {
		return nil
	}
//如果设置了homestead重定价哈希，请验证它
	if config.EIP150Block != nil && config.EIP150Block.Cmp(header.Number) == 0 {
		if config.EIP150Hash != (common.Hash{}) && config.EIP150Hash != header.Hash() {
			return fmt.Errorf("homestead gas reprice fork: have 0x%x, want 0x%x", header.Hash(), config.EIP150Hash)
		}
	}
//一切还好，归来
	return nil
}

