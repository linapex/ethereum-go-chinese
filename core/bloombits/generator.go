
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450078098067456>


package bloombits

import (
	"errors"

	"github.com/ethereum/go-ethereum/core/types"
)

var (
//如果用户尝试添加更多Bloom筛选器，则返回errSectionOutofBounds
//批处理的可用空间不足，或者如果尝试检索超过容量。
	errSectionOutOfBounds = errors.New("section out of bounds")

//如果用户尝试检索指定的
//比容量大一点。
	errBloomBitOutOfBounds = errors.New("bloom bit out of bounds")
)

//发电机接收许多布卢姆滤波器并生成旋转的布卢姆位
//用于批量过滤。
type Generator struct {
blooms   [types.BloomBitLength][]byte //每比特匹配的旋转花束
sections uint                         //要一起批处理的节数
nextSec  uint                         //添加花束时要设置的下一节
}

//NewGenerator创建一个旋转的Bloom Generator，它可以迭代填充
//批量布卢姆过滤器的钻头。
func NewGenerator(sections uint) (*Generator, error) {
	if sections%8 != 0 {
		return nil, errors.New("section count not multiple of 8")
	}
	b := &Generator{sections: sections}
	for i := 0; i < types.BloomBitLength; i++ {
		b.blooms[i] = make([]byte, sections/8)
	}
	return b, nil
}

//addbloom接受一个bloom过滤器并设置相应的位列
//在记忆中。
func (b *Generator) AddBloom(index uint, bloom types.Bloom) error {
//确保我们添加的布卢姆过滤器不会超过我们的容量
	if b.nextSec >= b.sections {
		return errSectionOutOfBounds
	}
	if b.nextSec != index {
		return errors.New("bloom filter with unexpected index")
	}
//旋转花束并插入我们的收藏
	byteIndex := b.nextSec / 8
	bitMask := byte(1) << byte(7-b.nextSec%8)

	for i := 0; i < types.BloomBitLength; i++ {
		bloomByteIndex := types.BloomByteLength - 1 - i/8
		bloomBitMask := byte(1) << byte(i%8)

		if (bloom[bloomByteIndex] & bloomBitMask) != 0 {
			b.blooms[i][byteIndex] |= bitMask
		}
	}
	b.nextSec++

	return nil
}

//位集返回属于给定位索引的位向量。
//花开了。
func (b *Generator) Bitset(idx uint) ([]byte, error) {
	if b.nextSec != b.sections {
		return nil, errors.New("bloom not fully generated yet")
	}
	if idx >= types.BloomBitLength {
		return nil, errBloomBitOutOfBounds
	}
	return b.blooms[idx], nil
}

