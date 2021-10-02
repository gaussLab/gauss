package simulation_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/gauss/gauss/v4/x/orderbook/simulation"
	"github.com/gauss/gauss/v4/x/orderbook/types"
)

var (
	ownerPk1   = ed25519.GenPrivKey().PubKey()
	ownerAddr1 = sdk.AccAddress(ownerPk1.Address())
)

func makeTestCodec() (cdc *codec.LegacyAmino) {
	cdc = codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	cryptocodec.RegisterCrypto(cdc)
	types.RegisterLegacyAminoCodec(cdc)
	return
}

func TestDecodeStore(t *testing.T) {
	cdc, _ := simapp.MakeCodecs()
	dec := simulation.NewDecodeStore(cdc)

	bondTime := time.Now().UTC()

	pool := types.NewPool()
	order := types.NewOrder()
	txPairStats := types.NewTxPairStats()

	kvPairs := kv.Pairs{
		Pairs: []kv.Pair{
			{Key: types.PoolKey, Value: cdc.MustMarshalBinaryBare(&pool)},
			{Key: types.TxPairStatsKey, Value: cdc.MustMarshalBinaryBare(&txPairStats)},
			{Key: types.OrderbookKey, Value: cdc.MustMarshalBinaryBare(&order)},
		},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"pool", fmt.Sprintf("%v\n%v", pool, pool)},
		{"order", fmt.Sprintf("%v\n%v", order, order)},
		{"tps", fmt.Sprintf("%v\n%v", txPairStats, txPairStats)},
		{"other", ""},
	}
	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { dec(kvPairs.Pairs[i], kvPairs.Pairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, dec(kvPairs.Pairs[i], kvPairs.Pairs[i]), tt.name)
			}
		})
	}
}
