package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// Implements DefiHooks interface
var _ types.DefiHooks = Keeper{}

// AfterDefiCreated - call hook if registered
func (k Keeper) AfterDefiCreated(ctx sdk.Context, defiAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterDefiCreated(ctx, defiAddr)
	}
}

// BeforeDefiModified - call hook if registered
func (k Keeper) BeforeDefiModified(ctx sdk.Context, defiAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDefiModified(ctx, defiAddr)
	}
}

// AfterDefiRemoved - call hook if registered
func (k Keeper) AfterDefiRemoved(ctx sdk.Context, defiAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterDefiRemoved(ctx, defiAddr)
	}
}

// AfterDefiBonded - call hook if registered
func (k Keeper) AfterDefiBonded(ctx sdk.Context, defiAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterDefiBonded(ctx, defiAddr)
	}
}

// AfterDefiBeginUnbonding - call hook if registered
func (k Keeper) AfterDefiBeginUnbonding(ctx sdk.Context, defiAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterDefiBeginUnbonding(ctx, defiAddr)
	}
}

// BeforeDelegationCreated - call hook if registered
func (k Keeper) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationCreated(ctx, delAddr, defiAddr)
	}
}

// BeforeDelegationSharesModified - call hook if registered
func (k Keeper) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationSharesModified(ctx, delAddr, defiAddr)
	}
}

// BeforeDelegationRemoved - call hook if registered
func (k Keeper) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.BeforeDelegationRemoved(ctx, delAddr, defiAddr)
	}
}

// AfterDelegationModified - call hook if registered
func (k Keeper) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) {
	if k.hooks != nil {
		k.hooks.AfterDelegationModified(ctx, delAddr, defiAddr)
	}
}
