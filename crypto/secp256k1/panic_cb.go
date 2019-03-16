
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:36</date>
//</624450085383573506>

//版权所有2015 Jeffrey Wilcke、Felix Lange、Gustav Simonsson。版权所有。
//此源代码的使用受BSD样式许可证的控制，该许可证可在
//许可证文件。

package secp256k1

import "C"
import "unsafe"

//将libsecp256k1内部故障转换为
//恢复性恐慌。

//出口secp256k1gopanicilegal
func secp256k1GoPanicIllegal(msg *C.char, data unsafe.Pointer) {
	panic("illegal argument: " + C.GoString(msg))
}

//导出secp256k1gopanicerror
func secp256k1GoPanicError(msg *C.char, data unsafe.Pointer) {
	panic("internal error: " + C.GoString(msg))
}

