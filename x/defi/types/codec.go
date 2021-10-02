package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/defi interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateDefi{}, "gauss/defi/MsgCreateDefi", nil)
	cdc.RegisterConcrete(&MsgEditDefi{}, "gauss/defi/MsgEditDefi", nil)
	cdc.RegisterConcrete(&MsgDefiDelegate{}, "gauss/defi/MsgDefiDelegate", nil)
	cdc.RegisterConcrete(&MsgDefiUndelegate{}, "gauss/defi/MsgDefiUndelegate", nil)
	cdc.RegisterConcrete(&MsgSetDefiWithdrawAddress{}, "gauss/defi/MsgSetDefiWithdrawAddress", nil)
	cdc.RegisterConcrete(&MsgWithdrawDefiDelegatorReward{}, "gauss/defi/MsgWithdrawDefiDelegatorReward", nil)
	cdc.RegisterConcrete(&MsgWithdrawDefiCommission{}, "gauss/defi/MsgWithdrawDefiCommission", nil)
	cdc.RegisterConcrete(&MsgFundDefiCommunityPool{}, "gauss/defi/MsgFundDefiCommunityPool", nil)
}

// RegisterInterfaces registers the x/defi interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateDefi{},
		&MsgEditDefi{},
		&MsgDefiDelegate{},
		&MsgDefiUndelegate{},
		&MsgSetDefiWithdrawAddress{},
		&MsgWithdrawDefiDelegatorReward{},
		&MsgWithdrawDefiCommission{},
		&MsgFundDefiCommunityPool{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/staking module codec. Note, the codec should
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
