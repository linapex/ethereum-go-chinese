
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450105067442176>


package nat

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/jackpal/go-nat-pmp"
)

//natpmpclient适应nat-pmp协议实现，因此它符合
//公共接口。
type pmp struct {
	gw net.IP
	c  *natpmp.Client
}

func (n *pmp) String() string {
	return fmt.Sprintf("NAT-PMP(%v)", n.gw)
}

func (n *pmp) ExternalIP() (net.IP, error) {
	response, err := n.c.GetExternalAddress()
	if err != nil {
		return nil, err
	}
	return response.ExternalIPAddress[:], nil
}

func (n *pmp) AddMapping(protocol string, extport, intport int, name string, lifetime time.Duration) error {
	if lifetime <= 0 {
		return fmt.Errorf("lifetime must not be <= 0")
	}
//注意，端口参数的顺序在
//addmapping和客户端的addportmapping。
	_, err := n.c.AddPortMapping(strings.ToLower(protocol), intport, extport, int(lifetime/time.Second))
	return err
}

func (n *pmp) DeleteMapping(protocol string, extport, intport int) (err error) {
//若要销毁映射，请发送一个内部端口为的添加端口。
//要销毁的内部端口、零和A的外部端口
//时间为零。
	_, err = n.c.AddPortMapping(strings.ToLower(protocol), intport, 0, 0)
	return err
}

func discoverPMP() Interface {
//在所有可能的网关上运行外部地址查找
	gws := potentialGateways()
	found := make(chan *pmp, len(gws))
	for i := range gws {
		gw := gws[i]
		go func() {
			c := natpmp.NewClient(gw)
			if _, err := c.GetExternalAddress(); err != nil {
				found <- nil
			} else {
				found <- &pmp{gw, c}
			}
		}()
	}
//返回第一个响应的。
//发现需要很快，所以我们不再关心
//在非常短的超时之后的任何响应。
	timeout := time.NewTimer(1 * time.Second)
	defer timeout.Stop()
	for range gws {
		select {
		case c := <-found:
			if c != nil {
				return c
			}
		case <-timeout.C:
			return nil
		}
	}
	return nil
}

var (
//局域网IP范围
	_, lan10, _  = net.ParseCIDR("10.0.0.0/8")
	_, lan176, _ = net.ParseCIDR("172.16.0.0/12")
	_, lan192, _ = net.ParseCIDR("192.168.0.0/16")
)

//托多：改进这个。我们目前假设（在大多数网络上）
//路由器是本地局域网范围内的X.X.X.1。
func potentialGateways() (gws []net.IP) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, iface := range ifaces {
		ifaddrs, err := iface.Addrs()
		if err != nil {
			return gws
		}
		for _, addr := range ifaddrs {
			if x, ok := addr.(*net.IPNet); ok {
				if lan10.Contains(x.IP) || lan176.Contains(x.IP) || lan192.Contains(x.IP) {
					ip := x.IP.Mask(x.Mask).To4()
					if ip != nil {
						ip[3] = ip[3] | 0x01
						gws = append(gws, ip)
					}
				}
			}
		}
	}
	return gws
}

