
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450118640209920>


package feed

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

//Kv模拟密钥值存储
type KV map[string]string

func (kv KV) Get(key string) string {
	return kv[key]
}
func (kv KV) Set(key, value string) {
	kv[key] = value
}

func compareByteSliceToExpectedHex(t *testing.T, variableName string, actualValue []byte, expectedHex string) {
	if hexutil.Encode(actualValue) != expectedHex {
		t.Fatalf("%s: Expected %s to be %s, got %s", t.Name(), variableName, expectedHex, hexutil.Encode(actualValue))
	}
}

func testBinarySerializerRecovery(t *testing.T, bin binarySerializer, expectedHex string) {
	name := reflect.TypeOf(bin).Elem().Name()
	serialized := make([]byte, bin.binaryLength())
	if err := bin.binaryPut(serialized); err != nil {
		t.Fatalf("%s.binaryPut error when trying to serialize structure: %s", name, err)
	}

	compareByteSliceToExpectedHex(t, name, serialized, expectedHex)

	recovered := reflect.New(reflect.TypeOf(bin).Elem()).Interface().(binarySerializer)
	if err := recovered.binaryGet(serialized); err != nil {
		t.Fatalf("%s.binaryGet error when trying to deserialize structure: %s", name, err)
	}

	if !reflect.DeepEqual(bin, recovered) {
		t.Fatalf("Expected that the recovered %s equals the marshalled %s", name, name)
	}

	serializedWrongLength := make([]byte, 1)
	copy(serializedWrongLength[:], serialized)
	if err := recovered.binaryGet(serializedWrongLength); err == nil {
		t.Fatalf("Expected %s.binaryGet to fail since data is too small", name)
	}
}

func testBinarySerializerLengthCheck(t *testing.T, bin binarySerializer) {
	name := reflect.TypeOf(bin).Elem().Name()
//使切片太小，无法包含元数据
	serialized := make([]byte, bin.binaryLength()-1)

	if err := bin.binaryPut(serialized); err == nil {
		t.Fatalf("Expected %s.binaryPut to fail, since target slice is too small", name)
	}
}

func testValueSerializer(t *testing.T, v valueSerializer, expected KV) {
	name := reflect.TypeOf(v).Elem().Name()
	kv := make(KV)

	v.AppendValues(kv)
	if !reflect.DeepEqual(expected, kv) {
		expj, _ := json.Marshal(expected)
		gotj, _ := json.Marshal(kv)
		t.Fatalf("Expected %s.AppendValues to return %s, got %s", name, string(expj), string(gotj))
	}

	recovered := reflect.New(reflect.TypeOf(v).Elem()).Interface().(valueSerializer)
	err := recovered.FromValues(kv)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(recovered, v) {
		t.Fatalf("Expected recovered %s to be the same", name)
	}
}

