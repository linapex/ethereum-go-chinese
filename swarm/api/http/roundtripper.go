
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450111920934912>


package http

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/swarm/log"
)

/*
注册BZZ URL方案的HTTP往返器
请参阅https://github.com/ethereum/go-ethereum/issues/2040
用途：

进口（
 “github.com/ethereum/go-ethereum/common/httpclient”
 “github.com/ethereum/go-ethereum/swarm/api/http”
）
客户端：=httpclient.new（）
//对于本地运行的（私有）Swarm代理
client.registerscheme（“bzz”，&http.roundtripper port:port）
client.registerscheme（“bzz不可变”，&http.roundtripper port:port）
client.registerscheme（“bzz raw”，&http.roundtripper port:port）

您给往返者的端口是Swarm代理正在监听的端口。
如果主机为空，则假定为localhost。

使用公共网关，上面的几条线为您提供了最精简的
BZZ方案感知只读HTTP客户端。你真的只需要这个
如果你需要本地群访问BZZ地址。
**/


type RoundTripper struct {
	Host string
	Port string
}

func (self *RoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	host := self.Host
	if len(host) == 0 {
		host = "localhost"
	}
url := fmt.Sprintf("http://%s:%s/%s/%s/%s”，主机，自身端口，请求协议，请求URL.host，请求URL.path）
	log.Info(fmt.Sprintf("roundtripper: proxying request '%s' to '%s'", req.RequestURI, url))
	reqProxy, err := http.NewRequest(req.Method, url, req.Body)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(reqProxy)
}

