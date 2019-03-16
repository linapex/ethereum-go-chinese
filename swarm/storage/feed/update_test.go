
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450119692980224>


package feed

import (
	"testing"
)

func getTestFeedUpdate() *Update {
	return &Update{
		ID:   *getTestID(),
		data: []byte("El que lee mucho y anda mucho, ve mucho y sabe mucho"),
	}
}

func TestUpdateSerializer(t *testing.T) {
	testBinarySerializerRecovery(t, getTestFeedUpdate(), "0x0000000000000000776f726c64206e657773207265706f72742c20657665727920686f7572000000876a8936a7cd0b79ef0735ad0896c1afe278781ce803000000000019456c20717565206c6565206d7563686f207920616e6461206d7563686f2c207665206d7563686f20792073616265206d7563686f")
}

func TestUpdateLengthCheck(t *testing.T) {
	testBinarySerializerLengthCheck(t, getTestFeedUpdate())
//如果更新太大，则测试失败
	update := getTestFeedUpdate()
	update.data = make([]byte, MaxUpdateDataLength+100)
	serialized := make([]byte, update.binaryLength())
	if err := update.binaryPut(serialized); err == nil {
		t.Fatal("Expected update.binaryPut to fail since update is too big")
	}

//如果数据为空或为零，则测试失败
	update.data = nil
	serialized = make([]byte, update.binaryLength())
	if err := update.binaryPut(serialized); err == nil {
		t.Fatal("Expected update.binaryPut to fail since data is empty")
	}
}

