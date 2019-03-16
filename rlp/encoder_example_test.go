
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450108108312576>


package rlp

import (
	"fmt"
	"io"
)

type MyCoolType struct {
	Name string
	a, b uint
}

//encoderlp将x写为rlp list[a，b]，省略name字段。
func (x *MyCoolType) EncodeRLP(w io.Writer) (err error) {
//注意：接收器可以是零指针。这允许你
//控制nil的编码，但这也意味着必须
//检查零接收器。
	if x == nil {
		err = Encode(w, []uint{0, 0})
	} else {
		err = Encode(w, []uint{x.a, x.b})
	}
	return err
}

func ExampleEncoder() {
var t *MyCoolType //T为零，指向mycoltype的指针
	bytes, _ := EncodeToBytes(t)
	fmt.Printf("%v → %X\n", t, bytes)

	t = &MyCoolType{Name: "foobar", a: 5, b: 6}
	bytes, _ = EncodeToBytes(t)
	fmt.Printf("%v → %X\n", t, bytes)

//输出：
//<nil>→C28080
//&foobar 5 6→C20506
}

