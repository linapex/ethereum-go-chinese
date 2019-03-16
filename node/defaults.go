
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450102219509760>


package node

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/nat"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
DefaultHTTPHost = "localhost" //HTTP RPC服务器的默认主机接口
DefaultHTTPPort = 8545        //HTTP RPC服务器的默认TCP端口
DefaultWSHost   = "localhost" //WebSocket RPC服务器的默认主机接口
DefaultWSPort   = 8546        //WebSocket RPC服务器的默认TCP端口
)

//默认配置包含合理的默认设置。
var DefaultConfig = Config{
	DataDir:          DefaultDataDir(),
	HTTPPort:         DefaultHTTPPort,
	HTTPModules:      []string{"net", "web3"},
	HTTPVirtualHosts: []string{"localhost"},
	HTTPTimeouts:     rpc.DefaultHTTPTimeouts,
	WSPort:           DefaultWSPort,
	WSModules:        []string{"net", "web3"},
	P2P: p2p.Config{
		ListenAddr: ":30303",
		MaxPeers:   25,
		NAT:        nat.Any(),
	},
}

//defaultdatadir是用于数据库和其他数据库的默认数据目录
//持久性要求。
func DefaultDataDir() string {
//尝试将数据文件夹放在用户的home目录中
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Ethereum")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Ethereum")
		} else {
			return filepath.Join(home, ".ethereum")
		}
	}
//因为我们无法猜测一个稳定的位置，所以返回空的，稍后再处理
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

