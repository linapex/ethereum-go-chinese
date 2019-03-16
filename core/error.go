
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450078668492800>


package core

import "errors"

var (
//当要导入的块在本地已知时，返回errknownBlock。
	ErrKnownBlock = errors.New("block already known")

//如果所需的气体量达到，则气体池将返回达到的ErrgasLimited。
//一个事务比块中剩余的事务高。
	ErrGasLimitReached = errors.New("gas limit reached")

//如果要导入的块在黑名单上，则返回errBlackListedHash。
	ErrBlacklistedHash = errors.New("blacklisted hash")

//如果事务的nonce高于
//下一个基于本地链的期望值。
	ErrNonceTooHigh = errors.New("nonce too high")
)

