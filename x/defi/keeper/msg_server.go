package keeper

import (
	"context"
	"time"

	metrics "github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gauss/gauss/v4/x/defi/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CreateDefi(goCtx context.Context, msg *types.MsgCreateDefi) (*types.MsgCreateDefiResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	defiAddr, err := sdk.ValAddressFromBech32(msg.DefiAddress)
	if err != nil {
		return nil, err
	}

	// check to see if the pubkey or sender has been registered before
	if _, found := k.GetDefi(ctx, defiAddr); found {
		return nil, types.ErrDefiOwnerExists
	}

	bondDenom := k.BondDenom(ctx)
	if msg.Value.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(types.ErrBadDenom, "got %s, expected %s", msg.Value.Denom, bondDenom)
	}

	if _, err := msg.Description.EnsureLength(); err != nil {
		return nil, err
	}

	defi, err := types.NewDefi(defiAddr, msg.Description)
	if err != nil {
		return nil, err
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	defi.MinSelfDelegation = msg.MinSelfDelegation

	k.SetDefi(ctx, defi)

	// call the after-creation hook
	k.AfterDefiCreated(ctx, defi.GetOperator())

	// move coins from the msg.Address account to a (self-delegation) delegator account
	// the defi account and global shares are updated within here
	// NOTE source will always be from a wallet which are unbonded
	_, err = k.Keeper.Delegate(ctx, delegatorAddress, msg.Value.Amount, types.Unbonded, defi, true)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateDefi,
			sdk.NewAttribute(types.AttributeKeyDefi, msg.DefiAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Value.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})

	return &types.MsgCreateDefiResponse{}, nil
}

func (k msgServer) EditDefi(goCtx context.Context, msg *types.MsgEditDefi) (*types.MsgEditDefiResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	defiAddr, err := sdk.ValAddressFromBech32(msg.DefiAddress)
	if err != nil {
		return nil, err
	}
	// defi must already be registered
	defi, found := k.GetDefi(ctx, defiAddr)
	if !found {
		return nil, types.ErrNoDefiFound
	}

	// replace all editable fields (clients should autofill existing values)
	description, err := defi.Description.UpdateDescription(msg.Description)
	if err != nil {
		return nil, err
	}

	defi.Description = description

	if msg.MinSelfDelegation != nil {
		if !msg.MinSelfDelegation.GT(defi.MinSelfDelegation) {
			return nil, types.ErrMinSelfDelegationDecreased
		}

		if msg.MinSelfDelegation.GT(defi.Tokens) {
			return nil, types.ErrSelfDelegationBelowMinimum
		}

		defi.MinSelfDelegation = (*msg.MinSelfDelegation)
	}

	k.SetDefi(ctx, defi)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEditDefi,
			sdk.NewAttribute(types.AttributeKeyMinSelfDelegation, defi.MinSelfDelegation.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DefiAddress),
		),
	})

	return &types.MsgEditDefiResponse{}, nil
}


func (k msgServer) Delegate(goCtx context.Context, msg *types.MsgDefiDelegate) (*types.MsgDefiDelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	defiAddr, valErr := sdk.ValAddressFromBech32(msg.DefiAddress)
	if valErr != nil {
		return nil, valErr
	}

	defi, found := k.GetDefi(ctx, defiAddr)
	if !found {
		return nil, types.ErrNoDefiFound
	}

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}

	bondDenom := k.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(types.ErrBadDenom, "got %s, expected %s", msg.Amount.Denom, bondDenom)
	}

	// NOTE: source funds are always unbonded
	_, err = k.Keeper.Delegate(ctx, delegatorAddress, msg.Amount.Amount, types.Unbonded, defi, true)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "delegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", msg.Type()},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegate,
			sdk.NewAttribute(types.AttributeKeyDefi, msg.DefiAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})

	return &types.MsgDefiDelegateResponse{}, nil
}

func (k msgServer) Undelegate(goCtx context.Context, msg *types.MsgDefiUndelegate) (*types.MsgDefiUndelegateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	addr, err := sdk.ValAddressFromBech32(msg.DefiAddress)
	if err != nil {
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	shares, err := k.DefiUnbondAmount(
		ctx, delegatorAddress, addr, msg.Amount.Amount,
	)
	if err != nil {
		return nil, err
	}

	bondDenom := k.BondDenom(ctx)
	if msg.Amount.Denom != bondDenom {
		return nil, sdkerrors.Wrapf(types.ErrBadDenom, "got %s, expected %s", msg.Amount.Denom, bondDenom)
	}

	completionTime, err := k.Keeper.Undelegate(ctx, delegatorAddress, addr, shares)
	if err != nil {
		return nil, err
	}

	if msg.Amount.Amount.IsInt64() {
		defer func() {
			telemetry.IncrCounter(1, types.ModuleName, "undelegate")
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", msg.Type()},
				float32(msg.Amount.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel("denom", msg.Amount.Denom)},
			)
		}()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnbond,
			sdk.NewAttribute(types.AttributeKeyDefi, msg.DefiAddress),
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.Amount.String()),
			sdk.NewAttribute(types.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	})

	return &types.MsgDefiUndelegateResponse{
		CompletionTime: completionTime,
	}, nil
}

func (k msgServer) SetWithdrawAddress(goCtx context.Context, msg *types.MsgSetDefiWithdrawAddress) (*types.MsgSetDefiWithdrawAddressResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	withdrawAddress, err := sdk.AccAddressFromBech32(msg.WithdrawAddress)
	if err != nil {
		return nil, err
	}
	err = k.SetWithdrawAddr(ctx, delegatorAddress, withdrawAddress)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	)

	return &types.MsgSetDefiWithdrawAddressResponse{}, nil
}

func (k msgServer) WithdrawDelegatorReward(goCtx context.Context, msg *types.MsgWithdrawDefiDelegatorReward) (*types.MsgWithdrawDefiDelegatorRewardResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	defiAddr, err := sdk.ValAddressFromBech32(msg.DefiAddress)
	if err != nil {
		return nil, err
	}
	delegatorAddress, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	amount, err := k.WithdrawDelegationRewards(ctx, delegatorAddress, defiAddr)
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "withdraw_reward"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DelegatorAddress),
		),
	)
	return &types.MsgWithdrawDefiDelegatorRewardResponse{}, nil
}

func (k msgServer) WithdrawDefiCommission(goCtx context.Context, msg *types.MsgWithdrawDefiCommission) (*types.MsgWithdrawDefiCommissionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	defiAddr, err := sdk.ValAddressFromBech32(msg.DefiAddress)
	if err != nil {
		return nil, err
	}
	amount, err := k.Keeper.WithdrawDefiCommission(ctx, defiAddr)
	if err != nil {
		return nil, err
	}

	defer func() {
		for _, a := range amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "withdraw_commission"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.DefiAddress),
		),
	)

	return &types.MsgWithdrawDefiCommissionResponse{}, nil
}

func (k msgServer) FundCommunityPool(goCtx context.Context, msg *types.MsgFundDefiCommunityPool) (*types.MsgFundDefiCommunityPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	depositer, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil, err
	}
	if err := k.Keeper.FundCommunityPool(ctx, msg.Amount, depositer); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Depositor),
		),
	)

	return &types.MsgFundDefiCommunityPoolResponse{}, nil
}
