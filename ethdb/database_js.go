
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450090370600960>


//+构建JS

package ethdb

import (
	"errors"
)

var errNotSupported = errors.New("ethdb: not supported")

type LDBDatabase struct {
}

//NewLdbDatabase返回一个LevelDB包装的对象。
func NewLDBDatabase(file string, cache int, handles int) (*LDBDatabase, error) {
	return nil, errNotSupported
}

//path返回数据库目录的路径。
func (db *LDBDatabase) Path() string {
	return ""
}

//Put将给定的键/值放入队列
func (db *LDBDatabase) Put(key []byte, value []byte) error {
	return errNotSupported
}

func (db *LDBDatabase) Has(key []byte) (bool, error) {
	return false, errNotSupported
}

//get返回给定的键（如果存在）。
func (db *LDBDatabase) Get(key []byte) ([]byte, error) {
	return nil, errNotSupported
}

//删除从队列和数据库中删除键
func (db *LDBDatabase) Delete(key []byte) error {
	return errNotSupported
}

func (db *LDBDatabase) Close() {
}

//Meter配置数据库度量收集器和
func (db *LDBDatabase) Meter(prefix string) {
}

func (db *LDBDatabase) NewBatch() Batch {
	return nil
}

