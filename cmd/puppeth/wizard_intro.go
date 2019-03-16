
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450070019837952>


package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/log"
)

//
func makeWizard(network string) *wizard {
	return &wizard{
		network: network,
		conf: config{
			Servers: make(map[string][]byte),
		},
		servers:  make(map[string]*sshClient),
		services: make(map[string][]string),
		in:       bufio.NewReader(os.Stdin),
	}
}

//
//
func (w *wizard) run() {
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("| Welcome to puppeth, your Ethereum private network manager |")
	fmt.Println("|                                                           |")
	fmt.Println("| This tool lets you create a new Ethereum network down to  |")
	fmt.Println("| the genesis block, bootnodes, miners and ethstats servers |")
	fmt.Println("| without the hassle that it would normally entail.         |")
	fmt.Println("|                                                           |")
	fmt.Println("| Puppeth uses SSH to dial in to remote servers, and builds |")
	fmt.Println("| its network components out of Docker containers using the |")
	fmt.Println("| docker-compose toolset.                                   |")
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println()

//确保我们有一个好的网络名称来使用fmt.println（）。
//Docker接受图像名称中的连字符，但不喜欢容器名称中的连字符
	if w.network == "" {
		fmt.Println("Please specify a network name to administer (no spaces, hyphens or capital letters please)")
		for {
			w.network = w.readString()
			if !strings.Contains(w.network, " ") && !strings.Contains(w.network, "-") && strings.ToLower(w.network) == w.network {
				fmt.Printf("\nSweet, you can set this via --network=%s next time!\n\n", w.network)
				break
			}
			log.Error("I also like to live dangerously, still no spaces, hyphens or capital letters")
		}
	}
	log.Info("Administering Ethereum network", "name", w.network)

//加载初始配置并连接到所有活动服务器
	w.conf.path = filepath.Join(os.Getenv("HOME"), ".puppeth", w.network)

	blob, err := ioutil.ReadFile(w.conf.path)
	if err != nil {
		log.Warn("No previous configurations found", "path", w.conf.path)
	} else if err := json.Unmarshal(blob, &w.conf); err != nil {
		log.Crit("Previous configuration corrupted", "path", w.conf.path, "err", err)
	} else {
//同时拨打所有以前已知的服务器
		var pend sync.WaitGroup
		for server, pubkey := range w.conf.Servers {
			pend.Add(1)

			go func(server string, pubkey []byte) {
				defer pend.Done()

				log.Info("Dialing previously configured server", "server", server)
				client, err := dial(server, pubkey)
				if err != nil {
					log.Error("Previous server unreachable", "server", server, "err", err)
				}
				w.lock.Lock()
				w.servers[server] = client
				w.lock.Unlock()
			}(server, pubkey)
		}
		pend.Wait()
		w.networkStats()
	}
//基本完成了，无限循环关于该做什么
	for {
		fmt.Println()
		fmt.Println("What would you like to do? (default = stats)")
		fmt.Println(" 1. Show network stats")
		if w.conf.Genesis == nil {
			fmt.Println(" 2. Configure new genesis")
		} else {
			fmt.Println(" 2. Manage existing genesis")
		}
		if len(w.servers) == 0 {
			fmt.Println(" 3. Track new remote server")
		} else {
			fmt.Println(" 3. Manage tracked machines")
		}
		if len(w.services) == 0 {
			fmt.Println(" 4. Deploy network components")
		} else {
			fmt.Println(" 4. Manage network components")
		}

		choice := w.read()
		switch {
		case choice == "" || choice == "1":
			w.networkStats()

		case choice == "2":
			if w.conf.Genesis == nil {
				fmt.Println()
				fmt.Println("What would you like to do? (default = create)")
				fmt.Println(" 1. Create new genesis from scratch")
				fmt.Println(" 2. Import already existing genesis")

				choice := w.read()
				switch {
				case choice == "" || choice == "1":
					w.makeGenesis()
				case choice == "2":
					w.importGenesis()
				default:
					log.Error("That's not something I can do")
				}
			} else {
				w.manageGenesis()
			}
		case choice == "3":
			if len(w.servers) == 0 {
				if w.makeServer() != "" {
					w.networkStats()
				}
			} else {
				w.manageServers()
			}
		case choice == "4":
			if len(w.services) == 0 {
				w.deployComponent()
			} else {
				w.manageComponents()
			}
		default:
			log.Error("That's not something I can do")
		}
	}
}

