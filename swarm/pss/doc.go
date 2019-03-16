
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450116513697792>


//PSS为Swarm节点提供了devp2p功能，无需在它们之间建立直接的TCP连接。
//
//消息封装在devp2p消息结构“pssmsg”中。这些胶囊使用普通的tcp devp2p从一个节点转发到另一个节点，直到到达目的地：能够成功解密消息的一个或多个节点。
//
//消息路由使用Swarm自己的Kademlia路由完成。可以选择关闭路由，强制将消息发送到所有对等端，类似于耳语协议的行为。
//
//PSS适用于大小有限的消息，通常最多为数千字节。消息本身可以是任何东西；复杂的数据结构或非描述字节序列。
//
//文档可以在自述文件中找到。
//
//有关PSS开发的当前状态和路线图，请参阅https://github.com/ethersphere/swarm/wiki/swarm-dev-progress。
//
//请在https://github.com/ethersphere/go-ethereum上报告问题
//
//请随时在https://gitter.im/ethersphere/pss中提问
//
//话题
//
//PSS消息的加密信封始终包含主题。这是PSS决定对消息采取什么行动的方式。主题仅对可以解密消息的节点可见。
//
//这个“主题”不像电子邮件的主题，而是一个类似于哈希的任意4字节值。可以使用'pss_u*totopic'API方法生成有效的主题。
//
//PSS中的身份
//
//PSS旨在实现完美的黑暗。这意味着使用PSS进行通信的两个节点的最低要求是共享秘密。这个秘密可以是任意字节片，也可以是ECDSA密钥对。
//
//对等密钥可以通过其api调用'pss'setpeerpublickey'和'pss'setsymetrickey'手动添加到pss节点。键总是与主题耦合，并且这些键只对这些主题有效。
//
//连接
//
//PSS中的“连接”是纯虚拟构造。没有适当的机制来确保远程对等机确实存在。实际上，“添加”一个对等节点只涉及节点认为对等节点在那里的观点。它可以将消息发送给远程对等机，发送给直接连接的对等机，然后由该对等机传递消息。但是，如果它不在网络上——或者如果没有到网络的路由——那么消息将永远无法通过转发到达目的地。
//
//在实现devp2p协议栈时，远程对等端的“添加”是实际启动协议通信的一方的先决条件。添加一个有效的对等机“运行”该对等机上的协议，并在主题和该对等机之间添加一个内部映射。它还可以使用devp2p中的主IO结构（p2p.msgreadwriter）发送和接收消息。
//
//在引擎盖下，PSS实现了自己的msgreadwriter，它将msgreadwriter.writemsg与pss.sendraw连接起来，并巧妙地添加了injectmsg方法，通过管道将传入的消息显示在msgreadwriter.readmsg通道上。
//
//一个传入的连接只不过是一个实际的pssmsg，它出现在一个特定的主题中。如果处理程序har已注册到该主题，则消息将传递给它。如果出现以下情况，则构成“新”连接：
//
//-PSS节点从未使用远程对等地址和主题的组合调用addpeer，以及
//
//-PSS节点以前从未从具有此特定主题的远程对等端接收到PSSMSG。
//
//如果是“新”连接，协议将在远程对等机上“运行”，其方式与预先添加协议的方式相同。
//
package pss

