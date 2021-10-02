package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/defi module sentinel errors
//
// TODO: Many of these errors are redundant. They should be removed and replaced
// by sdkerrors.ErrInvalidRequest.
//
// REF: https://github.com/cosmos/cosmos-sdk/issues/5450
var (
	ErrEmptyDefiAddr                   = sdkerrors.Register(ModuleName, 1, "empty defi address")
	ErrBadDefiAddr                     = sdkerrors.Register(ModuleName, 2, "defi address is invalid")
	ErrNoDefiFound                     = sdkerrors.Register(ModuleName, 3, "defi does not exist")
	ErrDefiOwnerExists                 = sdkerrors.Register(ModuleName, 4, "defi already exist for this operator address; must use new defi operator address")
	ErrSelfDelegationBelowMinimum      = sdkerrors.Register(ModuleName, 5, "defi's self delegation must be greater than their minimum self delegation")
	ErrMinSelfDelegationInvalid        = sdkerrors.Register(ModuleName, 6, "minimum self delegation must be a positive integer")
	ErrMinSelfDelegationDecreased      = sdkerrors.Register(ModuleName, 7, "minimum self delegation cannot be decrease")
	ErrEmptyDelegatorAddr              = sdkerrors.Register(ModuleName, 8, "empty delegator address")
	ErrBadDenom                        = sdkerrors.Register(ModuleName, 9, "invalid coin denomination")
	ErrBadDelegatorAddr                = sdkerrors.Register(ModuleName, 10, "delegator does not exist with address")
	ErrBadDelegationAmount             = sdkerrors.Register(ModuleName, 11, "invalid delegation amount")
	ErrNoDelegation                    = sdkerrors.Register(ModuleName, 12, "no delegation for (address, defi) tuple")
	ErrNoDelegatorForAddress           = sdkerrors.Register(ModuleName, 13, "delegator does not contain delegation")
	ErrInsufficientShares              = sdkerrors.Register(ModuleName, 14, "insufficient delegation shares")
	ErrDelegationDefiEmpty             = sdkerrors.Register(ModuleName, 15, "cannot delegate to an empty defi")
	ErrNotEnoughDelegationShares       = sdkerrors.Register(ModuleName, 16, "not enough delegation shares")
	ErrBadSharesAmount                 = sdkerrors.Register(ModuleName, 17, "invalid shares amount")
	ErrBadSharesPercent                = sdkerrors.Register(ModuleName, 18, "Invalid shares percent")
	ErrNotMature                       = sdkerrors.Register(ModuleName, 19, "entry not mature")
	ErrNoUnbondingDelegation           = sdkerrors.Register(ModuleName, 20, "no unbonding delegation found")
	ErrMaxUnbondingDelegationEntries   = sdkerrors.Register(ModuleName, 21, "too many unbonding delegation entries for (delegator, defi) tuple")
	ErrDelegatorShareExRateInvalid     = sdkerrors.Register(ModuleName, 22, "cannot delegate to defis with invalid (zero) ex-rate")
	ErrBothShareMsgsGiven              = sdkerrors.Register(ModuleName, 23, "both shares amount and shares percent provided")
	ErrNeitherShareMsgsGiven           = sdkerrors.Register(ModuleName, 24, "neither shares amount nor shares percent provided")
	ErrInvalidHistoricalInfo           = sdkerrors.Register(ModuleName, 25, "invalid historical info")
	ErrNoHistoricalInfo                = sdkerrors.Register(ModuleName, 26, "no historical info found")
	ErrEmptyDelegationDistInfo         = sdkerrors.Register(ModuleName, 27, "no delegation distribution info")
	ErrNoDefiDistInfo                  = sdkerrors.Register(ModuleName, 28, "no defi distribution info")
	ErrEmptyWithdrawAddr               = sdkerrors.Register(ModuleName, 29, "withdraw address is empty")
	ErrBadDistribution                 = sdkerrors.Register(ModuleName, 30, "community pool does not have sufficient coins to defi")
	ErrNoDefiCommission                = sdkerrors.Register(ModuleName, 31, "no defi commission to withdraw")
)
