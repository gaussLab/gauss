package token_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	tokenmodule "github.com/gauss/gauss/v4/x/token"
	tokenkeeper "github.com/gauss/gauss/v4/x/token/keeper"
	"github.com/gauss/gauss/v4/x/token/types"
	"github.com/gauss/gauss/v4/simapp"
)

const (
	isCheckTx = false
)

var (
	denom       = "ugauss"
	owner       = sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	initAmount  = sdk.NewIntWithDecimal(100000000, int(6))
	initCoin    = sdk.Coins{sdk.NewCoin(denom, initAmount)}
)

type HandlerSuite struct {
	suite.Suite

	cdc    codec.JSONMarshaler
	ctx    sdk.Context
	keeper tokenkeeper.Keeper
	bk     bankkeeper.Keeper
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerSuite))
}

func (suite *HandlerSuite) SetupTest() {
	app := simapp.Setup(isCheckTx)

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.TokenKeeper
	suite.bk = app.BankKeeper

	// set params
	suite.keeper.SetParam(suite.ctx, types.DefaultParams())

	// init tokens to addr
	err := suite.bk.MintCoins(suite.ctx, types.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func (suite *HandlerSuite) issueToken(token types.Token) {
	err := suite.keeper.AddToken(suite.ctx, token)
	suite.NoError(err)

	mintCoins := sdk.NewCoins(
		sdk.NewCoin(
			token.GetSmallestUnit(),
			sdk.NewIntWithDecimal(int64(token.GetInitialSupply()), 0)),
		),
	)

	err = suite.bk.MintCoins(suite.ctx, types.ModuleName, mintCoins)
	suite.NoError(err)

	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, mintCoins)
	suite.NoError(err)
}

func (suite *HandlerSuite) TestIssueUnlockedToken() {
	// issue token
	h := tokenmodule.NewHandler(suite.keeper)

	msg := types.NewMsgIssueToken("Ethereum Network", "eth", "ueth", 6, 1000000000, 10000000000, false, true, owner.String())
	_, err := h(suite.ctx, msg)
	suite.NoError(err)

	// check
	amount1 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount
	amount2 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount
	mode, coin, err := suite.keeper.GetTokenIssueParams(suite.ctx)
	suite.NoError(err)
	suite.Equal(amount1.Sub(coin.Amount), amount2)

	// check
	mintTokenAmount := sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Decimals))
	amount3 := suite.bk.GetBalance(suite.ctx, owner, msg.Symbol).Amount
	suite.Equal(amount3, mintTokenAmount)
}

func (suite *HandlerSuite) TestIssueLockedToken() {
	// issue token
	h := tokenmodule.NewHandler(suite.keeper)

	msg := types.NewMsgIssueToken("Ethereum Network", "eth", "ueth", 6, 1000000000, 10000000000, false, false, owner.String())
	_, err := h(suite.ctx, msg)
	suite.NoError(err)

	// check
	amount1 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount
	amount2 := suite.bk.GetBalance(suite.ctx, owner, denom).Amount
	mode, coin, err := suite.keeper.GetTokenIssueParams(suite.ctx)
	suite.NoError(err)
	suite.Equal(amount1.Sub(coin.Amount), amount2)

	// check
	mintTokenAmount := sdk.NewIntWithDecimal(int64(msg.InitialSupply), int(msg.Decimals))
	amount3 := suite.bk.GetBalance(suite.ctx, owner, msg.Symbol).Amount
	suite.Equal(amount3, mintTokenAmount)
}

func (suite *HandlerSuite) TestMintToken() {
	// issue token
	token := types.NewToken("Ethereum Network", "eth", "ueth", 6, 1000000000, 10000000000, true, owner)
	suite.issueToken(token)

	// before: native token amount
	nativeAmount1 := suite.bk.GetBalance(suite.ctx, token.GetOwnerString(), denom).Amount

	// begin: token amount
	tokenAmount1  := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetSmallestUnit()).Amount
	suite.Equal(sdk.NewIntWithDecimal(int64(token.GetInitialSupply()), int(token.GetDecimals())), tokenAmount1)

	// mint
	h := tokenmodule.NewHandler(suite.keeper)
	msgMintToken := types.NewMsgMintToken(token.GetSymbol(), token.GetOwner(), "", 1000)
	_, err := h(suite.ctx, msgMintToken)
	suite.NoError(err)

	// end: token amount
	tokenAmount2 := suite.bk.GetBalance(suite.ctx, token.GetOwnerString(), token.GetSmallestUnit()).Amount

	// check
	mintTokenAmount := sdk.NewIntWithDecimal(int64(msgMintToken.Amount), int(token.GetDecimals()))
	suite.Equal(tokenAmount1.Add(mintTokenAmount), tokenAmount2)

	// check
	coin, err := suite.keeper.GetTokenMintParams(suite.ctx)
	suite.NoError(err)

	nativeAmount2 := suite.bk.GetBalance(suite.ctx, token.GetOwner(), denom).Amount
	suite.Equal(nativeAmount1.Sub(coin.Amount), nativeAmount2)
}

func (suite *HandlerSuite) TestBurnToken() {
	// issue token
	token := types.NewToken("Ethereum Network", "eth", "ueth", 6, 1000000000, 10000000000, true, owner)
	suite.issueToken(token)

	// before
	tokenAmount1 := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetSmallestUnit()).Amount
	suite.Equal(sdk.NewIntWithDecimal(int64(token.GetInitialSupply()), int(token.GetDecimals())), tokenAmount1)

	// burn
	h := tokenmodule.NewHandler(suite.keeper)
	msgBurnToken := types.NewMsgBurnToken(token.Symbol, token.Owner, 200)
	_, err := h(suite.ctx, msgBurnToken)
	suite.NoError(err)

	// end
	tokenAmount2 := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetSmallestUnit()).Amount
	burnTokenAmount := sdk.NewIntWithDecimal(int64(msgBurnToken.Amount), int(token.GetDecimals()))

	// check
	suite.Equal(tokenAmount1.Sub(burnTokenAmount), tokenAmount2)
}
