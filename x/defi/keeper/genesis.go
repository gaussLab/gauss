package keeper

import (
	"fmt"
	"log"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// InitGenesis sets the pool and parameters for the provided keeper.  For each
// defi in data, it sets that defi in the keeper along with manually
// setting the indexes. In addition, it also sets any delegations found in
// data. Finally, it updates the bonded defis.
// Returns final defi set after applying all declaration and delegations
func (k Keeper) InitGenesis(
	ctx sdk.Context,  data *types.GenesisState,
) (res []abci.ValidatorUpdate) {

	bondedTokens := sdk.ZeroInt()
	notBondedTokens := sdk.ZeroInt()

	// We need to pretend to be "n blocks before genesis", where "n" is the
	// defi update delay, so that e.g. slashing periods are correctly
	// initialized for the defi set e.g. with a one-block offset - the
	// first TM block is at height 1, so state updates applied from
	// genesis.json are in block 0.
	// ctx = ctx.WithBlockHeight(1 - sdk.ValidatorUpdateDelay)

	k.SetFeePool(ctx, data.FeePool)
	k.SetParams(ctx, data.Params)

	// Defi-Delegation
	for _, defi := range data.Defis {
		k.SetDefi(ctx, defi)

		// Call the creation hook if not exported
		if !data.Exported {
			k.AfterDefiCreated(ctx, defi.GetOperator())
		}

		// update timeslice if necessary
		if defi.IsUnbonding() {
			k.InsertUnbondingDefiQueue(ctx, defi)	
		}

		switch defi.GetStatus() {
		case types.Bonded:
			bondedTokens = bondedTokens.Add(defi.GetTokens())
		case types.Unbonding, types.Unbonded:
			notBondedTokens = notBondedTokens.Add(defi.GetTokens())
		default:
			panic("invalid defi status")
		}
	}

	for _, delegation := range data.Delegations {
		delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		if err != nil {
			panic(err)
		}

		// Call the before-creation hook if not exported
		if !data.Exported {
			k.BeforeDelegationCreated(ctx, delegatorAddress, delegation.GetDefiAddr())
		}

		k.SetDelegation(ctx, delegation)
		// Call the after-modification hook if not exported
		if !data.Exported {
			k.AfterDelegationModified(ctx, delegatorAddress, delegation.GetDefiAddr())
		}
	}

	for _, ubd := range data.UnbondingDelegations {
		k.SetUnbondingDelegation(ctx, ubd)

		for _, entry := range ubd.Entries {
			k.InsertUBDQueue(ctx, ubd, entry.CompletionTime)
			notBondedTokens = notBondedTokens.Add(entry.Balance)
		}
	}

	bondedCoins := sdk.NewCoins(sdk.NewCoin(data.Params.BondDenom, bondedTokens))
	notBondedCoins := sdk.NewCoins(sdk.NewCoin(data.Params.BondDenom, notBondedTokens))

	// check if the unbonded and bonded pools accounts exists
	bondedPool := k.GetBondedPool(ctx)
	if bondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	// add coins if not provided on genesis
	if k.bankKeeper.GetAllBalances(ctx, bondedPool.GetAddress()).IsZero() {
		if err := k.bankKeeper.SetBalances(ctx, bondedPool.GetAddress(), bondedCoins); err != nil {
			panic(err)
		}

		k.authKeeper.SetModuleAccount(ctx, bondedPool)
	}

	notBondedPool := k.GetNotBondedPool(ctx)
	if notBondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	if k.bankKeeper.GetAllBalances(ctx, notBondedPool.GetAddress()).IsZero() {
		if err := k.bankKeeper.SetBalances(ctx, notBondedPool.GetAddress(), notBondedCoins); err != nil {
			panic(err)
		}

		k.authKeeper.SetModuleAccount(ctx, notBondedPool)
	}
	// final
	if data.Exported {
	}else{
		var err error
		res, err = k.ApplyAndReturnDefiSetUpdates(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}


	// Rewards
	var moduleHoldings sdk.DecCoins
	for _, dwi := range data.DelegatorWithdrawInfos {
		delegatorAddress, err := sdk.AccAddressFromBech32(dwi.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		withdrawAddress, err := sdk.AccAddressFromBech32(dwi.WithdrawAddress)
		if err != nil {
			panic(err)
		}
		k.SetDelegatorWithdrawAddr(ctx, delegatorAddress, withdrawAddress)
	}

	for _, rew := range data.OutstandingRewards {
		defiAddr, err := sdk.ValAddressFromBech32(rew.DefiAddress)
		if err != nil {
			panic(err)
		}
		k.SetDefiOutstandingRewards(ctx, defiAddr, types.DefiOutstandingRewards{Rewards: rew.OutstandingRewards})
		moduleHoldings = moduleHoldings.Add(rew.OutstandingRewards...)
	}

        for _, acc := range data.DefiAccumulatedCommissions {
		defiAddr, err := sdk.ValAddressFromBech32(acc.DefiAddress)
		if err != nil {
			panic(err)
		}
		k.SetDefiAccumulatedCommission(ctx, defiAddr, acc.Accumulated)
	}

	for _, his := range data.DefiHistoricalRewards {
		defiAddr, err := sdk.ValAddressFromBech32(his.DefiAddress)
		if err != nil {
			panic(err)
		}
		k.SetDefiHistoricalRewards(ctx, defiAddr, his.Period, his.Rewards)
	}

	for _, cur := range data.DefiCurrentRewards {
		defiAddr, err := sdk.ValAddressFromBech32(cur.DefiAddress)
		if err != nil {
			panic(err)
		}
		k.SetDefiCurrentRewards(ctx, defiAddr, cur.Rewards)
	}

	for _, del := range data.DelegatorStartingInfos {
		defiAddr, err := sdk.ValAddressFromBech32(del.DefiAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress, err := sdk.AccAddressFromBech32(del.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		k.SetDelegatorStartingInfo(ctx, defiAddr, delegatorAddress, del.StartingInfo)
	}

	moduleHoldings = moduleHoldings.Add(data.FeePool.CommunityPool...)
	moduleHoldingsInt, _ := moduleHoldings.TruncateDecimal()
	// check if the module account exists
	moduleAcc := k.GetDefiAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	balances := k.bankKeeper.GetAllBalances(ctx, moduleAcc.GetAddress())
	if balances.IsZero() {
		if err := k.bankKeeper.SetBalances(ctx, moduleAcc.GetAddress(), moduleHoldingsInt); err != nil {
			panic(err)
		}
		k.authKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	return res
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, defis, and bonds found in
// the keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	// defi
	var unbondingDelegations []types.UnbondingDelegation

	k.IterateUnbondingDelegations(ctx, func(_ int64, ubd types.UnbondingDelegation) (stop bool) {
		unbondingDelegations = append(unbondingDelegations, ubd)
		return false
	})

	// reward
	feePool := k.GetFeePool(ctx)
	params := k.GetParams(ctx)

	dwi := make([]types.DelegatorWithdrawInfo, 0)
	k.IterateDelegatorWithdrawAddrs(ctx, func(del sdk.AccAddress, addr sdk.AccAddress) (stop bool) {
		dwi = append(dwi, types.DelegatorWithdrawInfo{
			DelegatorAddress: del.String(),
			WithdrawAddress:  addr.String(),
		})
		return false
	})

	outstanding := make([]types.DefiOutstandingRewardsRecord, 0)
	k.IterateDefiOutstandingRewards(ctx,
		func(addr sdk.ValAddress, rewards types.DefiOutstandingRewards) (stop bool) {
			outstanding = append(outstanding, types.DefiOutstandingRewardsRecord{
				DefiAddress:   addr.String(),
				OutstandingRewards: rewards.Rewards,
			})
			return false
		},
	)

	acc := make([]types.DefiAccumulatedCommissionRecord, 0)
	k.IterateDefiAccumulatedCommissions(ctx,
		func(addr sdk.ValAddress, commission types.DefiAccumulatedCommission) (stop bool) {
			acc = append(acc, types.DefiAccumulatedCommissionRecord{
				DefiAddress: addr.String(),
				Accumulated:      commission,
			})
			return false
		},
	)

	his := make([]types.DefiHistoricalRewardsRecord, 0)
	k.IterateDefiHistoricalRewards(ctx,
		func(val sdk.ValAddress, period uint64, rewards types.DefiHistoricalRewards) (stop bool) {
			his = append(his, types.DefiHistoricalRewardsRecord{
				DefiAddress:      val.String(),
				Period:           period,
				Rewards:          rewards,
			})
			return false
		},
	)

	cur := make([]types.DefiCurrentRewardsRecord, 0)
	k.IterateDefiCurrentRewards(ctx,
		func(val sdk.ValAddress, rewards types.DefiCurrentRewards) (stop bool) {
			cur = append(cur, types.DefiCurrentRewardsRecord{
				DefiAddress:      val.String(),
				Rewards:          rewards,
			})
			return false
		},
	)

	dels := make([]types.DelegatorStartingInfoRecord, 0)
	k.IterateDelegatorStartingInfos(ctx,
		func(val sdk.ValAddress, del sdk.AccAddress, info types.DelegatorStartingInfo) (stop bool) {
			dels = append(dels, types.DelegatorStartingInfoRecord{
				DefiAddress:      val.String(),
				DelegatorAddress: del.String(),
				StartingInfo:     info,
			})
			return false
		},
	)

	return &types.GenesisState{
		FeePool:		       feePool,
		Params:                        params,
		Defis:                         k.GetAllDefis(ctx),
		Delegations:                   k.GetAllDelegations(ctx),
		UnbondingDelegations:          unbondingDelegations,
		DelegatorWithdrawInfos:        dwi,
		OutstandingRewards:            outstanding,
		DefiAccumulatedCommissions:    acc,
		DefiHistoricalRewards:         his,
		DefiCurrentRewards:            cur,
		DelegatorStartingInfos:        dels,
		Exported:                      true,
	}

}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate defis)
func ValidateGenesis(data *types.GenesisState) error {
	if err := validateGenesisStateDefis(data.Defis); err != nil {
		return err
	}
	if err := data.FeePool.ValidateGenesis(); err != nil {
		return err
	}

	return data.Params.Validate()
}

func validateGenesisStateDefis(defis []types.Defi) error {
	for i := 0; i < len(defis); i++ {
		defi := defis[i]
		if defi.DelegatorShares.IsZero() && !defi.IsUnbonding() {
			return fmt.Errorf("bonded/unbonded genesis defi cannot have zero delegator shares, defi: %v", defi)
		}
	}

	return nil
}
