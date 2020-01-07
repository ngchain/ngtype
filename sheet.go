package ngtype

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"math/big"
	"sort"
)

var (
	ErrAccountNotExists = errors.New("the account does not exist")
	ErrMalformedSheet   = errors.New("the sheet structure is malformed")
)

// NewSheet gets the rows from db and return the sheet for transport/saving
func NewSheet(prevVaultHash []byte, rows map[uint64]*Account) *Sheet {
	return &Sheet{
		PrevVaultHash: prevVaultHash,
		Rows:          rows,
	}
}

func NewEmptySheet(prevVaultHash []byte) *Sheet {
	rows := make(map[uint64]*Account)
	return &Sheet{
		PrevVaultHash: prevVaultHash,
		Rows:          rows,
	}
}

func (m *Sheet) ApplyBlock(block *Block) (*Sheet, error) {
	var err error
	var totalFee = Big0
	snapshot := m.Copy()

	for _, op := range block.Operations {
		if snapshot.Rows[op.From] == nil {
			err = errors.New("malformed op sender")
			return nil, err
		}

		value := new(big.Int).SetBytes(op.Value)
		fee := new(big.Int).SetBytes(op.Fee)
		totalExpense := new(big.Int).Add(fee, value)
		if new(big.Int).SetBytes(snapshot.Rows[op.From].Balance).Cmp(totalExpense) < 0 {
			err = errors.New("balance is not enough for op")
			return nil, err
		}

		totalFee = totalFee.Add(totalFee, fee)

		snapshot.Rows[op.From].Balance = new(big.Int).Sub(new(big.Int).SetBytes(snapshot.Rows[op.From].Balance), totalExpense).Bytes()

		snapshot.Rows[op.To].Balance = new(big.Int).Add(new(big.Int).SetBytes(snapshot.Rows[op.To].Balance), value).Bytes()
	}

	var accounts []*Account
	accounts, err = m.GetAccountByKeyBytes(block.Beneficiary)
	if len(accounts) == 0 {
		newAccount := NewAccount(new(big.Int).SetBytes(block.Hash).Uint64(), block.Beneficiary, totalFee, nil)
		snapshot.RegisterAccount(newAccount)
	} else {
		sort.Slice(accounts, func(i, j int) bool {
			return accounts[i].ID < accounts[j].ID
		})
		accounts[0].Balance = new(big.Int).Add(new(big.Int).SetBytes(accounts[0].Balance), totalFee).Bytes()
	}

	return snapshot, err
}

func (m *Sheet) RegisterAccount(account *Account) error {
	if m.Rows[account.ID] != nil {
		return errors.New("failed to register, account already exists")
	}

	m.Rows[account.ID] = account
	return nil
}

func (m *Sheet) GetAccountByID(accountID uint64) (*Account, error) {
	if !m.HasAccount(accountID) {
		return nil, errors.New("no such account")
	}

	return m.Rows[accountID], nil
}

func (m *Sheet) GetAccountByKey(publicKey ecdsa.PublicKey) ([]*Account, error) {
	accounts := make([]*Account, 0)
	bPublicKey := elliptic.Marshal(elliptic.P256(), publicKey.X, publicKey.Y)
	for i := range m.Rows {
		if bytes.Compare(m.Rows[i].Owner, bPublicKey) == 0 {
			accounts = append(accounts, m.Rows[i])
		}
	}

	return accounts, nil
}

func (m *Sheet) GetAccountByKeyBytes(bPublicKey []byte) ([]*Account, error) {
	accounts := make([]*Account, 0)
	for i := range m.Rows {
		if bytes.Compare(m.Rows[i].Owner, bPublicKey) == 0 {
			accounts = append(accounts, m.Rows[i])
		}
	}

	return accounts, nil
}

func (m *Sheet) HasAccount(accountID uint64) bool {
	return m.Rows[accountID] != nil
}

func (m *Sheet) DelAccount(accountID uint64) error {
	if !m.HasAccount(accountID) {
		return errors.New("no such account")
	}

	m.Rows[accountID] = nil
	return nil
}

func (m *Sheet) ExportAccounts() []*Account {
	accounts := make([]*Account, len(m.Rows))
	for i, row := range m.Rows {
		accounts[i] = row
	}
	return accounts
}

func (m *Sheet) Copy() *Sheet {
	s := *m
	return &s
}
