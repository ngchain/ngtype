package ngtype

import (
	"math"
	"math/big"
	"time"
)

const (
	GenesisPK      = "8ZzCXNDkTvFvND72AopxwVHtseV8m42sAnKup8zU7jfDXjwKp2ESdYCLNxcWPPnAL8vg2K1Z9kHuy6MZ24BgtvWEaKjSYQER3JH6kryzyeunBrfJEWPJVkegCT6knE61F1zY"
	GenesisBalance = math.MaxInt64
	GenesisData    = "NGIN TESTNET"
	GenesisHash    = "123"
)

var (
	MaxDiff           = [32]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255} // new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)) // Target = MaxDiff / diff
	GenesisDifficulty = new(big.Int).SetBytes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255})
	GenesisNonce      = new(big.Int).SetUint64(5577006791947779410)
	Big1              = big.NewInt(1)

	TargetTime = 1e10

	genesisTimestamp = time.Date(2020, time.February, 2, 2, 2, 2, 2, time.UTC).Unix()
)

// Units
var (
	OneNGIN = new(big.Int).SetUint64(1000000)
	//MinimalUnit = big.NewInt(1)
	OneBlockReward = new(big.Int).Mul(OneNGIN, big.NewInt(10)) // 10NG
)
