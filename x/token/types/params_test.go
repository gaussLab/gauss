package types_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateParams(t *testing.T) {
	tests := []struct {
		testCase string
		Params
		expectPass bool
	}{
		{
			"Minimum value",
			Params{
				CommunityTax:  sdk.ZeroDec(),
				IssueFee:      sdk.NewCoin(DefaultParamsDenom, sdk.ZeroInt()),
				MintFeeRatio:  sdk.ZeroDec(),
			},
			true,
		}, {
			"Maximum value",
			Params{
				CommunityTax: sdk.NewDec(1),
				IssueFee:     sdk.NewCoin(DefaultParamsDenom, sdk.NewInt(math.MaxInt64)),
				MintFeeRatio: sdk.NewDec(1),
			},
			true,
		}, {
			"CommunityTax less than the maximum.",
			Params{
				CommunityTax: sdk.NewDecWithPrec(-1, 1),
				IssueFee:     sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
				MintFeeRatio: sdk.NewDec(0),
			},
			false
		}, {
			"MintFeeRatio less than the maximum",
			Params{
				CommunityTax: sdk.NewDec(0),
				IssueFee:     sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
				MintFeeRatio: sdk.NewDecWithPrec(-1, 1),
			},
			false
		}, {
			"CommunityTax greater than the maximum.",
			Params{
				CommunityTax: sdk.NewDecWithPrec(11, 1),
				IssueFee:     sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
				MintFeeRatio: sdk.NewDec(1),
			},
			false
		}, {
			"MintFeeRatio greater than the maximum",
			Params{
				CommunityTax: sdk.NewDec(1),
				IssueFee:     sdk.NewCoin(defaultToken.Symbol, sdk.NewInt(1)),
				MintFeeRatio: sdk.NewDecWithPrec(11, 1),
			},
			false

		}, {
			"IssueTokenFee is negative",
			Params{
				CommunityTax: sdk.NewDec(1),
				IssueFee:     sdk.Coin{Denom: DefaultParamsDenom, Amount: sdk.NewInt(-1)},
				MintFeeRatio: sdk.NewDec(1),
			},
			false,
		},
	}

	for _, tc := range tests {
		if tc.expectPass {
			require.Nil(t, tc.Params.Validate(), "test: %v", tc.testCase)
		} else {
			require.NotNil(t, tc.Params.Validate(), "test: %v", tc.testCase)
		}
	}
}
