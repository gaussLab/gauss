package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// get a single defi
func (k Keeper) GetDefi(ctx sdk.Context, addr sdk.ValAddress) (defi types.Defi, found bool) {
	store := ctx.KVStore(k.storeKey)

	value := store.Get(types.GetDefiKey(addr))
	if value == nil {
		return defi, false
	}

	defi = types.MustUnmarshalDefi(k.cdc, value)
	return defi, true
}

func (k Keeper) mustGetDefi(ctx sdk.Context, addr sdk.ValAddress) types.Defi {
	defi, found := k.GetDefi(ctx, addr)
	if !found {
		panic(fmt.Sprintf("defi record not found for address: %X\n", addr))
	}

	return defi
}

// set the main record holding defi details
func (k Keeper) SetDefi(ctx sdk.Context, defi types.Defi) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalDefi(k.cdc, &defi)
	store.Set(types.GetDefiKey(defi.GetOperator()), bz)
}

// Update the tokens of an existing defi
func (k Keeper) AddDefiTokensAndShares(ctx sdk.Context, defi types.Defi,
	tokensToAdd sdk.Int) (defiOut types.Defi, addedShares sdk.Dec) {
	defi, addedShares = defi.AddTokensFromDel(tokensToAdd)
	k.SetDefi(ctx, defi)

	return defi, addedShares
}

// Update the tokens of an existing defi
func (k Keeper) RemoveDefiTokensAndShares(ctx sdk.Context, defi types.Defi,
	sharesToRemove sdk.Dec) (defiOut types.Defi, removedTokens sdk.Int) {
	defi, removedTokens = defi.RemoveDelShares(sharesToRemove)
	k.SetDefi(ctx, defi)

	return defi, removedTokens
}

// Update the tokens of an existing defi
func (k Keeper) RemoveDefiTokens(ctx sdk.Context,
	defi types.Defi, tokensToRemove sdk.Int) types.Defi {
	defi = defi.RemoveTokens(tokensToRemove)
	k.SetDefi(ctx, defi)

	return defi
}

// remove the defi record and associated indexes
// except for the bonded defi index which is only handled in ApplyAndReturnTendermintUpdates
// TODO, this function panics, and it's not good.
func (k Keeper) RemoveDefi(ctx sdk.Context, address sdk.ValAddress) {
	// first retrieve the old defi record
	defi, found := k.GetDefi(ctx, address)
	if !found {
		return
	}

	if !defi.IsUnbonded() {
		panic("cannot call RemoveDefi on bonded or unbonding defis")
	}

	if defi.Tokens.IsPositive() {
		panic("attempting to remove a defi which still contains tokens")
	}

	// delete the old defi record
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDefiKey(address))

	// call hooks
	k.AfterDefiRemoved(ctx, defi.GetOperator())
}

// get groups of defis

// get the set of all defis with no limits, used during genesis dump
func (k Keeper) GetAllDefis(ctx sdk.Context) (defis []types.Defi) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.DefisKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		defi := types.MustUnmarshalDefi(k.cdc, iterator.Value())
		defis = append(defis, defi)
	}

	return defis
}

// return a given amount of all the defis
func (k Keeper) GetDefis(ctx sdk.Context, maxRetrieve uint32) (defis []types.Defi) {
	store := ctx.KVStore(k.storeKey)
	defis = make([]types.Defi, maxRetrieve)

	iterator := sdk.KVStorePrefixIterator(store, types.DefisKey)
	defer iterator.Close()

	i := 0
	for ; iterator.Valid() && i < int(maxRetrieve); iterator.Next() {
		defi := types.MustUnmarshalDefi(k.cdc, iterator.Value())
		defis[i] = defi
		i++
	}

	return defis[:i] // trim if the array length < maxRetrieve
}

//_______________________________________________________________________
// GetUnbondingDefis returns a slice of mature defi addresses that
// complete their unbonding at a given time and height.
func (k Keeper) GetUnbondingDefis(ctx sdk.Context, endTime time.Time, endHeight int64) []string {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetDefiQueueKey(endTime, endHeight))
	if bz == nil {
		return []string{}
	}

	addrs := types.DefiAddresses{}
	k.cdc.MustUnmarshalBinaryBare(bz, &addrs)

	return addrs.Addresses
}

// SetUnbondingDefisQueue sets a given slice of defi addresses into
// the unbonding defi queue by a given height and time.
func (k Keeper) SetUnbondingDefisQueue(ctx sdk.Context, endTime time.Time, endHeight int64, addrs []string) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&types.DefiAddresses{Addresses: addrs})
	store.Set(types.GetDefiQueueKey(endTime, endHeight), bz)
}

// InsertUnbondingDefiQueue inserts a given unbonding defi address into
// the unbonding defi queue for a given height and time.
func (k Keeper) InsertUnbondingDefiQueue(ctx sdk.Context, defi types.Defi) {
	addrs := k.GetUnbondingDefis(ctx, defi.UnbondingTime, defi.UnbondingHeight)
	addrs = append(addrs, defi.OperatorAddress)
	k.SetUnbondingDefisQueue(ctx, defi.UnbondingTime, defi.UnbondingHeight, addrs)
}

// DeleteDefiQueueTimeSlice deletes all entries in the queue indexed by a
// given height and time.
func (k Keeper) DeleteDefiQueueTimeSlice(ctx sdk.Context, endTime time.Time, endHeight int64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetDefiQueueKey(endTime, endHeight))
}

// DeleteDefiQueue removes a defi by address from the unbonding queue
// indexed by a given height and time.
func (k Keeper) DeleteDefiQueue(ctx sdk.Context, defi types.Defi) {
	addrs := k.GetUnbondingDefis(ctx, defi.UnbondingTime, defi.UnbondingHeight)
	newAddrs := []string{}

	for _, addr := range addrs {
		if addr != defi.OperatorAddress {
			newAddrs = append(newAddrs, addr)
		}
	}

	if len(newAddrs) == 0 {
		k.DeleteDefiQueueTimeSlice(ctx, defi.UnbondingTime, defi.UnbondingHeight)
	} else {
		k.SetUnbondingDefisQueue(ctx, defi.UnbondingTime, defi.UnbondingHeight, newAddrs)
	}
}

// DefiQueueIterator returns an interator ranging over defis that are
// unbonding whose unbonding completion occurs at the given height and time.
func (k Keeper) DefiQueueIterator(ctx sdk.Context, endTime time.Time, endHeight int64) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return store.Iterator(types.DefiQueueKey, sdk.InclusiveEndBytes(types.GetDefiQueueKey(endTime, endHeight)))
}

// UnbondAllMatureDefis unbonds all the mature unbonding defis that
// have finished their unbonding period.
func (k Keeper) UnbondAllMatureDefis(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	blockTime := ctx.BlockTime()
	blockHeight := ctx.BlockHeight()

	// unbondingValIterator will contains all defi addresses indexed under
	// the DefiQueueKey prefix. Note, the entire index key is composed as
	// DefiQueueKey | timeBzLen (8-byte big endian) | timeBz | heightBz (8-byte big endian),
	// so it may be possible that certain defi addresses that are iterated
	// over are not ready to unbond, so an explicit check is required.
	unbondingDefiIterator := k.DefiQueueIterator(ctx, blockTime, blockHeight)
	defer unbondingDefiIterator.Close()

	for ; unbondingDefiIterator.Valid(); unbondingDefiIterator.Next() {
		key := unbondingDefiIterator.Key()
		keyTime, keyHeight, err := types.ParseDefiQueueKey(key)
		if err != nil {
			panic(fmt.Errorf("failed to parse unbonding key: %w", err))
		}

		// All addresses for the given key have the same unbonding height and time.
		// We only unbond if the height and time are less than the current height
		// and time.
		if keyHeight <= blockHeight && (keyTime.Before(blockTime) || keyTime.Equal(blockTime)) {
			addrs := types.DefiAddresses{}
			k.cdc.MustUnmarshalBinaryBare(unbondingDefiIterator.Value(), &addrs)

			for _, defiAddr := range addrs.Addresses {
				addr, err := sdk.ValAddressFromBech32(defiAddr)
				if err != nil {
					panic(err)
				}
				defi, found := k.GetDefi(ctx, addr)
				if !found {
					panic("defi in the unbonding queue was not found")
				}

				if !defi.IsUnbonding() {
					panic("unexpected defi in unbonding queue; status was not unbonding")
				}

				defi = k.UnbondingToUnbonded(ctx, defi)
				if defi.GetDelegatorShares().IsZero() {
					k.RemoveDefi(ctx, defi.GetOperator())
				}
			}

			store.Delete(key)
		}
	}
}

// -------------------------------------------
// rewards
// initialize rewards for a new defi
func (k Keeper) initializeDefi(ctx sdk.Context, defi types.DefiI) {
	// set initial historical rewards (period 0) with reference count of 1
	k.SetDefiHistoricalRewards(ctx, defi.GetOperator(), 0, types.NewDefiHistoricalRewards(sdk.DecCoins{}, 1))

	// set current rewards (starting at period 1)
	k.SetDefiCurrentRewards(ctx, defi.GetOperator(), types.NewDefiCurrentRewards(sdk.DecCoins{}, 1))

	// set accumulated commission
	k.SetDefiAccumulatedCommission(ctx, defi.GetOperator(), types.InitialDefiAccumulatedCommission())

	// set outstanding rewards
	k.SetDefiOutstandingRewards(ctx, defi.GetOperator(), types.DefiOutstandingRewards{Rewards: sdk.DecCoins{}})
}

// increment defi period, returning the period just ended
func (k Keeper) IncrementDefiPeriod(ctx sdk.Context, defi types.DefiI) uint64 {
	// fetch current rewards
	rewards := k.GetDefiCurrentRewards(ctx, defi.GetOperator())

	// calculate current ratio
	var current sdk.DecCoins
	if defi.GetTokens().IsZero() {

		// can't calculate ratio for zero-token defis
		// ergo we instead add to the community pool
		feePool := k.GetFeePool(ctx)
		outstanding := k.GetDefiOutstandingRewards(ctx, defi.GetOperator())
		feePool.CommunityPool = feePool.CommunityPool.Add(rewards.Rewards...)
		outstanding.Rewards = outstanding.GetRewards().Sub(rewards.Rewards)
		k.SetFeePool(ctx, feePool)
		k.SetDefiOutstandingRewards(ctx, defi.GetOperator(), outstanding)

		current = sdk.DecCoins{}
	} else {
		// note: necessary to truncate so we don't allow withdrawing more rewards than owed
		current = rewards.Rewards.QuoDecTruncate(defi.GetTokens().ToDec())
	}

	// fetch historical rewards for last period
	historical := k.GetDefiHistoricalRewards(ctx, defi.GetOperator(), rewards.Period-1).CumulativeRewardRatio

	// decrement reference count
	k.decrementReferenceCount(ctx, defi.GetOperator(), rewards.Period-1)

	// set new historical rewards with reference count of 1
	k.SetDefiHistoricalRewards(ctx, defi.GetOperator(), rewards.Period, types.NewDefiHistoricalRewards(historical.Add(current...), 1))

	// set current rewards, incrementing period by 1
	k.SetDefiCurrentRewards(ctx, defi.GetOperator(), types.NewDefiCurrentRewards(sdk.DecCoins{}, rewards.Period+1))

	return rewards.Period
}

// increment the reference count for a historical rewards value
func (k Keeper) incrementReferenceCount(ctx sdk.Context, defiAddr sdk.ValAddress, period uint64) {
	historical := k.GetDefiHistoricalRewards(ctx, defiAddr, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.SetDefiHistoricalRewards(ctx, defiAddr, period, historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func (k Keeper) decrementReferenceCount(ctx sdk.Context, defiAddr sdk.ValAddress, period uint64) {
	historical := k.GetDefiHistoricalRewards(ctx, defiAddr, period)
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		k.DeleteDefiHistoricalReward(ctx, defiAddr, period)
	} else {
		k.SetDefiHistoricalRewards(ctx, defiAddr, period, historical)
	}
}
