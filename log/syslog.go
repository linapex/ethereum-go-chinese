
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:39</date>
//</624450097945513984>

//+建设！窗户，！计划9

package log

import (
	"log/syslog"
	"strings"
)

//SyslogHandler通过调用
//并将所有记录写入其中。
func SyslogHandler(priority syslog.Priority, tag string, fmtr Format) (Handler, error) {
	wr, err := syslog.New(priority, tag)
	return sharedSyslog(fmtr, wr, err)
}

//syslognethandler通过网络打开与日志守护程序的连接并写入
//所有日志记录。
func SyslogNetHandler(net, addr string, priority syslog.Priority, tag string, fmtr Format) (Handler, error) {
	wr, err := syslog.Dial(net, addr, priority, tag)
	return sharedSyslog(fmtr, wr, err)
}

func sharedSyslog(fmtr Format, sysWr *syslog.Writer, err error) (Handler, error) {
	if err != nil {
		return nil, err
	}
	h := FuncHandler(func(r *Record) error {
		var syslogFn = sysWr.Info
		switch r.Lvl {
		case LvlCrit:
			syslogFn = sysWr.Crit
		case LvlError:
			syslogFn = sysWr.Err
		case LvlWarn:
			syslogFn = sysWr.Warning
		case LvlInfo:
			syslogFn = sysWr.Info
		case LvlDebug:
			syslogFn = sysWr.Debug
		case LvlTrace:
syslogFn = func(m string) error { return nil } //没有用于跟踪的系统日志级别
		}

		s := strings.TrimSpace(string(fmtr.Format(r)))
		return syslogFn(s)
	})
	return LazyHandler(&closingHandler{sysWr, h}), nil
}

func (m muster) SyslogHandler(priority syslog.Priority, tag string, fmtr Format) Handler {
	return must(SyslogHandler(priority, tag, fmtr))
}

func (m muster) SyslogNetHandler(net, addr string, priority syslog.Priority, tag string, fmtr Format) Handler {
	return must(SyslogNetHandler(net, addr, priority, tag, fmtr))
}

