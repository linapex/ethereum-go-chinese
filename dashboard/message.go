
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450087409422336>


package dashboard

import (
	"encoding/json"
	"time"
)

type Message struct {
	General *GeneralMessage `json:"general,omitempty"`
	Home    *HomeMessage    `json:"home,omitempty"`
	Chain   *ChainMessage   `json:"chain,omitempty"`
	TxPool  *TxPoolMessage  `json:"txpool,omitempty"`
	Network *NetworkMessage `json:"network,omitempty"`
	System  *SystemMessage  `json:"system,omitempty"`
	Logs    *LogsMessage    `json:"logs,omitempty"`
}

type ChartEntries []*ChartEntry

type ChartEntry struct {
	Time  time.Time `json:"time,omitempty"`
	Value float64   `json:"value,omitempty"`
}

type GeneralMessage struct {
	Version string `json:"version,omitempty"`
	Commit  string `json:"commit,omitempty"`
}

type HomeMessage struct {
 /*托多（Kurkomisi）*/
}

type ChainMessage struct {
 /*托多（Kurkomisi）*/
}

type TxPoolMessage struct {
 /*托多（Kurkomisi）*/
}

type NetworkMessage struct {
 /*托多（Kurkomisi）*/
}

type SystemMessage struct {
	ActiveMemory   ChartEntries `json:"activeMemory,omitempty"`
	VirtualMemory  ChartEntries `json:"virtualMemory,omitempty"`
	NetworkIngress ChartEntries `json:"networkIngress,omitempty"`
	NetworkEgress  ChartEntries `json:"networkEgress,omitempty"`
	ProcessCPU     ChartEntries `json:"processCPU,omitempty"`
	SystemCPU      ChartEntries `json:"systemCPU,omitempty"`
	DiskRead       ChartEntries `json:"diskRead,omitempty"`
	DiskWrite      ChartEntries `json:"diskWrite,omitempty"`
}

//logsmessage包装了一个日志块。如果源不存在，则块是流块。
type LogsMessage struct {
Source *LogFile        `json:"source,omitempty"` //日志文件的属性。
Chunk  json.RawMessage `json:"chunk"`            //包含日志记录。
}

//日志文件包含日志文件的属性。
type LogFile struct {
Name string `json:"name"` //文件名。
Last bool   `json:"last"` //指示实际日志文件是否是目录中的最后一个。
}

//请求表示客户端请求。
type Request struct {
	Logs *LogsRequest `json:"logs,omitempty"`
}

type LogsRequest struct {
Name string `json:"name"` //请求处理程序根据此文件名搜索日志文件。
Past bool   `json:"past"` //指示客户端是要上一个文件还是下一个文件。
}

