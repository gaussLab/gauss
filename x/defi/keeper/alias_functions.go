package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

//_______________________________________________________________________
// Defi Set

// iterate through the defi set and perform the provided function
func (k Keeper) IterateDefis(ctx sdk.Context, fn func(index int64, defi types.DefiI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.DefisKey)
	defer iterator.Close()

	i := int64(0)

	for ; iterator.Valid(); iterator.Next() {
		defi := types.MustUnmarshalDefi(k.cdc, iterator.Value())
		stop := fn(i, defi) // XXX is this safe will the defi unexposed fields be able to get written to?

		if stop {
			break
		}
		i++
	}
}

// Defi gets the Defi interface for a particular address
func (k Keeper) Defi(ctx sdk.Context, address sdk.ValAddress) types.DefiI {
	val, found := k.GetDefi(ctx, address)
	if !found {
		return nil
	}

	return val
}

//_______________________________________________________________________
// Delegation Set

// Returns self as it is both a defiset and delegationset
func (k Keeper) GetDefiSet() types.DefiSet {
	return k
}

// Delegation get the delegation interface for a particular set of delegator and defi addresses
func (k Keeper) Delegation(ctx sdk.Context, addrDel sdk.AccAddress, addrDefi sdk.ValAddress) types.DelegationI {
	bond, ok := k.GetDelegation(ctx, addrDel, addrDefi)
	if !ok {
		return nil
	}

	return bond
}

// iterate through all of the delegations from a delegator
func (k Keeper) IterateDelegations(ctx sdk.Context, delAddr sdk.AccAddress,
	fn func(index int64, del types.DelegationI) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	delegatorPrefixKey := types.GetDelegationsKey(delAddr)

	iterator := sdk.KVStorePrefixIterator(store, delegatorPrefixKey) // smallest to largest
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		del := types.MustUnmarshalDelegation(k.cdc, iterator.Value())

		stop := fn(i, del)
		if stop {
			break
		}
		i++
	}
}

// return all delegations used during genesis dump
// TODO: remove this func, change all usage for iterate functionality
func (k Keeper) GetAllSDKDelegations(ctx sdk.Context) (delegations []types.Delegation) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.DelegationKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		delegation := types.MustUnmarshalDelegation(k.cdc, iterator.Value())
		delegations = append(delegations, delegation)
	}

	return
}

func (k Keeper) GetDelegationAmount(ctx sdk.Context, defiAddr sdk.ValAddress, delegator sdk.AccAddress) (amount sdk.Dec) {
	defi := k.Defi(ctx, defiAddr)
	if defi == nil {
		return sdk.ZeroDec()
	}

	delegation := k.Delegation(ctx, delegator, defiAddr)
	if delegation == nil {
		return sdk.ZeroDec()
	}
	

	return defi.TokensFromSharesTruncated(delegation.GetShares())
}

//_______________________________________________________________________
// Rewards Set
// get outstanding rewards
func (k Keeper) GetDefiOutstandingRewardsCoins(ctx sdk.Context, defi sdk.ValAddress) sdk.DecCoins {
	return k.GetDefiOutstandingRewards(ctx, defi).Rewards
}

// get the community coins
func (k Keeper) GetFeePoolCommunityCoins(ctx sdk.Context) sdk.DecCoins {
	return k.GetFeePool(ctx).CommunityPool
}

// GetDefiAccount returns the defi ModuleAccount
func (k Keeper) GetDefiAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
}
