package keeper

import (
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	
	"github.com/gauss/gauss/v4/x/token/types"
)

// ViewKeeper defines a module interface that facilitates read only access to
// account balances.
type ViewKeeper interface {
	GetTokens(ctx sdk.Context, owner sdk.AccAddress) (tokens []types.TokenI)
	GetToken(ctx sdk.Context, denom string) (types.TokenI, error)
	HasToken(ctx sdk.Context, denom string) bool
	GetOwner(ctx sdk.Context, denom string) (sdk.AccAddress, error)
	GetBurntCoin(ctx sdk.Context, denom string) (sdk.Coin, bool)
        GetAllBurntCoins(ctx sdk.Context) sdk.Coins
	IsUnlocked(ctx sdk.Context, denom string) bool
}

var _ ViewKeeper = (*BaseViewKeeper)(nil)

// BaseViewKeeper implements a read only keeper implementation of ViewKeeper.
type BaseViewKeeper struct {
	cdc		codec.BinaryMarshaler
	storeKey	sdk.StoreKey
	bankKeeper	types.BankKeeper
}

// NewBaseViewKeeper returns a new BaseViewKeeper.
func NewBaseViewKeeper(cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, bk types.BankKeeper) BaseViewKeeper {
	return BaseViewKeeper{
		cdc:		cdc,
		storeKey:	storeKey,
		bankKeeper:	bk,
	}
}

// Logger returns a module-specific logger.
func (k BaseViewKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// HasToken asserts a token exists
func (k BaseViewKeeper) HasToken(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetSymbolKey(denom))
}

// HasToken asserts a token exists
func (k BaseViewKeeper) HasTokenWithUnit(ctx sdk.Context, denom string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetUnitKey(denom))
}

// GetToken returns the token of the specified symbol
func (k BaseViewKeeper) GetToken(ctx sdk.Context, symbol string) (types.TokenI, error) {
	token, err := k.getTokenBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetToken returns the token of the specified symbol
func (k BaseViewKeeper) GetTokenWithUnit(ctx sdk.Context, unit string) (types.TokenI, error) {
	symbol, err := k.getSymbolByUnit(ctx, unit)
	if err != nil {
		return nil, err
	}
	return k.GetToken(ctx, symbol)
}


// GetTokens returns all existing tokens
func (k BaseViewKeeper) GetTokens(ctx sdk.Context, owner sdk.AccAddress) (tokens []types.TokenI) {
	store := ctx.KVStore(k.storeKey)

	if owner == nil {
		return k.getAllTokens(store)
	}

	return k.getAllTokenOfOwner(ctx, store, owner)
}

func (k BaseViewKeeper) GetOwner(ctx sdk.Context, symbol string) (sdk.AccAddress, error) {
	token, err := k.GetToken(ctx, symbol)
	if err != nil {
		return nil, err
	}
	return token.GetOwner(), nil
}

// use smallest-unit
func (k BaseViewKeeper) GetBurntCoin(ctx sdk.Context, denom string) (coin sdk.Coin, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetBurntCoinKey(denom)

	bz := store.Get(key)
	if len(bz) == 0 {
		return coin, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &coin)
	return coin, true
}

func (k BaseViewKeeper) GetAllBurntCoins(ctx sdk.Context) sdk.Coins {
	var coins sdk.Coins

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.BurntCoinPrefix)
	for ; iter.Valid(); iter.Next() {
		var coin sdk.Coin
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &coin)
		coins = append(coins, coin)
	}
	
	return coins
}

func (k BaseViewKeeper) getTokenBySymbol(ctx sdk.Context, symbol string) (token types.Token, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetSymbolKey(symbol))
	if bz == nil {
		return token, sdkerrors.Wrap(types.ErrTokenNotExists, fmt.Sprintf("Token[%s] does not exist", symbol))
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &token)
	return token, nil
}

func (k BaseViewKeeper) getSymbolByUnit(ctx sdk.Context, unit string) (symbol string, err error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetUnitKey(unit))
	if bz == nil {
		return symbol, sdkerrors.Wrap(types.ErrTokenNotExists, fmt.Sprintf("Token[%s] does not exist", unit))
	}

	var symbolL gogotypes.StringValue	
	k.cdc.MustUnmarshalBinaryBare(bz, &symbolL)
	return symbolL.Value, nil
}


func (k BaseViewKeeper) getAllTokens(store sdk.KVStore) (tokens []types.TokenI) {

	iter := sdk.KVStorePrefixIterator(store, types.SymbolPrefix)
	defer iter.Close()
	
	for ; iter.Valid(); iter.Next() {
		var token types.Token
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &token)

		tokens = append(tokens, &token)
	}

	return
}

func (k BaseViewKeeper) getAllTokenOfOwner(
	ctx sdk.Context, store sdk.KVStore, owner sdk.AccAddress) (tokens []types.TokenI) {

	iter := sdk.KVStorePrefixIterator(store, types.GetOwnerSymbolKey(owner, ""))
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var symbol gogotypes.StringValue
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &symbol)		
	
		token, err := k.getTokenBySymbol(ctx, symbol.Value)
		if err != nil {
			continue
		}
		tokens = append(tokens, token)	
	}

	return
}

// use smallest-unit
func (k BaseViewKeeper) IsUnlocked(ctx sdk.Context, denom string) bool {
	return k.bankKeeper.SendEnabledCoin(ctx, sdk.Coin{Denom:denom, Amount: sdk.ZeroInt()})
}

// getTokenSupply queries the token supply from the total supply
func (k BaseViewKeeper) getTokenSupply(ctx sdk.Context, denom string) sdk.Int {
	return k.bankKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
}
