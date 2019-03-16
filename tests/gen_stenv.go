
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450122104705024>

//代码由github.com/fjl/gencodec生成。不要编辑。

package tests

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

var _ = (*stEnvMarshaling)(nil)

func (s stEnv) MarshalJSON() ([]byte, error) {
	type stEnv struct {
		Coinbase   common.UnprefixedAddress `json:"currentCoinbase"   gencodec:"required"`
		Difficulty *math.HexOrDecimal256    `json:"currentDifficulty" gencodec:"required"`
		GasLimit   math.HexOrDecimal64      `json:"currentGasLimit"   gencodec:"required"`
		Number     math.HexOrDecimal64      `json:"currentNumber"     gencodec:"required"`
		Timestamp  math.HexOrDecimal64      `json:"currentTimestamp"  gencodec:"required"`
	}
	var enc stEnv
	enc.Coinbase = common.UnprefixedAddress(s.Coinbase)
	enc.Difficulty = (*math.HexOrDecimal256)(s.Difficulty)
	enc.GasLimit = math.HexOrDecimal64(s.GasLimit)
	enc.Number = math.HexOrDecimal64(s.Number)
	enc.Timestamp = math.HexOrDecimal64(s.Timestamp)
	return json.Marshal(&enc)
}

func (s *stEnv) UnmarshalJSON(input []byte) error {
	type stEnv struct {
		Coinbase   *common.UnprefixedAddress `json:"currentCoinbase"   gencodec:"required"`
		Difficulty *math.HexOrDecimal256     `json:"currentDifficulty" gencodec:"required"`
		GasLimit   *math.HexOrDecimal64      `json:"currentGasLimit"   gencodec:"required"`
		Number     *math.HexOrDecimal64      `json:"currentNumber"     gencodec:"required"`
		Timestamp  *math.HexOrDecimal64      `json:"currentTimestamp"  gencodec:"required"`
	}
	var dec stEnv
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Coinbase == nil {
		return errors.New("missing required field 'currentCoinbase' for stEnv")
	}
	s.Coinbase = common.Address(*dec.Coinbase)
	if dec.Difficulty == nil {
		return errors.New("missing required field 'currentDifficulty' for stEnv")
	}
	s.Difficulty = (*big.Int)(dec.Difficulty)
	if dec.GasLimit == nil {
		return errors.New("missing required field 'currentGasLimit' for stEnv")
	}
	s.GasLimit = uint64(*dec.GasLimit)
	if dec.Number == nil {
		return errors.New("missing required field 'currentNumber' for stEnv")
	}
	s.Number = uint64(*dec.Number)
	if dec.Timestamp == nil {
		return errors.New("missing required field 'currentTimestamp' for stEnv")
	}
	s.Timestamp = uint64(*dec.Timestamp)
	return nil
}

