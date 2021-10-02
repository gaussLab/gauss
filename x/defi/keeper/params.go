package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// MintInflation
func (k Keeper) MintInflation(ctx sdk.Context) (res sdk.Coin) {
	k.paramstore.Get(ctx, types.KeyMintInflation, &res)
	return
}

// CommunityTax
func (k Keeper) CommunityTax(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyCommunityTax, &res)
	return
}

// CommisssionRate
func (k Keeper) CommissionRate(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyCommissionRate, &res)
	return
}

// CommisssionRate
func (k Keeper) MarketRate(ctx sdk.Context) (res sdk.Dec) {
	k.paramstore.Get(ctx, types.KeyMarketRate, &res)
	return
}

// UnbondingTime
func (k Keeper) UnbondingTime(ctx sdk.Context) (res time.Duration) {
	k.paramstore.Get(ctx, types.KeyUnbondingTime, &res)
	return
}

// MaxDefis - Maximum number of defis
func (k Keeper) MaxDefis(ctx sdk.Context) (res uint32) {
        k.paramstore.Get(ctx, types.KeyMaxDefis, &res)
        return
}

// MaxEntries - Maximum number of simultaneous unbonding
// delegations (per pair/trio)
func (k Keeper) MaxEntries(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyMaxEntries, &res)
	return
}

// HistoricalEntries = number of historical info entries
// to persist in store
func (k Keeper) HistoricalEntries(ctx sdk.Context) (res uint32) {
	k.paramstore.Get(ctx, types.KeyHistoricalEntries, &res)
	return
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx sdk.Context) (res string) {
	k.paramstore.Get(ctx, types.KeyBondDenom, &res)
	return
}

// DefiBondDenom 
func (k Keeper) DefiBondDenom(ctx sdk.Context) string {
	return k.MintInflation(ctx).Denom
}

// Get all parameteras as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
        return params
/*
	return types.NewParams(
		k.BondDenom(ctx),
		k.MintInflation(ctx),
		k.CommunityTax(ctx),
		k.CommissionRate(ctx),
		k.UnbondingTime(ctx),
		k.MaxDefis(ctx),
		k.MaxEntries(ctx),
		k.HistoricalEntries(ctx),
	)
*/
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
