
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450104371187712>


package enode

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enr"
)

func newLocalNodeForTesting() (*LocalNode, *DB) {
	db, _ := OpenDB("")
	key, _ := crypto.GenerateKey()
	return NewLocalNode(db, key), db
}

func TestLocalNode(t *testing.T) {
	ln, db := newLocalNodeForTesting()
	defer db.Close()

	if ln.Node().ID() != ln.ID() {
		t.Fatal("inconsistent ID")
	}

	ln.Set(enr.WithEntry("x", uint(3)))
	var x uint
	if err := ln.Node().Load(enr.WithEntry("x", &x)); err != nil {
		t.Fatal("can't load entry 'x':", err)
	} else if x != 3 {
		t.Fatal("wrong value for entry 'x':", x)
	}
}

func TestLocalNodeSeqPersist(t *testing.T) {
	ln, db := newLocalNodeForTesting()
	defer db.Close()

	if s := ln.Node().Seq(); s != 1 {
		t.Fatalf("wrong initial seq %d, want 1", s)
	}
	ln.Set(enr.WithEntry("x", uint(1)))
	if s := ln.Node().Seq(); s != 2 {
		t.Fatalf("wrong seq %d after set, want 2", s)
	}

//创建一个新实例，它应该重新加载序列号。
//因为新的记录是
//创建时不带“X”项。
	ln2 := NewLocalNode(db, ln.key)
	if s := ln2.Node().Seq(); s != 3 {
		t.Fatalf("wrong seq %d on new instance, want 3", s)
	}

//在同一数据库上使用不同的节点键创建新实例。
//这将重置序列号。
	key, _ := crypto.GenerateKey()
	ln3 := NewLocalNode(db, key)
	if s := ln3.Node().Seq(); s != 1 {
		t.Fatalf("wrong seq %d on instance with changed key, want 1", s)
	}
}

