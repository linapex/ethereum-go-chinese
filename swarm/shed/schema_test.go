
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450118094950400>


package shed

import (
	"bytes"
	"testing"
)

//testdb_schemafieldkey验证schemafieldkey的正确性。
func TestDB_schemaFieldKey(t *testing.T) {
	db, cleanupFunc := newTestDB(t)
	defer cleanupFunc()

	t.Run("empty name or type", func(t *testing.T) {
		_, err := db.schemaFieldKey("", "")
		if err == nil {
			t.Errorf("error not returned, but expected")
		}
		_, err = db.schemaFieldKey("", "type")
		if err == nil {
			t.Errorf("error not returned, but expected")
		}

		_, err = db.schemaFieldKey("test", "")
		if err == nil {
			t.Errorf("error not returned, but expected")
		}
	})

	t.Run("same field", func(t *testing.T) {
		key1, err := db.schemaFieldKey("test", "undefined")
		if err != nil {
			t.Fatal(err)
		}

		key2, err := db.schemaFieldKey("test", "undefined")
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(key1, key2) {
			t.Errorf("schema keys for the same field name are not the same: %q, %q", string(key1), string(key2))
		}
	})

	t.Run("different fields", func(t *testing.T) {
		key1, err := db.schemaFieldKey("test1", "undefined")
		if err != nil {
			t.Fatal(err)
		}

		key2, err := db.schemaFieldKey("test2", "undefined")
		if err != nil {
			t.Fatal(err)
		}

		if bytes.Equal(key1, key2) {
			t.Error("schema keys for the same field name are the same, but must not be")
		}
	})

	t.Run("same field name different types", func(t *testing.T) {
		_, err := db.schemaFieldKey("the-field", "one-type")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.schemaFieldKey("the-field", "another-type")
		if err == nil {
			t.Errorf("error not returned, but expected")
		}
	})
}

//testdb_schemaIndexPrefix验证schemaIndexPrefix的正确性。
func TestDB_schemaIndexPrefix(t *testing.T) {
	db, cleanupFunc := newTestDB(t)
	defer cleanupFunc()

	t.Run("same name", func(t *testing.T) {
		id1, err := db.schemaIndexPrefix("test")
		if err != nil {
			t.Fatal(err)
		}

		id2, err := db.schemaIndexPrefix("test")
		if err != nil {
			t.Fatal(err)
		}

		if id1 != id2 {
			t.Errorf("schema keys for the same field name are not the same: %v, %v", id1, id2)
		}
	})

	t.Run("different names", func(t *testing.T) {
		id1, err := db.schemaIndexPrefix("test1")
		if err != nil {
			t.Fatal(err)
		}

		id2, err := db.schemaIndexPrefix("test2")
		if err != nil {
			t.Fatal(err)
		}

		if id1 == id2 {
			t.Error("schema ids for the same index name are the same, but must not be")
		}
	})
}

