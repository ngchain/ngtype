package ngtype

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/NebulousLabs/fastrand"
	"github.com/ngin-network/secp256k1"
)

func TestNewAccount(t *testing.T) {
	privateKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		log.Error(err)
	}

	randUint64 := fastrand.Uint64n(math.MaxUint64)
	acc := NewAccount(
		randUint64,
		privateKey.PubKey().SerializeCompressed(),
		big.NewInt(0),
		nil,
	)
	fmt.Println(acc)
}
