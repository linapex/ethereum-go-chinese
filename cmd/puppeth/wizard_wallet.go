
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450070250524672>


package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
)

//
func (w *wizard) deployWallet() {
//在用户浪费输入时间之前做一些健全性检查
	if w.conf.Genesis == nil {
		log.Error("No genesis block configured")
		return
	}
	if w.conf.ethstats == "" {
		log.Error("No ethstats server configured")
		return
	}
//选择要与之交互的服务器
	server := w.selectServer()
	if server == "" {
		return
	}
	client := w.servers[server]

//从服务器检索任何活动节点配置
	infos, err := checkWallet(client, w.network)
	if err != nil {
		infos = &walletInfos{
			nodePort: 30303, rpcPort: 8545, webPort: 80, webHost: client.server,
		}
	}
	existed := err == nil

	infos.genesis, _ = json.MarshalIndent(w.conf.Genesis, "", "  ")
	infos.network = w.conf.Genesis.Config.ChainID.Int64()

//找出要监听的端口
	fmt.Println()
	fmt.Printf("Which port should the wallet listen on? (default = %d)\n", infos.webPort)
	infos.webPort = w.readDefaultInt(infos.webPort)

//图1部署ethstats的虚拟主机
	if infos.webHost, err = w.ensureVirtualHost(client, infos.webPort, infos.webHost); err != nil {
		log.Error("Failed to decide on wallet host", "err", err)
		return
	}
//找出用户希望存储持久数据的位置
	fmt.Println()
	if infos.datadir == "" {
		fmt.Printf("Where should data be stored on the remote machine?\n")
		infos.datadir = w.readString()
	} else {
		fmt.Printf("Where should data be stored on the remote machine? (default = %s)\n", infos.datadir)
		infos.datadir = w.readDefaultString(infos.datadir)
	}
//找出要监听的端口
	fmt.Println()
	fmt.Printf("Which TCP/UDP port should the backing node listen on? (default = %d)\n", infos.nodePort)
	infos.nodePort = w.readDefaultInt(infos.nodePort)

	fmt.Println()
	fmt.Printf("Which port should the backing RPC API listen on? (default = %d)\n", infos.rpcPort)
	infos.rpcPort = w.readDefaultInt(infos.rpcPort)

//在统计页面上设置要报告的正确名称
	fmt.Println()
	if infos.ethstats == "" {
		fmt.Printf("What should the wallet be called on the stats page?\n")
		infos.ethstats = w.readString() + ":" + w.conf.ethstats
	} else {
		fmt.Printf("What should the wallet be called on the stats page? (default = %s)\n", infos.ethstats)
		infos.ethstats = w.readDefaultString(infos.ethstats) + ":" + w.conf.ethstats
	}
//尝试在主机上部署钱包
	nocache := false
	if existed {
		fmt.Println()
		fmt.Printf("Should the wallet be built from scratch (y/n)? (default = no)\n")
		nocache = w.readDefaultYesNo(false)
	}
	if out, err := deployWallet(client, w.network, w.conf.bootnodes, infos, nocache); err != nil {
		log.Error("Failed to deploy wallet container", "err", err)
		if len(out) > 0 {
			fmt.Printf("%s\n", out)
		}
		return
	}
//一切正常，运行网络扫描以获取任何更改
	log.Info("Waiting for node to finish booting")
	time.Sleep(3 * time.Second)

	w.networkStats()
}

