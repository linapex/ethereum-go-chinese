
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450111136600064>

//

package storage

import (
	"fmt"
)

type Storage interface {
//按键存储值。0长度键导致无操作
	Put(key, value string)
//get返回以前存储的值，如果该值不存在或键的长度为0，则返回空字符串
	Get(key string) string
}

//短暂存储是一种内存存储，它可以
//不将值持久化到磁盘。主要用于测试
type EphemeralStorage struct {
	data      map[string]string
	namespace string
}

func (s *EphemeralStorage) Put(key, value string) {
	if len(key) == 0 {
		return
	}
	fmt.Printf("storage: put %v -> %v\n", key, value)
	s.data[key] = value
}

func (s *EphemeralStorage) Get(key string) string {
	if len(key) == 0 {
		return ""
	}
	fmt.Printf("storage: get %v\n", key)
	if v, exist := s.data[key]; exist {
		return v
	}
	return ""
}

func NewEphemeralStorage() Storage {
	s := &EphemeralStorage{
		data: make(map[string]string),
	}
	return s
}

