
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450083991064576>

//+建立AMD64，！通用ARM64，！通用的

package bn256

//此文件包含特定于体系结构的转发声明
//这些函数的程序集实现，前提是它们存在。

import (
	"golang.org/x/sys/cpu"
)

//诺林特
var hasBMI2 = cpu.X86.HasBMI2

//逃走
func gfpNeg(c, a *gfP)

//逃走
func gfpAdd(c, a, b *gfP)

//逃走
func gfpSub(c, a, b *gfP)

//逃走
func gfpMul(c, a, b *gfP)

