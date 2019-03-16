
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450077204680704>


package asm

import (
	"testing"

	"encoding/hex"
)

//测试拆解有效EVM代码的说明
func TestInstructionIteratorValid(t *testing.T) {
	cnt := 0
	script, _ := hex.DecodeString("61000000")

	it := NewInstructionIterator(script)
	for it.Next() {
		cnt++
	}

	if err := it.Error(); err != nil {
		t.Errorf("Expected 2, but encountered error %v instead.", err)
	}
	if cnt != 2 {
		t.Errorf("Expected 2, but got %v instead.", cnt)
	}
}

//测试反汇编无效EVM代码的指令
func TestInstructionIteratorInvalid(t *testing.T) {
	cnt := 0
	script, _ := hex.DecodeString("6100")

	it := NewInstructionIterator(script)
	for it.Next() {
		cnt++
	}

	if it.Error() == nil {
		t.Errorf("Expected an error, but got %v instead.", cnt)
	}
}

//测试反汇编空EVM代码的说明
func TestInstructionIteratorEmpty(t *testing.T) {
	cnt := 0
	script, _ := hex.DecodeString("")

	it := NewInstructionIterator(script)
	for it.Next() {
		cnt++
	}

	if err := it.Error(); err != nil {
		t.Errorf("Expected 0, but encountered error %v instead.", err)
	}
	if cnt != 0 {
		t.Errorf("Expected 0, but got %v instead.", cnt)
	}
}

