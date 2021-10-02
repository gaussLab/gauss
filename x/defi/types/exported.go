package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DelegationI delegation bond for a delegated proof of defi system
type DelegationI interface {
	GetDelegatorAddr() sdk.AccAddress // delegator sdk.AccAddress for the bond
	GetDefiAddr()      sdk.ValAddress // defi operator address
	GetShares() sdk.Dec               // amount of defi's shares held in this delegation
}

// DefiI expected defi functions
type DefiI interface {
	GetMoniker() string                                     // moniker of the defi
	GetStatus() BondStatus                                  // status of the defi
	IsBonded() bool                                         // check if has a bonded status
	IsUnbonded() bool                                       // check if has status unbonded
	IsUnbonding() bool                                      // check if has status unbonding
	GetOperator() sdk.ValAddress                            // operator address to receive/return defis coins
	GetTokens() sdk.Int                                     // validation tokens
	GetBondedTokens() sdk.Int                               // defi bonded tokens
	GetMinSelfDelegation() sdk.Int                          // defi minimum self delegation
	GetDelegatorShares() sdk.Dec                            // total outstanding delegator shares
	TokensFromShares(sdk.Dec) sdk.Dec                       // token worth of provided delegator shares
	TokensFromSharesTruncated(sdk.Dec) sdk.Dec              // token worth of provided delegator shares, truncated
	TokensFromSharesRoundUp(sdk.Dec) sdk.Dec                // token worth of provided delegator shares, rounded up
	SharesFromTokens(amt sdk.Int) (sdk.Dec, error)          // shares worth of delegator's bond
	SharesFromTokensTruncated(amt sdk.Int) (sdk.Dec, error) // truncated shares worth of delegator's bond
}
