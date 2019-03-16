
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450083387084800>


package vm

import (
	"fmt"
	"math/big"
)

//堆栈是用于基本堆栈操作的对象。弹出到堆栈的项是
//需要更改和修改。堆栈不负责新添加
//初始化的对象。
type Stack struct {
	data []*big.Int
}

func newstack() *Stack {
	return &Stack{data: make([]*big.Int, 0, 1024)}
}

//data返回基础的big.int数组。
func (st *Stack) Data() []*big.Int {
	return st.data
}

func (st *Stack) push(d *big.Int) {
//注：在basecheck中检查推送限制（1024）
//stackitem:=新建（big.int）.set（d）
//st.data=附加（st.data，stackitem）
	st.data = append(st.data, d)
}
func (st *Stack) pushN(ds ...*big.Int) {
	st.data = append(st.data, ds...)
}

func (st *Stack) pop() (ret *big.Int) {
	ret = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-1]
	return
}

func (st *Stack) len() int {
	return len(st.data)
}

func (st *Stack) swap(n int) {
	st.data[st.len()-n], st.data[st.len()-1] = st.data[st.len()-1], st.data[st.len()-n]
}

func (st *Stack) dup(pool *intPool, n int) {
	st.push(pool.get().Set(st.data[st.len()-n]))
}

func (st *Stack) peek() *big.Int {
	return st.data[st.len()-1]
}

//返回堆栈中的第n项
func (st *Stack) Back(n int) *big.Int {
	return st.data[st.len()-n-1]
}

func (st *Stack) require(n int) error {
	if st.len() < n {
		return fmt.Errorf("stack underflow (%d <=> %d)", len(st.data), n)
	}
	return nil
}

//打印转储堆栈的内容
func (st *Stack) Print() {
	fmt.Println("### stack ###")
	if len(st.data) > 0 {
		for i, val := range st.data {
			fmt.Printf("%-3d  %v\n", i, val)
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}

