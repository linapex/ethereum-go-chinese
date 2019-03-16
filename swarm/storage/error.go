
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450118564712448>


package storage

import (
	"errors"
)

const (
	ErrInit = iota
	ErrNotFound
	ErrUnauthorized
	ErrInvalidValue
	ErrDataOverflow
	ErrNothingToReturn
	ErrInvalidSignature
	ErrNotSynced
)

var (
	ErrChunkNotFound = errors.New("chunk not found")
	ErrChunkInvalid  = errors.New("invalid chunk")
)

