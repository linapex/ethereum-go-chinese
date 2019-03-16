
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:32</date>
//</624450065276080128>


package accounts

import (
	"reflect"
	"sort"
	"sync"

	"github.com/ethereum/go-ethereum/event"
)

//经理是一个主要的客户经理，可以与
//用于签署事务的后端。
type Manager struct {
backends map[reflect.Type][]Backend //当前注册的后端索引
updaters []event.Subscription       //所有后端的钱包更新订阅
updates  chan WalletEvent           //后端钱包更改的订阅接收器
wallets  []Wallet                   //缓存所有注册后端的所有钱包

feed event.Feed //钱包信息提示到达/离开

	quit chan chan error
	lock sync.RWMutex
}

//NewManager创建一个通用的客户经理，通过
//支撑的后端。
func NewManager(backends ...Backend) *Manager {
//从后端检索钱包的初始列表并按URL排序
	var wallets []Wallet
	for _, backend := range backends {
		wallets = merge(wallets, backend.Wallets()...)
	}
//从所有后端订阅钱包通知
	updates := make(chan WalletEvent, 4*len(backends))

	subs := make([]event.Subscription, len(backends))
	for i, backend := range backends {
		subs[i] = backend.Subscribe(updates)
	}
//召集客户经理并返回
	am := &Manager{
		backends: make(map[reflect.Type][]Backend),
		updaters: subs,
		updates:  updates,
		wallets:  wallets,
		quit:     make(chan chan error),
	}
	for _, backend := range backends {
		kind := reflect.TypeOf(backend)
		am.backends[kind] = append(am.backends[kind], backend)
	}
	go am.update()

	return am
}

//关闭将终止客户经理的内部通知进程。
func (am *Manager) Close() error {
	errc := make(chan error)
	am.quit <- errc
	return <-errc
}

//更新是钱包事件循环，监听来自后端的通知
//更新钱包的缓存。
func (am *Manager) update() {
//管理器终止时关闭所有订阅
	defer func() {
		am.lock.Lock()
		for _, sub := range am.updaters {
			sub.Unsubscribe()
		}
		am.updaters = nil
		am.lock.Unlock()
	}()

//循环直至终止
	for {
		select {
		case event := <-am.updates:
//
			am.lock.Lock()
			switch event.Kind {
			case WalletArrived:
				am.wallets = merge(am.wallets, event.Wallet)
			case WalletDropped:
				am.wallets = drop(am.wallets, event.Wallet)
			}
			am.lock.Unlock()

//通知事件的任何侦听器
			am.feed.Send(event)

		case errc := <-am.quit:
//经理终止，返回
			errc <- nil
			return
		}
	}
}

//后端从帐户管理器中检索具有给定类型的后端。
func (am *Manager) Backends(kind reflect.Type) []Backend {
	return am.backends[kind]
}

//钱包将返回在此帐户管理器下注册的所有签名者帐户。
func (am *Manager) Wallets() []Wallet {
	am.lock.RLock()
	defer am.lock.RUnlock()

	cpy := make([]Wallet, len(am.wallets))
	copy(cpy, am.wallets)
	return cpy
}

//Wallet检索与特定URL关联的钱包。
func (am *Manager) Wallet(url string) (Wallet, error) {
	am.lock.RLock()
	defer am.lock.RUnlock()

	parsed, err := parseURL(url)
	if err != nil {
		return nil, err
	}
	for _, wallet := range am.Wallets() {
		if wallet.URL() == parsed {
			return wallet, nil
		}
	}
	return nil, ErrUnknownWallet
}

//查找与特定帐户对应的钱包。自从
//帐户可以动态地添加到钱包中或从钱包中删除，此方法具有
//钱包数量的线性运行时间。
func (am *Manager) Find(account Account) (Wallet, error) {
	am.lock.RLock()
	defer am.lock.RUnlock()

	for _, wallet := range am.wallets {
		if wallet.Contains(account) {
			return wallet, nil
		}
	}
	return nil, ErrUnknownAccount
}

//订阅创建异步订阅以在
//经理检测到钱包从其任何后端到达或离开。
func (am *Manager) Subscribe(sink chan<- WalletEvent) event.Subscription {
	return am.feed.Subscribe(sink)
}

//merge是一种类似于append的钱包排序方式，其中
//通过在正确位置插入新钱包，可以保留原始列表。
//
//假定原始切片已按URL排序。
func merge(slice []Wallet, wallets ...Wallet) []Wallet {
	for _, wallet := range wallets {
		n := sort.Search(len(slice), func(i int) bool { return slice[i].URL().Cmp(wallet.URL()) >= 0 })
		if n == len(slice) {
			slice = append(slice, wallet)
			continue
		}
		slice = append(slice[:n], append([]Wallet{wallet}, slice[n:]...)...)
	}
	return slice
}

//drop是merge的coutterpart，它从排序后的
//缓存并删除指定的。
func drop(slice []Wallet, wallets ...Wallet) []Wallet {
	for _, wallet := range wallets {
		n := sort.Search(len(slice), func(i int) bool { return slice[i].URL().Cmp(wallet.URL()) >= 0 })
		if n == len(slice) {
//未找到钱包，可能在启动过程中发生
			continue
		}
		slice = append(slice[:n], slice[n+1:]...)
	}
	return slice
}

