//nolint
package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/token module sentinel errors
var (
	ErrInvalidName          = sdkerrors.Register(ModuleName, 2, "invalid token name")
        ErrInvalidSymbol        = sdkerrors.Register(ModuleName, 3, "invalid standard denom")
        ErrInvalidInitSupply    = sdkerrors.Register(ModuleName, 4, "invalid token initial supply")
	ErrInvalidTotalSupply   = sdkerrors.Register(ModuleName, 5, "invalid token maximum supply")
        ErrInvalidDecimals      = sdkerrors.Register(ModuleName, 6, "invalid token decimals")
        ErrSymbolAlreadyExists  = sdkerrors.Register(ModuleName, 7, "symbol already exists")
        ErrTokenNotExists       = sdkerrors.Register(ModuleName, 8, "token does not exist")
        ErrInvalidToAddress     = sdkerrors.Register(ModuleName, 9, "the new owner must not be same as the original owner")
        ErrInvalidOwner         = sdkerrors.Register(ModuleName, 10, "invalid token owner")
        ErrNotMintable          = sdkerrors.Register(ModuleName, 11, "token is not mintable")
        ErrInvalidAmount        = sdkerrors.Register(ModuleName, 12, "invalid amount")
        ErrInvalidIssueFee      = sdkerrors.Register(ModuleName, 13, "invalid issue token fee")
	ErrUnlockedToken	= sdkerrors.Register(ModuleName, 14, "token has been unlocked")
	ErrNotFoundToken	= sdkerrors.Register(ModuleName, 15, "token is not found")
)
