
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450098172006400>


package metrics

//disk stats是每个进程的磁盘IO状态。
type DiskStats struct {
ReadCount  int64 //执行的读取操作数
ReadBytes  int64 //读取的字节总数
WriteCount int64 //执行的写入操作数
WriteBytes int64 //写入的字节总数
}

