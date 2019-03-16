
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450069139034112>


package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

var (
//当服务容器不存在时，返回errServiceUnknown。
	ErrServiceUnknown = errors.New("service unknown")

//存在服务容器时返回errServiceOffline，但不存在
//跑步。
	ErrServiceOffline = errors.New("service offline")

//当服务容器正在运行时返回errServiceUnreachable，但是
//似乎对沟通尝试没有反应。
	ErrServiceUnreachable = errors.New("service unreachable")

//如果Web服务没有公开的端口，则返回errnotexposed，或者
//它前面的反向代理来转发请求。
	ErrNotExposed = errors.New("service not exposed, nor proxied")
)

//ContainerInfos是大型检验数据集的一个大大简化的版本。
//从Docker-Inspect返回，被puppeth解析成易于使用的形式。
type containerInfos struct {
running bool              //标记容器当前是否正在运行
envvars map[string]string //容器上设置的环境变量集合
portmap map[string]int    //从内部端口/协议组合到主机绑定的端口映射
volumes map[string]string //从容器到主机目录的卷装入点
}

//InspectContainer运行Docker根据正在运行的容器进行检查
func inspectContainer(client *sshClient, container string) (*containerInfos, error) {
//检查是否有正在运行服务的容器
	out, err := client.Run(fmt.Sprintf("docker inspect %s", container))
	if err != nil {
		return nil, ErrServiceUnknown
	}
//如果是，则提取各种配置选项
	type inspection struct {
		State struct {
			Running bool
		}
		Mounts []struct {
			Source      string
			Destination string
		}
		Config struct {
			Env []string
		}
		HostConfig struct {
			PortBindings map[string][]map[string]string
		}
	}
	var inspects []inspection
	if err = json.Unmarshal(out, &inspects); err != nil {
		return nil, err
	}
	inspect := inspects[0]

//检索到的信息，将上面的内容解析为有意义的内容
	infos := &containerInfos{
		running: inspect.State.Running,
		envvars: make(map[string]string),
		portmap: make(map[string]int),
		volumes: make(map[string]string),
	}
	for _, envvar := range inspect.Config.Env {
		if parts := strings.Split(envvar, "="); len(parts) == 2 {
			infos.envvars[parts[0]] = parts[1]
		}
	}
	for portname, details := range inspect.HostConfig.PortBindings {
		if len(details) > 0 {
			port, _ := strconv.Atoi(details[0]["HostPort"])
			infos.portmap[portname] = port
		}
	}
	for _, mount := range inspect.Mounts {
		infos.volumes[mount.Destination] = mount.Source
	}
	return infos, err
}

//TearDown通过ssh连接到远程计算机并终止Docker容器
//在指定网络中以指定名称运行。
func tearDown(client *sshClient, network string, service string, purge bool) ([]byte, error) {
//拆下正在运行（或暂停）的容器
	out, err := client.Run(fmt.Sprintf("docker rm -f %s_%s_1", network, service))
	if err != nil {
		return out, err
	}
//如果需要，也清除关联的Docker映像
	if purge {
		return client.Run(fmt.Sprintf("docker rmi %s/%s", network, service))
	}
	return nil, nil
}

//resolve通过返回
//实际的服务器名称和端口，或者最好是nginx虚拟主机（如果可用）。
func resolve(client *sshClient, network string, service string, port int) (string, error) {
//检查服务以从中获取各种配置
	infos, err := inspectContainer(client, fmt.Sprintf("%s_%s_1", network, service))
	if err != nil {
		return "", err
	}
	if !infos.running {
		return "", ErrServiceOffline
	}
//在线容器，提取任何环境变量
	if vhost := infos.envvars["VIRTUAL_HOST"]; vhost != "" {
		return vhost, nil
	}
	return fmt.Sprintf("%s:%d", client.server, port), nil
}

//checkport尝试连接到给定主机上的远程主机
func checkPort(host string, port int) error {
	log.Trace("Verifying remote TCP connectivity", "server", host, "port", port)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), time.Second)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

