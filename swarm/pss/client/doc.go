
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450116475949056>


//实现PSS功能的简单抽象
//
//PSS客户端库旨在简化在PSS上使用p2p.protocols包的过程。
//
//IO使用普通的p2p.msgreadwriter接口执行，该接口使用websockets作为传输层，使用swarm/pss包中pssapi类中的方法，通过rpc透明地与pss节点通信。
//
//
//最小ISH使用示例（需要具有WebSocket RPC的正在运行的PSS节点）：
//
//
//进口（
//“语境”
//“FMT”
//“操作系统”
//pss“github.com/ethereum/go-ethereum/swarm/pss/client”
//“github.com/ethereum/go-ethereum/p2p/协议”
//“github.com/ethereum/go-ethereum/p2p”
//“github.com/ethereum/go-ethereum/swarm/pot”
//“github.com/ethereum/go-ethereum/swarm/log”
//）
//
//FOOMSG结构类型
//条形图
//}
//
//
//func foohandler（msg interface）错误
//foomsg，确定：=msg.（*foomsg）
//如果OK {
//log.debug（“yay，刚收到一条消息”，“msg”，foomsg）
//}
//返回errors.new（fmt.sprintf（“未知消息”））
//}
//
//规格：=&protocols.spec
//姓名：“福”，
//版本：1，
//最大尺寸：1024，
//消息：[]接口
//FoMsg{}
//}
//}
//
//协议：=&p2p.协议
//名称：规格名称，
//版本：规范版本，
//长度：uint64（len（spec.messages）），
//运行：func（p*p2p.peer，rw p2p.msgreadwriter）错误
//p p：=protocols.newpeer（p，rw，spec）
//返回PP.RUN（Foohandler）
//}
//}
//
//func实现（）
//cfg：=pss.newclientconfig（）。
//psc：=pss.newclient（context.background（），nil，cfg）
//错误：=psc.start（））
//如果犯错！= nIL{
//log.crit（“无法启动PSS客户端”）
//退出（1）
//}
//
//log.debug（“连接到PSS节点”，“bzz addr”，psc.baseaddr）
//
//err=psc.runprotocol（协议）
//如果犯错！= nIL{
//log.crit（“无法在PSS WebSocket上启动协议”）
//退出（1）
//}
//
//地址：=pot.randomaddress（）//当然应该是一个真实地址
//psc.addpsspeer（地址，规格）
//
////使用协议
//
//停止（）
//}
//
//bug（测试）：由于蜂群蜂巢中的死锁问题，测试超时
package client

