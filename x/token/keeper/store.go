package keeper

import (
	gogotypes "github.com/gogo/protobuf/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gauss/gauss/v4/x/token/types"
)

// storeToken add token with symbol
func (k BaseSendKeeper) storeToken(ctx sdk.Context, token types.Token) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&token)
	store.Set(types.GetSymbolKey(token.GetSymbol()), bz)
}

// storeTokenWithUnit add symbol with unit
func (k BaseSendKeeper) storeTokenWithUnit(ctx sdk.Context, unit, symbol string) {
	store := ctx.KVStore(k.storeKey)
	
	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: symbol})
	store.Set(types.GetUnitKey(unit), bz)
}


// storeTokenWithOwner add symbol with owner
func (k BaseSendKeeper) storeTokenWithOwner(ctx sdk.Context, owner sdk.AccAddress, symbol string) {
	store := ctx.KVStore(k.storeKey)
	
	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: symbol})
	store.Set(types.GetOwnerSymbolKey(owner, symbol), bz)
}

// storeBurntCoin
func (k BaseSendKeeper) storeBurntCoin(ctx sdk.Context, coin sdk.Coin) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&coin)
	store.Set(types.GetBurntCoinKey(coin.Denom), bz)
}

// reset all indices by the new owner for token query
func (k BaseSendKeeper) resetTokenOwner(ctx sdk.Context, symbol string, oldOwner, newOwner sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)

	// delete the old key
	store.Delete(types.GetOwnerSymbolKey(oldOwner, symbol))
	// add the new key
	k.storeTokenWithOwner(ctx, newOwner, symbol)
}
