
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:46</date>
//</624450125288181760>


//包含与Whisper协议信封对象关联的测试。

package whisperv6

import (
	mrand "math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestEnvelopeOpenAcceptsOnlyOneKeyTypeInFilter(t *testing.T) {
	symKey := make([]byte, aesKeyLength)
	mrand.Read(symKey)

	asymKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed GenerateKey with seed %d: %s.", seed, err)
	}

	params := MessageParams{
		PoW:      0.01,
		WorkTime: 1,
		TTL:      uint32(mrand.Intn(1024)),
		Payload:  make([]byte, 50),
		KeySym:   symKey,
		Dst:      nil,
	}

	mrand.Read(params.Payload)

	msg, err := NewSentMessage(&params)
	if err != nil {
		t.Fatalf("failed to create new message with seed %d: %s.", seed, err)
	}

	e, err := msg.Wrap(&params)
	if err != nil {
		t.Fatalf("Failed to Wrap the message in an envelope with seed %d: %s", seed, err)
	}

	f := Filter{KeySym: symKey, KeyAsym: asymKey}

	decrypted := e.Open(&f)
	if decrypted != nil {
		t.Fatalf("Managed to decrypt a message with an invalid filter, seed %d", seed)
	}
}

