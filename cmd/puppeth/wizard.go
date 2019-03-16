
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450069743013888>


package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/crypto/ssh/terminal"
)

//config包含Puppeth需要保存的所有配置
//
type config struct {
path      string   //
bootnodes []string //所有节点始终连接到的引导节点
ethstats  string   //要为节点部署缓存的ethstats设置

Genesis *core.Genesis     `json:"genesis,omitempty"` //用于节点部署的缓存Genesis块
	Servers map[string][]byte `json:"servers,omitempty"`
}

//服务器检索按字母顺序排序的服务器列表。
func (c config) servers() []string {
	servers := make([]string, 0, len(c.Servers))
	for server := range c.Servers {
		servers = append(servers, server)
	}
	sort.Strings(servers)

	return servers
}

//flush将配置的内容转储到磁盘。
func (c config) flush() {
	os.MkdirAll(filepath.Dir(c.path), 0755)

	out, _ := json.MarshalIndent(c, "", "  ")
	if err := ioutil.WriteFile(c.path, out, 0644); err != nil {
		log.Warn("Failed to save puppeth configs", "file", c.path, "err", err)
	}
}

type wizard struct {
network string //要管理的网络名称
conf    config //以前运行的配置

servers  map[string]*sshClient //ssh连接到要管理的服务器
services map[string][]string   //已知正在服务器上运行的以太坊服务

in   *bufio.Reader //包装stdin以允许读取用户输入
lock sync.Mutex    //锁定以在并发服务发现期间保护配置
}

//读取从stdin中读取一行，如果从空格中删除，则进行剪裁。
func (w *wizard) read() string {
	fmt.Printf("> ")
	text, err := w.in.ReadString('\n')
	if err != nil {
		log.Crit("Failed to read user input", "err", err)
	}
	return strings.TrimSpace(text)
}

//
//非空性。
func (w *wizard) readString() string {
	for {
		fmt.Printf("> ")
		text, err := w.in.ReadString('\n')
		if err != nil {
			log.Crit("Failed to read user input", "err", err)
		}
		if text = strings.TrimSpace(text); text != "" {
			return text
		}
	}
}

//readdefaultstring从stdin读取一行，从空格中剪裁if。如果
//输入空行，返回默认值。
func (w *wizard) readDefaultString(def string) string {
	fmt.Printf("> ")
	text, err := w.in.ReadString('\n')
	if err != nil {
		log.Crit("Failed to read user input", "err", err)
	}
	if text = strings.TrimSpace(text); text != "" {
		return text
	}
	return def
}

//
//
//
func (w *wizard) readDefaultYesNo(def bool) bool {
	for {
		fmt.Printf("> ")
		text, err := w.in.ReadString('\n')
		if err != nil {
			log.Crit("Failed to read user input", "err", err)
		}
		if text = strings.ToLower(strings.TrimSpace(text)); text == "" {
			return def
		}
		if text == "y" || text == "yes" {
			return true
		}
		if text == "n" || text == "no" {
			return false
		}
		log.Error("Invalid input, expected 'y', 'yes', 'n', 'no' or empty")
	}
}

//
//
func (w *wizard) readURL() *url.URL {
	for {
		fmt.Printf("> ")
		text, err := w.in.ReadString('\n')
		if err != nil {
			log.Crit("Failed to read user input", "err", err)
		}
		uri, err := url.Parse(strings.TrimSpace(text))
		if err != nil {
			log.Error("Invalid input, expected URL", "err", err)
			continue
		}
		return uri
	}
}

//readint从stdin读取一行，从空格中剪裁if，强制执行
//
func (w *wizard) readInt() int {
	for {
		fmt.Printf("> ")
		text, err := w.in.ReadString('\n')
		if err != nil {
			log.Crit("Failed to read user input", "err", err)
		}
		if text = strings.TrimSpace(text); text == "" {
			continue
		}
		val, err := strconv.Atoi(strings.TrimSpace(text))
		if err != nil {
			log.Error("Invalid input, expected integer", "err", err)
			continue
		}
		return val
	}
}

//
//
//返回。
func (w *wizard) readDefaultInt(def int) int {
	for {
		fmt.Printf("> ")
		text, err := w.in.ReadString('\n')
		if err != nil {
			log.Crit("Failed to read user input", "err", err)
		}
		if text = strings.TrimSpace(text); text == "" {
			return def
		}
		val, err := strconv.Atoi(strings.TrimSpace(text))
		if err != nil {
			log.Error("Invalid input, expected integer", "err", err)
			continue
		}
		return val
	}
}

//
//
//
func (w *wizard) readDefaultBigInt(def *big.Int) *big.Int {
	for {
		fmt.Printf("> ")
		text, err := w.in.ReadString('\n')
		if err != nil {
			log.Crit("Failed to read user input", "err", err)
		}
		if text = strings.TrimSpace(text); text == "" {
			return def
		}
		val, ok := new(big.Int).SetString(text, 0)
		if !ok {
			log.Error("Invalid input, expected big integer")
			continue
		}
		return val
	}
}

/*



 
  
  
  如果犯错！= nIL{
   
  
  如果text=strings.trimspace（text）；text=“”
   持续
  }
  val，err：=strconv.parsefloat（strings.trimspace（text），64）
  如果犯错！= nIL{
   log.error（“输入无效，应为float”，“err”，err）
   持续
  
  返回瓦尔
 

*/


//readdefaultfloat从stdin读取一行，从空格中剪裁if，强制
//
func (w *wizard) readDefaultFloat(def float64) float64 {
	for {
		fmt.Printf("> ")
		text, err := w.in.ReadString('\n')
		if err != nil {
			log.Crit("Failed to read user input", "err", err)
		}
		if text = strings.TrimSpace(text); text == "" {
			return def
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
		if err != nil {
			log.Error("Invalid input, expected float", "err", err)
			continue
		}
		return val
	}
}

//readpassword从stdin读取一行，从尾随的new
//
func (w *wizard) readPassword() string {
	fmt.Printf("> ")
	text, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Crit("Failed to read password", "err", err)
	}
	fmt.Println()
	return string(text)
}

//
//发送到以太坊地址。
func (w *wizard) readAddress() *common.Address {
	for {
//从用户处读取地址
		fmt.Printf("> 0x")
		text, err := w.in.ReadString('\n')
		if err != nil {
			log.Crit("Failed to read user input", "err", err)
		}
		if text = strings.TrimSpace(text); text == "" {
			return nil
		}
//
		if len(text) != 40 {
			log.Error("Invalid address length, please retry")
			continue
		}
		bigaddr, _ := new(big.Int).SetString(text, 16)
		address := common.BigToAddress(bigaddr)
		return &address
	}
}

//
//将其转换为以太坊地址。如果输入空行，则默认
//
func (w *wizard) readDefaultAddress(def common.Address) common.Address {
	for {
//从用户处读取地址
		fmt.Printf("> 0x")
		text, err := w.in.ReadString('\n')
		if err != nil {
			log.Crit("Failed to read user input", "err", err)
		}
		if text = strings.TrimSpace(text); text == "" {
			return def
		}
//
		if len(text) != 40 {
			log.Error("Invalid address length, please retry")
			continue
		}
		bigaddr, _ := new(big.Int).SetString(text, 16)
		return common.BigToAddress(bigaddr)
	}
}

//
func (w *wizard) readJSON() string {
	var blob json.RawMessage

	for {
		fmt.Printf("> ")
		if err := json.NewDecoder(w.in).Decode(&blob); err != nil {
			log.Error("Invalid JSON, please try again", "err", err)
			continue
		}
		return string(blob)
	}
}

//
//
//用户输入格式（而不是返回go net.ip）与
//
func (w *wizard) readIPAddress() string {
	for {
//从用户读取IP地址
		fmt.Printf("> ")
		text, err := w.in.ReadString('\n')
		if err != nil {
			log.Crit("Failed to read user input", "err", err)
		}
		if text = strings.TrimSpace(text); text == "" {
			return ""
		}
//
		if ip := net.ParseIP(text); ip == nil {
			log.Error("Invalid IP address, please retry")
			continue
		}
		return text
	}
}

