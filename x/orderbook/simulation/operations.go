package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	gaussimappparams "github.com/gauss/gauss/v4/simapp/params"
	"github.com/gauss/gauss/v4/x/orderbook/keeper"
	"github.com/gauss/gauss/v4/x/orderbook/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgCreatePool      = "op_weight_msg_create_pool"
	OpWeightMsgAddPledge       = "op_weight_msg_add_pledge"
	OpWeightMsgRedeemPledge    = "op_weight_msg_redeem_pledge"
	OpWeightMsgPlaceOrder      = "op_weight_msg_place_order"
	OpWeightMsgRevokeOrder     = "op_weight_msg_revoke_order"
	OpWeightMsgAgreeOrderPair  = "op_weight_msg_agree_order_pair"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams, cdc codec.JSONMarshaler, ak types.AccountKeeper,
	bk types.BankKeeper, k keeper.Keeper,
) simulation.WeightedOperations {
	var (
		weightMsgCreatePool      int
		weightMsgAddPledge       int
		weightMsgRedeemPledge    int
		weightMsgPlaceOrder      int
		weightMsgRevokeOrder     int
		weightMsgAgreeOrderPair  int
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgCreatePool, &weightMsgCreatePool, nil,
		func(_ *rand.Rand) {
			weightMsgCreatePool = gaussimappparams.DefaultWeightMsgCreatePool
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgAddPledge, &weightMsgAddPledge, nil,
		func(_ *rand.Rand) {
			weightMsgAddPledge = gaussimappparams.DefaultWeightMsgAddPledge
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgRedeemPledge, &weightMsgRedeemPledge, nil,
		func(_ *rand.Rand) {
			weightMsgRedeemPledge = gaussimappparams.DefaultWeightMsgRedeemPledge
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgPlaceOrder, &weightMsgPlaceOrder, nil,
		func(_ *rand.Rand) {
			weightMsgPlaceOrder = gaussimappparams.DefaultWeightMsgPlaceOrder
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgRevokeOrder, &weightMsgRevokeOrder, nil,
		func(_ *rand.Rand) {
			weightMsgRevokeOrder = gaussimappparams.DefaultWeightMsgRevokeOrder
		},
	)

	appParams.GetOrGenerate(cdc, OpWeightMsgAgreeOrderPair, &weightMsgAgreeOrderPair, nil,
		func(_ *rand.Rand) {
			weightMsgAgreeOrderPair = gaussimappparams.DefaultWeightMsgAgreeOrderPair
		},
	)

	return simulation.WeightedOperations{
		simulation.NewWeightedOperation(
			weightMsgCreatePool,
			SimulateMsgCreatePool(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgAddPledge,
			SimulateMsgAddPledge(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgRedeemPledge,
			SimulateMsgRedeemPledge(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgPlaceOrder,
			SimulateMsgPlaceOrder(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgRevokeOrder,
			SimulateMsgRevokeOrder(ak, bk, k),
		),
		simulation.NewWeightedOperation(
			weightMsgAgreeOrderPair,
			SimulateMsgAgreeOrderPair(ak, bk, k),
		),
	}
}

// SimulateMsgCreatePool generates a MsgCreatePool with random values
// nolint: interfacer
func SimulateMsgCreatePool(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// msg := types.NewMsgCreatePool()
		msg := &types.MsgCreatePool{}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgAddPledge generates a MsgAddPledge with random values
// nolint: interfacer
func SimulateMsgAddPledge(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// msg := types.NewMsgAddPledge()
		msg := &types.MsgAddPledge{}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgRedeemPledge generates a MsgRedeemPledge with random values
// nolint: interfacer
func SimulateMsgRedeemPledge(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// msg := types.NewMsgRedeemPledge()
		msg := &types.MsgRedeemPledge{}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgPlaceOrder generates a MsgPlaceOrder with random values
// nolint: interfacer
func SimulateMsgPlaceOrder(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// msg := types.NewMsgPlaceOrder()
		msg := &types.MsgPlaceOrder{}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgRevokeOrder generates a MsgRevokeOrder with random values
// nolint: interfacer
func SimulateMsgRevokeOrder(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// msg := types.NewMsgRevokeOrder()
		msg := &types.MsgRevokeOrder{}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// SimulateMsgAgreeOrderPair generates a MsgAgreeOrderPair with random values
// nolint: interfacer
func SimulateMsgAgreeOrderPair(ak types.AccountKeeper, bk types.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		// msg := types.NewMsgAgreeOrderPair()
		msg := &types.MsgAgreeOrderPair{}
		return simtypes.NewOperationMsg(msg, true, ""), nil, nil
	}
}
