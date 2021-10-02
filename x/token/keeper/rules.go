//nolint
package keeper

import (
	"math"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/gauss/gauss/v4/x/token/types"
)

// factor = (ln(len({data}))/ln{base})^{exp}
const (
	FactorBase = 3
	FactorExp  = 4
)

// DoIssueTokenFee performs fee handling for issuing token
func (k BaseKeeper) DeductIssueTokenFee(ctx sdk.Context, owner sdk.AccAddress, symbol string) error {
	fees, err := k.GetIssueTokenFee(ctx, symbol)
	if err != nil {
		return err
	}
	return handleFees(ctx, k, owner, fees)
}

// DeductMintTokenFee performs fee handling for minting token
func (k BaseKeeper) DeductMintTokenFee(ctx sdk.Context, owner sdk.AccAddress, symbol string) error {
	fees, err := k.GetMintTokenFee(ctx, symbol)
	if err != nil {
		return err
	}
	return handleFees(ctx, k, owner, fees)
}

// GetIssueTokenFee returns the token issuance fee
func (k BaseKeeper) GetIssueTokenFee(ctx sdk.Context, symbol string) (sdk.Coin, error) {
	params := k.GetParams(ctx)
	return k.calcIssueTokenFee(params, symbol), nil
}

// GetMintTokenFee returns the token minting fee
func (k BaseKeeper) GetMintTokenFee(ctx sdk.Context, symbol string) (sdk.Coin, error) {
	params := k.GetParams(ctx)
	issueFees := k.calcIssueTokenFee(params, symbol)

	mintFees := sdk.NewDecFromInt(issueFees.Amount).Mul(params.MintFeeRatio).TruncateInt()
	return sdk.NewCoin(params.IssueFee.Denom, mintFees), nil
}

func (k BaseKeeper) calcIssueTokenFee(params types.Params, symbol string) sdk.Coin {
	issueTokenFee := params.IssueFee

	fees := calcDataCost(symbol, issueTokenFee.Amount)
	if fees.GT(sdk.NewDec(1)) {
		return sdk.NewCoin(issueTokenFee.Denom, fees.TruncateInt())
	}
	return sdk.NewCoin(issueTokenFee.Denom, sdk.OneInt())
}

// calcDataCost computes the actual cost according to the data and principal
// The larger the data, the lower the cost
func calcDataCost(data string, principal sdk.Int) sdk.Dec {
	factor := getFactor(data)
	return sdk.NewDecFromInt(principal).Quo(factor)
}

// getFactor computes the factor
// Note: make sure that the data size is examined before invoking the function
// factor = (ln(len({data}))/ln{base})^{exp}
func getFactor(data string) sdk.Dec {
	dataLen := len(data)
	if dataLen == 0 {
		panic("the length of data must be greater than 0")
	}

	numerator := math.Log(float64(dataLen))
	denominator := math.Log(FactorBase)
	factor := math.Pow(numerator/denominator, FactorExp)

	factorDec, err := sdk.NewDecFromStr(strconv.FormatFloat(factor, 'f', 2, 64))
	if err != nil {
		panic("invalid factor string")
	}

	return factorDec
}

// handleFees handles the fees of token
func handleFees(ctx sdk.Context, k BaseKeeper, from sdk.AccAddress, feesCoin sdk.Coin) error {
	params := k.GetParams(ctx)
	tokenTax := params.TokenTax

	tokenTaxCoin := sdk.NewCoin(feesCoin.Denom,
		sdk.NewDecFromInt(feesCoin.Amount).Mul(tokenTax).TruncateInt())
	burnedCoins := sdk.NewCoins(feesCoin.Sub(tokenTaxCoin))

	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, from, types.ModuleName, sdk.NewCoins(feesCoin)); err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, 
		sdk.NewCoins(tokenTaxCoin)); err != nil {
		return err
	}
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, burnedCoins)
}
