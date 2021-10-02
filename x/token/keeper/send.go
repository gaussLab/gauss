package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/gauss/gauss/v4/x/token/types"
)

// SendKeeper defines a module interface that facilitates the transfer of coins
// between accounts without the possibility of creating coins.

type SendKeeper interface {
	ViewKeeper

	AddToken(ctx sdk.Context, token types.Token) error
	AddIssuedToken(ctx sdk.Context, token types.Token) error
	AddBurnedCoin(ctx sdk.Context, coin sdk.Coin)

	GetParams(ctx sdk.Context) types.Params
	SetParams(ctx sdk.Context, params types.Params)
}

var _ SendKeeper = (*BaseSendKeeper)(nil)

// BaseSendKeeper only allows transfers between accounts without the possibility of
// creating coins. It implements the SendKeeper interface.
type BaseSendKeeper struct {
	BaseViewKeeper

	cdc		codec.BinaryMarshaler
	storeKey	sdk.StoreKey
	bankKeeper	types.BankKeeper
	paramSpace	paramtypes.Subspace
}

func NewBaseSendKeeper(
	cdc codec.BinaryMarshaler, storeKey sdk.StoreKey, bk types.BankKeeper, paramSpace paramtypes.Subspace,
) BaseSendKeeper {
	return BaseSendKeeper{
		BaseViewKeeper: NewBaseViewKeeper(cdc, storeKey, bk),
		cdc:		cdc,
		bankKeeper:	bk,
		storeKey:	storeKey,
		paramSpace:	paramSpace,
	}
}

// GetParams returns the total set of bank parameters.
func (k BaseSendKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of bank parameters.
func (k BaseSendKeeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// AddToken adds amt to the account balance given by the provided address. An
// error is returned if the initial amount is invalid or if any resulting new
// balance is negative.
func (k BaseSendKeeper) AddToken(ctx sdk.Context, token types.Token) error {
	if k.hasToken(ctx, token) {
		return sdkerrors.Wrapf(types.ErrSymbolAlreadyExists,
			"Token[%s-%s] already exists", token.GetSymbol(), token.GetSmallestUnit())
	}

	return k.addToken(ctx, token)
}
func (k BaseSendKeeper) AddIssuedToken(ctx sdk.Context, token types.Token) error {
	if k.hasIssuedToken(ctx, token) {
		return sdkerrors.Wrapf(types.ErrSymbolAlreadyExists,
			"Token[%s-%s] already exists", token.GetSymbol(), token.GetSmallestUnit())
	}
	return k.addToken(ctx, token)
}

func (k BaseSendKeeper) addToken(ctx sdk.Context, token types.Token) error {

	k.storeToken(ctx, token)
	k.storeTokenWithUnit(ctx, token.GetSmallestUnit(), token.GetSymbol())
	if len(token.Owner) != 0 {
		k.storeTokenWithOwner(ctx, token.GetOwner(), token.GetSymbol())
	}

	denomMetaData := banktypes.Metadata{
		Description: token.GetName(),
		Base: token.GetSmallestUnit(),
		Display: token.GetSymbol(),
		DenomUnits: []*banktypes.DenomUnit{
			{Denom: token.GetSmallestUnit(), Exponent: 0},
			{Denom: token.GetSymbol(), Exponent: token.GetDecimals()},
		},
	}
	k.bankKeeper.SetDenomMetaData(ctx, denomMetaData)

	return nil
}

func (k BaseSendKeeper) hasToken(ctx sdk.Context, token types.Token) bool {
	if k.hasIssuedToken(ctx, token) {
		return true
	}

	// Low efficiency?
	issuedTokens := k.bankKeeper.GetSupply(ctx).GetTotal()
	if issuedTokens.AmountOf(token.GetSymbol()).GT(sdk.ZeroInt()) {
		return true
	}
	if issuedTokens.AmountOf(token.GetSmallestUnit()).GT(sdk.ZeroInt()) {
		return true
	}

	return false
}

func (k BaseSendKeeper) hasIssuedToken(ctx sdk.Context, token types.Token) bool {
        if k.HasToken(ctx, token.GetSymbol()) {
		return true
	}
	if k.HasTokenWithUnit(ctx, token.GetSmallestUnit()) {
		return true
	}

	return false
}

// use smallest-unit
// AddBurnCoin saves the total amount of the burned tokens
func (k BaseSendKeeper) AddBurnedCoin(ctx sdk.Context, coin sdk.Coin) {
	var coinL = coin
	if hasCoin, found := k.GetBurntCoin(ctx, coin.Denom); found {
		coinL = coinL.Add(hasCoin)
	}

	k.storeBurntCoin(ctx, coinL)
}

// mintCoinsToAccount
func (k BaseSendKeeper) mintCoinsToAccount(ctx sdk.Context, to sdk.AccAddress, coins sdk.Coins) error {
	// mint coins into module account
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	// sent coins to account
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, coins)
}
// mintCoinsToModule
func (k BaseSendKeeper) mintCoinsToModule(ctx sdk.Context, to string, coins sdk.Coins) error {
	// mint coins into module account
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
		return err
	}

	// sent coins to account
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, to, coins)
}


// burnCoins
func (k BaseSendKeeper) burnCoins(ctx sdk.Context, from sdk.AccAddress, coins sdk.Coins) error {
	// burn coins
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, coins); err != nil {
		return err
	}
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, coins)
}

// unlockToken unlock the specialied token
func (k BaseSendKeeper) unlockToken(ctx sdk.Context, symbol string, unlocked bool){
	bkParams := k.bankKeeper.GetParams(ctx);
	k.bankKeeper.SetParams(ctx, bkParams.SetSendEnabledParam(symbol, unlocked))
}
