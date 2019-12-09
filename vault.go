package ngtype

import (
	"errors"
	"github.com/gogo/protobuf/proto"
	"golang.org/x/crypto/sha3"
	"math/big"
	"time"
)

func NewVault(newAccountID uint64, prevVault *Vault, hookBlock *Block, currentSheet *Sheet) *Vault {
	if !hookBlock.IsCheckpoint() {
		log.Error(errors.New("not proper time for building new vault"))
		return nil
	}

	if !hookBlock.IsSealed() {
		log.Error(ErrBlockIsUnsealing)
		return nil
	}

	prevVaultHash, err := prevVault.CalculateHash()
	if err != nil {
		log.Error(err)
	}
	hookBlockHash, err := hookBlock.CalculateHash()
	if err != nil {
		log.Error(err)
	}

	newAccount := NewRewardAccount(newAccountID, hookBlock.Beneficiary, big.NewInt(0))

	return &Vault{
		Height:        0,
		NewAccount:    newAccount,
		Timestamp:     time.Now().Unix(),
		PrevVaultHash: prevVaultHash,
		HookBlockHash: hookBlockHash,
		Sheet:         currentSheet,
	}
}

func (m *Vault) CalculateHash() ([]byte, error) {
	b, err := proto.Marshal(m)
	hash := sha3.Sum256(b)
	return hash[:], err
}

func GetGenesisVault() *Vault {
	var hookGenesisBlock = GetGenesisBlock()
	blockHash, _ := hookGenesisBlock.CalculateHash()
	return &Vault{
		Height:     0,
		Sheet:      nil,
		NewAccount: GetGenesisAccount(),

		Timestamp:     genesisTimestamp,
		PrevVaultHash: nil,
		HookBlockHash: blockHash,
	}
}
