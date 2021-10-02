package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/token interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*TokenI)(nil), nil)

	cdc.RegisterConcrete(&Token{}, "gauss/token/Token", nil)

	cdc.RegisterConcrete(&MsgIssueToken{}, "gauss/token/MsgIssueToken", nil)
	cdc.RegisterConcrete(&MsgEditToken{}, "gauss/token/MsgEditToken", nil)
	cdc.RegisterConcrete(&MsgMintToken{}, "gauss/token/MsgMintToken", nil)
	cdc.RegisterConcrete(&MsgBurnToken{}, "gauss/token/MsgBurnToken", nil)
	cdc.RegisterConcrete(&MsgUnlockToken{}, "gauss/token/MsgUnlockToken", nil)
	cdc.RegisterConcrete(&MsgTransferTokenOwner{}, "gauss/token/MsgTransferTokenOwner", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgIssueToken{},
		&MsgEditToken{},
		&MsgMintToken{},
		&MsgBurnToken{},
		&MsgUnlockToken{},
		&MsgTransferTokenOwner{},
	)
	registry.RegisterInterface(
		"gauss.token.TokenI",
		(*TokenI)(nil),
		&Token{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()

	// ModuleCdc references the global x/bank module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
