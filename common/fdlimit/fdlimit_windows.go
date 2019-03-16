
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450073459167232>


package fdlimit

import "errors"

//raise尝试最大化此进程的文件描述符允许量
//达到操作系统允许的最大硬限制。
func Raise(max uint64) error {
//该方法设计为NOP：
//*Linux/Darwin对应程序需要手动增加每个进程的限制
//*在Windows上，Go使用CreateFile API，该API仅限于16K个文件，非
//可从正在运行的进程中更改
//这样，我们就可以“请求”提高限额，这两种情况都有
//或者基于我们运行的平台没有效果。
	if max > 16384 {
		return errors.New("file descriptor limit (16384) reached")
	}
	return nil
}

//当前检索允许由此打开的文件描述符数
//过程。
func Current() (int, error) {
//请参阅“加薪”，了解我们为什么使用硬编码16K作为限制的原因。
	return 16384, nil
}

//最大检索此进程的最大文件描述符数
//允许自己请求。
func Maximum() (int, error) {
	return Current()
}

