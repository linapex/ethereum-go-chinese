
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:46</date>
//</624450125434982400>

//代码由github.com/fjl/gencodec生成。不要编辑。

package whisperv6

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

var _ = (*criteriaOverride)(nil)

//marshaljson将类型条件封送到json字符串
func (c Criteria) MarshalJSON() ([]byte, error) {
	type Criteria struct {
		SymKeyID     string        `json:"symKeyID"`
		PrivateKeyID string        `json:"privateKeyID"`
		Sig          hexutil.Bytes `json:"sig"`
		MinPow       float64       `json:"minPow"`
		Topics       []TopicType   `json:"topics"`
		AllowP2P     bool          `json:"allowP2P"`
	}
	var enc Criteria
	enc.SymKeyID = c.SymKeyID
	enc.PrivateKeyID = c.PrivateKeyID
	enc.Sig = c.Sig
	enc.MinPow = c.MinPow
	enc.Topics = c.Topics
	enc.AllowP2P = c.AllowP2P
	return json.Marshal(&enc)
}

//将JSON的类型条件取消标记为JSON字符串
func (c *Criteria) UnmarshalJSON(input []byte) error {
	type Criteria struct {
		SymKeyID     *string        `json:"symKeyID"`
		PrivateKeyID *string        `json:"privateKeyID"`
		Sig          *hexutil.Bytes `json:"sig"`
		MinPow       *float64       `json:"minPow"`
		Topics       []TopicType    `json:"topics"`
		AllowP2P     *bool          `json:"allowP2P"`
	}
	var dec Criteria
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.SymKeyID != nil {
		c.SymKeyID = *dec.SymKeyID
	}
	if dec.PrivateKeyID != nil {
		c.PrivateKeyID = *dec.PrivateKeyID
	}
	if dec.Sig != nil {
		c.Sig = *dec.Sig
	}
	if dec.MinPow != nil {
		c.MinPow = *dec.MinPow
	}
	if dec.Topics != nil {
		c.Topics = dec.Topics
	}
	if dec.AllowP2P != nil {
		c.AllowP2P = *dec.AllowP2P
	}
	return nil
}

