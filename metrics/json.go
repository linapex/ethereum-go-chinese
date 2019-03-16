
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450099191222272>

package metrics

import (
	"encoding/json"
	"io"
	"time"
)

//marshaljson返回一个字节片，其中包含所有
//注册表中的指标。
func (r *StandardRegistry) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.GetAll())
}

//WRITEJSON定期将给定注册表中的度量值写入
//指定IO.Writer为JSON。
func WriteJSON(r Registry, d time.Duration, w io.Writer) {
	for range time.Tick(d) {
		WriteJSONOnce(r, w)
	}
}

//writejsonce将给定注册表中的度量值写入指定的
//写JSON。
func WriteJSONOnce(r Registry, w io.Writer) {
	json.NewEncoder(w).Encode(r)
}

func (p *PrefixedRegistry) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.GetAll())
}

