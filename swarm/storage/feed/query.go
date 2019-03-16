
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450119260966912>


package feed

import (
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage/feed/lookup"
)

//查询用于在执行更新查找时指定约束
//TimeLimit表示搜索的上限。设置为0表示“现在”
type Query struct {
	Feed
	Hint      lookup.Epoch
	TimeLimit uint64
}

//FromValues从字符串键值存储中反序列化此实例
//用于分析查询字符串
func (q *Query) FromValues(values Values) error {
	time, _ := strconv.ParseUint(values.Get("time"), 10, 64)
	q.TimeLimit = time

	level, _ := strconv.ParseUint(values.Get("hint.level"), 10, 32)
	q.Hint.Level = uint8(level)
	q.Hint.Time, _ = strconv.ParseUint(values.Get("hint.time"), 10, 64)
	if q.Feed.User == (common.Address{}) {
		return q.Feed.FromValues(values)
	}
	return nil
}

//AppendValues将此结构序列化到提供的字符串键值存储区中
//用于生成查询字符串
func (q *Query) AppendValues(values Values) {
	if q.TimeLimit != 0 {
		values.Set("time", fmt.Sprintf("%d", q.TimeLimit))
	}
	if q.Hint.Level != 0 {
		values.Set("hint.level", fmt.Sprintf("%d", q.Hint.Level))
	}
	if q.Hint.Time != 0 {
		values.Set("hint.time", fmt.Sprintf("%d", q.Hint.Time))
	}
	q.Feed.AppendValues(values)
}

//newquery构造一个查询结构以在“time”或之前查找更新
//如果time==0，将查找最新更新
func NewQuery(feed *Feed, time uint64, hint lookup.Epoch) *Query {
	return &Query{
		TimeLimit: time,
		Feed:      *feed,
		Hint:      hint,
	}
}

//newquerylatest生成查找参数，以查找源的最新更新
func NewQueryLatest(feed *Feed, hint lookup.Epoch) *Query {
	return NewQuery(feed, 0, hint)
}

