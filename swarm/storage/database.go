
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450118451466240>


package storage

//这是以太坊ethdb/数据库早期状态的克隆
//不需要排队/缓存

import (
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const openFileLimit = 128

type LDBDatabase struct {
	db *leveldb.DB
}

func NewLDBDatabase(file string) (*LDBDatabase, error) {
//打开数据库
	db, err := leveldb.OpenFile(file, &opt.Options{OpenFilesCacheCapacity: openFileLimit})
	if err != nil {
		return nil, err
	}

	database := &LDBDatabase{db: db}

	return database, nil
}

func (db *LDBDatabase) Put(key []byte, value []byte) error {
	metrics.GetOrRegisterCounter("ldbdatabase.put", nil).Inc(1)

	return db.db.Put(key, value, nil)
}

func (db *LDBDatabase) Get(key []byte) ([]byte, error) {
	metrics.GetOrRegisterCounter("ldbdatabase.get", nil).Inc(1)

	dat, err := db.db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (db *LDBDatabase) Delete(key []byte) error {
	return db.db.Delete(key, nil)
}

func (db *LDBDatabase) NewIterator() iterator.Iterator {
	metrics.GetOrRegisterCounter("ldbdatabase.newiterator", nil).Inc(1)

	return db.db.NewIterator(nil, nil)
}

func (db *LDBDatabase) Write(batch *leveldb.Batch) error {
	metrics.GetOrRegisterCounter("ldbdatabase.write", nil).Inc(1)

	return db.db.Write(batch, nil)
}

func (db *LDBDatabase) Close() {
//关闭LevelDB数据库
	db.db.Close()
}

