
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450111543447552>


package api

import (
	"encoding/binary"
	"errors"

	"github.com/ethereum/go-ethereum/swarm/storage/encryption"
	"golang.org/x/crypto/sha3"
)

type RefEncryption struct {
	refSize int
	span    []byte
}

func NewRefEncryption(refSize int) *RefEncryption {
	span := make([]byte, 8)
	binary.LittleEndian.PutUint64(span, uint64(refSize))
	return &RefEncryption{
		refSize: refSize,
		span:    span,
	}
}

func (re *RefEncryption) Encrypt(ref []byte, key []byte) ([]byte, error) {
	spanEncryption := encryption.New(key, 0, uint32(re.refSize/32), sha3.NewLegacyKeccak256)
	encryptedSpan, err := spanEncryption.Encrypt(re.span)
	if err != nil {
		return nil, err
	}
	dataEncryption := encryption.New(key, re.refSize, 0, sha3.NewLegacyKeccak256)
	encryptedData, err := dataEncryption.Encrypt(ref)
	if err != nil {
		return nil, err
	}
	encryptedRef := make([]byte, len(ref)+8)
	copy(encryptedRef[:8], encryptedSpan)
	copy(encryptedRef[8:], encryptedData)

	return encryptedRef, nil
}

func (re *RefEncryption) Decrypt(ref []byte, key []byte) ([]byte, error) {
	spanEncryption := encryption.New(key, 0, uint32(re.refSize/32), sha3.NewLegacyKeccak256)
	decryptedSpan, err := spanEncryption.Decrypt(ref[:8])
	if err != nil {
		return nil, err
	}

	size := binary.LittleEndian.Uint64(decryptedSpan)
	if size != uint64(len(ref)-8) {
		return nil, errors.New("invalid span in encrypted reference")
	}

	dataEncryption := encryption.New(key, re.refSize, 0, sha3.NewLegacyKeccak256)
	decryptedRef, err := dataEncryption.Decrypt(ref[8:])
	if err != nil {
		return nil, err
	}

	return decryptedRef, nil
}

