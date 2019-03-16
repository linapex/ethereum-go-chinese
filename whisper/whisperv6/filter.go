
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:46</date>
//</624450125330124800>


package whisperv6

import (
	"crypto/ecdsa"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

//筛选器表示耳语消息筛选器
type Filter struct {
Src        *ecdsa.PublicKey  //邮件的发件人
KeyAsym    *ecdsa.PrivateKey //收件人的私钥
KeySym     []byte            //与主题关联的键
Topics     [][]byte          //筛选邮件的主题
PoW        float64           //耳语规范中所述的工作证明
AllowP2P   bool              //指示此筛选器是否对直接对等消息感兴趣
SymKeyHash common.Hash       //优化所需的对称密钥的keccak256hash
id         string            //唯一标识符

	Messages map[common.Hash]*ReceivedMessage
	mutex    sync.RWMutex
}

//筛选器表示筛选器的集合
type Filters struct {
	watchers map[string]*Filter

topicMatcher     map[TopicType]map[*Filter]struct{} //将主题映射到在消息与该主题匹配时感兴趣收到通知的筛选器
allTopicsMatcher map[*Filter]struct{}               //列出将通知新邮件的所有筛选器，无论其主题是什么

	whisper *Whisper
	mutex   sync.RWMutex
}

//newfilters返回新创建的筛选器集合
func NewFilters(w *Whisper) *Filters {
	return &Filters{
		watchers:         make(map[string]*Filter),
		topicMatcher:     make(map[TopicType]map[*Filter]struct{}),
		allTopicsMatcher: make(map[*Filter]struct{}),
		whisper:          w,
	}
}

//安装将向筛选器集合添加新筛选器
func (fs *Filters) Install(watcher *Filter) (string, error) {
	if watcher.KeySym != nil && watcher.KeyAsym != nil {
		return "", fmt.Errorf("filters must choose between symmetric and asymmetric keys")
	}

	if watcher.Messages == nil {
		watcher.Messages = make(map[common.Hash]*ReceivedMessage)
	}

	id, err := GenerateRandomID()
	if err != nil {
		return "", err
	}

	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	if fs.watchers[id] != nil {
		return "", fmt.Errorf("failed to generate unique ID")
	}

	if watcher.expectsSymmetricEncryption() {
		watcher.SymKeyHash = crypto.Keccak256Hash(watcher.KeySym)
	}

	watcher.id = id
	fs.watchers[id] = watcher
	fs.addTopicMatcher(watcher)
	return id, err
}

//卸载将删除已从中指定ID的筛选器
//筛选器集合
func (fs *Filters) Uninstall(id string) bool {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	if fs.watchers[id] != nil {
		fs.removeFromTopicMatchers(fs.watchers[id])
		delete(fs.watchers, id)
		return true
	}
	return false
}

//addTopicMatcher向主题匹配器添加一个筛选器。
//如果过滤器的主题数组为空，将对每个主题进行尝试。
//否则，将在指定的主题上进行尝试。
func (fs *Filters) addTopicMatcher(watcher *Filter) {
	if len(watcher.Topics) == 0 {
		fs.allTopicsMatcher[watcher] = struct{}{}
	} else {
		for _, t := range watcher.Topics {
			topic := BytesToTopic(t)
			if fs.topicMatcher[topic] == nil {
				fs.topicMatcher[topic] = make(map[*Filter]struct{})
			}
			fs.topicMatcher[topic][watcher] = struct{}{}
		}
	}
}

//removeFromTopicMatchers从主题匹配器中删除筛选器
func (fs *Filters) removeFromTopicMatchers(watcher *Filter) {
	delete(fs.allTopicsMatcher, watcher)
	for _, topic := range watcher.Topics {
		delete(fs.topicMatcher[BytesToTopic(topic)], watcher)
	}
}

//GetWatchersByTopic返回一个包含以下筛选器的切片：
//匹配特定主题
func (fs *Filters) getWatchersByTopic(topic TopicType) []*Filter {
	res := make([]*Filter, 0, len(fs.allTopicsMatcher))
	for watcher := range fs.allTopicsMatcher {
		res = append(res, watcher)
	}
	for watcher := range fs.topicMatcher[topic] {
		res = append(res, watcher)
	}
	return res
}

//get返回集合中具有特定ID的筛选器
func (fs *Filters) Get(id string) *Filter {
	fs.mutex.RLock()
	defer fs.mutex.RUnlock()
	return fs.watchers[id]
}

//通知观察程序通知已声明感兴趣的任何筛选器
//信封主题。
func (fs *Filters) NotifyWatchers(env *Envelope, p2pMessage bool) {
	var msg *ReceivedMessage

	fs.mutex.RLock()
	defer fs.mutex.RUnlock()

	candidates := fs.getWatchersByTopic(env.Topic)
	for _, watcher := range candidates {
		if p2pMessage && !watcher.AllowP2P {
			log.Trace(fmt.Sprintf("msg [%x], filter [%s]: p2p messages are not allowed", env.Hash(), watcher.id))
			continue
		}

		var match bool
		if msg != nil {
			match = watcher.MatchMessage(msg)
		} else {
			match = watcher.MatchEnvelope(env)
			if match {
				msg = env.Open(watcher)
				if msg == nil {
					log.Trace("processing message: failed to open", "message", env.Hash().Hex(), "filter", watcher.id)
				}
			} else {
				log.Trace("processing message: does not match", "message", env.Hash().Hex(), "filter", watcher.id)
			}
		}

		if match && msg != nil {
			log.Trace("processing message: decrypted", "hash", env.Hash().Hex())
			if watcher.Src == nil || IsPubKeyEqual(msg.Src, watcher.Src) {
				watcher.Trigger(msg)
			}
		}
	}
}

func (f *Filter) expectsAsymmetricEncryption() bool {
	return f.KeyAsym != nil
}

func (f *Filter) expectsSymmetricEncryption() bool {
	return f.KeySym != nil
}

//触发器将未知消息添加到筛选器的列表中
//收到的消息。
func (f *Filter) Trigger(msg *ReceivedMessage) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if _, exist := f.Messages[msg.EnvelopeHash]; !exist {
		f.Messages[msg.EnvelopeHash] = msg
	}
}

//检索将返回所有相关联的已接收消息的列表
//过滤器。
func (f *Filter) Retrieve() (all []*ReceivedMessage) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	all = make([]*ReceivedMessage, 0, len(f.Messages))
	for _, msg := range f.Messages {
		all = append(all, msg)
	}

f.Messages = make(map[common.Hash]*ReceivedMessage) //删除旧邮件
	return all
}

//MatchMessage检查筛选器是否匹配已解密的
//消息（即已经由
//匹配前一个筛选器选中的信封）。
//这里不检查主题，因为这是由主题匹配器完成的。
func (f *Filter) MatchMessage(msg *ReceivedMessage) bool {
	if f.PoW > 0 && msg.PoW < f.PoW {
		return false
	}

	if f.expectsAsymmetricEncryption() && msg.isAsymmetricEncryption() {
		return IsPubKeyEqual(&f.KeyAsym.PublicKey, msg.Dst)
	} else if f.expectsSymmetricEncryption() && msg.isSymmetricEncryption() {
		return f.SymKeyHash == msg.SymKeyHash
	}
	return false
}

//匹配信封检查是否值得解密消息。如果
//它返回“true”，要求客户端代码尝试解密
//然后调用matchmessage。
//这里不检查主题，因为这是由主题匹配器完成的。
func (f *Filter) MatchEnvelope(envelope *Envelope) bool {
	return f.PoW <= 0 || envelope.pow >= f.PoW
}

//isubkeyEqual检查两个公钥是否相等
func IsPubKeyEqual(a, b *ecdsa.PublicKey) bool {
	if !ValidatePublicKey(a) {
		return false
	} else if !ValidatePublicKey(b) {
		return false
	}
//曲线总是一样的，只要比较点
	return a.X.Cmp(b.X) == 0 && a.Y.Cmp(b.Y) == 0
}

