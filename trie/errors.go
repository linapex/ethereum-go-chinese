
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450122847096832>


package trie

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

//trie函数（tryget、tryupdate、trydelete）返回MissingNodeError
//如果本地数据库中没有trie节点。它包含
//检索丢失节点所需的信息。
type MissingNodeError struct {
NodeHash common.Hash //缺少节点的哈希
Path     []byte      //丢失节点的十六进制编码路径
}

func (err *MissingNodeError) Error() string {
	return fmt.Sprintf("missing trie node %x (path %x)", err.NodeHash, err.Path)
}

