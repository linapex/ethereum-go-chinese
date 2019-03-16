
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450112629772288>


//包bmt是基于hashsize段的简单非当前引用实现
//任意但固定的最大chunksize上的二进制merkle树哈希
//
//此实现不利用任何并行列表和使用
//内存远比需要的多，但很容易看出它是正确的。
//它可以用于生成用于优化实现的测试用例。
//在bmt_test.go中对引用散列器的正确性进行了额外的检查。
//＊TestFisher
//*测试bmthshercorrection函数
package bmt

import (
	"hash"
)

//refhasher是BMT的非优化易读参考实现
type RefHasher struct {
maxDataLength int       //c*hashsize，其中c=2^ceil（log2（count）），其中count=ceil（length/hashsize）
sectionLength int       //2＊尺寸
hasher        hash.Hash //基哈希函数（keccak256 sha3）
}

//NewRefHasher返回新的RefHasher
func NewRefHasher(hasher BaseHasherFunc, count int) *RefHasher {
	h := hasher()
	hashsize := h.Size()
	c := 2
	for ; c < count; c *= 2 {
	}
	return &RefHasher{
		sectionLength: 2 * hashsize,
		maxDataLength: c * hashsize,
		hasher:        h,
	}
}

//hash返回字节片的bmt哈希
//实现swarmhash接口
func (rh *RefHasher) Hash(data []byte) []byte {
//如果数据小于基长度（maxdatalength），我们将提供零填充。
	d := make([]byte, rh.maxDataLength)
	length := len(data)
	if length > rh.maxDataLength {
		length = rh.maxDataLength
	}
	copy(d, data[:length])
	return rh.hash(d, rh.maxDataLength)
}

//数据的长度maxdatalength=segmentsize*2^k
//哈希在给定切片的两半递归调用自身
//连接结果，并返回该结果的哈希值
//如果d的长度是2*SegmentSize，则只返回该节的哈希值。
func (rh *RefHasher) hash(data []byte, length int) []byte {
	var section []byte
	if length == rh.sectionLength {
//部分包含两个数据段（D）
		section = data
	} else {
//部分包含左右BMT子目录的哈希
//通过在数据的左半部分和右半部分递归调用哈希来计算
		length /= 2
		section = append(rh.hash(data[:length], length), rh.hash(data[length:], length)...)
	}
	rh.hasher.Reset()
	rh.hasher.Write(section)
	return rh.hasher.Sum(nil)
}

