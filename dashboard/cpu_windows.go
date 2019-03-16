
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:37</date>
//</624450087254233088>


package dashboard

//GetProcessCPutime在Windows上返回0，因为没有要解析的系统调用
//实际进程的CPU时间。
func getProcessCPUTime() float64 {
	return 0
}

