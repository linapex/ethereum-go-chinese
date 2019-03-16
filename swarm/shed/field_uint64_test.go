
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450117847486464>


package shed

import (
	"testing"

	"github.com/syndtr/goleveldb/leveldb"
)

//testuint64字段验证Put和Get操作
//在uint64字段中。
func TestUint64Field(t *testing.T) {
	db, cleanupFunc := newTestDB(t)
	defer cleanupFunc()

	counter, err := db.NewUint64Field("counter")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("get empty", func(t *testing.T) {
		got, err := counter.Get()
		if err != nil {
			t.Fatal(err)
		}
		var want uint64
		if got != want {
			t.Errorf("got uint64 %v, want %v", got, want)
		}
	})

	t.Run("put", func(t *testing.T) {
		var want uint64 = 42
		err = counter.Put(want)
		if err != nil {
			t.Fatal(err)
		}
		got, err := counter.Get()
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("got uint64 %v, want %v", got, want)
		}

		t.Run("overwrite", func(t *testing.T) {
			var want uint64 = 84
			err = counter.Put(want)
			if err != nil {
				t.Fatal(err)
			}
			got, err := counter.Get()
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Errorf("got uint64 %v, want %v", got, want)
			}
		})
	})

	t.Run("put in batch", func(t *testing.T) {
		batch := new(leveldb.Batch)
		var want uint64 = 42
		counter.PutInBatch(batch, want)
		err = db.WriteBatch(batch)
		if err != nil {
			t.Fatal(err)
		}
		got, err := counter.Get()
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("got uint64 %v, want %v", got, want)
		}

		t.Run("overwrite", func(t *testing.T) {
			batch := new(leveldb.Batch)
			var want uint64 = 84
			counter.PutInBatch(batch, want)
			err = db.WriteBatch(batch)
			if err != nil {
				t.Fatal(err)
			}
			got, err := counter.Get()
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Errorf("got uint64 %v, want %v", got, want)
			}
		})
	})
}

//TESTUINT64字段_inc验证inc操作
//在uint64字段中。
func TestUint64Field_Inc(t *testing.T) {
	db, cleanupFunc := newTestDB(t)
	defer cleanupFunc()

	counter, err := db.NewUint64Field("counter")
	if err != nil {
		t.Fatal(err)
	}

	var want uint64 = 1
	got, err := counter.Inc()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}

	want = 2
	got, err = counter.Inc()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}
}

//testuint64字段“incinbatch”验证incinbatch操作
//在uint64字段中。
func TestUint64Field_IncInBatch(t *testing.T) {
	db, cleanupFunc := newTestDB(t)
	defer cleanupFunc()

	counter, err := db.NewUint64Field("counter")
	if err != nil {
		t.Fatal(err)
	}

	batch := new(leveldb.Batch)
	var want uint64 = 1
	got, err := counter.IncInBatch(batch)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}
	err = db.WriteBatch(batch)
	if err != nil {
		t.Fatal(err)
	}
	got, err = counter.Get()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}

	batch2 := new(leveldb.Batch)
	want = 2
	got, err = counter.IncInBatch(batch2)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}
	err = db.WriteBatch(batch2)
	if err != nil {
		t.Fatal(err)
	}
	got, err = counter.Get()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}
}

//TESTUINT64字段“DEC验证DEC操作”
//在uint64字段中。
func TestUint64Field_Dec(t *testing.T) {
	db, cleanupFunc := newTestDB(t)
	defer cleanupFunc()

	counter, err := db.NewUint64Field("counter")
	if err != nil {
		t.Fatal(err)
	}

//测试溢出保护
	var want uint64
	got, err := counter.Dec()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}

	want = 32
	err = counter.Put(want)
	if err != nil {
		t.Fatal(err)
	}

	want = 31
	got, err = counter.Dec()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}
}

//testuint64字段“decinbatch”验证decinbatch操作
//在uint64字段中。
func TestUint64Field_DecInBatch(t *testing.T) {
	db, cleanupFunc := newTestDB(t)
	defer cleanupFunc()

	counter, err := db.NewUint64Field("counter")
	if err != nil {
		t.Fatal(err)
	}

	batch := new(leveldb.Batch)
	var want uint64
	got, err := counter.DecInBatch(batch)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}
	err = db.WriteBatch(batch)
	if err != nil {
		t.Fatal(err)
	}
	got, err = counter.Get()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}

	batch2 := new(leveldb.Batch)
	want = 42
	counter.PutInBatch(batch2, want)
	err = db.WriteBatch(batch2)
	if err != nil {
		t.Fatal(err)
	}
	got, err = counter.Get()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}

	batch3 := new(leveldb.Batch)
	want = 41
	got, err = counter.DecInBatch(batch3)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}
	err = db.WriteBatch(batch3)
	if err != nil {
		t.Fatal(err)
	}
	got, err = counter.Get()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got uint64 %v, want %v", got, want)
	}
}

