
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:31</date>
//</624450063669661696>


package accounts

import (
	"errors"
	"fmt"
)

//对于没有后端的任何请求操作，将返回errUnknownAccount。
//提供指定的帐户。
var ErrUnknownAccount = errors.New("unknown account")

//对于没有后端的任何请求操作，将返回errunknownwallet。
//提供指定的钱包。
var ErrUnknownWallet = errors.New("unknown wallet")

//从帐户请求操作时返回errnotsupported
//它不支持的后端。
var ErrNotSupported = errors.New("not supported")

//当解密操作收到错误消息时返回errInvalidPassphrase
//口令。
var ErrInvalidPassphrase = errors.New("invalid passphrase")

//如果试图打开钱包，则返回errwalletalreadyopen
//第二次。
var ErrWalletAlreadyOpen = errors.New("wallet already open")

//如果试图打开钱包，则返回errWalletClosed
//间隔时间。
var ErrWalletClosed = errors.New("wallet closed")

//后端返回authneedederror，用于在用户
//在签名成功之前需要提供进一步的身份验证。
//
//这通常意味着要么需要提供密码，要么
//某些硬件设备显示的一次性PIN码。
type AuthNeededError struct {
Needed string //用户需要提供的额外身份验证
}

//newauthneedederror创建一个新的身份验证错误，并提供额外的详细信息
//关于所需字段集。
func NewAuthNeededError(needed string) error {
	return &AuthNeededError{
		Needed: needed,
	}
}

//错误实现标准错误接口。
func (err *AuthNeededError) Error() string {
	return fmt.Sprintf("authentication needed: %s", err.Needed)
}

