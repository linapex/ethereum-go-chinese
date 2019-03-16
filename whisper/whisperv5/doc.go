
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:46</date>
//</624450124076027904>


/*
软件包Whisperv5实现了Whisper协议（版本5）。

Whisper结合了DHTS和数据报消息系统（如UDP）的各个方面。
因此，可以将其与两者进行比较，而不是与
物质/能量二元性（为明目张胆地滥用
基本而美丽的自然法则）。

Whisper是一个纯粹的基于身份的消息传递系统。低语提供了一个低层次
（非特定于应用程序）但不基于
或者受到低级硬件属性和特性的影响，
尤其是奇点的概念。
**/

package whisperv5

import (
	"fmt"
	"time"
)

const (
	EnvelopeVersion    = uint64(0)
	ProtocolVersion    = uint64(5)
	ProtocolVersionStr = "5.0"
	ProtocolName       = "shh"

statusCode           = 0 //由耳语协议使用
messagesCode         = 1 //正常低语信息
p2pCode              = 2 //对等消息（由对等方使用，但不再转发）
p2pRequestCode       = 3 //点对点消息，由DAPP协议使用
	NumberOfMessageCodes = 64

	paddingMask   = byte(3)
	signatureFlag = byte(4)

	TopicLength     = 4
	signatureLength = 65
	aesKeyLength    = 32
	AESNonceLength  = 12
	keyIdSize       = 32

MaxMessageSize        = uint32(10 * 1024 * 1024) //邮件的最大可接受大小。
	DefaultMaxMessageSize = uint32(1024 * 1024)
	DefaultMinimumPoW     = 0.2

padSizeLimit      = 256 //只是一个任意数字，可以在不破坏协议的情况下进行更改（不得超过2^24）
	messageQueueLimit = 1024

	expirationCycle   = time.Second
	transmissionCycle = 300 * time.Millisecond

DefaultTTL     = 50 //秒
SynchAllowance = 10 //秒
)

type unknownVersionError uint64

func (e unknownVersionError) Error() string {
	return fmt.Sprintf("invalid envelope version %d", uint64(e))
}

//mail server表示一个邮件服务器，能够
//存档旧邮件以供后续传递
//对同龄人。任何实施都必须确保
//函数是线程安全的。而且，他们必须尽快返回。
//delivermail应使用directmessagescode进行传递，
//以绕过到期检查。
type MailServer interface {
	Archive(env *Envelope)
	DeliverMail(whisperPeer *Peer, request *Envelope)
}

