
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450117713268736>


package shed

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
)

//structField是用于存储复杂结构的帮助程序
//以rlp格式对其进行编码。
type StructField struct {
	db  *DB
	key []byte
}

//NewstructField返回新的structField。
//它根据数据库模式验证其名称和类型。
func (db *DB) NewStructField(name string) (f StructField, err error) {
	key, err := db.schemaFieldKey(name, "struct-rlp")
	if err != nil {
		return f, err
	}
	return StructField{
		db:  db,
		key: key,
	}, nil
}

//将数据库中的数据解包到提供的VAL。
//如果找不到数据，则返回leveldb.errnotfound。
func (f StructField) Get(val interface{}) (err error) {
	b, err := f.db.Get(f.key)
	if err != nil {
		return err
	}
	return rlp.DecodeBytes(b, val)
}

//放入marshals提供的val并将其保存到数据库中。
func (f StructField) Put(val interface{}) (err error) {
	b, err := rlp.EncodeToBytes(val)
	if err != nil {
		return err
	}
	return f.db.Put(f.key, b)
}

//Putinbatch Marshals提供了VAL并将其放入批处理中。
func (f StructField) PutInBatch(batch *leveldb.Batch, val interface{}) (err error) {
	b, err := rlp.EncodeToBytes(val)
	if err != nil {
		return err
	}
	batch.Put(f.key, b)
	return nil
}

