
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450092388061184>


package jsre

import (
	"sort"
	"strings"

	"github.com/robertkrimen/otto"
)

//completekeywords返回给定行的潜在连续性。既然是
//经过评估，调用方需要确保评估行没有副作用。
func (jsre *JSRE) CompleteKeywords(line string) []string {
	var results []string
	jsre.Do(func(vm *otto.Otto) {
		results = getCompletions(vm, line)
	})
	return results
}

func getCompletions(vm *otto.Otto, line string) (results []string) {
	parts := strings.Split(line, ".")
	objRef := "this"
	prefix := line
	if len(parts) > 1 {
		objRef = strings.Join(parts[0:len(parts)-1], ".")
		prefix = parts[len(parts)-1]
	}

	obj, _ := vm.Object(objRef)
	if obj == nil {
		return nil
	}
	iterOwnAndConstructorKeys(vm, obj, func(k string) {
		if strings.HasPrefix(k, prefix) {
			if objRef == "this" {
				results = append(results, k)
			} else {
				results = append(results, strings.Join(parts[:len(parts)-1], ".")+"."+k)
			}
		}
	})

//附加左括号（用于函数）或点（用于对象）
//如果行本身是唯一完成的。
	if len(results) == 1 && results[0] == line {
		obj, _ := vm.Object(line)
		if obj != nil {
			if obj.Class() == "Function" {
				results[0] += "("
			} else {
				results[0] += "."
			}
		}
	}

	sort.Strings(results)
	return results
}

