package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/gauss/gauss/v4/x/token/types"
)

// Simulation parameter constants
const (
	CommuniteTax      = "communite_tax"
	IssueTokenFee     = "issue_token_fee"
	MintTokenFeeRatio = "mint_token_fee_ratio"
)

// RandomDec randomized sdk.RandomDec
func RandomDec(r *rand.Rand) sdk.Dec {
	return sdk.NewDec(r.Int63())
}

// RandomInt randomized sdk.Int
func RandomInt(r *rand.Rand) sdk.Int {
	return sdk.NewInt(r.Int63())
}

// RandomizedGenState generates a random GenesisState for bank
func RandomizedGenState(simState *module.SimulationState) {

	var communiteTax sdk.Dec
	var issueTokenFee sdk.Int
	var mintTokenFeeRatio sdk.Dec
	var tokens []types.Token

	simState.AppParams.GetOrGenerate(
		simState.Cdc, CommuniteTax, &communiteTax, simState.Rand,
		func(r *rand.Rand) {
			communiteTax = sdk.NewDecWithPrec(int64(r.Intn(5)), 1)
		},
	)
	simState.AppParams.GetOrGenerate(
		simState.Cdc, IssueTokenFee, &issueTokenFee, simState.Rand,
		func(r *rand.Rand) {
			issueTokenFee = sdk.NewInt(int64(10))

			for i := 0; i < 5; i++ {
				tokens = append(tokens, randToken(r, simState.Accounts))
			}
		},
	)
	simState.AppParams.GetOrGenerate(
		simState.Cdc, MintTokenFeeRatio, &mintTokenFeeRatio, simState.Rand,
		func(r *rand.Rand) { 
			mintTokenFeeRatio = sdk.NewDecWithPrec(int64(r.Intn(5)), 1)
		},
	)

	gs := types.NewGenesisState(
		types.NewParams(communiteTax, sdk.NewCoin(sdk.DefaultBondDenom, issueTokenFee),
			mintTokenFeeRatio,
		),
		tokens,
		sdk.Coins{},
	)

	bz, err := json.MarshalIndent(&gs, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", types.ModuleName, bz)

	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(gs)
}
