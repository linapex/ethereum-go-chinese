
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450087858212864>


package downloader

import (
	"context"
	"sync"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
)

//PublicDownloaderAPI提供了一个API，它提供有关当前同步状态的信息。
//它只提供对任何人都可以使用的数据进行操作的方法，而不存在安全风险。
type PublicDownloaderAPI struct {
	d                         *Downloader
	mux                       *event.TypeMux
	installSyncSubscription   chan chan interface{}
	uninstallSyncSubscription chan *uninstallSyncSubscriptionRequest
}

//新建PublicDownloaderAPI创建新的PublicDownloaderAPI。API有一个内部事件循环，
//通过全局事件mux从下载程序侦听事件。如果它收到
//这些事件会将其广播到通过
//InstallSyncSubscription频道。
func NewPublicDownloaderAPI(d *Downloader, m *event.TypeMux) *PublicDownloaderAPI {
	api := &PublicDownloaderAPI{
		d:                         d,
		mux:                       m,
		installSyncSubscription:   make(chan chan interface{}),
		uninstallSyncSubscription: make(chan *uninstallSyncSubscriptionRequest),
	}

	go api.eventLoop()

	return api
}

//EventLoop运行一个循环，直到事件mux关闭。它将安装和卸载新的
//将订阅和广播同步状态更新同步到已安装的同步订阅。
func (api *PublicDownloaderAPI) eventLoop() {
	var (
		sub               = api.mux.Subscribe(StartEvent{}, DoneEvent{}, FailedEvent{})
		syncSubscriptions = make(map[chan interface{}]struct{})
	)

	for {
		select {
		case i := <-api.installSyncSubscription:
			syncSubscriptions[i] = struct{}{}
		case u := <-api.uninstallSyncSubscription:
			delete(syncSubscriptions, u.c)
			close(u.uninstalled)
		case event := <-sub.Chan():
			if event == nil {
				return
			}

			var notification interface{}
			switch event.Data.(type) {
			case StartEvent:
				notification = &SyncingResult{
					Syncing: true,
					Status:  api.d.Progress(),
				}
			case DoneEvent, FailedEvent:
				notification = false
			}
//广播
			for c := range syncSubscriptions {
				c <- notification
			}
		}
	}
}

//同步提供此节点何时开始与以太坊网络同步以及何时完成同步的信息。
func (api *PublicDownloaderAPI) Syncing(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		statuses := make(chan interface{})
		sub := api.SubscribeSyncStatus(statuses)

		for {
			select {
			case status := <-statuses:
				notifier.Notify(rpcSub.ID, status)
			case <-rpcSub.Err():
				sub.Unsubscribe()
				return
			case <-notifier.Closed():
				sub.Unsubscribe()
				return
			}
		}
	}()

	return rpcSub, nil
}

//同步结果提供有关此节点当前同步状态的信息。
type SyncingResult struct {
	Syncing bool                  `json:"syncing"`
	Status  ethereum.SyncProgress `json:"status"`
}

//UninstallSyncSubscriptionRequest在API事件循环中卸载同步订阅。
type uninstallSyncSubscriptionRequest struct {
	c           chan interface{}
	uninstalled chan interface{}
}

//SyncStatusSubscription表示同步订阅。
type SyncStatusSubscription struct {
api       *PublicDownloaderAPI //register subscription in event loop of this api instance
c         chan interface{}     //事件广播到的频道
unsubOnce sync.Once            //确保取消订阅逻辑执行一次
}

//取消订阅将从DeloLoad事件循环中卸载订阅。
//传递给subscribeSyncStatus的状态通道不再使用。
//在这个方法返回之后。
func (s *SyncStatusSubscription) Unsubscribe() {
	s.unsubOnce.Do(func() {
		req := uninstallSyncSubscriptionRequest{s.c, make(chan interface{})}
		s.api.uninstallSyncSubscription <- &req

		for {
			select {
			case <-s.c:
//drop new status events until uninstall confirmation
				continue
			case <-req.uninstalled:
				return
			}
		}
	})
}

//订阅同步状态创建将广播新同步更新的订阅。
//给定的通道必须接收接口值，结果可以是
func (api *PublicDownloaderAPI) SubscribeSyncStatus(status chan interface{}) *SyncStatusSubscription {
	api.installSyncSubscription <- status
	return &SyncStatusSubscription{api: api, c: status}
}

