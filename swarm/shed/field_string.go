
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450117608411136>


package shed

import (
	"github.com/syndtr/goleveldb/leveldb"
)

//StringField是最简单的字段实现
//它在特定的leveldb键下存储一个任意字符串。
type StringField struct {
	db  *DB
	key []byte
}

//newStringField重新运行StringField的新实例。
//它根据数据库模式验证其名称和类型。
func (db *DB) NewStringField(name string) (f StringField, err error) {
	key, err := db.schemaFieldKey(name, "string")
	if err != nil {
		return f, err
	}
	return StringField{
		db:  db,
		key: key,
	}, nil
}

//get返回数据库中的字符串值。
//如果找不到该值，则返回空字符串
//没有错误。
func (f StringField) Get() (val string, err error) {
	b, err := f.db.Get(f.key)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return "", nil
		}
		return "", err
	}
	return string(b), nil
}

//将存储字符串放入数据库。
func (f StringField) Put(val string) (err error) {
	return f.db.Put(f.key, []byte(val))
}

//putinbatch将字符串存储在可以
//稍后保存在数据库中。
func (f StringField) PutInBatch(batch *leveldb.Batch, val string) {
	batch.Put(f.key, []byte(val))
}

