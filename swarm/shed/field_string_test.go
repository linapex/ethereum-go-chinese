
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450117658742784>


package shed

import (
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

//TestStringField验证Put和Get操作
//在StringField中。
func TestStringField(t *testing.T) {
	db, cleanupFunc := newTestDB(t)
	defer cleanupFunc()

	simpleString, err := db.NewStringField("simple-string")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get empty", func(t *testing.T) {
		got, err := simpleString.Get()
		if err != nil {
			t.Fatal(err)
		}
		want := ""
		if got != want {
			t.Errorf("got string %q, want %q", got, want)
		}
	})

	t.Run("put", func(t *testing.T) {
		want := "simple string value"
		err = simpleString.Put(want)
		if err != nil {
			t.Fatal(err)
		}
		got, err := simpleString.Get()
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("got string %q, want %q", got, want)
		}

		t.Run("overwrite", func(t *testing.T) {
			want := "overwritten string value"
			err = simpleString.Put(want)
			if err != nil {
				t.Fatal(err)
			}
			got, err := simpleString.Get()
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Errorf("got string %q, want %q", got, want)
			}
		})
	})

	t.Run("put in batch", func(t *testing.T) {
		batch := new(leveldb.Batch)
		want := "simple string batch value"
		simpleString.PutInBatch(batch, want)
		err = db.WriteBatch(batch)
		if err != nil {
			t.Fatal(err)
		}
		got, err := simpleString.Get()
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("got string %q, want %q", got, want)
		}

		t.Run("overwrite", func(t *testing.T) {
			batch := new(leveldb.Batch)
			want := "overwritten string batch value"
			simpleString.PutInBatch(batch, want)
			err = db.WriteBatch(batch)
			if err != nil {
				t.Fatal(err)
			}
			got, err := simpleString.Get()
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Errorf("got string %q, want %q", got, want)
			}
		})
	})
}

