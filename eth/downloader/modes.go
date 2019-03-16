
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450088206340096>


package downloader

import "fmt"

//SyncMode表示下载程序的同步模式。
type SyncMode int

const (
FullSync  SyncMode = iota //从完整块同步整个区块链历史
FastSync                  //快速下载邮件头，仅在链头完全同步
LightSync                 //只下载邮件头，然后终止
)

func (mode SyncMode) IsValid() bool {
	return mode >= FullSync && mode <= LightSync
}

//字符串实现字符串接口。
func (mode SyncMode) String() string {
	switch mode {
	case FullSync:
		return "full"
	case FastSync:
		return "fast"
	case LightSync:
		return "light"
	default:
		return "unknown"
	}
}

func (mode SyncMode) MarshalText() ([]byte, error) {
	switch mode {
	case FullSync:
		return []byte("full"), nil
	case FastSync:
		return []byte("fast"), nil
	case LightSync:
		return []byte("light"), nil
	default:
		return nil, fmt.Errorf("unknown sync mode %d", mode)
	}
}

func (mode *SyncMode) UnmarshalText(text []byte) error {
	switch string(text) {
	case "full":
		*mode = FullSync
	case "fast":
		*mode = FastSync
	case "light":
		*mode = LightSync
	default:
		return fmt.Errorf(`unknown sync mode %q, want "full", "fast" or "light"`, text)
	}
	return nil
}

