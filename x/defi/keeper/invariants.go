package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// RegisterInvariants registers all defi invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "module-accounts",
		ModuleAccountInvariants(k))
	ir.RegisterRoute(types.ModuleName, "positive-delegation",
		PositiveDelegationInvariant(k))
	ir.RegisterRoute(types.ModuleName, "delegator-shares",
		DelegatorSharesInvariant(k))
	ir.RegisterRoute(types.ModuleName, "nonnegative-outstanding",
		NonNegativeOutstandingInvariant(k))
	ir.RegisterRoute(types.ModuleName, "can-withdraw",
		CanWithdrawInvariant(k))
	ir.RegisterRoute(types.ModuleName, "reference-count",
		ReferenceCountInvariant(k))
}

// AllInvariants runs all invariants of the defi module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := DelegatorSharesInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		res, stop = PositiveDelegationInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		res, stop = CanWithdrawInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		res, stop = NonNegativeOutstandingInvariant(k)(ctx)
		if stop {
			return res, stop
		}
		res, stop = ReferenceCountInvariant(k)(ctx)
		if stop {
			return res, stop
		}

		return ModuleAccountInvariants(k)(ctx)
	}
}

// ModuleAccountInvariants checks that the bonded and notBonded ModuleAccounts pools
// reflects the tokens actively bonded and not bonded
func ModuleAccountInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		bonded := sdk.ZeroInt()
		notBonded := sdk.ZeroInt()
		bondedPool := k.GetBondedPool(ctx)
		notBondedPool := k.GetNotBondedPool(ctx)
		bondDenom := k.BondDenom(ctx)

		k.IterateDefis(ctx, func(_ int64, defi types.DefiI) bool {
			switch defi.GetStatus() {
			case types.Bonded:
				bonded = bonded.Add(defi.GetTokens())
			case types.Unbonding, types.Unbonded:
				notBonded = notBonded.Add(defi.GetTokens())
			default:
				panic("invalid defi status")
			}
			return false
		})

		k.IterateUnbondingDelegations(ctx, func(_ int64, ubd types.UnbondingDelegation) bool {
			for _, entry := range ubd.Entries {
				notBonded = notBonded.Add(entry.Balance)
			}
			return false
		})

		poolBonded := k.bankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom)
		poolNotBonded := k.bankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom)

		//------------------------------------
		// rewards
		var expectedCoins sdk.DecCoins
		k.IterateDefiOutstandingRewards(ctx, func(_ sdk.ValAddress, rewards types.DefiOutstandingRewards) (stop bool) {
			expectedCoins = expectedCoins.Add(rewards.Rewards...)
			return false
		})

		communityPool := k.GetFeePoolCommunityCoins(ctx)
		expectedInt, _ := expectedCoins.Add(communityPool...).TruncateDecimal()

		macc := k.GetDefiAccount(ctx)
		balances := k.bankKeeper.GetAllBalances(ctx, macc.GetAddress())
		// -----------------------------------

		broken := !poolBonded.Amount.Equal(bonded) || !poolNotBonded.Amount.Equal(notBonded) || !balances.IsEqual(expectedInt)

		// Bonded tokens should equal sum of tokens with bonded defis
		// Not-bonded tokens should equal unbonding delegations	plus tokens on unbonded defis
		return sdk.FormatInvariant(types.ModuleName, "bonded and not bonded module account coins", fmt.Sprintf(
			"\tPool's bonded tokens: %v\n"+
				"\tsum of bonded tokens: %v\n"+
				"not bonded token invariance:\n"+
				"\tPool's not bonded tokens: %v\n"+
				"\tsum of not bonded tokens: %v\n"+
				"module accounts total (bonded + not bonded):\n"+
				"\tModule Accounts' tokens: %v\n"+
				"\tsum tokens:              %v\n"+
				"\trewards tokens:  expected - %s, got - %s\n",
			poolBonded, bonded, poolNotBonded, notBonded, poolBonded.Add(poolNotBonded), bonded.Add(notBonded), expectedInt, balances)), broken
	}
}

// PositiveDelegationInvariant checks that all stored delegations have > 0 shares.
func PositiveDelegationInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			msg   string
			count int
		)

		delegations := k.GetAllDelegations(ctx)
		for _, delegation := range delegations {
			if delegation.Shares.IsNegative() {
				count++

				msg += fmt.Sprintf("\tdelegation with negative shares: %+v\n", delegation)
			}

			if delegation.Shares.IsZero() {
				count++

				msg += fmt.Sprintf("\tdelegation with zero shares: %+v\n", delegation)
			}
		}

		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "positive delegations", fmt.Sprintf(
			"%d invalid delegations found\n%s", count, msg)), broken
	}
}

// DelegatorSharesInvariant checks whether all the delegator shares which persist
// in the delegator object add up to the correct total delegator shares
// amount stored in each defi.
func DelegatorSharesInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var (
			msg    string
			broken bool
		)

		defis := k.GetAllDefis(ctx)
		for _, defi := range defis {
			defiTotalDelShares := defi.GetDelegatorShares()
			totalDelShares := sdk.ZeroDec()

			delegations := k.GetDefiDelegations(ctx, defi.GetOperator())
			for _, delegation := range delegations {
				totalDelShares = totalDelShares.Add(delegation.Shares)
			}

			if !defiTotalDelShares.Equal(totalDelShares) {
				broken = true
				msg += fmt.Sprintf("broken delegator shares invariance:\n"+
					"\tdefi.DelegatorShares: %v\n"+
					"\tsum of Delegator.Shares: %v\n", defiTotalDelShares, totalDelShares)
			}
		}

		return sdk.FormatInvariant(types.ModuleName, "delegator shares", msg), broken
	}
}

// NonNegativeOutstandingInvariant checks that outstanding unwithdrawn fees are never negative
func NonNegativeOutstandingInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		var msg string
		var count int
		var outstanding sdk.DecCoins

		k.IterateDefiOutstandingRewards(ctx, func(addr sdk.ValAddress, rewards types.DefiOutstandingRewards) (stop bool) {
			outstanding = rewards.GetRewards()
			if outstanding.IsAnyNegative() {
				count++
				msg += fmt.Sprintf("\t%v has negative outstanding coins: %v\n", addr, outstanding)
			}
			return false
		})
		broken := count != 0

		return sdk.FormatInvariant(types.ModuleName, "nonnegative outstanding",
			fmt.Sprintf("found %d defis with negative outstanding rewards\n%s", count, msg)), broken
	}
}

// CanWithdrawInvariant checks that current rewards can be completely withdrawn
func CanWithdrawInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		// cache, we don't want to write changes
		ctx, _ = ctx.CacheContext()

		var remaining sdk.DecCoins

		defiDelegationAddrs := make(map[string][]sdk.AccAddress)
		for _, del := range k.GetAllSDKDelegations(ctx) {
			defiAddr := del.GetDefiAddr().String()
			defiDelegationAddrs[defiAddr] = append(defiDelegationAddrs[defiAddr], del.GetDelegatorAddr())
		}

		// iterate over all defis
		k.IterateDefis(ctx, func(_ int64, defi types.DefiI) (stop bool) {
			_, _ = k.WithdrawDefiCommission(ctx, defi.GetOperator())

			delegationAddrs, ok := defiDelegationAddrs[defi.GetOperator().String()]
			if ok {
				for _, delAddr := range delegationAddrs {
					if _, err := k.WithdrawDelegationRewards(ctx, delAddr, defi.GetOperator()); err != nil {
						panic(err)
					}
				}
			}

			remaining = k.GetDefiOutstandingRewardsCoins(ctx, defi.GetOperator())
			if len(remaining) > 0 && remaining[0].Amount.IsNegative() {
				return true
			}

			return false
		})

		broken := len(remaining) > 0 && remaining[0].Amount.IsNegative()
		return sdk.FormatInvariant(types.ModuleName, "can withdraw",
			fmt.Sprintf("remaining coins: %v\n", remaining)), broken
	}
}

// ReferenceCountInvariant checks that the number of historical rewards records is correct
func ReferenceCountInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {

		defiCount := uint64(0)
		k.IterateDefis(ctx, func(_ int64, defi types.DefiI) (stop bool) {
			defiCount++
			return false
		})
		dels := k.GetAllSDKDelegations(ctx)

		// one record per defi (last tracked period), one record per
		// delegation (previous period), one record per slash (previous period)
		expected := defiCount + uint64(len(dels))
		count := k.GetDefiHistoricalReferenceCount(ctx)
		broken := count != expected

		return sdk.FormatInvariant(types.ModuleName, "reference count",
			fmt.Sprintf("expected historical reference count: %d = %v defis + %v delegations \n"+
				"total defi historical reference count: %d\n",
				expected, defiCount, len(dels), count)), broken
	}
}
