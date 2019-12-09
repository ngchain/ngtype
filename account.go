package ngtype

import (
	"math"
	"math/big"

	"github.com/mr-tron/base58"
)

// AccountID will receive an ID and PK then return a Account without SubState and Balance(0
func NewAccount(id uint64, ownerKey []byte, balance *big.Int, subState []byte) *Account {
	return &Account{
		ID:       id,
		Owner:    ownerKey,
		Balance:  balance.Bytes(),
		SubState: subState,
	}
}

func NewRewardAccount(id uint64, ownerKey []byte, totalFeeReward *big.Int) *Account {
	reward := new(big.Int).Add(OneBlockReward, totalFeeReward)
	return NewAccount(id, ownerKey, reward, nil)
}

//
func GetGenesisAccount() *Account {
	pk, _ := base58.FastBase58Decoding(GenesisPK)

	genesisAccount := &Account{
		ID:       1,
		Balance:  big.NewInt(math.MaxInt64).Bytes(), // Init balance
		Owner:    pk,
		Nonce:    0,
		SubState: []byte(`{'name':'NGIN OFFICIAL'}`),
	}

	return genesisAccount
}
