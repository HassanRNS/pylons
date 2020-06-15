package types

import (
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Reference: BlockHeader struct is available at github.com/tendermint/tendermint@v0.33.4/types/proto3/block.pb.go:
// type Header struct {
// 		AppHash            []byte `protobuf:"bytes,11,opt,name=app_hash,json=appHash,proto3" json:"app_hash,omitempty"`
//   	...
// }

// RandomSeed calculate random seed from context and entity count
func RandomSeed(ctx sdk.Context, entityCount int) int64 {
	header := ctx.BlockHeader()
	appHash := header.AppHash
	seedValue := 0
	for i, bytv := range appHash { // len(appHash) = 11
		intv := int(bytv)
		seedValue += (i*i + 1) * intv
	}
	fmt.Println("RandomSeed entityCount:", entityCount, "BlockHeight:", header.Height)
	return int64(seedValue + entityCount)
}

// Reader struct is for entropy set on uuid
type Reader struct{}

// NewEntropyReader create an entropy reader
func NewEntropyReader() *Reader {
	return &Reader{}
}

func (r Reader) Read(b []byte) (n int, err error) {
	entropy := []byte{}
	for i := 0; i < len(b); i++ {
		entropy = append(entropy, byte(rand.Intn(256)))
	}

	n = copy(b, entropy)
	return n, nil
}