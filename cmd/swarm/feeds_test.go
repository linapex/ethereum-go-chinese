
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450070820950016>


package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/swarm/api"
	swarm "github.com/ethereum/go-ethereum/swarm/api/client"
	swarmhttp "github.com/ethereum/go-ethereum/swarm/api/http"
	"github.com/ethereum/go-ethereum/swarm/storage/feed"
	"github.com/ethereum/go-ethereum/swarm/storage/feed/lookup"
	"github.com/ethereum/go-ethereum/swarm/testutil"
)

func TestCLIFeedUpdate(t *testing.T) {
	srv := swarmhttp.NewTestSwarmServer(t, func(api *api.API) swarmhttp.TestServer {
		return swarmhttp.NewServer(api, "")
	}, nil)
	log.Info("starting a test swarm server")
	defer srv.Close()

//
	privkeyHex := "0000000000000000000000000000000000000000000000000000000000001979"
	privKey, _ := crypto.HexToECDSA(privkeyHex)
	address := crypto.PubkeyToAddress(privKey.PublicKey)

	pkFileName := testutil.TempFileWithContent(t, privkeyHex)
	defer os.Remove(pkFileName)

//撰写主题。我们要引用米格尔·德·塞万提斯的话
	var topic feed.Topic
	subject := []byte("Miguel de Cervantes")
	copy(topic[:], subject[:])
	name := "quotes"

//为更新准备一些数据
	data := []byte("En boca cerrada no entran moscas")
	hexData := hexutil.Encode(data)

	flags := []string{
		"--bzzapi", srv.URL,
		"--bzzaccount", pkFileName,
		"feed", "update",
		"--topic", topic.Hex(),
		"--name", name,
		hexData}

//创建一个更新并期待一个没有错误的退出
	log.Info("updating a feed with 'swarm feed update'")
	cmd := runSwarm(t, flags...)
	cmd.ExpectExit()

//现在尝试使用客户端获取更新
	client := swarm.NewClient(srv.URL)

//创建与以前相同的主题，这次
//我们使用newtopic自动创建主题。
	topic, err := feed.NewTopic(name, subject)
	if err != nil {
		t.Fatal(err)
	}

//feed配置我们要查找的更新。
	fd := feed.Feed{
		Topic: topic,
		User:  address,
	}

//生成查询以获取最新更新
	query := feed.NewQueryLatest(&fd, lookup.NoClue)

//检索内容！
	reader, err := client.QueryFeed(query, "")
	if err != nil {
		t.Fatal(err)
	}

	retrieved, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

//检查我们是否检索了发送的信息
	if !bytes.Equal(data, retrieved) {
		t.Fatalf("Received %s, expected %s", retrieved, data)
	}

//现在检索下一次更新的信息
	flags = []string{
		"--bzzapi", srv.URL,
		"feed", "info",
		"--topic", topic.Hex(),
		"--user", address.Hex(),
	}

	log.Info("getting feed info with 'swarm feed info'")
	cmd = runSwarm(t, flags...)
_, matches := cmd.ExpectRegexp(`.*`) //regex hack提取stdout
	cmd.ExpectExit()

//验证是否可以将结果反序列化为有效的JSON
	var request feed.Request
	err = json.Unmarshal([]byte(matches[0]), &request)
	if err != nil {
		t.Fatal(err)
	}

//
	if request.Feed != fd {
		t.Fatalf("Expected feed to be: %s, got %s", fd, request.Feed)
	}

//
	flags = []string{
		"--bzzapi", srv.URL,
		"--bzzaccount", pkFileName,
		"feed", "create",
		"--topic", topic.Hex(),
	}

	log.Info("Publishing manifest with 'swarm feed create'")
	cmd = runSwarm(t, flags...)
	_, matches = cmd.ExpectRegexp(`[a-f\d]{64}`)
	cmd.ExpectExit()

manifestAddress := matches[0] //

//
	reader, err = client.QueryFeed(nil, manifestAddress)
	if err != nil {
		t.Fatal(err)
	}

	retrieved, err = ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(data, retrieved) {
		t.Fatalf("Received %s, expected %s", retrieved, data)
	}

//
	flags = []string{
		"--bzzapi", srv.URL,
		"feed", "create",
		"--topic", topic.Hex(),
"--user", "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", //不同用户
	}

	log.Info("Publishing manifest with 'swarm feed create' for a different user")
	cmd = runSwarm(t, flags...)
	_, matches = cmd.ExpectRegexp(`[a-f\d]{64}`)
	cmd.ExpectExit()

manifestAddress = matches[0] //

//
	flags = []string{
		"--bzzapi", srv.URL,
		"--bzzaccount", pkFileName,
		"feed", "update",
		"--manifest", manifestAddress,
		hexData}

//
	log.Info("updating a feed with 'swarm feed update'")
	cmd = runSwarm(t, flags...)
cmd.ExpectRegexp("Fatal:.*") //目前检测故障的最佳方法。
	cmd.ExpectExit()
	if cmd.ExitStatus() == 0 {
		t.Fatal("Expected nonzero exit code when updating a manifest with the wrong user. Got 0.")
	}
}

