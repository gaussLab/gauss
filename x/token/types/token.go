package types

import (
	"github.com/gogo/protobuf/proto"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ proto.Message = &Token{}
)

// TokenI defines an interface for Token
type TokenI interface {
	GetName() string
	GetSymbol() string
	GetSmallestUnit() string
	GetDecimals() uint32
	GetInitialSupply() uint64
	GetTotalSupply() uint64
	GetMintable() bool
	GetOwner() sdk.AccAddress
	GetOwnerString() string
}

// NewToken constructs a new Token instance
func NewToken(name, symbol string, smallestUnit string, decimals uint32, initialSupply, 
	totalSupply uint64, mintable bool, owner sdk.AccAddress,) Token {
	if totalSupply == 0 {
		if mintable {
			totalSupply = MaximumAmount
		} else {
			totalSupply = initialSupply
		}
	}

	return Token{Name: name, Symbol: symbol, SmallestUnit: smallestUnit, 
		Decimals: decimals, InitialSupply: initialSupply,
		TotalSupply: totalSupply, Mintable: mintable, Owner: owner.String()}
}

func (t Token) GetName() string {
	return t.Name
}

func (t Token) GetSymbol() string {
	return t.Symbol
}

func (t Token) GetSmallestUnit() string {
	return t.SmallestUnit
}

func (t Token) GetDecimals() uint32 {
	return t.Decimals
}

func (t Token) GetInitialSupply() uint64 {
	return t.InitialSupply
}

func (t Token) GetTotalSupply() uint64 {
	return t.TotalSupply
}

func (t Token) GetMintable() bool {
	return t.Mintable
}

func (t Token) GetOwner() sdk.AccAddress {
	owner, _ := sdk.AccAddressFromBech32(t.Owner)
	return owner
}

func (t Token) GetOwnerString() string {
	return t.Owner
}

func (t Token) String() string {
	bz, _ := yaml.Marshal(t)
	return string(bz)
}
