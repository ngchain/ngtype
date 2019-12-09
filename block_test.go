package ngtype

import (
	"bytes"
	"fmt"
	"github.com/NebulousLabs/fastrand"
	"github.com/gogo/protobuf/proto"
	"math/big"
	"runtime"
	"testing"
	"time"

	"github.com/mr-tron/base58"
	"github.com/ngin-network/cryptonight-go"
)

func TestBlock_GetHash(t *testing.T) {
	b := GetGenesisBlock()
	block, _ := b.CalculateHash()
	t.Log(len(block))
}

func TestGetGenesisBlockNonce(t *testing.T) {
	// new genesisBlock
	runtime.GOMAXPROCS(3)

	//coinBase := operation.GetCoinbaseOperation()
	//ops := []*operation.Operation{coinBase}
	pk, _ := base58.FastBase58Decoding(GenesisPK)
	//genesisMTree := [32]byte{} //operation.GetMTreeHash(ops)

	b := NewBareBlock(0, nil, nil, pk, GenesisDifficulty)
	b = b.ToUnsealing(nil)

	max := new(big.Int).SetBytes(MaxDiff[:])
	diff := new(big.Int).SetBytes(b.Difficulty)
	genesisTarget := new(big.Int).Div(max, diff)

	nCh := make(chan []byte, 1)
	stopCh := make(chan struct{}, 1)
	thread := 3
	for i := 0; i < thread; i++ {
		go calcHash(i, b, genesisTarget, nCh, stopCh)
	}

	select {
	case answer := <-nCh:
		stopCh <- struct{}{}
		b.Nonce = answer
		blob, err := b.GetBlob()
		if err != nil {
			log.Panic(err)
		}
		hash := cryptonight.Sum(blob, 0)
		fmt.Println("Nonce is ", answer, " Hash is ", base58.FastBase58Encoding(hash))
	}
}

func calcHash(id int, b *Block, target *big.Int, answerCh chan []byte, stopCh chan struct{}) {
	fmt.Println("thread ", id, " running")
	fmt.Println("target is ", target.String())

	t := time.Now()
	for {
		select {
		case <-stopCh:
			return
		default:
			random := fastrand.Bytes(8)
			b.Nonce = random
			blob, err := b.GetBlob()
			if err != nil {
				log.Panic(err)
			}
			hash := cryptonight.Sum(blob, 0)
			//fmt.Println(new(big.Int).SetBytes(hash).Uint64())
			if new(big.Int).SetBytes(hash).Cmp(target) < 0 {
				answerCh <- random
				fmt.Println("Found ", random, hash)
				elapsed := time.Since(t)
				fmt.Println("Elapsed: ", elapsed)
				return
			}
		}
	}
}

func TestBlock_Marshal(t *testing.T) {
	block, _ := GetGenesisBlock().Marshal()

	var _GBlock Block
	_ = proto.Unmarshal(block, &_GBlock)
	_block, _ := _GBlock.Marshal()
	if bytes.Compare(block, _block) != 0 {
		t.Fail()
	}
}
