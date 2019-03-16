
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450074717458432>


package consensus

import "errors"

var (
//当验证块需要祖先时返回errUnknownancestor
//这是未知的。
	ErrUnknownAncestor = errors.New("unknown ancestor")

//验证块需要祖先时返回errprunedancestor
//这是已知的，但其状态不可用。
	ErrPrunedAncestor = errors.New("pruned ancestor")

//当块的时间戳在将来时，根据
//到当前节点。
	ErrFutureBlock = errors.New("block in the future")

//如果块的编号不等于其父块的编号，则返回errInvalidNumber。
//加一。
	ErrInvalidNumber = errors.New("invalid block number")
)

