package simulation

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"

	"github.com/gauss/gauss/v4/x/token/keeper"
	"github.com/gauss/gauss/v4/x/token/types"
)

// Simulation operation weights constants
const (
	OpWeightMsgIssueToken         = "op_weight_msg_issue_token"
	OpWeightMsgEditToken          = "op_weight_msg_edit_token"
	OpWeightMsgMintToken          = "op_weight_msg_mint_token"
	OpWeightMsgTransferTokenOwner = "op_weight_msg_transfer_token_owner"
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	cdc codec.JSONMarshaler,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
) simulation.WeightedOperations {

	var weightIssue, weightEdit, weightMint, weightTransfer int
	appParams.GetOrGenerate(
		cdc, OpWeightMsgIssueToken, &weightIssue, nil,
		func(_ *rand.Rand) {
			weightIssue = 100
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgEditToken, &weightEdit, nil,
		func(_ *rand.Rand) {
			weightEdit = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgMintToken, &weightMint, nil,
		func(_ *rand.Rand) {
			weightMint = 50
		},
	)

	appParams.GetOrGenerate(
		cdc, OpWeightMsgTransferTokenOwner, &weightTransfer, nil,
		func(_ *rand.Rand) {
			weightTransfer = 50
		},
	)

	return simulation.WeightedOperations{
		//simtypes.NewWeightedOperation(
		//	weightIssue,
		//	SimulateIssueToken(k, ak),
		//),
		simulation.NewWeightedOperation(
			weightEdit,
			SimulateEditToken(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightMint,
			SimulateMintToken(k, ak, bk),
		),
		simulation.NewWeightedOperation(
			weightTransfer,
			SimulateTransferTokenOwner(k, ak, bk),
		),
	}
}

// SimulateIssueToken tests and runs a single msg issue a new token
func SimulateIssueToken(k keeper.Keeper, ak authkeeper.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		token, maxFees := genToken(ctx, r, k, ak, bk, accs)

		msg := types.NewMsgIssueToken(token.GetName(), token.GetSymbol(), token.GetSmallestUnit(), token.GetDecimals(), 
			token.GetInitialSupply(), token.GetTotalSupply(), token.GetMintable(), true, token.GetOwnerString())

		simAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), fmt.Sprintf("account[%s] does not found", token.GetOwnerString())), 
				nil, fmt.Errorf("account[%s] does not found", token.GetOwnerString())
		}

		owner, _ := sdk.AccAddressFromBech32(msg.Owner)
		account := ak.GetAccount(ctx, owner)
		fees, err := simtypes.RandomFees(r, ctx, maxFees)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate issue token"), nil, nil
	}
}

// SimulateEditToken tests and runs a single msg edit a existed token
func SimulateEditToken(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		token, _ := selectToken(ctx, k, ak, bk, false)

		msg := types.NewMsgEditToken(token.GetSymbol(), true, token.GetOwnerString())

		simAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), fmt.Sprintf("account[%s] does not found", token.GetOwnerString())), 
				nil, fmt.Errorf("account[%s] does not found", token.GetOwnerString())
		}

		owner, _ := sdk.AccAddressFromBech32(msg.Owner)
		account := ak.GetAccount(ctx, owner)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate edit token"), nil, nil
	}
}

// SimulateMintToken tests and runs a single msg mint a existed token
func SimulateMintToken(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		token, maxFee := selectToken(ctx, k, ak, bk, true)
		simToAccount, _ := simtypes.RandomAcc(r, accs)

		msg := types.NewMsgMintToken(token.GetSymbol(), token.GetOwnerString(), simToAccount.Address.String(), 100)

		ownerAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), fmt.Sprintf("account[%s] does not found", token.GetOwnerString())), 
				nil, fmt.Errorf("account[%s] does not found", token.GetOwnerString())
		}

		owner, _ := sdk.AccAddressFromBech32(msg.Owner)
		account := ak.GetAccount(ctx, owner)
		fees, err := simtypes.RandomFees(r, ctx, maxFee)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			ownerAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate mint token"), nil, nil
	}
}

// SimulateTransferTokenOwner tests and runs a single msg transfer to others
func SimulateTransferTokenOwner(k keeper.Keeper, ak types.AccountKeeper, bk types.BankKeeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		token, _ := selectToken(ctx, k, ak, bk, false)
		var simToAccount, _ = simtypes.RandomAcc(r, accs)
		for simToAccount.Address.Equals(token.GetOwner()) {
			simToAccount, _ = simtypes.RandomAcc(r, accs)
		}

		msg := types.NewMsgTransferTokenOwner(token.GetSymbol(), token.GetOwnerString(), simToAccount.Address.String())

		simAccount, found := simtypes.FindAccount(accs, token.GetOwner())
		if !found {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), fmt.Sprintf("account[%s] does not found", token.GetOwnerString())), 
				nil, fmt.Errorf("account[%s] does not found", token.GetOwnerString())
		}

		srcOwner, _ := sdk.AccAddressFromBech32(msg.OldOwner)
		account := ak.GetAccount(ctx, srcOwner)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate fees"), nil, err
		}

		txGen := simappparams.MakeTestEncodingConfig().TxConfig
		tx, err := helpers.GenTx(
			txGen,
			[]sdk.Msg{msg},
			fees,
			helpers.DefaultGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		if _, _, err = app.Deliver(txGen.TxEncoder(), tx); err != nil {
			return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(msg, true, "simulate transfer token"), nil, nil
	}
}

func genToken(ctx sdk.Context,
	r *rand.Rand,
	k keeper.Keeper,
	ak authkeeper.AccountKeeper,
	bk types.BankKeeper,
	accs []simtypes.Account,
) (types.Token, sdk.Coins) {

	var token types.Token
	token = randToken(r, accs)

	for k.HasToken(ctx, token.GetSymbol()) {
		token = randToken(r, accs)
	}

	issueFee, err := k.GetIssueTokenFee(ctx, token.GetSymbol())
	if err != nil {
		panic(err)
	}

	account, maxFees := filterAccount(ctx, r, ak, bk, accs, issueFee)
	token.Owner = account.String()

	return token, maxFees
}

func selectToken(
	ctx sdk.Context,
	k keeper.Keeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	mint bool,
) (token types.TokenI, maxFees sdk.Coins) {
	tokens := k.GetTokens(ctx, nil)
	if len(tokens) == 0 {
		panic("No token available")
	}

	for _, t := range tokens {
		if !mint {
			return t, nil
		}

		account := ak.GetAccount(ctx, t.GetOwner())
		spendable := bk.SpendableCoins(ctx, account.GetAddress())
		spendableStake := spendable.AmountOf(types.DefaultParamsDenom)
		if spendableStake.IsZero() {
			continue
		}

		maxFees = sdk.NewCoins(sdk.NewCoin(types.DefaultParamsDenom, spendableStake))
		token = t
		return
	}

	panic("No token mintable")
}

func randToken(r *rand.Rand, accs []simtypes.Account) types.Token {

	name := randString(r, 1, types.MaximumNameLen)
	symbol := randString(r, types.MinimumSymbolLen, types.MaximumSymbolLen)
	decimals := simtypes.RandIntBetween(r, 1, int(types.MaximumDecimals))
	initialSupply := r.Int63n(int64(100000000000))
	totalSupply := 2 * initialSupply
	simAccount, _ := simtypes.RandomAcc(r, accs)

	return types.Token{
		Name:          name,
		Symbol:        strings.ToLower(symbol),
		Decimals:       uint32(decimals),
		InitialSupply: uint64(initialSupply),
		TotalSupply:   uint64(totalSupply),
		Mintable:      true,
		Owner:         simAccount.Address.String(),
	}
}

func filterAccount(
	ctx sdk.Context,
	r *rand.Rand,
	ak authkeeper.AccountKeeper,
	bk types.BankKeeper,
	accs []simtypes.Account, fees sdk.Coin,
) (owner sdk.AccAddress, maxFees sdk.Coins) {
loop:
	simAccount, _ := simtypes.RandomAcc(r, accs)
	account := ak.GetAccount(ctx, simAccount.Address)
	spendable := bk.SpendableCoins(ctx, account.GetAddress())
	spendableStake := spendable.AmountOf(types.DefaultParamsDenom)
	if spendableStake.IsZero() || spendableStake.LT(fees.Amount) {
		goto loop
	}
	owner = account.GetAddress()
	maxFees = sdk.NewCoins(sdk.NewCoin(types.DefaultParamsDenom, spendableStake).Sub(fees))
	return
}

func randString(r *rand.Rand, min, max int) string {
	strLen := simtypes.RandIntBetween(r, min, max)
	randStr := simtypes.RandStringOfLength(r, strLen)
	return randStr
}


