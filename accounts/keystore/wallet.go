
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:31</date>
//</624450064969895936>


package keystore

import (
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/core/types"
)

//keystorewallet实现原始帐户的accounts.wallet接口
//密钥存储区。
type keystoreWallet struct {
account  accounts.Account //钱包里有一个账户
keystore *KeyStore        //帐户来源的密钥库
}

//url实现accounts.wallet，返回帐户的url。
func (w *keystoreWallet) URL() accounts.URL {
	return w.account.URL
}

//状态实现accounts.wallet，返回
//密钥存储钱包是否已解锁。
func (w *keystoreWallet) Status() (string, error) {
	w.keystore.mu.RLock()
	defer w.keystore.mu.RUnlock()

	if _, ok := w.keystore.unlocked[w.account.Address]; ok {
		return "Unlocked", nil
	}
	return "Locked", nil
}

//打开工具帐户。钱包，但是普通钱包的一个noop，因为那里
//访问帐户列表不需要连接或解密步骤。
func (w *keystoreWallet) Open(passphrase string) error { return nil }

//close实现了账户。钱包，但它是普通钱包的一个noop，因为
//不是有意义的开放式操作。
func (w *keystoreWallet) Close() error { return nil }

//帐户实现帐户。钱包，返回包含
//普通的Kestore钱包中包含的单个帐户。
func (w *keystoreWallet) Accounts() []accounts.Account {
	return []accounts.Account{w.account}
}

//包含implements accounts.wallet，返回特定帐户是否为
//或未被此钱包实例包装。
func (w *keystoreWallet) Contains(account accounts.Account) bool {
	return account.Address == w.account.Address && (account.URL == (accounts.URL{}) || account.URL == w.account.URL)
}

//派生实现了accounts.wallet，但对于普通的钱包来说是一个noop，因为
//对于普通的密钥存储帐户，不存在分层帐户派生的概念。
func (w *keystoreWallet) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	return accounts.Account{}, accounts.ErrNotSupported
}

//Selfderive实现了accounts.wallet，但对于普通的钱包来说是一个noop，因为
//对于普通密钥库帐户，没有层次结构帐户派生的概念。
func (w *keystoreWallet) SelfDerive(base accounts.DerivationPath, chain ethereum.ChainStateReader) {}

//sign hash实现accounts.wallet，尝试用
//给定的帐户。如果钱包没有包裹这个特定的账户，
//返回错误以避免帐户泄漏（即使在理论上我们可能
//能够通过我们的共享密钥库后端进行签名）。
func (w *keystoreWallet) SignHash(account accounts.Account, hash []byte) ([]byte, error) {
//确保请求的帐户包含在
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
//帐户似乎有效，请求密钥库签名
	return w.keystore.SignHash(account, hash)
}

//signtx实现accounts.wallet，尝试签署给定的交易
//与给定的帐户。如果钱包没有包裹这个特定的账户，
//返回一个错误以避免帐户泄漏（即使在理论上我们可以
//能够通过我们的共享密钥库后端进行签名）。
func (w *keystoreWallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
//确保请求的帐户包含在
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
//帐户似乎有效，请求密钥库签名
	return w.keystore.SignTx(account, tx, chainID)
}

//signhashwithpassphrase实现accounts.wallet，尝试对
//使用密码短语作为额外身份验证的给定帐户的给定哈希。
func (w *keystoreWallet) SignHashWithPassphrase(account accounts.Account, passphrase string, hash []byte) ([]byte, error) {
//确保请求的帐户包含在
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
//帐户似乎有效，请求密钥库签名
	return w.keystore.SignHashWithPassphrase(account, passphrase, hash)
}

//signtxwithpassphrase实现accounts.wallet，尝试对给定的
//使用密码短语作为额外身份验证的给定帐户的事务。
func (w *keystoreWallet) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
//确保请求的帐户包含在
	if !w.Contains(account) {
		return nil, accounts.ErrUnknownAccount
	}
//帐户似乎有效，请求密钥库签名
	return w.keystore.SignTxWithPassphrase(account, passphrase, tx, chainID)
}

