package simulation

// DONTCOVER

import (
	"bytes"
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/gauss/gauss/v4/x/token/types"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding token type
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.SymbolPrefix):
			var tokenA, tokenB types.Token
			cdc.MustUnmarshalBinaryBare(kvA.Value, &tokenA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &tokenB)
			return fmt.Sprintf("%v\n%v", tokenA, tokenB)
		case bytes.Equal(kvA.Key[:1], types.OwnerSymbolKey):
			var symbolA, symbolB gogotypes.Value
			cdc.MustUnmarshalBinaryBare(kvA.Value, &symbolA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &symbolB)
			return fmt.Sprintf("%v\n%v", symbolA, symbolB)
		default:
			panic(fmt.Sprintf("invalid %s key prefix %X", types.ModuleName, kvA.Key[:1]))
		}
	}
}
