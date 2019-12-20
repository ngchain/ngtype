package ngtype

import (
	"errors"
	"math/big"
	"time"

	"github.com/gogo/protobuf/proto"
	"golang.org/x/crypto/sha3"
)

var (
	ErrNotCheckpoint    = errors.New("not proper time for building new vault")
	ErrInvalidHookBlock = errors.New("the vault's hook_block is invalid")
	ErrMalformedVault   = errors.New("the vault structure is malformed")
)

func NewVault(newAccountID uint64, prevVault *Vault, hookBlock *Block, currentSheet *Sheet) *Vault {
	if !hookBlock.IsCheckpoint() {
		log.Error(ErrNotCheckpoint)
		return nil
	}

	if !hookBlock.IsSealed() {
		log.Error(ErrBlockIsUnsealing)
		return nil
	}

	if !hookBlock.VerifyHash() {
		log.Error(ErrInvalidHookBlock)
	}

	prevVaultHash, err := prevVault.CalculateHash()
	if err != nil {
		log.Error(err)
	}

	newAccount := NewRewardAccount(newAccountID, hookBlock.Beneficiary, big.NewInt(0))

	v := &Vault{
		Height:        0,
		NewAccount:    newAccount,
		Timestamp:     time.Now().Unix(),
		PrevVaultHash: prevVaultHash,
		HookBlockHash: hookBlock.Hash,
		Sheet:         currentSheet,
	}

	v.Hash, err = v.CalculateHash()
	if err != nil {
		log.Error(err)
	}

	return v
}

func (m *Vault) CalculateHash() ([]byte, error) {
	v := m.Copy()
	v.Hash = nil
	b, err := proto.Marshal(v)
	hash := sha3.Sum256(b)

	m.Hash = hash[:]
	return hash[:], err
}

func GetGenesisVault() *Vault {
	var hookGenesisBlock = GetGenesisBlock()
	blockHash, err := hookGenesisBlock.CalculateHash()
	if err != nil {
		log.Error(err)
	}

	v := &Vault{
		Height:     0,
		Sheet:      nil,
		NewAccount: GetGenesisAccount(),

		Timestamp:     genesisTimestamp,
		PrevVaultHash: nil,
		HookBlockHash: blockHash,
	}

	v.Hash, err = v.CalculateHash()
	if err != nil {
		log.Error(err)
	}

	return v
}

func (m *Vault) Copy() *Vault {
	v := *m
	return &v
}
