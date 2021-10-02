package types

import (
	"errors"
	"fmt"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	// DefaultUnbondingTime reflects three weeks in seconds as the default
	// unbonding time.
	// TODO: Justify our choice of default here.
	DefaultUnbondingTime time.Duration = time.Hour * 24 * 7 * 3

        // Default maximum number of bonded defis
	DefaultMaxDefis uint32 = 2 

	// Default maximum entries in a UBD/RED pair
	DefaultMaxEntries uint32 = 7

	// DefaultHistorical entries is 10000. Apps that don't use IBC can ignore this
	// value by not adding the staking module to the application module manager's
	// SetOrderBeginBlockers.
	DefaultHistoricalEntries uint32 = 10000
)

var (
	KeyBondDenom         = []byte("BondDenom")
	KeyMintInflation     = []byte("MintInflation")
	KeyCommunityTax      = []byte("CommunityTax")
	KeyCommissionRate    = []byte("CommissionRate")
	KeyMarketRate        = []byte("MarketRate")
	KeyUnbondingTime     = []byte("UnbondingTime")
	KeyMaxDefis          = []byte("MaxDefis")
	KeyMaxEntries        = []byte("MaxEntries")
	KeyHistoricalEntries = []byte("HistoricalEntries") 
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable for defi module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new parameter configuration for the defi module
func NewParams(bondDenom string, mintInflation sdk.Coin, communityTax, commissionRate, marketRate sdk.Dec, 
	unbondingTime time.Duration, maxDefis, maxEntries, historicalEntries uint32) Params {
	return Params{
		BondDenom: bondDenom,
		MintInflation: mintInflation,
		CommunityTax: communityTax,
		CommissionRate: commissionRate,
		MarketRate: marketRate,
		UnbondingTime: unbondingTime,
		MaxDefis: maxDefis,
		MaxEntries: maxEntries,
		HistoricalEntries: historicalEntries,
	}
}

// DefaultParams is the default parameter configuration for the defi module
func DefaultParams() Params {
	return NewParams(
		sdk.DefaultBondDenom,
		sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(0)},
		sdk.NewDecWithPrec(2, 2), // 2%
		sdk.NewDecWithPrec(90, 2), // 90%
		sdk.NewDecWithPrec(5, 2), // 5%
		DefaultUnbondingTime,
		DefaultMaxDefis,
		DefaultMaxEntries,
		DefaultHistoricalEntries,
	)
}

// ParamSetPairs implements params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyBondDenom, &p.BondDenom, validateBondDenom),
		paramtypes.NewParamSetPair(KeyMintInflation, &p.MintInflation, validateMintInflation),
		paramtypes.NewParamSetPair(KeyCommunityTax, &p.CommunityTax, validateCommunityTax),
		paramtypes.NewParamSetPair(KeyCommissionRate, &p.CommissionRate, validateCommissionRate),
		paramtypes.NewParamSetPair(KeyMarketRate, &p.MarketRate, validateMarketRate),
		paramtypes.NewParamSetPair(KeyUnbondingTime, &p.UnbondingTime, validateUnbondingTime),
		paramtypes.NewParamSetPair(KeyMaxDefis, &p.MaxDefis, validateMaxDefis),
		paramtypes.NewParamSetPair(KeyMaxEntries, &p.MaxEntries, validateMaxEntries),
		paramtypes.NewParamSetPair(KeyHistoricalEntries, &p.HistoricalEntries, validateHistoricalEntries),
	}
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// unmarshal the current staking params value from store key or panic
func MustUnmarshalParams(cdc *codec.LegacyAmino, value []byte) Params {
	params, err := UnmarshalParams(cdc, value)
	if err != nil {
		 panic(err)
	}
	
	return params
}

// unmarshal the current staking params value from store key
func UnmarshalParams(cdc *codec.LegacyAmino, value []byte) (params Params, err error) {
	err = cdc.UnmarshalBinaryBare(value, &params)
	return
}


// Validate all defi module parameters
func (p Params) Validate() error {
	if err := validateBondDenom(p.BondDenom); err != nil {
		return err
	}
	if err := validateMintInflation(p.MintInflation); err != nil {
		return err
	}
	if err := validateCommunityTax(p.CommunityTax); err != nil {
		return err
	}
	if err := validateCommissionRate(p.CommissionRate); err != nil {
		return err
	}
	if err := validateMarketRate(p.MarketRate); err != nil {
		return err
	}
	if err := validateUnbondingTime(p.UnbondingTime); err != nil {
		return err
	}
	if err := validateMaxDefis(p.MaxDefis); err != nil {
		return err
	}
	if err := validateMaxEntries(p.MaxEntries); err != nil {
		return err
	}
	if err := validateHistoricalEntries(p.HistoricalEntries); err != nil {
		return err
	}
	return nil
}

func validateMintInflation(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsValid() {
		return fmt.Errorf("invalid constant mint-inflation: %s", v)
	}

	return nil
}

func validateCommunityTax(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	
	if v.IsNil() {
		return fmt.Errorf("community tax must be not nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("community tax must be positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("community tax too large: %s", v)
	}

	return nil
}

func validateCommissionRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("commission rate must be not nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("commission rate must be positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("commission rate too large: %s", v)
	}

	return nil
}

func validateMarketRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("commission rate must be not nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("commission rate must be positive: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("commission rate too large: %s", v)
	}

	return nil
}

func validateUnbondingTime(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("unbonding time must be positive: %d", v)
	}

	return nil
}

func validateMaxDefis(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max defis must be positive: %d", v)
	}
	
	return nil
}

func validateMaxEntries(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max entries must be positive: %d", v)
	}

	return nil
}

func validateHistoricalEntries(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	return nil
}

func validateBondDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("bond denom cannot be blank")
	}

	if err := sdk.ValidateDenom(v); err != nil {
		return err
	}

	return nil
}
