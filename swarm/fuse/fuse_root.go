
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450112889819136>


//+构建Linux Darwin Freebsd

package fuse

import (
	"bazil.org/fuse/fs"
)

var (
	_ fs.Node = (*SwarmDir)(nil)
)

type SwarmRoot struct {
	root *SwarmDir
}

func (filesystem *SwarmRoot) Root() (fs.Node, error) {
	return filesystem.root, nil
}

