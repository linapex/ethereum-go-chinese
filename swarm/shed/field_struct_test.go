
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450117755211776>


package shed

import (
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

//teststructfield验证put和get操作
//结构字段的。
func TestStructField(t *testing.T) {
	db, cleanupFunc := newTestDB(t)
	defer cleanupFunc()

	complexField, err := db.NewStructField("complex-field")
	if err != nil {
		t.Fatal(err)
	}

	type complexStructure struct {
		A string
	}

	t.Run("get empty", func(t *testing.T) {
		var s complexStructure
		err := complexField.Get(&s)
		if err != leveldb.ErrNotFound {
			t.Fatalf("got error %v, want %v", err, leveldb.ErrNotFound)
		}
		want := ""
		if s.A != want {
			t.Errorf("got string %q, want %q", s.A, want)
		}
	})

	t.Run("put", func(t *testing.T) {
		want := complexStructure{
			A: "simple string value",
		}
		err = complexField.Put(want)
		if err != nil {
			t.Fatal(err)
		}
		var got complexStructure
		err = complexField.Get(&got)
		if err != nil {
			t.Fatal(err)
		}
		if got.A != want.A {
			t.Errorf("got string %q, want %q", got.A, want.A)
		}

		t.Run("overwrite", func(t *testing.T) {
			want := complexStructure{
				A: "overwritten string value",
			}
			err = complexField.Put(want)
			if err != nil {
				t.Fatal(err)
			}
			var got complexStructure
			err = complexField.Get(&got)
			if err != nil {
				t.Fatal(err)
			}
			if got.A != want.A {
				t.Errorf("got string %q, want %q", got.A, want.A)
			}
		})
	})

	t.Run("put in batch", func(t *testing.T) {
		batch := new(leveldb.Batch)
		want := complexStructure{
			A: "simple string batch value",
		}
		complexField.PutInBatch(batch, want)
		err = db.WriteBatch(batch)
		if err != nil {
			t.Fatal(err)
		}
		var got complexStructure
		err := complexField.Get(&got)
		if err != nil {
			t.Fatal(err)
		}
		if got.A != want.A {
			t.Errorf("got string %q, want %q", got, want)
		}

		t.Run("overwrite", func(t *testing.T) {
			batch := new(leveldb.Batch)
			want := complexStructure{
				A: "overwritten string batch value",
			}
			complexField.PutInBatch(batch, want)
			err = db.WriteBatch(batch)
			if err != nil {
				t.Fatal(err)
			}
			var got complexStructure
			err := complexField.Get(&got)
			if err != nil {
				t.Fatal(err)
			}
			if got.A != want.A {
				t.Errorf("got string %q, want %q", got, want)
			}
		})
	})
}

