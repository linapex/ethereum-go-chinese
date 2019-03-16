
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450102320173056>


/*
包节点设置多协议以太坊节点。

在这个包公开的模型中，节点是使用共享的服务的集合。
提供RPC API的资源。服务还可以提供DEVP2P协议，这些协议是有线的
在节点实例启动时到达devp2p网络。


节点管理的资源

节点实例使用的所有文件系统资源都位于名为
数据目录。可以通过附加节点覆盖每个资源的位置
配置。数据目录是可选的。如果没有设置，则
资源未指定，包节点将在内存中创建资源。

要访问devp2p网络，节点配置并启动p2p.server。上的每个主机
devp2p网络有一个唯一的标识符，即节点密钥。节点实例保持此密钥
重新启动。节点还加载静态和受信任的节点列表，并确保
关于其他主机是持久化的。

运行http、websocket或ipc的json-rpc服务器可以在节点上启动。RPC模块
由注册服务提供的服务将在这些端点上提供。用户可以限制任何
rpc模块子集的终结点。节点本身提供“debug”、“admin”和“web3”
模块。

服务实现可以通过服务上下文打开级别数据库。包裹
节点选择每个数据库的文件系统位置。如果节点配置为运行
如果没有数据目录，数据库将在内存中打开。

节点还创建加密的以太坊帐户密钥的共享存储。服务可以访问
客户经理通过服务上下文。


实例间共享数据目录

如果多个节点实例具有不同的实例，则它们可以共享一个数据目录。
名称（通过名称配置选项设置）。共享行为取决于
资源。

与devp2p相关的资源（节点键、静态/可信节点列表、已知主机数据库）是
存储在与实例同名的目录中。因此，多个节点实例
使用相同的数据目录将此信息存储在
数据目录。

leveldb数据库也存储在instance子目录中。如果多个节点
实例使用相同的数据目录，打开具有相同名称的数据库将
为每个实例创建一个数据库。

帐户密钥存储在使用相同数据目录的所有节点实例之间共享
除非通过keystoredir配置选项更改其位置。


数据目录共享示例

在本例中，两个名为a和b的节点实例是以相同的数据开始的
目录。节点实例A打开数据库“db”，节点实例B打开数据库
“DB”和“DB-2”。将在数据目录中创建以下文件：

   数据目录
        a/
            node key—实例A的devp2p节点键
            节点/--实例A的devp2p发现知识库
            db/——“db”的leveldb内容
        a.ipc——实例A的json-rpc unix域套接字端点
        B/
            node key——b节点的devp2p节点键
            节点/--实例B的devp2p发现知识库
            static-nodes.json--实例B的devp2p静态节点列表
            db/——“db”的leveldb内容
            db-2/——db-2的级别db内容
        b.ipc——实例b的json-rpc unix域套接字端点
        key store/--帐户密钥存储，由两个实例使用
**/

package node

