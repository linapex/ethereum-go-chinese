
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450090492235776>


package ethdb

//使用批处理的代码应该尝试向批处理中添加这么多的数据。
//该值是根据经验确定的。
const IdealBatchSize = 100 * 1024

//推杆包装批处理和常规数据库都支持的数据库写入操作。
type Putter interface {
	Put(key []byte, value []byte) error
}

//删除程序包装批处理数据库和常规数据库都支持的数据库删除操作。
type Deleter interface {
	Delete(key []byte) error
}

//数据库包装所有数据库操作。所有方法对于并发使用都是安全的。
type Database interface {
	Putter
	Deleter
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Close()
	NewBatch() Batch
}

//批处理是一个只写的数据库，它将更改提交到其主机数据库。
//当调用写入时。批处理不能同时使用。
type Batch interface {
	Putter
	Deleter
ValueSize() int //批中的数据量
	Write() error
//重置将批重置为可重用
	Reset()
}

