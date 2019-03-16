
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450109840560128>


package rpc

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"

	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

//API描述了通过RPC接口提供的一组方法
type API struct {
Namespace string      //暴露服务的RPC方法的命名空间
Version   string      //DAPP的API版本
Service   interface{} //保存方法的接收器实例
Public    bool        //是否必须将这些方法视为公共使用安全的指示
}

//回调是在服务器中注册的方法回调
type callback struct {
rcvr        reflect.Value  //方法接受者
method      reflect.Method //回调
argTypes    []reflect.Type //输入参数类型
hasCtx      bool           //方法的第一个参数是上下文（不包括在argtype中）
errPos      int            //当方法无法返回错误时，错误返回IDX，共-1个
isSubscribe bool           //指示回调是否为订阅
}

//服务表示已注册的对象
type service struct {
name          string        //服务的名称
typ           reflect.Type  //接收机类型
callbacks     callbacks     //已注册的处理程序
subscriptions subscriptions //可用订阅/通知
}

//ServerRequest是一个传入请求
type serverRequest struct {
	id            interface{}
	svcname       string
	callb         *callback
	args          []reflect.Value
	isUnsubscribe bool
	err           Error
}

type serviceRegistry map[string]*service //服务收集
type callbacks map[string]*callback      //RPC回调集合
type subscriptions map[string]*callback  //认购回拨集合

//服务器表示RPC服务器
type Server struct {
	services serviceRegistry

	run      int32
	codecsMu sync.Mutex
	codecs   mapset.Set
}

//rpc request表示原始传入的rpc请求
type rpcRequest struct {
	service  string
	method   string
	id       interface{}
	isPubSub bool
	params   interface{}
err      Error //批元素无效
}

//错误包装了RPC错误，其中除消息外还包含错误代码。
type Error interface {
Error() string  //返回消息
ErrorCode() int //返回代码
}

//ServerCodec实现对服务器端的RPC消息的读取、分析和写入
//一个RPC会话。由于可以调用编解码器，因此实现必须是安全的执行例程。
//同时执行多个go例程。
type ServerCodec interface {
//阅读下一个请求
	ReadRequestHeaders() ([]rpcRequest, bool, Error)
//将请求参数解析为给定类型
	ParseRequestArguments(argTypes []reflect.Type, params interface{}) ([]reflect.Value, Error)
//组装成功响应，期望响应ID和有效负载
	CreateResponse(id interface{}, reply interface{}) interface{}
//组装错误响应，需要响应ID和错误
	CreateErrorResponse(id interface{}, err Error) interface{}
//使用有关错误的额外信息通过信息组装错误响应
	CreateErrorResponseWithInfo(id interface{}, err Error, info interface{}) interface{}
//
	CreateNotification(id, namespace string, event interface{}) interface{}
//将消息写入客户端。
	Write(msg interface{}) error
//关闭基础数据流
	Close()
//当基础连接关闭时关闭
	Closed() <-chan interface{}
}

type BlockNumber int64

const (
	PendingBlockNumber  = BlockNumber(-2)
	LatestBlockNumber   = BlockNumber(-1)
	EarliestBlockNumber = BlockNumber(0)
)

//unmashaljson将给定的json片段解析为一个blocknumber。它支持：
//-“最新”、“最早”或“挂起”作为字符串参数
//-区块编号
//返回的错误：
//-当给定参数不是已知字符串时出现无效的块号错误
//-当给定的块号太小或太大时出现超出范围的错误。
func (bn *BlockNumber) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	switch input {
	case "earliest":
		*bn = EarliestBlockNumber
		return nil
	case "latest":
		*bn = LatestBlockNumber
		return nil
	case "pending":
		*bn = PendingBlockNumber
		return nil
	}

	blckNum, err := hexutil.DecodeUint64(input)
	if err != nil {
		return err
	}
	if blckNum > math.MaxInt64 {
		return fmt.Errorf("Blocknumber too high")
	}

	*bn = BlockNumber(blckNum)
	return nil
}

func (bn BlockNumber) Int64() int64 {
	return (int64)(bn)
}

