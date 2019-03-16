
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450095537983488>


package les

import (
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/p2p"
)

var (
 /*proptxninpacketsmeter=metrics.newmeter（“eth/prop/txns/in/packets”）。
  proptxninttrafficmeter=metrics.newmeter（“eth/prop/txns/in/traffic”）。
  proptxnoutpacketsmeter=metrics.newmeter（“eth/prop/txns/out/packets”）。
  NewMeter（“ETH/PROP/TXNS/OUT /流量”）
  propashinpacketsmeter=metrics.newmeter（“eth/prop/hashes/in/packets”）。
  NewMeter（“ET/PROP/散列/IN /流量”）
  propashoutpacketsmeter=metrics.newmeter（“eth/prop/hashes/out/packets”）。
  propashouttrafficmeter=metrics.newmeter（“eth/prop/hashes/out/traffic”）。
  propblockinpacketsmeter=metrics.newmeter（“eth/prop/blocks/in/packets”）。
  propBlockInTrafficMeter   = metrics.NewMeter("eth/prop/blocks/in/traffic")
  propblockoutpacketsmeter=metrics.newmeter（“eth/prop/blocks/out/packets”）。
  PropBlockOutTrafficMeter=metrics.newMeter（“eth/prop/blocks/out/traffic”）。
  reqhashinpacketsmeter=metrics.newmeter（“eth/req/hashes/in/packets”）。
  reqhashIntraffimeter=metrics.newmeter（“eth/req/hashes/in/traffic”）。
  reqhashoutpacketsmeter=metrics.newmeter（“eth/req/hashes/out/packets”）。
  reqhashouttrafficmeter=metrics.newmeter（“eth/req/hashes/out/traffic”）。
  reqBlockInPacketsMeter    = metrics.NewMeter("eth/req/blocks/in/packets")
  reqBlockIntraffimeter=metrics.newmeter（“eth/req/blocks/in/traffic”）。
  reqblockoutpacketsmeter=metrics.newmeter（“eth/req/blocks/out/packets”）。
  reqblockouttrafficmeter=metrics.newmeter（“eth/req/blocks/out/traffic”）。
  reqHeaderInpacketsMeter=metrics.newMeter（“eth/req/headers/in/packets”）。
  reqheaderIntrafficemeter=metrics.newmeter（“eth/req/headers/in/traffic”）。
  reqHeaderOutPacketsMeter=metrics.newMeter（“eth/req/headers/out/packets”）。
  reqHeaderOutTrafficMeter=metrics.newMeter（“eth/req/headers/out/traffic”）。
  reqbodyinpacketsmeter=metrics.newmeter（“eth/req/bodies/in/packets”）。
  reqbodyIntraffimeter=metrics.newmeter（“eth/req/bodies/in/traffic”）。
  reqbodyoutpacketsmeter=metrics.newmeter（“eth/req/body/out/packets”）。
  reqbodyouttrafficmeter=metrics.newmeter（“eth/req/bodys/out/traffic”）。
  reqstateinpacketsmeter=metrics.newmeter（“eth/req/states/in/packets”）。
  reqstateIntraffimeter=metrics.newmeter（“eth/req/states/in/traffic”）。
  reqStateOutPacketsMeter   = metrics.NewMeter("eth/req/states/out/packets")
  reqstateouttrafficmeter=metrics.newmeter（“eth/req/states/out/traffic”）。
  reqReceiptInPacketsMeter  = metrics.NewMeter("eth/req/receipts/in/packets")
  reqReceiptInterffimeter=metrics.newmeter（“eth/req/receipts/in/traffic”）。
  reqReceiptOutPacketsMeter=metrics.newMeter（“eth/req/receipts/out/packets”）。
  reqReceiptOutTrafficMeter=metrics.newMeter（“eth/req/receipts/out/traffic”*/

	miscInPacketsMeter  = metrics.NewRegisteredMeter("les/misc/in/packets", nil)
	miscInTrafficMeter  = metrics.NewRegisteredMeter("les/misc/in/traffic", nil)
	miscOutPacketsMeter = metrics.NewRegisteredMeter("les/misc/out/packets", nil)
	miscOutTrafficMeter = metrics.NewRegisteredMeter("les/misc/out/traffic", nil)
)

//meteredmsgreadwriter是p2p.msgreadwriter的包装器，能够
//基于数据流内容累积上述定义的度量。
type meteredMsgReadWriter struct {
p2p.MsgReadWriter     //将消息流打包到仪表
version           int //选择正确仪表的协议版本
}

//newmeteredmsgwriter使用计量支持包装p2p msgreadwriter。如果
//度量系统被禁用，此函数返回原始对象。
func newMeteredMsgWriter(rw p2p.MsgReadWriter) p2p.MsgReadWriter {
	if !metrics.Enabled {
		return rw
	}
	return &meteredMsgReadWriter{MsgReadWriter: rw}
}

//init设置流使用的协议版本，以知道要
//协议版本之间的消息ID重叠时递增。
func (rw *meteredMsgReadWriter) Init(version int) {
	rw.version = version
}

func (rw *meteredMsgReadWriter) ReadMsg() (p2p.Msg, error) {
//读取信息，并在出现错误时短路
	msg, err := rw.MsgReadWriter.ReadMsg()
	if err != nil {
		return msg, err
	}
//计算数据流量
	packets, traffic := miscInPacketsMeter, miscInTrafficMeter
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	return msg, err
}

func (rw *meteredMsgReadWriter) WriteMsg(msg p2p.Msg) error {
//计算数据流量
	packets, traffic := miscOutPacketsMeter, miscOutTrafficMeter
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

//将数据包发送到P2P层
	return rw.MsgReadWriter.WriteMsg(msg)
}

