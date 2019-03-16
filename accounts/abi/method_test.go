
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:31</date>
//</624450062130352128>


package abi

import (
	"strings"
	"testing"
)

const methoddata = `
[
	{ "type" : "function", "name" : "balance", "constant" : true },
	{ "type" : "function", "name" : "send", "constant" : false, "inputs" : [ { "name" : "amount", "type" : "uint256" } ] },
	{ "type" : "function", "name" : "transfer", "constant" : false, "inputs" : [ { "name" : "from", "type" : "address" }, { "name" : "to", "type" : "address" }, { "name" : "value", "type" : "uint256" } ], "outputs" : [ { "name" : "success", "type" : "bool" } ]  }
]`

func TestMethodString(t *testing.T) {
	var table = []struct {
		method      string
		expectation string
	}{
		{
			method:      "balance",
			expectation: "function balance() constant returns()",
		},
		{
			method:      "send",
			expectation: "function send(uint256 amount) returns()",
		},
		{
			method:      "transfer",
			expectation: "function transfer(address from, address to, uint256 value) returns(bool success)",
		},
	}

	abi, err := JSON(strings.NewReader(methoddata))
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range table {
		got := abi.Methods[test.method].String()
		if got != test.expectation {
			t.Errorf("expected string to be %s, got %s", test.expectation, got)
		}
	}
}

