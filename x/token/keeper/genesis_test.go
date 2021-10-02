package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gauss/gauss/v4/x/token/types"
	"github.com/gauss/gauss/v4/simapp"
)

func TestExportGenesis(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	
	// export genesis
	genesisState := token.ExportGenesis(ctx, app.TokenKeeper)
	
	require.Equal(t, types.DefaultParams(), genesisState.Params)
	for _, token := range genesisState.Tokens {
	}
}

func TestInitGenesis(t *testing.T) {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// add token
	addr := sdk.AccAddress(tmhash.SumTruncated([]byte("addrOne")))
	eth := types.NewToken("Ethereum Network", "eth", "eth", 1, 1, 1, true, addr)

	genesisState := types.NewGenesisState(
			Params:	types.DefaultParams(),
			Tokens: []types.Token{eth},
		)

	// initialize genesis
	app.TokenKeeper.InitGenesis(ctx, genesisState)

	// query all tokens
	var tokens = app.TokenKeeper.GetTokens(ctx, nil)
	require.Equal(t, len(tokens), 2)
	require.Equal(t, tokens[0], &eth)
}
