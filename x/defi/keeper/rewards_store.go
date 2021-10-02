package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// get the delegator withdraw address, defaulting to the delegator address
func (k Keeper) GetDelegatorWithdrawAddr(ctx sdk.Context, delAddr sdk.AccAddress) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDelegatorWithdrawAddrKey(delAddr))
	if b == nil {
		return delAddr
	}
	return sdk.AccAddress(b)
}

// set the delegator withdraw address
func (k Keeper) SetDelegatorWithdrawAddr(ctx sdk.Context, delAddr, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetDelegatorWithdrawAddrKey(delAddr), withdrawAddr.Bytes())
}

// delete a delegator withdraw addr
func (k Keeper) DeleteDelegatorWithdrawAddr(ctx sdk.Context, delAddr, withdrawAddr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegatorWithdrawAddrKey(delAddr))
}

// iterate over delegator withdraw addrs
func (k Keeper) IterateDelegatorWithdrawAddrs(ctx sdk.Context, handler func(del sdk.AccAddress, addr sdk.AccAddress) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DelegatorWithdrawAddrPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Value())
		del := types.GetDelegatorWithdrawInfoAddress(iter.Key())
		if handler(del, addr) {
			break
		}
	}
}

// get the global fee pool distribution info
func (k Keeper) GetFeePool(ctx sdk.Context) (feePool types.FeePool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.FeePoolKey)
	if b == nil {
		panic("Stored fee pool should not have been nil")
	}
	k.cdc.MustUnmarshalBinaryBare(b, &feePool)
	return
}

// set the global fee pool distribution info
func (k Keeper) SetFeePool(ctx sdk.Context, feePool types.FeePool) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(&feePool)
	store.Set(types.FeePoolKey, b)
}

// get the starting info associated with a delegator
func (k Keeper) GetDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) (period types.DelegatorStartingInfo) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDelegatorStartingInfoKey(val, del))
	k.cdc.MustUnmarshalBinaryBare(b, &period)
	return
}

// set the starting info associated with a delegator
func (k Keeper) SetDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress, period types.DelegatorStartingInfo) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(&period)
	store.Set(types.GetDelegatorStartingInfoKey(val, del), b)
}

// check existence of the starting info associated with a delegator
func (k Keeper) HasDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDelegatorStartingInfoKey(val, del))
}

// delete the starting info associated with a delegator
func (k Keeper) DeleteDelegatorStartingInfo(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDelegatorStartingInfoKey(val, del))
}

// iterate over delegator starting infos
func (k Keeper) IterateDelegatorStartingInfos(ctx sdk.Context, handler func(val sdk.ValAddress, del sdk.AccAddress, info types.DelegatorStartingInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DelegatorStartingInfoPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var info types.DelegatorStartingInfo
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &info)
		val, del := types.GetDelegatorStartingInfoAddresses(iter.Key())
		if handler(val, del, info) {
			break
		}
	}
}

// get historical rewards for a particular period
func (k Keeper) GetDefiHistoricalRewards(ctx sdk.Context, val sdk.ValAddress, period uint64) (rewards types.DefiHistoricalRewards) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDefiHistoricalRewardsKey(val, period))
	k.cdc.MustUnmarshalBinaryBare(b, &rewards)
	return
}

// set historical rewards for a particular period
func (k Keeper) SetDefiHistoricalRewards(ctx sdk.Context, val sdk.ValAddress, period uint64, rewards types.DefiHistoricalRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(&rewards)
	store.Set(types.GetDefiHistoricalRewardsKey(val, period), b)
}

// iterate over historical rewards
func (k Keeper) IterateDefiHistoricalRewards(ctx sdk.Context, handler func(val sdk.ValAddress, period uint64, rewards types.DefiHistoricalRewards) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DefiHistoricalRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards types.DefiHistoricalRewards
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &rewards)
		addr, period := types.GetDefiHistoricalRewardsAddressPeriod(iter.Key())
		if handler(addr, period, rewards) {
			break
		}
	}
}

// delete a historical reward
func (k Keeper) DeleteDefiHistoricalReward(ctx sdk.Context, val sdk.ValAddress, period uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDefiHistoricalRewardsKey(val, period))
}

// delete historical rewards for a defi
func (k Keeper) DeleteDefiHistoricalRewards(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetDefiHistoricalRewardsPrefix(val))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// delete all historical rewards
func (k Keeper) DeleteAllDefiHistoricalRewards(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DefiHistoricalRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// historical reference count (used for testcases)
func (k Keeper) GetDefiHistoricalReferenceCount(ctx sdk.Context) (count uint64) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DefiHistoricalRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards types.DefiHistoricalRewards
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &rewards)
		count += uint64(rewards.ReferenceCount)
	}
	return
}

// get current rewards for a defi
func (k Keeper) GetDefiCurrentRewards(ctx sdk.Context, val sdk.ValAddress) (rewards types.DefiCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDefiCurrentRewardsKey(val))
	k.cdc.MustUnmarshalBinaryBare(b, &rewards)
	return
}

// set current rewards for a defi
func (k Keeper) SetDefiCurrentRewards(ctx sdk.Context, val sdk.ValAddress, rewards types.DefiCurrentRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(&rewards)
	store.Set(types.GetDefiCurrentRewardsKey(val), b)
}

// delete current rewards for a defi
func (k Keeper) DeleteDefiCurrentRewards(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDefiCurrentRewardsKey(val))
}

// iterate over current rewards
func (k Keeper) IterateDefiCurrentRewards(ctx sdk.Context, handler func(val sdk.ValAddress, rewards types.DefiCurrentRewards) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DefiCurrentRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var rewards types.DefiCurrentRewards
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &rewards)
		addr := types.GetDefiCurrentRewardsAddress(iter.Key())
		if handler(addr, rewards) {
			break
		}
	}
}

// get accumulated commission for a defi
func (k Keeper) GetDefiAccumulatedCommission(ctx sdk.Context, val sdk.ValAddress) (commission types.DefiAccumulatedCommission) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetDefiAccumulatedCommissionKey(val))
	if b == nil {
		return types.DefiAccumulatedCommission{}
	}
	k.cdc.MustUnmarshalBinaryBare(b, &commission)
	return
}

// set accumulated commission for a defi
func (k Keeper) SetDefiAccumulatedCommission(ctx sdk.Context, val sdk.ValAddress, commission types.DefiAccumulatedCommission) {
	var bz []byte

	store := ctx.KVStore(k.storeKey)
	if commission.Commission.IsZero() {
		bz = k.cdc.MustMarshalBinaryBare(&types.DefiAccumulatedCommission{})
	} else {
		bz = k.cdc.MustMarshalBinaryBare(&commission)
	}

	store.Set(types.GetDefiAccumulatedCommissionKey(val), bz)
}

// delete accumulated commission for a defi
func (k Keeper) DeleteDefiAccumulatedCommission(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDefiAccumulatedCommissionKey(val))
}

// iterate over accumulated commissions
func (k Keeper) IterateDefiAccumulatedCommissions(ctx sdk.Context, handler func(val sdk.ValAddress, commission types.DefiAccumulatedCommission) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DefiAccumulatedCommissionPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var commission types.DefiAccumulatedCommission
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &commission)
		addr := types.GetDefiAccumulatedCommissionAddress(iter.Key())
		if handler(addr, commission) {
			break
		}
	}
}

// get defi outstanding rewards
func (k Keeper) GetDefiOutstandingRewards(ctx sdk.Context, val sdk.ValAddress) (rewards types.DefiOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetDefiOutstandingRewardsKey(val))
	k.cdc.MustUnmarshalBinaryBare(bz, &rewards)
	return
}

// set defi outstanding rewards
func (k Keeper) SetDefiOutstandingRewards(ctx sdk.Context, val sdk.ValAddress, rewards types.DefiOutstandingRewards) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryBare(&rewards)
	store.Set(types.GetDefiOutstandingRewardsKey(val), b)
}

// delete defi outstanding rewards
func (k Keeper) DeleteDefiOutstandingRewards(ctx sdk.Context, val sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDefiOutstandingRewardsKey(val))
}

// iterate defi outstanding rewards
func (k Keeper) IterateDefiOutstandingRewards(ctx sdk.Context, handler func(val sdk.ValAddress, rewards types.DefiOutstandingRewards) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.DefiOutstandingRewardsPrefix)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		rewards := types.DefiOutstandingRewards{}
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &rewards)
		addr := types.GetDefiOutstandingRewardsAddress(iter.Key())
		if handler(addr, rewards) {
			break
		}
	}
}
