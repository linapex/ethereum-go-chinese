
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450120653475840>


package mem

import (
	"testing"

	"github.com/ethereum/go-ethereum/swarm/storage/mock/test"
)

//TestGlobalStore正在为GlobalStore运行测试
//使用test.mockstore函数。
func TestGlobalStore(t *testing.T) {
	test.MockStore(t, NewGlobalStore(), 100)
}

//testmortexport正在运行用于导入和
//在两个GlobalStores之间导出数据
//使用test.importexport函数。
func TestImportExport(t *testing.T) {
	test.ImportExport(t, NewGlobalStore(), NewGlobalStore(), 100)
}

