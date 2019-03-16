
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450072393814016>


package utils

import (
	"os"
	"os/user"
	"testing"
)

func TestPathExpansion(t *testing.T) {
	user, _ := user.Current()
	tests := map[string]string{
		"/home/someuser/tmp": "/home/someuser/tmp",
		"~/tmp":              user.HomeDir + "/tmp",
		"~thisOtherUser/b/":  "~thisOtherUser/b",
		"$DDDXXX/a/b":        "/tmp/a/b",
		"/a/b/":              "/a/b",
	}
	os.Setenv("DDDXXX", "/tmp")
	for test, expected := range tests {
		got := expandPath(test)
		if got != expected {
			t.Errorf("test %s, got %s, expected %s\n", test, got, expected)
		}
	}
}

