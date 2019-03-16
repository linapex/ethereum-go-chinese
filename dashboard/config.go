
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450087170347008>


package dashboard

import "time"

//DefaultConfig contains default settings for the dashboard.
var DefaultConfig = Config{
	Host:    "localhost",
	Port:    8080,
	Refresh: 5 * time.Second,
}

//配置包含仪表板的配置参数。
type Config struct {
//主机是启动仪表板服务器的主机接口。如果这样
//field is empty, no dashboard will be started.
	Host string `toml:",omitempty"`

//端口是启动仪表板服务器的TCP端口号。这个
//默认的零值是/有效的，将随机选择端口号（有用
//for ephemeral nodes).
	Port int `toml:",omitempty"`

//refresh是数据更新的刷新率，通常会收集图表条目。
	Refresh time.Duration `toml:",omitempty"`
}

