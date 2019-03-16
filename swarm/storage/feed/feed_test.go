
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450118837342208>

package feed

import (
	"testing"
)

func getTestFeed() *Feed {
	topic, _ := NewTopic("world news report, every hour", nil)
	return &Feed{
		Topic: topic,
		User:  newCharlieSigner().Address(),
	}
}

func TestFeedSerializerDeserializer(t *testing.T) {
	testBinarySerializerRecovery(t, getTestFeed(), "0x776f726c64206e657773207265706f72742c20657665727920686f7572000000876a8936a7cd0b79ef0735ad0896c1afe278781c")
}

func TestFeedSerializerLengthCheck(t *testing.T) {
	testBinarySerializerLengthCheck(t, getTestFeed())
}

