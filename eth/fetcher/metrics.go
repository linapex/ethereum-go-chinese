
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450088667713536>


//包含由提取程序收集的度量。

package fetcher

import (
	"github.com/ethereum/go-ethereum/metrics"
)

var (
	propAnnounceInMeter   = metrics.NewRegisteredMeter("eth/fetcher/prop/announces/in", nil)
	propAnnounceOutTimer  = metrics.NewRegisteredTimer("eth/fetcher/prop/announces/out", nil)
	propAnnounceDropMeter = metrics.NewRegisteredMeter("eth/fetcher/prop/announces/drop", nil)
	propAnnounceDOSMeter  = metrics.NewRegisteredMeter("eth/fetcher/prop/announces/dos", nil)

	propBroadcastInMeter   = metrics.NewRegisteredMeter("eth/fetcher/prop/broadcasts/in", nil)
	propBroadcastOutTimer  = metrics.NewRegisteredTimer("eth/fetcher/prop/broadcasts/out", nil)
	propBroadcastDropMeter = metrics.NewRegisteredMeter("eth/fetcher/prop/broadcasts/drop", nil)
	propBroadcastDOSMeter  = metrics.NewRegisteredMeter("eth/fetcher/prop/broadcasts/dos", nil)

	headerFetchMeter = metrics.NewRegisteredMeter("eth/fetcher/fetch/headers", nil)
	bodyFetchMeter   = metrics.NewRegisteredMeter("eth/fetcher/fetch/bodies", nil)

	headerFilterInMeter  = metrics.NewRegisteredMeter("eth/fetcher/filter/headers/in", nil)
	headerFilterOutMeter = metrics.NewRegisteredMeter("eth/fetcher/filter/headers/out", nil)
	bodyFilterInMeter    = metrics.NewRegisteredMeter("eth/fetcher/filter/bodies/in", nil)
	bodyFilterOutMeter   = metrics.NewRegisteredMeter("eth/fetcher/filter/bodies/out", nil)
)

