
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450107844071424>


package params

import (
	"fmt"
)

const (
VersionMajor = 1          //当前版本的主要版本组件
VersionMinor = 8          //当前版本的次要版本组件
VersionPatch = 22         //当前版本的补丁版本组件
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

//archiveversion保存用于geth存档的文本版本字符串。
//例如，“1.8.11-DEA1CE05”用于稳定释放，或
//“1.8.13-不稳定-21C059B6”用于不稳定释放
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
	vsn := VersionWithMeta
	if len(gitCommit) >= 8 {
		vsn += "-" + gitCommit[:8]
	}
	return vsn
}

