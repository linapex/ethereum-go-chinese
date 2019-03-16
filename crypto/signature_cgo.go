
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450085547151360>


//+建设！NaCl，！JS！NOCGO

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

//ecrecover返回创建给定签名的未压缩公钥。
func Ecrecover(hash, sig []byte) ([]byte, error) {
	return secp256k1.RecoverPubkey(hash, sig)
}

//sigtopub返回创建给定签名的公钥。
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
	s, err := Ecrecover(hash, sig)
	if err != nil {
		return nil, err
	}

	x, y := elliptic.Unmarshal(S256(), s)
	return &ecdsa.PublicKey{Curve: S256(), X: x, Y: y}, nil
}

//sign计算ECDSA签名。
//
//此函数容易受到选中的可能泄漏的明文攻击
//有关用于签名的私钥的信息。来电者必须
//请注意，给定的哈希不能由对手选择。共同的
//解决方案是在计算签名之前散列任何输入。
//
//生成的签名采用[R V]格式，其中V为0或1。
func Sign(hash []byte, prv *ecdsa.PrivateKey) (sig []byte, err error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	seckey := math.PaddedBigBytes(prv.D, prv.Params().BitSize/8)
	defer zeroBytes(seckey)
	return secp256k1.Sign(hash, seckey)
}

//VerifySignature checks that the given public key created signature over hash.
//公钥应为压缩（33字节）或未压缩（65字节）格式。
//签名应采用64字节[r_s]格式。
func VerifySignature(pubkey, hash, signature []byte) bool {
	return secp256k1.VerifySignature(pubkey, hash, signature)
}

//解压缩PubKey以33字节的压缩格式解析公钥。
func DecompressPubkey(pubkey []byte) (*ecdsa.PublicKey, error) {
	x, y := secp256k1.DecompressPubkey(pubkey)
	if x == nil {
		return nil, fmt.Errorf("invalid public key")
	}
	return &ecdsa.PublicKey{X: x, Y: y, Curve: S256()}, nil
}

//compresspubkey将公钥编码为33字节的压缩格式。
func CompressPubkey(pubkey *ecdsa.PublicKey) []byte {
	return secp256k1.CompressPubkey(pubkey.X, pubkey.Y)
}

//s256返回secp256k1曲线的一个实例。
func S256() elliptic.Curve {
	return secp256k1.S256()
}

