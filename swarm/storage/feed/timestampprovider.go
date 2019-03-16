
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450119500042240>


package feed

import (
	"encoding/json"
	"time"
)

//TimestampProvider设置源包的时间源
var TimestampProvider timestampProvider = NewDefaultTimestampProvider()

//timestamp将时间点编码为unix epoch
type Timestamp struct {
Time uint64 `json:"time"` //unix epoch时间戳（秒）
}

//TimestampProvider接口描述时间戳信息的来源
type timestampProvider interface {
Now() Timestamp //返回当前时间戳信息
}

//unmashaljson实现json.unmarshaller接口
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Time)
}

//marshaljson实现json.marshaller接口
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time)
}

//DefaultTimestampProvider是使用系统时间的TimestampProvider
//作为时间来源
type DefaultTimestampProvider struct {
}

//NewDefaultTimestampProvider创建基于系统时钟的时间戳提供程序
func NewDefaultTimestampProvider() *DefaultTimestampProvider {
	return &DefaultTimestampProvider{}
}

//现在根据此提供程序返回当前时间
func (dtp *DefaultTimestampProvider) Now() Timestamp {
	return Timestamp{
		Time: uint64(time.Now().Unix()),
	}
}

