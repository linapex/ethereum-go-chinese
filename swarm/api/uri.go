
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450112499748865>


package api

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/swarm/storage"
)

//匹配十六进制群哈希
//托多：这很糟糕，不应该硬编码哈希值有多长
var hashMatcher = regexp.MustCompile("^([0-9A-Fa-f]{64})([0-9A-Fa-f]{64})?$")

//URI是对存储在Swarm中的内容的引用。
type URI struct {
//方案具有以下值之一：
//
//*BZZ-群清单中的条目
//*BZZ原始-原始群内容
//*BZZ不可变-群清单中某个条目的不可变URI
//（地址未解析）
//*BZZ列表-包含在Swarm清单中的所有文件的列表
//
	Scheme string

//addr是十六进制存储地址，或者是
//解析为存储地址
	Addr string

//addr存储解析的存储地址
	addr storage.Address

//路径是群清单中内容的路径
	Path string
}

func (u *URI) MarshalJSON() (out []byte, err error) {
	return []byte(`"` + u.String() + `"`), nil
}

func (u *URI) UnmarshalJSON(value []byte) error {
	uri, err := Parse(string(value))
	if err != nil {
		return err
	}
	*u = *uri
	return nil
}

//解析将rawuri解析为一个uri结构，其中rawuri应该有一个
//以下格式：
//
//＊方案>：
//*<scheme>：/<addr>
//*<scheme>：/<addr>/<path>
//＊方案>：
//*<scheme>：/<addr>
//*<scheme>：/<addr>/<path>
//
//使用方案一：bzz、bzz raw、bzz immutable、bzz list或bzz hash
func Parse(rawuri string) (*URI, error) {
	u, err := url.Parse(rawuri)
	if err != nil {
		return nil, err
	}
	uri := &URI{Scheme: u.Scheme}

//检查方案是否有效
	switch uri.Scheme {
	case "bzz", "bzz-raw", "bzz-immutable", "bzz-list", "bzz-hash", "bzz-feed":
	default:
		return nil, fmt.Errorf("unknown scheme %q", u.Scheme)
	}

//处理类似bzz://<addr>/<path>的uri，其中addr和path
//已按URL拆分。分析
	if u.Host != "" {
		uri.Addr = u.Host
		uri.Path = strings.TrimLeft(u.Path, "/")
		return uri, nil
	}

//uri类似于bzz:/<addr>/<path>so split the addr and path from
//原始路径（将是/<addr>/<path>）
	parts := strings.SplitN(strings.TrimLeft(u.Path, "/"), "/", 2)
	uri.Addr = parts[0]
	if len(parts) == 2 {
		uri.Path = parts[1]
	}
	return uri, nil
}
func (u *URI) Feed() bool {
	return u.Scheme == "bzz-feed"
}

func (u *URI) Raw() bool {
	return u.Scheme == "bzz-raw"
}

func (u *URI) Immutable() bool {
	return u.Scheme == "bzz-immutable"
}

func (u *URI) List() bool {
	return u.Scheme == "bzz-list"
}

func (u *URI) Hash() bool {
	return u.Scheme == "bzz-hash"
}

func (u *URI) String() string {
	return u.Scheme + ":/" + u.Addr + "/" + u.Path
}

func (u *URI) Address() storage.Address {
	if u.addr != nil {
		return u.addr
	}
	if hashMatcher.MatchString(u.Addr) {
		u.addr = common.Hex2Bytes(u.Addr)
		return u.addr
	}
	return nil
}

