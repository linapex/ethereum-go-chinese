
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450120787693568>


//包模拟定义了不同实现使用的类型
//模拟仓库。
//
//模拟存储的实现位于目录中
//在此包下：
//
//-db-leveldb后端
//-mem-内存映射后端
//-rpc-可以连接到其他后端的rpc客户端
//
//模拟存储可以实现导入和导出接口
//用于导入和导出它们包含的所有块数据。
//导出的文件是一个tar存档，其中所有文件的名称为
//块键和内容的十六进制表示法
//具有json编码的exportedchunk结构。导出格式
//应该在所有模拟存储实现中保留。
package mock

import (
	"errors"
	"io"

	"github.com/ethereum/go-ethereum/common"
)

//errnotfound表示找不到块。
var ErrNotFound = errors.New("not found")

//nodestore保存节点地址和对globalStore的引用
//为了只访问和存储一个节点的块数据。
type NodeStore struct {
	store GlobalStorer
	addr  common.Address
}

//new nodestore创建nodestore的新实例
//使用提供地址的GlobalStrer对数据进行分组。
func NewNodeStore(addr common.Address, store GlobalStorer) *NodeStore {
	return &NodeStore{
		store: store,
		addr:  addr,
	}
}

//get返回具有地址的节点的键的块数据
//在节点存储初始化时提供。
func (n *NodeStore) Get(key []byte) (data []byte, err error) {
	return n.store.Get(n.addr, key)
}

//Put为具有地址的节点的键保存块数据
//在节点存储初始化时提供。
func (n *NodeStore) Put(key []byte, data []byte) error {
	return n.store.Put(n.addr, key, data)
}

//删除删除具有地址的节点的键的块数据
//在节点存储初始化时提供。
func (n *NodeStore) Delete(key []byte) error {
	return n.store.Delete(n.addr, key)
}

//GlobalStrer定义模拟数据库存储的方法
//存储所有群节点的块数据。
//在测试中用来构造模拟节点库
//并跟踪和验证块。
type GlobalStorer interface {
	Get(addr common.Address, key []byte) (data []byte, err error)
	Put(addr common.Address, key []byte, data []byte) error
	Delete(addr common.Address, key []byte) error
	HasKey(addr common.Address, key []byte) bool
//newnodestore创建nodestore的实例
//用于单个群节点
//地址地址
	NewNodeStore(addr common.Address) *NodeStore
}

//导入程序定义导入模拟存储数据的方法
//从导出的tar存档。
type Importer interface {
	Import(r io.Reader) (n int, err error)
}

//导出器定义用于导出模拟存储数据的方法
//去焦油档案馆。
type Exporter interface {
	Export(w io.Writer) (n int, err error)
}

//exportedchunk是保存在tar存档中的结构，用于
//每个块都是JSON编码的字节。
type ExportedChunk struct {
	Data  []byte           `json:"d"`
	Addrs []common.Address `json:"a"`
}

