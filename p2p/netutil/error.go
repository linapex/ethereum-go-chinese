
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:41</date>
//</624450105226825728>


package netutil

//IsTemporaryError检查给定的错误是否应被视为临时错误。
func IsTemporaryError(err error) bool {
	tempErr, ok := err.(interface {
		Temporary() bool
	})
	return ok && tempErr.Temporary() || isPacketTooBig(err)
}

