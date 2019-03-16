
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450070112112640>


package main

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/log"
)

//
//连接到新服务器的选项。
func (w *wizard) manageServers() {
//
	fmt.Println()

	servers := w.conf.servers()
	for i, server := range servers {
		fmt.Printf(" %d. Disconnect %s\n", i+1, server)
	}
	fmt.Printf(" %d. Connect another server\n", len(w.conf.Servers)+1)

	choice := w.readInt()
	if choice < 0 || choice > len(w.conf.Servers)+1 {
		log.Error("Invalid server choice, aborting")
		return
	}
//
	if choice <= len(w.conf.Servers) {
		server := servers[choice-1]
		client := w.servers[server]

		delete(w.servers, server)
		if client != nil {
			client.Close()
		}
		delete(w.conf.Servers, server)
		w.conf.flush()

		log.Info("Disconnected existing server", "server", server)
		w.networkStats()
		return
	}
//
	if w.makeServer() != "" {
		w.networkStats()
	}
}

//makeserver从stdin读取一行并将其解释为
//
//
//
//如果连接成功，服务器将添加到向导配置中！
func (w *wizard) makeServer() string {
	fmt.Println()
	fmt.Println("What is the remote server's address ([username[:identity]@]hostname[:port])?")

//读取并拨号服务器以确保Docker存在
	input := w.readString()

	client, err := dial(input, nil)
	if err != nil {
		log.Error("Server not ready for puppeth", "err", err)
		return ""
	}
//
	w.servers[input] = client
	w.conf.Servers[input] = client.pubkey
	w.conf.flush()

	return input
}

//selectserver列出用户要从中选择的所有当前已知服务器，
//
func (w *wizard) selectServer() string {
//
	fmt.Println()
	fmt.Println("Which server do you want to interact with?")

	servers := w.conf.servers()
	for i, server := range servers {
		fmt.Printf(" %d. %s\n", i+1, server)
	}
	fmt.Printf(" %d. Connect another server\n", len(w.conf.Servers)+1)

	choice := w.readInt()
	if choice < 0 || choice > len(w.conf.Servers)+1 {
		log.Error("Invalid server choice, aborting")
		return ""
	}
//
	if choice <= len(w.conf.Servers) {
		return servers[choice-1]
	}
	return w.makeServer()
}

//
//
func (w *wizard) manageComponents() {
//
	fmt.Println()

	var serviceHosts, serviceNames []string
	for server, services := range w.services {
		for _, service := range services {
			serviceHosts = append(serviceHosts, server)
			serviceNames = append(serviceNames, service)

			fmt.Printf(" %d. Tear down %s on %s\n", len(serviceHosts), strings.Title(service), server)
		}
	}
	fmt.Printf(" %d. Deploy new network component\n", len(serviceHosts)+1)

	choice := w.readInt()
	if choice < 0 || choice > len(serviceHosts)+1 {
		log.Error("Invalid component choice, aborting")
		return
	}
//
	if choice <= len(serviceHosts) {
//找出要销毁和执行它的服务
		service := serviceNames[choice-1]
		server := serviceHosts[choice-1]
		client := w.servers[server]

		if out, err := tearDown(client, w.network, service, true); err != nil {
			log.Error("Failed to tear down component", "err", err)
			if len(out) > 0 {
				fmt.Printf("%s\n", out)
			}
			return
		}
//
		services := w.services[server]
		for i, name := range services {
			if name == service {
				w.services[server] = append(services[:i], services[i+1:]...)
				if len(w.services[server]) == 0 {
					delete(w.services, server)
				}
			}
		}
		log.Info("Torn down existing component", "server", server, "service", service)
		return
	}
//
	w.deployComponent()
}

//
//
func (w *wizard) deployComponent() {
//
	fmt.Println()
	fmt.Println("What would you like to deploy? (recommended order)")
	fmt.Println(" 1. Ethstats  - Network monitoring tool")
	fmt.Println(" 2. Bootnode  - Entry point of the network")
	fmt.Println(" 3. Sealer    - Full node minting new blocks")
	fmt.Println(" 4. Explorer  - Chain analysis webservice (ethash only)")
	fmt.Println(" 5. Wallet    - Browser wallet for quick sends")
	fmt.Println(" 6. Faucet    - Crypto faucet to give away funds")
	fmt.Println(" 7. Dashboard - Website listing above web-services")

	switch w.read() {
	case "1":
		w.deployEthstats()
	case "2":
		w.deployNode(true)
	case "3":
		w.deployNode(false)
	case "4":
		w.deployExplorer()
	case "5":
		w.deployWallet()
	case "6":
		w.deployFaucet()
	case "7":
		w.deployDashboard()
	default:
		log.Error("That's not something I can do")
	}
}

