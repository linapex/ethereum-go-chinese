
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:38</date>
//</624450092002185216>


//+构建GO1.6

package debug

import "runtime/debug"

//响亮的恐慌以一种方式让所有的血腥堆栈打印在stderr上。
func LoudPanic(x interface{}) {
	debug.SetTraceback("all")
	panic(x)
}

