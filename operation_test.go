package ngtype

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/gogo/protobuf/proto"
	"math/big"
	"testing"

	"github.com/mr-tron/base58"
)

func TestDeserialize(t *testing.T) {
	op := &Operation{
		Type:  0,
		From:  0,
		To:    1,
		Nonce: 0,
		Value: new(big.Int).Exp(OneNGIN, big.NewInt(1000), nil).Bytes(),
		//Fee:   big.NewInt(100),

		PrevVaultHash: nil,
	}

	raw, _ := proto.Marshal(op)

	result := base58.FastBase58Encoding(raw)
	t.Log(result)

	_ = proto.Unmarshal(raw, op)
	t.Log(op)
}

func TestOperation_Signature(t *testing.T) {
	o := NewUnsignedOperation(OpType_TX, 1, 2, 0, big.NewInt(0), big.NewInt(0), nil, nil)
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	_, _, _ = o.Signature(priv)
	if !o.Verify(priv.PublicKey) {
		t.Fail()
	}
}
