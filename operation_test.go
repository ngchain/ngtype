package ngtype

import (
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
