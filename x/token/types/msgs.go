package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// MsgRoute identifies transaction types
	MsgRoute = "token"

	TypeMsgIssueToken         = "issue_token"
	TypeMsgEditToken          = "edit_token"
	TypeMsgMintToken          = "mint_token"
	TypeMsgBurnToken          = "burn_token"
	TypeMsgUnlockToken        = "unlock_token"
	TypeMsgTransferTokenOwner = "transfer_token_owner"

	// DoNotModify used to indicate that some field should not be updated
	DoNotModify = "[do-not-modify]"
)

var (
	_ sdk.Msg = &MsgIssueToken{}
	_ sdk.Msg = &MsgEditToken{}
	_ sdk.Msg = &MsgMintToken{}
	_ sdk.Msg = &MsgBurnToken{}
	_ sdk.Msg = &MsgUnlockToken{}
	_ sdk.Msg = &MsgTransferTokenOwner{}
)

// NewMsgIssueToken - construct token issue msg.
func NewMsgIssueToken(
	name string, symbol string, smallestUnit string, 
	decimals uint32, initialSupply, totalSupply uint64,
	mintable bool, unlocked bool, owner string,
) *MsgIssueToken {
	return &MsgIssueToken{
		Name:          name,
		Symbol:        symbol,
		SmallestUnit:  smallestUnit,
		Decimals:      decimals,
		InitialSupply: initialSupply,
		TotalSupply:   totalSupply,
		Mintable:      mintable,
		Unlocked:      unlocked,
		Owner:         owner,
	}
}

// Route Implements Msg.
func (msg MsgIssueToken) Route() string { return MsgRoute }

// Type Implements Msg.
func (msg MsgIssueToken) Type() string { return TypeMsgIssueToken }

// GetSignBytes Implements Msg.
func (msg MsgIssueToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgIssueToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic Implements Msg.
func (msg MsgIssueToken) ValidateBasic() error {
	owner, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	return NewToken(msg.Name,
		msg.Symbol,
		msg.SmallestUnit,
		msg.Decimals,
		msg.InitialSupply,
		msg.TotalSupply,
		msg.Mintable,
		owner,
	).Validate()
}

// NewMsgEditToken creates a MsgEditToken
func NewMsgEditToken(symbol string, mintable bool, owner string) *MsgEditToken {
	return &MsgEditToken{
		Symbol:    symbol,
		Mintable:  mintable,
		Owner:     owner,
	}
}

// Route implements Msg
func (msg MsgEditToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgEditToken) Type() string { return TypeMsgEditToken }

// GetSignBytes implements Msg
func (msg MsgEditToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgEditToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgEditToken) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	return ValidateSymbol(msg.Symbol)
}

// NewMsgMintToken creates a MsgMintToken
func NewMsgMintToken(symbol, owner, to string, amount uint64) *MsgMintToken {
	return &MsgMintToken{
		Symbol: symbol,
		Owner:  owner,
		To:     to,
		Amount: amount,
	}
}

// Route implements Msg
func (msg MsgMintToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgMintToken) Type() string { return TypeMsgMintToken }

// GetSignBytes implements Msg
func (msg MsgMintToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgMintToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgMintToken) ValidateBasic() error {
	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	if err := ValidateAmountGTZero(msg.Amount); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	// check the reception
	if len(msg.To) > 0 {
		if _, err := sdk.AccAddressFromBech32(msg.To); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid mint reception address (%s)", err)
		}
	}

	return nil
}

// NewMsgBurnToken creates a MsgMintToken
func NewMsgBurnToken(symbol string, owner string, amount uint64) *MsgBurnToken {
	return &MsgBurnToken{
		Symbol: symbol,
		Sender: owner,
		Amount: amount,
	}
}

// Route implements Msg
func (msg MsgBurnToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgBurnToken) Type() string { return TypeMsgBurnToken }

// GetSignBytes implements Msg
func (msg MsgBurnToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgBurnToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgBurnToken) ValidateBasic() error {
	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	if err := ValidateAmountGTZero(msg.Amount); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	return nil
}

// NewMsgUnlockToken creates a MsgUnlockToken
func NewMsgUnlockToken(symbol string, owner string) *MsgUnlockToken {
	return &MsgUnlockToken{
		Symbol: symbol,
		Owner: owner,
	}
}

// Route implements Msg
func (msg MsgUnlockToken) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgUnlockToken) Type() string { return TypeMsgUnlockToken }

// GetSignBytes implements Msg
func (msg MsgUnlockToken) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgUnlockToken) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgUnlockToken) ValidateBasic() error {
	if err := ValidateSymbol(msg.Symbol); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Owner); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	return nil
}

// NewMsgTransferTokenOwner return a instance of MsgTransferTokenOwner
func NewMsgTransferTokenOwner(symbol, oldOwner, newOwner string) *MsgTransferTokenOwner {
	return &MsgTransferTokenOwner{
		Symbol:   symbol,
		OldOwner: oldOwner,
		NewOwner: newOwner,
	}
}

// Route implements Msg
func (msg MsgTransferTokenOwner) Route() string { return MsgRoute }

// Type implements Msg
func (msg MsgTransferTokenOwner) Type() string { return TypeMsgTransferTokenOwner }

// GetSignBytes implements Msg
func (msg MsgTransferTokenOwner) GetSignBytes() []byte {
	b, err := ModuleCdc.MarshalJSON(&msg)
	if err != nil {
		panic(err)
	}

	return sdk.MustSortJSON(b)
}

// GetSigners implements Msg
func (msg MsgTransferTokenOwner) GetSigners() []sdk.AccAddress {
	from, err := sdk.AccAddressFromBech32(msg.OldOwner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{from}
}

// ValidateBasic implements Msg
func (msg MsgTransferTokenOwner) ValidateBasic() error {
	if err := ValidateSymbol(msg.Symbol) ; err != nil {
		return err
	}

	oldOwner, err := sdk.AccAddressFromBech32(msg.OldOwner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid source owner address (%s)", err)
	}

	newOwner, err := sdk.AccAddressFromBech32(msg.NewOwner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid destination owner address (%s)", err)
	}

	if oldOwner.Equals(newOwner) {
		return ErrInvalidToAddress
	}

	return nil
}
