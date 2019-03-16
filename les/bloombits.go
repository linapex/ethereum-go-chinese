
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450093260476416>


package les

import (
	"time"

	"github.com/ethereum/go-ethereum/common/bitutil"
	"github.com/ethereum/go-ethereum/light"
)

const (
//BloomServiceThreads是以太坊全局使用的Goroutine数。
//实例到服务BloomBits查找所有正在运行的筛选器。
	bloomServiceThreads = 16

//BloomFilterThreads是每个筛选器本地使用的goroutine数，用于
//将请求多路传输到全局服务goroutine。
	bloomFilterThreads = 3

//BloomRetrievalBatch是要服务的最大Bloom位检索数。
//一批。
	bloomRetrievalBatch = 16

//BloomRetrievalWait是等待足够的Bloom位请求的最长时间。
//累积请求整个批（避免滞后）。
	bloomRetrievalWait = time.Microsecond * 100
)

//StartBloomHandlers启动一批Goroutine以接受BloomBit数据库
//从可能的一系列过滤器中检索并为数据提供满足条件的服务。
func (eth *LightEthereum) startBloomHandlers(sectionSize uint64) {
	for i := 0; i < bloomServiceThreads; i++ {
		go func() {
			for {
				select {
				case <-eth.shutdownChan:
					return

				case request := <-eth.bloomRequests:
					task := <-request
					task.Bitsets = make([][]byte, len(task.Sections))
					compVectors, err := light.GetBloomBits(task.Context, eth.odr, task.Bit, task.Sections)
					if err == nil {
						for i := range task.Sections {
							if blob, err := bitutil.DecompressBytes(compVectors[i], int(sectionSize/8)); err == nil {
								task.Bitsets[i] = blob
							} else {
								task.Error = err
							}
						}
					} else {
						task.Error = err
					}
					request <- task
				}
			}
		}()
	}
}

