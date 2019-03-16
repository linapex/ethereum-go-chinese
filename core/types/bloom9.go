
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450081134743552>


package types

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type bytesBacked interface {
	Bytes() []byte
}

const (
//BloomByteLength表示头日志Bloom中使用的字节数。
	BloomByteLength = 256

//BloomBitLength表示在头日志Bloom中使用的位数。
	BloomBitLength = 8 * BloomByteLength
)

//Bloom表示2048位Bloom过滤器。
type Bloom [BloomByteLength]byte

//bytestobloom将字节片转换为bloom过滤器。
//如果B的尺寸不合适，就会恐慌。
func BytesToBloom(b []byte) Bloom {
	var bloom Bloom
	bloom.SetBytes(b)
	return bloom
}

//StimeByb将B的内容设置为给定字节。
//如果d的大小不合适，它会恐慌。
func (b *Bloom) SetBytes(d []byte) {
	if len(b) < len(d) {
		panic(fmt.Sprintf("bloom bytes too big %d %d", len(b), len(d)))
	}
	copy(b[BloomByteLength-len(d):], d)
}

//向过滤器中添加d。将来调用test（d）将返回true。
func (b *Bloom) Add(d *big.Int) {
	bin := new(big.Int).SetBytes(b[:])
	bin.Or(bin, bloom9(d.Bytes()))
	b.SetBytes(bin.Bytes())
}

//big将b转换为大整数。
func (b Bloom) Big() *big.Int {
	return new(big.Int).SetBytes(b[:])
}

func (b Bloom) Bytes() []byte {
	return b[:]
}

func (b Bloom) Test(test *big.Int) bool {
	return BloomLookup(b, test)
}

func (b Bloom) TestBytes(test []byte) bool {
	return b.Test(new(big.Int).SetBytes(test))

}

//MarshalText将B编码为带0x前缀的十六进制字符串。
func (b Bloom) MarshalText() ([]byte, error) {
	return hexutil.Bytes(b[:]).MarshalText()
}

//将文本B取消标记为带0x前缀的十六进制字符串。
func (b *Bloom) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Bloom", input, b[:])
}

func CreateBloom(receipts Receipts) Bloom {
	bin := new(big.Int)
	for _, receipt := range receipts {
		bin.Or(bin, LogsBloom(receipt.Logs))
	}

	return BytesToBloom(bin.Bytes())
}

func LogsBloom(logs []*Log) *big.Int {
	bin := new(big.Int)
	for _, log := range logs {
		bin.Or(bin, bloom9(log.Address.Bytes()))
		for _, b := range log.Topics {
			bin.Or(bin, bloom9(b[:]))
		}
	}

	return bin
}

func bloom9(b []byte) *big.Int {
	b = crypto.Keccak256(b)

	r := new(big.Int)

	for i := 0; i < 6; i += 2 {
		t := big.NewInt(1)
		b := (uint(b[i+1]) + (uint(b[i]) << 8)) & 2047
		r.Or(r, t.Lsh(t, b))
	}

	return r
}

var Bloom9 = bloom9

func BloomLookup(bin Bloom, topic bytesBacked) bool {
	bloom := bin.Big()
	cmp := bloom9(topic.Bytes())

	return bloom.And(bloom, cmp).Cmp(cmp) == 0
}

