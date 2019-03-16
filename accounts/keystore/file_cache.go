
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:31</date>
//</624450064131035136>


package keystore

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/ethereum/go-ethereum/log"
)

//filecache是在扫描密钥库期间看到的文件的缓存。
type fileCache struct {
all     mapset.Set //从keystore文件夹中设置所有文件
lastMod time.Time  //上次修改文件时的实例
	mu      sync.RWMutex
}

//scan对给定的目录执行新扫描，与
//缓存文件名，并返回文件集：创建、删除、更新。
func (fc *fileCache) scan(keyDir string) (mapset.Set, mapset.Set, mapset.Set, error) {
	t0 := time.Now()

//列出keystore文件夹中的所有故障
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return nil, nil, nil, err
	}
	t1 := time.Now()

	fc.mu.Lock()
	defer fc.mu.Unlock()

//迭代所有文件并收集其元数据
	all := mapset.NewThreadUnsafeSet()
	mods := mapset.NewThreadUnsafeSet()

	var newLastMod time.Time
	for _, fi := range files {
		path := filepath.Join(keyDir, fi.Name())
//跳过文件夹中的任何非关键文件
		if nonKeyFile(fi) {
			log.Trace("Ignoring file on account scan", "path", path)
			continue
		}
//收集所有修改过的文件集
		all.Add(path)

		modified := fi.ModTime()
		if modified.After(fc.lastMod) {
			mods.Add(path)
		}
		if modified.After(newLastMod) {
			newLastMod = modified
		}
	}
	t2 := time.Now()

//更新跟踪文件并返回三组
deletes := fc.all.Difference(all)   //删除=上一个-当前
creates := all.Difference(fc.all)   //创建=当前-上一个
updates := mods.Difference(creates) //更新=修改-创建

	fc.all, fc.lastMod = all, newLastMod
	t3 := time.Now()

//报告扫描数据并返回
	log.Debug("FS scan times", "list", t1.Sub(t0), "set", t2.Sub(t1), "diff", t3.Sub(t2))
	return creates, deletes, updates, nil
}

//非关键文件忽略编辑器备份、隐藏文件和文件夹/符号链接。
func nonKeyFile(fi os.FileInfo) bool {
//跳过编辑器备份和Unix样式的隐藏文件。
	if strings.HasSuffix(fi.Name(), "~") || strings.HasPrefix(fi.Name(), ".") {
		return true
	}
//跳过其他特殊文件、目录（是，也可以跳过符号链接）。
	if fi.IsDir() || fi.Mode()&os.ModeType != 0 {
		return true
	}
	return false
}

