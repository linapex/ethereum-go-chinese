
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:46</date>
//</624450125170741248>


package whisperv6

//config表示耳语节点的配置状态。
type Config struct {
	MaxMessageSize                        uint32  `toml:",omitempty"`
	MinimumAcceptedPOW                    float64 `toml:",omitempty"`
	RestrictConnectionBetweenLightClients bool    `toml:",omitempty"`
}

//defaultconfig表示（shocker！）默认配置。
var DefaultConfig = Config{
	MaxMessageSize:                        DefaultMaxMessageSize,
	MinimumAcceptedPOW:                    DefaultMinimumPoW,
	RestrictConnectionBetweenLightClients: true,
}

