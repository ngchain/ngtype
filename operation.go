package ngtype

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/gogo/protobuf/proto"
	"golang.org/x/crypto/sha3"

	"github.com/cbergoon/merkletree"
	"github.com/mr-tron/base58"
)

var (
	ErrInvalidOpNonce      = errors.New("the nonce in operation is smaller than the account's record")
	ErrIsNotSigned         = errors.New("the operation is not signed")
	ErrBalanceInsufficient = errors.New("balance is insufficient for payment")
	ErrWrongSign           = errors.New("the signer of operation is not the own of the account")
	ErrMalformedOperation  = errors.New("the operation structure is malformed")
)

// Sign will re-sign the Op with private key
func (m *Operation) Signature(privKey *ecdsa.PrivateKey) (R, S *big.Int, err error) {
	b, err := proto.Marshal(m)
	if err != nil {
		log.Error(err)
	}

	R, S, err = ecdsa.Sign(rand.Reader, privKey, b)
	if err != nil {
		log.Panic(err)
	}

	m.R = R.Bytes()
	m.S = S.Bytes()

	return
}

// IsSigned will return whether the op has been signed
func (m *Operation) IsSigned() bool {
	if m.R == nil || m.S == nil {
		return false
	}
	return true
}

// Verify helps verify the operation whether signed by the public key owner
func (m *Operation) Verify(pubKey ecdsa.PublicKey) bool {
	if m.R == nil || m.S == nil {
		log.Panic("unsigned operation")
	}

	o := m.Copy()
	o.R = nil
	o.S = nil

	b, err := proto.Marshal(o)
	if err != nil {
		log.Error(err)
	}

	return ecdsa.Verify(&pubKey, b, new(big.Int).SetBytes(m.R), new(big.Int).SetBytes(m.S))
}

// ReadableID = txs in string
func (m *Operation) ReadableHex() string {
	b, err := proto.Marshal(m)
	if err != nil {
		log.Error(err)
	}
	return base58.FastBase58Encoding(b)
}

// CalculateHash mainly for calculating the tire root of ops
func (m *Operation) CalculateHash() ([]byte, error) {
	b, err := proto.Marshal(m)
	if err != nil {
		log.Error(err)
	}
	hash := sha3.Sum256(b)
	return hash[:], nil
}

// Equals mainly for calculating the tire root of ops
func (m *Operation) Equals(other merkletree.Content) (bool, error) {
	var equal = true
	o, ok := other.(*Operation)
	if !ok {
		return false, errors.New("invalid operation type")
	}

	equal = m.Type == o.Type
	equal = bytes.Compare(m.PrevVaultHash, o.PrevVaultHash) == 0
	equal = bytes.Compare(m.R, o.R) == 0
	equal = bytes.Compare(m.S, o.S) == 0
	equal = bytes.Compare(m.Value, o.Value) == 0
	equal = bytes.Compare(m.Fee, o.Fee) == 0
	equal = bytes.Compare(m.Extra, o.Extra) == 0
	equal = m.From == o.From
	equal = m.To == o.To
	equal = m.Nonce == o.Nonce

	return equal, nil
}

func (m *Operation) Copy() *Operation {
	o := *m
	return &o
}

// NewUnsignedOperation will return an Unsigned Operation, must using Signature()
func NewUnsignedOperation(t OpType, sender, target, n uint64, value, fee *big.Int, prevVaultHash, extraData []byte) *Operation {
	op := &Operation{
		Type:  t,
		From:  sender,
		To:    target,
		Nonce: n,
		Value: value.Bytes(),
		Fee:   fee.Bytes(),

		PrevVaultHash: prevVaultHash,
		Extra:         extraData,
	}

	return op
}

// TotalFee is a helper which helps calc the total fee among the ops
func TotalFee(ops []*Operation) (totalFee *big.Int) {
	totalFee = big.NewInt(0)
	for _, op := range ops {
		totalFee = new(big.Int).Add(totalFee, new(big.Int).SetBytes(op.Fee))
	}

	return
}

// Operations is an advanced type, aiming to get the trie root hash
type Operations struct {
	Ops []*Operation
}

func NewOperations(ops []*Operation) *Operations {
	return &Operations{
		Ops: ops,
	}
}

func (ops *Operations) Append(op *Operation) {
	ops.Ops = append(ops.Ops, op)
}

func (ops *Operations) Del(op *Operation) bool {
	for i := 0; i < len(ops.Ops); i++ {
		if ops.Ops[i] == op {
			ops.Ops = append(ops.Ops[:i], ops.Ops[i+1:]...)
			return true
		}
	}
	return false
}

func (ops *Operations) TrieRoot() []byte {
	var list []merkletree.Content
	for i := 0; i < len(ops.Ops); i++ {
		if ops.Ops[i] != nil {
			list = append(list, ops.Ops[i])
		}
	}
	trie, err := merkletree.NewTree(list)
	if err != nil {
		log.Error(err)
	}
	return trie.MerkleRoot()
}
