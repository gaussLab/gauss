package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// var _ paramtypes.ParamSet = (*Params)(nil)

const (
	DefaultParamsDenom  = sdk.DefaultBondDenom
)

var (
	KeyTokenTax           = []byte("TokenTax")
	KeyIssueTokenFee      = []byte("IssueTokenFee")
	KeyMintTokenFeeRatio  = []byte("MintTokenFeeRatio")
)

// ParamKeyTable for token module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new parameter configuration for the bank module
func NewParams(tokenTax sdk.Dec, issueFee sdk.Coin, mintFeeRatio sdk.Dec) Params {
	return Params{
		TokenTax: tokenTax,
		IssueFee: issueFee,
		MintFeeRatio: mintFeeRatio,
	}
}

// DefaultParams is the default parameter configuration for the token module
func DefaultParams() Params {
	return NewParams(
		sdk.NewDecWithPrec(4, 1), // 40%
		sdk.NewCoin(DefaultParamsDenom, sdk.NewInt(100000000)),
		sdk.NewDecWithPrec(3, 1), // 5
	)
}

// Validate all token module parameters
func (p Params) Validate() error {
	if err := validateTokenTax(p.TokenTax); err != nil {
		return err
	}
	if err := validateMintFeeRatio(p.MintFeeRatio); err != nil {
		return err
	}
	if err := validateIssueTokenFee(p.IssueFee); err != nil {
		return err
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyTokenTax, &p.TokenTax, validateTokenTax),
		paramtypes.NewParamSetPair(KeyIssueTokenFee, &p.IssueFee, validateIssueTokenFee),
		paramtypes.NewParamSetPair(KeyMintTokenFeeRatio, &p.MintFeeRatio, validateMintFeeRatio),
	}
}

func validateTokenTax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v.GT(sdk.NewDec(1)) || v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("Tax should be between [0, 1]", v.String())
	}

	return nil
}

func validateMintFeeRatio(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v.GT(sdk.NewDec(1)) || v.LT(sdk.ZeroDec()) {
		return fmt.Errorf("Rate should be between [0, 1]", v.String())
	}

	return nil
}

func validateIssueTokenFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("fees should not be negative")
	}
	return nil
}
