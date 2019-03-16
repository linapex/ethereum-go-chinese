
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450077393424384>


package asm

import (
	"reflect"
	"testing"
)

func lexAll(src string) []token {
	ch := Lex("test.asm", []byte(src), false)

	var tokens []token
	for i := range ch {
		tokens = append(tokens, i)
	}
	return tokens
}

func TestLexer(t *testing.T) {
	tests := []struct {
		input  string
		tokens []token
	}{
		{
			input:  ";; this is a comment",
			tokens: []token{{typ: lineStart}, {typ: eof}},
		},
		{
			input:  "0x12345678",
			tokens: []token{{typ: lineStart}, {typ: number, text: "0x12345678"}, {typ: eof}},
		},
		{
			input:  "0x123ggg",
			tokens: []token{{typ: lineStart}, {typ: number, text: "0x123"}, {typ: element, text: "ggg"}, {typ: eof}},
		},
		{
			input:  "12345678",
			tokens: []token{{typ: lineStart}, {typ: number, text: "12345678"}, {typ: eof}},
		},
		{
			input:  "123abc",
			tokens: []token{{typ: lineStart}, {typ: number, text: "123"}, {typ: element, text: "abc"}, {typ: eof}},
		},
		{
			input:  "0123abc",
			tokens: []token{{typ: lineStart}, {typ: number, text: "0123"}, {typ: element, text: "abc"}, {typ: eof}},
		},
	}

	for _, test := range tests {
		tokens := lexAll(test.input)
		if !reflect.DeepEqual(tokens, test.tokens) {
			t.Errorf("input %q\ngot:  %+v\nwant: %+v", test.input, tokens, test.tokens)
		}
	}
}

