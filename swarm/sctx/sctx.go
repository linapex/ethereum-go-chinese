
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450117272866816>

package sctx

import "context"

type (
	HTTPRequestIDKey struct{}
	requestHostKey   struct{}
)

func SetHost(ctx context.Context, domain string) context.Context {
	return context.WithValue(ctx, requestHostKey{}, domain)
}

func GetHost(ctx context.Context) string {
	v, ok := ctx.Value(requestHostKey{}).(string)
	if ok {
		return v
	}
	return ""
}

