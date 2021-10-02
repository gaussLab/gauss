package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/gauss/gauss/v4/x/defi/types"
)

// Implements DefiSet interface
var _ types.DefiSet = Keeper{}

// Implements DelegationSet interface
var _ types.DelegationSet = Keeper{}

// keeper of the staking store
type Keeper struct {
	storeKey           sdk.StoreKey
	cdc                codec.BinaryMarshaler
	authKeeper         types.AccountKeeper
	bankKeeper         types.BankKeeper
	tokenKeeper	   types.TokenKeeper
	paramstore         paramtypes.Subspace
	hooks              types.DefiHooks

	blockedAddrs       map[string]bool
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(
	cdc codec.BinaryMarshaler, key sdk.StoreKey, ps paramtypes.Subspace, ak types.AccountKeeper, bk types.BankKeeper, 
	tk types.TokenKeeper, blockedAddrs map[string]bool,
) Keeper {
	// ensure defi module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	// ensure bonded and not bonded module accounts are set
	if addr := ak.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := ak.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	return Keeper{
		storeKey:           key,
		cdc:                cdc,
		authKeeper:         ak,
		bankKeeper:         bk,
		tokenKeeper:	    tk,
		paramstore:         ps,
		hooks:              nil,
		blockedAddrs:       blockedAddrs,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// Set the defi hooks
func (k *Keeper) SetHooks(sh types.DefiHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set defi hooks twice")
	}

	k.hooks = sh

	return k
}

// SetWithdrawAddr sets a new address that will receive the rewards upon withdrawal
func (k Keeper) SetWithdrawAddr(ctx sdk.Context, delegatorAddr sdk.AccAddress, withdrawAddr sdk.AccAddress) error {
	if k.blockedAddrs[withdrawAddr.String()] {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive external funds", withdrawAddr)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSetWithdrawAddress,
			sdk.NewAttribute(types.AttributeKeyWithdrawAddress, withdrawAddr.String()),
		),
	)

	k.SetDelegatorWithdrawAddr(ctx, delegatorAddr, withdrawAddr)
	return nil
}

// withdraw rewards from a delegation
func (k Keeper) WithdrawDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress, defiAddr sdk.ValAddress) (sdk.Coins, error) {
	defi := k.Defi(ctx, defiAddr)
	if defi == nil {
		return nil, types.ErrNoDefiDistInfo
	}

	del := k.Delegation(ctx, delAddr, defiAddr)
	if del == nil {
		return nil, types.ErrEmptyDelegationDistInfo
	}

	// withdraw rewards
	rewards, err := k.withdrawDelegationRewards(ctx, defi, del)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, rewards.String()),
			sdk.NewAttribute(types.AttributeKeyDefi, defiAddr.String()),
		),
	)

	// reinitialize the delegation
	k.initializeDelegation(ctx, defiAddr, delAddr)
	return rewards, nil
}

// withdraw defi commission
func (k Keeper) WithdrawDefiCommission(ctx sdk.Context, defiAddr sdk.ValAddress) (sdk.Coins, error) {
	// fetch defi accumulated commission
	accumCommission := k.GetDefiAccumulatedCommission(ctx, defiAddr)
	if accumCommission.Commission.IsZero() {
		return nil, types.ErrNoDefiCommission
	}

	commission, remainder := accumCommission.Commission.TruncateDecimal()
	k.SetDefiAccumulatedCommission(ctx, defiAddr, types.DefiAccumulatedCommission{Commission: remainder}) // leave remainder to withdraw later

	// update outstanding
	outstanding := k.GetDefiOutstandingRewards(ctx, defiAddr).Rewards
	k.SetDefiOutstandingRewards(ctx, defiAddr, types.DefiOutstandingRewards{Rewards: outstanding.Sub(sdk.NewDecCoinsFromCoins(commission...))})

	if !commission.IsZero() {
		accAddr := sdk.AccAddress(defiAddr)
		withdrawAddr := k.GetDelegatorWithdrawAddr(ctx, accAddr)
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, commission)
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeWithdrawCommission,
			sdk.NewAttribute(sdk.AttributeKeyAmount, commission.String()),
		),
	)

	return commission, nil
}

// GetTotalRewards returns the total amount of fee defi rewards held in the store
func (k Keeper) GetTotalRewards(ctx sdk.Context) (totalRewards sdk.DecCoins) {
	k.IterateDefiOutstandingRewards(ctx,
		func(_ sdk.ValAddress, rewards types.DefiOutstandingRewards) (stop bool) {
			totalRewards = totalRewards.Add(rewards.Rewards...)
			return false
		},
	)

	return totalRewards
}

// FundCommunityPool allows an account to directly fund the community fund pool.
// The amount is first added to the distribution module account and then directly
// added to the pool. An error is returned if the amount cannot be sent to the
// module account.
func (k Keeper) FundCommunityPool(ctx sdk.Context, amount sdk.Coins, sender sdk.AccAddress) error {
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount); err != nil {
		return err
	}

	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(sdk.NewDecCoinsFromCoins(amount...)...)
	k.SetFeePool(ctx, feePool)

	return nil
}
