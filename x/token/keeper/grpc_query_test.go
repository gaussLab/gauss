package keeper_test

import (
	gocontext "context"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gauss/gauss/v4/x/token/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryToken() {
	app, ctx := suite.app, suite.ctx

	_, _, addr := testdata.KeyTestPubAddr()
	token := types.NewToken("Test Network", "ttk", "uttk", 6, 1000000000, 10000000000, true, addr)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.TokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_ = suite.app.TokenKeeper.AddToken(ctx, token)

	// Query token
	tokenResp, err := queryClient.Token(gocontext.Background(), &types.QueryTokenRequest{Denom: "uttk"})
	suite.Require().NoError(err)
	suite.Require().NotNil(tokenResp)

	// Query tokens
	tokensResp, err := queryClient.Tokens(gocontext.Background(), &types.QueryTokensRequest{})
	suite.Require().NoError(err)
	suite.Require().NotNil(tokensResp)
	suite.Len(tokensResp.Tokens, 2)
}

func (suite *KeeperTestSuite) TestGRPCQueryFees() {
	app, ctx := suite.app, suite.ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.TokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, err := queryClient.Fees(gocontext.Background(), &types.QueryFeesRequest{Symbol: "uttk"})
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestGRPCQueryParams() {
	app, ctx := suite.app, suite.ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.TokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	paramsResp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	params := app.TokenKeeper.GetParam(ctx)
	suite.Require().NoError(err)
	suite.Equal(params, paramsResp.Params)
}

func (suite *KeeperTestSuite) TestGRPCQueryTotalBurn() {
	app, ctx := suite.app, suite.ctx

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, app.TokenKeeper)
	queryClient := types.NewQueryClient(queryHelper)

	_, _, addr := testdata.KeyTestPubAddr()
	token := types.NewToken("Test Network", "ttk", "uttk", 6, 1000000000, 10000000000, true, addr)
	err := suite.app.TokenKeeper.AddToken(ctx, token)
	suite.Require().NoError(err)

	burnCoin := sdk.NewInt64Coin("uttk", 10000000)
	app.TokenKeeper.AddBurnCoin(ctx, burnCoin)

	resp, err := queryClient.TotalBurn(gocontext.Background(), &types.QueryTotalBurnRequest{})
	suite.Require().NoError(err)
	suite.Len(resp.BurnedCoins, 1)
	suite.EqualValues(burnCoin, resp.BurnedCoins[0])
}
