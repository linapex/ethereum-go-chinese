
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450113137283072>


//+构建Linux Darwin Freebsd

package fuse

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/api"
	"github.com/ethereum/go-ethereum/swarm/log"
)

var (
	errEmptyMountPoint      = errors.New("need non-empty mount point")
	errNoRelativeMountPoint = errors.New("invalid path for mount point (need absolute path)")
	errMaxMountCount        = errors.New("max FUSE mount count reached")
	errMountTimeout         = errors.New("mount timeout")
	errAlreadyMounted       = errors.New("mount point is already serving")
)

func isFUSEUnsupportedError(err error) bool {
	if perr, ok := err.(*os.PathError); ok {
		return perr.Op == "open" && perr.Path == "/dev/fuse"
	}
	return err == fuse.ErrOSXFUSENotFound
}

//mountinfo包含有关每个活动装载的信息
type MountInfo struct {
	MountPoint     string
	StartManifest  string
	LatestManifest string
	rootDir        *SwarmDir
	fuseConnection *fuse.Conn
	swarmApi       *api.API
	lock           *sync.RWMutex
	serveClose     chan struct{}
}

func NewMountInfo(mhash, mpoint string, sapi *api.API) *MountInfo {
	log.Debug("swarmfs NewMountInfo", "hash", mhash, "mount point", mpoint)
	newMountInfo := &MountInfo{
		MountPoint:     mpoint,
		StartManifest:  mhash,
		LatestManifest: mhash,
		rootDir:        nil,
		fuseConnection: nil,
		swarmApi:       sapi,
		lock:           &sync.RWMutex{},
		serveClose:     make(chan struct{}),
	}
	return newMountInfo
}

func (swarmfs *SwarmFS) Mount(mhash, mountpoint string) (*MountInfo, error) {
	log.Info("swarmfs", "mounting hash", mhash, "mount point", mountpoint)
	if mountpoint == "" {
		return nil, errEmptyMountPoint
	}
	if !strings.HasPrefix(mountpoint, "/") {
		return nil, errNoRelativeMountPoint
	}
	cleanedMountPoint, err := filepath.Abs(filepath.Clean(mountpoint))
	if err != nil {
		return nil, err
	}
	log.Trace("swarmfs mount", "cleanedMountPoint", cleanedMountPoint)

	swarmfs.swarmFsLock.Lock()
	defer swarmfs.swarmFsLock.Unlock()

	noOfActiveMounts := len(swarmfs.activeMounts)
	log.Debug("swarmfs mount", "# active mounts", noOfActiveMounts)
	if noOfActiveMounts >= maxFuseMounts {
		return nil, errMaxMountCount
	}

	if _, ok := swarmfs.activeMounts[cleanedMountPoint]; ok {
		return nil, errAlreadyMounted
	}

	log.Trace("swarmfs mount: getting manifest tree")
	_, manifestEntryMap, err := swarmfs.swarmApi.BuildDirectoryTree(context.TODO(), mhash, true)
	if err != nil {
		return nil, err
	}

	log.Trace("swarmfs mount: building mount info")
	mi := NewMountInfo(mhash, cleanedMountPoint, swarmfs.swarmApi)

	dirTree := map[string]*SwarmDir{}
	rootDir := NewSwarmDir("/", mi)
	log.Trace("swarmfs mount", "rootDir", rootDir)
	mi.rootDir = rootDir

	log.Trace("swarmfs mount: traversing manifest map")
	for suffix, entry := range manifestEntryMap {
if suffix == "" { //空后缀表示文件没有名称，即这是清单中的默认条目。因为我们不能有没有名字的文件，所以让我们忽略这个条目。
			log.Warn("Manifest has an empty-path (default) entry which will be ignored in FUSE mount.")
			continue
		}
		addr := common.Hex2Bytes(entry.Hash)
		fullpath := "/" + suffix
		basepath := filepath.Dir(fullpath)
		parentDir := rootDir
		dirUntilNow := ""
		paths := strings.Split(basepath, "/")
		for i := range paths {
			if paths[i] != "" {
				thisDir := paths[i]
				dirUntilNow = dirUntilNow + "/" + thisDir

				if _, ok := dirTree[dirUntilNow]; !ok {
					dirTree[dirUntilNow] = NewSwarmDir(dirUntilNow, mi)
					parentDir.directories = append(parentDir.directories, dirTree[dirUntilNow])
					parentDir = dirTree[dirUntilNow]

				} else {
					parentDir = dirTree[dirUntilNow]
				}
			}
		}
		thisFile := NewSwarmFile(basepath, filepath.Base(fullpath), mi)
		thisFile.addr = addr

		parentDir.files = append(parentDir.files, thisFile)
	}

	fconn, err := fuse.Mount(cleanedMountPoint, fuse.FSName("swarmfs"), fuse.VolumeName(mhash))
	if isFUSEUnsupportedError(err) {
		log.Error("swarmfs error - FUSE not installed", "mountpoint", cleanedMountPoint, "err", err)
		return nil, err
	} else if err != nil {
		fuse.Unmount(cleanedMountPoint)
		log.Error("swarmfs error mounting swarm manifest", "mountpoint", cleanedMountPoint, "err", err)
		return nil, err
	}
	mi.fuseConnection = fconn

	serverr := make(chan error, 1)
	go func() {
		log.Info("swarmfs", "serving hash", mhash, "at", cleanedMountPoint)
		filesys := &SwarmRoot{root: rootDir}
//开始为实际的文件系统提供服务；请参阅下面的注释
		if err := fs.Serve(fconn, filesys); err != nil {
			log.Warn("swarmfs could not serve the requested hash", "error", err)
			serverr <- err
		}
		mi.serveClose <- struct{}{}
	}()

 /*
    重要提示：fs.serve功能为阻塞；
    通过调用
    ATTR在每个swarmfile上运行，创建文件inode；
    专门调用swarm的lazysactionreader.size（）来设置文件大小。

    这可能需要一些时间，而且如果我们访问fuse文件系统
    太早了，我们可以使测试死锁。目前的假设是
    此时，fuse驱动程序没有完成初始化文件系统。

    过早地访问文件不仅会使测试死锁，而且会锁定访问
    完全的fuse文件，导致操作系统级别的资源被阻塞。
    即使是一个简单的'ls/tmp/testdir/testmountdir'也可能在shell中死锁。

    目前的解决方法是等待一段时间，给操作系统足够的时间来初始化
    Fuse文件系统。在测试过程中，这似乎解决了这个问题。

    但应注意的是，这可能只是一种影响，
    以及由于某些竞争条件而阻塞访问的其他原因导致的死锁
    （在bazil.org库和/或swarmroot、swarmdir和swarmfile实现中导致）
 **/

	time.Sleep(2 * time.Second)

	timer := time.NewTimer(mountTimeout)
	defer timer.Stop()
//检查装载过程是否有要报告的错误。
	select {
	case <-timer.C:
		log.Warn("swarmfs timed out mounting over FUSE", "mountpoint", cleanedMountPoint, "err", err)
		err := fuse.Unmount(cleanedMountPoint)
		if err != nil {
			return nil, err
		}
		return nil, errMountTimeout
	case err := <-serverr:
		log.Warn("swarmfs error serving over FUSE", "mountpoint", cleanedMountPoint, "err", err)
		err = fuse.Unmount(cleanedMountPoint)
		return nil, err

	case <-fconn.Ready:
//这表示来自保险丝的实际安装点。安装调用已就绪；
//虽然fs.serve中的文件系统实际上是完全建立起来的，但它并不表示
		if err := fconn.MountError; err != nil {
			log.Error("Mounting error from fuse driver: ", "err", err)
			return nil, err
		}
		log.Info("swarmfs now served over FUSE", "manifest", mhash, "mountpoint", cleanedMountPoint)
	}

	timer.Stop()
	swarmfs.activeMounts[cleanedMountPoint] = mi
	return mi, nil
}

func (swarmfs *SwarmFS) Unmount(mountpoint string) (*MountInfo, error) {
	swarmfs.swarmFsLock.Lock()
	defer swarmfs.swarmFsLock.Unlock()

	cleanedMountPoint, err := filepath.Abs(filepath.Clean(mountpoint))
	if err != nil {
		return nil, err
	}

	mountInfo := swarmfs.activeMounts[cleanedMountPoint]

	if mountInfo == nil || mountInfo.MountPoint != cleanedMountPoint {
		return nil, fmt.Errorf("swarmfs %s is not mounted", cleanedMountPoint)
	}
	err = fuse.Unmount(cleanedMountPoint)
	if err != nil {
		err1 := externalUnmount(cleanedMountPoint)
		if err1 != nil {
			errStr := fmt.Sprintf("swarmfs unmount error: %v", err)
			log.Warn(errStr)
			return nil, err1
		}
	}

	err = mountInfo.fuseConnection.Close()
	if err != nil {
		return nil, err
	}
	delete(swarmfs.activeMounts, cleanedMountPoint)

	<-mountInfo.serveClose

	succString := fmt.Sprintf("swarmfs unmounting %v succeeded", cleanedMountPoint)
	log.Info(succString)

	return mountInfo, nil
}

func (swarmfs *SwarmFS) Listmounts() []*MountInfo {
	swarmfs.swarmFsLock.RLock()
	defer swarmfs.swarmFsLock.RUnlock()
	rows := make([]*MountInfo, 0, len(swarmfs.activeMounts))
	for _, mi := range swarmfs.activeMounts {
		rows = append(rows, mi)
	}
	return rows
}

func (swarmfs *SwarmFS) Stop() bool {
	for mp := range swarmfs.activeMounts {
		mountInfo := swarmfs.activeMounts[mp]
		swarmfs.Unmount(mountInfo.MountPoint)
	}
	return true
}

