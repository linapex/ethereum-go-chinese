
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450122406694912>


package tests

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"
)

//rlp test是单个rlp测试的JSON结构。
type RLPTest struct {
//如果in的值为“无效”或“有效”，则测试
//检查是否可以将out解码为值
//键入interface。
//
//对于其他JSON值，in被视为
//调用rlp.stream。测试还验证编码
//In产生输入输出字节。
	In interface{}

//out是十六进制编码的rlp值。
	Out string
}

//运行执行测试。
func (t *RLPTest) Run() error {
	outb, err := hex.DecodeString(t.Out)
	if err != nil {
		return fmt.Errorf("invalid hex in Out")
	}

//处理没有实际值的简单解码测试。
	if t.In == "VALID" || t.In == "INVALID" {
		return checkDecodeInterface(outb, t.In == "VALID")
	}

//检查编码值是否产生相同的字节。
	in := translateJSON(t.In)
	b, err := rlp.EncodeToBytes(in)
	if err != nil {
		return fmt.Errorf("encode failed: %v", err)
	}
	if !bytes.Equal(b, outb) {
		return fmt.Errorf("encode produced %x, want %x", b, outb)
	}
//测试流解码。
	s := rlp.NewStream(bytes.NewReader(outb), 0)
	return checkDecodeFromJSON(s, in)
}

func checkDecodeInterface(b []byte, isValid bool) error {
	err := rlp.DecodeBytes(b, new(interface{}))
	switch {
	case isValid && err != nil:
		return fmt.Errorf("decoding failed: %v", err)
	case !isValid && err == nil:
		return fmt.Errorf("decoding of invalid value succeeded")
	}
	return nil
}

//translateJSON使测试JSON值可以用rlp编码。
func translateJSON(v interface{}) interface{} {
	switch v := v.(type) {
	case float64:
		return uint64(v)
	case string:
if len(v) > 0 && v[0] == '#' { //开始一个假大内景。
			big, ok := new(big.Int).SetString(v[1:], 10)
			if !ok {
				panic(fmt.Errorf("bad test: bad big int: %q", v))
			}
			return big
		}
		return []byte(v)
	case []interface{}:
		new := make([]interface{}, len(v))
		for i := range v {
			new[i] = translateJSON(v[i])
		}
		return new
	default:
		panic(fmt.Errorf("can't handle %T", v))
	}
}

//checkdecodefromjson从exp.exp引导的s解码驱动
//通过调用基于解码操作（uint、big、list…）的流
//每个值的类型。从RLP流解码的值
//必须与JSON值匹配。
func checkDecodeFromJSON(s *rlp.Stream, exp interface{}) error {
	switch exp := exp.(type) {
	case uint64:
		i, err := s.Uint()
		if err != nil {
			return addStack("Uint", exp, err)
		}
		if i != exp {
			return addStack("Uint", exp, fmt.Errorf("result mismatch: got %d", i))
		}
	case *big.Int:
		big := new(big.Int)
		if err := s.Decode(&big); err != nil {
			return addStack("Big", exp, err)
		}
		if big.Cmp(exp) != 0 {
			return addStack("Big", exp, fmt.Errorf("result mismatch: got %d", big))
		}
	case []byte:
		b, err := s.Bytes()
		if err != nil {
			return addStack("Bytes", exp, err)
		}
		if !bytes.Equal(b, exp) {
			return addStack("Bytes", exp, fmt.Errorf("result mismatch: got %x", b))
		}
	case []interface{}:
		if _, err := s.List(); err != nil {
			return addStack("List", exp, err)
		}
		for i, v := range exp {
			if err := checkDecodeFromJSON(s, v); err != nil {
				return addStack(fmt.Sprintf("[%d]", i), exp, err)
			}
		}
		if err := s.ListEnd(); err != nil {
			return addStack("ListEnd", exp, err)
		}
	default:
		panic(fmt.Errorf("unhandled type: %T", exp))
	}
	return nil
}

func addStack(op string, val interface{}, err error) error {
	lines := strings.Split(err.Error(), "\n")
	lines = append(lines, fmt.Sprintf("\t%s: %v", op, val))
	return errors.New(strings.Join(lines, "\n"))
}

