
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450071651422208>


package main

import (
	"bytes"
	"context"
	"crypto/md5"
	crand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	"os"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/swarm/api"
	"github.com/ethereum/go-ethereum/swarm/api/client"
	"github.com/ethereum/go-ethereum/swarm/spancontext"
	"github.com/ethereum/go-ethereum/swarm/testutil"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pborman/uuid"

	cli "gopkg.in/urfave/cli.v1"
)

func generateEndpoints(scheme string, cluster string, app string, from int, to int) {
	if cluster == "prod" {
		for port := from; port < to; port++ {
endpoints = append(endpoints, fmt.Sprintf("%s://%v.swarm-gateways.net“，方案，端口）
		}
	} else {
		for port := from; port < to; port++ {
endpoints = append(endpoints, fmt.Sprintf("%s://%s-%v-%s.stg.swarm gateways.net“，方案，应用程序，端口，群集）
		}
	}

	if includeLocalhost {
endpoints = append(endpoints, "http://本地主机：8500“）
	}
}

func cliUploadAndSync(c *cli.Context) error {
	log.PrintOrigins(true)
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(verbosity), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	metrics.GetOrRegisterCounter("upload-and-sync", nil).Inc(1)

	errc := make(chan error)
	go func() {
		errc <- uploadAndSync(c)
	}()

	select {
	case err := <-errc:
		if err != nil {
			metrics.GetOrRegisterCounter("upload-and-sync.fail", nil).Inc(1)
		}
		return err
	case <-time.After(time.Duration(timeout) * time.Second):
		metrics.GetOrRegisterCounter("upload-and-sync.timeout", nil).Inc(1)
		return fmt.Errorf("timeout after %v sec", timeout)
	}
}

func uploadAndSync(c *cli.Context) error {
	defer func(now time.Time) {
		totalTime := time.Since(now)
		log.Info("total time", "time", totalTime, "kb", filesize)
		metrics.GetOrRegisterResettingTimer("upload-and-sync.total-time", nil).Update(totalTime)
	}(time.Now())

	generateEndpoints(scheme, cluster, appName, from, to)
	seed := int(time.Now().UnixNano() / 1e6)
	log.Info("uploading to "+endpoints[0]+" and syncing", "seed", seed)

	randomBytes := testutil.RandomBytes(seed, filesize*1000)

	t1 := time.Now()
	hash, err := upload(&randomBytes, endpoints[0])
	if err != nil {
		log.Error(err.Error())
		return err
	}
	metrics.GetOrRegisterResettingTimer("upload-and-sync.upload-time", nil).UpdateSince(t1)

	fhash, err := digest(bytes.NewReader(randomBytes))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("uploaded successfully", "hash", hash, "digest", fmt.Sprintf("%x", fhash))

	time.Sleep(time.Duration(syncDelay) * time.Second)

	wg := sync.WaitGroup{}
	if single {
		rand.Seed(time.Now().UTC().UnixNano())
		randIndex := 1 + rand.Intn(len(endpoints)-1)
		ruid := uuid.New()[:8]
		wg.Add(1)
		go func(endpoint string, ruid string) {
			for {
				start := time.Now()
				err := fetch(hash, endpoint, fhash, ruid)
				if err != nil {
					continue
				}

				metrics.GetOrRegisterResettingTimer("upload-and-sync.single.fetch-time", nil).UpdateSince(start)
				wg.Done()
				return
			}
		}(endpoints[randIndex], ruid)
	} else {
		for _, endpoint := range endpoints[1:] {
			ruid := uuid.New()[:8]
			wg.Add(1)
			go func(endpoint string, ruid string) {
				for {
					start := time.Now()
					err := fetch(hash, endpoint, fhash, ruid)
					if err != nil {
						continue
					}

					metrics.GetOrRegisterResettingTimer("upload-and-sync.each.fetch-time", nil).UpdateSince(start)
					wg.Done()
					return
				}
			}(endpoint, ruid)
		}
	}
	wg.Wait()
	log.Info("all endpoints synced random file successfully")

	return nil
}

//
func fetch(hash string, endpoint string, original []byte, ruid string) error {
	ctx, sp := spancontext.StartSpan(context.Background(), "upload-and-sync.fetch")
	defer sp.Finish()

	log.Trace("http get request", "ruid", ruid, "api", endpoint, "hash", hash)

	var tn time.Time
	reqUri := endpoint + "/bzz:/" + hash + "/"
	req, _ := http.NewRequest("GET", reqUri, nil)

	opentracing.GlobalTracer().Inject(
		sp.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))

	trace := client.GetClientTrace("upload-and-sync - http get", "upload-and-sync", ruid, &tn)

	req = req.WithContext(httptrace.WithClientTrace(ctx, trace))
	transport := http.DefaultTransport

//transport.tlsclientconfig=&tls.config不安全的skipverify:true

	tn = time.Now()
	res, err := transport.RoundTrip(req)
	if err != nil {
		log.Error(err.Error(), "ruid", ruid)
		return err
	}
	log.Trace("http get response", "ruid", ruid, "api", endpoint, "hash", hash, "code", res.StatusCode, "len", res.ContentLength)

	if res.StatusCode != 200 {
		err := fmt.Errorf("expected status code %d, got %v", 200, res.StatusCode)
		log.Warn(err.Error(), "ruid", ruid)
		return err
	}

	defer res.Body.Close()

	rdigest, err := digest(res.Body)
	if err != nil {
		log.Warn(err.Error(), "ruid", ruid)
		return err
	}

	if !bytes.Equal(rdigest, original) {
		err := fmt.Errorf("downloaded imported file md5=%x is not the same as the generated one=%x", rdigest, original)
		log.Warn(err.Error(), "ruid", ruid)
		return err
	}

	log.Trace("downloaded file matches random file", "ruid", ruid, "len", res.ContentLength)

	return nil
}

//upload正在通过“swarm up”命令将文件“f”上载到“endpoint”
func upload(dataBytes *[]byte, endpoint string) (string, error) {
	swarm := client.NewClient(endpoint)
	f := &client.File{
		ReadCloser: ioutil.NopCloser(bytes.NewReader(*dataBytes)),
		ManifestEntry: api.ManifestEntry{
			ContentType: "text/plain",
			Mode:        0660,
			Size:        int64(len(*dataBytes)),
		},
	}

//将数据上载到bzz://并检索内容寻址清单哈希，十六进制编码。
	return swarm.Upload(f, "", false)
}

func digest(r io.Reader) ([]byte, error) {
	h := md5.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

//在堆缓冲区中生成随机数据
func generateRandomData(datasize int) ([]byte, error) {
	b := make([]byte, datasize)
	c, err := crand.Read(b)
	if err != nil {
		return nil, err
	} else if c != datasize {
		return nil, errors.New("short read")
	}
	return b, nil
}

