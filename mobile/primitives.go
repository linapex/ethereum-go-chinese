
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450101758136320>


//包含用于基元类型的各种包装器。

package geth

import (
	"errors"
	"fmt"
)

//字符串表示str的s切片。
type Strings struct{ strs []string }

//SIZE返回切片中str的数目。
func (s *Strings) Size() int {
	return len(s.strs)
}

//get从切片返回给定索引处的字符串。
func (s *Strings) Get(index int) (str string, _ error) {
	if index < 0 || index >= len(s.strs) {
		return "", errors.New("index out of bounds")
	}
	return s.strs[index], nil
}

//set在切片中的给定索引处设置字符串。
func (s *Strings) Set(index int, str string) error {
	if index < 0 || index >= len(s.strs) {
		return errors.New("index out of bounds")
	}
	s.strs[index] = str
	return nil
}

//字符串实现字符串接口。
func (s *Strings) String() string {
	return fmt.Sprintf("%v", s.strs)
}

