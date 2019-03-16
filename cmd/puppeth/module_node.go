
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450069529104384>


package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

//
var nodeDockerfile = `
FROM ethereum/client-go:latest

ADD genesis.json /genesis.json
{{if .Unlock}}
	ADD signer.json /signer.json
	ADD signer.pass /signer.pass
{{end}}
RUN \
  echo 'geth --cache 512 init /genesis.json' > geth.sh && \{{if .Unlock}}
	echo 'mkdir -p /root/.ethereum/keystore/ && cp /signer.json /root/.ethereum/keystore/' >> geth.sh && \{{end}}
	echo $'exec geth --networkid {{.NetworkID}} --cache 512 --port {{.Port}} --nat extip:{{.IP}} --maxpeers {{.Peers}} {{.LightFlag}} --ethstats \'{{.Ethstats}}\' {{if .Bootnodes}}--bootnodes {{.Bootnodes}}{{end}} {{if .Etherbase}}--miner.etherbase {{.Etherbase}} --mine --miner.threads 1{{end}} {{if .Unlock}}--unlock 0 --password /signer.pass --mine{{end}} --miner.gastarget {{.GasTarget}} --miner.gaslimit {{.GasLimit}} --miner.gasprice {{.GasPrice}}' >> geth.sh

ENTRYPOINT ["/bin/sh", "geth.sh"]
`

//
//
var nodeComposefile = `
version: '2'
services:
  {{.Type}}:
    build: .
    image: {{.Network}}/{{.Type}}
    container_name: {{.Network}}_{{.Type}}_1
    ports:
      - "{{.Port}}:{{.Port}}"
      - "{{.Port}}:{{.Port}}/udp"
    volumes:
      - {{.Datadir}}:/root/.ethereum{{if .Ethashdir}}
      - {{.Ethashdir}}:/root/.ethash{{end}}
    environment:
      - PORT={{.Port}}/tcp
      - TOTAL_PEERS={{.TotalPeers}}
      - LIGHT_PEERS={{.LightPeers}}
      - STATS_NAME={{.Ethstats}}
      - MINER_NAME={{.Etherbase}}
      - GAS_TARGET={{.GasTarget}}
      - GAS_LIMIT={{.GasLimit}}
      - GAS_PRICE={{.GasPrice}}
    logging:
      driver: "json-file"
      options:
        max-size: "1m"
        max-file: "10"
    restart: always
`

//
//Docker和Docker组合。如果具有指定网络名称的实例
//已经存在，将被覆盖！
func deployNode(client *sshClient, network string, bootnodes []string, config *nodeInfos, nocache bool) ([]byte, error) {
	kind := "sealnode"
	if config.keyJSON == "" && config.etherbase == "" {
		kind = "bootnode"
		bootnodes = make([]string, 0)
	}
//生成要上载到服务器的内容
	workdir := fmt.Sprintf("%d", rand.Int63())
	files := make(map[string][]byte)

	lightFlag := ""
	if config.peersLight > 0 {
		lightFlag = fmt.Sprintf("--lightpeers=%d --lightserv=50", config.peersLight)
	}
	dockerfile := new(bytes.Buffer)
	template.Must(template.New("").Parse(nodeDockerfile)).Execute(dockerfile, map[string]interface{}{
		"NetworkID": config.network,
		"Port":      config.port,
		"IP":        client.address,
		"Peers":     config.peersTotal,
		"LightFlag": lightFlag,
		"Bootnodes": strings.Join(bootnodes, ","),
		"Ethstats":  config.ethstats,
		"Etherbase": config.etherbase,
		"GasTarget": uint64(1000000 * config.gasTarget),
		"GasLimit":  uint64(1000000 * config.gasLimit),
		"GasPrice":  uint64(1000000000 * config.gasPrice),
		"Unlock":    config.keyJSON != "",
	})
	files[filepath.Join(workdir, "Dockerfile")] = dockerfile.Bytes()

	composefile := new(bytes.Buffer)
	template.Must(template.New("").Parse(nodeComposefile)).Execute(composefile, map[string]interface{}{
		"Type":       kind,
		"Datadir":    config.datadir,
		"Ethashdir":  config.ethashdir,
		"Network":    network,
		"Port":       config.port,
		"TotalPeers": config.peersTotal,
		"Light":      config.peersLight > 0,
		"LightPeers": config.peersLight,
		"Ethstats":   config.ethstats[:strings.Index(config.ethstats, ":")],
		"Etherbase":  config.etherbase,
		"GasTarget":  config.gasTarget,
		"GasLimit":   config.gasLimit,
		"GasPrice":   config.gasPrice,
	})
	files[filepath.Join(workdir, "docker-compose.yaml")] = composefile.Bytes()

	files[filepath.Join(workdir, "genesis.json")] = config.genesis
	if config.keyJSON != "" {
		files[filepath.Join(workdir, "signer.json")] = []byte(config.keyJSON)
		files[filepath.Join(workdir, "signer.pass")] = []byte(config.keyPass)
	}
//将部署文件上载到远程服务器（然后清理）
	if out, err := client.Upload(files); err != nil {
		return out, err
	}
	defer client.Run("rm -rf " + workdir)

//构建和部署引导或密封节点服务
	if nocache {
		return nil, client.Stream(fmt.Sprintf("cd %s && docker-compose -p %s build --pull --no-cache && docker-compose -p %s up -d --force-recreate --timeout 60", workdir, network, network))
	}
	return nil, client.Stream(fmt.Sprintf("cd %s && docker-compose -p %s up -d --build --force-recreate --timeout 60", workdir, network))
}

//
//各种配置参数。
type nodeInfos struct {
	genesis    []byte
	network    int64
	datadir    string
	ethashdir  string
	ethstats   string
	port       int
	enode      string
	peersTotal int
	peersLight int
	etherbase  string
	keyJSON    string
	keyPass    string
	gasTarget  float64
	gasLimit   float64
	gasPrice   float64
}

//报表将类型化结构转换为纯字符串->字符串映射，其中包含
//大多数（但不是全部）字段用于向用户报告。
func (info *nodeInfos) Report() map[string]string {
	report := map[string]string{
		"Data directory":           info.datadir,
		"Listener port":            strconv.Itoa(info.port),
		"Peer count (all total)":   strconv.Itoa(info.peersTotal),
		"Peer count (light nodes)": strconv.Itoa(info.peersLight),
		"Ethstats username":        info.ethstats,
	}
	if info.gasTarget > 0 {
//
		report["Gas price (minimum accepted)"] = fmt.Sprintf("%0.3f GWei", info.gasPrice)
		report["Gas floor (baseline target)"] = fmt.Sprintf("%0.3f MGas", info.gasTarget)
		report["Gas ceil  (target maximum)"] = fmt.Sprintf("%0.3f MGas", info.gasLimit)

		if info.etherbase != "" {
//
			report["Ethash directory"] = info.ethashdir
			report["Miner account"] = info.etherbase
		}
		if info.keyJSON != "" {
//
			var key struct {
				Address string `json:"address"`
			}
			if err := json.Unmarshal([]byte(info.keyJSON), &key); err == nil {
				report["Signer account"] = common.HexToAddress(key.Address).Hex()
			} else {
				log.Error("Failed to retrieve signer address", "err", err)
			}
		}
	}
	return report
}

//
//
func checkNode(client *sshClient, network string, boot bool) (*nodeInfos, error) {
	kind := "bootnode"
	if !boot {
		kind = "sealnode"
	}
//
	infos, err := inspectContainer(client, fmt.Sprintf("%s_%s_1", network, kind))
	if err != nil {
		return nil, err
	}
	if !infos.running {
		return nil, ErrServiceOffline
	}
//
	totalPeers, _ := strconv.Atoi(infos.envvars["TOTAL_PEERS"])
	lightPeers, _ := strconv.Atoi(infos.envvars["LIGHT_PEERS"])
	gasTarget, _ := strconv.ParseFloat(infos.envvars["GAS_TARGET"], 64)
	gasLimit, _ := strconv.ParseFloat(infos.envvars["GAS_LIMIT"], 64)
	gasPrice, _ := strconv.ParseFloat(infos.envvars["GAS_PRICE"], 64)

//
	var out []byte
	if out, err = client.Run(fmt.Sprintf("docker exec %s_%s_1 geth --exec admin.nodeInfo.enode --cache=16 attach", network, kind)); err != nil {
		return nil, ErrServiceUnreachable
	}
	enode := bytes.Trim(bytes.TrimSpace(out), "\"")

	if out, err = client.Run(fmt.Sprintf("docker exec %s_%s_1 cat /genesis.json", network, kind)); err != nil {
		return nil, ErrServiceUnreachable
	}
	genesis := bytes.TrimSpace(out)

	keyJSON, keyPass := "", ""
	if out, err = client.Run(fmt.Sprintf("docker exec %s_%s_1 cat /signer.json", network, kind)); err == nil {
		keyJSON = string(bytes.TrimSpace(out))
	}
	if out, err = client.Run(fmt.Sprintf("docker exec %s_%s_1 cat /signer.pass", network, kind)); err == nil {
		keyPass = string(bytes.TrimSpace(out))
	}
//运行健全性检查以查看是否可以访问devp2p
	port := infos.portmap[infos.envvars["PORT"]]
	if err = checkPort(client.server, port); err != nil {
		log.Warn(fmt.Sprintf("%s devp2p port seems unreachable", strings.Title(kind)), "server", client.server, "port", port, "err", err)
	}
//收集并返回有用的信息
	stats := &nodeInfos{
		genesis:    genesis,
		datadir:    infos.volumes["/root/.ethereum"],
		ethashdir:  infos.volumes["/root/.ethash"],
		port:       port,
		peersTotal: totalPeers,
		peersLight: lightPeers,
		ethstats:   infos.envvars["STATS_NAME"],
		etherbase:  infos.envvars["MINER_NAME"],
		keyJSON:    keyJSON,
		keyPass:    keyPass,
		gasTarget:  gasTarget,
		gasLimit:   gasLimit,
		gasPrice:   gasPrice,
	}
	stats.enode = string(enode)

	return stats, nil
}

