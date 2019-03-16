
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:42</date>
//</624450110104801280>


package rpc

import "testing"

func TestWSGetConfigNoAuth(t *testing.T) {
config, err := wsGetConfig("ws://示例.com:1234“，”）
	if err != nil {
		t.Logf("wsGetConfig failed: %s", err)
		t.Fail()
		return
	}
	if config.Location.User != nil {
		t.Log("User should have been stripped from the URL")
		t.Fail()
	}
	if config.Location.Hostname() != "example.com" ||
		config.Location.Port() != "1234" || config.Location.Scheme != "ws" {
		t.Logf("Unexpected URL: %s", config.Location)
		t.Fail()
	}
}

func TestWSGetConfigWithBasicAuth(t *testing.T) {
config, err := wsGetConfig("wss://testuser:test-pass_01@example.com:1234“，”）
	if err != nil {
		t.Logf("wsGetConfig failed: %s", err)
		t.Fail()
		return
	}
	if config.Location.User != nil {
		t.Log("User should have been stripped from the URL")
		t.Fail()
	}
	if config.Header.Get("Authorization") != "Basic dGVzdHVzZXI6dGVzdC1QQVNTXzAx" {
		t.Log("Basic auth header is incorrect")
		t.Fail()
	}
}

