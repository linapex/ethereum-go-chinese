
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450119294521344>


package feed

import (
	"testing"
)

func getTestQuery() *Query {
	id := getTestID()
	return &Query{
		TimeLimit: 5000,
		Feed:      id.Feed,
		Hint:      id.Epoch,
	}
}

func TestQueryValues(t *testing.T) {
	var expected = KV{"hint.level": "25", "hint.time": "1000", "time": "5000", "topic": "0x776f726c64206e657773207265706f72742c20657665727920686f7572000000", "user": "0x876A8936A7Cd0b79Ef0735AD0896c1AFe278781c"}

	query := getTestQuery()
	testValueSerializer(t, query, expected)

}

