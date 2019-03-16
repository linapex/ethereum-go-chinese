
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:35</date>
//</624450081499648000>


package types

import (
	"bytes"
	"fmt"
	"io"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
)

//go:generate gencodec -type Receipt -field-override receiptMarshaling -out gen_receipt_json.go

var (
	receiptStatusFailedRLP     = []byte{}
	receiptStatusSuccessfulRLP = []byte{0x01}
)

const (
//ReceiptStatusFailed是执行失败时事务的状态代码。
	ReceiptStatusFailed = uint64(0)

//ReceiptStatusSuccessful是执行成功时事务的状态代码。
	ReceiptStatusSuccessful = uint64(1)
)

//收据表示交易的结果。
type Receipt struct {
//共识领域
	PostState         []byte `json:"root"`
	Status            uint64 `json:"status"`
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed" gencodec:"required"`
	Bloom             Bloom  `json:"logsBloom"         gencodec:"required"`
	Logs              []*Log `json:"logs"              gencodec:"required"`

//实现字段（不要重新排序！）
	TxHash          common.Hash    `json:"transactionHash" gencodec:"required"`
	ContractAddress common.Address `json:"contractAddress"`
	GasUsed         uint64         `json:"gasUsed" gencodec:"required"`
}

type receiptMarshaling struct {
	PostState         hexutil.Bytes
	Status            hexutil.Uint64
	CumulativeGasUsed hexutil.Uint64
	GasUsed           hexutil.Uint64
}

//receiptrlp是收据的共识编码。
type receiptRLP struct {
	PostStateOrStatus []byte
	CumulativeGasUsed uint64
	Bloom             Bloom
	Logs              []*Log
}

type receiptStorageRLP struct {
	PostStateOrStatus []byte
	CumulativeGasUsed uint64
	Bloom             Bloom
	TxHash            common.Hash
	ContractAddress   common.Address
	Logs              []*LogForStorage
	GasUsed           uint64
}

//NewReceipt创建一个准交易凭证，复制init字段。
func NewReceipt(root []byte, failed bool, cumulativeGasUsed uint64) *Receipt {
	r := &Receipt{PostState: common.CopyBytes(root), CumulativeGasUsed: cumulativeGasUsed}
	if failed {
		r.Status = ReceiptStatusFailed
	} else {
		r.Status = ReceiptStatusSuccessful
	}
	return r
}

//encoderlp实现rlp.encoder，并扁平收据的共识字段
//进入RLP流。如果不存在后置状态，则假定为拜占庭分叉。
func (r *Receipt) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &receiptRLP{r.statusEncoding(), r.CumulativeGasUsed, r.Bloom, r.Logs})
}

//decoderlp实现rlp.decoder，并加载收据的共识字段
//来自RLP流。
func (r *Receipt) DecodeRLP(s *rlp.Stream) error {
	var dec receiptRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	if err := r.setStatus(dec.PostStateOrStatus); err != nil {
		return err
	}
	r.CumulativeGasUsed, r.Bloom, r.Logs = dec.CumulativeGasUsed, dec.Bloom, dec.Logs
	return nil
}

func (r *Receipt) setStatus(postStateOrStatus []byte) error {
	switch {
	case bytes.Equal(postStateOrStatus, receiptStatusSuccessfulRLP):
		r.Status = ReceiptStatusSuccessful
	case bytes.Equal(postStateOrStatus, receiptStatusFailedRLP):
		r.Status = ReceiptStatusFailed
	case len(postStateOrStatus) == len(common.Hash{}):
		r.PostState = postStateOrStatus
	default:
		return fmt.Errorf("invalid receipt status %x", postStateOrStatus)
	}
	return nil
}

func (r *Receipt) statusEncoding() []byte {
	if len(r.PostState) == 0 {
		if r.Status == ReceiptStatusFailed {
			return receiptStatusFailedRLP
		}
		return receiptStatusSuccessfulRLP
	}
	return r.PostState
}

//SIZE返回所有内部内容使用的大致内存。它被使用
//估计和限制各种缓存的内存消耗。
func (r *Receipt) Size() common.StorageSize {
	size := common.StorageSize(unsafe.Sizeof(*r)) + common.StorageSize(len(r.PostState))

	size += common.StorageSize(len(r.Logs)) * common.StorageSize(unsafe.Sizeof(Log{}))
	for _, log := range r.Logs {
		size += common.StorageSize(len(log.Topics)*common.HashLength + len(log.Data))
	}
	return size
}

//ReceiptForStorage是一个围绕收据的包装，它扁平化并解析
//收据的全部内容，而不是最初的共识字段。
type ReceiptForStorage Receipt

//encoderlp实现rlp.encoder，并平展收据的所有内容字段
//进入RLP流。
func (r *ReceiptForStorage) EncodeRLP(w io.Writer) error {
	enc := &receiptStorageRLP{
		PostStateOrStatus: (*Receipt)(r).statusEncoding(),
		CumulativeGasUsed: r.CumulativeGasUsed,
		Bloom:             r.Bloom,
		TxHash:            r.TxHash,
		ContractAddress:   r.ContractAddress,
		Logs:              make([]*LogForStorage, len(r.Logs)),
		GasUsed:           r.GasUsed,
	}
	for i, log := range r.Logs {
		enc.Logs[i] = (*LogForStorage)(log)
	}
	return rlp.Encode(w, enc)
}

//decoderlp实现rlp.decoder，并加载共识和实现
//来自RLP流的收据字段。
func (r *ReceiptForStorage) DecodeRLP(s *rlp.Stream) error {
	var dec receiptStorageRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	if err := (*Receipt)(r).setStatus(dec.PostStateOrStatus); err != nil {
		return err
	}
//分配共识字段
	r.CumulativeGasUsed, r.Bloom = dec.CumulativeGasUsed, dec.Bloom
	r.Logs = make([]*Log, len(dec.Logs))
	for i, log := range dec.Logs {
		r.Logs[i] = (*Log)(log)
	}
//分配实现字段
	r.TxHash, r.ContractAddress, r.GasUsed = dec.TxHash, dec.ContractAddress, dec.GasUsed
	return nil
}

//Receipts是一个围绕收据数组的包装器，用于实现可派生列表。
type Receipts []*Receipt

//Len返回此列表中的收货数量。
func (r Receipts) Len() int { return len(r) }

//getrlp返回列表中一个收据的rlp编码。
func (r Receipts) GetRlp(i int) []byte {
	bytes, err := rlp.EncodeToBytes(r[i])
	if err != nil {
		panic(err)
	}
	return bytes
}

