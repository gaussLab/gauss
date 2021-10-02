package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "token"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouteKey defines the module's message routing key
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName
)

var (
	// SymbolPrefix define a symbol prefix
	SymbolPrefix    = []byte{0x21}
	// UnitPrefix define a unit prefix
	UnitPrefix      = []byte{0x22}
	// OwnerSymbolsKey define a symbols prefix with owner
	OwnerSymbolKey  = []byte{0x23}
	// BurntCoinPrefix define a symbols prefix 
	BurntCoinPrefix = []byte{0x24}
)

// GetSymbolKey returns the key with the specified symbol
func GetSymbolKey(symbol string) []byte {
	return append(SymbolPrefix, []byte(symbol)...)
}

// GetUnitKey returns the key with the specified symbol
func GetUnitKey(unit string) []byte {
	return append(UnitPrefix, []byte(unit)...)
}

// GetOwnerSymbolKey returns the key of the specified owner and symbol. Intended for querying all symbols of an owner
func GetOwnerSymbolKey(owner sdk.AccAddress, symbol string) []byte {
	return append(append(OwnerSymbolKey, owner.Bytes()...), []byte(symbol)...)
}

// GetBurntCoinKey
func GetBurntCoinKey(symbol string) []byte {
	return append(BurntCoinPrefix, []byte(symbol)...)
}

