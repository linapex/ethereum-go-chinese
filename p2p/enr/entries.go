
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450104786423808>


package enr

import (
	"fmt"
	"io"
	"net"

	"github.com/ethereum/go-ethereum/rlp"
)

//条目由已知的节点记录条目类型实现。
//
//要定义要包含在节点记录中的新条目，
//创建满足此接口的Go类型。类型应该
//如果需要对值进行额外检查，还可以实现rlp.decoder。
type Entry interface {
	ENRKey() string
}

type generic struct {
	key   string
	value interface{}
}

func (g generic) ENRKey() string { return g.key }

func (g generic) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, g.value)
}

func (g *generic) DecodeRLP(s *rlp.Stream) error {
	return s.Decode(g.value)
}

//WithEntry用键名包装任何值。它可用于设置和加载任意值
//在记录中。值v必须由rlp支持。要在加载时使用WithEntry，值
//必须是指针。
func WithEntry(k string, v interface{}) Entry {
	return &generic{key: k, value: v}
}

//tcp是“tcp”密钥，它保存节点的tcp端口。
type TCP uint16

func (v TCP) ENRKey() string { return "tcp" }

//udp是“udp”密钥，它保存节点的udp端口。
type UDP uint16

func (v UDP) ENRKey() string { return "udp" }

//ID是“ID”键，它保存标识方案的名称。
type ID string

const IDv4 = ID("v4") //默认标识方案

func (v ID) ENRKey() string { return "id" }

//IP是“IP”密钥，它保存节点的IP地址。
type IP net.IP

func (v IP) ENRKey() string { return "ip" }

//encoderlp实现rlp.encoder。
func (v IP) EncodeRLP(w io.Writer) error {
	if ip4 := net.IP(v).To4(); ip4 != nil {
		return rlp.Encode(w, ip4)
	}
	return rlp.Encode(w, net.IP(v))
}

//decoderlp实现rlp.decoder。
func (v *IP) DecodeRLP(s *rlp.Stream) error {
	if err := s.Decode((*net.IP)(v)); err != nil {
		return err
	}
	if len(*v) != 4 && len(*v) != 16 {
		return fmt.Errorf("invalid IP address, want 4 or 16 bytes: %v", *v)
	}
	return nil
}

//keyError是一个与键相关的错误。
type KeyError struct {
	Key string
	Err error
}

//错误实现错误。
func (err *KeyError) Error() string {
	if err.Err == errNotFound {
		return fmt.Sprintf("missing ENR key %q", err.Key)
	}
	return fmt.Sprintf("ENR key %q: %v", err.Key, err.Err)
}

//IsNotFound报告给定的错误是否意味着键/值对
//记录中缺少。
func IsNotFound(err error) bool {
	kerr, ok := err.(*KeyError)
	return ok && kerr.Err == errNotFound
}

