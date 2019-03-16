
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450077141766144>


//为处理EVM装配说明（例如，拆卸它们）提供支持。
package asm

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/core/vm"
)

//反汇编EVM指令的迭代器
type instructionIterator struct {
	code    []byte
	pc      uint64
	arg     []byte
	op      vm.OpCode
	error   error
	started bool
}

//
func NewInstructionIterator(code []byte) *instructionIterator {
	it := new(instructionIterator)
	it.code = code
	return it
}

//如果有下一条指令并继续，则返回true。
func (it *instructionIterator) Next() bool {
	if it.error != nil || uint64(len(it.code)) <= it.pc {
//我们以前遇到过一个错误或结束。
		return false
	}

	if it.started {
//由于迭代已经开始，我们将转到下一条指令。
		if it.arg != nil {
			it.pc += uint64(len(it.arg))
		}
		it.pc++
	} else {
//我们从第一条指令开始迭代。
		it.started = true
	}

	if uint64(len(it.code)) <= it.pc {
//我们到了终点。
		return false
	}

	it.op = vm.OpCode(it.code[it.pc])
	if it.op.IsPush() {
		a := uint64(it.op) - uint64(vm.PUSH1) + 1
		u := it.pc + 1 + a
		if uint64(len(it.code)) <= it.pc || uint64(len(it.code)) < u {
			it.error = fmt.Errorf("incomplete push instruction at %v", it.pc)
			return false
		}
		it.arg = it.code[it.pc+1 : u]
	} else {
		it.arg = nil
	}
	return true
}

//返回可能遇到的任何错误。
func (it *instructionIterator) Error() error {
	return it.error
}

//返回当前指令的PC。
func (it *instructionIterator) PC() uint64 {
	return it.pc
}

//返回当前指令的操作码。
func (it *instructionIterator) Op() vm.OpCode {
	return it.op
}

//返回当前指令的参数。
func (it *instructionIterator) Arg() []byte {
	return it.arg
}

//将所有反汇编的EVM指令打印到stdout。
func PrintDisassembled(code string) error {
	script, err := hex.DecodeString(code)
	if err != nil {
		return err
	}

	it := NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			fmt.Printf("%05x: %v 0x%x\n", it.PC(), it.Op(), it.Arg())
		} else {
			fmt.Printf("%05x: %v\n", it.PC(), it.Op())
		}
	}
	return it.Error()
}

//以可读的格式返回所有反汇编的EVM指令。
func Disassemble(script []byte) ([]string, error) {
	instrs := make([]string, 0)

	it := NewInstructionIterator(script)
	for it.Next() {
		if it.Arg() != nil && 0 < len(it.Arg()) {
			instrs = append(instrs, fmt.Sprintf("%05x: %v 0x%x\n", it.PC(), it.Op(), it.Arg()))
		} else {
			instrs = append(instrs, fmt.Sprintf("%05x: %v\n", it.PC(), it.Op()))
		}
	}
	if err := it.Error(); err != nil {
		return nil, err
	}
	return instrs, nil
}

