
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450118199808000>


package state

import (
	"encoding"
	"encoding/json"
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
)

//如果没有从数据库返回结果，则返回errnotfound。
var ErrNotFound = errors.New("ErrorNotFound")

//存储区定义获取、设置和删除不同键的值所需的方法
//关闭基础资源。
type Store interface {
	Get(key string, i interface{}) (err error)
	Put(key string, i interface{}) (err error)
	Delete(key string) (err error)
	Close() error
}

//dbstore使用leveldb存储值。
type DBStore struct {
	db *leveldb.DB
}

//new dbstore创建dbstore的新实例。
func NewDBStore(path string) (s *DBStore, err error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}
	return &DBStore{
		db: db,
	}, nil
}

//newInMemoryStore返回dbstore的新实例。仅用于测试和模拟。
func NewInmemoryStore() *DBStore {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		panic(err)
	}
	return &DBStore{
		db: db,
	}
}

//get检索特定键的持久化值。如果没有结果
//返回errnotfound。提供的参数应为字节片或
//实现encoding.binaryUnmarshaler接口的结构
func (s *DBStore) Get(key string, i interface{}) (err error) {
	has, err := s.db.Has([]byte(key), nil)
	if err != nil || !has {
		return ErrNotFound
	}

	data, err := s.db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return ErrNotFound
	}

	unmarshaler, ok := i.(encoding.BinaryUnmarshaler)
	if !ok {
		return json.Unmarshal(data, i)
	}
	return unmarshaler.UnmarshalBinary(data)
}

//Put存储为特定键实现二进制的对象。
func (s *DBStore) Put(key string, i interface{}) (err error) {
	var bytes []byte

	marshaler, ok := i.(encoding.BinaryMarshaler)
	if !ok {
		if bytes, err = json.Marshal(i); err != nil {
			return err
		}
	} else {
		if bytes, err = marshaler.MarshalBinary(); err != nil {
			return err
		}
	}

	return s.db.Put([]byte(key), bytes, nil)
}

//删除删除存储在特定键下的条目。
func (s *DBStore) Delete(key string) (err error) {
	return s.db.Delete([]byte(key), nil)
}

//close释放基础级别db使用的资源。
func (s *DBStore) Close() error {
	return s.db.Close()
}

