package simulation_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/simapp"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/gauss/gauss/v4/x/orderbook/simulation"
	"github.com/gauss/gauss/v4/x/orderbook/types"
)

// TestWeightedOperations tests the weights of the operations.
func TestWeightedOperations(t *testing.T) {

	app, ctx := createTestApp(false)

	ctx.WithChainID("test-chain")

	cdc := app.AppCodec()
	appParams := make(simtypes.AppParams)

	weightesOps := simulation.WeightedOperations(appParams, cdc, app.AccountKeeper,
		app.BankKeeper, app.StakingKeeper,
	)

	s := rand.NewSource(1)
	r := rand.New(s)
	accs := simtypes.RandomAccounts(r, 3)

	expected := []struct {
		weight     int
		opMsgRoute string
		opMsgName  string
	}{{simappparams.DefaultWeightMsgCreatePool, types.ModuleName, types.TypeMsgCreatePool},
		{simappparams.DefaultWeightMsgAddPledge, types.ModuleName, types.TypeMsgAddPledge},
		{simappparams.DefaultWeightMsgRedeemPledge, types.ModuleName, types.TypeMsgRedeemPledge},
		{simappparams.DefaultWeightMsgPalceOrder, types.ModuleName, types.TypeMsgPlaceOrder},
		{simappparams.DefaultWeightMsgRevokeOrder, types.ModuleName, types.TypeMsgRevokeOrder},
		{simappparams.DefaultWeightMsgAgreeOrder, types.ModuleName, types.TypeMsgAgreeOrder},
	}

	for i, w := range weightesOps {
		operationMsg, _, _ := w.Op()(r, app.BaseApp, ctx, accs, ctx.ChainID())
		// the following checks are very much dependent from the ordering of the output given
		// by WeightedOperations. if the ordering in WeightedOperations changes some tests
		// will fail
		require.Equal(t, expected[i].weight, w.Weight(), "weight should be the same")
		require.Equal(t, expected[i].opMsgRoute, operationMsg.Route, "route should be the same")
		require.Equal(t, expected[i].opMsgName, operationMsg.Name, "operation Msg name should be the same")
	}
}

// TestSimulateMsgCreatePool tests the normal scenario of a valid message of type TypeMsgCreatePool
// Abonormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgCreatePool(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
}

// TestSimulateMsgAddPledge tests the normal scenario of a valid message of type TypeMsgAddPledge.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgAddPledge(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
}

// TestSimulateMsgRedeemPledge tests the normal scenario of a valid message of type TypeMsgRedeemPledge.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgRedeemPledge(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
}

// TestSimulateMsgPlaceOrder tests the normal scenario of a valid message of type TypeMsgPlaceOrder.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgPlaceOrder(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
}

// TestSimulateMsgRevokeOrder tests the normal scenario of a valid message of type TypeMsgRevokeOrder.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func TestSimulateMsgRevokeOrder(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(5)
	r := rand.New(s)

}

// TestSimulateMsgAgreeOrder tests the normal scenario of a valid message of type TypeMsgAgreeOrder.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func TestSimulateMsgAgreeOrder(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(5)
	r := rand.New(s)

}

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	// sdk.PowerReduction = sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	app := simapp.Setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, minttypes.DefaultParams())
	app.MintKeeper.SetMinter(ctx, minttypes.DefaultInitialMinter())

	return app, ctx
}

func getTestingAccounts(t *testing.T, r *rand.Rand, app *simapp.SimApp, ctx sdk.Context, n int) []simtypes.Account {
	accounts := simtypes.RandomAccounts(r, n)

	initAmt := sdk.TokensFromConsensusPower(200)
	initCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, initAmt))

	// add coins to the accounts
	for _, account := range accounts {
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, account.Address)
		app.AccountKeeper.SetAccount(ctx, acc)
		err := app.BankKeeper.SetBalances(ctx, account.Address, initCoins)
		require.NoError(t, err)
	}

	return accounts
}
