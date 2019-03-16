
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450112411668480>


package api

import (
	"context"
	"testing"
)

func testStorage(t *testing.T, f func(*Storage, bool)) {
	testAPI(t, func(api *API, toEncrypt bool) {
		f(NewStorage(api), toEncrypt)
	})
}

func TestStoragePutGet(t *testing.T) {
	testStorage(t, func(api *Storage, toEncrypt bool) {
		content := "hello"
		exp := expResponse(content, "text/plain", 0)
//exp：=expresponse（[]byte（content），“text/plain”，0）
		ctx := context.TODO()
		bzzkey, wait, err := api.Put(ctx, content, exp.MimeType, toEncrypt)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		err = wait(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		bzzhash := bzzkey.Hex()
//根据API检查Put
		resp0 := testGet(t, api.api, bzzhash, "")
		checkResponse(t, resp0, exp)

//检查存储获取
		resp, err := api.Get(context.TODO(), bzzhash)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		checkResponse(t, &testResponse{nil, resp}, exp)
	})
}

