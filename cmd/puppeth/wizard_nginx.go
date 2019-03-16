
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450070162444288>


package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"
)

//确保revirtualhost检查是否在指定的
//
//
//
//
//
func (w *wizard) ensureVirtualHost(client *sshClient, port int, def string) (string, error) {
	proxy, _ := checkNginx(client, w.network)
	if proxy != nil {
//反向代理正在运行，如果端口匹配，我们需要一个虚拟主机
		if proxy.port == port {
			fmt.Println()
			fmt.Printf("Shared port, which domain to assign? (default = %s)\n", def)
			return w.readDefaultString(def), nil
		}
	}
//
	fmt.Println()
	fmt.Println("Allow sharing the port with other services (y/n)? (default = yes)")
	if w.readDefaultYesNo(true) {
		nocache := false
		if proxy != nil {
			fmt.Println()
			fmt.Printf("Should the reverse-proxy be rebuilt from scratch (y/n)? (default = no)\n")
			nocache = w.readDefaultYesNo(false)
		}
		if out, err := deployNginx(client, w.network, port, nocache); err != nil {
			log.Error("Failed to deploy reverse-proxy", "err", err)
			if len(out) > 0 {
				fmt.Printf("%s\n", out)
			}
			return "", err
		}
//已部署反向代理，请再次请求虚拟主机
		fmt.Println()
		fmt.Printf("Proxy deployed, which domain to assign? (default = %s)\n", def)
		return w.readDefaultString(def), nil
	}
//
	return "", nil
}

