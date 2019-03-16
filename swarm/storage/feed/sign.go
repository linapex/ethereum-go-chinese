
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450119428739072>


package feed

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const signatureLength = 65

//签名是具有签名大小的静态字节数组的别名
type Signature [signatureLength]byte

//签名者签名源更新有效负载
type Signer interface {
	Sign(common.Hash) (Signature, error)
	Address() common.Address
}

//GenericSigner实现了签名者接口
//在大多数情况下，可能应该使用普通签名者。
type GenericSigner struct {
	PrivKey *ecdsa.PrivateKey
	address common.Address
}

//NewGenericSigner生成一个签名者，该签名者将使用提供的私钥对所有内容进行签名。
func NewGenericSigner(privKey *ecdsa.PrivateKey) *GenericSigner {
	return &GenericSigner{
		PrivKey: privKey,
		address: crypto.PubkeyToAddress(privKey.PublicKey),
	}
}

//在提供的数据上签名
//它包装了ethereum crypto.sign（）方法
func (s *GenericSigner) Sign(data common.Hash) (signature Signature, err error) {
	signaturebytes, err := crypto.Sign(data.Bytes(), s.PrivKey)
	if err != nil {
		return
	}
	copy(signature[:], signaturebytes)
	return
}

//地址返回签名者私钥的公钥
func (s *GenericSigner) Address() common.Address {
	return s.address
}

//getuseraddr提取源更新签名者的地址
func getUserAddr(digest common.Hash, signature Signature) (common.Address, error) {
	pub, err := crypto.SigToPub(digest.Bytes(), signature[:])
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pub), nil
}

