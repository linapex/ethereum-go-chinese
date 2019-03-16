
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450073329143808>


//+构建FreeBSD

package fdlimit

import "syscall"

//此文件与fdlimiteux.go基本相同，
//但是rlimit字段在freebsd上有int64类型，因此它需要
//额外的转换。

//raise尝试最大化此进程的文件描述符允许量
//达到操作系统允许的最大硬限制。
func Raise(max uint64) error {
//获取当前限制
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return err
	}
//尝试将限额更新为最大限额
	limit.Cur = limit.Max
	if limit.Cur > int64(max) {
		limit.Cur = int64(max)
	}
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return err
	}
	return nil
}

//当前检索允许由此打开的文件描述符数
//过程。
func Current() (int, error) {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return 0, err
	}
	return int(limit.Cur), nil
}

//最大检索此进程的最大文件描述符数
//允许自己请求。
func Maximum() (int, error) {
	var limit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		return 0, err
	}
	return int(limit.Max), nil
}

