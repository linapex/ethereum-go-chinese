
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450117503553536>


package shed

import (
	"io/ioutil"
	"os"
	"testing"
)

//testnewdb构造新的db
//并验证架构是否正确初始化。
func TestNewDB(t *testing.T) {
	db, cleanupFunc := newTestDB(t)
	defer cleanupFunc()

	s, err := db.getSchema()
	if err != nil {
		t.Fatal(err)
	}
	if s.Fields == nil {
		t.Error("schema fields are empty")
	}
	if len(s.Fields) != 0 {
		t.Errorf("got schema fields length %v, want %v", len(s.Fields), 0)
	}
	if s.Indexes == nil {
		t.Error("schema indexes are empty")
	}
	if len(s.Indexes) != 0 {
		t.Errorf("got schema indexes length %v, want %v", len(s.Indexes), 0)
	}
}

//testdb_persistence创建一个db，保存一个字段并关闭该db。
//然后，它构造另一个数据库并trues来检索保存的值。
func TestDB_persistence(t *testing.T) {
	dir, err := ioutil.TempDir("", "shed-test-persistence")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	db, err := NewDB(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	stringField, err := db.NewStringField("preserve-me")
	if err != nil {
		t.Fatal(err)
	}
	want := "persistent value"
	err = stringField.Put(want)
	if err != nil {
		t.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		t.Fatal(err)
	}

	db2, err := NewDB(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	stringField2, err := db2.NewStringField("preserve-me")
	if err != nil {
		t.Fatal(err)
	}
	got, err := stringField2.Get()
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("got string %q, want %q", got, want)
	}
}

//newtestdb是构造
//临时数据库并返回一个清理函数，该函数必须
//调用以删除数据。
func newTestDB(t *testing.T) (db *DB, cleanupFunc func()) {
	t.Helper()

	dir, err := ioutil.TempDir("", "shed-test")
	if err != nil {
		t.Fatal(err)
	}
	cleanupFunc = func() { os.RemoveAll(dir) }
	db, err = NewDB(dir, "")
	if err != nil {
		cleanupFunc()
		t.Fatal(err)
	}
	return db, cleanupFunc
}

