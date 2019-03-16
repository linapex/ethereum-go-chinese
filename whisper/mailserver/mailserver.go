
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450123711123456>


//包mailserver提供了一个简单的示例mailserver实现
package mailserver

import (
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	whisper "github.com/ethereum/go-ethereum/whisper/whisperv6"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

//wmailserver表示邮件服务器的状态数据。
type WMailServer struct {
	db  *leveldb.DB
	w   *whisper.Whisper
	pow float64
	key []byte
}

type DBKey struct {
	timestamp uint32
	hash      common.Hash
	raw       []byte
}

//newdbkey是一个帮助函数，它创建一个leveldb
//来自哈希和整数的键。
func NewDbKey(t uint32, h common.Hash) *DBKey {
	const sz = common.HashLength + 4
	var k DBKey
	k.timestamp = t
	k.hash = h
	k.raw = make([]byte, sz)
	binary.BigEndian.PutUint32(k.raw, k.timestamp)
	copy(k.raw[4:], k.hash[:])
	return &k
}

//init初始化邮件服务器。
func (s *WMailServer) Init(shh *whisper.Whisper, path string, password string, pow float64) error {
	var err error
	if len(path) == 0 {
		return fmt.Errorf("DB file is not specified")
	}

	if len(password) == 0 {
		return fmt.Errorf("password is not specified")
	}

	s.db, err = leveldb.OpenFile(path, &opt.Options{OpenFilesCacheCapacity: 32})
	if err != nil {
		return fmt.Errorf("open DB file: %s", err)
	}

	s.w = shh
	s.pow = pow

	MailServerKeyID, err := s.w.AddSymKeyFromPassword(password)
	if err != nil {
		return fmt.Errorf("create symmetric key: %s", err)
	}
	s.key, err = s.w.GetSymKey(MailServerKeyID)
	if err != nil {
		return fmt.Errorf("save symmetric key: %s", err)
	}
	return nil
}

//关闭前清理。
func (s *WMailServer) Close() {
	if s.db != nil {
		s.db.Close()
	}
}

//存档存储
func (s *WMailServer) Archive(env *whisper.Envelope) {
	key := NewDbKey(env.Expiry-env.TTL, env.Hash())
	rawEnvelope, err := rlp.EncodeToBytes(env)
	if err != nil {
		log.Error(fmt.Sprintf("rlp.EncodeToBytes failed: %s", err))
	} else {
		err = s.db.Put(key.raw, rawEnvelope, nil)
		if err != nil {
			log.Error(fmt.Sprintf("Writing to DB failed: %s", err))
		}
	}
}

//Delivermail根据
//邮件的所有者。
func (s *WMailServer) DeliverMail(peer *whisper.Peer, request *whisper.Envelope) {
	if peer == nil {
		log.Error("Whisper peer is nil")
		return
	}

	ok, lower, upper, bloom := s.validateRequest(peer.ID(), request)
	if ok {
		s.processRequest(peer, lower, upper, bloom)
	}
}

func (s *WMailServer) processRequest(peer *whisper.Peer, lower, upper uint32, bloom []byte) []*whisper.Envelope {
	ret := make([]*whisper.Envelope, 0)
	var err error
	var zero common.Hash
	kl := NewDbKey(lower, zero)
ku := NewDbKey(upper+1, zero) //leveldb是独占的，而whisper API是包含的
	i := s.db.NewIterator(&util.Range{Start: kl.raw, Limit: ku.raw}, nil)
	defer i.Release()

	for i.Next() {
		var envelope whisper.Envelope
		err = rlp.DecodeBytes(i.Value(), &envelope)
		if err != nil {
			log.Error(fmt.Sprintf("RLP decoding failed: %s", err))
		}

		if whisper.BloomFilterMatch(bloom, envelope.Bloom()) {
			if peer == nil {
//用于测试目的
				ret = append(ret, &envelope)
			} else {
				err = s.w.SendP2PDirect(peer, &envelope)
				if err != nil {
					log.Error(fmt.Sprintf("Failed to send direct message to peer: %s", err))
					return nil
				}
			}
		}
	}

	err = i.Error()
	if err != nil {
		log.Error(fmt.Sprintf("Level DB iterator error: %s", err))
	}

	return ret
}

func (s *WMailServer) validateRequest(peerID []byte, request *whisper.Envelope) (bool, uint32, uint32, []byte) {
	if s.pow > 0.0 && request.PoW() < s.pow {
		return false, 0, 0, nil
	}

	f := whisper.Filter{KeySym: s.key}
	decrypted := request.Open(&f)
	if decrypted == nil {
		log.Warn(fmt.Sprintf("Failed to decrypt p2p request"))
		return false, 0, 0, nil
	}

	src := crypto.FromECDSAPub(decrypted.Src)
	if len(src)-len(peerID) == 1 {
		src = src[1:]
	}

//如果你想核对签名，你可以在这里核对。例如。：
//如果！bytes.equal（peerid，src）
	if src == nil {
		log.Warn(fmt.Sprintf("Wrong signature of p2p request"))
		return false, 0, 0, nil
	}

	var bloom []byte
	payloadSize := len(decrypted.Payload)
	if payloadSize < 8 {
		log.Warn(fmt.Sprintf("Undersized p2p request"))
		return false, 0, 0, nil
	} else if payloadSize == 8 {
		bloom = whisper.MakeFullNodeBloom()
	} else if payloadSize < 8+whisper.BloomFilterSize {
		log.Warn(fmt.Sprintf("Undersized bloom filter in p2p request"))
		return false, 0, 0, nil
	} else {
		bloom = decrypted.Payload[8 : 8+whisper.BloomFilterSize]
	}

	lower := binary.BigEndian.Uint32(decrypted.Payload[:4])
	upper := binary.BigEndian.Uint32(decrypted.Payload[4:8])
	return true, lower, upper, bloom
}

