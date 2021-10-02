package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankexported "github.com/cosmos/cosmos-sdk/x/bank/exported"
)

// DistributionKeeper expected distribution keeper (noalias)
type DistributionKeeper interface {
	GetFeePoolCommunityCoins(ctx sdk.Context) sdk.DecCoins
	GetDefiOutstandingRewardsCoins(ctx sdk.Context, val sdk.ValAddress) sdk.DecCoins
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	IterateAccounts(ctx sdk.Context, process func(authtypes.AccountI) (stop bool))
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.AccountI // only used for simulation

	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, moduleName string) authtypes.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, authtypes.ModuleAccountI)
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	SetBalances(ctx sdk.Context, addr sdk.AccAddress, balances sdk.Coins) error
	LockedCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	GetSupply(ctx sdk.Context) bankexported.SupplyI

	SendCoinsFromModuleToModule(ctx sdk.Context, senderPool, recipientPool string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	UndelegateCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	DelegateCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error

	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}

type TokenKeeper interface {
	// to module
	MintTokenWithUnit(ctx sdk.Context, unit string, amount uint64, recipient string) error
}

// DefiSet expected properties for the set of all defi (noalias)
type DefiSet interface {
	// iterate through defis by operator address, execute func for each defi
	IterateDefis(sdk.Context,
		func(index int64, defi DefiI) (stop bool))

	Defi(sdk.Context, sdk.ValAddress) DefiI            // get a particular defi by operator address
	TotalBondedTokens(sdk.Context) sdk.Int             // total bonded tokens within the defi set
	DefiTokenSupply(sdk.Context) sdk.Int               // total staking token supply

	// Delegation allows for getting a particular delegation for a given defi
	// and delegator outside the scope of the staking module.
	Delegation(sdk.Context, sdk.AccAddress, sdk.ValAddress) DelegationI

	// MaxDefis returns the maximum amount of bonded defis
	MaxDefis(sdk.Context) uint32

	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation DelegationI) (stop bool))

	GetAllSDKDelegations(ctx sdk.Context) []Delegation
}

// DelegationSet expected properties for the set of all delegations for a particular (noalias)
type DelegationSet interface {
	GetDefiSet() DefiSet // defi set for which delegation set is based upon

	// iterate through all delegations from one delegator by defi-AccAddress,
	//   execute func for each defi
	IterateDelegations(ctx sdk.Context, delegator sdk.AccAddress,
		fn func(index int64, delegation DelegationI) (stop bool))

	GetAllSDKDelegations(ctx sdk.Context) []Delegation
}

//_______________________________________________________________________________
// Event Hooks
// These can be utilized to communicate between a defi keeper and another
// keeper which must take particular actions when defi/delegators change
// state. The second keeper must implement this interface, which then the
// defi keeper can call.

// DefiHooks event hooks for staking defi object (noalias)
type DefiHooks interface {
	AfterDefiCreated(ctx sdk.Context, defiAddr sdk.ValAddress)                           // Must be called when a defi is created
	BeforeDefiModified(ctx sdk.Context, defiAddr sdk.ValAddress)                         // Must be called when a defi's state changes
	AfterDefiRemoved(ctx sdk.Context, defiAddr sdk.ValAddress) // Must be called when a defi is deleted

	AfterDefiBonded(ctx sdk.Context, defiAddr sdk.ValAddress)         // Must be called when a defi is bonded
	AfterDefiBeginUnbonding(ctx sdk.Context, defiAddr sdk.ValAddress) // Must be called when a defi begins unbonding

	BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress)        // Must be called when a delegation is created
	BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) // Must be called when a delegation's shares are modified
	BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress)        // Must be called when a delegation is removed
	AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress)
}
