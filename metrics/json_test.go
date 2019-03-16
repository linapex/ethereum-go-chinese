
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450099220582400>

package metrics

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestRegistryMarshallJSON(t *testing.T) {
	b := &bytes.Buffer{}
	enc := json.NewEncoder(b)
	r := NewRegistry()
	r.Register("counter", NewCounter())
	enc.Encode(r)
	if s := b.String(); "{\"counter\":{\"count\":0}}\n" != s {
		t.Fatalf(s)
	}
}

func TestRegistryWriteJSONOnce(t *testing.T) {
	r := NewRegistry()
	r.Register("counter", NewCounter())
	b := &bytes.Buffer{}
	WriteJSONOnce(r, b)
	if s := b.String(); s != "{\"counter\":{\"count\":0}}\n" {
		t.Fail()
	}
}

