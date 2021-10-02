package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// Calculate the DefiUpdates for the current block
// Called in each EndBlock
func (k Keeper) BlockDefiUpdates(ctx sdk.Context) []abci.ValidatorUpdate {
	// Calculate defi set changes.
	//
	// NOTE: ApplyAndReturnDefiSetUpdates has to come before
	// UnbondAllMatureDefiQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the defi to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnDefiSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureDefiQueue).
	defiUpdates, err := k.ApplyAndReturnDefiSetUpdates(ctx)
	if err != nil {
		panic(err)
	}

	// unbond all mature defis from the unbonding queue
	k.UnbondAllMatureDefis(ctx)

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, ddPair := range matureUnbonds {
		addr, err := sdk.ValAddressFromBech32(ddPair.DefiAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress, err := sdk.AccAddressFromBech32(ddPair.DelegatorAddress)
		if err != nil {
			panic(err)
		}
		balances, err := k.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDefi, ddPair.DefiAddress),
				sdk.NewAttribute(types.AttributeKeyDelegator, ddPair.DelegatorAddress),
			),
		)
	}

	return defiUpdates
}

// Apply and return accumulated updates to the bonded defi set. Also,
// * Updates defi status' according to updated powers.
// * Updates the fee pool bonded vs not-bonded tokens.
// * Updates relevant indices.
// It gets called once after genesis, another time maybe after genesis transactions,
// then once at every EndBlock.
//
// CONTRACT: Only defis with non-zero power or zero-power that were bonded
// at the previous block height or were removed from the defi set entirely
// are returned to Tendermint.
func (k Keeper) ApplyAndReturnDefiSetUpdates(ctx sdk.Context) (updates []abci.ValidatorUpdate, err error) {
	return updates, err
}

// Defi state transitions

func (k Keeper) bondedToUnbonding(ctx sdk.Context, defi types.Defi) (types.Defi, error) {
	if !defi.IsBonded() {
		panic(fmt.Sprintf("bad state transition bondedToUnbonding, defi: %v\n", defi))
	}

	return k.beginUnbondingDefi(ctx, defi)
}

func (k Keeper) unbondingToBonded(ctx sdk.Context, defi types.Defi) (types.Defi, error) {
	if !defi.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, defi: %v\n", defi))
	}

	return k.bondDefi(ctx, defi)
}

func (k Keeper) unbondedToBonded(ctx sdk.Context, defi types.Defi) (types.Defi, error) {
	if !defi.IsUnbonded() {
		panic(fmt.Sprintf("bad state transition unbondedToBonded, defi: %v\n", defi))
	}

	return k.bondDefi(ctx, defi)
}

// UnbondingToUnbonded switches a defi from unbonding state to unbonded state
func (k Keeper) UnbondingToUnbonded(ctx sdk.Context, defi types.Defi) types.Defi {
	if !defi.IsUnbonding() {
		panic(fmt.Sprintf("bad state transition unbondingToBonded, defi: %v\n", defi))
	}

	return k.completeUnbondingDefi(ctx, defi)
}

// perform all the store operations for when a defi status becomes bonded
func (k Keeper) bondDefi(ctx sdk.Context, defi types.Defi) (types.Defi, error) {
	defi = defi.UpdateStatus(types.Bonded)

	// save the now bonded defi record to the two referenced stores
	k.SetDefi(ctx, defi)

	// delete from queue if present
	k.DeleteDefiQueue(ctx, defi)

	k.AfterDefiBonded(ctx, defi.GetOperator())

	return defi, nil
}

// perform all the store operations for when a defi begins unbonding
func (k Keeper) beginUnbondingDefi(ctx sdk.Context, defi types.Defi) (types.Defi, error) {
	params := k.GetParams(ctx)

	// sanity check
	if defi.Status != types.Bonded {
		panic(fmt.Sprintf("should not already be unbonded or unbonding, defi: %v\n", defi))
	}

	defi = defi.UpdateStatus(types.Unbonding)

	// set the unbonding completion time and completion height appropriately
	defi.UnbondingTime = ctx.BlockHeader().Time.Add(params.UnbondingTime)
	defi.UnbondingHeight = ctx.BlockHeader().Height

	// save the now unbonded defi record and power index
	k.SetDefi(ctx, defi)

	// Adds to unbonding defi queue
	k.InsertUnbondingDefiQueue(ctx, defi)

	k.AfterDefiBeginUnbonding(ctx, defi.GetOperator())

	return defi, nil
}

// perform all the store operations for when a defi status becomes unbonded
func (k Keeper) completeUnbondingDefi(ctx sdk.Context, defi types.Defi) types.Defi {
	defi = defi.UpdateStatus(types.Unbonded)
	k.SetDefi(ctx, defi)

	return defi
}
