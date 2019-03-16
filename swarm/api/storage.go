
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450112373919744>


package api

import (
	"context"
	"path"

	"github.com/ethereum/go-ethereum/swarm/storage"
)

type Response struct {
	MimeType string
	Status   int
	Size     int64
//内容[]字节
	Content string
}

//实现服务
//
//已弃用：请改用HTTP API
type Storage struct {
	api *API
}

func NewStorage(api *API) *Storage {
	return &Storage{api}
}

//将内容上传到群中，并提供一个简单的清单
//其内容类型
//
//已弃用：请改用HTTP API
func (s *Storage) Put(ctx context.Context, content string, contentType string, toEncrypt bool) (storage.Address, func(context.Context) error, error) {
	return s.api.Put(ctx, content, contentType, toEncrypt)
}

//get从bzzpath检索内容并完全读取响应
//它返回响应对象，该对象将包含
//响应正文作为内容字段的值
//注：如果错误为非零，则响应可能仍有部分内容
//实际大小以len（resp.content）表示，而预期大小
//尺寸为相应尺寸
//
//已弃用：请改用HTTP API
func (s *Storage) Get(ctx context.Context, bzzpath string) (*Response, error) {
	uri, err := Parse(path.Join("bzz:/", bzzpath))
	if err != nil {
		return nil, err
	}
	addr, err := s.api.Resolve(ctx, uri.Addr)
	if err != nil {
		return nil, err
	}
	reader, mimeType, status, _, err := s.api.Get(ctx, nil, addr, uri.Path)
	if err != nil {
		return nil, err
	}
	quitC := make(chan bool)
	expsize, err := reader.Size(ctx, quitC)
	if err != nil {
		return nil, err
	}
	body := make([]byte, expsize)
	size, err := reader.Read(body)
	if int64(size) == expsize {
		err = nil
	}
	return &Response{mimeType, status, expsize, string(body[:size])}, err
}

