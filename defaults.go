package ngtype

import (
	"math"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

const (
	MajorVersion = 0
	MinorVersion = 1
	PatchVersion = 0
)

const (
	GenesisPK      = "8ZzCXNDkTvFvND72AopxwVHtseV8m42sAnKup8zU7jfDXjwKp2ESdYCLNxcWPPnAL8vg2K1Z9kHuy6MZ24BgtvWEaKjSYQER3JH6kryzyeunBrfJEWPJVkegCT6knE61F1zY"
	GenesisBalance = math.MaxInt64
	GenesisData    = "NGIN TESTNET"
	GenesisHash    = "123"

	//Coinbase = "EVj8i3gn7CBvjowfpt8xYv6nifB7kepoq8MkxAwLbmUS"
)

var (
	//GenesisDifficulty = [32]byte{1, 134, 160, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0} // 100000
	MaxDiff           = [32]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255} // new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)) // Target = MaxDiff / diff
	GenesisDifficulty = new(big.Int).SetBytes([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255})
	GenesisNonce      = new(big.Int).SetUint64(5577006791947779410)
	Big1              = big.NewInt(1)

	TargetTime = 1e10

	genesisTimestamp = time.Date(2020, time.February, 2, 2, 2, 2, 2, time.UTC).Unix()
)

type Version struct {
	Major uint
	Minor uint
	Patch uint
}

// GetDefaultDataFolder will return the data dir for current user in string, use filepath.Join to operate the files
func GetDefaultDataFolder() string {
	dir, err := os.UserHomeDir() // warning: require go>=1.12
	if err != nil {
		log.Info(err)
		dir = os.TempDir()
	}

	dataDir := filepath.Join(dir, ".ngd")

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err = os.Mkdir(dataDir, os.ModePerm); err != nil {
			log.Error("failed to mkdir!")
		} else {
			log.Info("Local Data directory:", dataDir)
		}
	}

	return dataDir
}

// Units
var (
	OneNGIN = new(big.Int).SetUint64(1000000)
	//MinimalUnit = big.NewInt(1)
	OneBlockReward = new(big.Int).Mul(OneNGIN, big.NewInt(10)) // 10NG
)
