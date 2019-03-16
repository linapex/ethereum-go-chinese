
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450101212876800>


//包含go-ethereum根包中的所有包装。

package geth

import (
	"errors"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

//订阅表示事件订阅，其中
//通过数据通道传送。
type Subscription struct {
	sub ethereum.Subscription
}

//取消订阅取消向数据通道发送事件
//关闭错误通道。
func (s *Subscription) Unsubscribe() {
	s.sub.Unsubscribe()
}

//callmsg包含合同调用的参数。
type CallMsg struct {
	msg ethereum.CallMsg
}

//newcallmsg创建一个空的合同调用参数列表。
func NewCallMsg() *CallMsg {
	return new(CallMsg)
}

func (msg *CallMsg) GetFrom() *Address    { return &Address{msg.msg.From} }
func (msg *CallMsg) GetGas() int64        { return int64(msg.msg.Gas) }
func (msg *CallMsg) GetGasPrice() *BigInt { return &BigInt{msg.msg.GasPrice} }
func (msg *CallMsg) GetValue() *BigInt    { return &BigInt{msg.msg.Value} }
func (msg *CallMsg) GetData() []byte      { return msg.msg.Data }
func (msg *CallMsg) GetTo() *Address {
	if to := msg.msg.To; to != nil {
		return &Address{*msg.msg.To}
	}
	return nil
}

func (msg *CallMsg) SetFrom(address *Address)  { msg.msg.From = address.address }
func (msg *CallMsg) SetGas(gas int64)          { msg.msg.Gas = uint64(gas) }
func (msg *CallMsg) SetGasPrice(price *BigInt) { msg.msg.GasPrice = price.bigint }
func (msg *CallMsg) SetValue(value *BigInt)    { msg.msg.Value = value.bigint }
func (msg *CallMsg) SetData(data []byte)       { msg.msg.Data = common.CopyBytes(data) }
func (msg *CallMsg) SetTo(address *Address) {
	if address == nil {
		msg.msg.To = nil
		return
	}
	msg.msg.To = &address.address
}

//当节点与
//以太坊网络。
type SyncProgress struct {
	progress ethereum.SyncProgress
}

func (p *SyncProgress) GetStartingBlock() int64 { return int64(p.progress.StartingBlock) }
func (p *SyncProgress) GetCurrentBlock() int64  { return int64(p.progress.CurrentBlock) }
func (p *SyncProgress) GetHighestBlock() int64  { return int64(p.progress.HighestBlock) }
func (p *SyncProgress) GetPulledStates() int64  { return int64(p.progress.PulledStates) }
func (p *SyncProgress) GetKnownStates() int64   { return int64(p.progress.KnownStates) }

//主题是一组用于筛选事件的主题列表。
type Topics struct{ topics [][]common.Hash }

//newtopics创建一个未初始化主题的切片。
func NewTopics(size int) *Topics {
	return &Topics{
		topics: make([][]common.Hash, size),
	}
}

//newtopicsempty创建主题值的空切片。
func NewTopicsEmpty() *Topics {
	return NewTopics(0)
}

//SIZE返回集合内主题列表的数目
func (t *Topics) Size() int {
	return len(t.topics)
}

//get从切片返回给定索引处的主题列表。
func (t *Topics) Get(index int) (hashes *Hashes, _ error) {
	if index < 0 || index >= len(t.topics) {
		return nil, errors.New("index out of bounds")
	}
	return &Hashes{t.topics[index]}, nil
}

//set在切片中的给定索引处设置主题列表。
func (t *Topics) Set(index int, topics *Hashes) error {
	if index < 0 || index >= len(t.topics) {
		return errors.New("index out of bounds")
	}
	t.topics[index] = topics.hashes
	return nil
}

//附加将新的主题列表添加到切片的末尾。
func (t *Topics) Append(topics *Hashes) {
	t.topics = append(t.topics, topics.hashes)
}

//filterquery包含用于合同日志筛选的选项。
type FilterQuery struct {
	query ethereum.FilterQuery
}

//newfilterquery为合同日志筛选创建空的筛选器查询。
func NewFilterQuery() *FilterQuery {
	return new(FilterQuery)
}

func (fq *FilterQuery) GetFromBlock() *BigInt    { return &BigInt{fq.query.FromBlock} }
func (fq *FilterQuery) GetToBlock() *BigInt      { return &BigInt{fq.query.ToBlock} }
func (fq *FilterQuery) GetAddresses() *Addresses { return &Addresses{fq.query.Addresses} }
func (fq *FilterQuery) GetTopics() *Topics       { return &Topics{fq.query.Topics} }

func (fq *FilterQuery) SetFromBlock(fromBlock *BigInt)    { fq.query.FromBlock = fromBlock.bigint }
func (fq *FilterQuery) SetToBlock(toBlock *BigInt)        { fq.query.ToBlock = toBlock.bigint }
func (fq *FilterQuery) SetAddresses(addresses *Addresses) { fq.query.Addresses = addresses.addresses }
func (fq *FilterQuery) SetTopics(topics *Topics)          { fq.query.Topics = topics.topics }

