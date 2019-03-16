
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450104215998464>


package enode

import (
	"crypto/ecdsa"
	"fmt"
	"io"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enr"
	"github.com/ethereum/go-ethereum/rlp"
	"golang.org/x/crypto/sha3"
)

//已知安全身份方案的列表。
var ValidSchemes = enr.SchemeMap{
	"v4": V4ID{},
}

var ValidSchemesForTesting = enr.SchemeMap{
	"v4":   V4ID{},
	"null": NullID{},
}

//v4id是“v4”标识方案。
type V4ID struct{}

//signv4使用v4方案对记录进行签名。
func SignV4(r *enr.Record, privkey *ecdsa.PrivateKey) error {
//复制r以避免在签名失败时修改它。
	cpy := *r
	cpy.Set(enr.ID("v4"))
	cpy.Set(Secp256k1(privkey.PublicKey))

	h := sha3.NewLegacyKeccak256()
	rlp.Encode(h, cpy.AppendElements(nil))
	sig, err := crypto.Sign(h.Sum(nil), privkey)
	if err != nil {
		return err
	}
sig = sig[:len(sig)-1] //移除V
	if err = cpy.SetSig(V4ID{}, sig); err == nil {
		*r = cpy
	}
	return err
}

func (V4ID) Verify(r *enr.Record, sig []byte) error {
	var entry s256raw
	if err := r.Load(&entry); err != nil {
		return err
	} else if len(entry) != 33 {
		return fmt.Errorf("invalid public key")
	}

	h := sha3.NewLegacyKeccak256()
	rlp.Encode(h, r.AppendElements(nil))
	if !crypto.VerifySignature(entry, h.Sum(nil), sig) {
		return enr.ErrInvalidSig
	}
	return nil
}

func (V4ID) NodeAddr(r *enr.Record) []byte {
	var pubkey Secp256k1
	err := r.Load(&pubkey)
	if err != nil {
		return nil
	}
	buf := make([]byte, 64)
	math.ReadBits(pubkey.X, buf[:32])
	math.ReadBits(pubkey.Y, buf[32:])
	return crypto.Keccak256(buf)
}

//secp256k1是保存公钥的“secp256k1”密钥。
type Secp256k1 ecdsa.PublicKey

func (v Secp256k1) ENRKey() string { return "secp256k1" }

//encoderlp实现rlp.encoder。
func (v Secp256k1) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, crypto.CompressPubkey((*ecdsa.PublicKey)(&v)))
}

//decoderlp实现rlp.decoder。
func (v *Secp256k1) DecodeRLP(s *rlp.Stream) error {
	buf, err := s.Bytes()
	if err != nil {
		return err
	}
	pk, err := crypto.DecompressPubkey(buf)
	if err != nil {
		return err
	}
	*v = (Secp256k1)(*pk)
	return nil
}

//s256raw是一个未分析的secp256k1公钥条目。
type s256raw []byte

func (s256raw) ENRKey() string { return "secp256k1" }

//v4compatid是“v4”方案的一个较弱且不安全的版本，它只检查
//存在secp256k1公钥，但不验证签名。
type v4CompatID struct {
	V4ID
}

func (v4CompatID) Verify(r *enr.Record, sig []byte) error {
	var pubkey Secp256k1
	return r.Load(&pubkey)
}

func signV4Compat(r *enr.Record, pubkey *ecdsa.PublicKey) {
	r.Set((*Secp256k1)(pubkey))
	if err := r.SetSig(v4CompatID{}, []byte{}); err != nil {
		panic(err)
	}
}

//“空”是“空”的ENR身份方案。此方案存储节点
//记录中没有任何签名的ID。
type NullID struct{}

func (NullID) Verify(r *enr.Record, sig []byte) error {
	return nil
}

func (NullID) NodeAddr(r *enr.Record) []byte {
	var id ID
	r.Load(enr.WithEntry("nulladdr", &id))
	return id[:]
}

func SignNull(r *enr.Record, id ID) *Node {
	r.Set(enr.ID("null"))
	r.Set(enr.WithEntry("nulladdr", id))
	if err := r.SetSig(NullID{}, []byte{}); err != nil {
		panic(err)
	}
	return &Node{r: *r, id: id}
}

