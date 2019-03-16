
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450118044618752>


package shed

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
//用于存储架构的LevelDB键值。
	keySchema = []byte{0}
//所有字段类型的leveldb键前缀。
//将通过将名称值附加到此前缀来构造LevelDB键。
	keyPrefixFields byte = 1
//索引键起始的级别数据库键前缀。
//每个索引都有自己的键前缀，这个值定义了第一个。
keyPrefixIndexStart byte = 2 //问：或者可能是更高的数字，比如7，为潜在的特定性能提供更多的空间。
)

//架构用于序列化已知的数据库结构信息。
type schema struct {
Fields  map[string]fieldSpec `json:"fields"`  //键是字段名
Indexes map[byte]indexSpec   `json:"indexes"` //键是索引前缀字节
}

//fieldspec保存有关特定字段的信息。
//它不需要名称字段，因为它包含在
//架构。字段映射键。
type fieldSpec struct {
	Type string `json:"type"`
}

//indxspec保存有关特定索引的信息。
//它不包含索引类型，因为索引没有类型。
type indexSpec struct {
	Name string `json:"name"`
}

//SchemaFieldKey检索的完整级别数据库键
//一个特定的字段构成了模式定义。
func (db *DB) schemaFieldKey(name, fieldType string) (key []byte, err error) {
	if name == "" {
		return nil, errors.New("field name can not be blank")
	}
	if fieldType == "" {
		return nil, errors.New("field type can not be blank")
	}
	s, err := db.getSchema()
	if err != nil {
		return nil, err
	}
	var found bool
	for n, f := range s.Fields {
		if n == name {
			if f.Type != fieldType {
				return nil, fmt.Errorf("field %q of type %q stored as %q in db", name, fieldType, f.Type)
			}
			break
		}
	}
	if !found {
		s.Fields[name] = fieldSpec{
			Type: fieldType,
		}
		err := db.putSchema(s)
		if err != nil {
			return nil, err
		}
	}
	return append([]byte{keyPrefixFields}, []byte(name)...), nil
}

//SchemaIndexID检索的完整级别数据库前缀
//一种特殊的索引。
func (db *DB) schemaIndexPrefix(name string) (id byte, err error) {
	if name == "" {
		return 0, errors.New("index name can not be blank")
	}
	s, err := db.getSchema()
	if err != nil {
		return 0, err
	}
	nextID := keyPrefixIndexStart
	for i, f := range s.Indexes {
		if i >= nextID {
			nextID = i + 1
		}
		if f.Name == name {
			return i, nil
		}
	}
	id = nextID
	s.Indexes[id] = indexSpec{
		Name: name,
	}
	return id, db.putSchema(s)
}

//GetSchema从中检索完整的架构
//数据库。
func (db *DB) getSchema() (s schema, err error) {
	b, err := db.Get(keySchema)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(b, &s)
	return s, err
}

//PutSchema将完整的架构存储到
//数据库。
func (db *DB) putSchema(s schema) (err error) {
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return db.Put(keySchema, b)
}

