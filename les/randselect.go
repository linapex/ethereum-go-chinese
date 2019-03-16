
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450095894499328>


//包les实现轻以太坊子协议。
package les

import (
	"math/rand"
)

//WRSitem接口应由要从中选择的任何条目实现
//加权随机选择集。注意，重新计算单调递减项
//允许按需重量（无需不断调用更新）
type wrsItem interface {
	Weight() int64
}

//WeightedRandomSelect能够从一组项目中对随机选择进行加权
type weightedRandomSelect struct {
	root *wrsNode
	idx  map[wrsItem]int
}

//new weightedrandomselect返回新的weightedrandomselect结构
func newWeightedRandomSelect() *weightedRandomSelect {
	return &weightedRandomSelect{root: &wrsNode{maxItems: wrsBranches}, idx: make(map[wrsItem]int)}
}

//更新更新更新项目的权重，如果不存在则添加该权重，如果
//新的重量是零。请注意，不需要显式更新递减权重。
func (w *weightedRandomSelect) update(item wrsItem) {
	w.setWeight(item, item.Weight())
}

//移除从集合中移除项
func (w *weightedRandomSelect) remove(item wrsItem) {
	w.setWeight(item, 0)
}

//setweight将项目的权重设置为特定值（如果为零，则移除该值）
func (w *weightedRandomSelect) setWeight(item wrsItem, weight int64) {
	idx, ok := w.idx[item]
	if ok {
		w.root.setWeight(idx, weight)
		if weight == 0 {
			delete(w.idx, item)
		}
	} else {
		if weight != 0 {
			if w.root.itemCnt == w.root.maxItems {
//添加新的级别
				newRoot := &wrsNode{sumWeight: w.root.sumWeight, itemCnt: w.root.itemCnt, level: w.root.level + 1, maxItems: w.root.maxItems * wrsBranches}
				newRoot.items[0] = w.root
				newRoot.weights[0] = w.root.sumWeight
				w.root = newRoot
			}
			w.idx[item] = w.root.insert(item, weight)
		}
	}
}

//随机选择从集合中选择一个项目，其机会与其
//当前重量。如果所选元素的重量自
//最后一个存储值，以newweight/oldweight的概率返回它，否则
//更新其权重并选择另一个权重
func (w *weightedRandomSelect) choose() wrsItem {
	for {
		if w.root.sumWeight == 0 {
			return nil
		}
		val := rand.Int63n(w.root.sumWeight)
		choice, lastWeight := w.root.choose(val)
		weight := choice.Weight()
		if weight != lastWeight {
			w.setWeight(choice, weight)
		}
		if weight >= lastWeight || rand.Int63n(lastWeight) < weight {
			return choice
		}
	}
}

const wrsBranches = 8 //Wrsnode树中的最大分支数

//wrsnode是树结构的一个节点，可以存储wrsitems或其他wrsnodes。
type wrsNode struct {
	items                    [wrsBranches]interface{}
	weights                  [wrsBranches]int64
	sumWeight                int64
	level, itemCnt, maxItems int
}

//递归插入将新项插入树并返回项索引
func (n *wrsNode) insert(item wrsItem, weight int64) int {
	branch := 0
	for n.items[branch] != nil && (n.level == 0 || n.items[branch].(*wrsNode).itemCnt == n.items[branch].(*wrsNode).maxItems) {
		branch++
		if branch == wrsBranches {
			panic(nil)
		}
	}
	n.itemCnt++
	n.sumWeight += weight
	n.weights[branch] += weight
	if n.level == 0 {
		n.items[branch] = item
		return branch
	}
	var subNode *wrsNode
	if n.items[branch] == nil {
		subNode = &wrsNode{maxItems: n.maxItems / wrsBranches, level: n.level - 1}
		n.items[branch] = subNode
	} else {
		subNode = n.items[branch].(*wrsNode)
	}
	subIdx := subNode.insert(item, weight)
	return subNode.maxItems*branch + subIdx
}

//setweight更新某个项目（应该存在）的权重并返回
//存储在树中的最后一个权重值的更改
func (n *wrsNode) setWeight(idx int, weight int64) int64 {
	if n.level == 0 {
		oldWeight := n.weights[idx]
		n.weights[idx] = weight
		diff := weight - oldWeight
		n.sumWeight += diff
		if weight == 0 {
			n.items[idx] = nil
			n.itemCnt--
		}
		return diff
	}
	branchItems := n.maxItems / wrsBranches
	branch := idx / branchItems
	diff := n.items[branch].(*wrsNode).setWeight(idx-branch*branchItems, weight)
	n.weights[branch] += diff
	n.sumWeight += diff
	if weight == 0 {
		n.itemCnt--
	}
	return diff
}

//递归选择从树中选择一个项并返回其权重
func (n *wrsNode) choose(val int64) (wrsItem, int64) {
	for i, w := range n.weights {
		if val < w {
			if n.level == 0 {
				return n.items[i].(wrsItem), n.weights[i]
			}
			return n.items[i].(*wrsNode).choose(val)
		}
		val -= w
	}
	panic(nil)
}

