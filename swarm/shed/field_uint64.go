
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450117809737728>


package shed

import (
	"encoding/binary"

	"github.com/syndtr/goleveldb/leveldb"
)

//uint64字段提供了一种在数据库中使用简单计数器的方法。
//它透明地将uint64类型值编码为字节。
type Uint64Field struct {
	db  *DB
	key []byte
}

//newuint64字段返回新的uint64字段。
//它根据数据库模式验证其名称和类型。
func (db *DB) NewUint64Field(name string) (f Uint64Field, err error) {
	key, err := db.schemaFieldKey(name, "uint64")
	if err != nil {
		return f, err
	}
	return Uint64Field{
		db:  db,
		key: key,
	}, nil
}

//get从数据库中检索uint64值。
//如果在数据库中找不到该值，则为0
//返回，没有错误。
func (f Uint64Field) Get() (val uint64, err error) {
	b, err := f.db.Get(f.key)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return 0, nil
		}
		return 0, err
	}
	return binary.BigEndian.Uint64(b), nil
}

//放入编码uin64值并将其存储在数据库中。
func (f Uint64Field) Put(val uint64) (err error) {
	return f.db.Put(f.key, encodeUint64(val))
}

//Putinbatch在批中存储uint64值
//以后可以保存到数据库中。
func (f Uint64Field) PutInBatch(batch *leveldb.Batch, val uint64) {
	batch.Put(f.key, encodeUint64(val))
}

//inc在数据库中增加uint64值。
//此操作不是goroutine save。
func (f Uint64Field) Inc() (val uint64, err error) {
	val, err = f.Get()
	if err != nil {
		if err == leveldb.ErrNotFound {
			val = 0
		} else {
			return 0, err
		}
	}
	val++
	return val, f.Put(val)
}

//incinbatch在批中增加uint64值
//通过从数据库中检索值，而不是同一批。
//此操作不是goroutine save。
func (f Uint64Field) IncInBatch(batch *leveldb.Batch) (val uint64, err error) {
	val, err = f.Get()
	if err != nil {
		if err == leveldb.ErrNotFound {
			val = 0
		} else {
			return 0, err
		}
	}
	val++
	f.PutInBatch(batch, val)
	return val, nil
}

//DEC减少数据库中的uint64值。
//此操作不是goroutine save。
//该字段受到保护，不会溢出为负值。
func (f Uint64Field) Dec() (val uint64, err error) {
	val, err = f.Get()
	if err != nil {
		if err == leveldb.ErrNotFound {
			val = 0
		} else {
			return 0, err
		}
	}
	if val != 0 {
		val--
	}
	return val, f.Put(val)
}

//decinbatch在批处理中递减uint64值
//通过从数据库中检索值，而不是同一批。
//此操作不是goroutine save。
//该字段受到保护，不会溢出为负值。
func (f Uint64Field) DecInBatch(batch *leveldb.Batch) (val uint64, err error) {
	val, err = f.Get()
	if err != nil {
		if err == leveldb.ErrNotFound {
			val = 0
		} else {
			return 0, err
		}
	}
	if val != 0 {
		val--
	}
	f.PutInBatch(batch, val)
	return val, nil
}

//编码将uint64转换为8字节长
//以big-endian编码切片。
func encodeUint64(val uint64) (b []byte) {
	b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, val)
	return b
}

