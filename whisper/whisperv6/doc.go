
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:46</date>
//</624450125208489984>


/*
软件包Whisper实现了Whisper协议（版本6）。

Whisper结合了DHTS和数据报消息系统（如UDP）的各个方面。
因此，可以将其与两者进行比较，而不是与
物质/能量二元性（为明目张胆地滥用
基本而美丽的自然法则）。

Whisper是一个纯粹的基于身份的消息传递系统。低语提供了一个低层次
（非特定于应用程序）但不基于
或者受到低级硬件属性和特性的影响，
尤其是奇点的概念。
**/


//包含耳语协议常量定义

package whisperv6

import (
	"time"
)

//耳语协议参数
const (
ProtocolVersion    = uint64(6) //协议版本号
ProtocolVersionStr = "6.0"     //与字符串相同
ProtocolName       = "shh"     //GETH中协议的昵称

//耳语协议消息代码，根据EIP-627
statusCode           = 0   //由耳语协议使用
messagesCode         = 1   //正常低语信息
powRequirementCode   = 2   //战俘需求
bloomFilterExCode    = 3   //布卢姆过滤器交换
p2pRequestCode       = 126 //点对点消息，由DAPP协议使用
p2pMessageCode       = 127 //对等消息（由对等方使用，但不再转发）
	NumberOfMessageCodes = 128

SizeMask      = byte(3) //用于从标志中提取有效负载大小字段大小的掩码
	signatureFlag = byte(4)

TopicLength     = 4  //以字节为单位
signatureLength = 65 //以字节为单位
aesKeyLength    = 32 //以字节为单位
aesNonceLength  = 12 //以字节为单位；有关详细信息，请参阅cipher.gcmstandardnonconsize&aesgcm.nonconsize（）。
keyIDSize       = 32 //以字节为单位
BloomFilterSize = 64 //以字节为单位
	flagsLength     = 1

	EnvelopeHeaderLength = 20

MaxMessageSize        = uint32(10 * 1024 * 1024) //邮件的最大可接受大小。
	DefaultMaxMessageSize = uint32(1024 * 1024)
	DefaultMinimumPoW     = 0.2

padSizeLimit      = 256 //只是一个任意数字，可以在不破坏协议的情况下进行更改
	messageQueueLimit = 1024

	expirationCycle   = time.Second
	transmissionCycle = 300 * time.Millisecond

DefaultTTL           = 50 //秒
DefaultSyncAllowance = 10 //秒
)

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

