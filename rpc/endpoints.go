
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450108670349312>


package rpc

import (
	"net"

	"github.com/ethereum/go-ethereum/log"
)

//starthttpendpoint启动用cors/vhosts/modules配置的HTTP RPC终结点
func StartHTTPEndpoint(endpoint string, apis []API, modules []string, cors []string, vhosts []string, timeouts HTTPTimeouts) (net.Listener, *Server, error) {
//根据允许的模块生成白名单
	whitelist := make(map[string]bool)
	for _, module := range modules {
		whitelist[module] = true
	}
//注册服务公开的所有API
	handler := NewServer()
	for _, api := range apis {
		if whitelist[api.Namespace] || (len(whitelist) == 0 && api.Public) {
			if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
				return nil, nil, err
			}
			log.Debug("HTTP registered", "namespace", api.Namespace)
		}
	}
//所有已注册的API，启动HTTP侦听器
	var (
		listener net.Listener
		err      error
	)
	if listener, err = net.Listen("tcp", endpoint); err != nil {
		return nil, nil, err
	}
	go NewHTTPServer(cors, vhosts, timeouts, handler).Serve(listener)
	return listener, handler, err
}

//startwsendpoint启动WebSocket终结点
func StartWSEndpoint(endpoint string, apis []API, modules []string, wsOrigins []string, exposeAll bool) (net.Listener, *Server, error) {

//根据允许的模块生成白名单
	whitelist := make(map[string]bool)
	for _, module := range modules {
		whitelist[module] = true
	}
//注册服务公开的所有API
	handler := NewServer()
	for _, api := range apis {
		if exposeAll || whitelist[api.Namespace] || (len(whitelist) == 0 && api.Public) {
			if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
				return nil, nil, err
			}
			log.Debug("WebSocket registered", "service", api.Service, "namespace", api.Namespace)
		}
	}
//所有已注册的API，启动HTTP侦听器
	var (
		listener net.Listener
		err      error
	)
	if listener, err = net.Listen("tcp", endpoint); err != nil {
		return nil, nil, err
	}
	go NewWSServer(wsOrigins, handler).Serve(listener)
	return listener, handler, err

}

//StartIPCendPoint启动IPC终结点。
func StartIPCEndpoint(ipcEndpoint string, apis []API) (net.Listener, *Server, error) {
//注册服务公开的所有API。
	handler := NewServer()
	for _, api := range apis {
		if err := handler.RegisterName(api.Namespace, api.Service); err != nil {
			return nil, nil, err
		}
		log.Debug("IPC registered", "namespace", api.Namespace)
	}
//所有已注册的API，启动IPC侦听器。
	listener, err := ipcListen(ipcEndpoint)
	if err != nil {
		return nil, nil, err
	}
	go handler.ServeListener(listener)
	return listener, handler, nil
}

