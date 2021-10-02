package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/tendermint/tendermint/crypto/tmhash"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	"github.com/gauss/gauss/v4/x/token/types"
	"github.com/gauss/gauss/v4/simapp"
)

const (
	isCheckTx = false
)

var (
	denom		= types.GetNativeToken().Symbol
	owner		= sdk.AccAddress(tmhash.SumTruncated([]byte("tokenTest")))
	initAmount	= sdk.NewIntWithDecimal(100000000, int(6))
	initCoin	= sdk.Coins{sdk.NewCoin(denom, initAmount)}
)

type KeeperTestSuite struct {
	suite.Suite

	legacyAmino *codec.LegacyAmino
	ctx         sdk.Context
	keeper      keeper.Keeper
	bk          bankkeeper.Keeper
	app         *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(isCheckTx)

	suite.legacyAmino = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	suite.keeper = app.TokenKeeper
	suite.bk = app.BankKeeper
	suite.app = app

	// set params
	suite.keeper.SetParam(suite.ctx, types.DefaultParams())

	// init tokens to addr
	err := suite.bk.MintCoins(suite.ctx, types.ModuleName, initCoin)
	suite.NoError(err)
	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, initCoin)
	suite.NoError(err)
}

func TestKeeperSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) setToken(token types.Token) {
	err := suite.keeper.AddToken(suite.ctx, token)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) issueToken(token types.Token) {
	suite.setToken(token)

	mintCoins := sdk.NewCoins(
		sdk.NewCoin(
			token.GetSmallestUnit(),
			sdk.NewIntWithDecimal(int64(token.InitialSupply), 0),
		),
	)

	err := suite.bk.MintCoins(suite.ctx, types.ModuleName, mintCoins)
	suite.NoError(err)

	err = suite.bk.SendCoinsFromModuleToAccount(suite.ctx, types.ModuleName, owner, mintCoins)
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestIssueUnlockedToken() {
	token := types.NewToken("Ethereum Network", "eth", "ueth", 6, 10000000000, 100000000000, false, owner)

	err := suite.keeper.IssueToken(
		suite.ctx, token.GetName(), token.GetSymbol(), token.GetSmallestUnit(),
		token.GetDecimals(), token.GetInitialSupply(), token.GetTotalSupply(),
		token.GetMintable(), true, token.GetOwner(),
	)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, token.GetSymbol()))

	issuedToken, err := suite.keeper.GetToken(suite.ctx, token.GetSymbol())
	suite.NoError(err)

	suite.Equal(token.Owner, issuedToken.GetOwnerString())
	suite.EqualValues(&token, issuedToken.(*types.Token))
}

func (suite *KeeperTestSuite) TestIssueLockedToken() {
	token := types.NewToken("Ethereum Network", "eth", "ueth", 6, 10000000000, 100000000000, false, owner)

	err := suite.keeper.IssueToken(
		suite.ctx, token.GetName(), token.GetSymbol(), token.GetSmallestUnit(),
		token.GetDecimals(), token.GetInitialSupply(), token.GetTotalSupply(),
		token.GetMintable(), false, token.GetOwner(),
	)
	suite.NoError(err)

	suite.True(suite.keeper.HasToken(suite.ctx, token.GetSymbol()))

	issuedToken, err := suite.keeper.GetToken(suite.ctx, token.GetSymbol())
	suite.NoError(err)

	suite.Equal(token.Owner, issuedToken.GetOwnerString())
	suite.EqualValues(&token, issuedToken.(*types.Token))
}

func (suite *KeeperTestSuite) TestEditToken() {
	token := types.NewToken("Ethereum Network", "eth", "ueth", 6, 10000000000, 100000000000, false, owner)
	suite.setToken(token)

	symbol := "eth"
	name := "Ethereum ERC20"
	mintable := types.True
	totalSupply := uint64(200000000000)

	err := suite.keeper.EditToken(suite.ctx, symbol, name, totalSupply, mintable, owner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetToken(suite.ctx, symbol)
	suite.NoError(err)

	expectToken := types.NewToken("Ethereum ERC20", "eth", "ueth", 6, 12000000000, 120000000000, mintable.ToBool(), owner)
	suite.EqualValues(newToken.(*types.Token), &expectToken)
}

func (suite *KeeperTestSuite) TestMintToken() {
	token := types.NewToken("Ethereum Network", "eth", "ueth", 6, 1000000000, 2000000000, true, owner)
	suite.issueToken(token)

	amount := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetSmallestUnit())
	suite.Equal("100000000eth", amount.String())

	mintAmount := uint64(1000000000)
	recipient := sdk.AccAddress{}

	err := suite.keeper.MintToken(suite.ctx, token.GetSymbol(), mintAmount, recipient, token.GetOwner())
	suite.NoError(err)

	amount = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetSmallestUnit())
	suite.Equal("2000000000", amount.String())

	// mint token without owner
	err = suite.keeper.MintToken(suite.ctx, token.GetSymbol(), mintAmount, owner, sdk.AccAddress{})
	suite.Error(err, "can not mint token without owner when the owner exists")

	token = types.NewToken("Cosmos Hub", "atom", "uatom", 6, 1000000000, 2000000000, true, sdk.AccAddress{})
	suite.issueToken(token)

	err = suite.keeper.MintToken(suite.ctx, token.GetSymbol(), mintAmount, owner, sdk.AccAddress{})
	suite.NoError(err)
}

func (suite *KeeperTestSuite) TestBurnToken() {
	token := types.NewToken("Ethereum Network", "eth", "ueth", 6, 1000000000, 2000000000, true, owner)
	suite.issueToken(token)

	amount := suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetSmallestUnit())
	suite.Equal("1000000000eth", amount.String())

	burnedAmount := uint64(200000000)

	err := suite.keeper.BurnToken(suite.ctx, token.GetSymbol(), burnedAmount, token.GetOwner())
	suite.NoError(err)

	amt = suite.bk.GetBalance(suite.ctx, token.GetOwner(), token.GetSmallestUnit())
	suite.Equal("8000000000eth", amount.String())
}

func (suite *KeeperTestSuite) TestTransferToken() {
	token := types.NewToken("Ethereum Network", "eth", "ueth", 6, 10000000000, 100000000000, false, owner)
	suite.setToken(token)

	dstOwner := sdk.AccAddress(tmhash.SumTruncated([]byte("TokenNewOwner")))

	err := suite.keeper.TransferTokenOwner(suite.ctx, token.Symbol, token.GetOwner(), dstOwner)
	suite.NoError(err)

	newToken, err := suite.keeper.GetToken(suite.ctx, token.Symbol)
	suite.NoError(err)

	suite.Equal(dstOwner, newToken.GetOwner())
}
