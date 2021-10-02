package token

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	ibctransfertypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"

	"github.com/gauss/gauss/v4/x/token/keeper"
	"github.com/gauss/gauss/v4/x/token/types"
)


// ValidateTokenDecorator is responsible for restricting the token participation of the swap prefix
type ValidateTokenDecorator struct {
	keeper keeper.Keeper
}

// NewValidateTokenDecorator returns an instance of ValidateTokenDecorator
func NewValidateTokenDecorator(tk keeper.Keeper) ValidateTokenDecorator {
	return ValidateTokenDecorator {
		keeper: tk,
	}
}

// AnteHandle checks the transaction
func (vtd ValidateTokenDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, 
	simulate bool, next sdk.AnteHandler) (sdk.Context, error) {

	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *types.MsgBurnToken:
			if _, err := vtd.keeper.GetToken(ctx, msg.Symbol); err != nil {
				return ctx, sdkerrors.Wrap(
					sdkerrors.ErrInvalidRequest, "burn failed")
			}
		case *ibctransfertypes.MsgTransfer:
		case *govtypes.MsgSubmitProposal:
		case *govtypes.MsgDeposit:
		default:
			break
		}
	}

	return next(ctx, tx, simulate)
}
