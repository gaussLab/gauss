package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/gauss/gauss/v4/x/token/types"
)

// InitGenesis initializes the token module's state from a given genesis state.
func (k BaseKeeper) InitGenesis(ctx sdk.Context, gs *types.GenesisState) {
	if err := gs.Validate(); err != nil {
		panic(err.Error())
	}

	k.SetParams(ctx, gs.Params)

	// init tokens
	for _, token := range gs.Tokens {
		if err := k.AddIssuedToken(ctx, token); err != nil {
			panic(err.Error())
		}
	}

	for _, coin := range gs.BurnedCoins {
		k.AddBurnedCoin(ctx, coin)
	}
}

// ExportGenesis returns the bank module's genesis state.
func (k BaseKeeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	var tokens []types.Token
	for _, token := range k.GetTokens(ctx, nil) {
		t := token.(*types.Token)
		tokens = append(tokens, *t)
	}

	return types.NewGenesisState(
		k.GetParams(ctx),
		tokens,
		k.GetAllBurntCoins(ctx),
	)
}
