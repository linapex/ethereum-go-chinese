
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450074163810304>

//这是“gopkg.in/karalabe/cookiejar.v2/collections/prque”的一个复制和稍加修改的版本。

package prque

import (
	"container/heap"
)

//优先级队列数据结构。
type Prque struct {
	cont *sstack
}

//创建新的优先级队列。
func New(setIndex setIndexCallback) *Prque {
	return &Prque{newSstack(setIndex)}
}

//将具有给定优先级的值推入队列，必要时展开。
func (p *Prque) Push(data interface{}, priority int64) {
	heap.Push(p.cont, &item{data, priority})
}

//从堆栈中弹出优先级为greates的值并返回该值。
//目前还没有收缩。
func (p *Prque) Pop() (interface{}, int64) {
	item := heap.Pop(p.cont).(*item)
	return item.value, item.priority
}

//只从队列中弹出项目，删除关联的优先级值。
func (p *Prque) PopItem() interface{} {
	return heap.Pop(p.cont).(*item).value
}

//移除移除具有给定索引的元素。
func (p *Prque) Remove(i int) interface{} {
	if i < 0 {
		return nil
	}
	return heap.Remove(p.cont, i)
}

//检查优先级队列是否为空。
func (p *Prque) Empty() bool {
	return p.cont.Len() == 0
}

//返回优先级队列中的元素数。
func (p *Prque) Size() int {
	return p.cont.Len()
}

//清除优先级队列的内容。
func (p *Prque) Reset() {
	*p = *New(p.cont.setIndex)
}

