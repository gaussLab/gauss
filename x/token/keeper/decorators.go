package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	
	"github.com/gauss/gauss/v4/x/token/types"
)

type ValidateTokenDecorator struct {
	k  Keeper
	bk types.BankKeeper
}

func NewValidateTokenDecorator(k Keeper, bk types.BankKeeper) ValidateTokenDecorator {
	return ValidateTokenDecorator{
		k:  k,
		bk: bk,
	}
}

// AnteHandle returns an AnteHandler that checks if the balance of
// the fee payer is sufficient for token related fee
func (dtf ValidateTokenDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	feesMap := make(map[string]sdk.Coin)

	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *types.MsgIssueToken:
			issueFee, err := dtf.k.GetIssueTokenFee(ctx, msg.Symbol)
			if err != nil {
				return ctx, sdkerrors.Wrap(types.ErrInvalidIssueFee, err.Error())
			}
			
			if fees, ok := feesMap[msg.Owner]; ok {
				feesMap[msg.Owner] = fees.Add(issueFee)
			} else {
				feesMap[msg.Owner] = issueFee
			}
		case *types.MsgMintToken:
			mintFee, err := dtf.k.GetMintTokenFee(ctx, msg.Symbol)
			if err != nil {
				return ctx, sdkerrors.Wrap(types.ErrInvalidIssueFee, err.Error())
			}
			
			if fees, ok := feesMap[msg.Owner]; ok {
				feesMap[msg.Owner] = fees.Add(mintFee)
			} else {
				feesMap[msg.Owner] = mintFee
			}
		}
	}

	for addr, fees := range feesMap {
		owner, _ := sdk.AccAddressFromBech32(addr)
		balance := dtf.bk.GetBalance(ctx, owner, fees.Denom)
		if balance.IsLT(fees) {
			return ctx, sdkerrors.Wrapf(
				sdkerrors.ErrInsufficientFunds, "insufficient fees: balance[%s] < fees[%s]", balance, fees,
			)
		}
	}

	return next(ctx, tx, simulate)
}
