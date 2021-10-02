package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/gauss/gauss/v4/x/orderbook/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding staking type.
func NewDecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.PoolKey):
			var poolA, poolB types.Pool

			cdc.MustUnmarshalBinaryBare(kvA.Value, &poolA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &poolB)

			return fmt.Sprintf("%v\n%v", poolA, poolB)
		case bytes.Equal(kvA.Key[:1], types.TxPairStatsKey):
			var txPairStatsA, txPairStatsB types.TxPairStats

			cdc.MustUnmarshalBinaryBare(kvA.Value, &txPairStatsA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &txPairStatsB)

			return fmt.Sprintf("%v\n%v", txPairStatsA, txPairStatsB)
		case bytes.Equal(kvA.Key[:1], types.OrderbookKey):
			var orderA, orderB types.Order

			cdc.MustUnmarshalBinaryBare(kvA.Value, &orderA)
			cdc.MustUnmarshalBinaryBare(kvB.Value, &orderB)

			return fmt.Sprintf("%v\n%v", orderA, orderB)
		default:
			panic(fmt.Sprintf("invalid orderbook key prefix %X", kvA.Key[:1]))
		}
	}
}
