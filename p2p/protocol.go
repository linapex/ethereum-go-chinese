
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450105721753600>


package p2p

import (
	"fmt"

	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/p2p/enr"
)

//协议表示P2P子协议实现。
type Protocol struct {
//名称应包含官方协议名称，
//通常是三个字母的单词。
	Name string

//版本应包含协议的版本号。
	Version uint

//长度应包含使用的消息代码数
//按照协议。
	Length uint64

//当协议
//与同行协商。它应该读写来自
//RW。每个消息的有效负载必须完全消耗。
//
//当Start返回时，对等连接将关闭。它应该会回来
//任何协议级错误（如I/O错误），即
//遇到。
	Run func(peer *Peer, rw MsgReadWriter) error

//nodeinfo是用于检索协议特定元数据的可选助手方法
//关于主机节点。
	NodeInfo func() interface{}

//peerinfo是一个可选的帮助器方法，用于检索协议特定的元数据
//关于网络中的某个对等点。如果设置了信息检索功能，
//但返回nil，假设协议握手仍在运行。
	PeerInfo func(id enode.ID) interface{}

//属性包含节点记录的协议特定信息。
	Attributes []enr.Entry
}

func (p Protocol) cap() Cap {
	return Cap{p.Name, p.Version}
}

//cap是对等能力的结构。
type Cap struct {
	Name    string
	Version uint
}

func (cap Cap) String() string {
	return fmt.Sprintf("%s/%d", cap.Name, cap.Version)
}

type capsByNameAndVersion []Cap

func (cs capsByNameAndVersion) Len() int      { return len(cs) }
func (cs capsByNameAndVersion) Swap(i, j int) { cs[i], cs[j] = cs[j], cs[i] }
func (cs capsByNameAndVersion) Less(i, j int) bool {
	return cs[i].Name < cs[j].Name || (cs[i].Name == cs[j].Name && cs[i].Version < cs[j].Version)
}

func (capsByNameAndVersion) ENRKey() string { return "cap" }

