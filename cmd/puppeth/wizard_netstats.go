
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450070065975296>


package main

import (
	"encoding/json"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"github.com/olekukonko/tablewriter"
)

//
//配置集，用于向用户提供有关如何执行各种任务的提示。
func (w *wizard) networkStats() {
	if len(w.servers) == 0 {
		log.Info("No remote machines to gather stats from")
		return
	}
//清除一些以前的配置以从当前扫描中重新填充
	w.conf.ethstats = ""
	w.conf.bootnodes = w.conf.bootnodes[:0]

//遍历所有指定主机并检查其状态
	var pend sync.WaitGroup

	stats := make(serverStats)
	for server, pubkey := range w.conf.Servers {
		pend.Add(1)

//同时收集每个服务器的服务状态
		go func(server string, pubkey []byte) {
			defer pend.Done()

			stat := w.gatherStats(server, pubkey, w.servers[server])

//所有状态检查完成，报告并检查下一个服务器
			w.lock.Lock()
			defer w.lock.Unlock()

			delete(w.services, server)
			for service := range stat.services {
				w.services[server] = append(w.services[server], service)
			}
			stats[server] = stat
		}(server, pubkey)
	}
	pend.Wait()

//打印所有收集的统计数据并返回
	stats.render()
}

//GatherStats收集特定远程服务器的服务统计信息。
func (w *wizard) gatherStats(server string, pubkey []byte, client *sshClient) *serverStat {
//收集一些全局统计信息以提供给向导
	var (
		genesis   string
		ethstats  string
		bootnodes []string
	)
//确保到远程服务器的ssh连接有效
	logger := log.New("server", server)
	logger.Info("Starting remote server health-check")

	stat := &serverStat{
		services: make(map[string]map[string]string),
	}
	if client == nil {
		conn, err := dial(server, pubkey)
		if err != nil {
			logger.Error("Failed to establish remote connection", "err", err)
			stat.failure = err.Error()
			return stat
		}
		client = conn
	}
	stat.address = client.address

//客户端以某种方式连接，运行运行运行状况检查
	logger.Debug("Checking for nginx availability")
	if infos, err := checkNginx(client, w.network); err != nil {
		if err != ErrServiceUnknown {
			stat.services["nginx"] = map[string]string{"offline": err.Error()}
		}
	} else {
		stat.services["nginx"] = infos.Report()
	}
	logger.Debug("Checking for ethstats availability")
	if infos, err := checkEthstats(client, w.network); err != nil {
		if err != ErrServiceUnknown {
			stat.services["ethstats"] = map[string]string{"offline": err.Error()}
		}
	} else {
		stat.services["ethstats"] = infos.Report()
		ethstats = infos.config
	}
	logger.Debug("Checking for bootnode availability")
	if infos, err := checkNode(client, w.network, true); err != nil {
		if err != ErrServiceUnknown {
			stat.services["bootnode"] = map[string]string{"offline": err.Error()}
		}
	} else {
		stat.services["bootnode"] = infos.Report()

		genesis = string(infos.genesis)
		bootnodes = append(bootnodes, infos.enode)
	}
	logger.Debug("Checking for sealnode availability")
	if infos, err := checkNode(client, w.network, false); err != nil {
		if err != ErrServiceUnknown {
			stat.services["sealnode"] = map[string]string{"offline": err.Error()}
		}
	} else {
		stat.services["sealnode"] = infos.Report()
		genesis = string(infos.genesis)
	}
	logger.Debug("Checking for explorer availability")
	if infos, err := checkExplorer(client, w.network); err != nil {
		if err != ErrServiceUnknown {
			stat.services["explorer"] = map[string]string{"offline": err.Error()}
		}
	} else {
		stat.services["explorer"] = infos.Report()
	}
	logger.Debug("Checking for wallet availability")
	if infos, err := checkWallet(client, w.network); err != nil {
		if err != ErrServiceUnknown {
			stat.services["wallet"] = map[string]string{"offline": err.Error()}
		}
	} else {
		stat.services["wallet"] = infos.Report()
	}
	logger.Debug("Checking for faucet availability")
	if infos, err := checkFaucet(client, w.network); err != nil {
		if err != ErrServiceUnknown {
			stat.services["faucet"] = map[string]string{"offline": err.Error()}
		}
	} else {
		stat.services["faucet"] = infos.Report()
	}
	logger.Debug("Checking for dashboard availability")
	if infos, err := checkDashboard(client, w.network); err != nil {
		if err != ErrServiceUnknown {
			stat.services["dashboard"] = map[string]string{"offline": err.Error()}
		}
	} else {
		stat.services["dashboard"] = infos.Report()
	}
//
	w.lock.Lock()
	defer w.lock.Unlock()

	if genesis != "" && w.conf.Genesis == nil {
		g := new(core.Genesis)
		if err := json.Unmarshal([]byte(genesis), g); err != nil {
			log.Error("Failed to parse remote genesis", "err", err)
		} else {
			w.conf.Genesis = g
		}
	}
	if ethstats != "" {
		w.conf.ethstats = ethstats
	}
	w.conf.bootnodes = append(w.conf.bootnodes, bootnodes...)

	return stat
}

//serverstat是服务配置参数和运行状况的集合
//检查要打印给用户的报告。
type serverStat struct {
	address  string
	failure  string
	services map[string]map[string]string
}

//ServerStats是多个主机的服务器状态集合。
type serverStats map[string]*serverStat

//
//并打印到标准输出。
func (stats serverStats) render() {
//开始收集服务统计信息和配置参数
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Server", "Address", "Service", "Config", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColWidth(40)

//查找黑客分隔符的所有列的最长行
	separator := make([]string, 5)
	for server, stat := range stats {
		if len(server) > len(separator[0]) {
			separator[0] = strings.Repeat("-", len(server))
		}
		if len(stat.address) > len(separator[1]) {
			separator[1] = strings.Repeat("-", len(stat.address))
		}
		if len(stat.failure) > len(separator[1]) {
			separator[1] = strings.Repeat("-", len(stat.failure))
		}
		for service, configs := range stat.services {
			if len(service) > len(separator[2]) {
				separator[2] = strings.Repeat("-", len(service))
			}
			for config, value := range configs {
				if len(config) > len(separator[3]) {
					separator[3] = strings.Repeat("-", len(config))
				}
				for _, val := range strings.Split(value, "\n") {
					if len(val) > len(separator[4]) {
						separator[4] = strings.Repeat("-", len(val))
					}
				}
			}
		}
	}
//按字母顺序填写服务器报告
	servers := make([]string, 0, len(stats))
	for server := range stats {
		servers = append(servers, server)
	}
	sort.Strings(servers)

	for i, server := range servers {
//在所有服务器之间添加分隔符
		if i > 0 {
			table.Append(separator)
		}
//按字母顺序填写服务报告
		services := make([]string, 0, len(stats[server].services))
		for service := range stats[server].services {
			services = append(services, service)
		}
		sort.Strings(services)

		if len(services) == 0 {
			if stats[server].failure != "" {
				table.Append([]string{server, stats[server].failure, "", "", ""})
			} else {
				table.Append([]string{server, stats[server].address, "", "", ""})
			}
		}
		for j, service := range services {
//在所有服务之间添加空行
			if j > 0 {
				table.Append([]string{"", "", "", separator[3], separator[4]})
			}
//按字母顺序填写配置报告
			configs := make([]string, 0, len(stats[server].services[service]))
			for service := range stats[server].services[service] {
				configs = append(configs, service)
			}
			sort.Strings(configs)

			for k, config := range configs {
				for l, value := range strings.Split(stats[server].services[service][config], "\n") {
					switch {
					case j == 0 && k == 0 && l == 0:
						table.Append([]string{server, stats[server].address, service, config, value})
					case k == 0 && l == 0:
						table.Append([]string{"", "", service, config, value})
					case l == 0:
						table.Append([]string{"", "", "", config, value})
					default:
						table.Append([]string{"", "", "", "", value})
					}
				}
			}
		}
	}
	table.Render()
}

