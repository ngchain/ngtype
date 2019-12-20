// block is just a tx mass, but can be treat as a block
// it means that after the vault summary, the data in previous block chain can be throw
// we just need to keep the latest some block and treasuries to make the chain safe
package ngtype

import (
	"bytes"
	"errors"
	"math/big"
	"time"

	"github.com/gogo/protobuf/proto"

	"github.com/mr-tron/base58"
	"github.com/ngin-network/cryptonight-go"
	"github.com/whyrusleeping/go-logging"
)

var (
	ErrBlockIsBare           = errors.New("the block is bare")
	ErrBlockIsUnsealing      = errors.New("the block is unsealing")
	ErrBlockHeightInvalid    = errors.New("the block's height is invalid")
	ErrBlockMTreeInvalid     = errors.New("the merkle tree in block is invalid")
	ErrBlockPrevBlockHash    = errors.New("the block's previous block hash is invalid")
	ErrBlockPrevTreasuryHash = errors.New("the block's backend vault is invalid")
	ErrBlockDiffInvalid      = errors.New("the block's difficulty is invalid")
	ErrBlockHashInvalid      = errors.New("the block's hash is invalid")
	ErrBlockNonceInvalid     = errors.New("the block's Nonce is invalid")
)

var log = logging.MustGetLogger("block")

// IsCheckpoint will check whether the Block is the
func (m *Block) IsCheckpoint() bool {
	return m.Height%8 == 0
}

func (m *Block) IsUnsealing() bool {
	return m.TrieHash == nil
}

func (m *Block) IsSealed() bool {
	return m.Nonce != nil
}

func (m *Block) IsGenesisBlock() bool {
	return m.Height == 0
}

// ToUnsealing converts a bare block to an unsealing block
func (m *Block) ToUnsealing(ops []*Operation) *Block {
	b := m.Copy()
	ops = append(ops, b.Operations...)
	mTreeHash := NewOperations(ops).TrieRoot
	copy(b.TrieHash, mTreeHash)

	return b
}

// ToUnsealing converts an unsealing block to a sealed block
func (m *Block) ToSealed(nonce *big.Int, hash []byte) *Block {
	b := m.Copy()
	if b.IsUnsealing() == false {
		log.Error(ErrBlockIsBare)
		return nil
	}
	b.Nonce = nonce.Bytes()
	return b
}

// GetBlob will return a complete blob for block hash
// = PoWBlob + Nonce
func (m *Block) GetBlob() ([]byte, error) {
	b := m.Copy()
	if b.Nonce == nil {
		return nil, errors.New("the block need to get Nonce (mined) before getting blob")
	}

	b.Hash = nil // empty the hash prop
	return proto.Marshal(b)
}

// GetHash will help you get the hash of block
func (m *Block) CalculateHash() ([]byte, error) {
	blob, err := m.GetBlob()
	if err != nil {
		return nil, err
	}

	if m.Nonce == nil {
		return nil, errors.New("missing the Nonce")
	}

	m.Hash = cryptonight.Sum(blob, 0)
	return m.Hash, nil
}

// GetHash will help you get the hash of block
func (m *Block) VerifyHash() bool {
	blob, _ := m.GetBlob()

	return bytes.Compare(m.Hash, cryptonight.Sum(blob, 0)) == 0
}

// NewBareBlock will return an unsealing block and
// then you need to add ops and seal with the correct Nonce
func NewBareBlock(height uint64, prevBlockHash, prevVaultHash []byte, PublicKey []byte, diff *big.Int) *Block {
	block := &Block{
		Height:        height,
		Beneficiary:   PublicKey[:],
		Timestamp:     time.Now().Unix(),
		PrevBlockHash: prevBlockHash[:],
		PrevVaultHash: prevVaultHash[:],

		//Empty params
		//MTreeHash:     nil,
		//Ops:   nil,
		//Nonce: nil,
		//Hash: nil,

		Difficulty: diff.Bytes(),
	}

	return block
}

// GetGenesisBlock will return a complete sealed GenesisBlock
func GetGenesisBlock() *Block {
	bPK, _ := base58.FastBase58Decoding(GenesisPK)

	b := &Block{
		Height:        0,
		Timestamp:     1024,
		TrieHash:      NewOperations(nil).TrieRoot,
		PrevBlockHash: nil,
		PrevVaultHash: nil,
		Beneficiary:   bPK,
		Nonce:         GenesisNonce.Bytes(),
		Difficulty:    GenesisDifficulty.Bytes(),
	}

	return b
}

// CheckError will check the errors in block inner fields
func (m *Block) CheckError() error {
	if m.Nonce == nil {
		return ErrBlockNonceInvalid
	}

	mTreeHash := NewOperations(m.Operations).TrieRoot

	if bytes.Compare(mTreeHash[:], m.TrieHash) != 0 {
		return ErrBlockMTreeInvalid
	}

	return nil
}

func (m *Block) Copy() *Block {
	b := *m
	return &b
}
