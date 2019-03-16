
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450113196003328>


//+构建Linux Darwin Freebsd

package fuse

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/ethereum/go-ethereum/swarm/log"
)

func externalUnmount(mountPoint string) error {
	ctx, cancel := context.WithTimeout(context.Background(), unmountTimeout)
	defer cancel()

//尝试通用umount。
	if err := exec.CommandContext(ctx, "umount", mountPoint).Run(); err == nil {
		return nil
	}
//如果umount不工作，请尝试使用fuse-specific命令。
	switch runtime.GOOS {
	case "darwin":
		return exec.CommandContext(ctx, "diskutil", "umount", mountPoint).Run()
	case "linux":
		return exec.CommandContext(ctx, "fusermount", "-u", mountPoint).Run()
	default:
		return fmt.Errorf("swarmfs unmount: unimplemented")
	}
}

func addFileToSwarm(sf *SwarmFile, content []byte, size int) error {
	fkey, mhash, err := sf.mountInfo.swarmApi.AddFile(context.TODO(), sf.mountInfo.LatestManifest, sf.path, sf.name, content, true)
	if err != nil {
		return err
	}

	sf.lock.Lock()
	defer sf.lock.Unlock()
	sf.addr = fkey
	sf.fileSize = int64(size)

	sf.mountInfo.lock.Lock()
	defer sf.mountInfo.lock.Unlock()
	sf.mountInfo.LatestManifest = mhash

	log.Info("swarmfs added new file:", "fname", sf.name, "new Manifest hash", mhash)
	return nil
}

func removeFileFromSwarm(sf *SwarmFile) error {
	mkey, err := sf.mountInfo.swarmApi.RemoveFile(context.TODO(), sf.mountInfo.LatestManifest, sf.path, sf.name, true)
	if err != nil {
		return err
	}

	sf.mountInfo.lock.Lock()
	defer sf.mountInfo.lock.Unlock()
	sf.mountInfo.LatestManifest = mkey

	log.Info("swarmfs removed file:", "fname", sf.name, "new Manifest hash", mkey)
	return nil
}

func removeDirectoryFromSwarm(sd *SwarmDir) error {
	if len(sd.directories) == 0 && len(sd.files) == 0 {
		return nil
	}

	for _, d := range sd.directories {
		err := removeDirectoryFromSwarm(d)
		if err != nil {
			return err
		}
	}

	for _, f := range sd.files {
		err := removeFileFromSwarm(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func appendToExistingFileInSwarm(sf *SwarmFile, content []byte, offset int64, length int64) error {
	fkey, mhash, err := sf.mountInfo.swarmApi.AppendFile(context.TODO(), sf.mountInfo.LatestManifest, sf.path, sf.name, sf.fileSize, content, sf.addr, offset, length, true)
	if err != nil {
		return err
	}

	sf.lock.Lock()
	defer sf.lock.Unlock()
	sf.addr = fkey
	sf.fileSize = sf.fileSize + int64(len(content))

	sf.mountInfo.lock.Lock()
	defer sf.mountInfo.lock.Unlock()
	sf.mountInfo.LatestManifest = mhash

	log.Info("swarmfs appended file:", "fname", sf.name, "new Manifest hash", mhash)
	return nil
}

