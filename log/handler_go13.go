
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450097773547520>

//+建设！GO1.4

package log

import (
	"sync/atomic"
	"unsafe"
)

//swaphandler包装了另一个可以换出的处理程序
//在运行时以线程安全的方式动态运行。
type swapHandler struct {
	handler unsafe.Pointer
}

func (h *swapHandler) Log(r *Record) error {
	return h.Get().Log(r)
}

func (h *swapHandler) Get() Handler {
	return *(*Handler)(atomic.LoadPointer(&h.handler))
}

func (h *swapHandler) Swap(newHandler Handler) {
	atomic.StorePointer(&h.handler, unsafe.Pointer(&newHandler))
}

