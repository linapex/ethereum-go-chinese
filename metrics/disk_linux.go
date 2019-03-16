
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450098213949440>


//包含进程磁盘IO计数器检索的Linux实现。

package metrics

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

//readdiskstats检索属于当前进程的磁盘IO状态。
func ReadDiskStats(stats *DiskStats) error {
//打开进程磁盘IO计数器文件
	inf, err := os.Open(fmt.Sprintf("/proc/%d/io", os.Getpid()))
	if err != nil {
		return err
	}
	defer inf.Close()
	in := bufio.NewReader(inf)

//迭代IO计数器，提取我们需要的
	for {
//读取下一行并拆分为键和值
		line, err := in.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			return err
		}

//根据键更新计数器
		switch key {
		case "syscr":
			stats.ReadCount = value
		case "syscw":
			stats.WriteCount = value
		case "rchar":
			stats.ReadBytes = value
		case "wchar":
			stats.WriteBytes = value
		}
	}
}

