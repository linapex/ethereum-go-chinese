
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450095793836032>


//包les实现轻以太坊子协议。
package les

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rlp"
)

//用于匹配协议版本和消息的常量
const (
	lpv1 = 1
	lpv2 = 2
)

//支持的LES协议版本（第一个是主协议）
var (
	ClientProtocolVersions    = []uint{lpv2, lpv1}
	ServerProtocolVersions    = []uint{lpv2, lpv1}
AdvertiseProtocolVersions = []uint{lpv2} //客户端正在搜索列表中的第一个公告协议
)

//对应于不同协议版本的已实现消息数。
var ProtocolLengths = map[uint]uint64{lpv1: 15, lpv2: 22}

const (
	NetworkId          = 1
ProtocolMaxMsgSize = 10 * 1024 * 1024 //协议消息大小的最大上限
)

//LES协议消息代码
const (
//属于lpv1的协议消息
	StatusMsg          = 0x00
	AnnounceMsg        = 0x01
	GetBlockHeadersMsg = 0x02
	BlockHeadersMsg    = 0x03
	GetBlockBodiesMsg  = 0x04
	BlockBodiesMsg     = 0x05
	GetReceiptsMsg     = 0x06
	ReceiptsMsg        = 0x07
	GetProofsV1Msg     = 0x08
	ProofsV1Msg        = 0x09
	GetCodeMsg         = 0x0a
	CodeMsg            = 0x0b
	SendTxMsg          = 0x0c
	GetHeaderProofsMsg = 0x0d
	HeaderProofsMsg    = 0x0e
//属于lpv2的协议消息
	GetProofsV2Msg         = 0x0f
	ProofsV2Msg            = 0x10
	GetHelperTrieProofsMsg = 0x11
	HelperTrieProofsMsg    = 0x12
	SendTxV2Msg            = 0x13
	GetTxStatusMsg         = 0x14
	TxStatusMsg            = 0x15
)

type errCode int

const (
	ErrMsgTooLarge = iota
	ErrDecode
	ErrInvalidMsgCode
	ErrProtocolVersionMismatch
	ErrNetworkIdMismatch
	ErrGenesisBlockMismatch
	ErrNoStatusMsg
	ErrExtraStatusMsg
	ErrSuspendedPeer
	ErrUselessPeer
	ErrRequestRejected
	ErrUnexpectedResponse
	ErrInvalidResponse
	ErrTooManyTimeouts
	ErrMissingKey
)

func (e errCode) String() string {
	return errorToString[int(e)]
}

//一旦旧代码用完，XXX就会更改
var errorToString = map[int]string{
	ErrMsgTooLarge:             "Message too long",
	ErrDecode:                  "Invalid message",
	ErrInvalidMsgCode:          "Invalid message code",
	ErrProtocolVersionMismatch: "Protocol version mismatch",
	ErrNetworkIdMismatch:       "NetworkId mismatch",
	ErrGenesisBlockMismatch:    "Genesis block mismatch",
	ErrNoStatusMsg:             "No status message",
	ErrExtraStatusMsg:          "Extra status message",
	ErrSuspendedPeer:           "Suspended peer",
	ErrRequestRejected:         "Request rejected",
	ErrUnexpectedResponse:      "Unexpected response",
	ErrInvalidResponse:         "Invalid response",
	ErrTooManyTimeouts:         "Too many request timeouts",
	ErrMissingKey:              "Key missing from list",
}

type announceBlock struct {
Hash   common.Hash //正在公布的一个特定块的哈希
Number uint64      //公布的一个特定区块的编号
Td     *big.Int    //宣布一个特定区块的总难度
}

//公告数据是块公告的网络包。
type announceData struct {
Hash       common.Hash //正在公布的一个特定块的哈希
Number     uint64      //公布的一个特定区块的编号
Td         *big.Int    //宣布一个特定区块的总难度
	ReorgDepth uint64
	Update     keyValueList
}

//签名通过给定的私钥向块公告添加签名
func (a *announceData) sign(privKey *ecdsa.PrivateKey) {
	rlp, _ := rlp.EncodeToBytes(announceBlock{a.Hash, a.Number, a.Td})
	sig, _ := crypto.Sign(crypto.Keccak256(rlp), privKey)
	a.Update = a.Update.add("sign", sig)
}

//checksignature验证块通知是否具有给定pubkey的有效签名
func (a *announceData) checkSignature(id enode.ID) error {
	var sig []byte
	if err := a.Update.decode().get("sign", &sig); err != nil {
		return err
	}
	rlp, _ := rlp.EncodeToBytes(announceBlock{a.Hash, a.Number, a.Td})
	recPubkey, err := crypto.SigToPub(crypto.Keccak256(rlp), sig)
	if err != nil {
		return err
	}
	if id == enode.PubkeyToIDV4(recPubkey) {
		return nil
	}
	return errors.New("wrong signature")
}

type blockInfo struct {
Hash   common.Hash //正在公布的一个特定块的哈希
Number uint64      //公布的一个特定区块的编号
Td     *big.Int    //宣布一个特定区块的总难度
}

//GetBlockHeadersData表示块头查询。
type getBlockHeadersData struct {
Origin  hashOrNumber //从中检索邮件头的块
Amount  uint64       //要检索的最大头数
Skip    uint64       //要在连续标题之间跳过的块
Reverse bool         //查询方向（假=上升到最新，真=下降到创世纪）
}

//hashornumber是用于指定源块的组合字段。
type hashOrNumber struct {
Hash   common.Hash //要从中检索头的块哈希（不包括数字）
Number uint64      //要从中检索头的块哈希（不包括哈希）
}

//encoderlp是一个专门的编码器，用于hashornumber只对
//两个包含联合字段。
func (hn *hashOrNumber) EncodeRLP(w io.Writer) error {
	if hn.Hash == (common.Hash{}) {
		return rlp.Encode(w, hn.Number)
	}
	if hn.Number != 0 {
		return fmt.Errorf("both origin hash (%x) and number (%d) provided", hn.Hash, hn.Number)
	}
	return rlp.Encode(w, hn.Hash)
}

//decoderlp是一种特殊的译码器，用于hashornumber对内容进行译码。
//分块散列或分块编号。
func (hn *hashOrNumber) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	origin, err := s.Raw()
	if err == nil {
		switch {
		case size == 32:
			err = rlp.DecodeBytes(origin, &hn.Hash)
		case size <= 8:
			err = rlp.DecodeBytes(origin, &hn.Number)
		default:
			err = fmt.Errorf("invalid input size %d for origin", size)
		}
	}
	return err
}

//codedata是用于节点数据检索的网络响应包。
type CodeData []struct {
	Value []byte
}

type proofsData [][]rlp.RawValue

type txStatus struct {
	Status core.TxStatus
	Lookup *rawdb.TxLookupEntry `rlp:"nil"`
	Error  string
}

