package types

import (
	"bytes"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// staking message types
const (
	TypeMsgCreateDefi                  = "create_defi"
	TypeMsgEditDefi                    = "edit_defi"
	TypeMsgDefiDelegate                    = "delegate"
	TypeMsgDefiUndelegate                  = "begin_unbonding"
	TypeMsgSetDefiWithdrawAddress          = "set_withdraw_address"
	TypeMsgWithdrawDefiCommission      = "withdraw_defi_commission"
	TypeMsgWithdrawDefiDelegatorReward     = "withdraw-delegator-reward"
	TypeMsgFundDefiCommunityPool           = "fund_community_pool"
)

var (
	_ sdk.Msg                            = &MsgCreateDefi{}
	_ codectypes.UnpackInterfacesMessage = (*MsgCreateDefi)(nil)
	_ sdk.Msg                            = &MsgCreateDefi{}
	_ sdk.Msg                            = &MsgEditDefi{}
	_ sdk.Msg                            = &MsgDefiDelegate{}
	_ sdk.Msg                            = &MsgDefiUndelegate{}
	_ sdk.Msg                            = &MsgSetDefiWithdrawAddress{}
	_ sdk.Msg                            = &MsgWithdrawDefiCommission{}
	_ sdk.Msg                            = &MsgWithdrawDefiDelegatorReward{}
	_ sdk.Msg                            = &MsgFundDefiCommunityPool{}
)

// NewMsgCreateDefi creates a new MsgCreateDefi instance.
// Delegator address and defi address are the same.
func NewMsgCreateDefi(
	defiAddr sdk.ValAddress, //nolint:interfacer
	selfDelegation sdk.Coin, description Description, minSelfDelegation sdk.Int,
) (*MsgCreateDefi, error) {
	return &MsgCreateDefi{
		Description:       description,
		DelegatorAddress:  sdk.AccAddress(defiAddr).String(),
		DefiAddress:       defiAddr.String(),
		Value:             selfDelegation,
		MinSelfDelegation: minSelfDelegation,
	}, nil
}

// Route implements the sdk.Msg interface.
func (msg MsgCreateDefi) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgCreateDefi) Type() string { return TypeMsgCreateDefi }

// GetSigners implements the sdk.Msg interface. It returns the address(es) that
// must sign over msg.GetSignBytes().
// If the defi address is not same as delegator's, then the defi must
// sign the msg as well.
func (msg MsgCreateDefi) GetSigners() []sdk.AccAddress {
	// delegator is first signer so delegator pays fees
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	addrs := []sdk.AccAddress{delAddr}
	addr, err := sdk.ValAddressFromBech32(msg.DefiAddress)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(delAddr.Bytes(), addr.Bytes()) {
		addrs = append(addrs, sdk.AccAddress(addr))
	}

	return addrs
}

// GetSignBytes returns the message bytes to sign over.
func (msg MsgCreateDefi) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgCreateDefi) ValidateBasic() error {
	// note that unmarshaling from bech32 ensures either empty or valid
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return err
	}
	if delAddr.Empty() {
		return ErrEmptyDelegatorAddr
	}

	if msg.DefiAddress == "" {
		return ErrEmptyDefiAddr
	}

	defiAddr, err := sdk.ValAddressFromBech32(msg.DefiAddress)
	if err != nil {
		return err
	}
	if !sdk.AccAddress(defiAddr).Equals(delAddr) {
		return ErrBadDefiAddr
	}

	if !msg.Value.IsValid() || !msg.Value.Amount.IsPositive() {
		return ErrBadDelegationAmount
	}

	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if !msg.MinSelfDelegation.IsPositive() {
		return ErrMinSelfDelegationInvalid
	}

	if msg.Value.Amount.LT(msg.MinSelfDelegation) {
		return ErrSelfDelegationBelowMinimum
	}

	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateDefi) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return nil
}

// NewMsgEditDefi creates a new MsgEditDefi instance
//nolint:interfacer
func NewMsgEditDefi(defiAddr sdk.ValAddress, description Description, newMinSelfDelegation *sdk.Int) *MsgEditDefi {
	return &MsgEditDefi{
		Description:       description,
		DefiAddress:       defiAddr.String(),
		MinSelfDelegation: newMinSelfDelegation,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgEditDefi) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgEditDefi) Type() string { return TypeMsgEditDefi }

// GetSigners implements the sdk.Msg interface.
func (msg MsgEditDefi) GetSigners() []sdk.AccAddress {
	defiAddr, err := sdk.ValAddressFromBech32(msg.DefiAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{defiAddr.Bytes()}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgEditDefi) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgEditDefi) ValidateBasic() error {
	if msg.DefiAddress == "" {
		return ErrEmptyDefiAddr
	}

	if msg.Description == (Description{}) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "empty description")
	}

	if msg.MinSelfDelegation != nil && !msg.MinSelfDelegation.IsPositive() {
		return ErrMinSelfDelegationInvalid
	}

	return nil
}

// NewMsgDefiDelegate creates a new MsgDefiDelegate instance.
//nolint:interfacer
func NewMsgDefiDelegate(delAddr sdk.AccAddress, defiAddr sdk.ValAddress, amount sdk.Coin) *MsgDefiDelegate {
	return &MsgDefiDelegate{
		DelegatorAddress: delAddr.String(),
		DefiAddress:      defiAddr.String(),
		Amount:           amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgDefiDelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgDefiDelegate) Type() string { return TypeMsgDefiDelegate }

// GetSigners implements the sdk.Msg interface.
func (msg MsgDefiDelegate) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgDefiDelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDefiDelegate) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return ErrEmptyDelegatorAddr
	}

	if msg.DefiAddress == "" {
		return ErrEmptyDefiAddr
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return ErrBadDelegationAmount
	}

	return nil
}

// NewMsgDefiUndelegate creates a new MsgDefiUndelegate instance.
//nolint:interfacer
func NewMsgDefiUndelegate(delAddr sdk.AccAddress, defiAddr sdk.ValAddress, amount sdk.Coin) *MsgDefiUndelegate {
	return &MsgDefiUndelegate{
		DelegatorAddress: delAddr.String(),
		DefiAddress: defiAddr.String(),
		Amount:           amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgDefiUndelegate) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgDefiUndelegate) Type() string { return TypeMsgDefiUndelegate }

// GetSigners implements the sdk.Msg interface.
func (msg MsgDefiUndelegate) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgDefiUndelegate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgDefiUndelegate) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return ErrEmptyDelegatorAddr
	}

	if msg.DefiAddress == "" {
		return ErrEmptyDefiAddr
	}

	if !msg.Amount.IsValid() || !msg.Amount.Amount.IsPositive() {
		return ErrBadSharesAmount
	}

	return nil
}

func NewMsgSetDefiWithdrawAddress(delAddr, withdrawAddr sdk.AccAddress) *MsgSetDefiWithdrawAddress {
	return &MsgSetDefiWithdrawAddress{
		DelegatorAddress: delAddr.String(),
		WithdrawAddress:  withdrawAddr.String(),
	}
}

func (msg MsgSetDefiWithdrawAddress) Route() string { return ModuleName }
func (msg MsgSetDefiWithdrawAddress) Type() string  { return TypeMsgSetDefiWithdrawAddress }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgSetDefiWithdrawAddress) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

// get the bytes for the message signer to sign on
func (msg MsgSetDefiWithdrawAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgSetDefiWithdrawAddress) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return ErrEmptyDelegatorAddr
	}
	if msg.WithdrawAddress == "" {
		return ErrEmptyWithdrawAddr
	}

	return nil
}

// NewMsgWithdrawDefiDelegatorReward creates a new MsgWithdrawDefiDelegatorReward instance.
//nolint:interfacer
func NewMsgWithdrawDefiDelegatorReward(delAddr sdk.AccAddress, defiAddr sdk.ValAddress) *MsgWithdrawDefiDelegatorReward {
	return &MsgWithdrawDefiDelegatorReward{
		DelegatorAddress: delAddr.String(),
		DefiAddress: defiAddr.String(),
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgWithdrawDefiDelegatorReward) Route() string { return RouterKey }

// Type implements the sdk.Msg interface.
func (msg MsgWithdrawDefiDelegatorReward) Type() string { return TypeMsgWithdrawDefiDelegatorReward }

// GetSigners implements the sdk.Msg interface.
func (msg MsgWithdrawDefiDelegatorReward) GetSigners() []sdk.AccAddress {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{delAddr}
}

// GetSignBytes implements the sdk.Msg interface.
func (msg MsgWithdrawDefiDelegatorReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgWithdrawDefiDelegatorReward) ValidateBasic() error {
	if msg.DelegatorAddress == "" {
		return ErrEmptyDelegatorAddr
	}

	if msg.DefiAddress == "" {
		return ErrEmptyDefiAddr
	}

	return nil
}

func NewMsgWithdrawDefiCommission(defiAddr sdk.ValAddress) *MsgWithdrawDefiCommission {
	return &MsgWithdrawDefiCommission{
		DefiAddress: defiAddr.String(),
	}
}

func (msg MsgWithdrawDefiCommission) Route() string { return ModuleName }
func (msg MsgWithdrawDefiCommission) Type() string  { return TypeMsgWithdrawDefiCommission }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawDefiCommission) GetSigners() []sdk.AccAddress {
	defiAddr, err := sdk.ValAddressFromBech32(msg.DefiAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{defiAddr.Bytes()}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawDefiCommission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawDefiCommission) ValidateBasic() error {
	if msg.DefiAddress == "" {
		return ErrEmptyDefiAddr
	}
	return nil
}
// NewMsgFundDefiCommunityPool returns a new MsgFundDefiCommunityPool with a sender and
// a funding amount.
func NewMsgFundDefiCommunityPool(amount sdk.Coins, depositor sdk.AccAddress) *MsgFundDefiCommunityPool {
	return &MsgFundDefiCommunityPool{
		Amount:    amount,
		Depositor: depositor.String(),
	}
}

// Route returns the MsgFundDefiCommunityPool message route.
func (msg MsgFundDefiCommunityPool) Route() string { return ModuleName }

// Type returns the MsgFundDefiCommunityPool message type.
func (msg MsgFundDefiCommunityPool) Type() string { return TypeMsgFundDefiCommunityPool }

// GetSigners returns the signer addresses that are expected to sign the result
// of GetSignBytes.
func (msg MsgFundDefiCommunityPool) GetSigners() []sdk.AccAddress {
	depoAddr, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{depoAddr}
}

// GetSignBytes returns the raw bytes for a MsgFundDefiCommunityPool message that
// the expected signer needs to sign.
func (msg MsgFundDefiCommunityPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic performs basic MsgFundDefiCommunityPool message validation.
func (msg MsgFundDefiCommunityPool) ValidateBasic() error {
	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	if msg.Depositor == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, msg.Depositor)
	}

	return nil
}
