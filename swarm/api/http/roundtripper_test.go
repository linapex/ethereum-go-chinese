
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:43</date>
//</624450111962877952>


package http

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRoundTripper(t *testing.T) {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "text/plain")
			http.ServeContent(w, r, "", time.Unix(0, 0), strings.NewReader(r.RequestURI))
		} else {
			http.Error(w, "Method "+r.Method+" is not supported.", http.StatusMethodNotAllowed)
		}
	})

	srv := httptest.NewServer(serveMux)
	defer srv.Close()

	host, port, _ := net.SplitHostPort(srv.Listener.Addr().String())
	rt := &RoundTripper{Host: host, Port: port}
	trans := &http.Transport{}
	trans.RegisterProtocol("bzz", rt)
	client := &http.Client{Transport: trans}
resp, err := client.Get("bzz://COM/PATH“测试”
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}
	if string(content) != "/HTTP/1.1:/test.com/path" {
		t.Errorf("incorrect response from http server: expected '%v', got '%v'", "/HTTP/1.1:/test.com/path", string(content))
	}

}

