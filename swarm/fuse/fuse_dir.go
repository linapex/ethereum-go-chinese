
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450112789155842>


//+构建Linux Darwin Freebsd

package fuse

import (
	"os"
	"path/filepath"
	"sync"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/ethereum/go-ethereum/swarm/log"
	"golang.org/x/net/context"
)

var (
	_ fs.Node                = (*SwarmDir)(nil)
	_ fs.NodeRequestLookuper = (*SwarmDir)(nil)
	_ fs.HandleReadDirAller  = (*SwarmDir)(nil)
	_ fs.NodeCreater         = (*SwarmDir)(nil)
	_ fs.NodeRemover         = (*SwarmDir)(nil)
	_ fs.NodeMkdirer         = (*SwarmDir)(nil)
)

type SwarmDir struct {
	inode       uint64
	name        string
	path        string
	directories []*SwarmDir
	files       []*SwarmFile

	mountInfo *MountInfo
	lock      *sync.RWMutex
}

func NewSwarmDir(fullpath string, minfo *MountInfo) *SwarmDir {
	log.Debug("swarmfs", "NewSwarmDir", fullpath)
	newdir := &SwarmDir{
		inode:       NewInode(),
		name:        filepath.Base(fullpath),
		path:        fullpath,
		directories: []*SwarmDir{},
		files:       []*SwarmFile{},
		mountInfo:   minfo,
		lock:        &sync.RWMutex{},
	}
	return newdir
}

func (sd *SwarmDir) Attr(ctx context.Context, a *fuse.Attr) error {
	sd.lock.RLock()
	defer sd.lock.RUnlock()
	a.Inode = sd.inode
	a.Mode = os.ModeDir | 0700
	a.Uid = uint32(os.Getuid())
	a.Gid = uint32(os.Getegid())
	return nil
}

func (sd *SwarmDir) Lookup(ctx context.Context, req *fuse.LookupRequest, resp *fuse.LookupResponse) (fs.Node, error) {
	log.Debug("swarmfs", "Lookup", req.Name)
	for _, n := range sd.files {
		if n.name == req.Name {
			return n, nil
		}
	}
	for _, n := range sd.directories {
		if n.name == req.Name {
			return n, nil
		}
	}
	return nil, fuse.ENOENT
}

func (sd *SwarmDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Debug("swarmfs ReadDirAll")
	var children []fuse.Dirent
	for _, file := range sd.files {
		children = append(children, fuse.Dirent{Inode: file.inode, Type: fuse.DT_File, Name: file.name})
	}
	for _, dir := range sd.directories {
		children = append(children, fuse.Dirent{Inode: dir.inode, Type: fuse.DT_Dir, Name: dir.name})
	}
	return children, nil
}

func (sd *SwarmDir) Create(ctx context.Context, req *fuse.CreateRequest, resp *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	log.Debug("swarmfs Create", "path", sd.path, "req.Name", req.Name)

	newFile := NewSwarmFile(sd.path, req.Name, sd.mountInfo)
newFile.fileSize = 0 //0意味着，文件还没有在Swarm中，只是创建了

	sd.lock.Lock()
	defer sd.lock.Unlock()
	sd.files = append(sd.files, newFile)

	return newFile, newFile, nil
}

func (sd *SwarmDir) Remove(ctx context.Context, req *fuse.RemoveRequest) error {
	log.Debug("swarmfs Remove", "path", sd.path, "req.Name", req.Name)

	if req.Dir && sd.directories != nil {
		newDirs := []*SwarmDir{}
		for _, dir := range sd.directories {
			if dir.name == req.Name {
				removeDirectoryFromSwarm(dir)
			} else {
				newDirs = append(newDirs, dir)
			}
		}
		if len(sd.directories) > len(newDirs) {
			sd.lock.Lock()
			defer sd.lock.Unlock()
			sd.directories = newDirs
		}
		return nil
	} else if !req.Dir && sd.files != nil {
		newFiles := []*SwarmFile{}
		for _, f := range sd.files {
			if f.name == req.Name {
				removeFileFromSwarm(f)
			} else {
				newFiles = append(newFiles, f)
			}
		}
		if len(sd.files) > len(newFiles) {
			sd.lock.Lock()
			defer sd.lock.Unlock()
			sd.files = newFiles
		}
		return nil
	}
	return fuse.ENOENT
}

func (sd *SwarmDir) Mkdir(ctx context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	log.Debug("swarmfs Mkdir", "path", sd.path, "req.Name", req.Name)
	newDir := NewSwarmDir(filepath.Join(sd.path, req.Name), sd.mountInfo)
	sd.lock.Lock()
	defer sd.lock.Unlock()
	sd.directories = append(sd.directories, newDir)

	return newDir, nil
}

