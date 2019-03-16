
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:40</date>
//</624450100806029312>


//包含来自帐户包的所有包装器以支持客户端密钥
//移动平台管理。

package geth

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
//标准加密是加密算法的n个参数，使用256MB
//在现代处理器上占用大约1秒的CPU时间。
	StandardScryptN = int(keystore.StandardScryptN)

//StandardScryptP是加密算法的P参数，使用256MB
//在现代处理器上占用大约1秒的CPU时间。
	StandardScryptP = int(keystore.StandardScryptP)

//lightscryptn是加密算法的n个参数，使用4MB
//在现代处理器上占用大约100毫秒的CPU时间。
	LightScryptN = int(keystore.LightScryptN)

//lightscryptp是加密算法的p参数，使用4MB
//在现代处理器上占用大约100毫秒的CPU时间。
	LightScryptP = int(keystore.LightScryptP)
)

//帐户表示存储的密钥。
type Account struct{ account accounts.Account }

//帐户代表帐户的一部分。
type Accounts struct{ accounts []accounts.Account }

//SIZE返回切片中的帐户数。
func (a *Accounts) Size() int {
	return len(a.accounts)
}

//GET从切片返回给定索引处的帐户。
func (a *Accounts) Get(index int) (account *Account, _ error) {
	if index < 0 || index >= len(a.accounts) {
		return nil, errors.New("index out of bounds")
	}
	return &Account{a.accounts[index]}, nil
}

//set在切片中的给定索引处设置帐户。
func (a *Accounts) Set(index int, account *Account) error {
	if index < 0 || index >= len(a.accounts) {
		return errors.New("index out of bounds")
	}
	a.accounts[index] = account.account
	return nil
}

//GetAddress检索与帐户关联的地址。
func (a *Account) GetAddress() *Address {
	return &Address{a.account.Address}
}

//GetURL检索帐户的规范URL。
func (a *Account) GetURL() string {
	return a.account.URL.String()
}

//keystore管理磁盘上的密钥存储目录。
type KeyStore struct{ keystore *keystore.KeyStore }

//newkeystore为给定目录创建一个keystore。
func NewKeyStore(keydir string, scryptN, scryptP int) *KeyStore {
	return &KeyStore{keystore: keystore.NewKeyStore(keydir, scryptN, scryptP)}
}

//hasAddress报告具有给定地址的密钥是否存在。
func (ks *KeyStore) HasAddress(address *Address) bool {
	return ks.keystore.HasAddress(address.address)
}

//getaccounts返回目录中存在的所有密钥文件。
func (ks *KeyStore) GetAccounts() *Accounts {
	return &Accounts{ks.keystore.Accounts()}
}

//如果密码短语正确，则删除帐户匹配的密钥。
//如果a不包含文件名，则地址必须与唯一键匹配。
func (ks *KeyStore) DeleteAccount(account *Account, passphrase string) error {
	return ks.keystore.Delete(account.account, passphrase)
}

//signhash为给定哈希计算ECDSA签名。生成的签名
//格式为[R V]，其中V为0或1。
func (ks *KeyStore) SignHash(address *Address, hash []byte) (signature []byte, _ error) {
	return ks.keystore.SignHash(accounts.Account{Address: address.address}, common.CopyBytes(hash))
}

//signtx用请求的帐户对给定的事务进行签名。
func (ks *KeyStore) SignTx(account *Account, tx *Transaction, chainID *BigInt) (*Transaction, error) {
if chainID == nil { //从移动应用程序传递的空值
		chainID = new(BigInt)
	}
	signed, err := ks.keystore.SignTx(account.account, tx.tx, chainID.bigint)
	if err != nil {
		return nil, err
	}
	return &Transaction{signed}, nil
}

//如果与给定地址匹配的私钥可以
//用给定的密码短语解密。生成的签名在
//[R_S_V]格式，其中V为0或1。
func (ks *KeyStore) SignHashPassphrase(account *Account, passphrase string, hash []byte) (signature []byte, _ error) {
	return ks.keystore.SignHashWithPassphrase(account.account, passphrase, common.CopyBytes(hash))
}

//signtxpassphrase如果私钥与
//给定的地址可以用给定的密码短语解密。
func (ks *KeyStore) SignTxPassphrase(account *Account, passphrase string, tx *Transaction, chainID *BigInt) (*Transaction, error) {
if chainID == nil { //从移动应用程序传递的空值
		chainID = new(BigInt)
	}
	signed, err := ks.keystore.SignTxWithPassphrase(account.account, passphrase, tx.tx, chainID.bigint)
	if err != nil {
		return nil, err
	}
	return &Transaction{signed}, nil
}

//解锁无限期地解锁给定的帐户。
func (ks *KeyStore) Unlock(account *Account, passphrase string) error {
	return ks.keystore.TimedUnlock(account.account, passphrase, 0)
}

//lock从内存中删除具有给定地址的私钥。
func (ks *KeyStore) Lock(address *Address) error {
	return ks.keystore.Lock(address.address)
}

//timedunlock使用密码短语解锁给定的帐户。帐户保持不变
//在超时期间解锁（纳秒）。超时0将解锁
//帐户，直到程序退出。帐户必须与唯一的密钥文件匹配。
//
//如果帐户地址在一段时间内已解锁，则TimedUnlock将扩展或
//缩短活动解锁超时。如果地址以前是解锁的
//无限期地超时不会改变。
func (ks *KeyStore) TimedUnlock(account *Account, passphrase string, timeout int64) error {
	return ks.keystore.TimedUnlock(account.account, passphrase, time.Duration(timeout))
}

//newaccount生成一个新密钥并将其存储到密钥目录中，
//用密码短语加密。
func (ks *KeyStore) NewAccount(passphrase string) (*Account, error) {
	account, err := ks.keystore.NewAccount(passphrase)
	if err != nil {
		return nil, err
	}
	return &Account{account}, nil
}

//更新帐户更改现有帐户的密码。
func (ks *KeyStore) UpdateAccount(account *Account, passphrase, newPassphrase string) error {
	return ks.keystore.Update(account.account, passphrase, newPassphrase)
}

//exportkey作为json密钥导出，用newpassphrase加密。
func (ks *KeyStore) ExportKey(account *Account, passphrase, newPassphrase string) (key []byte, _ error) {
	return ks.keystore.Export(account.account, passphrase, newPassphrase)
}

//importkey将给定的加密JSON密钥存储到密钥目录中。
func (ks *KeyStore) ImportKey(keyJSON []byte, passphrase, newPassphrase string) (account *Account, _ error) {
	acc, err := ks.keystore.Import(common.CopyBytes(keyJSON), passphrase, newPassphrase)
	if err != nil {
		return nil, err
	}
	return &Account{acc}, nil
}

//importecdsakey将给定的加密JSON密钥存储到密钥目录中。
func (ks *KeyStore) ImportECDSAKey(key []byte, passphrase string) (account *Account, _ error) {
	privkey, err := crypto.ToECDSA(common.CopyBytes(key))
	if err != nil {
		return nil, err
	}
	acc, err := ks.keystore.ImportECDSA(privkey, passphrase)
	if err != nil {
		return nil, err
	}
	return &Account{acc}, nil
}

//importpresalekey解密给定的以太坊预售钱包和商店
//密钥目录中的密钥文件。密钥文件使用相同的密码短语加密。
func (ks *KeyStore) ImportPreSaleKey(keyJSON []byte, passphrase string) (ccount *Account, _ error) {
	account, err := ks.keystore.ImportPreSaleKey(common.CopyBytes(keyJSON), passphrase)
	if err != nil {
		return nil, err
	}
	return &Account{account}, nil
}

