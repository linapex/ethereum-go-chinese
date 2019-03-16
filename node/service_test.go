
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450102689271808>


package node

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

//测试数据库是否基于
//配置的服务上下文。
func TestContextDatabases(t *testing.T) {
//创建临时文件夹并确保其中不包含任何数据库
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("failed to create temporary data directory: %v", err)
	}
	defer os.RemoveAll(dir)

	if _, err := os.Stat(filepath.Join(dir, "database")); err == nil {
		t.Fatalf("non-created database already exists")
	}
//请求打开/创建一个数据库，并确保它一直存在于磁盘上
	ctx := &ServiceContext{config: &Config{Name: "unit-test", DataDir: dir}}
	db, err := ctx.OpenDatabase("persistent", 0, 0)
	if err != nil {
		t.Fatalf("failed to open persistent database: %v", err)
	}
	db.Close()

	if _, err := os.Stat(filepath.Join(dir, "unit-test", "persistent")); err != nil {
		t.Fatalf("persistent database doesn't exists: %v", err)
	}
//请求打开/创建临时数据库，并确保它不会持久化。
	ctx = &ServiceContext{config: &Config{DataDir: ""}}
	db, err = ctx.OpenDatabase("ephemeral", 0, 0)
	if err != nil {
		t.Fatalf("failed to open ephemeral database: %v", err)
	}
	db.Close()

	if _, err := os.Stat(filepath.Join(dir, "ephemeral")); err == nil {
		t.Fatalf("ephemeral database exists")
	}
}

//已经构造了服务的测试可以由以后的测试检索。
func TestContextServices(t *testing.T) {
	stack, err := New(testNodeConfig())
	if err != nil {
		t.Fatalf("failed to create protocol stack: %v", err)
	}
//定义一个验证器，确保noopa在它之前，noopb在它之后。
	verifier := func(ctx *ServiceContext) (Service, error) {
		var objA *NoopServiceA
		if ctx.Service(&objA) != nil {
			return nil, fmt.Errorf("former service not found")
		}
		var objB *NoopServiceB
		if err := ctx.Service(&objB); err != ErrServiceUnknown {
			return nil, fmt.Errorf("latters lookup error mismatch: have %v, want %v", err, ErrServiceUnknown)
		}
		return new(NoopService), nil
	}
//注册服务集合
	if err := stack.Register(NewNoopServiceA); err != nil {
		t.Fatalf("former failed to register service: %v", err)
	}
	if err := stack.Register(verifier); err != nil {
		t.Fatalf("failed to register service verifier: %v", err)
	}
	if err := stack.Register(NewNoopServiceB); err != nil {
		t.Fatalf("latter failed to register service: %v", err)
	}
//启动协议栈并确保服务按顺序构造
	if err := stack.Start(); err != nil {
		t.Fatalf("failed to start stack: %v", err)
	}
	defer stack.Stop()
}

