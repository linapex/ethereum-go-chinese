
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450072704192512>

//版权所有2013 Go作者。版权所有。
//此源代码的使用受BSD样式的控制
//可以在许可文件中找到的许可证。

//改编自：https://golang.org/src/crypto/cipher/xor.go

//包bitutil实现快速按位操作。
package bitutil

import (
	"runtime"
	"unsafe"
)

const wordSize = int(unsafe.Sizeof(uintptr(0)))
const supportsUnaligned = runtime.GOARCH == "386" || runtime.GOARCH == "amd64" || runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le" || runtime.GOARCH == "s390x"

//xorbytes xor是a和b中的字节。假定目标具有足够的
//空间。返回字节数xor'd。
func XORBytes(dst, a, b []byte) int {
	if supportsUnaligned {
		return fastXORBytes(dst, a, b)
	}
	return safeXORBytes(dst, a, b)
}

//FastXorBytes大容量XORS。它只适用于支持
//未对齐的读/写。
func fastXORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	w := n / wordSize
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] = aw[i] ^ bw[i]
		}
	}
	for i := n - n%wordSize; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}

//安全字节一个接一个。它适用于所有体系结构，独立于
//它是否支持未对齐的读/写。
func safeXORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return n
}

//AndBytes和A和B中的字节。假定目标具有足够的
//空间。返回字节数和'd'。
func ANDBytes(dst, a, b []byte) int {
	if supportsUnaligned {
		return fastANDBytes(dst, a, b)
	}
	return safeANDBytes(dst, a, b)
}

//FastAndBytes和批量。它只适用于支持
//未对齐的读/写。
func fastANDBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	w := n / wordSize
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] = aw[i] & bw[i]
		}
	}
	for i := n - n%wordSize; i < n; i++ {
		dst[i] = a[i] & b[i]
	}
	return n
}

//安全字节和一个接一个。它适用于所有体系结构，独立于
//它是否支持未对齐的读/写。
func safeANDBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] & b[i]
	}
	return n
}

//ORBYTES或A和B中的字节。假定目标具有足够的
//空间。返回字节数或'd'。
func ORBytes(dst, a, b []byte) int {
	if supportsUnaligned {
		return fastORBytes(dst, a, b)
	}
	return safeORBytes(dst, a, b)
}

//FastOrbytes或大量。它只适用于支持
//未对齐的读/写。
func fastORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	w := n / wordSize
	if w > 0 {
		dw := *(*[]uintptr)(unsafe.Pointer(&dst))
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			dw[i] = aw[i] | bw[i]
		}
	}
	for i := n - n%wordSize; i < n; i++ {
		dst[i] = a[i] | b[i]
	}
	return n
}

//安全字节或一个接一个。它适用于所有体系结构，独立于
//它是否支持未对齐的读/写。
func safeORBytes(dst, a, b []byte) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] | b[i]
	}
	return n
}

//测试字节测试输入字节片中是否设置了任何位。
func TestBytes(p []byte) bool {
	if supportsUnaligned {
		return fastTestBytes(p)
	}
	return safeTestBytes(p)
}

//FastTestBytes批量测试设置位。它只适用于那些
//支持未对齐的读/写。
func fastTestBytes(p []byte) bool {
	n := len(p)
	w := n / wordSize
	if w > 0 {
		pw := *(*[]uintptr)(unsafe.Pointer(&p))
		for i := 0; i < w; i++ {
			if pw[i] != 0 {
				return true
			}
		}
	}
	for i := n - n%wordSize; i < n; i++ {
		if p[i] != 0 {
			return true
		}
	}
	return false
}

//safetestBytes一次测试一个字节的设置位。它适用于所有
//体系结构，独立于是否支持未对齐的读/写。
func safeTestBytes(p []byte) bool {
	for i := 0; i < len(p); i++ {
		if p[i] != 0 {
			return true
		}
	}
	return false
}

