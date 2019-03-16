
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450118682152960>


package feed

import (
	"bytes"
	"context"
	"time"

	"github.com/ethereum/go-ethereum/swarm/storage"
)

const (
	hasherCount            = 8
	feedsHashAlgorithm     = storage.SHA3Hash
	defaultRetrieveTimeout = 100 * time.Millisecond
)

//cacheEntry缓存特定群源的最后一次已知更新。
type cacheEntry struct {
	Update
	*bytes.Reader
	lastKey storage.Address
}

//实现Storage.LazySectionReader
func (r *cacheEntry) Size(ctx context.Context, _ chan bool) (int64, error) {
	return int64(len(r.Update.data)), nil
}

//返回源的主题
func (r *cacheEntry) Topic() Topic {
	return r.Feed.Topic
}

