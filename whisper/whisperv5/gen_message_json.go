
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:46</date>
//</624450124562567168>

//代码由github.com/fjl/gencodec生成。不要编辑。

package whisperv5

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var _ = (*messageOverride)(nil)

func (m Message) MarshalJSON() ([]byte, error) {
	type Message struct {
		Sig       hexutil.Bytes `json:"sig,omitempty"`
		TTL       uint32        `json:"ttl"`
		Timestamp uint32        `json:"timestamp"`
		Topic     TopicType     `json:"topic"`
		Payload   hexutil.Bytes `json:"payload"`
		Padding   hexutil.Bytes `json:"padding"`
		PoW       float64       `json:"pow"`
		Hash      hexutil.Bytes `json:"hash"`
		Dst       hexutil.Bytes `json:"recipientPublicKey,omitempty"`
	}
	var enc Message
	enc.Sig = m.Sig
	enc.TTL = m.TTL
	enc.Timestamp = m.Timestamp
	enc.Topic = m.Topic
	enc.Payload = m.Payload
	enc.Padding = m.Padding
	enc.PoW = m.PoW
	enc.Hash = m.Hash
	enc.Dst = m.Dst
	return json.Marshal(&enc)
}

func (m *Message) UnmarshalJSON(input []byte) error {
	type Message struct {
		Sig       *hexutil.Bytes `json:"sig,omitempty"`
		TTL       *uint32        `json:"ttl"`
		Timestamp *uint32        `json:"timestamp"`
		Topic     *TopicType     `json:"topic"`
		Payload   *hexutil.Bytes `json:"payload"`
		Padding   *hexutil.Bytes `json:"padding"`
		PoW       *float64       `json:"pow"`
		Hash      *hexutil.Bytes `json:"hash"`
		Dst       *hexutil.Bytes `json:"recipientPublicKey,omitempty"`
	}
	var dec Message
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Sig != nil {
		m.Sig = *dec.Sig
	}
	if dec.TTL != nil {
		m.TTL = *dec.TTL
	}
	if dec.Timestamp != nil {
		m.Timestamp = *dec.Timestamp
	}
	if dec.Topic != nil {
		m.Topic = *dec.Topic
	}
	if dec.Payload != nil {
		m.Payload = *dec.Payload
	}
	if dec.Padding != nil {
		m.Padding = *dec.Padding
	}
	if dec.PoW != nil {
		m.PoW = *dec.PoW
	}
	if dec.Hash != nil {
		m.Hash = *dec.Hash
	}
	if dec.Dst != nil {
		m.Dst = *dec.Dst
	}
	return nil
}

