package types

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// MaximumNameLen is the maximum limitation for the length of the token's name
	MaximumNameLen = 32
	// MinimumSymbolLen is the minimum limitation for the length of the token's symbol
	MinimumSymbolLen = 3
	// MaximumSymbolLen is the maximum limitation for the length of the token's symbol
	MaximumSymbolLen = 64
	// MaximumDecimals is the maximum limitation for token decimals
	MaximumDecimals = uint32(18)
	// MaximumAmount is the maximum limitation for the token supply
	MaximumAmount = math.MaxUint64
)

var (
	keywords = strings.Join([]string{
		"peg", "ibc", "swap",
	}, "|")

	regexpKeywordsFmt = fmt.Sprintf("^(%s).*", keywords)
	regexpKeyword     = regexp.MustCompile(regexpKeywordsFmt).MatchString

	regexpSymbolFmt = fmt.Sprintf("^[a-z][a-z0-9]{%d,%d}$", MinimumSymbolLen-1, MaximumSymbolLen-1)
	regexpSymbol    = regexp.MustCompile(regexpSymbolFmt).MatchString
)

// ValidateToken checks if the given token is valid
func (token Token) Validate() error {
	if err := ValidateName(token.GetName()); err != nil {
		return err
	}
	if err := ValidateSymbol(token.GetSymbol()); err != nil {
		return err
	}
	if err := ValidateSmallestUnit(token.GetSmallestUnit()); err != nil {
		return err
	}
	if token.GetTotalSupply() < token.GetInitialSupply() {
		return sdkerrors.Wrapf(ErrInvalidTotalSupply, "invalid total supply %d, only accepts value [%d, %d]", 
			token.GetTotalSupply(), token.GetInitialSupply(), uint64(MaximumAmount))
	}
	if err := ValidateDecimals(token.GetDecimals()); err != nil {
		return err
	}
	if len(token.Owner) > 0 {
		if _, err := sdk.AccAddressFromBech32(token.Owner); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
		}
	}
	return nil
}

// ValidateName verifies whether the given name is legal
func ValidateName(name string) error {
	if (len(name) == 0) || (len(name) > MaximumNameLen) {
		return sdkerrors.Wrapf(ErrInvalidName, "invalid token name %s, only accepts length (0, %d]", 
			name, MaximumNameLen)
	}
	return nil
}

// ValidateSymbol checks if the given symbol is valid
func ValidateSymbol(symbol string) error {
	if !regexpSymbol(symbol) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid symbol: %s, only accepts lowercase ascii " + 
			"letters and numbers, length [%d, %d], and begin with an letter, regexp: %s", 
			symbol, MinimumSymbolLen, MaximumSymbolLen, regexpSymbolFmt)
	}
	if regexpKeyword(symbol) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid symbol: %s, can not begin with keyword: (%s)",
                        symbol, keywords)
	}

	return nil
}

// ValidateSymbol checks if the given symbol is valid
func ValidateSmallestUnit(smallestUnit string) error {
	if !regexpSymbol(smallestUnit) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid smallestUnit: %s, only accepts lowercase ascii " + 
			"letters and numbers, length [%d, %d], and begin with an letter, regexp: %s", 
			smallestUnit, MinimumSymbolLen, MaximumSymbolLen, regexpSymbolFmt)
	}
	if regexpKeyword(smallestUnit) {
		return sdkerrors.Wrapf(ErrInvalidSymbol, "invalid smallestUnit: %s, can not begin with keyword: (%s)",
                        smallestUnit, keywords)
	}

	return nil
}

// ValidateDecimals verifies whether the given decimals is legal
func ValidateDecimals(decimals uint32) error {
	if decimals > MaximumDecimals {
		return sdkerrors.Wrapf(ErrInvalidDecimals, "invalid token decimals %d, only accepts value [0, %d]", 
			decimals, MaximumDecimals)
	}
	return nil
}

// ValidateAmountGTZero checks if the given amount
func ValidateAmountGTZero(amount uint64) error {
	if amount == 0 {
		return sdkerrors.Wrapf(ErrInvalidAmount, "invalid amount [%d]: only accepts value (0, %d]",
			amount, uint64(MaximumAmount))
	}
	return nil
}

