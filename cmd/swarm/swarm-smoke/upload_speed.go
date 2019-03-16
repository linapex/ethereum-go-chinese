
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450071710142464>


package main

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/swarm/testutil"

	cli "gopkg.in/urfave/cli.v1"
)

var endpoint string

//只需使用第一个端点
func generateEndpoint(scheme string, cluster string, app string, from int) {
	if cluster == "prod" {
endpoint = fmt.Sprintf("%s://%v.swarm-gateways.net“，方案，来自）
	} else {
endpoint = fmt.Sprintf("%s://%s-%v-%s.stg.swarm gateways.net“，方案，应用程序，来自，群集）
	}
}

func cliUploadSpeed(c *cli.Context) error {
	log.PrintOrigins(true)
	log.Root().SetHandler(log.LvlFilterHandler(log.Lvl(verbosity), log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	metrics.GetOrRegisterCounter("upload-speed", nil).Inc(1)

	errc := make(chan error)
	go func() {
		errc <- uploadSpeed(c)
	}()

	select {
	case err := <-errc:
		if err != nil {
			metrics.GetOrRegisterCounter("upload-speed.fail", nil).Inc(1)
		}
		return err
	case <-time.After(time.Duration(timeout) * time.Second):
		metrics.GetOrRegisterCounter("upload-speed.timeout", nil).Inc(1)
		return fmt.Errorf("timeout after %v sec", timeout)
	}
}

func uploadSpeed(c *cli.Context) error {
	defer func(now time.Time) {
		totalTime := time.Since(now)

		log.Info("total time", "time", totalTime, "kb", filesize)
		metrics.GetOrRegisterCounter("upload-speed.total-time", nil).Inc(int64(totalTime))
	}(time.Now())

	generateEndpoint(scheme, cluster, appName, from)
	seed := int(time.Now().UnixNano() / 1e6)
	log.Info("uploading to "+endpoint, "seed", seed)

	randomBytes := testutil.RandomBytes(seed, filesize*1000)

	t1 := time.Now()
	hash, err := upload(&randomBytes, endpoint)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	metrics.GetOrRegisterCounter("upload-speed.upload-time", nil).Inc(int64(time.Since(t1)))

	fhash, err := digest(bytes.NewReader(randomBytes))
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("uploaded successfully", "hash", hash, "digest", fmt.Sprintf("%x", fhash))
	return nil
}

