package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ types.DefiHooks = Hooks{}

// Create new distribution hooks
func (k Keeper) Hooks() Hooks { return Hooks{k} }

// initialize defi distribution record
func (h Hooks) AfterDefiCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	val := h.k.Defi(ctx, valAddr)
	h.k.initializeDefi(ctx, val)
}

// cleanup for after defi is removed
func (h Hooks) AfterDefiRemoved(ctx sdk.Context, defiAddr sdk.ValAddress) {
	// fetch outstanding
	outstanding := h.k.GetDefiOutstandingRewardsCoins(ctx, defiAddr)

	// force-withdraw commission
	commission := h.k.GetDefiAccumulatedCommission(ctx, defiAddr).Commission
	if !commission.IsZero() {
		// subtract from outstanding
		outstanding = outstanding.Sub(commission)

		// split into integral & remainder
		coins, remainder := commission.TruncateDecimal()

		// remainder to community pool
		feePool := h.k.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
		h.k.SetFeePool(ctx, feePool)

		// add to defi account
		if !coins.IsZero() {
			accAddr := sdk.AccAddress(defiAddr)
			withdrawAddr := h.k.GetDelegatorWithdrawAddr(ctx, accAddr)

			if err := h.k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, coins); err != nil {
				panic(err)
			}
		}
	}

	// add outstanding to community pool
	feePool := h.k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(outstanding...)
	h.k.SetFeePool(ctx, feePool)

	// delete outstanding
	h.k.DeleteDefiOutstandingRewards(ctx, defiAddr)

	// remove commission record
	h.k.DeleteDefiAccumulatedCommission(ctx, defiAddr)

	// clear historical rewards
	h.k.DeleteDefiHistoricalRewards(ctx, defiAddr)

	// clear current rewards
	h.k.DeleteDefiCurrentRewards(ctx, defiAddr)
}

// increment period
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) {
	defi := h.k.Defi(ctx, defiAddr)
	h.k.IncrementDefiPeriod(ctx, defi)
}

// withdraw delegation rewards (which also increments period)
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) {
	defi := h.k.Defi(ctx, defiAddr)
	del := h.k.Delegation(ctx, delAddr, defiAddr)

	if _, err := h.k.withdrawDelegationRewards(ctx, defi, del); err != nil {
		panic(err)
	}
}

// create new delegation period record
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) {
	h.k.initializeDelegation(ctx, defiAddr, delAddr)
}

func (h Hooks) BeforeDefiModified(_ sdk.Context, _ sdk.ValAddress)                         {}
func (h Hooks) AfterDefiBonded(_ sdk.Context,  _ sdk.ValAddress)                           {}
func (h Hooks) AfterDefiBeginUnbonding(_ sdk.Context,  _ sdk.ValAddress)                   {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)  {}
