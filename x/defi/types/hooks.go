package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// combine multiple staking hooks, all hook functions are run in array sequence
type MultiDefiHooks []DefiHooks

func NewMultiDefiHooks(hooks ...DefiHooks) MultiDefiHooks {
	return hooks
}

func (h MultiDefiHooks) AfterDefiCreated(ctx sdk.Context, defiAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterDefiCreated(ctx, defiAddr)
	}
}
func (h MultiDefiHooks) BeforeDefiModified(ctx sdk.Context, defiAddr sdk.ValAddress) {
	for i := range h {
		h[i].BeforeDefiModified(ctx, defiAddr)
	}
}
func (h MultiDefiHooks) AfterDefiRemoved(ctx sdk.Context, defiAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterDefiRemoved(ctx, defiAddr)
	}
}
func (h MultiDefiHooks) AfterDefiBonded(ctx sdk.Context, defiAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterDefiBonded(ctx, defiAddr)
	}
}
func (h MultiDefiHooks) AfterDefiBeginUnbonding(ctx sdk.Context, defiAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterDefiBeginUnbonding(ctx, defiAddr)
	}
}
func (h MultiDefiHooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) {
	for i := range h {
		h[i].BeforeDelegationCreated(ctx, delAddr, defiAddr)
	}
}
func (h MultiDefiHooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) {
	for i := range h {
		h[i].BeforeDelegationSharesModified(ctx, delAddr, defiAddr)
	}
}
func (h MultiDefiHooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) {
	for i := range h {
		h[i].BeforeDelegationRemoved(ctx, delAddr, defiAddr)
	}
}
func (h MultiDefiHooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) {
	for i := range h {
		h[i].AfterDelegationModified(ctx, delAddr, defiAddr)
	}
}
