package ngtype

// NewSheet gets the rows from db and return the sheet for transport/saving
func NewSheet(prevVaultHash []byte, rows map[uint64]*Account) *Sheet {
	return &Sheet{
		PrevVaultHash: prevVaultHash,
		Rows:          rows,
	}
}

func GetEmptySheet(prevVaultHash []byte) *Sheet {
	rows := make(map[uint64]*Account)
	return &Sheet{
		PrevVaultHash: prevVaultHash,
		Rows:          rows,
	}
}

func (m *Sheet) Copy() *Sheet {
	s := *m
	return &s
}
