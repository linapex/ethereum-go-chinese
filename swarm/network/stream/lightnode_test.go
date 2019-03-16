
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:44</date>
//</624450115519647744>

package stream

import (
	"testing"

	p2ptest "github.com/ethereum/go-ethereum/p2p/testing"
)

//此测试检查服务器的默认行为，即
//当它提供检索请求时。
func TestLigthnodeRetrieveRequestWithRetrieve(t *testing.T) {
	registryOptions := &RegistryOptions{
		Retrieval: RetrievalClientOnly,
		Syncing:   SyncingDisabled,
	}
	tester, _, _, teardown, err := newStreamerTester(t, registryOptions)
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	node := tester.Nodes[0]

	stream := NewStream(swarmChunkServerStreamName, "", false)

	err = tester.TestExchanges(p2ptest.Exchange{
		Label: "SubscribeMsg",
		Triggers: []p2ptest.Trigger{
			{
				Code: 4,
				Msg: &SubscribeMsg{
					Stream: stream,
				},
				Peer: node.ID(),
			},
		},
	})
	if err != nil {
		t.Fatalf("Got %v", err)
	}

	err = tester.TestDisconnected(&p2ptest.Disconnect{Peer: node.ID()})
	if err == nil || err.Error() != "timed out waiting for peers to disconnect" {
		t.Fatalf("Expected no disconnect, got %v", err)
	}
}

//此测试在服务检索时检查服务器的lightnode行为
//请求被禁用
func TestLigthnodeRetrieveRequestWithoutRetrieve(t *testing.T) {
	registryOptions := &RegistryOptions{
		Retrieval: RetrievalDisabled,
		Syncing:   SyncingDisabled,
	}
	tester, _, _, teardown, err := newStreamerTester(t, registryOptions)
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	node := tester.Nodes[0]

	stream := NewStream(swarmChunkServerStreamName, "", false)

	err = tester.TestExchanges(
		p2ptest.Exchange{
			Label: "SubscribeMsg",
			Triggers: []p2ptest.Trigger{
				{
					Code: 4,
					Msg: &SubscribeMsg{
						Stream: stream,
					},
					Peer: node.ID(),
				},
			},
			Expects: []p2ptest.Expect{
				{
					Code: 7,
					Msg: &SubscribeErrorMsg{
						Error: "stream RETRIEVE_REQUEST not registered",
					},
					Peer: node.ID(),
				},
			},
		})
	if err != nil {
		t.Fatalf("Got %v", err)
	}
}

//此测试检查服务器的默认行为，即
//启用同步时。
func TestLigthnodeRequestSubscriptionWithSync(t *testing.T) {
	registryOptions := &RegistryOptions{
		Retrieval: RetrievalDisabled,
		Syncing:   SyncingRegisterOnly,
	}
	tester, _, _, teardown, err := newStreamerTester(t, registryOptions)
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	node := tester.Nodes[0]

	syncStream := NewStream("SYNC", FormatSyncBinKey(1), false)

	err = tester.TestExchanges(
		p2ptest.Exchange{
			Label: "RequestSubscription",
			Triggers: []p2ptest.Trigger{
				{
					Code: 8,
					Msg: &RequestSubscriptionMsg{
						Stream: syncStream,
					},
					Peer: node.ID(),
				},
			},
			Expects: []p2ptest.Expect{
				{
					Code: 4,
					Msg: &SubscribeMsg{
						Stream: syncStream,
					},
					Peer: node.ID(),
				},
			},
		})

	if err != nil {
		t.Fatalf("Got %v", err)
	}
}

//此测试检查服务器的lightnode行为，即
//当同步被禁用时。
func TestLigthnodeRequestSubscriptionWithoutSync(t *testing.T) {
	registryOptions := &RegistryOptions{
		Retrieval: RetrievalDisabled,
		Syncing:   SyncingDisabled,
	}
	tester, _, _, teardown, err := newStreamerTester(t, registryOptions)
	defer teardown()
	if err != nil {
		t.Fatal(err)
	}

	node := tester.Nodes[0]

	syncStream := NewStream("SYNC", FormatSyncBinKey(1), false)

	err = tester.TestExchanges(p2ptest.Exchange{
		Label: "RequestSubscription",
		Triggers: []p2ptest.Trigger{
			{
				Code: 8,
				Msg: &RequestSubscriptionMsg{
					Stream: syncStream,
				},
				Peer: node.ID(),
			},
		},
		Expects: []p2ptest.Expect{
			{
				Code: 7,
				Msg: &SubscribeErrorMsg{
					Error: "stream SYNC not registered",
				},
				Peer: node.ID(),
			},
		},
	}, p2ptest.Exchange{
		Label: "RequestSubscription",
		Triggers: []p2ptest.Trigger{
			{
				Code: 4,
				Msg: &SubscribeMsg{
					Stream: syncStream,
				},
				Peer: node.ID(),
			},
		},
		Expects: []p2ptest.Expect{
			{
				Code: 7,
				Msg: &SubscribeErrorMsg{
					Error: "stream SYNC not registered",
				},
				Peer: node.ID(),
			},
		},
	})

	if err != nil {
		t.Fatalf("Got %v", err)
	}
}

