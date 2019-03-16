
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:45</date>
//</624450121848852480>


package version

import (
	"fmt"
)

const (
VersionMajor = 0          //当前版本的主要版本组件
VersionMinor = 3          //当前版本的次要版本组件
VersionPatch = 10         //当前版本的补丁版本组件
VersionMeta  = "unstable" //要附加到版本字符串的版本元数据
)

//version保存文本版本字符串。
var Version = func() string {
	return fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
}()

//versionWithMeta保存包含元数据的文本版本字符串。
var VersionWithMeta = func() string {
	v := Version
	if VersionMeta != "" {
		v += "-" + VersionMeta
	}
	return v
}()

//archiveversion保存用于swarm存档的文本版本字符串。
//例如，“0.3.0-DEA1CE05”用于稳定释放，或
//“0.3.1-不稳定-21C059B6”用于不稳定释放
func ArchiveVersion(gitCommit string) string {
	vsn := Version
	if VersionMeta != "stable" {
		vsn += "-" + VersionMeta
	}
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	return vsn
}

func VersionWithCommit(gitCommit string) string {
	vsn := Version
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	return vsn
}

