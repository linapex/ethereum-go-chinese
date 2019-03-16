
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450088852262912>


package filters

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/bloombits"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/rpc"
)

type Backend interface {
	ChainDb() ethdb.Database
	EventMux() *event.TypeMux
	HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error)
	HeaderByHash(ctx context.Context, blockHash common.Hash) (*types.Header, error)
	GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error)
	GetLogs(ctx context.Context, blockHash common.Hash) ([][]*types.Log, error)

	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
	SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription
	SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription
	SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription

	BloomStatus() (uint64, uint64)
	ServiceFilter(ctx context.Context, session *bloombits.MatcherSession)
}

//筛选器可用于检索和筛选日志。
type Filter struct {
	backend Backend

	db        ethdb.Database
	addresses []common.Address
	topics    [][]common.Hash

block      common.Hash //如果筛选单个块，则阻止哈希
begin, end int64       //过滤多个块时的范围间隔

	matcher *bloombits.Matcher
}

//newrangefilter创建一个新的过滤器，它在块上使用bloom过滤器来
//找出一个特定的块是否有趣。
func NewRangeFilter(backend Backend, begin, end int64, addresses []common.Address, topics [][]common.Hash) *Filter {
//将地址和主题筛选子句展平为单个bloombits筛选器
//系统。因为bloombits不是位置的，所以不允许使用任何主题，
//它被压扁成一个零字节的片。
	var filters [][][]byte
	if len(addresses) > 0 {
		filter := make([][]byte, len(addresses))
		for i, address := range addresses {
			filter[i] = address.Bytes()
		}
		filters = append(filters, filter)
	}
	for _, topicList := range topics {
		filter := make([][]byte, len(topicList))
		for i, topic := range topicList {
			filter[i] = topic.Bytes()
		}
		filters = append(filters, filter)
	}
	size, _ := backend.BloomStatus()

//创建通用筛选器并将其转换为范围筛选器
	filter := newFilter(backend, addresses, topics)

	filter.matcher = bloombits.NewMatcher(size, filters)
	filter.begin = begin
	filter.end = end

	return filter
}

//newblockfilter创建一个新的过滤器，它直接检查
//用来判断它是否有趣的块。
func NewBlockFilter(backend Backend, block common.Hash, addresses []common.Address, topics [][]common.Hash) *Filter {
//创建通用筛选器并将其转换为块筛选器
	filter := newFilter(backend, addresses, topics)
	filter.block = block
	return filter
}

//newfilter创建一个通用筛选器，该筛选器可以基于块哈希进行筛选，
//或者基于范围查询。需要显式设置搜索条件。
func newFilter(backend Backend, addresses []common.Address, topics [][]common.Hash) *Filter {
	return &Filter{
		backend:   backend,
		addresses: addresses,
		topics:    topics,
		db:        backend.ChainDb(),
	}
}

//日志在区块链中搜索匹配的日志条目，从
//包含匹配项的第一个块，相应地更新筛选器的开头。
func (f *Filter) Logs(ctx context.Context) ([]*types.Log, error) {
//如果我们进行单例块过滤，执行并返回
	if f.block != (common.Hash{}) {
		header, err := f.backend.HeaderByHash(ctx, f.block)
		if err != nil {
			return nil, err
		}
		if header == nil {
			return nil, errors.New("unknown block")
		}
		return f.blockLogs(ctx, header)
	}
//找出过滤范围的限制
	header, _ := f.backend.HeaderByNumber(ctx, rpc.LatestBlockNumber)
	if header == nil {
		return nil, nil
	}
	head := header.Number.Uint64()

	if f.begin == -1 {
		f.begin = int64(head)
	}
	end := uint64(f.end)
	if f.end == -1 {
		end = head
	}
//收集所有索引日志，并使用非索引日志完成
	var (
		logs []*types.Log
		err  error
	)
	size, sections := f.backend.BloomStatus()
	if indexed := sections * size; indexed > uint64(f.begin) {
		if indexed > end {
			logs, err = f.indexedLogs(ctx, end)
		} else {
			logs, err = f.indexedLogs(ctx, indexed-1)
		}
		if err != nil {
			return logs, err
		}
	}
	rest, err := f.unindexedLogs(ctx, end)
	logs = append(logs, rest...)
	return logs, err
}

//indexedlogs返回与基于bloom的筛选条件匹配的日志
//在本地或通过网络可用的索引位。
func (f *Filter) indexedLogs(ctx context.Context, end uint64) ([]*types.Log, error) {
//创建Matcher会话并从后端请求服务
	matches := make(chan uint64, 64)

	session, err := f.matcher.Start(ctx, uint64(f.begin), end, matches)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	f.backend.ServiceFilter(ctx, session)

//迭代匹配项，直到耗尽或上下文关闭
	var logs []*types.Log

	for {
		select {
		case number, ok := <-matches:
//如果满足所有匹配，则中止
			if !ok {
				err := session.Error()
				if err == nil {
					f.begin = int64(end) + 1
				}
				return logs, err
			}
			f.begin = int64(number) + 1

//检索建议的块并提取任何真正匹配的日志
			header, err := f.backend.HeaderByNumber(ctx, rpc.BlockNumber(number))
			if header == nil || err != nil {
				return logs, err
			}
			found, err := f.checkMatches(ctx, header)
			if err != nil {
				return logs, err
			}
			logs = append(logs, found...)

		case <-ctx.Done():
			return logs, ctx.Err()
		}
	}
}

//indexedlogs返回与基于原始块的筛选条件匹配的日志
//迭代和开花匹配。
func (f *Filter) unindexedLogs(ctx context.Context, end uint64) ([]*types.Log, error) {
	var logs []*types.Log

	for ; f.begin <= int64(end); f.begin++ {
		header, err := f.backend.HeaderByNumber(ctx, rpc.BlockNumber(f.begin))
		if header == nil || err != nil {
			return logs, err
		}
		found, err := f.blockLogs(ctx, header)
		if err != nil {
			return logs, err
		}
		logs = append(logs, found...)
	}
	return logs, nil
}

//block logs返回与单个块中的筛选条件匹配的日志。
func (f *Filter) blockLogs(ctx context.Context, header *types.Header) (logs []*types.Log, err error) {
	if bloomFilter(header.Bloom, f.addresses, f.topics) {
		found, err := f.checkMatches(ctx, header)
		if err != nil {
			return logs, err
		}
		logs = append(logs, found...)
	}
	return logs, nil
}

//checkmatches检查属于给定头的收据是否包含
//匹配筛选条件。当布卢姆滤波器发出潜在匹配信号时，调用此函数。
func (f *Filter) checkMatches(ctx context.Context, header *types.Header) (logs []*types.Log, err error) {
//获取块的日志
	logsList, err := f.backend.GetLogs(ctx, header.Hash())
	if err != nil {
		return nil, err
	}
	var unfiltered []*types.Log
	for _, logs := range logsList {
		unfiltered = append(unfiltered, logs...)
	}
	logs = filterLogs(unfiltered, nil, nil, f.addresses, f.topics)
	if len(logs) > 0 {
//我们有匹配的日志，检查是否需要通过Light客户端解析完整的日志
		if logs[0].TxHash == (common.Hash{}) {
			receipts, err := f.backend.GetReceipts(ctx, header.Hash())
			if err != nil {
				return nil, err
			}
			unfiltered = unfiltered[:0]
			for _, receipt := range receipts {
				unfiltered = append(unfiltered, receipt.Logs...)
			}
			logs = filterLogs(unfiltered, nil, nil, f.addresses, f.topics)
		}
		return logs, nil
	}
	return nil, nil
}

func includes(addresses []common.Address, a common.Address) bool {
	for _, addr := range addresses {
		if addr == a {
			return true
		}
	}

	return false
}

//FieldLtG创建一个与给定标准匹配的日志片段。
func filterLogs(logs []*types.Log, fromBlock, toBlock *big.Int, addresses []common.Address, topics [][]common.Hash) []*types.Log {
	var ret []*types.Log
Logs:
	for _, log := range logs {
		if fromBlock != nil && fromBlock.Int64() >= 0 && fromBlock.Uint64() > log.BlockNumber {
			continue
		}
		if toBlock != nil && toBlock.Int64() >= 0 && toBlock.Uint64() < log.BlockNumber {
			continue
		}

		if len(addresses) > 0 && !includes(addresses, log.Address) {
			continue
		}
//如果到筛选的主题大于日志中的主题数量，则跳过。
		if len(topics) > len(log.Topics) {
			continue Logs
		}
		for i, sub := range topics {
match := len(sub) == 0 //empty rule set == wildcard
			for _, topic := range sub {
				if log.Topics[i] == topic {
					match = true
					break
				}
			}
			if !match {
				continue Logs
			}
		}
		ret = append(ret, log)
	}
	return ret
}

func bloomFilter(bloom types.Bloom, addresses []common.Address, topics [][]common.Hash) bool {
	if len(addresses) > 0 {
		var included bool
		for _, addr := range addresses {
			if types.BloomLookup(bloom, addr) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	for _, sub := range topics {
included := len(sub) == 0 //空规则集==通配符
		for _, topic := range sub {
			if types.BloomLookup(bloom, topic) {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}
	return true
}

