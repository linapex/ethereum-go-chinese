
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450112919179264>


package fuse

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/swarm/api"
)

const (
	Swarmfs_Version = "0.1"
	mountTimeout    = time.Second * 5
	unmountTimeout  = time.Second * 10
	maxFuseMounts   = 5
)

var (
swarmfs     *SwarmFS //swarm文件系统singleton
	swarmfsLock sync.Once

inode     uint64 = 1 //全局索引
	inodeLock sync.RWMutex
)

type SwarmFS struct {
	swarmApi     *api.API
	activeMounts map[string]*MountInfo
	swarmFsLock  *sync.RWMutex
}

func NewSwarmFS(api *api.API) *SwarmFS {
	swarmfsLock.Do(func() {
		swarmfs = &SwarmFS{
			swarmApi:     api,
			swarmFsLock:  &sync.RWMutex{},
			activeMounts: map[string]*MountInfo{},
		}
	})
	return swarmfs

}

//inode编号必须是唯一的，它们用于在fuse中缓存
func NewInode() uint64 {
	inodeLock.Lock()
	defer inodeLock.Unlock()
	inode += 1
	return inode
}

