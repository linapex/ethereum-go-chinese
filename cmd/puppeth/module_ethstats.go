
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450069281640448>


package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/log"
)

//ethstatsDockerfile是构建ethstats后端所需的Dockerfile。
//以及相关的监测站点。
var ethstatsDockerfile = `
FROM puppeth/ethstats:latest

RUN echo 'module.exports = {trusted: [{{.Trusted}}], banned: [{{.Banned}}], reserved: ["yournode"]};' > lib/utils/config.js
`

//ethstatscomposefile是部署和
//维护一个ethstats监控站点。
var ethstatsComposefile = `
version: '2'
services:
  ethstats:
    build: .
    image: {{.Network}}/ethstats
    container_name: {{.Network}}_ethstats_1{{if not .VHost}}
    ports:
      - "{{.Port}}:3000"{{end}}
    environment:
      - WS_SECRET={{.Secret}}{{if .VHost}}
      - VIRTUAL_HOST={{.VHost}}{{end}}{{if .Banned}}
      - BANNED={{.Banned}}{{end}}
    logging:
      driver: "json-file"
      options:
        max-size: "1m"
        max-file: "10"
    restart: always
`

//deployethstats通过ssh将新的ethstats容器部署到远程计算机，
//Docker和Docker组合。如果具有指定网络名称的实例
//已经存在，将被覆盖！
func deployEthstats(client *sshClient, network string, port int, secret string, vhost string, trusted []string, banned []string, nocache bool) ([]byte, error) {
//生成要上载到服务器的内容
	workdir := fmt.Sprintf("%d", rand.Int63())
	files := make(map[string][]byte)

	trustedLabels := make([]string, len(trusted))
	for i, address := range trusted {
		trustedLabels[i] = fmt.Sprintf("\"%s\"", address)
	}
	bannedLabels := make([]string, len(banned))
	for i, address := range banned {
		bannedLabels[i] = fmt.Sprintf("\"%s\"", address)
	}

	dockerfile := new(bytes.Buffer)
	template.Must(template.New("").Parse(ethstatsDockerfile)).Execute(dockerfile, map[string]interface{}{
		"Trusted": strings.Join(trustedLabels, ", "),
		"Banned":  strings.Join(bannedLabels, ", "),
	})
	files[filepath.Join(workdir, "Dockerfile")] = dockerfile.Bytes()

	composefile := new(bytes.Buffer)
	template.Must(template.New("").Parse(ethstatsComposefile)).Execute(composefile, map[string]interface{}{
		"Network": network,
		"Port":    port,
		"Secret":  secret,
		"VHost":   vhost,
		"Banned":  strings.Join(banned, ","),
	})
	files[filepath.Join(workdir, "docker-compose.yaml")] = composefile.Bytes()

//将部署文件上载到远程服务器（然后清理）
	if out, err := client.Upload(files); err != nil {
		return out, err
	}
	defer client.Run("rm -rf " + workdir)

//构建和部署ethstats服务
	if nocache {
		return nil, client.Stream(fmt.Sprintf("cd %s && docker-compose -p %s build --pull --no-cache && docker-compose -p %s up -d --force-recreate --timeout 60", workdir, network, network))
	}
	return nil, client.Stream(fmt.Sprintf("cd %s && docker-compose -p %s up -d --build --force-recreate --timeout 60", workdir, network))
}

//ethstatsInfo从ethstats状态检查返回以允许报告
//各种配置参数。
type ethstatsInfos struct {
	host   string
	port   int
	secret string
	config string
	banned []string
}

//报表将类型化结构转换为纯字符串->字符串映射，其中包含
//大多数（但不是全部）字段用于向用户报告。
func (info *ethstatsInfos) Report() map[string]string {
	return map[string]string{
		"Website address":       info.host,
		"Website listener port": strconv.Itoa(info.port),
		"Login secret":          info.secret,
		"Banned addresses":      strings.Join(info.banned, "\n"),
	}
}

//check ethstats对ethstats服务器执行运行状况检查，以验证
//它正在运行，如果是，收集有关它的有用信息。
func checkEthstats(client *sshClient, network string) (*ethstatsInfos, error) {
//检查主机上可能的ethstats容器
	infos, err := inspectContainer(client, fmt.Sprintf("%s_ethstats_1", network))
	if err != nil {
		return nil, err
	}
	if !infos.running {
		return nil, ErrServiceOffline
	}
//从主机或反向代理解析端口
	port := infos.portmap["3000/tcp"]
	if port == 0 {
		if proxy, _ := checkNginx(client, network); proxy != nil {
			port = proxy.port
		}
	}
	if port == 0 {
		return nil, ErrNotExposed
	}
//从反向代理解析主机并配置连接字符串
	host := infos.envvars["VIRTUAL_HOST"]
	if host == "" {
		host = client.server
	}
	secret := infos.envvars["WS_SECRET"]
	config := fmt.Sprintf("%s@%s", secret, host)
	if port != 80 && port != 443 {
		config += fmt.Sprintf(":%d", port)
	}
//检索IP黑名单
	banned := strings.Split(infos.envvars["BANNED"], ",")

//运行健全检查以查看端口是否可访问
	if err = checkPort(host, port); err != nil {
		log.Warn("Ethstats service seems unreachable", "server", host, "port", port, "err", err)
	}
//容器可用，组装并返回有用的信息
	return &ethstatsInfos{
		host:   host,
		port:   port,
		secret: secret,
		config: config,
		banned: banned,
	}, nil
}

