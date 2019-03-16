
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450089942781952>


//包跟踪程序是JavaScript事务跟踪程序的集合。
package tracers

import (
	"strings"
	"unicode"

	"github.com/ethereum/go-ethereum/eth/tracers/internal/tracers"
)

//全部按名称包含所有内置的javascript跟踪程序。
var all = make(map[string]string)

//camel将snake-cased输入字符串转换为camel-cased输出。
func camel(str string) string {
	pieces := strings.Split(str, "_")
	for i := 1; i < len(pieces); i++ {
		pieces[i] = string(unicode.ToUpper(rune(pieces[i][0]))) + pieces[i][1:]
	}
	return strings.Join(pieces, "")
}

//init检索go-ethereum中包含的javascript事务跟踪程序。
func init() {
	for _, file := range tracers.AssetNames() {
		name := camel(strings.TrimSuffix(file, ".js"))
		all[name] = string(tracers.MustAsset(file))
	}
}

//跟踪程序按名称检索特定的javascript跟踪程序。
func tracer(name string) (string, bool) {
	if tracer, ok := all[name]; ok {
		return tracer, true
	}
	return "", false
}

