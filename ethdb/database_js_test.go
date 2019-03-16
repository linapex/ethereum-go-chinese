
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450090412544000>


//+构建JS

package ethdb_test

import (
	"github.com/ethereum/go-ethereum/ethdb"
)

var _ ethdb.Database = &ethdb.LDBDatabase{}

