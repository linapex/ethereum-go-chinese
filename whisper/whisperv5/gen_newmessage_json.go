
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:46</date>
//</624450124604510208>

//代码由github.com/fjl/gencodec生成。不要编辑。

package whisperv5

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var _ = (*newMessageOverride)(nil)

func (n NewMessage) MarshalJSON() ([]byte, error) {
	type NewMessage struct {
		SymKeyID   string        `json:"symKeyID"`
		PublicKey  hexutil.Bytes `json:"pubKey"`
		Sig        string        `json:"sig"`
		TTL        uint32        `json:"ttl"`
		Topic      TopicType     `json:"topic"`
		Payload    hexutil.Bytes `json:"payload"`
		Padding    hexutil.Bytes `json:"padding"`
		PowTime    uint32        `json:"powTime"`
		PowTarget  float64       `json:"powTarget"`
		TargetPeer string        `json:"targetPeer"`
	}
	var enc NewMessage
	enc.SymKeyID = n.SymKeyID
	enc.PublicKey = n.PublicKey
	enc.Sig = n.Sig
	enc.TTL = n.TTL
	enc.Topic = n.Topic
	enc.Payload = n.Payload
	enc.Padding = n.Padding
	enc.PowTime = n.PowTime
	enc.PowTarget = n.PowTarget
	enc.TargetPeer = n.TargetPeer
	return json.Marshal(&enc)
}

func (n *NewMessage) UnmarshalJSON(input []byte) error {
	type NewMessage struct {
		SymKeyID   *string        `json:"symKeyID"`
		PublicKey  *hexutil.Bytes `json:"pubKey"`
		Sig        *string        `json:"sig"`
		TTL        *uint32        `json:"ttl"`
		Topic      *TopicType     `json:"topic"`
		Payload    *hexutil.Bytes `json:"payload"`
		Padding    *hexutil.Bytes `json:"padding"`
		PowTime    *uint32        `json:"powTime"`
		PowTarget  *float64       `json:"powTarget"`
		TargetPeer *string        `json:"targetPeer"`
	}
	var dec NewMessage
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.SymKeyID != nil {
		n.SymKeyID = *dec.SymKeyID
	}
	if dec.PublicKey != nil {
		n.PublicKey = *dec.PublicKey
	}
	if dec.Sig != nil {
		n.Sig = *dec.Sig
	}
	if dec.TTL != nil {
		n.TTL = *dec.TTL
	}
	if dec.Topic != nil {
		n.Topic = *dec.Topic
	}
	if dec.Payload != nil {
		n.Payload = *dec.Payload
	}
	if dec.Padding != nil {
		n.Padding = *dec.Padding
	}
	if dec.PowTime != nil {
		n.PowTime = *dec.PowTime
	}
	if dec.PowTarget != nil {
		n.PowTarget = *dec.PowTarget
	}
	if dec.TargetPeer != nil {
		n.TargetPeer = *dec.TargetPeer
	}
	return nil
}

