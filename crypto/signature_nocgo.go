
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450085589094400>


//+构建Nacl JS Nocgo

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
)

//ecrecover返回创建给定签名的未压缩公钥。
func Ecrecover(hash, sig []byte) ([]byte, error) {
	pub, err := SigToPub(hash, sig)
	if err != nil {
		return nil, err
	}
	bytes := (*btcec.PublicKey)(pub).SerializeUncompressed()
	return bytes, err
}

//sigtopub返回创建给定签名的公钥。
func SigToPub(hash, sig []byte) (*ecdsa.PublicKey, error) {
//转换为btcec输入格式，开始时使用“recovery id”v。
	btcsig := make([]byte, 65)
	btcsig[0] = sig[64] + 27
	copy(btcsig[1:], sig)

	pub, _, err := btcec.RecoverCompact(btcec.S256(), btcsig, hash)
	return (*ecdsa.PublicKey)(pub), err
}

//sign计算ECDSA签名。
//
//此函数容易受到选中的可能泄漏的明文攻击
//有关用于签名的私钥的信息。来电者必须
//请注意，给定的哈希不能由对手选择。共同的
//解决方案是在计算签名之前散列任何输入。
//
//生成的签名采用[R V]格式，其中V为0或1。
func Sign(hash []byte, prv *ecdsa.PrivateKey) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash is required to be exactly 32 bytes (%d)", len(hash))
	}
	if prv.Curve != btcec.S256() {
		return nil, fmt.Errorf("private key curve is not secp256k1")
	}
	sig, err := btcec.SignCompact(btcec.S256(), (*btcec.PrivateKey)(prv), hash, false)
	if err != nil {
		return nil, err
	}
//转换为末尾带有'recovery id'v的以太坊签名格式。
	v := sig[0] - 27
	copy(sig, sig[1:])
	sig[64] = v
	return sig, nil
}

//VerifySignature检查给定的公钥是否通过哈希创建了签名。
//公钥应为压缩（33字节）或未压缩（65字节）格式。
//签名应采用64字节[r_s]格式。
func VerifySignature(pubkey, hash, signature []byte) bool {
	if len(signature) != 64 {
		return false
	}
	sig := &btcec.Signature{R: new(big.Int).SetBytes(signature[:32]), S: new(big.Int).SetBytes(signature[32:])}
	key, err := btcec.ParsePubKey(pubkey, btcec.S256())
	if err != nil {
		return false
	}
//拒绝可延展签名。libsecp256k1执行此检查，但btcec不执行。
	if sig.S.Cmp(secp256k1halfN) > 0 {
		return false
	}
	return sig.Verify(hash, key)
}

//解压缩PubKey以33字节的压缩格式解析公钥。
func DecompressPubkey(pubkey []byte) (*ecdsa.PublicKey, error) {
	if len(pubkey) != 33 {
		return nil, errors.New("invalid compressed public key length")
	}
	key, err := btcec.ParsePubKey(pubkey, btcec.S256())
	if err != nil {
		return nil, err
	}
	return key.ToECDSA(), nil
}

//compresspubkey将公钥编码为33字节的压缩格式。
func CompressPubkey(pubkey *ecdsa.PublicKey) []byte {
	return (*btcec.PublicKey)(pubkey).SerializeCompressed()
}

//s256返回secp256k1曲线的一个实例。
func S256() elliptic.Curve {
	return btcec.S256()
}

