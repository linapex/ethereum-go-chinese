
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:34</date>
//</624450073949900800>


package math

import (
	"fmt"
	"strconv"
)

//整数限制值。
const (
	MaxInt8   = 1<<7 - 1
	MinInt8   = -1 << 7
	MaxInt16  = 1<<15 - 1
	MinInt16  = -1 << 15
	MaxInt32  = 1<<31 - 1
	MinInt32  = -1 << 31
	MaxInt64  = 1<<63 - 1
	MinInt64  = -1 << 63
	MaxUint8  = 1<<8 - 1
	MaxUint16 = 1<<16 - 1
	MaxUint32 = 1<<32 - 1
	MaxUint64 = 1<<64 - 1
)

//hexordecimal64将uint64封送为十六进制或十进制。
type HexOrDecimal64 uint64

//UnmarshalText实现encoding.textUnmarshaller。
func (i *HexOrDecimal64) UnmarshalText(input []byte) error {
	int, ok := ParseUint64(string(input))
	if !ok {
		return fmt.Errorf("invalid hex or decimal integer %q", input)
	}
	*i = HexOrDecimal64(int)
	return nil
}

//MarshalText实现Encoding.TextMarshaler。
func (i HexOrDecimal64) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%#x", uint64(i))), nil
}

//ParseUInt64以十进制或十六进制语法解析为整数。
//可接受前导零。空字符串解析为零。
func ParseUint64(s string) (uint64, bool) {
	if s == "" {
		return 0, true
	}
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		v, err := strconv.ParseUint(s[2:], 16, 64)
		return v, err == nil
	}
	v, err := strconv.ParseUint(s, 10, 64)
	return v, err == nil
}

//mustParseUInt64作为整数进行分析，如果字符串无效，则会恐慌。
func MustParseUint64(s string) uint64 {
	v, ok := ParseUint64(s)
	if !ok {
		panic("invalid unsigned 64 bit integer: " + s)
	}
	return v
}

//注意：以下方法需要使用位检查或ASM进行优化

//SAFESUB返回减法结果以及是否发生溢出。
func SafeSub(x, y uint64) (uint64, bool) {
	return x - y, x < y
}

//safeadd返回结果以及是否发生溢出。
func SafeAdd(x, y uint64) (uint64, bool) {
	return x + y, y > MaxUint64-x
}

//safemul返回乘法结果以及是否发生溢出。
func SafeMul(x, y uint64) (uint64, bool) {
	if x == 0 || y == 0 {
		return 0, false
	}
	return x * y, y > MaxUint64/x
}

