package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Validate performs basic validation of supply genesis data returning an
// error for any failed validation criteria.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	
	// validate token
	for _, token := range gs.Tokens {
		if err := token.Validate(); err != nil {
			return err
		}
	}

	for _, coin := range gs.BurnedCoins {
		if err := coin.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, tokens []Token, burntCoins sdk.Coins) *GenesisState {
	return &GenesisState{
		Params:	params,
		Tokens:	tokens,
		BurnedCoins: burntCoins,
	}
}

// DefaultGenesisState returns a default bank module genesis state.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), []Token{}, sdk.Coins{})
}


