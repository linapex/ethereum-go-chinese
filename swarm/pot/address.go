
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450116119433216>


//包装罐见Go医生
package pot

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var (
	zerosBin = Address{}.Bin()
)

//地址是common.hash的别名
type Address common.Hash

//newAddressFromBytes从字节片构造地址
func NewAddressFromBytes(b []byte) Address {
	h := common.Hash{}
	copy(h[:], b)
	return Address(h)
}

func (a Address) String() string {
	return fmt.Sprintf("%x", a[:])
}

//marshaljson地址序列化
func (a *Address) MarshalJSON() (out []byte, err error) {
	return []byte(`"` + a.String() + `"`), nil
}

//取消标记JSON地址反序列化
func (a *Address) UnmarshalJSON(value []byte) error {
	*a = Address(common.HexToHash(string(value[1 : len(value)-1])))
	return nil
}

//bin返回地址的二进制表示的字符串形式（仅前8位）
func (a Address) Bin() string {
	return ToBin(a[:])
}

//Tobin将字节片转换为字符串二进制表示形式
func ToBin(a []byte) string {
	var bs []string
	for _, b := range a {
		bs = append(bs, fmt.Sprintf("%08b", b))
	}
	return strings.Join(bs, "")
}

//字节以字节片的形式返回地址
func (a Address) Bytes() []byte {
	return a[:]
}

//procmp比较距离a->target和b->target。
//如果a接近目标返回-1，如果b接近目标返回1
//如果相等，则为0。
func ProxCmp(a, x, y interface{}) int {
	return proxCmp(ToBytes(a), ToBytes(x), ToBytes(y))
}

func proxCmp(a, x, y []byte) int {
	for i := range a {
		dx := x[i] ^ a[i]
		dy := y[i] ^ a[i]
		if dx > dy {
			return 1
		} else if dx < dy {
			return -1
		}
	}
	return 0
}

//randomaddressat（地址，代理）生成随机地址
//在接近顺序，相对于地址的代理
//如果prox为负，则生成随机地址。
func RandomAddressAt(self Address, prox int) (addr Address) {
	addr = self
	pos := -1
	if prox >= 0 {
		pos = prox / 8
		trans := prox % 8
		transbytea := byte(0)
		for j := 0; j <= trans; j++ {
			transbytea |= 1 << uint8(7-j)
		}
		flipbyte := byte(1 << uint8(7-trans))
		transbyteb := transbytea ^ byte(255)
		randbyte := byte(rand.Intn(255))
		addr[pos] = ((addr[pos] & transbytea) ^ flipbyte) | randbyte&transbyteb
	}
	for i := pos + 1; i < len(addr); i++ {
		addr[i] = byte(rand.Intn(255))
	}

	return
}

//random address生成随机地址
func RandomAddress() Address {
	return RandomAddressAt(Address{}, -1)
}

//newAddressFromString从二进制表示的字符串创建字节片
func NewAddressFromString(s string) []byte {
	ha := [32]byte{}

	t := s + zerosBin[:len(zerosBin)-len(s)]
	for i := 0; i < 4; i++ {
		n, err := strconv.ParseUint(t[i*64:(i+1)*64], 2, 64)
		if err != nil {
			panic("wrong format: " + err.Error())
		}
		binary.BigEndian.PutUint64(ha[i*8:(i+1)*8], n)
	}
	return ha[:]
}

//BytesAddress是一个接口，用于按字节片寻址的元素
type BytesAddress interface {
	Address() []byte
}

//tobytes将val转换为字节
func ToBytes(v Val) []byte {
	if v == nil {
		return nil
	}
	b, ok := v.([]byte)
	if !ok {
		ba, ok := v.(BytesAddress)
		if !ok {
			panic(fmt.Sprintf("unsupported value type %T", v))
		}
		b = ba.Address()
	}
	return b
}

//defaultpof返回一个接近顺序比较运算符函数
func DefaultPof(max int) func(one, other Val, pos int) (int, bool) {
	return func(one, other Val, pos int) (int, bool) {
		po, eq := proximityOrder(ToBytes(one), ToBytes(other), pos)
		if po >= max {
			eq = true
			po = max
		}
		return po, eq
	}
}

//接近顺序返回两个参数：
//1。一个和另一个参数的相对接近顺序；
//2。布尔值，指示是否发生完全匹配（一个=另一个）。
func proximityOrder(one, other []byte, pos int) (int, bool) {
	for i := pos / 8; i < len(one); i++ {
		if one[i] == other[i] {
			continue
		}
		oxo := one[i] ^ other[i]
		start := 0
		if i == pos/8 {
			start = pos % 8
		}
		for j := start; j < 8; j++ {
			if (oxo>>uint8(7-j))&0x01 != 0 {
				return i*8 + j, false
			}
		}
	}
	return len(one) * 8, true
}

//标签以二进制格式显示节点的密钥
func Label(v Val) string {
	if v == nil {
		return "<nil>"
	}
	if s, ok := v.(fmt.Stringer); ok {
		return s.String()
	}
	if b, ok := v.([]byte); ok {
		return ToBin(b)
	}
	panic(fmt.Sprintf("unsupported value type %T", v))
}

